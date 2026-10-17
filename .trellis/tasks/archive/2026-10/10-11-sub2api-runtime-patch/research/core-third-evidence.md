# Research: Core C13-C17 Evidence

- Query: Compare the current `internal/eval` fingerprint, ModelTrace, Pelican, progress, and admission behavior with the original `internal/plugin` implementation and the pinned Sub2API v0.2.16 native account-test path.
- Scope: mixed internal source and pinned upstream source
- Date: 2026-10-15

## Findings

### C13: Fingerprint collection and attribution

The current `internal/eval/fingerprint.go:13-27` keeps the original quick, standard, and strict modes: 4x15, 8x25, and 16x25 samples. `:61-86` still selects a random prompt, sends low effort, classifies the response, and preserves numeric usage. `:101-172` still shuffles jobs, uses the two-worker default and six-worker maximum, shares six global sample slots, stops after eight consecutive request errors, drains admitted workers, and computes the same attribution tables. The original orchestration in `internal/plugin/fingerprint.go:16-28,128-209` has the same mode sizes, prompt selection, shuffle, concurrency, global slot, threshold, cancellation, usage aggregation, and attribution flow. The pure data/algorithm files (`fingerprint_algorithm.go` and embedded probe/baseline resources) are equivalent after the package-boundary change.

Current-only API surface is `FingerprintModes` at `internal/eval/fingerprint.go:98-99` and `FingerprintResult.Summary` at `:55-58`; repository search finds no production caller for either. They are API/list-detail plumbing, not fingerprint semantics. The current implementation calls `publish` during collection and comparison at `:165-168`; the original writes progress into its package-global state at `internal/plugin/fingerprint.go:199-207`.

Recommendation: preserve the prompts, probe data, normalization, scoring/attribution, mode sizes, low effort, retry/threshold behavior, six-slot limit, token sums, and settled cancellation behavior. Keep only the progress and list/detail surface required by the native page. Remove unconsumed exported mode/summary helpers if the final native handler does not need them; do not add a second fingerprint API or alter the eight-error rule.

### C14: ModelTrace collection and partial outcomes

`internal/eval/modeltrace.go:12-44` retains the three randomized challenges, lengths 292-332, default effort, one-to-three concurrency with default three, and six global slots; `:47-110` retains sample fields, retry count, text cap of 64,000, parsing and validity threshold. `:113-196` preserves challenge-indexed collection, usage aggregation, completion counts, sample order, attribution, and partial/failure selection. The original `internal/plugin/modeltrace.go:15-46,48-82,139-225` has the same challenge construction, sample fields, request/retry path, order restoration and analyzer. `modeltrace_algorithm.go`, the embedded bank, and `traceBankRevision` are equivalent at `internal/eval/modeltrace_algorithm.go:51,266` and the corresponding original file.

The original outcome distinction is material: `internal/plugin/modeltrace.go:214-223` marks the result failed when analysis fails, and marks it partial when analysis succeeds but fewer than three samples are valid. The current code keeps that distinction at `internal/eval/modeltrace.go:182-194`. Current `TraceResult.Done/Total/Valid/Errors` fields (`:56-61`) are consumed to publish progress at `:155-169`; they are not an independent UI history source. `Summary` at `:76-85` removes samples/rankings for list responses while the full result remains available for detail.

Recommendation: preserve the challenge bank, random lengths, prompt wording, default effort, validity threshold, 64,000-character cap, usage and sample order, bank revision, and partial/failure semantics. Keep one progress snapshot with the counts needed by the page; do not create a second persisted counter source. If final handlers derive list output from summaries, retain the small summary projection; otherwise remove the unused export rather than adding another endpoint.

### C15: Pelican prompt, extraction, and detail data

`internal/eval/pelican.go:12-17` is the original no-tools Pelican prompt and `:58-77` preserves the extraction order: labeled HTML code block, first markup code block, then bare markup. `:12` and `:43-55` retain the 256 KiB HTML limit, missing-HTML failure, duration and numeric usage. The original `internal/plugin/pelican.go:12-17,83-124` has the same prompt, extractor, and bound.

The current result adds explicit `ResultInfo`, status/error propagation and cancellation handling at `:19-55`, plus `Summary` at `:80-81` that omits HTML from list payloads. These are host/API plumbing. The original page already separated list records from detail HTML, so the summary shape is compatible with that existing list/detail behavior. Empty or oversized HTML is an evaluation failure and must not be repaired, rewritten, or retried as a different account.

Recommendation: retain the prompt, extraction rules, 256 KiB bound, HTML detail payload, and usage/duration. Keep explicit transport/cancellation status and the list summary only where the native response handler consumes them. Do not introduce a Pelican-specific cache, sanitizer, or alternate result format.

### C16: Progress callback boundary

`internal/eval/eval.go:33-47` defines a copied `Progress` value and `OnProgress func(Progress)`, with the callback excluded from JSON. `:137-200` emits an initial snapshot, collector/comparison snapshots, and one final snapshot after the evaluation result is known. `internal/eval/runs.go:428-447` installs the manager callback, folds Candy iterations, and updates the in-memory run snapshot under `stateMu`; the callback does not call the resolver, executor, or store.

The original package updates progress maps directly under its package mutex: Fingerprint uses `internal/plugin/fingerprint.go:199-207`; ModelTrace uses `internal/plugin/modeltrace.go:190-203` (and analogous runner state for the other types). There is no separate event-bus consumer. The native page reads manager state through polling: `internal/eval/web/app.js:175-187` maps `data.active[*].progress`, and `internal/eval/web/components.js:167-179` renders phase and `done/total` with the stop action.

Recommendation: retain one synchronous, copied callback as the boundary between the pure evaluation and RunManager. It must return promptly and remain free of I/O, persistence, resolver calls, executor calls, or waits. Preserve `cancelling` until workers settle and emit final counts only after settlement. Remove duplicate final-count assembly or fields with no page/handler consumer; do not add an event bus, background publisher, progress log, or new UI flow.

### C17: Administrator admission and native `/test` eligibility

Current admission in `internal/eval/runs.go:234-250` resolves each exact requested target before locking. `:270-294` deduplicates and reserves by `AccountID`, so one account cannot run twice across groups or evaluation types; `:321` records the same account-level busy key and `:412-419` releases it. `StartSummary` exposes `Started`, `Unavailable`, and `Busy` at `:51-55`. The current unavailable branch (`:280-289`) also appends `syntheticResult` (`:351-370`) to durable history without making an executor call.

The original cross-type exclusion is explicit in `internal/plugin/runs.go:9-12`; it uses one credential identity across Candy, Fingerprint, ModelTrace, and Pelican. Its catalog behavior at `:43-49` creates skipped records only when the caller supplies a known model catalog that excludes the selected model, and otherwise reports unchecked. The approved native administrator design removes that CPA catalog/all-wildcard path, so it should not be recreated in the Sub2API endpoint.

Pinned Sub2API v0.2.16 native evidence is in `backend/internal/server/routes/admin.go:393` and `backend/internal/handler/admin/account_handler.go:1289-1318`: `/api/v1/admin/accounts/:id/test` is protected by native admin authorization, parses the path ID/body, and calls `AccountTestService.TestAccountConnection`; the ordinary handler then invokes `RecoverAccountAfterSuccessfulTest` on success (`:1313-1317`). `backend/internal/service/account_test_service.go:356-427` first calls `accountRepo.GetByID(ctx, accountID)`, returns `Account not found` when absent, and routes by the stored account platform/protocol/type. `backend/internal/repository/account_repo.go:265-278` shows that lookup is an exact account read. `backend/internal/service/account.go:191-212` defines `IsSchedulable`, but that helper is used by normal gateway scheduling and is not called by `TestAccountConnection`; the native test entrypoint therefore has no explicit status, schedulable, expiry, quota, or model-catalog gate. Group identity is a page filter/display relation, not a sender authorization or billing identity in the approved administrator channel.

Recommendation: start only after native administrator authorization, exact account lookup, and the approved OpenAI/protocol capability check. Preserve account-level busy exclusion and collapse duplicate account IDs within one request. Keep the requested target in `Unavailable` with a safe admission error and never substitute another account or group. Do not reintroduce group-active/member, `IsSchedulable`, expiry, user quota, gateway billing, internal-user identity, or client model-catalog checks. Unavailable admission failures should remain summary-only unless the final UI/history contract explicitly requires fabricated records; the current synthetic history branch is not justified by the native `/test` contract and should be removed from the approved design. A request that actually starts and later fails in the native sender should produce the normal evaluation result with the received error and usage.

## Files Found

- `internal/eval/fingerprint.go` — current fingerprint orchestration, result summary, and progress publication.
- `internal/eval/fingerprint_algorithm.go` and embedded fingerprint data — current pure normalization, attribution, probes, and baselines.
- `internal/plugin/fingerprint.go` and `internal/plugin/fingerprint_algorithm.go` — original fingerprint orchestration and pure algorithm.
- `internal/eval/modeltrace.go` and `internal/eval/modeltrace_algorithm.go` — current ModelTrace orchestration, outcome status, summary, and bank revision.
- `internal/plugin/modeltrace.go` and `internal/plugin/modeltrace_algorithm.go` — original ModelTrace orchestration, partial semantics, and bank.
- `internal/eval/pelican.go` and `internal/plugin/pelican.go` — current/original Pelican prompt, extraction, bound, and result plumbing.
- `internal/eval/eval.go` — progress callback and evaluation dispatch boundary.
- `internal/eval/runs.go` — target resolution, account-level busy admission, progress wiring, synthetic unavailable result, and worker cleanup.
- `internal/eval/web/app.js` and `internal/eval/web/components.js` — polling and progress rendering consumers.
- `internal/plugin/runs.go` — original cross-evaluation busy and model-catalog skip behavior.
- `backend/internal/server/routes/admin.go`, `backend/internal/handler/admin/account_handler.go`, `backend/internal/service/account_test_service.go`, `backend/internal/repository/account_repo.go`, `backend/internal/service/account.go` — pinned native administrator `/test` registration, exact account lookup, dispatch, and scheduling helper distinction.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/native-connectivity-dispatch.md` — prior pinned-source account-test dispatch evidence and side-effect notes.

## Related Specs

- `.trellis/spec/backend/evaluation-core.md` — preserved evaluation algorithms, usage, cancellation, and progress contracts; later approved changes supersede only explicitly withdrawn transport limits or duplicate plumbing.
- `.trellis/spec/backend/run-state.md` — account-level exclusion, admission summaries, persistence, cancellation, and history ownership; C17's synthetic unavailable-history question remains a planning decision.
- `.trellis/spec/backend/host-target-adapter.md` — historical group/account resolver contract; its gateway-style active/member/schedulable checks are superseded for the administrator account-test channel.
- `.trellis/spec/backend/runtime-installation.md` — administrator route and native runtime ownership boundaries.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md` and `design.md` — current administrator endpoint, native sending, UI, and no-user-billing requirements.

## Caveats / Not Found

- The pinned Sub2API source was previously inspected and recorded in the task research; it is not present as a checked-out backend tree in this workspace, so this pass did not re-run source retrieval or a build.
- No product code, build, test, deployment, cleanup, or live provider request was run for this evidence.
- `FingerprintModes`, result `Summary` methods, and some result counters are current API/plumbing additions; their final retention depends on the native handler/list-detail implementation. No independent production consumer was found for the exported mode/summary helpers.
- Native `/test` can still have account-state side effects such as success recovery and 401/429 handling; the new full-response Codex eval path must follow the separately approved S04 decision rather than blindly calling the ordinary handler.
- Removing synthetic unavailable history requires the final history/UI contract to accept summary-only admission failures; no approved requirement currently demands fabricated records.
