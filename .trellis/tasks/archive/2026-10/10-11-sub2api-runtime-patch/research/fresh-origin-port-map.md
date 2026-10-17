# Research: Fresh Sub2API Native Port Map

- Query: Identify the exact v0.2.16 native constructors, routes, response transport, dependency-injection files, and embedded-frontend bypass points required for the approved administrator Codex evaluation port.
- Scope: internal
- Date: 2026-10-15

## Findings

This map was researched against the existing original Sub2API v0.2.16 tree that completed the pristine build. Use that source and this repository's internal/plugin directly; commit/archive identities from prior research are historical evidence, not mandatory source pins or repeated validation gates. The map below is limited to the approved administrator account endpoint, background evaluation service, and standalone `/codex-eval` page.

### Native administrator route boundary

- `backend/internal/server/routes/admin.go:26-32` creates the `/api/v1/admin` group and applies the native administrator middleware, panel limiter, audit middleware, and compliance guard. All new evaluation-control routes belong under this group; no second authentication mechanism is needed.
- `backend/internal/server/routes/admin.go:357-394` defines the account route group and the existing `POST /api/v1/admin/accounts/:id/test` at line 393. Add `POST /api/v1/admin/accounts/:id/codex-eval` beside that route, targeting the account-level complete-request handler. Keep `/test` unchanged.
- A separate task handler can be registered from the same `RegisterAdminRoutes` body as a new `registerCodexEvalRoutes(admin, h)` call. Its start/state/cancel/records/clear-history routes should remain under `/api/v1/admin/codex-eval/...`, inheriting the already-mounted administrator middleware.

### Account sender and constructor chain

- `backend/internal/handler/admin/account_handler.go:51-72` owns `accountTestService`; `NewAccountHandler` accepts it at `:87-120`. The direct account endpoint should be an additional method on this handler (or a deliberately small handler method sharing this injected service), not a loopback HTTP call.
- Existing `AccountHandler.Test` at `backend/internal/handler/admin/account_handler.go:1289-1318` parses the path ID and calls `AccountTestService.TestAccountConnection`. The new method should parse the same path ID, accept the approved complete streaming Responses body, and call the new native sender method. Do not route through `Test`, because that method emits simplified test events and performs the ordinary success recovery at `:1313-1317`.
- `backend/internal/service/account_test_service.go:139-160` defines the sender owner. `NewAccountTestService` at `:260-281` currently receives the account repository, token providers, HTTP upstream, config, and TLS fingerprint service. The native complete-request method belongs on this service so it can reuse account lookup, credential shadow resolution, model mapping, account header overrides, proxy, TLS, and plugin transport.
- The existing OpenAI test path at `backend/internal/service/account_test_service.go:799-991` is the source boundary for those pieces: account/model selection (`:804-837`), OAuth/API-key URL and credentials (`:840-877`), request headers and account overrides (`:902-946`), proxy and `doOpenAIAccountTestUpstream` (`:948-955`), OAuth snapshot and 401/429 handling (`:960-986`). Its fixed payload and `processOpenAIStream` summary are not the complete evaluation contract and must not be reused as the response envelope.
- `backend/internal/service/openai_plugin_transport.go:21-45` contains `doOpenAIAccountTestUpstream`, which first uses the native OpenAI OAuth plugin manager and otherwise calls `HTTPUpstream.DoWithTLS`/`Do`. The new complete sender should share this path and return the native `*http.Response` (or a context-aware equivalent) to the evaluation parser.
- `backend/internal/service/wire.go:255-283` is the source provider for `AccountTestService`; `backend/cmd/server/wire_gen.go:210-212` shows the generated construction order. Any new evaluation service dependency must be added to the normal provider set and generated constructor, rather than creating a second runtime graph.
- `backend/internal/handler/admin/grok_import_probe.go:195-234` is the existing pattern for production-only `ProvideAccountHandler`; `backend/internal/handler/wire.go:12-97` assembles `AdminHandlers`. If a distinct `CodexEvalHandler` owns task routes, add it to these two provider/assembly points and to `handler.AdminHandlers`; keep the account direct method on `AccountHandler` if that keeps the route surface local.

### Responses and Chat Completions conversion

- `backend/internal/pkg/apicompat/chatcompletions_responses_bridge.go:36-45` exposes `ResponsesToChatCompletionsRequest`; the same file contains the stream conversion helpers and tool/usage handling. The complete sender must use these existing native conversion functions when an explicitly supported account only accepts Chat Completions. No second conversion implementation is justified by this map.

### Standalone page and SPA bypass

- The native router is assembled in `backend/internal/server/router.go:22-95`; `SetupRouter` installs the embedded frontend middleware before route registration (`:73-93`). `registerRoutes` mounts the API/admin groups at `:116-135` and currently calls `handler.RegisterPageRoutes` at `:134`.
- `backend/internal/web/embed_on.go:86-117` serves embedded assets and falls back to `index.html` for any non-file path. `shouldBypassEmbeddedFrontend` at `:355-386` only bypasses known API/upstream paths. Add the exact `/codex-eval` shell path (and its trailing-slash form if served) to this bypass predicate, otherwise the frontend middleware will consume the request before the dedicated page handler runs. Keep API paths bypassed as they are.
- Existing page route registration is `backend/internal/handler/page_handler.go:260-284`; it is for data-backed markdown pages and is not the standalone evaluation shell. Register the evaluation shell from `backend/internal/server/router.go` through a dedicated native handler, alongside the route groups, after the bypass predicate is updated.
- `backend/internal/server/middleware/security_headers.go:108-151` applies `X-Frame-Options: DENY` globally and generates nonce CSP for non-API paths. This is already decided by F01/F04: the native custom menu embeds the page at the same origin, and original inline scripts/styles/renderers are preserved. The dedicated shell handler must override only its response with same-origin frame permission and the approved original inline capability; no HTML/style-to-class/nonce rewriting or global middleware change. SecurityHeaders sets headers before c.Next() (:128-150), so a route-local override is available without editing that middleware. The evaluation page's own resources should be embedded under the fresh evaluation package and served by the dedicated handler.
- The original page source is this repository's `internal/plugin/web/ui.html` and its embedded scripts/styles/vendor files listed in `internal/plugin/plugin.go:31-77`. Port only the approved standalone page assets into a fresh Sub2API package; replace CPA management calls with native administrator APIs while preserving layout, renderer assets, task controls, polling, history and theme behavior.

### Planned fresh files (names are implementation choices, not existing source)

The following new files are the smallest ownership split consistent with the native constructors above:

- `backend/internal/codexeval/` (or an equivalently named fresh package): migrated pure algorithms, request construction, SSE/result parsing, data and page assets; no CPA host callbacks.
- `backend/internal/service/codex_eval_service.go`: memory task maps, account busy coordination, four runners, history persistence and calls to the native account sender.
- `backend/internal/handler/admin/codex_eval_handler.go`: administrator start/state/cancel/records/clear-history handlers and standalone shell serving. `AccountHandler` gains only the direct single-account endpoint method.
- `backend/internal/service/wire.go`, `backend/internal/handler/wire.go`, `backend/internal/handler/handler.go`, and generated `backend/cmd/server/wire_gen.go`: ordinary provider/assembly updates only.
- `backend/internal/server/routes/admin.go`, `backend/internal/server/router.go`, and `backend/internal/web/embed_on.go`: account/task route registration, shell route registration, and frontend bypass.

## Files Found

- `backend/internal/server/routes/admin.go` — native admin group and account `/test` route.
- `backend/internal/handler/admin/account_handler.go` — existing account handler, constructor and `/test` behavior.
- `backend/internal/service/account_test_service.go` — account lookup, credential/model preparation, native test sender and response handling.
- `backend/internal/service/openai_plugin_transport.go` — plugin/HTTP upstream handoff returning `*http.Response`.
- `backend/internal/service/wire.go` — service provider for `AccountTestService`.
- `backend/internal/handler/admin/grok_import_probe.go` — production handler provider pattern.
- `backend/internal/handler/wire.go` and `backend/internal/handler/handler.go` — admin handler assembly.
- `backend/cmd/server/wire_gen.go` — generated native construction order.
- `backend/internal/pkg/apicompat/chatcompletions_responses_bridge.go` — existing Responses/Chat conversion.
- `backend/internal/server/router.go` — middleware and route assembly order.
- `backend/internal/web/embed_on.go` — embedded frontend fallback and bypass predicate.
- `backend/internal/server/middleware/security_headers.go` — global frame/CSP behavior.
- `backend/internal/handler/page_handler.go` — existing markdown page routes, intentionally separate from the evaluation shell.
- `internal/plugin/plugin.go` plus `internal/plugin/web/` — current repository's original standalone page embed list.

## Related Specs

- `.trellis/spec/backend/patch-native-installation.md` — current native administrator and fixed delivery boundary.
- `.trellis/spec/backend/runtime-installation.md` — historical runtime/UI contracts; the fresh plan supersedes its retired activation/private-transport sections.
- `.trellis/spec/backend/evaluation-core.md` — pure evaluation and streaming contracts, subject to the current fresh-source plan.
- `.trellis/spec/backend/quality-guidelines.md` — scope, ordinary Wire construction, and no added race gate.
- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md` — approved native route, task, page and source boundaries.

## Caveats / Not Found

- This is a source map only. No product code, generated Wire file, build, test, deployment, or provider request was changed or run.
- Exact new Go type and filename names remain implementation choices; the existing constructor and route owners above are established by the available original source.
- The original source has no pre-existing `/codex-eval` route or page; page-scoped same-origin embedding and original inline capability follow approved F01/F04; implement the route-local response policy without weakening global headers or rewriting output.
- Reuse the existing source tree recorded in the completed baseline report; no new source manifest, fixed-source export or Git clone is required. Conflicting older spec/research test matrices do not add requirements to the current task.
