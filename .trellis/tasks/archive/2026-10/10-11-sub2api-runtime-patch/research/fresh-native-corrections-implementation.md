# U1–U4 implementation progress

Status: implementation, affected build/checks, fixed package, actual native
runtime/API, readiness reproduction and stopped-container rollback completed.
Main owns task/spec/
commit; implementation retains independent reviewer fixes and existing baseline.

Workspace root: `/Users/minami/Code/C/cpa-plugin-codex-candy-eval`.
Original source: `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/pristine-sub2api-v0.2.16.ifVBya/extracted/sub2api-01842576b3d5c1e2545b5a99ad555d36e28ed306`.
Previous materialized source/evidence root:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8`.
Reusable disposable harness: `native-smoke-final-complete.py` and
`native-ui-runtime-check.cjs` under that root; no fixture currently remains.
Prior independent findings: `research/fresh-native-check.md`.

U1: automatically supply one native administrator settings custom-menu item
for same-origin /codex-eval and existing sidebar/iframe. No Vue/database writes.
U2: native schedulable flag controls disabled/switch; separate account status
and temporary restrictions.
U3: administrator-selected model bypasses evaluation-only mapping/config;
keep native credentials/protocol/transport.
U4: readiness inspects selected container state/address on retries, reports
bounded elapsed/terminal failure and preserves rollback.

## Changes and root causes

- U1: settings GetSettings originally returned only stored menus, so no sidebar
  entry was delivered. `codex_eval_menu.go` appends one admin `账户评估` using
  the existing native request-origin helper and same-origin absolute URL. It
  preserves other entries, replaces the canonical/stale duplicate and sorts
  after administrator custom entries. `native/sub2api.patch` adds one settings
  response expression to its existing three original-file changes. Native Vue,
  custom route and iframe implementation remain unchanged; no DB menu write.
- U2: schedulable false originally became generic unavailable while the switch
  controlled status. Standalone state now derives disabled from schedulable;
  the switch uses native POST `/accounts/:id/schedulable`. It shows `已停用`
  with switch off and distinguishes inactive/error/quota status, future rate
  reset, overload and temporary unschedulable fields. Switch changes do not
  activate errored/inactive accounts or erase runtime restrictions.
- U3: shared sender no longer applies account GetMappedModel or OAuth model
  normalization; selected model reaches upstream directly for either protocol.
  Existing model options become datalist suggestions for editable model inputs
  on all four tabs; empty/config-derived catalogs cannot prohibit a direct name.
  Existing effort controls and native credential/proxy/TLS/header/plugin/bridge
  handling are retained. Existing focused sender test now asserts requested
  model despite a conflicting configured mapping for Responses and CC.
- U4: readiness takes selected container ID, inspects current state/address
  each attempt, and reports elapsed/deadline, state, restart count, address and
  request path/error. Restarting retries remain bounded; exited/dead/paused
  states terminate with exit/OOM/runtime error diagnostics. Failed installs and
  Ctrl-C retain reviewer-fixed restoration. Stopped rollback requires no exec.

Changed delivery documentation: README and docs/sub2api-operations. All previous
reviewer streaming/Pelican/theme/interruption fixes were retained.

## Build and artifacts

- `workroot/corrections-build.log`: final build `--check`, relevant evaluation/
  service/admin tests, touched-package vet and actual embed/static amd64 build
  exited 0 (session 2816). Original frontend and caches reused.
- `patch -p1 --dry-run --batch` passed on original after final patch text cleanup;
  cleanup removes only patch-file whitespace context, generated Go unchanged.
- Final package (session 55770) exited 0; inventory exactly binary, installer
  and five license texts. Python/changed JavaScript syntax and diff check pass.
- Binary SHA256 `dc2f1ca622930197754379758f7a17715677604dedb75b20e658ff60a82c2f38`.
- Archive SHA256 `48dc870fb0515999b53605eeff2aa732bab6550e4ba90b47f5c27a24f59b0c92`.

## Actual disposable runtime evidence

Final binary harness `workroot/native-smoke-corrections.py`, session 33874,
exited 0. Evidence `workroot/smoke-corrections/`:

- `outcome.json`: installer replacement, same-origin page/CSP, native admin auth,
  automatic menu settings, existing `/custom/codex-eval` native SPA response,
  direct sender and background Candy with deliberately conflicting model map.
- `ui-corrections-checks.json`: actual served standalone JS against real HTTP
  APIs in existing Node VM fake DOM; all four tab handlers, row/renderers,
  schedulable POST/off/disabled and separate restriction facts, direct editable
  model even when mapped native catalog returned no suggestions.
- Same VM ran unchanged original AppSidebar admin filtering/append fragment and
  original TypeScript buildEmbeddedUrl helper against actual settings response:
  final item `/custom/codex-eval`, absolute same-origin iframe URL, native dark/
  language/embedded parameters and reachable evaluation shell.
- Fixture seeded one stale canonical and two unrelated menu entries. Actual
  response preserves unrelated admin/user entries and has one rewritten current
  canonical entry. DB query after reads retained original seed unchanged,
  proving response generation did not persist the generated menu.
- Local upstream logs show evaluation model `gpt-5.5` and
  `direct-unconfigured` unchanged. Original background capability probe used its
  configured mapped model, as expected; that ordinary native path is unchanged.
- `native-address-change.json`, `native-restart-readiness.log`: actual native
  restart changed published port 49153→49283; readiness followed new address.
- `installer-failure.stderr`: malformed candidate exits133; stopped state was
  detected after0.1s, backup restored, restarted readiness passed and evaluation
  page remained reachable. Native fixture resources removed in finally.

## Readiness reproduction and limits

Controlled Docker HTTP listener calls resetAndDestroy. Docker Desktop converts
that reset to an actual EOF/RemoteDisconnected observed by host Python. Initial
prepare harness had an overly strict final assertion expecting literal reset
text after its native restart and Docker bounded-timeout/restart checks had
already passed. The corrected Docker-only harness accepts the actual proxy
behavior; session10748 exited0. No product edit was needed for this
fixture distinction.

`corrections-real-reset.py` additionally exercised an actual host TCP SO_LINGER
RST with the selected Docker container's real inspected state and a controlled
local proxy address. `real-tcp-reset.log` contains actual macOS `[Errno 54]
Connection reset by peer`; readiness retried and completed health/setup. This
is real transport error evidence, not an injected exception or errno104 Linux
claim. Docker timeout stayed bounded at2.0s; restart recovery followed a changed
published port (`reset-*.log`/`reset-restart-result.json`).

CUA inventory remains unavailable (Sky native pipe startup failed, no browsers).
Source/served-JavaScript/API checks do not claim visual browser layout or an
actual rendered iframe screenshot. Native binary ran via amd64 binfmt/Rosetta
in ARM64 Docker VM, not N150. No production writes, real providers, paid calls,
pruning, business/history deletion, commit or sub-agent spawning occurred.
Final Docker-only run recovered across published port49999→50034 and bounded
failure at2.0s. Its owned container was removed. Final exact-label `docker ps
-a -q --filter label=codex-eval-check=w8ax6xp8` returned empty. No owned fixture
resources remain, and unrelated resources were preserved.
Next: main's independent Trellis check and spec/task/commit coordination.
