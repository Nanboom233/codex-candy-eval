# Independent U10 Dropdown Check

## Findings fixed

No additional product fix was made by this review.
The implementer repaired the account menu placement after the browser check found clipped options near the iframe bottom.
The menu now opens above the trigger when needed and limits its scroll area to the visible viewport.

## Findings not fixed

The main browser check reported one inline `__APP_CONFIG__` nonce error after a normal parent-page reload.
This is outside U10 and does not come from the native sidebar or evaluation iframe patch.
The fresh no-cache response contains one script nonce that matches its CSP header.
Do not change the global CSP for this report.

The main session owns the final browser recheck and the final artifact update.
This review did not create, remove, or restart a Docker container.

## Verification

- `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`: passed.
- The focused UI check covers all five model selectors, custom literal values, empty catalogs, paginated account inventory, disabled-account discovery, wildcard persistence, special account names, active discovery errors, save/reload, dirty edits, stale responses, keyboard escape, and viewport menu placement.
- `node --check` passed for `web/app.js`, `web/catalog.js`, `web/components.js`, and `web/inspection.js`.
- `git diff --check` passed for the changed web and UI test files.

The implementation preserves the existing API arrays and literal model values.
It uses the native select style for model choices and checkbox menus for multiple account selection.
The U10 browser recheck remains required after the implementer rebuilds the fixed artifact.

The implementer rebuilt the static v0.2.17 binary after the viewport-placement repair.
The fixed archive was regenerated and contains only the binary, installer, and license texts.
The main session must load this candidate into the retained container before its final browser check.
