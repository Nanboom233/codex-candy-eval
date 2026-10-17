# Plugin Migration to Sub2API

Status: in_progress. Implementation and delivery checks through U18 passed. The earlier history cancellation report remains open.

Use approved A1-A12 for acceptance. Hardware and production execution limits remain evidence statements.
They do not add acceptance gates. Identifying an old event's cause is not a new closure requirement.
The user's cancellation report remains part of the task. Passed lifecycle tests do not establish that its cause was fixed.

## Current requirement U22: auditable patches and separate tests

The user requested a patch-based repository, rejected a single large patch, and required all tests under `tests/`.
Use multiple production patches with an explicit application order under `patches/`.
Keep all test source, test-only data, fixtures, installer checks, and UI checks under `tests/`.
Production patches contain only runtime source, runtime data, third-party assets, and upstream integration changes.
Remove the production `src/` copy when the patch-series conversion is implemented.
Place the existing installer at `scripts/install.py`. Preserve its filename.
Update repository callers to that path. Preserve `install.py` as the packaged installation entry.
Keep only `scripts/build.py` and `scripts/install.py` as command entries.
The build entry applies patches, builds, and packages. Its `--check` option runs the existing checks.
Remove the redundant shell wrappers and separate packaging entry during conversion.
Keep user instructions in `README.md`. Merge useful rollback and history paths from `docs/sub2api-operations.md`, then remove that duplicate guide.
Remove `docs/` from the repository layout.
Preserve current behavior, fixed output names, installation commands, licenses, and existing verification scope.
The user approved this layout. Implementation and independent review passed.
Production composition, test separation, existing focused checks, consolidated commands, and the seven-entry archive passed.
The existing binary remains unchanged. Evidence: `research/U22-patch-series-implementation.md` and `research/U22-patch-series-check.md`.

## Previous correction U21: make repository roles clear

The user rejected the remaining unclear layout and generated clutter.
Use `src/backend/` and `src/frontend/` for the current source additions.
Use `patches/sub2api.patch` for upstream integration edits.
Keep developer commands in `scripts/`, checks in `tests/`, guides in `docs/`, and licenses in `licenses/`.
Keep the fixed delivery outputs under `dist/build/` and root `install.py` as the installation entry.
Remove generated Python caches, Finder metadata, and old browser output through macOS Trash.
Use `licenses/THIRD_PARTY_NOTICES.txt` as the sole notice document. Remove its duplicate root Markdown copy.
Preserve private files, Git data, Trellis and agent integration directories, licenses, and runtime behavior.
Update current callers and documentation together. Do not add compatibility wrappers or rebuild unchanged product inputs.
Update the packaged notice's source locations and regenerate only the archive with the unchanged executable.
Earlier references to `native/backend`, `native/frontend`, and `native/sub2api.patch` map to these new paths.
U21 is complete and independently checked. Current source bytes and modes remain unchanged after the moves.
The root no longer contains old browser output, Python caches, Finder metadata, or the duplicate notice.
The existing UI harness, installer tests, patch dry run, current links, and archive checks passed.
Evidence: `research/U21-repository-layout-implementation.md` and `research/U21-repository-layout-check.md`.

## Previous correction U20: remove old CPA contents

The user requested removal of old CPA contents after the Sub2API migration.
This instruction replaces earlier requirements to retain the original CPA source and command paths.
Keep the completed native implementation, data, renderer assets, tests, current commands, and fixed delivery outputs.
Move obsolete CPA source, ABI entry points, build/deployment tools, examples, documentation, and unused assets to macOS Trash.
Remove the root Go module and placeholder configuration if they serve only the old CPA build.
Preserve required licenses and attribution, private configuration, unrelated files, and current task research.
Update current documentation and specifications to describe the Sub2API-only repository.
Do not create another archive. Do not change evaluation behavior or rebuild the unchanged binary.
If archive metadata changes, regenerate only the installation archive and check its contents.
Check affected references, syntax, and package inputs. Reuse valid unchanged product evidence.
U20 is implemented and independently checked. Obsolete CPA contents are absent and native product inputs remain unchanged.
The archive was regenerated with `LICENSE/Project-MIT.txt` and all required license texts.
Evidence: `research/U20-CPA-cleanup-implementation.md` and `research/U20-CPA-cleanup-check.md`.

## Previous correction U19: organize the repository before commits

The user requested repository organization before the work commit.
Keep the active Sub2API source, build, installation, tests, and documentation easy to identify.
The original CPA retention decision in this step is superseded by U20.
Move the complete task `history/` directory to macOS Trash, as the user explicitly requested.
This includes abandoned implementations and older task snapshots. Do not create a replacement archive in the repository.
Preserve unrelated files, private configuration, current task artifacts, and research records.
Keep fixed delivery outputs, current command entry points, and runtime behavior unchanged.
Ignore generated caches and browser output. Update affected current documentation and references.
Reuse valid build and runtime evidence. Check only the paths affected by this organization.
This request authorizes organization within the current migration task.
The earlier U19 move into task history was rejected because that directory does not follow the task layout.
The implementation owner moved the whole directory to Trash and removed current archive references.
Earlier U19 reports retain the rejected intermediate state under their correction notes.

## Previous correction U18: use the existing account status bar

Move the next inspection time into the existing Candy account header beside `账户` and `仅启用`.
Remove the separate inspection runtime card below the account list.
Use the existing header or error paths for necessary active feedback and actual errors.
Do not add another column, card, explanatory caption, or move the settings inspection buttons.
Preserve the U17 copy removals and all runtime behavior.
This correction is approved for implementation.
U18 passed independent review, fixed build/package, and desktop/narrow browser checks in the retained container.
Evidence: `research/U18-header-status-browser.md`.

## Previous correction U17: remove added status copy

Remove saving and unsaved-change captions, per-account inspection summaries, and the `巡检完成` caption.
Do not replace them with other captions, tooltips, or explanatory text.
Keep actual errors visible, including notification failures. Keep active progress and the next scheduled time.
Hide completed-only status containers when they have no useful content.
Preserve autosave and the settings-bottom inspection buttons from U16.
This correction is approved for implementation.
U17 passed independent source review, the fixed build/package, and final retained-container browser checks.
See `research/U17-status-copy-browser.md` for the final evidence.

## Previous correction U16: inspection controls and save copy

Restore `立即巡检` and `停止当前巡检` to the settings tab.
The user did not authorize moving these buttons to Candy.
Keep them available while automatic inspection is off.
Keep account results and history in Candy. Do not restore a settings account list.
Remove the successful-save text `已保存`. Keep automatic persistence and actual validation/write errors.
This correction is approved for implementation.
U16 passed independent review and final retained-container browser checks. The fixed installation package includes the correction.

## Current requirement U15: autosave and shared Candy account selection

This section replaces conflicting inspection account-display, selection, and explicit-save requirements below.
The user approved these requirements and instructed the agent to update this task and implement them.

1. Remove the settings save button. Valid edits save automatically through native settings storage.
2. Show actual validation and save failures. Do not silently lose edits or overwrite newer values with old responses.
3. Remove the account list, records, and duplicate result display from the settings tab.
4. Scheduled Candy uses the same selected group and account filters as the Candy page.
5. Its group/filter scope matches Candy `测试全部`. Do not select a separate all-account inventory or apply another Free filter.
6. Keep results and history in Candy. Reuse existing runtime and sender primitives, with no settings-page account display.
7. Keep the independent inspection count and flow: first test, conditional confirmations, recovery wait, priority adjustment, and notifications.
8. Keep backend scheduling active when the browser closes. Persist shared account-selection inputs through native settings.
9. Preserve manual fixed test counts, row/selected/all entry points, native admission, U13 group isolation, and the other three tabs.

The user clarified that reuse concerns account selection, not unifying inspection and manual test counts or orchestration.
Manual fixed runs and inspection `max_tests_per_account` remain independent.
Keep inspection's independent model and effort settings. The clarification only merges account selection, not evaluation configuration or orchestration.
Retain existing native inspection availability, including expiry and quota checks. Do not change manual admission to match those checks.

The settings tab keeps inspection model, effort, maximum tests, scheduling, priority, and notification controls.
The existing Candy page owns the shared account group and filters.
Do not add another inspection account catalog, configuration store, or result system.
Reuse existing checks. Add only focused checks for autosave and the shared scheduled batch.
The earlier history `context canceled` report remains open pending its triggering operation.
It does not block the authorized U15 work or authorize changing real cancellation into success.

U15 passed independent source review, focused checks, a fixed v0.2.17 build, and retained-container browser/timer verification.
See `research/U15-shared-Candy-browser.md` for the delivered candidate and actual account scope.

## Goal and existing inputs

Migrate the plugin functions and management page to Sub2API v0.2.17. Use native administrator account-test services and settings storage.
Deliver a prebuilt Linux/amd64 replacement for N150/8GB iStoreOS. U15 corrects settings persistence and inspection account selection within this migration.

The user changed the target from v0.2.16 to v0.2.17 during approved implementation. Apply the integration to actual v0.2.17 source. Preserve earlier reports as v0.2.16 evidence. Check changed native contracts and the final v0.2.17 paths without adding a repeated pristine baseline.

## Whole-task scope and status

| Workstream | Current state | Decision or evidence |
| --- | --- | --- |
| Original Sub2API compile and baseline checks | v0.2.16 evidence recorded. Reuse only where applicable. | `research/pristine-v0.2.16-baseline.md` |
| v0.2.17 source integration and delivery | Built and checked from actual v0.2.17 source. | `research/native-v0217-integrated-check.md` |
| Four original evaluations and native request transport | Implemented and reviewed within the recorded scope. | `research/fresh-native-implementation.md`, `research/fresh-native-check.md` |
| Original page, native accounts, task controls, and result history | Implemented. Later corrections control current behavior. | F01-F06, C02-C22, U1-U9 |
| Native sidebar, icon, iframe, port handling, and literal models | Implemented. Actual Chrome and local Docker evidence exists. | `research/native-browser-verification.md`, `research/native-browser-model-check.md` |
| Installer progress, `/tmp` backups, and disabled model discovery | Implemented and reviewed. | `research/installer-disabled-model-repair.md`, `research/installer-disabled-model-check.md` |
| Native Candy inspection, settings tab, and database settings | U15 autosave and shared account selection passed source, browser, and timer checks. | Current U15 reports |
| Final integrated delivery | U18 v0.2.17 package rebuilt after independent review and checked on desktop and narrow screens in the retained container. | `research/U18-header-status-check.md`, `research/U18-header-status-browser.md` |
| Repository organization | The complete task history directory moved to Trash. Current archive pointers removed. Product inputs remain unchanged. | Current U19 correction |
| Old CPA removal | Old-only sources, tools, module, examples, screenshots, and template removed. Current references and package checked. | `research/U20-CPA-cleanup-check.md` |
| Commit and task closure | Pending. No commit approval or task completion is inferred. | `task.json` |

The complete decision history and the current interpretation are in `research/migration-history-and-decisions.md`. Keep the original approval files and snapshots. Their old planning labels describe their original stage. They do not cancel later execution approval or prove that later work is implemented.

Use original Sub2API v0.2.17 source, this repository's `src/` additions, and `patches/sub2api.patch` for current builds. The initial port used the original plugin directly. U20 removes that obsolete CPA source copy. Restore upstream source from the official archive when temporary inputs are absent. Do not clone Git history or add source pins, commit exports, manifests, or archive/blob parity gates. Reuse existing dependency caches. Earlier v0.2.16 inputs remain comparison evidence.

Reuse completed pristine build/check evidence from research/pristine-v0.2.16-baseline.md, retaining failures/skips and platform limits; do not repeat the baseline. Old installation is not accepted functionality. Migrated internal/eval, adapter, compatibility, old patches/scripts/tests and generated output are defect evidence only, never code inputs.

## Requirements

R1. Preserve original Candy/Fingerprint/ModelTrace/Pelican prompts, data, algorithms, renderer assets, sampling/concurrency, retries, application limits, partial outcomes, token/duration and result behavior. Apply only approved C02—C22 minimal host changes; no extra transport caps or model support guessing.

R2. Preserve F01—F06: original standalone page at /codex-eval through the native /admin/codex-eval sidebar/icon/iframe view, native same-origin auth_token/direct APIs, group/account/raw-evaluation-model/status operations including native schedulable POST, original type/subscription/Pro filters and model suggestions with direct model input, toolbar/layout/theme adaptation, Markdown/math/Pelican rendering, polling/progress/history/detail/clear and original stop controls. Replace CPA glue directly, without legacy API translation, activation or extra login.

R3. Add separate native-admin POST /api/v1/admin/accounts/:id/codex-eval beside /test. Complete streaming Responses input and output, single call, path account only, actual HTTP/SSE error and observed usage. S01 excludes nonstream/media/compact expansion. Original /test remains native.

R4. HTTP endpoint and four background evaluations share a sending method based on original AccountTestService. Send the selected model literally, bypassing evaluation model configuration. Preserve native credentials/URL/proxy/TLS/header/plugin transport and existing Responses/Chat Completions conversion. No internal HTTP, fake Gin or retained admin JWT. Native snapshots/401/429 handling remain; new endpoint omits broad successful-test recovery, per S04.

R5. Administrator test semantics only: no internal user/Key/owner/expiry/selection header, user billing/group pricing, gateway auth/scheduler/sticky/fallback/partial billing/parser changes. Group is page filter/result origin, not permission/routing. Exact account/protocol support and original per-account busy remain. Failures never send to another account. Supplier calls still consume supplier quota.

R6. Original simple memory-only tasks/progress/cancel and four-map history. Keep latest20 per account/type, per-result save/errors, unreadable/null-file no-overwrite protection, atomic temporary rename and clear failure rollback. No ActiveIntent/schema/start log/restart recovery/auto-resume/synthetic unstarted history/Flush/migration. No real history/file/DB deletion is authorized by this plan or implementation of a clear control.

R7. Ordinary patch and completed native evaluation module using the existing original Sub2API source, developed in a standard system temporary workspace. U20 removes the original CPA copy after migration. Reuse upstream tools/locks and installed frontend dependencies. Rebuild changed native frontend inputs before embedding them. Build the patched native binary. No AST/multiversion/type-export/universal-profile/old-overlay pipeline, and no build/source fetch on user's host.

R8. Fixed outputs dist/build/sub2api-codex-eval and dist/build/sub2api-codex-eval.tar.gz. Fresh standard-library install.py selects explicit Docker container, backs up/replaces standard binary, restarts only that container, reports progress/failures and restores backup on failure. No required nc/buildx/helper image/docker --format or maintenance framework. Package the binary, installer and necessary licenses. No standalone cleanup script or cleanup packaging/testing requirement.

R9. Keep one reusable Docker test container named `sub2api-codex-eval-test`. Run its test dependencies inside that container. Reuse it for later checks. Remove earlier containers created by this task. Preserve unrelated containers, volumes, images, caches, and real application data.

## Acceptance

A1—A6. Completed pure baseline is existing evidence, not a repeated acceptance gate. Earlier test/check matrices are superseded by this acceptance section; no race requirement.
A7. Current delivery implements approved functions through the completed native integration and original Sub2API tree. The initial fresh port used original plugin source. U20 removes obsolete CPA files after migration. No mandatory source identity/parity verification applies.
A8. Native-admin only, complete selected-account call, no fallback; original /test and normal gateway preserved.
A9. Four original evaluations and standalone page preserve approved behavior. Use relevant existing tests for changed code and only necessary local checks of native sending, page/API integration and cancellation; reuse evidence across these checks.
A10. History retains20, load/save errors are visible, corrupt-file protection and clear rollback hold; restart issues no calls or synthetic record.
A11. Deliver the static Linux/amd64 binary with the native frontend at the fixed package paths and include necessary licenses, excluding dev/test content and standalone cleanup.
A12. Check the changed binary replacement path, including failure restoration, using local disposable resources. Build/test only what the actual changes require; no repeated pristine baseline, exhaustive cross-platform matrix, resource parity, coverage threshold or extra full-suite requirement. Report actual build/function/install results and existing limitations. No paid calls or CPA deployment; local checks do not claim N150 execution.

## Decisions and current boundary

F01—F06/S01—S05/C02—C22 are approved function contracts, not old-code repair recipes. C01/C05 are covered by S02/S05. User R02 permits redundancy removal without extra questions; R03 requires reconstruction from originals. Authoritative ledger: research/patch-scope-approval-ledger.json. Exact decisions/anchors survive in frontend-first-approval.md, admin-account-codex-eval-approval.md, core-admin-channel-approval.md, core-second/third/fourth-approval.md and evidence under research/.

Previous PRD/design/checklist/context/task records remain in `research/before-fresh-origin-*`. The user authorized moving the task history directory to Trash under U19. Old inventories and source manifests are evidence. Later decisions control when these records conflict. R04 removes new source pins, standalone cleanup, and added acceptance gates. The original execution approval remains valid for the approved core port and corrections. The user approved remaining inspection/settings implementation with `批准实施` after the complete integrated plan. No production request, paid provider call, or external menu write is authorized.

## Current user corrections

The user reported four gaps after trying the delivered package. These supersede conflicting earlier details.

U1. Deliver the native administrator sidebar entry and iframe embedding for /codex-eval; manual optional configuration is insufficient. Reuse the existing native menu/iframe mechanism without Vue changes or an additional menu framework.
U2. Align the evaluation page with native account disabled semantics. The native account switch controls schedulable; false must display as disabled with the switch off, rather than generic unavailable. Preserve account status and runtime restrictions as distinct facts, using the native matching update API.
U3. Evaluation requests use the administrator's chosen model directly, bypassing Sub2API model configuration/mapping restrictions in the evaluation path only. Keep credentials, URL/proxy/TLS/plugin transport and native protocol conversion. Do not invent another configuration system.
U4. Reproduce and repair the installer waiting path reported as connection reset by peer. Follow the selected container's current restarted address/state and actual readiness. Retain failure restoration; increasing timeout alone is insufficient.

Run focused checks for these corrections on the final rebuilt artifact, including real disposable Docker/native menu/iframe API and disabled-state/model bypass/readiness behavior. Reuse unchanged baseline and earlier unrelated checks. No production settings writes or paid model calls are needed.

## Latest browser and model correction

U5 supersedes U1's zero-Vue/menu-response implementation restriction. The user now authorizes a direct native sidebar entry with a native icon and native iframe view for /codex-eval. Preserve the actual browser origin including nonstandard port; patch the native URL/origin handling responsible for the reported parent frame-src 'self' violation. Do not weaken global CSP to hide a wrong-origin URL. Rebuild the changed native frontend rather than embedding stale output. Remove superseded generated-menu glue where no longer needed.

U6. The user reports model bypass still does not work. Trace the actual page-start/model-discovery/capability/sender paths and locate configuration-dependent calls that still reject or remap the administrator-selected model. Repair the evaluation path rather than treating the previous sender-only fixture as sufficient. Retain ordinary native test/gateway behavior.

Validate the final artifact in an actual browser at a nonstandard port, including clicking the native sidebar/icon, successful iframe content and no frame-src rejection. Verify the page-start request and upstream model with conflicting configuration using the same disposable native runtime. Existing Node-VM evidence remains partial contract evidence, not browser acceptance. No production or paid calls.

## Latest installer feedback and disabled model discovery

U7. The user confirmed installation succeeds after the first connection-reset line; the remaining issue is misleading waiting output. Show normal application startup/readiness progress and a clear ready result. Keep raw transport details for actual final failure, preserving the existing bounded readiness and rollback behavior. Do not redesign working waiting logic based on the first transient reset.

U8. Installer-generated backup/temporary files must default under `/tmp`, rather than the unpacked installation directory. Preserve backup-based rollback and report the actual backup location. Do not download tools or add a cleanup subsystem. Do not delete existing user backups as part of this repair.

U9. Accounts whose native calling/scheduling switch is off must not produce model-fetch error notices. Skip disabled accounts during evaluation model discovery and maintain visible errors for active accounts' actual failures; preserve account rows, enable controls and direct model input.

Verify only these changed paths and rebuild the fixed output package. Reuse the passed browser-origin/model-forwarding and unchanged installer rollback evidence where applicable.

## Latest UI correction U10

U11 below replaces the account-picker and status-copy requirements in this section.

Remove assistant-added small-print explanations from the evaluation page and settings tab.
Keep necessary field labels, units, actual errors, progress, result data, and concise save feedback.
Use dropdown controls for models and account selection instead of plain inputs or account-pattern textareas.
Account include/exclude selection uses dropdowns that support multiple accounts.
Model dropdowns show actual native account model suggestions and retain the existing literal/custom-model capability.
Preserve saved shell wildcard arrays and native settings storage without changing the backend API.
Match native Sub2API controls and use the existing single reusable Docker test container.
Rebuild the fixed v0.2.17 binary and archive after the page change.

## Latest settings simplification U11

Remove the inspection include-account and exclude-account controls and their unused filtering configuration and logic.
Keep the Free switch and native disabled/unavailable/busy selection rules.
Keep model dropdowns and independent raw model discovery.
Hide the normal and degraded priority fields while automatic priority adjustment is off.
Default normal priority to 10 and require a minimum value of 1.
Remove idle `等待巡检` and `自动巡检已关闭` status copy.
Keep actual running progress, errors, results, and a scheduled next-run timestamp when available.
Remove notification configured-status text while preserving blank-secret save behavior.
Rebuild the fixed binary/archive and verify the changed controls in the retained container.

## Latest settings collapse U12

Collapse automatic inspection configuration while its switch is off, as with the priority fields.
Remove visible `启用` and `自动调整` switch captions. Keep accessible switch names.
Keep manual start and stop available when automatic inspection is off.
Place manual actions and actual progress/errors outside the collapsed configuration.
Retain all U11 removals, model controls, priority values, secret handling, and save behavior.

## Historical account-list correction U13, replaced in part by U15

Automatic inspection must reuse the existing evaluation account-list implementation and presentation.
Replace the separate simplified inspection result rows with the shared account cards, native identity/state, and Candy records.
Keep actual inspection outcomes, priority transitions, test counts, and errors available with their matching accounts.
Avoid another independent account-list renderer or inconsistent account-state conversion.
Preserve U11/U12 settings behavior and the existing native inspection scheduling and selection rules.
Do not bind autonomous server inspection to an unpersisted browser filter.
The user also reports accounts from other groups in manual test lists.
All manual evaluation tabs must show only accounts belonging to the selected group.
Pagination, independent inspection catalogs, and late responses must not mix manual group inventories.
Submitted manual targets must match the selected group's displayed eligible accounts.

## Settings correction U14: persistence and cancellation

The user reports settings that do not save or reload correctly.
Check the existing inspection settings and evaluation preference paths to identify the affected contract.
Saved values must return through the established storage and page controls without silent replacement.
Preserve native secret redaction and blank-secret handling.

The user identifies history as the location of the `context canceled` error.
Trace actual task and request lifetimes before changing cancellation behavior.
Page navigation, polling, and a completed start request must not cancel a background evaluation.
Explicit task stop must still cancel supplier work and record the actual outcome.
Preserve historical records. Do not rewrite existing failures or turn cancellation into success.
Keep U13 account-list behavior and all previous approved contracts.

## Historical inspection implementation before U15

The user approved native Candy inspection and the settings design with `批准实施` after the complete integrated plan. Implementation, browser verification, and independent review are complete. Earlier requirements, approvals, and evidence remain part of this migration.

The inspection runs inside the existing Sub2API process and Docker container. It does not use cron, a second container, Python, Docker socket access, or a long-lived administrator token. It reuses the existing native Candy evaluation sender, task state, history, account filtering, and account update service.

The original `examples/inspection.py` runs only Candy. It performs one first test, then adds at most two confirmation tests only when the first result indicates a possible priority transition. It waits 240 seconds before recovery confirmation, requires three valid matching outcomes, skips disabled/unavailable/busy/Free accounts according to configuration, updates only the two configured priority values, and sends ntfy/Bark notifications only after successful updates. A run normally polls every 15 seconds and has a 30-minute completion bound; these are execution limits, not the automatic schedule.

Native scheduling adds a persisted settings tab to `/codex-eval`, styled after the original Sub2API administrator settings tabs and cards. The settings include automatic inspection, raw model, effort, maximum tests per account, minimum and maximum wait minutes, skip Free, automatic priority changes, normal/degraded priority values, and ntfy/Bark fields. U11 removes include/exclude account patterns. Maximum tests means the original `1 + up to N-1` behavior, with default 3. Minimum and maximum wait default to 15 and 35 minutes; each completed schedule draws a new random delay. Equal bounds produce a fixed delay.

Save inspection settings through the native `SettingService` and `SettingRepository` in the existing PostgreSQL `settings` table. Use inspection-specific setting keys. Do not create `inspection.json`, another table, or a separate storage service. Follow native secret-field behavior. Read responses show configured status, and empty secret updates preserve the saved value. Keep the next run time in memory. A browser refresh reads this time without changing it. A process restart draws a fresh delay from the saved settings.

The settings tab has no explanatory paragraphs. It uses visible field labels, units, native controls, inline validation, and save feedback. The automatic inspection card shows the current state, next run, current run progress, and immediate-run/stop-current-run actions. Existing evaluation tabs retain their existing controls; their refresh and clear actions remain on evaluation tabs.

Inspection results remain tied to the existing Candy history and run IDs. A disabled or unavailable account is skipped without a model-fetch error; active account errors remain visible. A cancelled, failed, incomplete, or inconsistent run never changes priority. Sub2API priority order is ascending, so lower numeric values are more preferred; do not copy CPA's larger-is-preferred constants. Only accounts whose current priority equals one of the two configured values can be changed. Automatic priority changes default off; when off, the page shows recommendations without writing accounts.

Container lifecycle is part of the design: the scheduler starts with the process, stops cleanly with it, and after a restart loads settings and schedules a fresh random wait when enabled. An interrupted run is not resumed and does not create synthetic history. Notification failure is reported separately and does not rerun a test or undo a successful priority update.

Inspection acceptance covers persisted settings, random bounded scheduling, no overlapping runs, original `1 + up to N-1` confirmation behavior, disabled-account filtering, ascending priority updates, notification outcomes, stop/restart behavior, and the settings-tab layout. Reuse existing native evaluation checks and add only checks required by these new paths. Implementation, local checks, and delivery builds are authorized. Production deployment, paid calls, and production setting writes remain outside this work.

Response model comparison is outside the plugin contract. Its absence is not a defect or a required test. Preserve the native gateway response-model audit behavior. The researcher confirmed that its relevant source files match the original source.
