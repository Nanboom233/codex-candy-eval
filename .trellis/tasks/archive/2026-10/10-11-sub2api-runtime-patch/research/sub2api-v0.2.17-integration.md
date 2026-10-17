# Research: Sub2API v0.2.17 Integration

- Query: Check the v0.2.17 source changes that affect the native account evaluation integration, inspection settings, startup and shutdown, frontend patch, and response-model audit.
- Scope: internal
- Date: 2026-10-16

## Findings

The current target is the supplied Sub2API v0.2.17 source at `/tmp/sub2api-codex-eval-v0217-wi8Ct1/sub2api-0.2.17`. The prior source comparison found that the existing patch targets are unchanged from v0.2.16. The native account-test, model-list, settings-route, priority, shutdown, and frontend URL/sidebar integration needs no version-specific hunk rebase. The version-specific behavior found in v0.2.17 is OpenAI model-retirement filtering. The archive also retains a stale fallback version value.

| Area | v0.2.16 evidence | v0.2.17 source anchor | Finding |
| --- | --- | --- | --- |
| Account test | `research/core-third-evidence.md:47` records `account_handler.go:1289-1318` and `account_test_service.go:356-427`. | `backend/internal/handler/admin/account_handler.go:1289-1318`; `backend/internal/service/account_test_service.go:356-427` | `/test` still resolves the account by ID and dispatches by its platform/protocol/type. The ordinary handler still calls `RecoverAccountAfterSuccessfulTest` after success. The evaluation endpoint must remain separate to keep its approved no-recovery behavior. |
| Raw model discovery | `research/fresh-origin-port-map.md:21-22` records the native test picker and OpenAI sender boundaries. | `backend/internal/service/openai_models_list.go:16-20`; `backend/internal/service/account_test_service.go:180-220`; current overlay `native/backend/internal/service/account_codex_eval_service.go:19-43` | The shared model-list method returns the raw upstream catalog. The ordinary account-test picker projects account mappings at `account_test_service.go:191-202`. The evaluation method calls `FetchOpenAIModelsList` directly and returns its model entries, so it avoids that projection. `SendCodexEvaluation` keeps the selected model literal at `account_codex_eval_service.go:63-68,115-116`. |
| Model retirement | No v0.2.16 retirement helper is present in the recorded source inventory. | `backend/internal/service/openai_model_retirement.go:8-18,40-47,84-91`; `backend/internal/service/openai_account_scheduler.go:1815-1817`; `backend/internal/service/gateway_service.go:1464-1492` | v0.2.17 filters retired GPT-5.5 entries from ordinary OpenAI model listings and rejects retired subscription models in ordinary scheduling. The direct evaluation sender does not use the account scheduler or projected gateway catalog. Keep retirement behavior in the native gateway; do not route evaluation through it. |
| Settings | `research/native-candy-inspection-plan.md:48` records the original settings repository contract and table. | `backend/internal/server/routes/admin.go:567-590`; `backend/internal/service/setting_service.go:99-105`; `backend/internal/repository/setting_repo.go:85-115`; overlay `native/backend/internal/service/codex_inspection_settings.go:11-70` | v0.2.17 retains the native administrator settings route group and `GetMultiple`/`SetMultiple` repository contract. The inspection extension uses namespaced setting keys and the existing table. The patch adds only the inspection GET/PUT routes; it does not require a version-specific storage adaptation. |
| Priority | `research/migration-history-and-decisions.md` records the inspection design decision that lower numeric priority is preferred. | `backend/internal/repository/account_repo.go:2035`; `backend/internal/service/openai_account_scheduler.go:703`; overlay `native/backend/internal/handler/admin/codex_inspection_handler.go:14-20,45-54` | Native selection sorts ascending by priority. The inspection adapter reads the current priority and updates it through `AdminService.UpdateAccount`. Keep ascending priority semantics; do not carry over CPA's larger-is-preferred values. |
| Shutdown | `research/pristine-v0.2.16-upstream-contracts.md:19` and `research/source-audit.md:118` record the normal signal-driven server shutdown boundary. | `backend/cmd/server/main.go:171-194`; current patch `native/sub2api.patch:44-49` | The upstream process starts the server, waits for SIGINT/SIGTERM, then calls `Server.Shutdown`. The patch inserts `codexeval.StopInspection()` before that shutdown. The source context is otherwise unchanged, so the scheduler stop stays at the same lifecycle boundary. |
| Frontend patch | `research/response-model-validation-followup.md:9` records the small patch boundary and the current v0.2.16 URL helper anchor. | `frontend/src/utils/embedded-url.ts:16-42`; `frontend/src/components/layout/AppSidebar.vue:795-810`; `frontend/src/router/index.ts:516-545`; `backend/internal/web/embed_on.go:355-370`; patch `native/sub2api.patch:37-72` | The unchanged upstream helper parses the iframe URL with `new URL(baseUrl)` at line 25. The patch supplies `window.location.origin` as its base, so a relative evaluation path resolves with the browser's scheme, host, and port. It also adds the native sidebar and protected route. The separate page bypass is scoped to `/codex-eval`. |
| Response-model audit | `research/response-model-validation-followup.md:11,25` records the v0.2.16 observer and usage audit as unchanged by the patch. | `backend/internal/service/upstream_response_model.go:15-22,71-74,274-294`; `backend/internal/service/gateway_usage_billing.go:1275-1295`; `backend/internal/service/openai_gateway_usage.go:360-394` | The completed v0.2.16/v0.2.17 source comparison found `upstream_response_model.go`, `openai_gateway_response_handling.go`, `openai_gateway_usage.go`, `gateway_usage_billing.go`, and the existing `account_test_service.go` unchanged. `openai_gateway_forward.go` changed only legacy-ingress normalization, not the audit. The ordinary gateway still records mismatch metadata without rejecting forwarding. The v0.2.17 evaluation sender is a direct account-test transport and does not enter those usage-recording paths. Response-model comparison remains outside the plugin contract. |

The v0.2.17 source archive has `backend/cmd/server/VERSION` set to `0.2.16`. `backend/scripts/resolve-version.sh:7-23` prefers an exact Git tag when available, then reads this file. A source archive without Git metadata therefore resolves to `0.2.16`. The repository build entry point avoids this fallback: `scripts/build-sub2api.py:38-42` explicitly links `main.Version=0.2.17` and `main.BuildType=release`. Use that build path or another explicit version linker setting for archive-based builds.

## Files Found

- `backend/internal/service/openai_model_retirement.go`: v0.2.17 retirement predicate and catalog filtering.
- `backend/internal/service/openai_account_scheduler.go`: ordinary scheduling check for retired or unsupported request models.
- `backend/internal/service/account_test_service.go`: native account test and ordinary model-picker projection.
- `backend/internal/service/openai_models_list.go`: raw per-account OpenAI model discovery.
- `backend/internal/server/routes/admin.go`: native administrator settings route group.
- `backend/internal/service/setting_service.go`: native multi-key setting repository contract.
- `backend/internal/repository/setting_repo.go`: native multi-key settings reads and writes.
- `backend/cmd/server/main.go`: signal handling and HTTP server shutdown.
- `backend/cmd/server/VERSION`: stale source-archive fallback value `0.2.16`.
- `backend/scripts/resolve-version.sh`: tag-first, VERSION-file fallback logic.
- `backend/internal/service/upstream_response_model.go`: ordinary gateway response-model observer and mismatch calculation.
- `scripts/build-sub2api.py`: patched-source integration and explicit v0.2.17 linker override.
- `native/sub2api.patch`: evaluation routes, page bypass, settings routes, sidebar/router, shutdown, and relative iframe URL patch.
- `native/backend/internal/service/account_codex_eval_service.go`: raw evaluation model discovery and literal selected-model sender.
- `native/backend/internal/service/codex_inspection_settings.go`: inspection settings stored through the native repository.
- `native/backend/internal/handler/admin/codex_inspection_handler.go`: native account inventory, current priority, and native priority update adapter.

## Code Patterns

- Ordinary `/test` success recovery remains at `backend/internal/handler/admin/account_handler.go:1307-1317`.
- Account evaluation reads the selected model at `native/backend/internal/service/account_codex_eval_service.go:63-68` and writes it back unchanged at `:115-116`.
- Evaluation model discovery uses raw cached discovery at `native/backend/internal/service/account_codex_eval_service.go:21-42`.
- Native settings use multi-key reads/writes at `backend/internal/repository/setting_repo.go:85-115`.
- Inspection priority is read from the native account and updated through `AdminService` at `native/backend/internal/handler/admin/codex_inspection_handler.go:14-20,45-54`.
- Inspection shutdown precedes HTTP shutdown in `native/sub2api.patch:44-49`.
- The ordinary response-model observer documents non-interference at `backend/internal/service/upstream_response_model.go:15-22`; usage records the mismatch at `backend/internal/service/gateway_usage_billing.go:1285-1295` and `backend/internal/service/openai_gateway_usage.go:384-394`.

## External References

None. The review used the supplied source archive and repository files.

## Related Specs

- `.trellis/spec/backend/runtime-installation.md`: current administrator evaluation and inspection contracts.
- `.trellis/spec/backend/index.md`: current v0.2.17 source and delivery scope.

## Caveats / Not Found

- The v0.2.16 source tree is not present in the current temporary workspace. The version-to-version target-file comparison and hunk applicability were completed in the preceding source review; this report cites its persisted v0.2.16 anchors and the available v0.2.17 source.
- No build, test, runtime request, or external documentation lookup was performed for this report.
- No response-model comparison feature or extra plugin test is proposed. Preserve the native ordinary gateway audit behavior.
