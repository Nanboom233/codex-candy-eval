# Stage 4B Run-State Quality Review

- Date: 2026-10-12
- Scope: `internal/eval/runs.go`, `state.go`, `runs_test.go`, and `state_test.go`.
- Verdict: **PASS** for the delivered Stage 4B manager and file-state boundary.

## Findings (fixed)

- None. No Stage 4B source or test changes were needed.

## Findings (not fixed)

- File: `internal/eval/state.go`
- Issue: No test injects a failure after the temporary file is created and then
  asserts that the owned temporary file is removed. Existing tests verify
  cleanup after successful rename, atomic replacement, and preservation of the
  destination for startup and recovery failures. The failure path is implemented
  by `defer os.Remove(tmp)` immediately after `CreateTemp`.
- Why not fixed: The current `StateStore` interface does not expose filesystem
  operations, and a deterministic post-create failure fixture would require a
  new filesystem seam or a timing-dependent filesystem race. Neither is needed
  for this local quality gate; the cleanup ownership is explicit in the source.
- Remaining acceptance: Docker lifecycle, authenticated browser integration,
  deployment, and live billing were not run. The task records no approved
  target and an unavailable Docker daemon. These are later migration gates, not
  local Stage 4B failures.

## Verification

- Lint: pass. `gofmt -l` returned no files; no separate linter is configured.
- TypeCheck: pass. `go vet ./...` and copied-package `go vet` passed.
- Tests: pass. `go test -race -count=1 ./internal/eval`, root
  `go test -race -count=1 ./...`, and copied-package
  `go test -race -count=1 ./internal/codexeval` passed.
- Formatting, whitespace, and source parity: changed Go files are formatted;
  `git diff --no-index --check /dev/null <file>` emitted no whitespace
  diagnostics for the four new Go files; `diff -rq internal/eval
  dist/sub2api-build/workspace/backend/internal/codexeval` reported no
  differences. The no-index diff returns status 1 because the files are
  untracked.
- Frozen patch hashes match: 001
  `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`, 002
  `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`, and
  003 `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.

The review verified save-before-worker admission, rollback after failed startup
saves, request-context independence after commit, cross-group/type account
exclusion, unavailable-target identity, recovery without paid calls, protected
load/recovery failures, explicit flush, clear rollback, retention and deep
copies, all-type cancellation with observed usage, and deadline-aware shutdown
while startup, terminal, or unavailable-only saves are pending. State locks are
released during storage I/O, while `saveMu` serializes durable transitions.
