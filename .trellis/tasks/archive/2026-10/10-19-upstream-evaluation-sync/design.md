# Native Upstream Evaluation Sync

## Source and Ownership

Adapt the approved behavior from upstream `7a49068` into Sub2API v0.2.17.
The commit inventory and exclusions remain in [Git research](research/git-integration.md).
Keep the existing seven-patch order and original Sub2API Go module.

| Owner | Change |
| --- | --- |
| `patches/0001-evaluation-core.patch` | Detector, request builder, retry runner, task dispatch, and history |
| `patches/0002-evaluation-data.patch` | Detector resources and captured Codex context |
| `patches/0004-admin-channel.patch` | Native sender, HTTP status, and focused service/repository transport error hunks |
| `patches/0005-inspection.patch` | Participation checks before retry sends and unchanged confirmation counting |
| `patches/0006-evaluation-ui.patch` | Three test pages, detector controls, results, and preferences |
| `tests/backend/`, `tests/ui/` | Adapted focused tests and test-only data |
| `scripts/build.py` | Replace retired algorithm test names in the existing focused selection |
| `licenses/`, `README.md` | Detector attribution and Chinese user instructions |
| `.trellis/spec/backend/` | Current detector, history, retry, and three-page contracts |

Retain vendor assets, native routing, frontend embedding, installer source, and release ownership.
Do not restore the standalone CPA source tree, module, launcher, installation scripts, or documentation tree.

## Request and Result Flow

The native administrator page submits selected numeric account IDs through the existing start API.
The runtime keeps one active run per account across all test types.
Cancellation continues to use the exact native `run_id`.
Manual tests retain native admission, group membership, raw model discovery, and literal model sending.

All requests use `AccountTestService.SendCodexEvaluation`.
Preserve API-key and OAuth credentials, custom account User-Agent, proxy rules, Chat Completions conversion, and stream cancellation.
Do not add another HTTP client or a CPA host call.

Use this default User-Agent:

```text
codex-tui/0.163.0 (Debian 13.0.0; x86_64) xterm-ghostty (codex-tui; 0.163.0)
```

Use the upstream captured context with `history_initialization: new`.
Place challenge developer instructions in the corresponding `generic.developer_instructions` context field.
Keep the native sender's protocol-header allowlist and account override behavior.

### Native API and Runtime Adapter

Keep `Target`, `Info`, `Sender`, `Resolver`, `run`, and `StartSummary` under their current native owners.
Keep numeric `group_id` and `account_id`, native administrator envelopes, and account-ID history keys.
Do not import upstream credential IDs, management routes, `model_catalog`, `skipped`, or `unchecked` admission.

Change `StartRequest` to this field set:

```go
type StartRequest struct {
    Type        string   `json:"type"`
    Targets     []Target `json:"targets"`
    Model       string   `json:"model"`
    Effort      string   `json:"effort"`
    Runs        int      `json:"runs"`
    Concurrency int      `json:"concurrency"`
}
```

Remove `Mode`, its lookup, and all `modeltrace` dispatch branches from `startChecked`.
For fingerprint, default zero concurrency to three and retain the upstream range of one through three.
Keep existing Candy run bounds and model-name checks. Add no unrelated input checks.
Set fingerprint progress to `done: 0, total: 3` through the native active run.
Keep separate account-level and global challenge concurrency.

| Native path | Final contract |
| --- | --- |
| Start | `candy`, `fingerprint`, or `pelican`, with existing `started`, `busy`, and `unavailable` lists |
| State | `active`, `candy`, `fingerprint`, `pelican`, and `storage_error` |
| Record | Fingerprint and Pelican details, with existing account key and record ID lookup |
| Clear | `candy`, `fingerprint`, `pelican`, or `all`, with existing failed-save rollback |
| Cancel | Exact native `run_id`, with active stream closure and shared busy-state cleanup |

Manual run contexts remain independent of the completed start HTTP request.
Inspection runs remain children of the existing inspection context.
Keep `active[account]` as the single busy-state owner across all three test types.

## Detector

Replace the old probe fingerprint implementation and ModelTrace with the upstream detector.
Select three distinct challenges from the 36-challenge suite for each run.
Port the upstream number parsing, sequence checks, ranking, aliases, and probability calibration without changing their meaning.
Rank against the 60 enrolled candidates.
After shape acceptance, three count-valid answers permit probabilities.
One or two permit partial rankings without probabilities.
Use native detector statuses `completed`, `partial`, `failed`, and `cancelled`.
Keep native `busy` and `unavailable` start results.
Do not import CPA model-catalog admission or create its skipped records.
An empty model catalog must not block a literal custom model.

The sequence check must pass before either result can rank candidates.
Its failure returns `attribution.decision: "not_model"` without rankings or probabilities.
This attribution decision remains separate from the run's completion status.
Numeric-count acceptance can retry. Shape rejection does not retry.
Keep nullable probability fields. Do not replace null with zero or retain ModelTrace family probability aggregates.

The full decoder, ranker, and collector contract is in [detector code research](research/code-detector.md).
Replace `fingerprint_algorithm.go` with `fingerprint_detector.go` and remove both ModelTrace Go files.
Move the existing `mustDecode` helper into `codex.go` before removing its old owner.
Codex context, enrollment-suite, and detector decoding share that helper.

Use the upstream per-account concurrency setting, bounded by the three challenges actually available.
Keep six global challenge slots. A challenge retains its slot during retries.
Each run normally sends three requests and can send at most nine with retries.
Keep samples, attempts, diagnostics, and full rankings in record details as upstream does.
Keep summary responses small through the existing summary/detail boundary.

### Collector and Algorithm Contract

`runFingerprint` selects three distinct suite entries before its workers start.
Each worker acquires one global slot through a cancellation-aware select.
It retains that slot through `collectFingerprintSample`, all retries, and both retry waits.
It releases the slot after the logical challenge returns.
Collect indexed events and assemble samples in original challenge order.
Accumulate usage from every event, including failed attempts.

`fingerprintText` limits text to 64,000 bytes and parses the longest numeric run.
It accepts values from 1 through 355.
Its minimum count is `max(80, ceil(expected * 0.55))`.
It truncates excess numbers to the expected count before analysis.
Only this acceptance check belongs in the retry validator.

`analyzeFingerprint` then runs `sequenceShape` across count-valid answers.
The median lag-one autocorrelation must be at most `-0.23`.
The median share of steps no larger than 20 must be at most `0.05`.
Retain the complete upstream decoder, features, stable ranking, calibration, and math helpers.
Do not combine these helpers with the retired ModelTrace algorithm.

| Count-valid answers | Shape | Attribution | Run status after normal collection |
| --- | --- | --- | --- |
| None | Not evaluated | No attribution | `failed` |
| One or two | Fails | `not_model`, without ranking or probabilities | `partial` |
| Three | Fails | `not_model`, without ranking or probabilities | `completed` |
| One or two | Passes | Partial ranking, with null probabilities | `partial` |
| Three | Passes | Ranking with calibrated probabilities | `completed` |

Explicit failure and cancellation take precedence over attribution in the page.
If cancellation occurs with fewer than three completed challenge events, record `cancelled`.
Do not replace this status with ordinary `partial` because an earlier sample permits attribution.

Use native `fingerprintResult` metadata and the upstream sample and attribution fields.
Keep `Info`, record ID, timestamp, literal requested model, duration, status, error, and observed usage.
Add `detector_revision`, concurrency, samples, and optional attribution.
Each sample retains instructions, prompt, expected count, attempts, text, error, parsed count, and acceptance.
Summaries clear samples, full candidate rows, and diagnostics. Details retain them.
Existing non-omitting JSON tags serialize cleared slices as null, so the page must accept that summary shape.

Keep the same comparison aliases in Go and JavaScript:
`gpt-6.1-sol` maps to `gpt-6-astra`.
`deepseek-flash` maps to `deepseek-v4.1-flash-1017`.
Normalize only for comparison. Send the original model string.

Use the upstream resources from `Ikaleio/lm-detector` revision `654c8334bb53951aaf04cbd61252d163649e1369`.
The detector uses only the existing Go module and standard library.
Resources total about 4.55 MB. Decoded floats occupy about 8.2 MiB.
Keep golden data under `tests/backend/`, outside production patches.

| Resource | Patch 0002 destination under `data/` | Bytes |
| --- | --- | ---: |
| Captured context | `codex_context.json` | 70,018 |
| Challenge suite | `lm-detector/enrollment-suite.json` | 152,315 |
| Candidate metadata | `lm-detector/models.json` | 9,921 |
| Frozen detector | `lm-detector/shared_detector.bin` | 4,385,348 |
| Runtime license | `lm-detector/LICENSE` | 1,064 |
| Source reference | `lm-detector/README.md` | 1,227 |

Keep `shared_detector.bin` as a Git binary hunk. The existing build already uses `git apply --binary`.
Remove `fingerprint_baselines.json`, `fingerprint_probes.json`, and `data/modeltrace/`.
Keep the 54,780-byte golden fixture under `tests/backend/internal/codexeval/testdata/` only.
The decoder reconstructs byte-plane float32 data and exposes the upstream float64 vectors.
Keep the upstream decoder checks and revision metadata. Add no separate resource-integrity gate.

## Retry Runner

Use one shared bounded runner for Candy, fingerprint challenges, and Pelican requests.
Build the request once and send at most three attempts.
Wait two seconds before attempt two and four seconds before attempt three.
Accumulate observed usage and actual attempts across the logical request.
Save optional attempt fields only on new records. Do not rewrite existing Candy or Pelican record fields.

Move `askCodex` from `codex.go` to `model.go` with this native signature:

```go
func askCodex(
    ctx context.Context,
    auth credential,
    model, effort, instructions, prompt string,
    limit time.Duration,
    accept func(string) error,
    beforeSend func() error,
) (modelResponse, int, error)
```

Call `codexTurn` once before the loop. Reuse the same request identity, body, and headers across attempts.
Keep `executeModel` as the owner of response reading, usage extraction, and body closure.
Remove `executeProbe` after the retired callers disappear.

Use a fresh three-minute timeout for each fingerprint attempt.
Use a fresh ten-minute timeout for each Candy or Pelican attempt.
Parent cancellation and existing inspection deadlines remain effective throughout the runner.
Measure record duration around the whole logical request, including retry waits.

Candy supplies no answer validator. Its final text still goes through `hasStandalone21` once.
Fingerprint supplies only the numeric-count validator.
Pelican supplies the existing HTML extraction and size rules as its validator.
Do not tighten the accepted HTML format.

| Outcome | Action |
| --- | --- |
| Marked direct transport failure, eligible timeout, or HTTP server error | Retry while the run context remains active |
| Plugin error before any HTTP response | Keep native no-replay behavior and stop |
| Empty, failed, incomplete, or tool-only stream | Retry while the run context remains active |
| Fingerprint answer rejected by upstream parsing limits | Retry |
| Pelican answer without acceptable HTML | Retry |
| Candy text with a wrong answer | Complete after one attempt |
| HTTP 401, HTTP 429, or other permanent HTTP failure | Stop |
| Native account lookup, credential, or request setup failure | Retain the native error without a retry |
| Cancellation | Stop the wait, active stream, and later sends |
| Inspection participation becomes false | Stop later sends and retain the unchanged account outcome |

Keep HTTP status and errors distinct at the native sender boundary.
Use existing typed errors and context causes where available. Do not classify errors through message text.
Add no generic retries around account updates, settings saves, or notifications.

Keep retry classification inside `executeModel` and the shared runner.
Preserve `context.AfterFunc`, deferred body closure, and the native cancel-on-close response wrapper.
Cancel each attempt context as soon as that attempt returns.
The runner must release a response before it waits or starts another attempt.

### Observed Status and Transport Boundary

`doOpenAIAccountTestUpstream` still performs host validation and client setup through `HTTPUpstream`.
It is not a sufficient boundary for marking every returned error as retryable.
Keep proxy parsing, host validation, client creation, credentials, and request conversion errors final.
Mark actual transport calls at `repository.doUpstreamRequest`, after setup has completed.
Keep the underlying error available through `Unwrap`.
Use a small marker interface in `codexeval` to avoid a dependency on the service package.
Define `HTTPUpstreamRequestError` beside `HTTPUpstream` in `service/http_upstream_port.go`.
Give it `Error`, `Unwrap`, and `RequestAttempted() bool` methods.
Wrap only the actual `servertiming.Do` error in `repository/http_upstream.go`.
These small native transport hunks belong in patch `0004`.
Do not change the `HTTPUpstream` method signatures or its existing setup checks.
`executeModel` recognizes `RequestAttempted() bool` through `errors.As`.
No second sender-error wrapper is required.

Preserve native plugin no-replay behavior before a valid HTTP response exists.
`PluginTransportError.RequestSent` describes possible sending, rather than error transience.
Do not use it as permission to replay a request.
The current plugin error loses the structured RPC cause, so it cannot identify pre-send transience exactly.
Keep all pre-response plugin errors final. Do not expand the plugin error protocol for this task.
After valid HTTP headers, the approved HTTP-status and unusable-stream rules apply.
Plugin timeout normalization can return a plain context error without the possible-send bit.
Do not use an attempt deadline alone to override this pre-response no-replay boundary.
Retry timeouts only at the marked direct transport boundary or after a response starts.
Parent cancellation always stops the runner.

When the sender cannot read a non-200 body, return the observed response together with the read error.
Retain current body closure, redaction, and status-specific account handling.
In `executeModel`, capture any response status and arrange body closure before processing a sender error.
Keep HTTP 401 and 429 final even when their bodies fail to read.
Use the observed 408 or 5xx status for their existing retry category.
For status 200, classify unusable stream errors through the shared runner.
Never replace a known response status with zero because body reading failed.

## Candy Inspection

Automatic inspection and immediate inspection remain Candy-only.
One logical Candy test creates one final result after all internal attempts.
Only completed answers enter the existing confirmation logic.
A valid wrong answer counts once and does not trigger request retries.
Check participation before every supplier send, including after retry backoff.
Intentional opt-out must not become a synthetic failed evaluation.

Keep the participation closure in `inspection.testPhase` and pass it through the Candy worker to `askCodex`.
Call it after retry backoff and before each send, without holding the global evaluation mutex.
Use one internal stop sentinel for intentional opt-out.
Use `evaluateCandyChecked(..., check) (candyResult, bool)` to distinguish a result from an intentional stop.
The worker appends history and increments `Done` only when the second return value permits a save.
Ordinary lookup, account-state, and priority-change errors remain visible through the existing inspection result owner.
Do not cancel the round context when one account opts out.

Add `Attempts int` with `json:"attempts,omitempty"` to new Candy and Pelican results.
Keep every existing result field and old-record decoder.
Do not introduce a shared record rewrite merely to match CPA embedding.
The inspection collector still counts saved logical results, never HTTP attempts.
Keep the existing 30-minute inspection phase wait limit and its active-run cancellation path.
The phase limit can stop a later attempt before its individual ten-minute timeout.
Do not extend inspection timing to guarantee three full attempts.

Keep the current saved-verdict reread, transition confirmations, and recovery wait.
After each detection establishes the account state, apply existing model control and optional priority handling.
Retain stable-state handling, atomic model-map updates, the `blacklist:blacklist` placeholder, recorded removals, and recovery rules.
Manual fingerprint results do not change account verdicts, model mappings, or priority.

## Native History

Keep `<data-dir>/codex-eval/history.json` and native account ID keys.
Retain the `results` and `pelicans` maps, their record formats, and the latest-20 retention rule.
Use `fingerprints` for new detector records. Remove the `modeltraces` runtime and saved-state owner.

Add native root field `schema_version: 1` to distinguish new fingerprint records from retired records.
An unversioned file loads only Candy and Pelican records, then immediately saves the new versioned snapshot.
Its new fingerprint history starts empty.
Decode the fingerprint payload only after the marker identifies the new format.
Do not decode old records as the new detector type.

Use a small read envelope with typed `results` and `pelicans`, plus raw `fingerprints` and the root marker.
Decode the raw fingerprint field into current types only when `schema_version` is one.
For unversioned data, ignore both retired payloads without per-record inspection.
Assign loaded maps only after all retained fields decode successfully.
`currentState` writes the marker and the three current maps.
Remove `ModelTraces` from persistence, retention, detail lookup, clear, and rollback snapshots.
After successful local decoding, assign retained maps and call the existing `saveStateLocked` once for an unversioned file.
This minimal migration writes the marker and current three-map snapshot without retired records.
Keep Candy and Pelican payload fields unchanged. Do not synthesize IDs, statuses, or attempts for old records.
On save failure, retain loaded memory records and expose the existing storage error. The atomic saver protects the original file.
Version-1 files need no migration write. Missing or invalid files follow the existing load behavior.
Do not copy the CPA migration, rename `results` to `candies`, or create a history backup.

Keep existing unreadable, malformed, and null history protection.
Keep same-directory temporary writes, mode `0600`, atomic rename, visible save errors, and clear rollback.
The accepted history loss is irreversible after a successful save removes the retired records.
Installer executable backups do not restore evaluation history.

## Page and Preferences

Show Candy, fingerprint detection, Pelican, and Settings in the current native layout.
Replace the fingerprint mode control with the upstream per-account concurrency control.
Retain the native model dropdown and custom-model interaction.
Render detector progress, result status, attempts, rankings, probabilities, and record details from the new data.
Preserve account selection, groups, enable controls, degraded tags, and shared busy state.

Keep the Settings-first initialization owner and shared group refresh.
Preserve Candy and Pelican preferences. Stop reading retired fingerprint-mode and ModelTrace preferences.
A saved ModelTrace tab falls back to Candy. The existing fingerprint tab opens the new detector.
Add no explanatory small print, principle panel, or tooltips.
Keep actual errors, field labels, units, and result data.
Inspect the actual rendered page and exercise the changed controls.

Use the exact native symbol and markup replacement map in [UI code research](research/code-ui-history.md).
Retain the current page shell, account table, shared model picker, and Settings handlers.
Change only test-specific controls, progress, result rows, and detail views.
Show `not_model` as `不像语言模型回答`.
Show null probability as `—`, without a numeric probability or fabricated percentage.
Keep Candy and Pelican result details and add attempts only where the new result supplies them.

| Frontend owner | Concrete change |
| --- | --- |
| `tabNames`, `evaluationState`, `stateCredentials`, and `BUSY_TESTS` | Remove ModelTrace entries and projections. Keep native `active` and account controls. |
| `fpSavePrefs`, `runFingerprint`, and `initializeFingerprint` | Keep `fp-` IDs, model selection, and concurrency. Remove mode preferences and the `mode` request field. |
| `fpOutcome`, `renderFingerprintRow`, and `showFingerprintDetail` | Read the detector decision, samples, nullable probabilities, aliases, and status precedence. |
| All `mt*` handlers and projections | Remove selection, history, detail, clear bindings, preferences, and busy labels. |
| `credentialRow`, `historyEntry`, `historyPanel`, and collection helpers | Reuse the native shared components. |
| `start()` and `switchTab()` | Preserve Settings-first sequencing and later settings refresh. |
| `web.go` and `ui.html` | Remove ModelTrace embeds, config, script/license placeholders, panel, and dialog. Replace old fingerprint mode config. |
| `style.css` | Remove only retired test-specific styles. Retain shared layout and result styles. |

Keep started-run JSON field `id` and saved record field `run_id` distinct.
Continue using the existing record wrapper `{type, fingerprint}` or `{type, pelican}` for detail responses.

## Documentation and Distribution

Update Chinese README usage for the three test pages and the new fingerprint controls.
Keep installation, manual restoration, inspection, build prerequisites, build commands, output paths, and Credits.
Keep product decisions, implementation details, version descriptions, and upgrade sections outside the user README.

Replace unused ModelTrace resources, test data, license, and attribution with the detector's MIT material.
Update the sole third-party notice and Chinese Credits together.
The existing build includes required `licenses/*.txt` files in the installer.
Keep the native Linux/amd64 server and self-extracting installer outputs and release asset path.

Rename `licenses/ModelTrace-MIT.txt` for the detector instead of adding another identical MIT text.
Update `licenses/THIRD_PARTY_NOTICES.txt` and Chinese Credits to the selected `Ikaleio/lm-detector` source.
Keep the separate embedded license used by the page.
`scripts/build.py` already packages the root license files.
Keep `.github/workflows/release.yml`, installer sources, vendor patch `0003`, and native integration patch `0007` unchanged.

Remove the README's old 60/200/400-request table and unsupported 14,000-input-token estimate.
Describe three challenges per detector run, with retries where needed.
Any timeout guidance must state the per-attempt limit, rather than promise a full-run duration.
Keep the existing Build section and its commands.

Update current specs after implementation passes its focused checks:

| Spec | Required update |
| --- | --- |
| `runtime-installation.md` | Three tests, fingerprint request/progress/results, native summary/detail fields, and retained cancellation |
| `database-guidelines.md` | Root marker, raw fingerprint read boundary, three maps, and retained Candy/Pelican records |
| `error-handling.md` | Shared bounded runner, setup/transport distinction, final 401/429, validators, and cancellation |
| `native-inspection.md` | Callback before each attempt and unchanged confirmation/model-control timing |
| `directory-structure.md` and `index.md` | Current detector assets, retired owners, and test/page/map counts |

Keep build and installer specs unchanged unless the actual delivery path changes.

## Git Integration and Recovery

Create `sync/upstream-main-2026-10-19` from local `main` after final plan approval.
Refresh upstream and fork refs before the merge.
The researched upstream tip is `7a49068`.
If refreshed refs add commits, inspect their changes and reconcile this plan before integration.
The execution fetch found `e120868` and `4f83676` after the reviewed tip.
Defer their new Pelican variant and schema-2 changes because the approved scope keeps Pelican unchanged.
Merge `7a49068` for this task.
Use a normal two-parent merge with `--no-commit` so conflict resolutions and native adaptations receive one review.
Resolve CPA delete/modify conflicts by native ownership. Review renamed tests and licenses by their actual consumers.
Do not use a blanket `ours` merge or rewrite signed local history.

Keep the merge pending until checks pass and the user approves the exact merge commit files, message, and signing behavior.
Run Trellis finish after work commits. Include archive and session commits when fast-forwarding local `main`.
Recheck the fork before a normal `git push -u origin main`.
If the fork advanced, integrate that work before publication. Do not force-push.
Do not create a release as part of this task.

Before aborting a pending merge, retain any task edits needed for later work.
Use the recorded pre-merge state to protect unrelated work.
After commit, use a normal corrective commit for defects. Preserve the installer backup and restoration route for executable replacement.
