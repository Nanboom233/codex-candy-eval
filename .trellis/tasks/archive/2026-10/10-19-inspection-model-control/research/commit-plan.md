# Work commit record

This proposal covers the completed model-control feature and its Go installer delivery in the same task.
The complete strict-style README rewrite passed author checks and independent documentation review. Its final content is staged.
See [the rewrite evidence](readme-strict-implementation.md) and [the review](readme-strict-check.md).
The Chinese README gives user instructions for installation, evaluation, inspection, recovery, and practical builds.
The staged README includes the user's title change and removal of the separate Docker-permission sentence.
Its usage and acknowledgments use `cdd4fb9:README.md` as their source, adapted to the current Sub2API flow.
It omits product version descriptions, upgrade descriptions, product decisions, and internal implementation details. Focused documentation review passed.
The Settings-entry data-loading repair passed the UI regression, actual browser flow, and final delivery build.
Independent review passed. See [the repair evidence](settings-entry-implementation.md) and [the review](settings-entry-check.md).
The installer correction and scope review are complete.
Model control and priority now apply after each account completes detection and existing confirmation rules establish its state.
This includes stable states and skips matching writes.
The correction passed focused tests, affected-package vet, the final build, and independent full-scope review.
See [the implementation evidence](inspection-timing-implementation.md) and [the independent review](inspection-timing-check.md).
The installer candidate excludes argument validation, added SIGTERM behavior, and duplicate payload checks. It retains the original HTTP redirect behavior.
The user approved this exact message and work scope, including the final README they staged.
Git created signed work commit `c4cc8e8`. Its signature verified successfully.

## Committed files

- `.github/workflows/release.yml`
- `.trellis/spec/backend/database-guidelines.md`
- `.trellis/spec/backend/directory-structure.md`
- `.trellis/spec/backend/index.md`
- `.trellis/spec/backend/minimal-linux-delivery.md`
- `.trellis/spec/backend/native-inspection.md`
- `.trellis/spec/backend/quality-guidelines.md` (installer hunk only)
- `.trellis/spec/backend/runtime-installation.md` (Settings-entry lifecycle hunks only)
- `AGENTS.md`
- `README.md`
- `licenses/THIRD_PARTY_NOTICES.txt`
- `patches/0005-inspection.patch`
- `patches/0006-evaluation-ui.patch`
- `patches/0007-sub2api-integration.patch`
- `scripts/build.py`
- `scripts/install.py` (delete)
- `scripts/installer/main.go` (add)
- `tests/backend/internal/codexeval/inspection_test.go`
- `tests/backend/internal/handler/admin/codex_inspection_handler_test.go`
- `tests/backend/internal/service/codex_inspection_models_test.go`
- `tests/backend/internal/service/codex_inspection_settings_test.go`
- `tests/installer/installer_test.go` (add)
- `tests/installer/test_native_version.py` (delete)
- `tests/ui/sub2api-dropdowns.cjs`

Only model-control and Settings-entry repair changes are staged from the shared UI patch and UI test file.
The earlier degraded-tag changes remain outside this commit.
The existing return-to-main guideline hunk and degraded-tag runtime-spec changes remain unstaged.
The user's separate deletion of the repository instructions in `AGENTS.md` remains unstaged.
Task artifacts remain under the Trellis lifecycle owner.

## Exact message

```text
feat(codex-eval)!: control degraded models and add a Go installer

- Apply model control and priority after confirmed inspections.
- Restore removed models on recovery and skip unchanged writes.
- Use native model controls and fix Settings-first group loading.
- Replace Python/archive delivery with a self-extracting Go installer.
- Keep numbered container selection, backups, and failure recovery.
- Write a Chinese user guide with original usage and credits.

BREAKING CHANGE: Replace the Python installer and archive asset.
Run sub2api-codex-eval-installer without installation arguments.
```

## Signing and state

The effective Git setting is `commit.gpgSign=true`.
Use the configured signing behavior without an override.
The current branch is `feat/inspection-model-control`.
The staged set contains 24 files. It excludes generated artifacts and Trellis bookkeeping.
The corrected installer entry and redirect checks passed. Other valid installer and UI evidence remains applicable.
The confirmed-state correction passed focused `CodexInspection` tests and vet in the three affected backend packages.
The final normal build regenerated the raw server and embedded installer.
The new installer was not installed into an actual Docker target. Transaction checks simulated Docker and used local HTTP.
Final user approval and the signed work commit are complete.
Follow the native Trellis finish workflow and return the completed task commits to local `main`.

## Excluded work and finish order

Keep these unstaged hunks outside the work commit:

- `.trellis/spec/backend/quality-guidelines.md`: the existing Return to Main section.
- `AGENTS.md`: the separate deletion of the repository instructions.
- `.trellis/spec/backend/runtime-installation.md`: degraded-tag behavior.
- `patches/0006-evaluation-ui.patch`: the degraded tag.
- `tests/ui/sub2api-dropdowns.cjs`: degraded-tag assertions.

Keep `.trellis/tasks/10-18-inspection-degraded-tag/` active. Do not archive it during this task's finish.
After the approved signed work commit, archive only `10-19-inspection-model-control` through Trellis.
Record the session through Trellis, then return its completed commits to local `main`.
Keep unrelated working changes and the existing signing policy. Do not publish or push.
