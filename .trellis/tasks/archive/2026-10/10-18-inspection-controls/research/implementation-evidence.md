# Implementation Evidence

## Candidate

Production changes remain in patches `0001`, `0005`, `0006`, and `0007`.
Tests remain under `tests/backend/` and `tests/ui/`.
The retained development source is `/tmp/inspection-controls-implementation`.
This source has no attached Go tests. Its frontend dependency link uses the existing original-source installation.
The original source remains `/tmp/sub2api-codex-eval-v0217-check`.

## Changes

- Keep successful inventory readiness during a same-group refresh. Invalidate readiness after an actual inventory failure.
- Use current priority after the first inspection test. Use that priority before confirmation sends and priority updates.
- Store participation in `extra.codex_eval_inspection_enabled`. An absent value means enabled.
- Add the atomic administrator account update endpoint. Preserve the participation key during ordinary full account edits.
- Check inspection participation before each Candy request. Keep public manual entry points and fixed counts.
- End intentional opt-out without a synthetic test error. Keep completed records.
- Add the Candy account switch column and row-local save feedback. Preserve pending state during refresh.
- Remove runtime and notification error prompts. Keep original records, Settings errors, and `下次巡检：` with its time.

## Validation

The following entry completed with exit code 0:

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check
```

The entry passed frontend checks and build, installer-version checks, the UI harness, focused Go tests, and changed-package vet.
It then built the Linux/amd64 executable and archive.
See `implementation-build.log` for the command output.
The temporary build source was removed by the standard entry after completion.

The UI harness covers one delayed successful refresh, independent switches, pending refresh, save rollback, retained record errors, and removed prompts.
The coordinator tests cover stable and stale-list degraded-wrong accounts with adjustment on and off.
They also cover initial opt-out, opt-out during a blocked confirmation, and manual fixed counts for the opted-out account.
Existing tests supply unchanged confirmation, recovery, cancellation, membership, and failure evidence.
The handler test covers required Boolean binding, false saves, absent-value defaults, and key-only updates.
The service test covers participation preservation during a full account edit.

`gofmt`, UI harness syntax, `git diff --check`, and patch parsing passed.
After the build started, only unused patch index metadata was removed. Effective source did not change.

## Deliverables and Handoff

- `dist/build/sub2api-codex-eval`
- `dist/build/sub2api-codex-eval.tar.gz`

The main session owns the one browser acceptance pass, independent Trellis review, spec integration, and commit workflow.
No fixture or browser operation ran in this implementation assignment.
The build emitted existing frontend dependency and chunk warnings. No check failed.

## Narrow Layout Correction

The browser showed different vertical positions for the native account switch and inspection switch at 390px.
Candy now aligns both controls at the start of their shared metadata row.
The rule applies only to Candy at narrow widths. Other evaluation tabs and the desktop grid retain their layout.
The inspection cell keeps its row-local error output.

The existing UI harness passed after the CSS correction. Backend checks remain valid for unchanged source.
The standard build entry completed without `--check` and refreshed both fixed outputs.
See `layout-build.log`. The final executable contains the Candy-only alignment rule.
The final selector scope was copied to the active build source before Go embedded the page.
The main session owns the affected browser recheck.
