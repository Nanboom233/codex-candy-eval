# Host Topology Coordinator Specification Review

Date: 2026-10-12. Verdict: PASS for patch 006 against the frozen coordinator contract.

## Context and Scope

The native hook supplied a saved full-context artifact. I read it in bounded sections, then read `check.jsonl`, the PRD, design, current execution plan, frozen coordinator plan and delivery, backend index, and the relevant directory, persistence, error, logging, run-state, quality, and deployment-evidence guides. Context mode: native full hook context, completed from the saved artifact.

Reviewed only:

- `dist/sub2api-build/workspace/backend/internal/service/eval_topology_runtime.go`
- `dist/sub2api-build/workspace/backend/internal/service/eval_topology_runtime_test.go`
- `patches/sub2api/v0.2.15/006-codex-eval-topology-runtime.patch`
- This review report

No product file required correction. No existing upstream source, evaluation core, patch 001-005, lockfile, task manifest, or unrelated `.codex/` or `.devin/` file was changed.

## Findings

### Fixed

None.

### Not Fixed

None. The production current-facts provider, protected evidence writer, actual listener source, and real loopback probe remain absent by design and are later integration owners. This patch rejects missing dependencies and does not claim readiness.

## API and Contract Review

The exported coordinator surface matches the frozen plan:

```go
func NewEvalTopologyRuntime(EvalTopologyRuntimeDeps) (*EvalTopologyRuntime, error)
func (*EvalTopologyRuntime) Begin(context.Context, int64, EvalProbeBeginRequest) (EvalProbeBeginResult, error)
func (*EvalTopologyRuntime) Observe(*gin.Context) (bool, error)
func (*EvalTopologyRuntime) Complete(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Abort(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Status(context.Context, int64, string) (EvalProbeSessionStatus, error)
```

Exported DTOs are `EvalProbeCaseDTO`, `EvalProbeBeginRequest`, `EvalProbeExternalTicket`, `EvalProbeBeginResult`, `EvalProbeSessionStatus`, and `EvalSessionEvidence`. Their JSON tags match the plan. Every case tuple encodes `kind`, `ingress_id`, `header`, and `source`; successful status values encode empty case arrays. The evidence DTO carries the exact inventory digest, independent deployment facts, completion/proof expiry, lowercase digest of exact marshaled proof bytes, and required cases. It carries no proof, ticket, raw token, or Key lookup.

Exported dependency types are `EvalTopologyProvisioner`, `EvalDeploymentCurrentFactsProvider`, `EvalBoundListener`, `EvalBoundListenerSource`, `EvalProbeLocalTicket`, `EvalTopologyLocalProbe`, `EvalSessionEvidenceWriter`, and `EvalTopologyRuntimeDeps`. Constructor nil and typed-nil checks cover all mandatory ports. The listener contract is an actual server-owned TCP listener plus `http` or `https` scheme; the implementation rejects unsupported schemes, non-TCP addresses, and invalid ports.

The reviewed behavior follows the frozen state and ownership rules: actor-bound owner reservations, release of only a provisional wrong-owner restart reservation, fixed ten-minute deadline capped by inventory expiry, digest-only retained tickets, exact method/path and case-insensitive header validation, consume-plus-lease under the mutex, unlocked native observation and drain, same-generation publication, single-flight completion, and context checks around state claims.

Completion reloads independent facts and the exact inventory before provisioning. It calls `CheckAvailable` before ticket creation and again before token loading, sealing, or initialization. Every `Initialize` error enters recovery-required while preserving the owner reservation; missing row identity remains discoverable through `errors.Is`. Only successful initialization followed by evidence failure enters completion-pending. Retry reloads facts, inventory, intent, token, and the exact owned identity, validates proof and lookup bindings, and invokes only `Save`. Status is owner/session-bound and sanitized. The coordinator adds no readiness flag or production HTTP/lifecycle wiring.

## Producer and Consumer Ownership

| Value | Producer and future owner | Consumer |
| --- | --- | --- |
| Protected deployment inventory | Privileged operator/installer after independent topology inspection | Patch-005 reader, then coordinator Begin/Complete/retry |
| Current deployment facts | Trusted runtime/installer provenance provider | Reader comparisons and coordinator reloads |
| Actual bound listener and scheme | Server listener capture and TLS integration | Begin and loopback probe |
| Local probe requests | Real loopback HTTP transport using the captured listener and same Observe hook | Coordinator collector |
| External tickets and destinations | Patch-006 coordinator from accepted inventory | Future external runner and authenticated admin response |
| Actor ID | Future administrator middleware from `AuthSubject.UserID` | Begin/Complete/Abort/Status |
| Native observations and topology proof | Existing collector through the coordinator | Existing provisioner and evidence digest |
| Persisted owner and Key identity | Existing store/provisioner | Restart binding and evidence-only retry |
| `EvalSessionEvidence` domain value | Patch-006 coordinator after successful Initialize | Future protected atomic evidence writer and readiness consumer |
| Availability and idempotent evidence Save | Future protected writer | Begin/Complete preflight and Save-only retry |
| Full readiness | Later integration composing facts, inventory, proof, evidence, executor, and lifecycle | Future paid admission and status decision |

The future outer marker guard, admin routes/authentication, no-store HTTP response policy, production providers, listener/probe implementations, Wire/bootstrap, executor gate, and shutdown drain remain outside patch 006.

## Verification

- `gofmt -d` on both Go files: pass, no output.
- Focused coordinator unit race: pass.
- Coordinator, native topology, and deployment-reader unit race: pass in maintained and fresh replay trees.
- Default service test command: pass as compilation only, no matching tests without `unit` tags.
- Default and `unit` service vet: pass in maintained and fresh replay trees.
- Root `go vet ./...` and `go test -race ./... -count=1`: pass.
- Maintained and fresh-replay evaluation-core vet and race tests: pass.
- Full Sub2API service unit race suite: not run in this review. The pristine upstream suite has a separately reproduced race blocker; focused passes do not clear it.

Independent pristine replay used `dist/sub2api-build/scratch-coordinator-spec-review-HRoTBp`, extracted from pinned commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. Strict sequential `git apply --check --whitespace=error-all` and apply passed for patches 001 through 006. Both new coordinator files were present. The staged 28-file core, all 3,260 backend files, all 931 frontend files, and the complete backend/frontend contents matched the maintained workspace byte-for-byte.

Patch 006 SHA-256: `e9684efbcd470a95e4315ac66ef4f2f7d4ac8e3c979f26ff464cf71234131a39`. Patches 001-005 retained their frozen checksums. `backend/go.mod`, `backend/go.sum`, and `frontend/pnpm-lock.yaml` matched pristine source, maintained workspace, and fresh replay. The core raw-digest hash remained `3eabe2af794582dc8e41dfd637e57fc55940b9e57bc13e9606c1d748f9a8a7ba`; its separate hex-ASCII hash remained `b8477eb71da52e832aefa690b7bdf4fa02c257a792f1669ff3cbb5126ec67af5`.

No Linux execution, release build, Docker lifecycle, live topology, browser, paid request, or charge acceptance was performed. This is a local specification PASS only, not migration completion or live readiness.

## Handoff

Main should register this report in `check.jsonl`. After this specification PASS, a different `trellis-check` should perform the quality review and a fresh independent final replay. No product correction is pending.
