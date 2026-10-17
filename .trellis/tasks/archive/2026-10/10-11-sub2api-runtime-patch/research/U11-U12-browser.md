# U11 and U12 Final Browser Verification

## Candidate

Playwright CLI controlled the existing headed browser and native evaluation iframe at `http://127.0.0.1:60198`.
The final rebuilt v0.2.17 binary ran in the same retained `sub2api-codex-eval-test` container.
PostgreSQL, Redis, and the synthetic supplier remained inside that container.
No provider evaluation or production request ran during these UI checks.

## Simplified settings

The include and exclude account controls are absent.
The page retains the Free switch inside automatic inspection configuration.
The idle waiting/off messages and notification configured-status captions are absent.
The notification form retains its labels, blank secret inputs, and save action.

With priority adjustment off, both priority inputs are hidden.
Turning it on reveals normal priority 10, minimum 1, and degraded priority 99.
Entering zero failed browser validity with `rangeUnderflow:true`.
A valid save and reload retained normal priority 10.
Turning adjustment off hid its fields again.

These initial U11 controls were checked before the independent review's invalid-loaded-pair repair.
The final U12 candidate includes that repair and retains the same normal priority and lower bound.
Focused runtime, native service, and UI tests cover invalid-load blocking and explicit-save recovery.

## Collapse and retained values

Both header switches have accessible names and no visible captions.
With automatic inspection off, model and scheduling fields are hidden.
With priority adjustment off, its two inputs are hidden.
Manual start remains enabled, and the stop control remains visible outside the collapsed configuration.

Turning both switches on reveals the expected model, timing, and priority controls.
The model remains `inspection-literal`, and normal priority remains 10 with minimum 1.
The browser changed the wait range to 16/36, collapsed the configuration, and saved.
After reload, the fields remained collapsed. Expanding them showed the saved range and model unchanged.
The browser restored 15/35 and saved with both switches off.

The main session inspected desktop screenshots for collapsed and expanded settings.
It also inspected the 430-pixel layout, including the relocated manual and save actions.
The narrow page had no horizontal document overflow.

## Evidence and limits

- `output/playwright/sub2api-U11-settings-off.png`
- `output/playwright/sub2api-U11-settings-on.png`
- `output/playwright/sub2api-U11-notifications.png`
- `output/playwright/sub2api-U12-collapsed.png`
- `output/playwright/sub2api-U12-expanded.png`
- `output/playwright/sub2api-U12-narrow.png`

Independent reviews: `U11-settings-check.md` and `U12-settings-collapse-check.md`.
The final binary and archive are at the fixed `dist/build/sub2api-codex-eval` paths.
The previously recorded upstream parent-page cached-304 CSP nonce error recurred on normal reload.
The evaluation iframe and all changed settings interactions worked.
N150/iStoreOS and production HTTPS/proxy execution remain unperformed.
