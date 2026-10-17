# Host Topology Coordinator Quality Review

Date: 2026-10-12. Verdict: PASS for the final patch 006 quality gate.

## Context and Scope

The native hook marker was present and pointed to a saved full-context artifact. I read it in bounded sections, then read `check.jsonl`, the PRD, design, current execution plan, frozen coordinator plan, final delivery, specification PASS, backend index, full quality guide, and relevant directory, persistence, error, logging, run-state, and deployment-evidence guides. Context mode: native full hook context, completed from the saved artifact.

Reviewed only:

- `dist/sub2api-build/workspace/backend/internal/service/eval_topology_runtime.go`
- `dist/sub2api-build/workspace/backend/internal/service/eval_topology_runtime_test.go`
- `patches/sub2api/v0.2.15/006-codex-eval-topology-runtime.patch`
- This quality report

No product source needed correction. No source outside this scope, patch 001-005, core file, lockfile, or unrelated `.codex/` and `.devin/` file was changed.

## Findings (Fixed)

None.

## Findings (Not Fixed)

None. The production facts provider, protected evidence writer, bound-listener source, loopback HTTP probe, routes, outer marker guard, Wire/bootstrap, readiness composition, and lifecycle wiring remain later integration owners. Patch 006 requires the four injected ports and fails closed when they are absent.

## API and Contract Review

The public methods match the frozen contract:

```go
func NewEvalTopologyRuntime(EvalTopologyRuntimeDeps) (*EvalTopologyRuntime, error)
func (*EvalTopologyRuntime) Begin(context.Context, int64, EvalProbeBeginRequest) (EvalProbeBeginResult, error)
func (*EvalTopologyRuntime) Observe(*gin.Context) (bool, error)
func (*EvalTopologyRuntime) Complete(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Abort(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Status(context.Context, int64, string) (EvalProbeSessionStatus, error)
```

Exported DTOs are `EvalProbeCaseDTO`, `EvalProbeBeginRequest`, `EvalProbeExternalTicket`, `EvalProbeBeginResult`, `EvalProbeSessionStatus`, and `EvalSessionEvidence`. All case tuples encode `kind`, `ingress_id`, `header`, and `source`. Evidence contains independent current facts and exact inventory/proof digests, but no proof, ticket, raw token, or Key lookup.

The required dependency ports are `EvalTopologyProvisioner` (`Initialize`, `Load`), `EvalDeploymentCurrentFactsProvider` (`CurrentEvalDeploymentFacts`), `EvalBoundListenerSource` (`ActualEvalListener`), `EvalTopologyLocalProbe` (`ProbeLoopback`), and `EvalSessionEvidenceWriter` (`CheckAvailable`, `Save`). Constructor checks reject nil and typed-nil dependencies. The listener must be TCP with a server-owned `http` or `https` scheme.

The coordinator binds sessions and tickets to the authenticated actor and generation. It retains ticket digests only, validates the exact method/path and case-insensitive header key, consumes tickets and leases under the mutex, and runs native observations outside it. Abort drains leased work; late results cannot publish into invalid generations. Status is read-only and sanitized, including synthetic expiry without changing state.

Completion reloads trusted facts and the exact inventory before writer availability, token loading, sealing, and initialization. Every `Initialize` error enters recovery-required and preserves the owner binding. Only successful initialization followed by evidence-save failure enters completion-pending. Retry reloads the exact intent and identity, verifies token-derived lookup and proof bindings, and calls only `Save`.

## Producer and Consumer Mapping

| Value | Producer and remaining owner | Consumer |
| --- | --- | --- |
| Protected deployment inventory | Privileged operator/installer after independent topology inspection | Patch-005 reader, then coordinator Begin/Complete/retry |
| Current deployment facts | Trusted runtime/installer provenance provider | Reader comparisons and coordinator reloads |
| Actual bound TCP listener and scheme | Server listener capture and TLS integration | Begin and loopback probe |
| Local probe requests | Real loopback HTTP transport through the captured listener and same Observe hook | Coordinator collector |
| External tickets and destinations | Patch-006 coordinator from accepted inventory | Future external runner and authenticated admin response |
| Actor ID | Admin middleware from `AuthSubject.UserID` | Begin/Complete/Abort/Status |
| Native observations and topology proof | Existing collector through coordinator | Existing provisioner and evidence digest |
| Persisted owner and exact Key identity | Existing store and provisioner | Restart binding and evidence-only retry |
| `EvalSessionEvidence` | Patch-006 coordinator after successful Initialize | Future protected atomic evidence writer and readiness consumer |
| Evidence availability and idempotent Save | Future protected writer | Begin/Complete preflight and Save-only retry |
| Full readiness | Later integration composing facts, inventory, proof, evidence, executor, and lifecycle | Future paid admission and status decision |

## Verification

- Lint: no standalone linter is configured; `gofmt -d` on both coordinator files passed with no output.
- TypeCheck: service packages compiled in maintained and replay trees; `go vet ./internal/service` and `go vet -tags=unit ./internal/service` passed in both.
- Tests: focused coordinator/native-topology/deployment-reader unit race passed in maintained and replay trees.
- Default service race test command passed as compilation only; no matching tests ran without `unit` tags, in maintained and replay trees.
- Root `go vet ./...` and `go test -race ./... -count=1` passed.
- Maintained and replay evaluation-core vet and race tests passed.
- Full Sub2API service unit race suite: not run. The pristine upstream baseline has a separately reproduced race blocker, and unrelated fixture repair was declined. Focused passes do not clear that blocker.

## Pristine Replay

Created a new isolated Git replay at `dist/sub2api-build/scratch-coordinator-quality-PMsqqJ` from pinned commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. Applied patches 001 through 006 sequentially; each `git apply --check --whitespace=error-all` and application passed. Both coordinator files were present. The unchanged 28-file evaluation core was staged into the replay.

Full tree comparisons matched: all 3,260 backend files, all 931 frontend files, and all 28 core files. Patch SHA-256 values are:

| Patch | SHA-256 |
| --- | --- |
| `001-codex-eval-auth.patch` | `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b` |
| `002-codex-eval-scheduling.patch` | `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5` |
| `003-codex-eval-cancellation.patch` | `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4` |
| `004-codex-eval-target-adapter.patch` | `850af726f80bb2f379e0cf9b3f395f28e9d8e89356106f0643c39ee5558f4088` |
| `005-codex-eval-deployment-evidence.patch` | `3402198f3f75e73ebf46e748ef4a289001a13d0632692c3bb493a630055650cb` |
| `006-codex-eval-topology-runtime.patch` | `e9684efbcd470a95e4315ac66ef4f2f7d4ac8e3c979f26ff464cf71234131a39` |

The three lockfiles matched pristine source, maintained workspace, and replay:

| File | SHA-256 |
| --- | --- |
| `backend/go.mod` | `809e6806dc248bd44cea4198ddc21fd6ed04f80003300c3668367b91bae319a8` |
| `backend/go.sum` | `b17668c6300c04027bd397650e60e28e514dc4abb9ab113e3cea634fdd708109` |
| `frontend/pnpm-lock.yaml` | `1d8e54cc8e1f1bd3e7593805f607d9314ff81dab98ab5a9a6396475e148b26ce` |

No Linux execution, release build, frontend build, Wire generation, Docker lifecycle, live topology, browser, paid request, or charge acceptance was performed. This is a local quality PASS for patch 006, not migration completion or live readiness.

## Handoff

Main can register this report and update the source-backed coordinator spec after this final quality gate. Runtime integration remains blocked on the four production ports and the later HTTP, readiness, and lifecycle owners listed above.
