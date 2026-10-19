# Research: native-model-ui

- Query: Inspect Sub2API's native model whitelist UI and the embedded evaluation page to choose a model-control UI that matches the native interaction pattern.
- Scope: internal
- Date: 2026-10-19

## Findings

### Source observations

- `frontend/src/components/account/ModelWhitelistSelector.vue:1-33` renders selected models as compact removable tags inside one bordered field. The field shows a dropdown affordance, but it also renders a model count on the bottom row. The count is native behavior but must be omitted for this feature because the requested UI has no small explanatory text.
- `frontend/src/components/account/ModelWhitelistSelector.vue:35-91` renders a sticky search input and a scrollable model list. Each model row has a real `button` for toggling selection and a separate copy button. The selected state is a checkbox-like visual, not a native checkbox. The list is capped with `max-h-52`.
- `frontend/src/components/account/ModelWhitelistSelector.vue:217-237` filters model options by the supplied platform set and then by the search query. `:232-237` matches both model IDs and labels.
- `frontend/src/components/account/ModelWhitelistSelector.vue:240-255` toggles a model and removes it from the selected array. `:261-279` supports a separate custom-model text input and Enter handling: trim the value, ignore an empty value, reject a duplicate, append it, then clear the input. The four requested IDs are defaults, not the complete allowed set, so this custom-add behavior is required.
- `frontend/src/views/admin/RiskControlView.vue:776-819` embeds the same selector for a model filter and adds a separate model-count summary at `:816-818`. This confirms the project uses the selector for editable model sets, but the new control should omit the summary.
- `frontend/src/components/common/Toggle.vue:1-14` is the native switch shape: a `button` with `role="switch"`, `aria-checked`, and the primary-color track. `frontend/src/views/admin/SettingsView.vue:1454-1467` shows the native setting-row pattern: a label block at left and `Toggle` at right, separated by a border when the row follows another setting. The new card header should follow the embedded page's equivalent rather than introduce a different switch.
- `frontend/src/views/admin/SettingsView.vue:1469-1527` uses an inline-entry tag field. This is a different control and is not the selected reference.

### Existing embedded page and reusable boundaries

- `patches/0006-evaluation-ui.patch:2329-2350` defines the settings page as independent `.settings-card` sections. The priority section at `:2344-2349` has a header with `h2`, an unlabeled switch, and a hidden body containing fields. Add model control as the next sibling card with the same structure and no subtitle.
- `patches/0006-evaluation-ui.patch:1737-1754` defines the card, header, body, two-column grid, and the narrow-screen collapse to one column. Reuse these classes; do not add Vue components or a second card inside the model-control card.
- `patches/0006-evaluation-ui.patch:75-88` fills the inspection model picker from the cached model catalog. `:900-932` shows that the embedded page currently uses native `<select>` controls for single-model choices and reveals a text input only for the explicit `自定义模型…` option. The multi-model control should remain separate from this single-model picker.
- `patches/0006-evaluation-ui.patch:770-845` loads accounts for the selected group, calls each enabled account's `/codex-eval/models` endpoint, deduplicates model IDs, filters image models, and stores the result in `catalogCache.models.ids`. Use this already-loaded current-group inventory for dropdown options; do not add a second discovery request.
- `patches/0006-evaluation-ui.patch:1330-1365` and `:1421-1477` provide the existing field-dirty revision, visibility, validation, and autosave path. The new switch and selection field must call `markInspectionEdit` through the same form event path. Success should remain silent; `inspectionMessage` is reserved for actual validation or save errors.
- `tests/ui/sub2api-dropdowns.cjs:294-362` verifies collapsed settings, unlabeled header switches, reopening, and persistence of edits. In particular, `:319-322` requires the switch label element to have no visible caption while retaining an accessible label. Extend this test for model-control visibility and selection persistence.
- `tests/ui/sub2api-dropdowns.cjs:1-36` creates the embedded page with JSDOM and exposes helpers from the actual bundled files. This is the correct harness boundary for interaction checks; mounting `ModelWhitelistSelector.vue` is not an option because the production page is vanilla HTML/JS.

### Preferred interface

Copy `ModelWhitelistSelector.vue` directly into the embedded page's HTML, CSS, and event handlers.
The user's implementation feedback requires the native layout, not a new layout with similar interactions.

1. Add a sibling `.settings-card` after the priority card. Header copy is only `模型控制`; put the switch on the right with `role="switch"`, `aria-label="模型控制"`, and no visible switch caption. Keep the body hidden while off.
2. Show `降级移除模型` with the native two-column tag grid and 6 px gaps. Keep gray 12 px tags, 8 by 4 px padding, and 14 px model icons. Use the four default IDs, with removable tags and the native bottom separator and right-aligned arrow.
3. Keep custom entry below and outside the selection field. Copy the `自定义模型名称` label, input, and `填入` text button. Trim on click or Enter, ignore empty input, prevent duplicates, append the value, and clear the input. An empty catalog must not disable custom entry.
4. Keep the native absolute dropdown below the field, with a 4 px gap, 8 px corners, and shadow. Copy its bordered search area and 208 px maximum list height. Use current-group IDs and retain selected IDs absent from that catalog.
5. Keep the native option selected state and model ordering from `compareModels`. Omit counts, quick actions, copy actions, sync actions, help text, tooltips, and success text.
6. Preserve selection order and the existing autosave path. Keep the native empty-field state without an explanatory paragraph.

Keep the visible native layout while adding a keyboard-accessible trigger and focusable options. Do not nest buttons. Tab reaches controls, Enter/Space activates them, and Escape closes the dropdown and returns focus to the trigger.

### Integration owners and bounded checks

- `patches/0006-evaluation-ui.patch` owns the embedded markup, vanilla event handlers, and CSS. Keep the new control in the existing inspection form and reuse `catalogCache.models.ids`, `inspectionDirtyFields`, `updateInspectionVisibility`, and `saveInspectionSettings`.
- The backend/inspection worker owns only the persisted setting and account update behavior already settled by the parent task. This UI report does not change the identity whitelist or `blacklist:blacklist` placeholder design.
- Use the bounded check table in [implement.md](../implement.md) as the test scope. Extend existing settings cases for the new widget and array field.
- Use the existing JSDOM harness for selection, search, custom entry, and persistence assertions. Inspect wrapping and keyboard behavior in the actual render.
- Reuse the catalog and save fixtures. These source observations do not require separate discovery, reload, browser, or save-concurrency suites.
- Keep actual validation/save errors in the existing status element. Do not add helper paragraphs, model counts, selection counts, syntax instructions, or explanatory tooltips.

## Caveats / Not Found

- `ModelWhitelistSelector.vue` does not implement keyboard navigation for its clickable container or option list beyond browser behavior of the inner buttons. The embedded adaptation must add the missing trigger focus and Escape return explicitly.
- The Vue component depends on `ModelIcon`, i18n, Vue state, and account APIs (`ModelWhitelistSelector.vue:147-158`, `:293-351`). It cannot be mounted directly in the vanilla evaluation page; only the interaction pattern and existing page data may be translated.
- No separate native component exposes a multi-select with the exact current-group catalog contract. The proposed control therefore uses the native whitelist visual pattern with the embedded page's existing catalog and autosave primitives. The catalog is a suggestion source, not an allowlist: an empty or unavailable catalog must not disable the custom-ID input.
- No browser or provider request was used for this inspection, and no product code or build was changed.
