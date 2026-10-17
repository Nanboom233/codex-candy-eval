# U5–U6 native browser/model repair

Status: implementation, final frontend/native build/package, actual browser
verification and final owned rollback/teardown passed. Ready for independent review.
Main owns actual Playwright browser and independent review; implementer owns
native patch/module/build/docs and disposable fixture. No commit or nested agents.

Original input:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/pristine-sub2api-v0.2.16.ifVBya/extracted/sub2api-01842576b3d5c1e2545b5a99ad555d36e28ed306`.
Reusable workspace/evidence root:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8`.
Plan: direct native administrator sidebar/icon and iframe using browser origin
including actual port, remove generated settings menu; rebuild changed Vue with
existing dependencies then static native binary/package; repair remaining
evaluation-only configuration coupling and verify browser click→upstream model.
No CSP loosening, baseline/race/full suites, paid or production requests.


## Concrete traces and repairs

U5: the retired settings-response menu used native deriveGeminiRedirectURI.
When X-Forwarded-Host contains only a hostname, that native helper overwrites
Request.Host including its external port; it does not restore a forwarded port.
The generated absolute iframe URL then carries that lost authority. Native
buildEmbeddedUrl previously used new URL(baseUrl) and retained the bad absolute
URL; parent frame-src self correctly rejects a different port/origin. We did not
observe the user's private proxy headers, so this source trace is not a claim
about their exact header configuration.

The fix removes server-generated menu/origin guessing. Native AppSidebar now
has an administrator-only beaker icon/translated label beside account management,
and native router mounts /admin/codex-eval with requiresAuth/requiresAdmin. New
native/frontend/src/views/admin/CodexEvalView.vue embeds /codex-eval through the
existing native browser URL helper. That helper now resolves relative inputs
against full window.location.origin, preserving protocol/hostname/port. Existing
absolute third-party custom URLs remain unchanged. Parent CSP was not modified.
The retired codex_eval_menu.go was preserved under previous-delivery/u1-menu;
settings GetSettings is restored to its original expression.

U6: page refreshCatalog used normal /accounts/:id/models. Its native
GetAvailableModels calls FetchOpenAIAccountModels, which takes raw discovery
and applies projectAccountModelsBody: mapped aliases appear, unconfigured names
are deliberately removed. Editing the sender alone left this display/admission
coupling in front of it. New administrator GET
/accounts/:id/codex-eval/models calls CodexEvalModels, which reuses raw native
FetchOpenAIModelsList without projection, mapping or local allowlist expansion.
The page calls this evaluation-specific path and retains editable model fields.
Existing focused test now contrasts raw evaluation names with unchanged normal
projected picker names under the same conflicting map/cache.

Start → codexeval.Start → original runner/codexTurn → executeModel →
SendCodexEvaluation still passes the selected model literally and shares exact
account transport/protocol conversion. doOpenAIAccountTestUpstream itself does
not map models. Ordinary create/update capability probes independently select a
configured mapped model; they remain ordinary native behavior and are distinct
from stream:true evaluation calls. No normal /test or gateway behavior changed.

## Build and fixture status

Final build log: workroot/native-browser-build.log; session38201 exited0.
It applies the ordinary backend+frontend patch and copies native/backend and
native/frontend. Original frozen dependencies are reused through a temporary
workspace symlink, then pnpm9.15.9 build runs the native i18n check, vue-tsc and
Vite. New CodexEvalView chunk and changed native app/sidebar are rebuilt before
embed. Relevant Go tests (including raw-vs-projected model test), touched-package
vet and static Linux/amd64 embed build passed. No stale native dist reused.
Release workflow installs frozen dependencies then invokes this single builder.
Package session97049 exited0 at fixed dist/build outputs.

Final-artifact disposable fixture: workroot/native-browser-fixture.py,
session50589, evidence workroot/native-browser-fixture/. It passed real installer,
admin auth, rebuilt native route, direct sender and Candy sanity checks, then
held on its continue control file for main's actual browser. Verified origin
http://127.0.0.1:55182; group2/account1; model browser-unconfigured is returned by
raw upstream catalog while configured mapping points it to wrong-remap. Upstream
logs record only path/model/stream. Login and mutable endpoint details are in
browser-fixture.json using disposable-only credentials; no JWT is logged.

Main completed real Chrome verification through the working Playwright CLI.
Actual native login → sidebar 账户评估 → /admin/codex-eval succeeded. Sidebar
SVG has the native beaker path; iframe contentDocument loaded with parent/frame
origins both http://127.0.0.1:55182, path /codex-eval, port55182 retained. Browser
console had zero errors/warnings, including no frame-src violation. Actual raw
suggestions were gpt-5.5/browser-unconfigured; normal authenticated catalog was
empty while the evaluation catalog exposed both. Main filled the real frame's
model combobox, clicked Test, and observed completed Candy history showing
browser-unconfigured(low), correct1/1. Upstream latest stream:true CC call had
model browser-unconfigured; independent baseline mapped probe was stream:false.
Screenshot output/playwright/sub2api-native-port.png was viewed by main.
Main's browser evidence owner is research/native-browser-verification.md.

Final rollback/teardown: session50589 exited0 with `Installer rollback check
passed.` The malformed executable caused the selected disposable container to
exit; the installer restored its actual previous executable from the stopped
container and verified the restored application. Evidence outcome.json reports
passed, including failed replacement restoration. installer-failure.stderr
confirms restoration and records exit133, OOMKilled=False. Exact owned-container
query `docker ps -a -q --filter label=codex-eval-check=w8ax6xp8` returned empty
after harness teardown. No unrelated containers or business data were removed.
Prior Node-VM checks were partial contract evidence and did not satisfy browser
acceptance. Original pristine/race/full baseline is not repeated; actual N150
execution remains unperformed.


Fixed artifact SHA256:
- dist/build/sub2api-codex-eval:
  5a570c1cde16b88492c81fbcef771077f45e70e9065e933a985a94a8c34cd173.
- dist/build/sub2api-codex-eval.tar.gz:
  c69d522f7c52badd3b65b5e50c76c306fe806b308491c55fc6c1ad473f23e716.
Archive inventory remains binary, installer and five license texts only.
