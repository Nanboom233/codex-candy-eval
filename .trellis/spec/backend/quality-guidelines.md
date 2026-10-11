# Quality Guidelines

> Preserve host contracts, shared-state safety, and existing verification gates.

---

## Overview

<!--
Document your project's quality standards here.

Questions to answer:
- What patterns are forbidden?
- What linting rules do you enforce?
- What are your testing requirements?
- What code review standards apply?
-->

Follow adjacent Go code and the standard-library helpers already used here. Check `go.mod` before adding imports or dependencies.

The module declares Go 1.24. The release workflow configures no separate backend lint framework or frontend package-manager build.

Use `gofmt` on changed Go files. Update behavior tests when code behavior changes.

Evidence: [repository checks](../../../AGENTS.md), [module metadata](../../../go.mod), and [release check job](../../../.github/workflows/release.yml).

---

## Forbidden Patterns

<!-- Patterns that should never be used and why -->

- Do not call the host while holding the histories mutex. Host requests can block and interfere with result updates or shutdown.
- Do not bypass `credentialBusyLocked` when starting evaluations. One credential cannot run multiple test types at once.
- Do not start tracked workers without matching `tasks.Add` and `tasks.Done`.
- Do not use `t.Parallel` in tests that replace package globals through `setupTest`.
- Do not return a different management error shape from the one consumed by the browser.
- Do not edit vendored renderers or add a framework as an unrelated cleanup.
- Do not change ABI callbacks or C buffer ownership without checking the CPA host contract.
- Do not expose secrets in tests, logs, fixtures, or deployment reports.

Evidence: [stateResponse and Quiesce](../../../internal/plugin/plugin.go), [startRuns](../../../internal/plugin/runs.go), [test setup](../../../internal/plugin/plugin_test.go), and [C ABI bridge](../../../cmd/cpa-codex-candy-eval/cabi.go).

---

## Required Patterns

<!-- Patterns that must always be used -->

### Shared state and shutdown

Protect result maps, progress, configured credentials, and shutdown state with `mu`. Read the locking contract of each helper before calling it.

`startRuns` selects credentials before locking. Its callbacks run with `mu` held. A start callback launches a worker that eventually calls `tasks.Done`.

`Quiesce` rejects new runs, cancels collection work, releases `mu`, and waits for tracked tasks. Do not wait for workers while holding their required lock.

Evidence: [run callback contract](../../../internal/plugin/runs.go), [candy worker cleanup](../../../internal/plugin/candy.go), and `TestQuiesceDrainsBothTests` in [plugin tests](../../../internal/plugin/plugin_test.go).

### Host and embedding contracts

Keep JSON tags aligned across Go handlers, browser consumers, and persisted records. Register and handle each management route together.

The C bridge uses the `cshared` build tag. Plugin response buffers use C allocation and the plugin free callback. Host response buffers use the host free callback.

When changing embedded assets, check Go embed declarations, replacement markers, script order, and inline-script safety. The browser sources share one assembled page.

Evidence: [C ownership](../../../cmd/cpa-codex-candy-eval/cabi.go), [UI assembly](../../../internal/plugin/plugin.go), and `TestScriptsInlineSafely` in [plugin tests](../../../internal/plugin/plugin_test.go).

---

## Testing Requirements

<!-- What level of testing is expected -->

Use `setupTest(t)` for tests that touch plugin globals. It assigns a temporary state path, installs a failing default host mock, and restores prior globals after workers finish.

Use synthetic host responses rather than live model requests. `streamHost` and `mockModelResponse` demonstrate the streaming host contract.

Use channels and bounded waits to verify concurrent runs and cancellation. Do not let asynchronous workers outlive test cleanup.

Evidence: [plugin test setup](../../../internal/plugin/plugin_test.go), [model tests](../../../internal/plugin/model_test.go), and [fingerprint concurrency tests](../../../internal/plugin/fingerprint_test.go).

For Go code changes, run:

```sh
go vet ./...
go test -race ./...
```

For embedded UI changes, also run:

```sh
node scripts/test-ui.cjs
```

The UI checks use Node built-ins and a mocked DOM through `node:vm`. They do not replace browser verification.

Plugin code, UI, build, or installer changes also require deployment and authenticated end-to-end verification under `AGENTS.md`. Documentation and local setup changes are exempt.

Evidence: [release checks](../../../.github/workflows/release.yml), [UI harness](../../../scripts/test-ui.cjs), and [deployment requirements](../../../AGENTS.md).

---

## Code Review Checklist

<!-- What reviewers should check -->

- [ ] The change stays within its approved scope and preserves unrelated work.
- [ ] Existing helpers own repeated scheduling, model execution, history, and response behavior.
- [ ] Host calls occur outside `mu`, and worker tracking permits shutdown to finish.
- [ ] Persistence failures preserve the contracts in [JSON persistence guidelines](./database-guidelines.md).
- [ ] Go JSON tags, management routes, UI consumers, and response errors agree.
- [ ] Tests isolate package globals and cover changed success and failure behavior.
- [ ] Embedded scripts cannot introduce an early HTML closing script tag.
- [ ] Release tags match `pluginVersion`.
- [ ] Linux builds retain the glibc 2.17 baseline when build work is in scope.
- [ ] Checks appropriate to the changed files have completed, with blockers reported.

Release evidence: [workflow](../../../.github/workflows/release.yml) and [macOS cross-build instructions](../../../AGENTS.md).

Do not claim default Go tests cover the C bridge. Default Go checks exclude that file unless the `cshared` tag is present.

## Commit and Finish Work

Use a hybrid workflow. Trellis owns development phases and managed task artifacts. `/git-commit` owns the review and approval of ordinary work commits.

### Ownership Boundary

| Change | Owner |
| --- | --- |
| Product code, project specs, and reviewed runtime or AI integration changes | Trellis work phase with `/git-commit` review |
| Task planning artifacts, status, and archive moves | Trellis task workflow |
| Session journals and generated workspace indexes | Trellis session workflow |
| Initial installation of portable Trellis and AI integration files | Reviewed work commit |
| Local identity, session pointers, caches, credentials, and machine-only settings | Local files, not a blanket work commit |

Do not mix task or workspace bookkeeping into an unrelated work commit. Do not rewrite an earlier commit to change this boundary.

### Ordinary Work Commits

1. Record source-backed findings in the relevant spec before drafting the commit.
2. Inspect the working tree and staged files.
3. Separate work changes from Trellis bookkeeping and unrelated changes.
4. Stage only the verified work paths.
5. Review the complete staged diff.
6. Draft an English Conventional Commit with an imperative subject.
7. Keep the subject within 74 characters.
8. Wrap body lines at 72 characters to follow `AGENTS.md`.
9. Check the effective Git signing policy.
10. If the signing policy is unset, ask the user to choose signed or unsigned.
11. Show the exact file list, message, and signing behavior.
12. Obtain final approval immediately before each non-temporary work commit.

Useful checks:

```sh
git status --short
git diff --cached --name-status
git diff --cached
git config --get --bool commit.gpgSign
```

For the signing-policy check, exit status `1` means unset. A configured value takes precedence unless the user explicitly chooses otherwise.

Task-creation consent and an earlier commit request do not replace final work-commit approval.

### Trellis-Managed Artifacts

Follow `.trellis/workflow.md` for task planning, completion, archival, and session records. Use the supported Trellis commands for lifecycle operations.

Create work commits before archive and journal commits. Once the workflow authorizes finish-work, use its native artifact and auto-commit steps.

Do not add `/git-commit` confirmation before every generated bookkeeping commit. Do not replace native artifact handling with manual staging or custom commits.

Honor `session_auto_commit` and `session_commit_message` in `.trellis/config.yaml`. Do not disable auto-commit merely to force the ordinary commit workflow.

The archive command commits task changes. The session command commits workspace records. Both use scoped Git pathspecs and preserve unrelated staged work.

Inspect their results after each step. Verify signatures when Git signing is enabled. Stop on signing or lifecycle failure.

Do not change Git configuration or downgrade signing to bypass a failure. Do not push without approval.

References:

- [Trellis workflow](../../workflow.md): work commits, archive, and session order.
- [Lifecycle configuration](../../config.yaml): artifact auto-commit settings.
- [Archive implementation](../../scripts/common/task_store.py): scoped task commits.
- [Session implementation](../../scripts/add_session.py): scoped workspace commits.
- [Repository commit rules](../../../AGENTS.md): imperative subjects and body wrapping.

The user clarified this ownership boundary during the bootstrap task. Native artifact handling is not an exception to Trellis. It is part of the hybrid workflow.
