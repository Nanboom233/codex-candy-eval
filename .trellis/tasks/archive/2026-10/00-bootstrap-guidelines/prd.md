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

Keep planning and implementation separate. After approval, select the existing task with the supported Trellis runtime. Verify the guidelines and report checks and remaining lifecycle blockers.

Use `/git-commit` review for ordinary work commits. Use the native Trellis workflow for task artifacts, archive commits, and session records.

### R6: Complete the Trellis installation commit

Correct the spec to describe this hybrid ownership boundary. Include the existing portable Trellis runtime and AI integration files in the next work commit. Preserve local configuration until its inclusion is resolved.

Do not change plugin code, Git signing configuration, or generated runtime behavior. Keep ignored identity, session state, caches, and private configuration out of the commit.

## Acceptance Criteria

- [x] AC1 / R1: The five existing guidelines and index contain usable repository-specific guidance.
- [x] AC2 / R2: Each guideline cites real source or tests and includes relevant local examples and mistakes to avoid.
- [x] AC3 / R3: Persistence guidance describes JSON storage, lock ownership, bounded history, atomic replacement, and read-error protection.
- [x] AC4 / R3: Diagnostics guidance matches current error reporting and does not claim a nonexistent logger.
- [x] AC5 / R1-R2: No unfinished content markers remain in `.trellis/spec/backend/`. Index links resolve to the final guideline files.
- [x] AC6 / R4: Product files and pre-existing unrelated changes remain unchanged by this task.
- [x] AC7 / R5: Required planning artifacts exist before activation, and final verification results and lifecycle state are explicit.
- [x] AC8 / R5-R6: The spec separates ordinary work commits from native Trellis artifact handling.
- [x] AC9 / R6: Signed work commit `f84347d` includes the approved Trellis and AI installation files without private or runtime state.

## Out of Scope

- Product features, bug fixes, and frontend redesign.
- New spec layers, databases, and logging dependencies.
- Trellis runtime changes and hook changes.
- Deployments, releases, and live quota use.
- Disabling or replacing Trellis native lifecycle auto-commits.

## Review Status

The user approved the final planning summary. The supported runtime selected this task and changed its status to `in_progress`.

The user then requested `/output` style for documents and Chinese replies. This changes wording, not the approved scope.

The user approved the signed work commit `2e01ddb`. Its signature passed local verification.

The user then clarified the hybrid workflow. Trellis-managed artifacts retain their native lifecycle. The next work commit must include Trellis and AI integration files.

The spec now records this boundary. The user chose portable files only.

The next work commit excludes the Codex machine-only MCP block and both Devin local configuration files. Task and session records remain in the native Trellis lifecycle.

The user approved the 152-file work commit and native finish-work. Commit `f84347d` passed local GPG verification.

A separate parallel task appeared during verification and remains outside this task. Native finish-work will archive only the bootstrap task and record this session.
