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
- [x] Record the combined Trellis and `/git-commit` convention in the quality guide.
- [ ] Obtain separate approval for commit, archive, and session-record operations.

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

Only edit `.trellis/spec/backend/` and this existing task's artifacts during implementation. Leave business code, vendored files, Go dependencies, platform integrations, generated runtime scripts, deployment settings, and the user's existing `AGENTS.md` changes untouched.

Do not read `.env`, run deployments, invoke live model requests, commit, or archive without the required approval. Do not delete existing files or reset the working tree for recovery.

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

No production deployment, live evaluation, or C shared-library build ran. This task changes documentation only.

## Completion State

The guideline documents are complete. Commit, archive, and the finish-work journal remain pending separate approval. The task stays selected with status `in_progress`.

Do not run lifecycle auto-commits before that approval. Do not describe the pending lifecycle steps as complete.
