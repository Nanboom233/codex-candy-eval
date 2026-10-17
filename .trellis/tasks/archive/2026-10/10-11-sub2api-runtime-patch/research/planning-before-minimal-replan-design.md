# Pristine Baseline Evidence Design

## Source boundary

Resolve the official v0.2.16 annotated tag to its full commit. Materialize a fresh exact checkout/archive in system-temp.
Retain source/archive hashes and dependency-lock hashes. Verify source cleanliness before and after commands.
Apply no local plugin code, adapters, production edits or development overlays.

## Verification boundary

Read upstream Dockerfile, workflows, package scripts and testing instructions before choosing commands.
Match declared versions/tags. Distinguish native host tests, cross compilation and any emulated runtime execution.
Use actual declared service dependencies in isolated owned Docker resources where required.
Record environment-sensitive required settings without exposing private values. Do not use repository private .env.

## Evidence boundary

Keep complete command logs, process exit codes, durations, source/tool identities and failure test names.
Compare prerequisite/environment errors with assertion/race/compiler failures using targeted reproducible investigation.
No passing claim for an unexecuted, skipped or blocked check.
Do not infer upstream attribution from patched test history.
Preserve canonical delivery bytes and all maintained product code. Retain baseline workroot until evidence is reviewed.

## Decision boundary

This remains planning. The first deliverable is a baseline report, not a replacement plugin package.
Use actual results to determine the next planning step with the user.

## Individual patch approval

The user now requires removal of assistant-added race acceptance requirements and an item-by-item scope audit.
The complete baseline report remains historical evidence; race results do not block this task.
Review the structural adapter and runtime additions that actually enter the current binary before judging historical patch files.
Each approval item must identify concrete behavior and dependencies. Source facts come from code; retention/removal decisions belong to the user.
Keep pending product code and the existing package intact. Do not invent additional acceptance requirements or implement a batch cleanup.

## Breaking-change revision

The latest user instruction rejects compatibility configuration and treats changes as breaking. The current proposal is research/sub2api-breaking-change-list.md.
Freeze the target at pristine v0.2.16 Linux/amd64. Replace multi-version adaptation and legacy representations with one explicit new interface/state format and an ordinary version-specific patch.
Do not maintain old API translation, state/identity migrations, conditional shims or fallback implementations. Do not infer permission to delete user data from permission to remove compatibility code.
Review remaining feature boundaries and concrete deletion consequences individually before implementation. Prior inventories/specs are evidence only.

## Frontend first

Latest user steering pauses full-stack replanning. Follow research/frontend-first-approval.md.
Exclude all unapproved frontend behavior from the new proposal. Native v0.2.16 already provides administrator custom menus and URL-embedded pages, so F01 proposes one settings entry with zero native frontend source edits.
Do not infer approval for old UI transformations or the remaining backend/install plan from this entry decision. Concrete code rollback and new implementation must be described accurately and individually reviewed.

F01 is now approved in planning: /codex-eval through the existing administrator custom-menu settings, with no native Vue source changes. External settings remain untouched.
F02 is approved in planning: read the same-origin auth_token for each Bearer request, remove CPA key/login remnants, stop polling on missing token/401/403, and distinguish login expiry from missing permission. Native management owns session renewal; the standalone evaluation page adds no refresh or cross-window framework. Request/response redesign remains outside F02. Product edits remain unapplied.
F03a is approved in planning after user correction: the current standalone full-width group card misses the existing native-select styling and toolbar spacing. Place the group field first in each evaluation tab toolbar, sharing one selected group and reusing existing control styles and responsive grid. The API is confirmed; initial-selection and group-switch behavior are not approved by this display proposal.

F03b is revised by explicit user instruction: retain the existing interactive account status switch and use the native PUT /api/v1/admin/accounts/:id endpoint with only status active/inactive. Remove the permanently-false switch guard and placeholder operation; preserve original pending/run locks and failure restoration. Server-returned status updates the shared account display. Native toggleStatus confirms this contract; schedulable is separate and is not changed by this switch. The read-only proposal is rejected; no product edits have been applied.

F03c is approved in planning: reuse GET /api/v1/admin/accounts/:id/models for current-group accounts and take the model-id union for the original picker. Preserve original picker layout/sorting/preferences and evaluation-specific display rules; use transient page options, expose lookup errors, and discard stale responses after switching group. The group allowlist-candidates API only combines default and mapped names, so it is not presented as upstream discovery. Per-account execution/support checks remain outside this decision.

F04a is approved in planning: ordinary result markup currently undergoes nonceStyledHTML style-to-class rewriting because the evaluation shell introduced nonce-only style-src. Native default CSP allows inline styles. Restore original ordinary result rendering with that style permission and remove the transformation and ordinary row/result call sites as one reviewed item. Keep script nonces and same-origin embedding. Pelican HTML/script rewriting and theme choices remain separate, unresolved items.

F04b is approved in planning: remove added sandboxHTML, preserve the original escaped srcdoc and allow-scripts-only sandbox, list freeze and dialog cleanup. Original plugin.go:215 used inline-script permission; restore that permission only for the evaluation shell together with F04a style support. Remove the shell nonce-only script restriction rather than repairing/rewriting model HTML. No new preview endpoint or URL-token auth; keep existing same-origin framing and API administrator checks. The approved F04b decision supersedes F04a script-nonce retention only for the evaluation shell.

Current review follows the user’s five-decision batch instruction. F05a—F05e are approved in planning: direct operation call sites for /start, /state, /cancel, /records/:type/:account_id/:run_id and /clear-history under the evaluation base. Keep original controls and polling intervals; use numeric targets and actual task ids, avoid legacy api-wrapper request/response rewriting. Detail payloads use explicit typed result fields. History clear remains global per type/all and does not cancel runs. Current backend error compression remains an unresolved later review item; frontend cannot recover discarded diagnostics.

F06a—F06e are approved in planning after the full batch explanation. Use native account pagination/group/lite projection, keep original account/subscription controls using native fields, preserve existing group memory with stale-read invalidation, remove extra manual activation UI only together with a separately reviewed replacement of backend manual readiness dependency, and map minimal native light/dark tokens without flattening original subscription emphasis. No native Vue source edits or CPA compatibility paths. The existing Start readiness/identity guards make F06d a real dependency, not an implementation-ready isolated deletion.

The approved backend batch is research/backend-first-approval.md (B01—B05); product edits remain paused. Frontend activation removal does not approve automatically calling old Activate. Removing the singleton-owner permission must be paired with an explicitly reviewed initiating-actor accounting design; never silently bill another administrator. Native gateway error forwarding is reviewed separately from restoring details discarded by added evaluation layers.

Approved planning identity/lifecycle detail is in research/backend-identity-approval.md (B06—B10); product edits remain paused. It uses the native api_keys unique key for stable per-user lookup and an atomically past expiry, not a new mapping file or schema. Background workers must explicitly carry authenticated initiating UserID/APIKeyID, because manager tasks derive from the manager lifetime. Reserved identifiers are never network credentials; public ingress must reject them before simple/read-only shortcuts. Past expiry alone is not complete isolation on an original unpatched binary. No hidden recreation after native disable/delete, no automatic old identity migration, and no whole-database cleanup.

Latest user direction supersedes the expanded B11—B15 design. Authoritative current boundary: research/backend-minimal-key-header.md. Special internal-key validation plus one hard account-selection header; preserve actual initiating-user native accounting. Reuse native eligibility/slot/retry/forward/error/usage/disconnect/shutdown behavior. A target is mandatory, not sticky preference; failed/excluded target cannot fall back to a different account. Remove added diagnostic collection, CC parser alterations, extra fresh eligibility pass, custom sticky guards, partial billing hooks, upstream cancellation changes and producer/shutdown framework. Special-key representation/binding is the remaining technical detail; no product code changes or execution approval.
