# U13 Account List Browser Check

The final v0.2.17 candidate passed actual browser checks in the retained `sub2api-codex-eval-test` container.
Its origin was `http://127.0.0.1:60198`.
The fixture helper reused the container and loaded `dist/build/sub2api-codex-eval`.
No additional container or paid supplier was used.

## Native groups and manual lists

The fixture used actual native groups and account memberships.

| Group | ID | Members |
| --- | ---: | --- |
| U13 Group A | 3 | `u13-a-only` (6), `u13-shared` (8) |
| U13 Group B | 4 | `u13-b-only` (7), `u13-shared` (8) |

Candy, fingerprint, ModelTrace, and Pelican each showed only the selected group's two members.
The shared account appeared in both groups.
The other group's exclusive account did not appear.

A native account update moved account 6 from A to B while the browser remained on A.
Ordinary polling removed account 6 without a forced refresh.
Each of the four lists then showed only account 8.
Each actual batch request body contained only `{group_id:3,account_id:8}`.
The browser intercepted those starts to avoid unnecessary algorithm runs.
The native API restored account 6 to A after the check.

## Final candidate state failure

The final candidate repeated the move while the browser returned HTTP 503 for the evaluation-state endpoint.
Native account inventory requests continued through the real API.
Ordinary polling still removed account 6 from all four lists.
Each actual batch request again contained only `{group_id:3,account_id:8}`.
The page showed `读取测试记录失败：u13 synthetic state outage`.
The shared account retained its two prior Candy marks.
After removing the interception and restoring membership, ordinary polling restored the A-only row.

The independent focused UI checks cover inventory-fetch failure, row starts, recovery, pagination, obsolete responses, and inspection model-cache changes.
See `U13-account-list-check.md`.

## Shared inspection display

A manual Candy request completed through the actual local synthetic supplier.
An inspection then completed through the page's `立即巡检` control.
The final rebuilt candidate also completed an inspection after its last container reload.
It displayed eight native accounts, including the disabled account.
Six eligible accounts had inspection summaries.
The disabled account showed `已停用` and had no test record.

Inspection reused native account cards, priority badges, Candy verdicts, metrics, marks, and expandable history.
The final shared account showed three history cards after expansion.
Inspection rows had no manual account-start controls.
The automatic inspection and priority configuration bodies remained collapsed while their switches were off.
Manual inspection controls remained available.

The main session inspected the actual desktop and narrow screenshots.
At a 430-pixel browser width, the iframe body width and viewport width were both 381 pixels.
The layout had no horizontal overflow.

Screenshots:

- `output/playwright/sub2api-U13-final-desktop.png`
- `output/playwright/sub2api-U13-final-narrow.png`

## Delivery and limits

The final binary has 125595808 bytes.
The archive has 38700060 bytes and seven runtime/license entries.
The build passed the native frontend locale tests, TypeScript checks, Vite build, and static Linux/amd64 compilation.
The independent review has no unresolved findings.

Injected HTTP 503 failures produced expected browser network errors.
The previously recorded native parent-page cached-HTML/CSP nonce issue remains outside this change.
N150/iStoreOS and production HTTPS execution were not performed.
No commit or production deployment ran.
