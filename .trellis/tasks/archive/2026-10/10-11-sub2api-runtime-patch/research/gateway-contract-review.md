# Research: Gateway Contract Review (corrected)

- **Query**: Verify and correct the previous gateway-contract review against pinned sub2api `v0.2.15` commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` (`dist/sub2api-v0.2.15/backend/`).
- **Scope**: internal (pinned upstream source); no live builds/tests.
- **Date**: 2026-10-11 (previous report incorrectly dated 2025).

Status legend: **[V]** verified against source, **[I]** inferred from code shape, **[U]** unresolved product/lifecycle question.

---

## 1. Authentication ordering — corrected §2.1 [V]

`internal/server/middleware/api_key_auth.go` (449 lines) actual order:

1. Abuse limiter + header-size + query-key deprecation checks (`:37-96`).
2. `apiKeyService.GetByKey` at `:100`; `SetOpsFallbackAPIKey` at `:119`.
3. Key disabled check `:124-130`; IP allow/deny `:134-146`.
4. `apiKey.User == nil` → 401 `:149-152`; **`!apiKey.User.IsActive()` → 401 `:155-159`** (active-user check).
5. Group checks: `abortIfAPIKeyGroupUnavailable` `:160` (GROUP_DELETED/GROUP_DISABLED → 403) then `abortIfAPIKeyGroupNotAllowed` `:163` (GROUP_NOT_ALLOWED → 403).
6. `ctxkey.UserID` set `:166`; `skipBilling` endpoints computed `:168-172`.
7. SimpleMode early return `:176-189` — sets `ContextKeyAPIKey`, `AuthSubject`, `ContextKeyUserRole`, `setGroupContext` there.
8. Normal mode: subscription load `:193-212` — **missing subscription returns 403 `SUBSCRIPTION_NOT_FOUND` at `:205`** (previous report said 402).
9. Billing enforcement `:216-268`; context set + `c.Next()` `:272-286`. `AuthSubject` set at `:276-279`, not at any `setAuthSubject` symbol (no such function exists).

Correction: there is **no `setAuthSubject` at :160**; `:160` is the group-availability check. The directed-target hook belongs after the active-user check (`:159`) and before/at the group checks (`:160-165`), matching patch-contracts.md §2 ("after active-user authentication, before group availability and permission checks") — anchor to `api_key_auth.go:160-165`, not to a nonexistent `setAuthSubject`.

`APIKeyService.Create` (`api_key_service.go:492-598`) **permits `GroupID == nil`**: group binding permission (`canUserBindGroup` `:481-489`) runs only when `req.GroupID != nil` (`:517-527`). `ExistsByKey`/`ErrAPIKeyExists` (`:544-551`) is **duplicate detection for `CustomKey`, not idempotent recovery** — re-running `Create` with the same key returns `ErrAPIKeyExists` and increments the error counter (`:550`). Recovery of interrupted provisioning must be keyed on the persisted intent + exact lookup identity, tested for owner/status/journal cases.

## 2. Private-token threat model — corrected facts and confirmed K6

- **[V]** The *unpatched* binary cannot authenticate the RAW private token (absent from DB), but the deterministic stored lookup value **IS** a valid credential: `GetByKey` (`api_key_service.go:739-791`) does a cached deterministic DB lookup of whatever string is presented. The patch must map raw token → lookup value and **reject the lookup value itself as a gateway credential** while patched (patch-contracts.md §3 stands).
- **[V]** An `s2eval-<digest>` row created via normal `Create` is an ordinary active Key. Do **not** reserve the whole `s2eval` prefix or reject pre-existing ordinary keys; matching on exact configured row/lookup identity (immutable Key ID + owner) suffices.
- **[V]** Unbound-key inference is governed by `RequireGroupAssignment` (`server/middleware/middleware.go:143-160`): `GroupID == nil` → 403 only when `IsUngroupedKeySchedulingAllowed` is false; it **passes when true**. Owner binding of the row is a separate choice. Distinguish **key authentication** (works regardless of group) from **inference authorization** (requires bound group when ungrouped scheduling is disabled).
- **[V]** With the required data volume preserved, the raw-token file and lookup row survive recreation. The original executable lacks the token-to-lookup mapping.
- **Confirmed K6**: the user accepts original administrator visibility and possible normal-Key use from approved local sources. This resolves the risk decision, not the source-IP/proxy checks or a stronger privacy guarantee.
- Test owner/status/journal recovery paths; `ExistsByKey` conflict must abort-and-recover by identity, not retry blindly.

## 3. Responses failover-loop anchors — corrected [V]

`handler/openai_gateway_handler.go` (the OpenAI POST /v1/responses handler; note `handler/gateway_handler_responses.go` is the *Anthropic-platform* `/v1/responses` handler, 420 lines — previous anchors belong to the OpenAI file):

- `:676` `SelectAccountWithSchedulerForCapability` inside the failover `for` loop.
- `:740` `account := selection.Account` — local capture.
- `:741-762` unsupported HTTP continuation: `previous_response_id` with non-API-key account → exclude + `continue` (`:762`) re-entering selection.
- `:768` `acquireResponsesAccountSlot`; `:769-776` profit veto → `recordOpenAIProfitVeto` + `continue` (reselection), exhaustion → `handleOpenAIProfitVetoExhausted` at `:2241-2249` which writes 503 and returns — **terminal**, not a reselection.
- `:798` `h.gatewayService.Forward(c.Request.Context(), c, account, attemptBody)` forwards the **local** `account`.
- `:916-933` same-account retry → `continue` **re-enters the whole selection loop** (relies on sticky/session re-picking, not a pinned same-account call).
- `:936-965` account switch: record, `failedAccountIDs[account.ID]`, `continue` at `:965`.
- `:967-985` non-failover error: `:969-971` `ensureForwardErrorResponse` writes fallback error — **not** a credential reselect.
- Profit-veto exhausted in the generic `GatewayHandler.Responses` is at `gateway_handler_responses.go:260-264` (503 terminal) — do not conflate the two handlers.

Directed-request guard requirement (corrected): terminate on unsupported continuation (`:741-762`), profit veto (`:769-776`), failover switch (`:965`), and same-account retry (`:933`) — each `continue` re-enters `SelectAccountWithSchedulerForCapability`, so the directed pin must be enforced inside selection, not only at loop entry.

## 4. F19 — public selector enters legacy directly [V]

`OpenAIGatewayService.SelectAccountWithLoadAwareness` (`openai_gateway_scheduling.go:1129-1136`) calls `selectAccountWithLoadAwareness` (`:1138`) **directly**, bypassing `selectAccountWithScheduler`/`selectAccountWithSchedulerOnce` entirely (it installs quota-pause, privacy, and profit-gate ctx itself at `:1130-1134`). Therefore "patch the common `Once` at `openai_account_scheduler.go:2298` covers all selections" is **false**: it covers only `SelectAccountWithSchedulerForCapability` (`:2116` → `selectAccountWithScheduler` `:2163` → `Once` `:2298`), i.e. the normal POST Responses selector. Callers of the public `SelectAccountWithLoadAwareness` (e.g. `gateway_handler_responses.go:190`, chat-completions handler) reach the legacy selector without touching `Once`. The directed interception must live in `selectAccountWithSchedulerOnce` **and** be enforced or explicitly rejected on the public legacy entry, or be placed in the shared `selectAccountWithLoadAwareness` internals (e.g. `getSchedulableAccount`/`recheckSelectedOpenAIAccountFromDB` call sites `:999,1219,1386-1479`).

## 5. F5 — slot-admission swaps `selection.Account`; Responses forwards stale local [V]

`acquireOpenAIAccountSlot` (`openai_gateway_handler.go:2269-2402`) calls `ProfitControlVetoLatest` at `:2295`, `:2334`, `:2390` and on success assigns `account = latest` + `selection.Account = latest` (`:2303-2304`, `:2342-2343`, `:2398-2399`). But the Responses loop captured `account := selection.Account` at `:740` **before** `:768`, and `Forward` at `:798` uses that local — so the hydrated/`latest` pointer propagated into `selection.Account` is **not** used for forwarding or usage accounting (`:826`, `:835` use local `account`). Additionally `ProfitControlVetoLatest` (`openai_profit_control.go:333-338` → `profitControlVetoLatest`) re-reads the cached account and re-evaluates **only the rate gate**; it does not necessarily re-check group membership, runtime block, or model/capability eligibility after the wait.

Directed-design consequence: after any wait (WaitPlan path `:2371-2390`), a directed request needs a **directed fresh check after the wait** (re-verify target account identity, `openAIAccountMatchesSchedulingGroup`, runtime/model eligibility via `recheckSelectedOpenAIAccountFromDB`, `openai_gateway_scheduling.go:1597`), and the local `latest` pointer must be propagated to the forward/usage call sites. Assert `latest.ID == target.AccountID` and group equality — do not assert pointer inequality as the success signal.

## 6. Quarantine fail-open wrapper [V]

`selectAccountWithScheduler` (`openai_account_scheduler.go:2163-2195`) retries `selectAccountWithSchedulerOnce` with `withOpenAIProxyStreamQuarantineBypass(ctx)` at `:2194` when the only failure is quarantined proxies (`:2179-2193`). For directed requests a **directed early return before the wrapper's bypass retry** is sufficient for this wrapper, provided the directed helper itself also rejects a pre-existing `openAIProxyStreamQuarantineBypassed(ctx)` context. Do not claim "only stripping the bypass works" — both points are required together.

## 7. F15 — POST /v1/responses CAN reach the upstream WebSocket transport [V] — most important correction

Route trace: handler `:798` → `OpenAIGatewayService.Forward(ctx, c, account, body)` (`openai_gateway_forward.go:21`, `ctx` = `c.Request.Context()`) → WS transport decision at `:786` (`OpenAIUpstreamTransportResponsesWebsocketV2`) → `forwardOpenAIWSV2(ctx, ...)` at `:887` inside a reconnect loop `:884-965` with per-attempt backoff that **does** observe `ctx.Done()` (`:955`).

Inside `openai_ws_forwarder_v2.go`:
- Conn-pool acquire uses `acquireCtx = WithTimeout(ctx, acquireTimeout)` (`:202-205`) — cancellable by client.
- Write uses `lease.WriteJSONWithContextTimeout(ctx, ...)` (`:345`) — cancellable.
- Read loop uses `upstreamReadCtx`, initially `ctx` (`:390`), but `markClientDisconnected` swaps it to `context.WithoutCancel(ctx)` at **`:402`** so draining continues after client cancel; reads happen at `:542` via `ReadMessageWithContextTimeout(upstreamReadCtx, ...)`.
- Lease ownership: single `defer` at `:263-267` — `MarkBroken()` on error else `lease.Release()` back to the **shared pool**. A drained-but-canceled turn can release a conn whose next turn binds another response's state (`:824-828`).

Corrections:
- `openai_ws_http_bridge.go:493` `detachUpstreamContext` sits inside `proxyOpenAIWSHTTPBridgeTurn` (`:409`) which builds an **HTTP upstream request** — it is the **WS-ingress → HTTP-upstream** bridge, not HTTP→WS. Do not list it as the POST→WS detach point.
- POST-only ingress does **not** exclude `forwarder_v2.go:402`. The cancellation patch (patch-contracts §5) must cover: (a) `markClientDisconnected`/`upstreamReadCtx` detach for trusted eval context — cancel should stop the read rather than drain; (b) close ownership — on directed cancel, `MarkBroken` instead of returning the leased conn to the pool, since partial turn state may be bound (`:824-828`); (c) partial-usage semantics — draining exists so observed usage reaches accounting; for directed cancel decide whether `result.ClientDisconnect` partial usage (`handler :862-868`) is still recorded via `detachedBillingContext` (keep billing detach unchanged).
- Preserve ordinary draining behavior for non-directed traffic — no flag flips on shared helpers; gate on trusted eval context only.

## 8. Runtime lifecycle — corrected [V]

The runtime is the **original official image unchanged**; patching is **container-local binary replacement** (stage beside `/app/sub2api`, atomic rename, restart only that container). Do not write "patched image ships". `web.HasEmbeddedFrontend()` (`internal/web/embed_on.go:408`, `embed_off.go:46`) is **not** recreation proof — the original binary also embeds the frontend under `-tags embed`. Post-restart verification requires the **actual binary hash** plus a **patch-status endpoint**; a health or SPA response proves nothing.

## 9. Regression matrix — corrected assertions (anchored to existing fixtures)

| # | Assertion | Owner / fixture anchor |
|---|---|---|
| R1 | Directed target resolved only after active-user check, before group checks; missing subscription still 403 | `api_key_auth_test.go`; `api_key_auth.go:155-165,205` |
| R2 | Original host accepts lookup-value authentication. Patched host rejects external lookup-value replay before mapping. Raw token remains absent from DB/DTO | `api_key_auth_test.go`, `api_key_service.go:739`, `dto/mappers.go:80-111` |
| R3 | `Create` with `GroupID=nil` succeeds; ungrouped-key inference blocked iff `IsUngroupedKeySchedulingAllowed=false` | `api_key_service.go:517-527`, `middleware.go:143-160` |
| R4 | `ExistsByKey` conflict on recreate → recover by exact lookup identity + journal; never re-enable externally disabled Key | `api_key_service.go:544-551` |
| R5 | Directed request: every `continue` at handler `:762,:776,:933,:965` terminates instead of reselecting; profit-veto exhausted `:2241-2249` is terminal 503 | `openai_profit_veto_budget_test.go`, `openai_gateway_credential_failover*_test.go` |
| R6 | `SelectAccountWithLoadAwareness` (public, `:1129`) directed context → directed path or hard reject; does not silently enter legacy | `openai_gateway_service_test.go` (`TestOpenAISelectAccountWithLoadAwareness_*`) |
| R7 | After WaitPlan acquisition, directed fresh check re-verifies ID/group/model; forward uses propagated `latest` local | `openai_profit_slot_recheck_test.go`; handler `:2295,:2390,:798` |
| R8 | Directed ctx present → quarantine bypass retry at `:2194` skipped and pre-bypass ctx rejected | scheduler tests; `openai_account_scheduler.go:2178-2194` |
| R9 | Directed cancel on WS transport → `upstreamReadCtx` cancels (no `WithoutCancel` at `forwarder_v2.go:402`), conn `MarkBroken` not pooled, partial usage still billed | `openai_ws_forwarder_client_cancel_test.go`, `gateway_handler_cancellation_test.go` |
| R10 | `openai_ws_http_bridge.go:493` unchanged (WS→HTTP bridge); ordinary HTTP detach at `openai_gateway_forward.go:1038` unchanged for non-directed | — |
| R11 | Status proves patch by binary hash + patch endpoint, not `HasEmbeddedFrontend` | installer contract (patch-contracts §8) |

## Caveats / Not Found

- This corrected rewrite supersedes the first report's inaccurate anchors and overbroad conclusions. The main session reviewed the earlier report before replacement.
- Main-session inspection of `gateway_profit_control.go:83-101` confirms the helper refreshes a scheduler snapshot and checks only the profit rate gate. It does not perform membership/runtime/model checks. This is a verified directed post-wait gap, not an assumed complete eligibility recheck.
- Line numbers verified at commit `f2669c8`; they will drift on any upstream bump.
