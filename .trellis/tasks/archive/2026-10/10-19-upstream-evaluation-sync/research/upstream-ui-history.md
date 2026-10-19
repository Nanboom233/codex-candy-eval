# Research: Upstream UI and Saved State

- Query: Compare upstream `7a49068` with shared base `6b46d4f` and local `main` `2aedffe`. Map UI, saved test preferences, result history, state migration, and user instructions to the native Sub2API contracts.
- Scope: mixed
- Date: 2026-10-19

## Findings

### Upstream replaces two test pages with one

At the shared base, the CPA page has four test tabs: Candy, probe-based fingerprint, ModelTrace, and Pelican (`6b46d4f:internal/plugin/web/ui.html:28-32`). Commit `63098bd` replaces the two fingerprint methods with one `指纹检测` page. The current upstream page has three tabs: Candy, fingerprint detection, and Pelican (`7a49068:internal/plugin/web/ui.html:28-32`). It removes the ModelTrace page and the old fingerprint sampling-mode selector. The new fingerprint page exposes model and per-credential concurrency controls (`7a49068:internal/plugin/web/ui.html:67-85`). Its implementation details and request cost belong to the separate evaluation-algorithm research line.

The page describes three numeric challenges per detection and replaces the old 60/200/400 mode choices with concurrency control (`7a49068:internal/plugin/web/ui.html:67-82`).

Commit `d093ee8` makes test names consistent across the page and management API. The final page uses `TESTS = ["candy", "fingerprint", "pelican"]`; its shared `running` record names the active test (`7a49068:internal/plugin/web/app.js:9-12,192-198`). Batch selection uses the visible, available credentials. A row-level stop posts credential IDs to `/cancel` (`7a49068:internal/plugin/web/app.js:348-385,387-408`). The CPA page restores a saved tab only if it remains in `TESTS`; a saved `modeltrace` tab falls back to Candy (`7a49068:internal/plugin/web/app.js:324-333,533-535`).

This is not a drop-in replacement for the native page. The native page has four evaluation tabs plus `设置`, native account IDs and groups, Sub2API bearer authentication, and native evaluation routes. Its manual fingerprint and ModelTrace controls have separate saved preferences (`patches/0006-evaluation-ui.patch:8-14,75-87,1210-1213,1824-1827`). Its settings tab shares the Candy group and filters. Initial Settings entry must load settings before shared groups (`patches/0006-evaluation-ui.patch:351-352,560`; `tests/ui/sub2api-dropdowns.cjs:87-112`). Keep that startup owner, the Settings tab, account-level inspection controls, model control, priority, and degraded tags when adapting upstream page components. Do not import the upstream CPA routes, credential identity, or authentication flow.

Upstream commit `f2aeab2` removes helper text from Candy and Pelican effort selectors (`7a49068:internal/plugin/web/ui.html:49-52,97-100`). The native selectors already have no effort tooltip. This change needs no native behavior change. Keep the current no-small-print UI contract.

### CPA state migration does not apply to native history

Upstream state schema 1 stores three maps: `candies`, `fingerprints`, and `pelicans`, each keyed by a CPA credential ID. Every record has an ID and status (`7a49068:internal/plugin/state.go:35-46,47-61`). The migration treats a file without `schema_version` as version 0, migrates it, makes a byte-for-byte `.v0.json` backup, and saves the new layout (`7a49068:internal/plugin/migration.go:11-47`; `7a49068:internal/plugin/state.go:82-113`).

The migration converts old `results` to `candies`, deriving IDs from timestamps and statuses from `ok`, `error`, and `skipped`. It copies `detections` to `fingerprints` and adds statuses to Pelican records. It clears every other root key (`7a49068:internal/plugin/migration.go:50-79,82-94`). Therefore it drops pre-0.4 probe-fingerprint and ModelTrace histories. The migration test asserts this for a legacy fixture (`7a49068:internal/plugin/state_test.go:43-83`). Malformed, null, invalid-version, and newer-version files remain unchanged and block saves (`7a49068:internal/plugin/state.go:95-113,115-118`; `7a49068:internal/plugin/state_test.go:85-100`).

The native file is a different format and data set. It lives at `<data-dir>/codex-eval/history.json`; it stores `results`, `fingerprints`, `modeltraces`, and `pelicans`, keyed by decimal Sub2API account ID (`patches/0001-evaluation-core.patch:1951-1957,2195-2207`; `.trellis/spec/backend/database-guidelines.md:21-23,60-71`). Manual evaluation history keeps the latest 20 records per account and type. Native cancellation uses a run ID and saves the cancelled result (`patches/0001-evaluation-core.patch:1979-1981,2098-2114`; `tests/backend/internal/codexeval/runtime_test.go:111-151`). The existing runtime test covers all four test types and verifies history retention and corrupt-file protection (`tests/backend/internal/codexeval/runtime_test.go:79-103,145-171`).

Do not run upstream `readState` or `migrateUnversionedState` on the native file. The CPA migration expects CPA-specific record fields and credential IDs. It also clears old root keys, so applying it could silently discard native `fingerprints` and `modeltraces`. The upstream migration does not convert native settings or account metadata. Native inspection settings remain in Sub2API PostgreSQL settings; verdicts and removed model names remain in native account `extra` (`.trellis/spec/backend/native-inspection.md:42-57,153-178`).

The user approved a breaking replacement of the old fingerprint and ModelTrace histories.
Do not migrate, convert, or back up these retired records. Do not add a legacy page.
Keep native Candy and Pelican records and formats unchanged.
Use the raw read envelope and native marker defined in [code research](code-ui-history.md).
The final user adjustment requires immediate minimal migration after an unversioned file loads.
Save retained Candy/Pelican records in the versioned snapshot through the existing atomic saver and omit retired records.

### Native preference and runtime boundaries

Upstream stores per-page preferences in browser local storage and uses a CPA credential ID. The native page stores its manual test preferences under `sub2api-codex-eval.prefs` and related keys. Candy inspection settings instead use native PostgreSQL and account `extra`; those settings are not upstream page preferences (`patches/0006-evaluation-ui.patch:8-14,604-608`; `.trellis/spec/backend/database-guidelines.md:1-23`; `.trellis/spec/backend/native-inspection.md:42-57,153-178`).

The user approved the three-test layout with native Settings.
Keep automatic inspection, `立即巡检`, model control, and priority independent of manual fingerprint tests.
Keep native evaluation entry points and Candy/Pelican behavior.
Retain one active run per account and cancellation by exact `run_id`.
Sources: `.trellis/spec/backend/native-inspection.md:73-89,129-147` and `patches/0001-evaluation-core.patch:1979-1981,2098-2114`.

### User README changes

The current Chinese README describes four test pages, the old fingerprint modes (60, 200, and 400 requests), and ModelTrace (`README.md:32-90`). If the replacement is approved, update those sections to the three upstream names, remove the obsolete mode choices, and describe the new fingerprint controls and history from the actual native implementation. Keep the native installer, account-evaluation route, Candy inspection steps, model control, priority behavior, recovery, and build commands (`README.md:5-30,92-196`). Do not copy upstream CPA marketplace, plugin-install, or external `crontab` instructions.

Upstream changes project attribution from ModelTrace to LM Fingerpoint Detector (`7a49068:README.md:91-95`). If the native replacement adopts that detector, update the Chinese Credits and third-party notices with its source and MIT license. Remove ModelTrace attribution only if its code and data are no longer shipped.

### Commit inventory for this line

The 11 upstream commits after the shared base divide as follows:

- `0de3802`, `c363f6b`, and `7a49068`: release commits. No separate UI or persistence design change is assigned to this line.
- `63098bd`: replaces probe fingerprint and ModelTrace tests with fingerprint detection; changes UI, histories, test assets, and README. The detector internals belong to the algorithm owner.
- `009c26e` and `96deb3c`: update README screenshot assets.
- `6573867`: adds retries for failed or unusable answers. Retry and transport semantics belong to the other research owner; the UI shows attempt counts through shared record metrics (`7a49068:internal/plugin/web/components.js:68-73`).
- `f2aeab2`: removes effort tooltips from the two remaining effort controls.
- `d093ee8`: versions CPA state, migrates unversioned records, and unifies test names, page IDs, routes, and active-test state.
- `dedb1bd`: changes repository instructions, not user behavior.
- `b4c6b60`: updates the Codex TUI request transport; the other research owner covers it.

## Files Found

- `.trellis/tasks/10-19-upstream-evaluation-sync/prd.md` — approved scope and breaking history decision.
- `.trellis/spec/backend/native-inspection.md` — native Candy inspection, settings, account metadata, and model-control contracts.
- `.trellis/spec/backend/runtime-installation.md` — native evaluation APIs, task state, cancellation, page entry, and history behavior.
- `.trellis/spec/backend/database-guidelines.md` — native settings, account metadata, and four-map history format.
- `patches/0001-evaluation-core.patch` — native evaluation runtime, IDs, persisted history, and cancellation.
- `patches/0006-evaluation-ui.patch` — native tabs, manual preferences, Settings-first loading, and controls.
- `tests/backend/internal/codexeval/runtime_test.go` — four test types, cancellation, retention, and invalid-history coverage.
- `tests/ui/sub2api-dropdowns.cjs` — Settings-first entry, per-test preferences, model pickers, inspection controls, and current UI fixtures.
- `README.md` — Chinese installation, usage, inspection, restoration, build, and credits instructions.
- Upstream `internal/plugin/web/ui.html`, `app.js`, `components.js`, `state.go`, `migration.go`, `state_test.go`, `README.md`, and `scripts/test-ui.cjs` — replacement UI, page state, migration, tests, and user workflow.

## Code Patterns

- CPA page tab set and fallback: `7a49068:internal/plugin/web/app.js:9-12,324-333,533-535`.
- CPA batch targeting and per-credential cancel: `7a49068:internal/plugin/web/app.js:348-408`.
- CPA schema and record fields: `7a49068:internal/plugin/state.go:35-61`.
- CPA migration and preservation backup: `7a49068:internal/plugin/migration.go:11-47,50-79`; `7a49068:internal/plugin/state.go:82-118`.
- CPA legacy-history loss test: `7a49068:internal/plugin/state_test.go:43-83`.
- Native history path and four maps: `patches/0001-evaluation-core.patch:1951-1957,2195-2207`.
- Native run-ID cancellation and four-type runtime coverage: `patches/0001-evaluation-core.patch:1979-1981,2098-2114`; `tests/backend/internal/codexeval/runtime_test.go:79-103,111-171`.
- Native Settings-first loading and controls: `patches/0006-evaluation-ui.patch:351-352,560`; `tests/ui/sub2api-dropdowns.cjs:87-112,365-455`.

## External References

- Upstream source: `haowang02/cpa-plugin-codex-candy-eval`, commit `7a49068` (`v0.4.1`), compared with shared base `6b46d4f`.
- Detector attribution and license: `7a49068:README.md:91-95`; upstream `internal/plugin/data/lm-detector/LICENSE`.

## Related Specs

- `.trellis/spec/backend/native-inspection.md`
- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/database-guidelines.md`
- `.trellis/spec/backend/quality-guidelines.md`

## Caveats / Not Found

- The CPA migration is not a migration contract for the native Sub2API JSON file. Its different paths, map keys, record fields, settings store, account IDs, and cancellation API make direct reuse unsafe.
- Native history has no schema marker today. The accepted replacement uses the minimal marker defined in `code-ui-history.md`.
- The exact model-detector request cost and result semantics are left to the separate algorithm researcher.
- This report reads source and tests only. It does not run checks, builds, installations, provider calls, or write product files.
