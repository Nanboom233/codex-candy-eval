# Independent U1–U4 final check

Status: completed. Independent Trellis reviewer resumed after the frozen
correction implementation. Current PRD U1–U4 supersede earlier manual-menu,
status-only switch and model-mapping details. Main owns spec/task updates and
commit. Reviewer made no product changes during this pass.

## Findings (fixed)

No further in-scope defect found. Earlier independent streaming, Pelican status,
cancellation copy, theme and installer interruption fixes remain present.

## Full-scope review

- **U1 menu/iframe:** `native/sub2api.patch` now modifies four original files:
  router registration, administrator routes, standalone page bypass and one
  `GetSettings` custom-menu response expression. `codex_eval_menu.go` preserves
  unrelated stored entries, removes the stale administrator canonical entry and
  appends the admin-only `账户评估` with same-origin absolute URL and final sort
  order. It uses the existing native origin helper and performs no repository
  write. Original `AppSidebar.vue` consumes administrator settings, filters
  `visibility=admin` and creates `/custom/<id>`; original `CustomPageView.vue`
  resolves that menu and calls unchanged `buildEmbeddedUrl`, then renders its
  iframe. Exact `/codex-eval` bypass and page-local SAMEORIGIN/CSP remain intact.
  No native Vue source or embedded frontend build input changed.
- **U2 disabled/status:** page `stateCredentials` derives `disabled` from
  `schedulable===false`. The switch calls the existing native POST
  `/api/v1/admin/accounts/:id/schedulable` with its boolean body and adopts the
  returned flag. The existing native handler calls `SetAccountSchedulable` and
  returns runtime-aware account data. The page shows `已停用` separately from
  inactive/error/quota status, future rate-limit reset, overload and temporary
  unschedulable facts. Optimistic state/failure restoration, busy locks and
  enabled-only filtering remain coherent. This introduces no backend gateway
  eligibility check into administrator evaluation admission.
- **U3 literal model:** shared evaluation sender no longer calls account model
  mapping or OAuth normalization. Its complete request keeps the chosen model
  before native Responses/Chat conversion; credentials, shadow resolution,
  endpoint validation, proxy, TLS, headers and plugin transport are retained.
  All four model controls accept direct names with catalog datalist suggestions;
  an empty/config-derived catalog cannot block the model input. Original `/test`,
  model-discovery/capability probing and ordinary gateway mappings remain outside
  this patch. The corresponding sender assertions now expect the chosen model
  despite a deliberately conflicting configured mapping.
- **U4 readiness:** installer receives the exact selected container ID, refreshes
  inspect state/address each retry, handles restarting separately, and exits on
  stopped/dead/paused application state with exit/OOM/runtime diagnostics. HTTP
  requests and retry sleep use the remaining readiness time; progress includes
  elapsed limit, state, restarts, endpoint and real error. Backup precedes mutation;
  failure/Ctrl-C restore remains, including Docker cp restoration without exec
  when the candidate stopped the container. No nc/buildx/helper/source-build
  dependency returned.
- Four original algorithms/resources/request shapes/options/retries, memory-only
  tasks/busy/cancellation and four history-map loading/save/clear behavior were
  rechecked against the previous independent review and current caller changes.
  No new Key/user/billing/scheduler/sticky/recovery/compatibility framework,
  standalone cleanup or redundant baseline/test requirement was introduced.
- Current build/check/package entry and release/documentation wiring use the
  fresh native patch/module and fixed outputs; retired directories are not inputs.

## Verification and evidence

Reused completed pristine baseline and unaffected evidence from
`fresh-native-check.md`, without repeating baseline, race, full suite or coverage
checks. Reused correction build/runtime evidence only after inspecting the
current code and actual output files:

- `workroot/corrections-build.log`: relevant evaluation/service/admin tests,
  touched-package vet and final embed/static Linux/amd64 build all completed
  successfully using the existing original frontend and caches.
- `workroot/smoke-corrections/outcome.json`: final binary install, native auth,
  automatic settings menu/custom route, direct sender/background Candy and actual
  malformed-candidate stopped-container restore passed.
- `workroot/smoke-corrections/ui-corrections-checks.json`: served standalone JS
  used real local APIs, existing native sidebar and URL-builder source, schedulable
  POST/off/disabled facts and editable unmapped model. Seeded unrelated menu items
  survived; original DB seed stayed unchanged after settings reads.
- `workroot/smoke-corrections/native-restart-readiness.log`: actual native restart
  produced RemoteDisconnected at the refreshed endpoint and then recovered.
  `native-address-change.json` records changed published port.
- `real-tcp-reset.log`: actual TCP reset (`Errno54` on macOS) retried and recovered.
  `reset-timeout.log` and `reset-restart-result.json` record bounded 2-second
  failure and recovery across actual Docker restart/published-port changes.
- Reviewer rechecked changed JavaScript syntax, Python syntax, Go formatting and
  diff whitespace; current archive contains the exact current installer.
- Fixed current outputs remain `dist/build/sub2api-codex-eval` and
  `dist/build/sub2api-codex-eval.tar.gz`. This no-change review required no rebuild
  or repeat of passing Docker scenarios.

Workroot:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8`.

## Findings (not fixed) and limits

No outstanding in-scope code defect or unresolved functional decision.
Native menu/iframe integration evidence is source/served-API/Node-VM execution;
browser tooling remains unavailable, so this does not establish a rendered
iframe screenshot or visual layout. Binary execution was amd64 through
binfmt/Rosetta on an ARM64 Docker VM, not N150/iStoreOS. No production/paid-provider
request, production settings write, resource pruning or actual business/history
deletion occurred.

Main may proceed to source-backed spec synchronization and the reviewed commit.
