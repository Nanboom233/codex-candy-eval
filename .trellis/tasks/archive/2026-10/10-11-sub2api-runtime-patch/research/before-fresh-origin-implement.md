# Administrator Evaluation Implementation Plan

Status: planning. Do not edit product code or run task.py start.
Approved: S01—S05 and C02/C03/C04/C06/C07. C08—C12 approved by user: 无误 继续. C13—C17 approved by user: 是. Next core batch: research/core-fourth-approval.md (C18—C22, pending).
Previous checklist remains in research/before-admin-channel-implement.md.

## 1. Complete planning

- [x] Reuse completed pristine v0.2.16 source/build/check evidence.
- [x] Record user choice of administrator account testing and new /accounts/:id/codex-eval endpoint.
- [x] Supersede internal user/Key/expiry/accounting/target-header gateway requirements.
- [x] Research native account-test and native protocol-conversion sending primitives.
- [x] Approve S01—S05: request, native sending/protocols, response, state effects and task integration.
- [x] Approve C02/C03/C04/C06/C07 per research/core-admin-channel-approval.md.
- [x] Approve C08—C12 per research/core-second-approval.md.
- [x] Approve C13—C17 per research/core-third-approval.md.
- [ ] Review remaining core/history/build/installer differences in complete five-item batches.
- [ ] Map retained/deleted items to current authorization and converge task artifacts.
- [ ] Present final implementation summary and await explicit approval, then task.py start.

No product builds/tests/deployment/cleanup in this stage. Do not repeat baseline checks
without reason or add race/paid/full-suite requirements.

## 2. Native administrator sending after approval

Ownership: Trellis native integration worker.
Files: fixed patch, new account handler/service files, evaluation runtime/executor integration.
Do not edit standalone UI or installer files in this workstream.

- Add the approved POST /api/v1/admin/accounts/:id/codex-eval route under native adminAuth.
- Base its service method on original OpenAI account connectivity-test sending code.
- Reuse selected-account credentials/model/URL/proxy/TLS/header/plugin primitives.
- Implement only approved full payload/response and protocol conversion contracts.
- Apply approved native state effects and task context/response ownership.
- Connect original evaluation tasks to the same sending method without internal HTTP/JWTs.
- Delete superseded user/Key/expiry/selector/diagnostic/partial-billing/maintenance frameworks.
- Restore unrelated native gateway middleware, scheduler, parsers, main and shutdown.
- Review upstream diff before compiling; preserve /test behavior.

Check with local fixture upstreams: administrator authorization, complete request,
full response/usage/error, account A-only failure paths, native selected protocol,
approved state changes and cancellation. Do not create billing rows for tests.

## 3. Approved standalone UI and evaluation core

Ownership: Trellis evaluation core/page worker.
Files: internal/eval and required embedded resources; coordinate executor interface.

- Apply approved F01—F06 without native Vue edits or extra activation/login flows.
- Preserve four question banks, algorithms, data, protocol signals and result/Pelican behavior.
- Apply only approved core/history differences and preserve received errors under B05.
- Remove withdrawn compatibility/catalog/body rewriting and identity machinery.
- Use existing relevant core/UI checks and local fixtures; no real history deletion.

## 4. Fixed release build

Ownership: Trellis build/delivery worker after approved backend/core inputs exist.
Files: fixed build scripts, patch/export inputs and packaging.

- Materialize pinned pristine v0.2.16 in a standard temporary development workspace.
- Apply only approved ordinary patch and required evaluation resources.
- Preserve upstream locks/declared tools and run applicable required/approved checks.
- Build static Linux/amd64 server with actual native embedded frontend.
- Output dist/build/sub2api-codex-eval and record source/build identity.

Do not run source retrieval/compilation/build-image preparation on the user's host.
Report failures/skips/platform differences accurately; reuse completed pure baseline evidence.

## 5. Installation, packaged cleanup and feature check

Ownership: Trellis build/delivery worker and subsequent check agent.
Files: install.py, cleanup.py, package script and package-facing documentation.

- Remove obsolete server activation/maintenance dependencies from installation.
- Retain approved explicit container selection/backup/replacement/restart/rollback.
- Include standalone cleanup.py, metadata/licenses and prebuilt binary.
- Output dist/build/sub2api-codex-eval.tar.gz.
- Verify package scripts and actual evaluation features on disposable deployments.
- Verify cleanup preserves unrelated images/caches, business volumes, database rows and history.

No private .env, CPA deployment, paid supplier request or unapproved remote deployment.
Report N150/iStoreOS verification as incomplete unless actually performed there.

## Agents and dependencies

Use Trellis implement/check agents only after task-wide final approval.
Assign file ownership, preserve concurrent work, and load current curated JSONL context.
Dependencies: full scope approval → native/core integration → fixed build → packaged checks.
One task jointly produces one package; old user-gateway records are historical, not implementation contracts.
