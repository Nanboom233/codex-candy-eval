# Guideline Bootstrap Design

## Boundaries

Keep the current single-repository layout and `backend` spec layer. This task produces documentation, not changes to the plugin architecture.

| Guideline | Content | Primary evidence |
| --- | --- | --- |
| `directory-structure.md` | C ABI entry point, plugin package, embedded UI and data, tests, scripts, release artifacts | `AGENTS.md`, `cmd/cpa-codex-candy-eval/`, `internal/plugin/plugin.go` |
| `database-guidelines.md` | JSON snapshots, mutex ownership, history limits, atomic writes, read-error protection, failed-write rollback | `internal/plugin/state.go`, `internal/plugin/plugin_test.go` |
| `error-handling.md` | Host envelopes, management JSON errors, request validation, panic boundary, worker errors | `internal/plugin/plugin.go`, `internal/plugin/runs.go`, `cmd/cpa-codex-candy-eval/cabi.go` |
| `logging-guidelines.md` | Existing user-visible diagnostics, host failures, storage errors, sensitive-data exclusions | `internal/plugin/state.go`, evaluation runners, `AGENTS.md` |
| `quality-guidelines.md` | Go conventions, shared state, test isolation, embedded UI checks, ABI and release checks | Go tests, `scripts/test-ui.cjs`, `.github/workflows/release.yml`, `AGENTS.md` |
| `index.md` | Scope, links, pre-development checklist, applicable quality checks | Final guideline set |

Retain `database-guidelines.md` for link compatibility. Its title and content will identify JSON persistence. No SQL, ORM, or migration rules apply.

## Evidence Flow

1. Read existing project rules.
2. Inspect the implementation and tests.
3. Write each guideline from this evidence.
4. Cite the source for each important rule.
5. Use short snippets only when they clarify local behavior.

The backend index will explain that the UI ships as embedded plain JavaScript, HTML, and CSS. This task does not introduce a separate frontend spec tree.

## Installation Follow-up

The user requested a second work commit for the existing Trellis runtime and AI integration files. Preserve their current behavior.

Include portable runtime scripts, workflow templates, guides, agents, skills, and reviewed platform configuration. Keep task and session bookkeeping in the native Trellis lifecycle.

The user chose portable files only. Stage the Trellis part of `.codex/config.toml`, not its machine-only MCP block. Leave `.devin/config.local.json` and `.devin/hooks.v1.json` untracked.

Include the static workspace root guide with the installation. Leave the developer journal and generated personal index for the native session command.

Normalize Devin agent line endings through `.gitattributes`. Preserve their rules and permission lists. Do not edit Git configuration.

## Compatibility

Preserve file paths, module layout, plugin versions, management contracts, and the generated thinking guides. No runtime data or deployment state changes.

The current session uses the main-session workflow. Do not dispatch subagents or create subagent manifests for this documentation task.

## Task Lifecycle

The existing bootstrap record returns to `planning` because its initial status did not reflect completed planning. Activate it only after approval of the final summary.

Use supported Trellis commands to select the task. Check session identity before activation. If the current host lacks an identity, request the required local setup rather than edit runtime state or install hooks.

## Verification

Review all guideline claims against source and tests. Check unfinished markers, index links, cited repository paths, task metadata, and the change boundary.

Run these local checks when the required tools are available:

1. `go vet ./...`
2. `go test -race ./...`
3. `node scripts/test-ui.cjs`

Report unavailable tools or pre-existing failures separately. Documentation and local setup are exempt from deployment under `AGENTS.md:62-65`.

## Risks and Recovery

- Existing uncommitted files belong to the user. Preserve them and review only this task's edits.
- Template headings can suggest nonexistent technologies. Replace their content with observed patterns rather than introduce those technologies.
- Trellis owns archive and journal commits. Use its native finish-work gates, not extra `/git-commit` review for each generated commit.
- If a guideline is wrong, correct only that document. Do not reset the repository or overwrite unrelated content.
