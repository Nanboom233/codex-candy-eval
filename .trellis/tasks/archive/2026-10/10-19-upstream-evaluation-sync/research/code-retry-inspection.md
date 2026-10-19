# Research: Native Retry, Request, and Inspection Contracts

- Query: Trace native model requests and Candy inspection, then map the approved upstream retry behavior to the existing Sub2API code.
- Scope: internal / mixed
- Date: 2026-10-19

## Findings

### Critical plan correction

The current native sender does not expose one error category for every failure. `SendCodexEvaluation` returns `(*http.Response, error)`. It returns a non-200 response with a nil error after reading and restoring the body. It returns a nil response and an error for both setup failures and transport failures. See `/tmp/sub2api-settings-entry-final/source/backend/internal/service/account_codex_eval_service.go:56-72,84-137,183-212` and `patches/0004-admin-channel.patch:211-365`.

The current `executeProbe` treats status `0`, `429`, and all 5xx statuses as retryable. Status `0` includes JSON encoding errors, account lookup failures, missing credentials, invalid base URLs, request conversion failures, and native transport errors. This does not satisfy the approved rule that setup failures are final and 429 is final. See `patches/0001-evaluation-core.patch:998-1053`.

The minimal valid correction is one native error type in the existing HTTP upstream port. Declare `HTTPUpstreamRequestError` in `service/http_upstream_port.go` with `Error`, `Unwrap`, and `RequestAttempted() bool`. In `repository.doUpstreamRequest`, wrap only an error returned by the actual `servertiming.Do` call. `HTTPUpstream.Do` and `DoWithTLS` perform host validation and acquire/build the client first; proxy normalization, TLS/client setup, and host-policy failures must stay unmarked and final. `codexeval` can recognize the marker through a small interface and `errors.As`, without importing `service`. Do not add another HTTP client, gateway retry layer, or generic retry service.

`PluginTransportError.RequestSent` describes replay safety, not transience. A plugin error with `RequestSent=true` before a valid HTTP response must remain final, preserving the native no-replay contract. The current plugin error frame and `normalizePluginRPCError` do not expose a structured transient category, so unknown pre-response plugin setup/frame errors also remain final. Once valid HTTP response headers arrive, an unusable or interrupted response stream follows the approved stream-retry rule. The runner must stop when its parent context is cancelled. A per-attempt timeout is retryable only when it carries the direct HTTP request-attempt marker or occurs after valid HTTP response headers; `normalizePluginRPCError` can return plain `ctx.Err()` after plugin context expiry and lose the plugin replay metadata, so do not retry every pre-response attempt deadline globally.

### Current native call graph

The native administrator path is:

```text
admin start handler
  -> codexeval.Start
  -> startChecked
  -> active[account] and run context
  -> goroutine for one account
     -> runCandyChecked / runFingerprint / runModelTrace / runPelican
        -> request runner
           -> executeModel
              -> Sender(account ID, JSON body, Codex headers)
                 -> AccountTestService.SendCodexEvaluation
                    -> account lookup and credential/setup checks
                    -> doOpenAIAccountTestUpstream
                       -> PluginManager.RoundTripOpenAIOAuth or HTTPUpstream.DoWithTLS
```

The dispatch and admission code is in `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/runtime.go:84-201`. It validates the request, resolves the account name, rejects an active account, creates a context, and starts one worker. The active map is shared across test types. Cancellation calls the run cancel function at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/runtime.go:203-213`.

The current Candy path is:

```text
runCandyChecked
  -> evaluateCandy
     -> askCodex
        -> codexTurn once
        -> executeModel once
```

`runCandyChecked` checks participation before each logical Candy test, evaluates one answer, saves one result, and increments `Done` once. See `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/candy.go:44-87`.

The current fingerprint path is:

```text
runFingerprint worker
  -> acquire fingerprintSlots
  -> collectFingerprintSample
     -> codexTurn
     -> executeProbe
        -> executeModel for up to three status-based attempts
  -> release fingerprintSlots
```

The current ModelTrace path has the same `executeProbe` dependency. `collectFingerprintSample` is at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/fingerprint.go:162-182`. `collectTraceSample` is at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/modeltrace.go:193-211`. These are the only current `executeProbe` callers. The current `askCodex` callers are Candy and Pelican at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/candy.go:75-87` and `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/pelican.go:53-76`.

After the approved manual test replacement, remove `modeltrace.go` and its callers. Replace the fingerprint caller with the new detector's three challenges. Remove `executeProbe` after all callers move to the shared `askCodex` runner. Do not keep a compatibility wrapper that retains old status-0 retry semantics.

### Current request and response signatures

The native model package currently defines:

```go
type modelResponse struct {
    Answer          string
    InputTokens     int64
    OutputTokens    int64
    ReasoningTokens int64
}

func executeProbe(
    context.Context, credential, string, any, http.Header, chan struct{},
) (modelResponse, int, error)

func executeModel(
    context.Context, credential, string, any, http.Header,
) (modelResponse, int, error)

func askCodex(credential, string, string, string) (modelResponse, time.Duration, error)
```

See `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/model.go:17-85,87-146` and `codex.go:130-189`.

The upstream shared runner has the correct logical contract:

```go
func askCodex(
    ctx context.Context,
    auth credential,
    model, effort, instructions, prompt string,
    limit time.Duration,
    accept func(string) error,
) (out modelResponse, attempts int, err error)
```

It creates the request once, creates a fresh timeout context per attempt, counts attempts before sending, accumulates usage, keeps the last answer, and calls `accept` only after a successful answer. See upstream `7a49068:internal/plugin/model.go:17-98`.

Use the upstream signature with one native-only pre-send callback added at the end, or use a small request options struct if the implementation already needs a callback for all callers. The callback must run before every supplier send, not only before the logical test. Keep the callback out of the upstream detector validator. The minimal direct shape is:

```go
func askCodex(
    ctx context.Context,
    auth credential,
    model, effort, instructions, prompt string,
    limit time.Duration,
    accept func(string) error,
    beforeSend func() error,
) (modelResponse, int, error)
```

A callback error must have a distinct internal stop result. It must not become a saved failed Candy result. A small `inspectionStopError` sentinel, handled by an `evaluateCandyChecked` helper that returns `(result, save bool)`, is sufficient. Do not cancel the whole inspection context for account opt-out because that would cancel unrelated accounts and mark the round cancelled.

### Proposed runner behavior

Keep request construction before the loop:

```text
body, headers := codexTurn(...)
for attempt = 1..3:
    if parent context is done: return cancellation, attempts so far
    if beforeSend exists and returns stop: return inspection stop, no save
    requestCtx = WithTimeoutCause(parent, limit, timeout cause)
    response, status, err = executeModel(requestCtx, ...)
    cancel requestCtx
    add response usage
    keep response answer
    if err is nil and accept exists:
        err = accept(answer)
        if err != nil: mark validator rejection retryable
    if err is nil: return
    if parent context is done: return cancellation
    if err is final: return
    if attempt == 3: return
    wait 2s, then 4s, with parent cancellation select
```

The first wait is two seconds. The second wait is four seconds. The upstream loop multiplies the delay by two after each wait at `7a49068:internal/plugin/model.go:73-97`.

The runner must preserve the same request body and headers on every attempt. The upstream test asserts byte-identical request bodies at `7a49068:internal/plugin/model_test.go:88-132`. A new UUID, thread, or installation must not be generated inside the retry loop.

### Status and error classification

Use this decision table at `executeModel` and the runner boundary:

| Boundary result | Current evidence | Approved action |
| --- | --- | --- |
| JSON encoding error before sender | `model.go:94-99` | Final. No supplier request was made. |
| Account lookup, unsupported account, missing token/key, invalid base URL, request conversion, or agent setup error | `account_codex_eval_service.go:60-137,152-161` | Final. Preserve native error. |
| Direct HTTP upstream returns `HTTPUpstreamRequestError` from the actual `client.Do` call | `http_upstream.go:299-310`; the marker is owned by `http_upstream_port.go` | Retry up to the limit unless parent context is done. |
| Plugin returns a pre-response `PluginTransportError` | `plugin_runtime.go:274-335,357-387`; `RequestSent=true` protects against replay | Final. Do not infer transient status from `RequestSent` or plugin error text. |
| Non-200 response with status 408 | Sender preserves response status at `account_codex_eval_service.go:197-212` | Retry as a timeout. |
| Non-200 response with status 5xx | Same | Retry as a server failure. |
| Non-200 response with status 401 | Sender records account error at `account_codex_eval_service.go:207-209` | Final. Do not resend. |
| Non-200 response with status 429 | Sender reconciles cooldown at `account_codex_eval_service.go:204-206` | Final. Do not resend. |
| Other 4xx response | Sender preserves status | Final. |
| Response stream ends early, emits an error, is incomplete, failed, blank, or tool-only | `model.go:115-145,195-276`; upstream retryable wrapper at `7a49068:model.go:143-153` | Mark retryable and retry. |
| Valid fingerprint or Pelican answer rejected by validator | Upstream wraps `accept` error at `7a49068:model.go:84-89` | Retry. |
| Valid Candy answer, including wrong text | Native `hasStandalone21` grading at `candy.go:75-87`; upstream Candy has nil validator at `7a49068:candy.go:67-81` | Complete one logical result. Do not retry. |
| Attempt timeout during a marked direct HTTP request or after valid HTTP headers | Per-attempt `WithTimeoutCause` in upstream `7a49068:model.go:72-80` | Retry while the parent context remains active. |
| Attempt timeout before plugin response headers | `normalizePluginRPCError` can return plain `ctx.Err()` and discard `PluginTransportError` at `plugin_runtime.go:379-387` | Final. Do not globally retry every attempt deadline before a response. |
| Parent cancellation or explicit stop | `executeModel` closes response with `context.AfterFunc` at `model.go:100-106`; run cancel at `runtime.go:203-213` | Stop active stream, wait, and later sends. Save only the native cancellation result where the existing caller owns it. |

Do not retain `status == 0` as a blanket retry condition. A plain sender error is final unless the sender marks it as a transport failure. Do not classify by English or Chinese error text.

### Native sender boundary and setup distinction

`SendCodexEvaluation` has a clean setup region before `doOpenAIAccountTestUpstream` and a clean upstream-call boundary at lines 183-189. Early returns cover malformed JSON, missing model or stream, account lookup, account type, shadow resolution, missing token or key, invalid base URL, JSON conversion, agent identity authentication, and request construction. Leave these errors unmarked and final.

The upstream call can use the plugin manager or native HTTP transport at `/tmp/sub2api-settings-entry-final/source/backend/internal/service/openai_plugin_transport.go:21-44`. For HTTP transport, `HTTPUpstream.Do` and `DoWithTLS` run host validation and client acquisition before `doUpstreamRequest`; client acquisition normalizes the proxy and constructs or selects transport configuration at `/tmp/sub2api-settings-entry-final/source/backend/internal/repository/http_upstream.go:203-239,242-297,655-730`. Therefore, wrapping `doOpenAIAccountTestUpstream` in the account service would incorrectly mark those setup errors. The exact small boundary is the `servertiming.Do` call at `repository/http_upstream.go:299-310`, where a failed `client.Do` is wrapped as `service.HTTPUpstreamRequestError`. The repository already imports the service package. `codexeval` checks the marker interface, so it does not gain a service dependency.

For the plugin path, `PluginTransportError.RequestSent` is explicitly used to prevent replay after ambiguous metadata delivery at `plugin_runtime.go:298-302` and `PLUGIN_DEVELOPMENT.md:67`. It is not a transient classifier. `normalizePluginRPCError` stores the underlying gRPC error in a message string rather than an unwrap-able cause, while plugin error frames expose only arbitrary `code`, `message`, and `request_sent` fields (`plugin_runtime.go:312-335,357-387`; `plugin.proto:138-149`). No existing structured field distinguishes transient pre-response plugin failures from permanent plugin/setup errors. Keep pre-response plugin errors final, including `RequestSent=true`; permit the approved retry for a stream failure only after valid HTTP response headers have been returned.

For non-200 responses, `SendCodexEvaluation` currently reads and restores the body, returning the response with nil error; a body-read error currently discards both response and status at `account_codex_eval_service.go:197-212`. Preserve the native change narrowly: return `resp` with the read error. In `executeModel`, capture status and arrange body closure before handling that error. The retry runner must prioritize status: 401 and 429 are final even when reading their body failed; 408 and 5xx remain retryable; other HTTP failures remain final. Successful-response stream read errors are retried as unusable streams. This avoids classifying body read errors by text or losing the response's status.

### Inspection timeout interaction

The existing `answerTimeout` is ten minutes and currently wraps the whole one-shot Candy call in `askCodex` at `codex.go:182-189`. The approved retry runner creates a fresh ten-minute timeout for each attempt; with three attempts and the approved 2s/4s delays, one logical Candy call could take at most 30 minutes 6 seconds when its parent remains active.

Automatic inspection keeps `phaseLimit` at 30 minutes (`inspection.go:94-98`). It starts account workers with child contexts of the inspection run context in `runtime.go:154-159`; the `phaseCtx` is created later and limits only the coordinator's wait in `testPhase` (`inspection.go:635-677`). At the 30-minute coordinator-wait deadline, it cancels each still-active account run and waits for those workers to exit. This is a 30-minute wait limit, not a strict 30-minute wall-clock cap: worker startup precedes it, and cancellation/drain follows it. A single retrying logical Candy request has a theoretical 30-minute-6-second per-attempt budget, but inspection normally cancels active work when this phase wait expires. The phase deadline is not inherited directly by each request, so retain the account run context as the per-attempt timeout parent and preserve its explicit cancellation behavior.

Each `testPhase` has its own 30-minute coordinator-wait limit. The inspection can run an initial test phase and up to two confirmation phases, with the existing 240-second recovery wait before normal-state confirmation. If both confirmation groups run, the configured wait budgets total 90 minutes plus the four-minute recovery wait, in addition to worker startup and any cancellation/drain time. This is not a strict wall-clock bound because each expired phase has an unbounded drain loop until its workers exit. Do not change `phaseLimit`, recovery delay, inspection scheduling, or parent-context ownership in this task.

### Stream lifetime and cancellation

The native `executeModel` defers `response.Body.Close`, registers `context.AfterFunc` to close the body when the attempt context ends, reads usage before grading, and returns the context cause if cancellation wins. See `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/model.go:100-145`.

The native sender returns a tracked body. `cancelOnCloseBody.Close` cancels the upstream request before closing and waits for an active read, preventing a blocked stream from surviving a retry or explicit stop. See `/tmp/sub2api-settings-entry-final/source/backend/internal/repository/http_upstream.go:299-347`. Preserve this ownership. The retry runner must not close a body that `executeModel` already owns, and it must call the per-attempt timeout cancel immediately after `executeModel` returns.

For a cancelled retry wait, the timer select returns the parent context error and sends no later request. For cancellation while waiting for a fingerprint slot, the worker's existing select at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/fingerprint.go:111-123` returns without calling the supplier. Keep this select. Once a worker owns a slot, it must retain that slot through all attempts and backoff, then release it after `askCodex` returns. This matches the approved three-challenge concurrency and upstream's logical-challenge slot behavior.

The current `executeProbe` releases `fingerprintSlots` after every attempt because it owns the slot around one call. Deleting it and keeping the slot around `collectFingerprintSample` changes the lifetime correctly without adding a slot parameter to `askCodex`.

### Attempts, usage, and result counting

The current `modelResponse` carries input, output, and reasoning token counts. `executeProbe` already accumulates those fields across attempts at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/model.go:66-71`, but current `askCodex` discards the attempt count and performs one request at `codex.go:182-189`.

The upstream shared record has common usage fields and stores `Attempts` only on new Candy, fingerprint-sample, and Pelican records. See `7a49068:state.go:47-56`, `7a49068:candy.go:22-28`, `7a49068:fingerprint.go:76-93`, and `7a49068:pelican.go:21-27`.

Native current Candy and Pelican records store usage directly but have no `Attempts` field. See `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/candy.go:23-37` and `pelican.go:17-31`. Add `Attempts int json:"attempts,omitempty"` to new records and assign the shared runner's count. Keep existing fields and old JSON records readable. Do not rewrite old Candy or Pelican records.

One logical Candy test still saves one final record. `runCandyChecked` increments `candyRunning.Done` once after `evaluateCandy`, regardless of internal attempts. The inspection collector counts records, not HTTP sends: `r.Tests += len(results)` and `out = append(out, r.OK)` at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/inspection.go:604-633`. Therefore a retry sequence that ends in a valid wrong answer produces one completed logical result with `OK=false`, one `Tests` count, and one confirmation answer.

An error after all attempts produces one failed result with accumulated usage and attempt count. It is excluded from confirmation because `collect` requires `Status == "completed"`, no error, and not skipped at `inspection.go:613-633`. An opt-out stop is different: it must produce no synthetic result and leave the inspection outcome `unchanged`.

Fingerprint samples store per-challenge `Attempts`, accepted text, parsed count, and error. The new record totals usage from each event. See upstream `7a49068:fingerprint.go:76-83,135-231`. Three logical challenges can therefore create at most nine supplier requests while the result still contains three logical samples.

### Participation callback and no synthetic failure

The existing inspection callback is created in `startChecked` and passed into `runCandyChecked` only before each logical test at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/runtime.go:183-192`. `testPhase` defines the callback and updates the result at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/inspection.go:574-603`.

Current opt-out behavior is correct between logical tests: the callback sets `Outcome = "unchanged"`, returns an error, and `runCandyChecked` exits before evaluating or saving a result. The existing test asserts no synthetic error at `tests/backend/internal/codexeval/inspection_test.go:540-625`.

The approved retry change requires the same callback before every attempt. Do not call `hooks.Get` while holding the global evaluation mutex. Keep the callback closure in `testPhase`, but pass it into the shared runner through the Candy evaluation path. Use a sentinel stop error for intentional opt-out. Handle that sentinel before converting the runner error into `candyResult.Error`, and return from the logical worker without saving a record. Keep account-state or priority-change errors as the existing backend result error because those are not intentional opt-out.

The callback must run after each retry delay and immediately before the next sender call. A test must block the first response, change `AutoInspectionEnabled` to false, release the response, and assert that the second attempt is not sent. The existing participation test only covers the next logical test (`inspection_test.go:540-625`), so the new runner test must cover the retry boundary.

### Request identity and native account behavior

The native request builder currently uses `codex-tui/0.160.0` and a four-argument `codexTurn` at `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/codex.go:16-27,130-179`. Upstream `7a49068:codex.go:16-20,131-185` changes the default identity to `codex-tui/0.163.0`, adds `history_initialization: "new"`, and inserts nonempty fingerprint instructions as the first `generic.developer_instructions` content part.

The native sender copies only Codex protocol headers at `/tmp/sub2api-settings-entry-final/source/backend/internal/service/account_codex_eval_service.go:143-150`. It then preserves OAuth identity normalization and account-specific User-Agent overrides at lines 163-175, applies account header overrides at 179, selects the configured account proxy at 183-187, and sends through the native plugin/HTTP transport. Keep all of these steps. The new default identity must not overwrite a stored account User-Agent or proxy.

### Focused checks

Use the existing owners and the approved bounded set:

- Add `tests/backend/internal/codexeval/model_test.go` with the native shared-runner retry cases. Cover one server/transport failure followed by success, one unusable stream followed by success, validator rejection followed by success, accumulated usage, exact three-attempt cap, and cancellation during 2-second backoff. Keep test delays injectable or use the existing upstream-style short timeout hook. Do not import CPA host helpers.
- Extend `tests/backend/internal/codexeval/runtime_test.go` to assert one Candy record for a wrong answer, `Attempts == 1`, and accumulated usage after a retryable failure. Extend stream semantics only for changed retry status behavior.
- Extend `tests/backend/internal/codexeval/inspection_test.go` to assert valid wrong Candy text makes one confirmation answer and to cover opt-out between retry attempts. Keep the existing confirmation and participation assertions.
- Replace old fingerprint statistics and ModelTrace parity tests with the approved detector golden/parsing tests in `tests/backend/internal/codexeval/fingerprint_test.go`. Use the upstream algorithm golden data only. Do not import CPA `TestAskCodex*`, `TestEvaluateCandy`, or host infrastructure wholesale.
- Extend `tests/backend/internal/service/account_codex_eval_service_test.go` only for sender status preservation and the retry marker. Reuse existing 401 and 429 tests, which already assert native account state updates. Keep `TestCodexEvalBackgroundNativeStreamLifetime` for active stream cancellation.
- Update `scripts/build.py` focused selector to include the new detector algorithm names and native model/inspection tests. Do not add race, coverage, full-suite, architecture, benchmark, or supplier tests.

## Files Found

- `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/model.go` — native status handling, SSE parsing, stream closure, and current status-based retries.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/codex.go` — current request builder, default identity, and one-shot `askCodex`.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/candy.go` — logical Candy result and inspection worker.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/fingerprint.go` — old fingerprint worker and `executeProbe` caller.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/modeltrace.go` — retired second `executeProbe` caller.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/runtime.go` — account admission, run contexts, cancellation, and worker dispatch.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/inspection.go` — participation callback, result collection, confirmation counting, and priority flow.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/service/account_codex_eval_service.go` — native account setup, status preservation, cooldown updates, auth, headers, and proxy.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/service/openai_plugin_transport.go` — plugin or native HTTP transport boundary.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/service/http_upstream_port.go` — native upstream interface; owner for the request-attempt marker type.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/repository/http_upstream.go` — request cancellation, body tracking, and transport errors.
- `/tmp/sub2api-settings-entry-final/source/backend/internal/service/plugin_runtime.go` — `PluginTransportError.RequestSent` classification.
- `patches/0001-evaluation-core.patch` — checked-in native evaluator source and current `executeProbe` contract.
- `patches/0004-admin-channel.patch` — checked-in native sender implementation.
- `patches/0005-inspection.patch` — checked-in automatic Candy inspection and participation flow.
- `tests/backend/internal/codexeval/runtime_test.go` — current logical result, cancellation, and history checks.
- `tests/backend/internal/codexeval/inspection_test.go` — current confirmation and opt-out checks.
- `tests/backend/internal/service/account_codex_eval_service_test.go` — current native sender, 401, and stream tests.
- `tests/backend/internal/service/codex_eval_lifetime_test.go` — native active-stream cancellation and background lifetime.
- Upstream `7a49068:internal/plugin/model.go`, `candy.go`, `fingerprint.go`, `pelican.go`, `codex.go`, and `model_test.go` — approved runner, result, request, and focused behavior.

## Code Patterns

- Native sender response status is preserved for non-200 results at `account_codex_eval_service.go:197-212`.
- HTTP upstream validation and client setup precede the actual `servertiming.Do` call at `http_upstream.go:203-239,242-297,299-310`; only the latter error is eligible for the typed attempted-request marker.
- `PluginTransportError.RequestSent` protects against replay, while plugin error payloads have no structured transient category at `plugin_runtime.go:298-335,357-387` and `plugin.proto:138-149`.
- Inspection starts account workers from the round context, then applies a separate 30-minute polling deadline per `testPhase` at `runtime.go:154-159` and `inspection.go:635-677`.
- Native response bodies cancel transport work when closed at `http_upstream.go:299-347`.
- Native inspection counts completed logical records at `inspection.go:613-633`.
- Native opt-out currently stops before the next logical Candy test at `runtime.go:183-192` and `inspection.go:574-603`.
- Upstream shared retry and validator handling is at `7a49068:model.go:47-98`.
- Upstream status and stream retry classification is at `7a49068:model.go:100-183`.
- Upstream attempt metrics are at `7a49068:candy.go:22-81`, `fingerprint.go:76-231`, and `pelican.go:21-81`.

## External References

- Upstream repository at commit `7a49068`: `https://github.com/haowang02/cpa-plugin-codex-candy-eval/tree/7a49068`.
- Native source snapshot: `/tmp/sub2api-settings-entry-final/source`.

## Caveats / Not Found

- No implementation, build, test, branch, merge, supplier request, or Git operation was performed.
- Plugin pre-response errors do not expose a structured transient category; this plan keeps them final and preserves the `RequestSent=true` no-replay contract. Stream failures after valid response headers remain eligible for the approved unusable-stream retry.
- The inspection's 30-minute phase deadline cancels request workers explicitly rather than being their direct context deadline; its post-cancel worker-drain loop has no separate timeout.
- Existing native source still contains old fingerprint and ModelTrace code. The report maps their deletion but does not change it.
- The existing participation test covers opt-out between logical tests, not during retry backoff. The new focused assertion remains required.
