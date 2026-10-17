# Complete Plugin Migration Execution Plan

Status: in_progress. Changes through U18 are implemented and checked. Work-commit approval and the earlier cancellation report remain open.

## Current U22 implementation

The revised layout uses production patches and places every test input under `tests/`.
Before conversion, assign each production file and upstream edit to one functional patch.
Separate runtime data from test-only data. Keep existing tests and assertions without adding a new suite.
Apply the complete series to an original upstream tree and attach Go tests only for the check operation.
Update UI checks to read the patched temporary tree instead of a production source copy.
Move the existing `install.py` to `scripts/install.py` and update tests, packaging inputs, and repository references.
Preserve the package entry name and installation command. Do not rename it to `installer.py`.
Rename the build entry to `scripts/build.py` and merge the existing packaging code into it.
Remove `build-sub2api.sh`, `test-sub2api.sh`, and `package-sub2api.py` after updating their callers.
Preserve existing checks through `build.py --check`. Do not add wrappers or verification requirements.
Merge useful rollback commands and the history path into `README.md`, then remove the duplicate operations guide.
Remove `docs/`. Keep only the build and installer command entries under `scripts/`.
Verify patch composition and current test entry points before removing the `src/` copy.
The user approved implementation. The implementation owner converts the layout, then the independent checker reviews the result.
Main synchronizes current specifications and the staged work candidate after the check.
Reuse U21 and earlier runtime evidence when effective source contents remain unchanged.
U22 implementation and independent review passed. Main synchronized current specifications and the exact 100-path staged candidate.
The final layout has seven production patches, separate tests, two script entries, and README instructions without `docs/`.
Existing focused tests, patch composition, byte preservation, command orchestration, syntax, links, and archive checks passed.
No product rebuild or Docker operation ran. The archive uses the existing binary.
Evidence: `research/U22-patch-series-implementation.md` and `research/U22-patch-series-check.md`.

## Previous U21 repository layout

The implementation owner moves current additions into `src/backend/`, `src/frontend/`, and `patches/sub2api.patch`.
Update their build, test, fixture, documentation, and specification consumers together.
Trash generated caches, Finder metadata, old browser output, and the duplicate root notice.
Preserve current package inputs and fixed artifacts, private data, unrelated work, Git, and agent integration directories.
Run the existing UI harness after its path change and check affected Python/shell entry points.
Verify new source/patch composition without another full build or Docker run.
The independent checker reviews final roles, links, callers, and checks.
Main updates task records and the candidate commit plan after the final paths settle.
U21 implementation, independent review, and current spec synchronization are complete.
Main refreshed staging to the exact 138-path work candidate, including current source paths and CPA removals.
The stale staged source paths and duplicate notice are gone. No commit occurred.

## Previous U20 CPA cleanup

The existing implementation owner traces current dependencies and removes the obsolete CPA graph through macOS Trash.
Include exclusive ABI/source/module files, tools, examples, documentation, placeholder configuration, and unused screenshots.
Keep the native implementation, current tests and commands, required licenses, private data, and unrelated changes.
Update current AGENTS and specs to describe only Sub2API ownership and commands.
If package metadata changes, regenerate only the archive with the unchanged executable.
Check current references, affected script syntax, archive entries, and documentation links.
The independent checker reviews removals and current entry points after implementation.
Main synchronizes task records and the complete commit candidate. Do not stage or commit in this cleanup.
U20 implementation and independent review passed.
The source-only cleanup and regenerated archive need no binary rebuild or Docker operation.
The refreshed work candidate contains 83 added/modified paths and 56 old CPA removals before rename detection.
The index remains at its earlier 80-file proposal. No new staging or commit occurred.

## Previous U19 repository organization

Research the exact active and retired file inventories and current references before moving files.
The existing implementation owner organizes the confirmed paths and updates affected current documentation.
Keep current Sub2API command entry points, fixed artifacts, and unrelated changes intact. U20 supersedes original CPA retention.
Move the whole task history directory to macOS Trash and remove current archive references.
Do not create a replacement archive. Retain current task artifacts and research.
The independent checker reviews the final tree, current references, unchanged entry points, and Git scope.
Run affected syntax and path checks only. Reuse existing product build and browser evidence when inputs remain unchanged.
Main refreshes the business versus Trellis-managed commit plan after organization.
The first U19 archive layout was rejected. Apply the user's whole-directory Trash correction and check its affected references.
The whole task history directory moved to macOS Trash. Independent review confirmed the absent path and corrected current references.
Keep the current directory and quality specs accurate and preserve the hybrid commit workflow.
The refreshed `research/work-commit-review.md` lists 81 candidate work files, including `.gitignore` and current documentation changes.
The existing 80-file index remains unchanged. Final work-commit approval remains pending.

## Previous U18 correction

The existing implement owner moves next-run status into the Candy account header and removes the separate lower runtime card.
Update only affected UI assertions and remove dead styling or adapters.
Preserve real errors, active feedback, settings-bottom controls, scheduling, and U17 copy removal.
The independent checker reviews source before the fixed build/package.
The main session verifies desktop and narrow layouts in the retained container.
Reuse valid backend and timer evidence.
U18 implementation, independent source review, fixed build/package, and desktop/narrow browser checks passed.
Evidence: `research/U18-header-status-implementation.md`, `research/U18-header-status-check.md`, and `research/U18-header-status-browser.md`.

## Previous U17 correction

The existing implement owner removes saving/unsaved captions, account inspection summaries, and the completed-phase caption.
Keep real errors, active progress, next schedule, autosave, and U16 button locations.
Update existing focused UI assertions. Do not add backend work or broad checks.
The independent checker reviews source before the fixed build and package.
The main session checks the final page in the retained container and records the evidence.
Reuse valid U15 backend and timer evidence.
U17 implementation, independent review, build/package, and retained-container browser checks passed.
Evidence: `research/U17-status-copy-implementation.md`, `research/U17-status-copy-check.md`, and `research/U17-status-copy-browser.md`.

## Previous U16 correction

The existing implement owner restores the two inspection buttons to settings and removes successful-save text.
Update affected UI assertions. Reuse valid U15 backend and timer evidence.
The independent checker reviews this small correction before the fixed build/package.
The main session verifies the final controls and silent successful saves in the retained container.

U16 implementation, independent review, fixed build/package, and actual browser checks passed.
Evidence: `research/U16-controls-save-copy-implementation.md`, `research/U16-controls-save-copy-check.md`, and `research/U16-controls-save-copy-browser.md`.

## Current U15 execution

The existing Trellis implement owner owns this coherent frontend/backend change.
Read the U15 PRD/design before historical sections and manifest references.
Trace existing Candy batch selection, native settings, coordinator, history, and page controls first.
Record the minimal shared configuration and count semantics before editing.
Remove the explicit save action, independent settings account list, duplicate model discovery, and unused adapters.
Connect valid autosave edits and current Candy group/filter selection to native persisted settings.
Keep independent inspection model/effort/count controls and reuse the current group's model suggestions.
Make scheduled selection match the Candy page's current group and filter semantics.
Preserve independent inspection counts/confirmation/recovery flow and every manual fixed-count entry point.
Keep priority/notification rules and real cancellation. Do not combine manual and inspection orchestration.
Add focused regressions for autosave, settings read/restart, selected-group/filter equality, and shared Candy results.
Report source-ready changes and check evidence before rebuilding.
The independent Trellis checker reviews source before the implementation owner builds the fixed outputs.
The main session loads the final binary into the same retained container and verifies the actual page and timer path.
Update current specs and evidence. Preserve previous reports for unchanged behavior.
Do not create containers, run paid calls, add broad suites, calculate hashes, stage, commit, or deploy from a child.

U15 implementation, independent source review, fixed build/package, and retained-container browser checks are complete.
The browser confirmed automatic persistence, group isolation, Candy history, and scheduled execution without an open application page.
Inspection counts and confirmation flow remain independent from manual fixed-count tests.
Evidence: `research/U15-shared-Candy-schedule-implementation.md`, `research/U15-shared-Candy-schedule-check.md`, and `research/U15-shared-Candy-browser.md`.
The earlier history cancellation report remains open. No work commit has been approved or created.

## Whole-task execution state

Keep the previous execution approval and completed checklists below. They cover the reviewed core port and later corrections. They do not complete the whole migration or authorize the remaining inspection implementation.

The user approved the complete integrated plan with `批准实施`. Implement the remaining inspection, native database settings, settings tab, focused checks, and rebuilt delivery. See `research/migration-history-and-decisions.md` for decisions, replacements, and source records.

| Workstream | State |
| --- | --- |
| Pristine Sub2API baseline | Recorded. Reuse evidence, failures, skips, and limits. |
| Four evaluations, native sender, original page, tasks/history | Implemented and checked within the recorded scope. |
| U1-U9 corrections | Implemented and checked. Later reports replace affected earlier results. |
| Native Candy inspection and database settings | Implemented and checked, including actual browser controls. |
| Integrated checks and final artifacts after remaining work | v0.2.17 package rebuilt after review and checked in the retained container. |
| Work commit and task closure | Pending. |

## Planning

- [x] Preserve completed pristine baseline with failures/skips.
- [x] Record approved F01—F06/S01—S05/C02—C22.
- [x] Record user authorization for redundant removal without batch questions.
- [x] Reuse the already compiled original Sub2API and this Git plugin; no fixed-source export or repeated identity verification.
- [x] Finish original-to-native file/constructor/route map in research/fresh-origin-port-map.md.
- [x] Converge requirements/design/current manifests; page policy follows already approved F01/F04, no new functional choice.
- [x] Remove standalone cleanup delivery and checks; restrict verification to actual changes without new acceptance requirements.
- [x] Present updated plan; user approved execution and task.py start completed.

Old planning/context/task saved as research/before-fresh-origin-*. Preserve actual workspace/product/artifacts/resources during planning.

## 1. Fresh inputs after approval: completed port stage

Owner: Trellis implement agent. Use original v0.2.17 source and this repository's internal/plugin in a standard system temporary workspace. The user changed the target during implementation. Reuse dependency caches and retain applicable v0.2.16 evidence. Select original algorithms/runners/data/web/history and replace CPA glue. No fixed commit export, source manifest, or repeated source verification. Do not use old internal/eval/adapter/compatibility/patches/build/installer as source. No reset/clean or real-data/resource deletion.

## 2. Native evaluation

Owner: native/evaluation implementer; fixed native patch and fresh backend module/resources. Add dedicated admin endpoint/shared native sender using original dependencies. Preserve /test and approved request/protocol/state behavior. Connect original four runners directly, simple memory tasks/progress/context and history-only storage. Review diff against pristine; exclude gateway/Key/billing/scheduler/sticky/parser/listener/shutdown/maintenance edits.

## 3. Original standalone page

Owner: page implementer if separately delegated; fresh web assets and necessary serving/SPA bypass, coordinate API with previous owner. Start from original page, with the approved native sidebar/icon/iframe view and URL helper patch. Apply approved native session/direct APIs/group/account/raw-model/status/filter/theme. Keep layouts/renderers/Pelican/history/polling/stop. No activation/login/API translation/ignored catalogs/HTML-CSS rewriting/cache framework.

## 4. Fresh delivery

Owner: delivery implementer once native/page inputs ready; simple build/package entry, install.py/docs. Ordinary patch, original resources, existing upstream locks/tools/caches and built frontend. Build patched static Linux/amd64 at dist/build/sub2api-codex-eval; rebuild frontend only if its build inputs change. Write small explicit-selection/backup/replacement/restart/readiness/progress/rollback installer; no on-host source/build, buildx/nc/helper-image/maintenance or --format prerequisite. Package binary/installer/necessary licenses into dist/build/sub2api-codex-eval.tar.gz; no dev/test payload. Remove standalone cleanup from active build/package references; delete the obsolete script only during approved implementation, without executing cleanup or deleting caches/business/history.

## 5. Check

Owner: Trellis check agent; main coordinates/specs/final delivery. Native context injection preferred, child-side full loading fallback; assign file ownership and preserve others' changes.

Review only the new changes against approved behavior. Format changed Go files and run the relevant existing package tests/checks once. Add a focused local sender/page integration check only for behavior not covered by those tests; reuse the same fake-upstream evidence for streaming/account targeting/cancellation instead of adding parallel suites. Check the built binary/package and changed installer replacement/failure-restoration path. No standalone cleanup checks, repeated pristine tests/builds, source/resource checksum or parity gates, coverage thresholds, added race/full-suite/platform matrix. Historical spec/research checklists do not add requirements to this task. Report actual failures/skips/build/install/function outcomes; local execution does not prove N150 execution.

Dependencies: final approval → native/module from existing original inputs → original page port → patched binary/package → checks for the changes. One task's ordered workstreams, without repeated redundant-removal approvals.

## Execution results

- [x] Reconstruct native module from existing original inputs.
- [x] Connect native administrator sender, background tasks and standalone page.
- [x] Build fixed binary/package and simple replacement installer.
- [x] Complete independent Trellis check and repair findings.
- [x] Synchronize current source-backed specs.
- [ ] Approve and create the reviewed work commit under the configured signing policy.
- [ ] Run Trellis finish-work archival and session recording after the work commit.

Implementation evidence: research/fresh-native-implementation.md.
Independent check and final artifact evidence: research/fresh-native-check.md.
At this initial check, local checks passed. Visual browser and actual N150 execution were unperformed. Later U5-U6 reports record actual Chrome checks. N150 execution remains unperformed.

## User-reported corrections U1-U4

- [x] Dispatch Trellis implement for native sidebar/iframe, disabled state, model-configuration bypass and installer readiness.
- [x] Rebuild fixed artifact and verify the corrected paths in disposable Docker.
- [x] Dispatch independent Trellis check and update specs.
- [x] Regenerate the commit review.
- [ ] Receive final signed-commit approval.

Correction implementation and artifact evidence: research/fresh-native-corrections-implementation.md. Independent Trellis check passed without further product edits: research/fresh-native-corrections-check.md. At this stage, visual browser and N150 execution were unperformed. Later U5-U6 reports record actual Chrome checks.

Prior work-commit review and approval question are superseded; no commit has been created.

## Latest U5-U6 repair

- [x] Trellis implement: direct native sidebar/icon and iframe with nonstandard-port origin patch; trace and repair the remaining evaluation model restriction.
- [x] Rebuild changed frontend, native binary and fixed package; verify the actual browser/sidebar/iframe/CSP and page-start model path in disposable Docker.
- [x] Independent Trellis check and synchronize specs.
- [x] Replace the superseded commit review.
- [ ] Obtain final signed-commit approval.

Previous U1-U4 passing evidence is retained; the user's browser/model failure reopens those affected acceptance paths. No final commit approval was given.

Final evidence: research/native-browser-model-repair.md, research/native-browser-verification.md and research/native-browser-model-check.md. Actual Chrome, related builds/checks and final stopped-container rollback passed. N150/iStoreOS and production HTTPS/proxy execution remain unperformed.

## Latest U7-U9 repair

- [x] Trellis implement: clear normal startup/readiness output, default backup/temp path under /tmp, disabled-account model discovery/error filtering.
- [x] Focused changed-path checks and fixed binary/archive update.
- [x] Independent Trellis check and spec/task synchronization.
- [x] Update the commit review.
- [ ] Receive final signed-commit approval.

The user confirmed installation eventually succeeds. Treat the reset line as a progress-copy issue; preserve successful waiting behavior. The prior commit review/approval request is superseded again; no commit has been approved or created.

Final U7-U9 evidence: research/installer-disabled-model-repair.md and research/installer-disabled-model-check.md. Focused startup/reset, mixed disabled/enabled model discovery/toggle, default /tmp backup and restoration checks passed. Earlier 3-second synthetic fixture restoration timed out; the recorded 8-second fixture rerun passed without changing product bounds. Current fixed binary/archive rebuilt; unchanged prior browser/Go evidence reused.

## Remaining migration work: native Candy inspection and settings

This section is part of the same complete migration. The user approved its implementation after the integrated plan review.

### Agent order and scope

One `trellis-implement` agent owns the remaining cross-layer change from settings storage through the page and delivery. It first loads the spec/research manifest and task artifacts. It records the concrete API fields before it edits their consumers. Backend, page, and delivery preparation are ordered parts of this implementation, not separate simultaneous implement assignments.

After implementation and focused checks, the main session dispatches an independent `trellis-check` agent. The check covers the complete affected task scope and reads the relevant spec Quality Check sections. It fixes mechanical issues directly and reports design findings to the main session.

The main session then applies `trellis-update-spec`. New API/storage/lifecycle contracts require the seven code-spec sections. Work commit review follows Phase 3.4. Task archival and session recording follow the work commit.

Context injection uses native `SubagentStart` when available. Without its complete marker, the agent loads its own manifest and all referenced files. Every dispatch starts with the explicit active-task path. Implement/check agents do not dispatch another implement/check agent.

1. Add a backend inspection coordinator beside the existing native evaluation runtime. Reuse the Candy sender, run IDs, cancellation, state, and native account update service. Keep the next run time in memory.
2. Add typed inspection settings methods to the native `SettingService`. Read and write inspection keys through `SettingRepository` in the existing PostgreSQL `settings` table. Use native administrator settings routes and secret-field behavior. Add the fifth `/codex-eval` settings tab. Match native administrator settings tabs, cards, and controls.
3. Implement the original `1 + up to N-1` Candy confirmation flow, 240-second recovery wait, include/exclude/Free filtering, disabled-account skip, ascending priority update boundary, and ntfy/Bark notification result handling.
4. Implement the random inclusive delay between configured minimum and maximum minutes, no overlap, immediate-run and stop-current-run actions, and clean process lifecycle. Reload settings after container restart without resuming an interrupted run.
5. Check the changed settings, scheduling, confirmation, priority, notification, cancellation, restart, and page paths. Reuse existing evaluation checks and disposable local Docker fixtures where applicable. Do not add race, coverage, full-suite, cross-platform, paid-provider, production deployment, or cleanup requirements.

Execution authorization: `批准实施` after the complete integrated planning summary. Use Trellis implement and independent check agents. Local Docker/browser checks and delivery builds are authorized. Production deployment and paid calls are outside this work.

### Single reusable Docker container

The user now authorizes removal of earlier task containers and requires one retained test container.
The existing implement agent owns fixture consolidation and the fixed reusable entry point.
The main session owns the scoped residual-container removal.
Use `sub2api-codex-eval-test`, including its PostgreSQL and Redis test dependencies.
Keep the environment available for the browser and independent check consumers, then retain it for later reuse.
Do not create another random fixture container or prune unrelated Docker resources.

### Final integrated evidence

- [x] Implement native settings, inspection, page controls, and lifecycle.
- [x] Build actual v0.2.17 and generate the fixed binary/archive.
- [x] Complete independent review and repair failed-initial-settings-load admission.
- [x] Rebuild after the repair and complete the final retained-container browser smoke.
- [x] Remove 11 earlier task containers and retain one reusable container with its test dependencies.
- [x] Synchronize native inspection, persistence, delivery, and container-use specs.
- [ ] Obtain final approval for the exact reviewed signed work commit.
- [ ] Commit, then run Trellis finish-work.

Evidence: `research/native-inspection-implementation.md`, `research/native-v0217-integrated-check.md`, and `research/native-v0217-browser.md`.
The retained container is `sub2api-codex-eval-test` at `http://127.0.0.1:60198`.
No N150/iStoreOS or production HTTPS/proxy execution is claimed.

### U10 execution

The existing Trellis implement agent owns the complete UI correction and rebuilt delivery.
Remove assistant-added small-print explanations and use dropdowns for models and inspection account selection.
Preserve native settings fields, saved wildcard arrays, direct literal model entry, actual errors, and current task controls.
Inspect native Sub2API control patterns before implementing the standalone equivalents.
Use focused frontend checks and the current build entry. Reuse unchanged backend and installer evidence.
Then hand the candidate to an independent Trellis check agent.
The main session verifies the actual dropdown interactions in the retained Docker container and updates specs and commit review.

- [x] Remove added model and effort explanations and replace five model fields with dropdowns.
- [x] Replace account-pattern textareas with searchable multiple-account dropdowns.
- [x] Preserve literal models, saved patterns, disabled-account state, and newer form edits.
- [x] Repair the browser-reported menu clipping and rebuild the fixed binary/archive.
- [x] Complete independent Trellis review and actual desktop/narrow browser verification in the retained container.
- [x] Synchronize the final dropdown contracts and evidence.
- [ ] Obtain approval for the refreshed exact work commit.

Evidence: `research/U10-dropdown-implementation.md`, `research/U10-dropdown-check.md`, and `research/U10-dropdown-browser.md`.
The upstream parent-page 304/CSP nonce issue remains recorded. The evaluation iframe and U10 interactions passed.

Do not add a response-model comparison or a test requirement for its absence. Preserve the existing native gateway audit behavior.

### U11 execution

The existing Trellis implement agent owns the complete settings simplification and rebuilt delivery.
Follow the current U11 PRD/design instead of the superseded include/exclude controls and idle copy.
Remove the unused account filtering settings, page controls, shared picker code, and affected tests.
Keep raw model discovery and native secret handling.
Change normal priority default/minimum and handle the existing stored zero default on read.
Use focused checks for changed settings and visibility behavior, then rebuild the fixed outputs.
The independent Trellis check agent reviews the changed contracts and reuses unchanged task evidence.
The main session verifies the final settings page in the same container and synchronizes specs and commit review.

### U12 execution

The existing Trellis implement agent owns automatic inspection collapse, caption removal, and the rebuilt fixed outputs.
Keep manual actions and actual runtime feedback available outside collapsed configuration.
Use focused frontend checks and reuse the passed U11 backend/settings evidence.
The independent checker reviews the final source. The main session verifies the final built page in the same container.

- [x] Remove include/exclude settings and unused picker/filter code.
- [x] Set normal priority default 10 and minimum 1, with collapsed priority fields when off.
- [x] Remove idle and notification configured-status copy.
- [x] Repair invalid loaded-priority scheduling and preserve explicit-save recovery.
- [x] Collapse automatic inspection fields and remove visible switch captions.
- [x] Preserve manual actions, runtime feedback, and saved hidden values.
- [x] Rebuild fixed outputs and complete independent U11/U12 review.
- [x] Inspect final desktop/narrow render and hidden-value save/reload in the retained container.
- [x] Synchronize current settings and delivery documentation.
- [ ] Obtain final work-commit approval, then commit and finish the task.

Evidence: `research/U11-settings-implementation.md`, `research/U11-settings-check.md`, `research/U12-settings-collapse-implementation.md`, `research/U12-settings-collapse-check.md`, and `research/U11-U12-browser.md`.

### U13 execution

The existing Trellis implement agent owns shared account-list integration and rebuilt delivery.
Inspect the current account card, Candy row/history, list filtering/state, and inspection consumers before editing.
Reuse the existing owners with minimal parameterization and preserve existing test callers.
Replace the independent inspection result rows and retain actual inspection details on their matching accounts.
Run focused list/state/UI checks and rebuild fixed outputs. Reuse unaffected backend/settings/installer evidence.
The independent Trellis checker reviews the shared callers. The main session verifies the final list in the retained container.
The same owner also fixes the reported manual group leak through the actual native query and catalog consumers.
Use two synthetic groups to check each manual tab, shared/exclusive membership, group switching, and submitted targets.
Do not change unrelated backend account authorization or scheduler behavior to compensate for a page catalog defect.

U13 implementation, independent review, rebuilt package, and final browser checks are complete.
Evidence: `research/U13-account-list-implementation.md`, `research/U13-account-list-check.md`, and `research/U13-account-list-browser.md`.

### U14 execution

The existing Trellis implement owner investigates settings save/reload and `context canceled` in history as one coherent follow-up.
Identify and reproduce the responsible paths before editing.
Check existing preference and native settings contracts without adding persistence mechanisms.
Trace task context ownership and preserve explicit stop behavior and stored history.
Use focused regressions for demonstrated failures and reuse unaffected U13 evidence.
Report source findings before rebuilding so independent review can catch remaining defects before packaging.
Then rebuild the fixed outputs and verify the actual corrected flow in the retained container.
The main session owns container loading and browser checks. The independent checker owns source review.

U14 settings fixes passed independent review, final build/package, and actual same-container browser checks.
Evidence: `research/U14-settings-history-implementation.md`, `research/U14-settings-history-check.md`, and `research/U14-settings-browser.md`.
The reported history cancellation remains open because its triggering operation is unknown.
The earlier clarification asked whether that run received an explicit Stop. It did not receive an answer.
The passed tests establish the tested request-lifetime behavior. They do not prove a fix for the reported event.
Preserve the user's report and actual cancellation behavior.
Apply approved A1-A12. Do not add hardware, production, or historical-cause acceptance gates.
