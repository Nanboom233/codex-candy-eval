# Sub2API v0.2.17 Browser Verification

## Candidate and environment

The browser used the built v0.2.17 application inside `sub2api-codex-eval-test`.
The origin was `http://127.0.0.1:60198`.
PostgreSQL, Redis, the synthetic supplier, and Sub2API ran in that single container.
The main session removed all 11 identified earlier task containers.
Unrelated containers, images, volumes, and caches remained in place.

Playwright CLI controlled a real headed browser.
The fixture provided restricted local authentication state without printing credentials.
The native frontend also requires `auth_user`.
The browser loaded that user through authenticated native `/api/v1/auth/me` before opening the administrator panel.

## Actual page and controls

The native sidebar displayed `账户评估` with its beaker icon.
The application displayed v0.2.17.
Clicking the sidebar opened `/admin/codex-eval` and loaded its evaluation iframe.
The iframe URL retained the complete origin, including port 60198.
The page displayed all four evaluation tabs and the fifth `设置` tab.
The disabled synthetic account retained its off switch and `已停用` state without a model-fetch error.

The settings tab used native-style cards, labels, switches, fields, and a right-aligned save action.
No explanatory paragraphs appeared in the form.
Evaluation refresh and clear controls disappeared while the settings tab was active.

The main session inspected actual screenshots at desktop width and a 430-pixel viewport.
Fields used two columns on desktop and one column on the narrow screen.
Labels, units, and controls remained readable.

## Save, validation, and scheduling

Entering minimum 20 and maximum 19 minutes blocked saving.
The maximum field reported `最长等待须不小于最短等待`.

Saving equal 17-minute bounds with automatic inspection enabled returned `已保存`.
The page showed the next inspection time.
The exact runtime timestamp was `2026-10-16T11:27:31.520590715Z` before and after browser reload.
Both saved wait fields remained 17.
Notification secret fields remained empty while their configured indicators remained `已配置`.

An immediate run selected only `inspection-correct` and used the synthetic supplier.
The start control became disabled while the stop control became available.
The round completed with one valid test and an unchanged priority of zero.
It scheduled the next run 17 minutes after completion.

A second immediate run followed by the stop action produced phase `cancelled` and `running:false`.
The cancelled result showed its actual context cancellation and made no priority change.

The browser restored automatic inspection off and the default 15/35-minute bounds after these checks.
Automatic priority changes remained off.
Only synthetic local accounts and a synthetic supplier were used.

## Artifacts

- `output/playwright/sub2api-v0217-settings.png`: sidebar, iframe, settings cards, and main controls.
- `output/playwright/sub2api-v0217-settings-save.png`: notification fields, configured flags, save feedback, and result summary.
- `output/playwright/sub2api-v0217-settings-narrow.png`: single-column narrow layout.

After reload and the evaluated flows, the browser console reported zero errors and warnings.
An earlier dashboard load reported one error before reload. It did not recur in the evaluation flow.
No frame-src rejection occurred.

## Final rebuilt candidate

The earlier form checks preceded the independent review's backend settings-load repair.
That repair changes startup admission, not page assets or these form contracts.

The main session then ran `python3 tests/installer/native_inspection_fixture.py`.
It replaced the tested binary and restarted the same retained container without creating another container.
The application reported Sub2API 0.2.17.
Generated browser storage included both native authentication fields and loaded successfully.

The final browser smoke opened the native evaluation iframe and started inspection from the settings tab.
The round completed with `error:""`, `running:false`, one valid test, and an unchanged priority of zero.
Automatic inspection remained off, and `next_run_at` remained null.
The browser console reported zero errors and warnings.
The fixed container remains available for reuse at the same origin.

N150/iStoreOS and the production HTTPS proxy were not exercised.
