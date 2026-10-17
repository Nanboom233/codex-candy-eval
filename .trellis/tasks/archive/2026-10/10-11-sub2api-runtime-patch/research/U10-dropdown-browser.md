# U10 Browser Verification

## Candidate and environment

Playwright CLI controlled the existing headed browser at `http://127.0.0.1:60198/admin/codex-eval`.
The test used the rebuilt v0.2.17 binary in `sub2api-codex-eval-test`.
The native sidebar and same-origin iframe remained available at the retained port.
No additional container or paid provider request was used.

## Controls and persistence

All four evaluation tabs displayed model dropdowns with `inspection-literal` and `自定义模型…`.
The browser selected the catalog model, opened the custom input, entered a literal, and returned to the catalog selection on each tab.

The settings model dropdown displayed the saved model and the custom-model action.
Saving `inspection-custom-U10` retained that exact model after browser reload.

The account dropdown displayed searchable checkboxes, including `inspection-disabled · 已停用`.
The browser selected two included accounts and the disabled excluded account.
It added the unmatched wildcard `unmatched-U10*` through the dropdown search action.
After save and reload, both included accounts, the unmatched wildcard, and the excluded account remained selected.
The browser then restored model `inspection-literal`, include `inspection-correct`, and an empty exclusion list.
Automatic inspection and automatic priority changes remained off.

The form displayed labels, units, current state, and save feedback without explanatory paragraphs.
The added model-mapping and effort explanation tooltips were absent.
Original evaluation descriptions and result data remained available.

## Visual finding

The first desktop screenshot exposed a clipped account menu near the iframe bottom.
Only its search box was visible before the menu extended outside the frame.
The implementer repaired the menu placement and rebuilt the fixed binary and archive.
The main session loaded that candidate into the same retained container.
The final desktop screenshot shows all options above the trigger near the iframe bottom.
The 430-pixel screenshot shows a complete menu within the field width and iframe viewport.
The narrow page has no horizontal document overflow.
The selected account remains checked, and the disabled account remains clearly labeled.

Initial evidence: `output/playwright/sub2api-U10-account-dropdown.png`.
Final evidence: `output/playwright/sub2api-U10-dropdown-final.png` and `output/playwright/sub2api-U10-dropdown-narrow.png`.

## Native shell diagnostic

One normal browser reload reported a blocked inline configuration script on `/admin/codex-eval:13`.
The evaluation iframe and the save/reload interactions remained functional.
A fresh `cache:no-store` request returned an inline configuration nonce that matched the CSP header.
The independent checker traced this configuration injection to the unchanged upstream native shell.
Its 304 branch can reuse cached HTML while the response middleware generates a fresh CSP nonce.
That path explains the observed cached-body and response-header mismatch.
The global CSP was not changed.

## Result

The U10 model/account interactions, persistence, copy removal, and corrected menu rendering passed.
The final fixed binary and archive include the reviewed menu repair.
The existing backend, installer, and evaluation evidence remains applicable.

N150/iStoreOS and production HTTPS/proxy execution remain unperformed.
