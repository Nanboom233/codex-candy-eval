# Directory Structure

> Place code beside the package that owns its behavior.

---

## Overview

<!--
Document your project's backend directory structure here.

Questions to answer:
- How are modules/packages organized?
- Where does business logic live?
- Where are API endpoints defined?
- How are utilities and helpers organized?
-->

This repository contains one Go module: `cpa-codex-candy-eval`. It builds a CPA shared-library plugin, not a standalone HTTP server.

The `main` package owns the C ABI bridge. The `plugin` package owns host dispatch, management responses, evaluation runners, shared state, and embedded assets.

Evidence: [go.mod](../../../go.mod), [C ABI bridge](../../../cmd/cpa-codex-candy-eval/cabi.go), and [plugin entry point](../../../internal/plugin/plugin.go).

---

## Directory Layout

<!-- Replace with your actual structure -->

```text
cmd/cpa-codex-candy-eval/  Main package and cshared C ABI bridge
internal/plugin/         Go handlers, runners, algorithms, and adjacent tests
internal/plugin/web/     Embedded HTML, plain JavaScript, and CSS
internal/plugin/web/vendor/  Vendored Markdown and math renderers
internal/plugin/data/    Embedded probes and attribution baselines
scripts/                Deployment helper and Node UI checks
docs/                   Usage guides and README images
examples/               Scheduled inspection script
.github/workflows/      Release checks, shared-library builds, and packaging
dist/                   Ignored build artifacts
.trellis/spec/          Project coding guidelines
```

The root installers are `install.sh` and `install.ps1`. Deployment settings stay in ignored `.env`, not in source or specs.

Evidence: [repository layout and private configuration](../../../AGENTS.md).

---

## Module Organization

<!-- How should new features/modules be organized? -->

- Keep management registration and routing in `plugin.go`.
- Keep credential inventory, model execution, run scheduling, and persistence in their existing owner files.
- Place test-specific request types, results, and runners beside that test's handler.
- Keep pure attribution algorithms separate from host calls and worker orchestration.
- Add Go tests beside the implementation with the `_test.go` suffix.
- Keep browser logic in `web/`. The Go entry point embeds and assembles these assets.
- Inspect the CPA and management-console source before changing their contracts. Their sibling paths appear in `AGENTS.md`.

Do not create a router framework, service hierarchy, or new package solely to match a generic architecture template.

Examples: [shared run scheduling](../../../internal/plugin/runs.go), [state ownership](../../../internal/plugin/state.go), [fingerprint algorithms](../../../internal/plugin/fingerprint_algorithm.go), and [ModelTrace orchestration](../../../internal/plugin/modeltrace.go).

---

## Naming Conventions

<!-- File and folder naming rules -->

Use the existing `package plugin` for internal behavior. Use descriptive topic filenames such as `credentials.go`, `model.go`, and `inspection.go`.

Follow adjacent names: `candyRunRequest`, `candyResult`, `candyRunResponse`, and `runCandy`. Shared helpers include `startRuns`, `executeModel`, and `appendRecent`.

Export only symbols needed across package boundaries, such as `HandleMethod`, `SetHostCall`, `Quiesce`, and `ABIVersion`.

Evidence: [candy handler and runner](../../../internal/plugin/candy.go) and [plugin exports](../../../internal/plugin/plugin.go).

---

## Examples

<!-- Link to well-organized modules as examples -->

| Work | Owner | Related verification |
| --- | --- | --- |
| Add or change a management route | `plugin.go` registration and `handleManagement` | `TestRegisteredRoutesAreHandled` in `plugin_test.go` |
| Change candy run behavior | `candy.go` and shared `runs.go` | `plugin_test.go` and `credentials_test.go` |
| Change persisted history | `state.go` | `TestStatePersistence` and `TestUnreadableStateIsPreserved` |
| Change embedded page assembly | `plugin.go` and `web/` | `TestScriptsInlineSafely` and `scripts/test-ui.cjs` |
| Change C buffer ownership | `cmd/cpa-codex-candy-eval/cabi.go` | Inspect the host ABI and use the cshared build flags |

References: [plugin tests](../../../internal/plugin/plugin_test.go), [credential tests](../../../internal/plugin/credentials_test.go), and [UI checks](../../../scripts/test-ui.cjs).

When adding an embedded asset, update its Go embed declaration and page assembly. The plugin embeds browser files instead of using a separate development server.
