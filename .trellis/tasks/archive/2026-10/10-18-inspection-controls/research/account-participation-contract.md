# Research: Account Inspection Participation Contract

- Query: Confirm the per-account Candy inspection switch against Sub2API v0.2.17 code. Trace persistence, list projection, account edits, cloning, shadows, inspection admission, later sends, priority updates, and Candy row behavior.
- Scope: internal
- Date: 2026-10-18

## Findings

### Recommended contract

Store `codex_eval_inspection_enabled` as a Boolean in the native account `Extra` map. Treat only Boolean `false` as opted out. Treat an absent value as enabled so existing accounts keep their current behavior. A dedicated endpoint should update only this key.

Use `PUT /api/v1/admin/accounts/:id/codex-eval-inspection` with `{"enabled": false}`. Return the native response envelope with `{account_id, enabled}`. Define the request field as `*bool` or validate it explicitly. A required plain Go `bool` rejects `false` under common `binding:"required"` validation because false is its zero value. Use native administrator route middleware. Validate that the target is an OpenAI OAuth or API-key account, including an OpenAI shadow account.

The handler should call `AdminService.UpdateAccountExtra` with only the participation key. This service calls the repository's key-level JSONB merge. The update is atomic and preserves unrelated `Extra` keys. A database migration is not needed. A post-write read is not needed because the response value is the exact Boolean just saved.

The proposed key is not removed by the current `UpdateAccountExtra` sanitizer. The service removes only other explicitly managed keys at `backend/internal/service/admin_account.go:928-951`. The key must be preserved by ordinary full account edits. `UpdateAccount` replaces `Extra` when the request contains a non-nil map. Add this key to its managed-key preservation list at `backend/internal/service/admin_account.go:664-711`. The normal admin update forwards `req.Extra` at `backend/internal/handler/admin/account_handler.go:1150-1195`.

The upstream `AccountService.Update` also replaces `Extra` when `req.Extra` is non-nil at `/tmp/sub2api-codex-eval-v0217-check/backend/internal/service/account_service.go:321-350`, but this is not an active admin entry point in the v0.2.17 server. `NewAccountService` appears only in the service Wire provider set at `backend/internal/service/wire.go:857-870`; the generated server graph does not construct or inject it. Active account routes use `AdminService` through `admin.ProvideAccountHandler`, shown in `backend/cmd/server/wire_gen.go:190-212`. Therefore the implementation owner only needs to preserve the key in `adminServiceImpl.UpdateAccount`. Do not widen this task to the legacy service unless a future caller is added. Its unit tests remain useful only if that service is later reactivated.

Bulk account updates already merge `Extra` keys. `BulkUpdateAccounts` forwards the partial map at `backend/internal/service/admin_account.go:954-967` and `:1114-1120`. The repository uses `COALESCE(extra, '{}'::jsonb) || payload` at `backend/internal/repository/account_repo.go:3131-3143`. An unrelated bulk edit therefore keeps this key without adding it to the bulk preservation list.

### Read and inspection flow

The current lite account DTO already carries `Extra`. `AccountFromServiceShallow` copies redacted `Extra` at `backend/internal/handler/dto/mappers.go:235-257`. `redactAccountManagedExtra` removes five named usage-state keys and passes other keys through at `:417-434`. `AccountListItemFromAccount` copies `Extra` at `:462-480`. The lite list path uses this projection at `backend/internal/handler/admin/account_handler.go:800-857`. This path needs no DTO change or separate regression test.

The current `InspectionAccount.Enabled` field means native `Schedulable`, not inspection participation. `inspectionAccount` fills it from `a.Schedulable` at `patches/0005-inspection.patch:786-790`. `inspectionEligible` uses it only for the existing `enabled_only` filter at `:275-294`. Add a separate participation field, such as `AutoInspectionEnabled`, and check it unconditionally in `inspectionEligible`. Do not reuse `Enabled`, `Schedulable`, or the global settings `Enabled` field.

Both list selection and fresh account reads pass through `inspectionAccounts`. The account list hook is at `patches/0005-inspection.patch:808-829`; the per-ID read hook is at `:831-837`. Map the participation value in `inspectionAccount` so scheduled rounds, immediate coordinator rounds, later test checks, and priority checks see the same account-level policy. The timer's global switch is checked separately in `startRun(automatic)` at `:226-258`. Manual Candy starts call `Start`, not the inspector, through `patches/0001-evaluation-core.patch:1972-1976`; leave their admission unchanged.

### Critical opt-out gap during multi-test runs

The current plan's later-admission check is not sufficient for a confirmation task with `Runs > 1`. `testPhase` rereads the account once, then calls `start` with the full remaining run count at `patches/0005-inspection.patch:515-549`. The Candy worker loops over every run internally at `patches/0001-evaluation-core.patch:50-71`. That loop checks only `auth.ctx.Err()` before each send. An account setting update does not cancel this context. Therefore, opt-out during the first confirmation request does not stop later requests in the same worker.

Preferred minimal fix: add an inspection-only pre-run admission callback to the internal Candy credential or start path. `testPhase` supplies the callback from `InspectionHooks.Get`. `runCandy` calls it before each request, including the first request. It checks context under the runtime lock, then invokes the callback without holding the global runtime mutex. A false eligibility result ends the worker before another send. Manual `Start` leaves the callback nil, so its configured run count is unchanged. Do not cancel the in-flight context when the endpoint saves false. The current request can finish, and the next loop admission prevents another request.

Alternative: have the coordinator start and await each Candy request with `Runs: 1`, reread the account before each subsequent request, and aggregate all run IDs and results. This makes each reread explicit, but it adds coordinator state and more run records. It is larger than the callback approach.

After opt-out stops a confirmation batch, the recorded test count must reflect completed requests. The batch is incomplete, so it cannot produce a recommendation, priority change, or notification. The integrated design selects the existing `unchanged` outcome for intentional opt-out, without an evaluation-error message.

The coordinator already rereads an account before priority update and checks `inspectionEligible` at `patches/0005-inspection.patch:447-465`. Adding the participation field to that predicate blocks a fresh opt-out read. There is a small read/write race if the switch changes after `Get` and before `UpdatePriority`. The current update hook has no conditional-write contract. The PRD requires a fresh reread, not an atomic database compare, so keep the existing hook unless acceptance requires closing that race.

### Metadata, groups, deletion, duplicate, and shadow behavior

`UpdateAccountExtra` is an existing key-level service API at `backend/internal/service/admin_service.go:91-95`. Its implementation sanitizes the update and calls `AccountRepository.UpdateExtra` at `backend/internal/service/admin_account.go:928-951`. The repository atomically merges the JSONB object at `backend/internal/repository/account_repo.go:2733-2807`, specifically `:2763` and `:2770-2807`. Store a JSON Boolean, not a string.

Group membership is not part of `Extra`. The inventory reads accounts by group, while the saved key stays on the native account row. One account in several groups therefore has one setting. Reload and process restart also read the same account `Extra` value.

Deleting an account removes it from the active account inventory. A newly created account with no participation key defaults to enabled. Parent deletion also deletes its shadow through `DeleteAccount` at `backend/internal/service/admin_account.go:1285-1299`.

`DuplicateAccount` clones `Extra` and removes only named runtime keys at `backend/internal/service/admin_account.go:150-165` and `:254-318`. Data export and import also carry `Extra` at `backend/internal/handler/admin/account_data.go:202-215` and `:436-453`. The integrated plan retains these native copy paths for this persistent configuration. Duplicate and backup/restore therefore copy the switch. Each copied account can then change independently. A fresh account without copied `Extra` uses the absent-value default.

Shadow creation builds its own `Extra` map and does not copy the parent's map at `backend/internal/service/admin_account.go:1440-1455`. The inspection mapper also uses the shadow's own account object, while parent enrichment only supplies a fallback plan at `patches/0005-inspection.patch:792-805`. Recommended behavior: the shadow has an independent switch, defaults enabled, and does not inherit later parent changes. Parent-linked inspection participation would be a different product contract.

### Candy row integration

The Candy column header is created by `credentialCard` at `patches/0006-evaluation-ui.patch:1110-1117`; the current shared header adds only the native `启用` column. The Candy card calls it with its columns at `:462-465`. `renderCandyRow` passes Candy-only markup to the shared `credentialRow` at `:602-623`. `credentialRow` currently places the native schedulable switch after shared row metadata at `:967-977`. Add a second Candy-only switch without changing other tabs.

Current CSS defines separate panel grid tracks at `patches/0006-evaluation-ui.patch:1833-1841`. Candy needs one additional fixed switch track and matching header cell. On narrow screens the list header is hidden and row metadata becomes a flex footer at `:2065-2083`; label both switches distinctly there. The outer `.row-main` click handler ignores label clicks at `:658-681`, so a native `<label class="switch">` avoids opening the history row. The delegated change handler currently branches only on `data-enable` and account selection at `:533-553`. Add a distinct data attribute for inspection participation. Do not route it through `setCredentialEnabled`, which writes `/schedulable` at `:370-385`.

Keep a per-account pending value so a background inventory refresh cannot replace an in-flight switch state. On save failure, restore the last saved value and render the error beside that switch. Do not reuse the header error summary.

### Check handoff

The check table in `../implement.md` defines the required cases and their owners.
Cover the new Boolean endpoint, default participation, key-only save, and changed full-edit preservation.
Cover initial opt-out and opt-out between worker sends in existing coordinator tests. Extend one existing manual-count case with opt-out.
Use the existing UI harness for independent switches, pending refresh, and local save failures.
Use one browser pass for column layout, switch interaction, and reload.
Native storage and unchanged copy, export/import, bulk-update, and DTO paths need no separate test matrices.

## Files Found

- `patches/0001-evaluation-core.patch`: native Candy worker loop, runtime start path, run identity, manual entry point, and cancellation.
- `patches/0005-inspection.patch`: inspection account projection, eligibility, round coordinator, per-phase admission, fresh account reread, priority update, settings routes, and native route handlers.
- `patches/0006-evaluation-ui.patch`: Candy list markup, keyed renderer, account schedulable toggle, event delegation, row click handling, and responsive grid.
- `patches/0007-sub2api-integration.patch`: administrator route registration for current Candy and inspection endpoints.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/service/admin_service.go`: native admin service contract, including `UpdateAccountExtra`.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/service/admin_account.go`: full, bulk, duplicate, create, delete, and shadow account semantics.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/repository/account_repo.go`: atomic `Extra` key merge for single and bulk account updates.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/handler/dto/mappers.go`: account redaction and full/lite DTO projections.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/handler/admin/account_handler.go`: account listing, full account update, and account handler conventions.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/handler/admin/grok_media_eligibility_handler.go`: current dedicated account-level GET/PUT handler pattern using `UpdateAccountExtra`.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/server/routes/admin.go`: native account route registration.
- `/tmp/sub2api-codex-eval-v0217-check/backend/internal/handler/admin/account_data.go`: account backup export/import with `Extra`.
- `tests/backend/internal/handler/admin/codex_inspection_handler_test.go`: current shadow parent-plan test and inspection mapper assertions.
- `tests/backend/internal/codexeval/inspection_test.go`: coordinator admission, priority, and notification tests.
- `tests/backend/internal/codexeval/runtime_test.go`: current Candy fixed-run and task lifecycle tests.
- `tests/ui/sub2api-dropdowns.cjs`: DOM, asynchronous API, row identity, and inspection settings harness.
- `.trellis/spec/backend/native-inspection.md`: current global settings, eligibility, confirmation, and manual-test contracts.
- `.trellis/spec/backend/runtime-installation.md`: current native account-list projection and manual evaluation contracts.
- `.trellis/spec/backend/database-guidelines.md`: native persistence conventions.

## Caveats / Not Found

- No external references are needed. The design uses existing Sub2API v0.2.17 APIs and repository patterns.
- The current worker checks admission only before a task. The integrated design now requires the inspection-only callback before each request.
- The chosen plan copies the persistent switch through duplicate-account creation and account backup/restore. Shadows keep independent values and do not inherit the parent switch.
- The integrated design uses the existing `unchanged` outcome for intentional opt-out and keeps completed test records.
- A fresh read before priority update is not atomic with the priority write. The existing native hook does not provide a conditional update by metadata value.
- This research read the original Sub2API source and repository patches only. It did not modify source, tests, task plans, or specs, and it ran no product checks.
