# Directory Structure

> Keep production changes in functional patches and all tests under `tests/`.

---

## Overview

Active integration lives in a numbered production patch series under `patches/`.
The standalone build downloads the selected upstream tag and applies top-level `patches/*.patch` files in lexical filename order.
This repository has no root Go module or shared-library entry point.
The user package replaces the server binary without compilation on the installation host.

Evidence: [native build](../../../scripts/build.py), [patches](../../../patches/), and [installer source](../../../scripts/installer/).

---

## Directory Layout

```text
patches/                 Contextual production patches in lexical filename order
scripts/build.py         Download source, apply patches, test, build, and package
scripts/installer/       Go self-extracting installation and failure restoration
tests/test_*_test.go     Actual native Go test cases with package-context annotations
tests/test_ui.cjs         Actual native UI assertions
tests/test_build.py       Common source-composition and test-discovery cases
tests/test_installer_test.go  Common installer transaction tests
tests/support.py          Native test execution support
tests/native_*            Manual test fixture and support programs
assets/                   Data resources, including golden JSON
licenses/                Original third-party license files
dist/build/<os>-<arch>/  Ignored platform server and installer executables
dist/work/              Ignored workspace retained within one CI build job
README.md AGENTS.md LICENSE  Project documentation and license
.github/                 Single build job and optional Release publication
.trellis/                Workflow, specifications, tasks, and session records
.agents/ .codex/ .claude/ .devin/  Tool integration and skills
```

The self-contained installer carries native and required Linux container payloads.
Users select Docker/Compose, a native program, or supported Apple container installation and retain the original program backup.
Builds download upstream source and use frozen dependencies in a temporary workspace, rebuilding the patched frontend before embedding it.
Only compressed application programs and required license texts enter the installer payload.
The build assembles the installer in a temporary Go module.
Prepare PTerm dependencies there. Keep terminal messages under `scripts/installer/`, separate from installation decisions.
Names use `sub2api-plugins-[installer-]<upstream-tag>-<os>-<arch>`, where `installer-` is present only for installers.
Windows outputs add `.exe`. CI uploads each file directly, with its basename as the artifact name.
The parent and standalone kits carry aligned common tests. Execute common cases once per build.
Keep installer tests out of production source and keep the raw server as a local build output.

Private configuration and reusable dependency caches remain outside source and package inputs.
Keep user instructions in `README.md`. Do not create a separate `docs/` tree or source overlay.

Apply the production patches in this order:

```text
0001-evaluation-core.patch
0002-evaluation-data.patch
0003-vendor-assets.patch
0004-admin-channel.patch
0005-inspection.patch
0006-evaluation-ui.patch
0007-sub2api-integration.patch
```

Keep test source and test-only data out of production patches.
Keep all test cases and support code in flat `tests/`. Do not place test programs in assets or embedded source strings.
Native suites use `// sub2api-package: backend/internal/codexeval` or their corresponding native package path.
The runner stages native suites in the temporary package to preserve private access and executes discovered cases once.
Manual `tests/native_*` helpers remain outside automatic `test_*` discovery.

Evidence: [repository layout and private configuration](../../../AGENTS.md).

---

## Native Module Organization

- Keep evaluation tasks and account admission in `codexeval/runtime.go`.
- Keep each evaluation's request types, results, runner, and algorithm in its existing package.
- Keep the LM detector in `fingerprint.go` and `fingerprint_detector.go`. Retired ModelTrace and probe algorithms have no runtime owner.
- Keep detector assets under `codexeval/data/lm-detector/` in patch `0002`, including its Git binary hunk.
- Keep the detector golden fixture at `assets/lm-detector-golden.json`, outside production patches.
- Keep shared request retries in `model.go` and the shared context decoder in `codex.go`.
- Keep corresponding native Go tests directly under flat `tests/` with explicit package context.
- Keep history persistence in `codexeval/state.go`.
- Keep inspection scheduling and notifications in `codexeval/inspection*.go`.
- Keep native settings persistence and upstream account sending in `internal/service/` additions.
- Keep administrator HTTP adaptation in `internal/handler/admin/` additions.
- Keep page logic in `codexeval/web/` and assembly in `codexeval/web.go`.
- Keep the iframe view with the evaluation UI patch and upstream edits with the integration patch.

Do not create a router framework, service hierarchy, or new package solely to match a generic architecture template.

Examples: [evaluation core](../../../patches/0001-evaluation-core.patch) and [native inspection](../../../patches/0005-inspection.patch).

---

## Native Examples

| Work | Owner | Existing verification |
| --- | --- | --- |
| Change tasks, cancellation, or history | `runtime.go` and `state.go` | `tests/test_codexeval_runtime_test.go` |
| Change automatic Candy inspection | `inspection*.go` | `tests/test_codexeval_inspection_test.go` |
| Change native account transport | `account_codex_eval_service.go` | `tests/test_service_account_codex_eval_service_test.go` and lifetime suite |
| Change settings storage | `codex_inspection_settings.go` | `tests/test_service_codex_inspection_settings_test.go` |
| Change embedded page behavior | `web.go` and `web/` | `tests/test_ui.cjs` |

When adding an embedded asset, update its Go embed declaration and page assembly.
Keep source, tests, and package ownership aligned with [binary delivery](./minimal-linux-delivery.md).
