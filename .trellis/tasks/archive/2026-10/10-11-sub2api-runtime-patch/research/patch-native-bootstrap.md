# Research: Patch-Native Bootstrap and In-Process Evaluation Dispatch

- Query: Can a patch-controlled, server-only evaluation dispatcher remove public topology probing and manual credential artifacts while preserving native authorization, billing, directed account selection, streaming/cancellation, and the accepted original-runtime K6 boundary?
- Scope: internal source/spec review with architecture analysis
- Date: 2026-10-12

## Findings

### Existing code already has a capability propagation seam

The current production patch carries the directed target through an unexported `context` key (`patches/sub2api/v0.2.15/production/001-codex-eval-auth.patch:44-61`). The patched API-key middleware currently obtains the target from externally supplied headers, validates the owned row, copies the Key/Group, installs the context value, and strips the headers before forwarding (`001-codex-eval-auth.patch:188-258`). The scheduler and cancellation patches already branch on `evalroute.FromContext` and disable ordinary affinity/failover behavior for directed requests (`002-codex-eval-scheduling.patch:17-20, 70-91, 122-136`; `003-codex-eval-cancellation.patch:63-77`).

The evaluation worker still crosses an HTTP boundary solely because the current design authenticates the normal gateway with the private token. `internal/eval/http.go:27-56, 90-110` builds a loopback `/v1/responses` request carrying `Authorization: Bearer <private token>` and target headers. `dist/sub2api-build/workspace/backend/internal/service/eval_host_runtime.go:36-70` obtains the actual listener and uses that HTTP executor. This is the replaceable part: the evaluation worker can call an internal gateway dispatcher that creates the context capability directly and never serializes a token, target header, probe ticket, or public origin.

### Minimum direct-dispatch architecture

The narrowest viable design keeps the normal database Key row for accounting but removes it from the request authentication path:

1. The authenticated administrator starts a run through the existing admin handler. The manager continues to resolve and validate the selected OpenAI group/account (`dist/sub2api-build/workspace/backend/internal/service/eval_target_adapter.go:54-80`; `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:792-817`).
2. `EvalHostRuntime` exposes a private dispatcher handle created at process startup. The handle owns an unexported capability value (for example, a package-private pointer or private context key) and a reference to the normal OpenAI gateway service. Only the runtime worker can obtain the handle. A public HTTP request can never manufacture a Go `context` value, so no external HTTP or WebSocket request can set this capability.
3. The dispatcher loads the exact owned row identity and owner, validates the requested target against the normal group/account rules, and creates a request-scoped Key/Group copy plus the authenticated owner/subscription context. It calls a refactored common OpenAI Responses admission/forwarding function with `context.WithValue` carrying the opaque capability and target. Do not expose a general exported `WithTarget` constructor for the public ingress path; retain a private dispatcher-only constructor or capability check.
4. The common pipeline must preserve the existing native stages: request validation, user concurrency, billing eligibility and in-flight balance reservation (`openai_gateway_handler.go:593-631`), directed scheduler selection and post-slot recheck (`openai_gateway_handler.go:643-814`), upstream forwarding, cancellation-aware streaming, partial-usage handling, and normal usage submission (`openai_gateway_handler.go:816-850`; `gateway_handler_responses.go:342-375`). The internal dispatcher returns the same sanitized `ModelResponse` shape consumed by `codexeval.ModelExecutor`.
5. The dispatcher must reject a target-account failure rather than fall back. Existing directed scheduler code already returns an unavailable result for account/group/capability mismatch (`002-codex-eval-scheduling.patch:173-247, 249-312`). Directed contexts must remain attached to the upstream context so cancellation and observed usage follow the current directed cancellation contract (`003-codex-eval-cancellation.patch:63-77`).

Calling `OpenAIGatewayHandler.Responses` directly with a synthetic `gin.Context` is possible in principle but is not the minimum safe change. That handler owns an HTTP writer, request body, request metadata, error response shape, stream heartbeats, and asynchronous usage callbacks (`openai_gateway_handler.go:390-500, 593-631, 816-850`). A synthetic context would need to reproduce those contracts and could accidentally misattribute client IP or detach the wrong cancellation context. A service-level common pipeline or an internal response sink is safer and easier to audit.

The generic `GatewayHandler.Responses` path explicitly rejects directed contexts (`dist/sub2api-build/workspace/backend/internal/handler/gateway_handler_responses.go:184-187`). The dispatcher must therefore target the OpenAI gateway path only, as the current evaluation target contract already restricts directed evaluations to OpenAI groups. It must fail closed for generic Responses/other providers rather than silently falling back.

### External ingress and WebSocket boundaries after the change

Every externally reachable API-key entry point must reject the exact owned stored lookup and raw private token before normal lookup, independent of ClientIP or forwarded-proxy headers. The patched main middleware already has the exact lookup rejection and raw-token mapping hooks (`001-codex-eval-auth.patch:188-213, 401-435`); the Google middleware rejects both values before `GetByKey` and rejects reserved target headers (`dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth_google.go:62-74`). The same exact checks must cover all main/Google aliases and any future middleware wrapper.

The public Responses WebSocket handler obtains its Key and user only after the gateway middleware (`dist/sub2api-build/workspace/backend/internal/handler/openai_gateway_handler.go:2445-2475`). Exact owned lookup/raw-token rejection there prevents a public WebSocket from obtaining the evaluation identity. The existing WS turn refresh deliberately calls `GetByKey(conn.Key)` (`dist/sub2api-build/workspace/backend/internal/handler/openai_gateway_handler.go:2168-2204`), so an internal evaluation WebSocket would need a separate capability-aware refresh path; the current evaluations use HTTP SSE and should not add an internal WS path merely to simplify installation.

The internal capability is network-unforgeable, not a cryptographic claim about arbitrary code already executing in the process. It is unavailable to external HTTP/WebSocket requests because those requests enter through `net/http`/Gin with fresh contexts and can only provide bytes, headers, and connection metadata. Keep the capability type/key private to the dispatcher package and never copy it into request headers, URLs, persistent evidence, logs, or an exported API response.

### Unix socket variant

A private Unix listener could retain the current HTTP executor shape, but it adds a second listener and routing path. The socket would need a short path beneath the protected evaluation directory, mode `0600`, ownership UID 1000, no host/Compose publication, and peer-credential checks before accepting a request. Even then, every process with the same container UID can connect, and a host operator with Docker/filesystem privilege can reach or replace the socket. It is therefore weaker and more complex than an in-process capability for the stated network threat model. It also requires lifecycle handling across restart, cleanup of stale sockets, and a dedicated non-public router. No existing production socket-dispatch seam was found.

### What this does and does not remove from K4

For the patched runtime, direct dispatch makes the public topology proof unnecessary for evaluation execution: no public request can authenticate the owned row, and no external path can inject the internal target capability. The activation gate can instead bind build/image/binary facts, protected state, exact owned-row identity, and the dispatcher capability. The external probe session, ingress inventory, `ingress_revision`, probe ticket header, and separate operator probe command can be removed from the default setup for this mode.

That conclusion does not preserve the current K6 boundary after unplanned recreation. The original executable can still see and accept the ordinary stored lookup value in the database. The accepted research explicitly records that an IP whitelist is only a privacy/exposure boundary, depends on forwarded-IP/proxy behavior, and cannot guarantee revocation after patch removal (`research/private-key-lifecycle.md:15-21, 54-60, 73-91, 135-151`). An in-process capability disappears with the patched process; it cannot constrain the original binary that starts after a writable-layer recreation. Therefore:

- K4 can be removed for patched evaluation dispatch while K6 remains unchanged only if the ordinary stored row keeps a native source restriction whose correctness is proven independently. That returns the source-IP/proxy topology problem for the original-runtime residual.
- Removing topology proof and keeping the row requires an explicit K6 revision that accepts original-runtime lookup visibility/use from broader sources. This is a weaker security contract and must be approved; it cannot be implied by the patch.
- Deleting/revoking the row before every possible recreation is not reliable for an unplanned event. A browser-authorized planned deactivation satisfies K5 but does not cover crash, image restart, or writable-layer loss before deactivation.

### Row-free internal identity is a separate architectural option

The strongest way to remove both public topology proof and the original-runtime lookup residual is to stop creating a persistent API-Key row. The dispatcher would use a protected owner/group evaluation principal held in server state and pass a request-scoped internal identity into native billing and scheduling. This is a substantial contract change, not a free patch:

- Current usage billing updates API-Key quota/rate state by `APIKey.ID` and invalidates the auth cache by `APIKey.Key` (`gateway_usage_billing.go:174-181, 296, 383-405`). Usage rows require an `api_key_id` foreign-key value (`ent/schema/usage_log.go:36, 200-202`).
- The existing handler passes the APIKey, user, group, and subscription into usage recording (`openai_gateway_handler.go:342-375`; WS equivalent `openai_gateway_handler.go:3203-3221`).
- A row-free principal would need an explicit billing/usage identity and quota model that preserves normal user/group subscriptions, balance, charges, user/account concurrency, and audit attribution. It would revise R4/K3's “normal Key” contract and probably the usage schema or its service boundary.

This option is technically coherent but broader than removing the current installer/probe machinery. It is the only route found that can make K6 disappear rather than merely accepting a weaker original-runtime residual.

### Browser-authorized setup and cleanup

Existing browser authentication can drive setup without a JWT file: the native page stores the access token in `localStorage`, and the API client attaches it as a same-origin header (`dist/sub2api-build/workspace/frontend/src/stores/auth.ts:104-125`; `dist/sub2api-build/workspace/frontend/src/api/client.ts:41-47`). The evaluation shell can call an authenticated activation endpoint that generates protected state, creates/reuses the owned identity, and merges its menu entry without returning a token. Internal JSON can remain generated server state rather than a user artifact.

For rollback/uninstall without a JWT file, the browser should first call authenticated deactivation/revocation, persist an owner-bound completion state, and then allow the host command to perform only the image/checksum-bound binary restore. The current runtime already drains and revokes the owned row before deleting private state (`eval_host_runtime.go:243-283`; `007-codex-eval-runtime.patch:874-885`). A host-only restore without that authenticated completion state would not satisfy K5.

## Proposed minimum architecture

1. Replace `evalLazyExecutor`'s loopback `HTTPExecutor` with an internal `EvalGatewayDispatcher` owned by `EvalHostRuntime`.
2. Refactor the OpenAI Responses admission/forward/usage pipeline behind a service entry point that accepts a private capability, an owner-scoped Key/Group/subscription snapshot, the selected target, and the evaluator context. Keep the public HTTP handler as an adapter into that common pipeline.
3. Keep external main/Google/WS authentication rejecting exact owned lookup/raw token values and reserved target headers before `GetByKey`.
4. Remove patched-mode dependency on inventory, ingress revision, external tickets, and runner commands. Readiness checks protected build/image/binary facts, exact owner/row identity (if the row model remains), and dispatcher availability.
5. Use the existing browser session for activation/menu setup and deactivation. Keep host installation binary-only and generate all protected state internally.
6. Decide separately whether to retain the ordinary row with the K6 residual/native ACL or undertake the row-free internal billing identity described above.

## Exact contracts requiring revision

| Contract | Required change if direct dispatch is adopted |
| --- | --- |
| R4 | Keep the normal row only for the row-backed option; the row-free option replaces it with an internal billing identity. |
| R7/K1-K3 | Preserve admin-only APIs, raw-token absence, normal authorization/billing and exact external replay rejection. |
| R14/K4 | Replace patched execution's source-IP/proxy proof with an internal-capability and external-replay gate; retain a separate gate only if K6 still relies on loopback ACL after recreation. |
| A16 | Replace public topology fixtures with exact lookup/raw-token replay rejection across HTTP/Google/WS plus direct-dispatch authorization, target, billing, concurrency, stream and cancellation tests. |
| K5/K7 | Keep exact owned-row drain/revocation/conflict behavior for the row-backed option; browser-authorized deactivation must precede host restore. |
| K6 | Either retain verified native source restriction, broaden the residual explicitly, or remove the persistent row. An in-process capability alone does not satisfy the current approved-local-only wording. |

## Feasibility verdict

The patch can automate the default setup and eliminate public probe/inventory/JWT artifacts for the patched runtime by using an internal capability and a shared gateway service path. The current code has a usable context/scheduler/cancellation seam, but direct dispatch requires a deliberate service refactor; a synthetic Gin context is not sufficiently small or safe.

Strict preservation of the accepted K6 boundary is the blocker. With a persistent normal Key row, source-IP/proxy controls remain necessary for the original binary after unplanned recreation. A zero-configuration design that also removes that residual needs a row-free internal billing identity, which changes R4/K3 and usage attribution contracts. No implementation, build, test, deployment, or live recreation evidence was available for this report.

## External references

- Sub2API v0.2.15, pinned source commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` from the task PRD.

## Related specs

- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/deployment-evidence.md`
- `.trellis/spec/backend/host-coordinator.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/private-key-lifecycle.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/patch-contracts.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/host-runtime-wiring-plan.md`
