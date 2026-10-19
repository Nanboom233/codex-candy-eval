# README Strict Review

## Findings (fixed)

- File: `README.md`
  - Issue: The rewrite said the evaluation sender bypasses model configuration restrictions. The current contract supports literal model handling, account-mapping bypass, and OAuth-normalization bypass.
  - Fix: Narrowed the sentence to the current contract. Changed installer automatic actions to bullets, corrected the stopped-container wording, and removed added container-path details while keeping the original fixture command.

## Findings (not fixed)

- None.

## Verification

- Lint: Pass. Reused the author's full-document style, command, and reference checks. The reviewer changes use short active sentences. `git diff --check -- README.md` passed.
- TypeCheck: Not applicable. This review changed documentation only.
- Tests: Not run. This review changed documentation only. Product checks and builds remain covered by the existing implementation and independent-review records.
- Semantics: Compared installation, request, inspection, restoration, build, and license claims with the current installer, build entry, and project contracts. The readme rewrite evidence records the retained commands and checked local references.
