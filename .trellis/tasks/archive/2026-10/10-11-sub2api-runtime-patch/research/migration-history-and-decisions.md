# Complete Migration History and Decisions

Date: 2026-10-16.

## Task identity and recovery

The task is the complete plugin migration to Sub2API. It first targeted v0.2.16. The user changed delivery to v0.2.17 during approved implementation. Native Candy inspection is one remaining part. It does not replace the four evaluations, page, native transport, task/history design, delivery plan, or previous approvals.

The main session incorrectly described the whole task as complete and presented the inspection plan as the complete migration plan. This recovery restores the whole-task scope and the links between requirements, decisions, implementation, and evidence.

The original approval files, history snapshots, and execution reports still exist. No missing source approval file was identified during this recovery. The main documents now include the detailed original contracts and the later corrections. Original records retain their original text and stage labels.

## Recovery sources

The local conversation index recovered an earlier full migration plan from Codex session `01a12946-c24c-7ff3-bd1e-294c86280cff`. Read it with `trellis mem extract 01a12946-c24 --grep '从纯净 Sub2API'`.

The recovered plan covers source strategy, all four algorithms, the native administrator sender, the page, tasks/history, native patch boundaries, delivery, checks, and Trellis execution. Later user decisions replace several parts of that early plan. Do not restore superseded source pins, cleanup, status switching, model mapping, or the zero-Vue restriction as active requirements.

The conversation tool reported partial older assistant history and unreadable encrypted inter-agent messages. This recovery does not claim a full transcript. The original task approval files and later check reports provide the available detailed evidence.

## Decision sequence

| Stage | Decision and later effect | Source |
| --- | --- | --- |
| Initial universal build/install approach | The user reported missing tools, Git transfer failures, Docker incompatibilities, slow builds, and ineffective integration. These failures remain defect evidence. | `history/`, early research reports |
| Binary replacement direction | Target prebuilt v0.2.16 Linux/amd64 for N150/8GB iStoreOS. Do not compile on the user host. Use fixed output paths. | `fixed-output-plan.md`, R04 |
| Pristine baseline | Compile and check original Sub2API before migration. Retain failures, skips, and platform limits. Reuse this evidence. | `pristine-v0.2.16-baseline.md` |
| R00-R01 | Remove the added race gate. Use breaking changes instead of compatibility configuration. | `patch-scope-approval-ledger.json` |
| F01-F06 | Review the frontend from the original plugin. Preserve real operations, renderers, filters, and task controls. | `frontend-first-approval.md` |
| B01-B10 | Earlier backend approvals removed runtime gates and proposed user accounting. Later administrator-test and literal-model decisions replace conflicting mechanisms. | `backend-first-approval.md`, `backend-identity-approval.md`, approval ledger |
| B11-B15 and internal-key proposals | The user rejected expanded gateway work. A special-key/header proposal followed, then a shared five-minute sliding key proposal. The later administrator-test direction replaced both. | `backend-dispatch-approval.md`, `minimal-replan.md`, `shared-key-expiry.md`, approval ledger |
| S01-S05 | Add a dedicated complete administrator account-test sender. Retire user/owner/key/billing and target-header designs. | `admin-account-codex-eval-approval.md` |
| C02-C22 | Preserve original evaluation behavior. Remove redundant transport, state, recovery, and storage frameworks. | Four core approval reports |
| R02-R04 | Remove redundancy without repeated questions. Rebuild from original inputs. Reuse existing source and baseline. Remove standalone cleanup and added gates. | `fresh-origin-rebuild.md`, approval ledger |
| Execution approval | The user approved execution and required Trellis agents. This approval applies to the reviewed migration scope. | Approval ledger, `implement.md` |
| Initial native implementation | Build the native package, page, sender, installer, and fixed artifacts. Review findings and repairs. | `fresh-native-implementation.md`, `fresh-native-check.md` |
| U1-U6 | Repair sidebar/iframe, disabled semantics, literal model selection, startup waiting, and nonstandard-port origin handling. | Correction and browser reports |
| U7-U9 | Clarify startup progress, use `/tmp` backups, and skip disabled-account model discovery. | Installer/disabled-model reports |
| Remaining inspection design | Move original Candy automation into the native process. Add configurable random waits, maximum tests, priority changes, and notifications. | `native-candy-inspection-plan.md` |
| Settings correction | Use native database settings. Replace the proposed JSON settings file. Keep next-run time in memory. | Current PRD/design, native settings evidence |
| Native response-model audit review | Relevant native gateway files match the original source. Plugin model comparison is outside scope and is not an extra test requirement. | Response-model research reports |
| Integrated execution approval | The user approved the complete twelve-part migration plan with `批准实施`. Remaining inspection/settings implementation and local delivery checks are authorized. | Current `task.json` and `implement.md` |
| Orchestration correction | Add real specs to both context manifests. One implement agent owns remaining cross-layer implementation and delivery. Independent check follows. | `trellis-orchestration-correction.md` |
| Target version update | The user requested Sub2API v0.2.17. Use its actual source for the final build. Keep v0.2.16 reports as historical evidence. | Current task/design and version integration research |

## Frontend decisions and current mechanisms

| ID | Required behavior | Later correction |
| --- | --- | --- |
| F01 | Serve `/codex-eval` as the independent evaluation page. | U5 replaces custom-menu-only and zero-Vue restrictions with the native sidebar/icon/iframe route. |
| F02 | Use native same-origin administrator `auth_token`. Distinguish missing/expired login from HTTP 403. | No second login or retained backend JWT. |
| F03a | Put the shared group selector first in each original evaluation toolbar. | Keep original control and responsive grid styles. |
| F03b | Preserve the account calling switch and use an actual native update operation. | U2 replaces the earlier `status` mechanism with native `schedulable`. |
| F03c | Preserve model picker behavior and temporary current-group suggestions. Show lookup errors. | U6 uses evaluation-specific raw models. U9 skips disabled-account discovery. |
| F04a | Preserve ordinary Markdown/math output and original styles. | No style-to-class rewrite or new renderer restriction. |
| F04b | Preserve original Pelican HTML, frozen list preview, detail playback, and cleanup. | Keep script-only sandbox without same-origin permission. |
| F05a | Start selected evaluations and use actual start summary arrays. | Native administrator start API. |
| F05b | Preserve progress and 2.5-second active/20-second idle polling. | State polling does not fetch model catalogs. |
| F05c | Cancel the matching run and show cancelling until the worker stops. | Do not treat the cancel response as worker completion. |
| F05d | Preserve history, details, and previews. Ignore stale detail responses. | Typed record API. |
| F05e | Preserve explicit per-type/all-history clear controls and failure handling. | Clearing history does not cancel active work. |
| F06a | Read all current-group account pages and show native account facts. | Do not synthesize missing priority or account fields. |
| F06b | Keep type/subscription/Pro ALL filters, badges, and native plan fields. | Separate disabled calling from runtime restrictions. |
| F06c | Remember a valid group. Clear selections and stale group responses on change. | Group changes do not cancel backend tasks or delete history. |
| F06d | Remove activation controls and manual preparation from normal usage. | S01-S05 remove the internal identity prerequisite. |
| F06e | Use native light/dark colors while keeping original layout and plan accents. | U5 permits the necessary native frontend edits. |

The original text and source anchors remain in `frontend-first-approval.md`. Its planning status records that stage, not the current implementation state.

## Earlier backend approvals and later replacements

These approvals remain historical records. Apply S01-S05 and U3/U6 when an earlier mechanism conflicts with them.

| ID | Earlier approved decision | Current treatment |
| --- | --- | --- |
| B01 | Remove menu and deployment-proof readiness gates. Keep actual dependency checks. | Retain the reduced runtime boundary. |
| B02 | Use native administrator permissions without a first-activator owner. Retain initiating actor for proposed accounting. | Native administrator permissions remain. S01-S05 remove user accounting. |
| B03 | Remove the redundant `/targets` list and DTO. Validate targets on the server. | Use native account lookup and C17 administrator-test admission. |
| B04 | Use native account model eligibility without client-catalog admission or a new paid probe. | U3/U6 require literal evaluation models and raw suggestions. |
| B05 | Preserve actual errors discarded by evaluation layers. Review native rewriting separately. | S03/C06 preserve actual errors with native sensitive-data redaction. |
| B06 | Prepare one expired native accounting row per initiating administrator. | S01-S05 remove this row and the user-accounting design. |
| B07 | Prepare accounting during normal initialization or the first valid start. Remove manual activation. | Manual activation remains removed. No accounting preparation is needed. |
| B08 | Use task-bound user/key IDs and local dispatch context. Remove identity files. | S01-S05 remove user/key IDs and the gateway dispatch mechanism. |
| B09 | Bind groups per request and use initiating-user permissions, billing, and limits. | Groups describe page filters and result origin. User billing and gateway limits do not apply. |
| B10 | Remove lifecycle/maintenance endpoints, socket, and CLI. Keep cancellation and proposed native key controls. | Cancellation remains. S01-S05 remove the key controls. |

The original authorizations remain in `patch-scope-approval-ledger.json`. The internal-key and five-minute expiry designs also remain as superseded history.

## Administrator sending decisions

| ID | Current contract |
| --- | --- |
| S01 | Separate `POST /api/v1/admin/accounts/:id/codex-eval`. Complete streaming Responses input. One selected-account call. |
| S02 | Use the native account-test service, credentials, shadows, URL validation, proxy, TLS, headers, plugin transport, and native protocol conversion. U3/U6 replace evaluation mapping with literal model sending. |
| S03 | Return complete stream/terminal/observed usage and actual HTTP/SSE errors. Apply native sensitive-data redaction. No user billing. |
| S04 | Keep native snapshots and 401/429 handling. Omit broad successful-test recovery only on the new path. Ordinary `/test` is unchanged. |
| S05 | HTTP and four background evaluations share one native sender. Backend tasks use their context directly. No local HTTP, fake Gin, or administrator JWT storage. |

See `admin-account-codex-eval-approval.md` and `native-connectivity-dispatch.md` for the original reasoning and source anchors.

## Core decisions

| ID | Contract retained in the complete design |
| --- | --- |
| C01 | Required host boundary is covered by S02/S05. |
| C02 | One model/effort pair and account/group/name/run/status metadata. Remove duplicate requested/effective fields. |
| C03 | Separate model and effort. Restore original effort choices. Remove CPA suffix parsing and model-support guesses. |
| C04 | Keep original per-evaluation parameter rules in one preparation path. |
| C05 | Remove the user-gateway executor. S02/S05 provide the native sender. |
| C06 | Keep actual errors, original panic handling, response closure, and cancellation. |
| C07 | Read the Responses stream incrementally with original terminal/text/item semantics and native `response.done`. |
| C08 | Remove added transport byte/event caps. Keep original algorithm answer limits. |
| C09 | Keep observed token values and retry-attempt sums. Remove unused usage flags and fabricated billing. |
| C10 | Keep original probe retries, waits, slots, and timeouts. Remove unused execution wrappers and repeated input copies. |
| C11 | Keep original Codex prompt/context/tools/headers/client signals. Derive stable installation identity from the native account ID. |
| C12 | Keep Candy prompt, standalone `21` scoring, complete-answer scoring, 4,000-character saved display, and wrong-answer semantics. |
| C13 | Keep Fingerprint modes, samples, baseline, normalization, attribution, concurrency, and error-stop rules. |
| C14 | Keep ModelTrace bank, challenge generation/order, valid-sample rules, partial outcomes, and attribution. |
| C15 | Keep Pelican prompt, extraction, 256 KiB limit, and complete HTML details without rewriting. |
| C16 | Keep visible progress and cancelling state. No extra progress bus, persistent log, or duplicate result counters. |
| C17 | Keep administrator account-test admission and account-wide busy state. Fold duplicates and preserve real start summaries. Unstarted accounts create no history. |
| C18 | Keep serial Candy batches and save each actual result. Stop further rounds on cancellation. |
| C19 | Keep active tasks/busy/progress in memory and persist results only. Remove ActiveIntent, recovery schema, and synthetic restart records. |
| C20 | Keep JSON/null/read-error overwrite protection and latest 20 records per account/type. Do not add a strict business schema or migration. |
| C21 | Keep same-directory temporary writes, mode 0600, close, and atomic rename. Remove repeated path and directory-permission policies. |
| C22 | Keep original detail/clear APIs and clear failure rollback. Remove unused Flush and record/run-ID aliases. |

The original rationale and source anchors remain in `core-admin-channel-approval.md`, `core-second-approval.md`, `core-third-approval.md`, and `core-fourth-approval.md`.

## Corrections and source reports

| IDs | Current correction | Implementation and check evidence |
| --- | --- | --- |
| U1-U4 | Sidebar/iframe entry, actual disabled semantics, literal models, and selected-container readiness. | `fresh-native-corrections-implementation.md`, `fresh-native-corrections-check.md` |
| U5-U6 | Direct native sidebar/icon/iframe, complete browser origin including port, raw model discovery, and real UI-to-supplier model checks. | `native-browser-model-repair.md`, `native-browser-verification.md`, `native-browser-model-check.md` |
| U7-U9 | Normal startup messages, clear ready result, default `/tmp` backups, and disabled-account model discovery skip. | `installer-disabled-model-repair.md`, `installer-disabled-model-check.md` |

The U5/U6 Chrome evidence replaces the earlier unperformed-browser status for those affected paths. It does not prove N150/iStoreOS or production HTTPS execution. The failed three-second installer fixture and passed eight-second rerun remain in the U7-U9 reports.

## Remaining migration work

Native inspection and settings now have implementation, v0.2.17 build, Docker/browser verification, and independent review evidence. The design includes:

- Add a fifth settings tab, styled after native Sub2API settings. Keep labels, units, status, and field errors. Omit helper paragraphs.
- Store settings through native `SettingService` and `SettingRepository` in PostgreSQL `settings`.
- Schedule Candy inside the application process. Default random waits are 15-35 minutes after each round.
- Interpret maximum tests as one first test plus up to `N-1` confirmation tests. Default N is 3.
- Keep the original 240-second recovery wait, include/exclude/Free filtering, managed priority values, and ntfy/Bark behavior.
- Use lower-is-preferred Sub2API priority semantics. Default automatic updates off.
- Keep runtime scheduling in memory. Restore configuration after restart and draw a new delay. Do not resume interrupted tests.

The detailed subplan is `native-candy-inspection-plan.md`. The original generic inspection API supported more evaluation types. That fact does not expand the approved Candy scheduling behavior.

The v0.2.17 binary and archive include inspection and settings. Independent review repaired startup settings-load admission. The rebuilt candidate passed a final browser smoke in the retained container. Commit and task closure remain pending.

## Historical designs that stay historical

| Proposal | Current treatment |
| --- | --- |
| Fixed new source exports, source manifests, repeated archive/blob identity gates | R04 removed these requirements. Existing evidence stays readable. |
| Compatibility profiles, AST adapters, old overlays | R01/R03 replace them with original inputs and a small ordinary patch. |
| User owner, internal key, expiry renewal, hard-selection header, gateway billing/sticky hooks | S01-S05 replace them with native administrator account testing. |
| Custom-menu-only entry, zero native Vue edits | U5 replaces this mechanism with native sidebar/icon/iframe edits. |
| Calling switch implemented with `status` | U2 requires `schedulable` instead. |
| Evaluation model mapping | U3/U6 require literal model sending and raw evaluation suggestions. |
| Standalone cleanup and broad build-cache removal | R04 removed these from delivery. The later user request authorizes only scoped residual task-container removal. |
| Separate inspection JSON file or external cron/Python runtime | Native process and native database settings replace these proposals. |
| Response-model comparison as a plugin feature or missing-test defect | The user explicitly excluded it. Preserve native audit behavior. |

## Preserved snapshot index

- `research/patch-scope-approval-ledger.json`: original approval sequence and later correction records.
- `research/before-fresh-origin-prd.md`, `before-fresh-origin-design.md`, `before-fresh-origin-implement.md`, and `before-fresh-origin-task.json`: prior stage before original-source reconstruction.
- `research/before-fresh-origin-implement.jsonl` and `before-fresh-origin-check.jsonl`: prior agent context.
- `research/before-admin-channel-minimal-replan.md` and `before-shared-expiry-minimal-replan.md`: earlier transport/identity proposals.
- `history/`: delivery, baseline, context, and coordinator snapshots from earlier phases.
- `research/fresh-origin-rebuild.md` and `fresh-origin-port-map.md`: current source strategy and native integration map.

These records preserve earlier designs and their reasoning. Later user decisions determine which mechanisms apply now. A discussion of remaining work never replaces this whole-task history.

## Restoration checks at the earlier planning stage

The task metadata and approval ledger parse as JSON. Both agent context manifests parse as JSONL, and all referenced context files exist. Core execution authorization remains recorded. Inspection implementation authorization remains false. This restoration changes task documents only. It does not run product builds or tests, deploy, commit, or close the task.

## Later execution and Docker decision

The user later approved the complete integrated plan with `批准实施`.
Inspection implementation authorization is now true.
The user changed the delivery target to Sub2API v0.2.17.
The native build, focused checks, Docker integration, and actual browser review passed.
Independent review also repaired startup admission after a failed settings read.
The rebuilt candidate completed inspection through the actual browser with no console errors.

The user now requires one reusable Docker test container, including its dependencies.
Use the fixed name `sub2api-codex-eval-test`.
The implementation agent consolidated PostgreSQL, Redis, the synthetic supplier, and Sub2API into the retained application container.
The main session removed all 11 identified earlier task containers.
The fixed container remains available at `http://127.0.0.1:60198` for later reuse.
This authorization does not include unrelated containers, images, volumes, caches, or real application data.
