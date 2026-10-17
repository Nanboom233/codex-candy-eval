# Execution Plan

## Implementation Gate

- [x] Obtain user approval of the latest planning summary.
- [x] Start this task with `python3 .trellis/scripts/task.py start 10-18-inspection-controls`.
- [x] Dispatch one Trellis implementation owner with the task artifacts and curated context.

## Ordered Work

1. Use the completed research probes as evidence. Convert needed assertions into existing tests instead of retaining parallel regression suites.
2. Repair both readiness resets for a valid same-group refresh. Reread current priority before confirmation selection and later sends.
3. Add account participation persistence, its administrator update endpoint, and inspection admission checks through the existing native account owners.
   Add the Candy row column and independent switches. Add an inspection-only pre-send callback for each multi-test Candy worker iteration.
   Preserve full-update metadata and manual entry points. Treat intentional opt-out as unchanged.
   Update existing `InspectionAccount` fixtures with participation values so their original admission checks remain valid.
4. Remove inspection runtime and notification error prompts, their aggregation, and their DOM writers. Do not relocate them.
   Keep account test records, settings save errors, account switch save errors, and the `下次巡检：` timestamp.
5. Extend existing tests only for the gaps in the check table below. Format changed Go source and keep production patches valid.
6. Run the existing build/check entry once for the final candidate, using the original source and installed dependencies.
7. Inspect the candidate in one browser pass through the retained fixture. Check the added column at normal and narrow widths.
8. Dispatch the Trellis check owner with the existing results. Resolve material findings without repeating valid checks.
9. Update affected specs with verified behavior and record validation results in this task.
10. Present the verified work diff and exact commit proposal under the project commit workflow.

## Checks

The research probes already confirmed the readiness flash and stale-priority extra calls. Do not rerun them as acceptance gates.
The retained Go probe asserts the old three-call defect. Transfer the corrected one-call assertion into the existing coordinator tests.

Use the existing original tree at `/tmp/sub2api-codex-eval-v0217-check`.
The source directory and its installed `frontend/node_modules` are available.

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check
```

The entry runs existing UI, installer-version, focused backend, and changed-package vet checks before building and packaging.
It already includes the affected test files. Do not schedule a separate passing run of each file before this entry.
Use an individual check for a concrete development failure when needed. Reuse valid unchanged installer and integration evidence.

For browser checks, use `python3 tests/installer/native_inspection_fixture.py` and `sub2api-codex-eval-test`.
Inspect the added column, distinct switch labels, switch interaction, reload, and visible refresh stability.
Confirm the removed prompt leaves no empty row and `下次巡检：` still shows its scheduled time.
Backend tests own sender counts. The DOM harness owns synthetic error and delayed-response cases.
Do not repeat those scenarios in the browser. Use the local fake supplier and exclude private values from reports.

## Check Ownership

This table is the complete task-specific check plan. Research suggestions do not add gates beyond this table.

| Requirement | Owner | Necessary change |
| --- | --- | --- |
| R1, AC1–AC2 | Existing UI harness | Add one delayed successful refresh case. Assert stable button, row, selection, and expansion. Reuse membership, failure, and late-group checks. |
| R2, AC3 | Existing coordinator tests | Add stable degraded-wrong and stale-list/degraded-current cases. Assert one sender call/result, unchanged priority, and no notification with adjustment on and off. |
| R2, AC4–AC5 | Existing coordinator/runtime tests | Reuse confirmation, recovery-wait, cancellation, and fixed-count checks. Extend one manual case to use an opted-out account. |
| R3, AC6–AC7 | Existing admin/service tests | Check Boolean false binding, required-field validation, absent-value default, key-only save, and full-edit preservation. |
| R3, AC7 | Existing coordinator tests | Check initial opt-out and opt-out during a blocked confirmation. Count completed sends and assert no later send, priority write, notification, or synthetic error. |
| R3, AC6–AC7 | Existing UI harness | Check separate switches, account/global isolation, pending refresh, and failed-save rollback with a row-local error. |
| R4, AC9 | Existing UI harness | Replace old summary assertions. Keep result/settings errors, remove runtime/notification prompts, and retain the next inspection label/time. |
| AC8 | Existing build/check entry and one browser pass | Verify the patched candidate and the visible controls described above. |

Store participation in native account metadata. That existing store supplies cross-group and restart persistence without new lifecycle tests.
Keep existing copy, export/import, DTO, bulk-update, and account-type behavior. Do not add separate tests for those unchanged paths.
Do not add custom-priority, repeated-round, browser-engine, race, coverage, full-suite, or installer-restoration matrices.
Repeat a passed check only after a relevant change or failure. The check owner reviews valid evidence before running anything again.

## Implementation Context

Read these research reports before editing their owners:

- `research/inspection-count-contract.md`: current-priority confirmation rules and executed call counts.
- `research/account-participation-contract.md`: account endpoint, metadata persistence, worker callback, and native copy behavior.
- `research/refresh-error-contract.md`: readiness reproduction, runtime error prompt removal, retained next-run time, and shared row layout.
