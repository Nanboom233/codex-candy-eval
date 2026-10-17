# Research: verdict-state-contract

- Query: Separate the automatic Candy inspection verdict from the native account priority, with one-test stability after a confirmed degraded result even when automatic priority adjustment is disabled.
- Scope: internal
- Date: 2026-10-18

## Findings

### Current defect and call path

The current coordinator has no persisted verdict field. `InspectionAccount` exposes only the native `Priority`, and `InspectionHooks` has no metadata-write hook (`patches/0005-inspection.patch:18-36`). The `run` path performs one first test, rereads the account, and selects confirmation groups by the reread native priority (`patches/0005-inspection.patch:387-416`). This makes the native priority both the observation and the verdict.

The update loop skips all transition side effects when `AdjustPriority` is false (`patches/0005-inspection.patch:459-490`). The notification loop consumes only the accounts that received a successful priority update (`patches/0005-inspection.patch:492-513`). Therefore a normal-priority account that is confirmed wrong while adjustment is off has no durable degraded state, and the next round again follows the normal-priority confirmation path.

The existing checked Candy path is the right admission boundary. `runCandyChecked` calls its callback before every send (`patches/0001-evaluation-core.patch:49-77`); `startChecked` wires that callback only for inspection calls and leaves public/manual `Start` unchanged (`patches/0001-evaluation-core.patch:1979-1987`, `patches/0001-evaluation-core.patch:2078-2087`). The callback in `testPhase` currently checks participation and, for confirmation sends, compares the fresh native priority with the phase account priority (`patches/0005-inspection.patch:531-559`). That actual-priority baseline check must remain, but it must not decide the logical verdict.

### Minimal state contract

Use one string value in native account `extra`:

```text
key: codex_eval_inspection_verdict
values: "normal" | "degraded"
default: normal
```

Use a two-state internal type:

```go
type InspectionVerdict string

const (
    InspectionVerdictNormal   InspectionVerdict = "normal"
    InspectionVerdictDegraded InspectionVerdict = "degraded"
)
```

Add `Verdict InspectionVerdict` to `InspectionAccount` and add one coordinator hook:

```go
PersistVerdict func(context.Context, int64, InspectionVerdict) error
```

The handler maps the account `Extra` value to the type. `PersistVerdict` calls the existing `AdminService.UpdateAccountExtra` with only the verdict key. No new route, table, file, scheduler, or public DTO is needed.

The current native service already provides the required storage boundary. `UpdateAccountExtra` forwards key-only updates to the account repository (`/tmp/inspection-controls-implementation/backend/internal/service/admin_account.go:927-952`). The repository performs an atomic JSONB merge with `COALESCE(extra, '{}')::jsonb || $1::jsonb` in one SQL update (`/tmp/inspection-controls-implementation/backend/internal/repository/account_repo.go:2733-2774`). Add the new key to the managed-key preservation list used by full account edits (`patches/0007-sub2api-integration.patch:99-102`). Do not add it to the duplicate discard list: the existing Extra copy/export path should copy inspection configuration, while each shadow account keeps its own Extra value.

### Verdict default

Read the saved normal or degraded state. An account without a saved verdict starts as normal, regardless of actual priority.
The user accepts this breaking change. Do not infer verdicts from old priorities or add a migration.
Custom priorities do not create an unknown state or bypass confirmations.
Save a verdict only after a complete, consistent confirmed transition. A stable default normal state needs no initialization write.
Errors, incomplete work, mixed answers, cancellation, or opt-out cannot create a transition.

### Run and transition algorithm

1. Keep the first test at one send for every eligible account.
2. After it completes, call `Get` and capture `priorityBaseline` from the current native priority. Read the verdict, with normal as default.
3. A degraded verdict with a wrong first answer is stable. Stop after one test, regardless of `AdjustPriority`.
4. A normal verdict with a correct first answer is stable. Do not start confirmation.
5. A normal-to-degraded or degraded-to-normal answer starts the existing confirmation count and recovery wait. Confirmation admission checks current participation/eligibility and requires `current.Priority == priorityBaseline`; it does not compare native priority with the logical verdict or the account-list snapshot.
6. After a complete matching transition, persist the new verdict first. A persistence error prevents priority writes and notifications for that transition.
7. Only after a confirmed verdict change, check the adjustment switch and use the existing guarded priority-update path when enabled.
   Retain existing final eligibility, busy-state, and priority-baseline checks on that priority-update path.
   Do not add final priority or busy-state checks as verdict-save conditions. A failed priority write retains the saved verdict without automatic retry.
8. If adjustment is disabled, never call `UpdatePriority`; the verdict write still occurs.

A stable verdict does not trigger a priority write. Enabling adjustment alone does not synchronize actual priority.
This follows the user's clarified transition-only scope.

### Notification boundary

Keep notifications on the existing successful-priority-update path. The clarified scope does not change notification rules, text, or error handling.

No new result column is required. Preserve the historical `InspectionResult.Priority` and existing target/outcome fields; only the transition bookkeeping changes.

### Existing test owners

`TestCodexInspectionConfirmationAndBoundaries` already covers transitions, stable outcomes, mixed answers, provider errors, priorities, and notification failures.
`TestCodexInspectionCancelAndSchedule` covers cancellation. `TestCodexInspectionParticipationStopsLaterSends` covers opt-out and incomplete confirmation.
`TestCodexInspectionDegradedWrongOnce` has a priority-based matrix that must become the single adjustment-off lifecycle regression.
The handler and service test files already contain metadata and full-edit fixtures.
Use the check table in `implement.md` as the complete validation scope. This report adds no test requirements.

## Related specs

- `.trellis/spec/backend/native-inspection.md:136-173` already requires native participation storage, fresh actual-priority baselines, independent test counts, and no transition from incomplete work. It needs a small verdict-state update during the implementation phase.
- `.trellis/spec/backend/database-guidelines.md:25-39` already requires key-only atomic account Extra updates and full-edit preservation. Add the verdict key to that contract during spec update.
- `.trellis/tasks/10-18-inspection-verdict-state/prd.md` records the approved goal and the missing initial-normal/adjustment-off regression.

## Caveats / Not Found

- The previous implementation has no verdict key, no `PersistVerdict` hook, and no persisted inspection transition state. The exact storage key and type above are a new contract for this task.
- Actual priority no longer classifies inspection state. Accounts without saved verdicts start as normal under the user's accepted breaking change.
- Native priority and Extra are separate writes. The recommended ordering persists the logical verdict before an optional priority update, so a failed priority write does not erase a confirmed verdict. A cross-field transaction would be a larger change and is not required for the stated behavior.
