# Final Task Audit

## Authorization

The user approved the final planning summary with `无误` on 2026-10-20.
The same message requested a reread to remove redundant tests and added requirements.
All task artifacts, review notes, manifests, and current test owners were reread before task activation.
The approved scope remains unchanged. This audit bounds verification and does not add product behavior.

## Scope Decisions

- Keep the four approved behavior groups in the PRD, including missing-variant compatibility and the migration copy.
- Keep production work in patches `0001` and `0006`, with necessary adjacent tests and documentation.
- Use the upstream version-based recovery filename. Remove the invented unique-name scheme and its collision-preservation test requirement.
- Do not add an XML parser, stricter output validation, automatic image scoring, new history filters, or periodic backups.
- Keep current transport, retries, inspection, installer transactions, and release workflow.
- Keep native Sub2API version 0.2.17. The upstream CPA release number is not a native host version change.

## Verification Ownership

| Owner | Necessary additions or updates | Existing evidence to reuse |
| --- | --- | --- |
| Existing Pelican extraction test | Selected-language preference, SVG page preparation, and the unchanged code-size boundary | Retain current HTML extraction cases. No separate parser suite. |
| `TestCodexEvalRuntime` | Explicit variant routing, unknown-variant admission, saved version fields, schema-1 retention, and original-file copy outcomes | Existing omitted-variant starts, cancellation, invalid/null protection, history limits, and failed-clear rollback. Add version assertions there. |
| Existing UI harness | Variant configuration, startup preference preservation, prompt/copy updates, start payloads, and record-derived labels | Existing account targeting, refresh, Settings startup, model dropdown, and selection checks. Extend their fixtures where relevant. |
| Standard build/check entry | Run once for the final candidate and keep new tests in its focused selector | Existing shared retry, transport, inspection, vet, frontend, and installer checks run through that entry. Do not run them separately for reassurance. |
| Actual browser | Inspect the changed toolbar and version metadata at normal/narrow widths, both preview types, and detail animation/close behavior | DOM assertions establish preferences, copy text, payloads, and labels. Do not repeat that assertion list manually. |

The new backup failure path needs a focused assertion because it blocks later writes.
Reuse the current save-failure path and fixture to check migration-save failure. Do not create another storage harness.
Keep reload assertions with migration coverage. Do not add a restart matrix or re-prove original file protection in a second suite.
Reuse `TestCodexEvalRetries` for request reuse, attempt limits, waits, usage accumulation, final HTTP failures, and backoff cancellation.
Do not duplicate these cases for each Pelican version.
Use a small output sequence only when needed to connect the new variant validator to the existing retry owner.

## Excluded Gates

Do not add race, coverage, architecture, benchmark, source-resource parity, digest, real-model, or production-installation gates.
Do not run a pristine baseline, a second unchanged build, or a duplicate browser regression checklist.
The reviewer reuses valid passing evidence. Rerun only checks affected by a repair or concrete failure.
The standard successful build establishes both fixed executable outputs. No separate packaging test is required.

## Document Authority

`prd.md` defines behavior and acceptance. `design.md` defines the native adaptation.
`implement.md` and this audit define execution and verification ownership.
`plan-review.md` preserves approved discussion and source evidence. Its earlier evidence lists do not require separate tests or runs.
Existing specs describe the current product. The task's explicit approved changes supersede their older schema-1 and no-backup descriptions.
