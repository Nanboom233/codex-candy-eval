# Stage 4B Run-State Specification Review

- Date: 2026-10-12
- Scope: `internal/eval/runs.go`, `state.go`, `runs_test.go`, and `state_test.go`.
- Verdict: **PASS** for the delivered Stage 4B manager and file-state boundary.

## Findings (fixed)

- None. No Stage 4B source or test changes were needed.

## Findings (not fixed)

- File: `.trellis/spec/backend/run-state.md`
- Issue: The status sentence says both independent reviews remain pending. This review passes the specification gate, while the quality review remains pending.
- Why not fixed: The assigned edit ownership covers the four Stage 4B source and test files plus this report. The main session should update the status after recording this verdict.
- Remaining acceptance: Docker lifecycle, authenticated browser integration, deployment, and live billing were not run. The task records no approved target and an unavailable Docker daemon. These are later migration gates, not local Stage 4B failures.

## Verification

- Formatting and whitespace: `gofmt -l` returned no files; `git diff --check` passed.
- Focused race: `go test -race -count=1 ./internal/eval` passed.
- Focused stress: `go test -race -count=20 ./internal/eval` passed in 155.985 seconds.
- Root static check: `go vet ./...` passed.
- Root race suite: `go test -race -count=1 ./...` passed for the C ABI command package, `internal/eval`, and the existing `internal/plugin` package.
- Pinned copied-package checks: `go vet ./internal/codexeval` and `go test -race -count=1 ./internal/codexeval` passed in the staged Sub2API backend.
- Source parity: `diff -rq internal/eval dist/sub2api-build/workspace/backend/internal/codexeval` reported no differences.
- Frozen patch hashes match the recorded values: patch 001 `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`, patch 002 `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`, and patch 003 `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.
- No separate lint command is configured. `go vet` passed for the root and copied package.

The source review verified that startup saves precede reservations and workers; failed saves preserve prior state; request cancellation after commit does not own worker lifetime; account reservations span groups and evaluation types; unavailable and recovered records preserve requested identity; history is capped and returned records are copied; protected load and recovery errors cannot be flushed over; terminal save failures retain in-memory results and block starts; cancellation covers all four types; and shutdown remains deadline-aware through blocked persistence and unavailable-only commits.
