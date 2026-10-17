# Refined Review: Sub2API Host Runtime Wiring Plan

- Date: 2026-10-12
- Scope: follow-up review of the host runtime wiring, probe session, startup,
  and metadata-adapter plans against the pinned Sub2API source.
- Pinned source: Sub2API `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.
- Verdict: **NEEDS-REFINEMENT** for host runtime implementation.

## Findings (fixed)

- None in product code. This was a read-only plan review; this report is the
  only file added.

## Findings (not fixed)

### 1. Put the marker guard outside Gin route dispatch

- File: `research/host-runtime-wiring-plan.md`, probe session protocol and
  `internal/server/middleware/eval_probe.go` ownership.
- Issue: The plan places the marker guard in the gateway middleware chain and
  requires marked requests on aliases, unsupported methods, and no-route paths
  to fail before normal routing. Gin 1.9.1 performs trailing-slash and fixed-
  path redirects inside `Engine.handleHTTPRequest` before it calls the route
  handler chain (`github.com/gin-gonic/gin@v1.9.1/gin.go`, lines 624-631).
  A route middleware therefore cannot guard a marked request that takes one of
  those redirect paths. The marker could receive a redirect response or reach
  another handler without the promised exact-route rejection semantics.
- Required refinement: wrap the actual Gin handler in a narrow `net/http`
  handler at server composition time. If the reserved marker is absent, call
  the original handler unchanged. If it is present, reject unless the request
  is exactly `POST /v1/responses`; only then invoke the ticket-aware probe hook.
  Keep ordinary Gin redirect behavior unchanged for requests without the
  marker. Add no-route, trailing-slash/fixed-path, unsupported-method, alias,
  WebSocket, and Google-route tests asserting no auth lookup, forwarding, or
  billing.
- Why not fixed: this changes the HTTP ownership boundary and must be designed
  with the later server/Wire lifecycle patch rather than silently adding a
  second handler path in the plan review.

### 2. Add an actor-bound abort/replacement path for failed probe sessions

- File: `research/host-runtime-wiring-plan.md`, `Probe Session Protocol`.
- Issue: The candidate consumes a ticket before `Observe` and deliberately
  consumes it on failed network observations. `Begin` permits only one live
  session and blocks another begin until completion or the ten-minute expiry.
  A failed observation therefore leaves an incomplete session that the owning
  administrator cannot finish and cannot immediately replace. Waiting for the
  TTL is an avoidable availability failure, and a different administrator must
  not be able to clear the actor-bound session.
- Required refinement: define an authenticated, same-actor abort operation (or
  an equivalent replacement transaction) that atomically invalidates the
  session and all tickets, records a sanitized failure state, and permits a
  fresh `Begin`. It must accept the actor from the existing admin auth context,
  never from request JSON, and must be safe to race with ticket consumption and
  completion. Completion, expiry, configuration drift, and process restart
  should use the same invalidation primitive.
- Why not fixed: this is a public control-plane lifecycle decision and requires
  the host handler/session API to be finalized before implementation.

### 3. Make unprovisioned executor construction executable

- File: `research/host-runtime-wiring-plan.md`, `Wire and dependency ownership`.
- Issue: `NewRunManager` rejects a nil `ModelExecutor` and retains the executor
  for worker execution (`internal/codexeval/runs.go`, lines 100-104). The plan
  mentions either a lazy executor factory or an unavailable executor, but an
  unavailable placeholder has no activation mechanism while the manager is
  alive. Rebuilding the manager after initialization would complicate active
  state, shutdown ownership, and handler references.
- Required refinement: construct one nonnil lazy executor proxy at startup.
  Its calls must resolve the provisioned private token and verified literal
  loopback listener origin at execution time, create/close the real
  `HTTPExecutor`, and return a typed not-ready error when bootstrap/evidence is
  not valid. Keep the raw token inside the runtime service. A startup manager
  must remain usable for sanitized status and ordinary host startup without
  claiming evaluation readiness.
- Why not fixed: the proxy's ownership and token refresh/close semantics belong
  to the host runtime service and Wire/lifecycle design, not this plan review.

### 4. Keep the metadata slice Wire-safe and update the stale gate text

- File: `research/host-metadata-adapter-plan.md`.
- Issue: The proposed adapter correctly limits the next chunk to target
  resolution and credential-free DTO projection, but its constructor is stated
  to take three private local interfaces while the later Wire graph is built in
  `cmd/server`. A `wire.Bind` in that external package cannot name those
  unexported interface types. The same document also says Stage 4B quality
  review is pending, while `research/eval-run-state-quality-review.md` records
  a PASS.
- Required refinement: keep private reader interfaces and a private test
  constructor, but expose a production constructor/provider that accepts the
  existing exported dependency types (`service.AdminService`,
  `service.GroupRepository`, and `*service.APIKeyService`), or add an exported
  service-owned provider wrapper that performs the private-interface assembly.
  Do not export interfaces merely for Wire. Update the adapter gate to cite the
  Stage 4B quality PASS. The implementation remains safe only when limited to
  the two adapter files and patch 004; do not include handlers, probes,
  bootstrap, Wire graph changes, lifecycle, UI, or installer changes in that
  chunk.
- Why not fixed: this is a design/documentation adjustment for the next
  implementation slice, not a mechanical change to current product code.

## Accepted plan evidence

The candidate's actor-bound begin/complete ownership, one-shot external runner,
actual-listener loopback probe, protected inventory/session evidence, strict
TLS/redirect/proxy destination rules, ticket secrecy, normal gateway policy,
and manager-before-HTTP shutdown ordering are source-backed and retained. The
shutdown section correctly accounts for the current deferred cleanup path and
requires retrying manager drain while the local gateway and its dependencies
remain available.

## Verification

- Source review: pass for the named plan, PRD/design, Stage 4B manager/store,
  Gin 1.9.1 dispatch, gateway route registration, admin identity, and metadata
  adapter boundaries.
- Lint: not run; no product code was changed.
- TypeCheck: not run; no product code was changed.
- Tests: not run; this is a read-only plan review and no live Sub2API target is
  configured.
