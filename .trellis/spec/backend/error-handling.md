# Error Handling

> Keep host-call errors, management errors, and evaluation outcomes separate.

---

## Overview

<!--
Document your project's error handling conventions here.

Questions to answer:
- What error types do you define?
- How are errors propagated?
- How are errors logged?
- How are errors returned to clients?
-->

The plugin has two response boundaries. Host method dispatch returns an `Envelope`. Management handlers return a `managementResponse` with an HTTP status and JSON body.

An evaluation can finish with an incorrect answer without a transport error. Do not convert an incorrect model answer into a host or management failure.

Evidence: [HandleMethod and response helpers](../../../internal/plugin/plugin.go), [evaluateCandy](../../../internal/plugin/candy.go), and `TestEvaluate` in [plugin tests](../../../internal/plugin/plugin_test.go).

---

## Error Types

<!-- Custom error classes/types -->

`EnvelopeError` carries `code`, `message`, and optional `http_status`. Its `Error` method preserves the host code in the error string.

Use `errors.As` when callers need the upstream status. `executeModel` uses it to distinguish an `EnvelopeError` from failures without an upstream response.

Wrap contextual Go errors with `%w` when callers must retain the underlying error. `credentials` and `loadState` use this pattern.

Use `errors.Is` for cancellation and filesystem conditions. Do not infer cancellation solely from an error-message string.

Evidence: [host envelope type](../../../internal/plugin/plugin.go), [model execution](../../../internal/plugin/model.go), [credential loading](../../../internal/plugin/credentials.go), and [state loading](../../../internal/plugin/state.go).

---

## Error Handling Patterns

<!-- Try-catch patterns, error propagation -->

Validate JSON and request fields at the management handler before starting workers. Reuse `validModelName`, `startRuns`, and existing credential selection helpers.

Use narrow recovery boundaries. `HandleMethod` reports `plugin_panic`. `executeModel` converts model-request panics into errors. Do not add recovery to every helper.

`executeProbe` retries the same request at most twice after the initial attempt. It retries status `0`, `429`, or server errors, not invalid answer content.

Preserve context cancellation, timeout causes, and stream closure. `executeModel` closes the host stream on completion or cancellation through `sync.OnceFunc`.

Record asynchronous failures in their result or sample fields. Candy and ModelTrace truncate request errors through `truncate(err.Error(), 500)`.

Evidence: [request validation and result errors](../../../internal/plugin/candy.go), [probe retries and stream cleanup](../../../internal/plugin/model.go), and [trace sample errors](../../../internal/plugin/modeltrace.go).

---

## API Error Responses

<!-- Standard error response format -->

Use `jsonError` for management failures. Do not return a flat string or a host envelope as the management body.

Example management response body:

```json
{"error":{"message":"Request failed"}}
```

Use `errorEnvelope` at the host boundary. A failed management operation can still sit inside a successful host-dispatch envelope.

| Condition | Existing HTTP status | Evidence |
| --- | --- | --- |
| Malformed JSON or invalid fields | 400 | `candyRunResponse`, `syncCredentialsResponse` |
| Unknown route or missing record | 404 | `handleManagement`, `recordResponse` |
| Host credential lookup failure | 502 | `startRuns`, `stateResponse` |
| Plugin quiescing | 503 | `startRuns` |
| Failed history removal save | 500 | `clearHistoryResponse` |

The browser `api` helper consumes `error.message`. Keep this shape aligned with [app.js](../../../internal/plugin/web/app.js).

References: [management response helpers](../../../internal/plugin/plugin.go), [run scheduling](../../../internal/plugin/runs.go), and [history handlers](../../../internal/plugin/state.go).

---

## Common Mistakes

<!-- Error handling mistakes your team has made -->

- Do not grade a failed model request, even when the host error text contains the expected candy answer.
- Do not retry malformed or incorrect answers as if they were transport failures.
- Do not hide save failures. Runners expose `storageError`. History removal handlers also restore the previous histories.
- Do not hold `mu` during host calls. `stateResponse` performs host lookups before locking the histories.
- Do not replace an upstream status with a generic success code before deciding whether a retry is allowed.

Verification:

- [Plugin tests](../../../internal/plugin/plugin_test.go): `TestEvaluate`, `TestStatePersistence`, and `TestRegisteredRoutesAreHandled`.
- [Model tests](../../../internal/plugin/model_test.go): probe retries and stream failures.
