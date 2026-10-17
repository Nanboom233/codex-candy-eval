# Repository Instructions

## References and Scope

Active delivery targets Sub2API v0.2.17 on Linux/amd64. Use `scripts/install.py` and
the current Trellis task/spec contracts.

- Build and release conventions: `.github/workflows/release.yml`.
- Inspect relevant source before changing integration contracts. Keep changes
  focused and preserve unrelated work.

## Project Layout

```text
patches/                 Ordered production patches and series list
scripts/build.py         Build, check, and package entry
scripts/install.py       Linux binary replacement and failure restoration
tests/backend/           Go tests and test data, by upstream package path
tests/installer/         Installer tests and retained-container fixtures
tests/ui/                Page checks against the patched upstream tree
licenses/                Required license texts and the sole third-party notice
dist/build/              Fixed ignored executable and archive outputs
README.md                User and developer commands
AGENTS.md                Repository instructions
LICENSE                  Project MIT license
.github/                 Release workflow
.trellis/                Task workflow and project specifications
.agents/ .codex/ .claude/ .devin/  Tool integration and skills
```

Production code, runtime data, assets, and integration edits live in seven patches under `patches/`.
Apply them in `patches/series` order. Keep tests and test-only data outside production patches.
`scripts/build.py --check` copies Go tests into temporary upstream packages before the existing checks.
The UI harness reads the effective patched tree. Keep the original source tree unchanged.
The build entry creates the binary and archive without a separate packaging command.
Fixed outputs are `dist/build/sub2api-codex-eval` and `dist/build/sub2api-codex-eval.tar.gz`.

## Sub2API Test Container

Sub2API checks use one reusable Docker container named
`sub2api-codex-eval-test`, including its PostgreSQL and Redis test dependencies.
Reuse it for later checks. Do not create random fixture containers.
Remove only identified earlier task containers when cleanup is authorized.
Preserve unrelated containers, volumes, images, caches, and real application data.

## Sub2API Build and Checks

Use the current Trellis task and `.trellis/spec/backend/index.md` for check scope.
Use an existing original Sub2API source tree and installed frontend dependencies.

```sh
python3 scripts/build.py --source /path/to/original-sub2api
python3 scripts/build.py --source /path/to/original-sub2api --check
```

Reuse completed checks for unchanged inputs. Directory and documentation changes
need affected reference and syntax checks. They do not require another product build.

The build uses the existing upstream Go module and compiles with CGO disabled.
There is no root Go module or shared-library entry point.
Run `gofmt` on changed Go files. Check changed behavior through existing focused tests.
Do not add race, coverage, or full-suite gates.

## Private Configuration

Preserve ignored `.env` and machine configuration. The active installer does
not read dotenv files. Keep private files untracked and out of uploads.
Keep `.env` at mode `600`.
Never expose credentials or private deployment values in logs, screenshots,
commits, or reports.

## Installation and Verification

The archive places `install.py` beside the executable. Run `python3 install.py` after extraction.
Its repository source is `scripts/install.py`. It replaces the selected
container executable, retains a backup, and restores that backup after failure.
Use the active task and native delivery spec for approved runtime checks.
Keep production data and provider quota outside local fixture checks.

## Commits and Releases

Use an imperative subject. For nontrivial changes, add a blank line and a body
wrapped at 72 characters. Small changes need only a subject.
Follow the hybrid commit rules in `.trellis/spec/backend/quality-guidelines.md`.
Follow `.github/workflows/release.yml` for releases.
<!-- TRELLIS:START -->
# Trellis Instructions

These instructions are for AI assistants working in this project.

This project is managed by Trellis. The working knowledge you need lives under `.trellis/`:

- `.trellis/workflow.md` — development phases, when to create tasks, skill routing
- `.trellis/spec/` — package- and layer-scoped coding guidelines (read before writing code in a given layer)
- `.trellis/workspace/` — per-developer journals and session traces
- `.trellis/tasks/` — active and archived tasks (PRDs, research, jsonl context)

If a Trellis command is available on your platform (e.g. `/trellis:finish-work`, `/trellis:continue`), prefer it over manual steps. Not every platform exposes every command.

If you're using Codex or another agent-capable tool, additional project-scoped helpers may live in:
- `.agents/skills/` — reusable Trellis skills
- `.codex/agents/` — optional custom subagents

Managed by Trellis. Edits outside this block are preserved; edits inside may be overwritten by a future `trellis update`.

<!-- TRELLIS:END -->
