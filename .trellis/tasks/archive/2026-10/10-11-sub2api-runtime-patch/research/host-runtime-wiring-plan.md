# Research: Sub2API Host Runtime Wiring

- Query: Define source-backed host seams for administrator metadata/API wiring,
  private-Key bootstrap, topology probe delivery and provenance, marker
  dispatch, probe-session abort/concurrency, lazy executor activation, Stage 4B
  ownership, provisioner-error recovery, protected inventory loading, Wire
  providers, and shutdown ordering.
- Scope: internal
- Date: 2026-10-12

## Findings

### Current boundary and implementation status

Stage 4A is frozen in `internal/eval`. Its `Target` is the positive numeric
group/account identity, `AccountMetadata` is display-only, and `ResultInfo`
already records the numeric target, run ID, account name, requested/effective
model and effort, and terminal status ([eval.go](../../../../internal/eval/eval.go#L13-L79)).
The Stage 4B contract in [eval-run-state-plan.md](eval-run-state-plan.md)
therefore remains the right boundary. The pinned copied workspace now contains
`TargetResolver`, `NewRunManager`, `StateStore`, and `RunManager.Shutdown`
([runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L19-L23),
[runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L98-L138),
[state.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/state.go#L31-L44),
[runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L656-L688)).
Use `Flush` explicitly and do not add group names, owner IDs, or artifact status
to `ResultInfo`. Stage 4B implementation and its independent specification and
quality gates have passed ([eval-run-state-spec-review.md](eval-run-state-spec-review.md),
[eval-run-state-quality-review.md](eval-run-state-quality-review.md)). Host
adapter and lifecycle wiring remain separate integration work.

The pinned v0.2.15 workspace has the reviewed Stage 2A, 2B, and 2C code from
the three frozen patches, but it has no production caller for the Stage 2A
topology collector or key provisioner. There is no admin evaluation handler,
probe-session route, bootstrap call, or Wire binding for these primitives in
the current source. The generated workspace is evidence for later patching
only; the repository does not yet contain the host runtime integration.

The Stage 4B worker must be independent of the HTTP request after its durable
start transaction. Its parent context is the manager lifetime. `Flush` is the
only recovery path for transient terminal-write errors. A malformed state file
or a failed recovery save preserves the original bytes and blocks new starts;
`Flush` must refuse to overwrite those protected states. Shutdown must expose
the distinction between workers still running and workers drained with a
terminal storage error.

### Administrator authentication and route seam

The existing router creates `/api/v1` and calls `RegisterAdminRoutes` beside
the ordinary user and gateway routes ([router.go](../../../../dist/sub2api-build/workspace/backend/internal/server/router.go#L119-L134)).
`RegisterAdminRoutes` creates `/admin`, applies `AdminAuthMiddleware`, panel
rate limiting, audit logging, and the compliance guard before registering
feature groups ([admin.go](../../../../dist/sub2api-build/workspace/backend/internal/server/routes/admin.go#L13-L35)).
The evaluation endpoints belong in this group. They must not create a second
auth path or accept a gateway API key.

`NewAdminAuthMiddleware` accepts either `x-api-key` containing the configured
admin API key or `Authorization: Bearer <JWT>` ([admin_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/admin_auth.go#L14-L79)).
The JWT path loads the user, verifies active status, token version, session
binding, and admin role before setting `ContextKeyUser`, `ContextKeyUserRole`,
email, and session ID ([admin_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/admin_auth.go#L156-L219)).
The admin API-key path resolves the first admin and sets the same subject and
role ([admin_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/admin_auth.go#L120-L153)).
An evaluation handler should read `middleware.GetAuthSubjectFromContext(c)` and
use its positive `UserID` as the initialization owner; it should not infer an
owner from configuration, an environment variable, a browser value, or a
previously asserted status payload.

Use the standard host response envelope. `response.Response` is
`{code,message,reason,metadata,data}`, `response.Success` returns HTTP 200 and
`response.Accepted` returns HTTP 202 ([response.go](../../../../dist/sub2api-build/workspace/backend/internal/pkg/response/response.go#L14-L57)).
`response.ErrorFrom` maps application errors through the host error mapper and
redacts logged internal details ([response.go](../../../../dist/sub2api-build/workspace/backend/internal/pkg/response/response.go#L69-L95)).
The embedded API adapter should decode `data` and preserve `message` for errors.

### Metadata adapter and exact target admission

The existing `AdminService` exposes the narrow read operations needed by an
adapter: `GetAllGroupsByPlatform`, `GetGroup`, `GetAccount`, and
`GetAccountsByIDs` ([admin_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_service.go#L34-L42),
[admin_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_service.go#L73-L83)).
`GetAllGroupsByPlatform` delegates to the repository's active-platform query
([admin_group.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_group.go#L52-L58)).
`GetAccountsByIDs` is a single repository lookup and returns only service
objects ([admin_account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_account.go#L55-L69)).
The service `Account` model contains display name, platform, type, status,
schedulability, group IDs, and runtime fields as well as a `Credentials` map
([account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/account.go#L24-L68)).
The host DTO for `/targets` must construct a new safe metadata shape and never
marshal the service object or its credentials. `AccountMetadata` in the core
should remain name-only.

Use `APIKeyService.ResolveDirectedTargetGroup` as the exact target admission
primitive. It rejects nonpositive IDs and missing groups, inactive/deleted
groups, non-OpenAI platforms, and accounts absent from the group's membership
set ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L701-L725)).
The resolver returns the requested target unchanged. A `TargetResolver` adapter
can call that method, then call `AdminService.GetAccount` for the display name
and verify the returned account ID and OpenAI platform. A missing, changed, or
mismatched account becomes the requested target's typed unavailable record; it
must never trigger another group or account lookup as a fallback.

For `/targets`, enumerate active OpenAI groups through
`GetAllGroupsByPlatform(ctx, service.PlatformOpenAI)`, obtain membership IDs
through the existing group repository/service path, and batch hydrate names
with `GetAccountsByIDs`. Return only positive IDs, group name, account name,
account type, and safe status/model capability fields selected by the host DTO.
Do not use a paginated admin HTTP endpoint as an implicit authorization check,
and do not expose credentials, tokens, proxy settings, or raw `Extra` values.
The final target request still re-resolves membership through
`ResolveDirectedTargetGroup`, because the target list is a snapshot and groups
or accounts can change between listing and start.

The patched ingress preserves normal policy after it maps the private token to
the stored lookup. `resolveEvalIngress` rejects external stored-lookup replay,
rejects directed headers from ordinary Keys, requires a valid topology proof,
rejects simple mode and non-`POST /v1/responses`, then parses positive target
headers ([api_key_auth_eval.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth_eval.go#L19-L63)).
After `GetByKey`, the normal middleware still checks IP ACL, user existence and
status, group availability/permission, subscriptions, balance, quota, and rate
limits before installing the directed request context
([api_key_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth.go#L99-L187),
[api_key_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth.go#L211-L306)).
The adapter must not duplicate or bypass those checks.

### Private-Key initialization, persistence, and readiness

The reviewed Stage 2A key lifecycle is already source-backed. The private raw
token is generated once with a random 32-byte suffix and stored outside the
database. `EvalDeriveLookup` stores a namespaced SHA-256 lookup value, while
`EvalKeyStore` owns `private-token` and `provisioning.json`
([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L18-L32),
[eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L67-L109)).
The storage directory is mode `0700`; token and intent writes enforce the
existing restrictive file contract ([private-key-lifecycle.md](private-key-lifecycle.md#relevant-source-contract)).

`EvalKeyProvisioner.Initialize` serializes provisioning, validates the supplied
topology, loads or creates the private token, validates the proof against the
current config, creates or adopts the normal API-Key row, verifies owner,
lookup, active status, unbound group, and exact IP whitelist, then installs the
runtime identity ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L365-L457),
[eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L480-L530)).
The API-Key row is created through `APIKeyService.Create`, so existing user,
IP, group-binding, duplicate, creation-limit, status, quota, and cache rules
remain in force ([api_key_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/api_key_service.go#L491-L597)).

Startup must call `BootstrapRuntime` only as a recovery/guard operation. It
must verify the persisted row identity, owner, active status, unbound group,
whitelist, raw-token-to-lookup derivation, and live topology proof. Any missing
row, owner/status/binding/IP/lookup conflict, expired proof, or config drift
installs identity-only guards and leaves evaluation unavailable
([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L569-L625)).
`Load` is stricter and requires a current valid proof before reporting the
runtime usable ([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L533-L566)).
Planned uninstall/rollback must call `RevokeOwnedKey`, which re-verifies the
exact owned row before deleting it and clearing runtime state
([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L627-L655)).

There is no boolean "asserted ready" that can replace these checks. The host
should expose a sanitized status assembled from the provisioner result,
`EvalRuntimeProofValid`, manager storage state, and local gateway availability.
At process startup the default must be `not-ready`; startup may serve ordinary
Sub2API traffic and the admin status page while evaluation workers remain
disabled. Do not read the old CPA `.env`, invent loopback/proxy inventory, or
trust an installer assertion as a proof. The authenticated probe
begin/complete operations must bind ownership to the initiating admin user and
return no raw token. Completion means only that the supplied observed proof
and normal Key row passed validation.

The topology primitive fingerprints forwarded-IP trust, recognized headers,
trusted proxies, whitelist, declared ingress peers, and server-observed
requests. `Observe` captures method/path, `RemoteAddr`, resolved security
client IP, and forwarded values; `Seal` requires every direct, public, spoof,
and declared-ingress case before HMAC sealing ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L45-L155),
[eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L310-L397),
[eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L525-L615)).
Validation rejects expiry, future skew, fingerprint drift, MAC changes, missing
coverage, and inconsistent observations ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L618-L670)).
The collector cannot detect an undeclared reverse proxy. Activation therefore
requires deployment evidence that the supplied inventory covers every real
proxy hop and its expected peer source; `Seal` proves the declared set only.
Completeness is an operator/deployment acceptance obligation, not a runtime
fact this primitive can discover.

The evidence also does not identify every public route independently. The
observation's `IngressID` is copied from the server-selected `EvalProbeCase`,
while `checkObservation` verifies the declared peer source and canonical
`/v1/responses` path; it receives no public hostname or upstream proxy route
identity ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L116-L135),
[eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L400-L453)).
Two external paths that arrive from the same peer and rewrite to the same path
cannot be distinguished by this proof. Use a fixed, protected operator-installed
inventory file as the concrete external acceptance input, not admin JSON or an
invented signing key. Candidate path: `EvalKeyStore.Dir()/deployment-evidence/`;
the directory must be owned by the Sub2API runtime UID and mode `0700`, and the
inventory file must be a regular non-symlink file owned by that UID and mode
`0400` or `0600`. Resolve this path from the existing private eval-store root,
never from request data or the old CPA `.env`. A missing file, wrong owner/mode,
symlink, malformed content, or stale binding fails closed. Installing or
replacing this file requires host/container filesystem privilege; that
privileged operation, recorded by the deployment's existing operator process,
is the attestation that external inventory verification was performed. A
human-name field inside JSON is descriptive only and is not an authenticity
check.

The deployment operator is the inventory producer and is responsible for
reviewing/installing the complete artifact. Patch 005 is a read-only consumer
that validates its fixed protected path and compares it with current facts.
Those current facts must come from a later trusted runtime/deployment provider,
which owns collecting and validating them; neither the inventory, an admin
request, nor artifact metadata can supply the expected current values.

The admin API supplies no inventory, destination, or `ready` boolean. `Begin`
loads the fixed protected file, computes SHA-256 over its exact bytes, validates
its bindings, and returns that digest with the required external case
descriptors and origins. The operator-run client receives a local copy of the
same artifact and must compare its digest before using any returned destination.
At minimum, the installed artifact must contain:

- schema/version and stable artifact ID (the artifact has no self-referential
  digest field), deployment/environment ID, source-system revision or immutable
  export digest, creation time, validity deadline, and the operator-process
  audit reference;
- the exact runtime/source/plugin and host-image provenance, active security
  configuration fingerprint/revision, and ingress-configuration revision it
  describes; if the running host cannot expose a value needed for comparison,
  it cannot claim this artifact is current;
- a complete entry for every externally reachable ingress hostname/listener,
  stable ingress ID, canonical request path and rewrite to `/v1/responses`, TLS
  hostname/policy, each proxy hop in order, and the expected immediate peer
  source observed by this service;
- one exact HTTPS origin and designated external runner network location for
  each remotely executed direct-public/ingress case, plus its associated
  ingress ID and expected peer/hop path. Runtime derives every direct-loopback
  origin from the actual listener and does not accept one from the file.

After a complete probe, write a separate runtime-owned session-evidence record
under the same protected directory. Bind it to the inventory digest, current
runtime/source/config revisions, session ID, completion/expiry times, exact
required case set, and a digest of the sealed proof/observations persisted by
the existing provisioning store. It may record each case's accepted result and
observed peer/client-source values, but never ticket values/hashes or the raw
evaluation token. On bootstrap, compare this record against the currently
installed artifact and the loaded provisioning proof; any mismatch, expiry, or
missing file keeps deployment-wide readiness false. The session record is
runtime-generated evidence, not a substitute for the operator-installed
inventory attestation.

The digest contract is deliberately byte-based: compute SHA-256 over the exact
protected file bytes before JSON parsing, and pass that digest as the
`evidenceDigest` request field. The host compares it with its own read; the
runner compares the same digest over its local artifact bytes; whitespace or
formatting changes therefore require a new artifact ID/approval. Do not put a
digest of the file inside the file itself or hash a non-specified JSON
serialization. The runtime-owned session-evidence record stores this exact
file-byte digest, not a reserialized JSON hash.

This file contract gives software concrete checks: owner/mode/path safety,
exact-file-byte digest, schema, current runtime/source/config revision, expiry,
complete required-case coverage, destination syntax, exact runner-to-case
mapping, and binding of the current session observations to the sealed proof.
It cannot prove that the
operator's source export is actually authoritative or complete, or that the
public origin corresponds to the claimed route. The authorized deployment
operator must independently verify every direct listener, DNS name, CDN/WAF,
load balancer, reverse proxy, rewrite, and peer hop before installing the file,
and invalidate/reinstall it when any such configuration changes. The running
application must keep local proof validity and installed inventory/session
evidence as distinct status fields; an admin request cannot set either by
assertion. Until the protected artifact and current session evidence are
present, status cannot claim deployment-wide topology readiness and evaluation
start remains unavailable. The collector does not discover omitted proxies.

The server's actual forwarded-IP behavior is live config: the snapshot reads
`ForwardedClientIPSettings()` and `TrustedProxies`, including the atomic live
settings override ([config.go](../../../../dist/sub2api-build/workspace/backend/internal/config/config.go#L730-L794)).
Do not change those host settings during evaluation initialization. A changed
trust flag, header list, or trusted-proxy list invalidates the sealed proof and
requires a new real probe. The probe producer must send the exact spoof value
for each recognized header and must not claim path identity beyond the observed
peer/path fields that the primitive records.

### Minimal administrator endpoint set

The execution plan already names the required operations. Keep them under the
existing admin middleware chain, with explicit status and no browser secret:

| Method | Path | Host responsibility |
| --- | --- | --- |
| `GET` | `/api/v1/admin/codex-eval/status` | Sanitized patch/runtime provenance, separate local-proof/inventory/session-evidence status, manager storage state, active runs, quiescing state, and probe-session progress. Never return raw token, stored lookup, ticket hashes, or ticket values. |
| `POST` | `/api/v1/admin/codex-eval/probes/begin` | Authenticated admin supplies only the protected inventory artifact ID and digest; host loads the fixed evidence file, derives actor/owner ID from context, creates one bounded ephemeral session and per-external-case tickets, executes all local direct-loopback cases itself, and returns external origins/tickets with `Cache-Control: no-store`. |
| `POST` | `/v1/responses` | A marker guard around the final `http.Server.Handler` admits marked traffic only for the exact method/path. That request reaches the ticket-aware hook at the start of the `/v1` chain, before API-key auth, target selection, forwarding, or billing. No marker follows normal routing unchanged. |
| `POST` | `/api/v1/admin/codex-eval/probes/complete` | Only the same authenticated admin `UserID` that began the session may complete it. Require all unique observations and unchanged fingerprint, seal internally, call `Initialize` with that same owner ID, atomically write session evidence bound to the inventory digest and proof digest, return sanitized identity/status, then discard session/tickets. |
| `POST` | `/api/v1/admin/codex-eval/probes/abort` | Only the same authenticated admin `UserID` that began the session may abort or replace it. Invalidate the collector and every ticket, retain a sanitized failure state, and allow that actor to begin again without revoking or rebinding the owned Key. |
| `GET` | `/api/v1/admin/codex-eval/targets` | Safe OpenAI group/account metadata. Re-check target membership at start. |
| `GET` | `/api/v1/admin/codex-eval/state` | Deep-copied manager snapshot, active runs, typed histories, storage errors, and quiescing state. |
| `POST` | `/api/v1/admin/codex-eval/start` | Validate request, require readiness, call manager `Start`, and return started/unavailable/busy summary. |
| `POST` | `/api/v1/admin/codex-eval/cancel` | Cancel selected run IDs; no network work in the handler. |
| `GET` | `/api/v1/admin/codex-eval/records/:type/:account_id/:run_id` | Return one typed record after positive ID/type validation and account/run lookup. |
| `POST` | `/api/v1/admin/codex-eval/clear-history` | Call manager `ClearHistory`; retain active intents and busy reservations. |
| `POST` | `/api/v1/admin/codex-eval/flush` | Explicitly retry a transient state save. Refuse protected corrupt/recovery-blocked states. |

The status endpoint must distinguish: runtime not provisioned, local topology
proof unverified/expired, protected inventory missing/stale/accepted, session
evidence missing/stale/accepted, external Key conflict, manager storage error,
probe session pending/failed/completion-pending/provisioning-recovery-required/
expired/aborted,
actor-bound replacement reserved,
workers running, and fully drained. It should not claim readiness merely
because the HTTP server is healthy. `probes/begin` and
`probes/complete` are separate operations because probes may have to travel
through multiple deployment paths between them. The authenticated abort route
is available when a ticket was consumed but its observation failed, or a
partial `Begin` did not finish. `complete` must refuse an incomplete or expired
session and must not accept an owner ID from request data. Neither operation
runs a paid evaluation or returns the private token. The `start` endpoint
should reject simple mode and any status that is not fully ready.

### Wire and dependency ownership

The current Wire graph is layered: `service.ProviderSet` supplies
`APIKeyService`, `AdminService`, gateway services, and settings
([service/wire.go](../../../../dist/sub2api-build/workspace/backend/internal/service/wire.go#L832-L872));
`handler.ProviderSet` builds top-level/admin handlers and the `Handlers` root
([handler/wire.go](../../../../dist/sub2api-build/workspace/backend/internal/handler/wire.go#L180-L230),
[handler/wire.go](../../../../dist/sub2api-build/workspace/backend/internal/handler/wire.go#L232-L295));
`server.ProviderSet` builds the router and HTTP server
([http.go](../../../../dist/sub2api-build/workspace/backend/internal/server/http.go#L23-L45),
[http.go](../../../../dist/sub2api-build/workspace/backend/internal/server/http.go#L112-L125)).
`cmd/server/wire.go` composes these providers and currently exposes only
`Server`, `PromptAudit`, `PluginManager`, and `Cleanup` through `Application`
([wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go#L27-L62)).

Keep ownership narrow and inject dependencies through Wire rather than globals:

1. Add a host service/adapter that owns `EvalKeyStore`, provisioner, topology
   probe producer, readiness snapshot, and the `TargetResolver` implementation.
   It may depend on `APIKeyService`, `AdminService`, `GroupRepository` or a
   narrow membership interface, `config.Config`, and the existing server
   address/config values. It must not expose raw credentials to handlers.
2. Add an evaluation admin handler containing the manager, runtime adapter, and
   safe DTO mapping. It should depend on the narrow host service interfaces,
   not on ent repositories or the entire gateway service.
3. Add the handler to `AdminHandlers`, register routes in `RegisterAdminRoutes`,
   and add it to `handler.ProviderSet`/`ProvideAdminHandlers`.
4. Add the manager and runtime host service to the application graph. Extend
   `Application` with the manager/runtime shutdown owner, update the Wire
   provider set, and regenerate `cmd/server/wire_gen.go`; do not hand-edit the
   generated file.

The local `HTTPExecutor` must target an explicit literal loopback origin and
must call the normal `POST /v1/responses` route. `ServerConfig.Address()` is
just the configured `Host:Port`, backed by `ServerConfig.Host` and `Port`
([config.go](../../../../dist/sub2api-build/workspace/backend/internal/config/config.go#L688-L690),
[config.go](../../../../dist/sub2api-build/workspace/backend/internal/config/config.go#L1539-L1541));
it is not proof of a dynamically assigned listener's actual port. The adapter
must derive a safe loopback origin from an explicit listener fact and refuse
ambiguous/nonlocal origins. For a configured fixed port, use the configured
port only after confirming the application listener uses it; for port `0`,
the server must expose/retain its actual bound listener address or readiness
stays false. Do not read the old CPA `.env`. If the route cannot be reached
through the local listener, keep the manager not-ready rather than relaxing
the executor's origin validation.

The host's `ProvideHTTPServer` currently begins with the Gin router, optionally
wraps it in `http.MaxBytesHandler`, then assigns the final `httpHandler` to
`server.Handler` ([http.go](../../../../dist/sub2api-build/workspace/backend/internal/server/http.go#L112-L163)).
Place the marker wrapper around that final composed handler, after existing
body-limit handling. Gin route middleware is insufficient: in v1.9.1,
`Engine.handleHTTPRequest` calls `c.Next()` only after a route match and returns
directly from trailing-slash and fixed-path redirect branches
(`gin.go:617-631`).

The current `runMainServer` calls `app.Server.ListenAndServe()` and therefore
does not retain the listener object or its actual bound address
([main.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/main.go#L171-L178)).
The lifecycle integration must bind with `net.Listen`, retain the resulting
`net.Listener` and its `Addr()`, then serve with `app.Server.Serve(listener)`.
Derive a literal `127.0.0.1:<actual-port>` candidate from the bound port only
after the listener is successfully established, and activate the local
executor only after a real self-request confirms that origin reaches the same
server. If either fact is unavailable, keep evaluation not-ready.

`NewHTTPExecutor` rejects an empty token and stores the raw token in the
executor ([http.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/http.go#L27-L56)).
`NewRunManager` rejects a nil `ModelExecutor`
([runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L98-L105),
[model.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/model.go#L28-L42)),
so its startup dependency must be a nonnil lazy proxy. The proxy returns a
typed `ModelError{Status: 503, Code: "not_ready"}` while inactive, and resolves
the current runtime readiness/delegate on each `Execute`; this keeps worker
execution fail-closed if proof or deployment evidence drifts after a start
request. Wire must not require
`NewHTTPExecutor(origin, token)` at startup. The runtime service may activate
the proxy with a real executor only after `BootstrapRuntime` or authenticated
`Initialize` verifies the owned identity and current proof, then independently
validates the protected inventory/session evidence and confirms that a
successfully bound `net.Listener` supplies a reachable literal loopback
origin. `ServerConfig.Address()` alone is not a listener fact, and an unknown
or unreachable loopback origin leaves the proxy inactive. The runtime
service owns the loaded raw token and real executor; handlers, state DTOs, logs,
and browser code never receive either. Synchronize activation/deactivation with
in-flight execution, retire idle connections only when safe, and close the
executor after manager drain and HTTP shutdown.

Evaluation manager/storage initialization must not make ordinary Sub2API
startup fatal. `NewRunManager` can return a nonnil inspectable manager together
with a load or recovery error
([runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L100-L138));
retain it and publish its protected storage status. If `NewFileStateStore`
rejects the path or cannot be constructed
([state.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/state.go#L49-L85)),
keep the runtime unavailable and report a sanitized status error rather than
substituting an empty store or overwriting existing bytes. The application
provider should expose a statusable blocked runtime even when no manager could
be constructed, so ordinary routes and the admin status endpoint remain
available while evaluation state/start operations stay disabled.

The collector's server observation contract is path-specific. `Observe` records
the actual request method/path, remote peer, resolved client IP, and recognized
forwarded values; `checkObservation` accepts only `POST /v1/responses`
([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L525-L571),
[eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L400-L455)).
The production probe producer therefore has to send each case through the real
listener and exact Responses route, with a server-side hook that calls
`Observe` before the normal gateway forwards or bills the request. A separate
admin-only probe URL, a claimed method/path in the payload, or a synthetic
`gin.Context` cannot satisfy this primitive. The probe label only selects a
predeclared case; it does not prove which proxy path handled the request.
The current `/v1` route stack attaches body limiting, request IDs, ingress
logging, endpoint normalization, API-key auth, group-model policy, target
platform routing, and then the Responses dispatch
([gateway.go](../../../../dist/sub2api-build/workspace/backend/internal/server/routes/gateway.go#L179-L220)).
There is no probe interception in that stack today. The host patch needs an
explicit narrowly scoped hook before forwarding; probe requests must terminate
there and never proceed to upstream selection/accounting. It must preserve the
normal path for ordinary requests, and readiness must remain false until every
planned case has been observed through the actual listener.

Producer authorization is a separate unresolved contract. `Observe(c, kase)`
checks that `kase` is in the collector's required set, the actual method/path
and network observations satisfy that case, and that the same case was not
already observed ([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L525-L571),
[eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L400-L455)).
It does not authenticate who supplied `kase`. `Seal` checks complete unique
coverage and config drift, not producer identity
([eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go#L574-L615));
`EvalIngressTopology.Validate` binds claims to the sealed inventory but also
does not authenticate the observer
([eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L336-L362)).
Consequently, a public request must never choose an arbitrary label and cause
the route to short-circuit normal forwarding.

### Probe Runner and Destination Contract (Candidate; Re-review Required)

Use a temporary operator-run client built only on the language standard
library, invoked for the initialization workflow and then exited. It is a
client process, not an installed daemon or sidecar. Prefer a statically built
Go client run from the approved external runner environment; a `docker exec`
invocation is only for a delegated local case when the official image lacks a
scripting runtime. The command must leave no resident process. It calls
authenticated `probes/begin` and `probes/complete` itself, using the same
administrator credential/identity for both requests; the existing admin
middleware continues to derive the actor and owner. Keep admin credentials in a
protected prompt or existing credential facility and send them only in the
normal authorization header over trusted TLS. Never place credentials in URL
userinfo, command arguments, or diagnostics. The client does not need or
receive the private evaluation token.

Divide probe execution by the network fact being tested. The host runtime may
exercise its direct-loopback case itself against the actual bound listener,
using a server-owned ticket that never leaves process memory. The operator
client exercises direct-public and declared-ingress cases from the approved
external network location named by the inventory artifact. A Docker-host
loopback or a public hairpin through the same proxy cannot establish the
server's direct-loopback source; likewise, running a public case from the
container cannot prove the external source/peer path named in the artifact. If
the runtime cannot perform its local probe, a one-shot delegated client must
run in the Sub2API container network namespace, not on the Docker host. The
runner's source location and build/version are provenance fields, and this
source limitation remains an external acceptance obligation.

The operator supplies a case-to-origin map from the current provenance artifact.
One approved external runner location must be able to reach the admin API and
every mapped public origin in the session. The client requires exactly one
origin for every externally-run case returned by the server and rejects extra,
missing, duplicate, or artifact-mismatched case entries. A host-owned
direct-loopback case has no external origin and its ticket is not returned; the
runtime records its own accepted result before completion. The server's `Begin`
implementation must not fetch or dial public origins; it only stores the
inventory and case bindings in the bounded in-memory session. The client appends
the fixed `/v1/responses` path to each origin and sends a minimal `POST` with the
one dedicated ticket header and only the case-specific spoof header when
required. It sends no `Authorization`, admin API key, private evaluation token,
or gateway API key to the probe destination. It does not follow redirects and
disables ambient HTTP proxy configuration. The deployed reverse-proxy path is
reached by connecting to its declared public origin; it is not a reason to
inherit an unrelated workstation proxy setting.

Destination validation is strict and case-aware:

- The server-owned direct-loopback case uses the actual bound listener port and
  a literal loopback origin selected by the runtime (for example,
  `http://127.0.0.1:<bound-port>`). The runtime sends its server-owned ticket
  through a real HTTP request to this address and keeps that ticket inside
  process memory. An origin derived only from configured `Host:Port` is
  insufficient, especially when the configured port is zero. If internal
  self-probing is infeasible and a one-shot client must perform it, that client
  must run in the Sub2API container network namespace, never on the Docker host.
- Direct-public and declared ingress cases use the exact origin assigned by
  the artifact. Require `https`, hostname verification against trusted system
  roots (or an explicitly approved deployment CA bundle with hostname
  verification), and a non-loopback destination; reject literal loopback IPs,
  URL userinfo, query/fragment, non-root path prefixes, unsupported schemes,
  and ports other than 443 unless the artifact explicitly records an approved
  TLS port. The request path remains exactly `/v1/responses`, as
  required by `Observe`.
- Fail the probe on any redirect. A redirect could change the destination or
  route while preserving a ticket-bearing request. Do not retry an ambiguous
  timed-out ticket; the request may already have consumed it, so discard the
  session and begin a fresh authorized session.

Each public ingress ID maps to its own declared origin, even if two origins
currently resolve to the same peer. This gives the operator client a concrete
destination for each required case, but does not make the collector prove
public-path identity: the provenance artifact and independent deployment
acceptance remain responsible for completeness and route-to-peer correctness.
The approved external runner/network performs the direct-public and ingress
cases, so the direct-public observation comes from an external source rather
than a container-to-public-origin hairpin. Exit the operator client after
completion; it is not a resident process. No browser-only workflow is required.
In particular, a
browser POST with the dedicated custom
ticket header is cross-origin and normally triggers a CORS preflight; host CORS
policy may reject it or cause the preflight to traverse a different route, and
the browser cannot verify the runner's direct transport policy. Keep tickets
out of browser code and avoid making CORS configuration part of topology proof.

### Probe Session Protocol (Candidate; Re-review Required)

Keep session control in the existing `/api/v1/admin` group and probe receipt on
the existing gateway route. A narrow runtime service interface candidate is:

```text
Begin(ctx, actorUserID, evidenceArtifactID, evidenceDigest) -> sessionID, expiresAt, case/ticket list
Observe(ginContext) -> whether request was an authorized probe, safe response/error
Complete(ctx, actorUserID, sessionID) -> sanitized readiness / owned-key identity
Abort(ctx, actorUserID, sessionID, optionalReplacementArtifact) -> sanitized session state
```

The admin handler must get `actorUserID` from
`middleware.GetAuthSubjectFromContext(c)` after `AdminAuthMiddleware`, require a
positive `UserID`, and never accept an owner ID in request JSON. `Begin` binds
the single active in-memory session and every ticket record to that user. The
gateway request is a bearer-capability request and does not independently
authenticate the human actor; only the begin/complete control-plane calls do.
`Complete` must compare its authenticated `UserID` to the recorded actor, then
pass exactly that ID to `EvalKeyProvisioner.Initialize`. A different admin may
inspect sanitized global status, but cannot complete, replace, or retrieve the
active session's ticket list. `Begin` rejects another live session until the
current session completes, is aborted by its owner, or expires.

Use one coordinator for session lifecycle and ticket state, with explicit
`building`, `active`, `failed`, `completing`, `aborting`, and terminal
`completion_pending`, `provisioning_recovery_required`, `completed`, and
`invalidated` states. `completion_pending` means `Initialize` succeeded but
writing session evidence failed. `provisioning_recovery_required` means
`Initialize` returned an error whose persistence effects may be partial; it is
an owner-bound recovery state, not a completed provisioning result. Reserve the
sole session slot before `Begin` reads the protected artifact or creates
tickets, then release the coordinator lock while doing file, network, or
service work.
Do not return any external tickets until local direct-loopback cases have all
succeeded and the session has changed from `building` to `active`. The local
self-requests need the same ticket observation path, so never hold the session
lock while sending them. A failed or cancelled partial `Begin` runs the same
invalidation primitive as an abort:
drop the collector and every ticket, retain only a sanitized failure code, and
release or actor-reserve the slot so the initiating admin can retry. Do not
leave an incomplete live session waiting for the ten-minute TTL.

Ticket lookup, one-time consumption, and `collector.Observe` must be serialized
against abort and completion through an observation lease. Under the coordinator,
atomically consume a ticket and increment the in-flight count; call `Observe`
without the mutex, then reacquire it and publish only if the same active
generation remains current. An observation must not be accepted after
invalidation. If an
observation fails, transition the session to `failed`, consume that ticket, and
invalidate its remaining tickets so no partial session can later complete.
Abort first marks `aborting`, which prevents new ticket consumption, then waits
for that count to reach zero before invalidating the collector and ticket map.
Do not hold the coordinator lock while calling `Observe`, doing network/file
I/O, or calling the provisioner. `Complete` atomically claims a single-flight
transition from `active` to `completing` only after all observations settle,
then releases the lock while it seals or provisions. Concurrent completion
requests return a safe conflict/busy result; abort, expiry, config-drift
cleanup, and replacement cannot mutate that claimed transition. Recheck session
generation, expiry, and the live fingerprint before publishing its result.
If sealing or constructing the topology fails before `Initialize` is called,
discard the session and all tickets, record a sanitized failure for the same
actor, release the completion claim, and require a fresh probe. No provisioner
side effect has occurred in this branch. Treat every `Initialize` error as
potentially partially committed: keep readiness false, discard the session and
all tickets, record `provisioning_recovery_required` for its owner, and release
the completion claim. Only that actor may begin an explicit recovery probe.
Before retrying, revalidate the exact inventory byte digest and artifact ID,
current source/runtime/deployment/config facts, topology fingerprint, and proof
expiry. Then call the existing `Initialize` path with the same owner and exact
protected inventory. Do not run a background retry or a paid model request.

The current provisioner can recover only when its persisted intent identifies
the same lookup and owner and the exact owned API-Key row exists and passes the
status, unbound-group, and whitelist checks. Its existing-intent branch calls
`adoptExistingRow`; it does not fall through to create a replacement row. In
particular, a `row_pending` intent whose API-Key creation failed and whose row
is absent returns `ErrEvalKeyRowMissing` on retry. Keep that state blocked for
explicit repair rather than promising automatic recovery or risking a
duplicate. Lookup, owner, status, binding, and IP-whitelist conflicts are also
blocked for explicit repair. Never revoke, rebind, reenable, or otherwise
rewrite the Key to make recovery pass. An authenticated abort may clear only
ephemeral session and ticket state; it preserves the persisted intent and any
owned row. A fresh session for the same owner may retry an existing-row
recovery after all evidence checks pass.

Only successful `Initialize` followed by a session-evidence write failure uses
`completion_pending`. Preserve its owner-bound completion record and proof for
an idempotent evidence-write retry, discard tickets, and do not reopen
observations or rerun `Initialize`. Before that retry, reload the exact
protected artifact and current facts, compare artifact ID and exact-byte digest,
validate the proof MAC, live fingerprint, and expiry, and strictly reload the
exact owned Key row with its owner, status, binding, and whitelist checks. The
retry performs evidence persistence only; it does not call `Initialize`, create
or change a Key, or issue a model request. Other actors and replacement remain
blocked. An expired or drifted session invalidates through the same coordinator
once no completion owns the transition. Process restart starts with an empty
ephemeral ticket table and invalidates any residual in-memory session state
during runtime cleanup; tickets are never persisted. After restart, derive
recovery ownership from the persisted provisioning intent's `OwnerUserID`. Do
not infer successful completion from that intent when runtime session evidence
is missing: require a fresh probe by that same owner unless a distinct durable
completion binding was persisted. Only a validated durable completion binding
may support evidence-only retry after restart.

The abort request takes only a session identifier (and, for a replacement
transaction, the new protected artifact ID/digest); it never takes an actor or
owner ID. The handler derives the actor from `GetAuthSubjectFromContext`. Under
the lifecycle coordinator, abort verifies both session ID and the exact
positive initiating `UserID`, prevents new ticket consumption, waits for any
in-flight observation to settle, clears all ticket records and collector state,
and stores a non-sensitive reason such as `observation_failed`, `aborted`,
`expired`, or `config_drift`. A failed observation consumes its ticket and
leaves that sanitized failure visible until the same actor aborts or replaces
the session. Preserve a short-lived replacement reservation for that actor
after explicit abort, partial-`Begin` cleanup, or a provisioner error, so
another actor cannot take the slot between invalidation and the owner's retry.
Another admin cannot abort or replace it. A mismatched actor gets a denial
without revealing the ticket list or failure details. Use a same-actor
replacement transaction or this bounded slot reservation.
Abort/replacement never calls `RevokeOwnedKey`, changes Key ownership, binds the
Key to a group, or clears a successfully provisioned identity. Only the
separately reviewed uninstall path may revoke the owned Key.

Candidate sequence:

1. The operator-run client submits the protected provenance artifact ID and
   exact-file-byte SHA-256 digest to authenticated `POST
   /api/v1/admin/codex-eval/probes/begin`; the host does not accept inventory or
   destination data from the request. It loads the fixed evidence file, snapshots
   the current config fingerprint, creates `EvalTopologyCollector`, obtains
   `RequiredCases()`, validates exact external-case map coverage and destination
   rules, then creates random one-use tickets for the external cases. It derives
   one ticket per required direct-loopback source from the actual bound listener,
   sends real local HTTP requests through the same marker hook, and retains
   those tickets only inside process memory. Proposed bounds are one live
   session, at most 256 cases, and a 10-minute session/ticket TTL; the numeric
   limits are provisional and require operational/security review. The response
   is `Cache-Control: no-store` and contains the session ID, expiry, non-secret
   external case descriptors, validated public-origin mapping, and each
   external ticket once. It contains no local ticket and no raw evaluation
   token. If the real local request fails, begin fails closed and exposes no
   external ticket set for a partially viable session.
2. The same operator client that called `Begin` keeps external tickets only in
   memory and sends each only in the dedicated `X-Codex-Eval-Probe-Ticket`
   header on `POST /v1/responses` to the exact mapped origin. The server never
   accepts a case label from the probe request. A server-side ticket lookup
   supplies the collector/session, actor binding, fingerprint, exact case,
   expected spoof value where relevant, and expiry.
3. Wrap the final `http.Server.Handler` after normal handler composition with a
   narrow `net/http` marker guard. Detect the reserved marker name
   case-insensitively and require exactly one syntactically valid ticket value.
   Encode 32 random ticket bytes with unpadded base64url; reject an empty,
   repeated, noncanonical, or malformed value. With no marker, call the
   original handler unchanged, preserving existing Gin routing and redirects.
   With a marker, reject unless `Request.Method == POST` and
   `Request.URL.Path == "/v1/responses"`; only an exact match may enter Gin and
   reach the ticket-aware probe hook. This outer boundary is required because
   Gin 1.9.1 can emit trailing-slash or fixed-path redirects inside
   `Engine.handleHTTPRequest` before route handlers run. Put the probe hook at
   the start of the `/v1` middleware chain, before body/auth/target/forwarding
   middleware, and return immediately after a valid observation. All marked
   methods, paths, aliases, providers, protocols, and missing/repeated/malformed/
   unknown/expired/consumed tickets fail before API-key auth, account selection,
   forwarding, or billing. On the exact route, validate and consume the ticket
   under the coordinator and increment the session's in-flight observation
   count. Release the coordinator mutex before calling
   `collector.Observe(c, serverBoundCase)`. Reacquire it afterward and publish
   the result only if the same session generation remains current. Failed
   observations consume the ticket and retain a sanitized failed-session state
   for same-actor abort/replacement. Abort prevents new observation leases and
   waits for outstanding leases before invalidating the collector. The probe
   response terminates with no normal gateway call.
4. After every external case returns an accepted observation and the runtime's
   direct-loopback observation is recorded, the same client calls
   `POST /api/v1/admin/codex-eval/probes/complete`; it requires the same positive
   authenticated `UserID`, the current session ID, unexpired session, unchanged
   fingerprint, and exactly one successful observation for every required
   case. A single-flight lifecycle transition owns `Seal`, `Initialize`, and
   session-evidence persistence. Incomplete sessions return a conflict with
   remaining case identifiers but no ticket values and do not call `Seal` or
   `Initialize`; the initiating actor can use the authenticated abort route to
   invalidate the failed session and start again. On complete
   coverage, the runtime service loads/ensures the private token internally,
   calls `Seal(token)`, builds `EvalIngressTopology` from the protected inventory
   and proof, and calls `Initialize(ctx, actorUserID, topology)`. A sealing or
   topology-construction failure invalidates the attempt and requires the same
   actor to re-probe. Every `Initialize` error records
   `provisioning_recovery_required`, keeps readiness false, clears ephemeral
   tickets, and releases the completion claim while retaining the actor-bound
   reservation. A same-actor explicit retry must reload the same exact-byte
   artifact, verify all current revisions and fingerprint/expiry, seal a fresh
   proof, and re-enter `Initialize`. Recovery may adopt only the exact existing
   owned row. If that row is absent from a persisted `row_pending` intent, or a
   lookup/owner/status/binding/whitelist check conflicts, return a blocked
   repair status without creating, revoking, rebinding, or reenabling a Key.
   After successful initialization, atomically write the runtime-owned
   session-evidence record bound to the artifact digest and sealed-proof digest.
   If that write fails, report evaluation not-ready and permit only an
   idempotent evidence-write retry by the same actor. Before retrying, reload
   the exact artifact/current facts, validate artifact ID and exact-byte digest,
   proof MAC/fingerprint/expiry, and strictly reload the exact owned Key row.
   Retry only evidence persistence; never call `Initialize` or issue a model
   request. Do not allow abort/replacement while this transition is in flight.
   At restart, recover ownership from persisted `OwnerUserID`, and require a
   fresh same-owner probe when evidence is missing unless a separate durable
   completion binding was persisted. Do not infer completion from provisioning
   intent alone. Return only sanitized status/identity, never the token or
   stored lookup.
   Complete, abort, expiry, config drift, partial-Begin cleanup, and process
   restart all invalidate through the same lifecycle primitive; ticket values
   and hashes remain memory-only. Re-probing after a failed observation starts
   a new authorized session for the same actor.

Ticket generation should use `crypto/rand` with at least 32 bytes of entropy.
Store only `SHA-256(ticket)` and bounded metadata in memory: session ID,
actor `UserID`, collector instance/fingerprint, server-selected case, expected
spoof value, expiry, and consumed/available state. Atomically remove or mark
consumed before observation. Do not write ticket values or hashes to
`provisioning.json`, the sealed proof, run state, audit records, request bodies,
URLs, command arguments, or logs. Return the opaque ticket only in the
authenticated no-store response and require the runner to send it in the
dedicated header. The approved ingress deployment must also redact that header;
the application cannot inspect proxy logs.

The marker name must be compared case-insensitively against
`collector.RecognizedHeaders()` at session creation and again when validating
the active session. If it collides with any configured/built-in forwarded-IP
header, fail session creation closed without changing global trust settings.
This is required because `Observe` copies all recognized header values into the
observation, which is later persisted in the sealed proof. Also assert that no
marker or ticket value enters `EvalIngressObservation.Forwarded` or any other
proof field. Invalid-marker rejections should be marked through the existing
ingress-rejection path so Ops error-body capture skips them; valid probes return
before authentication and do not create billing or usage records.

The control/API sequence is source-compatible with the existing collector and
provisioner, but all session/ticket/middleware and runner types above are
candidates only. No current source provides producer authentication, session
expiry/replay state, ticket distribution, destination-map validation, or a
marker guard. The proposed operator manifest resolves the missing runner input;
it does not change `EvalIngressPeer`, and neither destination origins nor
attestation data belong in the sealed proof or provisioning file. Independent
re-review is required before implementation.

Required host assertions: another admin cannot complete, abort, replace, or
retrieve tickets; failed/partial `Begin` does not strand a live session;
ticket consumption/`Observe` races safely with abort and completion; two
completion requests cannot run `Initialize` concurrently; abort does not revoke
or rebind a provisioned Key; request-supplied actor/case labels cannot override
server state; tickets bind one session/fingerprint/actor/case and are
single-use; missing, empty, malformed, repeated, expired, wrong-session,
wrong-route, consumed, and case-mismatched markers add no observation;
collision with any recognized forwarded header prevents session creation;
invalid markers cause no auth lookup, upstream call, usage row, or charge;
valid probe tickets also cause no auth lookup/forwarding/billing. Exercise the
final handler through an actual bound listener, not only `gin.Engine.ServeHTTP`:
cover a trailing-slash redirect, fixed-path alias when enabled, a no-route 404,
an unsupported-method 404/405, `/v1/responses/*subpath`, root `/responses` and
`/api/v3/responses` aliases, `/backend-api/codex/responses`, Google
`/v1beta/models` and `/antigravity/v1beta/models`, and the `GET /v1/responses`
and `GET /responses` WebSocket upgrade paths.
For every marked non-exact request, assert the outer marker guard rejects it
before auth lookup, probe observation, forwarding, usage creation, or billing.
For exact marked `POST /v1/responses`, assert one ticket causes at most one
observation and never reaches normal auth/forwarding/billing. Also send
unmarked variants through the same listener and assert existing redirects,
404/405, Google, alias, and WebSocket behavior is unchanged. Capture logs and
persisted proof to prove ticket values, hashes, and raw evaluation tokens are
absent.

### Startup and shutdown ordering

`runMainServer` initializes the application, starts long-lived services, starts
`app.Server.ListenAndServe` in a goroutine, waits for a signal, then currently
calls `app.Server.Shutdown` with a five-second context before deferred
`app.Cleanup` ([main.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/main.go#L134-L195)).
`provideCleanup` stops application services in parallel and closes Redis and
Ent only afterward ([wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go#L84-L145),
[wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go#L397-L450)).

Evaluation workers issue HTTP requests through the local gateway. Therefore the
manager must be drained before `app.Server.Shutdown`; otherwise its executor
loses the local route while cancellation, stream closure, usage accounting, or
terminal persistence is still in flight. The host shutdown sequence should be:

1. Mark the manager quiescing and cancel all active runs.
2. Wait for the manager with a bounded drain context while the local HTTP
   server, OpenAI gateway, WebSocket pool, Redis, and Ent remain available.
3. If that context expires, treat it only as a drain retry signal. Stay inside
   the process lifecycle, leave the server accepting ordinary traffic, keep the
   manager quiescing so new evaluations are rejected, log a sanitized
   “workers still draining” status, wait a short bounded backoff, and call
   `Shutdown` again with a fresh deadline context. Repeat until drain completes
   or the process receives an external hard kill. Never return through
   `runMainServer`'s cleanup path, call `Server.Shutdown`, close the WebSocket
   pool, or run `app.Cleanup` while workers remain tracked.
4. Once drain completes, call `app.Server.Shutdown` with its own bounded context.
5. Run cleanup only after the completed drain and HTTP shutdown; the existing
   cleanup may then close runtime/executor resources and application services,
   followed by Redis and Ent. Do not move worker draining into the parallel
   cleanup list, because it would race the existing parallel gateway-pool close.

This retry behavior is required by both the current host and manager source.
`runMainServer` installs `defer app.Cleanup()` before it starts serving, and
today returns after a single five-second `Server.Shutdown` attempt
([main.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/main.go#L151-L195)).
That cleanup closes the OpenAI WebSocket pool in the parallel application
steps, then Redis and Ent in the infrastructure steps
([wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go#L335-L339),
[wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go#L397-L450)).
Therefore merely skipping `Server.Shutdown` on timeout is insufficient; the
function must not return and trigger its defer. Prefer explicit cleanup after
the drain loop and server shutdown, or structure the loop so every timeout
continues within `runMainServer`.

`RunManager.Shutdown(ctx)` marks quiescing, cancels active runs, cancels its
lifetime context, and waits on `drained`; that channel closes only after both
tracked runs and pending starts settle. Its deadline branch returns `ctx.Err()`
directly, while its drained branch returns the stored terminal error
([runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L656-L688)).
That stored error wraps `Store.Save` failures and may itself satisfy
`errors.Is(err, context.DeadlineExceeded)` or
`errors.Is(err, context.Canceled)`. Therefore the host must not classify the
result with `errors.Is` alone. For each bounded drain call, retain the exact
context and inspect its `Err()` before canceling it. A result identical to that
call's direct `ctx.Err()` is a drain retry signal. If the context/error state
does not distinguish a deadline return from a context-shaped storage error,
retry `Shutdown` immediately with a fresh live bounded context. A non-nil error
returned while that fresh context is still live is a settled persistence error
under the current manager control flow, even if it wraps a context sentinel;
report “drained with storage error” using sanitized `State().StorageError`, then
continue HTTP shutdown and cleanup. If the fresh call returns its own direct
`ctx.Err()`, workers or pending starts still have not drained; remain in the
retry loop. Do not label a context-wrapped storage failure “workers still
running” or “shutdown timed out.”

Do not infer drain solely from `ManagerState.Active`: pending starts are tracked
separately and are not exposed in that snapshot. A startup save blocked before
worker launch can leave `Active` empty while `Shutdown` is still waiting for
`pendingStarts`; only successful `Shutdown` return after the drain signal proves
both classes have settled. Keep each drain context alive through result
classification, then cancel it.

Host-level fixtures: block an executor past the first drain deadline; assert the
process loop has not returned, `app.Cleanup` has not run, HTTP `POST
/v1/responses` and its required dependencies remain usable, and manager starts
are quiesced. Release the worker; call `Shutdown` with a fresh context and assert
the order manager drained -> HTTP shutdown -> gateway/dependency cleanup. Also
block the initial startup `Store.Save` before a worker exists; assert
`ManagerState.Active` is empty but the host still keeps HTTP/dependencies alive
until the pending-start transaction settles. Finally make a terminal save wrap
`context.DeadlineExceeded` or `context.Canceled` while the drain context remains
live; assert the host classifies it as a settled storage error, not a drain
timeout, and cleans up only afterward. Exercise an ambiguous first-call deadline
with a fresh-context retry as well. Stage 4B's independent specification and
quality gates are recorded as passed in
[eval-run-state-spec-review.md](eval-run-state-spec-review.md) and
[eval-run-state-quality-review.md](eval-run-state-quality-review.md); these
host drain fixtures remain a separate integration acceptance gate.

### Ordered implementation chunks and gates

The following order keeps each seam independently reviewable:

1. **Stage 4B core (complete).** Implementation and distinct independent
   specification and quality reviews passed according to
   [eval-run-state-plan.md](eval-run-state-plan.md). Do not repeat this chunk.
   The host still must preserve save-before-worker/ack, manager-owned contexts,
   exact target resolver identity, all-type cancellation, `Flush`, corrupt-file
   protection, retention, recovery, and shutdown deadline/drain behavior.
2. **Host metadata adapter.** Add a narrow resolver/target DTO service using
   the existing admin/group/account APIs. Gate on positive numeric IDs, exact
   OpenAI membership, account/group mismatch failures, credential-free JSON,
   and a target-list snapshot that cannot authorize a later start by itself.
3. **Protected inventory reader (next isolated patch 005).** Add only
   `internal/service/eval_deployment_evidence.go` and its focused test. Resolve
   the fixed `EvalKeyStore.Dir()/deployment-evidence/inventory.json` path; check
   private path ownership/modes/type/symlinks; hash exact bytes before strict
   bounded JSON decoding; compare explicit current revision facts; and validate
   exact external-case coverage and safe HTTPS destinations. Do not add a file
   writer, session evidence, tickets, routes, Wire, bootstrap, or lifecycle
   changes. Replay patches 001-005 on a pristine pinned source tree and compare
   the new files byte-for-byte with the maintained workspace.
4. **Actor-bound session coordinator.** In a separate patch after the reader,
   add only the coordinator and its focused tests. Gate on actor ownership,
   begin reservation, single-use ticket leases, generation checks, abort wait,
   completion single-flight, expiry/config-drift invalidation, and the
   owner-bound recovery states below. Do not add HTTP, runner, gateway, Key
   bootstrap, or lifecycle wiring.
5. **Topology probe producer.** Add the server-only probe workflow and later
   runtime-owned session-evidence writer under the private evidence directory.
   Gate on file owner/mode/path safety, source/runtime/config revision and
   expiry, exact case-to-origin coverage, direct-loopback listener facts,
   external runner network/source provenance, remote TLS/redirect/proxy policy,
   direct-public and forwarded-loopback spoof cases, and every declared ingress
   peer/header case.
   Require observed `RemoteAddr` and resolved client IP, plus authenticated
   one-use case authorization (candidate above, pending review). Readiness
   requires external deployment evidence that the inventory covers every
   actual proxy hop, and remains false when that evidence is absent, config
   drifts, or proof expires. The collector cannot discover omitted hops. Do not
   modify host trust settings.
6. **Key bootstrap/readiness.** Wire `BootstrapRuntime` during application
   startup and `Initialize` behind the authenticated probe-complete operation.
   Gate on owner/status/binding/lookup/IP conflicts, idempotent renewal,
   raw-token redaction, stored-lookup replay rejection, and normal API-Key
   policy. A failed bootstrap must leave ordinary host startup available but
   evaluation unavailable.
7. **Admin handler/routes.** Add status, probe begin/complete, targets, state, start,
   cancel, records, clear-history, and flush operations under the existing
   `/api/v1/admin` middleware chain. Gate on `{code,message,data}`, admin-only
   rejection, auth expiry, no token in responses/URLs/logs, and explicit
   unavailable/busy/error states.
8. **Wire/application lifecycle.** Add provider declarations and generated Wire
   output, extend `Application`, and put manager drain before HTTP shutdown.
   Gate on a test executor that proves the local gateway remains reachable
   until workers settle, deadline behavior, terminal save errors, and final
   dependency close order.

### Next Source Ownership (Candidate; Not Yet Implemented)

The next host patch should own these files in the pinned Sub2API workspace;
names for new files are proposed and require the same review as the protocol:

- Patch 005 owns only new `internal/service/eval_deployment_evidence.go` and
  `internal/service/eval_deployment_evidence_test.go`: the concrete exported
  reader/schema, protected inventory reading, and pure validation. Enforce the
  1 MiB file cap, 256 derived-case cap, strict duplicate-free JSON, Unix-second
  creation/expiry rules, and anchored `O_NOFOLLOW` traversal of every absolute
  path component. It has no writer, API, readiness, or startup side effects.
  The operator produces the artifact; this patch consumes it; a later trusted
  runtime/deployment provider owns current-fact sourcing and validation.
- A later evidence-integration patch may own runtime session-evidence
  persistence, artifact/proof digest binding, and sanitized readiness inputs.
  Extend `internal/service/eval_key_runtime.go` only where bootstrap must
  compare current protected evidence with the persisted proof.
- `internal/service/eval_topology_runtime.go` (new): actor-bound begin/complete
  session, collector ownership, external ticket state, and the real local
  loopback self-probe. Do not change the frozen observation semantics in
  `eval_topology.go` without a separate review.
- `internal/server/middleware/eval_probe.go` (new) plus
  `internal/server/routes/gateway.go`: ticket lookup/consumption and exact
  `/v1/responses` observation hook at the start of the gateway chain. The outer
  marker guard wraps the final handler in `internal/server/http.go` so it runs
  before Gin redirects, no-route, and method handling.
- `internal/handler/admin_codex_eval.go` (new),
  `internal/server/routes/admin.go`, and handler/service Wire declarations:
  authenticated status, begin, complete, same-actor abort/replacement, and
  existing evaluation operations.
- `cmd/codex-eval-probe/main.go` (new, standard library only): one-shot static
  operator client that checks the local artifact digest, obtains tickets from
  authenticated begin in memory, and probes only the artifact's public origins.
- `internal/server/http.go`, `cmd/server/main.go`, `cmd/server/wire.go`, and
  regenerated `cmd/server/wire_gen.go`: expose the actual `net.Listener` fact to
  the runtime, keep local HTTP alive through manager drain, then shut down and
  clean dependencies in order.

The file-readiness gate must reject missing, stale, unsafe, or revision-mismatched
inventory/evidence before returning fully ready; no admin request may upload
inventory or assert readiness. The probe gate must prove that the local
self-request uses the actual bound port and that the operator client enforces
artifact digest, exact origin/path, trusted TLS, redirects-off, ambient
proxies-off, and credential isolation. The lifecycle gate must prove a blocked
worker and a blocked pending-start save both keep the listener/dependencies
available, and must distinguish direct drain-context errors from wrapped
terminal-save context errors. The operator still has to verify the real
deployment inventory and external source locations before installing the
protected artifact; these code gates cannot establish that external fact.

### Focused fixtures and limits

Required host fixtures should include:

- admin JWT and admin API-key requests, non-admin rejection, malformed
  envelopes, and sanitized status/target DTOs;
- target-list membership changes between `/targets` and `/start`, exact account
  and group IDs in directed request headers, and normal billing/subscription
  rejection responses;
- topology probes with each configured forwarded header, public clients that
  spoof loopback, direct loopback, declared proxy peers, missing/partial
  inventory evidence, config fingerprint drift, proof expiry, and inventory
  conflicts; also unauthorized case selection, ticket expiry/replay/duplicate,
  case mismatch, and ordinary route pass-through;
- actor lifecycle races: another admin cannot abort/replace, consumed failed
  probes can be cleared by the initiating actor, partial `Begin` cleans all
  tickets, abort cannot race a committed observation, completion/provisioning
  is single-flight, and replacement preserves Key ownership/binding;
- actual-listener probe guard checks for marked trailing-slash/fixed-path
  redirects, no-route 404, unsupported-method 404/405, Responses subpath
  aliases, Google/Antigravity Google paths, and WebSocket upgrades; assert no
  auth lookup, `Observe`, forwarding, usage, or billing for non-exact requests,
  and unchanged behavior for their unmarked counterparts;
- reject a missing/incorrectly owned/mode-open/symlinked/malformed/stale
  protected evidence file; reject an admin request that attempts to upload
  inventory or assert readiness; verify session evidence digest binds the
  loaded file, exact current cases, and sealed proof, and that an evidence write
  failure leaves start unavailable; verify an evidence-only retry reloads and
  validates current facts, proof, and exact Key without calling `Initialize`;
- restart recovery derives the actor only from persisted `OwnerUserID`; absent
  session evidence requires a fresh same-owner probe unless a distinct durable
  completion binding exists, and provisioning intent alone never proves
  completion;
- runner case-map coverage and provenance-revision mismatch; reject public
  literal-loopback destinations, URL userinfo, query/fragment, path prefixes,
  unsupported schemes/ports, untrusted TLS roots, and redirects; prove ambient
  proxies are disabled and no admin/private/gateway credential reaches the
  probe origin; timeouts with ambiguous ticket consumption require a new
  session. Keep the browser/CORS path out of acceptance criteria;
- bootstrap with missing intent, missing row, owner/status/binding/IP/lookup
  conflicts, expired proof, idempotent renewal, and owned-key revoke failure;
- manager startup save failure, HTTP request cancellation after commit, all
  four evaluation cancellations, terminal write failure plus successful
  `Flush`, recovery-save protection, shutdown timeout/drain, pending-start drain
  while `Active` is empty, and context-wrapped storage errors distinguished from
  the exact drain-context error;
- Wire startup checks with no provisioned token: the nonnil proxy remains
  typed-not-ready, the manager still loads state and completes recovery, the
  actual listener can be captured, and ordinary host/admin status routes serve;
  after verified bootstrap/initialization plus listener evidence, activation
  uses only the literal loopback origin. Also verify manager load/recovery
  errors remain statusable and an invalid `FileStateStore` path does not create
  an empty fallback or overwrite prior bytes;
- Wire construction checks that include the manager and runtime service without
  changing ordinary service provider ownership.

These fixtures prove host contracts and local lifecycle only. They do not prove
deployment proxy inventory, a production container's source-IP behavior,
live usage accounting, or a paid model result. Real deployment evidence is
indispensable for the actual source/peer inventory and the final authenticated
gateway accounting check.

## Files Found

- [host-runtime-plan-refined-review.md](host-runtime-plan-refined-review.md):
  independent review identifying the outer marker-guard, same-actor
  abort/replacement, and executable lazy-executor refinements.
- [implement.md](../implement.md): Stage 4B contract, endpoint list, and
  shutdown requirement.
- [eval-run-state-plan.md](eval-run-state-plan.md): authoritative manager,
  state, `Flush`, recovery, and drain interfaces.
- [private-key-lifecycle.md](private-key-lifecycle.md): raw-token storage,
  normal Key registration, conflict rules, and topology boundary.
- [patch-contracts.md](patch-contracts.md): frozen Stage 2A/2B/2C patch
  owners and deferred admin/bootstrap wiring.
- [eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go):
  private-Key store, provisioner, bootstrap, target membership, and revoke.
- [eval_topology.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_topology.go):
  observed ingress proof, inventory binding, freshness, and live-config checks.
- [deployment-evidence-reader-plan.md](deployment-evidence-reader-plan.md):
  isolated patch 005 scope, schema, source facts, and validation gates.
- [host-runtime-refinement-response.md](host-runtime-refinement-response.md):
  completion-error recovery and observation-lease refinements from the final
  plan review.
- [runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go):
  current manager shutdown return branches and pending-start drain ownership.
- [model.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/model.go):
  exported `ModelExecutor` contract and safe typed model-error shape for a
  not-ready lazy proxy.
- [state.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/state.go):
  guarded `FileStateStore` constructor and load/save behavior; unsafe storage
  must not be replaced with an empty fallback.
- [runs_test.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs_test.go):
  manager shutdown/deadline and unavailable-only pending-start drain fixtures.
- [http.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/http.go):
  token-dependent executor construction, literal-loopback validation, and the
  fixed `/v1/responses` request path.
- [gateway.go](../../../../dist/sub2api-build/workspace/backend/internal/server/routes/gateway.go):
  current `/v1` middleware and Responses dispatch order; no probe hook exists.
- [evalroute.go](../../../../dist/sub2api-build/workspace/backend/internal/pkg/evalroute/evalroute.go):
  canonical route path and directed target header contract.
- [api_key_auth_eval.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth_eval.go):
  private-token mapping, directed-header validation, and readiness gate.
- [api_key_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth.go):
  normal Key ACL, group, subscription, quota, balance, and context setup.
- [admin_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/admin_auth.go):
  administrator authentication and trusted request identity.
- [admin.go](../../../../dist/sub2api-build/workspace/backend/internal/server/routes/admin.go):
  shared administrator middleware chain and route registration owner.
- [router.go](../../../../dist/sub2api-build/workspace/backend/internal/server/router.go):
  Gin composition, shared admin route registration, and gateway route setup.
- [http.go](../../../../dist/sub2api-build/workspace/backend/internal/server/http.go):
  final `http.Server.Handler` composition and global body-limit wrapper; the
  outer marker guard belongs at this server boundary.
- [admin_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_service.go),
  [admin_account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_account.go),
  [admin_group.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_group.go):
  metadata and group/account lookup seams.
- [wire.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/wire.go),
  [main.go](../../../../dist/sub2api-build/workspace/backend/cmd/server/main.go),
  [http.go](../../../../dist/sub2api-build/workspace/backend/internal/server/http.go):
  application Wire graph, current `ListenAndServe` startup, shutdown, and
  dependency close order; runtime integration must retain the actual listener.
- [go.mod](../../../../dist/sub2api-build/workspace/backend/go.mod):
  pinned upstream Go version and Gin v1.9.1 dependency; the proposed probe
  client uses only the standard library.

## Related Specs

- [evaluation-core.md](../../../../.trellis/spec/backend/evaluation-core.md):
  Stage 4A API and ownership boundary.
- [run-state.md](../../../../.trellis/spec/backend/run-state.md):
  pending-start tracking, shutdown error distinction, and host drain obligation.
- [error-handling.md](../../../../.trellis/spec/backend/error-handling.md):
  typed host/management/evaluation errors and safe API responses.
- [quality-guidelines.md](../../../../.trellis/spec/backend/quality-guidelines.md):
  checks and persistence/error visibility expectations.

## External References

- Gin v1.9.1 `gin.go` (`Engine.handleHTTPRequest`, lines 592-649), pinned in
  the Sub2API `go.mod`: matching routes call `c.Next`, while trailing-slash and
  fixed-path redirect branches return before route handlers; this is why the
  marker guard must wrap the actual `http.Server.Handler`.

## Caveats / Not Found

- The pinned workspace contains the Stage 4B manager core but no production
  host runtime/provider that owns it, administrator evaluation handler, probe
  producer, topology inventory source, or Wire binding. The proposed seams
  above are source-backed integration points, not implemented runtime claims.
- The topology collector cannot authenticate the probe producer, bind an
  observation to a one-use probe ticket, or discover an undeclared proxy. The
  source-compatible ticket/session and operator-runner scheme above, including
  the outer marker guard, actor-bound abort lifecycle, and lazy executor
  activation, remains a candidate for independent design review; do not treat
  it as an approved interface. No
  current source defines or verifies the provenance artifact/trust root or
  destination map. Activation requires an authorized external acceptance
  record covering every proxy hop, peer source, and per-case destination bound
  to the current deployment revision. Never use the old CPA `.env` or an
  unverified installer value.
- No source evidence supports changing global trusted-proxy settings,
  creating a new Key schema, changing business group bindings, or asserting a
  global host-ready boolean. Those remain out of scope.
- Docker is available through `desktop-linux` (Docker Desktop 29.8.2,
  linux/arm64). No approved target, live probe, live accounting request, or paid
  model evaluation was used for this read-only research.
