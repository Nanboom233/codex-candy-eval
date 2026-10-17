# Recovery Review: Sub2API Host Runtime Plan

- Date: 2026-10-12
- Scope: review of the revised provisioner-error and probe-session recovery
  contracts against pinned Sub2API v0.2.15 source
  `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.
- Verdict: **NEEDS-REFINEMENT** for the candidate completion-pending recovery
  path. The provisioner-error retry and observation-lease corrections are
  source-compatible; freshness and owner continuity still need explicit rules.
- Deployment readiness: not evaluated and cannot be inferred from this plan
  review.

## Findings (fixed)

- File: `research/host-runtime-wiring-plan.md`
- Issue: `Initialize` can partially persist before returning an error, while
  the earlier plan did not define a recovery state.
- Fix: The revised plan retains owner-bound `provisioning_recovery_required`,
  discards ephemeral tickets, requires a same-owner fresh probe and current
  proof, and recovers only through the existing persisted-intent path. It
  correctly treats a `row_pending` intent with no row as blocked
  `ErrEvalKeyRowMissing` rather than retrying creation. It also leaves lookup,
  owner, status, binding, and whitelist conflicts blocked for repair, with no
  automatic retry, revoke, rebind, or reenable.
- Source check: `Initialize` persists a row-pending intent before API-Key
  creation and can save the ready intent before runtime installation fails
  ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L420-L456),
  [eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L517-L530)).
  Its existing-intent path calls `adoptExistingRow`; a missing row returns
  `ErrEvalKeyRowMissing` and does not create a replacement
  ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L399-L417),
  [eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L480-L500)).

- File: `research/host-runtime-wiring-plan.md`
- Issue: The observation lease description could be read as holding the
  coordinator mutex during `Observe`.
- Fix: The revised algorithm consumes the ticket and increments the in-flight
  count while locked, calls `Observe` without the mutex, and publishes under
  lock only if the same generation remains active. Abort waits for leased
  observations before invalidating the session.
- Source check: `EvalTopologyCollector.Observe` has its own mutex around the
  duplicate check and observation insertion ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L525-L571));
  the coordinator lease can therefore guard lifecycle state without holding
  its mutex across collector work.

- File: `research/host-runtime-wiring-plan.md`
- Issue: Abort or process cleanup might revoke a Key or erase persistent
  provisioning state.
- Fix: The revised plan limits abort to ephemeral session/ticket state and
  explicitly preserves the intent and owned row. Only the separately reviewed
  uninstall path may revoke the Key.

## Findings (not fixed)

### Completion-pending retry must revalidate before readiness

- File: `research/host-runtime-wiring-plan.md`, `completion_pending` evidence
  write retry.
- Issue: The plan says to preserve the owner-bound completion record and proof
  for an idempotent evidence-write retry, without reopening observations or
  rerunning `Initialize`. It does not explicitly require that retry to reload
  the currently installed artifact and verify its exact digest/revision
  bindings, confirm the proof remains unexpired and matches live security
  configuration, and reverify the exact owned row before writing evidence or
  allowing readiness. The source proof validator rejects expiry and fingerprint
  drift ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L625-L670));
  `Load` verifies the intent/row and live proof ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L533-L566)).
- Required refinement: before an evidence-only retry, require the same actor,
  reload the protected artifact bytes, compare ID and digest to the retained
  completion record, validate current runtime/source/deployment/config
  revisions, validate the proof's MAC/fingerprint/expiry, and verify the exact
  owned row through the strict provisioner load path. Atomically write the
  session evidence only after these checks; recompute readiness from current
  artifact, session evidence, proof, Key identity, manager storage, and local
  listener state. Do not call `Initialize` in this evidence-only retry and do
  not issue a model request. If the artifact or proof is stale, keep readiness
  false and require a fresh same-owner probe and explicit provisioning update.
- Why not fixed: this is a readiness/recovery contract, not a mechanical source
  correction.

### Reconstruct recovery ownership after process restart

- File: `research/host-runtime-wiring-plan.md`, `provisioning_recovery_required`
  and process-restart behavior.
- Issue: Probe tickets and sessions are intentionally ephemeral, but the plan
  does not state how the “only that actor” recovery restriction is restored
  after restart. The provisioning intent persists `OwnerUserID`; the in-memory
  replacement reservation does not. If the runtime forgets to rebuild its
  recovery reservation from the intent, another administrator could begin a
  recovery probe after restart, only to fail later at `Initialize` with an owner
  conflict. This would violate the stated actor boundary and expose a confusing
  partial workflow.
- Required refinement: when startup finds an existing provisioned/recovery
  intent, derive the recovery owner restriction from its persisted
  `OwnerUserID`, not lost session memory. A restarted `completion_pending`
  state with missing session evidence must stay not-ready and be recoverable
  only by that owner. If current artifact bytes or proof freshness cannot be
  matched to a persisted completion record, require a fresh same-owner probe;
  do not infer the missing evidence from the provisioning intent alone. Keep
  row-missing and identity-conflict cases blocked for explicit repair. For
  failures that left no provisioning intent or owned row, retain the normal
  unprovisioned initialization path without claiming a prior owner was
  persisted.
- Why not fixed: the host startup/recovery coordinator has not been
  implemented, and this assignment is read-only.

## Readiness and Guard Identity

The plan keeps these states distinct. It makes `BootstrapRuntime` a recovery
and identity-guard operation, leaves startup `not-ready`, and assembles
evaluation readiness from provisioner proof validity, protected inventory and
session evidence, manager storage state, and local gateway availability. This
matches the source: `BootstrapRuntime` can return a nonnil identity while
installing identity-only guards on missing-row, owner/status/binding/IP, token,
or proof failures ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L569-L624)).
That returned identity must never be treated as proof of evaluation readiness.

## Verification

- Plan/source review: pass for partial-intent recovery semantics, exact-row
  recovery, missing-row blocking, abort preservation, observation-lease locking,
  and guard/readiness separation. Candidate verdict remains
  **NEEDS-REFINEMENT** for the two continuity findings above.
- Lint, type-check, and tests: not run; no product source was changed.
- Context: no complete injected-hook artifact was supplied in this dispatch;
  direct child-side reads were used for the active task artifacts and pinned
  source. Metadata patch 004's two gates are reported passed by the parent; the
  metadata reader plan was not reviewed.
- No Docker deployment, approved live target, topology probe, paid evaluation,
  or external application action was performed.
