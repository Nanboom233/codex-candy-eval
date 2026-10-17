# Review Evidence

## Backend and Contract Review

- `patches/0001-evaluation-core.patch` keeps public manual start entry points and adds the inspection-only pre-send callback.
- `patches/0005-inspection.patch` rereads current priority after the first test and uses that value for confirmation and update baselines.
- A degraded account with a wrong first answer does not enter a confirmation group, so it keeps one test and its existing priority.
- The account participation check returns before the next send when the saved value is false. That branch sets `unchanged` and does not create a test error.
- `patches/0007-sub2api-integration.patch` registers the account endpoint and preserves `codex_eval_inspection_enabled` during ordinary full account edits.
- The dedicated endpoint accepts an explicit Boolean, uses the existing account-extra merge, and does not change native `schedulable` state.
- Manual `Start` remains independent from account participation.

## Scope Review

The changed tests match the task check table. They cover the stable and stale-priority degraded-wrong cases, participation admission and opt-out, manual fixed counts, endpoint binding, and metadata preservation. No restart, copy, import/export, DTO, custom-priority, repeated-round, race, coverage, full-suite, or installer-restoration matrix was added.

## Reused Validation

The implementation owner recorded a successful `python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check` run in `implementation-build.log`. It passed the existing UI harness, focused Go tests, vet, build, and package steps. This review did not repeat that valid run.

## Resolved Browser Item

The browser review found a narrow-screen alignment defect between the native account switch and the automatic-inspection switch.
The implementation owner corrected Candy's narrow-screen CSS and rebuilt the package after the affected UI check passed.
The main session inspected the corrected package in the retained fixture. Both visible rows have a measured switch offset of zero pixels.
See `browser-evidence.md` for the screenshots and observations. No review item remains open.
