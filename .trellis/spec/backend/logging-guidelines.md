# Diagnostics and Logging Guidelines

> Report failures through existing error and state fields. This plugin has no dedicated logging framework.

---

## Overview

<!--
Document your project's logging conventions here.

Questions to answer:
- What logging library do you use?
- What are the log levels and when to use each?
- What should be logged?
- What should NOT be logged (PII, secrets)?
-->

The Go plugin returns errors to CPA and exposes evaluation and storage diagnostics to the management page. It does not configure `log`, `slog`, or a structured logging dependency.

Do not add stdout prints to a host-loaded shared library as a substitute for its response contract.

Evidence: [module dependencies](../../../go.mod), [host and management responses](../../../internal/plugin/plugin.go), and [result errors](../../../internal/plugin/candy.go).

---

## Diagnostic Categories

<!-- When to use each level: debug, info, warn, error -->

There is no plugin log-level policy. Use the existing category that owns the failure:

| Category | Existing channel | Owner |
| --- | --- | --- |
| Host dispatch or ABI callback failure | `EnvelopeError`, `errorEnvelope`, or returned Go error | `plugin.go`, `cabi.go` |
| Invalid management request | HTTP status and `error.message` | Management handlers and `jsonError` |
| Evaluation request failure | Result or sample `error` field | `candy.go`, `modeltrace.go` |
| State read or save failure | `storageError` and `/state` response `storage_error` | `state.go`, `stateResponse` |
| Incorrect model answer | Evaluation score, validity, or attribution outcome | Evaluation and algorithm files |

Do not describe an incorrect answer as a plugin crash. Transport success and answer correctness are separate outcomes.

References: [response helpers](../../../internal/plugin/plugin.go), [C ABI errors](../../../cmd/cpa-codex-candy-eval/cabi.go), and [trace outcomes](../../../internal/plugin/modeltrace.go).

---

## Structured Diagnostic Fields

<!-- Log format, required fields -->

Preserve the JSON field names consumed by the UI. The state endpoint returns `storage_error`, and the page renders it through `setNotice`.

Candy records include model, effort, time, duration, token counts, answer, the `ok` flag, and request error. These are result fields, not a separate event-log schema.

Use `truncate` for bounded user-visible request errors where adjacent runners do so. Candy and trace samples retain at most 500 runes of those errors.

This bound is not secret redaction. Do not treat a truncated host response as safe to publish.

Evidence: [candyResult and evaluateCandy](../../../internal/plugin/candy.go), [collectTraceSample](../../../internal/plugin/modeltrace.go), and [page notices](../../../internal/plugin/web/app.js).

---

## What to Report

<!-- Important events to log -->

- Preserve contextual error messages for failed host calls, response decoding, and interrupted streams.
- Expose state-file failures through `storageError`, including read failures that pause future saves.
- Preserve upstream HTTP status when retry decisions depend on it.
- Retain cancellation, incomplete response, and invalid-answer distinctions in the owning evaluation flow.
- Use test failure messages and synthetic host mocks when investigating local behavior.

Examples: `%w` context in [credentials](../../../internal/plugin/credentials.go), stream errors in [model execution](../../../internal/plugin/model.go), and visible save failures in [state tests](../../../internal/plugin/plugin_test.go).

---

## What NOT to Report

<!-- Sensitive data, PII, secrets -->

- Do not log or copy `.env` values, deployment passwords, management keys, API keys, or authorization headers.
- Do not publish raw credential configuration, token contents, account metadata, full model payloads, or private deployment paths.
- Use synthetic credentials in tests. Keep production account data out of screenshots, journals, and reports.
- Use existing masked key previews and stable credential IDs instead of adding secret fields to synchronization payloads.

The credential sync contract sends metadata, not API keys. `previewAPIKey` masks display names. It does not authorize disclosure of the original value.

Evidence: [private configuration rules](../../../AGENTS.md), [syncCredentialsResponse and previewAPIKey](../../../internal/plugin/credentials.go), and credential assertions in [UI checks](../../../scripts/test-ui.cjs).

Current errors can include upstream text. Check them for sensitive content before you share a report. This repository does not provide a universal redaction layer.
