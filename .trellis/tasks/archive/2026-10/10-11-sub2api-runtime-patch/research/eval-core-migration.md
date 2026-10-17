# Research: Sub2API Evaluation Core and Durable Run State

- Query: Define an implementation-ready boundary for moving the four reusable CPA evaluations into `internal/eval`, adapting model execution to the normal Sub2API Responses endpoint, and adding durable run state without changing CPA production behavior.
- Scope: internal
- Date: 2026-10-12

## Findings

### Source and copy boundary

The pinned build helper already defines the intended source/destination contract. `compatibility/sub2api-v0.2.15.json:64-68` declares `internal/eval` as the required source and `backend/internal/codexeval` as the generated destination. `scripts/sub2api_build.py:303-318` requires the source directory to contain Go files. `stage_eval_core` at `scripts/sub2api_build.py:456-470` resolves both paths, rejects source symlinks and symlink components, rotates an existing destination into a workspace backup, and copies with `shutil.copytree`. It does not transform imports, select files, or generate a package boundary. The current destination package therefore needs to be self-contained and buildable after a straight tree copy.

`scripts/build-sub2api.sh:167-175` stages the tree after applying the versioned upstream patches and before Wire generation; `scripts/build-sub2api.sh:197-202` builds the patched backend with `CGO_ENABLED=0`, the manifest's `embed` tag, and the pinned Go version. `scripts/test-sub2api.sh:53-69` treats the generated destination and patch marker as prerequisites for integration checks. The source tree in `internal/eval` should therefore have no import of the CPA `plugin` package, no host ABI dependency, and no path that depends on the ignored upstream checkout.

Recommended maintained source layout (all copied by the existing helper):

- `internal/eval/candy.go`: candy prompt, result shape, verdict predicate, and pure evaluation operation.
- `internal/eval/fingerprint.go` and `fingerprint_algorithm.go`: collection orchestration and pure normalization/statistics/attribution.
- `internal/eval/modeltrace.go` and `modeltrace_algorithm.go`: challenge collection, retry aggregation, scoring, attribution, and bank revision.
- `internal/eval/pelican.go`: prompt, HTML extraction/size validation, and pure result operation.
- `internal/eval/codex.go`: Codex context/request construction and model/effort normalization.
- `internal/eval/model.go` (or a split transport file): request encoding, local Responses executor contract, incremental SSE parser, terminal/error classification, and response-to-`ModelResponse` conversion.
- `internal/eval/state.go` and `runs.go`: the 4B run manager, active intent, history snapshots, and recovery. These should be introduced after the 4A interface is fixed.
- `internal/eval/data/codex_context.json`, `fingerprint_probes.json`, `fingerprint_baselines.json`, and `modeltrace/unified_bank.json`: unchanged embedded inputs.
- `internal/eval/data/modeltrace/LICENSE` and `README.md`: retain the ModelTrace MIT notice, pinned source commit, and bank provenance.
- `internal/eval/web/` and `web/vendor/`: retain the existing page, Markdown/Math assets, and vendor notices when the embedded page is moved. Admin route/menu/theme/bootstrap remains Stage 5; copying the assets is a packaging responsibility, not a claim that the page is host-integrated.

The original `internal/plugin` package and its tests remain the CPA implementation. No source move or shared import should be used: changing that package would alter the existing CPA host behavior and would make the generated Sub2API tree depend on CPA globals. A later core test suite can be copied or rewritten under `internal/eval`, but it must use a local fake executor and fake store rather than CPA's host callback.

### Four evaluations and preserved assets

The reusable algorithm behavior is source-backed by these anchors:

- Candy: `internal/plugin/candy.go:17-24` contains the Chinese prompt and table; `:100-111` calls `askCodex`, records token counts/duration, and accepts an answer only when `hasStandalone21` at `:114-122` finds a standalone `21`. The prompt, verdict, result fields, and ten-run cap (`:10`) are part of evaluation semantics.
- Fingerprint: `internal/plugin/fingerprint.go:16-28` defines quick/standard/strict modes, `:91-125` validates the request and chooses per-credential concurrency, and `:128-209` runs shuffled probes, gathers samples, tracks errors/cancellation, and computes attribution. `internal/plugin/fingerprint_algorithm.go:18-45` embeds the two JSON datasets; all normalization, comparison, permutation, and attribution functions in that file are pure and should remain in the same package after migration.
- ModelTrace: `internal/plugin/modeltrace.go:15-26` fixes three challenges per account and concurrency limits, `:29-46` generates the challenge prompts, and `:139-225` runs three probes concurrently, preserves challenge order, sums usage across retries, and computes attribution. `internal/plugin/modeltrace_algorithm.go:16-20` embeds the bank and MIT license, and `:50-51` derives a SHA-256 bank revision used in results. Keep this revision in persisted records so a later bank update does not make old results ambiguous.
- Pelican: `internal/plugin/pelican.go:12-17` defines the prompt and 256 KiB HTML limit; `:83-101` extracts and validates HTML while preserving usage/error metadata; `:104-123` implements the code-block/bare-markup parser. The preview's opaque-origin sandbox is in `internal/plugin/web/pelican.js:28-30`; it belongs to the later page integration but must not be weakened during the move.

The vendor files are functional inputs, not optional UI decoration: `internal/plugin/plugin.go:31-102` embeds the HTML, Markdown-it, CJK, TeX math, Temml, and four evaluation scripts into one page. `internal/plugin/web/vendor/markdown-it-cjk-friendly.min.js` and `texmath.js` carry MIT notices; `temml.css` identifies Temml 0.13.5 and MIT; the ModelTrace license is explicitly embedded by `modeltrace_algorithm.go`. Keep their notices and do not edit minified vendor files during this migration.

### CPA-specific behavior to remove

The current model boundary is CPA-specific. `internal/plugin/model.go:90-163` sends `host.model.execute_stream` with `forced_provider` and `auth_id`, receives a `stream_id`, polls `host.model.stream_read`, and closes it through `host.model.stream_close`. `internal/plugin/codex.go:183-193` creates the request and then starts `withModelTimeout(context.Background(), answerTimeout)`; this independent background context defeats a caller's cancellation. `internal/plugin/credentials.go:18-34` and `:105-169` use CPA credential files and host callbacks, while `:180-195` filters disabled/unavailable credentials before a run. None of these callback or credential types should cross the new core boundary.

Keep the following data but replace its ownership:

- Keep a stable internal account reference and display metadata, but obtain it from a Sub2API adapter. Numeric group/account IDs are API-facing; the core may use a stable account namespace keyed by the numeric account ID.
- Keep `codex_context.json`, `codexTurn`, `codexPrefix`, the fixed Codex headers, and the prompt bodies. The executor must own authorization and routing headers.
- Keep `executeProbe`'s retry policy from `internal/plugin/model.go:44-85`: at most three attempts, 2 s then 4 s waits, retry only status 0, 429, or 5xx, and stop immediately on caller cancellation. Each retry must reuse the same selected group/account target and request body; it must not reselect an account or group.
- Keep the distinction between requested and effective model/effort. `validModelName` and `credentialSupportsModel` show that CPA permits a trailing parenthesized effort suffix and strips it for catalog comparison (`internal/plugin/model.go:24-27`, `internal/plugin/credentials.go:200-219`). The Sub2API wire request must normalize a known CPA suffix such as `model(high)` to a plain model plus a separate reasoning effort. Record both the user-requested values and the effective wire values. If an explicit effort conflicts with a suffix, reject the request or apply one documented deterministic rule; do not silently send parentheses to Sub2API.

### 4A executor and SSE contract

Stage 4A should expose a small pure-core contract that Stage 4B can use directly. A suitable shape is:

```go
type Target struct {
    GroupID   int64
    AccountID int64
}

type ResponsesRequest struct {
    Target Target
    Model  string
    Effort string
    Body   []byte
    Header http.Header // Codex request headers only; no Authorization or routing headers
}

type ModelResponse struct {
    Answer          string
    InputTokens     int64
    OutputTokens    int64
    ReasoningTokens int64
}

type ModelExecutor interface {
    Execute(ctx context.Context, req ResponsesRequest) (ModelResponse, error)
}
```

The exact exported names can follow local Sub2API style, but the ownership is material. The core receives a parent run context and a target; it never receives or stores the private evaluation token. The Sub2API adapter owns a configured normal HTTP client, adds `Authorization` and exactly one `X-Codex-Eval-Group-ID`/`X-Codex-Eval-Account-ID` pair, and POSTs to `/v1/responses`. The adapter must use a fixed, explicitly configured container-local base URL, disable environment proxy lookup (`Transport.Proxy = nil`), and reject redirects rather than forwarding Authorization to another origin. It must pass the caller's context to `http.NewRequestWithContext`, apply the per-request timeout as a child context, and close the response body on every path. A constructor may accept the private token, but the token remains in the adapter closure/configuration and is never present in `ResponsesRequest`, result records, logs, browser DTOs, or error strings.

This adapter does not decide readiness or topology. K4 source-IP/proxy proof and private-Key initialization remain later bootstrap/wiring work. The core should fail if no executor is supplied, but it must not call a direct provider or bypass normal authentication and billing.

The parser must preserve the observable semantics of `internal/plugin/model.go:186-249` while making memory and stream behavior explicit:

1. Read SSE incrementally across arbitrary chunk boundaries. A line may be split between reads; a new `data:`, `event:`, `id:`, `retry:`, or comment field may begin in the next chunk.
2. Enforce a maximum buffered response/event size. The current CPA implementation accumulates an unbounded `events` slice (`:187-210`); the port must replace this with an incremental scanner/reader and a documented bound. Oversized lines/events are errors and never successful answers.
3. Treat `response.completed`, `response.incomplete`, and `response.failed` as terminal response events. Only a completed response with acceptable status is successful. `response.incomplete`/`response.failed` return an error with any authoritative usage observed so far.
4. Treat an SSE `error` event, non-2xx HTTP status, body read error, or context cancellation as failure. No missing terminal event may be turned into a completed answer.
5. Preserve `response.output_item.done` items and reconstruct output from those items when the terminal response has an empty `output` array, as the existing code does (`:227-249`). Collect text content only from message parts of type `output_text` or `text`.
6. If no text exists, distinguish tool-only output (an output item type ending in `_call`) from a generic empty response. Both are failed evaluations, with the tool-only diagnostic preserved. Do not treat tool calls as a valid answer or fabricate text from arguments.
7. Preserve authoritative usage fields and do not invent totals when the stream ends before usage is available. The core result can carry partial usage internally for the gateway/accounting layer, but a normal completed evaluation record is created only by a terminal-success path.

The 4A tests should use an `httptest`-style local fixture or an injected round-tripper, never a paid model request. Required assertions are terminal success, chunk-split SSE, incomplete/failed terminal errors, SSE error, non-2xx, no-terminal EOF, oversized event, tool-only output, item-only reconstruction, cancellation before first output, cancellation after partial output, and that the Authorization/routing headers are present only on the local request. Include an explicit redirect/proxy test proving the private token cannot leave the configured local origin.

### 4B run manager and durable state

Stage 4B should consume only the 4A `ModelExecutor` and metadata interfaces. It owns run contexts, account exclusion, active intents, final records, and the file store. It should not duplicate HTTP parsing or model request construction.

The current cross-type exclusion is explicit in `internal/plugin/runs.go:9-12`: one credential is busy when any candy, fingerprint, ModelTrace, or pelican run exists. Preserve this as `accountBusyLocked(accountID)` over a unified active map or four maps keyed by the stable numeric account ID. Group selection is request-scoped and must not be part of mutable global busy state; the same account remains excluded even if two requests name different groups.

The current `startRuns` calls `selectedCredentials` before locking and silently drops disabled/unavailable entries (`internal/plugin/runs.go:24-57`, `internal/plugin/credentials.go:180-195`). The new target resolver should retain explicitly requested unavailable accounts and append one failed record for each, including selected group/account IDs and the requested/effective model. It may omit unavailable accounts from an `all` selection according to documented UI semantics, but an explicit ID must never disappear without a record. A selected account that is not a member of the selected group, fails model/capability eligibility, or is missing must likewise become a failure for that target; no alternate group/account is permitted.

Use a persisted state shape with a version field and separate active intents from final histories. A minimal conceptual shape is:

```go
type persistedState struct {
    SchemaVersion int                         `json:"schema_version"`
    Results       map[string][]Result         `json:"results"`
    Fingerprints  map[string][]Fingerprint    `json:"fingerprints"`
    ModelTraces   map[string][]ModelTrace     `json:"modeltraces"`
    Pelicans      map[string][]Pelican        `json:"pelicans"`
    Active        map[string]ActiveRunIntent  `json:"active"`
}
```

The exact record types can retain the current JSON fields, but every final result must also carry the selected group/account identity and a run ID. `ActiveRunIntent` should capture evaluation type, run ID, selected targets, requested/effective model and effort, batch controls, and creation time. It does not contain authorization material or full request bodies.

The startup sequence must be transactional:

1. Resolve and validate metadata outside the state mutex.
2. Under the state mutex, reject quiescing/busy conflicts and construct active intent entries for every worker, including explicit unavailable failures if they are part of the request.
3. Write the complete snapshot atomically before acknowledging `/run`. If writing fails, restore the prior in-memory state, start no workers, and return a storage error.
4. Only after the write succeeds, add tracked workers and launch them. Worker completion removes its active entry, appends a final record, trims to the latest 20 for that account/evaluation type, and writes another atomic snapshot.

This ordering prevents a successful acknowledgement with no durable evidence of a paid run. It also means an asynchronous worker must not be started from a callback that runs before the initial `saveState` completes.

On startup, every active intent from the previous process is reconciled to an `interrupted` final record (or one interrupted record per target) and removed from `Active`; no executor call or retry is made. The recovery write must itself be atomic. If the state file is malformed, unreadable, a symlink, or cannot be safely decoded, preserve it and set a visible storage error. Do not replace it with an empty state. This follows the existing protection in `internal/plugin/state.go:27-66` and `:87-107`, where `stateLoadError` blocks future saves and writes use a temporary file plus rename.

Retain `historyLimit = 20` semantics from `internal/plugin/state.go:68-79`, independently for each account and evaluation type. `clearHistory` must preserve active intents and in-flight maps, and restore the previous snapshot on a write failure as the existing `clearHistoryResponse` does at `:134-158`. Use a dedicated state path under `/app/data/codex-eval/`, created with restrictive directory/file permissions by the later runtime wiring; the core's store should accept an injected path so tests never touch deployment data.

Cancellation and shutdown must cover all four evaluation types. Fingerprint and ModelTrace already have cancellable contexts (`internal/plugin/fingerprint.go:120-125` and `modeltrace.go:126-131`), but candy only checks `quiescing` between requests (`candy.go:75-98`) and pelican has no cancellation path (`pelican.go:73-81`). The new run manager must create a context per run, pass it to every `askCodex`/executor call, and expose cancellation for candy, fingerprint, ModelTrace, and pelican. Cancellation before output, after partial output, during probe retry backoff, and during pelican generation must settle a non-completed record with the observed usage and cancellation status. Quiesce must mark all active runs cancelling, cancel their contexts, wait for tracked workers, persist terminal interruption/cancellation state, and only then let the local HTTP gateway shut down. This is required because the worker's executor calls the local `/v1/responses` route.

### Recommended ordered Stage 4 dispatch

Split the stage at the existing contract boundary, with independent specification and quality review after each part:

**Stage 4A: transport and pure evaluation core.** Copy/re-home the four algorithms, Codex context builder, assets/licenses, model/effort normalization, and the `ModelExecutor` plus bounded SSE parser. Add fake-executor tests for all four pure result paths and real HTTP fixture tests for SSE/cancellation. Do not add admin routes, target discovery, active state, or host readiness. The 4A entry point is a reusable `RunOne`/probe operation that accepts a caller context, immutable account metadata, model/effort, and executor; 4B can call it directly.

**Stage 4B: run manager and durable state.** Add target metadata adapter interfaces, selected-unavailable failure records, per-account cross-type exclusion, active intent transaction, recovery, retention, clear-history rollback, cancellation, and quiesce. Wire it to the 4A executor without reimplementing parsing. Add store and lifecycle tests using temporary directories and fake executors. Administrator routes/UI, key bootstrap, topology proof, Wire, and installer remain later stages.

This split is meaningful because 4A can be validated entirely with deterministic fixtures and has no durable scheduling state, while 4B has the high-risk crash/recovery and concurrency contract. A single large dispatch would make failures ambiguous and would encourage temporary scheduler interfaces that are later discarded. The split does not create a second public API: 4A's executor and pure operation interfaces are the stable inputs of 4B.

## Source-Supported Facts vs Assumptions

Source-supported facts:

- The build helper performs a guarded copy of `internal/eval` into `backend/internal/codexeval`; it does not rewrite package imports (`scripts/sub2api_build.py:456-470`).
- CPA model execution uses streamed Responses events, requires a terminal event, reconstructs output items, and rejects incomplete/tool-only/no-text responses (`internal/plugin/model.go:90-163`, `:186-249`).
- CPA Codex requests carry captured context, fixed headers, and a separate effort field; `askCodex` currently owns an independent background context (`internal/plugin/codex.go:131-193`).
- Fingerprint/ModelTrace retry the same probe up to three attempts with 2 s/4 s backoff; candy/pelican perform one request per evaluation (`internal/plugin/model.go:44-85`, runner files above).
- Existing state is atomic final-history JSON with latest-20 trimming and corrupt-file preservation, but no active-run persistence (`internal/plugin/state.go:13-107`).
- Existing busy detection spans all four evaluation types (`internal/plugin/runs.go:9-12`), while explicit unavailable targets are currently filtered out (`internal/plugin/credentials.go:180-195`).
- The pinned Sub2API build is Go 1.27.2, static, and frontend-embedded (`dist/sub2api-build/workspace/backend/go.mod:1-4`, `scripts/build-sub2api.sh:177-202`).

Implementation assumptions requiring later verification:

- The local Sub2API listener address/port and whether it is reachable at loopback from the server process are deployment facts. The 4A adapter should require an explicit configured local URL; K4 topology/bootstrap work must prove it before activation.
- The exact Sub2API administrator metadata DTOs and account/group service interfaces are Stage 5/2A integration concerns. The 4B core should use narrow internal metadata interfaces and numeric IDs rather than importing handler or repository types.
- A redirect/proxy-free local HTTP client is the intended privacy boundary, but the final adapter must be checked against the pinned runtime's listener and middleware behavior. No core test can prove a real deployment's proxy topology.
- The persisted active-intent schema version and interrupted-record presentation need to be chosen during implementation; they must preserve existing UI JSON fields and the latest-20 invariant.

## Caveats / Not Found

- No product files were changed and no tests were run during this research-only pass.
- No live Sub2API container, database, account, balance, subscription, or local listener was available. Billing, target eligibility, and real `/v1/responses` accounting remain later integration/live checks.
- The upstream source tree contains many provider-specific SSE adapters, but this migration is OpenAI/Codex-only. No provider expansion or reuse of unrelated adapters is justified by this research.
- The current browser/UI source still assumes CPA routes and authentication. Its route, menu, theme, and iframe work belongs to the planned administrator integration stage; this report only records the assets and sandbox behavior that must survive the core copy.
