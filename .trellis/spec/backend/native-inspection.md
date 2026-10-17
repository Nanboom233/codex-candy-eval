# Native Candy Inspection

## 1. Scope / Trigger

Use this contract for automatic Candy inspection inside the Sub2API process.
The coordinator reuses the native Candy runtime, account services, history, and cancellation.
Inspection uses the Candy page's saved group and account filters with current native membership.
Its model, effort, count, confirmations, and recovery flow remain independent from manual evaluations.
Manual row, selected-account, and fixed-count batch tests retain their behavior.

Owners:

- `patches/0005-inspection.patch`: coordinator, settings type, notifications, native settings service, and administrator handlers.
- `patches/0001-evaluation-core.patch`: inspection-only admission before each Candy request.
- `patches/0006-evaluation-ui.patch`: settings form, account participation switches, and runtime controls.
- `patches/0007-sub2api-integration.patch`: routes, account metadata preservation, initialization, and shutdown integration.
- `tests/backend/`: corresponding package tests and test-only data.

## 2. Signatures

All routes use native administrator authentication and the native response envelope.

| Method | Path | Result |
| --- | --- | --- |
| GET | `/api/v1/admin/settings/codex-eval-inspection` | Public settings and secret configured flags |
| PUT | `/api/v1/admin/settings/codex-eval-inspection` | Save settings and return public settings |
| GET | `/api/v1/admin/codex-eval/inspection/state` | In-memory inspection state |
| POST | `/api/v1/admin/codex-eval/inspection/start` | Start one round or return 409 |
| POST | `/api/v1/admin/codex-eval/inspection/stop` | Cancel the current round |
| PUT | `/api/v1/admin/accounts/:id/codex-eval-inspection` | Save one account's inspection participation |

The account endpoint accepts `{"enabled":true}` or `{"enabled":false}` and returns `{account_id, enabled}` in the native envelope.
Bind `enabled` as `*bool`. An absent field is invalid, while explicit false is valid.

Native service signatures:

```go
GetCodexInspectionSettings(context.Context) (codexeval.InspectionSettings, error)
UpdateCodexInspectionSettings(context.Context, codexeval.InspectionSettings) (codexeval.InspectionSettings, error)
```

The service uses `SettingRepository.GetMultiple` and `SetMultiple` in the existing PostgreSQL `settings` table.
Each persisted key uses `codex_eval_inspection_` plus its JSON field name.
Each value contains the JSON representation of that field.

## 3. Contracts

### Settings and page

| Fields | Type and default |
| --- | --- |
| `enabled`, `adjust_priority` | Boolean, false |
| `model`, `effort` | String, `gpt-6.1-sol` and `low` |
| `max_tests_per_account` | Integer, 3 |
| `min_wait_minutes`, `max_wait_minutes` | Integer, 15 and 35 |
| `group_id` | Integer, 0 until the Candy page selects a group |
| `account_type`, `account_plan` | String, `all` for both. Shared Candy filters. |
| `enabled_only` | Boolean, false. Shared Candy enabled-account filter. |
| `normal_priority`, `degraded_priority` | Integer, 10 and 99. Normal priority must be at least 1. |
| `ntfy_server` | String, `https://ntfy.sh` |
| `ntfy_topic`, `ntfy_username`, `ntfy_password`, `bark_url` | String, empty |

Public responses clear `ntfy_password` and `bark_url`.
They return `ntfy_password_configured` and `bark_url_configured` booleans.
These flags are not persisted.
Empty secret updates retain the saved values.
A failed save leaves active settings unchanged.

The fifth page tab is `设置`.
Use native settings cards, visible labels, units, and inline validation.
Keep explanatory paragraphs out of the form.
Use a visible model dropdown with raw suggestions and a custom literal-model path.
Use current-group model suggestions from the manual catalog without another all-account model discovery.
Preserve newer unsaved model edits during discovery.
Include/exclude account controls, fields, persistence keys, and filtering logic are removed.
Remove the independent Free switch. Use the Candy plan filter and native availability and busy rules.
Hide normal and degraded priority fields while automatic adjustment is off.
Hide automatic inspection configuration fields while scheduling is off.
Remove visible switch captions. Keep accessible names on both header switches.
Keep collapsed configuration values during automatic saves.
The settings page has no save button or account list.
Show account progress, results, and existing Candy history only on the Candy tab.
Keep inspection runtime controls available when automatic scheduling is off.
Keep immediate-inspection and stop buttons on the settings tab, outside its collapsed configuration.
Automatic saves show no saving, saved, or unsaved-change captions. Preserve actual validation and save failures.
Keep evaluation errors in the original Candy result records and history.
Remove inspection runtime and notification error prompts, including the account-header summary. Do not move them to Settings.
Keep settings validation/save errors in Settings and account-switch save errors beside the affected switch.
Do not show a completed-phase caption or stale completed progress. Hide empty status containers.
Read an existing nonpositive stored normal priority as default 10 before validating loaded settings.
New writes below 1 fail validation. The read does not write the database.
Hide empty idle status and runtime failure captions.
Keep running progress and Candy history/results.
Keep `下次巡检：` and its scheduled time in the Candy account header with `账户` and `仅启用`.
Add only the account-level `启用自动巡检` column to the Candy list. Keep the native account enabled column separate.
Do not add a separate inspection card or result columns.
Keep other evaluation headers unchanged. Wrap the Candy header on narrow screens without horizontal overflow.
Align both Candy row switches on the same baseline at narrow widths, regardless of history length.
Do not render notification configured-status text. Keep secret redaction and blank-secret preservation.
Catalog and save responses must preserve newer unsaved form edits.
Read native settings on each settings-tab entry and retain fields with unsaved edits.
Valid input edits save automatically after a short debounce. Serialize writes through one owner.
Invalid input stays visible without replacing saved valid settings. Show actual validation and write errors.
A save invalidates older reads. A successful save clears only submitted fields that have not changed again.
Apply returned public values to those fields, including clearing submitted secret inputs.
Preserve edits made during the save, including when another field changed.
Evaluation refresh and clear controls stay on their evaluation tabs.

### Scheduling and state

Load saved settings at process start.
Complete the initial settings read before exposing the coordinator to HTTP requests.
Validate the loaded settings before configuring or scheduling the coordinator.
If reading fails, retain the backend error and reject starts until a successful settings save configures the coordinator.
An invalid loaded priority pair follows the same blocked-start path and creates no timer.
Keep its explicit values available through the settings API for correction.
Only a successful repository read can select defaults for absent settings keys.
Enabled inspection draws an inclusive random delay between the configured minute bounds.
Equal bounds give a fixed delay.
After a round ends, draw a fresh delay.
Browser refresh does not redraw it.
Unrelated settings saves retain the pending timestamp.
Changes to the enabled switch or wait bounds draw a new delay.
An immediate run replaces the pending slot. A running round blocks another start.
Disabling automatic inspection leaves the current round running.
Stop cancels its requests and recovery wait.
Shutdown stops the coordinator and waits for its workers.
Restart reloads configuration and schedules a fresh wait without restoring an interrupted round.

State fields are `enabled`, `running`, `cancelling`, `phase`, `next_run_at`, `started_at`, `finished_at`, `total`, `done`, `error`, and `results`.
Phases are `idle`, `testing`, `recovery_wait`, `completed`, `cancelled`, and `failed`.
Results contain account identity, original priority, optional target priority, run IDs, test count, outcome, and separate evaluation/notification errors.
State stays in memory. Candy records use the existing history file.
Keep backend errors, authentication handling, and polling recovery after removing the page error prompts.

### Account participation

Store Boolean `codex_eval_inspection_enabled` in the native account's `extra` object. An absent value means enabled.
Save only that key through `AdminService.UpdateAccountExtra`, using the existing atomic JSONB merge.
Preserve the key during unrelated full account edits. Existing DTO, bulk-update, copy, and export/import paths need no changes.
One native account ID has one value across groups, reloads, and restarts. Shadows keep their own value.
Keep the global timer switch and native account `schedulable` independent from account inspection participation.
Apply participation to scheduled and immediate coordinator rounds. Manual Candy tests retain their entry points and fixed counts.
An in-flight inspection request can finish after opt-out. Check participation before each later send and before a priority update.
Use an inspection-only pre-send callback in the Candy worker. Do not hold the runtime mutex during account reads.
Intentional opt-out retains completed results and ends as `unchanged`, without a synthetic evaluation error.
An incomplete confirmation batch cannot produce a recommendation, priority write, or notification.
Keep per-account pending UI state during refresh. A failed save restores the saved value and shows a row-local error.

### Account tests and priority changes

Read current members of the saved native group and apply the shared Candy type and plan filters.
Use the native parent-plan enrichment for shadow accounts during selection and the pre-update account reread.
An account's own plan takes precedence. Otherwise, use its parent's plan, as the Candy account API does.
Apply native availability and the existing account busy state. Do not replace a missing group with all accounts.
Pass the configured model literally through the existing native sender.

Each eligible account gets one first test.
After that test, reread current priority before selecting confirmations. Do not classify from the account-list snapshot.
If the current priority is degraded and the answer remains wrong, stop after one test and keep priority unchanged.
This rule applies with automatic priority adjustment on or off. Send no priority-change notification for that account.
Keep the original result priority as historical data. Use the fresh classified priority as the confirmation/update baseline.
Only a possible normal/degraded transition starts up to `max_tests_per_account - 1` confirmations.
Recovery confirmation waits 240 seconds.
Require the full configured count of valid, complete, matching outcomes before recommending or changing priority.
Errors, cancellation, incomplete results, and mixed outcomes prevent changes.

Lower Sub2API priority values are preferred.
Only accounts at the configured normal or degraded priority can change.
Recheck account availability, priority, and busy state before an update.
With automatic changes off, show recommendations without updating accounts.
After successful updates, send configured ntfy/Bark notifications.
Notification failures remain separate and do not reverse successful updates.

## 4. Validation & Error Matrix

| Condition | Result |
| --- | --- |
| Missing administrator authentication | Native authentication failure |
| Negative group ID, or automatic inspection enabled with group ID 0 | Reject settings with 400 |
| Start without a selected group | Reject the start. Never select all accounts. |
| Account type is not `all`, `oauth`, or `apikey`, or plan is empty | Reject settings with 400 |
| Invalid model or test count outside 1–10 | Reject settings with 400 |
| Wait below one minute, reversed bounds, or unrepresentable duration | Reject settings with 400 |
| Normal priority is not lower than degraded priority | Reject settings with 400 |
| Normal priority below 1 in an update | Reject settings with 400 |
| Notification address lacks an HTTP/HTTPS scheme or host | Reject settings with 400 |
| Settings repository read/write fails | Return 500. Keep active settings unchanged after a failed write. |
| Settings change in another session | The next settings-tab entry reads the current saved values |
| A settings read returns after a newer save | Ignore the obsolete read |
| A field changes during a save | Keep that newer edit and clear only unchanged submitted fields |
| Initial settings read fails | Keep the runtime error and reject start with 409 until a successful settings save |
| Loaded settings fail validation | Keep the runtime error, create no schedule, and permit correction through valid automatic persistence |
| Another round runs | Immediate start returns 409 |
| Account is outside the saved group/filter, disabled, unavailable, or busy | Skip automatic selection |
| Account state or priority changes before update | Retain the backend account error without changing priority |
| Account participation is false | Skip inspection selection and later sends. Keep manual tests and normal traffic usable. |
| Account switch save fails | Keep the saved value and show the error beside that switch |
| Account participation field is absent in the PUT body | Reject the request with 400 |
| Notification fails after update | Keep the update and backend `notification_error`. Show no page notification error prompt. |

## 5. Good / Base / Bad Cases

Good: a normal-priority account returns three valid wrong answers. Enabled adjustment changes its priority from 10 to 99.

Base: a normal-priority account returns one correct answer. No confirmation or priority change follows.

Base: an already degraded account returns one wrong answer. The round records one test and leaves its priority unchanged.

Bad: a wrong first answer has an incomplete confirmation. Do not change priority from that partial result.

## 6. Tests Required

Use the existing build/check entry and focused inspection, runtime, admin, service, and UI tests.
Assert one call/result for stable and stale-list degraded-wrong cases, with adjustment on and off.
Cover participation defaults, Boolean false saves, metadata preservation, initial opt-out, and opt-out between confirmation sends.
Retain manual fixed counts for opted-out accounts. Reuse existing confirmation, cancellation, scheduling, settings, and notification coverage.
Use the UI harness for independent switches, pending saves, local errors, removed prompts, and the retained next inspection time.
Use one browser pass for the added column, switch interaction, reload, and narrow layout in the retained native fixture.
Reuse valid evidence for unchanged inputs. Add no separate restart, copy, DTO, custom-priority, or repeated-round matrices.

## 7. Wrong vs Correct

Wrong: apply inspection availability filters to every manual administrator evaluation.

Correct: filter automatic inspection through native account hooks and preserve manual evaluation admission.

Wrong: save notification credentials and settings in a new `inspection.json` file.

Correct: save settings through native `SettingService` and the existing PostgreSQL `settings` table.
