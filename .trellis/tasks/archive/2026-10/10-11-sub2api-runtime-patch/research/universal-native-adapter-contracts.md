# Research: Universal Native Adapter Contracts

- Query: Compare the native gateway, auth, billing, cancellation, route, and wiring seams across actual Sub2API source revisions; verify how the generated AST adapter anchors them and identify semantic drift risks.
- Scope: internal
- Date: 2026-10-13

## Findings

### Revision comparison

The retained source worktrees are Sub2API v0.2.14 at `0363b8cdba8cec3e2ba4b2dbd49c4481143fa55d` and v0.2.15 at `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. They differ materially across roughly 278 files. The implementer reports that the 73-rule adapter compiles, exported-symbol checks pass, and focused/composed native tests pass against both revisions. This is evidence for these two revisions, not arbitrary forks or future releases.

| Seam | Stable contract across v0.2.14 and v0.2.15 | Adapter hook and semantic constraint |
| --- | --- | --- |
| Main API-key authentication | `apiKeyAuthWithSubscription` retains the public credential extraction/lookup path, active-user check, group checks, and native subscription/billing flow. | `adapter/rules.json:550-580` inserts `AuthenticateEvalCapability` before public parsing, keeps public parsing in the `!internal` branch, inserts exact target application after the active-user check, and exempts expiry only for the private capability (`:581-625`). Exact lookup/private-token rejection also belongs before Google public lookup (`:627-638`). These three ingress owners, including the WebSocket path that shares main auth, must stay covered. |
| Responses selection | `SelectAccountWithSchedulerForCapability` retains its request shape and is called inside the existing Responses loop; account selection still precedes the response handler's account-slot admission. | Rules `:680-730` route a directed context through `selectDirectedEvaluationAccount`; that path rechecks the named account, group membership, fresh schedulability, DB state, model/capability, and transport. `:407-456` adds a post-slot recheck before Forward and updates both local `account` and `selection.Account`. Any error must terminate this target rather than select a replacement. |
| Scheduler fallback | Normal selection still calls `selectAccountWithSchedulerOnce`, with the existing proxy-stream quarantine fail-open retry in its wrapper. | Rule `:680-730` returns immediately for directed contexts before the quarantine retry. Rules later in the recipe guard legacy previous-response, sticky, and load-awareness helpers as unsupported. Keep directed requests out of ordinary fallback and sticky-session mutation. |
| Retry and failover | The handler keeps its existing bounded same-account retry, failed-account set, and account-switch paths. | Rules `:420-483` stop unsupported continuation/account-switch failover for directed work while retaining retry on the selected account. Repeated target exclusion must fail closed. Partial usage is submitted before directed failover returns (`:459-483`). Do not move the target branch ahead of native account eligibility or slot admission. |
| Responses forwarding | `OpenAIGatewayService.Forward` is still the production dispatch point and receives the selected account and request context. | Preserve the whole current function's protocol resolution and forwarding branches. In v0.2.14 OpenCode Go and Anthropic protocol branches are explicit; v0.2.15 routes these through `upstreamRoutingModel` and `resolveUpstreamProtocolFor`, and also propagates `ResponsesLite` and recognizes a broader encrypted-content retry. Rules `:738-757` add partial-result/cancellation handling at a local streaming error branch; they do not replace the protocol dispatch. |
| Billing identity | `OpenAIGatewayService.RecordUsage` and the shared usage-billing machinery remain the accounting path. Usage rows and billing commands use `APIKey.ID`; quota/rate mutations and cache invalidation use the same native service. | Preserve the `APIKey`, `User`, `Account`, `Subscription`, pricing, usage result, and API-key service passed by native `RecordUsage`. The owned API-key row remains the accounting foreign key. Do not introduce direct charges, substitute a synthetic ID, or bypass `applyUsageBilling`. In v0.2.15, see `internal/service/openai_gateway_usage.go:154-177,380-382,506` and `internal/service/gateway_usage_billing.go:288-384,647-667,906`. |
| HTTP cancellation | Native streaming handlers detach ordinary stream contexts so they can drain usage after a client disconnect. | Rules `:654-678,700-737,759-765` route context creation through directed-aware helpers. Directed requests retain cancellation; ordinary requests still use native `WithoutCancel` behavior. Rules `:738-757,766-791` preserve partial Results/passthrough usage returned with stream errors for native billing. |
| WebSocket cancellation | Native WSv2 tracks client disconnect, drains usage, and releases or poisons the upstream lease according to its terminal/read state. | Rules `:840-928` give directed work a cancelable read context, cancel it on directed disconnect, and prevent the ordinary detached-context cancellation branch from continuing the drain. Confirm `MarkBroken` remains on failed/canceled reads and `Release` remains in the lease defer. Rule `:733-736` preserves a WS partial result/error for the handler rather than writing the normal fallback response. |
| Protocol and route dispatch | `Forward`'s signature and route registration remain stable, but source code changed: v0.2.14 has a platform switch; v0.2.15 uses `domain.UsesOpenAIGateway` and platform catalogs. | The directed path must enter the registered `/v1/responses` gateway handler and preserve the runtime-selected account. Do not inject a second static platform allowlist. Current/previous route evidence: `source-current/backend/internal/server/routes/gateway.go:49-60`; `source-previous/.../routes/gateway.go:48-67`. |
| Runtime/admin wiring | Wire, handler provider sets, admin routes, router construction, and frontend menu integration are the composition owners. | Rules `:394-548` and `:930-1005` bind the runtime/handler/router/admin path; the generated Wire graph is changed through exact statement anchors. `adapter/ui-rules.json` provides three frontend files of exact substring anchors, not Go AST matches. Compilation and browser/UI checks remain separate requirements. |

### Stable implementation anchors

- `backend/internal/handler/openai_gateway_handler.go`: `Responses` selects at `:676`, binds `account := selection.Account` at `:740`, acquires the account slot at `:768`, calls `Forward` at `:798`, and records usage at `:830-849`. The post-slot recipe insertion is immediately after the native `slotResult != openAISlotAcquireOK` guard. Current v0.2.15 also has WebSocket per-turn billing-group refresh at `:2124-2221` and uses the selected turn key in `RecordUsage` at `:3162-3168`; this is a revision-specific behavior the adapter must leave intact.
- `backend/internal/handler/openai_gateway_handler.go`: `acquireOpenAIAccountSlot` at `:2269-2403` handles already-acquired slots, quick acquisition, wait queue, the post-wait profit veto, and release wrapping. Each successful veto recheck writes `latest` into `selection.Account`. The directed post-slot fresh-account recheck must return the same target ID and must not lose/release the slot incorrectly.
- `backend/internal/service/openai_account_scheduler.go`: `SelectAccountWithSchedulerForCapability` at `:2116-2135`, the quarantine wrapper at `:2163-2195`, and `selectAccountWithSchedulerOnce` at `:2298-2325` are the stable service-level hooks. The directed branch is inserted after platform normalization and before general scheduler setup.
- `backend/internal/service/openai_gateway_forward.go`: `Forward` begins at `:21`; current protocol dispatch is at `:180-195`, WSv2 at `:786-904`, detached HTTP upstream context at `:1036-1046`, and partial streaming response handling at `:1225-1264`. v0.2.14/v0.2.15 protocol blocks differ; use the local loop/error anchors, not a full-function replacement.
- `backend/internal/service/openai_ws_forwarder_v2.go`: lease ownership is established at `:263-268`; ordinary disconnect detachment is at `:395-413`; reads use `ReadMessageWithContextTimeout` at `:542`; read failures mark the lease broken at `:577-600`. The directed additions preserve the existing cleanup owner while changing only the cancellation choice.
- `backend/internal/server/routes/gateway.go`: `RegisterGatewayRoutes` keeps route topology, while the compatible-platform predicate changed from explicit switch to domain function between revisions. Its platform set is upstream-owned and should not be duplicated by the adapter.

### Structural matching and fail-closed behavior

- `adapter/main.go:20-29,173-258` represents rules by file, receiver/name declaration key, AST node kind, optional AST field, canonical old/new syntax, scope, and adjacent sibling anchors. `apply` requires exactly one declaration and one match per rule (`:268-285`), parses/formats all staged files before writing them (`:301-321`), and exits on missing, duplicate, overlapping, or syntactically invalid anchors.
- The main auth rewrite replaces 14 exact sibling statements to branch private capability from ordinary credential parsing (`adapter/rules.json:550-560`). `APIKeyService.Create` replaces 15 exact statements with the expiry-aware helper (`:641-651`). Those are substantial anchors: drift should fail the adaptation, not be fuzzed or patched approximately.
- The rule engine is AST-shaped but does not resolve each selected hook with `go/types`. `adapter/typecheck.go:14-68` checks that several named methods exist, without asserting their exact signatures. The generated Go is then compiler-checked by the selected source build (`adapter/build.py:76-83`), which proves generated syntax and call types compile, but compilation alone does not prove billing, cancellation, or exact-target semantics.
- UI adaptation in `adapter/build.py:36-48` uses exact substring count checks. It fails on missing/ambiguous anchors, but it is not AST/type-aware and should be described as a separate constrained text adapter.
- The current adapter is source-driven with no release/image branch in `adapter/build.py:15-35,51-86`; module paths and build facts derive from the selected source/profile. The preserved compatibility evidence is specifically v0.2.14 and v0.2.15. Keep compatibility claims bounded to the successfully adapted revisions and future sources that satisfy all hooks/contracts.

### Critical semantic constraints for patch review

1. Keep exact target group/account identity through native `getSchedulableAccount`, membership and platform checks, fresh eligibility, `recheckSelectedOpenAIAccountFromDB`, post-slot recheck, and retries. Never let `excludedIDs`, previous-response lookup, quarantine retry, or failover choose a sibling account.
2. Keep account/user concurrency slots, profit-gate selection context, post-wait `ProfitControlVetoLatest`, and release behavior. Recheck after waiting before forwarding and retain `selection.Account = latest` when account data is refreshed.
3. Preserve native `RecordUsage` on both successful and partial/error-with-usage paths. Keep the owned API-key database ID on usage logs and billing commands; preserve quota/rate updates, billing repository application, and cache invalidation.
4. Directed request cancellation must cancel HTTP and WS upstream work. Ordinary streaming must keep the source's detached usage-drain behavior. For WS, canceled reads must terminate the directed operation and poison an uncertain upstream lease while still surfacing any already-observed usage.
5. Keep current protocol routing, retries, response limits, and handlers intact. Insert narrowly at the observed calls/branches; do not replace the current `Forward` function or route catalogs.

## Files Found

- `.trellis/workflow.md`: Trellis planning, research persistence, and execution workflow.
- `.trellis/spec/backend/universal-linux-delivery.md`: approved source-driven build and adapter constraints.
- `.trellis/spec/backend/patch-native-installation.md`: private capability, exact target, native accounting, cleanup, and cancellation contract.
- `.trellis/spec/backend/host-target-adapter.md`: exact group/account target resolution and scalar metadata contract.
- `adapter/rules.json`: 73 Go AST rules, imports, and injected declarations for native integration.
- `adapter/main.go`: canonical AST matching, exact-one checks, staging, and application.
- `adapter/typecheck.go`: exported compiler-type lookup and method-presence checks.
- `adapter/build.py`: source-module rewrite, adapter application, resource copy, source facts, type graph check, and build.
- `adapter/ui-rules.json`: exact text-anchor updates for administrator UI integration.
- `adapter/runtime/backend/internal/service/openai_eval_directed.go.in`: exact directed account re-resolution and scheduler-selection implementation.
- `dist/sub2api-build/universal-20261013/source-previous/backend/internal/{handler/openai_gateway_handler.go,service/openai_gateway_forward.go,server/routes/gateway.go}`: v0.2.14 comparison sources.
- `dist/sub2api-build/universal-20261013/source-current/backend/internal/{handler/openai_gateway_handler.go,service/openai_gateway_forward.go,service/openai_account_scheduler.go,service/openai_ws_forwarder_v2.go,service/openai_gateway_usage.go,service/gateway_usage_billing.go,server/routes/gateway.go}`: v0.2.15 comparison sources.

## Related Specs

- `.trellis/spec/backend/universal-linux-delivery.md`
- `.trellis/spec/backend/patch-native-installation.md`
- `.trellis/spec/backend/host-target-adapter.md`
- `.trellis/spec/backend/runtime-installation.md`

## External References

None used. Findings come from the two retained upstream source revisions and the current adapter implementation.

## Caveats / Not Found

- Successful compilation on two adjacent revisions does not establish compatibility with arbitrary forks or later upstream changes. Every required hook and semantic acceptance test must remain part of pre-replacement admission.
- The AST rule matcher uses canonical syntax/scope and exact sibling anchors. It is not a general type-aware program analysis. The compiler build checks generated types, while separate composed tests must continue to prove native behavior.
- Frontend rules use text anchors and need their own missing/duplicate checks and browser verification.
- No source-level claim here substitutes for real billing-provider behavior, a production container lifecycle, or deployed account evaluation.
