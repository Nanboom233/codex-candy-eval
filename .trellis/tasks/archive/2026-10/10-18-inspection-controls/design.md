# Stable Refresh and Shared Inspection State

## Change Boundary

Keep this work in one task. All four requirements share account rendering and inspection state.
One implementation owner must maintain the UI, runtime contract, tests, and production patch hunks together.

Production owners are patches `0001`, `0005`, `0006`, and `0007` in the existing series.
Keep regression tests under `tests/backend/` and `tests/ui/`.

| Applied source | Required change |
| --- | --- |
| `codexeval/web/app.js`, `catalog.js` | Separate request activity from valid same-group inventory. |
| `codexeval/inspection.go` | Use fresh priority for confirmation selection. Add account participation and per-request inspection checks. |
| `codexeval/runtime.go`, `candy.go` | Carry an inspection-only pre-send check without changing manual entry points. |
| `handler/admin/codex_inspection_handler.go` | Project account participation and save it through a focused account endpoint. |
| Native `service/admin_account.go` | Preserve the owned account metadata key during full account edits. |
| Native `server/routes/admin.go` | Register the account-level administrator update route. |
| `codexeval/web/components.js`, `candy.js`, `app.js`, `style.css` | Add the Candy-only account switch column, events, pending state, and narrow-screen labels. |
| `codexeval/web/inspection.js`, `components.js` | Remove inspection runtime and notification error prompts. Keep the next inspection timestamp. |

Paths in the table are relative to the applied upstream `backend/internal/` directory.

## Background Refresh

Separate an active inventory request from the validity of the last successful same-group inventory.
Keep unchanged rendered controls stable while that request runs.
Commit refreshed account membership before rendering new evaluation results.
Invalidate admission after a real account-read failure, initial load, group change, or session reset.
Retain current generation and cancellation checks so obsolete responses cannot replace a newer group.
Keep keyed row updates and existing model discovery caching.

The delayed DOM probe confirmed that the same row survives while its button changes from enabled to disabled and back.
Repair both readiness resets in `load` and `refreshCatalog`. Keep actual inventory failures separate from model discovery failures.

## Inspection Counts

The user requires one test when a degraded account still answers incorrectly.
Apply this rule to each inspection round, whether automatic priority adjustment is on or off.
Keep the account eligible for the initial test in later rounds.

| Account priority | First answer | Required behavior |
| --- | --- | --- |
| Normal | Correct | End after one test. Keep the priority. |
| Normal | Wrong | Run existing downgrade confirmations up to the configured total. |
| Degraded | Wrong | End after one test. Keep the priority. |
| Degraded | Correct | Wait for the existing recovery delay, then run recovery confirmations up to the configured total. |

The unchanged coordinator passed a stable degraded-wrong probe with one actual sender call.
A normal list snapshot followed by a degraded current-account read produced three calls in both adjustment modes.
`run` classifies confirmations from the list snapshot. `testPhase` rereads eligibility but discards the current priority.
The final pre-update priority check runs after the extra calls and cannot prevent them.

After the initial test, reread the account before selecting a confirmation group.
Use that fresh priority for the transition decision. Recheck it before confirmation admission and before a priority write.
Keep the initial result priority as historical data. Carry the fresh classified priority as the transition baseline.
If the current priority is degraded and the initial answer is wrong, finish with one test and no recommendation.
Keep settings stable within the existing round snapshot and preserve exact configured priority comparisons.

Keep run identity, cancellation, history, and manual fixed counts unchanged.
Do not add cumulative-progress fields or change the progress display for this requirement.

## Automatic Inspection Switch

Add an `启用自动巡检` column with a switch in each Candy account row.
Keep the existing account enabled switch as a separate control. Keep the other evaluation tabs unchanged.
Keep the global Settings switch, collapsed configuration, and immediate inspection controls unchanged.

Store a Boolean `codex_eval_inspection_enabled` in each native account's existing `extra` object.
Read an absent value as true to preserve current inspection participation. No schema migration is required.
Expose the saved value through the existing account inventory projection.
Save only this field through an administrator account endpoint and `AdminService.UpdateAccountExtra`.
That service uses the existing repository JSONB merge. Do not submit a whole account or `extra` snapshot from the page.
Preserve the participation key during ordinary account edits that do not explicitly update it.
Use a per-account pending state, preserve row nodes during refresh, and show save failures beside the row switch.

Carry participation into `InspectionAccount` and check it at initial selection, later test admission, and before priority updates.
An in-flight request can finish after opt-out. Do not submit later inspection tests or write a priority after a fresh opt-out read.
Apply participation to both scheduled and immediate coordinator rounds. Manual Candy evaluation remains independent.
The global setting controls the timer. Account participation controls membership without changing that timer or native schedulability.
An account that appears in multiple groups has one saved participation value.

The account endpoint is `PUT /api/v1/admin/accounts/:id/codex-eval-inspection`.
Its body is `{"enabled": false}` or `{"enabled": true}`.
Bind the field as `*bool` so false is valid and an absent value fails validation.
Accept native OpenAI OAuth and API-key accounts, including their supported shadow accounts.
Return `{account_id, enabled}` through the native response envelope after the atomic metadata merge succeeds.
Do not send the whole account object or update native `schedulable`.

Add the metadata key to the native full-update managed-key preservation list.
The active server uses `AdminService.UpdateAccount`. The older `AccountService.Update` has no active route caller and needs no change.
The existing bulk JSONB merge, account DTO redaction, and lightweight DTO already retain the key.
Keep native copy semantics: duplicate and export/import copy persistent configuration. A shadow keeps its own independent value.
A fresh account without the key follows the absent-value default. Do not add a separate migration or copy mechanism.

Add `InspectionAccount.AutoInspectionEnabled` separately from the existing `Enabled` field, which means native schedulability.
Map an absent metadata value to true and an explicit Boolean false to false.
Use that field in initial selection and fresh eligibility checks.

`testPhase` starts a multi-test confirmation task once. `runCandy` then loops internally without another eligibility read.
Checking only `testPhase` cannot enforce opt-out between those sends.
Add a private checked-start path that carries a per-account pre-send callback into the Candy worker.
Keep public `Start` and the existing private `start` entry usable under their original contracts, with no callback for manual tests.
Invoke the callback before each inspection request without holding the runtime mutex.
The initial-request callback checks current participation and eligibility without requiring the old list priority.
For confirmation requests, it also checks the fresh transition-priority baseline.
An explicit opt-out ends remaining work without another request or synthetic error record.
Keep completed records and use the existing `unchanged` inspection outcome for intentional opt-out.
An actual account-read failure remains an error and blocks further sends.
An incomplete batch cannot produce a recommendation, priority update, or notification.

## Redundant Error Summary

Remove the header-wide `inspection-runtime-error` output from the Candy account header below `账户` and `仅启用`.
Stop combining every account test failure into that header summary.
Keep the original test error in each account's existing Candy result and history.
Delete inspection runtime and notification error output. Do not move that output into Settings or another page area.
Remove the corresponding error aggregation and DOM writers in runtime-state polling and runtime actions.
Keep authentication handling and polling recovery. Keep backend error state, notification behavior, and original Candy test records.
Suppress runtime failure captions that would recreate the removed error prompt.
Keep global settings errors in Settings. Keep account switch save errors beside the affected account control.
Update all element references and UI assertions so removal cannot break state polling or settings saves.
Keep running progress and the exact `下次巡检：` label with its scheduled timestamp. Remove any empty error-row spacing.

## Delivery Boundary

Use existing administrator authentication and account metadata services.
Keep production edits in the seven-patch delivery and leave the original upstream tree unchanged.
Use the existing build/check entry and retained fixture for local verification.
The check scope is defined in `implement.md`. Unchanged native copy, DTO, bulk-update, and installer paths need no additional tests.
No production deployment belongs to this task.

## Evidence and Remaining Investigation

See `research/source-findings.md` for exact source anchors.
The delayed-refresh probe confirms readiness changes during an unchanged successful poll.
See `research/inspection-count-contract.md` for the executed stale-priority probe and its limits.
See `research/account-participation-contract.md` for the native API, metadata, copy, and per-send admission paths.
See `research/refresh-error-contract.md` for the executed delayed-refresh probe and header-error owners.
The probes establish code behavior. The exact production timing behind the user's count report remains unobserved.
Actual browser acceptance and checks of the repaired candidate belong to implementation.
