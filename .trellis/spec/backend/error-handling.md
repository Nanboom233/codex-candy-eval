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

`executeProbe` retries the same request at most twice after the initial attempt. It retries status `0`, `429`, or server errors, not invalid answer content.

Preserve context cancellation, timeout causes, and stream closure.
`executeModel` closes the response on completion and uses `context.AfterFunc` to close it on cancellation.
Keep background task lifetime separate from the completed HTTP start request.
Explicit stop must cancel the task and supplier request.

Record asynchronous failures in their result or sample fields. Candy and ModelTrace truncate request errors through `truncate(err.Error(), 500)`.

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
- Do not retry malformed or incorrect answers as if they were transport failures.
- Do not hide save failures. Runners expose `storageError`. History removal handlers also restore the previous histories.
- Do not hold `mu` during supplier requests.
- Do not replace an upstream status with a generic success code before deciding whether a retry is allowed.
- Keep notification errors separate from successful priority changes.

Verification:

- [Runtime tests](../../../tests/backend/internal/codexeval/runtime_test.go): admission, history, cancellation, and stream outcomes.
- [Sender tests](../../../tests/backend/internal/service/account_codex_eval_service_test.go): native transport and conversion errors.
- [Inspection tests](../../../tests/backend/internal/codexeval/inspection_test.go): confirmations, cancellation, priorities, and notification failures.
