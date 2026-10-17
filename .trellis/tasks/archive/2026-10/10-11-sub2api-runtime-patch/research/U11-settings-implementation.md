# U11 Settings Simplification

## Scope and decisions

U11 replaces the earlier include/exclude account-filter design.
The page and coordinator no longer expose or apply those filters.
The Free switch and native availability and busy checks remain.
All five model dropdowns retain literal names and independent raw model discovery.

## Removed filters

Removed `IncludePatterns` and `ExcludePatterns` from `InspectionSettings`.
Removed wildcard validation, the name/email matcher, and its coordinator selection branch.
Removed the now-unused inspection account email field and native adapter assignment.
The settings service derives its read/write fields from this type, so it no longer reads or writes either filter key.
The settings type now defines 15 persisted keys and two public configured flags.
Old database rows remain untouched and have no effect on account selection.
No database deletion, implicit write, or migration framework was added.

Removed the include/exclude settings card, picker helpers, placement code, CSS, and obsolete picker assertions.
Moved the existing Free switch and model-catalog error into the automatic inspection card.
Removed the obsolete filter field from the development fixture settings update.
The complete OpenAI inventory still supplies raw model suggestions independently of the current evaluation group.

## Priority

Normal priority now defaults to 10.
New settings writes reject normal priorities below 1.
The native read converts a stored nonpositive normal priority to the new default in memory.
That read leaves the saved database value unchanged.
Only an explicit successful save writes the selected value.
Both priority inputs hide while automatic adjustment is off and appear when it is on.
Hiding does not discard their current values.
The normal input has minimum 1 and default 10.

## Status and notifications

Removed idle `等待巡检`, `自动巡检已关闭`, and the unscheduled fallback explanation.
The status row stays hidden when it has no phase, next timestamp, or actual progress.
Scheduled timestamps, current progress, errors, completed results, and action feedback remain.
Removed notification configured-status spans and their UI update code.
The backend still redacts secrets and preserves saved secrets when the submitted inputs are blank.
Public configured flags remain part of that native response contract but have no page status text.

## Changed files

- `native/backend/internal/codexeval/inspection_settings.go`: removed filter fields and validation, changed priority default and minimum.
- `inspection.go`: removed matching and the unused account email.
- `inspection_test.go`: confirms all eligible names are tested, with the new managed normal priority.
- `native/backend/internal/service/codex_inspection_settings.go`: normalizes an old stored nonpositive priority on read.
- `codex_inspection_settings_test.go`: checks read-only normalization, removed-key exclusion, explicit writes, and preserved secret behavior.
- `native/backend/internal/handler/admin/codex_inspection_handler.go`: removed unused account email extraction.
- `native/backend/internal/codexeval/web/{ui.html,inspection.js,components.js,style.css,app.js,catalog.js}`: simplified controls and status.
- `tests/ui/sub2api-dropdowns.cjs`: retains model checks and adds changed settings visibility, validation, and state checks.
- `tests/installer/native_inspection_fixture.py`: removed obsolete filter configuration.

## Checks

- `go test ./native/backend/internal/codexeval -run '^TestCodexInspection' -count=1` passed.
- Native `go test -tags=unit ./internal/service ./internal/handler/admin -run '^TestCodexInspection' -count=1` passed.
- The native administrator package compiled. That test filter matched no administrator test cases.
- `go vet ./native/backend/internal/codexeval` passed.
- Native `go vet ./internal/service ./internal/handler/admin` passed.
- The focused jsdom settings checks passed for priorities, idle/progress/errors, removed fields, and blank notification inputs.
- Existing dropdown checks passed for five literal/custom models, cross-group inventory, disabled discovery, save/reload, dirty edits, and stale responses.
- JavaScript syntax, Go formatting, and diff whitespace checks passed.
- The development fixture passed Python syntax checks.
- The native build passed through `scripts/build-sub2api.py`.
- The frontend build passed its three locale tests, TypeScript checks, and Vite build.
- Packaging passed through `scripts/package-sub2api.py`.

Service checks used a standard `/tmp/sub2api-U11-check.*` copy of the original native backend with the current overlay.
The build used `/tmp/sub2api-codex-eval-v0217-check` as its source directory.
The build log is `/tmp/sub2api-U11-build.log`.
The final executable is a stripped, static Linux ELF for amd64.

| Output | Size |
| --- | ---: |
| `dist/build/sub2api-codex-eval` | 125587616 bytes |
| `dist/build/sub2api-codex-eval.tar.gz` | 38701207 bytes |

The archive has seven entries: the executable, the installer, and license files.
The archive contains no development or test files.
No Docker, production, paid-request, staging, commit, or cleanup action ran.
The main session owns candidate loading and actual browser verification.
Actual browser verification and the independent Trellis check remain pending.
