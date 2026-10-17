# U13 Shared Account Lists and Group Membership

## Scope

Inspection now uses the existing evaluation account card, list patching, native state conversion, and Candy record display.
The independent inspection result rows and their styles are removed.
Manual evaluation keeps its four entry points and account filters.
The backend scheduler and account-only native sender remain unchanged.

## Shared inspection display

`credentialCard`, `renderList`, `stateCredentials`, and `renderCandyRow` now accept the inspection context.
Their defaults preserve existing manual callers.
The inspection context reads every native OpenAI account, independently of the selected browser group.
It reads Candy records and active tasks through the existing evaluation state endpoint.
The shared conversion supplies account identity, subscription, priority, disabled state, restrictions, and active progress.

Inspection results match native account IDs.
The native account name remains authoritative instead of the inspection result name snapshot.
Each matching account retains its outcome, priority transition, test count, evaluation error, and notification error.
The shared Candy renderer supplies verdicts, metrics, record marks, tooltips, history, Markdown, and answer expansion.
Inspection answer elements have their own ID prefix.
The shared keyed list keeps row nodes during updates.

Inspection uses the existing type, plan, and enabled-only display filters.
These filters do not change autonomous backend membership.
The inspection view does not add account mutation, individual test starts, selection, or history clearing.
Its native enable switch shows account state without accepting changes.
Manual inspection start and stop remain in the existing settings footer.
Session reset clears inspection account records and stops its pending state request.

## Mixed-group investigation and reproduced defect

The native v0.2.17 query parameter is `group`, not `group_id`.
The existing page already uses `group`.
`AccountHandler.List` passes it through `AdminService.ListAccounts` to the repository account-group predicate.
The lite account response uses `dto.AccountFromServiceShallow`, which explicitly copies `GroupIDs`.
The JSON field is `group_ids`. Empty memberships can omit it.

The old page cached the account inventory after its first successful load.
Ordinary polling called `refreshCatalog`, which returned immediately when `catalogLoaded` was true.
Polling then refreshed only evaluation state, not native membership.
An account moved from group A to group B could remain on the open group A page.
That stale row could also produce a group A target.
This is a reproduced stale-membership defect. It does not prove the user's exact production sequence.

The focused fake API returns copied JSON objects and has exclusive, shared, and ungrouped accounts.
After the first group A load, it moves account 11 to group B.
With the repair, an ordinary refresh removes account 11 from all four manual lists and their targets.
Restoring only the old early-return branch in memory makes the same test fail.
The old result is `[11, 33]`. The required result is `[33]`.
This control check did not modify repository files.

## Group repair

Ordinary polling now refreshes the native inventory for the selected group.
Raw model discovery remains cached when the enabled account IDs are unchanged.
Returned `group_ids` membership restricts the manual inventory, state conversion, and target builder.
Ungrouped accounts do not enter a selected-group list.
Shared accounts can appear in each group to which they belong.

Catalog cache state now records its group.
Each request captures its generation and selected group before committing results.
A late old-group response cannot replace a newer group's inventory or models.
Manual and inspection inventory remain separate.
The common paginated native account reader serves both scopes without combining their results.

## Independent review repairs

The first review found that a failed native inventory refresh still permitted starts through retained rows.
Manual inventory now tracks readiness, group, and read error separately from model discovery.
An incomplete or failed inventory refresh blocks row buttons, batch buttons, target creation, and the actual start-request entry.
The native catalog error remains visible until recovery or dismissal.
A successful inventory refresh restores starts with current memberships.
A raw-model error does not invalidate a successfully refreshed inventory or block a custom literal model.

The review also found that inspection polls did not refresh model suggestions after account changes.
Each poll now compares the sorted enabled account ID set with the successful discovery set.
An unchanged set reuses model suggestions.
A changed set refreshes suggestions through the existing discovery function.
Disabled and removed accounts no longer supply suggestions.
Failed discovery keeps its actual error and retries on a later poll.
Successful recovery clears that error.
Catalog responses retain newer custom-model input and its editor state.
The catalog error now sits outside the folded configuration body.

Focused regressions passed for both failures, recovery, preserved custom inputs, and model cache reuse.

The next review found a separate evaluation-state failure path.
A successful account refresh could remove a member, while a failed evaluation-state read left the old rendered credentials intact.
Targets now read membership directly from the current native account inventory.
The page projects fresh native accounts through the shared converter before requesting evaluation state.
It uses the last successful evaluation-state response for existing records until a new response succeeds.
Account identity, membership, priority, and native restrictions therefore remain current when record reads fail.
The visible record-read error remains separate, and session reset clears the retained response.

A focused regression moves account 11 out of A and updates account 33's native name and priority.
It then returns the new account inventory while failing the evaluation-state request.
All four account lists and batch inputs retain only account 33.
The actual row and four batch start requests submit only account 33 under A.
The previous Candy record remains available, and the read error remains visible until a successful state refresh.

## Changed files

- `native/backend/internal/codexeval/web/catalog.js`: common pagination, live manual membership, group-scoped commits, and separate inspection inventory.
- `web/app.js`: shared conversion and filters, group-restricted targets, shared inspection layout, and session reset.
- `web/components.js`: shared read-only card and account row context, with attached inspection summaries.
- `web/candy.js`: shared row/history/click handling and scoped tooltip lookup.
- `web/inspection.js`: native account and evaluation-state loading through the shared list.
- `web/ui.html`: replaced independent inspection result markup with the shared account card container.
- `web/style.css`: removed independent result styles and added the shared account summary layout.
- `tests/ui/sub2api-dropdowns.cjs`: shared list/history and two-group membership regression checks.
- `tests/installer/native_inspection_fixture.py`: idempotent two-group browser fixture and a manifest with actual IDs.

## Focused checks

The settings and account UI checks passed.
They cover all five model controls, U11/U12 settings behavior, four manual submission paths, and disabled/unavailable/busy account admission.
They also cover membership moves, shared accounts, pagination, late group responses, inspection isolation, record interactions, and session reset.
The old early-return control check failed at the expected moved-account assertion.
JavaScript syntax, Python fixture syntax, and diff whitespace checks passed.
Reuse unchanged backend, settings, and installer evidence.
No broad Go suite or race requirement was added.

The final U13 build passed through `scripts/build-sub2api.py` after all three review repairs.
This candidate includes the evaluation-state failure repair and its focused regression.
The frontend passed three locale tests, TypeScript checks, and its Vite build.
Packaging passed through `scripts/package-sub2api.py`.
The executable is a stripped, static Linux ELF for amd64.

| Output | Size |
| --- | ---: |
| `dist/build/sub2api-codex-eval` | 125595808 bytes |
| `dist/build/sub2api-codex-eval.tar.gz` | 38700060 bytes |

The archive has seven entries: the executable, the installer, and license files.
It contains no development or test files.
Build source: `/tmp/sub2api-codex-eval-v0217-check`.
Build log: `/tmp/sub2api-U13-build.log`.

## Retained-container browser fixture

The fixture helper creates or reuses native OpenAI standard groups `U13 Group A` and `U13 Group B`.
It creates or updates `u13-a-only` in A, `u13-b-only` in B, and `u13-shared` in both groups.
These accounts use the existing local fake supplier at `http://127.0.0.1:18081`.
The helper prints `Group fixture manifest: /tmp/sub2api-inspection-fixture-*/u13-groups.json`.
That file contains the origin and actual group/account IDs, without credentials.
It uses mode 600.

For the browser check, move `a_only` to B through the native account API while its A page remains open.
Wait for ordinary polling and check that A removes it from rows and submitted targets.
Restore its original A membership after the check.
The pure UI fixture uses group IDs 101 and 202 and account IDs 11, 22, and 33.
Actual native IDs are available only after the main session runs the helper.

The main session owns container loading, actual browser verification, and independent review.
These checks remain pending at this handoff.
No Docker, production, paid-request, staging, commit, or cleanup action ran in this implementation turn.
