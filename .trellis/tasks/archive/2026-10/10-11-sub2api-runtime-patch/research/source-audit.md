# Pinned Source Audit

## Scope and State

This audit continues the existing planning task. Product implementation remains subject to artifact review and activation.

The source checkout lives in ignored `dist/sub2api-v0.2.15/`. Do not commit the checkout or run the old CPA deployment.

## Verified Baseline

- Repository: https://github.com/Wei-Shaw/sub2api
- Release tag: `v0.2.15`.
- Annotated tag object: `86a80c13dcba86f52f9ca815b5a471cecc236227`.
- Peeled source commit: `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.
- Source commit date: `2026-10-09 11:19:17 +0800`.
- Local toolchain: `go1.27.2 darwin/arm64`.
- Docker CLI is unavailable. Container lifecycle tests cannot run in this environment yet.
- Default SSH clone failed host-key verification. An explicit public HTTPS clone succeeded without Git configuration changes.

Verify provenance with `git rev-parse v0.2.15` and `git rev-parse 'v0.2.15^{commit}'`. These identify different objects.

The pinned root Dockerfile builds with Go 1.27.2 and `CGO_ENABLED=0`. It requires `-tags embed` and a compiled frontend.

`Dockerfile:83-98` copies the frontend output before the backend build. A backend-only build would not preserve the management console.

`Dockerfile:134-161` establishes UID/GID 1000, `/app/sub2api`, `/app/data`, and the entrypoint. Installer checks must account for these facts.

## Local Baseline Checks

- `go test -race ./...`: passed. Plugin tests completed in 15.345 seconds.
- `go vet ./...`: passed.
- `node scripts/test-ui.cjs`: passed.

These checks cover the existing CPA implementation. They do not verify a Sub2API patch.

## Initial Source Findings

`backend/internal/server/middleware/api_key_auth.go:148-165` checks the user before group availability and binding permission.

Apply request-scoped group selection after the active-user check and before both group checks. Preserve normal authentication and billing after that point.

`backend/internal/server/routes/gateway.go:179-220` applies authentication, model allowlists, composite-target resolution, and required-group middleware before Responses dispatch.

Restrict directed evaluation requests to the exact supported Responses endpoint and OpenAI groups. Do not allow composite-target resolution to replace their selected group.

`internal/plugin/model.go:90-162` owns CPA stream execution. Replace its transport rather than sending direct upstream requests.

`internal/plugin/model.go:186-249` requires a final Responses event. Preserve this contract in the HTTP SSE adapter.

`internal/plugin/state.go:13-66` persists only final histories. Active runs are not persisted, so restart recovery needs more than the existing snapshot format.

`internal/plugin/runs.go:24-57` excludes simultaneous evaluation types per credential. Preserve this exclusion per Sub2API account.

## Confirmed Design Gaps

1. Advanced and legacy OpenAI schedulers both select accounts. Patch the common service boundary before either affinity path.
2. `dto/mappers.go:80-111` returns `APIKey.Key`, including in nested usage DTOs. File permissions alone cannot keep a registered Key private.
3. `CustomPageView.vue:117-121` uses a normal iframe navigation. It cannot send an administrator Authorization header.
4. `api_key_auth.go:176-188` bypasses billing in simple mode. Do not claim charged evaluations work in that deployment mode.
5. `openai_gateway_handler.go:668-670` documents detached stream draining. Canceling the worker does not imply immediate upstream cancellation.
6. `state.go:13-18` stores only final records. Persist running state before acknowledging a run, then reconcile interrupted runs at startup.

`openai_account_scheduler.go:1775-1831` already owns model, privacy, runtime, quota, shadow-parent, channel, capability, and profit eligibility.

Reuse these checks for a forced target. Do not replace them with `GetByID` followed by direct forwarding.

`openai_gateway_handler.go:768-780` performs account-slot and profit admission. `817-859` records usage, including partial-result accounting after stream failure.

Retain both stages. Guard profit-veto retries, unsupported continuation, same-account retry, and account-switch branches for directed requests.

The generic Responses path uses other provider adapters. For this OpenAI-only scope, reject directed context in that path instead of adding unsupported-provider routing.

## Registry and Frontend Corrections

The official GitHub release API reports publication at `2026-10-09T03:35:49Z`. Its release assets include both Linux architectures and checksums.

The GHCR `0.2.15` OCI index has digest `sha256:da0fd94daebf890419e21e956821fd6c2008fc666f9617f2eb917f4492245da3`.

A direct registry response and a local SHA-256 calculation agree. Its Linux image manifests are:

- amd64: `sha256:71d627c5f0d383be3f9098957ac20e0adc2c65ffbc758d7a7b1f7c0084bc5745`.
- arm64: `sha256:8b81de4ebe8fae8773e88b4f452851c75dd206ada98c6559a9efa8322397ec21`.

Generated web search returned different architecture-tag digests. Do not equate architecture-tag indexes with manifests inside the combined index.

Search trace: planning `954b59569d92`, search `b98cfccee84d`. Primary release and registry responses supersede unverified generated claims.

`frontend/src/utils/embedded-url.ts:25-31` appends the administrator JWT as a `token` URL parameter. `CustomPageView.vue:259` requires an absolute HTTP URL.

The previous no-frontend-change assumption is incomplete. Add a narrow source exception for the same-origin evaluation resource, not a compiled-bundle rewrite.

For that resource, omit JWT and source-page query parameters. Read the current same-origin `auth_token` only for authenticated API headers.

Use a non-sensitive static page shell, then require normal admin middleware for all data and operations. Do not accept JWT query authentication.

`security_headers.go:130-148` applies `X-Frame-Options: DENY` and CSP nonces. The resource needs the existing nonce and a route-local `SAMEORIGIN` embedding response.

Do not weaken the global security policy. Preserve `frame-ancestors 'self'` and the opaque-origin pelican preview.

## Additional Integration Boundaries

`api_key_service.go:492-597` creates a normal Key with binding permissions and creation limits. Reuse this service directly on the server.

`api_key_service.go:739-790` performs cached Key authentication. Match the configured owned Key by immutable identity, never by its display name.

`api_key_auth.go:426-434` defers subscription-group access to subscription checks. `Create` additionally calls `canUserBindGroup` for normal binding permission.

Use both applicable rules for directed target groups. Do not infer subscription access from administrator role.

`gateway_usage_billing.go:573-587` uses `context.WithoutCancel` for upstream contexts. Preserve cancellation only for trusted evaluation context.

Keep `detachedBillingContext` unchanged so already observed usage can still reach normal accounting after cancellation.

`backend/internal/web/embed_on.go:86-106` intercepts unrecognized non-API paths and serves the SPA. `/codex-eval/ui` needs an exact frontend-bypass entry.

Alternatively, use an `/api/` resource path, but that changes the existing planned interface. The proposed design retains `/codex-eval/ui`.

`cmd/server/main.go:185-190` shuts down HTTP before deferred cleanup. Stop evaluation workers before HTTP shutdown because they call the local gateway.

Existing handler failover fixtures run in simple mode. They cannot verify charged evaluation behavior without normal-mode billing and concurrency fixtures.

## Verified Image Configs

Both registry configs identify Linux, the expected architecture, `/app/docker-entrypoint.sh`, `/app/sub2api`, version `0.2.15`, and the pinned source revision.

- amd64 config/image ID: `sha256:6c51fac82f6348083b2dad597177185e64ebf7c734592116edcb017014c2edc6`.
- arm64 config/image ID: `sha256:b7cf64a3ed962aa518565c09197ad16ca64df4a20ef1edca76b1d3c79d90c539`.

Keep these config identities separate from image-index and platform-manifest digests. Installation must inspect the real target before claiming compatibility.

## Pinned-Host Baseline Tests

Run from `dist/sub2api-v0.2.15/backend/`:

```sh
go test -tags=unit -race ./internal/server/middleware ./internal/server/routes ./internal/handler ./internal/service -run 'TestAPIKeyAuth|TestAdminAuth|TestOpenAISelectAccountWithLoadAwareness|TestOpenAIGatewayHandlerResponses_Failover' -count=1
```

Middleware, handler, and service checks passed. The routes package reported no matching tests in this command, so it was not counted as tested.

Run the separate route and DTO checks:

```sh
go test -tags=unit -race ./internal/server/routes ./internal/handler/dto -run 'TestGatewayRoutesGroupModelAllowlist|TestAPIKeyFromService' -count=1 -v
```

All six route-allowlist tests and both selected DTO tests passed. Both packages completed with exit status zero.

Additional unmodified-host checks passed after agent review:

```sh
go test -tags=unit -race ./internal/service ./internal/handler -run 'TestForwardOpenAIWSV2_ClientCancellationDrainsWithoutSyntheticFailure|TestAcquireResponsesAccountSlotProfitRecheck|TestOpenAIProxyStreamCircuit|TestOpenAIProxyStreamQuarantineBypassContext' -count=1
```

This verifies the existing WebSocket cancellation-drain, profit-slot recheck, and proxy-circuit baselines. It does not prove new evaluation-specific cancellation or isolation.

These are unmodified-host regression baselines. No new forced-routing, token privacy, billed-concurrency, or installer test has passed yet.

## Concrete Implementation Proposal

The [patch contracts](patch-contracts.md) define exact source owners, proposed service signatures, negative tests, and installer transactions.

A source scan confirmed the administrator DTO returns database Key values. Container recreation therefore defeats plaintext-Key privacy even when the patch adds DTO redaction.

The recommended proposal separates the private random token from its deterministic registered lookup value. This preserves the normal Key row without storing the token there.

`api_key_auth_cache_impl.go:269-331` keys and materializes authentication from its lookup argument. Keep that argument consistent for cache invalidation and Key updates.

`openai_gateway_usage.go:202-203` resolves pricing from request Key group fields. Preserve the selected-group copy through accounting instead of changing the persisted binding.

`internal/plugin/codex.go:189` creates an independent background context. Porting cancellation requires passing the run context into `askCodex`, not only replacing transport.

The user confirmed the conditional Key boundary K1-K7. Product implementation and the final planning summary still require separate approval.

The UI source exception and Key scheme are not applied patches or verified production behavior.

## UI Continuation Constraint and Test Gap

The user requires Sub2API visual adaptation without changing the main evaluation structure. This is recorded as requirement R13 and acceptance criterion A15.

Two `trellis-research` agents independently examined UI adaptation and gateway contracts. Their reports are scoped to separate task research files.

The UI report verified native teal/slate tokens and theme-change iframe reload. The main session checked `tailwind.config.js` and `CustomPageView.vue` anchors.

The gateway report required a corrective pass for inaccurate anchors and overbroad authentication/cancellation claims. Its expired session handle was replaced using the saved report.

Main-session source reads confirmed `openai_gateway_forward.go:887` can use upstream WebSocket from the ordinary HTTP Responses path.

POST-only ingress therefore does not exclude WebSocket cancellation work. The corrected contract covers `openai_ws_forwarder_v2.go:402` and interrupted lease disposal.

The stored lookup value remains a normal Key authentication value on the unpatched host. Raw-token privacy is not proof that all derivatives become unusable.

The [private-Key lifecycle report](private-key-lifecycle.md) compares concrete alternatives. `ent/schema/api_key.go:37-40` and native authentication both limit Key values to 128 bytes.

Google ingress accepts query `key` values, so non-header control markers are not a general solution. No storage or input-policy exception has been approved.

A loopback-scoped normal row is only conditional defense. Native source-IP checks depend on forwarded-header trust and actual reverse-proxy topology.

`config/config.go:2115` enables legacy forwarded-IP trust by default. `pkg/ip/ip.go:164-171` shows its effect on security-sensitive source resolution.

Do not change that policy automatically or claim remote isolation without deployment evidence. The original browser can still display the stored lookup value.

Keep token mapping/replay rejection at ingress boundaries. Internal `GetByKey` must remain usable for cache invalidation, recovery, and billing refresh.

An oversized database lookup marker cannot solve this within the existing schema. Do not widen native limits or alter security policy.

`RequireGroupAssignment` rejects unbound inference only when ungrouped scheduling is disabled. Do not conflate Key authentication, gateway inference, and directed capability.

`README` or product-code changes remain outside this planning continuation. The agents may not activate the task or edit implementation files.

`scripts/test-ui.cjs:46-47` excludes browser initialization before running mocked JavaScript. Its environment does not exercise actual theme listeners or parent-window synchronization.

The current UI check therefore does not prove Sub2API theme adaptation. Add focused theme-environment tests and browser checks after implementation approval.

Keep the opaque-origin pelican assertion at `scripts/test-ui.cjs:188`. Theme changes must not grant model-generated HTML same-origin privileges.

## Investigation Errors and Recovery

- The initial public clone used SSH and failed host-key verification. Explicit HTTPS succeeded without configuration changes.
- A nested Python `exec` registry query triggered the local destructive-operation hook. A plain read-only urllib query succeeded without policy changes.
- A guessed `handler/handlers.go` path did not exist. Source search found the owner at `handler/handler.go`.
- Incorrect root-level artifact paths (`design.md` and `task.json`) were rejected by the edit tool. No root files changed.

## Verification Still Required

- Review and approve the revised private lookup, scheduler, UI, and lifecycle contracts.
- Write failing regression tests, apply source patches, and prove red-to-green behavior.
- Run full patched-host unit/race checks, frontend checks, and both static Linux builds.
- Run PostgreSQL, Redis, Docker lifecycle, and one paid evaluation against an approved Sub2API target.

Docker CLI and target configuration are unavailable in this environment. Their acceptance checks remain incomplete, not passed.
