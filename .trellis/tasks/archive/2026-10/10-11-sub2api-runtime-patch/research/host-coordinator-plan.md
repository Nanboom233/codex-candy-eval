# Research: Host Topology Session Coordinator

- Query: Freeze the source-backed patch-006 boundary for actor-bound topology sessions, one-use observation tickets, exact-owner provisioning recovery, and session-evidence retry.
- Scope: internal
- Date: 2026-10-12

## Findings

### Decision and patch boundary

Patch 006 is one service-level coordinator and adjacent tests in the maintained Sub2API workspace:

- `backend/internal/service/eval_topology_runtime.go`
- `backend/internal/service/eval_topology_runtime_test.go`
- Export as `patches/sub2api/v0.2.15/006-codex-eval-topology-runtime.patch`

It owns actor/session state, ticket issuance and consumption, collector observations, proof sealing, provisioner calls, and the result of writing session evidence. It adds no routes, middleware, runner, listener wiring, Wire/bootstrap, manager lifecycle, UI, installer, or deployment readiness decision.

Four injected dependencies have no production implementation: current-facts provider, protected evidence writer, bound-listener source, and loopback HTTP probe. Require these plus config, store, and provisioner at construction; nil/unavailable dependencies fail closed before `Initialize`, with no permissive defaults. Test fakes cover transitions only. Production providers and readiness composition belong to later integration.

### Proposed service API

The authenticated handler supplies the actor ID. No request field may choose or override the owner.

```go
type EvalProbeCaseDTO struct { Kind EvalProbeKind `json:"kind"`; IngressID string `json:"ingress_id"`; Header string `json:"header"`; Source string `json:"source"` }
type EvalProbeBeginRequest struct { ArtifactID string `json:"artifact_id"`; SHA256 string `json:"sha256"` }
type EvalProbeExternalTicket struct { Case EvalProbeCaseDTO `json:"case"`; Origin string `json:"origin"`; RunnerLocationID string `json:"runner_location_id"`; Ticket string `json:"ticket"`; ExpiresAt int64 `json:"expires_at"` }
type EvalProbeBeginResult struct { SessionID string `json:"session_id"`; ExpiresAt int64 `json:"expires_at"`; Cases []EvalProbeExternalTicket `json:"cases"` }
type EvalProbeSessionStatus struct { SessionID string `json:"session_id,omitempty"`; State string `json:"state"`; ExpiresAt int64 `json:"expires_at,omitempty"`; RequiredCases []EvalProbeCaseDTO `json:"required_cases"`; ObservedCases []EvalProbeCaseDTO `json:"observed_cases"`; PendingCases []EvalProbeCaseDTO `json:"pending_cases"`; ErrorCode string `json:"error_code,omitempty"`; APIKeyID int64 `json:"api_key_id,omitempty"` }
type EvalTopologyRuntimeDeps struct { Config *config.Config; Store *EvalKeyStore; Provisioner EvalTopologyProvisioner; CurrentFacts EvalDeploymentCurrentFactsProvider; EvidenceWriter EvalSessionEvidenceWriter; ListenerSource EvalBoundListenerSource; LocalProbe EvalTopologyLocalProbe }

func NewEvalTopologyRuntime(EvalTopologyRuntimeDeps) (*EvalTopologyRuntime, error)
func (*EvalTopologyRuntime) Begin(context.Context, int64, EvalProbeBeginRequest) (EvalProbeBeginResult, error)
func (*EvalTopologyRuntime) Observe(*gin.Context) (handled bool, err error)
func (*EvalTopologyRuntime) Complete(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Abort(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Status(context.Context, int64, string) (EvalProbeSessionStatus, error)
```

Required ports:

```go
type EvalTopologyProvisioner interface { Initialize(context.Context, int64, EvalIngressTopology) (*EvalKeyIdentity, error); Load(context.Context) (*EvalKeyIdentity, error) }
type EvalDeploymentCurrentFactsProvider interface { CurrentEvalDeploymentFacts(context.Context) (EvalDeploymentCurrentFacts, error) }
type EvalBoundListener struct { Listener net.Listener; Scheme string }
type EvalBoundListenerSource interface { ActualEvalListener() (EvalBoundListener, error) }
type EvalProbeLocalTicket struct { Case EvalProbeCaseDTO; Ticket string }
type EvalTopologyLocalProbe interface { ProbeLoopback(context.Context, EvalBoundListener, []EvalProbeLocalTicket) error }
type EvalSessionEvidenceWriter interface { CheckAvailable(context.Context) error; Save(context.Context, EvalSessionEvidence) error }
```

The listener must be the live TCP listener used by the server, with its actual bound address and server-owned `http` or `https` scheme. Never infer it from config or accept a request URL. The loopback probe makes real requests through that listener and the same `Observe` hook; it cannot call the collector directly.

Patch 006 owns the domain DTO; a later writer owns the on-disk schema/path and protected atomic implementation. Use Unix seconds, explicit tags for every current fact, and all four case fields. `ProofSHA256` is lowercase SHA-256 of the exact `json.Marshal(proof)` bytes. Persist no proof, ticket, raw token, or Key `Lookup` in this evidence record:

```go
const EvalSessionEvidenceSchemaVersion = 1
type EvalSessionEvidence struct {
    SchemaVersion int `json:"schema_version"`
    SessionID string `json:"session_id"`
    OwnerUserID int64 `json:"owner_user_id"`
    ArtifactID string `json:"artifact_id"`
    InventorySHA256 string `json:"inventory_sha256"`
    SecurityFingerprint string `json:"security_fingerprint"`
    Version string `json:"version"`
    BuildType string `json:"build_type"`
    DeploymentID string `json:"deployment_id"`
    SourceRevision string `json:"source_revision"`
    RuntimePatchSHA256 string `json:"runtime_patch_sha256"`
    ImageIndexDigest string `json:"image_index_digest"`
    PlatformManifestDigest string `json:"platform_manifest_digest"`
    ImageConfigDigest string `json:"image_config_digest"`
    BinarySHA256 string `json:"binary_sha256"`
    IngressRevision string `json:"ingress_revision"`
    CompletedAt int64 `json:"completed_at"`
    ProofExpiresAt int64 `json:"proof_expires_at"`
    ProofSHA256 string `json:"proof_sha256"`
    RequiredCases []EvalProbeCaseDTO `json:"required_cases"`
}
```

`CheckAvailable` is read-only: validate protected storage and effective write permission without creating files. Call at Begin before tickets, and at Complete/evidence retry after live checks but before token load, Seal, or Initialize. Failure returns `storage_unavailable` and preserves the active session to its fixed deadline. Only a later `Save` failure enters `completion_pending`; identical `(SessionID, InventorySHA256, ProofSHA256)` writes must be idempotent and conflicting content rejected.

### Begin and dependency ordering

1. Constructor rejects nil config/store and nil or typed-nil ports with a fixed unavailable error. Detect typed nil through a private nil-able-kind reflection helper. Each call checks actor, artifact/digest, standard run mode, and context.
2. Under the mutex expire stale state, enforce reservation rules, allocate session ID/generation, set the fixed deadline, and publish `building`. Unlock before all I/O, providers, probes, collector calls, and provisioner calls.
3. Resolve the actual live TCP listener and load intent. A persisted positive `OwnerUserID` is authoritative. If a restart has no in-memory owner reservation and this request's provisional `building` reservation conflicts with persisted ownership, release only that unverified reservation before returning the generic owner error; the real owner can Begin immediately. If a valid owner reservation already exists, a different actor's request must leave it unchanged. Neither error reveals intent or owner data.
4. Load independent current facts and exact deployment evidence using the request artifact ID/SHA. Retain the reader's artifact ID and exact-byte digest. Any missing fact/read failure clears partial session state and returns no ticket.
5. Build the collector from inventory loopback sources and ingress peers. Require its fingerprint to match, 1–256 cases, and no case-insensitive collision with recognized forwarding headers. Call read-only writer `CheckAvailable` before minting tickets.
6. Mint one 32-byte random canonical base64url ticket per server-derived case. Retain only SHA-256 digests with generation, owner, exact case, class, expiry, and use state. Map external cases to validated destinations; valid unique unused runner descriptors remain allowed.
7. In `building`, send only `direct_allowed` tickets through the real probe to the captured listener and same observation hook. Require every local observation. Probe failure, cancellation, drift, or incomplete cases clears all tickets/collector and returns no external ticket.
8. Publish `active` only for the same generation. Return external tuples/origins/runner IDs/tickets with expiry `min(sessionStartedAt + 10m, inventory.ExpiresAt)`. The proof's separate maximum age is 24 hours. Future responses use `Cache-Control: no-store`.

The fixed session/reservation deadline is `min(Begin start + 10m, inventory expiry)`; evidence retry is also capped by proof expiry. Abort/failure/retry never extends it. Begin during building/active/completing/aborting or pending retry is busy for the owner and generic owner-conflict for others. The same actor may replace aborted/invalidated/recovery-required state before expiry, inheriting its deadline; completed can be replaced immediately. Expiry releases an unpersisted reservation, while persisted intent retains its owner binding. Expired evidence cannot be retried; a fresh same-owner Begin starts a new bounded session.

Begin cancellation before reservation changes nothing. Cancellation or dependency/read failure after reservation invalidates that generation and clears tickets/collector; the fixed owner reservation remains until deadline. No ticket is returned.

### State and observation leases

One mutex owns the current session, monotonically increasing generation, fixed owner-reservation deadline, ticket digests, observed cases, state, and in-flight count. Published session data is otherwise immutable. Never hold the mutex across I/O or external calls. Wait on a replaceable notification channel after unlocking. Context cancellation before an operation claims state makes no change; cancellation after a claim follows that operation's rules below.

```text
building -> active -> completing -> completed
    |          |             |-> completion_pending
    |          |             |-> provisioning_recovery_required
    |          |             |-> invalidated
    |                        completion_pending --Complete retry--> completing
    |                        completing --Save success--> completed
    |                        completing --failure--> completion_pending
    |          |-> failed -> aborting -> invalidated
    |          |-> aborting -> invalidated
    |          |-> invalidated
```

- `building` leases local observations only; `active` leases external observations only. Other states accept no ticket. The header is exactly `X-Codex-Eval-Probe-Ticket`.
- With no case-insensitive match for the ticket header, `Observe` returns `(false, nil)`. Iterate raw header-map keys with `strings.EqualFold`; require one matching key and exactly one nonempty value. Reject repeated values, duplicate case variants, malformed/canonicalization failures, unknown/expired/replayed tickets, and anything except exact `POST /v1/responses`. Ignore caller case/owner labels.
- Before leasing, reject an already-cancelled request context without consuming the ticket. Under lock validate state/generation/expiry, consume its digest, and increment in-flight; unlock for `collector.Observe`. Once leased, always settle the observation and release the lease even if the client disconnects. Publish only if the same generation still permits it; discard late results after invalidation/abort.
- `Abort` on `building`/`active` enters `aborting`, blocks leases, and waits unlocked. On cancellation it remains `aborting` for the same actor to resume. Success clears tickets/collector only and retains the owner reservation to its original deadline; it never changes Key or intent. Abort during `completing` returns busy. Abort during `completion_pending` returns that state and preserves the only evidence retry context. Recovery-required/completed states are idempotent no-ops for their owner.
- `Status(ctx, actorID, sessionID)` is the only patch-006 status surface: per-session, owner-checked, read-only. Unknown IDs and wrong owners return the same `session_owner_conflict` without details; canceled context returns its context error. It exposes no ticket/hash/Lookup; APIKeyID is present only for the owner. There is no global status API.
- Check context and expiry before leasing and each state transition. Expiry/config drift invalidates the generation, blocks new leases, drains existing leases, then clears the collector. No timer goroutine is needed.

The future outer HTTP marker guard must reject marked wrong methods/paths before Gin redirect, auth, billing, or forwarding. That guard is outside patch 006. The coordinator does not write HTTP responses or authenticate callers.

### Complete, partial provisioning, and evidence retry

`Complete` claims `completing` only for the owner/session with full coverage, no lease, and time remaining. `building` returns `topology_unverified`; completing/aborting is `session_busy`; invalidated/expired is `session_expired`; recovery-required needs fresh Begin; completed is sanitized/idempotent. Pending evidence claims single-flight retry before session/proof expiry. Wrong/unknown owner is generic `session_owner_conflict`. Success clears tickets, proof, collector, and private identity; only sanitized completed status remains.

Before `Initialize`, reload facts and exact inventory; compare digest, artifact, every fact, fingerprint, cases, generation, and expiry. Then call writer `CheckAvailable`, load the private token, seal, validate topology, and call Initialize once. No model request occurs. Drift/expiry and Seal/topology failure or cancellation after Seal invalidate and clear tickets/collector. A transient facts/writer/token preflight error or cancellation before Seal restores `active`, retaining observations and consumed-ticket state to the same deadline; tickets are never reissued. Check context immediately before Initialize. Once invoked, any error/cancellation may follow durable effects and enters recovery-required.

Treat every `Initialize` error, including cancellation, as potentially partially committed. Enter `provisioning_recovery_required`; preserve intent/row and the owner reservation, but clear tickets, collector, and sealed transient data. Do not auto-retry, revoke, delete, duplicate-create, rebind, reenable, or issue a paid request. Complete in this state returns the fixed recovery error; the same owner must Begin a fresh proof. Existing Initialize adoption verifies the exact row. A missing row under `row_pending` returns `ErrEvalKeyRowMissing` and cannot recreate it; owner/status/lookup/group/IP conflicts need operator repair.

Reserve `completion_pending` only for successful Initialize plus failed Save, including a Save canceled after Initialize. Retain original owner/identity, artifact/digest, sealed proof/digest; clear tickets. Until the fixed session/proof deadline, retry preflights writer, facts/inventory/intent, calls Load, and binds its owner/Key ID/Lookup to the original Initialize identity. Verify `EvalDeriveLookup(privateToken) == Lookup`, ready intent, proof digest/MAC, fingerprint, topology, and expiry. Retry idempotent Save only. Preflight/Save failure restores pending data; mismatch stays not-ready with no Lookup/error. Never Initialize, Seal, reissue tickets, or run a model. At expiry return `session_expired`; a fresh same-owner Begin is required.

After restart, bind any new probe to `LoadIntent().OwnerUserID`. Intent/proof or `EvalRuntimeProofValid` alone never proves session evidence/readiness. No durable completion binding exists, so missing evidence requires a fresh complete same-owner probe; patch 006 does not invent one.

Return only fixed error codes such as `session_busy`, `session_owner_conflict`, `session_expired`, `topology_unverified`, `storage_unavailable`, `provisioning_recovery_required`, `completion_pending`, and `owned_key_repair_required`. Preserve cancellation/deadline identity. Never expose paths, artifact contents, raw token, ticket/hash, provider messages, backend text, or Key `Lookup`. Status is session/provisioning state, not a global readiness boolean.

### Focused tests

- Reject nil/typed-nil ports; inject provider, listener, probe, and `CheckAvailable` failures at Begin/Complete/retry with no tickets or `Initialize`. Only `Save` failure after successful `Initialize` enters `completion_pending`.
- Cover actor/digest/intent binding, including restart with wrong-owner Begin followed immediately by the persisted owner's successful Begin; also prove another actor cannot release an existing valid reservation. Cover exact facts/cases, local-ticket omission, ticket replay/expiry, header variants, exact route, ignored labels, and header collision.
- Block `collector.Observe` through a private per-runtime test seam; verify unlocked calls, lease drain, abort resume, and stale-generation rejection. Inject provisioning/evidence failures; verify exact-row recovery, no duplicate Key or mutation, Save-only retry with identity binding, expiry/restart ownership, and sanitized Status. Use bounded waits; clock-hook tests are serial.

No live host, Docker, external runner, billing, or model request belongs in patch-006 tests.

### Handoff sequence

Implement only the two coordinator files and patch 006, reusing `Load`, `LoadIntent`, `LoadPrivateToken`, and `EvalDeriveLookup`. Run gofmt, focused unit race tests, service vet with/without `unit`, root vet/race, and separate spec/quality/replay reviews. Replay 001–006 on pinned commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`; assert both files, compare backend/frontend/core/lockfiles, verify patch 005 hash `3402198f3f75e73ebf46e748ef4a289001a13d0632692c3bb493a630055650cb`, and record patch-006 hash. This cannot establish live readiness; provider/listener/probe, routes, Wire/bootstrap, executor gate, and shutdown remain later integration.

## Files Found

- `dist/sub2api-build/workspace/backend/internal/service/eval_topology.go`: collector, required cases, observations, proof seal/validation.
- `dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go`: persisted owner intent, exact-row adoption, strict `Load`.
- `dist/sub2api-build/workspace/backend/internal/service/eval_deployment_evidence.go`: independent inventory facts and exact-byte digest.
- `dist/sub2api-build/workspace/backend/internal/service/eval_identity.go`: proof validity does not establish runtime readiness.
- `dist/sub2api-v0.2.15/backend/cmd/server/main.go`, `dist/sub2api-v0.2.15/backend/internal/server/http.go`: `ListenAndServe` hides the bound listener; the patched workspace has no listener capture.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/host-runtime-wiring-plan.md`: scope, integration ownership, and coordinator lifecycle.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/host-runtime-refinement-response.md`: refined lease/recovery findings and remaining operational caveats.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/deployment-evidence-reader-plan.md`: patch-005 input contract and artifact identity.

## Code Patterns

- `eval_topology.go:105-110,310-325,464-522`: native case tuple has `omitempty`; collector derives required cases and snapshot fingerprint.
- `eval_topology.go:525-571,574-670`: collector observation is mutex-protected; Seal requires full coverage and emits a 24-hour MAC proof checked against live fingerprint.
- `eval_key_runtime.go:196-241,374-456,480-566`: intent binds owner/Lookup and can partially persist; existing-row adoption and strict `Load` expose private `Lookup` internally.
- `eval_deployment_evidence.go:30-41,109-151`, `eval_identity.go:75-84`: inventory compares independently supplied facts/exact digest; runtime proof validity alone is not readiness.
- `main.go:121-129,171-179`, `http.go:112-125`: server uses `ListenAndServe`, not an exposed listener; `.trellis/spec/backend/run-state.md:52-97` requires I/O outside state locks and visible persistence errors.

## Related Specs

- `.trellis/spec/backend/deployment-evidence.md`
- `.trellis/spec/backend/quality-guidelines.md:194-294`
- `.trellis/spec/backend/run-state.md:52-97`
- `.trellis/spec/backend/host-target-adapter.md`

## External References

- None. Source review is against Sub2API v0.2.15 pinned at `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` and reviewed local patch-005 artifacts.

## Caveats / Not Found

- Patch 006 and its tests do not exist yet. No production facts provider, evidence writer/reader, actual listener source, local probe, handler, ticket guard, or readiness composition exists; required dependencies must fail closed before Initialize.
- Host `BuildInfo` provides only version/build type. `Load` also requires identity and token-derived Lookup binding during evidence retry; keep Lookup private.
- Inventory and installed proof do not establish deployment readiness. No live target, listener, browser, billing, model, or deployment check ran.
