# Research: Native Accounting Identity for Patch-Native Dispatch

- Query: Compare persistent expired, disabled, row-free, and active native Key identities for a private in-process evaluation dispatcher and simplified installation.
- Scope: mixed internal source/spec review with contract analysis
- Date: 2026-10-12

## Findings

### Native creation and expiry semantics

`APIKeyService.Create` validates `ExpiresInDays` as either nil or strictly greater than zero (`dist/sub2api-build/workspace/backend/internal/service/api_key_service.go:213-229,259-268`). It creates an active row with the supplied owner, optional group, IP rules, quota, and rate limits (`:568-582`), then computes `ExpiresAt = time.Now().AddDate(0, 0, *req.ExpiresInDays)` (`:584-588`) before one repository insert (`:590-592`). The repository writes `ExpiresAt` through the existing nullable column (`dist/sub2api-build/workspace/backend/internal/repository/api_key_repo.go:45-75`), and the schema indexes it without making it part of the key identity (`dist/sub2api-build/workspace/backend/ent/schema/api_key.go:72-76,136-147`).

The native public create path therefore cannot atomically create an already-expired row. A future row followed by `Update` is a two-step exposure window and can leave an active, never-expiring row after a crash. A safe expired-row design needs a narrow internal service/repository create method that reuses owner/group/IP/limit checks but supplies a past `ExpiresAt` in the single insert. It must retain the existing intent-first recovery sequence instead of writing the database directly.

The existing provisioner already persists a `row_pending` intent before lookup/create, adopts an exact deterministic lookup, and recovers `ErrAPIKeyExists` by rereading and verifying the row (`dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go:374-456`). Bootstrap can recover a row created before an intent save completed, but deliberately installs guards and returns `ErrEvalNotProvisioned` rather than activating an uncommitted identity (`:569-589`). Missing rows remain repair-required. The new expiry field must be part of the immutable expected identity and must never be repaired from the database after an external change.

### Main API-key lookup, cache, and expiry behavior

`GetByKey` checks the L1/L2 authentication caches first, then loads `GetByKeyForAuth` and materializes the presented lookup as `APIKey.Key` (`dist/sub2api-build/workspace/backend/internal/service/api_key_service.go:738-790`; `dist/sub2api-build/workspace/backend/internal/service/api_key_auth_cache_impl.go:269-331`). The auth snapshot explicitly stores `ExpiresAt` along with status, quota, IP rules, user, and group (`api_key_auth_cache_impl.go:334-350`). The repository lookup filters only soft-deleted rows, not expiry, so expired rows are still materialized (`dist/sub2api-build/workspace/backend/internal/repository/api_key_repo.go:40-43,111-127,130-147`).

The main middleware accepts a lookup from Bearer, `x-api-key`, or `x-goog-api-key`, and rejects query credentials before lookup (`dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth.go:49-96`). It calls `GetByKey` (`:99-122`), rejects disabled or unknown statuses during the base authorization stage (`:130-137`), and then, for normal billing endpoints, rejects `StatusAPIKeyExpired` and `apiKey.IsExpired()` (`:234-255`). `APIKey.IsExpired` is a runtime `time.Now().After(*ExpiresAt)` check (`dist/sub2api-build/workspace/backend/internal/service/api_key.go:67-82`).

The normal middleware intentionally skips the billing-stage expiry check for `/v1/usage`, `/v1/sub2api/billing`, and asynchronous image-task reads (`api_key_auth.go:188-193`). An expired row can therefore still authenticate those read-only or introspection paths after the user and base status checks. In standard mode it cannot issue an ordinary paid request. In simple mode the middleware returns before the billing-stage expiry checks (`api_key_auth.go:194-209`), so an active row with a past `ExpiresAt` would not be a complete external block if simple mode were enabled. The existing design rejects simple mode for directed evaluation and must keep that restriction.

Native updates can clear or change expiry. The user handler maps an empty `expires_at` to `ClearExpiration` and parses a new RFC3339 time (`dist/sub2api-build/workspace/backend/internal/handler/api_key_handler.go:229-287`). `APIKeyService.Update` writes only declared columns, clears expiry when requested, accepts an arbitrary supplied time, and can reactivate a row whose status is `expired` when clearing or extending it (`dist/sub2api-build/workspace/backend/internal/service/api_key_service.go:793-891,927-944`). These are legitimate native edits, so a private evaluation identity must record the expected expiry and report an exact conflict on any external change. It must never silently restore the past value or status.

### Cache materialization and invalidation limits

Auth-cache entries preserve the row's expiry, and `GetByKey` applies the cached snapshot without a database recheck (`dist/sub2api-build/workspace/backend/internal/service/api_key_auth_cache_impl.go:199-240,318-331`). Native `Create`, `Update`, delete, and quota-exhaustion paths invalidate by the stored lookup and publish cross-instance invalidations (`api_key_service.go:590-595,932-944,959-969,1172-1200`; `api_key_auth_cache_invalidate.go:5-12`; `api_key_auth_cache_impl.go:254-267`).

This is sufficient for native service mutations but not for privileged SQL edits or a failed Redis subscriber. A different instance may retain an old active/no-expiry snapshot until L1/L2 TTL expiry if invalidation is missed. Internal bootstrap and each evaluation start should read the exact row by ID from the repository, bypass auth cache, and fail on expiry/status/owner/binding drift. Atomic creation of the already-expired row avoids first populating caches with a temporary active value.

### Google authentication and query/header extraction

The Google middleware extracts `x-goog-api-key`, Bearer, `x-api-key`, and, for `/v1beta` or `/antigravity/v1beta`, query `key` (`dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth_google.go:234-271`). It rejects the private token and exact stored evaluation lookup before `GetByKey` (`:62-74`), then applies disabled-status checks (`:95-103`) and normal-mode status plus runtime expiry checks (`:161-179`). As with the main middleware, simple mode returns before billing expiry checks (`:147-159`).

Thus a past `ExpiresAt` blocks standard Google routes, does not block simple mode, and does not replace the patched exact-lookup replay rejection. The query-key route also means any proposed non-printable stored marker would need a separate query-path proof; an expired ordinary lookup has no such storage/parser exception.

### WebSocket initial auth and refresh

The OpenAI gateway group installs `apiKeyAuth` before `GET /v1/responses`, and the root/codex aliases reuse the same middleware chain (`dist/sub2api-build/workspace/backend/internal/server/routes/gateway.go:179-187,213-231,354-390`). The WebSocket handler reads the authenticated Key from context only after this middleware (`dist/sub2api-build/workspace/backend/internal/handler/openai_gateway_handler.go:2445-2475`). An already-expired row therefore fails the initial handshake in standard mode through the same expiry checks.

Subsequent turns are different. The refresh seam calls `GetByKey(conn.Key)` through L1/L2 auth cache (`openai_gateway_handler.go:2168-2204`). `refreshOpenAIWSTurnBillingAPIKey` only accepts a same-ID, same-group, same-platform/subscription snapshot and replaces the group used for pricing (`:2229-2248`). It does not check `Status`, `ExpiresAt`, or user/key limits. The WebSocket's later billing check is balance/subscription/rate-limit oriented and does not re-run API-key expiry in standard mode (`openai_gateway_handler.go:2700-2734`; `billing_cache_service.go:731-781`). Expiry before handshake blocks the connection. Expiry or status changes after handshake do not revoke the existing connection, and a refresh can still return the connection snapshot when the latest lookup is ineligible. This residual applies to expired and disabled-row strategies alike.

### Billing identity and foreign-key requirements

The gateway records the native row ID in each usage log (`dist/sub2api-build/workspace/backend/internal/service/gateway_usage_billing.go:1185-1218`), and the unified billing command carries `APIKeyID` as a deduplication/accounting identity (`:289-300`). Usage schema requires `api_key_id` and a required APIKey edge (`dist/sub2api-build/workspace/backend/ent/schema/usage_log.go:31-37,192-204`). Post-usage billing updates the row's quota and rate-window counters by `APIKey.ID` (`gateway_usage_billing.go:173-181`), and unified billing invalidates auth cache when quota exhaustion changes status (`:373-386`).

The row can be expired and still remain a valid foreign-key target. Native billing does not itself reject expiry; the rejection currently occurs in ingress middleware. A private dispatcher could therefore retain the row ID and all quota/rate/usage behavior, but only by explicitly bypassing the external credential expiry/status gate before entering the common gateway pipeline. A dispatcher that simply calls the existing HTTP middleware cannot use the expired row.

A row-free identity would require either nullable/synthetic `api_key_id` values or a new internal identity table and changes to usage-log foreign keys, billing deduplication, quota/rate updates, dashboard queries, and cache keys. A synthetic process-local ID would break restart and multi-instance attribution. This is a schema and accounting redesign, not an installation simplification.

### Exact owned-row and deactivation conflict behavior

The current provisioner verifies lookup, owner, active status, unbound group, and IP whitelist exactly (`dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go:669-710`). Bootstrap and `RevokeOwnedKey` use the same exact row checks (`:595-665`). A past-expiry strategy must add the expected `ExpiresAt` comparison while retaining `StatusActive`; it must not reinterpret an externally disabled row as an internal row. A disabled-row strategy would require changing every verify/revoke/bootstrap guard and would blur an administrator's ordinary disabled state with the evaluation identity.

Planned deactivation drains the manager first, verifies the actor owner, revokes the exact owned row, then removes protected session, intent, and token files (`dist/sub2api-build/workspace/backend/internal/service/eval_host_runtime.go:243-283`). A disabled or externally rebound row already causes revocation to fail as a conflict. Expiry changes must be handled identically: report conflict and retain recoverable transaction state rather than clearing expiry, re-enabling status, or deleting an unrelated row.

### Four identity strategies

| Strategy | Native external behavior | Billing/accounting | Recovery and residuals | Assessment |
|---|---|---|---|---|
| Active ordinary row | Full normal Key use wherever the original binary accepts the lookup. Patched replay rejection still applies only while patched. | Best fit. Existing row ID, auth cache, quota/rate counters, usage FK, dashboard and WS refresh all work. | Existing intent/recovery works. After recreation, the original runtime can use the lookup from any permitted source. | Fully compatible with current K1-K7 but leaves the largest K6 residual. |
| Disabled row | Main and Google base auth reject it even in simple mode. Existing WS connections and refresh are still not forcibly revoked. | Row ID and counters remain usable only if the private dispatcher bypasses disabled status. | Current provisioner and revoke guards reject it as `ErrEvalKeyStatusConflict`; every lifecycle path would need a special exemption. | Ambiguous and invasive. It changes the meaning of native status and risks violating K3. |
| Active row with past `ExpiresAt` | Standard-mode paid HTTP and initial WS/Google auth reject it at runtime expiry. Usage/billing reads can still authenticate, and simple mode is not blocked. | Best fit after a private dispatcher bypasses only the external expiry gate. Row ID and native quota/rate/usage paths remain intact. | Requires atomic expired creation and an expected-expiry conflict guard. Native clear/extend edits must block reapplication. Existing WS sessions can outlive expiry. | Preferred for stronger original-runtime paid-request denial, but it requires an explicit K3/R4/R5 contract exception for server-only billing dispatch. |
| Row-free internal identity | No original-runtime lookup exists. | Requires a new billing identity schema or nullable/synthetic FK design across usage, deduplication, quota, dashboards, caches and restarts. | New migration/recovery and multi-instance semantics. | Strongest credential isolation but disproportionate scope and risk for this task. |

### Original runtime, multi-instance, and lookup exposure

An expired row is safer than an active row only for native code that actually enforces `ExpiresAt`. The pinned v0.2.15 main/Google middleware does so for standard paid requests, but it still exposes the row to usage/billing reads and does not terminate established WebSockets. A deliberately cleared expiry, direct database edit, a simple-mode deployment, or a divergent older binary can restore or preserve ordinary lookup use. The original worker can also display the stored lookup through native management surfaces; K6 must not claim raw-token secrecy from expiry.

Atomic creation means no instance should see a temporary active/no-expiry row. After a native update, Redis-backed invalidation normally removes other instances' cache entries, but missed subscribers leave a bounded stale window. Exact row/expiry checks in patched bootstrap and dispatcher admission catch the drift regardless of cache state. An original unpatched worker cannot participate in those conflict guards, so unplanned recreation remains a residual rather than an automatically repaired state.

### Credential-free Docker host control

The binary replacement itself can remain credential-free because image/source/binary checks and staged replacement are host/container operations (`research/installation-simplification-options.md`; `.trellis/spec/backend/runtime-installation.md:23-27,75-81`). A trusted host could use `docker exec` with a script supplied on stdin to invoke a new server-local maintenance operation for deactivation, avoiding secrets in argv, environment, or logs.

No existing product CLI or local-control endpoint was found that safely performs native `RevokeOwnedKey`, cache invalidation, actor/conflict checks, and manager drain. Direct SQL or a shell mutation would bypass those contracts. A future Docker cleanup mode must call the same server-side `Deactivate` owner, pass no token in argv, preserve conflict reporting, and restore the original binary only after revocation succeeds. Until that local channel exists, browser-authenticated deactivation remains the fully specified path.

## Recommendation and required contract revisions

For the requested one-command install plus browser activation, use a **persistent native row with `StatusActive` and an already-past `ExpiresAt`, paired with a private in-process dispatcher**, only if the product accepts the explicit internal expiry exemption below. This keeps the native row ID for usage foreign keys, quota/rate counters, billing deduplication, dashboards, cache identity, and WS accounting while making ordinary standard-mode original-runtime paid requests fail closed after recreation. Do not use a disabled row, and do not introduce a row-free schema in this task. If K1-K7 must remain byte-for-byte unchanged, keep the active ordinary row with the existing accepted K6 residual instead; expired billing-only dispatch is a contract change.

Minimal wording changes required before implementation:

- **R4:** “Persist one native accounting Key row with an exact owner/row identity and an atomically past `ExpiresAt`. Keep the raw evaluation token only in protected server storage. Evaluation requests use a private in-process capability and never serialize the lookup or token.”
- **R5:** “The private capability may select the exact owned accounting row despite its external credential expiry, but must re-run normal user/group permission, subscription, balance, quota, model, account eligibility, target-group membership, user/account concurrency, billing, and cancellation checks. No public HTTP or WebSocket request receives this exemption.”
- **K3:** “Expiry is an external credential gate for the evaluation row. A server-only capability may bypass only that gate for the exact owned row during native billing dispatch. Status, owner, group binding, IP configuration, subscriptions, quotas, charges, and concurrency remain authoritative. Any external field drift is a conflict, never a repair.”
- **K4:** “For direct in-process dispatch, activation binds protected deployment/build facts and construction of the opaque capability. Public ingress/proxy topology proof is not an authorization input. Retain the current K4 probe contract for any legacy loopback/external HTTP executor mode.”
- **K6:** “After unplanned recreation, the original pinned standard runtime may display the stored lookup, but its past expiry must reject normal paid HTTP and new WebSocket authentication. Usage/billing introspection, existing WebSocket sessions, deliberate expiry clearing, simple mode, privileged database edits, and divergent binaries remain residuals. Automatic revocation is not promised.”
- **K7:** Keep the existing rule and add `ExpiresAt` to the exact owned-row conflict set. Never clear or rewrite an administrator's expiry/status edit during reapplication.

## Caveats / Not Found

- No product code was changed, and no build, test, deployment, live billing, Docker lifecycle, or browser check was run for this research.
- The direct dispatcher still requires a service-level extraction that preserves native streaming, cancellation, usage callbacks, and target scheduling. A synthetic Gin context remains unsafe without that refactor.
- The report relies on the pinned local Sub2API v0.2.15 workspace at commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`; no external documentation lookup was needed.
- The exact admin/browser activation and any trusted Docker `exec` cleanup RPC remain design work. Do not treat host privilege as permission to bypass native revocation or conflict checks.

## Related specs and research

- `.trellis/spec/backend/quality-guidelines.md` (private-Key, directed scheduling, cancellation contracts)
- `.trellis/spec/backend/runtime-installation.md` (protected lifecycle and deactivation)
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md` (R4-R5, K1-K7, installation revision)
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/patch-native-bootstrap.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/private-key-lifecycle.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/installation-simplification-options.md`
