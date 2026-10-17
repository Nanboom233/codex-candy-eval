# Execution Plan

## Implementation Gate

- [x] Obtain approval of the latest final planning summary.
- [x] Start this task with `python3 .trellis/scripts/task.py start 10-18-inspection-verdict-state`.
- [x] Dispatch one Trellis implementation owner with these artifacts and the curated context.

## Ordered Work

1. Add the two-state verdict type, account field, native Extra projection, and persistence hook. Default absent state to normal.
2. Preserve the verdict key during full account edits in patch `0007`.
3. Replace priority-based confirmation selection with the effective verdict contract in `design.md`.
4. Save confirmed verdict transitions in both adjustment modes. Reuse existing participation and cancellation checks.
5. On a confirmed verdict change, check the priority switch and use the existing guarded update path only when enabled.
6. Extend existing tests only for the missing assertions in the table below.
7. Format changed Go source and update the affected production patch hunks.
8. Run the existing build/check entry for the final candidate.
9. Dispatch the Trellis check owner with the candidate and valid check evidence.
10. Resolve material findings and update the native inspection and persistence specs.
11. Record results and present the work diff with an exact commit proposal.

Do not add stable-state priority synchronization or automatic write retries. Keep existing notifications and their text unchanged.
Remove priority-based classification without a migration or compatibility path. Custom actual priorities do not affect state transitions.

## Check Ownership

This table defines the complete task-specific checks. Research and historical specs do not add acceptance gates.
Extend existing tests. Add a focused test only for a new persistence failure boundary without an existing case.

| Requirements | Existing test owner | Missing assertions |
| --- | --- | --- |
| R1–R4, AC1–AC2 | `TestCodexInspectionDegradedWrongOnce` | Replace its old adjustment/list-priority matrix with one adjustment-off lifecycle. Confirm degradation, check one wrong test next round, confirm recovery, then check one correct test. Assert durable verdicts and unchanged actual priority. |
| R1–R4, AC3, AC5 | `TestCodexInspectionConfirmationAndBoundaries`, `TestCodexInspectionCancelAndSchedule`, `TestCodexInspectionParticipationStopsLaterSends` | Add verdict-save assertions to the existing transition, stable, mixed, provider-error, cancellation, and opt-out cases. Use explicit verdicts for former priority-derived fixtures. Keep existing priority and notification assertions. |
| R1, R4, AC6 | Existing coordinator test file | Check the two new write boundaries: verdict-save failure prevents priority writes, and priority-write failure retains the saved verdict. Do not repeat the full lifecycle in each failure case. |
| R1, R5, AC4, AC7 | Existing handler test file and `TestCodexInspectionParticipationSurvivesAccountEdit` | Check saved verdict projection, absent-state normal default, and key-only persistence with unrelated keys intact. Add the verdict key to the existing full-edit preservation assertion. |
| R2–R3, R5, AC4, AC8 | Existing build/check entry | Reuse manual fixed-count, participation, scheduling, cancellation, UI, and native integration coverage. |

The adjustment-off lifecycle is necessary because the earlier pre-degraded fixture missed the defect.
Keep all four required phases in that one test. Do not add a second adjustment-on lifecycle or a switch-toggle sequence.
Existing adjustment-on stable-case assertions and code review cover the absence of stable-state priority synchronization.
The incomplete-confirmation case already exists in the opt-out test. Do not add another incomplete-confirmation scenario.
Use the handler's absent-state assertion with an actual degraded priority. Do not add a priority-value default matrix.
Reuse current notification assertions. Do not add notification scenarios or a separate priority-change scenario.
Use the fake sender and short injected recovery wait. Do not spend provider quota or wait 240 real seconds.
The existing native atomic Extra merge supplies persistence. Do not add restart, copy/export, or database-migration test matrices.

## Validation Command

Use the existing original tree and installed frontend dependencies:

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check
```

The entry applies the patches, attaches existing tests, runs focused tests and vet, and builds the fixed package outputs.
Use an individual development check only for a concrete failure. Reuse valid completed checks for unchanged inputs.
The page remains unchanged. Reuse the prior browser evidence instead of adding browser acceptance checks.
Do not run the existing individual suites again before or after a passing final build/check without a relevant change or failure.
Do not add race, coverage, full-suite, browser-engine, installer-restoration, or custom-priority matrices.

No container is needed for planning or the unchanged page.
If a concrete runtime check needs a container, use only `sub2api-codex-eval-test` and retain it afterward.
Preserve the pending container-retention spec edit from the preceding user instruction.

## Acceptance and Handoff

Read `prd.md`, `design.md`, and `research/verdict-state-contract.md` before implementation.
The approved task contracts supersede older spec statements that derive the verdict solely from actual priority.
Update `.trellis/spec/backend/native-inspection.md` and `.trellis/spec/backend/database-guidelines.md` with the verified contract before the work commit.
Keep all unrelated code, private configuration, real data, and the original upstream tree unchanged.
Do not deploy, push, or create a work commit during implementation without the applicable user authorization.
