# U22 Patch Series Implementation

Date: 2026-10-18

## Scope and owners

The user approved production patches, separate tests, two script entries, and removal of `docs/`.
This change replaces the U21 source overlay without changing effective runtime source.
The implementation owner changed product layout, commands, tests, README, AGENTS, the notice, and the release workflow.
The main session owns task metadata, specifications, independent review, staging, and commits.
Unrelated work and private configuration remain unchanged.

## Production patch inventory

`patches/series` lists these seven patches in this order:

| Patch | Files | Contents |
| --- | ---: | --- |
| `0001-evaluation-core.patch` | 10 | Evaluation types, four runners, algorithms, task runtime, and history persistence |
| `0002-evaluation-data.patch` | 6 | Runtime context, fingerprint probes and baselines, ModelTrace bank, attribution, and license |
| `0003-vendor-assets.patch` | 5 | Existing Markdown and math renderers |
| `0004-admin-channel.patch` | 2 | Native account sender and administrator request adapter |
| `0005-inspection.patch` | 5 | Inspection coordinator, notifications, settings types, native settings persistence, and administrator adaptation |
| `0006-evaluation-ui.patch` | 13 | Page assembly, embedded page assets, and the native administrator iframe view |
| `0007-sub2api-integration.patch` | 9 | Existing native route, lifecycle, embed bypass, sidebar, router, origin, and translation edits |

The first six patches add the former 41 production files.
The last patch retains the former integration patch bytes unchanged.
The complete series affects 50 production paths.
No production patch contains Go test source or `testdata/` paths.

### Core files

All paths below use `backend/internal/codexeval/`:

```text
candy.go
codex.go
fingerprint.go
fingerprint_algorithm.go
model.go
modeltrace.go
modeltrace_algorithm.go
pelican.go
runtime.go
state.go
```

### Runtime data

All paths below use `backend/internal/codexeval/data/`:

```text
codex_context.json
fingerprint_baselines.json
fingerprint_probes.json
modeltrace/LICENSE
modeltrace/README.md
modeltrace/unified_bank.json
```

### Vendor assets

All paths below use `backend/internal/codexeval/web/vendor/`:

```text
markdown-it-cjk-friendly.min.js
markdown-it.min.js
temml.css
temml.min.js
texmath.js
```

### Administrator channel

```text
backend/internal/handler/admin/account_codex_eval_handler.go
backend/internal/service/account_codex_eval_service.go
```

### Inspection

```text
backend/internal/codexeval/inspection.go
backend/internal/codexeval/inspection_notify.go
backend/internal/codexeval/inspection_settings.go
backend/internal/handler/admin/codex_inspection_handler.go
backend/internal/service/codex_inspection_settings.go
```

### Page UI

```text
backend/internal/codexeval/web.go
backend/internal/codexeval/web/app.js
backend/internal/codexeval/web/candy.js
backend/internal/codexeval/web/catalog.js
backend/internal/codexeval/web/components.js
backend/internal/codexeval/web/fingerprint.js
backend/internal/codexeval/web/inspection.js
backend/internal/codexeval/web/modeltrace.js
backend/internal/codexeval/web/pelican.js
backend/internal/codexeval/web/style.css
backend/internal/codexeval/web/theme.css
backend/internal/codexeval/web/ui.html
frontend/src/views/admin/CodexEvalView.vue
```

### Upstream integration

```text
backend/internal/server/router.go
backend/internal/server/routes/admin.go
backend/cmd/server/main.go
backend/internal/web/embed_on.go
frontend/src/components/layout/AppSidebar.vue
frontend/src/router/index.ts
frontend/src/utils/embedded-url.ts
frontend/src/i18n/locales/en/common.ts
frontend/src/i18n/locales/zh/common.ts
```

## Separate test inventory

Direct moves preserved all 11 Go test files and the single golden JSON file.
Paths under `tests/backend/` match their upstream backend package paths:

```text
internal/codexeval/candy_test.go
internal/codexeval/codex_test.go
internal/codexeval/fingerprint_test.go
internal/codexeval/inspection_test.go
internal/codexeval/modeltrace_test.go
internal/codexeval/pelican_test.go
internal/codexeval/runtime_test.go
internal/codexeval/testdata/modeltrace-golden.json
internal/handler/admin/codex_inspection_handler_test.go
internal/service/account_codex_eval_service_test.go
internal/service/codex_eval_lifetime_test.go
internal/service/codex_inspection_settings_test.go
```

The existing installer fixtures and UI checks retain their `tests/installer/` and `tests/ui/` owners.
The installer test imports `scripts/install.py`.
The UI harness reads assets from its effective patched upstream tree.
An explicit tree argument takes precedence over `SUB2API_SOURCE`.
Existing assertions remain unchanged.

## Command consolidation

Only `scripts/build.py` and `scripts/install.py` remain under `scripts/`.
The installer moved directly from root without source changes.
Its packaged filename remains `install.py`, beside the executable.

`build.py --source ...` copies the original tree into a standard temporary workspace.
`prepare_tree()` applies each production patch in series order.
It uses `git apply --binary --unidiff-zero --whitespace=nowarn`.
Git is a developer build dependency. The installation host still needs only Python and Docker CLI.

The build retains the native frontend build and static Linux/amd64 Go compilation.
Its `--check` option copies `tests/backend/` into corresponding temporary upstream packages.
That option retains the existing focused Go test expression and touched-package vet list.
It also invokes the existing installer tests and UI harness through the single check entry.
Corepack setup from the removed check wrapper remains in that check flow.
No new test source, broader Go test selection, race check, or coverage gate was added.

After compilation, the same entry packages the executable, installer, and required license texts.
The release workflow now invokes this combined entry once.
The separate packaging script and both shell wrappers went to macOS Trash.

## Composition evidence and corrections

The original source tree is `/tmp/sub2api-codex-eval-v0217-check`.
Its existing frontend dependencies remain available.
The available `dist/` contents did not contain an original v0.2.17 source tree.
No source download ran.

Both production and check preparation used this standard temporary workspace:

```text
/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-u22-jurh19qo/
```

The production tree is `production/source/`. The check tree is `checks/source/`.
All seven patches applied successfully to both copies.
All 41 production additions matched their former source bytes and `0644` modes.
All 12 moved test inputs matched their former bytes and modes.
Those test paths are absent from the production copy and present in the check copy.

A later comparison used the former staged `src/` files as an independent retained reference.
All 53 files matched their effective production or moved test contents.
The integration patch also matched its former staged `patches/sub2api.patch` bytes.
The comparison did not stage or change the index.

The first text-patch composition failed only for `texmath.js`.
That asset has 368 CRLF lines. The patch editing tool normalized those line endings.
A standard Git binary patch now preserves its exact 14071 bytes.
The other 40 production additions use ordinary text patches.
Git binary metadata comes from the ordinary full-index diff format.
No delivery checksum or source manifest was added.

Four original vendor/UI lines contain trailing spaces. Their bytes remain unchanged.
The apply command suppresses whitespace warnings without fixing source whitespace.

The first mocked CLI assertion compared a `/tmp` path with its resolved `/private/tmp` path.
The corrected assertion uses `resolve()` and passes. No product correction was needed for that check.

## Documentation and removals

README now contains the useful manual rollback commands, backup lifetime, and native history location.
It also describes the seven-patch series, separate tests, two entries, and automatic packaging.
AGENTS describes the same current layout.
The sole third-party notice identifies `patches/` and `patches/series` as the production source locations.

After successful composition checks, these paths went to macOS Trash:

```text
src/
docs/
scripts/build-sub2api.sh
scripts/test-sub2api.sh
scripts/package-sub2api.py
```

The integration patch moved to `patches/0007-sub2api-integration.patch`.
The build script moved to `scripts/build.py` before its focused edits.
The root installer moved to `scripts/install.py` unchanged.
No replacement source archive, compatibility wrapper, or duplicate operations guide was added.

## Archive-only update

The packaged notice changed its source locations, so the archive required regeneration.
The merged `package()` function generated the archive with the existing binary.
No frontend build or server compilation ran for this layout change.

The binary remains 125599904 bytes, mode `0755`, with modification time `1792248387913425168` nanoseconds.
The archive is 38700708 bytes at `dist/build/sub2api-codex-eval.tar.gz`.
It contains exactly these seven entries:

```text
sub2api-codex-eval/sub2api-codex-eval
sub2api-codex-eval/install.py
sub2api-codex-eval/LICENSE/Project-MIT.txt
sub2api-codex-eval/LICENSE/ModelTrace-MIT.txt
sub2api-codex-eval/LICENSE/Renderers-MIT.txt
sub2api-codex-eval/LICENSE/Sub2API-LGPL-3.0.txt
sub2api-codex-eval/LICENSE/THIRD_PARTY_NOTICES.txt
```

Every entry matches its source bytes, including the unchanged executable and installer.
The binary size, modification time, and mode remained unchanged through packaging.

## Checks and limits

| Check | Result |
| --- | --- |
| Production and check patch composition | Both passed. Exact source bytes and modes match. |
| Prior staged source comparison | All 53 former source files and the integration patch match. |
| Test separation | All 12 task test paths appear only in the check tree. Production patches contain no test paths. |
| Existing installer tests | Both passed with the moved installer import. |
| Existing UI harness | Passed against the patched check tree with reused dependencies. |
| Existing focused Go tests | Passed in `internal/codexeval`, `internal/service`, and `internal/handler/admin`. |
| Build orchestration | Both check branches forwarded the expected commands and build environment through mocked process calls. |
| CLI orchestration | Preserved source/check arguments and automatic packaging through mocked build/package calls. |
| Command help | Both approved entries returned their help without product changes. |
| Python, shell, and Node syntax | Passed for affected scripts and existing fixture files. |
| Current Markdown links | All 66 local links passed across README, AGENTS, and nine backend specs. |
| Current obsolete-path search | No old source, documentation, or command references remain outside historical evidence. |
| Package inventory and bytes | Exactly seven entries, all matching their source bytes. |
| `git diff --check` | Passed for the current tracked working diff. New patch source retains original whitespace. |

No product build, repeated baseline, broad suite, race check, Docker operation, deployment, staging, or commit ran.
The existing focused Go tests compiled their test packages. No server binary was compiled.
No new product behavior or runtime verification is claimed.
Earlier runtime evidence remains applicable because effective production bytes remain unchanged.
The earlier history cancellation report remains open.
The main session must complete independent review and final work-commit preparation.
