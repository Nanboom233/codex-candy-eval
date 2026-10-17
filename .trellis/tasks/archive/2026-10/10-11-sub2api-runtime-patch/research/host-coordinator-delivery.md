# Host Topology Coordinator Delivery

Date: 2026-10-12. Scope: isolated Sub2API patch 006.

## Context and Ownership

The native SubagentStart hook identified this task and provided a saved full context artifact.
The visible hook and some tool outputs were truncated. The saved artifact and direct task documents supplied the context.
The implementer also loaded the implement manifest, its four documents, PRD, design, and current execution plan.
The frozen primary contract is `research/host-coordinator-plan.md`, not the oversized historical runtime plan.
Directory, persistence, error, logging, and run-state guides were read directly.
Context mode was reported to main before implementation. No product review PASS is claimed here.

Changed files:

- `dist/sub2api-build/workspace/backend/internal/service/eval_topology_runtime.go`.
- `dist/sub2api-build/workspace/backend/internal/service/eval_topology_runtime_test.go`.
- `patches/sub2api/v0.2.15/006-codex-eval-topology-runtime.patch`.
- This delivery report.

No existing generated source, core, patches 001-005, upstream lock, build helper, route, middleware, Wire, lifecycle, UI, or installer changed.
Unrelated `.codex/` and `.devin/` edits remain untouched.
No `.env`, CPA deployment, live/paid request, container operation, commit, or push was used.
Isolated ignored replay and check artifacts were created under `dist/sub2api-build/`.

## Exact Final API and Exported Types

```go
func NewEvalTopologyRuntime(EvalTopologyRuntimeDeps) (*EvalTopologyRuntime, error)
func (*EvalTopologyRuntime) Begin(context.Context, int64, EvalProbeBeginRequest) (EvalProbeBeginResult, error)
func (*EvalTopologyRuntime) Observe(*gin.Context) (bool, error)
func (*EvalTopologyRuntime) Complete(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Abort(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Status(context.Context, int64, string) (EvalProbeSessionStatus, error)
```

The following declarations are copied from final source. All four case members always encode.
Every successful status collection is a JSON array, including empty collections.
`EvalTopologyRuntime` exports no public test hook or global status surface.

```go
const EvalSessionEvidenceSchemaVersion = 1

type EvalProbeCaseDTO struct {
	Kind      EvalProbeKind `json:"kind"`
	IngressID string        `json:"ingress_id"`
	Header    string        `json:"header"`
	Source    string        `json:"source"`
}
type EvalProbeBeginRequest struct {
	ArtifactID string `json:"artifact_id"`
	SHA256     string `json:"sha256"`
}
type EvalProbeExternalTicket struct {
	Case             EvalProbeCaseDTO `json:"case"`
	Origin           string           `json:"origin"`
	RunnerLocationID string           `json:"runner_location_id"`
	Ticket           string           `json:"ticket"`
	ExpiresAt        int64            `json:"expires_at"`
}
type EvalProbeBeginResult struct {
	SessionID string                    `json:"session_id"`
	ExpiresAt int64                     `json:"expires_at"`
	Cases     []EvalProbeExternalTicket `json:"cases"`
}
type EvalProbeSessionStatus struct {
	SessionID     string             `json:"session_id,omitempty"`
	State         string             `json:"state"`
	ExpiresAt     int64              `json:"expires_at,omitempty"`
	RequiredCases []EvalProbeCaseDTO `json:"required_cases"`
	ObservedCases []EvalProbeCaseDTO `json:"observed_cases"`
	PendingCases  []EvalProbeCaseDTO `json:"pending_cases"`
	ErrorCode     string             `json:"error_code,omitempty"`
	APIKeyID      int64              `json:"api_key_id,omitempty"`
}
type EvalSessionEvidence struct {
	SchemaVersion          int                `json:"schema_version"`
	SessionID              string             `json:"session_id"`
	OwnerUserID            int64              `json:"owner_user_id"`
	ArtifactID             string             `json:"artifact_id"`
	InventorySHA256        string             `json:"inventory_sha256"`
	SecurityFingerprint    string             `json:"security_fingerprint"`
	Version                string             `json:"version"`
	BuildType              string             `json:"build_type"`
	DeploymentID           string             `json:"deployment_id"`
	SourceRevision         string             `json:"source_revision"`
	RuntimePatchSHA256     string             `json:"runtime_patch_sha256"`
	ImageIndexDigest       string             `json:"image_index_digest"`
	PlatformManifestDigest string             `json:"platform_manifest_digest"`
	ImageConfigDigest      string             `json:"image_config_digest"`
	BinarySHA256           string             `json:"binary_sha256"`
	IngressRevision        string             `json:"ingress_revision"`
	CompletedAt            int64              `json:"completed_at"`
	ProofExpiresAt         int64              `json:"proof_expires_at"`
	ProofSHA256            string             `json:"proof_sha256"`
	RequiredCases          []EvalProbeCaseDTO `json:"required_cases"`
}
type EvalTopologyProvisioner interface {
	Initialize(context.Context, int64, EvalIngressTopology) (*EvalKeyIdentity, error)
	Load(context.Context) (*EvalKeyIdentity, error)
}
type EvalDeploymentCurrentFactsProvider interface {
	CurrentEvalDeploymentFacts(context.Context) (EvalDeploymentCurrentFacts, error)
}
type EvalBoundListener struct {
	Listener net.Listener
	Scheme   string
}
type EvalBoundListenerSource interface {
	ActualEvalListener() (EvalBoundListener, error)
}
type EvalProbeLocalTicket struct {
	Case   EvalProbeCaseDTO
	Ticket string
}
type EvalTopologyLocalProbe interface {
	ProbeLoopback(context.Context, EvalBoundListener, []EvalProbeLocalTicket) error
}
type EvalSessionEvidenceWriter interface {
	CheckAvailable(context.Context) error
	Save(context.Context, EvalSessionEvidence) error
}
type EvalTopologyRuntimeDeps struct {
	Config         *config.Config
	Store          *EvalKeyStore
	Provisioner    EvalTopologyProvisioner
	CurrentFacts   EvalDeploymentCurrentFactsProvider
	EvidenceWriter EvalSessionEvidenceWriter
	ListenerSource EvalBoundListenerSource
	LocalProbe     EvalTopologyLocalProbe
}

```

## Runtime Contracts

The coordinator requires config, store, provisioner, facts provider, writer, listener source, and loopback probe.
Construction rejects nil and typed-nil ports with a fixed sanitized unavailable error.
The four production port implementations remain absent. Test fakes are not production readiness.
The listener source must provide the actual server TCP listener with server-owned `http` or `https` scheme.
The coordinator rejects missing listeners, unsupported schemes, non-TCP addresses, and invalid bound ports.
It never infers listener identity from config or a request URL.

One mutex owns the current session, generation, owner reservation, state, deadline, digests, observations, and lease count.
Ticket owner, generation, expiry, and local/external class are bound through their owning session and exact native case.
Only SHA-256 digests and use state are retained. Raw tickets exist only in detached outgoing probe/result values.
Providers, storage calls, probe calls, collector Observe/Seal, and provisioner calls execute outside that mutex.
A replaceable notification channel permits unlocked, context-aware draining without a timer goroutine.

Begin establishes a fixed deadline at entry plus ten minutes, capped by inventory expiry.
Replacement after abort, invalidation, or recovery inherits the unexpired reservation deadline.
Retry and abort never extend it. Completed sessions can be replaced immediately.
An expired session needs a fresh Begin. Persisted intent ownership remains authoritative.
A provisional wrong-owner restart reservation is released before returning the generic conflict.
Another actor cannot release an existing valid owner reservation.

Begin reloads independent facts and exact protected inventory before collecting observations.
It derives required native cases and rejects ticket/forwarding-header collisions.
It checks writer availability before issuing tickets and returns only external case tickets.
The injected local probe receives only direct-allowed tickets and must use the captured listener and Observe hook.
Every local observation is required before active publication. Failure returns no external ticket.
Valid unused inventory runners remain accepted through the frozen reader.

Observe scans raw header keys case-insensitively for exactly one `X-Codex-Eval-Probe-Ticket` key and one value.
It requires one canonical unpadded base64url 32-byte ticket and exact `POST /v1/responses`.
Duplicate values, case-variant duplicate keys, malformed, unknown, expired, consumed, and wrong-class tickets fail closed.
No marker returns `(false, nil)`. Request labels cannot choose case or owner.
Cancellation before leasing does not consume. Consumption and lease acquisition happen under the mutex.
Native observation runs unlocked and settles after disconnect. Publication requires the same valid generation.
Abort and invalidation discard late observations. Collector clearing waits for the last lease.

Abort blocks new leases and waits outside the mutex for existing leases and building work.
A canceled Abort retains aborting state for the same actor to resume.
It never mutates Key or provisioning intent. Completing is busy.
Pending evidence context survives Abort. Recovery-required and completed states are owner-idempotent.
Status is context-bound and actor/session-bound. Unknown IDs and wrong owners return identical conflicts.
Only the owner can see the sanitized APIKeyID. Status exposes no ticket, digest, token, lookup, or global readiness.

Complete requires full coverage, no lease, the same owner, and remaining time.
It reloads exact inventory and every independent current fact before writer availability, token load, seal, and initialization.
Transient provider/writer/token preflight failure restores active state without new tickets.
Drift, expiry, Seal/topology failure, and cancellation after Seal require fresh owner-bound proof.
Initialize is invoked once per claimed completion. Every error may follow durable effects and enters recovery-required.
No automatic Initialize retry, duplicate creation, rebinding, reenabling, revocation, or model request is added.
Missing rows retain `errors.Is(err, ErrEvalKeyRowMissing)` behind `owned_key_repair_required`.
Native row-pending recovery cannot recreate the missing row.

Only successful Initialize plus an evidence Save failure enters completion-pending.
Cancellation after successful Initialize also retains pending evidence context without automatic Save retry.
Pending retry is single-flight and bounded by the original session deadline and proof expiry.
It reloads facts, exact inventory, token, intent, and the exact owned identity through native Load.
It checks original owner/Key ID/Lookup, token-derived lookup, ready intent, whitelist, proof digest, MAC, fingerprint, and freshness.
It invokes only idempotent Save. It never Initialize/Seal/probe/model-calls on that retry.
The writer receives detached evidence cases. Successful completion clears collector, proof, tickets, and private identity.
Fixed error strings preserve only cancellation/deadline and missing-row identities. Dependency text stays private.
Status, Abort, and Complete recheck cancellation after mutex acquisition before claiming or inspecting session state.

## Producer and Consumer Mapping

| Value or port | Producer and remaining owner | Consumer |
| --- | --- | --- |
| Protected deployment inventory | Later privileged operator/installer after independent deployment inspection | Frozen patch-005 reader, then Begin/Complete/retry |
| Current deployment facts | Later trusted runtime/installer provenance provider | Patch-005 comparisons and coordinator reloads |
| Actual bound TCP listener and scheme | Later server listener capture and TLS policy integration | Begin and injected real loopback probe |
| Local probe requests | Later real loopback HTTP transport through the captured listener | Same future Observe hook |
| External tickets and exact destinations | This coordinator, from accepted inventory | Later external runner and authenticated admin Begin response |
| Authenticated actor ID | Later normal administrator middleware, using AuthSubject.UserID | Begin/Complete/Abort/Status |
| Native observations and topology proof | Frozen native collector, invoked by coordinator | Native provisioner and domain evidence digest |
| Persisted ownership and exact Key identity | Frozen store and provisioner | Begin recovery and evidence-only retry |
| EvalSessionEvidence domain DTO | This coordinator after successful Initialize | Later protected atomic evidence writer and readiness consumer |
| Evidence storage availability and idempotent Save | Later writer; fixed disk schema/path and protected persistence remain its scope | Begin preflight, completion preflight, Save-only retry |
| Full readiness | Later runtime integration composing independent facts, inventory, proof, session evidence, executor, and lifecycle | Later paid admission and status decision |

Patch 006 defines no on-disk evidence schema/path and provides no evidence reader or writer implementation.
The later writer must validate protected storage without creating files in CheckAvailable.
Identical `(SessionID, InventorySHA256, ProofSHA256)` writes must be idempotent. Conflicting content must fail.
The domain record stores no raw proof, raw token, ticket, or Key lookup.
ProofSHA256 hashes exact `json.Marshal(proof)` bytes with lowercase SHA-256.
Restart without a durable completion binding requires fresh same-owner probe coverage.
Inventory or native proof validity alone never establishes readiness.

The later outer HTTP marker guard must reject marked wrong routes/methods before Gin redirects, auth, billing, or forwarding.
Admin routes, authentication, no-store response headers, real probe transport, providers, Wire/bootstrap, executor gating, and shutdown remain integration work.

## Checks and Test Evidence

Toolchain: Go 1.27.2 on macOS arm64. Both new files were formatted with gofmt.
The coordinator adds 15 top-level unit tests with table-driven failure and identity coverage.
Tests use native collector and real native Key provisioner over existing in-memory repository doubles.
No live host, external runner, upstream model, billing system, or Docker fixture is invoked.
Channel barriers and bounded waits verify single-flight completion, unlocked collector calls, and abort/drain/resume.
Cleanup releases barriers and waits for workers. Clock-hook tests are serial.

Final checks in the maintained and pristine replay backend trees:

- `go test -race ./internal/service -run '^TestEvalTopologyRuntime' -count=1`: PASS, compilation only, no coordinator tests without unit tag.
- `go test -race -tags unit ./internal/service -run '^(TestEvalTopology|TestEvalDeploymentEvidence)' -count=1`: PASS, meaningful native-topology/reader/coordinator race coverage.
- `go vet ./internal/service`: PASS.
- `go vet -tags unit ./internal/service`: PASS.
- `go vet ./internal/codexeval`: PASS.
- `go test -race ./internal/codexeval -count=1`: PASS.

Final root checks:

- `go vet ./...`: PASS.
- `go test -race ./... -count=1`: PASS, including internal/eval and original plugin.

Focused coordinator-only unit race also passed after the last source/test correction.
Tests cover absent/typed-nil dependencies, Begin failures, writer preflight, canonical ticket validation, actor/restart binding, DTO isolation, and expiry.
They cover transient versus drift preflight, native partial provisioning before/after row creation, cancellation, row-pending missing-row repair, and strict Save-only retry.
They cover owner/ID/lookup/token/intent/proof/status/group/IP mismatches with no duplicate or automatic mutation.
They cover building abort, canceled Begin, leased disconnect, late expiry settlement, concurrent Complete, and detached writer DTOs.

The documented pristine upstream full-service unit race failures remain unresolved.
That full suite was not relabeled or repaired. Focused service PASS does not resolve the blocker.
Default focused compilation is not behavior coverage.
No Linux execution, release binary, frontend build, Wire generation, Docker lifecycle, live topology, browser, paid evaluation, or actual charge acceptance ran.
These remain later integration/full-scope checks. This delivery is not migration completion or live readiness.

## Final Patch and Pristine Replay

Patch 006 contains only the two new coordinator files, in bare unified `a/backend/` and `b/backend/` format.
Both new files were asserted after strict sequential `git apply --check --whitespace=error-all` and `git apply --whitespace=error-all` for 001 through 006.
The replay has its own isolated Git root, preventing nested-parent new-file skipping.
All 28 maintained evaluation-core files were copied unchanged after applying the patch series.
The final replay follows the last product correction.

Final 006 SHA-256: `e9684efbcd470a95e4315ac66ef4f2f7d4ac8e3c979f26ff464cf71234131a39`.

Fresh pinned replay: `dist/sub2api-build/scratch-coordinator-final-0tj2fbh_`.
Pinned source commit: `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.

All 3260 backend files, 931 frontend files, and 28 core files matched the maintained tree byte-for-byte.
The frozen input checksums are:

| Input | SHA-256 |
| --- | --- |
| 001-codex-eval-auth.patch | `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b` |
| 002-codex-eval-scheduling.patch | `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5` |
| 003-codex-eval-cancellation.patch | `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4` |
| 004-codex-eval-target-adapter.patch | `850af726f80bb2f379e0cf9b3f395f28e9d8e89356106f0643c39ee5558f4088` |
| 005-codex-eval-deployment-evidence.patch | `3402198f3f75e73ebf46e748ef4a289001a13d0632692c3bb493a630055650cb` |
| 006-codex-eval-topology-runtime.patch | `e9684efbcd470a95e4315ac66ef4f2f7d4ac8e3c979f26ff464cf71234131a39` |

Core raw-digest tree SHA-256: `3eabe2af794582dc8e41dfd637e57fc55940b9e57bc13e9606c1d748f9a8a7ba`.
Core hexadecimal-ASCII tree SHA-256: `b8477eb71da52e832aefa690b7bdf4fa02c257a792f1669ff3cbb5126ec67af5`.
Both algorithms match their frozen values. Their different encodings are not drift.

All three lockfiles match pristine source, maintained workspace, and replay:

| File | SHA-256 |
| --- | --- |
| backend/go.mod | `809e6806dc248bd44cea4198ddc21fd6ed04f80003300c3668367b91bae319a8` |
| backend/go.sum | `b17668c6300c04027bd397650e60e28e514dc4abb9ab113e3cea634fdd708109` |
| frontend/pnpm-lock.yaml | `1d8e54cc8e1f1bd3e7593805f607d9314ff81dab98ab5a9a6396475e148b26ce` |

Machine-readable replay evidence is at `dist/sub2api-build/coordinator-implement-replay.json` and the replay's `coordinator-replay-evidence.json`.
Final check logs use `/tmp/codex-eval-maintained-check-{0..5}.log`, `/tmp/codex-eval-replay-check-{0..5}.log`, and `/tmp/codex-eval-root-*-final.log`.

## Final Review Registration

Specification and distinct independent quality reviews passed without product changes.
Read [specification](host-coordinator-spec-review.md) and [quality](host-coordinator-quality-review.md).
The final quality replay is `dist/sub2api-build/scratch-coordinator-quality-PMsqqJ`.
Main captured [the coordinator contract](../../../spec/backend/host-coordinator.md).
Final006 remains `e9684efbcd470a95e4315ac66ef4f2f7d4ac8e3c979f26ff464cf71234131a39`.

## Required Handoff

Main registers this persisted report in check.jsonl, then dispatches the specification trellis-check.
After specification PASS, main dispatches a distinct quality/replay trellis-check.
A product correction requires a new patch identity and pristine replay after its last change.
Main updates the source-backed coordinator spec only after the final local product gate.
The next integration chunk must implement all four mandatory production ports and preserve the producer/consumer boundaries above.
No local readiness default can substitute for missing current facts, protected evidence persistence, actual listener capture, or real probes.
