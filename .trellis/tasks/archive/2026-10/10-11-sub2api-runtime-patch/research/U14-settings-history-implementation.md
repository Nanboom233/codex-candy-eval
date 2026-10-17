# U14 Settings Persistence and Cancellation History

## Scope

This correction covers saved evaluation preferences, native inspection settings, and the reported cancellation text in history.
The change uses the existing browser preferences and native settings repository.
It adds no store, setting, cancellation policy, or timeout.
U13 account membership and shared result displays remain part of the focused UI check.

## Reproduced settings failures

The four manual model controls saved preferences on `change` only.
A custom model edit did not reach local storage before blur.
The regression types `typed/model`, sends its `input` event, and creates a new page with the same local storage.
Before the repair, Candy restored `saved/literal` instead of `typed/model`.
The same event gap affected existing number and effort controls.

Inspection settings had a separate stale-read path.
`inspectionSettingsLoaded` prevented all later settings GET requests within the page.
The regression changes the synthetic server values after the first settings read, then reopens settings.
Before the repair, the page kept `saved/literal` instead of `saved/other-session`.
An in-memory omission of the first assertion exposed this second independent failure.
That control check did not change repository files or product behavior.

## Settings changes

Existing manual preferences now save on `input` as well as `change`.
Existing change handlers still normalize number controls and render results.
The saved fields remain unchanged:

| Evaluation | Existing saved fields |
| --- | --- |
| Candy | Model, effort, runs |
| Fingerprint | Model, mode, concurrency |
| ModelTrace | Model, concurrency |
| Pelican | Model, effort |

Each settings-tab entry now reads native inspection settings again.
The form applies saved values only to fields without unsaved edits.
A save invalidates any older settings read, so a late GET cannot replace the completed save.
Settings reads do not run during an active save.

Dirty fields now retain their edit revision.
A successful save clears only submitted fields that the user has not edited again.
The response then applies the saved public values while preserving later edits.
This also clears a saved secret input when another field changes during the save.
A failed save retains the form edits and shows the error.

The native service directly reads and writes the existing PostgreSQL settings through `GetMultiple` and `SetMultiple`.
It does not use a settings cache for this feature.
It applies coordinator settings only after a successful repository write.
Existing settings tests passed for restart reads, secret preservation, write failures, initial-load failure, and priority validation.
The investigation found no settings-service or repository failure that requires a backend product change.

## Cancellation investigation

Manual `Start` creates a task context from `context.Background()`.
Inspection tasks inherit the coordinator run context.
`SendCodexEvaluation` retains that task context through the native HTTP transport and response body.
The native `doUpstreamRequest` helper cancels an attempt when its body closes, not when response headers return.
The Chat Completions bridge reads from that same body.
The existing probe and answer timeout limits remain unchanged.

The retained fixture's old cancellation record came from the earlier explicit-stop scenario.
That record is not evidence of unintended cancellation.
No existing history record was changed or removed.

The new integration regression uses actual native components:

1. The native `CodexEvalStart`, `CodexEvalState`, and `CodexEvalCancel` handlers serve local HTTP requests.
2. The runtime calls `AccountTestService.SendCodexEvaluation` through the configured sender.
3. `repository.NewHTTPUpstream` sends to a local supplier that flushes headers, then waits before sending response bytes.
4. The test waits until the HTTP start request context ends and completes a separate state request.
5. Releasing the supplier produces successful Candy records with actual token counts in `history.json`.
6. A second blocked request receives explicit stop and produces one cancelled record with `context canceled`.

The check covers the Responses stream and the native Chat Completions conversion.
The normal run completes three Candy calls after the start request ends.
Explicit stop cancels the active supplier request and prevents the remaining Candy calls.
This delayed-body check did not reproduce unintended cancellation.
It does not prove the cause of an unobserved production error.
Cancellation product code remains unchanged pending contrary evidence.

## Changed files

- `native/backend/internal/codexeval/web/candy.js`: save existing preferences on input.
- `native/backend/internal/codexeval/web/fingerprint.js`: save existing preferences on input.
- `native/backend/internal/codexeval/web/modeltrace.js`: save existing preferences on input.
- `native/backend/internal/codexeval/web/pelican.js`: save existing preferences on input.
- `native/backend/internal/codexeval/web/inspection.js`: reread native settings and protect edits from asynchronous responses.
- `tests/ui/sub2api-dropdowns.cjs`: reload-before-blur, native reread, and delayed read/save regressions.
- `native/backend/internal/service/codex_eval_lifetime_test.go`: actual native handler, transport, delayed stream, cancellation, and saved-history regression.

## Focused checks

The UI check passed:

```sh
SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs
```

It checks all four custom models, Candy runs, both concurrency fields, and Candy/Pelican effort after reload before blur.
It also checks native settings rereads, preserved dirty fields, cleared submitted fields, secret redaction, and old GET responses after saves.
Existing U11/U12 and U13 assertions passed in the same check.

Go checks use the retained native check workspace without a delivery build.
The relevant runtime, model, sender, handler, and settings product files match the current source.
The overlay maps the new lifetime test and current inspection settings tests into that workspace.
Its path is `/tmp/sub2api-U14-go-overlay.json`.
Go resolves the workspace under `/private/tmp`, so the overlay uses that absolute path.

```sh
cd /tmp/sub2api-U11-check.CqGH5I/backend
go test -v -tags=unit -overlay /tmp/sub2api-U14-go-overlay.json ./internal/service -run '^TestCodex(InspectionSettings|EvalBackgroundNativeStreamLifetime)' -count=1 -timeout=60s
go vet -tags=unit -overlay /tmp/sub2api-U14-go-overlay.json ./internal/service
```

The focused settings and lifetime tests passed.
The changed package passed vet.
Changed JavaScript files passed `node --check`.
The new Go test passed `gofmt`.
`git diff --check` passed.

## Final build and package

The independent review passed without product edits.
See `U14-settings-history-check.md` for its findings and checks.
The final build then used the existing original source and dependency caches:

```sh
python3 scripts/build-sub2api.py --source /tmp/sub2api-codex-eval-v0217-check > /tmp/sub2api-U14-build.log 2>&1
python3 scripts/package-sub2api.py
```

Both commands passed.
The native frontend passed three locale tests, its TypeScript check, and its Vite build.
Vite completed in 9.34 seconds.
The Go build produced the static Linux/amd64 executable with native version `0.2.17` and release build type.
The final files contain the reviewed U14 correction.

| Output | Size |
| --- | ---: |
| `dist/build/sub2api-codex-eval` | 125595808 bytes |
| `dist/build/sub2api-codex-eval.tar.gz` | 38700266 bytes |

Executable SHA-256: `b6f4fec67ed0395d4330a734423d8451190a2456a1ced5f891567dfee8d24b5d`.
Archive SHA-256: `380f1c2f74667dd8eb6c11ca0c5bcc116d8c984a0e5cebeff000fa8e5caf5452`.
`file` identifies a stripped, statically linked x86-64 ELF executable.
The archive has seven entries: the executable, installer, and five license files.
It contains no development or test files.
The build did not change product source or repeat the focused checks.
The earlier historical cancellation event remains unattributed.

## Handoff

Source review, final build, and packaging are complete.
Actual settings controls and persistence still require the main session's final retained-container browser check.
The main session owns the retained-container and browser checks.
No Docker, browser, paid request, production deployment, staging, commit, or cleanup action ran in this turn.
No race check, broad suite, coverage gate, or additional delivery requirement was added.
