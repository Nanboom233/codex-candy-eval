# Bootstrap Project Guidelines

## Goal

Complete the existing Trellis initialization task with source-backed coding guidelines for this repository. Future sessions must find usable project rules instead of empty templates.

The user approved reuse of `00-bootstrap-guidelines`. Do not create a replacement task. Implementation requires approval of the final planning summary.

## Background

- `AGENTS.md:11-31` defines the project layout and Go checks.
- `go.mod:1-5` defines one Go module, Go 1.24, and the `golang.org/x/text` dependency.
- `cmd/cpa-codex-candy-eval/cabi.go:1-5,70-137` provides the CGO shared-library bridge and host-buffer ownership.
- `internal/plugin/plugin.go:16-23,31-103,141-202` defines versions, embedded UI assembly, and host method dispatch.
- `internal/plugin/runs.go:9-57` defines shared locking and run scheduling.
- `internal/plugin/state.go:13-107,134-159` uses JSON persistence, limits history, protects unreadable files, and rolls back failed history removal.
- `.github/workflows/release.yml:12-24` runs Go vet, race tests, and the Node UI checks.
- At intake, `.trellis/spec/backend/index.md:15-21` listed five template guidelines. A source scan found 29 unfinished markers across these files and the index.
- The generated task had status `in_progress`, but no selected session task and no completed planning artifacts.

## Requirements

### R1: Project-specific guidelines

Fill the five existing backend guideline files and their index in English. Cover directory ownership, JSON persistence, error boundaries, diagnostics, and quality checks. Include embedded UI and C ABI boundaries where they affect backend work.

Use `/output` style for documents. Use short sentences, active voice, and consistent terms. Keep conversation replies in Chinese with the same style.

### R2: Evidence

Support important rules with real repository paths, named symbols, and short examples. Inspect tests before documenting behavior. Describe current patterns, not new architecture or hypothetical tools.

### R3: Accurate persistence and diagnostics

Explain that this repository has no database or ORM. Keep the existing `database-guidelines.md` path and document the actual JSON storage contract there. Describe the actual diagnostic mechanisms without inventing a logging framework.

### R4: Scope and safety

Preserve unrelated working-tree changes. Do not modify product source, embedded assets, dependency files, release workflows, installers, deployment scripts, or platform configuration. Do not read private `.env` values or run live evaluations.

### R5: Workflow completion

Keep planning and implementation separate. After approval, select the existing task with the supported Trellis runtime. Verify the guidelines and report checks and any remaining lifecycle blockers. Do not commit or run auto-commit lifecycle operations without separate approval.

## Acceptance Criteria

- [x] AC1 / R1: The five existing guidelines and index contain usable repository-specific guidance.
- [x] AC2 / R2: Each guideline cites real source or tests and includes relevant local examples and mistakes to avoid.
- [x] AC3 / R3: Persistence guidance describes JSON storage, lock ownership, bounded history, atomic replacement, and read-error protection.
- [x] AC4 / R3: Diagnostics guidance matches current error reporting and does not claim a nonexistent logger.
- [x] AC5 / R1-R2: No unfinished content markers remain in `.trellis/spec/backend/`. Index links resolve to the final guideline files.
- [x] AC6 / R4: Product files and pre-existing unrelated changes remain unchanged by this task.
- [x] AC7 / R5: Required planning artifacts exist before activation, and final verification results and lifecycle state are explicit.

## Out of Scope

- Product features, bug fixes, and frontend redesign.
- New spec layers, databases, and logging dependencies.
- Trellis runtime changes and hook changes.
- Deployments, releases, and live quota use.
- Automatic Git commits.

## Review Status

The user approved the final planning summary. The supported runtime selected this task and changed its status to `in_progress`.

The user then requested `/output` style for documents and Chinese replies. This changes wording, not the approved scope.

The user requested spec records and the combined Trellis and `/git-commit` process. The quality guide now records that convention. Final commit approval and the signing choice remain pending.
