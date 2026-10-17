# Minimal Native Integration Design

Status: planning. This design describes proposed behavior, not implemented behavior.
Current detailed plan: research/minimal-replan.md.
Historical design and every earlier source anchor remain in research/planning-before-minimal-replan-design.md.

## Architecture

Keep one Trellis task because every workstream contributes to one installable package.
Sequence approval, native integration, evaluation UI/core, build, then packaging.
Do not create a separate coordination framework or independently deployable service.

The page uses native administrator authentication and approved management APIs.
The evaluation manager captures the actual initiating user and the prepared native accounting row.
The executor calls the native local HTTP listener with a special key and target header.
Native middleware supplies the standard APIKey/User/Group context.
Native scheduling returns only the target. Native forwarding owns the remaining request lifecycle.

## Key and accounting

Preserve B06's per-user native accounting row with fixed past expiry and reserved database lookup identifier.
Keep native creation limits and B09 user/group billing rules.
Prepare the row only after a valid start request. Generate a random key for each native HTTP invocation.
Bind UserID/APIKeyID, task context, target, body digest, allowed Codex headers and the existing request deadline.
Require the exact route, raw loopback socket peer and matching request before atomic single-use consumption.
Reuse the existing runtime lock only for map access and final consume. Do not hold it during body reads or DB calls.
Cancel/expiration rejects pending keys. The executor removes remaining entries on return.
Never persist keys, return them to the browser, or generate an extra token file.
Details and first-use theft limits are in research/internal-key-abuse.md.
The database identifier is not an authentication secret.

On incoming POST /v1/responses, recognize Authorization: Bearer <internal-key>.
Read the bound native row and user. Verify identity and the approved accounting-row constraints.
Copy the APIKey/Group for the selected group and enter the native middleware checks.
Skip only the fixed past-expiry check for that verified internal row.
Normal status, user, ACL, subscription, balance, quota and concurrency rules still apply.

Normal credentials use their native path. Reject a target header supplied with normal credentials.
Reject reserved database identifiers as public credentials before native early-return paths.
Strip internal credentials and the target header before upstream forwarding.
Use request target metadata, not a capability/generation/ticket protocol.

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
