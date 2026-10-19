# Sub2API account model storage

## Result

Account model restrictions use `accounts.credentials.model_mapping` in PostgreSQL.
The frontend combines its whitelist and model mapping controls into this one JSON object.
The whitelist is not a separate database list.
The upstream model catalog is a separate discovery result.

## Evidence scope

This review reads the original source at `/tmp/sub2api-codex-eval-v0217-check` and the current production patches.
The existing source archive is `/tmp/sub2api-v0.2.17.tar.gz`, with root `sub2api-0.2.17/`.
The source has a stale `backend/cmd/server/VERSION` value of `0.2.16`.
The archived runtime task records that upstream issue in [its integration report](../../archive/2026-10/10-11-sub2api-runtime-patch/research/sub2api-v0.2.17-integration.md).
The release workflow selects `v0.2.17` at `.github/workflows/release.yml:19`.

Evidence covers source-defined storage, serialization, reads, and update behavior.
Live database rows were not inspected. No Docker containers were running during this review.

## Stored fields

Source anchors in this report are relative to the original source root unless they name a production patch.

| Object | Stored location | Evidence |
| --- | --- | --- |
| Account credentials | `accounts.credentials`, JSONB | `backend/migrations/001_init.sql:60`, `backend/ent/schema/account.go:79` |
| Account metadata | `accounts.extra`, JSONB | `backend/ent/schema/account.go:85` |
| Whitelist and normal model mappings | `credentials.model_mapping` | `frontend/src/composables/useModelWhitelist.ts:597`, `frontend/src/components/account/EditAccountModal.vue:4260` |
| Compact-only model mappings | `credentials.compact_model_mapping` | `frontend/src/components/account/EditAccountModal.vue:4274`, `backend/internal/service/account.go:1023` |
| Automatic passthrough | `extra.openai_passthrough`, with the legacy `openai_oauth_passthrough` fallback | `backend/internal/service/account.go:2125` |
| Inspection verdict | `extra.codex_eval_inspection_verdict` | `patches/0005-inspection.patch:899` |
| Global inspection settings | `settings.key` and `settings.value` | `backend/ent/schema/setting.go:35`, `patches/0005-inspection.patch:996` |

The settings value column uses PostgreSQL `text`.
The inspection service stores each field as JSON text under `codex_eval_inspection_` plus the field name.
It uses the existing repository's `GetMultiple` and `SetMultiple` methods.

## Frontend serialization

`buildModelMappingObject` writes each selected whitelist model as `model -> model`.
It then adds explicit mappings as `from -> to`.
An explicit mapping with the same source name replaces the identity entry.
The following example illustrates the stored format. It is not a database row.

```json
{
  "model_mapping": {
    "gpt-6.1-sol": "gpt-6.1-sol",
    "gpt-6-sol": "gpt-6-sol",
    "custom-alias": "gpt-6.1-sol"
  }
}
```

`splitModelMappingObject` reads equal keys and values into `allowedModels`.
It reads unequal keys and values into `modelMappings`.
See `frontend/src/composables/useModelWhitelist.ts:571`.

The OpenAI edit form uses the combined builder.
If it returns no entries, the form deletes `credentials.model_mapping`.
The form retains existing mappings while automatic passthrough is on.
See `frontend/src/components/account/EditAccountModal.vue:4247` and `:4260`.

## Save and read path

1. The account handler passes `credentials` and `extra` to `AdminService.UpdateAccount`.
   See `backend/internal/handler/admin/account_handler.go:1179`.
2. The service reads the account and applies the incoming credential object.
   See `backend/internal/service/admin_account.go:577` and `:648`.
3. `MergePreservingSensitiveCreds` copies incoming fields and retains absent sensitive credential fields.
   It does not retain absent ordinary credential fields.
   See `backend/internal/service/account_credentials_redact.go:38`.
4. The repository writes the credential object through `SetCredentials`.
   It writes account metadata through `SetExtra`.
   See `backend/internal/repository/account_repo.go:541`.
5. The update includes a scheduler outbox event in its transaction.
   After its own transaction commits, it refreshes the account's scheduler snapshot.
   See `backend/internal/repository/account_repo.go:499` and `:513`.
6. `GetByID` reads the account entity. `accountEntityToService` copies the stored credential and metadata objects.
   See `backend/internal/repository/account_repo.go:265` and `:3615`.

Submitting only a new `model_mapping` through the generic credential update can remove other ordinary credential fields.
The lower-level `UpdateCredentials` method also replaces the complete credential JSONB object.
See `backend/internal/repository/account_repo.go:843` and `:869`.

`UpdateAccountExtra` uses a separate key-level merge.
The repository executes `COALESCE(extra, '{}'::jsonb) || $1::jsonb`.
This path does not update `credentials.model_mapping`.
See `backend/internal/service/admin_account.go:928` and `backend/internal/repository/account_repo.go:2733`.

## Runtime model availability

`GetModelMapping` reads `Credentials["model_mapping"]` and keeps string-valued entries.
See `backend/internal/service/account.go:624` and `:663`.

The scheduler checks `supportsOpenAIRequestModel` before it accepts an account.
That function applies the subscription model retirement rule, then the account mapping rule.
See `backend/internal/service/openai_account_scheduler.go:1815` and `backend/internal/service/openai_model_retirement.go:86`.

The mapping rule has these behaviors:

| Condition | Current behavior | Evidence |
| --- | --- | --- |
| Exact request-side key exists | The mapping allows that name | `backend/internal/service/account.go:879` |
| A wildcard key matches | The mapping allows that name | `backend/internal/service/account.go:886` |
| No mapping entries remain | Native default admission applies | `backend/internal/service/account.go:928` |
| Automatic passthrough is on | The mapping restriction does not apply | `backend/internal/service/account.go:925` |

For OpenAI API-key accounts, an empty mapping permits all names at the mapping layer.
For OpenAI OAuth accounts, the empty-mapping rule excludes known foreign model families.
See `backend/internal/service/openai_model_mapping.go:68`.
The separate subscription retirement rule still applies.

These rules establish three concrete deletion limits.
Removing the final entry returns the account to native default admission.
Removing an exact entry leaves a matching wildcard effective.
Removing entries does not restrict an account while automatic passthrough remains on.

Request-side keys and upstream targets have different meanings.
Removing `gpt-6.1-sol` does not remove a separate `custom-alias -> gpt-6.1-sol` entry.
That alias remains a separate request-side name under the existing mapping rules.

## Catalogs and shadow accounts

`FetchOpenAIModelsList` gets the raw upstream catalog through the shared model discovery cache.
The cache stores entries in a process-memory map.
See `backend/internal/service/openai_models_list.go:19` and `backend/internal/service/openai_codex_models_service.go:1627`.

The normal account test picker projects the catalog through account mappings after discovery.
See `backend/internal/service/account_test_service.go:183` and `backend/internal/service/openai_models_list.go:179`.
The resulting picker list is not a separate saved account whitelist.

Shadow accounts store their own model mappings in their own credential JSONB object.
Their allowed credential keys are `model_mapping` and `compact_model_mapping`.
Authentication credentials come from the parent account.
See `backend/internal/service/account_credentials_persistence.go:35` and `backend/internal/service/credential_shadow.go:13`.

The evaluation model picker uses the raw upstream catalog.
The evaluation sender uses the native account-test transport and bypasses ordinary gateway scheduling.
See `patches/0004-admin-channel.patch:176` and `:213`.
An account model restriction therefore does not require removing the existing recovery-test entry point.

## Implications for this task

The user confirmed that this task targets identity whitelist entries, with the format `model -> model`.
Use the saved account mapping as the source for those changes.
Do not treat the upstream catalog or frontend `allowedModels` array as another stored list.
Record the removed identity entries before removal, so recovery can restore them.
Keep authentication credentials out of recovery metadata.
Preserve other credential fields and use the native scheduler update path.

The existing mapping format has no exclusion entry in the inspected admission path.
The user proposed `"blacklist": "blacklist"` as a placeholder after removal empties the whitelist.
The inspected native code retains this string-valued entry and treats the resulting mapping as a nonempty whitelist.
Real model names do not match that exact key. A literal request for `blacklist` would match it.
This is a normal mapping entry, not a native deny keyword.
The placeholder avoids a new scheduler block rule and preserves unrelated native empty-mapping behavior.
Wildcard exclusions, aliases, automatic passthrough, and compact-only mappings remain outside this task's confirmed scope.

## Validation

This was a source investigation. No product code, database rows, or running services changed.
Existing tests corroborate empty-mapping and wildcard admission in `backend/internal/service/account_wildcard_test.go:158`.
The bulk-edit test expects an empty mapping to restore all-model behavior in `frontend/src/components/account/__tests__/BulkEditAccountModal.spec.ts:152`.
These tests were read, not executed during this investigation.
