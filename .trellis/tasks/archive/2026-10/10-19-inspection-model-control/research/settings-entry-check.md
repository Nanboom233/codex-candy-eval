# Settings-entry independent check

## Findings fixed

- File: `patches/0006-evaluation-ui.patch`
  - The initial tab restore now calls `switchTab(..., { loadSettings: false })`.
  - `switchTab` keeps `loadSettings: true` by default, so later user entries to
    Settings still call `loadInspectionSettings()`.
  - `start()` remains the owner of the first settings read and waits for it
    before `load({ refresh: true })` starts group and account discovery.

The source trace matches the reported failure. Before the correction, restored
Settings started a request while `inspectionSettingsLoading` was true. The
`start()` call returned early, and the catalog captured group `0`. The settings
response then changed `selectedGroup`, so the obsolete group response was
discarded. The final source keeps the existing obsolete-response guard and
removes only the competing initial read.

## Findings not fixed

None. No design, interface, or out-of-scope issue was found.

## Verification

- Regression: passed in the existing UI harness. The pre-fix harness failed with
  the selected group shown as unavailable and no account rows; the final case
  loads group `202` and rows `22, 33, 55, 66, 77` across Candy, Fingerprint,
  ModelTrace, and Pelican. It also checks one startup settings read and a new
  read after a later Settings entry.
- Browser: the implementation evidence records an actual Chrome check through
  the native sidebar, `/admin/codex-eval`, the Settings-first landing, all four
  evaluation tabs, and native re-entry. The group and account requests use the
  same-origin synthetic API fixture; no provider request was made.
- Delivery: the implementation evidence records a passing normal server and
  installer build. The unchanged backend and installer checks remain valid.
- Static syntax: `node --check` passes for the effective `app.js`,
  `inspection.js`, and `tests/ui/sub2api-dropdowns.cjs`.
- Diff hygiene: `git diff --check` and the cached check pass for the affected
  patch and UI test.

## Limits

The browser check used synthetic API responses because the retained Docker
fixture was unavailable and its `/app/sub2api` path is a dependency launcher.
No live account, provider quota, or production data was used. No additional
build or full check was needed for this unchanged behavior outside the repair.
