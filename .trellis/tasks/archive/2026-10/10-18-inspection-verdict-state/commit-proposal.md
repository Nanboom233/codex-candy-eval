# Work Commit Proposal

The final build/check and Trellis review passed. Git has `commit.gpgSign=true`.
Use the configured signing policy for both commits. Do not push.
The user must approve these exact work commits before creation.

## 1. Independent Inspection Verdict

```text
fix(codex-eval)!: separate inspection verdicts from priority

Persist confirmed verdicts before optional priority writes. Keep stable
rounds to one test.

BREAKING CHANGE: Accounts without a saved verdict start as normal.
Actual priority no longer determines the inspection verdict.
```

Files:

- `patches/0005-inspection.patch`
- `patches/0007-sub2api-integration.patch`
- `tests/backend/internal/codexeval/inspection_test.go`
- `tests/backend/internal/handler/admin/codex_inspection_handler_test.go`
- `tests/backend/internal/service/codex_inspection_settings_test.go`
- `.trellis/spec/backend/native-inspection.md`
- `.trellis/spec/backend/database-guidelines.md`

## 2. Retained Test Container

This separate change records the user's preceding container-retention instruction.

```text
docs: retain the reusable Sub2API test container
```

File:

- `.trellis/spec/backend/minimal-linux-delivery.md`

After both approved work commits, use Trellis for task archival and the session record.
Keep task and workspace artifacts out of these work commits.
