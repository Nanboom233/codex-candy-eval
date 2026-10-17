# Research: Response Model Validation Follow-up

- Query: Recheck whether the Sub2API patch causes a response-model validation failure when an upstream response model ID differs from the requested model ID.
- Scope: internal
- Date: 2026-10-16

## Findings

The build input boundary is narrow. `scripts/build-sub2api.py:18-24` copies the original source, applies `native/sub2api.patch`, and overlays `native/backend` and `native/frontend`. The patch hunks at `native/sub2api.patch:1-84` add the evaluation routes, page bypass, sidebar/router entry, and relative embedded URL handling. They do not touch gateway response handling, SSE model extraction, usage billing, or model mismatch auditing. The native backend overlay adds new evaluation files; it does not replace any existing response-model implementation.

The maintained comparison of the pristine v0.2.16 tree and the materialized patched tree reported identical SHA1 values for the existing response-related files:

| File | SHA1 in both trees |
| --- | --- |
| `backend/internal/service/upstream_response_model.go` | `548b69976860a8ed5474c54b735f33b87e0e93d9` |
| `backend/internal/service/openai_gateway_usage.go` | `da64f4368e20424a108b424b167f849172787e08` |
| `backend/internal/service/gateway_usage_billing.go` | `21116e0e6fb13166d8a71a206cb61bf97b553f51` |
| `backend/internal/service/openai_gateway_response_handling.go` | `55ebb62337c4a9cfd53b90a889fbde43558076fc` |
| `backend/internal/service/account_test_service.go` | `29d15453710fb26c7fe12291a0464fbe3c51d149` |

The evaluation sender validates the incoming request only: `native/backend/internal/service/account_codex_eval_service.go:56-68` requires a nonempty model and `stream:true`. It writes that selected model back into the outgoing payload at `:115-116`, sends through the native account-test transport at `:187-188`, and returns the upstream response. There is no comparison between the selected model and a response model. For Chat Completions accounts, `:214-225` initializes the conversion state with the requested model and wraps the stream; this path validates stream structure and errors but does not reject a model ID mismatch.

The Responses parser also has no response-model equality check. `native/backend/internal/codexeval/model.go:156-167` defines status, error, output, and usage fields without a model field. `:195-276` parses completion events and usage, checks stream completion, and returns output or an error. A valid completed response with a different upstream model ID is therefore accepted by this parser.

The existing ordinary gateway observer is a separate audit feature present in pristine v0.2.16. `backend/internal/service/upstream_response_model.go:15-22` documents that observation does not affect forwarding. Its model declaration and conflict data are consumed by usage logging at `backend/internal/service/gateway_usage_billing.go:1285-1295` and `backend/internal/service/openai_gateway_usage.go:384-394`; those paths record mismatch metadata rather than reject a response. The direct evaluation sender does not enter the ordinary Gin gateway observer or usage-recording path.

The focused evaluation test verifies literal model forwarding for Responses and Chat Completions at `native/backend/internal/service/account_codex_eval_service_test.go:19-75`. The user confirmed that response-model comparison is outside the plugin contract. Its absence is not a defect or an additional test requirement.

## Files Found

- `scripts/build-sub2api.py`: materializes the original source, applies the small patch, and overlays native evaluation files.
- `native/sub2api.patch`: route, page, sidebar, router, and embedded URL changes only.
- `native/backend/internal/service/account_codex_eval_service.go`: direct selected-account evaluation sender.
- `native/backend/internal/codexeval/model.go`: Responses stream parser without a model field.
- `native/backend/internal/service/account_codex_eval_service_test.go`: literal request-model and protocol-path tests.
- `backend/internal/service/upstream_response_model.go`: pre-existing ordinary gateway response-model observer.
- `backend/internal/service/gateway_usage_billing.go`: pre-existing mismatch audit persistence.
- `backend/internal/service/openai_gateway_usage.go`: pre-existing mismatch audit persistence.
- `backend/internal/service/account_test_service.go`: pre-existing account-test stream handling.

## Code Patterns

- Evaluation input contract: `account_codex_eval_service.go:56-68`.
- Literal outgoing model: `account_codex_eval_service.go:115-116`.
- Native upstream dispatch: `account_codex_eval_service.go:187-188`.
- Chat conversion seeded with requested model: `account_codex_eval_service.go:214-225`.
- Responses stream fields and completion handling: `codexeval/model.go:156-167,195-276`.
- Existing ordinary observer is audit-only: `upstream_response_model.go:15-22,71-74,274-294`.

## External References

None. The conclusion uses the pinned pristine v0.2.16 source, the maintained patch, and the native overlay.

## Related Specs

- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md:13,39` requires literal selected-model forwarding and bypass of evaluation model mapping.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md:21,53` requires the native account-test sender and literal model input.
- `.trellis/spec/backend/runtime-installation.md:42-49,103-105` defines the complete Responses request and literal model behavior.

## Conclusion

The patch did not introduce an erroneous response-model mismatch rejection. Sub2API's ordinary gateway still records requested-versus-upstream model differences for audit purposes, but those files are byte-identical to pristine v0.2.16 and that observer does not reject forwarding. The new evaluation route sends the administrator-selected model literally and accepts a valid completed response even when the upstream response model ID differs.

## Caveats / Not Found

- No additional plugin response-model test is required. Preserve the existing native gateway observer and audit behavior.
- No external documentation was required.
