# Fingerprint Stop Fixture Review

## Verdict

**PASS.** The reopened Stage 4A Fingerprint stop gate is closed for the maintained evaluation core. The correction tests the source-preserving threshold contract: cancel when the aggregator handles the eighth consecutive error, drain admitted work, retain returned usage, and publish final progress only after worker settlement.

This verdict covers the maintained `internal/eval` package and its byte-identical staged copy. It does not change the separate pristine upstream service baseline or deployment gates.

## Findings (fixed)

- `internal/eval/eval_test.go:428-445` now bounds the concurrency-one immediate-error fixture to eight or nine executor calls and counted errors. It requires the threshold failure, `Done == Errors`, no valid samples, exact returned input usage, and final progress equality. This replaces the unsupported exact-eight assumption.
- `internal/eval/eval_test.go:662-732` applies the same bounded threshold and progress contract to the failure tail of the cancellation/progress fixture. It does not assert zero calls after the threshold callback.
- `internal/eval/eval_test.go:461-569` adds a deterministic two-worker drain fixture. The first admitted request publishes its request context and waits for cancellation. The other worker returns immediate non-retryable errors with known usage. The fixture verifies threshold cancellation, bounded calls (9–10), bounded counted errors (8–9), exclusion of the cancelled held sample from `Done`, usage totals from every executor return, zero active calls, held-request settlement, and final progress/result equality. Startup, completion, and cleanup waits have two-second timeouts.
- The threshold callback checks the held request context without blocking. It permits a post-callback handoff call, so it does not encode an exact post-cancellation call guarantee.
- The recorded temporary-copy negative control removed only the threshold `cancel()` call. It failed with missing threshold cancellation and worker-drain timeout, then cleaned up without modifying maintained sources or leaving a race.

No production files, run-state files, frozen patches, or unrelated tests required correction during this review.

## Findings (not fixed)

- None in the scoped fixture review.
- The complete pristine upstream service race baseline remains a separately documented fixture-only blocker. Docker lifecycle, authenticated browser integration, deployment, and live billing remain unverified by design.

## Verification

- Formatting: `gofmt -l internal/eval/eval_test.go` clean.
- Diff hygiene: `git diff --check -- internal/eval/eval_test.go` clean.
- Focused race stress: the three stop/progress fixtures passed at `GOMAXPROCS=4` and `GOMAXPROCS=8`, `-count=300`.
- Root vet: `go vet ./...` passed.
- Root race: `go test -race -count=1 ./...` passed. Evaluation core completed in 10.273 seconds and the existing plugin in 16.837 seconds.
- Copied package vet: `go vet ./internal/codexeval` passed in the staged pinned backend.
- Copied package race: `go test -race -count=1 ./internal/codexeval` passed in 11.631 seconds.
- Source parity: `diff -rq internal/eval dist/sub2api-build/workspace/backend/internal/codexeval` reported no differences.
- Frozen patch hashes remain unchanged: 001 `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`, 002 `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`, and 003 `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.
