# Native Responses Cancellation Review

## Scope

This review corrects an earlier reachability claim for the native Anthropic
Responses conversion. It does not expand the directed evaluation scope.

The approved directed target is an OpenAI group and an OpenAI account. Native
Anthropic conversion remains an ordinary OpenAI-gateway provider path unless
the selector contract changes in a later task.

## Reachability Evidence

The directed authentication boundary resolves the target group in
`backend/internal/service/eval_key_runtime.go:701-724`. It rejects every group
whose platform is not `PlatformOpenAI`.

The directed scheduler fixes `platform := PlatformOpenAI` and checks the target
account at `backend/internal/service/openai_eval_directed.go:39-63`.
`NormalizeOpenAICompatiblePlatform` preserves registered OpenAI-gateway
provider IDs such as Kimi, so a Kimi account does not satisfy this check.

`backend/internal/service/account.go:1431-1445` returns
`APIProtocolChatCompletions` for a `PlatformOpenAI` account because it is not a
multi-protocol API-key provider. A stored `api_protocol=anthropic` value cannot
make that account enter the native Anthropic branch.

The ordinary `Forward` dispatch at
`backend/internal/service/openai_gateway_forward.go:187-192` can still enter
`forwardResponsesViaNativeAnthropic` for a supported multi-protocol provider.
Direct service tests using a Kimi account therefore exercise the conversion
branch, but they bypass the trusted evaluation selector and do not prove a live
directed evaluation reaches it.

## Corrected Verdict

The withdrawn native conversion fixes were defensive controls for direct
service calls. They were not required by the current directed evaluation route.
The final source restores native conversion to the pinned upstream version.

They must not be reported as required reachable directed-evaluation coverage.
The real directed cancellation surface is the selected OpenAI account's HTTP,
passthrough, Responses conversion, and upstream WebSocket paths, plus handler
partial-usage submission and failover termination.

The native direct-service fixtures were removed with those changes.
They remain outside the current directed selector acceptance boundary.
No live Docker or billing acceptance follows from this source review.

Final specification verdict: **PASS** for the retained Stage 2C scope.
The final patch was replayed on pristine+001+002 and matched the workspace.
Patch 003 SHA-256:
`6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.
The follow-up adds a handler regression for failover-typed partial usage. Directed
service and handler race tests pass on the final patch.

## Uncovered Boundary

There is no handler-level native fixture through the real directed middleware,
because the current middleware and scheduler reject that account/provider
combination before `Forward`. A future provider-expansion task would need to
change the group/platform contract first, then add an end-to-end native test.
