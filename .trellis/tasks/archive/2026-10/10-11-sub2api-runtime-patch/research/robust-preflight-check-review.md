# Independent Check: Robust Preflight

Date: 2026-10-13

## Findings (fixed)

- `install.py`: Maintenance parsing bounded response characters instead of
  UTF-8 bytes and called `.get()` before checking the decoded JSON shape.
  Legal non-object JSON could escape as an uncaught `AttributeError`. It also
  accepted truthy non-boolean fields and did not require successful drain and
  owned-Key revocation on deactivate. The parser now bounds encoded response
  bytes, requires an object, validates the integer schema, strict `ok`, owned
  deployment ID, supported state, and boolean cleanup fields. Deactivate
  requires both cleanup confirmations to be exactly `true`.
- `tests/installer/test_preflight.py`: Added regressions for malformed JSON
  shapes, strict field types, unsupported state, UTF-8 byte overflow, and
  incomplete or incorrectly typed deactivate confirmations. Valid status and
  deactivation responses remain accepted, including unrelated provenance.
  Corrected the existing no-binary-argument assertion to inspect each CLI
  argument string instead of comparing tuple members for equality.

## Findings (not fixed)

- None in the reviewed product or active specification scope. The main session
  updated the Active Sub2API Delivery section to name the structural adapter,
  functional probes/builds and removed binary gates. It also synchronized the
  task artifacts and added the strict maintenance parser contract and tests to
  `patch-native-installation.md`. I clarified that the installer checks byte
  length before JSON parsing, after its Docker helper decodes UTF-8 output.

## Remaining Verification Boundaries

- The broad upstream native race, frontend, and previous integration gates
  remain documented historical failures. They were not rerun or changed in
  this installer-only parser review. Focused native checks and the 18-point
  disposable arm64 lifecycle are implementation-stage evidence, not new parser
  review executions.
- Published bootstrap, live deployment, production billing/concurrency, and
  deployed browser acceptance remain unverified. No CPA deployment or paid
  request was run.

## Verification

- Lint: No dedicated Python lint target/configuration or linter executable is
  available. `git diff --check`, shell syntax, and UI checks pass.
- TypeCheck: Go compilation and vet pass in the repository check. Python
  installer and test modules import successfully in both source and archive
  preflight runs.
- Tests: Pass. `scripts/test-sub2api.sh check` completed with 72 helper tests
  (1 historical skip), 48 transaction tests, 8 universal tests, 16 preflight
  tests, Go vet/race, shell syntax, and UI checks.
- Evidence provenance: Independent review created fresh current and previous
  adapted workspaces at `/tmp/codex-robust-independent-review-20261013`.
  Their `production-byte-parity.json` files were generated at 18:03 after
  `scripts/universal_dev.py prepare` completed structural adaptation and
  overlay byte parity. This is fresh two-source preparation evidence.
  The retained native type, focused race, and DTO/repository test logs under
  `/tmp/codex-robust-preflight-20261013` are timestamped 17:23–17:30, before
  that independent preparation, so they are implementation-stage test
  evidence, not fresh native gates from independent review.
  The parser-only review then reran the root check and archived preflight; it
  did not rerun two-source native type or race gates.
- Archive: Repacked functional-only v0.4.0 package. Installer extraction and
  `validate_bundle` pass for all 62 payload files. Extracted bytes match fresh
  staging byte-for-byte. The archive contains 62 files and 22 directory/root
  entries. The archived installer passes all 16 preflight tests.
- Functional source identity: `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759`
  (unchanged).
- Archive identity: 525,565 bytes, SHA-256
  `f9980775d4aefe08906b9c0a480420f5454eb5e03164780b98643c4535176cfd` at
  `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`.
- Temporary evidence remains under `/tmp/codex-robust-preflight-20261013`
  and `/tmp/codex-robust-independent-review-20261013`. No cleanup was run.
