# Code Research: Detector Port

## Source Identity

This report reads upstream `7a49068` and the current native production patches.
Native readable source comes from `/tmp/sub2api-settings-entry-final/source/backend/internal/codexeval/`.
The patches remain authoritative delivery inputs.
No implementation or check ran for this report.

## Required File Changes

| Effective native file | Action |
| --- | --- |
| `fingerprint.go` | Replace repeated probe collection with the upstream three-challenge collector |
| `fingerprint_algorithm.go` | Remove old probe normalization, JSD, permutation, and baseline attribution |
| `fingerprint_detector.go` | Add the upstream decoder, sequence check, feature extraction, ranker, and calibration |
| `modeltrace.go` | Remove the retired runner, progress maps, and sample types |
| `modeltrace_algorithm.go` | Remove the retired bank decoder, features, family aggregation, and attribution |
| `codex.go` | Retain the shared `mustDecode` helper and add challenge instructions through the request builder |

All these production changes belong in `patches/0001-evaluation-core.patch`.
Resource ownership belongs in `0002`; runtime and page integration belong to the other code reports.

The native `fingerprint_algorithm.go:47` currently defines `mustDecode`.
`codex.go:39` calls it during captured-context initialization.
Deleting the old algorithm without retaining this helper breaks package compilation.
Move the unchanged helper to `codex.go`, which already imports `encoding/json`.
The new detector and enrollment-suite decoder reuse the same helper.
Do not retain the old algorithm file only for this function.

## Collector Call Graph

```text
Start(type="fingerprint", targets, model, concurrency)
  -> native account admission and active-run registration
  -> runFingerprint(ctx, auth, model, concurrency, progress)
     -> fingerprintChallenges()
     -> per-account workers
        -> acquire one of six fingerprintSlots
        -> collectFingerprintSample(ctx, auth, model, challenge)
           -> askCodex(... instructions, prompt, challengeTimeout, accept)
              -> fingerprintText(answer, expectedCount)
        -> release slot
        -> indexed event with sample, usage, and error
     -> assemble samples in challenge order
     -> analyzeFingerprint(outputs)
     -> append native result and save history
```

Keep native admission, active-run cleanup, account metadata, and cancellation owners.
Do not import upstream `startRuns`, CPA credential lookup, or `records` globals.
The exact native runtime adapter belongs to [UI and history research](code-ui-history.md).

Upstream collector anchors:

- `7a49068:internal/plugin/fingerprint.go:29`: suite decoder and challenge instruction/prefix construction.
- `fingerprint.go:53`: three distinct random suite entries.
- `fingerprint.go:61`: comparison aliases.
- `fingerprint.go:135`: bounded worker collector and indexed result events.
- `fingerprint.go:170`: global slot acquisition through a cancellation-aware select.
- `fingerprint.go:185`: sample order and usage aggregation.
- `fingerprint.go:205`: attribution and final run status.
- `fingerprint.go:221`: per-challenge runner and answer validator.

The suite has 36 entries, including 12 developer-instruction entries and nine user-prefix entries.
For a prefixed entry, concatenate `UserPrefix`, `\n\nFinal task:\n`, and `Prompt` exactly as upstream does.
Do not append developer instructions to the user prompt.
Fingerprint requests omit an explicit reasoning effort.

## Answer Acceptance and Attribution Are Separate

`fingerprintText` limits answer text to 64,000 bytes.
`parseNumbers` finds the longest numeric run separated by non-letter content.
It accepts integers from 1 through 355 and ignores values outside that range.
Letters terminate the current run. An equally long later run does not replace the first longest run.

The minimum accepted count is `max(80, ceil(expectedCount * 0.55))`.
For an expected count of 300, the minimum is 165.
When a response contains extra numbers, `fingerprintText` keeps only the requested count before scoring.
An answer below the count limit can retry. Shape rejection does not retry.

`analyzeFingerprint` parses each retained answer and records its diagnostic.
It returns an error when no answer meets the count limit.
It then applies the sequence check before ranking candidates.

| Accepted answers | Shape check | Attribution |
| --- | --- | --- |
| Zero | No check | Error, no attribution |
| One to three | Fails | `decision: "not_model"`, sequence diagnostics, no ranking or probability |
| One or two | Passes | `decision: "partial"`, ranked candidates, null probabilities |
| Three | Passes | Ranked candidates with calibrated probabilities |

For a `not_model` attribution, the run status can still be `completed` or `partial`.
Do not confuse `record.status` with `attribution.decision`.
A repeated `42,` answer can pass the numeric-count validator but fail the shape check.
It must not consume extra attempts only because the final attribution is `not_model`.

Evidence: `7a49068:internal/plugin/fingerprint.go:237` and `fingerprint_detector.go:203,233,278`.

## Sequence Check and Ranking

`sequenceShape` computes each accepted answer's lag-one autocorrelation and share of steps no larger than 20.
It takes the median of each quantity across answers.
The sequence passes when autocorrelation is at most `-0.23` and small-step share is at most `0.05`.
Keep these upstream limits unchanged.

`detector.rank` combines frequency and positional scores.
The frequency score combines standardized LDA, nearest-reference, and baseline scores with weights `0.5`, `0.25`, and `0.25`.
LDA uses the first 128 numbers. The nearest-reference comparison uses full-answer features and seven nearest references.
Full features combine Hellinger value frequencies and ordered blocks with shares `0.75` and `0.25`.
The positional component combines value frequency, ordered blocks, position kernels, and sequence features.
Its block shares are `0.6`, `0.15`, `0.125`, and `0.125`.
The final blend reads its weights from the frozen detector.

For three accepted, model-like answers, `softmax` applies the detector's calibration temperature to all 60 candidate scores.
Candidate order uses a stable descending score sort.
Partial results retain `*float64` probability fields as nil. Do not substitute zero.
The new detector does not produce the old ModelTrace family probability aggregates.

Evidence: `7a49068:internal/plugin/fingerprint_detector.go:309,327,366,395,410,422,452,476,580`.

Port the full upstream algorithm and its small math helpers together.
Do not mix the old ModelTrace feature implementation with the new ranker.
Their similar names do not establish equivalent feature or calibration contracts.

## Resource Decoder

`decodeDetector` reads a little-endian 32-bit JSON-header length.
The float section starts at the next four-byte boundary after that header.
The file stores float32 bytes by byte plane, rather than as ordinary adjacent float32 values.
The decoder reconstructs those values and converts them to `float64`.
`detectorVector.UnmarshalJSON` resolves each JSON `$f32` offset and length into that shared array.
It uses a full slice expression to restrict each vector's capacity.
The decoded vectors retain the array after the temporary global decoder reference becomes nil.

The decoder checks upstream detector/calibration schemas, model order, bank timestamp, binding metadata, and response counts.
These checks belong to the imported loader. Add no separate source-integrity or runtime validation framework.
Precompute reference norms through the upstream loop.
Keep the upstream `detector_revision` field and its existing revision calculation as result metadata.
Do not add a separate task hash gate.

Evidence: `7a49068:internal/plugin/fingerprint_detector.go:21,30,38,115,119`.

## Comparison Aliases

Use the same alias table in Go and the page:

| Requested model | Enrolled candidate |
| --- | --- |
| `gpt-6.1-sol` | `gpt-6-astra` |
| `deepseek-flash` | `deepseek-v4.1-flash-1017` |

Comparison trims spaces, converts to lowercase, removes a provider prefix, then applies the alias.
This comparison must not change the literal model sent through the native sender.
Evidence: `7a49068:internal/plugin/fingerprint.go:61`.

## Focused Acceptance Evidence

Reuse upstream `TestFingerprintUpstreamParity` and its four golden cases.
Those cases cover two calibrated rankings, one partial ranking, and one shape rejection.
They compare candidate order, scores, nullable probabilities, diagnostics, and sequence measures at tolerance `1e-9`.
Reuse `TestFingerprintParseAndChallenges` for parsing, count limits, suite construction, distinct challenges, and both aliases.
Adapt only their native package and fixture ownership.

Do not copy CPA `fingerprintTestHost`, `setupTest`, `hostCall`, or management routes into native tests.
Use the existing native Sender fixtures for collector, persistence, cancellation, and busy-state coverage.
The exact test selection belongs to [delivery research](code-delivery.md).

## Corrections to the Earlier Plan

- Three count-valid answers alone do not guarantee rankings or probabilities. The shape check must also pass.
- One or two count-valid answers can also return `not_model`, rather than a partial ranking.
- Removing `fingerprint_algorithm.go` also removes a helper used by Codex context initialization. Retain that helper under its live owner.
- Retain nullable probability fields and the new sequence diagnostic. Remove retired family probability rendering.

These corrections clarify the approved upstream behavior. They do not introduce another product decision.
