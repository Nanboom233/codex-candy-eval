# Complete Plugin Migration Design

Status: in_progress. Changes through U18 passed independent review and final browser checks. U15 timer evidence remains applicable.

## U22 approved layout: production patches and separate tests

Store production changes as a numbered patch series with application order in `patches/series`.
Split patches by evaluation core, runtime data, vendor assets, administrator transport, inspection, page UI, and upstream integration.
Keep test source and test-only data out of every production patch.
Store Go tests under `tests/backend/` with paths matching their upstream packages.
Store golden test data beside those test packages. Runtime probe and attribution data remain production inputs.
Keep installer fixtures and UI checks under their existing `tests/installer/` and `tests/ui/` owners.
For checks, copy Go tests and their test data into corresponding packages in the temporary patched upstream tree.
This preserves package-private test access without changing production exports.
Production builds apply only the production series. Remove the source-overlay copy step and the separate `src/` tree.
Update existing build and test callers together. Do not add compatibility wrappers or broader verification requirements.
Move the repository installer from root `install.py` to `scripts/install.py` during conversion.
Update installer test imports, package inputs, and repository documentation together.
The package still places `install.py` beside the executable. Preserve the extracted-package command `python3 install.py`.
Use `scripts/build.py` as the single build, check, and package entry.
Merge the existing packaging code into that entry. Remove both shell wrappers and the separate packaging script.
Keep `--check` for existing checks without adding another command layer.
Merge manual rollback commands and the native history path into `README.md`.
Remove `docs/sub2api-operations.md` and the `docs/` directory after that merge.
The approved conversion is complete. Applied production and moved test contents match their previous bytes and modes.
Independent review passed. The archive uses the unchanged executable and current license notice.

## Previous U21 source and generated-file layout

Move the current backend and frontend additions from `native/` to `src/` without changing source contents.
Move the integration patch to `patches/sub2api.patch`.
Update the build copy paths and existing UI/fixture paths with the move.
The upstream tree still owns the Go module and frontend dependency installation.
Keep the same build, test, package, installer commands, and fixed output names.
Trash generated Python caches, Finder metadata, and old browser output. Keep Git and agent integration files.
Keep `licenses/THIRD_PARTY_NOTICES.txt` and remove the duplicate root Markdown notice.
Update that notice's source locations and regenerate the archive with the existing executable. Do not compile again.
Check source/patch composition and existing consumers at the new paths.
Historical reports keep their old source paths. The current PRD and this mapping control new work.
The final layout passed independent review. Required source and package inputs remain complete.

## Previous U20 CPA removal

The completed `native/` integration is the sole current local implementation.
The build overlays it onto the existing original Sub2API tree.
Trace active build, test, installer, package, and workflow reads before removing old files.
Move the obsolete CPA source graph and its exclusive tools, documentation, module files, and assets to macOS Trash.
Do not copy them into task records or another archive.
Remove current CPA development and deployment instructions. Preserve the hybrid commit workflow.
Preserve required MIT, upstream, algorithm, and renderer license attribution.
Use a neutral project-license entry name when packaging the root MIT license.
Keep private configuration and unrelated local files. Do not access private deployment values for this cleanup.
Regenerate the archive only if its metadata changes. Reuse the unchanged binary and product checks.
Earlier original-CPA-retention clauses describe the completed port stage and are superseded by U20.
The implemented boundary passed independent review. Current source uses only `native/` and the upstream module.
The regenerated archive contains the unchanged executable, current installer, and five license/notice entries.

## Previous U19 repository organization

Retain `native/` as the active Sub2API integration boundary.
Retain the four existing Sub2API command scripts and root `install.py`.
Retain the current native installer and UI tests in their existing directories.
Original CPA retention from this step is superseded by U20.
Move the complete task `history/` directory to macOS Trash under the user's corrected U19 instruction.
Do not create a replacement archive or move its files into another repository directory.
Remove current pointers to that directory. Keep current task artifacts and research at their existing paths.
Ignore generated Python caches and browser output instead of mixing them with source changes.
Do not add wrappers, change application behavior, touch private configuration, or rebuild unchanged product inputs.
The earlier archive boundary and its review are superseded by this correction.
Current specs retain the hybrid commit workflow.

## U18 correction

The Candy account header owns the next inspection timestamp.
Remove the separate runtime card below the account list and its obsolete visibility adapter.
Keep active feedback and actual errors in existing page structures without new columns or cards.
Keep settings-bottom inspection actions and silent autosave.

## U17 correction

Remove nonerror save captions and per-account inspection summaries from the embedded page.
Remove the completed-phase caption and stale completion progress. Keep active progress and scheduled time.
Use existing error locations for actual save, runtime, evaluation, and notification failures.
Remove dead summary rendering and styles. Preserve runtime data, history, settings, and button locations.

## U16 correction

Keep immediate-inspection and stop buttons on the settings tab, outside collapsed configuration fields.
Candy continues to display account progress, results, and history.
Successful automatic saves clear their transient status instead of displaying `已保存`.
Keep actual validation, save, and action failures visible at their relevant controls.
Do not change persistence, selection, scheduling, or inspection flow.

## U15 current design: shared account selection and automatic persistence

U15 replaces conflicting independent inspection inventory, result rendering, and explicit-save details below.
Keep historical sections as decision evidence, not active implementation requirements.

Use one shared account-selection configuration for selected group and current Candy filters.
Persist scheduling inputs through native `SettingService` and `SettingRepository`, so the backend can run without a browser.
Resolve current group membership at execution time rather than saving a fixed all-account or stale target snapshot.
Scheduled group/type/plan scope must match Candy `测试全部`, using current native membership and the same filter semantics.
Retain inspection's existing native availability, expiry, quota, and busy checks without changing manual admission.
The user's clarification keeps inspection model, effort, counts, and orchestration independent from manual tests.
Preserve all manual row, selected-account, and fixed-count batch behavior.
Keep inspection's first test, bounded conditional confirmations, recovery wait, priority changes, and notifications.
Do not introduce a new manual batch endpoint merely to combine these independent flows.

The Candy page remains the only account/results/history display.
Remove the inspection account-card container, separate inventory, context adapters used only there, and independent model discovery.
Keep necessary scheduling state and actual runtime errors without restoring idle filler text.
Keep independent inspection model, effort, and maximum-test controls.
Reuse the current group's raw model suggestions without a separate all-account discovery catalog.
Remove the independent Free-selection override. The shared Candy plan filter controls account selection.
Keep the other three manual evaluation tabs unchanged.

Automatically save valid settings edits and shared Candy account-selection inputs.
Use a short input debounce and one ordered save owner to retain newer edits and prevent overlapping writes.
Keep invalid in-progress input visible without writing an invalid configuration.
Use native secret redaction and blank-secret preservation, with actual save errors in the existing feedback location.
Read saved values on entry and restart without allowing old reads or save responses to replace newer edits.
Do not add a save button, another store, or a general compatibility layer.

The implement owner first traces existing batch/filter/settings boundaries and records the minimal shared contract.
Then it implements the coherent frontend/backend change and focused checks before independent review and delivery build.

The complete history is in `research/migration-history-and-decisions.md`. The source approval files retain their original text. Use the later decisions when they replace an earlier mechanism.

## Original module: completed port stage

The initial port used original v0.2.17 source and the original plugin from this Git repository. U20 removes the obsolete CPA copy after completion. Current builds use `native/`. Retain earlier v0.2.16 reports as comparison evidence. Work in a standard system temporary workspace. Restore absent upstream source from the official archive without cloning Git history or adding provenance gates. Keep all four native algorithms/runners/data/history/web/vendor files. Exclude CPA ABI, management envelope, credential inventory, and host-RPC loader. Do not copy retired adapters or test overlays.

Preserve original prompts/Codex context/wire/UA/UUID, pure algorithm/data bytes and page layout/renderers. Necessary changes are native account ID/name/group origin, direct context-aware native sending, approved incremental stream semantics and native session/API/theme/account/model glue. No general host compatibility framework.

## Administrator sending and routes

Add dedicated AccountHandler/AccountTestService methods beside /test, using existing dependencies and normal construction. HTTP endpoint and background runners share complete selected-account Responses sending; no ordinary-handler invocation, internal HTTP/JWT/private Key/fake Gin. Reuse native request/credential/proxy/TLS/plugin and Responses/CC conversion primitives while passing the selected model literally in evaluation. Keep full response/errors/usage, native snapshots/401/429, omit broad success recovery only here.

Register approved admin-only start/state/cancel/detail/clear operations and POST /accounts/:id/codex-eval. Exact account lookup/protocol support only; original per-account cross-type busy, duplicate folding and real unavailable summary. Group is display origin, not gateway admission.

## Tasks, history and page

Use one ordinary evaluation instance with original running maps/lock/sample slots and a simple value progress callback. Tasks are memory-only; cancel stops next requests/closes responses and busy clears after worker exit. No special process-control/drain framework.

Persist four original histories at the native data directory, latest20, per-result save and actual errors, unreadable-file write protection, atomic same-dir temporary rename and clear rollback. No task logs/schema/recovery/interrupted/Flush; no actual old-data migration or deletion.

Serve original independent /codex-eval shell with necessary backend route/SPA bypass and page-scoped original renderer/script/style capability. Native sidebar with beaker icon opens /admin/codex-eval; its native iframe view resolves relative /codex-eval against the complete browser origin including port. Original SecurityHeaders sets headers before c.Next, so the shell handler uses a page-only same-origin frame/approved inline policy, without modifying global middleware or rewriting renderer output. Native administrator session, direct native group/account APIs and evaluation-specific raw model suggestions and schedulable POST replace CPA glue; approved filters/theme/controls/results/polling/history remain.

## Fresh build and delivery

One ordinary patch plus fresh module/resources. Reuse the existing upstream tools, locks and dependency caches; rebuild the patched native frontend before embedding it. Build the patched static Linux/amd64 binary. No AST/type-export/profile/old-overlay pipeline. Reuse relevant existing tests; add checks only for changed behavior lacking coverage.

Fresh small standard-library installer uses explicit selection, standard /app/sub2api layout, original backup/replacement/restart/readiness/progress/raw errors/rollback. No old maintenance/activation/socket/deployment journal/owned-listener scanner/helper-image code. Fixed dist/build outputs include binary/installer/necessary licenses. Standalone cleanup is removed from delivery and checks; this planning update does not delete existing files or caches.

## Review and limits

Review new diff against originals. Native edits only cover required admin sending/routing/ordinary evaluation construction/page serving; normal gateway/main/listener/shutdown remain original. Determine necessary standard Wire edits from original constructors, not old runtime graph templates. Exact file mapping: research/fresh-origin-port-map.md.

Reuse completed baseline evidence and retain its failures/skips without repeating it. Limit validation to changed packages, native sender/page integration and binary replacement/failure restoration. Use existing tests and shared local fixtures where applicable. No source/resource parity gate, cleanup tests, new coverage/race/full-suite/platform matrix, paid provider calls or CPA deployment. Current task decisions supersede conflicting historical verification requirements. Prior execution approval covers the reviewed core port and U1-U9 through Trellis agents. The user approved remaining inspection/settings implementation with `批准实施` after the complete integrated plan.

## Current correction boundary

Current PRD U1-U4 supersede manual-only menu setup, status-only enable mapping and account model-mapping behavior in evaluation. Use the native administrator settings response to supply the existing sidebar/iframe entry, preserving unrelated menu items without writing production settings. Align the page with native schedulable updates and separate status display. The shared sender passes the selected model without configuration remapping. Installer readiness follows current selected-container state/address across restart. Preserve the small ordinary patch and existing native frontend output.

Latest U5-U6 supersede the menu-response/zero-Vue decision: directly patch the native sidebar/icon and iframe view, use the complete browser origin with nonstandard port, and rebuild changed frontend inputs. Trace remaining evaluation model-configuration dependencies, including discovery/probing before the shared sender. Actual browser integration and page-start/upstream evidence determine the result; previous Node-VM checks did not cover parent CSP enforcement.

U7-U9 narrow the next repair to startup progress wording (the user confirmed eventual successful readiness), default installer backup/temp storage under /tmp, and excluding native schedulable:false accounts from model discovery/error notices. Retain readiness bounds, rollback, active-account error visibility and direct model input. No new maintenance/cleanup or download path.

## Restored evaluation details and decision design

### Four evaluation contracts

| Evaluation | Required original behavior | Decision sources |
| --- | --- | --- |
| Candy | Keep the original prompt and standalone `21` check. Score the complete answer. Save up to 4,000 characters for display. Manual runs remain serial, from 1 to 10 tests per account. A wrong answer is a completed test with `ok=false`. | C04, C12, C18 |
| Fingerprint | Keep probes, baselines, normalization, and attribution. Keep 60/200/400 samples for quick/standard/strict modes. Use `low`, default concurrency 2, maximum 6, six global slots, and the original stop after eight consecutive request errors. | C04, C13 |
| ModelTrace | Keep the bank and attribution algorithm. Generate three distinct challenges with 292-332 digits. Keep challenge order, 1-3 workers with default 3, six global slots, and the 64,000-byte answer limit. Valid attribution with fewer than three valid samples remains `partial`. | C04, C14 |
| Pelican | Keep the original prompt, HTML extraction order, and 256 KiB limit. Keep full HTML in details. Do not rewrite the generated HTML. Keep the original frozen list preview and detail playback. | C15, F04b |

Keep the original probe retry rules. For status 0, 429, or at least 500, use at most three attempts with 2/4-second waits. Release a sampling slot before a retry wait. Each attempt keeps the same account. Probe timeout remains three minutes. Candy and Pelican answer timeout remains ten minutes. Do not apply probe retries to ordinary answers. Task cancellation stops a wait or request.

### Request, result, and stream contracts

C02 keeps one `model` and one `effort`, plus group/account/name/run/status fields. These fields describe the evaluation request. They do not prove the model used by the supplier. Remove duplicate requested/effective model fields.

C03 keeps model and effort separate. Keep `none`, `low`, `medium`, `high`, `xhigh`, and `max`. Empty effort means unspecified. Do not parse CPA model suffixes or guess model support. U3/U6 later require literal model names in evaluation sending and discovery. These corrections replace the earlier mapping wording in S02 and C02/C03/C10/C11. Native credential and protocol preparation still applies.

C06-C09 keep actual HTTP/SSE errors, terminal status, and observed usage. Read SSE incrementally using the original answer semantics. Accept native `response.done`. Keep completed output items when final output is empty. Do not use tool arguments or incomplete deltas as answers. An otherwise complete terminal event does not require an extra final blank line. Preserve cancellation and request-boundary panic handling.

Keep actual input/output/reasoning token values. Replace a usage snapshot within one request. Add observed usage across retry attempts. Do not estimate missing usage or create user billing. Remove added response byte/event caps and unused observation fields. Keep original algorithm-specific answer limits.

Response-model comparison is outside the plugin contract. The existing native gateway observer and mismatch audit remain unchanged. The independent source review is in `research/response-model-validation-followup.md`. No plugin comparison feature or extra mismatch test is required.

### Frontend and native management contracts

F01-F06 preserve the original page and controls. F02 uses same-origin `auth_token` for administrator requests. Missing token or HTTP 401 stops polling and asks for native login. HTTP 403 shows missing administrator permission. The page has no second login or refresh-token system.

F03a places the shared group selector first in each evaluation toolbar. F06c restores the last valid group or selects the first available group. An empty group list disables starts. Switching groups clears selections and ignores stale account/model responses. It does not cancel existing backend tasks or delete history.

F06a reads all pages of the current native group. Show actual name, type, priority, subscription, status, and runtime restrictions. Do not invent a zero priority or a subscription. F06b keeps type, subscription, enabled-only, and Pro ALL filters. F06e keeps native light/dark colors and original plan badge differences.

U2 replaces the earlier F03b status-toggle mechanism. The calling switch controls native `schedulable`. Status and runtime restrictions remain separate. Keep request locks and failure restoration. U6 adds evaluation-specific raw model discovery. U9 skips disabled accounts in that discovery and keeps active-account errors visible. Manual model input remains available.

U5 replaces the original zero-Vue/custom-menu proposal. The native beaker sidebar entry opens `/admin/codex-eval`. The native iframe resolves `/codex-eval` against the complete browser origin, including its port. Keep global CSP unchanged. The evaluation shell has its own required response policy.

F04 keeps Markdown/math renderers and the original Pelican `srcdoc` preview with `sandbox="allow-scripts"`. Do not grant same-origin access or rewrite generated content. F05 keeps start/state/cancel/detail/clear controls. Poll every 2.5 seconds while active and every 20 seconds while idle. State polling does not fetch models.

### Administrator endpoint and task API

| Method | Path | Purpose |
| --- | --- | --- |
| POST | `/api/v1/admin/accounts/:id/codex-eval` | One complete streaming Responses call to the path account. |
| GET | `/api/v1/admin/accounts/:id/codex-eval/models` | Raw evaluation model suggestions. |
| POST | `/api/v1/admin/codex-eval/start` | Start selected evaluation tasks. |
| GET | `/api/v1/admin/codex-eval/state` | Active tasks, progress, and result summaries. |
| POST | `/api/v1/admin/codex-eval/cancel` | Cancel the selected run ID. |
| GET | `/api/v1/admin/codex-eval/records/:type/:account_id/:record_id` | Fingerprint, ModelTrace, or Pelican details. |
| POST | `/api/v1/admin/codex-eval/clear-history` | Clear one history type or all types. |

S01-S05 keep native administrator authentication. HTTP and backend tasks share `AccountTestService.SendCodexEvaluation`. Account selection comes from the path or admitted task. Keep credentials, shadows, URL validation, header overrides, proxy, TLS, and native plugin transport. Reuse native Responses/Chat conversion. Do not forward the administrator token to a supplier.

The direct endpoint returns the complete stream, observed usage, and actual errors with native sensitive-data redaction. Preserve native snapshots and 401/429 handling. Do not call broad successful-test recovery. Ordinary `/test` keeps its original recovery behavior. No account fallback, user billing, internal key, or selection header is involved.

C17 permits the original administrator test semantics. Only check the request, selected account, and supported OpenAI protocol/type. Group is display origin. Do not add gateway availability or membership checks to manual evaluation starts. Inspection has its own automatic selection rules below. Those rules do not change manual tests.

C18-C22 keep memory-only active tasks and four history maps. A single account is busy across all evaluation types. Fold duplicate targets. Return actual `started`, `busy`, and `unavailable` arrays. Unstarted accounts create no synthetic records. Cancel stops further calls and closes the response. Release busy state after workers exit.

Save the latest 20 records per account/type to the existing evaluation `history.json`. Save each completed result. Keep memory results and report a save error. Preserve unreadable/null-file overwrite protection, a same-directory temporary file with mode 0600, atomic rename, and clear failure rollback. Clear does not cancel active tasks. Candy details stay in state. Other details use record IDs, not a second run-ID alias.

### Delivery and replacement

R03/R04 use the existing original Sub2API source and this Git plugin. Use standard temporary directories and existing dependency caches. Do not reintroduce the old adapter, compatibility profiles, AST pipeline, or source identity gates. Rebuild native frontend output when its inputs change.

Deliver the prebuilt static Linux/amd64 binary and archive under the fixed `dist/build/sub2api-codex-eval` names. Include `install.py` and required licenses. The user host does not fetch source or compile. No standalone cleanup script is included.

The installer selects a container explicitly, backs up and replaces `/app/sub2api`, and restarts only that container. Keep readiness checks and failure/interruption restoration. U7 treats normal connection resets as startup progress and prints a clear ready result. Keep raw errors for final failure. U8 creates default backups under `/tmp` and prints their actual path. Do not require `nc`, buildx, helper images, or Docker `--format`.

The v0.2.17 package supports upgrading the existing v0.2.16 installation and reinstalling v0.2.17. Check the candidate's actual version before replacement. Build from v0.2.17 source rather than changing only a version label.

## Historical inspection design before U15

### Boundary

Add a small inspection coordinator inside the existing native backend. The coordinator owns persisted inspection settings, one scheduled loop, one optional current run, and an in-memory inspection summary. It calls the existing Candy evaluation start/cancel/state paths and the existing native account update service. It does not create a second transport, user API key, scheduler service, container, or external script.

### Schedule

Read inspection settings through the native `SettingService` at process start. If enabled, draw a random delay in `[min_wait_minutes, max_wait_minutes]`. Keep the next run time in memory and return it with runtime state. A browser refresh does not change it. After a run ends, draw a new delay if automatic inspection is enabled. A manual run replaces the current scheduled slot. A current run blocks another automatic or manual run. Stop-current-run cancels its evaluation tasks and recovery timer. Automatic inspection remains enabled unless the administrator turns it off.

### Account and result flow

At the start of each run, read the native account list. Skip `schedulable=false`, non-active/runtime-restricted, busy, and configured Free accounts. U11 removes include/exclude name and email matching. Start Candy once per eligible account with the configured literal model, effort, and one test. Use returned run IDs and the existing state endpoint to collect only this run's records.

For an account at the normal priority, a valid wrong first result starts up to `max_tests - 1` confirmation tests. For an account at the degraded priority, a valid correct first result starts the same confirmation phase after the fixed 240-second recovery delay. Other first results finish after one test. A transition is valid only when all tests are complete, valid, and agree with the expected verdict. Errors, skips, cancellation, timeout, incomplete output, or mixed answers produce no update. The setting is named `max_tests_per_account`; values below 1 are rejected, and the default is 3.

### Priority and notifications

Use the native account update endpoint/service and its ascending priority semantics. Only update an account if its current priority equals the configured normal or degraded value. Keep all other values unchanged. Automatic updates default off; recommendation rows are still recorded when off. After a successful update, group notifications into recovery and degradation messages. Send configured ntfy and Bark notifications from the backend. Report notification errors in the inspection summary without retrying evaluation or reversing the account update.

### Settings UI

Add a fifth top-level tab named `设置` after the four evaluation tabs. Match native `SettingsView.vue`: a tablist with keyboard navigation, compact card sections, visible labels, native switches/selects/number inputs, and a right-aligned save action. Use two-column fields on wide screens and one column on narrow screens. Keep labels and units in the controls; omit helper paragraphs and feature explanations. Show only field-level validation and save status. Keep clear-history and refresh controls on the evaluation tabs.

### Persistence and lifecycle

Add typed inspection read and update methods to the native `SettingService`. Use its existing `SettingRepository.GetMultiple` and `SetMultiple` methods. Store inspection-specific keys in the existing PostgreSQL `settings` table. Use administrator settings routes and the native response format. Do not add a JSON file, a table, direct SQL, or a second storage service.

Follow the native secret-field contract. Read responses return configured status instead of saved credentials. Empty secret input preserves the saved value. Apply this rule to the ntfy credential and the Bark address, which contains the notification key. The backend reads the saved values when it sends a notification.

Validate settings before saving them. After a successful save, update the coordinator settings in memory. A failed save returns an error and does not change the active settings. A process restart reads the saved settings and draws a fresh delay. It does not restore an interrupted run. Stop the scheduler context and timer during server shutdown. No migration from the old Python configuration is required.

Complete the initial settings read before exposing the coordinator to requests.
A failed read blocks starts and retains its error until a successful native settings save configures the coordinator.
Absent keys use defaults only after a successful repository read.

Native source evidence: `setting_service.go:99-106` defines the repository contract. `repository/setting_repo.go:85-115` reads and writes setting keys. `ent/schema/setting.go:28-44` defines the existing table. `setting_update.go:196-202` and `setting_parse.go:338-344` show the native secret-field behavior. `server/routes/admin.go:567-571` registers native administrator settings routes.

### Deferred boundaries

The scope is Candy inspection only. Fingerprint, ModelTrace, and Pelican remain manual evaluations. There is no cron export, notification retry framework, historical inspection database, fixed-clock schedule, account enable/disable mutation, or new administrator authentication mechanism.

## Reusable Docker test environment

The user requires one test container for this task and future checks.
Use the fixed name `sub2api-codex-eval-test`.
Run PostgreSQL, Redis, and the tested Sub2API binary inside that disposable development container.
Retain its published browser port and synthetic state where practical.
Keep this development layout separate from the delivered installer and production deployment.
Remove earlier task containers after their required synthetic inputs are available in the retained container.
Do not prune volumes, images, caches, or unrelated containers.

## U10: dropdown controls and concise page copy

U11 replaces this section's account selection and account-menu requirements. Model dropdown and asynchronous edit behavior remain current.

The latest user instruction removes assistant-added explanatory small print.
Keep field labels, units, actual errors, progress, and save feedback.
Model selection must have a visible working dropdown with native raw catalog choices and the existing custom-model path.
Inspection account include/exclude fields use multiple-account dropdown selection.
Load the native account inventory needed by these controls instead of restricting inspection configuration to the currently displayed evaluation group.
Preserve saved wildcard array values and literal model names at the existing API boundary.
Use native visual conventions and lightweight page controls without a new frontend framework or backend configuration schema.
The custom-model option reveals the literal input inside the same field.
Account menus escape selected literal names and retain existing wildcard entries.
Account menus choose their direction and scroll height from the visible iframe viewport.
Asynchronous catalog and save responses preserve newer form edits.

## U11: remove unused settings and idle copy

U11 replaces U10 account selection and the earlier include/exclude inspection contract.
Remove include/exclude fields from the page, settings type, persistence field list, and coordinator selection.
Remove the account picker helpers, styles, and assertions that no longer serve a control.
Keep the independently loaded OpenAI inventory for raw model suggestions.
Keep the Free switch and native availability and busy checks.
The automatic priority switch controls visibility of both priority fields.
Use normal priority 10 by default and reject new settings below 1.
Read an existing stored nonpositive normal priority as the new default before validating the loaded settings.
This read does not write the database. The next explicit save stores the current value.
Validate loaded settings before scheduling. An invalid pair blocks starts and retains its error until an explicit valid save.
Keep explicit degraded values available for correction instead of silently replacing them.
When priority inputs are hidden, show any server save error in the existing visible save status.
Remove notification configured-status spans without changing secret redaction or blank-value preservation.
Hide idle status when there is no next run or actual error.
Keep actual progress, results, errors, and scheduled timestamps.

## U12: collapse automatic inspection configuration

The automatic inspection switch controls visibility of its configuration fields.
Remove visible switch captions from automatic inspection and priority headers.
Use accessible names on the switches without replacement small print or tooltips.
Keep current configuration values while collapsed and serialize them during save.
Move manual start/stop and actual runtime feedback outside the collapsed fields, near the existing settings actions.
This preserves manual inspection when scheduling is disabled without another card or explanatory text.

## Historical U13: evaluation account-list reuse before U15

The current inspection result renderer uses independent simplified rows.
Replace that renderer with the existing evaluation account-card and Candy row/record implementation.
Share account identity/state conversion and list behavior instead of maintaining a second version.
Use the inspection inventory and match runtime results to native account IDs.
Retain inspection-specific outcome, priority transition, counts, and actual errors on the matching account.
Preserve existing evaluation entry points and their filtering conditions.
Keep autonomous backend selection unchanged. Browser group or display filters must not redefine the scheduled target set.
Trace the original native account-list query contract before fixing the reported manual group leak.
Keep group-filtered manual account data separate from the inspection inventory.
Check every page and reject obsolete group responses before they update shared state.
Refresh native membership on ordinary polls. A loaded catalog must not retain accounts moved to another group.
Reuse raw model suggestions when the eligible discovery account IDs have not changed.
Use exclusive and shared memberships across two synthetic groups to verify lists and submitted target IDs.

## Historical U14: settings persistence and cancellation history

Trace form serialization, native settings writes and reads, coordinator configuration, and form refill.
Also inspect the existing per-evaluation preferences when reproducing the reported settings loss.
Repair the established owner rather than adding another store or compatibility layer.

Trace history errors through evaluation admission, task context creation, sender transport, and result persistence.
Distinguish explicit stop, process shutdown, timeout, and unintended request-context cancellation using the actual execution path.
Background evaluations must keep their task-owned cancellation lifetime after the HTTP start request completes.
Preserve explicit cancellation and real supplier errors.
Use local synthetic delayed responses to reproduce the identified failure without paid calls.
