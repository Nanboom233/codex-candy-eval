# Research: Consolidated Sub2API Check Plan

- Query: Prepare one actionable full-scope `trellis-check` handoff for the Sub2API runtime patch task, mapping requirements and acceptance to existing evidence, missing integration, or live checks.
- Scope: internal
- Date: 2026-10-12

## Findings

### Current boundary

The frozen target is Sub2API v0.2.15 at commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`; toolchain pins are Go 1.27.2 and pnpm 9.15.9. Patches 001-006 and the unchanged 28-file evaluation core have local stage reviews and frozen identities in `implement.md`. These are component gates, not migration acceptance.

The maintained tree has evaluation primitives and a topology coordinator, but no production evaluation handler/routes, Wire/bootstrap, current-facts provider, protected session-evidence writer, actual listener/probe adapters, evaluation frontend, or Sub2API installer. `router.go` registers existing route sets at lines 116-135; the Wire injection sets in `cmd/server/wire.go:34-62` do not include evaluation integration. The coordinator spec explicitly excludes these owners (`host-coordinator.md:5-11`).

### Files found

- `.trellis/workflow.md`, task `prd.md`, and `implement.md` - workflow, all criteria, current handoff, frozen identities, and known gaps.
- Task research `host-coordinator-delivery.md`, `deployment-evidence-reader-delivery.md`, and `upstream-race-baseline.md` - reviewed evidence, replay procedure, and known race blocker.
- `.trellis/spec/backend/{quality-guidelines,evaluation-core,run-state,host-target-adapter,deployment-evidence,host-coordinator}.md` - applicable contracts.
- `internal/eval/` - evaluation core, run manager, durable state, and tests copied into the pinned backend.
- `backend/internal/service/eval_key_runtime.go`, `eval_topology.go`, `eval_deployment_evidence.go`, and `eval_topology_runtime.go` in the maintained workspace - token/key, topology, inventory, and coordinator primitives.
- Maintained backend `middleware/api_key_auth_eval.go`, `service/openai_eval_directed.go`, and `handler/openai_gateway_handler.go` - directed auth/scheduling and normal Responses gateway owners.
- Maintained backend `internal/server/router.go` and `cmd/server/wire.go` - production route and dependency registration owners.
- `scripts/test-sub2api.sh`, `scripts/build-sub2api.sh`, and `scripts/sub2api_build.py` - local checks, pinned build, replay inputs, and amd64/arm64 artifacts.
- Maintained `frontend/package.json` - frontend typecheck/test/build commands.
- `README.md`, root `install.sh`/`install.ps1`, and `.github/workflows/release.yml` - still CPA-specific public install/release surfaces.

### Code patterns and evidence

- `internal/eval/` implements all four types and bounded local Responses transport (`evaluation-core.md:15-86`); `runs.go`/`state.go` cover durable admission, recovery, 20-record retention, cancellation, and shutdown (`run-state.md:36-98`).
- `eval_key_runtime.go:67-91,123-193` stores a random private token in a `0600` file and derives a separate lookup. `eval_identity.go:28-91` holds runtime identity/token in memory. `api_key_auth_eval.go:19-63` rejects lookup replay and restricts directed headers to private-token Responses requests with valid proof.
- `openai_eval_directed.go:29-153` reloads the exact group/account; adjacent service/handler tests cover no fallback and directed cancellation. They do not exercise an integrated evaluation producer/API.
- `eval_topology_runtime.go:80-113,181-185` requires facts, writer, listener, probe, provisioner, store, and config; several production ports are absent. `router.go:92-135` and `cmd/server/wire.go:34-62` have no evaluation route/coordinator registration.
- Frontend search found no evaluation feature surface; the existing evaluation UI is the CPA page in `internal/plugin/web/`. The test/build scripts provide root, backend, frontend, and static-build commands but no installer (`test-sub2api.sh:39-82`, `build-sub2api.sh:104-219`).

### One consolidated check sequence

The checker should complete this sequence once and label each result `PASS`, `FAIL`, `BLOCKED`, `NOT RUN`, or `LIVE ONLY`.

1. Run `scripts/build-sub2api.sh check`. In a fresh isolated replay from the pinned archive, apply 001-006 in order with `git apply --check --whitespace=error-all` and strict apply; assert patch 005/006 files exist; stage the unchanged 28-file core; compare backend/frontend/core bytes and three lockfiles with the maintained tree. Verify frozen patch identities and compare core hashes using the same algorithm (`3eabe2...` raw digest or `b8477e...` hex text) recorded in `implement.md`.
2. Run `scripts/test-sub2api.sh check`. In the maintained backend, run `go vet ./internal/codexeval` and `go test -race ./internal/codexeval -count=1` for copied-core checks.
3. Run focused unit-race tests for `EvalTopology`, `EvalDeploymentEvidence`, `EvalKey`, `EvalTargetAdapter`, and `Directed` service tests; run all `internal/server/middleware`, `internal/handler`, and `internal/handler/dto` tests with `-tags=unit -race`. Mark default-tag runs compilation-only where tests are unit-tagged.
4. Run `scripts/test-sub2api.sh integration` once on the patch-verified workspace. It runs backend vet, full unit race, integration-tag tests, and frontend checks. The full service race gate is expected to hit the documented pristine fixture failures. Do not rerun that pristine baseline, repair unrelated fixtures, suppress races, or report focused PASS as full-service PASS. If it stops there, run `(cd dist/sub2api-build/workspace/backend && go test -tags=integration ./...)` and step 5 separately.
5. Run `(cd dist/sub2api-build/workspace/frontend && pnpm install --frozen-lockfile && pnpm typecheck && pnpm test:run && pnpm run build)`. Once workspace patch/core state is verified, run `scripts/build-sub2api.sh build --arch all`; inspect both static artifacts and provenance/checksums. This does not prove installed-container compatibility.
6. Audit every R/A/K ID below against source, tests, and artifacts. Treat absent production owners as gaps. Do not run CPA deployment/install scripts, paid requests, or an unapproved live deployment as Sub2API checks.

### Requirement and acceptance mapping

| Requirements | Acceptance | Existing evidence | Gap or disposition |
| --- | --- | --- |
| R1, R2 | A1 | Four core implementations, prompts/data, transport, and fixtures. | No Sub2API evaluation handler/UI; controlled gateway completion is incomplete. |
| R3 | A2, A3 | Directed gateway primitives preserve normal admission. | Usage, charges, and shared concurrency are not assembled; real charges need a live target. |
| R4, R5 | A4, A5, A6 | Token/auth primitives, target adapter, normal gateway boundary. | Admin endpoints/bootstrap and concurrent-group flow are absent; authorization/quota/billing cases need integration tests. |
| R6 | A7 | Exact-target/no-fallback fixtures. | Run manager must record selected-target failure after integration. |
| R7 | A8 | Lookup replay rejection, owned-Key redaction tests, admin auth primitive. | Evaluation APIs/menu and assembled browser/token-output checks are absent. |
| R8 | A9 | State tests cover retention, interruption, persistence, cancellation. | Host/API/UI wiring and refresh/restart acceptance are incomplete. |
| R9, R10 | A10, A11, A12 | Pinned build tooling only. | Installer/status/rollback/uninstall, restore, history retention, and recreation detection are absent. |
| R11 | A13 | Manifest pins release/commit/digests and amd64/arm64; builder validates source/toolchain/locks. | Build outputs remain; target refusal needs installer checks. |
| R12 | A14 | Sub2API build workflow exists. | README, root installers, and release workflow still describe/build CPA; docs/release migration is incomplete. |
| R13 | A15 | CPA UI and VM harness exist. | Sub2API evaluation UI, theme, and browser behavior are absent. |
| R14 | A16, A17 | Topology proof and Key conflict primitives have focused tests. | Facts/listener/probe/routes/readiness/revocation/reapply integration is absent. |

### Key boundary mapping

| IDs | Current evidence | Outstanding check |
| --- | --- | --- |
| K1 | Protected token file and server-side mapping; separate in-memory identity. | Verify no DB/cache/browser/URL/log/process-argument leak in assembled runtime. |
| K2 | Lookup redaction and external replay rejection have unit tests. | Verify all browser/admin responses and original-runtime boundary. |
| K3 | Directed routing uses a Key copy and normal gateway admission. | Integration/live proof for permissions, subscriptions, balance, quotas, billing, concurrency. |
| K4 | Loopback/proxy proof primitives exist. | Real facts/probes are absent; unknown topology must keep activation unavailable. |
| K5 | Coordinator state exists; installer does not. | Stop work and revoke only owned Key; preserve history and report failures. |
| K6 | Accepted residual in PRD. | Do not claim stronger original-runtime secrecy or auto-revocation. |
| K7 | Provisioning checks exact owner/status/binding/lookup conflicts. | Installer tests must show reapply never silently rebinds or reenables. |

### Related specs and external references

- Related specs: `.trellis/spec/backend/index.md`, `quality-guidelines.md`, `evaluation-core.md`, `run-state.md`, `host-target-adapter.md`, `deployment-evidence.md`, and `host-coordinator.md`.
- External source pin recorded by the PRD: Wei-Shaw/Sub2API v0.2.15, commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`, annotated tag object `86a80c13dcba86f52f9ca815b5a471cecc236227`.
- Toolchain pins recorded by `implement.md`: Go 1.27.2 and pnpm 9.15.9. No external documentation was fetched for this audit.

## Caveats / Not Found

- This is a research handoff. No tests, build, Docker command, deployment, browser request, or live request was run here.
- No approved live Sub2API target is configured. Authenticated browser behavior, actual proxy topology, one real candy result, usage records, and charges are `LIVE ONLY` and remain unverified.
- Docker Desktop is reachable through `desktop-linux`, but lifecycle coverage needs the missing installer/integration code. No production evaluation UI/routes/Wire, evidence writer, current-facts provider, listener/probe, or Sub2API installer was found.
- The pristine upstream service race failure is documented and unrelated fixture repair was declined. Do not rerun the pristine baseline or claim focused checks as a full-service PASS. Old CPA deployment scripts and private `.env` values are outside scope.
