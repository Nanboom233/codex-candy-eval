# Research: Native Build, Patch, Test, and License Delivery

- Query: Map the upstream detector, retry, and request-context changes to the existing Sub2API patch series, test inputs, build selector, and redistributed notices.
- Scope: internal
- Date: 2026-10-19

## Findings

### Exact upstream resource inventory

The current upstream tree is `7a49068`. These are all files under its plugin `data/` and `testdata/` paths.

| Upstream path | Bytes | Use | Native destination |
| --- | ---: | --- | --- |
| `internal/plugin/data/codex_context.json` | 70,018 | Runtime request context. Embedded by `codex.go:36`. | `backend/internal/codexeval/data/codex_context.json` in patch 0002 |
| `internal/plugin/data/lm-detector/enrollment-suite.json` | 152,315 | Runtime challenge suite. Embedded by `fingerprint.go:31`. | `backend/internal/codexeval/data/lm-detector/enrollment-suite.json` in patch 0002 |
| `internal/plugin/data/lm-detector/models.json` | 9,921 | Runtime candidate names and families. Embedded by `fingerprint_detector.go:24`. | `backend/internal/codexeval/data/lm-detector/models.json` in patch 0002 |
| `internal/plugin/data/lm-detector/shared_detector.bin` | 4,385,348 | Runtime ranker and calibration. Embedded by `fingerprint_detector.go:21`. | `backend/internal/codexeval/data/lm-detector/shared_detector.bin` in patch 0002 |
| `internal/plugin/data/lm-detector/LICENSE` | 1,064 | Runtime license text. Embedded by `fingerprint_detector.go:27` and inserted into the page by `plugin.go:86`. | `backend/internal/codexeval/data/lm-detector/LICENSE` in patch 0002 |
| `internal/plugin/data/lm-detector/README.md` | 1,227 | Source and pin documentation. Not embedded. | Keep beside the detector data in patch 0002 because the upstream detector source comment points to it. |
| `internal/plugin/testdata/lm-detector-golden.json` | 54,780 | Test fixture. Read by upstream fingerprint tests from `testdata/`. | `tests/backend/internal/codexeval/testdata/lm-detector-golden.json`. `build.py` copies it into the temporary package. |

The three detector runtime data assets total 4,547,584 bytes. The detector code also embeds its license. The golden file is test-only and must not enter a production patch or installer payload.

### Patch ownership and complete delivery boundary

Keep the existing seven-patch order. Do not add a root Go module, CPA source overlay, or new delivery path.

| Owner | Affected paths or action |
| --- | --- |
| `patches/0001-evaluation-core.patch` | Update `backend/internal/codexeval/candy.go`, `codex.go`, `fingerprint.go`, `model.go`, and `pelican.go` for the approved request and retry behavior. Replace `fingerprint_algorithm.go` with the detector implementation, add `fingerprint_detector.go`, and remove `modeltrace.go` and `modeltrace_algorithm.go`. Update `runtime.go` and `state.go` only for the new result owner and retired ModelTrace type. |
| `patches/0002-evaluation-data.patch` | Update `data/codex_context.json`. Remove `data/fingerprint_baselines.json`, `data/fingerprint_probes.json`, and `data/modeltrace/{LICENSE,README.md,unified_bank.json}`. Add the five `data/lm-detector/` files in the inventory above. |
| `patches/0004-admin-channel.patch` | Preserve observed status in `backend/internal/service/account_codex_eval_service.go`. Add the small transport-attempt marker beside `HTTPUpstream` and wrap only the actual network call in `repository/http_upstream.go`. Keep setup errors final and request sending on the native service. |
| `patches/0005-inspection.patch` | Update `backend/internal/codexeval/inspection.go` only where the inspection participation check must run before every retry send. Keep confirmation counting and account control here. |
| `patches/0006-evaluation-ui.patch` | Update `backend/internal/codexeval/web/{app.js,candy.js,components.js,fingerprint.js,pelican.js,style.css,ui.html}` and `backend/internal/codexeval/web.go` for the three-page UI, attempt details, and detector license. Remove `backend/internal/codexeval/web/modeltrace.js`. Preserve the native Settings tab and `frontend/src/views/admin/CodexEvalView.vue` integration. Do not copy the CPA page or routes. |
| `patches/0003-vendor-assets.patch` | No upstream change requires edits. Keep the existing vendor assets and binary hunk. |
| `patches/0007-sub2api-integration.patch` | No route or sidebar change is required by these upstream commits. Keep native integration. |

The upstream package uses `//go:embed` for the context, suite, ranker, model list, and license. The native package path is `backend/internal/codexeval`, so each embed pattern must resolve there after all patches apply. Keep the golden file outside `backend/` until `scripts/build.py` copies `tests/backend/` into the check tree (`scripts/build.py:19-30`).

`shared_detector.bin` is a 4.4 MB binary. Patch 0002 can carry it as a Git binary patch. The established build applies patches with `git apply --binary` (`scripts/build.py:22-27`), and patch 0003 already carries a Git binary hunk for `texmath.js`. Do not hex-encode the detector into Go source or put the raw asset in a root build directory.

The detector code adds no Go module dependency. It uses standard-library `slices.Sorted` and `slices.Values`, which require Go 1.23 or later. The retained Sub2API v0.2.17 source tree declares `go 1.27.2` in `backend/go.mod`. The release workflow selects Go from that file (`.github/workflows/release.yml:22-24`). Do not copy the CPA root `go.mod` or its `golang.org/x/text` requirement.

### Focused build test selection

The current `scripts/build.py:53` selector is:

```text
Codex|FingerprintNormalization|FingerprintStatistics|FingerprintAttribution|ModelTraceUpstreamParity|ModelTraceParseAndThreshold|ExtractHTML|HasStandalone21
```

Replace the retired test names with this focused expression:

```text
Codex|FingerprintUpstreamParity|FingerprintParseAndChallenges|ExtractHTML|HasStandalone21
```

The only upstream fingerprint tests to adapt are `TestFingerprintUpstreamParity` and `TestFingerprintParseAndChallenges`. Keep CPA collector, retry, Candy, Pelican, and persistence suites as source references. Do not copy them into native tests. Extend native `TestCodexEvalRuntime`, `TestCodexEvalStreamSemantics`, and the existing sender and inspection tests. Add `tests/backend/internal/codexeval/model_test.go` with one `TestCodexEvalRetries` for shared retry behavior that those tests do not cover. The `Codex` fragment selects these native request, sender, inspection, participation, and retry tests. The final two fragments retain the existing HTML and Candy-answer checks.

The old fragments `FingerprintNormalization`, `FingerprintStatistics`, `FingerprintAttribution`, `ModelTraceUpstreamParity`, and `ModelTraceParseAndThreshold` describe algorithms being removed. Keep only the detector golden data needed by parity testing. Remove retired probe and ModelTrace data. Keep the existing UI check and build/vet commands in `scripts/build.py:51-54`. No new full-suite, race, or coverage gate is needed.

### License and root repository artifacts

Upstream `internal/plugin/data/lm-detector/LICENSE` and local `licenses/ModelTrace-MIT.txt` contain the same MIT text, including `Copyright (c) 2026 xqy2006`. Reuse and rename the existing root license file to an LM-detector name. Do not add a duplicate license text. Attribute the detector source to `Ikaleio/lm-detector` at `654c8334bb53951aaf04cbd61252d163649e1369`. Retain the copyright and terms from the embedded license. Replace the ModelTrace section in `licenses/THIRD_PARTY_NOTICES.txt` and its Chinese README credit with the selected detector attribution.

`prepare_installer` copies every root `licenses/*.txt` file into the installer (`scripts/build.py:71-75`). Renaming the file therefore changes the packaged path without a build-script change. The detector binary also embeds the license into the admin page. Retain this separate runtime use when removing the old ModelTrace page license.

The root work owners are patches `0001`, `0002`, `0004`, `0005`, and `0006`, affected tests, and `scripts/build.py`.
Documentation owners are the renamed license, the sole notice, and Chinese README.
Keep `.github/workflows/release.yml`, installer source, and the asset upload path unchanged.
The workflow uses `python3 scripts/build.py --source upstream --check` and publishes `sub2api-codex-eval-installer`.
Source: `.github/workflows/release.yml:34-39,48-55`.

### README usage and cost statements

The current README still describes the retired four-page workflow and old request estimates (`README.md:51-90`). Update it with the accepted three-page workflow:

- Keep Candy and Pelican instructions and their visible effort controls.
- Replace the old fingerprint sampling modes and 60/200/400 request table with model selection and per-account concurrency from 1 to 3.
- Remove the ModelTrace section and references.
- State the user-visible detector fact: each run selects three challenges. Do not call this a fixed request count because retries can send up to three attempts for each challenge (`7a49068:internal/plugin/fingerprint.go:17-21,122-129`). See `7a49068:internal/plugin/model.go:43-51` for the retry limit.
- Remove the claim that Candy, ModelTrace, and Pelican each use about 14,000 input tokens. The captured context changed and no current measurement supports that estimate. Do not replace it with a guessed number.
- Keep the latest-20 record statement only if the native history owner retains it. Keep stop and cancellation instructions.
- Keep the existing visible timeout guidance only after implementation confirms the native mapping. Upstream defines a three-minute limit for each fingerprint challenge attempt and a ten-minute limit for each Candy or Pelican attempt. See `7a49068:internal/plugin/model.go:43-45`, `7a49068:internal/plugin/fingerprint.go:220-223`, `7a49068:internal/plugin/candy.go:71`, and `7a49068:internal/plugin/pelican.go:70`. The limit is per attempt, not the full logical run. Retry waits add up to six seconds before the final attempt.

These edits are user documentation only. They do not require a new benchmark or token-cost test. The build and installer paths remain the existing commands and outputs in `README.md:169-196`.

## Files Found

- `.trellis/tasks/10-19-upstream-evaluation-sync/prd.md` — approved delivery, compatibility, documentation, and check requirements.
- `.trellis/tasks/10-19-upstream-evaluation-sync/design.md` — approved patch owners, resource pin, test approach, and release boundaries.
- `.trellis/tasks/10-19-upstream-evaluation-sync/research/upstream-evaluations.md` — algorithm and retry behavior research.
- `.trellis/tasks/10-19-upstream-evaluation-sync/research/upstream-ui-history.md` — UI and saved-state research.
- `patches/0001-evaluation-core.patch` — current native runtime and embedded request-context owner.
- `patches/0002-evaluation-data.patch` — current data and third-party asset patch.
- `patches/0003-vendor-assets.patch` — existing Git binary patch precedent.
- `patches/0004-admin-channel.patch` — native account sender boundary.
- `patches/0005-inspection.patch` — native inspection participation and confirmation boundary.
- `patches/0006-evaluation-ui.patch` — page source and Go embed declarations.
- `patches/0007-sub2api-integration.patch` — native routes and page integration.
- `scripts/build.py` — patch application, test copying, test selector, package, and installer assembly.
- `.github/workflows/release.yml` — current source version, Go setup, build entry, and release asset.
- `tests/backend/internal/codexeval/` — runtime, detector, Candy, Pelican, inspection tests, and test data.
- `licenses/ModelTrace-MIT.txt` and `licenses/THIRD_PARTY_NOTICES.txt` — current bundled license and attribution.
- `README.md` — Chinese user instructions and credits.
- Upstream tree `7a49068` — exact runtime resources, tests, Go embeds, and upstream license.

## Related Specs

- `.trellis/spec/backend/directory-structure.md`
- `.trellis/spec/backend/minimal-linux-delivery.md`
- `.trellis/spec/backend/quality-guidelines.md`
- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/database-guidelines.md`

## Caveats / Not Found

- This report maps files and checks. It does not change patches, licenses, tests, or build commands.
- No patch was applied. No build, test, installer, release, network publication, branch, or merge was run.
- The changed server and installer sizes are unknown until the approved implementation is built.
- Detailed algorithm, page behavior, and history decisions remain in their separate research reports.
