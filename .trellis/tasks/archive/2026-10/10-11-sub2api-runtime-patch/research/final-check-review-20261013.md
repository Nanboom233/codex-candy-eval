# Final Independent Check Review

Date: 2026-10-13
Scope: full-scope review of the approved patch-native revision, including a fresh strict replay, focused native checks, frontend checks, release guards, prior lifecycle evidence, and acceptance mapping.

## Findings (fixed)

- File: `.github/workflows/release.yml:16-23`.
  Issue: the release workflow could publish the pinned `v0.2.15` bundle from a different `v*` tag.
  Fix: compare `GITHUB_REF_NAME` with the manifest `release` field before delivery. Workflow dispatch runs remain available without a tag.

## Findings (not fixed)

- Full upstream unit-race gate remains failed in the authoritative final revision. `integration-cleanup.log` reports 16 race warnings and 9 top-level failures, with zero compiler failures or timeouts. Failures are the documented Grok quota/probe, CyberPolicy scope, tool-schema rewrite, usage cleanup, and concurrent AuthService alias fixtures. Focused passes do not replace this gate.
- Live acceptance remains open: real paid billing, production concurrency, deployed browser and payload inspection, published online bootstrap, native Windows execution, and native Linux/amd64 Docker attestation.
- Native amd64 Docker lifecycle remains blocked by Rosetta executable identity. The arm64 lifecycle passed, but amd64 loaded-process attestation requires a native Linux/amd64 runtime.

The approved patch-native revision retires public topology probes and ingress-inventory readiness. No external-ingress-inventory gap is claimed here. The remaining live gaps concern deployed behavior and platform/runtime availability.

## Reviewed implementation boundaries

- Production contains six patches: 001, 002, 003, 004, 007 and 008 under `patches/sub2api/v0.2.15/production/`. Patches 005 and 006 are retired. The reusable production evaluation core is `internal/eval/`, staged at `backend/internal/codexeval/`.
- Runtime composition is owned by `backend/internal/service/eval_host_runtime.go`, `eval_runtime_ports.go`, `eval_identity.go`, `eval_key_runtime.go`, `eval_http_transport.go`, `eval_maintenance.go`, `eval_target_adapter.go`, and `openai_eval_directed.go`.
- The reviewed runtime interfaces are `EvalHostRuntime.BindRouter(http.Handler)`, `EvalHostRuntime.Activate(context.Context, int64)`, `EvalHostRuntime.Deactivate(context.Context, int64)`, `EvalHostRuntime.ServeMaintenance(context.Context)`, and private `evalServeHTTPTransport.RoundTrip(*http.Request)`.
- `backend/internal/server/http.go:ProvideEvalHTTPServer` binds the runtime to the real native router before the outer guard. `backend/internal/server/middleware/api_key_auth.go`, `api_key_auth_eval.go`, and `api_key_auth_google.go` own capability admission and external replay rejection.
- `backend/internal/handler/admin_codex_eval.go` and `backend/internal/server/routes/admin.go` expose administrator-authenticated status, activation, targets, state, start, cancel, records, history, flush, and deactivation APIs under `/api/v1/admin/codex-eval`. The static shell is `/codex-eval/ui`.
- `backend/cmd/server/main.go` and `maintenance_client.go` implement the same-binary maintenance mode. The live runtime owns `/app/data/codex-eval/control.sock`; the client sends only bounded deployment and binary identity facts.
- Internal evaluation dispatch is restricted to canonical `POST /v1/responses` through the genuine router. No private seed, lookup, or directed target header is sent on the model request.

## Fresh replay

Reviewer workspace: `dist/sub2api-build/reviewer-final-20261013/`.

- Source: `dist/sub2api-v0.2.15`, verified clean at commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`, annotated tag object `86a80c13dcba86f52f9ca815b5a471cecc236227`.
- Production prepare/build: passed with strict patch, lockfile, Wire, frontend, and static build guards.
- Production patch-set SHA-256: `12a1865ff561c4c9f1c84df1c407fefb405db4093ec118d589bc651c065d9437`.
- Test workspace creation: passed with add-only overlay and unchanged production bytes for 2,633 production paths.
- Test overlay patch-set SHA-256: `2fe1d542fca6ca070b2e721bc47c29c07b45485d063e6cc518f86af210335f38`.
- Test augmented-tree SHA-256: `36dfe8542bb09cd496f9999d645243cd9e13c11038829ab4c7307224fdae068c`.
- Lockfile baselines: production and test both passed.
- Linux amd64 binary SHA-256: `fd3109997b4d971e524a13ddd87b4b4c8f1cd8aad213036b0dae14ba33be213f`.
- Linux arm64 binary SHA-256: `dd3760228f18cf63ab487ef6d3d433344d5c1b3d8ec89255f38a7deeb1bd129a`.

The fresh replay is evidence for the current patch and source inputs. It is not a pristine A/B attribution run. The known full-race failures must not be described as proven baseline-equivalent solely from failure counts or source overlap.

## Previously completed full-scope checks

- `scripts/test-sub2api.sh check` passed 68 build-helper checks, 22 installer/bootstrap checks, shell syntax, root `go vet ./...`, root `go test -race ./...`, and the UI harness. The complete failed integration gate is retained separately and remains failed.
- Strict production/test payload inspection passed: 15-file release allowlist, archive contents, checksums, production-only provenance, lock identities, and production/test byte parity. The authoritative cleanup archive is `dist/sub2api-build/installation-revision-20261013/sub2api-codex-eval-v0.2.15-linux-amd64-arm64.tar.gz` with SHA-256 `fde738f78d88f522c94f82a99b01bfe6ea84209bb8abbe6d65c86f9a3671db0a`.
- The arm64 disposable lifecycle passed 14 checkpoints in `lifecycle-arm64-cleanup.log`, with evidence in `codexeval-fixture-63cba0a6f0/evidence.json`: explicit target selection, uninitialized install, repeat activation, maintenance peer denial, expiry drift, recreation, original paid rejection, native read-only residual, legacy migration, Redis invalidation, interrupted retry, exact revocation, restoration, uninstall, and candidate recovery.
- The retained PostgreSQL revocation retry in `revocation-retry.log` verified exact owned-row deletion after live-runtime drain, preserved row counters and expiry, preserved the revoked marker, and restored the original process. SQL was used for inspection only, not cleanup.
- Browser evidence in `browser-validated.json` and `browser-cleanup-upgrade.log` passed four tabs, math rendering, light/dark themes, preference/history retention, 11 authenticated same-origin calls, an `allow-scripts`-only pelican sandbox, and zero captured CSP errors.
- Native amd64 lifecycle evidence is incomplete because Docker Desktop exposed Rosetta at `/proc/1/exe`; the strict executable identity check rejected it and the candidate restored the official original. This is not a successful amd64 rollback attestation.
- `actionlint` was unavailable in the local environment. Ruby YAML parsing, manifest tag checks, and `git diff --check` were used instead.

## Verification

- Source identity: pass (`verify-source` returned `source ok`).
- Unit vet: pass for service, middleware, handler, DTO, evalroute, and codexeval.
- Focused native unit race: pass for service, middleware, handler, DTO, evalroute, and codexeval using `TestEval|TestProbe|TestProvideHTTPServer|TestHTTP|TestGlobal|TestAPIKey`.
- Evaluation core default race: pass.
- Evaluation core unit race: pass.
- Frontend lint: pass.
- Frontend typecheck: pass.
- Frontend tests: pass, 363 files and 2,797 tests.
- Release workflow YAML parse: pass.
- `git diff --check`: pass.

## Acceptance mapping

The implementation report's detailed R/A/K matrix was reviewed against the current source paths and the fresh replay. The following records each ID's result and exception boundary.

| IDs | Result and exception boundary |
| --- | --- |
| R1-R2 / A1 | Local composed native gateway and four evaluation flows pass with preserved prompts, resources, scoring, model, reasoning, and batch controls. Live inference is not claimed. |
| R3 / A2-A3 | Local usage, charge attribution, and shared user/account slot fixtures pass. Real paid billing and production concurrency remain live-only. |
| R4-R6 / A4-A7 | Exact owned identity, native permissions, target membership, no fallback, fresh recheck, and cancellation fixtures pass. |
| R7 / A8 | Administrator-only APIs, DTO redaction, private-token exclusion, and main/Google/WebSocket replay rejection pass locally. |
| R8 / A9 | Durable history, retention, restart, interruption, and browser refresh checks pass. |
| R9-R10 / A10-A12 | Arm64 official-image preservation, install, idempotence, recreation, rollback, uninstall, exact cleanup, and recovery pass. Native amd64 lifecycle attestation remains blocked by Rosetta. |
| R11 / A13 | Pinned source/image facts and both static architectures pass. Native Windows execution and native amd64 lifecycle remain open. |
| R12 / A14 | README, operations docs, installers, release workflow, payload, and licenses target Sub2API. Deployed payload inspection remains open. |
| R13 / A15 | Local browser tabs, controls, details/history, theme reload, renderers, and sandbox pass. Deployed browser behavior remains open. |
| R14 / A16-A17 | Capability-only admission, forged-header rejection, original-runtime paid expiry rejection, exact revocation, conflict handling, and incomplete cleanup behavior pass. Approved public topology/probe readiness is retired in this revision. |
| R15 / A18 | Independent production/test inputs, path/hash parity, test exclusion, provenance, and release allowlists pass. |
| R16 / A19-A20 | Mandatory container choice, explicit `--container`, cancellation/invalid-choice guards, install-uninitialized state, native administrator activation, repeat activation, and crash recovery pass locally. Published online bootstrap and native Windows remain open. |
| R4/R14 / A21 | Atomic first-insert past expiry and original-runtime paid rejection pass. Native read-only access retains the documented unbound-group residual. |
| R4/R10 / A22 | Stable row ID/counters, legacy migration, Redis invalidation, drift refusal, durable retry, and no replacement identity pass. |
| R10/R12/R16/R17 / A23 | Numbered process output, failure/recovery states, exit codes, secret/dump/traceback exclusions, and lifecycle output pass. Native Windows execution remains open. |
| K1 | Protected seed/intent state, private permissions, no-symlink checks, and credential-free internal requests pass. |
| K2 | Exact owned lookup/private seed rejection and changed-row redaction pass on external main/Google/WebSocket paths. |
| K3 | Only the owned row's fixed expiry is exempted; status, owner, binding, subscriptions, quotas, limits, pricing, slots, and charges remain native. |
| K4 | Protected build facts, owned identity, genuine-router binding, and dispatcher readiness pass without public topology inputs. |
| K5 | Live-runtime drain, exact revocation, durable retry, retained history, and strict restore pass. Unverifiable legacy menu entries remain untouched and incomplete. |
| K6 | Original normal-mode HTTP/Google/new-WS paid rejection passes. Read-only usage/billing visibility and preexisting WebSocket session limits remain accepted. |
| K7 | Missing, disabled, rebound, expiry-drifted, and owner-conflicted identities fail closed. Verified legacy migration preserves the row and counters without replacement. |

No deployment, CPA `.env` access, live quota use, commit, push, or publication occurred during this review.

Evidence logs are retained beside the reviewer workspace: `prepare-production.log`, `build-production.log`, `create-test-workspace.log`, `focused-unit-vet.log`, `focused-unit-race.log`, `core-race.log`, `core-unit-race.log`, and `frontend-checks.log`.

