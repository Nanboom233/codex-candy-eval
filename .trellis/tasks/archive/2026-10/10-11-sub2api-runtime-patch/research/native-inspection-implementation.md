# Native Candy Inspection Implementation

## Change boundary

Extend `native/backend/internal/codexeval` with one process coordinator.
Reuse Candy admission, run IDs, result history, and cancellation.
Keep automatic account selection separate from manual evaluation admission.

Native `SettingService` reads and writes inspection keys through `SettingRepository`.
Native administrator settings handlers protect configuration access.
The server starts the coordinator after normal construction and stops it during shutdown.

The fifth settings tab uses the existing page controls and native settings card layout.
The fixed binary and archive include the completed change.

## API contract

All routes use the native administrator response envelope.

| Method | Path | Data |
| --- | --- | --- |
| GET, PUT | `/api/v1/admin/settings/codex-eval-inspection` | Persisted inspection settings. |
| GET | `/api/v1/admin/codex-eval/inspection/state` | Runtime state. |
| POST | `/api/v1/admin/codex-eval/inspection/start` | Start immediately. Return 409 when a run exists. |
| POST | `/api/v1/admin/codex-eval/inspection/stop` | Cancel the current run. Keep the automatic setting. |

Settings fields: `enabled`, `model`, `effort`, `max_tests_per_account`, `min_wait_minutes`, `max_wait_minutes`, `include_patterns`, `exclude_patterns`, `skip_free`, `adjust_priority`, `normal_priority`, `degraded_priority`, `ntfy_server`, `ntfy_topic`, `ntfy_username`, `ntfy_password`, and `bark_url`.

Pattern arrays retain the original script's case-sensitive shell wildcard matching against account name and email.
Defaults retain model `gpt-6.1-sol`, effort `low`, and Free filtering.
The test limit defaults to three. Waiting defaults to 15–35 minutes.
Automatic inspection and priority changes default off.
Native normal and degraded priorities default to zero and 99.

Reads clear secret values and return `ntfy_password_configured` and `bark_url_configured`.
Empty secret updates retain saved values.

Runtime fields: `enabled`, `running`, `cancelling`, `phase`, `next_run_at`, `started_at`, `finished_at`, `total`, `done`, `error`, and `results`.
Each result contains account identity, original priority, optional target priority, run IDs, test count, outcome, error, and notification error.
Runtime state stays in memory. Configuration saves do not cancel a current round.

## Implementation and checks

The native coordinator, settings service, administrator handlers, settings tab, and lifecycle patch are complete.
Automatic inspection uses native availability checks. Manual evaluation keeps its existing admission rules.
The coordinator completes all confirmation groups before it changes account priorities.
Notifications follow successful changes. Notification failures remain separate from evaluation results.

The target source is actual Sub2API v0.2.17.
The upstream archive retains a stale v0.2.16 VERSION file. Release linker flags set the built version to v0.2.17.
The installer accepts installed v0.2.16 or v0.2.17 binaries. It checks the candidate version before replacement.

Completed checks:

- Native frontend i18n tests, vue-tsc, and Vite build passed.
- Focused native codexeval, service, and administrator handler tests passed.
- Changed-package Go vet and the static Linux/amd64 build passed.
- Focused root inspection runtime tests passed.
- Installer version tests passed for v0.2.16 upgrade and wrong-candidate rejection.
- JavaScript syntax, Python compilation, shell syntax, patch dry-run, and diff whitespace checks passed.

The actual Docker integration passed administrator authentication, literal-model sending, filtering, and confirmation counts.
Counts were wrong=3, mixed=3, correct=1, disabled=0, and Free=0.
The wrong account changed priority from 0 to 99. The mixed account retained priority zero.
Synthetic ntfy and Bark failures appeared separately after the successful priority change.
Stop cancelled the current run. Restart restored all 17 settings keys and secrets without resuming the interrupted run.
The runtime reported v0.2.17. An actual installer reinstall completed and reported ready after 1.1 seconds.

The Docker environment restarted after these checks. Its host temporary source, supplier, credential, and log files disappeared.
This removed the browser route before the pending browser check. No browser pass is claimed here.
Earlier failed fixture runs exposed missing timezone data, native compliance acceptance, and incompatible fake supplier streaming.
The fixture now uses UTC, native compliance acceptance, and both Responses and Chat Completions streams.

## Delivery

The completed build produced `dist/build/sub2api-codex-eval` and `dist/build/sub2api-codex-eval.tar.gz`.
The archive contains the binary, installer, and necessary license texts.
Development fixtures, tests, and a standalone cleanup script are excluded.
The existing source build and package results remain valid during fixture consolidation.

## Reusable Docker fixture

The user requires one retained container named `sub2api-codex-eval-test`.
The retained application container keeps its existing browser port, 60198.
Development startup and supplier files live under `tests/installer/native_fixture_*`.
The startup wrapper runs PostgreSQL, Redis, the fake supplier, and the existing built application inside that container.
The fixture entry point prepares or reuses this fixed container. It does not create random service containers.
The main session owns removal of obsolete task containers after synthetic database migration.
Recovery completed with the retained container ID `bd6c9583c3e6` and origin `http://127.0.0.1:60198`.
The PostgreSQL dump restored the previous synthetic database inside the retained container.
Its five test accounts now use the internal fake supplier at `http://127.0.0.1:18081`.
The application reports v0.2.17. Administrator login and native settings reads and saves passed after recovery.
The first migration attempt lacked the PostgreSQL `pg_trgm` extension. The development wrapper now installs `postgresql18-contrib`.
No product change was required for recovery.

Run `python3 tests/installer/native_inspection_fixture.py` to prepare or reuse this container.
The command prints the browser origin and the path to a mode-600 credential file under `/tmp`.
It preserves the container ID, port, and database. It resets only synthetic account endpoints and inspection test settings.
The command does not create containers. It stops and starts the retained test container during preparation.
For a simple restart without test-setting changes, run `docker restart sub2api-codex-eval-test`.
All test dependencies restart inside that same container.
Keep this development wrapper separate from installer replacement checks because the wrapper is not the delivered executable.

Python compilation, shell syntax, and diff whitespace checks passed for the consolidated fixture.
The main session can now remove the previous PostgreSQL and Redis containers.
The browser check remains the main session's responsibility.

## Changed files

- `native/backend/internal/codexeval/inspection*.go`: settings, coordinator, notifications, and focused runtime tests.
- `native/backend/internal/codexeval/runtime.go`: shared private start path with inspection cancellation and ungrouped automatic targets.
- `native/backend/internal/codexeval/web.go` and `web/{app.js,inspection.js,ui.html,style.css}`: fifth settings tab and controls.
- `native/backend/internal/service/codex_inspection_settings*.go`: native database settings and service tests.
- `native/backend/internal/handler/admin/codex_inspection_handler.go`: native administrator settings and runtime handlers.
- `native/sub2api.patch`: routes, startup, and shutdown integration.
- `scripts/build-sub2api.py`, `scripts/test-sub2api.sh`, and `.github/workflows/release.yml`: v0.2.17 target.
- `install.py` and `tests/installer/test_native_version.py`: supported installed versions and candidate-version checks.
- `tests/installer/native_inspection_fixture.py` and `native_fixture_*`: one retained development runtime.

The main session owns README, operations documentation, specs, and task metadata.

## Remaining work and limits

The independent Trellis check and actual browser settings checks remain pending.
The main session owns specification updates, task status, work-commit review, and task closure.
No production calls, paid requests, deployment, or cache/image/volume pruning ran.
These local checks do not claim N150 or iStoreOS execution.
