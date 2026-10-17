# Account Refresh and Error Display Research

## Scope and Source

This report covers R1 and R4. It also identifies row-layout dependencies for R3.
The research used unchanged production patches from the current working tree.
The isolated patched package is `/tmp/inspection-controls-research-20261018/backend/internal/codexeval`.
All source paths below refer to that applied package unless an upstream or repository path is explicit.

Production ownership remains `patches/0006-evaluation-ui.patch`.
The probe reuses the existing UI harness and installed original-source dependencies.
It uses synthetic API responses. It makes no browser, provider, or container mutation.

## R1: Confirmed Readiness Flicker

The exact path is:

1. `web/app.js:223` calls `load` for ordinary polling and foreground refresh.
2. `web/app.js:227` sets `accountListReady` to false before `render`.
3. `web/catalog.js:57` clears the same value again in `refreshCatalog`.
4. `web/catalog.js:96` makes `manualCatalogReady` false during that request.
5. `web/app.js` routes this through `pageBlocker` and `runBlocker` into each row's button state.
6. The account read succeeds. The catalog sets readiness true at `web/catalog.js:79`.
7. `load` renders the same rows again, restoring button admission.

The keyed renderer already preserves row nodes. Replacing it is not the required repair.
The source polls every 2.5 seconds during active work and every 20 seconds otherwise at `web/app.js:261`.

Executed probe:

```sh
node .trellis/tasks/10-18-inspection-controls/research/ui-refresh-probe.cjs /tmp/inspection-controls-research-20261018
```

| Observation | Before poll | During delayed account read | After success |
| --- | --- | --- | --- |
| Test button disabled | false | true | false |
| Account selected | true | true | true |
| History expanded | true | true | true |
| Existing row retained | yes | yes | yes |
| Button reason | none | `正在读取账户列表，请稍候` | none |

The probe confirms a control-state flash during an unchanged successful refresh.
It does not establish whole-table replacement or the exact production browser sequence.

## R1: Minimal Repair Boundary

Use the last successful same-group inventory while a background request runs.
Remove the unconditional readiness invalidations from `load` and the successful-snapshot path in `refreshCatalog`.
Treat request activity and inventory validity as separate facts.
Initial load, changed group, session reset, and a real inventory failure still invalidate admission.
After a real failure, update the error and readiness together before rendering.
Retain the current request generation and group checks before any result commit.
Keep model-discovery errors separate from inventory failures because a valid custom model can still run.
Retain model caching when enabled account IDs are unchanged.

Add one delayed-success case to `tests/ui/sub2api-dropdowns.cjs`.
Assert stable button state, row identity, selection, and expanded records in that case.
Reuse the harness's existing inventory failure, recovery, membership, and obsolete-group checks.
The check table in `../implement.md` defines the complete check scope. These observations do not require another suite.

## R4: Confirmed Header Error Aggregation

`web/components.js` creates `inspection-runtime-error` and `candy-settings-status` inside `credentialCard` for Candy only.
The patch owner is `patches/0006-evaluation-ui.patch:1113`.
`web/inspection.js:78` reads the runtime error element.
Lines 79–85 combine the round error, every account error, and every notification error into one text value.
`web/inspection.js:104` uses the same element for a failed inspection-state read.
`inspectionMessage` also writes Settings errors into `candy-settings-status`.
The rule `.candy-head > .error` gives the error element a full header row.

The same probe rendered synthetic failures and observed this header text:

```text
runtime failure；A only：record failure；A only：通知：notification failure
```

This confirms the exact block named by the user.

## R4: Minimal Repair Boundary

Remove both error placeholders from the Candy header and update every writer together.
Keep `inspection-status`, `inspection-progress`, and `inspection-next` in that header.
Delete runtime and notification error prompts rather than moving the runtime feedback element into Settings.
Keep `inspection-save-status` for actual settings validation and save failures.
Make `inspectionMessage` write only Settings feedback after removing `candy-settings-status`.
Put account participation save failures beside the new row switch through its own per-account state.

Test errors remain in Candy result records and history.
Remove visible coordinator, notification, state-read, and runtime-action error prompts from both Candy and Settings.
Keep their backend state and actual error handling. Keep authentication handling and future polling attempts.
Remove failure captions that would recreate the deleted prompt. Do not add a replacement summary or card.
Keep running progress and the exact `下次巡检：` label with its scheduled timestamp.
The user explicitly retained that label and time after requesting direct removal of runtime and notification error prompts on 2026-10-18.
Remove the unused full-width error CSS rule only if no retained header error uses it.

Update existing assertions at `tests/ui/sub2api-dropdowns.cjs:266`, `:321`, and `:493` with the removal and timestamp-retention contract.
Check that polling and settings saves work when the old header elements are absent.

## R3: Shared Row Layout Dependencies

`credentialCard` currently appends the native account `启用` column for all four evaluation lists.
`credentialRow` inserts `enableSwitch` in `row-meta`. Its switch controls native `schedulable`.
Add an optional Candy-only inspection cell and its matching column title.
Use a separate dataset key and per-account pending state for inspection participation.
Keep `setCredentialEnabled` and its `/schedulable` request unchanged.

The Candy desktop grid in `web/style.css` needs one additional inspection track.
Do not change the other tab grids.
The narrow layout hides headers and displays `row-meta` as a wrapping flex row.
Its current generic `.row-meta > .switch::before` caption is `启用`.
Give the new control its own visible narrow-screen caption and avoid applying the account caption to it.
The `.row-meta.untested` rule must not hide the new inspection switch or its actual save error.
Keep label clicks out of history expansion and test-start handlers.

## Research Limits

No product source changed. No full build or actual browser acceptance check ran in this research pass.
Column layout and switch interaction need one browser pass on the final candidate.
