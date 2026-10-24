# Quality Guidelines

## Current Sub2API scope

Use [administrator evaluation](./runtime-installation.md), [native inspection](./native-inspection.md), and [binary delivery](./minimal-linux-delivery.md) for current contracts.
The delivery targets Sub2API v0.2.18 on Linux/amd64.
Installer version mismatches or unresolved probes produce advisory warnings and explicit confirmation, without a release allowlist.
Use the original Sub2API source tree and this repository's native evaluation module.
Active integration lives in contextual `patches/*.patch` files. All test cases and support code live in flat `tests/`.
Keep `assets/` for data resources only.
Use [Directory Structure](./directory-structure.md) for file ownership.

For native behavior changes, use the approved focused checks:

```sh
python3 scripts/build.py --source /path/to/original-sub2api --check
```

Reuse completed baseline and unchanged integration evidence.
Add only checks needed for changed behavior.
Format changed Go files with `gofmt`.
Check affected scripts, resources, and references.
Directory and documentation changes do not require another product build.
Report actual failures, skipped checks, and execution limits.

Use `tests/test_ui.cjs` directly for existing native page assertions.
Native Go cases live in `tests/test_*_test.go` with `// sub2api-package:` context annotations.
Use native `go test` against the temporary target package. Execute discovered cases once and stop on failure.
Keep execution support in `tests/support.py` and manual fixture helpers in flat `tests/native_*` files.
The common build repository owns the authoritative fake-Docker and local HTTP installer suite.
Standalone `--check` runs this plugin's backend/UI tests and vet, then compiles the standalone server and installer.
Keep host test execution separate from Linux/amd64 release flags. Preserve standalone use without requiring a parent checkout.
No retained Docker container or manual browser session is required for the approved build migration.
Do not create another suite for retired interfaces.

Preserve native administrator authentication, account transport, cancellation, history errors, and settings storage.
Keep tests synthetic. Do not expose credentials in fixtures, logs, or reports.
Use existing result and error owners.
Review handler, service, runtime, and page fields together when a contract changes.
Keep generated output outside the user package and work commits.

Compatibility profiles, legacy migrations, shims, and multi-version adapters are retired mechanisms.
Gateway private keys, directed billing, activation, maintenance, and universal builders are retired mechanisms.
Race, coverage, full-suite, architecture matrices, and source-resource parity are not current acceptance requirements.

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

### Return to Main

If work uses a task branch, merge its completed commits into local `main` before the final handoff.
Include the work, archive, and session commits. Switch to local `main` after the merge.

References:

- [Trellis workflow](../../workflow.md): work commits, archive, and session order.
- [Lifecycle configuration](../../config.yaml): artifact auto-commit settings.
- [Archive implementation](../../scripts/common/task_store.py): scoped task commits.
- [Session implementation](../../scripts/add_session.py): scoped workspace commits.
- [Repository commit rules](../../../AGENTS.md): imperative subjects and body wrapping.

The user clarified this ownership boundary during the bootstrap task. Native artifact handling is not an exception to Trellis. It is part of the hybrid workflow.
