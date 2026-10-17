# U14 Settings Browser Check

The final U14 v0.2.17 binary passed actual settings checks in the retained `sub2api-codex-eval-test` container.
The fixture reused its existing origin, `http://127.0.0.1:60198`.

## Reload before blur

The browser selected the existing custom-model option and typed a value in each manual evaluation tab.
Each input still had focus immediately before navigation reloaded the page.

| Tab | Restored model |
| --- | --- |
| Candy | `u14/candy` |
| Fingerprint | `u14/fingerprint` |
| ModelTrace | `u14/modeltrace` |
| Pelican | `u14/pelican` |

All four values survived reload without blur.
The focused UI checks also cover the existing number, effort, and mode preferences.

## Native settings save and reread

The browser edited the wait fields to 17 and 37 minutes, collapsed automatic inspection, and clicked `保存设置`.
The actual native GET returned both saved values with `enabled:false`.

While the browser showed Candy, another native API update saved 18 and 38 minutes.
Returning to settings read both new values.

The browser then entered an unsaved notification username.
Another native API update saved 19 and 39 minutes.
Returning to settings read the new wait values and preserved the unsaved username.
Both secret inputs remained empty.
The native API restored the original settings after these checks.

The source review also covers delayed GET/PUT ordering, submitted-field revisions, secret clearing, and failed saves.
See `U14-settings-history-check.md`.

## Cancellation history remains open

The actual native delayed-stream tests passed for Responses and Chat Completions.
Ending the start and state requests did not cancel the background evaluation.
Explicit Stop cancelled the supplier stream and retained the actual cancellation outcome in history.

The cause of the user's reported history error remains unknown.
The pending user clarification asks whether that run received an explicit Stop.
No cancellation product code or old history was changed.

## Delivery

The final static Linux/amd64 executable is `dist/build/sub2api-codex-eval`, with 125595808 bytes.
The archive is `dist/build/sub2api-codex-eval.tar.gz`, with 38700266 bytes and seven runtime/license entries.
Independent review, focused UI/native checks, native frontend build, static build, and packaging passed.
No additional Docker container, production deployment, paid supplier call, or commit was used.
