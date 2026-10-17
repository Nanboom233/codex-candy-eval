# Work commit review

This proposal contains the complete Sub2API v0.2.17 migration through U22.
U19 removes retired attempts and older task snapshots through macOS Trash under the user's corrected layout instruction.
Main synchronized the work candidate after independent U22 review.
The user approved the exact candidate and message. Signed work commit: `6cac3a2788750b06c8b553315b8593c96bc35a3b`.
Git verified a good GPG signature. No push ran.
U20 removes obsolete CPA sources, tools, module files, examples, screenshots, and the old configuration template.
U22 replaces the U21 source overlay with seven functional production patches and a series list.
All tests live under `tests/`. Only `scripts/build.py` and `scripts/install.py` remain as command entries.
User instructions live in README. The `docs/` tree is removed.
The sole notice source is `licenses/THIRD_PARTY_NOTICES.txt`.
The candidate contains 44 added or modified paths and 56 removed CPA paths before Git rename detection.
Required licenses remain. The archive uses `LICENSE/Project-MIT.txt` and the unchanged native executable.
Review those exact staged changes before final work-commit approval.
The settings tab uses automatic persistence and has no save button or account list.
Inspection shares Candy account selection. Its model, effort, count, and flow remain independent.
The account header contains the next inspection time. The lower inspection card is removed.
Both inspection buttons remain at the bottom of settings, outside collapsed fields.
Saving captions, completion captions, and account inspection summaries are removed.
The staged set includes the U14 request-lifetime and U15 shadow-parent regression tests.
The review found no additional product defect in the checked paths.
It did not establish a fix for the user's cancellation report.
See `final-U18-wrap-check.md` for the complete review and evidence boundaries.
The earlier history `context canceled` report remains open because its trigger is unknown.
This review replaces all earlier commit proposals and records the approved work commit.

## Exact commit message

```text
feat(sub2api): port the complete evaluation plugin to v0.2.17

Port Candy, Fingerprint, ModelTrace, and Pelican evaluations.
Integrate native administrator sending, sidebar, and iframe access.
Preserve account controls, test execution, cancellation, and history.
Add native Candy inspection with autosaved PostgreSQL settings.
Include scheduling, priority adjustment, and ntfy/Bark notifications.
Ship a prebuilt Linux/amd64 package with binary replacement and rollback.
Add focused regression tests and installation documentation.
Remove obsolete CPA sources, tools, and deployment documentation.
```

## Signing

Effective Git configuration: commit.gpgSign=true.
Use plain git commit with the existing signing configuration; no override or push.

## Exact candidate work files

- `M .github/workflows/release.yml`
- `M .gitignore`
- `M .trellis/spec/backend/database-guidelines.md`
- `M .trellis/spec/backend/directory-structure.md`
- `M .trellis/spec/backend/error-handling.md`
- `M .trellis/spec/backend/index.md`
- `M .trellis/spec/backend/logging-guidelines.md`
- `A .trellis/spec/backend/minimal-linux-delivery.md`
- `A .trellis/spec/backend/native-inspection.md`
- `M .trellis/spec/backend/quality-guidelines.md`
- `A .trellis/spec/backend/runtime-installation.md`
- `M AGENTS.md`
- `M README.md`
- `A licenses/ModelTrace-MIT.txt`
- `A licenses/Renderers-MIT.txt`
- `A licenses/Sub2API-LGPL-3.0.txt`
- `A licenses/THIRD_PARTY_NOTICES.txt`
- `A patches/0001-evaluation-core.patch`
- `A patches/0002-evaluation-data.patch`
- `A patches/0003-vendor-assets.patch`
- `A patches/0004-admin-channel.patch`
- `A patches/0005-inspection.patch`
- `A patches/0006-evaluation-ui.patch`
- `A patches/0007-sub2api-integration.patch`
- `A patches/series`
- `A scripts/build.py`
- `A scripts/install.py`
- `A tests/backend/internal/codexeval/candy_test.go`
- `A tests/backend/internal/codexeval/codex_test.go`
- `A tests/backend/internal/codexeval/fingerprint_test.go`
- `A tests/backend/internal/codexeval/inspection_test.go`
- `A tests/backend/internal/codexeval/modeltrace_test.go`
- `A tests/backend/internal/codexeval/pelican_test.go`
- `A tests/backend/internal/codexeval/runtime_test.go`
- `A tests/backend/internal/codexeval/testdata/modeltrace-golden.json`
- `A tests/backend/internal/handler/admin/codex_inspection_handler_test.go`
- `A tests/backend/internal/service/account_codex_eval_service_test.go`
- `A tests/backend/internal/service/codex_eval_lifetime_test.go`
- `A tests/backend/internal/service/codex_inspection_settings_test.go`
- `A tests/installer/native_fixture_start.sh`
- `A tests/installer/native_fixture_supplier.py`
- `A tests/installer/native_inspection_fixture.py`
- `A tests/installer/test_native_version.py`
- `A tests/ui/sub2api-dropdowns.cjs`

## Exact cleanup removals

The original paths below leave the repository. Git can present migrated native files as renames.
The existing `install.sh` and `install.ps1` deletions are old CPA installers and belong to this cleanup.

- `D .env.example`
- `D cmd/cpa-codex-candy-eval/cabi.go`
- `D cmd/cpa-codex-candy-eval/main.go`
- `D docs/images/example.png`
- `D docs/inspection.md`
- `D examples/inspection.py`
- `D go.mod`
- `D go.sum`
- `D install.ps1`
- `D install.sh`
- `D internal/plugin/candy.go`
- `D internal/plugin/codex.go`
- `D internal/plugin/codex_test.go`
- `D internal/plugin/credentials.go`
- `D internal/plugin/credentials_test.go`
- `D internal/plugin/data/codex_context.json`
- `D internal/plugin/data/fingerprint_baselines.json`
- `D internal/plugin/data/fingerprint_probes.json`
- `D internal/plugin/data/modeltrace/LICENSE`
- `D internal/plugin/data/modeltrace/README.md`
- `D internal/plugin/data/modeltrace/unified_bank.json`
- `D internal/plugin/fingerprint.go`
- `D internal/plugin/fingerprint_algorithm.go`
- `D internal/plugin/fingerprint_test.go`
- `D internal/plugin/inspection.go`
- `D internal/plugin/inspection_test.go`
- `D internal/plugin/model.go`
- `D internal/plugin/model_test.go`
- `D internal/plugin/modeltrace.go`
- `D internal/plugin/modeltrace_algorithm.go`
- `D internal/plugin/modeltrace_test.go`
- `D internal/plugin/pelican.go`
- `D internal/plugin/pelican_test.go`
- `D internal/plugin/plugin.go`
- `D internal/plugin/plugin_test.go`
- `D internal/plugin/runs.go`
- `D internal/plugin/state.go`
- `D internal/plugin/testdata/modeltrace-golden.json`
- `D internal/plugin/web/app.js`
- `D internal/plugin/web/candy.js`
- `D internal/plugin/web/catalog.js`
- `D internal/plugin/web/components.js`
- `D internal/plugin/web/credentials.js`
- `D internal/plugin/web/fingerprint.js`
- `D internal/plugin/web/modeltrace.js`
- `D internal/plugin/web/pelican.js`
- `D internal/plugin/web/style.css`
- `D internal/plugin/web/theme.css`
- `D internal/plugin/web/ui.html`
- `D internal/plugin/web/vendor/markdown-it-cjk-friendly.min.js`
- `D internal/plugin/web/vendor/markdown-it.min.js`
- `D internal/plugin/web/vendor/temml.css`
- `D internal/plugin/web/vendor/temml.min.js`
- `D internal/plugin/web/vendor/texmath.js`
- `D scripts/deploy.sh`
- `D scripts/test-ui.cjs`

## Reviewable diff

Use `git diff --cached` and `git diff` to inspect the candidate and its current staging boundary.
Main replaced stale U21 staged paths with the current patch series, scripts, and separate tests.
The duplicate root notice is absent from the index. Its sole source appears in the candidate list above.
The staged set matches the exact 100-path candidate, including 56 original CPA removals before rename detection.
Include current specs, documentation, package metadata, `.gitignore`, and the exact removal list in the candidate review.
These changes establish Sub2API-only owners and commands, required attribution, and the approved hybrid commit workflow.
The generated binary/archive stay in ignored `dist/build`.

The original renderer/license line endings, original page whitespace and unified patch context are preserved.
A broad staged whitespace check flags those inherited bytes and patch syntax; it is not a new acceptance gate.
Authored native code is formatted and relevant syntax/tests/vet/build passed.

## Verification and limits

U20 removed the old CPA source graph and current references after dependency review.
Independent cleanup review: `U20-CPA-cleanup-check.md`.
The regenerated archive contains seven entries and the unchanged U18 executable.
The current candidate is 44 added/modified and 56 removed paths, excluding local machine settings.
Git can display matching original and current test or license files as renames.
U22 implementation and independent review: `U22-patch-series-implementation.md` and `U22-patch-series-check.md`.
Effective production and test contents remain unchanged. Patch composition, existing focused tests, commands, and archive checks passed.
The archive was regenerated with the existing executable. No product rebuild or Docker operation ran for U22.

Final static Linux/amd64 binary and package: `dist/build/sub2api-codex-eval` and the same-name `.tar.gz`.
Earlier core and U1-U9 reports remain evidence for unchanged evaluation, installer progress, `/tmp` backup, and rollback paths.
See `fresh-native-check.md`, `native-browser-model-check.md`, and `installer-disabled-model-check.md` for their scope and limits.
Final v0.2.17 implementation: `native-inspection-implementation.md`.
Independent full-scope review: `native-v0217-integrated-check.md`.
Actual browser and final rebuilt-candidate smoke: `native-v0217-browser.md`.
The review repaired initial settings-load admission and aligned the local check entry with pnpm 9.15.9.
The final scoped check and package build passed after those changes.
The final application reported v0.2.17 and completed synthetic Candy inspection from the actual settings page without errors.
The browser console reported zero errors and warnings.
The main session removed 11 earlier task containers and retained `sub2api-codex-eval-test` with all its test dependencies.
U10 implementation and build evidence: `U10-dropdown-implementation.md`.
Independent U10 review: `U10-dropdown-check.md`.
Actual model/account selection, custom-model and wildcard save/reload, and final viewport repair: `U10-dropdown-browser.md`.
U11 supersedes the account-picker and wildcard configuration from that earlier report.
U11 implementation and independent review: `U11-settings-implementation.md` and `U11-settings-check.md`.
U12 implementation and independent review: `U12-settings-collapse-implementation.md` and `U12-settings-collapse-check.md`.
Final built-page collapse, switch captions, priority values, saved hidden fields, and narrow layout: `U11-U12-browser.md`.
The fixed binary and archive include all corrections through U18.
U14 settings ordering and both supported request-lifetime paths passed focused tests.
U15 shared selection, independent inspection flow, and native persistence passed Go tests and touched-package vet.
Its retained-container timer run completed after the browser closed and selected only current Candy-group accounts.
U16 keeps inspection buttons at the bottom of settings.
U17 removes added save captions, completion captions, and account summaries while retaining actual errors.
U18 moves the next-run time into the existing Candy account header and removes the separate lower runtime card.
The final U18 desktop and narrow screenshots were inspected. Neither layout had horizontal document overflow.
See `U14-settings-history-check.md`, `U15-shared-Candy-schedule-check.md`, and `U15-shared-Candy-browser.md`.
See the U16, U17, and U18 implementation, independent-check, and browser reports for later corrections.
The final full-scope review is `final-U18-wrap-check.md`. No product fix or repeated suite was required.
Current archive inventory contains seven runtime/license entries. The binary is a static Linux/amd64 ELF executable.
The main session confirmed only one task container remains running. Unrelated `kevbot` resources remain outside this task.
The unchanged native shell can report an inline configuration nonce error after a cached 304 reload.
Fresh no-store HTML uses a matching CSP nonce. The evaluation iframe and U10 controls passed.
N150/iStoreOS and production HTTPS/proxy execution were not performed.
The historical `context canceled` event remains unresolved. Ending start/state requests preserves background evaluations in both protocol tests.
Explicit stop cancels supplier work and records the actual cancellation. These tests do not identify the earlier event.
Preserve that report and its history. Do not claim the reported cause was fixed.
Apply approved A1-A12 without adding hardware, production, or historical-cause acceptance gates.

## Files outside this work commit

The following existing dirty paths remain excluded from the work commit:

```text
.codex/config.toml
.devin/config.local.json
.devin/hooks.v1.json
```

Generated caches, browser output, and `dist/` remain ignored and outside work commits.
The active task directory contains Trellis-managed plans and research reports. The rejected history directory goes to macOS Trash.
Use native Trellis task handling for those records. Do not stage them into the product work commit.

## Workflow completion

The completed U20, U21, and U22 reviews supported the exact 100-path work candidate.
The user approved the displayed work message, existing signature policy, and native archive/session sequence.
The signed work commit is `6cac3a2788750b06c8b553315b8593c96bc35a3b`. Git verified its good signature.
Native Trellis task archival and session recording follow the work commit.
Use the original acceptance criteria for closure. Evidence limits do not add requirements.
