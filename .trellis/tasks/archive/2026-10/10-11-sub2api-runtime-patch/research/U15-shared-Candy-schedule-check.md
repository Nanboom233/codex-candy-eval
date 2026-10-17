# U15 Shared Candy Account Selection Check

## Findings fixed

- `native/backend/internal/handler/admin/codex_inspection_handler.go`
  - The implementation now uses the native account-list conversion for both
    scheduled-list and single-account recheck paths.
  - Shadow accounts use the current parent plan when their own plan is empty.
    An account plan on the shadow account remains authoritative.
  - The regression test covers batch parent lookup, own-plan precedence, and a
    second lookup after the parent plan changes.

No further product edits were required during this independent check.

## Findings not fixed

- No product finding remains in the reviewed U15 paths.
- Docker inventory and browser/build checks were not run by this check. The
  task contract requires one retained `sub2api-codex-eval-test` container; this
  check created no container and changed no Docker resource. The main session
  must use the existing container evidence or perform the authorized inventory
  check if the final record needs a current runtime readback.
- Build and package checks were not repeated. U15 source changes are ready for
  the main session's delivery step.

## Verification

- Go tests: pass
  - `go test -tags=unit -overlay /tmp/sub2api-U15-go-overlay.json ./internal/codexeval ./internal/service ./internal/handler/admin -run 'TestCodexInspection|TestCodexEvalBackgroundNativeStreamLifetime' -count=1`
  - This includes the native lifetime test and the shadow-parent regression.
- Go vet: pass
  - `go vet -tags=unit -overlay /tmp/sub2api-U15-go-overlay.json ./internal/codexeval ./internal/service ./internal/handler/admin`
- UI checks: pass
  - `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`
  - Covers autosave ordering, shared selection, missing-group preservation,
    disabled-account model discovery, manual entry points, and Candy history.
- JavaScript syntax: pass (`inspection.js`, `catalog.js`).
- Go formatting: pass (`gofmt -d` on changed Go files).
- Diff whitespace: pass (`git diff --check`).

No race, coverage, full-suite, browser, Docker, production, or paid-provider
gate was added.
