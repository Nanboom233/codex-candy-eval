# Minimal Native Integration Design

Latest steering: pause C01—C05 and investigate extending native account connectivity tests.
See research/native-connectivity-dispatch.md. The Key/local-gateway requirements below
are the prior planning branch under review, not authorization to proceed. The administrator-test
alternative changes user billing and scheduler semantics and awaits user choice.

Status: planning. This design describes proposed behavior, not implemented behavior.
Current detailed plan: research/minimal-replan.md.
Historical design and every earlier source anchor remain in research/planning-before-minimal-replan-design.md.

## Architecture

Keep one Trellis task because every workstream contributes to one installable package.
Sequence approval, native integration, evaluation UI/core, build, then packaging.
Do not create a separate coordination framework or independently deployable service.

The page uses native administrator authentication and approved management APIs.
The evaluation manager uses the shared internal Key and its native accounting owner.
The executor calls the native local HTTP listener with a special key and target header.
Native middleware supplies the standard APIKey/User/Group context.
Native scheduling returns only the target. Native forwarding owns the remaining request lifecycle.

## Key and accounting

Use one shared internal Key and one native accounting row.
Native administrator authentication controls evaluation management operations independently.
Usage belongs to the shared Key's native owner. Remove initiating-user Key mappings.
Do not retain a first-activator-only permission gate.

Before each evaluation HTTP call, the trusted evaluation module resets the shared Key's native ExpiresAt to now + 5 minutes. Keep native expiry checks.
Native standard-mode Responses authentication already calls IsExpired().
Simple mode returns before that check; the internal-Key validation must call the same
native expiry predicate before any early return. Ordinary-key behavior stays native.
Do not create an already-expired row and then bypass expiry.
Gateway requests cannot extend expiry. Do not generate keys per request.
Reuse APIKeyService.Update with ExpiresAt; it updates the existing row, reactivates an
expired status and invalidates the native auth cache. It does not reset quota or enable
a disabled Key. Renewal failure prevents the evaluation HTTP call.
Remove one-time consumption, body/header bindings and task credential states.
After inactivity, the next authorized evaluation call renews the same expired Key.
There is no timer, heartbeat or separate renewal endpoint. Already admitted requests
retain native completion behavior even if their admission credential expires.
The shared native owner preparation must also be stated before implementation.

Internal authentication accepts only the shared Key on the fixed local Responses endpoint.
Normal keys cannot use the target header. Native row/user/group/ACL/balance/quota checks remain.
Strip internal credentials and the target header before upstream forwarding.
Do not add credential files, maintenance protocols, migrations or a separate authentication framework.

Latest source of scope: research/shared-key-expiry.md.
Earlier B06/B08/B09 per-user identity and expiry-exception behavior is superseded.

## Hard selection

Parse one X-Codex-Eval-Target header as positive group_id:account_id.
At the primary native selector, validate the selected group and OpenAI platform.
Keep existing quota/privacy/profit context preparation and channel pricing checks.

Reuse native selectBySessionHash with the explicit target ID and its existing request fields.
Preserve native qualification, DB recheck, concurrency acquisition and WaitPlan.
Set native DisableStickyEscape and PreserveStickyBinding for this direct selection call.
Return this branch directly. Do not call the full candidate-pool Select afterward.
No new qualification chain or post-slot database pass is added.

The selected target is subject to native ExcludedIDs.
Native same-account retries retain the target. Cross-account failure excludes it and leaves no candidate.
Keep the original Responses handler's retry/error/usage flow.
Keep native quarantine policy, profit admission, session records and protocol conversion.

Source anchors: openai_account_scheduler.go:72, :495, :2298.
Native retry exclusion is passed at openai_gateway_handler.go:676.
Native session signals are generated at openai_gateway_scheduling.go:164.
The current evaluation protocol supplies its original session/prompt-cache signals.

## Transport and cancellation

Use a standard HTTP client against the application's own listener and /v1/responses.
Build the destination from loaded native Server configuration, never client Host or Origin.
The current deployment listens on 0.0.0.0, so use 127.0.0.1 with its native port.
Do not use environment proxies or follow redirects for internal credentialed requests.
No user-defined gateway URL or compatibility resolver is added.

Remove evalServeHTTPTransport, io.Pipe, producer registry/barriers and capability generations.
Native HTTP server tracks the real connection and retains original shutdown behavior.
Task cancel closes its local request and prevents later probes/retries.
Native detachment may continue upstream reading for billing. This remains native behavior.
Do not claim immediate supplier cancellation or complete billing persistence.

## UI, core and history

Apply already approved frontend decisions without native Vue edits.
Preserve original result/Pelican rendering and evaluation algorithms.
Keep errors actually received by the evaluation layer under B05.
Do not add raw-error side channels or modify native HTTP/SSE/CC parsers.

Remaining core/history differences require individual decisions before implementation.
This plan does not authorize all current history validation, recovery or persistence machinery.
No legacy migration, automatic resume, user-history deletion or native database cleanup is introduced.

## Build and delivery

Use the pinned v0.2.16 source and one ordinary patch with necessary evaluation resources.
Replace structural multi-version adaptation and compatibility manifests.
Build the static Linux/amd64 server with real native embedded frontend in the development environment.
Reuse upstream declared versions and locks. Development temp work uses standard system directories.

Canonical outputs: dist/build/sub2api-codex-eval and dist/build/sub2api-codex-eval.tar.gz.
Install only the prebuilt binary with necessary scripts, metadata and licenses.
Remove maintenance/activation calls from installation together with their server removal.
Standalone cleanup.py targets identifiable old build resources and preserves business resources.

## Removed scopes

Remove native Vue menu/route modifications, activation/proof/ticket/intent/private-token frameworks and legacy migrations.
Remove alternate-selector, WS, custom sticky, extra freshness, native diagnostic and partial-billing hooks.
Restore native upstream contexts, HTTP parsers, ListenAndServe and shutdown.
Do not replace these deletions with another service, queue, framework or configuration layer.

## Validation and approval

Reuse completed pristine baseline evidence with its original failures and skips.
Validate approved behavior with local fixture upstreams and disposable Docker resources.
Distinguish actual feature checks from binary compilation and installer checks.
Race and paid production evaluation are not acceptance requirements.

This backend proposal is reviewable now. Core/history/delivery scope still needs individual approval.
The final task-wide planning review happens only after those decisions converge.
