# Research: Response Model Validation in the Native Patch

- Query: Determine whether the current native Sub2API patch added or changed a check that rejects a response when the upstream response model ID differs from the requested model ID.
- Scope: internal
- Date: 2026-10-16

## Findings

### Patch scope

The maintained patch does not change response handling, SSE parsing, model billing, or gateway forwarding. `native/sub2api.patch:1-84` only adds the `/codex-eval` page and administrator routes, the native sidebar/router entry, the embedded-page bypass, and relative embedded URL resolution. There is no added `upstreamModelMismatch`, response-model observer, model comparison, or response rejection code in this patch.

### Native evaluation sender

The new sender checks only the incoming evaluation request. `native/backend/internal/service/account_codex_eval_service.go:56-68` requires a non-empty `model` and `stream:true`. It writes the selected model back to the outgoing payload at `:115-116`; this preserves the literal model selected by the administrator. It sends through the existing account-test transport at `:187-188`, then returns the upstream response. The success path does not compare a response model with the request model.

For API-key accounts, the sender chooses Responses or Chat Completions from the account capability at `:100-113`. The Chat Completions request is converted from the Responses request at `:120-134`. After the upstream response returns, the bridge is created with the requested model at `:214-225`. The bridge validates stream lifecycle, tool arguments, and upstream errors at `:249-284`; it has no model equality check.

### Responses SSE parsing

The evaluation parser in `native/backend/internal/codexeval/model.go:195-276` reads response events, status, output items, errors, and usage. The `responseBody` type has no model field (`:156-167`), and `readResponseStream` never compares a response model to the requested model. A response with a different model ID therefore continues when status, text output, and stream completion are valid.

### Chat Completions conversion

The baseline bridge state is initialized with the requested model in `backend/internal/pkg/apicompat/chatcompletions_responses_bridge.go:1636-1649` as used by the new sender. The converter updates `state.Model` from a Chat Completions chunk only when the state model is empty (`:1694-1699`). Because the evaluation path supplies the requested model, an upstream chunk with another model ID is not rejected and is not adopted into the converted Responses events. This is output normalization, not validation.

### Existing account-test path

The pristine v0.2.16 account-test stream handler at `backend/internal/service/account_test_service.go:2946-3014` handles text deltas, completion, failure, and error events. It does not compare model IDs. The shared transport at `backend/internal/service/openai_plugin_transport.go:21-44` only selects the plugin or HTTP transport and returns `http.Response`; it does not inspect response model IDs.

### Existing ordinary gateway audit

The pristine v0.2.16 source already has a separate response-model observer. Its comment states that observation never affects forwarding (`backend/internal/service/upstream_response_model.go:15-22`). The observer records model declarations from OpenAI response payloads at `:71-74`, and `upstreamModelMismatch` returns a nullable audit flag at `:274-294`. Normal usage logging stores this flag at `backend/internal/service/gateway_usage_billing.go:1285-1295` and `backend/internal/service/openai_gateway_usage.go:384-394`. These paths record mismatch information for ordinary gateway usage; they do not reject a successful response. The native evaluation sender does not create a Gin gateway observer or call these usage-recording functions.

### Related specifications

The task design requires the selected model to pass literally and bypass evaluation model configuration (`.trellis/tasks/10-11-sub2api-runtime-patch/design.md:13`, `:39`). The PRD requires the new endpoint to use the native account-test sender and preserve the selected model literally (`prd.md:19-21`). The runtime specification identifies `SendCodexEvaluation` as the shared sender and requires the complete Responses stream contract (`.trellis/spec/backend/runtime-installation.md:13-15`, `:42-47`). The observed code follows these contracts.

## Conclusion

There is no evidence that this patch introduced an erroneous response-model mismatch check. The ordinary Sub2API gateway already records requested-versus-response model differences for usage auditing, but the patch does not modify that code and the new evaluation route does not enter that path. The evaluation sender sends the requested model literally and accepts a valid completed response even when the upstream response model ID differs.

The evaluation result does not expose the upstream model ID. Responses parsing ignores it. Chat conversion keeps the requested model in the generated Responses events. This matches the plugin scope. Sending the requested model literally and comparing a response model are separate operations.

## Caveats / Not Found

- The user confirmed that response-model comparison is outside the plugin contract. Its absence is not a defect or a required test. The existing sender test covers literal request-model forwarding and both protocol paths at `native/backend/internal/service/account_codex_eval_service_test.go:19-75`.
- No external references were needed. The conclusion is based on the pinned pristine v0.2.16 source and the current `native/sub2api.patch` and native evaluation files.
