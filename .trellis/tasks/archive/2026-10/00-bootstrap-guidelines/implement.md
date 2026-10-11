# Guideline Bootstrap Execution Plan

## Planning Gate

- [x] Reuse the existing task with user consent.
- [x] Inspect project rules, module metadata, plugin entry points, persistence, scheduling, and release checks.
- [x] Write requirements, design, and this execution plan.
- [x] Present the final planning summary and receive explicit approval.
- [x] Check session identity and activate `00-bootstrap-guidelines` with the supported runtime.

## Execution

- [x] Load `trellis-before-dev` and read the existing guideline files.
- [x] Inspect representative Go tests, evaluation runners, credential handling, embedded UI modules, and diagnostic paths.
- [x] Fill `directory-structure.md` with real ownership boundaries and file placement examples.
- [x] Fill `database-guidelines.md` with JSON persistence rules and relevant tests.
- [x] Fill `error-handling.md` and `logging-guidelines.md` with observed error and diagnostic behavior.
- [x] Fill `quality-guidelines.md` with source-backed conventions, test isolation, concurrency rules, ABI boundaries, and local checks.
- [x] Update `index.md` with correct links, scope, pre-development checks, and quality checks.
- [x] Load `trellis-check` and review the documents.
- [x] Update the PRD acceptance checklist and this plan from verified evidence.
- [x] Create and verify signed work commit `2e01ddb`.
- [x] Correct the spec to keep Trellis-managed artifacts in their native workflow.
- [x] Resolve the inclusion of machine-specific AI configuration. The user chose portable files only.
- [x] Review and validate portable Trellis and AI installation files.
- [x] Approve and create signed ordinary work commit `f84347d`.
- [x] Hand off task archival and session records to native Trellis finish-work.

## Validation

- Confirm all six existing backend spec documents are complete and use English.
- Search `.trellis/spec/backend/` for `To fill`, `To be filled`, and `TODO: fill`.
- Verify that visible sections contain project-specific content. Preserve existing HTML comments.
- Verify local links and each referenced source or test path.
- Confirm no new runtime framework, logger, ORM, or spec layer is claimed.
- Run `python3 ./.trellis/scripts/task.py validate 00-bootstrap-guidelines` where applicable. Main-session execution does not need JSONL manifests.
- Run `go vet ./...`.
- Run `go test -race ./...`.
- Run `node scripts/test-ui.cjs`.
- Run `git diff --check` and review tracked changes and task-owned untracked files.
- Verify the selected task and its status after activation.

## Scope Checks

Edit the project specs and current task artifacts to record the hybrid workflow. The user approved review of existing Trellis and AI installation files for the next commit.

Do not change business code, dependencies, deployment settings, or runtime behavior. Do not stage task or workspace bookkeeping in the ordinary work commit.

Include only reviewed installation files. Resolve machine-specific configuration before staging it. Preserve ignored identity, session pointers, private values, and caches.

Do not read `.env` or run deployments. Do not invoke live model requests. Follow `/git-commit` approval for ordinary work commits and native Trellis gates for lifecycle steps.

Do not delete existing files or reset the working tree for recovery.

## Verification Results

| Check | Result |
| --- | --- |
| Six backend spec documents | Complete, with source and test references |
| Unfinished content markers and semicolons | No matches in backend specs |
| Spec index and source links | Reviewed against the repository file list |
| `/output` wording | Reviewed for short sentences, active voice, and consistent terms |
| `go vet ./...` | Passed |
| `go test -race -count=1 ./...` | Passed |
| `node scripts/test-ui.cjs` | Passed |
| `task.py validate 00-bootstrap-guidelines` | Passed. Missing JSONL manifests were skipped. |
| Tracked diff whitespace | Passed |
| Task-owned untracked Markdown whitespace | Passed through per-file Git checks |
| Product source, dependencies, installers, and release workflow | No tracked diff |
| Task activation | Current task selected with status `in_progress` |

Existing HTML comments remain unchanged. The new content fills their sections. The optional Metabase style include was unavailable. The review used `/output` rules.

No production deployment, live evaluation, or C shared-library build ran. The follow-up tracks existing Trellis runtime files and portable AI integration.

## Installation Verification

- `verify_installation.py` checks the staged installation, not live AI sessions.
- Python syntax passed for 34 files.
- JSON parsing passed for three files.
- TOML parsing passed for four files.
- Spec links resolve to files tracked in the proposed commit.
- Task records, personal journals, local Devin settings, and the Codex machine-only MCP block are not staged.
- Git whitespace checks passed after agent line-ending normalization and template formatting fixes.
- Template differences were reviewed as platform command variants, existing model pins, Codex hook variants, and stale memory-tool documentation.
- Memory-tool documentation now matches `trellis mem list --help` platform filters and its OpenCode phase limitation.
- Go vet, uncached race tests, and Node UI checks passed again.

The static workspace root guide belongs to the installation commit. The generated developer index and journal remain for the native session command.

A separate parallel task appeared during verification. It is excluded from the work commit and bootstrap archival. No cross-platform AI CLI end-to-end tests ran.

## Completion State

The work commits are complete: `2e01ddb` and `f84347d`. Both passed local GPG verification. Machine-only settings remain outside the commits.

The user authorized native Trellis finish-work. Task and workspace changes stay in that lifecycle. Native commands own the archive move and journal commit.

Read `task.json` and the session journal for the final lifecycle state. The separate parallel task remains untouched.
