# Diagnostics and Logging Guidelines

> Report failures through existing native API, result, and state fields.

---

## Overview

The native module exposes evaluation, inspection, and storage diagnostics to the administrator page.
It does not add a separate logging framework or event store.
Use the response or state owner for each failure.

Evidence: [administrator channel](../../../patches/0004-admin-channel.patch), [evaluation core](../../../patches/0001-evaluation-core.patch), and [inspection](../../../patches/0005-inspection.patch).

---

## Diagnostic Categories

Use the existing category that owns the failure:

| Category | Existing channel | Owner |
| --- | --- | --- |
| Invalid administrator request | Native HTTP status and response message | Administrator handlers |
| Supplier HTTP or stream failure | Actual upstream status and error | Native account sender and `executeModel` |
| Evaluation request failure | Result or sample `error` field | `candy.go`, `modeltrace.go` |
| State read or save failure | `storageError` and `/state` response `storage_error` | `state.go`, `State` |
| Inspection or notification failure | State `error`, account result `error`, or `notification_error` | Inspection coordinator and page |
| Settings persistence failure | Native response message and form error | Native settings service and page |
| Incorrect model answer | Evaluation score, validity, or attribution outcome | Evaluation and algorithm files |

Do not describe an incorrect answer as an application crash.
Transport success and answer correctness are separate outcomes.

References: [administrator channel](../../../patches/0004-admin-channel.patch) and [evaluation UI](../../../patches/0006-evaluation-ui.patch).

---

## Structured Diagnostic Fields

Preserve the JSON field names consumed by the UI. The state endpoint returns `storage_error`, and the page renders it through `setNotice`.

Candy records include model, effort, time, duration, token counts, answer, the `ok` flag, and request error. These are result fields, not a separate event-log schema.

Use `truncate` for bounded user-visible request errors where adjacent runners do so. Candy and trace samples retain at most 500 runes of those errors.

This bound is not secret redaction. Do not treat truncated supplier text as safe to publish.

Evidence: results and trace samples in the [evaluation core](../../../patches/0001-evaluation-core.patch), and notices in the [evaluation UI](../../../patches/0006-evaluation-ui.patch).

---

## What to Report

- Preserve contextual errors for native account lookup, supplier requests, response decoding, and interrupted streams.
- Expose state-file failures through `storageError`, including read failures that pause future saves.
- Preserve upstream HTTP status when retry decisions depend on it.
- Retain cancellation, incomplete response, and invalid-answer distinctions in the owning evaluation flow.
- Use test failure messages and synthetic senders when investigating local behavior.

Examples: stream errors in the [evaluation core](../../../patches/0001-evaluation-core.patch) and save failures in [runtime tests](../../../tests/backend/internal/codexeval/runtime_test.go).

---

## What NOT to Report

- Do not log or copy `.env` values, deployment passwords, management keys, API keys, or authorization headers.
- Do not publish raw credential configuration, token contents, account metadata, full model payloads, or private deployment paths.
- Use synthetic credentials in tests. Keep production account data out of screenshots, journals, and reports.
- Use stable native account IDs and existing masked display fields. Do not add secrets to page payloads.

Inspection settings reads omit notification secrets and return configured flags.
Empty secret updates retain saved values. Do not expose saved values in screenshots or reports.

Evidence: [private configuration rules](../../../AGENTS.md), [native inspection](../../../patches/0005-inspection.patch), and [UI checks](../../../tests/ui/sub2api-dropdowns.cjs).

Current errors can include upstream text. Check them for sensitive content before you share a report. This repository does not provide a universal redaction layer.
