# Independent U11 Settings Check

## Findings (fixed)

- File: `native/backend/internal/service/codex_inspection_settings_test.go`
- Issue: A previously valid stored pair with normal priority 0 and degraded priority 1–9 becomes invalid after the required read-only normalization to normal priority 10. The old degraded value must not be rewritten.
- Fix: Preserve the degraded value and test that the read does not write. Validate loaded settings before the coordinator configures or schedules them. An invalid load keeps a visible error and blocks starts until an explicit valid save. This prevents an expired schedule from being retried in a tight loop.

- File: `native/backend/internal/codexeval/web/inspection.js`
- Issue: The page rejected an invalid priority pair in the browser even when automatic priority adjustment was off and both fields were hidden. This prevented the API error from reaching the visible save status.
- Fix: Check the priority pair in the page only while automatic adjustment is on. With the fields hidden, the save reaches the API and displays its validation error. Turning adjustment on reveals the fields and enables the client-side check.

- Files: `native/backend/internal/codexeval/inspection_test.go`, `native/backend/internal/service/codex_inspection_settings_test.go`, `tests/ui/sub2api-dropdowns.cjs`
- Issue: No focused regression covered the invalid legacy priority pair across settings load, coordinator startup, and the hidden-field save path.
- Fix: Added assertions for read-only normalization, preservation of the stored degraded value, blocked scheduling/start on invalid loaded settings, visible API validation feedback, and correction through the visible priority fields.

## Findings (not fixed)

No U11 code finding remains. The main session owns the post-rebuild browser recheck in the retained container. The previously reported parent-page 304/CSP nonce issue is outside U11 and remains unchanged.

The active-delivery section in `.trellis/spec/backend/quality-guidelines.md` still describes the retired universal adapter pipeline. It conflicts with the current backend index and the R04 task decisions. I left it for the main session's spec synchronization.

## Verification

- Lint: No dedicated linter is configured for the embedded standalone scripts. `node --check` passed for the changed scripts.
- TypeCheck: Passed through the native frontend build (`vue-tsc -b`).
- Tests: Passed. `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs` passed. `scripts/test-sub2api.sh check` passed the focused Go tests, touched-package `go vet`, frontend locale tests/build, and static Linux/amd64 build.
- Package: `python3 scripts/package-sub2api.py` passed. The binary is a stripped static Linux/x86-64 ELF. The archive contains the binary, installer, and five license files.
- Source formatting and whitespace: `gofmt` and `git diff --check` passed.
