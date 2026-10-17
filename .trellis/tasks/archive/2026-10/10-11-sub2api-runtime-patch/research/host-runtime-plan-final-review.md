# Final Review: Sub2API Host Runtime Wiring Plan

- Date: 2026-10-12
- Scope: final plan review against pinned Sub2API v0.2.15 source at
  `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.
- Verdict: **NEEDS-REFINEMENT** for the host implementation candidate. The
  three findings in the prior review are resolved in the updated plan; one
  provisioner-error lifecycle gap remains.
- Deployment acceptance: not evaluated. A code-plan review cannot establish
  that the installed proxy inventory is complete or that real deployment paths
  produce the recorded sources.

## Findings (closed from prior review)

- File: `research/host-runtime-wiring-plan.md`
- Issue: A Gin route middleware could not reject marked requests that Gin
  redirects before route dispatch.
- Fix: The plan now puts the guard around the finalized `http.Server.Handler`,
  outside the existing body-limit wrapper. It specifies unchanged pass-through
  for unmarked requests and permits a marked request into Gin only for exact
  `POST /v1/responses`; tests exercise redirects, no-route/method handling,
  aliases, Google routes, and WebSockets through a bound listener. This matches
  the pinned `ProvideHTTPServer` composition and Gin v1.9.1 dispatch behavior.

- File: `research/host-runtime-wiring-plan.md`
- Issue: Consumed failed tickets could strand the only session until expiry.
- Fix: The plan defines actor-bound abort/replacement, reservation of the
  session slot during begin, partial-begin cleanup, ticket invalidation, and
  in-flight observation and completion coordination. It explicitly prevents
  another actor from aborting/replacing the session and states abort never
  revokes or rebinds a provisioned Key. The required fixtures cover these
  races and ownership boundaries.

- File: `research/host-runtime-wiring-plan.md`
- Issue: `NewRunManager` requires a nonnil executor, while an unprovisioned
  startup has no token or verified local origin.
- Fix: The plan now requires a retained nonnil lazy `ModelExecutor` proxy with a
  typed not-ready error, per-execution readiness resolution, and synchronized
  delegate lifecycle. It also defines actual-listener capture and self-request
  validation, keeps manager/store errors visible without making ordinary host
  startup fatal, and prohibits empty-store fallback or state overwrite.

## Findings (not fixed)

### Provisioner errors can leave completion ownership unresolved

- File: `research/host-runtime-wiring-plan.md`, probe completion protocol.
- Issue: The session state machine defines `completion_pending` for a successful
  `Initialize` followed by a session-evidence write failure, but does not define
  the state or retry path when `Seal` or `EvalKeyProvisioner.Initialize`
  returns an error. This matters because `Initialize` is not side-effect-free:
  it can persist a row-pending intent before Key creation, create the Key row,
  save a ready intent, and then fail while installing runtime identity
  ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L420-L456),
  [eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L517-L530)).
  The session could remain stuck in `completing`, or a generic abort/replacement
  could discard the only retry context for a partially provisioned owned row.
- Required refinement: specify failure transitions separately. A `Seal` failure
  before provisioning may invalidate the attempt and require same-actor
  re-probing. Treat every `Initialize` error as potentially partially
  committed: discard probe tickets and retain an owner-bound recovery state.
  A same-actor explicit retry must revalidate the protected artifact digest,
  proof fingerprint, expiry, and current runtime/source/config revisions, then
  use the existing persisted-intent path to recover only the exact owned row
  and identity. Conflict, disabled, owner, status, whitelist, or binding errors
  remain blocked for explicit repair; do not silently revoke, rebind, or
  reenable a Key. An authenticated abort may clear ephemeral session/ticket
  state while preserving any persisted intent and owned row; a fresh session
  may then retry `Initialize` as the same owner and recover that exact identity.
  No hidden control retry or paid evaluation retry is implied. Add
  failure-injection gates around intent creation, API-Key creation, ready-intent
  save, runtime installation, and evidence save; assert retries do not create
  duplicate rows or change ownership and readiness remains false until recovery
  completes.
- Why not fixed: this is a runtime state-machine and recovery contract; the
  current assignment is a read-only plan review.

### Clarify the observation coordinator lock boundary

- File: `research/host-runtime-wiring-plan.md`, probe session protocol.
- Issue: The detailed concurrency rule correctly says to reserve/consume a
  ticket and increment an in-flight count under the coordinator, release the
  lock for `collector.Observe`, then reacquire it to publish only if the same
  generation remains current. The request sequence later says to call
  `Observe` “under the session coordinator,” which can be read as holding the
  coordinator mutex across the call and conflicts with the explicit lock rule.
- Required refinement: state that the coordinator owns the observation lease
  and lifecycle transition: consume the ticket and increment the in-flight
  count under lock, call `Observe` without holding that mutex, then reacquire it
  to publish only if the same session generation remains current. Make abort
  wait for all leased observations.
- Why not fixed: the intended algorithm is already present; the later shorthand
  should be made consistent before implementation to avoid two readings.

## Next Smallest Implementation Chunk

After patch 004's independent specification, quality, and replay checks pass,
the next isolated chunk should be the protected deployment-inventory reader,
not the combined probe/session/API integration:

- Files in the pinned workspace: new
  `backend/internal/service/eval_deployment_evidence.go` and
  `backend/internal/service/eval_deployment_evidence_test.go` only.
- Repository artifact: export as the next ordered patch, for example
  `patches/sub2api/v0.2.15/005-codex-eval-deployment-evidence.patch`.
- Contract: resolve the fixed `EvalKeyStore.Dir()/deployment-evidence/` path
  internally; reject unsafe directory/file ownership, permissions, symlinks,
  and non-regular files; hash the exact bytes read before strict decoding; and
  validate schema/version, deployment/runtime/source/config revisions, expiry,
  exact required external-case mapping, and case-specific HTTPS destination
  rules. Return an immutable inventory plus its byte digest. Do not add an API,
  readiness claim, session-evidence writer, ticket state, Wire binding, probe
  runner, key bootstrap, or listener lifecycle in this chunk.
- Test gate: missing and unsafe files fail closed; malformed/trailing or stale
  artifacts fail closed; whitespace changes alter the digest; duplicate,
  missing, or extra case destinations fail; unsafe URL forms fail; and returned
  values cannot mutate cached validation state. Replay patches 001-005 on
  pristine pinned source and compare the new source/test bytes. Keep the
  external operator attestation and actual deployment path verification as
  later acceptance requirements.

Keep the actor-bound coordinator as a separate subsequent chunk only after the
completion-error transition above is resolved. Its owner should be limited to
`internal/service/eval_topology_runtime.go`, its focused tests, and the next
ordered patch. Test actor ownership, begin reservation, single-use tickets,
observation leases/generation checks, abort waiting, completion single-flight,
expiry/config-drift invalidation, and same-actor recovery reservation without
HTTP routes, public runner, listener wiring, Key revocation, or automatic model
requests. The gateway hook, outer marker wrapper, admin routes, one-shot client,
and application lifecycle remain a later integration chunk with the actual-
listener and no-auth/no-billing fixtures.

This loader is only a local artifact-validation boundary. It neither proves the
operator's source export authoritative nor verifies that every real ingress was
listed. The complete probe producer remains a later stage after this boundary
and the remaining completion-error contract are resolved.

## Verification

- Plan/source review: pass for the updated marker guard, actor-bound session
  races, lazy executor/startup behavior, protected-evidence contract, and
  shutdown ordering; the verdict remains **NEEDS-REFINEMENT** for the
  provisioner-error and lock-wording findings above.
- Pinned-source evidence checked: `ProvideHTTPServer` final handler composition,
  Gin v1.9.1 pre-middleware redirects, retained executor requirement in
  `NewRunManager`, and partial-persistence branches in
  `EvalKeyProvisioner.Initialize`.
- Lint, type-check, and tests: not run; this was a read-only plan review and no
  product code changed.
- Metadata patch 004 exists and is delivered; its independent reviews are
  ongoing. No deployment target or live topology acceptance was used.
