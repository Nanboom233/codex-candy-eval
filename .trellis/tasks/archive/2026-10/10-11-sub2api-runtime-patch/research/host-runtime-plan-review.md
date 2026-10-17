# Review: Host Runtime Wiring Plan

- Date: 2026-10-12
- Scope: source-backed review of the host wiring plan and topology probe candidate.
- Pinned source: Sub2API `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` in `dist/sub2api-build/workspace/backend`.
- Verdict: **NEEDS-REFINEMENT**. The ticket scheme is viable for producer authorization, but the host plan needs executable session, ingress, and shutdown contracts before implementation.

## Candidate Assessment

The candidate can close the collector's producer-authentication gap without sending the raw evaluation token to the probe runner. A ticket is a short-lived, one-use capability for one predeclared observation. It is not the evaluation credential, does not authenticate a normal gateway request, and needs no database row or persistent token file.

Keep only ticket hashes in bounded in-memory session state. Bind each entry to one collector instance and its security fingerprint, one server-selected case, and an expiry. Consume it atomically before calling `Observe`. Consume tickets even when a network observation fails. Discard the entire session on completion, expiry, configuration drift, or process restart. `Seal` and `Initialize` retain the raw token inside the host service.

The existing primitives support the proposed network assertions: `Observe` reads the live request method, path, peer, resolved security client IP, and recognized forwarding headers; `Seal` requires the unique case set and current configuration fingerprint ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L525)). Neither method authenticates the producer or makes the request-supplied case label trustworthy ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L400)).

## Required Plan Refinements

### 1. Make Session Creation and Completion Explicit

The plan proposes one authenticated `POST /initialize` that accepts an inventory, runs the probes, and then calls `Initialize` ([host-runtime-wiring-plan.md](host-runtime-wiring-plan.md#minimal-administrator-endpoint-set)). An external probe runner must make multiple requests through different ingress paths between those actions. The current synchronous endpoint description does not define who runs them, how the runner receives each ticket, or how the admin completes the session.

Define a bounded lifecycle under the existing admin middleware chain: an authenticated admin starts a session with a deployment inventory; the host creates the collector and one ticket per required case; an authorized runner submits each ticket through the actual listener; an authenticated completion/status operation seals only complete observations and then calls `Initialize`. The API may use separate routes or explicit states on existing routes, but it must define the wire sequence and expiry behavior. An admin operation must not claim completion while required external probes remain outstanding.

Derive the session actor and Key owner from `GetAuthSubjectFromContext`; do not accept either as request data. Bind the session to that actor and the same owner passed to `EvalKeyProvisioner.Initialize`. The existing admin middleware provides a stable `UserID`; JWT requests use the authenticated user, while admin API-key requests resolve the configured first admin ([admin_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/admin_auth.go#L120), [admin_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/admin_auth.go#L156)). `RegisterAdminRoutes` already places routes under administrator auth, panel limits, audit logging, and compliance checks ([admin.go](../../../../dist/sub2api-build/workspace/backend/internal/server/routes/admin.go#L26)).

Specify how a ticket reaches a runner that can exercise each declared path. The current `EvalIngressPeer` carries only an ID and expected peer source, not a destination URL or another way for the host to initiate that path ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L132)). The persistent provisioning intent stores the sealed proof ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L196)). Keep the runner inside the existing admin/browser or server workflow. Do not add a sidecar or route the raw evaluation token through it. If an external runner is required, document its authorization and ticket-delivery boundary before implementation.

Ticket values should be returned only through an authenticated response with `Cache-Control: no-store`, held only as needed by the runner, and sent in one dedicated header. Never put them in a URL, request body, audit body, process argument, or diagnostic. They are active bearer capabilities even though they are not evaluation tokens. Do not persist ticket values or hashes in the proof, `provisioning.json`, or run state.

### 2. Guard the Marker Across Gateway Routes

The exact `/v1/responses` route currently runs body limiting, request-ID assignment, Ops error capture, endpoint normalization, and API-key auth before dispatch ([gateway.go](../../../../dist/sub2api-build/workspace/backend/internal/server/routes/gateway.go#L179)). Root aliases and `/backend-api/codex/responses` have separate chains, and Google ingress has another authentication chain ([gateway.go](../../../../dist/sub2api-build/workspace/backend/internal/server/routes/gateway.go#L337)). The candidate only places a one-use hook at `/v1/responses`.

Use one reserved header as the probe marker. Install a marker guard before normal gateway authentication on every ingress surface. With no marker, the middleware must pass through without changing the ordinary route. Any marker on an unsupported method, path, alias, provider, or protocol must fail closed and must not reach authentication, forwarding, or an upstream provider. Only a valid ticket on the exact `POST /v1/responses` route may call `Observe`, and that request must terminate before API-key auth, account selection, and billing. The marker must never select a case; the server-side ticket record selects it.

Add fixtures for missing, empty, malformed, repeated, expired, consumed, wrong-session, and wrong-route markers. Assert that invalid markers cause no observation, auth lookup, upstream call, usage row, or charge. A valid probe must also cause no auth lookup or billing. Compare ordinary no-marker behavior on `/v1/responses` with the existing path, and cover aliases plus Google routes to prove marker rejection and no header forwarding.

### 3. Preserve Ticket Secrecy in Observations and Logs

The default access logger records method, path, source, status, and request ID but not header values ([logger.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/logger.go#L13), [request_logger.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/request_logger.go#L16)). Ops error capture skips requests marked as ingress rejections ([ops_error_logger.go](../../../../dist/sub2api-build/workspace/backend/internal/handler/ops_error_logger.go#L1077)). These server defaults are compatible with a dedicated ticket header, provided the hook never logs its value and invalid-marker responses are marked for the existing rejection path.

There is a persistence collision to prevent. `evalRecognizedForwardedHeaders` includes both configured forwarding headers and built-in headers ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L292)); `Observe` copies every recognized header value into the observation ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L542)); the sealed proof is persisted with the provisioning intent ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L196)). Therefore session creation must reject a probe-marker header that collides case-insensitively with any recognized forwarding header, or the observation path must explicitly exclude the marker. Prefer failing closed on collision without changing global forwarding settings.

Required assertions: no raw evaluation token appears in ticket responses, probe headers, API responses, audit records, application logs, or persisted observations; no ticket value appears in logs or persisted proof; a configured header-name collision prevents session creation; and the invalid-marker path is absent from Ops error-body capture. Standard server logs do not establish what an undeclared reverse proxy logs, so deployment acceptance must also establish that the approved ingress paths do not log the ticket header.

### 4. Keep Public-Path Identity and Inventory Claims Honest

`Observe` records the server-side request path, peer, resolved client IP, and forwarding values. `checkObservation` requires the canonical `/v1/responses` path and checks a declared ingress case against its expected peer source ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L400)). The `IngressID` comes from the server-selected case, not from independently observed public routing metadata. Two public paths that reach the same peer source and rewrite to `/v1/responses` are indistinguishable to this proof.

Keep the candidate's explicit limitation: it cannot discover undeclared proxies, and a ticket label cannot prove which public path was traversed ([host-runtime-wiring-plan.md](host-runtime-wiring-plan.md#private-key-initialization-persistence-and-readiness)). Before claiming K4, identify the authoritative source of the complete deployment inventory and require independent evidence that every declared public path was exercised through its actual destination. An authenticated admin-submitted list alone is not completeness evidence. Keep readiness false when inventory provenance or per-path evidence is missing. Do not claim the collector discovers unknown hops or proves distinct paths that share one peer.

The current upstream proof can bind only its existing peer inventory and observations. Do not add a database schema or a second persistent evaluation token to bridge this gap. Whether external deployment evidence is sufficient for the accepted K4 boundary, and who is authorized to attest to it, remain required deployment/authorization inputs. No live proxy inventory was available in this review.

### 5. Complete Shutdown Ownership After a Deadline

The plan correctly requires manager drain before HTTP shutdown and dependency closure, and it says a timeout must leave the gateway available ([host-runtime-wiring-plan.md](host-runtime-wiring-plan.md#startup-and-shutdown-ordering)). The pinned host does not currently satisfy that consequence: `runMainServer` installs unconditional `defer app.Cleanup()` before starting the server ([main.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/main.go#L151)). Returning after a manager drain deadline still invokes cleanup. `provideCleanup` closes the OpenAI WebSocket pool and later Redis and Ent ([wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go#L335), [wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go#L397)). Avoiding `Server.Shutdown` alone is not enough.

Define one host lifecycle owner that keeps the process and local gateway alive while manager workers remain tracked, even after the first drain context expires. It may report the deadline and retry/wait for drain, but it must not return through the deferred cleanup path, shut down the HTTP server, close the WebSocket pool, or close Redis/Ent until workers settle. After settlement, persist/report terminal storage errors separately, then shut down HTTP and run cleanup. A hard process kill remains an external runtime limit, not a successful drain.

Add a host-level fixture where a worker remains blocked past the first manager deadline. Assert the local Responses route and all required dependencies remain available, and that cleanup has not run. Release the worker and assert the final order: manager drained, HTTP server shutdown, then gateway/dependency cleanup. Also assert the terminal-save-error case drains workers without misreporting persistence success.

## Decisions Already Supported by Source

- Use the existing `/api/v1/admin` chain and authenticated subject for initialize ownership; do not create another authentication route.
- Keep the process default not-ready and make `BootstrapRuntime` a recovery/guard operation. A failed bootstrap leaves normal host startup available.
- Keep the executor lazy. Its constructor requires a nonempty raw token, so Wire cannot require a real executor before initialization ([http.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/http.go#L27)).
- Resolve the local executor origin from the actual listener. `ServerConfig.Address()` is only configured `Host:Port`; port `0` does not provide the actual bound address ([config.go](../../../../dist/sub2api-build/workspace/backend/internal/config/config.go#L688), [config.go](../../../../dist/sub2api-build/workspace/backend/internal/config/config.go#L1539)). Keep readiness false until a literal loopback request reaches the exact local Responses route.
- Preserve normal Key authentication, permissions, subscriptions, balance, quota, concurrency, and billing for ordinary no-marker requests. Probe capabilities terminate before those checks and never forward.

## Verification

- Source review: complete for the named plan, PRD/design, current implementation checkpoint, backend quality/logging contracts, frozen topology primitives, gateway routes, admin identity, and shutdown source.
- Plan verdict: **NEEDS-REFINEMENT**, not a runtime/code gate.
- Product code, tests, live probes, deployment, and paid evaluations: not run or changed. The main session owns plan curation and the later implementation gate.
