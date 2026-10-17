# Research: Fingerprint eight-error stop ordering

- Query: Reproduce the concurrency-one nine-error failure and compare current core ordering with the original plugin's stopping behavior and tests.
- Scope: internal
- Date: 2026-10-12

## Findings

### Reproduction

- A first focused run passed: `go test -race -run '^TestFingerprintEightErrorsAndModelTracePartial$' -count=1 ./internal/eval`.
- `env GOMAXPROCS=1 go test -race -run '^TestFingerprintEightErrorsAndModelTracePartial$' -count=300 ./internal/eval` passed.
- The same focused command with `GOMAXPROCS=4` and `GOMAXPROCS=8`, each at `-count=300`, failed repeatedly. A captured failure had `Concurrency:1`, `Done:9`, `Errors:9`, `calls=9`, and nine input tokens; the result remained `Status:failed` with the expected eight-error message.
- `env GOMAXPROCS=4 go test -run '^TestFingerprintEightErrorsAndModelTracePartial$' -count=1000 ./internal/eval` also failed repeatedly without the race detector.
- `env GOMAXPROCS=4 go test -race -run '^TestProgressCancellationAndFailureSettleWorkers$' -count=100 ./internal/eval` failed its second strict-eight assertion at `eval_test.go:608` (`failure progress mismatch`). This is the same over-count reflected through final progress.
- A broader `Fingerprint` race subset at `-count=50` passed, confirming that a passing rerun does not resolve this schedule-dependent failure.
- The Stage 4B implementer separately reports the initial root `go test -race -count=1 ./...` failure with the same nine-call result and a passing copied-package run that round. Their report says Stage 4B left all 24 Stage 4A sources unchanged; the reopened-gate record is at `implement.md:1225-1232`.

### Confirmed ordering window

- In [fingerprint.go](../../../../internal/eval/fingerprint.go#L122), `events` is unbuffered. Each worker checks `ctx.Err()`, claims an index, collects a sample, and sends it at lines 129-135. The send completes when the aggregator receives the value; it does not wait for the aggregator to increment the error count.
- The single aggregator loop in [fingerprint.go](../../../../internal/eval/fingerprint.go#L141) adds observed usage, increments `Done`/`Errors` and the consecutive count, then cancels at eight at lines 148-157. Thus the worker can resume from the eighth send, observe a still-live context, and enter sample nine before the aggregator reaches `cancel()`.
- `executeProbe` and `safeExecute` check cancellation around slot acquisition and executor calls in [model.go](../../../../internal/eval/model.go#L74). Those checks preserve ordinary cancellation, but they do not order the worker after aggregation: the context has not yet been cancelled in the observed window. The test executor returns immediately and ignores its context, making a ninth call visible when the worker wins scheduling.
- This is a logical scheduling race, not evidence of an unsynchronized memory access. The mutable result and consecutive count belong to the single aggregator; `go test -race` exposes the interleaving but reports no data-race stack.
- The same send-before-aggregation ordering is present in the complete working reference [plugin fingerprint runner](../../../../internal/plugin/fingerprint.go#L157) at lines 164-171 and 177-191. The migrated core preserved that implementation pattern.

### Historical contract versus new assertion

- The original [plugin test](../../../../internal/plugin/fingerprint_test.go#L266) `TestFingerprintUpstreamErrors` runs with the default concurrency of two and accepts `r.Errors` from 8 through 11 while requiring a failed result and no valid samples (lines 271-278). No original plugin test found here requires exactly eight calls or errors at concurrency one.
- The core spec says Fingerprint “stops collection after eight consecutive request errors” in [evaluation-core.md](../../../../.trellis/spec/backend/evaluation-core.md#L67) and requires a test for that behavior at lines 108-115. It does not define whether already dispatched requests count in the final result or require exactly eight calls.
- The new [core regression](../../../../internal/eval/eval_test.go#L428) makes the stronger assertion `calls == 8`, `Done == 8`, `Errors == 8`, and input tokens equal eight at concurrency one (lines 429-437). [ProgressCancellationAndFailureSettleWorkers](../../../../internal/eval/eval_test.go#L545) repeats the strict final `Done == 8` and `Errors == 8` requirement at lines 600-608. These are the only exact-eight assertions found.
- The Stage 4A specification and quality reviews recorded passing root and copied-package race runs ([spec review](eval-core-spec-review.md#L31), [quality review](eval-core-quality-review.md#L32)); their recorded checks did not catch this interleaving. The later Stage 4B run reopened the gate. Passing review runs do not invalidate the reproduced failure.
- Source-backed judgment: the observed ninth call violates the new tests' strict concurrency-one contract, but exact-eight was not part of the reference test contract. Under the historical interpretation, cancellation is decided when the aggregator handles the eighth error and already dispatched work may settle; under a stricter interpretation, the core must prevent a sole worker from dispatching the next sample until that aggregation decision. The spec wording leaves this distinction open, so source alone cannot call the legacy behavior a semantic regression.

### Repair options and regression design

1. Preserve the reference's threshold-plus-in-flight behavior. Keep status failed once the aggregator observes eight consecutive errors, and test the stop decision while allowing already dispatched results and their observed usage to settle. For concurrency one, the current code can expose one extra call in the handoff window; with higher concurrency, tests must allow work already dispatched before cancellation. This aligns with the original `8..11` fixture evidence, but requires changing the new exact-eight assertions and clarifying the spec. Do not merely weaken the checks to `>=8`; also prove no new work starts after the stop decision and bound/settle work already in flight.
2. Adopt the stronger exact-eight guarantee at concurrency one. The smallest coordination change is to send an event containing an acknowledgment channel and have each worker wait for the aggregator to finish accounting and apply threshold cancellation before looping for another sample. Acknowledging after the threshold decision lets one worker stop at eight; with concurrency greater than one, requests dispatched before the eighth result is aggregated remain in flight and must still be drained, with their observed usage retained. The aggregator must continue reading until workers settle and acknowledge cancelled events so no worker blocks or leaks.

For option 2, retain a deterministic concurrency-one regression with an immediate non-retryable executor, `Concurrency=1`, exact call/count/token assertions, and final failed progress equality. Add a gated concurrency-greater-than-one fixture: keep multiple requests active, arrange the eighth error while other requests are already active, then have those requests settle through cancellation with usage. Assert no replacement dispatch after the threshold decision, preserve all returned usage, wait for active executor calls to reach zero before final progress, and use timeouts to detect blocked workers. Existing [progress tests](../../../../internal/eval/eval_test.go#L454) already check serialized callbacks, monotonic counts, final/result equality, and zero live requests at final progress; retain those constraints and do not add a lock around callbacks.

## Files Found

- `internal/eval/fingerprint.go`: worker event send and aggregator threshold cancellation.
- `internal/eval/model.go`: per-probe cancellation, retry, and executor boundary.
- `internal/eval/eval_test.go`: strict-eight and progress-settlement assertions, plus callback/concurrency tests.
- `internal/plugin/fingerprint.go`: original worker and aggregator ordering.
- `internal/plugin/fingerprint_test.go`: original test allows 8-11 errors at default concurrency.
- `.trellis/spec/backend/evaluation-core.md`: threshold wording and required evaluation test contracts.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/eval-core-migration.md`: migration source map for the original Fingerprint runner.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/eval-core-spec-review.md` and `eval-core-quality-review.md`: Stage 4A review and verification records.
- `.trellis/tasks/10-11-sub2api-runtime-patch/implement.md`: Stage 4A freeze and later gate reopening evidence.

## Caveats / Not Found

- No external references were needed; the behavior and intended compatibility boundary are defined by repository source, tests, specs, and task records.
- No prior plugin test was found that pins concurrency-one stop counts. The original range test is explicit evidence for tolerated overrun at concurrency two, not a proof of the exact maximum for every configured concurrency.
- This investigation made no product or test changes. No full suite was rerun here; focused race and non-race commands above were used to reproduce the failure without treating a pass as resolution.
