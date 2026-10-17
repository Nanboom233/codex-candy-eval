# Research: Host Runtime Plan Refinement Response

- Query: Close the final host-runtime plan review findings for provisioner errors, observation leases, evidence-only retry validation, and restart ownership, while accounting for exact-intent recovery behavior.
- Scope: internal
- Date: 2026-10-12

## Findings

The [host runtime wiring plan](host-runtime-wiring-plan.md) now separates seal failures, provisioner errors, and session-evidence write failures.

- A `Seal` or topology-construction failure occurs before `Initialize`. Invalidate the attempt, discard all tickets/collector state, preserve a sanitized same-actor failure, release the completion claim, and require a fresh probe. No Key provisioning side effect has been invoked by this branch.
- Treat every `Initialize` error as potentially partially committed. Keep readiness false, clear ephemeral session/ticket state, publish the owner-bound `provisioning_recovery_required` state, release the active completion claim, and retain a bounded reservation for the initiating actor. A same-actor explicit retry must reload the fixed artifact and verify its ID and exact-byte digest, current source/runtime/deployment/config revisions, live security fingerprint, and fresh proof expiry before invoking `Initialize` again. No background retry or paid evaluation request is permitted.
- Recovery is limited by the current provisioner implementation. `Initialize` persists a `row_pending` intent before looking up or creating the API-Key row (`eval_key_runtime.go:397-456`). When an intent already exists, it verifies lookup/owner/peer identity, saves the new proof, and calls `adoptExistingRow` (`:399-417`). That method looks up the existing row and returns `ErrEvalKeyRowMissing` when it is absent (`:480-498`); it does not resume row creation. A failure after row creation can recover through the same exact intent only if the row still exists and passes `verifyEvalOwnedRow`. A row-missing intent, lookup mismatch, owner/status/binding/whitelist conflict, or other ownership conflict remains blocked for explicit repair. Do not promise automatic recovery for every failure-injection point, create a duplicate row, or mutate the Key to make adoption pass.
- An authenticated same-actor abort clears only ephemeral session/ticket state and preserves any persisted intent and owned row. A fresh same-owner session may retry a verified existing-row recovery. It must not revoke, rebind, reenable, or clear the provisioned identity.
- Reserve `completion_pending` for the one distinct success case: `Initialize` returned success, but saving runtime-owned session evidence failed. Retain its owner-bound completion record and proof, discard tickets, and allow only an idempotent evidence-write retry. Before retrying, reload the exact artifact and current facts, validate artifact ID and exact-byte digest, proof MAC/fingerprint/expiry, and strictly reload the exact Key row with owner/status/binding/whitelist checks. The retry writes evidence only: it must not call `Initialize`, change a Key, reopen observations, or issue a model request.
- After restart, derive recovery ownership from the persisted provisioning intent's `OwnerUserID`. Missing runtime session evidence cannot be treated as successful completion based on that intent; require a fresh same-owner probe unless a distinct durable completion binding exists. Only a validated persisted completion binding can support evidence-only retry after restart.

The observation rule is uniform throughout the plan: under the coordinator lock, validate and consume the one-use ticket and increment the in-flight count; release the lock while calling `collector.Observe`; reacquire it and publish only if the same generation is still current. Abort blocks new leases and waits for existing leases before invalidating the collector. The coordinator mutex is not held across `Observe`, network/file work, or provisioner calls.

Focused lifecycle gates should inject errors at intent creation, API-Key creation, ready-intent save, runtime installation, and session-evidence save. Verify readiness remains false until all evidence and runtime identity are valid; exact existing-row recovery does not create another row or change ownership; row-missing and Key conflicts stay blocked; and evidence-save retry does not repeat provisioning. Include abort/replacement ownership checks and a race where abort waits for a leased observation. These are future implementation gates; this plan review ran no product tests.

## Files Found

- `dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go`: `Initialize` persistence stages and `adoptExistingRow` recovery behavior (`:374-456`, `:480-530`).
- `dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime_test.go`: exact Key ownership/conflict fixtures and store helpers.
- `research/host-runtime-plan-final-review.md`: original completion-error and lock-wording findings.
- `research/host-runtime-wiring-plan.md`: corrected state machine, completion sequence, observation lease, and later patch boundaries.
- `research/deployment-evidence-reader-plan.md`: frozen inventory reader schema, file checks, current-facts boundary, and parser contract.

## Related Specs

- `.trellis/spec/backend/run-state.md`: persistence failure visibility and explicit recovery boundaries for the separate evaluation manager.
- `.trellis/spec/backend/quality-guidelines.md`: focused test and error-state reporting expectations.
- `research/deployment-evidence-reader-plan.md`: separate protected-inventory reader patch and current-facts requirements.

## External References

- No external references were needed. The recovery and lock contracts are based on the pinned Sub2API source and the final review report.

## Caveats / Not Found

- The current `row_pending` intent with no API-Key row has no supported automatic repair/resume operation. The later implementation must surface it as blocked and define an explicit operator repair process before claiming full recovery coverage.
- The pinned source does not implement the evidence-only retry or persist a separate durable completion binding. Until that integration exists, a restart with missing session evidence requires a fresh probe by the intent's persisted owner.
- No runtime coordinator, completion handler, session-evidence writer, or associated failure-injection tests exist in the pinned workspace. The report refines the plan; it does not establish implementation or deployment acceptance.
