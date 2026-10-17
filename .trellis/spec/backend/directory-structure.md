# Directory Structure

> Keep production changes in functional patches and all tests under `tests/`.

---

## Overview

Active integration lives in a numbered production patch series under `patches/`.
The build applies `patches/series` to an existing original Sub2API source tree and uses that tree's Go module.
This repository has no root Go module or shared-library entry point.
The user package replaces the server binary without compilation on the installation host.

Evidence: [native build](../../../scripts/build.py), [patch order](../../../patches/series), and [installer](../../../scripts/install.py).

---

## Directory Layout

```text
patches/                 Numbered production patches and explicit series order
scripts/build.py         Apply patches, check, build, and package
scripts/install.py       Linux binary replacement and failure restoration
tests/backend/           Go tests and test-only data, matching upstream package paths
tests/installer/         Installer checks and native fixture
tests/ui/                Page checks against the patched source tree
licenses/                License texts and sole third-party notice source
dist/build/              Fixed ignored executable and installation archive
README.md AGENTS.md LICENSE  Project documentation and license
.github/                 Release workflow
.trellis/                Workflow, specifications, tasks, and session records
.agents/ .codex/ .claude/ .devin/  Tool integration and skills
```

The packaged installer uses the adjacent prebuilt binary and Docker CLI.
Users select a running container, retain its backup and open 账户评估 in the native administrator sidebar.
Builds use the existing original Sub2API tree and installed dependencies in a standard temporary workspace, rebuilding the patched native frontend before embedding it.
Only the executable, installer and necessary license texts enter the archive.

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
For checks, copy `tests/backend/` into the temporary tree's `backend/` directory to preserve package-private access.

Evidence: [repository layout and private configuration](../../../AGENTS.md).

---

## Native Module Organization

- Keep evaluation tasks and account admission in `codexeval/runtime.go`.
- Keep each evaluation's request types, results, runner, and algorithm in its existing package.
- Keep corresponding Go tests under `tests/backend/` with the same package path.
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
| Change tasks, cancellation, or history | `runtime.go` and `state.go` | `runtime_test.go` |
| Change automatic Candy inspection | `inspection*.go` | `inspection_test.go` |
| Change native account transport | `account_codex_eval_service.go` | Sender and lifetime tests under `tests/backend/internal/service/` |
| Change settings storage | `codex_inspection_settings.go` | Settings tests under `tests/backend/internal/service/` |
| Change embedded page behavior | `web.go` and `web/` | `tests/ui/sub2api-dropdowns.cjs` |

When adding an embedded asset, update its Go embed declaration and page assembly.
Keep source, tests, and package ownership aligned with [binary delivery](./minimal-linux-delivery.md).
