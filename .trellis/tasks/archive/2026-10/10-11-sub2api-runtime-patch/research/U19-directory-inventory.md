# Research: U19 directory inventory

- Query: Identify active Sub2API delivery paths, retired implementation candidates, direct callers, current specification references, and route risks before U19 repository organization.
- Scope: internal
- Date: 2026-10-18

## Findings

### Active delivery boundary

The current Sub2API integration is confined to the following paths.

| Path | Role and evidence |
| --- | --- |
| `native/backend/` | Fresh native evaluation, inspection, sender, handler, history, page assets, and adjacent Go tests. The directory contains 11 native test files. |
| `native/frontend/` | Native administrator iframe view. `CodexEvalView.vue:23` embeds the relative `/codex-eval` page. |
| `native/sub2api.patch` | Small upstream patch for native routes, shutdown, embedded-page bypass, sidebar, router, origin resolution, and translations (`native/sub2api.patch:1-97`). |
| `scripts/build-sub2api.py` | Copies an existing source tree to a temporary workspace, applies the patch, overlays `native/backend` and `native/frontend`, rebuilds the frontend, runs focused checks, and emits `dist/build/sub2api-codex-eval` (`scripts/build-sub2api.py:18-42`). |
| `scripts/build-sub2api.sh` | Thin wrapper around the active Python build entry (`scripts/build-sub2api.sh:1-3`). |
| `scripts/package-sub2api.py` | Packages the fixed binary, root `install.py`, and license texts into `dist/build/sub2api-codex-eval.tar.gz` (`scripts/package-sub2api.py:8-16`). |
| `scripts/test-sub2api.sh` | Requires an existing original v0.2.17 tree and calls the active build with `--check` (`scripts/test-sub2api.sh:22-33`). |
| `install.py` | Selects a running container, validates v0.2.16 or v0.2.17, stages and replaces `/app/sub2api`, waits for native readiness, and restores the backup on failure (`install.py:134-193`). |
| `.github/workflows/release.yml` | Checks out `Wei-Shaw/sub2api` at `v0.2.17`, then calls only `scripts/build-sub2api.py` and `scripts/package-sub2api.py` (`.github/workflows/release.yml:15-41`). |
| `dist/build/sub2api-codex-eval*` | Fixed ignored binary and archive outputs consumed by packaging and the native fixture. |

Current native fixture ownership is also clear. `tests/installer/native_inspection_fixture.py:11,77-101` uses the retained `sub2api-codex-eval-test` container and copies the fixed binary. `native_fixture_start.sh`, `native_fixture_supplier.py`, `seed_browser_fixture.py`, and `test_native_version.py` support that fixture. `tests/ui/sub2api-dropdowns.cjs:8-31` loads the native web assets and uses `/codex-eval`, so it matches the active route.

The original CPA implementation remains a separate owner. Keep `cmd/cpa-codex-candy-eval/`, `internal/plugin/`, `scripts/deploy.sh`, `scripts/test-ui.cjs`, `examples/inspection.py`, and CPA documentation. The C ABI bridge imports `internal/plugin` (`cmd/cpa-codex-candy-eval/cabi.go:67-137`), while the retained UI harness reads `internal/plugin/web` (`scripts/test-ui.cjs:7`).

### Retired path map

The following candidates are confirmed abandoned implementation inputs. Move each source to the matching path under `.trellis/tasks/10-11-sub2api-runtime-patch/history/retired/` and preserve the source contents.

| Former path | Proposed history path |
| --- | --- |
| `adapter/` | `history/retired/adapter/` |
| `internal/eval/` | `history/retired/internal/eval/` |
| `patches/sub2api/v0.2.15/` | `history/retired/patches/sub2api/v0.2.15/` |
| `compatibility/sub2api-v0.2.15.json` | `history/retired/compatibility/sub2api-v0.2.15.json` |
| `scripts/sub2api_build.py` | `history/retired/scripts/sub2api_build.py` |
| `scripts/test_sub2api_build.py` | `history/retired/scripts/test_sub2api_build.py` |
| `scripts/universal_dev.py` | `history/retired/scripts/universal_dev.py` |
| `scripts/universal_payload.py` | `history/retired/scripts/universal_payload.py` |
| `scripts/prebuilt_dev.py` | `history/retired/scripts/prebuilt_dev.py` |
| `scripts/prebuilt_test.py` | `history/retired/scripts/prebuilt_test.py` |
| `scripts/export-sub2api-overlays.py` | `history/retired/scripts/export-sub2api-overlays.py` |
| `scripts/export-sub2api-production.py` | `history/retired/scripts/export-sub2api-production.py` |
| `scripts/generate-wire.sh` | `history/retired/scripts/generate-wire.sh` |
| `scripts/pnpm` | `history/retired/scripts/pnpm` |
| `scripts/pnpm-support.sh` | `history/retired/scripts/pnpm-support.sh` |
| `tests/host-overlay/` | `history/retired/tests/host-overlay/` |
| `tests/overlays/codexeval/` | `history/retired/tests/overlays/codexeval/` |
| `tests/installer/current_suite.py` | `history/retired/tests/installer/current_suite.py` |
| `tests/installer/docker_lifecycle.py` | `history/retired/tests/installer/docker_lifecycle.py` |
| `tests/installer/test_builder_progress.py` | `history/retired/tests/installer/test_builder_progress.py` |
| `tests/installer/test_cleanup.py` | `history/retired/tests/installer/test_cleanup.py` |
| `tests/installer/test_docker_boundary.py` | `history/retired/tests/installer/test_docker_boundary.py` |
| `tests/installer/test_docker_compatibility.py` | `history/retired/tests/installer/test_docker_compatibility.py` |
| `tests/installer/test_http_transport.py` | `history/retired/tests/installer/test_http_transport.py` |
| `tests/installer/test_preflight.py` | `history/retired/tests/installer/test_preflight.py` |
| `tests/installer/test_prebuilt.py` | `history/retired/tests/installer/test_prebuilt.py` |
| `tests/installer/test_source_progress.py` | `history/retired/tests/installer/test_source_progress.py` |
| `tests/installer/test_transactions.py` | `history/retired/tests/installer/test_transactions.py` |
| `tests/installer/test_universal.py` | `history/retired/tests/installer/test_universal.py` |
| `tests/installer/universal_ci.py` | `history/retired/tests/installer/universal_ci.py` |

Inventory sizes are 27 files in `adapter/`, 38 in `internal/eval/`, 18 in `patches/sub2api/v0.2.15/`, 8 in `tests/host-overlay/`, 8 in `tests/overlays/codexeval/`, and 13 installer files in this initial inventory.

This initial inventory predates the final classification correction. The complete
relocation also includes `tests/go.mod`,
`tests/installer/seed_browser_fixture.py`, `tests/ui/sub2api-browser.js`, and
the pre-U19 `quality-guidelines.md`. The final 42-row mapping is in
`history/retired/README.md`; the correction below records why these files are
historical rather than active checks.

### Direct callers and coupling

The retired paths form one obsolete build and overlay graph.

- `scripts/pnpm` reads `compatibility/sub2api-v0.2.15.json` and invokes `scripts/sub2api_build.py`.
- `scripts/test_sub2api_build.py` invokes `scripts/pnpm-support.sh` and creates `internal/eval` fixtures.
- `scripts/export-sub2api-overlays.py` reads `tests/host-overlay/` and writes the old `patches/sub2api/v0.2.15/test-overlay/` directory.
- `scripts/universal_dev.py` reads the old test overlay and `tests/overlays/codexeval/`.
- `scripts/prebuilt_dev.py` invokes `adapter/build.py`.
- `tests/installer/universal_ci.py` invokes `scripts/package-sub2api.py` and `scripts/universal_dev.py`.
- `tests/installer/test_transactions.py` reads the compatibility manifest and imports `universal_payload`.
- `scripts/export-sub2api-production.py` only prints a retired-workflow message that names `scripts/universal_dev.py`.

No active release workflow, native build, package, installer, retained native fixture, or active Node UI harness calls these retired helpers. The old scripts must move together because their internal relative references will no longer resolve from repository root after relocation. The moved copies are historical evidence and do not need to remain runnable.

### Active specification references to update

Eight backend specifications point to the retired graph and should move to task history after the active references are corrected.

| Specification | Retired references |
| --- | --- |
| `.trellis/spec/backend/build-guards.md` | Adapter/build owners, compatibility manifest, wire and pnpm helpers, universal builder, overlays, and old release counts (`:18-19,41,47,51-52,83,132,180,185-188,207,250`). |
| `.trellis/spec/backend/deployment-evidence.md` | v0.2.15 patch 005 (`:12`). |
| `.trellis/spec/backend/evaluation-core.md` | `internal/eval`, overlay tests, and old copy helper (`:15,17,23,123`). |
| `.trellis/spec/backend/host-coordinator.md` | v0.2.15 patch 006 (`:12`). |
| `.trellis/spec/backend/host-target-adapter.md` | v0.2.15 patch 004 (`:9`). |
| `.trellis/spec/backend/patch-native-installation.md` | Adapter runtime templates and structural rules (`:26,157,298`). |
| `.trellis/spec/backend/run-state.md` | `internal/eval` and overlay tests (`:9,129`). |
| `.trellis/spec/backend/universal-linux-delivery.md` | Adapter builder, universal payload, and old package contracts (`:39,48-50,65,125,128,162-163`). |

The active `directory-structure.md:61` already labels these paths historical but has no relocation link. Add the inventory or relocated-history link there. `quality-guidelines.md:28-33` still describes the adapter and `internal/eval` as Sub2API production inputs even though the section is marked retired. Rewrite that paragraph for the native boundary. Its historical scenario links at `:241,338,431` also point to the old patch directory and need relocated-history links. The active `index.md` already states that historical specifications do not add requirements (`:23-26`) and needs no path change.

### Browser route caveat

The current patch registers only `/codex-eval` and `/codex-eval/` (`native/sub2api.patch:7-10`), and the native iframe view embeds `/codex-eval` (`native/frontend/src/views/admin/CodexEvalView.vue:23`). The active route contract also lists those two paths (`.trellis/spec/backend/runtime-installation.md:34`). `tests/ui/sub2api-browser.js:10` and `:21` still navigate to `/codex-eval/ui?theme=...` and expect four tabs. Treat this file as active pending route ownership review. Do not move it with the retired overlays. `tests/ui/sub2api-dropdowns.cjs:21` is aligned with `/codex-eval`.

### External references and versions

The release workflow is the repository's upstream version reference. It checks out `Wei-Shaw/sub2api` at `v0.2.17` (`.github/workflows/release.yml:15-20`). The native build injects `main.Version=0.2.17` (`scripts/build-sub2api.py:40-41`). The active delivery specification records the same release and the accepted v0.2.16 upgrade source (`.trellis/spec/backend/minimal-linux-delivery.md:5,42-46,62`). No external web lookup was needed for this inventory.

## Related specs

- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md:5-15` defines the U19 active boundary, retained paths, history relocation, and generated-file exclusion.
- `.trellis/tasks/10-11-sub2api-runtime-patch/implement.md:5-13` requires the inventory, path moves, current documentation updates, and affected path checks.
- `.trellis/spec/backend/directory-structure.md` describes current native and retained CPA owners.
- `.trellis/spec/backend/quality-guidelines.md` defines current checks and marks installation and universal-adapter sections as retired.
- `.trellis/spec/backend/runtime-installation.md` defines active native API and page routes.
- `.trellis/spec/backend/native-inspection.md` defines native inspection ownership and tests.
- `.trellis/spec/backend/minimal-linux-delivery.md` defines the fixed build, package, installer, and fixture paths.

## Caveats / Not Found

- Generated Python caches, `.playwright-cli/`, `output/`, and ignored `dist/` artifacts were excluded from the move inventory.
- CPA sources share evaluation names with retired Sub2API copies. Classification used owner paths and direct callers, so `internal/plugin/` remains retained.
- Historical moved scripts will retain old root-relative imports and source paths. Do not rewrite their contents during relocation.
- `tests/ui/sub2api-browser.js` route ownership is unresolved because it targets `/codex-eval/ui`; run a focused route review before retaining it unchanged.
- This research performed read-only inspection. No files were moved, deleted, built, tested, or committed.

## U19 implementation classification correction

The main session resolved the browser ownership caveat before implementation.
`seed_browser_fixture.py` depends on the removed `installer.Installer` class and old `codexeval-fixture-*` containers.
It does not support the current retained native fixture.
`sub2api-browser.js` expects `/codex-eval/ui` and four tabs. The current page uses `/codex-eval` and five tabs.
No current entry point calls either file. Both files move to history without a replacement suite.

`tests/go.mod` declares `codex-eval-development-overlays` and belongs only to the retired overlay tests.
It moves with those tests. Current Python fixtures and Node checks do not use it.
The full pre-U19 quality guide also moves to history. A focused current guide replaces it at the active path.
The [relocation inventory](../history/retired/README.md) records these moves and the original inventory paths.
