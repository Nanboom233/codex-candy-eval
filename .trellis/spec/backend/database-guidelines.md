# Native Persistence Guidelines

## Active Sub2API storage

Sub2API stores inspection settings in its existing PostgreSQL `settings` table.
The three evaluation histories remain in one JSON file.
Tasks, progress, inspection summaries, and the next run time stay in memory.

### 1. Scope / Trigger

Use native settings persistence for Candy inspection configuration.
Use native account `extra` for inspection participation, verdict, and actual removed-model names.
Use the existing JSON store for evaluation history only.

### 2. Signatures

`SettingService.GetCodexInspectionSettings` and `UpdateCodexInspectionSettings` call native `SettingRepository.GetMultiple` and `SetMultiple`.
Keys use `codex_eval_inspection_` plus each settings field name.
See [Native Candy Inspection](./native-inspection.md) for fields and administrator API signatures.

Native history owners are `backend/internal/codexeval/state.go` and `runtime.go` in the applied evaluation-core patch.
`Initialize` selects `<native-data-directory>/codex-eval/history.json`.
`loadState`, `saveStateLocked`, and `clearHistoryResponse` own file operations.

### 3. Contracts

Store each inspection setting as its JSON field value under an inspection-specific key.
Do not add a table, direct SQL path, separate storage service, or `inspection.json`.
Public reads clear notification secrets and return configured flags.
Empty secret updates preserve saved values.
Change active coordinator settings only after the repository write succeeds.
Initial-read failure must not allow inspection with default settings.
Keep starts blocked until a successful native settings save configures the coordinator.

Store account participation as Boolean `extra.codex_eval_inspection_enabled`. An absent value means enabled.
Use `AdminService.UpdateAccountExtra` for a key-only atomic merge. Preserve that key during unrelated full account edits.
Native account storage supplies cross-group and restart persistence. Keep it separate from global settings and evaluation history.

Store account verdict as `normal` or `degraded` in `extra.codex_eval_inspection_verdict`. An absent verdict means normal.
Actual priority does not initialize the verdict. Do not add old-priority inference or migration.
Use a key-only atomic merge for verdict transitions without model changes. Preserve the key during unrelated full account edits.
Apply model control after each completed detection establishes the account state, including stable states.
Skip writes when model data, pending removals, and verdict already match. An absent normal verdict needs no initialization write.
Persist the verdict before an optional actual-priority update. A failed priority write does not reverse the saved verdict.
Keep verdicts out of the evaluation history snapshot and global settings.

Store actual removed model names in `extra.codex_eval_inspection_removed_models` as a string array.
Merge new actual removals without duplicate names.
An absent or empty array means no pending restoration. Ordinary full account edits preserve this key.
Whitelist identity entries use `credentials.model_mapping` with equal keys and values.
Save changed model mappings, pending removals, and verdict through one native `AccountRepository.Update`.
Copy mutable maps and preserve unrelated credentials and metadata before that call.
Use the existing repository transaction, scheduler outbox, and snapshot refresh.
Do not send a partial credentials object through generic administrator updates. They preserve omitted sensitive fields only.
Recovery reads the saved removal list even when the global list changes or model control becomes off.
Keep current administrator mappings, remove the reserved identity placeholder, and clear pending removals in the same write.
No extra table, snapshot, scheduling flag, original scheduling value, or restore phase is necessary.
See [Native Candy Inspection](./native-inspection.md) for confirmed-state application and placeholder behavior.

Write root field `schema_version: 2` and history keys `results`, `fingerprints`, and `pelicans`.
Each map uses a native account ID as a decimal string.
Pelican records store `variant` as `original` or `all-in-one`. Records can include optional `attempts`.
Read fingerprints as raw JSON until the root marker identifies the current format.
An absent marker selects legacy loading. Decode only Candy and Pelican. Start the fingerprint map empty.
Discard obsolete fingerprint and ModelTrace payloads without per-record conversion.
Schema 1 retains all three current maps. Both older formats mark retained Pelican records as `original`.
Assign loaded maps only after retained fields decode successfully.
Before migration, copy the original bytes to same-directory `history.v0.json` or `history.v1.json` with mode `0600`.
If the copy fails, retain memory records, expose the error, and block replacement and later saves.
Use the existing atomic saver for migration. A failed save retains memory records, the original file, and the recovery copy.
Schema 2 needs no migration copy or write. Invalid and future markers keep the original file protected.
An executable downgrade needs a compatible old history snapshot. That snapshot excludes records saved after the copy.
Retain the latest 20 records per account and evaluation type.
Hold `mu` during shared history access and `saveStateLocked`.
`loadState` acquires its own lock and reads once per process.
A missing file means no saved history.
Other read/decode errors prevent later overwrite. A JSON `null` snapshot is invalid.

Save all three maps together.
Write `statePath + ".tmp"` with mode `0600`, then rename it over the destination.
Keep this temporary file beside the destination.
Atomic replacement does not provide an fsync durability guarantee.
Keep memory results and expose save errors.
A failed clear restores previous history maps. Clear does not cancel running tasks.
Restart reloads saved history/settings without restoring tasks or the previous inspection schedule.

### 4. Validation & Error Matrix

| Condition | Result |
| --- | --- |
| Missing history file | Start with empty histories |
| Unversioned native history | Load Candy/Pelican, mark Pelican original, copy original bytes, then save schema 2 with empty fingerprints |
| Version-1 native history | Retain all maps, mark Pelican original, copy original bytes, then save schema 2 |
| Version-2 native history | Load all maps without a migration copy or write |
| Migration copy fails | Retain memory records and block later saves |
| Unreadable, malformed, null, or future-schema history | Show the error and prevent overwrite |
| History save fails | Retain memory results and show the error |
| History clear save fails | Restore previous maps and return failure |
| Settings read fails | Return the repository error |
| Settings validation fails | Reject before writing |
| Settings write fails | Return the error without changing active settings |
| Account verdict save fails | Keep the previous verdict and prevent the optional priority write |
| Compound account update fails | Keep stored mapping, removed-model list, and verdict unchanged. Prevent the optional priority write. |

### 5. Good / Base / Bad Cases

Good: blank notification secret inputs retain existing native settings values.

Good: one account update saves the changed mapping, actual removal list, and degraded verdict together.

Base: absent history starts empty. Absent inspection keys use defaults.

Bad: a null history file must not become an empty snapshot after the next result.

### 6. Tests Required

Use existing native evaluation persistence tests for load protection, latest-20 retention, saves, and clear rollback.
Extend that fixture for retained maps, original variants, migration copies, write failures, and schema-2 reload.
Do not add a separate migration suite.
Use `tests/test_service_codex_inspection_settings_test.go` for settings persistence, secrets, validation, and write failures.
Reuse existing persistence evidence. Add focused checks only for changed storage behavior without coverage.
Account participation needs key-only-save and full-edit-preservation assertions, without a separate restart or copy test matrix.
Extend the same metadata fixtures for the verdict key and absent-state normal default.
Extend the full-edit fixture for pending removals. Use account-helper tests for compound model writes and recovery.
Use the focused coordinator tests for verdict-save failure. Do not repeat lifecycle checks in storage tests.
Tests use temporary files and synthetic data. Do not delete real history or settings.

### 7. Wrong vs Correct

Wrong: put inspection credentials in the history snapshot because both features share the evaluation runtime.

Correct: retain history in JSON and save inspection settings through the native settings service.
