# U10 Dropdown Implementation

## Scope

The user requested fewer explanations and visible dropdown controls.
This change affects the standalone evaluation page only.
Native algorithms, transport, API fields, settings storage, and installer behavior remain unchanged.

I inspected native `Select.vue`, `GroupSelector.vue`, and `SettingsView.vue` before implementation.
The page reuses its existing native select style and semantic theme tokens.
The UI/UX Pro Max keyboard-navigation guidance applies to the new account dropdown.
No framework or dependency was added.

## Model controls

All four evaluation tabs and the settings tab now have a visible native select.
The select contains the existing raw native model suggestions and the current literal value.
The final option, `自定义模型…`, opens the existing model input inside the same field.
The original model input remains the request and preference value owner.
Catalog updates retain selected, saved, and unfinished custom values.
An empty or failed catalog does not prevent a saved or custom literal model.
Names remain literal at the API boundary. No alias, mapping, or catalog membership requirement was added.

## Account selection

The include and exclude textareas became searchable dropdowns with native checkboxes.
The controls independently read every page of the native OpenAI account list without an evaluation-group filter.
The option label contains account identity and the actual disabled state when applicable.
Selected names become escaped literal shell patterns. Account IDs are not persisted.
Existing wildcard, email, unmatched, and escaped pattern values remain selected and survive an unchanged save.
The dropdown search also permits an explicit custom pattern through its add action.
The server still receives the existing `include_patterns` and `exclude_patterns` string arrays.

The settings model catalog uses this complete inventory.
Both model-discovery paths skip `schedulable:false` accounts.
Active-account catalog errors remain visible.
An account enable change or explicit page refresh invalidates the inspection catalog.
Changing the evaluation group does not restrict inspection choices.

The account dropdown supports pointer use, Tab traversal, ArrowDown opening, Escape closing, and outside dismissal.
Selection does not replace the focused checkbox.
Long labels wrap inside the menu. The trigger truncates long selected labels.
The menu uses a bounded scroll area and remains inside the responsive field width.

## Removed copy

- Removed the model-mapping and model-restriction explanation tooltips from all four evaluation tabs.
- Removed the `none` and upstream-effort explanation tooltips from Candy and Pelican controls.
- Removed the account-pattern `每行一项` instructions with the old textareas.
- Removed the unused textarea style.

No replacement helper paragraph, footnote, or explanatory tooltip was added.
Visible labels, units, actual errors, progress, secret configuration state, and save feedback remain.
Original prompts, evaluation principles, result data, attribution, and destructive-action confirmation remain.

## Async form behavior

Account/model catalog updates change choices without rewriting the settings form.
Initial settings reads preserve fields the administrator already edited.
A save response preserves edits made after that request started.
The save status identifies a newer unsaved edit.
Aborted or obsolete inventory responses cannot replace the current inventory.

## Changed files

- `native/backend/internal/codexeval/web/ui.html`: five model dropdowns and two account dropdowns.
- `web/components.js`: shared model selection and account multiple-selection controls.
- `web/catalog.js`: shared raw model discovery and independent paginated inspection inventory.
- `web/inspection.js`: array selection, save serialization, and dirty-field preservation.
- `web/app.js`: picker initialization, model preference preservation, and catalog invalidation.
- `web/style.css`: native account dropdown layout and focus/overflow behavior.
- `tests/ui/sub2api-dropdowns.cjs`: focused tests with the existing upstream jsdom dependency.

## Validation

`SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs` passed.
It covers all model selectors, custom values, empty catalogs, all-account pagination, and disabled-account discovery.
It also covers saved patterns, special account names, dropdown selection, catalog errors, save/reload, dirty edits, and stale responses.
JavaScript syntax and diff whitespace checks passed.

`python3 scripts/build-sub2api.py --source /tmp/sub2api-codex-eval-v0217-check` passed.
This included three native i18n tests, vue-tsc, Vite, and the static Linux/amd64 Go build.
`python3 scripts/package-sub2api.py` produced the fixed archive.
The current binary is `dist/build/sub2api-codex-eval` at 125599904 bytes after the positioning repair.
The current archive is `dist/build/sub2api-codex-eval.tar.gz` at 38698426 bytes.
Its seven entries contain the binary, installer, and license texts. No development fixture or test is included.

The existing fixture entry updated the same retained `sub2api-codex-eval-test` container.
The origin remains `http://127.0.0.1:60198`, and the runtime reports v0.2.17.
Its native settings save and authenticated fixture preparation passed.
An HTTP check confirmed the served rebuilt page contains all five model selects and both account dropdowns.
The served page contains no datalist or removed model/effort explanation tooltip.

The main session owns the actual browser dropdown check and independent Trellis review.
Neither check is claimed complete by this implementation report.
The existing backend, installer, and earlier native integration evidence remains applicable.
No production or paid request ran. No extra container was created.

## Browser follow-up: clipped account options

The main session found that an account menu near the iframe bottom showed only its search field.
The first implementation always placed the menu below the trigger.
The iframe viewport clipped its option list even though the settings card allowed overflow.

`positionAccountPicker` now compares available space above and below the trigger inside the current iframe viewport.
It opens upward when the lower space is insufficient and the upper space is larger.
It subtracts the measured search field, padding, border, and gap from the available option height.
The option list scrolls within the remaining space. The menu also has a viewport-based maximum height.
Resize, viewport changes, and page scrolling update this placement.
The control retains its existing responsive field width.

The focused dropdown checks now cover upward opening, a 260-pixel iframe height, and repositioning after scrolling.
All dropdown checks, JavaScript syntax, and affected-file whitespace checks passed.
These geometry assertions use synthetic layout measurements. They do not replace the main session's browser recheck.
No Docker action ran for this focused repair.
The existing build entry rebuilt the binary successfully. Its automatic native frontend build also passed.
The fixed archive was regenerated after the repair. The build log is `/tmp/sub2api-U10-position-build.log`.
The main session can reload this candidate into the retained container and recheck the reported menu position.
Earlier browser save, reload, custom-model, and four-tab evidence remains valid for this geometry-only change.
