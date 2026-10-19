# Research: Upstream Evaluation Changes

- Query: Compare upstream evaluation changes with the native Sub2API patches and tests.
- Scope: mixed
- Date: 2026-10-19

## Findings

### Commit inventory

| Commit | Effect on this research line |
| --- | --- |
| `0de3802` | Release version and UI changes. No runner or model-request change. |
| `63098bd` | Replaces the old fingerprint test and ModelTrace with LM fingerprint detection. |
| `c363f6b` | Release version only. |
| `009c26e` | README screenshot only. |
| `6573867` | Adds bounded retries for failed requests and unusable answers. |
| `f2aeab2` | Removes UI tooltips. No runner change. |
| `96deb3c` | README screenshot only. |
| `d093ee8` | Versions saved state and renames detection records and tests to fingerprint records. The history migration belongs to the history owner. |
| `dedb1bd` | Condenses repository instructions. No runner change. |
| `b4c6b60` | Updates the Codex request identity and captured context to `codex-tui 0.163.0`. |
| `7a49068` | Release version only. |

The main evaluation changes are `63098bd`, `6573867`, `d093ee8`, and `b4c6b60`.

### The new detector changes the old test contract

The shared baseline has three inspection choices: Candy, the old fingerprint test, and ModelTrace. The old fingerprint test compares repeated samples against enrolled baselines. Its quick, standard, and strict modes collect 60, 200, and 400 logical samples. The new upstream change removes that test and ModelTrace, then replaces them with a three-challenge LM fingerprint detector. It does not preserve the old repeated-baseline comparison. Treat this as a new classifier, not an exact port of either old test. See `6b46d4f:internal/plugin/inspection.go:24`, `6b46d4f:internal/plugin/fingerprint.go:21`, `6b46d4f:internal/plugin/fingerprint.go:28`, and `6b46d4f:internal/plugin/fingerprint.go:128`.

Each run selects three challenges from a 36-challenge suite. Twelve challenges add developer instructions, and nine add a user prefix. Each challenge allows at most three attempts, or two retries, so a run can make at most nine supplier requests. The saved `fingerprintSample` includes its actual per-challenge `Attempts` count as optional JSON field `attempts`; it is not only an in-memory counter. The detector accepts a sequence only when it has at least `max(80, ceil(expected_count * 0.55))` numbers from 1 through 355. See `7a49068:internal/plugin/fingerprint.go:53`, `7a49068:internal/plugin/fingerprint.go:76`, `7a49068:internal/plugin/fingerprint.go:220`, and `7a49068:internal/plugin/fingerprint_detector.go:203`.

After numeric-count acceptance, the detector checks sequence shape.
It requires median lag-one autocorrelation at or below `-0.23`.
The median share of steps no larger than 20 must be at or below `0.05`.
Shape rejection produces `not_model`, without rankings or probabilities, and does not retry.
After shape acceptance, three count-valid answers permit calibrated probabilities across 60 candidates.
One or two permit partial rankings with null probabilities.
Sources: `7a49068:internal/plugin/fingerprint_detector.go:233,278` and [detector code research](code-detector.md).

The bank maps `gpt-6.1-sol` to candidate `gpt-6-astra`. The detector compares the requested model with the top candidate. Its probabilities compare only enrolled candidates. They do not establish model identity. See `7a49068:internal/plugin/fingerprint.go:61` and `7a49068:internal/plugin/data/lm-detector/README.md:22`.

Inspection verdicts are not interchangeable. The shared baseline accepts Candy, old fingerprint, or ModelTrace and treats any conclusive degraded record since the inspection start as degraded. Current upstream accepts Candy or the new fingerprint detector; any completed/partial attributed mismatch can mark an account degraded without an additional confirmation round. Native automatic inspection runs Candy only. It reads the saved verdict after the first result, uses matching complete Candy results to confirm only possible state transitions, and applies model control and optional priority changes after the account's state is established. Preserve that native behavior. The new detector can replace the old manual fingerprint workflow, but it must not become the automatic inspection input without a separate decision. See `6b46d4f:internal/plugin/inspection.go:24`, `6b46d4f:internal/plugin/inspection.go:83`, `6b46d4f:internal/plugin/inspection.go:103`, `7a49068:internal/plugin/inspection.go:23`, `7a49068:internal/plugin/inspection.go:82`, `7a49068:internal/plugin/inspection.go:103`, `patches/0005-inspection.patch:435`, `patches/0005-inspection.patch:446`, and `patches/0005-inspection.patch:564`.

### Retry unusable results, not valid wrong answers

Upstream `askCodex` makes at most three total attempts. It waits two seconds before the first retry and four seconds before the second. It builds one request before the retry loop, sums usage across attempts, and keeps the last answer. This is broader than the shared baseline's existing `executeProbe`, which retries only request status `0`, `429`, or `5xx`; that code does not retry an unusable answer. See `7a49068:internal/plugin/model.go:47`, `7a49068:internal/plugin/model.go:67`, `6b46d4f:internal/plugin/model.go:29`, and `6b46d4f:internal/plugin/model.go:44`.

Upstream retries timeouts and server errors rejected before a stream starts. It does not retry `401` or `429` at that boundary because CPA handles pinned-account retries and cooldowns. After a stream starts, it retries stream, incomplete-response, blank-answer, and tool-only errors. A caller-supplied validator can also reject an unusable answer. Cancellation stops the wait or the next attempt. See `7a49068:internal/plugin/model.go:62`, `7a49068:internal/plugin/model.go:103`, and `7a49068:internal/plugin/model.go:156`.

Candy passes no answer validator. A text answer is a completed test, whether it is correct or wrong. A wrong answer must not retry. A missing text answer or a failed response is unusable and can retry. Fingerprint detection retries answers that are too long or contain too few usable numbers. Pelican retries answers without acceptable HTML. See `7a49068:internal/plugin/candy.go:67`, `7a49068:internal/plugin/fingerprint.go:220`, and `7a49068:internal/plugin/pelican.go:66`.

This distinction preserves native confirmation counts. One logical Candy test still creates one result after its internal attempts. Only completed Candy results enter confirmation answers. A wrong answer remains `completed` with `OK=false`. Request errors remain invalid confirmation results. See `patches/0005-inspection.patch:618` and `patches/0005-inspection.patch:635`.

Native inspection checks account participation before each logical test. The retry loop must also call that check before every new supplier request. If participation changes during a retry wait, it must stop before sending again. It must keep the account outcome `unchanged` and avoid a synthetic evaluation error. Current coverage checks opt-out between logical tests, not between retries. See `patches/0001-evaluation-core.patch:54`, `patches/0001-evaluation-core.patch:2082`, and `tests/backend/internal/codexeval/inspection_test.go:540`.

### Adapt retries to the direct native sender

The upstream retry rules depend on CPA's `host.model.execute_stream` retry layer. Native Sub2API sends each request through `AccountTestService.SendCodexEvaluation` to the selected account. It returns HTTP status and stream data directly. It does not use that CPA host layer. The native sender also reconciles `429` account state before it returns. See `patches/0004-admin-channel.patch:211`, `patches/0004-admin-channel.patch:342`, and `patches/0004-admin-channel.patch:352`.

Keep model execution on this native sender. Do not add another HTTP client or bypass its account, proxy, API-key, OAuth, Chat Completions, stream-conversion, or cancellation handling. The sender accepts only Codex protocol headers. OAuth requests pass through native identity normalization. See `patches/0004-admin-channel.patch:290`, `patches/0004-admin-channel.patch:298`, and `patches/0004-admin-channel.patch:318`.

The current `executeProbe` retries status `0`, `429`, and server errors.
The user approved immediate stop for native `401`, `429`, and other permanent HTTP responses.
Keep existing native account error and cooldown handling.
Replace blanket status-0 retries with the source-backed distinction in [retry code research](code-retry-inspection.md).
Sources: `patches/0001-evaluation-core.patch:998,1015` and `.trellis/spec/backend/error-handling.md:39`.

Separate permanent setup errors from retryable transport errors. The native sender can return account lookup, credential, and request-validation errors before any supplier request. Do not retry these as network failures. After a stream starts, preserve native stream closure and retry only when the run context remains active. Reuse the current lifetime test for cancellation. See `patches/0004-admin-channel.patch:213`, `tests/backend/internal/service/codex_eval_lifetime_test.go:42`, and `tests/backend/internal/service/account_codex_eval_service_test.go:19`.

### Update the request identity and detector resources

Upstream updates the user agent to `codex-tui/0.163.0 (Debian 13.0.0; x86_64) xterm-ghostty (codex-tui; 0.163.0)`. It also sets `history_initialization` to `new`. The context file comes from a Linux capture. Detection instructions enter the first Codex context message as `generic.developer_instructions`. See `7a49068:internal/plugin/codex.go:16`, `7a49068:internal/plugin/codex.go:20`, `7a49068:internal/plugin/codex.go:131`, and `7a49068:internal/plugin/codex_test.go:34`.

Port the request fields and context into the existing `codexeval` request builder. Continue to pass the request through the native sender and its header allowlist. Test the request JSON, developer-instruction placement, user agent input, and native OAuth/API-key handling. Do not copy CPA's `hostCall` protocol into Sub2API.

The detector uses the MIT-licensed `Ikaleio/lm-detector` port pinned at `654c8334bb53951aaf04cbd61252d163649e1369`. The binary is 4,385,348 bytes. The suite is 152,315 bytes. The 60-model file is 9,921 bytes. The Go code embeds these files and the license. It adds no Go module dependency. See `7a49068:internal/plugin/data/lm-detector/README.md:3` and `7a49068:internal/plugin/fingerprint_detector.go:21`.

Detector decoding runs during package initialization. The 4-byte length prefix stores `70,712`. The decoder expands the remaining binary floats into 1,078,658 `float64` values, about 8.2 MiB. No detector benchmark or production runtime measurement exists in the reviewed commits. A normal run has three logical challenges and initially sends one request per challenge. Each challenge allows three total attempts, so the maximum is nine supplier requests, not nine retries. The run allows three concurrent challenges per account and six across accounts. See `7a49068:internal/plugin/fingerprint_detector.go:119`, `7a49068:internal/plugin/fingerprint.go:17`, `7a49068:internal/plugin/fingerprint.go:76`, and `7a49068:internal/plugin/model.go:47`.

Production resources belong in `patches/0002-evaluation-data.patch`. Detector code belongs in `patches/0001-evaluation-core.patch`. Keep the golden detector data under `tests/backend/`. Replace the ModelTrace notice with the new MIT notice in the project's third-party notice. The current owners are `patches/0001`, `patches/0002`, `tests/backend/internal/codexeval/fingerprint_test.go`, and `licenses/THIRD_PARTY_NOTICES.txt`.

### Focused test delta

Replace the old fingerprint statistics and ModelTrace tests with two adapted upstream algorithm tests.
Use `TestFingerprintUpstreamParity`, `TestFingerprintParseAndChallenges`, and their golden fixture.
Use upstream collector and retry tests as references only.
Extend existing native owners and add one `TestCodexEvalRetries` owner for the missing shared-runner coverage.
The exact selector and test ownership are in [delivery code research](code-delivery.md).

Add focused native checks for these changed contracts:

1. Retry a parse-rejected fingerprint answer, then accept the next usable answer. Count all attempts and usage.
2. Record a valid wrong Candy answer after one attempt. Keep it as one inspection answer.
3. Change inspection participation during retry backoff. Assert that no later request is sent and no synthetic failure replaces `unchanged`.
4. Cover retryable and final native statuses, stream failures, and cancellation. Reuse `TestCodexEvalBackgroundNativeStreamLifetime` for stream lifetime.
5. Keep complete confirmation counts unchanged. Extend `TestCodexInspectionConfirmationAndBoundaries` only where the retry path changes its evidence.

Update the focused `-run` expression in `scripts/build.py:53`. Replace the retired ModelTrace test names with the new detector and runner names. Reuse `TestCodexEvalNativeSender`, `TestCodexInspectionConfirmationAndBoundaries`, `TestCodexInspectionParticipationStopsLaterSends`, and `TestCodexEvalBackgroundNativeStreamLifetime`. Do not add a full-suite, race, or coverage gate.

### Accepted decisions

1. Replace the old fingerprint baseline test and remove ModelTrace. Keep the new detector as a manual candidate-attribution test. Its output has different meaning from the old baseline comparison.
2. Keep automatic inspection on Candy. This preserves the existing wrong-answer, confirmation, model-control, priority, and recovery contracts.
3. Stop immediately after native `401` or `429`. Preserve native account state updates.
4. Start new fingerprint history empty. Drop retired histories without conversion, backup, or a legacy page.
5. Keep native Candy and Pelican history formats and all prior records.

## Related Local Owners

- `patches/0001-evaluation-core.patch`: model request runner, response parser, Candy, fingerprint, and Pelican.
- `patches/0002-evaluation-data.patch`: Codex context and embedded evaluation data.
- `patches/0004-admin-channel.patch`: native account sender and request boundary.
- `patches/0005-inspection.patch`: automatic Candy inspection, confirmations, and participation checks.
- `tests/backend/internal/codexeval/`: runner and inspection tests.
- `tests/backend/internal/service/`: native sender and stream lifetime tests.
- `scripts/build.py`: focused Go test selection.
- `.trellis/spec/backend/error-handling.md`: current retry contract.

## External References

- [Ikaleio/lm-detector at the pinned commit](https://github.com/Ikaleio/lm-detector/tree/654c8334bb53951aaf04cbd61252d163649e1369)
- Upstream repository commits `63098bd`, `6573867`, `d093ee8`, and `b4c6b60`.

## Caveats / Not Found

- No product patch, build, install, or supplier request ran during this research.
- The final server and installer size after embedding remain unmeasured.
- This report excludes page design and saved-state migration. The assigned owners must report those findings separately.
