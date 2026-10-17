# Independent Inspection Verdict

## Change Boundary

One implementation owner must maintain the coordinator, native metadata integration, and focused tests together.
Production changes belong to `patches/0005-inspection.patch` and `patches/0007-sub2api-integration.patch`.
Tests remain under `tests/backend/`. The frontend and manual runtime need no changes.

| Applied source | Change |
| --- | --- |
| `backend/internal/codexeval/inspection.go` | Add the verdict type and persistence hook. Classify from verdicts and separate optional priority writes. |
| `backend/internal/handler/admin/codex_inspection_handler.go` | Read verdict metadata and connect the existing native account Extra writer. |
| `backend/internal/service/admin_account.go` | Preserve the verdict key during unrelated full account edits. |

## Native Metadata Contract

Store `extra.codex_eval_inspection_verdict` as `normal` or `degraded`.
Use `InspectionVerdict` with normal and degraded values. Normal is the default state.
Add `Verdict InspectionVerdict` to `InspectionAccount`.
Add `PersistVerdict func(context.Context, int64, InspectionVerdict) error` to `InspectionHooks`.

The handler reads the saved verdict. An account without a saved verdict starts as normal.
Persist only the verdict key through `AdminService.UpdateAccountExtra` and its existing atomic JSONB merge.
Keep unrelated Extra keys, including `codex_eval_inspection_enabled`.
If the saved verdict already matches the valid outcome, do not write Extra again.
Add the verdict key to the existing full-edit preservation list in patch `0007`.
Keep existing copy, export/import, and shadow-account storage behavior.

## Initial State

All accounts without a saved verdict start as normal. Actual priority does not affect this default.
An account at the former degraded priority must complete normal-to-degraded confirmation before its independent state becomes degraded.
Custom priorities do not create a separate unknown state or bypass confirmation.
Save a verdict when a confirmed transition changes it. A stable default normal state needs no initialization write.
The user accepts this breaking change. Do not add priority-based inference or migration.

## Inspection Flow

1. Run one first test through the existing inspection-only admission path.
2. Reread the account after the result. Capture its actual priority as the confirmation baseline.
3. Read the verdict from fresh metadata, with normal as the default.
4. For a degraded verdict and wrong answer, finish after one test.
5. For a normal verdict and correct answer, finish after one test.
6. For a possible transition, retain the configured confirmation count and existing 240-second recovery wait.
7. Require the complete configured count of valid, matching answers before a transition.
8. Apply the existing participation and cancellation checks before the state commit.
9. Persist the verdict. Only a confirmed verdict change can enter the existing optional priority-update path.

Use the actual priority for historical results and the fresh-account baseline. It does not classify verdicts or select confirmations.
Preserve participation checks before each send. Intentional opt-out retains completed results without a synthetic evaluation error.
Errors, mixed answers, cancellation, or incomplete confirmation work cannot create a new verdict.
Check cancellation before the save. A successful save remains durable if cancellation occurs after it completes.
Keep existing final priority and busy-state checks on the priority-update path. Do not add them as verdict-save conditions.

## Optional Priority Writes

After a confirmed verdict change, check `AdjustPriority` in the existing round settings snapshot.
If disabled, finish after the verdict save. If enabled, use the existing guarded priority-update path with the new verdict's target.
A stable verdict never enters this path. Turning the switch on does not synchronize an unchanged verdict.
Retain fresh-account and baseline checks. A custom current priority does not block the configured target write.
An already matching actual priority needs no write.

Verdict persistence and priority updates use separate native operations.
A failed verdict save prevents the priority write.
A failed priority write retains the durable verdict and records the existing backend error.
Do not add synchronization, automatic retries, a cross-field transaction, or priority rollback.

Keep `InspectionResult.Priority` as historical data. Reuse the existing target, outcome, and error fields.
Do not add public result fields or a page status column.

## Existing Notifications

Keep notifications on the existing successful-priority-update path. Keep notification text, configuration, and backend errors unchanged.
Do not add verdict-only notifications or restore the removed page prompts.

## Delivery and Recovery

This repair needs no database migration, new route, new table, or new scheduler.
The existing native account store retains verdicts across groups and process restarts.
Use the existing seven-patch build entry and fixed output paths.
