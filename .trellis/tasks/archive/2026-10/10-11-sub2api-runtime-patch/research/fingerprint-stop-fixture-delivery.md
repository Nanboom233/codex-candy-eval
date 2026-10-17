# Fingerprint Stop Fixture Delivery

- Date: 2026-10-12, checks completed by 02:37 UTC.
- Scope: compatibility fixture correction after the independent stop-contract PASS.
- Status: implementation complete. Independent correction review remains pending.

## Changes

Only `internal/eval/eval_test.go` and this report changed under this assignment.
No production worker ordering, run-state sources, host patches, or build scripts changed.

`TestFingerprintEightErrorsAndModelTracePartial` now permits eight or nine executor calls and counted errors.
It requires the threshold failure message, equal Done/Errors, no valid samples, exact returned input usage, and matching final progress.
The ModelTrace partial-result assertions remain unchanged.

The failure tail of `TestProgressCancellationAndFailureSettleWorkers` uses the same bounded contract.
It now checks returned input usage and every final progress count against the result.
Existing cancellation, callback serialization, and monotonicity fixtures remain intact.

`TestFingerprintThresholdCancellationDrainsUsage` uses two workers and holds the first request until its context is cancelled.
Other executor calls wait for the held request's publication, then return immediate non-retryable errors with known usage.
The first progress callback with at least eight errors checks the held request's cancellation without blocking.
Callbacks do not wait for completion or forbid already admitted calls across the cancellation boundary.

The fixture bounds calls to nine or ten and counted errors to eight or nine.
It checks failed status, the threshold message, Done/Errors equality, zero valid samples, and exclusion of the held cancelled sample from Done.
It checks all executor returns, zero active calls, held-request settlement, and final progress/result equality.
It retains input, output, and reasoning usage from every return, including the held cancelled request.
Both startup and completion waits have two-second bounds.
Cleanup cancels the parent and waits for the run goroutine to exit, including failure paths.

## Red-to-Green Evidence

The existing strict-eight fixtures failed repeatedly under multicore stress before correction.
The independent reproduction and historical source comparison are recorded in `fingerprint-stop-race.md`.
Those failures exposed overly strict assertions, rather than changed production ordering.

A separate temporary module copied the final evaluation sources and root module files.
Only that copy's threshold `cancel()` call was removed.
The new fixture failed deterministically with `threshold progress preceded request cancellation` and `threshold cancellation did not drain workers`.
The negative control took 2.630 seconds and used `-race -count=1 -timeout=15s` with GOMAXPROCS=4.
Cleanup settled the held work, with no cleanup timeout or data-race report.
The temporary directory was removed, and maintained production sources were never mutated.

## Verification

The focused pattern selected both affected fixtures and the new drain fixture:

```text
^(TestFingerprintEightErrorsAndModelTracePartial|TestProgressCancellationAndFailureSettleWorkers|TestFingerprintThresholdCancellationDrainsUsage)$
```

| Check | Result |
| --- | --- |
| GOMAXPROCS=4 focused race, count=300 | PASS, 5.901s |
| GOMAXPROCS=8 focused race, count=300 | PASS, 6.028s |
| Root `go vet ./...` | PASS |
| Fresh root `go test -race -count=1 ./...` | PASS, eval 9.112s, existing plugin 15.603s |
| Safe `stage-eval-core` helper | PASS, previous generated tree preserved |
| Pinned backend `go vet ./internal/codexeval` | PASS |
| Pinned backend `go test -race -count=1 ./internal/codexeval` | PASS, 10.472s |
| Complete maintained/staged tree equality | PASS, all 28 files |
| Current tree versus previous generated backup | Only `eval_test.go` differs |
| `gofmt -l internal/eval/eval_test.go` | Clean |
| `git diff --check` and changed-file trailing whitespace | PASS |

The staging command was:

```sh
python3 scripts/sub2api_build.py stage-eval-core compatibility/sub2api-v0.2.15.json . dist/sub2api-build/workspace
```

The helper preserved its prior tree at `.codex-eval-backups/codexeval.prev-20261012T023603Z` beneath the ignored pinned workspace.
All production files and the four Stage 4B sources match that backup byte-for-byte.
The final 28-file tree digest is `3eabe2af794582dc8e41dfd637e57fc55940b9e57bc13e9606c1d748f9a8a7ba`.
This digest hashes sorted relative paths, a NUL byte, and each file's binary SHA-256 digest.
The checks used Go 1.27.2 on darwin/arm64.

Frozen patch hashes remain unchanged:

- 001: `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`
- 002: `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`
- 003: `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`

## Limits and Follow-ups

Independent correction review must close the reopened stop-fixture gate before Stage 4B reviews proceed.
The focused bounds describe these controlled schedules, not every possible concurrency configuration.
The complete pristine upstream race suite remains separately baseline-blocked and was not rerun or repaired here.
Docker lifecycle, authenticated browser integration, deployment, and live accounting remain unverified.
No private configuration, CPA deployment, paid model call, commit, or push was performed.
