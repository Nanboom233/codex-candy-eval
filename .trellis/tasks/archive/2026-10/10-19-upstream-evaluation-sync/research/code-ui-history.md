# Research: Native Detector UI and History Adapter

- Query: Map the upstream fingerprint result contract, native page state, model catalog boundary, and history adapter for the approved detector replacement.
- Scope: mixed
- Date: 2026-10-19

## Findings

### Upstream page and result contract

Upstream exposes three test pages: Candy, fingerprint detection, and Pelican. It removes the ModelTrace page and the old fingerprint mode control (`7a49068:internal/plugin/web/app.js:9-12`; `7a49068:internal/plugin/web/ui.html:28-32,67-85`). A saved ModelTrace tab is not valid and falls back to Candy (`7a49068:internal/plugin/web/app.js:324-333,533-535`). The native page must keep its fourth Settings tab and native group, account, authentication, and run-ID controls.

The upstream fingerprint record stores the detector revision, concurrency, three challenge samples, and optional attribution. Summary responses remove samples and full ranking diagnostics. Detail requests return the complete record (`7a49068:internal/plugin/fingerprint.go:85-104`; `7a49068:internal/plugin/state.go:146-169`). The native State/Record handlers must keep this summary/detail split while using numeric Sub2API account IDs.

The detector decision is separate from the run status:

- A count-valid answer can produce `attribution.decision: "not_model"`. It has no prediction, ranking, or probability. The run can still be `completed` or `partial`.
- One or two accepted, model-like answers produce a ranking with `probability: null`. Null must stay null. Do not render it as `0` or add a probability bar value to the text.
- Three accepted, model-like answers produce rankings with calibrated probabilities.
- Failed analysis with no accepted answer produces `failed`.

The UI maps `not_model` to a warning state and the text `不像语言模型回答`. Its detail says that the number transitions are not model-like and that no ranking was produced. For partial rankings, show the prediction and ranking but omit the probability value. Candidate rows may use zero width for a null bar, but the displayed value remains `—` (`7a49068:internal/plugin/web/fingerprint.js:7-37,63-83`; `7a49068:internal/plugin/web/components.js:52-56`).

The server and page compare model names after trimming, lower-casing, removing the last provider prefix, and applying these aliases: `gpt-6.1-sol` to `gpt-6-astra`, and `deepseek-flash` to `deepseek-v4.1-flash-1017` (`7a49068:internal/plugin/fingerprint.go:61-74`; `7a49068:internal/plugin/web/fingerprint.js:7-19`). The alias changes comparison only. The native sender must send the administrator's literal model unchanged.

Run completion and cancellation have a required precedence. The upstream runner starts as `completed`, counts `done` for each non-cancelled challenge event, and changes to `partial` only when attribution succeeds with fewer than three accepted answers. Its deferred cleanup changes the status to `cancelled` when the context is cancelled and `done < 3` (`7a49068:internal/plugin/fingerprint.go:135-152,185-217`). Therefore a cancelled run with fewer than three completed challenge events must not be shown as a normal partial result. A run that receives all three challenge events can remain `completed` or `partial`, even when one or more answers are invalid. The row renderer must give explicit `failed` and `cancelled` statuses precedence over attribution (`7a49068:internal/plugin/web/fingerprint.js:22-37`).

### Native history discriminator

The current native file stores `results`, `fingerprints`, `modeltraces`, and `pelicans` in `<data-dir>/codex-eval/history.json`. It decodes all four maps directly, returns summaries from `State`, serves fingerprint and ModelTrace details from `Record`, and rolls all maps back when a clear save fails (`patches/0001-evaluation-core.patch:2195-2341`). The old fingerprint and ModelTrace Go types cannot decode the new detector records as the same contract.

Use a native root marker `schema_version: 1` and a small raw envelope for reading. The marker separates an actually present unversioned native file, which contains retired fingerprint data, from a file written by the new detector. For an unversioned file, decode only the existing Candy `results` and Pelican `pelicans` maps. Leave the old `fingerprints` and `modeltraces` bytes untouched in memory and do not unmarshal them into the new detector type. For a version-1 file, decode the current Candy, fingerprint, and Pelican maps; only the fingerprint field needs a `json.RawMessage` read boundary before current typed decoding.

Do not use `detector_revision` as a per-record migration filter. It remains useful record metadata, but a per-record filter would require mixed-record handling and would silently accept a state that was never written by the native format. No extra future-schema or mixed-file matrix is required.

The final user decision requires minimal migration during load.
After successfully decoding an unversioned file, save retained Candy and Pelican records with `schema_version: 1` and an empty fingerprint map.
Use the existing atomic saver. Do not alter retained record fields or create a backup.
Version-1 files need no migration write. Existing malformed, null, and unreadable protection remains.
Save uses a same-directory temporary file, mode `0600`, and atomic rename.
A failed migration keeps loaded records in memory and exposes the save error. A failed clear restores all previous maps.
Sources: `patches/0001-evaluation-core.patch:2209-2289,2316-2341` and `.trellis/spec/backend/database-guidelines.md:60-71`.

After replacement, `State` returns `active`, `candy`, `fingerprint`, `pelican`, and `storage_error`. `Record` serves only fingerprint and Pelican details. `Clear` accepts Candy, fingerprint, Pelican, and all scopes. The retired ModelTrace scope and map must disappear. Keep the latest-20 retention rule and account-ID keys.

### Native admission and model catalog boundary

CPA's `startRuns` accepts a per-credential `model_catalog`. A known catalog entry without the requested model creates a skipped record; an unknown entry starts the run and increments `unchecked` (`7a49068:internal/plugin/runs.go:23-42,55-87`). This branch has no native equivalent and must not be copied.

Native `StartRequest` contains numeric `targets`, a literal `model`, effort or run settings, and concurrency. Native admission resolves each account, reports `unavailable` for invalid targets or resolver errors, and reports `busy` when any test already runs for that account (`patches/0001-evaluation-core.patch:1960-2057`). `SendCodexEvaluation` accepts the literal model and sends it through the native account transport (`patches/0004-admin-channel.patch:211-271`). Do not add `model_catalog`, credential-ID admission, `skipped`, or `unchecked` fields to native requests or summaries.

The native catalog is a suggestion source, not an authorization check. It reads current group membership, discovers raw evaluation model names, and keeps the account inventory ready even when model discovery returns an error (`patches/0006-evaluation-ui.patch:771-820,833-846`). A current inventory plus an empty or failed model suggestion list must keep the custom literal-model path usable. An actual inventory failure still blocks targets until a fresh group read succeeds. Preserve this existing distinction and do not add a new per-account model prerequisite. Existing model-picker code retains a selected literal when suggestions disappear and supports an empty catalog (`patches/0006-evaluation-ui.patch:901-930`; `tests/ui/sub2api-dropdowns.cjs:117-140,709-715`).

### Native page state and Settings-first loading

The current page keeps `tabNames` for Candy, old fingerprint, ModelTrace, Pelican, and Settings, with separate result maps and preference keys (`patches/0006-evaluation-ui.patch:7-29`). Replace the old fingerprint and ModelTrace owners with the new fingerprint page, remove ModelTrace state and preferences, and keep the Settings tab. A saved ModelTrace tab must open Candy. Keep `evaluationState.active` as the source for the native active run, then match by numeric account ID and group ID. Cancel must continue to post the exact native `run_id` (`patches/0006-evaluation-ui.patch:174-211,274-297,435-455`).

The initial Settings path is already serialized: `start()` reads inspection settings before it refreshes groups and the catalog, while later Settings entries refresh settings through `switchTab` (`patches/0006-evaluation-ui.patch:324-332,342-352`; `tests/ui/sub2api-dropdowns.cjs:87-115`). Keep this owner when replacing page scripts. Do not start a second settings read when the saved tab is Settings. Preserve native account filters, enable switches, degraded tags, model control, priority, and the no-small-print UI.

### Focused checks and exclusions

Adapt the existing UI fixture from four evaluation lists to Candy, fingerprint, Pelican, and Settings. Keep the Settings-first group-read assertion, custom literal model coverage, native active-run and run-ID cancellation coverage, summary/detail loading, clear rollback, and storage-error display. Add only focused assertions for the changed contracts: `not_model`, null probability, alias comparison, cancelled status when `done < 3`, unversioned history loading, version-1 save shape, and removal of ModelTrace state. Do not copy CPA catalog, skipped, credential-ID cancellation, or full CPA state-migration suites.

Retain `mustDecode` under its live Codex context owner when removing the old fingerprint algorithm. The captured Codex context still uses it, and the new detector resources also need the helper (`patches/0001-evaluation-core.patch:536-593`; `patches/0001-evaluation-core.patch:145-155`).

## Files Found

- `patches/0001-evaluation-core.patch` — native task admission, result types, State, Record, Clear, and file persistence.
- `patches/0004-admin-channel.patch` — native account lookup, literal model sender, authentication, proxy, and stream boundary.
- `patches/0006-evaluation-ui.patch` — current tabs, model picker, catalog readiness, Settings-first startup, and retired result pages.
- `tests/backend/internal/codexeval/runtime_test.go` — native busy, cancellation, retention, and corrupt-history coverage.
- `tests/ui/sub2api-dropdowns.cjs` — group loading, literal model fallback, Settings-first startup, and native account controls.
- `.trellis/spec/backend/runtime-installation.md` — native task, route, state, catalog, and page contracts.
- `.trellis/spec/backend/database-guidelines.md` — native history path, retention, atomic save, and clear rollback.
- `.trellis/spec/backend/native-inspection.md` — Settings, model control, priority, degraded state, and recovery contracts.
- Upstream `7a49068:internal/plugin/fingerprint.go` — detector records, collector status, and cancellation precedence.
- Upstream `7a49068:internal/plugin/fingerprint_detector.go` — decision, ranking, and nullable probability fields.
- Upstream `7a49068:internal/plugin/web/fingerprint.js` — result mapping, aliases, detail rendering, and attempts display.
- Upstream `7a49068:internal/plugin/state.go` and `migration.go` — upstream state shape, which is not a native migration contract.
- Upstream `7a49068:internal/plugin/runs.go` — CPA catalog admission and skipped branch, excluded from native integration.

## Code Patterns

- Native result summaries use a typed `summary` method and detail lookups under the records lock (`patches/0001-evaluation-core.patch:2148-2153,2291-2314`).
- Native clear operations snapshot all maps before mutation and restore them after a failed atomic save (`patches/0001-evaluation-core.patch:2316-2341`).
- Native catalog readiness is group-specific. Model discovery errors do not invalidate a current account inventory (`patches/0006-evaluation-ui.patch:793-820,833-846`).
- Upstream detector status is based on logical challenge completion, while attribution decision and accepted-answer count are separate fields (`7a49068:internal/plugin/fingerprint.go:185-217`; `7a49068:internal/plugin/fingerprint_detector.go:186-199,233-275`).
- Upstream UI compares aliases only for the displayed consistency label. It keeps the requested model in the record and request (`7a49068:internal/plugin/fingerprint.go:61-74`; `7a49068:internal/plugin/web/fingerprint.js:7-19`).

## External References

- Upstream source: `haowang02/cpa-plugin-codex-candy-eval`, commit `7a49068`.
- Detector source and data: `Ikaleio/lm-detector`, pinned by the upstream detector resource metadata.

## Related Specs

- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/database-guidelines.md`
- `.trellis/spec/backend/native-inspection.md`
- `.trellis/spec/guides/cross-layer-thinking-guide.md`
- `.trellis/tasks/10-19-upstream-evaluation-sync/research/code-detector.md`
- `.trellis/tasks/10-19-upstream-evaluation-sync/research/code-delivery.md`

## Caveats / Not Found

- This report reads source, patches, and tests only. It does not edit product code or run tests, builds, installations, or provider calls.
- The upstream CPA migration backup and credential model catalog do not apply to the native Sub2API history or account contracts.
- The native version marker is an adapter boundary for the accepted breaking history change. It is not a general migration framework.
- Detector algorithm thresholds and resource parity remain in `research/code-detector.md` and are not repeated here.

## Implementation Appendix

### Native request, progress, and record shapes

Keep the existing native request and run ownership. Remove only the retired `Mode` field and the `fingerprint` mode switch; use the selected model and concurrency as the detector inputs. A focused declaration shape is:

```go
type StartRequest struct {
    Type        string   `json:"type"`
    Targets     []Target `json:"targets"`
    Model       string   `json:"model"`
    Effort      string   `json:"effort"`
    Runs        int      `json:"runs"`
    Concurrency int      `json:"concurrency"`
}

type fingerprintChallenge struct {
    Instructions  string `json:"instructions,omitempty"`
    Prompt        string `json:"prompt"`
    ExpectedCount int    `json:"expected_count"`
}

type fingerprintSample struct {
    fingerprintChallenge
    Attempts int    `json:"attempts,omitempty"`
    Text     string `json:"text,omitempty"`
    Error    string `json:"error,omitempty"`
    Parsed   int    `json:"parsed_numbers"`
    Valid    bool   `json:"accepted"`
}

type fingerprintResult struct {
    Info
    DetectorRevision string                  `json:"detector_revision"`
    ID               string                  `json:"id"`
    Time             time.Time               `json:"time"`
    Model            string                  `json:"model"`
    Status           string                  `json:"status"`
    Error            string                  `json:"error,omitempty"`
    DurationMS       int64                   `json:"duration_ms"`
    InputTokens      int64                   `json:"input_tokens"`
    OutputTokens     int64                   `json:"output_tokens"`
    ReasoningTokens  int64                   `json:"reasoning_tokens"`
    Concurrency      int                     `json:"concurrency,omitempty"`
    Samples          []fingerprintSample     `json:"samples"`
    Attribution      *fingerprintAttribution `json:"attribution,omitempty"`
}

type fingerprintProgress struct {
    Model       string `json:"model"`
    Done        int    `json:"done"`
    Total       int    `json:"total"` // 3 selected challenges
    Concurrency int    `json:"concurrency"`
    cancel      context.CancelFunc
}
```

`fingerprintAttribution` follows the upstream JSON contract: `decision`, `prediction`, `prediction_name`, nullable `probability`, `used_outputs`, `family_prediction_name`, optional `sequence`, `results`, and `diagnostics`. Candidate `probability` is `*float64`; do not coerce null to zero. `fingerprintSequence` reports `autocorrelation`, `small_step_share`, and `model_like`. Keep the native embedded `Info` fields (`group_id`, numeric `account_id`, `account_name`, and `run_id`); do not replace them with CPA credential identity. `StartRequest` remains the shared Candy/fingerprint/Pelican request. Candy keeps `runs` and `effort`; Pelican keeps its existing fields. Only fingerprint dispatch uses `Concurrency`.

Use the following native-owner signatures; the detector internals retain the upstream challenge and attribution contracts:

```go
func fingerprintChallenges() []fingerprintChallenge
func collectFingerprintSample(context.Context, credential, string, fingerprintChallenge) (fingerprintSample, modelResponse, error)
func analyzeFingerprint([]fingerprintOutput) (fingerprintAttribution, error)
func runFingerprint(context.Context, credential, string, int, *fingerprintProgress)
```

The `runFingerprint` context remains the native active-run context. `Cancel` still matches the exact `run_id`; the runner sets `cancelled` when the context ends before all three logical challenge events complete. It saves one result through the existing latest-20 history helper. Do not bring upstream `skipped` or catalog admission into this request path.

### History adapter pseudocode and JSON

Read the root marker before typed fingerprint decoding. An absent marker means the existing native file: load only Candy and Pelican; leave old fingerprint and ModelTrace fields unclaimed. Marker version 1 means decode the current fingerprint map after the raw-field boundary. A normal save writes one current snapshot with the marker and three maps:

```go
var envelope struct {
    SchemaVersion json.RawMessage `json:"schema_version"`
    Results       json.RawMessage `json:"results"`
    Fingerprints  json.RawMessage `json:"fingerprints"`
    Pelicans      json.RawMessage `json:"pelicans"`
}

var candy, fingerprints, pelicans map[string][]Record
var version *int
if len(envelope.SchemaVersion) > 0 {
    if err = json.Unmarshal(envelope.SchemaVersion, &version); err != nil || version == nil || *version != 1 {
        return protectedLoadError
    }
}
switch {
case len(envelope.SchemaVersion) == 0:
    err = decode(envelope.Results, &candy)
    if err == nil { err = decode(envelope.Pelicans, &pelicans) }
case *version == 1:
    err = decode(envelope.Results, &candy)
    if err == nil { err = decode(envelope.Fingerprints, &fingerprints) }
    if err == nil { err = decode(envelope.Pelicans, &pelicans) }
default:
    return protectedLoadError
}
if err != nil { return protectedLoadError }
candyResults, pelicanResults = candy, pelicans
if len(envelope.SchemaVersion) > 0 { fingerprintResults = fingerprints }
if len(envelope.SchemaVersion) == 0 { return saveStateLocked() }
```

The snippet uses `Record` as a placeholder for each map's native slice type.
Assign globals only after all retained fields decode successfully.
A missing marker selects the minimal migration. A present null or unsupported value remains an invalid format.
After latest-20 trimming, write the new snapshot through `saveStateLocked` while holding the existing lock.
This writes retained Candy/Pelican payloads unchanged and omits retired data. Do not add backups or per-record conversion.
Preserve existing load protection, atomic rename, `0600` mode, and clear rollback.

Existing unversioned native history has this relevant shape (old fingerprint records are intentionally ignored):

```json
{
  "results": {"42": [{"status": "completed", "time": "...", "model": "gpt-6.1-sol"}]},
  "fingerprints": {"42": [{"mode": "standard", "attribution": {"status": "consistent"}}]},
  "modeltraces": {"42": [{"status": "completed"}]},
  "pelicans": {"42": [{"status": "completed", "time": "...", "model": "gpt-6.1-sol"}]}
}
```

The new saved shape is:

```json
{
  "schema_version": 1,
  "results": {"42": [{"status": "completed", "time": "...", "model": "gpt-6.1-sol"}]},
  "fingerprints": {"42": [{"id": "...", "time": "...", "group_id": 7, "account_id": 42, "account_name": "Example", "run_id": "...", "model": "gpt-6.1-sol", "status": "partial", "samples": [{"prompt": "challenge 1", "expected_count": 300, "attempts": 1, "parsed_numbers": 180, "accepted": true}, {"prompt": "challenge 2", "expected_count": 300, "attempts": 3, "parsed_numbers": 0, "accepted": false, "error": "有效数字不足"}, {"prompt": "challenge 3", "expected_count": 300, "attempts": 3, "parsed_numbers": 0, "accepted": false, "error": "有效数字不足"}], "attribution": {"decision": "partial", "probability": null, "used_outputs": 1, "results": [{"model": "gpt-6-astra", "probability": null}]}}]},
  "pelicans": {"42": [{"status": "completed", "time": "...", "model": "gpt-6.1-sol"}]}
}
```

The ranking array above is abbreviated to one candidate for readability; a stored attribution contains the full ranking. The saved partial record has three challenge samples, with one accepted answer. The State summary clears sample and ranking fields; these fields serialize as `null` under the native JSON tags, for example:

```json
{"fingerprint":{"42":[{"id":"...","group_id":7,"account_id":42,"account_name":"Example","run_id":"...","model":"gpt-6.1-sol","status":"partial","samples":null,"attribution":{"decision":"partial","probability":null,"used_outputs":1,"results":null,"diagnostics":null}}]}}
```

Retain the native `results` Candy key and Candy/Pelican record payloads. Do not rename them to CPA's `candies` or change `Info` identities. The sample is full detail data; the state response uses a summary copy with sample and full ranking fields cleared to `null`.

### Native API JSON boundary

Keep the existing native endpoints and response envelope (`response.Success`):

```json
POST /api/v1/admin/codex-eval/start
{"type":"fingerprint","targets":[{"group_id":7,"account_id":42}],"model":"openai/gpt-6.1-sol","concurrency":3}
```

The response `data` keeps native `started`, `busy`, and `unavailable` arrays. A started entry is a native `run` and exposes its identifier as `id`; saved result records retain `Info.run_id`. Polling exposes progress under `active`:

```json
{"data":{"active":[{"id":"run-id","type":"fingerprint","target":{"group_id":7,"account_id":42},"status":"running","progress":{"model":"openai/gpt-6.1-sol","done":1,"total":3,"concurrency":3}}],"candy":{"42":[]},"fingerprint":{"42":[]},"pelican":{"42":[]},"storage_error":""}}
```

After completion, the fingerprint summary in `state.fingerprint["42"]` clears samples and ranking diagnostics (serialized as `null` by these JSON tags) but retains status, decision, nullable probability, and prediction. Detail stays on the native `records/:type/:account_id/:record_id` endpoint:

```json
GET /api/v1/admin/codex-eval/records/fingerprint/42/record-id
{"data":{"type":"fingerprint","fingerprint":{"id":"record-id","account_id":42,"group_id":7,"run_id":"run-id","model":"openai/gpt-6.1-sol","status":"partial","samples":[{"prompt":"...","expected_count":300,"attempts":2,"parsed_numbers":180,"accepted":true}],"attribution":{"decision":"partial","probability":null,"used_outputs":1,"results":[]}}}}
```

`POST /api/v1/admin/codex-eval/clear-history` retains `{"scope":"fingerprint"}`. `State` removes `modeltrace`; `Record` accepts fingerprint and Pelican only; `Clear` accepts Candy, fingerprint, Pelican, and all. A successful all-clear saves the same schema-versioned three-map snapshot.

### Frontend symbol replacement map

| Current native symbol or owner | Change for the upstream detector |
| --- | --- |
| `tabNames` in `web/app.js` | Remove `modeltrace`; keep `candy`, `fingerprint`, `pelican`, `settings`. The saved fingerprint tab remains valid; a saved ModelTrace tab fails membership and falls back to Candy. |
| `evaluationState` | Keep `active`, `candy`, `fingerprint`, `pelican`; remove `modeltrace`. Update initial and logout/reset objects consistently. |
| `stateCredentials` | Keep `data.fingerprint[id]` as fingerprint history. Remove `data.modeltrace[id]` and ModelTrace running projections. Keep Candy/Pelican projections. |
| Prefix sets (`""`, `"fp-"`, `"mt-"`, `"pl-"`) and `BUSY_TESTS` | Remove only `mt-`/ModelTrace entries. The old native fingerprint page uses `fp-` DOM IDs; these can remain to limit native markup churn while the page behavior is replaced. |
| `fpSelected`, `fpExpanded`, `fpSavePrefs`, `fpOutcome`, `renderFingerprintRow`, `renderFingerprints`, `runFingerprint`, `initializeFingerprint` | Retain page/list ownership and shared collection actions. Remove `fp-mode`, old mode label/history metadata, `mode` request field, and mode preference read/write. Render upstream `decision`, alias comparison, partial ranking, and nullable probability. Use the detector's three-challenge progress and detail payload. Keep native `model` and concurrency controls. |
| `showFingerprintDetail`, old comparison/probe table renderer | Replace the old fingerprint statistics UI. Render upstream challenge prompts, accepted counts, retries, optional sequence decision, candidate ranking, nullable probabilities, and errors from the new record shape. Do not retain old probe-specific comparison tables. |
| `mtSelected`, `mtExpanded`, `mtSavePrefs`, `mtOutcome`, `renderModelTraceRow`, `renderModelTrace`, `runModelTrace`, `showModelTraceDetail`, `initializeModelTrace` | Delete. Also remove ModelTrace model/concurrency preference reads, record arrays, detail and clear event bindings, busy labels, page markup, result styles, and its type from `StartRequest` dispatch. |
| `credentialRow`, `historyEntry`, `historyPanel`, `collectionButton`, `bindCollectionActions`, shared state loading and run cancellation | Reuse. Keep numeric account IDs, group scoping, native `run_id` cancellation, filters, account enable switches, degraded tags, model catalog suggestions, and the literal custom-model path. |
| `start()` and `switchTab()` | Keep Settings-first initialization and the shared group refresh. Keep the fourth Settings tab; do not import the upstream three-tab-only navigation. |
| Preference keys under `PREF_STORE` | Continue using the fingerprint page's native preference owner for model and concurrency only; do not read a retired mode value. Keep Candy/Pelican preferences. Remove `.modeltrace` reads and writes. |

The native fingerprint module is not a drop-in for upstream `fingerprint.js`: its selected/expanded sets, detail functions, and outcome function use `fp*` names and the old eight-probe result schema. Either adapt those functions to the upstream record contract or replace their page module while retaining the native shared helpers. In either case, do not copy CPA account/catalog APIs or replace `credentialRow` and native layout.

### Embedded page changes

In `web/ui.html`, keep the native Settings panel, its model-control and priority controls, its account/group controls, and the existing fingerprint detail dialog location. Remove the ModelTrace tab/panel, `mt-*` detail dialog, `MODELTRACE_SCRIPT` placeholder, and ModelTrace-only explanatory content. Remove the old fingerprint mode select and old mode-dependent markup. The new fingerprint controls are the native model picker plus per-account concurrency; preserve the no-small-print constraint. Keep Candy and Pelican panels unchanged except for consuming the adapted shared state.

In `web.go`, remove the ModelTrace Go embed, config marshaling, injected license comment, and `MODELTRACE_SCRIPT` replacement. Keep the single fingerprint script/config and inject only values required by the detector UI (three requested challenges and model aliases if the page uses server-owned config). Do not retain old `fpModes`, `FP_CONFIG`, sampling-mode strings, or an unused modeltrace license variable. In the HTML template, remove the ModelTrace license placeholder and script element. Do not remove unrelated embedded markdown, Math, Candy, catalog, components, Pelican, or inspection assets.

### Removal and reuse boundary

Remove the old eight-probe statistical fingerprint implementation and its mode controls, and remove ModelTrace algorithm, UI, state, styling, license, and attribution owners. Add the upstream LM detector and its resources as the new fingerprint implementation. Do not migrate either retired result type. Keep Candy and Pelican records and native account IDs unchanged. Reuse `credentialRow`, common history, account controls, Settings-first loading, native authentication, and sender transport. The upstream native panel's explanatory principle text is not copied; keep only functional labels and rendered results.
