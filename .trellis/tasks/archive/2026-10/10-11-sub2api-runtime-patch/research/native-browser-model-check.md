# Independent U5–U6 browser/model check

Status: full-scope review completed against frozen final implementation and
actual Chrome evidence. Main owns spec/task updates and commit; reviewer owns
this report and safe local check fixes. Latest U5–U6 override earlier zero-Vue
and generated settings-menu contracts.

## Findings (fixed)

`scripts/test-sub2api.sh` still told developers to supply a source tree whose
frontend was already built. The builder now rebuilds the patched native Vue
frontend and requires its frozen-lockfile dependencies. Updated only that usage
text and checked shell syntax. This script is not packaged; the final
browser-tested binary/archive remain unchanged.

Earlier independent review did not establish parent CSP/browser integration or
the configured catalog projection ahead of sending. Its no-defect conclusion
was incomplete for those paths. This review uses actual rendered Chrome and the
complete discovery → controls → start → sender chain.

## Code review

- Native sidebar has one authenticated administrator entry beside account
  management, using the existing `Icon` component's beaker icon and translated
  navigation label. Native router requires authentication and administrator role
  for `/admin/codex-eval`; the new `CodexEvalView.vue` uses ordinary `AppLayout`
  and an iframe, with theme observation cleaned up on unmount.
- The iframe input is relative `/codex-eval`. Existing `buildEmbeddedUrl` now
  resolves relative URLs against complete `window.location.origin`, preserving
  nonstandard port and protocol. Existing absolute inputs still retain their
  absolute origin. Native URL helper keeps theme/language/embedded parameters;
  evaluation auth reads the same-origin native `localStorage.auth_token` rather
  than a URL token or additional login.
- Parent CSP remains native and unchanged. Exact standalone page bypass plus
  its page-only SAMEORIGIN/frame-ancestors/inline-renderer policy remain. The
  same full parent/frame origin satisfies parent `frame-src 'self'` and child
  framing policy, without relaxing global CSP to permit a wrong port.
- Superseded `codex_eval_menu.go` and settings-response patch are removed from
  current product source. Unrelated custom settings remain native; no automatic
  settings persistence or menu compatibility framework remains.
- New native admin `GET /accounts/:id/codex-eval/models` uses existing injected
  `AccountTestService`; it validates the same exact OpenAI target scope and calls
  native `FetchOpenAIModelsList`, then returns its standard raw upstream names.
  Original `FetchOpenAIAccountModels` calls that same raw discovery followed by
  `projectAccountModelsBody`, the config-dependent mapping/allowlist step. The
  evaluation endpoint omits that projection; normal `/models`, `/test`, account
  create/update probes and ordinary gateway keep their original behavior.
- Page catalog calls the evaluation-specific models endpoint and keeps editable
  inputs/datalist suggestions. Start has no model-catalog/config admission gate;
  original runners create the request and call the shared native exact-account
  sender directly. Sender retains literal selected model before native protocol
  conversion, and native plugin/TLS transport does not remap it. Credentials,
  shadow account resolution, proxy, URL validation, headers, snapshots and
  401/429 processing retain the previously reviewed behavior.
- Build entry applies backend/frontend ordinary patch, copies new view/module,
  reuses installed frozen dependencies and actually builds the changed native
  frontend before Go embed compilation. Release workflow installs dependencies
  once, then invokes the same build/check entry. It no longer relies on original
  stale embedded assets. The fixed package remains prebuilt binary/installer/
  licenses; no source build is performed by the user's installer.
- Reused prior independent scope checks for four original algorithms/resources,
  options/retries/results, memory tasks/busy/cancellation, four-map history,
  schedulable/status separation and installer reset/address/rollback behavior.
  No unrelated billing/gateway/Key/recovery/compatibility/cleanup feature entered
  these repairs.

## Actual verification

Evidence inspected: `native-browser-model-repair.md`, main's
`native-browser-verification.md`, final `workroot/native-browser-build.log`,
native tests/current source and screenshot
`output/playwright/sub2api-native-port.png` (also viewed independently by reviewer).

- Final native `pnpm build` ran the existing locale checks, Vue TypeScript check
  and Vite; log includes newly built `CodexEvalView` chunk. Relevant Go tests,
  touched-package vet and static Linux/amd64 embed build passed.
- Existing sender test checks unmapped literal model for both Responses and CC.
  The added focused catalog test contrasts raw evaluation names and unchanged
  normal projected names from the same conflicting configuration and cache.
- Main drove actual native login and clicked the rendered sidebar. Chrome at
  `http://127.0.0.1:55182` loaded `/admin/codex-eval`; child origin was the same
  complete origin/port55182, with readable loaded contentDocument. Native beaker
  SVG was present. Console recorded zero errors/warnings after loading and run.
- Authenticated normal catalog returned `[]`; evaluation catalog/suggestions
  returned `gpt-5.5` and `browser-unconfigured` despite conflicting stored map.
  Main filled the real iframe combobox, clicked its account Test button and
  observed completed Candy `browser-unconfigured(low)`, correct1/1. Local
  upstream log recorded `/v1/chat/completions`, model `browser-unconfigured`,
  `stream:true`. Ordinary mapped capability probe remained a separate
  `stream:false` request.
- Screenshot shows actual native header/sidebar/selected beaker entry, loaded
  evaluation layout, model input, account controls and successful result/history.
  This is rendered browser evidence, superseding the earlier Node-VM-only claim
  for these affected paths.
- Reviewer checked shell/Python syntax, whitespace and the current installer in
  the final archive. No binary-affecting review change, so no redundant rebuild
  or repeat of the passing browser run was needed.

Workroot:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8`.
Final outputs remain `dist/build/sub2api-codex-eval` and its `.tar.gz` archive.
Implementer owns final exact-fixture rollback/teardown record.

## Findings (not fixed) and limits

No remaining in-scope code defect or unresolved functional decision identified.
Actual Chrome was exercised on the development machine at a nonstandard local
HTTP port; N150/iStoreOS and the user's production HTTPS/proxy configuration were
not executed. Supplier was local fake upstream, so execution/model forwarding was
verified without claiming real model quality or spending provider quota.
No production write, resource prune, real history deletion or commit occurred.
No race/full-suite/coverage/parity requirement was added.

Main may synchronize current specs/task and prepare the updated commit review.
