# Upstream Comparison Evidence

## Source and Revision

Requested comparison:
https://github.com/Nanboom233/codex-candy-eval/compare/main...haowang02%3Acpa-plugin-codex-candy-eval%3Amain

The GitHub page identifies two commits but cannot render its diff.
The existing `upstream` remote points to `haowang02/cpa-plugin-codex-candy-eval`.
`git fetch upstream main` confirmed tip `4f8367642cc1e629c4ee15a975e5e0fc53c4ac16` on 2026-10-19.
`git log origin/main..upstream/main` lists only these commits:

| Commit | Change |
| --- | --- |
| `e120868f1f5bb4cb10c175fcab9eec7738c4908d` | Add Pelican variants, SVG presentation, saved version metadata, and schema 2. |
| `4f8367642cc1e629c4ee15a975e5e0fc53c4ac16` | Change the CPA plugin release version from `0.4.1` to `0.4.2`. |

The comparison changes 14 files, with 148 insertions and 60 deletions.
The shared ancestor is `7a490680a2e75660250faa235ebeb521d37d8564`.
The prior synchronization design explicitly deferred these two commits to keep its Pelican scope unchanged.

## Upstream Contract

`internal/plugin/pelican.go` defines `original` and `all-in-one` with labels, exact prompts, and HTML/SVG language selection.
The all-in-one prompt requests the upstream composite scene and the Candy problem's answer in the drawing.
The evaluator stores `variant` before sending and wraps SVG in a viewport-filling HTML page.
The code limit remains 256 KiB. Extraction prefers the requested language and retains the prior markup fallback.
Upstream requires an explicit request version. Native existing requests omit it, so native omission must still select `original`.

`internal/plugin/migration.go` changes schema 1 to 2 and marks prior Pelican records as `original`.
Its generic migration framework belongs to CPA. The native product already has a different minimal JSON loader.
Schema-1 native fingerprint records contain the current detector data and must survive this upgrade.

`internal/plugin/plugin.go` injects variants into the page from Go definitions.
`web/pelican.js` adds selection, question updates, request payloads, and saved version labels.
`web/ui.html`, `components.js`, `app.js`, and `style.css` add the control, icon, picture labels, and toolbar column.
The release-only version change does not change the native Sub2API version.

## Native Owner Map

| Owner | Anchor | Required adaptation |
| --- | --- | --- |
| `patches/0001-evaluation-core.patch` | `:1629` | Native `pelican.go`, its result, extraction, and runner. |
| `patches/0001-evaluation-core.patch` | `:1734` | Native `runtime.go` and `StartRequest`. |
| `patches/0001-evaluation-core.patch` | `:1996` | Native `state.go`, typed maps, schema read boundary, and atomic save. |
| `patches/0006-evaluation-ui.patch` | `:1` | Native `web.go` page assembly and injected configuration. |
| `patches/0006-evaluation-ui.patch` | `:1860` | Native `web/pelican.js`, `pl-` controls, preferences, and metadata. |
| `patches/0006-evaluation-ui.patch` | `:2171` | Native Pelican toolbar includes an extra group field. |
| `patches/0006-evaluation-ui.patch` | `:2549` | Native Pelican question and toolbar markup. |
| `tests/backend/internal/codexeval/pelican_test.go` | `:7` | Current HTML extraction cases. |
| `tests/backend/internal/codexeval/runtime_test.go` | `:84` | Existing starts omit a version. Preserve their contract. |
| `tests/backend/internal/codexeval/runtime_test.go` | `:184` | Existing unversioned migration fixture and retained-field assertions. |
| `tests/ui/sub2api-dropdowns.cjs` | `:31` | Inject the new configuration before initializing Pelican in the harness. |
| `scripts/build.py` | `:53` | Focused Go selector currently includes `Codex` and `ExtractHTML`. |

Native page code uses `pl-` identifiers and native account targets, not upstream credential lists.
Native history uses `results`, `fingerprints`, and `pelicans`. Keep these map names.
The current schema is 1 and rejects any other explicit marker.
Unversioned files currently discard retired fingerprint and ModelTrace data, then save immediately.
Keep that retirement rule while marking old Pelican records and advancing the marker.

## Validation Inputs

The retained original source directory `/tmp/sub2api-codex-eval-v0217-check` has frontend dependencies and no Git metadata.
The existing build copies source before applying the seven patches and attaching native tests.
Use the standard build entry. Do not apply patches directly to the retained original source.
The previous effective source is `/tmp/upstream-sync-native/source`.
Previous checks establish unchanged transport and inspection behavior. New variant and migration behavior need fresh checks.

## Research Limits

No product code changed during this investigation. No build or product test ran.
The IDE file tools could not access this repository because the IDE has another project open.
Harness shell reads and Git inspection supplied the local source evidence.
