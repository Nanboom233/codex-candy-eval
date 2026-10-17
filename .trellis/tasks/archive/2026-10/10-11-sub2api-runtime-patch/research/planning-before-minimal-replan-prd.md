# Sub2API Recovery Planning: Review Every Patch

## Goal

Use the completed unmodified official Sub2API v0.2.16 baseline to review each maintained plugin/native integration change and remove unnecessary scope.
Present concrete changes individually for the user's approval before changing product code.
The user explicitly returns this task to planning. Prior implementation approval no longer covers further product edits.

## Background and confirmed evidence

The user reports that the installed evaluation remains unusable after repeated installation compatibility fixes.
The latest package is not accepted as a working product.
The previous patched runtime came from v0.2.16 commit 01842576b3d5c1e2545b5a99ad555d36e28ed306 and targets Linux/amd64.
Previous patched full native unit race failed with 51 race warnings and 153 top-level failures.
No pristine comparison establishes attribution. Prior compiler or installer/lifecycle passes cannot answer this baseline question.
Historical detailed requirements, evidence and contracts remain preserved under history/before-pristine-baseline-planning-2026-10-14/.

## Current requirements

B1. Verify the official annotated v0.2.16 tag and exact source commit again.
B2. Use a fresh unmodified upstream tree. Apply no plugin adapters, patches, overlays, generated extension files or installer facts.
B3. Inspect upstream build/CI/testing contracts and use their declared Go/Node/package-manager versions, tags and configuration.
B4. Compile the real frontend and static Linux/amd64 server, preserving dependency lock files and source identity.
B5. Run exact upstream unit/integration and frontend critical CI gates, builds and declared lint/helper checks.
B5a. Use only actual upstream required checks and explicitly approved checks as acceptance requirements. Remove the assistant-added race gate. Completed extra-check logs remain historical evidence, not requirements or blockers.
B6. Record exact commands, tool versions, host/target architecture, exit codes, durations, complete stdout/stderr and per-test failures.
B7. Classify concrete failures with evidence. Do not declare failures upstream defects or environment defects without establishing attribution.
B8. Preserve baseline source/locks and diagnostic outputs for independent review. Do not clean evidence while investigating.
B9. Report build results separately from tests. A successful build cannot justify a passing test claim.
B10. Inventory actual structural rules, added runtime files, evaluation-core changes and build/installation additions. Distinguish current build inputs from historical patch files.
B11. Give each independently reviewable change an ID, exact source anchors, behavior, requirement provenance, dependencies, removal effect and retain/remove proposal. Do not treat an earlier spec or implementation as current approval.
B12. Record each concrete decision. Latest user steering requests batches of five decisions with complete explanations of current problems, proposed behavior, scope and limitations; batch approval may cover the explicitly displayed items. Do not infer product implementation approval or execute blanket deletions from batch presentation.
B13. The user rejects all compatibility configuration and accepts breaking changes. The new plan targets only v0.2.16 Linux/amd64; no legacy API/data/identity migration, dual-read/write, multi-version adapter or alternate implementation paths.
B14. Use research/sub2api-breaking-change-list.md as the current proposed scope. The previous inventory describes existing code, not retained requirements. Concrete lifecycle/data/cleanup changes still receive individual approval; no database/history deletion is implied.
B15. Latest steering supersedes B14's execution order: reset frontend assumptions and review frontend first. research/frontend-first-approval.md is the current proposal; backend/installer redesign is deferred.
B16. No prior frontend modification is retained by default. Prefer existing Sub2API custom-menu/embedded-page capabilities; approve the entry, session/data integration and each further UI behavior separately.
B17. F01 is approved in planning by the user's “继续”: /codex-eval, one existing administrator custom-menu URL entry named 账户评估 at the end, with zero native frontend source edits. No settings write or product rollback has been performed.
B18. F02 is approved in planning by the user's “是”: use the existing same-origin auth_token for each Bearer request, remove separate CPA key/login remnants, stop polling on missing token/401/403 and distinguish expired login from missing administrator permission. Native management owns session renewal. Product changes remain unapplied.
B19. The user clarifies F03a concerns the group control in the current evaluation UI; the API is correct. Review its placement and styling against the original evaluation toolbar, without broad frontend redesign. F03a is approved by “对 继续”: a group field first in each existing tab toolbar, reusing existing select styles with one shared selected group. Initial-selection policy and account/model behavior remain unapproved.

B20. The user rejects read-only replacement and instructs restoring the existing interactive account enable switch through a native API. F03b retains original row/filter behavior and calls PUT /api/v1/admin/accounts/:id with status active/inactive, using the existing administrator session. Preserve original in-progress/submitting locks and failure restoration; reflect the server-returned status. Do not conflate status with schedulable. This is recorded in planning, not implemented.

B21. F03c is approved in planning by the user’s “对”: retain the original model picker UI/sorting/preferences, using per-account native models responses to build a deduplicated current-group option list in page memory. No persistent catalog/TTL/signature compatibility framework; report lookup failures without treating unknown availability as unsupported. Group allowlist candidates are not upstream discovery. Execution/pre-start behavior remains a separate decision. Product implementation remains pending the final planning review.

B22. F04a is approved in planning by “对”: restore original ordinary-result styles for candy math and ModelTrace probability bars, replacing assistant-added style-to-class rewriting with the inline-style capability present in native default CSP. Preserve original content/layout/interactions. Pelican preview scripts/isolation and themes are separate unapproved items; do not remove a rendering workaround without its matching style-policy correction.

B23. F04b is approved in planning by “对，接下来5个决策一批次展示”: original Pelican srcdoc rendering without added HTML reparsing/nonce/style rewriting. Preserve allow-scripts-only sandbox, original list freeze and dialog playback/cleanup. Restore original evaluation-page inline-script permission only for its shell so model and freeze scripts run without rewriting; global native CSP remains untouched. Product implementation remains pending the final planning review.

B24. Next review batch contains exactly five decisions F05a—F05e: preserve original start, progress, stop, record-detail and clear-history UI and wire directly to explicit evaluation endpoints. Batch scope/observable behavior is recorded in research/frontend-first-approval.md. User approved all five after the complete explanation with “批准 继续”; product changes remain unapplied. Global clear-history scope and returned-error limitations must be visible; no actual user-record deletion is authorized by planning.

B25. Frontend review batch F06a—F06e approves native paginated account reads with actual display fields, restored type/subscription filters, existing remembered-group/first-group selection and clean switching, removal of extra manual activation UI paired with a separately reviewed backend prerequisite solution, and minimal native light/dark token mapping preserving original evaluation layout and subscription emphasis. All five are approved in planning by “继续 这不是能想明白吗 之前搞这么垃圾”. Product implementation remains paused; the backend replacement for F06d is not approved.

B26. Backend decisions B01—B05 are approved in planning by “批准 继续” in research/backend-first-approval.md: runtime readiness without menu/deployment proof, native administrator permissions with authenticated initiating actor, redundant target-list removal with server validation, native model eligibility, and original error preservation in evaluation layers. Accounting identity/normal initialization and necessary native error forwarding hooks remain separate unapproved decisions.

B27. Approved planning batch B06—B10 in research/backend-identity-approval.md proposes per-initiating-user expired native accounting rows, ordinary initialization/first-valid-start preparation, task-bound process-local dispatch without extra token/identity files, request-local group binding with native user billing/limits, and deletion of extra lifecycle maintenance protocols. Native API-key creation limits apply. Disabled/deleted identities are not automatically recreated; no native schema redesign, legacy migration or implicit old-data deletion. User approved this batch with “批准 继续”; product implementation remains paused.

B28. User supersedes unapproved expanded B11—B15 with exactly two gateway concerns: a special internal key and a special selection header. Validate the key and hard-select the specified account; all remaining gateway behavior belongs to native Sub2API. Current scope is research/backend-minimal-key-header.md. Withdraw added raw-diagnostic/CC-parser/post-slot DB recheck/sticky/partial-billing/upstream-cancellation/producer-barrier/shutdown-coordination designs. Restore upstream behavior where old patches added them. B08 context-capability authentication form is superseded; initiating-user native accounting and B09 group/billing remain. Exact key binding must be clarified technically without calling a predictable accounting identifier secret. No implementation approval.

## Acceptance criteria

A1. The tested source hashes/revision and source diff prove an untouched official v0.2.16 baseline.
A2. The Linux/amd64 artifact is produced by the actual upstream build and its version/architecture are verified.
A3. Every applicable upstream required check has an explicit pass/fail/blocked result with original diagnostics.
A4. A failed prerequisite is resolved where possible without modifying upstream source, then affected checks are rerun.
A5. Unresolved failures are retained with exact reproduction and a proposed investigation step for planning.
A6. Main reviews evidence before proposing any plugin implementation or new release.
A7. Every product change receives an individual recorded approval against the current inventory before removal or retention is implemented.

## Scope boundaries

Authorized now: planning, source/tool/CI inspection and real pristine baseline compilation/testing in isolated disposable workspaces.
Keep canonical dist/build package and binary untouched. Baseline artifacts stay in standard system-temp research workspaces.
Product code remains unchanged while items await individual approval. The race acceptance requirement is explicitly removed by the user now.
No upstream test deletion, fixture rewrite, race suppression, weakened assertions or skipped failing checks to obtain a green baseline.
No production deployment, paid request, commit, publication or task archival.
Further remediation and integration planning depends on the actual baseline results.
