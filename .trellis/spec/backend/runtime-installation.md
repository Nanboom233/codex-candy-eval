# Sub2API Administrator Evaluation

## 1. Scope / Trigger

Use this contract for the current v0.2.18 administrator evaluation integration.
Production source lives in [contextual patches](../../../patches/), applied in lexical filename order.
The [integration patch](../../../patches/0007-sub2api-integration.patch) adds native routes, initialization, sidebar, URL resolution, page bypass, and navigation translations.
The current implementation replaces the historical internal-key, activation, billing, maintenance and recovery designs.

## 2. Signatures

```go
func (s *AccountTestService) SendCodexEvaluation(ctx context.Context, id int64, body []byte, headers http.Header) (*http.Response, error)
func (s *AccountTestService) CodexEvalAccount(ctx context.Context, id int64) (string, error)
func (s *AccountTestService) CodexEvalModels(ctx context.Context, id int64) ([]openai.Model, error)

func Initialize(dataDir string, send Sender, lookup Resolver)
func Start(ctx context.Context, req StartRequest) (StartSummary, error)
func Cancel(runID string) bool
func State() ([]byte, error)
func Record(kind string, accountID int64, recordID string) ([]byte, int)
func Clear(scope string) ([]byte, int)
```

| Method | Path |
| --- | --- |
| POST | `/api/v1/admin/accounts/:id/codex-eval` |
| GET | `/api/v1/admin/accounts/:id/codex-eval/models` |
| GET | `/api/v1/admin/codex-eval/state` |
| POST | `/api/v1/admin/codex-eval/start` |
| POST | `/api/v1/admin/codex-eval/cancel` |
| GET | `/api/v1/admin/codex-eval/records/:type/:account_id/:record_id` |
| POST | `/api/v1/admin/codex-eval/clear-history` |
| GET | `/codex-eval` and `/codex-eval/` |
| Native administrator view | `/admin/codex-eval` |

## 3. Contracts

### Administrator sending

Native administrator middleware protects all evaluation APIs.
The account endpoint accepts a complete Responses JSON request with nonempty `model` and `stream:true`.
It returns the actual upstream status and stream.
Task operations use native `{code,message,data}` envelopes.

The endpoint and background runners call the same native sender directly.
AccountTestService owns credentials, URL validation, proxy, TLS, headers and plugin transport.
Evaluation passes the chosen model literally, bypassing account model mappings and OAuth normalization only in this sender.
The evaluation model catalog reuses native `FetchOpenAIModelsList` raw discovery and omits `projectAccountModelsBody`; the normal account catalog retains its configured projection.
Chat Completions accounts use the existing native request and response conversion functions.
Preserve observed usage and actual error text, including incomplete CC streams and response-copy failures.
Keep native snapshots and 401/429 handling.
The new sender omits broad successful-test recovery.
Ordinary `/test` and the user gateway retain their existing behavior.

Evaluation uses default identity `codex-tui/0.163.0 (Debian 13.0.0; x86_64) xterm-ghostty (codex-tui; 0.163.0)`.
Retain account-specific User-Agent overrides.
Use the current captured context and `history_initialization: new`.
Place fingerprint developer instructions in `generic.developer_instructions`.

### Tasks and history

`StartRequest` contains `type`, `targets`, `model`, `effort`, `runs`, `concurrency`, and optional `variant`.
Supported types are `candy`, `fingerprint`, and `pelican`.
Pelican accepts `variant: original` or `all-in-one`. Missing or empty values select `original`.
Reject unknown nonempty variants before account lookup or supplier calls. Other evaluation types ignore this field.
Fingerprint selects three distinct challenges. Per-account concurrency ranges from 1 through 3, with default 3.
Six global challenge slots include each challenge's retry waits.
Each target contains numeric `group_id` and `account_id`.
Groups describe page filters and result origin.
An account remains busy across evaluation types until its worker exits.
Start returns `started`, `busy` and `unavailable` arrays.
Cancel accepts `run_id` and interrupts the matching task context and active response.

State contains `active`, `candy`, `fingerprint`, `pelican`, and `storage_error`.
Three history maps persist at `<data-dir>/codex-eval/history.json` with root marker `schema_version: 2`.
The [persistence contract](./database-guidelines.md) defines old-format copies, retained maps, and migration write protection.
Keep the latest 20 records per account and type.
Save each result through a same-directory temporary file and rename.
Unreadable or null history blocks overwrites.
A failed clear restores the previous maps.
Run state, progress and cancellation remain in memory.
Record details support fingerprint and Pelican with exact record IDs.
Keep wrappers `{type, fingerprint}` and `{type, pelican}`.
Started entries expose `id`. Saved records retain native `Info.run_id`.
Candy results remain inline in state.

Pelican definitions in `pelican.go` supply exact upstream prompts and public page configuration through `web.go`.
Successful, failed, and cancelled records retain the admitted variant. Later toolbar changes affect only subsequent starts.
Prefer the selected language block without case sensitivity, then the first block whose trimmed content starts with `<`.
Accept bare markup only when no fenced block exists. An empty preferred block fails acceptance.
Limit extracted bytes to 256 KiB before the all-in-one SVG viewport wrapper. Keep original HTML unchanged.
Do not add XML validation or grade picture content and its Candy answer.
Use the existing shared retry runner, usage accounting, and native result statuses.

Fingerprint progress has `done`, `total: 3`, model, and concurrency.
Records retain native identity, detector revision, samples, usage, status, and optional attribution.
Summaries clear samples, full ranking, and diagnostics. The page accepts their null JSON values.
Details retain challenge instructions, prompts, answers, parsed counts, acceptance, and attempt counts.
After sequence acceptance, three valid answers produce probabilities. One or two produce partial rankings with null probabilities.
Sequence rejection produces `attribution.decision: not_model` without rankings or probabilities.
Parse the longest numeric run within 64,000 bytes. Keep values from 1 through 355.
Require `max(80, ceil(expected * 0.55))` values and truncate excess values to the expected count.
Sequence acceptance requires median lag-one autocorrelation at most -0.23 and median small-step share at most 0.05.
Explicit failed or cancelled status takes precedence over attribution.
Cancellation before three completed challenge events remains cancelled.
Comparison aliases do not change the literal sent model.
Use `gpt-6.1-sol` to `gpt-6-astra` and `deepseek-flash` to `deepseek-v4.1-flash-1017` for comparison only.

### Standalone page

Serve the original page and renderers through the standalone shell.
Read native same-origin `auth_token` for authenticated API calls.
Use native groups/accounts and evaluation-specific raw model suggestions directly.
When the saved tab is Settings, let `start()` own the initial settings read before refreshing shared groups.
Restoring the tab must not start a competing settings read that later invalidates the first group response.
Later user entries to Settings still refresh saved settings through the existing owner.
Manual account inventories must contain only selected-group members on every page.
Read `/api/v1/admin/accounts?platform=openai&group=<id>&lite=1&page=<n>&page_size=100` and its `group_ids` field.
Refresh account membership on ordinary polls and reject obsolete group responses before updating the catalog.
Keep the last valid same-group inventory usable during background refresh. Preserve unchanged row controls, selection, and expanded history.
Block starts before the initial successful read, after a group/session change, or after an actual inventory failure.
Keep an inventory error visible until recovery. Do not clear valid readiness merely because another request starts.
Current inventory determines rows and targets even when the subsequent evaluation-state request fails.
An evaluation-state failure must not restore accounts removed from the current inventory.
In all three evaluation lists, `credentialView` adds one `降级` tag after priority when `extra.codex_eval_inspection_verdict === "degraded"`.
Use the existing neutral `.tag` style for matching alignment and wrapping.
Normal or absent verdicts show no degraded tag, regardless of actual priority or automatic priority adjustment.
The existing account inventory refresh removes the tag after a saved recovery.
Reuse raw model suggestions while enabled account IDs stay unchanged.
Model discovery errors must not block a literal model when the account inventory is current.
Scheduled Candy uses the Candy page's persisted group and filters for account selection.
Inspection model, effort, counts, and flow remain independent. Manual fixed-count behavior stays unchanged.
Do not maintain a separate inspection account catalog or settings-page result list.
Manual batch targets must come from the displayed eligible selected-group accounts.
All three model controls use a visible dropdown with native raw suggestions.
Keep the direct custom-model path inside the dropdown interaction, including when the native catalog is empty.
Selecting `自定义模型…` reveals the literal model input inside the same field.
Retain a selected literal model when refreshed suggestions omit it.
Save valid preferences automatically on input, before blur.
Shared group and Candy account filters use native settings so the timer can run after the browser closes.
Candy and Pelican retain their browser preferences and manual entry points.
Pelican places version after effort, defaults to original, and restores a valid browser preference across group changes.
Version changes immediately update the visible and copied prompt. Merge only that preference during asynchronous model restoration.
Latest results, history, and details show each record's saved version. Toolbar changes do not filter or relabel history.
Both variants share the latest-20 limit. Keep 4:3 frozen previews, animated details, sandbox isolation, and stale-detail guards.
At narrow widths, pair group/model and effort/version. Put the test action on its own full-width row.
Fingerprint retains model and concurrency preferences. Do not read retired mode or ModelTrace preferences.
A saved ModelTrace tab falls back to Candy.
Render `not_model` as `不像语言模型回答` and null probability as `—`.
Account enable uses `POST /api/v1/admin/accounts/:id/schedulable` with `{schedulable:boolean}`.
Derive disabled/switch-off from `schedulable:false`; display account status and future rate-limit, overload and temporary restrictions separately.
Retain all account rows, but fetch model suggestions only for `schedulable !== false` accounts.
Keep enabled-account model errors visible across polling; mark the catalog loaded only after discovery succeeds.
The existing enable toggle forces catalog refresh, clearing disabled-account notices and restoring genuine errors when reenabled.
Missing token or 401 stops polling and requests native login.
403 reports missing administrator permission.

Native administrator sidebar provides `账户评估` beside account management with the existing beaker icon.
The authenticated administrator route `/admin/codex-eval` uses native `AppLayout` and the evaluation iframe.
Pass relative `/codex-eval` to `buildEmbeddedUrl`, which resolves it against the full browser `window.location.origin` including nonstandard port; retain theme/language/embedded parameters and existing absolute URLs.
Use the same-origin session in the frame rather than guessing its external origin from proxy headers.
Page-only headers permit same-origin framing and original inline scripts/styles.
Native parent/global headers remain unchanged. Rebuild the patched native frontend before embedding it in the server.
Retain original layout, controls, Markdown/math and Pelican sandbox behavior.
Remove assistant-added explanatory small print and tooltips.
Keep field labels, units, actual errors, result data, progress, and concise action feedback.
Use native light/dark theme variables.

## 4. Validation & Error Matrix

| Condition | Outcome |
| --- | --- |
| Missing administrator authorization | Native authentication rejects the API call |
| Invalid account/model/stream input | Actual validation error, no alternative account |
| Configured model mapping or empty model catalog | Send the administrator's chosen model literally |
| Raw model exists but configured account catalog filters it | Show raw evaluation suggestion; normal catalog remains filtered |
| Browser uses a nonstandard port | Frame keeps the full browser origin and satisfies parent `frame-src 'self'` |
| Account `schedulable:false` | Show disabled and switch off; use native schedulable POST to toggle |
| Disabled account would fail model lookup | Skip its discovery request and error notice; retain its row/control |
| Enabled account model lookup fails | Show its actual error, including after subsequent polling |
| Native entry restores Settings before a test-tab switch | Load shared groups after initial settings, then show the selected group's accounts |
| Background inventory refresh has a valid same-group snapshot | Keep unchanged controls usable without a readiness flash |
| Manual account inventory refresh fails | Block row and batch starts until a successful inventory refresh |
| Inventory refresh succeeds but evaluation state fails | Keep the state error. Remove moved accounts from rows and targets. |
| A saved preference changes before blur | A page reload restores the input value through the existing preference store |
| Account already running | Target appears in `busy` |
| Lookup or account protocol unavailable | Target appears in `unavailable` with the actual reason |
| HTTP/SSE/CC failure or premature stream EOF | Actual failure and observed usage, no fabricated completion |
| Task cancelled | Interrupt active request and prevent subsequent calls |
| Unknown nonempty Pelican variant | Return 400 before account lookup or supplier calls |
| Pelican selected markup absent or oversized | Persist `failed` and the admitted variant consistently in summaries and details |
| History unreadable/null or save fails | Visible storage error; preserve the existing file |
| History clear cannot save | Restore previous in-memory maps |

## 5. Good / Base / Bad Cases

Good: an administrator tests one account through the shared native sender and sees its result and history.
Base: opening the shell without a session requests native login before fetching evaluation data.
Bad: treating a CC stream failure as successful EOF hides the upstream error.

## 6. Tests Required

Use existing relevant tests for changed code and add focused assertions only for uncovered behavior.
Pelican extraction tests own language preference, page preparation, and the byte boundary.
Runtime tests own variant admission, saved versions, and migration. Reuse shared retry tests without a per-variant matrix.
The existing UI harness owns preferences, prompt/copy updates, payloads, and record-derived labels.
Use the actual browser for toolbar layout, previews, and detail animation. Do not repeat the DOM assertion list.
The focused UI checks use two groups with exclusive and shared accounts, membership moves, pagination, and obsolete responses.
Check submitted targets after inventory failures and after successful inventory refresh followed by evaluation-state failure.
Check manual preference reload before blur through actual input events.
Use the existing UI initialization harness for Settings-first entry followed by an evaluation-tab switch.
Check the selected group and account rows without a manual refresh. Exercise that flow in an actual browser.
The native delayed-stream lifetime test covers start/state request completion and explicit stop through both supported response protocols.
Normal HTTP request completion preserves background work. Explicit stop cancels supplier work and retains the actual history outcome.
[Runtime tests](../../../tests/test_codexeval_runtime_test.go) cover admission, cancellation, history and stream outcomes.
[Sender tests](../../../tests/test_service_account_codex_eval_service_test.go) cover native transport and CC failures with observed usage.
Native suites live directly under flat `tests/`. Package-context annotations preserve private access in the temporary upstream packages.
The [independent review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/fresh-native-check.md) records actual results and limits.
The [browser/model review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/native-browser-model-check.md) records actual Chrome sidebar/icon/iframe/CSP and raw model discovery through UI start to literal upstream sending against the final artifact.
The [installer/discovery review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/installer-disabled-model-check.md) records mixed disabled/enabled lookup, toggle refresh and retained direct model input checks.
Reuse completed pristine and local integration evidence unless a change or concrete finding requires another check.
No race, coverage, full-suite, architecture matrix or source-resource parity gate is required.

## 7. Wrong vs Correct

Wrong: create an internal user/key, loopback request or alternate gateway to send administrator tests.
Correct: call the existing account-test sender with the exact selected account.

Wrong: persist task recovery or alias a run ID to a history record ID.
Correct: keep tasks in memory and use exact record IDs for saved details.

Wrong: generate an iframe URL from a hostname-only forwarded header or accept Node-VM execution as parent CSP verification.
Correct: resolve the frame URL in the browser with its complete origin and verify the native navigation/frame in an actual browser at a nonstandard port.

Wrong: remove model mapping only in sending while the page still fetches the configured model projection.
Correct: use raw evaluation discovery for suggestions and pass the chosen model literally through the UI start and sender.
