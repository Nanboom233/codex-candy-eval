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
AccountTestService.PersistCodexInspectionVerdict(context.Context, int64, codexeval.InspectionVerdict, codexeval.InspectionSettings) error
```

The service uses `SettingRepository.GetMultiple` and `SetMultiple` in the existing PostgreSQL `settings` table.
Each persisted key uses `codex_eval_inspection_` plus its JSON field name.
Each value contains the JSON representation of that field.

## 3. Contracts

### Settings and page

| Fields | Type and default |
| --- | --- |
| `enabled`, `adjust_priority`, `control_models` | Boolean, false |
| `degraded_remove_models` | String array, `gpt-6.1-sol`, `gpt-6-sol`, `gpt-5.6-sol`, `gpt-6-astra` |
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

The fourth page tab is `设置`.
Use native settings cards, visible labels, units, and inline validation.
Keep explanatory paragraphs out of the form.
Use a visible model dropdown with raw suggestions and a custom literal-model path.
Use current-group model suggestions from the manual catalog without another all-account model discovery.
Preserve newer unsaved model edits during discovery.
Include/exclude account controls, fields, persistence keys, and filtering logic are removed.
Remove the independent Free switch. Use the Candy plan filter and native availability and busy rules.
Hide normal and degraded priority fields while automatic adjustment is off.
Hide automatic inspection configuration fields while scheduling is off.
Remove visible switch captions. Keep accessible names on header switches.
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

Place `模型控制` after `优先级` as a sibling settings card with an independent switch.
Show `降级移除模型` only while model control is on. Keep its values while collapsed.
Copy the native `ModelWhitelistSelector.vue` layout into the embedded page's HTML and JavaScript.
Keep its two-column gray tags, icons, bottom separator, absolute dropdown, and external custom-model input row.
The custom row uses `自定义模型名称` and the native `填入` button.
Omit counts, small explanatory text, tooltips, copy actions, and quick actions.
Keep controls keyboard accessible without changing the native visible layout.
Reuse the current-group catalog and existing automatic save owner. Add no discovery request.
Retain selected IDs absent from the catalog. An empty catalog still permits custom entry.
Trim model names and remove duplicates while preserving order.
Absent removal settings select defaults. An explicit empty list stays empty.
Use the existing model-name rules. Reject wildcard entries and the reserved name `blacklist`.

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
Round state stays in memory. Candy records use the existing history file. Account verdicts use native account storage.
Keep backend errors, authentication handling, and polling recovery after removing the page error prompts.

### Account participation

Store Boolean `codex_eval_inspection_enabled` in the native account's `extra` object. An absent value means enabled.
Save only that key through `AdminService.UpdateAccountExtra`, using the existing atomic JSONB merge.
Preserve the key during unrelated full account edits. Existing DTO, bulk-update, copy, and export/import paths need no changes.
One native account ID has one value across groups, reloads, and restarts. Shadows keep their own value.
Keep the global timer switch and native account `schedulable` independent from account inspection participation.
Apply participation to scheduled and immediate coordinator rounds. Manual Candy tests retain their entry points and fixed counts.
An in-flight inspection request can finish after opt-out. Check participation before each later send and before a priority update.
This includes retry sends after backoff. Run the callback outside the global evaluation mutex.
Use an inspection-only pre-send callback in the Candy worker. Do not hold the runtime mutex during account reads.
Intentional opt-out retains completed results and ends as `unchanged`, without a synthetic evaluation error.
Do not append a synthetic history result or increment progress for an intentional pre-send stop.
Do not cancel unrelated accounts when one account opts out.
An incomplete confirmation batch cannot produce a recommendation, priority write, or notification.
Keep per-account pending UI state during refresh. A failed save restores the saved value and shows a row-local error.

### Account verdict

Store `normal` or `degraded` in native `extra.codex_eval_inspection_verdict`.
An account without a saved verdict starts as normal, regardless of its actual priority.
Read the saved verdict for both scheduled and immediate rounds. Manual tests do not change it.
Do not infer verdicts from old priorities or migrate priority-based state. This is an accepted breaking change.
Use `InspectionAccount.Verdict` and `InspectionHooks.PersistVerdict(context.Context, int64, InspectionVerdict, InspectionSettings) error`.
The callback receives the round settings and delegates to `AccountTestService.PersistCodexInspectionVerdict`.
When only the verdict changes, save its key through native `AccountRepository.UpdateExtra`.
With model changes, save credentials, pending removals, and verdict through one native `AccountRepository.Update`.
Preserve inspection metadata during unrelated full account edits.
After each account completes detection and existing confirmation rules establish its state, apply model control and optional priority adjustment.
Include stable states. Apply each ready account without waiting for unrelated confirmations.
Keep the existing confirmation-group order and recovery wait. Apply controls only after detection and required confirmations establish that account's state.
Save changed model data and verdict before an optional priority update. Skip persistence when all stored data already matches.
A failed verdict save prevents the priority write. A failed priority write retains the saved verdict and existing backend error.

### Model removal and recovery

After detection establishes a degraded state, apply the round's removal list only when `control_models` is on.
Remove a stored `credentials.model_mapping` entry only when its key and value both equal the configured name.
Skip absent entries and automatic-passthrough accounts. Keep unrestricted accounts unrestricted.
Preserve aliases, wildcard entries, compact mappings, other credentials, and native `schedulable`.
Merge actual removals into `extra.codex_eval_inspection_removed_models` without duplicate names.
If removal empties the complete mapping, insert `{"blacklist":"blacklist"}`.
The nonempty mapping blocks ordinary real-model names. Literal `blacklist` still matches because it is an ordinary exact key.
The direct evaluation sender bypasses gateway model admission and remains available for recovery inspection.

After detection establishes a normal state, restore names from the saved removal list.
Use that list even after global configuration changes or switch-off.
Restore absent names as identity entries. Keep any currently present administrator replacement.
Remove the reserved identity placeholder and clear the pending list in the same account update.
Use copies of mutable account maps so a failed write does not change the loaded account object.
Stable states also apply model control. Enabling control for an already degraded account takes effect after its next valid detection.
Use the round settings snapshot. A later round uses changed settings.
Disabling control or shortening the list preserves pending removals until a normal result restores them.
Model control adds no scheduler state, background synchronization, write retries, or result field.

### Account tests and priority changes

Read current members of the saved native group and apply the shared Candy type and plan filters.
Use the native parent-plan enrichment for shadow accounts during selection and the pre-update account reread.
An account's own plan takes precedence. Otherwise, use its parent's plan, as the Candy account API does.
Apply native availability and the existing account busy state. Do not replace a missing group with all accounts.
Pass the configured model literally through the existing native sender.

Each eligible account gets one first test.
One logical test can contain up to three request attempts. Only its final completed answer counts toward confirmation.
Valid wrong Candy text does not retry. Request failures do not become confirmation answers.
Keep the existing 30-minute phase coordinator wait, then cancel and drain active workers.
Do not extend that wait to guarantee three full ten-minute attempts.
After that test, reread the account verdict before selecting confirmations. Do not classify from actual priority or the account-list snapshot.
If the verdict is degraded and the answer remains wrong, stop after one test.
If the verdict is normal and the answer is correct, stop after one test.
Both rules apply with automatic priority adjustment on or off.
Keep the original result priority as historical data. Preserve the fresh actual-priority baseline for existing confirmation and update checks.
Only a possible normal/degraded transition starts up to `max_tests_per_account - 1` confirmations.
Recovery confirmation waits 240 seconds.
Require the full configured count of valid, complete, matching outcomes before changing the verdict.
Errors, cancellation, incomplete results, mixed outcomes, or intentional opt-out prevent a verdict transition.
Keep existing final priority and busy checks on the priority-update path. They do not add verdict-save conditions.

Lower Sub2API priority values are preferred.
Each completed detection with an established state can enter the optional priority-update path.
If `adjust_priority` is off, save the verdict without changing actual priority.
If enabled, use that state's configured target priority. A custom current priority does not block that write.
Recheck account availability, priority, and busy state before an update.
An already matching actual priority needs no write.
Stable states also use their configured target priority, including after adjustment becomes enabled.
Do not add automatic priority-write retries.
After actual priority updates succeed, send configured ntfy/Bark notifications.
Notification text must not claim three tests for a stable-state correction.
Notification failures remain separate and do not reverse successful updates.

## 4. Validation & Error Matrix

| Condition | Result |
| --- | --- |
| Missing administrator authentication | Native authentication failure |
| Negative group ID, or automatic inspection enabled with group ID 0 | Reject settings with 400 |
| Start without a selected group | Reject the start. Never select all accounts. |
| Account type is not `all`, `oauth`, or `apikey`, or plan is empty | Reject settings with 400 |
| Invalid model or test count outside 1–10 | Reject settings with 400 |
| Invalid removal model, wildcard, or reserved `blacklist` entry | Reject settings with 400 |
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
| Verdict save fails | Keep the previous verdict. Do not write priority or send a success notification. |
| Compound model update fails | Keep mapping, pending removals, and verdict unchanged. Do not write priority. |
| Priority update fails after a verdict save | Retain the saved verdict and backend error without automatic retry |
| Account participation is false | Skip inspection selection and later sends. Keep manual tests and normal traffic usable. |
| Account switch save fails | Keep the saved value and show the error beside that switch |
| Account participation field is absent in the PUT body | Reject the request with 400 |
| Notification fails after update | Keep the update and backend `notification_error`. Show no page notification error prompt. |

## 5. Good / Base / Bad Cases

Good: a normal-state account returns three valid wrong answers. Save degraded state, then change priority if adjustment is enabled.

Good: removal deletes the final mapping entry. Save the placeholder and actual removals, then restore them after confirmed recovery.

Base: a normal-state account returns one correct answer. Apply normal model control and priority without more confirmation.

Base: a degraded-state account returns one wrong answer. Apply degraded model control and priority without more confirmation.

Base: the stored model data, verdict, and priority match the target. No account write follows.

Bad: a wrong first answer has an incomplete confirmation. Do not change state or priority from that partial result.

Bad: recovery overwrites a model mapping that an administrator changed. Keep that current entry instead.

## 6. Tests Required

Use the existing build/check entry and focused inspection, runtime, admin, service, and UI tests.
Use one adjustment-off lifecycle for confirmed degradation, a stable wrong round, confirmed recovery, and a stable correct round.
Extend existing transition, mixed-answer, provider-error, cancellation, and opt-out cases with verdict assertions.
Extend those cases for stable-state controls, matching-write suppression, and application before unrelated accounts complete.
Check the new verdict-save and priority-write failure boundaries without repeating that lifecycle.
Reuse native metadata fixtures for the normal default, key-only persistence, and full-edit preservation.
Retain manual fixed counts for opted-out accounts. Reuse existing scheduling, settings, notification, and unchanged UI coverage.
Follow the active task's bounded check scope. Add no separate restart, copy, DTO, custom-priority, or browser matrices for backend-only changes.
Reuse valid evidence for unchanged inputs.
Extend settings and metadata fixtures for model-control fields and pending-removal preservation.
Use focused account-helper cases for partial/final removal, native admission, recovery, skipped entries, and failed compound writes.
Reuse the existing adjustment-off lifecycle and unchanged sender evidence. Add no separate restart or sender matrix.
Extend existing UI cases for the array field and widget. Compare the actual render with the native selector.

## 7. Wrong vs Correct

Wrong: apply inspection availability filters to every manual administrator evaluation.

Correct: filter automatic inspection through native account hooks and preserve manual evaluation admission.

Wrong: save notification credentials and settings in a new `inspection.json` file.

Correct: save settings through native `SettingService` and the existing PostgreSQL `settings` table.

Wrong: save an empty model mapping after deleting its final entry. Native empty mappings select default admission.

Correct: save the reserved identity placeholder and restore only recorded removals after confirmed recovery.
