# Error Handling

> Keep native API errors, upstream failures, and evaluation outcomes separate.

---

## Overview

Administrator handlers use native Sub2API response helpers.
The direct account sender returns the actual upstream status and response stream.
Background runners record request errors in their result or sample fields.

An evaluation can finish with an incorrect answer without a transport error.
Do not convert an incorrect model answer into an API or transport failure.

Evidence: [administrator channel](../../../patches/0004-admin-channel.patch) and [evaluation core](../../../patches/0001-evaluation-core.patch).

---

## Error Types

`executeModel` returns the result, observed HTTP status, and error separately.
Keep the upstream status available for retry decisions.
Wrap contextual errors with `%w` when callers must retain the underlying error.

Use `errors.Is` for cancellation and filesystem conditions. Do not infer cancellation solely from an error-message string.

Evidence: model execution and history persistence in the [evaluation core](../../../patches/0001-evaluation-core.patch).

---

## Error Handling Patterns

Validate JSON in the administrator handler and evaluation fields in `Start` before workers start.
Use the existing `validModelName` and native account admission paths.

`executeModel` converts model-request panics into errors. Do not add recovery to every helper.

`askCodex` builds one request and sends at most three attempts. Wait two seconds, then four seconds, between attempts.
Accumulate observed usage and attempts. Measure duration across the logical request, including waits.
Each fingerprint attempt has a three-minute limit. Each Candy or Pelican attempt has a ten-minute limit.

`HTTPUpstreamRequestError` wraps only errors from the actual native `servertiming.Do` transport call.
It provides `Error`, `Unwrap`, and `RequestAttempted() bool` beside `HTTPUpstream` in the service package.
Evaluation recognizes the marker through an interface and `errors.As`, without a service dependency.
Keep host, proxy, client, account, credential, and conversion setup errors unmarked and final.
Pre-response plugin errors remain final. Possible sending does not establish transience or replay permission.
An attempt deadline alone cannot override that plugin boundary.

Capture response status and arrange body closure before processing a sender error.
Return the observed response with a non-200 body-read error. Do not replace its status with zero.
Retry marked direct transport failures, HTTP 408/5xx, and unusable successful streams while the parent context remains active.
Keep HTTP 401, 429, and other permanent statuses final, including after body-read failure.
Retry empty, incomplete, failed, or tool-only successful streams.
Fingerprint retries numeric-count rejection. Sequence-shape rejection does not retry.
Pelican retries failed selected-language markup acceptance through the same runner. Valid wrong Candy text completes without a retry.
Do not add retries around account updates, settings saves, or notifications.

Preserve context cancellation, timeout causes, and stream closure.
`executeModel` closes the response on completion and uses `context.AfterFunc` to close it on cancellation.
Keep background task lifetime separate from the completed HTTP start request.
Explicit stop must cancel the task and supplier request.

Record asynchronous failures in their result or sample fields. Keep existing error truncation.
Cancellation stops active streams, retry waits, and later sends. Release each response before the next wait or attempt.
Inspection checks participation before every attempt, outside the runtime mutex.
An intentional opt-out does not create a failed result or cancel unrelated accounts.

Evidence: [evaluation core](../../../patches/0001-evaluation-core.patch) and [native lifetime test](../../../tests/backend/internal/service/codex_eval_lifetime_test.go).

---

## API Error Responses

Use native `response.BadRequest`, `response.ErrorFrom`, and `response.Error` at administrator boundaries.
Successful task operations use native `{code,message,data}` envelopes.
Existing record and clear helpers can return `error.message` bodies through the administrator handler.
The page accepts these bodies and native `message` fields.
The direct account endpoint forwards upstream HTTP and SSE outcomes.

| Condition | Existing HTTP status | Evidence |
| --- | --- | --- |
| Malformed JSON or invalid fields | 400 | Administrator handlers and `Start` |
| Missing history record | 404 | `recordResponse` |
| Unclassified native sender or model discovery failure | 502 | `CodexEval`, `CodexEvalModels` |
| Failed history clear save | 500 | `clearHistoryResponse` |
| Invalid inspection settings | 400 | Native inspection settings handler |
| Inspection unavailable or already running | 409 | Native inspection start handler |
| Settings repository failure | 500 | Native inspection settings handler |

References: [evaluation UI](../../../patches/0006-evaluation-ui.patch), [administrator channel](../../../patches/0004-admin-channel.patch), and [inspection](../../../patches/0005-inspection.patch).

---

## Common Mistakes

- Do not grade a failed model request, even when its error text contains the expected Candy answer.
- Use each evaluation's acceptance rule. A valid wrong Candy answer is complete. Detector count rejection can retry.
- Do not hide save failures. Runners expose `storageError`. History removal handlers also restore the previous histories.
- A failed migration copy blocks replacement and later saves. Keep decoded records readable in memory.
- Do not hold `mu` during supplier requests.
- Do not replace an upstream status with a generic success code before deciding whether a retry is allowed.
- Keep notification errors separate from successful priority changes.

Verification:

- [Runtime tests](../../../tests/backend/internal/codexeval/runtime_test.go): admission, history, cancellation, and stream outcomes.
- [Sender tests](../../../tests/backend/internal/service/account_codex_eval_service_test.go): native transport and conversion errors.
- [Inspection tests](../../../tests/backend/internal/codexeval/inspection_test.go): confirmations, cancellation, priorities, and notification failures.
- [Shared runner](../../../tests/backend/internal/codexeval/model_test.go): retry limits, request reuse, accumulated usage, validators, final errors, and cancellation.
