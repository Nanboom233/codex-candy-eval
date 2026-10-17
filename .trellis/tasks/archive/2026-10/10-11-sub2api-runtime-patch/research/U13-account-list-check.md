# U13 Independent Check

## Findings (fixed)

1. A failed or incomplete manual account refresh now blocks row and batch controls, target creation, and the `/start` request entry. Recovery reloads current memberships and restores starts. A model-only discovery error does not mark the successfully loaded account inventory unavailable, so a literal model remains usable. The focused UI test covers failure, recovery, and model-only failure.

2. Inspection model discovery is now keyed by the sorted enabled account IDs. Unchanged membership reuses the successful model list. Changed membership refreshes it and skips disabled accounts. A discovery error stays visible and retries on the next poll. A successful retry clears the error, and an in-progress custom model edit is preserved. The focused UI test covers these cases.

3. Target creation now checks IDs against the current native account inventory. After a successful inventory refresh, the page projects current account identity, group membership, priority, restrictions, and the last successful evaluation records before requesting new evaluation state. If that state request fails, stale group members are removed, existing records remain available, and the record error stays visible. A successful retry replaces the retained state and clears the error. The new regression covers all four lists, batch targets, row and batch starts, retained Candy history, and recovery.

## Findings (not fixed)

None.

## Verification

- UI test: pass after the final source delta. `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`
- JavaScript syntax: pass for `app.js` and `catalog.js` with `node --check`; inspection syntax passed before this delta, which did not change that file.
- Fixture syntax: pass with `python3 -m py_compile tests/installer/native_inspection_fixture.py`.
- Diff whitespace: pass with `git diff --check` after the final source delta.
- Lint/type-check: no separate lint or type-check is configured for these embedded JavaScript page files. The implementation report records the prior build and frontend TypeScript checks; a build including this final delta remains with the implementation owner.
- Browser/container verification: not run by this checker. The main session owns verification in the retained `sub2api-codex-eval-test` container.
