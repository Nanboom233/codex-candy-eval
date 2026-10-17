# Fingerprint Stop Contract Review

## Verdict

**PASS for the compatibility contract.** The source-backed contract is the original threshold cancellation with worker drain and observed-usage retention. The reopened implementation gate still requires replacing the two strict-eight assertions with fixtures that prove this behavior.

## Evidence

- `internal/plugin/fingerprint.go:157-174` and `internal/eval/fingerprint.go:122-138` use the same unbuffered event channel, worker loop, and worker-wait close. A worker's send completes when the aggregator receives its sample, before the aggregator updates the consecutive-error count or cancels.
- The aggregators in `internal/plugin/fingerprint.go:175-202` and `internal/eval/fingerprint.go:139-169` count samples serially and cancel when the eighth consecutive error is handled. They continue receiving until workers settle. The core adds progress publication but retains this ordering.
- The original `TestFingerprintUpstreamErrors` in `internal/plugin/fingerprint_test.go:266-279` runs at the default concurrency of two and accepts 8 through 11 errors. It does not define an exact request count at concurrency one.
- The only strict-eight assertions are `TestFingerprintEightErrorsAndModelTracePartial` at `internal/eval/eval_test.go:428-438` and the failure tail of `TestProgressCancellationAndFailureSettleWorkers` at `internal/eval/eval_test.go:600-609`.
- `internal/eval/model.go:74-115,136-148` retains response token counts when cancellation is observed, while classifying a cancelled probe separately. The fingerprint aggregator adds usage before skipping cancelled samples at `internal/eval/fingerprint.go:141-147`.
- The spec says Fingerprint “stops collection after eight consecutive request errors” in `.trellis/spec/backend/evaluation-core.md`. It does not promise exactly eight executor calls or define the treatment of admitted work.

The repeated concurrency-one failure recorded in `research/fingerprint-stop-race.md` is therefore a failure of the new exact-count assertions, not evidence that the migrated pipeline changed the original dispatch ordering.

## Stopping Contract

When the aggregator handles the eighth consecutive non-cancelled request error, it marks the result failed, records the threshold error, and cancels the collection context immediately. Workers then stop dispatching when they observe cancellation. The aggregator drains their result events and waits for worker settlement before returning final progress and results.

Work admitted before cancellation may still reach the executor or remain in flight. Those results must settle. Retain token usage returned by every drained sample, including samples categorized as cancelled. Count only non-cancelled samples in `Done`, `Errors`, and `Valid`; do not require those counters to include work cancelled during settlement. Final progress must match result counters and publish only after no workers remain.

With one worker and an immediate non-retryable executor, the send handoff can admit one next sample before the aggregator applies cancellation. The deterministic bound for that fixture is 8–9 executor calls. Do not generalize the original 8–11 test allowance into a universal bound for other concurrency settings.

Do not assert that executor calls remain zero after the eighth-error progress callback. That callback runs after the aggregator cancels, but a worker can have passed its last context check before cancellation and reach the executor across that boundary. Assert cancellation, bounded work, complete drain, usage retention, and final-progress equality instead.

## Fixture Design

1. Update the concurrency-one immediate-error fixture to use a non-retryable status and report one input token per executor return. Require a non-nil error, failed status with the eight-error stop message, 8–9 executor calls, 8–9 counted errors, `Done == Errors`, and input tokens equal to usage returned by the calls. Require final progress status/phase and counters to match the result.
2. Add a concurrency-two drain fixture. Make the first executor invocation publish its request context to a buffered channel, hold until that context is cancelled, then return observed usage and the cancellation error. Let the other worker return immediate non-retryable errors with known usage. The controlled schedule admits the held call plus eight threshold errors, with at most one post-threshold handoff call, so executor calls are bounded to 9–10 and counted errors to 8–9. In the first progress callback with at least eight errors, nonblockingly verify that the held request context is already cancelled. After `RunOne` returns, require the held call to have settled, active calls to be zero, failed status, `Done == Errors`, final progress/result counter equality, and total input tokens equal the atomically recorded sum returned by the executor. Use a timeout while waiting for the held-call and completion channels. Do not assert zero executor calls after the threshold callback.

Use channels to synchronize the held request and atomic counters for active calls and returned usage. Keep callbacks prompt and nonblocking. Existing callback serialization and monotonicity coverage can remain unchanged.

## Alternative

An acknowledgment from the aggregator before a worker starts its next sample could make the concurrency-one fixture stop at exactly eight calls. That changes the original worker pipeline and still cannot prevent already-admitted requests at higher concurrency. Adopt it only with an explicit requirement for this stricter dispatch boundary and tests defining its concurrency and drain semantics. The current spec and original fixture do not establish that requirement.

## Scope and Verification

This review changed only this research note. No product or test files were edited, and no test suite was run as part of this contract-only review. The focused race reproductions and source comparisons are recorded in `research/fingerprint-stop-race.md`.
