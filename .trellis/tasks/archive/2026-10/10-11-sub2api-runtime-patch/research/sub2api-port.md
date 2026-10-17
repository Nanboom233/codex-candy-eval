# Researcher Evidence: Sub2API Evaluation Port

## Provenance

This report preserves findings returned by `native_plugin_research` during the preceding planning conversation on 2026-10-11.

The researcher resumed after account-concurrency failures. The report is a transferred research record, not a new implementation verification.

Target: official [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api), researched release v0.2.15.

Line references below belong to the source inspected during that conversation. Recheck them against the exact resolved commit before patching.

Reported short SHAs were not reconciled. Resolve annotated tags and actual commits explicitly instead of copying an uncertain SHA.

## Native Plugins

`backend/pkg/pluginapi/v1/plugin.proto` exposes transport forwarding and host services for KV, account metadata, and outbound identity.

References: lines 21–39, 100–113, and 206–274.

`ForwardRequestStart.account_id` identifies the account already selected by the host. The protocol does not expose request-scoped account selection.

`PluginManager.ShouldRouteOpenAIOAuth` uses the selected account and rollout bucket. `RoundTrip` forwards that account to the transport plugin.

References: `backend/internal/service/plugin_manager.go` around 947–989 and `plugin_runtime.go` around 274–301.

The initial manifest capability targets OpenAI OAuth outbound transport. It does not provide general evaluation management or account selection.

Reference: `backend/internal/service/plugin_manifest.go` around 147–175.

Conclusion: a transport plugin or sidecar alone cannot force a normal gateway request onto the requested account.

## Responses Routing and Accounting

`backend/internal/server/routes/gateway.go:179–189` applies API Key authentication before group model restrictions and required-group checks.

Lines 213–227 dispatch POST `/v1/responses` by group platform to the OpenAI or generic handler.

OpenAI path: `backend/internal/handler/openai_gateway_handler.go`:

- Lines 599–630 acquire the user slot and check billing eligibility and inflight balance reservation.
- Lines 667–689 call `SelectAccountWithSchedulerForCapability` without a forced-account argument.
- Lines 768–780 acquire the selected account slot.
- Lines 817–859 record usage with Key, user, account, subscription, and pricing context.
- Failure branches can return to the selection loop and fail over.

`backend/internal/service/openai_account_scheduler.go:74–95` defines the schedule request without a forced-account field.

Its public selector appears around 2100–2134. Session-based scheduling and fallback appear around 2298 onward.

Generic path: `backend/internal/handler/gateway_handler_responses.go`:

- Lines 135–168 enforce user concurrency and billing.
- Lines 183–224 select an account through load-aware scheduling.
- Lines 228–274 enforce account concurrency and profit admission.
- Lines 336–369 record normal usage.

Conclusion: implement account constraints inside these existing paths. Keep their downstream accounting and concurrency checks.

## Request-Scoped Group Selection

The API Key service structure stores `GroupID` and `Group` on the authenticated Key object.

Reference: `backend/internal/service/api_key.go:30–48`.

The authentication middleware checks group availability and binding permissions, loads the target subscription, then populates request context.

Reference: `backend/internal/server/middleware/api_key_auth.go` around 148–286 and 401–448.

Researcher recommendation:

1. Identify the evaluation Key after Key/IP/user authentication.
2. Apply the target group before common group availability, permissions, subscriptions, and billing checks.
3. Copy the request API Key and Group. Never mutate cached or persisted objects.
4. Keep context group, subscription, and pricing aligned with those request copies.
5. Preserve account eligibility and slots. Disable fallback and alternate-account selection for directed requests.

The group copy must be separate before changing fallback fields. A shallow API Key copy still shares its original Group pointer.

## Docker and UI Integration

The official Dockerfile builds a `CGO_ENABLED=0` Linux binary at `/app/sub2api`, embedding the frontend.

Researcher confirmed that compiling the evaluation module into a version-matched patched binary is feasible.

The existing administrator custom menu supports iframe resources. Native plugin management primarily controls `.s2plugin` transport packages.

Relevant frontend files:

- `frontend/src/components/layout/AppSidebar.vue`.
- `frontend/src/views/user/CustomPageView.vue`.
- `frontend/src/api/admin/plugins.ts`.

Runtime binary replacement leaves the original image and Compose files unchanged. Container recreation discards writable-layer replacement.

Stage the candidate and rename atomically before restarting. Preserve the original binary and verify rollback before live installation.

## Current Repository Integration Boundaries

The planning investigation identified these local source boundaries:

- `internal/plugin/plugin.go`: CPA registration, routes, and embedded UI.
- `internal/plugin/credentials.go`: CPA account discovery and selection.
- `internal/plugin/model.go`: CPA streaming execution, reads, closure, and response parsing.
- `internal/plugin/state.go`: atomic JSON snapshots and per-account retention.
- `internal/plugin/data/` and algorithm files: reusable evaluation context, probes, and baselines.

These local anchors came from the main planning investigation. Reinspect them during implementation instead of treating this report as a code audit.

## Source-Audit Follow-Up

The [pinned source audit](source-audit.md) verifies the release commit, GHCR identities, and current source behavior.

The [code-level contracts](patch-contracts.md) replace these broad recommendations with owner files, service boundaries, and regression assertions.

They identify missing private-token storage, legacy scheduling, detached cancellation, iframe URL authentication, SPA bypass, and durable run-state requirements.

Keep this report as transferred evidence. Use the newer audit and contracts for implementation review.

## Limits and Deferred Verification

- No patched binary was built or installed during research.
- No Sub2API database, user balance, subscription, or deployed configuration was inspected.
- No live evaluation verified usage or concurrency.
- Exact source commit, image digest, compatibility checks, and rollback behavior remain implementation verification items.
