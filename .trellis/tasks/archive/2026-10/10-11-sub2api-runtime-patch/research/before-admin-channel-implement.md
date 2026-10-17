# Minimal Implementation Plan

Latest steering: pause C01—C05 and investigate extending native account connectivity tests.
See research/native-connectivity-dispatch.md. The Key/local-gateway requirements below
are the prior planning branch under review, not authorization to proceed. The administrator-test
alternative changes user billing and scheduler semantics and awaits user choice.

Status: planning. Do not edit product code or run task.py start yet.
Full user-facing plan: research/minimal-replan.md.
Historical planning checklist is preserved in research/planning-before-minimal-replan-implement.md.

## 1. Complete scope review

- [x] Read completed pristine v0.2.16 evidence and existing approval records.
- [x] Replace rejected expanded B11—B15 with the internal-key and target-header proposal.
- [x] Record user rejection of per-user and per-request Key designs.
- [x] Record user-approved five-minute expiry, reset before each evaluation call using the same Key.
- [x] Identify native Update/ExpiresAt and authentication cache invalidation for renewal.
- [ ] Determine concrete shared native owner preparation.
- [x] Record concrete native hooks to restore and required integration boundaries.
- [ ] Obtain review of this minimal backend proposal.
- [ ] Review remaining core/history/build/installer items in complete five-item batches.
- [ ] Map every retained/deleted item to current authorization in the ledger.
- [ ] Converge PRD/design/checklist and present a final implementation summary.
- [ ] Await explicit final approval, then run task.py start.

Keep existing product artifacts and user data intact during this step.
Do not rerun the pure baseline or add race/paid/full-suite requirements.

## 2. Implement minimal backend after approval

Ownership: native runtime/authentication/selection worker.
Files: fixed native patch, required backend runtime/handler resources, evaluation executor integration.
Do not edit frontend resources or installation scripts in this workstream.

- Replace capability/private identity files with one shared internal Key and native accounting row.
- Preserve native expiry checks and reject missing/unlimited expiry for this internal Key.
- Reuse native IsExpired() in internal-Key validation before the simple-mode early return.
- Before each trusted evaluation HTTP call, set the same Key's ExpiresAt to now + 5 minutes
  through native Update/cache invalidation; stop dispatch if renewal fails.
- Gateway bearer requests never renew expiry. Do not issue per-request keys, consume tokens
  or bind body hashes. Native expiry renewal does not repair disabled/deleted data or reset quota.
- Preserve native owner accounting and request-local group binding.
- Accept the special key and one target header only on native POST /v1/responses.
- Reuse native targeted selection and terminate candidate fallback.
- Replace the custom pipe transport with native local HTTP requests.
- Delete withdrawn selector/WS/sticky/freshness/diagnostic/usage/cancel/shutdown modifications.
- Restore native main, middleware continuation, forward parsers and billing behavior where altered.
- Review the resulting upstream diff before compiling.

Use local fixture upstreams to check A-only forwarding, A failure without B, invalid key/header rejection and ordinary-key behavior.
Check shared Key owner native accounting and limit behavior without paid supplier requests.

## 3. Apply approved frontend/core work

Ownership: evaluation core and standalone page worker.
Files: internal/eval and necessary embedded-page resources.
Coordinate its executor interface with the backend worker before integration.

- Apply approved F01—F06 without native Vue changes.
- Restore original rendering, controls, account operations and direct approved API calls.
- Preserve original four question banks, algorithms, data and protocol signals.
- Apply only individually approved core/history differences.
- Preserve received errors under B05 without adding native diagnostic capture.
- Remove compatibility, catalog, activation UI and HTML rewriting remnants.

Use existing relevant core/UI checks and a local fixture deployment.
Do not create new acceptance frameworks or perform real history deletion.

## 4. Produce the fixed release

Ownership: build/delivery worker after approved backend/core inputs exist.
Files: fixed build scripts, fixed patch/export inputs and release packaging.

- Materialize pinned pristine v0.2.16 in a standard temporary development workspace.
- Apply only the approved ordinary patch and required evaluation resources.
- Preserve upstream lockfiles and declared tool versions.
- Run actual applicable native required checks and already approved relevant checks.
- Build the embedded static Linux/amd64 server in the development environment.
- Output dist/build/sub2api-codex-eval.

Reuse completed pure-baseline evidence instead of repeating it without reason.
Record failures, skips, platform differences and build identity separately.
No compilation, git retrieval or build image setup runs on the user's host.

## 5. Package and check installation/cleanup

Ownership: build/delivery worker, followed by Trellis check agent review.
Files: install.py, cleanup.py, package script and package-facing documentation.

- Remove obsolete server maintenance/activation dependencies from install.py.
- Retain explicitly approved container selection, backup, replacement, restart and rollback behavior.
- Include standalone cleanup.py and necessary metadata/licenses.
- Output dist/build/sub2api-codex-eval.tar.gz.
- Check package-contained scripts against disposable Docker resources.
- Verify cleanup avoids unrelated images/caches, business volumes, databases, keys and history.
- Verify the feature through the packaged binary, not only installer success.

Use existing relevant package tests and local fixture/UI flows.
No private .env, CPA deployment, real supplier evaluation or unapproved remote deployment.
Report N150/iStoreOS testing as unverified unless actually performed on that host.

## Agent and dependency contract

Use Trellis implement/check agents only after task-wide approval.
Every agent has file ownership and must preserve unrelated concurrent edits.
Dependencies: scope approval → backend/core integration → fixed build → packaged lifecycle/feature checks.
Keep one task because these workstreams jointly produce one package.
implement.jsonl/check.jsonl contain real current research and approval entries.
Historical expanded gateway designs do not authorize implementation.
