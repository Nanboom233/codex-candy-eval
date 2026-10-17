# Installation Simplification Implementation

## Current authoritative delivery

All paths below are under
`dist/sub2api-build/installation-revision-20261013/`.
Use `production-cleanup/`, `test-cleanup/`, `out-cleanup/` and
`release-cleanup/`. Earlier final/release/validated directories and checkpoints
are superseded historical evidence. The final full-gate log is
`integration-cleanup.log`, completed with exit 1. Native vet/integration and
frontend lint/typecheck/2797 tests plus 3 locale tests/build passed. Unit race
failed with 16 race warnings and 9 top-level failures, zero compiler failures
or timeouts. This result has no pristine-baseline equivalence claim.

Both static binaries built successfully. amd64 SHA-256 is
`fd3109997b4d971e524a13ddd87b4b4c8f1cd8aad213036b0dae14ba33be213f`;
arm64 is
`dd3760228f18cf63ab487ef6d3d433344d5c1b3d8ec89255f38a7deeb1bd129a`.
The current allowlisted archive in that directory is
`sub2api-codex-eval-v0.2.15-linux-amd64-arm64.tar.gz`, SHA-256
`fde738f78d88f522c94f82a99b01bfe6ea84209bb8abbe6d65c86f9a3671db0a`.
The historical root archive is preserved.

Final focused race, both builds, 68 helper checks, 22 installer/bootstrap
checks, arm64 14-checkpoint lifecycle, real PostgreSQL durable cleanup retry,
payload/lock/parity checks and final-byte browser smoke passed. Native amd64
lifecycle is blocked by Rosetta; loaded-process restoration was not verified.
Real paid billing, production concurrency, published bootstrap, deployed
browser/payload and native Windows acceptance remain open. Full evidence and
R/A/K mapping are in the final sections below.

## Change Boundary

Implement the approved fixed-expiry identity and private genuine-router dispatch.
Keep native scheduling, billing, concurrency, cancellation, evaluation data and history.
Replace topology activation with native administrator activation and protected installation facts.
Add maintenance mode to the same server binary and a private live-runtime control socket.
Require explicit container selection and short process output in both installation entries.
Package server binaries, installation files, metadata, operations documentation and licenses only.

Main owns task planning, role manifests, specs and the independent check.
This agent owns product sources, patches, separate test overlays and delivery evidence.
Preserve all prior ignored workspaces and unrelated workspace changes.
Do not read private deployment configuration or change live resources.
Do not commit, publish or repair unrelated native fixtures.

## Historical initial checkpoint

Loaded the complete native hook artifact, approved PRD, design and implementation plan.
Read the output skill, evaluation-core contract and target-adapter contract.
Inspected current production/test patch boundaries, build helpers and existing workspaces.
Fresh implementation workspace: `dist/sub2api-build/installation-revision-20261013/workspace`.
Pristine comparison source: `dist/sub2api-build/installation-revision-20261013/pristine`.
Both came from the pinned archive. Prior workspaces remain intact.

Added a native shared creation helper with an internal fixed-expiry argument.
Public Create still accepts its original request and expiry rules.
Added version 2 intents, row-preserving migration and native repository conditional expiry/revocation mutations.
Added a runtime-private context capability and pipe-backed genuine-router transport.
Bound the dispatcher before the external guard. No internal request carries credentials or target headers.
Added protected live-runtime maintenance and same-binary client mode before server startup.
Retired topology runtime sources and the operator executable from the implementation workspace.
Native compilation is in progress. Tests, UI, installers, patch export and artifacts remain outstanding.

## Verification

No implementation checks have run for this revision yet.
Record exact commands, results, final payload hashes and external acceptance gaps here before handoff.

## Current revision checkpoint (2026-10-13)

The checkpoint above is historical. Product implementation now includes atomic fixed expiry, durable version 2 migration, private capability dispatch through the real router, bounded pipe transport, tracked producer cancellation, native activation, protected maintenance, same-binary client mode, static administrator entry and short installation output.

The first production/test replay passed. A second independent replay is `production-final` and `test-final` below `dist/sub2api-build/installation-revision-20261013/`. Their production bytes match. Root `scripts/test-sub2api.sh check` passed. Evidence: `repo-check-final.log`. This includes 68 build-helper checks, installer checks, shell syntax, root vet/race and UI checks.

Native compilation passed after generator replay. The final-named native focused race passed service, middleware, handler, DTO and evalroute. Evidence: `test-final/focused-final.log`. Complete evaluation-core race passed. Evidence: `test-final/core-final.log`. The two pending narrowed/core sessions from the earlier checkpoint also returned zero.

The first final-named double-architecture build passed. Evidence: `build-final.log`. Those artifacts are now superseded by the activation/menu distinction below. Do not use their hashes as final delivery identities. Full upstream integration is running against that superseded revision in `integration-final.log`. A final-revision full gate will follow the new replay.

The first arm64 lifecycle exposed a real Docker limitation: UID 0 cannot read the runtime's `/proc/1/exe` under the container's capability policy. The installer now verifies that executable as the declared runtime UID 1000. The failed log remains `lifecycle-arm64.log`. Only its exact task-owned app/pg/redis containers, network and volume were removed. A new task fixture runs in `lifecycle-arm64-2.log`.

Installer recovery now handles absent backups under `backup_pending`, recorded revocation/restoration retries, stopped-container incomplete output, steps before restart/check and partial changes without a false no-change claim. Direct Python entry checks the complete payload/checksum allowlist. Installer unit discovery currently passes 18 checks.

Legacy menu deletion has a stricter boundary than activation. Activation/readiness permits the exact non-origin tuple and canonical evaluation path without adopting or changing that menu. The browser verifies its current same origin and reports foreign URLs. Duplicate or edited tuples remain conflicts. Maintenance has no browser origin. It leaves every absolute candidate untouched and reports cleanup incomplete. A legacy `menu_owned` flag alone cannot authorize removal. Separate regression coverage checks activation and cleanup behavior.

Pending work: export the last menu correction, run a fresh independent production/test replay and both builds, complete focused/frontend/full gates, run both-architecture lifecycle checks, verify local browser first activation/theme/history/CSP, run bootstrap negative checks, inspect final payload/provenance/locks and record final hashes.

No CPA configuration, deployment, paid request, commit, push or publication occurred. Native Windows, real billing, production concurrency and deployed browser/payload checks remain external acceptance.

## Local environment failures and final retry boundary

The final revision is now `production-release`, `out-release` and `test-release`. `build-release.log` records both successful static builds. `release-current` is the current allowlisted package. The native menu distinction is in this revision.

`integration-final.log` completed the superseded unit-race phase. Its redundant integration phase was stopped with only that script's descendants. `integration-release.log` then lost the entire shared Go cache during compilation. The environment failure also affected `test-release/focused-release.log`. No main-agent cleanup caused it. Those logs are retained and do not establish a passed gate.

`integration-isolated.log` now runs the complete gate with GOCACHE under the task directory, GOMAXPROCS=2 and GOFLAGS=-p=2. `test-release/focused-isolated.log` runs final native focused race using that cache.

The amd64 lifecycle is blocked by measured Docker Rosetta behavior. `/proc/1/exe` resolves to `/run/rosetta/rosetta`, whose checksum differs from the target server. The strict verifier rejects it. `lifecycle-amd64-debug.log` records the repeated loaded-executable mismatch. The candidate restored the official original, whose file hash is `9e31dc9693c550ab6e78a5f503e5a458f169a5fa8ed50238f7f60dcdef0430e0`; health was separately confirmed. Strict attestation remains unchanged. Native amd64 lifecycle needs a native Linux/amd64 runtime.

The original-runtime replay test's first 500 response came from the harness closing nc stdin before the database read finished. Keeping stdin open produced native `403 API_KEY_EXPIRED`. Only the separate replay harness changed. `lifecycle-arm64-complete.log` now runs that corrected harness against a fresh fixture.

Only exact task-owned fixture containers/networks/volumes were removed after failed attempts. Other Docker resources remain untouched. A persistent CLI controller now keeps the local browser alive across tool calls.

## Established final gate checkpoint

The complete final-revision gate finished in `integration-isolated.log`. Native vet passed. Native integration passed. Frontend lint, typecheck, 363 test files/2797 tests plus 3 locale tests and production build passed. Full unit race failed with 29 race warnings and 51 top-level failures. There were zero compiler failures. The final native focused race passed service, middleware, handler and DTO in `test-release/focused-isolated.log`. Do not repeat this unchanged broad gate.

The earlier baseline was 23 warnings/112 native top-level failures. The final run used bounded parallelism. Counts alone do not prove a pristine-baseline match. All failures remain in the complete log. No unrelated native fixture repair occurred.

Arm64 tests verified atomic expiry, repeat activation, maintenance peer denial, expiry conflicts, recreation and exact native main/Google/new-WS expiry rejection. Google returns its normal `PERMISSION_DENIED` envelope with the expiry message. The separate harness now checks that envelope. Original billing auth bypasses expiry but then returns its native unbound-group rejection. It does not return 200 for the immutable unbound accounting row. The test now records this exact read-only residual without weakening native checks. The complete corrected lifecycle is running in `lifecycle-arm64-final.log`.

The local native browser shows the static administrator entry before activation. CLI storage-state replacement did not preserve a working browser context. A protected file now seeds only the disposable native administrator session without printing its token. The persistent CLI controller keeps the browser alive. First-entry activation, retained history, themes, renderers and sandbox checks remain in progress.

Remaining dependencies: complete corrected arm64 cleanup/recovery lifecycle, finish browser checks, validate final payload/locks/provenance and produce the cohesive handoff. The amd64 Rosetta run restored the original file and passed health, but strict loaded-process verification also remains blocked by the emulator. It is not a verified successful rollback. Native amd64 lifecycle remains open.

## SQL mutation regression and corrected final revision

Arm64 planned cleanup found a plugin regression: the new JSON equality predicate passed a Go string slice directly to database/sql. The driver rejected the argument before PostgreSQL. The row stayed active and expired. The durable intent reached revoke_pending. This was not an upstream or environment failure.

The repository now serializes whitelist JSON and binds driver-compatible JSON strings. Exact identity/configuration predicates and counters remain intact. A separate repository overlay test checks real database/sql argument conversion. New independent sources are `production-validated` and `test-validated`. Their production bytes match. `build-validated.log` records fresh static builds of both architectures. `release-validated` is the current payload.

The corrected final complete gate is `integration-validated.log`. Final focused checks are `test-validated/focused-validated.log`. The extended `lifecycle-arm64-validated.log` will exercise real PostgreSQL conditional expiry migration, stable counters/row ID, actual native Redis invalidation, drift refusal, durable migration retry and exact revocation before restore.

The browser fixture upgraded to the corrected final binary with active identity and retained synthetic history. `browser-upgrade.log` records successful strict installation. Native first-use entry and activation passed. Earlier browser checks passed four tabs, math rendering, light/dark themes, preferences, retained history, same-origin authenticated API calls and isolated pelican script/storage behavior with zero captured CSP errors. Final-byte browser replay is in progress.

Shell/PowerShell bootstrap messages now distinguish tool checks, download, checksum, extraction, payload and release identity failures. Both embed identical Python. Installer/bootstrap checks pass 22 tests. Rewritten README/operation prose uses STE English. Product UI retains its Chinese language.

## Exact revocation verification and final cleanup revision

The retained `codexeval-fixture-bdec3f2d3d` PostgreSQL row proved that the native
revocation mutation succeeded. ID 2 retained quota_used=19, usage_5h=7 and fixed
2000 expiry; deleted_at was set and lookup became the exact revoked marker.
Verification failed because Ent's native soft-delete traversal added
deleted_at IS NULL to the query that already required deleted_at IS NOT NULL.
The protected intent remained revoke_pending. This second cleanup failure was
a plugin regression, not a database or emulator failure.

The regression test imports the genuine Ent runtime/interceptor. It failed
before the fix (`test-revocation-red/revocation-red.log`). The exact verifier
now uses `mixins.SkipSoftDelete(ctx)` only for the ID/owner/revoked-marker/fixed
expiry query. Ordinary GetByID stays filtered and is explicitly asserted.
No bootstrap, admission, mutation or ordinary authentication bypass changed.

The final independent sources/artifacts are under
`dist/sub2api-build/installation-revision-20261013/`:

- Editable source: `workspace/`; pristine comparison: `pristine/`.
- Production: `production-cleanup/`; test source: `test-cleanup/`.
- Production outputs: `out-cleanup/`; allowlisted bundle: `release-cleanup/`.
- Archive: `sub2api-codex-eval-v0.2.15-linux-amd64-arm64.tar.gz`, with adjacent
  `.sha256`. The historical archive at repository `dist/` remains untouched.

Binary SHA-256:

- amd64: `fd3109997b4d971e524a13ddd87b4b4c8f1cd8aad213036b0dae14ba33be213f`
- arm64: `dd3760228f18cf63ab487ef6d3d433344d5c1b3d8ec89255f38a7deeb1bd129a`
- Final archive: `fde738f78d88f522c94f82a99b01bfe6ea84209bb8abbe6d65c86f9a3671db0a`

`build-cleanup.log` records strict fresh replay, isolated Wire generation,
frozen frontend build, preserved lockfiles and both static ELF builds.
`payload-cleanup.log` records complete production/test byte parity, 15-file
release allowlist, compact compatibility metadata, production-only provenance,
archive contents and checksums. Test logs live outside the final test source.
Production patch identity is
`12a1865ff561c4c9f1c84df1c407fefb405db4093ec118d589bc651c065d9437`.
Production staged identity is
`4d28366dbdac4a597aa2f4c9bb509bc2feefa22ace8b98eb2c790251fd550fac`.

Final focused repository/service/middleware/handler/DTO race passed in
`focused-cleanup.log`. This includes all four composed native administrator,
internal-router, gateway, usage and concurrency/cancellation fixtures.
Helper checks: 68 passed (`helper-cleanup.log`). Installer/bootstrap checks:
22 passed (`installer-cleanup.log`). Root vet/race/UI evidence remains
`repo-check-final.log`; the last fix changes only native repository patch and
separate tests. Shell syntax and git diff --check passed.

Real PostgreSQL durable retry passed in `revocation-retry.log`. The retained
disposable fixture received the final verified binary and installation facts
through a fixture-only upgrade, because ordinary install correctly refuses a
deleted-row conflict. After restart, ordinary bootstrap remained
repair_required. Live maintenance resumed revoke_pending, finalized exact
revocation and the installer strictly verified original process restoration.
ID/counters/expiry/revoked marker remained intact. No direct SQL cleanup
replaced the native effect; SQL was used only to inspect that retained effect.

The fresh arm64 lifecycle passed all 14 checkpoints in
`lifecycle-arm64-cleanup.log`. Evidence is
`codexeval-fixture-63cba0a6f0/evidence.json`. It covers explicit target,
uninitialized install, repeat install/activation, peer denial, expiry drift,
recreation, original main/Google/new-WS paid rejection, native read-only
unbound-group rejection, real PostgreSQL legacy migration/counter retention,
actual Redis invalidation, interrupted migration retry, exact revocation,
strict original restoration, uninitialized uninstall and candidate recovery.
The native original unbound Key did not populate L2. The harness seeded the
exact native hashed Redis key with a credential-free stale sentinel and then
verified native deletion; this is not evidence of an original populated auth
snapshot. The additional post-rollback SQL assertion is maintained for future
replays; the completed lifecycle loaded the harness before that assertion.
The retained retry inspection separately proves the same exact SQL state.

The final-byte browser fixture upgrade passed in
`browser-cleanup-upgrade.log`. Refresh shows active state, four tabs, retained
pelican history, dark theme and an allow-scripts-only sandbox. The complete
unchanged-UI browser evidence is `browser-validated.json`: four tabs, math,
light/dark themes, retained preferences/history, isolated script/storage,
11 authenticated same-origin calls and zero captured CSP errors. Native
first-use sidebar and browser activation passed before that replay.

The preceding `integration-validated.log` completed with native vet and
integration passing, frontend lint/typecheck/2797 tests plus 3 locale tests and
build passing. Full unit race failed with 23 warnings and 14 top-level
failures, zero compiler failures. Expected negative-case frontend stderr was
not a test failure. The actual final-source full gate is
`integration-cleanup.log`, now completed with exit 1. Native vet/integration
and frontend lint/typecheck/2797 tests plus 3 locale tests/build passed. Full
unit race failed with 16 warnings and 9 top-level failures, zero compiler
failures or timeouts. `gate-cleanup-summary.json` records the exact failure
names: Grok import/quota probes, quota snapshot attribution, CyberPolicy stale
scope, Grok scheduler/admin quota, tool parameter rewrite counts, usage cleanup
dashboard recompute (two tests), and concurrent AuthService email aliases.
Failing packages are native handler/admin and service. No unrelated fixtures
were repaired. Do not present focused passes or count comparisons as a
pristine-baseline result. Do not repeat the unchanged final broad gate.

## Acceptance mapping for independent review

“Local” below means controlled source, native fixture, disposable Docker or
browser evidence. It does not mean production acceptance.

| Requirement | Evidence / remaining boundary |
| --- | --- |
| R1 | Local four-flow composed gateway test and UI passed; live accounts untested. |
| R2 | Existing evaluation core and resources preserved; complete core race and four-flow assertions passed. |
| R3 | Local composed native usage/charges and shared-slot tests passed; real billing and production concurrency remain open. |
| R4 | Atomic native creation and actual PostgreSQL fixed expiry/migration passed; protected seed/intent remain server-only. |
| R5 | Private capability and exact identity/business checks passed focused native tests. |
| R6 | Directed no-fallback and freshness regressions passed focused checks. |
| R7 | Native admin activation, static entry, redaction and nonadmin tests/browser passed. |
| R8 | Durable history/core tests and browser refresh retention passed. |
| R9 | Arm64 official-image inspection, reversible binary replacement and restore passed. |
| R10 | Arm64 install/inspect/status/reapply/rollback/uninstall/recovery and real durable cleanup retry passed. |
| R11 | Both static builds/pins/checksums passed; native amd64 lifecycle remains blocked by Rosetta attestation. |
| R12 | Public docs/install/build/release now target Sub2API; old CPA deployment was not invoked. |
| R13 | Four tabs/controls/themes/math/history/sandbox browser checks passed. |
| R14 | Internal capability readiness/external replay tests passed; original-runtime expiry rejection passed arm64. |
| R15 | Fresh independent production/test parity, exact payload and production provenance passed. |
| R16 | Mandatory selection/explicit container/bootstrap tests and native first activation passed; published entry and native Windows remain untested. |
| R17 | 22 installer/bootstrap tests and arm64 lifecycle process output passed; no secret-bearing output was required. |

| Criterion | Evidence / remaining boundary |
| --- | --- |
| A1 | Local composed administrator/internal router/native gateway/upstream fixture passed all four flows. |
| A2 | Controlled native attribution/token/charge assertions passed; one real paid candy evaluation remains unperformed. |
| A3 | Controlled native shared user/account concurrency and cancellation passed; production concurrency remains open. |
| A4 | Native exact identity/current-row and concurrent request-local-group regression checks passed. |
| A5 | Forged/public-context/authorization/target-membership focused regressions passed. |
| A6 | Native business checks remain in the pipeline and directed rejection fixtures passed. |
| A7 | Native no-fallback/ordinary scheduling/freshness regressions passed. |
| A8 | Admin isolation/redaction/external main/Google/WS tests and browser activation passed locally. |
| A9 | Core persistence/restart/interruption/retention tests and browser refresh passed. |
| A10 | Exact arm64 official image/config/backup facts and container-local replacement passed. |
| A11 | Arm64 repeat install, candidate recovery and history-retaining cleanup passed. |
| A12 | Arm64 recreation detection and exact-row reapplication passed. |
| A13 | Both static builds passed. Complete broad gate is reported separately; native amd64 lifecycle remains open. |
| A14 | Sub2API public documentation/release inputs and source/payload inspection passed; no CPA deploy occurred. |
| A15 | Complete unchanged UI browser replay plus final-byte smoke passed. |
| A16 | Capability-only admission/external rejection focused tests and actual original arm64 paid rejection passed. |
| A17 | Actual exact native revocation, deleted-row verification, durable retry, expiry-conflict refusal and strict restore passed. |
| A18 | Independent source replay/production-test parity/15-file release/archive allowlists passed; deployed production payload remains open. |
| A19 | Download/checksum/extraction/payload/selection negative tests and explicit Docker installation passed locally. Published-release online entry and native Windows remain open. |
| A20 | Native logged-in owner activation/repeat activation and interrupted identity tests passed locally. |
| A21 | Real PostgreSQL initial atomic expiry and original paid rejection passed; read-only unbound-group residual recorded accurately. |
| A22 | Real PostgreSQL stable-ID/counter migration, actual cache deletion, drift refusal and migration retry passed. Sentinel-seeded cache limit described above. |
| A23 | Process-message/exit/recovery/negative-output tests and actual arm64 operations passed; Windows execution remains open. |

| Boundary | Evidence / remaining boundary |
| --- | --- |
| K1 | Protected state/permissions/no-symlink tests and credential-free internal requests passed; seed/lookup excluded from product output. |
| K2 | Patched public credential rejection/redaction/header isolation tests passed. |
| K3 | Immutable fixed expiry/private capability/current-row/native business checks passed locally. |
| K4 | Native build/facts/router/identity readiness and uninitialized activation passed without topology inputs. |
| K5 | Actual live maintenance drain/exact revocation/durable retry/restore passed; unverifiable legacy menus remain incomplete and untouched. |
| K6 | Original normal-mode main/Google/new-WS expiry rejection passed; usage/billing reads and preexisting WS remain accepted residuals. |
| K7 | Exact current-row conflicts and actual stable-row legacy migration/retry passed; no replacement or automatic repair. |

The named browser `codexeval-revision` is closed. Final lifecycle group
`codexeval-fixture-63cba0a6f0` was removed by its harness. Exact retained groups
`codexeval-fixture-bdec3f2d3d` and `codexeval-fixture-7e03190d38` were removed:
their app/pg/redis containers, named networks and matching data volumes only.
All source/artifact/evidence files remain available. Unrelated Docker and
workspace resources remain untouched.

No CPA `.env`, production deployment, live billing, commit, push, publication
or task archival occurred. Native amd64 Docker lifecycle, native Windows,
published online bootstrap, real paid evaluation, production concurrency and
deployed production browser/payload remain external acceptance. The complete
upstream race gate cannot be replaced by the focused passes above.
