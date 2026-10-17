# Administrator Account Evaluation Design

Status: planning. S01—S05 approved by user: 无误 继续. Remaining core/history/delivery review is pending.
Approved core: C02/C03/C04/C06/C07 in research/core-admin-channel-approval.md. C08—C12 approved by user: 无误 继续. C13—C17 approved by user: 是. Next review: research/core-fourth-approval.md (C18—C22, pending).
Prior design and source anchors are preserved in research/before-admin-channel-design.md
and research/planning-before-minimal-replan-design.md.

## Architecture

Use the existing administrator account-test infrastructure, not ordinary user gateway ingress.
Add POST /api/v1/admin/accounts/:id/codex-eval under native adminAuth beside /test.
Add a dedicated AccountHandler method and AccountTestService method based on the
existing OpenAI account-test sending code. Reuse existing dependencies; do not
refactor /test or introduce another identity/service framework.

The endpoint performs a single model request. Four original evaluations retain
approved start/state/cancel/records/clear-history management operations and algorithms.
The proposed task executor and new HTTP endpoint call the same sending method;
background tasks do not retain administrator JWTs or loop through local HTTP routing.
S05 approves this integration; task-wide implementation approval is still required.

## Sending and protocol proposal

Approved S01 specifies complete streaming Responses JSON as the request body and necessary
Codex protocol headers. Account ID comes solely from the path/accepted task.
No prompt wrapper, fixed hi, second account selector or internal target header.
The existing /test retains its connectivity payload and simplified output.
Non-streaming/media/compact scope is not silently included in this proposal.

Approved S02 reuses native selected-account credential/URL/model mapping, proxy, TLS,
account header overrides and doOpenAIAccountTestUpstream/plugin transport.
Upstream credentials stay server-side; incoming administrator Authorization is not forwarded.
No account pool, sticky binding, API-key authentication or user-billing row is involved.

Keep OpenAI OAuth/API-key target scope. For Chat Completions accounts, use the
existing apicompat request/event conversion functions; do not reduce full Codex
input to one prompt or rewrite native conversion rules. Full fixture equivalence
is unverified until implementation and checks of the approved scope.

Approved S03 specifies complete Responses SSE and actual upstream HTTP errors, with the
native sensitive-body redaction where applicable. Preserve completion/usage data
without simplification into test_complete. Supplier usage is evaluation data, not user billing.
Core error/SSE parser details remain individually reviewable.

Native anchors: routes/admin.go:29—36,:393; admin/account_handler.go:1265—1317;
account_test_service.go:800—1000,:2118—2186,:2947—3015;
openai_plugin_transport.go:23—48; apicompat/chatcompletions_responses_bridge.go:
:39,:1534,:1637,:1687,:1800. Full evidence is in research/native-connectivity-dispatch.md.

## State and cancellation proposal

Approved S04 retains retaining native account header snapshots and 401/429 handling,
but omitting the original handler's success-triggered broad state recovery
from the new full-response endpoint. This choice is approved under S04; /test remains native.
Native recovery anchor: admin/account_handler.go:1313—1317 and ratelimit_service.go:2202.

Approved S05 specifies forwarding the admitted task context, closing its response on cancel
and stopping further evaluation calls. Reuse original sampling limits; do not
claim native gateway Redis slots/WaitPlan or user-billing drains are inherited.
HTTPUpstream's connection management is not business concurrency enforcement.
Do not reintroduce io.Pipe, producer registries, fake Gin routers or custom shutdown.

## UI, history and delivery

Preserve approved F01—F06 without native Vue edits, original result/Pelican rendering,
account/model operations and task controls. Group selection remains a page filter,
not a user-billing identity. Remaining core/history changes require individual approval.
No automatic history migration/resume/deletion or database cleanup is introduced.

Build pinned v0.2.16 with native embedded frontend into a static Linux/amd64 binary.
Use an ordinary fixed patch; preserve upstream locks/versions. Development temp work
uses standard system directories. Final paths: dist/build/sub2api-codex-eval and
 dist/build/sub2api-codex-eval.tar.gz. Include install.py, standalone cleanup.py,
metadata/licenses; installer never fetches source/builds on the user's host.

## Retired integration

Remove evaluation-specific user/owner/Key/accounting/expiry/renewal logic,
selection headers, gateway authentication/selector/sticky/qualification/diagnostic/
partial-billing hooks, private capabilities, maintenance/activation and lifecycle frameworks.
Restore native gateway/main/shutdown code. No actual user data or runtime resources
are removed during planning, and code removal requires final implementation approval.

## Validation and approval

Reuse completed pristine evidence with original failures and skips.
Use local supplier fixtures and disposable Docker resources for approved behavior.
Distinguish compilation, installer operation and actual feature verification.
Race and paid supplier requests are excluded from acceptance.
S01—S05 and C02/C03/C04/C06/C07 are approved. Remaining core/history/delivery must converge before final task-wide review.
