# Independent fresh-native check

Status: independent review and local fixes completed. Reviewer owns local fixes in `native/`, the current
installer/build/package/check entries and delivery documentation. Main owns
task/spec updates and commit. Existing unrelated edits are preserved.

Current approved task and F01–F06/S01–S05/C02–C22 override historical gates.
Completed baseline, focused checks and disposable runtime evidence in
`fresh-native-implementation.md` are reused; no baseline/race/full-suite repeat.

## Findings fixed

- Active `scripts/test-sub2api.sh` was already narrowed by the previous checker
  to the current `build-sub2api.py --source ... --check` contract. Retired universal
  helper and root-suite requirements are no longer called.
- Chat Completions errors and premature EOF reached a Go reader error while the
  direct handler ignored `io.Copy` errors, losing the actual SSE failure reason.
  `native/backend/internal/service/account_codex_eval_service.go` now emits an
  actual error event with native observed usage. Upstream `error:null` is ignored.
  `native/backend/internal/handler/admin/account_codex_eval_handler.go` reports a
  non-cancellation copy failure as an error event. Existing native conversion
  functions remain unchanged.
- Pelican invalid/missing HTML only corrected status in list summaries; persisted
  detail still reported completed. `native/backend/internal/codexeval/pelican.go`
  now records failed at evaluation time; summaries only omit HTML.
- Existing stop text promised in-flight completion despite context cancellation.
  `native/backend/internal/codexeval/web/components.js` now explains interruption
  and worker completion. Its unused CPA filter placeholder was removed.
- `native/backend/internal/codexeval/web/theme.css` retained CPA host themes beneath
  a partial native override. It now contains one native light/dark token mapping
  using Sub2API's teal primary and slate dark colors, original spacing/components
  and subscription badge styling. No native Vue edits or layout redesign.
- `native/backend/internal/codexeval/web/app.js` had a stale CPA login comment;
  removed that inaccurate comment without changing the approved session behavior.
- `install.py` Ctrl-C after replacement bypassed restoration; the replacement
  exception branch now also handles KeyboardInterrupt before propagating it.

## Review result

Reviewed the new ordinary patch, native sender/handler, four original collectors,
history/task synchronization, standalone page/direct APIs and delivery entries
against approved F/S/C behavior and actual original callers. The patch remains
limited to native router registration, administrator routes and exact page bypass.
The ordinary `/test`, gateway, native conversion library and Vue source stay
outside this patch. Original algorithm differences are package namespace only;
original request identity/effort/context adaptations match approved decisions.

No unresolved public-interface or functional design change was identified.
Historical adapter/universal/migrated/test directories remain unreferenced by
the active build/check/package path; they were not deleted or used as inputs.

## Verification

- `gofmt` applied to changed Go files; Python AST, shell syntax, standalone
  JavaScript syntax, CSS parse using existing upstream PostCSS, and diff whitespace
  checks passed.
- Current `build-sub2api.py --check`: relevant evaluation/service/admin tests,
  touched-package vet and static Linux/amd64 compile passed. Existing runtime
  test now checks missing/oversize Pelican HTML. Native sender tests check CC error
  and premature EOF with actual error text and observed usage, without fabricated
  success. Log: original implementation workroot `check-rebuild.log`.
- Subsequent CSS-only token cleanup was followed by a cached final static build
  using the completed pristine source/frontend/module cache; no frontend build
  or repeated Go tests. Log: same workroot `check-final-build.log`.
- Focused installer interruption simulation passed: after replacement raises
  KeyboardInterrupt, the exact backup is restored, readiness is checked and the
  interrupt propagates. No user container was operated for this check.
- Regenerated `dist/build/sub2api-codex-eval` and
  `dist/build/sub2api-codex-eval.tar.gz`. `file` reports stripped statically linked
  x86-64 ELF. Archive has only executable, current installer and five license
  texts; no cleanup/development/test payload.
- Reused completed disposable Docker install/auth/exact-account/four-task/history/
  enable-PUT and malformed-ELF stopped-container restore evidence in
  `fresh-native-implementation.md`; no redundant repeat of that complete fixture.

## Findings not fixed / limits

No outstanding in-scope code defect. Browser tooling was unavailable in the
implementation session; Node VM evidence validates API/JavaScript behavior, not
visual CSS rendering. The narrowed native token mapping was checked against
upstream colors and parsed, but was not visually inspected in a browser. Actual
N150/iStoreOS execution and production/paid-provider calls were not performed.
These limits remain explicit; no extra acceptance gates were added.

Main owns spec synchronization, commit and Trellis finish after this handoff.
