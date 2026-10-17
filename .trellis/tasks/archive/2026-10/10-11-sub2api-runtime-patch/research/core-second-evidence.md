# Research: Core C08-C12 Evidence

- Query: Compare the evaluation-core additions for transport limits, observed usage, retries, Codex wire construction, and Candy result semantics with the original plugin and the pristine Sub2API v0.2.16 native runtime.
- Scope: internal
- Date: 2026-10-15

## Findings

### C08: transport and stream limits

The current evaluation package adds fixed parser limits in `internal/eval/http.go:12-19`: 1 MiB request, 16 MiB cumulative stream, 1 MiB event and line, 65,536 events, and 4 MiB retained output. `internal/eval/sse.go:52-67,96-99,134-145,152-161` enforces these limits while reading the Responses SSE stream. The original `internal/plugin/model.go`/stream reader has no equivalent set of request, line, event, output, or cumulative-stream constants.

The pristine native server already has global HTTP boundaries. `backend/internal/config/config.go:688-700` defines `ServerConfig.MaxHeaderBytes` and `MaxRequestBodySize`; defaults at `:2047-2050` are 64 KiB headers and 256 MiB request body. `backend/internal/server/http.go:115-133` applies `MaxHeaderBytes` to the server and wraps the entire Gin router with `http.MaxBytesHandler`, so the administrator route is covered as well. Application-specific limits remain separate: for example, the evaluation Pelican HTML cap is `internal/eval/pelican.go:12,43-49`.

Minimal recommendation: remove the evaluation-specific transport/parser caps when the request is sent through native Sub2API HTTP handling. Keep native global HTTP limits and preserve algorithm-specific limits such as the Pelican HTML cap. Do not silently remove application bounds while reviewing C08.

### C09: observed usage

`internal/eval/sse.go:20-31` decodes usage fields as pointers. `mergeUsage` at `:34-50` accepts only nonnegative observed values, writes explicit zero values, and marks `UsageObserved`. `internal/eval/model.go:99-112` sums numeric usage for every probe attempt and ORs `UsageObserved` across attempts, including failed attempts.

Production consumers use the numeric fields directly: Fingerprint initializes and sums them at `internal/eval/fingerprint.go:105-106,141-144`; ModelTrace initializes and sums them at `internal/eval/modeltrace.go:115,155-159`; Candy and Pelican copy them into result records at `internal/eval/candy.go:30-38` and `internal/eval/pelican.go:32-40`. State cloning only handles Fingerprint pointer result fields at `internal/eval/state.go:296-299`. A repository-wide search finds no production consumer of `UsageObserved` beyond its parser assignment and retry OR; remaining references are tests.

Minimal recommendation: retain pointer decoding and nonnegative presence filtering so missing usage and explicit zero continue to have stable parser semantics, retain numeric summation across failed attempts, and remove the unused `UsageObserved` result flag plus its retry OR. This does not change serialized numeric token values or native Sub2API billing, which is outside this evaluation result type.

### C10: retry, timeout, and request ownership

The original plugin retry loop at `internal/plugin/model.go:29-85` allows three attempts, retries status 0, 429, and status >=500, waits two seconds then four seconds, checks the parent context before admission and before each attempt, releases the slot before waiting, and derives probe timeout from a three-minute child context (`:64`). Answer requests use the ten-minute timeout constant at `:35-38`.

The current evaluation loop keeps those counts, waits, and context checks at `internal/eval/model.go:80-133`, but adds an exported `ExecuteProbe` wrapper (`:70-72`) and clones body/header data before the loop (`:78-79`) and again for every attempt (`:102-105`). Production code calls only the unexported helper from `internal/eval/fingerprint.go:75` and `internal/eval/modeltrace.go:95`; `ExecuteProbe` is found only in its declaration and tests.

The current HTTP executor does not mutate the caller-owned `ResponsesRequest`: `internal/eval/http.go:52-69` creates a new `http.Request` with `bytes.NewReader(req.Body)`, copies the allowlisted values into the request's newly allocated header map, and sets content negotiation headers there. The native Sub2API account test follows the same ownership boundary: pristine `backend/internal/service/account_test_service.go:902-946` creates a fresh request from `payloadBytes` and then sets headers; pristine `backend/internal/service/openai_plugin_transport.go:23-44` passes that request to the selected transport. No evidence shows the native sender mutating a shared evaluation body or header map.

Minimal recommendation: preserve three attempts, 2/4-second waits, retry status semantics selected by the approved C10 contract, slot release before backoff, and parent cancellation. Remove the unused exported wrapper. Once the native sender remains read-only with respect to the request value, remove duplicate per-attempt `bytes.Clone`/`Header.Clone`; keep a copy only at a concrete mutation boundary if a future sender actually mutates input.

### C11: Codex request wire

`internal/eval/codex.go:16-27,36-55` embeds the captured Codex tools, base instructions, messages, environment context, and the original `codex-tui/0.160.0` user agent. `codexTurn` at `:130-179` creates fresh thread/turn IDs, preserves the original header set and metadata, and derives a stable installation ID from the account namespace supplied by the caller. `internal/eval/eval.go:276-277` supplies `sub2api-account:<AccountID>` as that namespace. Preserve this wire and namespace behavior; do not alter the captured architecture string in the user agent as part of C11.

### C12: Candy behavior

The current Candy implementation at `internal/eval/candy.go:8-54` matches the original algorithm and prompt in `internal/plugin/candy.go:100-123`: one Codex request, answer truncated to 4000 runes, standalone `21` detection, duration and numeric usage copied into the result. A transport/cancellation error sets the result status and error; a successful but incorrect answer remains a completed result with `OK=false`. Preserve this distinction and do not turn a wrong answer into a retry or transport failure.

## Files Found

- `internal/eval/http.go` — current evaluation executor and added transport limits.
- `internal/eval/sse.go` — current SSE parser, usage decoding, and limit enforcement.
- `internal/eval/model.go` — current retry, timeout, clone, panic recovery, and executor wrapper.
- `internal/eval/fingerprint.go` — production numeric usage consumer and probe caller.
- `internal/eval/modeltrace.go` — production numeric usage consumer and probe caller.
- `internal/eval/candy.go` / `internal/eval/pelican.go` — single-request result consumers.
- `internal/eval/state.go` — persisted result cloning.
- `internal/eval/codex.go` / `internal/eval/eval.go` — Codex wire and account namespace.
- `internal/plugin/model.go` / `internal/plugin/candy.go` — original plugin retry and Candy behavior.
- pristine `backend/internal/config/config.go` / `backend/internal/server/http.go` — native global HTTP limits.
- pristine `backend/internal/service/account_test_service.go` / `openai_plugin_transport.go` — native request construction and transport handoff.

## Related Specs

- `.trellis/spec/backend/evaluation-core.md` — evaluation-core contracts and historical limit/retry signatures; current user approval governs which additions remain.
- `.trellis/spec/backend/quality-guidelines.md` — native Sub2API scope and removed race acceptance gate.

## Caveats / Not Found

- No product code, build, test, deployment, or live provider request was run for this evidence.
- The pristine native implementation was inspected at the pinned v0.2.16 commit supplied by the active task; transport internals may still evolve after this planning review.
- `UsageObserved` remains present in the current test fixtures; removing it later requires updating those tests and any copied-package manifests consistently.
