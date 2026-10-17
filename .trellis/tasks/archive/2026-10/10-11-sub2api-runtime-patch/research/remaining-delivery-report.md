# Remaining delivery implementation report

Active task: `.trellis/tasks/10-11-sub2api-runtime-patch`.
Status: concrete local delivery, including joined native administrator/gateway four-flow verification, is ready for the main session's single full-scope checker. The task is not complete: the retained full upstream unit-race gate fails and live acceptance has no approved target.

This report supersedes all earlier draft checkpoints. Product sources are final in the production patches/core; tests and verification helpers remain separate development inputs.

## Post-review identity notice

The full-scope checker corrected mapped-unspecified inventory validation in production005 and its separate test overlay.
Use [remaining-delivery-check-review.md](remaining-delivery-check-review.md) for final production/dev hashes, binary checksums and refreshed release payload.
Implementation-stage hashes and27/105 source classification below remain pre-review evidence, not the current gate result.
Final reviewer full-race result is23 warnings/112 top-level failures with zero compiler errors.

## Files and behavior delivered

- `patches/sub2api/v0.2.15/production/001–008`: production-only authentication, exact scheduling/cancellation, metadata adapter, inventory/coordinator and concrete runtime/API/UI wiring. Unchanged historical flat 001–006 remain as historical evidence. The production coordinator has no test-only callback. Test additions are add-only patches under `test-overlay/` plus maintained `tests/overlays/codexeval/` and `tests/host-overlay/`.
- `internal/eval/`: explicit production core/resource allowlist, HTTP/core/run-state preservation, cached embedded nonce UI. Four tabs, account rows, group selection, batch/details/history, Markdown/math and pelican sandbox remain. JWT is read from current native login storage for each request and goes only in headers. Native enable state is shown; changes direct operators to native account management. Generated renderer styles move into nonce stylesheets; model-page style attributes and style/script tags are prepared before DOMParser to preserve CSP. Embedded font data is allowed without `unsafe-inline`.
- 007 owners: `eval_runtime_ports.go`, `eval_host_runtime.go`, `admin_codex_eval.go`, native providers/routes/router/server/Application/main/Wire graph and `cmd/codex-eval-probe`. Protected facts bind independently inspected image/platform/source/ingress provenance and the actual running executable. Session evidence uses strict JSON/private owner/mode/no-symlink descriptor access, exact binding, atomic sync and idempotence. Begin/startup/deactivate preflight safe ancestry before old KeyStore access. Existing `row_pending` lookup rows recover identity guards only, never activate/repair the intent; missing rows require operator repair. Owned-Key conflicts become a fixed status code and block installation verification. DTO redaction uses the exact owned row ID even if external identity fields change. Shutdown drains before HTTP/dependency teardown. Native revocation precedes removing exact owned token/intent/session inputs; history survives.
- `sub2api_install.py`, root `install.sh`/`install.ps1`: selected-container inspect/install/status/rollback/uninstall, immutable original backup/private journal/facts, staged replacement, restart and health/admin/gateway checks, scoped menu ownership and Key cleanup. Docker's containerd image store reports indexes as IDs; the installer uses the container's selected platform descriptor and independently hashes exported OCI manifest/config blobs. Pristine-original comparison uses the selected platform explicitly. Failed and stopped candidates restore the verified original before final journaling; restored readiness is checked. No native compliance acknowledgement is auto-accepted in product operations.
- `scripts/sub2api_build.py`, build/Wire/pnpm wrappers: strict patch/path/whitespace/lock guards, pinned tooling, isolated production/test workspaces, explicit copied resources, source parity and distinct identities. Compiler caches are omitted from source hashes; runtime frontend assets remain included. Production frontend invokes the pinned compiler/build directly because upstream `pnpm build` starts an i18n test excluded from production inputs. `tests/go.mod` excludes overlay-only packages from root package discovery.
- `scripts/test-sub2api.sh`: root/helper/installer checks and independent vet/full-race/integration/frontend gates. Every independent gate is attempted and any failure still fails the overall integration command. No race suppression.
- `scripts/package-sub2api.py`, release workflow, README, operations and licenses: explicit release payload, compact installer compatibility metadata without test/build/core inputs, production-only provenance/checksums, complete Sub2API LGPL-3.0/ModelTrace MIT/renderer notices, Sub2API-only release operations and disposable Docker CI check.
- Development-only verification: maintained composed native admin/gateway four-flow and cancellation/concurrency fixture, installer transaction units, official-image Docker lifecycle driver, synthetic browser seeder, Playwright UI driver, expanded helper isolation tests and new runtime/operator/DTO/frontend overlays. None ship in production artifacts.

## Authoritative local checks

| Check | Command/result | Evidence |
| --- | --- | --- |
| Root and helpers | `scripts/test-sub2api.sh check`: 68 helper tests, five installer units, shell syntax, root vet/race and historical UI harness passed | `root-review-final.log`; `installer-unit-final.log` |
| Strict final replay/parity | `sub2api_build.py create-test-workspace ... test-composed-final`: fresh pinned source, eight production patches, copied generated runtime bytes, add-only overlay, unchanged 2,630 production path/hash entries before/after overlay | `test-composed-final-evidence.json` |
| Final focused race | `go test -tags=unit -race ./internal/service ./internal/server ./cmd/codex-eval-probe ./internal/handler/dto ./internal/codexeval -run 'TestEval|TestProbe|TestProvideHTTPServer|TestHTTP|TestGlobal|TestAPIKey' -count=1`: all five packages passed, including pending guards, safe preflight, real listener, four flows/history/drift, exact conflict/deactivation and changed-row redaction | `test-review-final/focused.log` |
| Joined native admin/gateway and full affected race | `go test -tags=unit -race ./internal/handler ./internal/server/middleware ./internal/pkg/evalroute ./internal/codexeval`: all four packages passed, including composed four-flow results, billing, both shared concurrency limits, partial-usage cancellation, history and drift rejection | `test-composed-final/gateway-core-race.log`; `composed-final-results.json` |
| Final affected runtime race | service/server/operator/DTO focused race passed in the fresh final workspace | `test-composed-final/focused-runtime-race.log` |
| Final affected unit vet | `go vet -tags=unit ./internal/handler ./internal/server/middleware ./internal/pkg/evalroute ./internal/codexeval`: passed | `test-composed-final/unit-vet.log` |
| Final backend vet | `go vet ./...`: passed | `test-review-final/vet.log` |
| Backend integration | final `go test -tags=integration ./...` passed (exit 0 confirmed) | `test-final-check/backend-integration.log`; `test-review-final/integration.log` |
| Frontend | pinned frozen pnpm, lint/typecheck, 362 files/2,795 tests and build passed | `frontend-final.log` |
| Static artifacts | `scripts/build-sub2api.sh build --arch all`: final amd64/arm64 server and operator clients passed; all are stripped static ELF, CGO=0 and no parent VCS metadata | `final-build.log`; `artifact-final-evidence.json` |
| Payload inspection | explicit allowlist and binary checksums passed; unexpected fixture and development identity negatively rejected; compact compatibility metadata has only installer-required fields | `release-platform-final/`; `artifact-final-evidence.json` |
| Real Docker arm64 | 12 checks passed: image inspect/preservation, install, idempotence, native menu conflict preserved, native owned-Key redaction/external disable detected without reenabling, native scoped revoke/unrelated Key preservation, rollback, uninstall, recreation/reapply, running and stopped candidate failure restoration | `docker-platform-final.log`; preceding `docker-review-final.log` |
| Real Docker amd64 | same 12 checks passed under Docker Desktop emulation with exact selected amd64 manifest and original/patched binaries | `docker-amd64-final.log` |
| Public entry points | packaged POSIX `install.sh status` passed; packaged `install.ps1 inspect` passed through local PowerShell. Windows-host execution itself was not available/tested. | `packaged-entry-status.log`; `powershell-entry-final.log` |
| Actual browser | authenticated assembled native fixture passed four tabs, math with no inline styles, correct light/dark tokens, preferences and history reload, 11 header-auth API calls, pelican scripts/styles with storage blocked and exact `allow-scripts` sandbox, zero new console/CSP errors | `playwright-browser-review-final.log`; driver `tests/ui/sub2api-browser.js` |

Paths in the evidence column are relative to `dist/sub2api-build/` unless stated otherwise. Browser histories are explicit synthetic local data; they do not claim model completion or live accounting. The maintained composed fixture joins native administrator JWT authentication and Start/State/Record APIs, the real topology coordinator, protected facts/session persistence, an actual local listener, normal API Key middleware, directed native Responses scheduling, billing/slots and a controlled local HTTP upstream. Its four results complete. It also checks normal user/group/account/token/cost records, balance deductions, independent shared user/account limits, cancellation with observed usage, history after manager reconstruction, and drift rejection without new upstream calls. All repository/transport adapters and the readiness-checked public HTTPExecutor wrapper are development-only. The production lazy executor is separately covered by the real host fixture. The joined fixture uses in-memory native repository ports. Image/ingress declarations and external peer observations are synthetic local test inputs. The executable hash and loopback HTTP requests are real. Docker database inference and live ingress topology are not established by this fixture. Browser histories remain synthetic and live acceptance is not inferred.

## Full upstream race gate

Final command: `go test -tags=unit -race ./...` in `dist/sub2api-build/test-review-final/backend`.

Final result: **FAILED — 27 race warnings, 105 top-level failing tests, zero compiler errors.** Failing packages are native `internal/handler/admin` and `internal/service`. New maintained overlay tests have no failures. The HTTP constructor compatibility regression from the first run was fixed by preserving native `ProvideHTTPServer(cfg, router)` and using a separate DI guard provider; native server tests now pass.

Evidence: `test-review-final/full-unit-race.log`, `full-race-final-analysis.json`. The historical pre-integration audit had 18 warnings/58 failed tests. This run's counts differ; no pristine baseline rerun or unrelated fixture repair was authorized/performed. Do not describe these different counts as identical historical failures, suppress the race gate, or declare the overall task passing. Source-only classification found 105 failures with immediate race-only output and no independent assertion mismatches. Their test files match pinned upstream bytes. Forty-one failed names recur in the saved historical failure set. The 27 direct races fall into nine native groups; eight have historical frame evidence. Two scheduler frame files changed in this task, but their three participating functions remain byte-identical to pinned source and the stack uses ordinary Grok scheduling. See [full-race-source-classification.md](full-race-source-classification.md) and `full-race-source-classification.json` for all frames/test owners and comparison evidence. This ownership classification does not prove identical historical execution or remove the failed gate.

## Artifact identities

Production patch set: `647e57ff63306274de0c726db9626f79c30023815678a2d26c64ea98922db06c`.
Production core/resource tree: `71cebfd65c9167eacaff687d677e51a032ed12503c876185a0d6236ff2f7bd58`.
Complete production staged tree: `c2afbc8c9c4601620493c8da4dea2f82e2bd7fc1fdfb068208d528914856e91d`.
Development overlay patch set: `50ca105b2508b39b11330aa8885a05d8fc6cae68a5825b8be31afc859438a545`.
Development augmented tree: `dceb21ec39c8445f09ed7d2138f63f93a7f3941be85bd4296e9d7d6dd1669c39`.

- `sub2api-v0.2.15-linux-amd64`: `b8688d794b846cc240c4e6faf4a5348dea4ebf291b75a264163b4c490a4a1c80`
- `codex-eval-probe-v0.2.15-linux-amd64`: `58eaf2ebfe3f3b18ea87f490060202702baf2bfe72de7f3d78eeb7338ac5fde1`
- `sub2api-v0.2.15-linux-arm64`: `8124959e5de78da30056eab3ec09e43394b0f8942100c5c5886c330db7f30548`
- `codex-eval-probe-v0.2.15-linux-arm64`: `3733f8b2720d2e573e6a6f274df70e6445b0ec5f9ab9d798bcfe76ba1a92b88f`

Final production output: `dist/sub2api-build/out/`.
Reviewed installer payload: `dist/sub2api-build/release-platform-final/`.
Final disposable overlay: `dist/sub2api-build/test-composed-final/`. The saved full-race and complete backend integration logs remain in `test-review-final/`.

## Requirement and acceptance matrix

| Requirement / acceptance | Evidence and boundary |
| --- | --- |
| R1–R2 / A1 | Joined native administrator/gateway four-flow fixture, real runtime/HTTP/persistence checks, preserved prompts/scoring/context/resources, batch controls and browser structure. Controlled local inference passes; no Docker-database or paid/live inference is claimed. |
| R3 / A2–A3 | Joined native administrator/gateway fixture records actual native usage service output, selected user/group/account/tokens/costs and matching fixture balance deductions. Evaluation and ordinary HTTP traffic share enforcing user/account slots. Partial observed usage is charged during cancellation. No paid request was issued; real A2 charges and production concurrency remain live-only. |
| R4–R6 / A4–A7 | Existing exact Key/group/account admission, immutable request-local auth objects, no alternate routing/fallback and post-slot recheck/cancellation fixtures passed. |
| R7 / A8 | Native actor-derived admin middleware/API, browser header-only current JWT, actual native row redaction, changed-row redaction regression and replay/ingress fixtures. No browser private token or ticket output. |
| R8 / A9 | Core/history/recovery/drain tests, real four-flow JSON restart read, synthetic native browser history retained across binary reinstall/reload. |
| R9 / A10 | Actual official-image original digest verification, selected manifest/config and image unchanged during both platform lifecycle checks; original deployment configuration and data mount unchanged. |
| R10 / A11–A12 | Both architecture Docker lifecycle runs, durable journal, idempotence, running/stopped failure restoration, planned restore, recreation detection/reapply; scoped owned menu conflict preserved. |
| R11 / A13 | Both static builds and both actual Docker platform lifecycle checks passed; incompatible source/image/platform/startup validation remains fail closed. |
| R12 / A14 | Public Sub2API docs/install/release and explicit artifact payload; no CPA environment or deployment was used. |
| R13 / A15 | Actual browser tabs/rows/details/history, native theme reload/preferences, Markdown/math and pelican runtime sandbox/CSP checks passed with synthetic records. Live deployed browser still pending. |
| R14 / A16–A17 | Real listener/collector/proof/inventory fixtures, native Docker fail-closed unsealed proof, native external Key disable conflict and planned scoped revoke/unrelated preservation. Complete real proxy inventory/external named-runner topology remains operator/live acceptance. |
| R15 / A18 | Production patch/core/resource allowlists, add-only separate overlays, fresh pre/post production path/hash equality, separate identities, negative payload inspection, compact runtime compatibility metadata and binary fixture-marker inspection passed. |

## Key boundary matrix

| Contract | Evidence / limitation |
| --- | --- |
| K1 | Raw token only private server files; operator tickets only memory/stdin Begin flow; payload/history/status/browser/log/URL outputs exclude them. Native fixture captures no token in reports/process arguments. |
| K2 | Patched exact-row DTO redaction even under external row identity changes; external lookup replay rejection remains native fixture-covered. Ordinary same-prefix rows remain unchanged. |
| K3 | Normal native Key/user/group/subscription/quota/pricing/slots controls remain authoritative. Fixture operator explicitly repairs only its synthetic Key status; product installer never repairs/re-enables/rebinds it. |
| K4 | Missing/unsealed/unsafe/drifted proof remains not-ready; real listener and spoof fixtures pass. Completeness of actual external ingress inventory is not claimed. |
| K5 | Both actual platform transactions verify native owned Key deletion, preserve unrelated Key, drain before restore, and keep pending failure visible. |
| K6 | Official binary restoration/recreation and reapply tested. Original-runtime lookup visibility/local normal-Key use remains the user-approved residual; automatic unplanned revoke is not promised. |
| K7 | Actual native disabled-row and menu conflicts are detected/preserved, exact owner/status/binding/IP/lookup conflicts remain covered in unit fixtures, missing `row_pending` rows require explicit operator repair. |

## Handoff and remaining acceptance

- Main session should dispatch the single full-scope `trellis-check` over this cohesive delivery, including the composed fixture repository/topology boundary and source-classified full-race blocker. No check agent was spawned by the implementer.
- No approved live target exists. Actual charged candy usage/account/group/tokens, deployed browser behavior, production concurrency, full independently inspected ingress inventory and deployed payload inspection remain unverified.
- The browser's only retained containers are task-owned `codexeval-fixture-b55373e4a7-{app,pg,redis}`, network/volume with the same prefix, at `http://127.0.0.1:49254`. All other completed Docker fixtures were removed by their drivers. Private synthetic credentials remain under ignored mode-700 fixture directories; do not upload those directories or logs containing private fixture input files. After checker/browser work, the main session can close the named Playwright session and remove those exact task-owned resources.
- No CPA `.env`, old deployment, paid request, commit/push/amend or task archive was performed. Unrelated `.codex/config.toml` and `.devin/` edits remain untouched. The task remains in progress.
