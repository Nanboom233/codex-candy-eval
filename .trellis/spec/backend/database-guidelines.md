# Native Persistence Guidelines

## Active Sub2API storage

Sub2API stores inspection settings in its existing PostgreSQL `settings` table.
The four evaluation histories remain in one JSON file.
Tasks, progress, inspection summaries, and the next run time stay in memory.

### 1. Scope / Trigger

Use native settings persistence for Candy inspection configuration.
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

Keep history keys `results`, `fingerprints`, `modeltraces`, and `pelicans`.
Each map uses a native account ID as a decimal string.
Retain the latest 20 records per account and evaluation type.
Hold `mu` during shared history access and `saveStateLocked`.
`loadState` acquires its own lock and reads once per process.
A missing file means no saved history.
Other read/decode errors prevent later overwrite. A JSON `null` snapshot is invalid.

Save all four maps together.
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
| Unreadable, malformed, or null history | Show the error and prevent overwrite |
| History save fails | Retain memory results and show the error |
| History clear save fails | Restore previous maps and return failure |
| Settings read fails | Return the repository error |
| Settings validation fails | Reject before writing |
| Settings write fails | Return the error without changing active settings |

### 5. Good / Base / Bad Cases

Good: blank notification secret inputs retain existing native settings values.

Base: absent history starts empty. Absent inspection keys use defaults.

Bad: a null history file must not become an empty snapshot after the next result.

### 6. Tests Required

Use existing native evaluation persistence tests for load protection, latest-20 retention, saves, and clear rollback.
Use `tests/backend/internal/service/codex_inspection_settings_test.go` for settings persistence, secrets, validation, and write failures.
Use the disposable native fixture to check PostgreSQL rows and restart loading.
Tests use temporary files and synthetic data. Do not delete real history or settings.

### 7. Wrong vs Correct

Wrong: put inspection credentials in the history snapshot because both features share the evaluation runtime.

Correct: retain history in JSON and save inspection settings through the native settings service.
