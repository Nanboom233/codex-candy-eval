# U15 Browser and Background Inspection Check

The final v0.2.17 candidate passed these checks in `sub2api-codex-eval-test`.
The existing container ID starts with `bd6c9583c3e6`. Its origin is `http://127.0.0.1:60198`.
The fixture helper replaced its binary and restarted this container. It created no container.
The Linux runtime reported Sub2API 0.2.17.

## Automatic persistence and page contents

The settings tab has no save button and no account list.
Typing a notification username saved it through the native settings API without a save action.
Both secret inputs remained empty.
A fresh page load restored the username, selected group, account type, and enabled-only filter.
Manual Candy retained its separate two-run preference.

The actual group lists were:

| Selected group | Visible accounts |
| --- | --- |
| U13 Group A, ID 3 | u13-a-only, ID 6. u13-shared, ID 8. |
| U13 Group B, ID 4 | u13-b-only, ID 7. u13-shared, ID 8. |

The native settings API saved group 3, type `apikey`, plan `all`, and `enabled_only:true`.
Manual `Test all` ran twice per account. Each matching account history grew from two records to four.
No other-group account appeared in that batch.

## Timer without an open application page

The browser enabled inspection with waits of one minute and a maximum of two tests per account.
Its independent model remained `inspection-literal`, with effort `low`.
Editing a notification field retained the same pending timestamp: `2026-10-17T10:12:35.698261217Z`.
The browser then navigated to `about:blank`, which removed the application page and its polling.

The backend completed the scheduled round:

| Field | Observed value |
| --- | --- |
| Started | 2026-10-17T10:12:35.706649967Z |
| Finished | 2026-10-17T10:12:50.735955669Z |
| Phase | completed |
| Account progress | 2 / 2 |
| Account IDs | 6 and 8 |
| Tests per account | 1 |
| Outcome per account | unchanged |
| Runtime error | empty |

The correct first answer ended each account's independent inspection flow after one test.
Group B's exclusive account, ID 7, was absent from the round.
Reopening Candy showed both inspection summaries and the existing account history.
Each selected account now had five Candy records.

## Render and restored fixture

The desktop render shows settings cards, collapsed optional fields, and model dropdowns.
Candy shows inspection summaries below the matching account rows and retains the history control.
The narrow settings render uses the existing native menu and fits the iframe width.
The iframe measured 396 pixels wide. Its body measured 381 pixels, without horizontal overflow.

Screenshots:

- `output/playwright/U15-settings-desktop.png`
- `output/playwright/U15-candy-desktop.png`
- `output/playwright/U15-settings-narrow.png`

The fixture now has automatic inspection off, waits 15/35 minutes, and maximum tests 3.
Its notification username is empty. Shared filters are `all`, with enabled-only off.
The final browser read confirmed the restored settings.

## Delivery and remaining issue

The binary is `dist/build/sub2api-codex-eval`, with 125599904 bytes.
The archive is `dist/build/sub2api-codex-eval.tar.gz`, with 38701721 bytes and seven runtime/license entries.
Source review and focused checks are in `U15-shared-Candy-schedule-check.md`.
The build log is `/tmp/sub2api-U15-build.log`.

The earlier user-reported history cancellation still lacks its triggering operation.
Focused native stream tests cover background completion and explicit Stop cancellation.
No cancellation result or historical record was rewritten.
N150/iStoreOS execution and production deployment remain unperformed.
