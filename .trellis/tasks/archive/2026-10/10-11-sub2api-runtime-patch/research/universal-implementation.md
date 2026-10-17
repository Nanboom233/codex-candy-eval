# Universal implementation — final handoff

## Post-review package identity

The [independent review](universal-check-review.md) supersedes the installer/archive identity in the handoff below.
The current archive has 62 files, 514,343 bytes and SHA-256 ff096a9de4a6a7d6b69dd6f76eb107e34aab31c65549e53aa35711bed0f9c8b1.
It uses the new v0.4.0 plugin tag and corrected adjacent-bundle discovery.
Functional identity f15183c55e894cd9dd1f114d3f5aabb1b364e8a33de59a24e8b2eefdbf2c6561 and native runtime sources remain unchanged.
The earlier v0.3.0 tag already exists for historical CPA delivery. It is not the universal release target.

## Current outcome

The universal source installer is implemented. The final archive is:
`dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`.
It has 62 regular files and 517,634 bytes. SHA-256:
`14249f68b995247a2d0a5ddde451db9a7a97020955f97e9108d1451082c34f61`.
Latest TOOL_FACTS: release_tag v0.3.0 and functional plugin SHA-256
`f15183c55e894cd9dd1f114d3f5aabb1b364e8a33de59a24e8b2eefdbf2c6561`.
Exact allowlist, source-identical bytes, actual installer extraction and
functional-hash admission passed. No fixed-version host binary, wrappers,
manifest, checksum sidecar, Markdown, tests or development artifacts are in it.

Both real repeated-source-build Linux arm64 lifecycle paths exited 0 with all
15 assertions: dedicated Buildx and standard Docker CLI without Buildx.
Raw terminal evidence lives in codexeval-fixture-535416f93f/evidence.json and
codexeval-fixture-92a37c4d1f/evidence.json under the work root. These full runs
use the prior 5af4052c... functional identity. Both produced the same actual
source-built binary SHA 9204ba16bd4d062e6e7dc1caf370d72a3d80a88946032b729fa2ecfbd1395d05.
They are not claimed as a fresh 15-case rerun of the later packaging change.

The final change aggregates ModelTrace's license only under LICENSE/. The builder
copies that exact legal input into the native host embed path. The text is
byte-identical to the preserved runtime license, and now participates in the
functional hash. A real final-source build/install/loaded-process check passed
with the 62-file package and latest hash. Native source remains current release
f2669c8cf62555cd92389b3f55920e9e6e7c6ff2. Final arm64 binary SHA is
7c82d4071d035680f6b92be9f36c1c4b3fa27e273ee71fc72ea3067ffa47c9d0.
Raw results: final-source-build.log, final-source-build-facts.json,
final-source-install.log and final-source-active-status.log.

The final runtime browser check passed four tabs, math, native light/dark colors,
preferences and retained synthetic history, executable pelican sandbox with
blocked origin storage, 11 authenticated same-origin API calls and zero CSP errors.
The real web Activate button succeeded. The current installer then verified
patched/active state. Logs: browser-final-source.log, browser-activate-final.log.
No inference request was sent. The native authenticated administrator menu passed: same-origin iframe at
/codex-eval/ui, only lang/theme query keys, no token/user/source leakage
(browser-native-menu-assertions.log).
The CLI VM disallows filesystem imports, so attempted direct browser login
helpers could not read fixture passwords. Native login had already succeeded
through the fixture auth API. Protected token/user storage loads that real
session for menu verification without printing credentials.

Required broad gates remain failed and block release publication. The complete
per-profile gate matrix appears below. No claim of universal future-fork support,
native amd64 lifecycle, live paid billing or production concurrency is made.
The same structural rules/signatures are proven on actual v0.2.14 and v0.2.15.
Native amd64 embed cross-compilation passed.

All new raw inputs, workspaces, final binaries, browser fixture and logs remain
for the independent checker. Both source-build lifecycle fixtures were removed
by their harness. The retained browser fixture remains isolated on loopback.
Cleanup follows checker review and the deletion hook. The old historical
minimal-python archive remains unchanged. Nothing was deployed, paid, committed,
pushed, published or archived.


## Context loaded

Reloaded the seven current implement.jsonl sources, then PRD, design and
implement.md in that order (124,695 bytes). The user-approved automatic Docker
source-build design supersedes all fixed-target payload/identity restrictions.
Main owns task/spec artifacts. This agent owns product/tests/package/docs and
does not spawn implement/check children. Unrelated changes are preserved.

The working tree still contains the fixed installer and historical production
patches. No earlier universal product implementation exists. The retained
historical archive is not overwritten. Old cleaned logs/workspaces are not
available evidence.

Implementation starts with the two Docker reproductions and executable native
source/adapter feasibility. Production adapters must retain exact selection,
normal billing/shared slots, cancellation, capability isolation and live exact
cleanup. Missing semantic capabilities must fail before replacement.

Fresh task work root: dist/sub2api-build/universal-20261013/.
No production deployment, paid call, commit, publication or archival is allowed.

## Initial implementation checkpoint (historical)

The GitHub releases/latest endpoint returned v0.2.15, published 2026-10-09.
Actual current source: f2669c8cf62555cd92389b3f55920e9e6e7c6ff2.
Second actual source (v0.2.14): 0363b8cdba8cec3e2ba4b2dbd49c4481143fa55d.
Both immutable development worktrees are source-current/source-previous below
the fresh work root. These identities are development inputs, not product gates.

tests/installer/test_universal.py independently reproduces both reported errors
against the old product. repro-before.log records No Sub2API container is
available and Selected container platform evidence unavailable, before mutation.

adapter/main.go implements AST node/statement-sequence matching and narrow edits.
The development extraction prototype derives 73 structural edits from reviewed
native behavior. Functional native additions live in adapter/runtime as .go.in
sources. adapter/rules.json and ui-rules.json are functional edit recipes, not
host-version manifests. No source-version branch selects a recipe.

The same 73 rules match and produce syntactically valid edits on both actual
revisions. adapter-current-initial.log and adapter-previous-initial.log record
matched native symbols. This is syntax evidence only, not semantic acceptance.
Native hook ambiguity failed closed during development; stream/nonstream hooks
now use enclosing if conditions and then/else control-flow scopes, never ordinal
selection. Optional-node edits discard narrower child edits to avoid overlap.

Commands executed:

```sh
python3 tests/installer/test_universal.py
go run ./adapter extract <source-current/backend> <reference-adapted/backend> <rules-initial.json>
go run ./adapter apply <adapter-proof-current/backend> <rules-initial.json>
go run ./adapter apply <adapter-proof-previous/backend> <rules-initial.json>
go build -o ../proof-server ./cmd/server
```

The last command ran independently in each adapter-proof backend with
proof-build.log. Initial embed-tag compilation correctly required missing
frontend build outputs; this no-embed preliminary check only validates native
types. Production will build the matching frontend before embed compilation.
Both no-embed native server compilations passed. The product Go type engine
imports the selected compiler's export graph and verifies native receiver
methods, including gateway Forward/recheck, runtime Execute/activate/deactivate,
router binding and atomic-expired Key creation. types-current.log and
types-previous.log pass. No universal final package acceptance is claimed yet.

Functional runtime readiness now compares schema-3 protected per-install facts
with compiled source/version/plugin/image/data-path facts, instead of fixed host
pins. Functional host adapter additions are .go.in templates. Development rule
extraction lives in adapter/export_dev.go under the adapter_dev build tag and
is excluded from the explicit release allowlist. The product engine contains
AST adaptation and native type checking only.

Separate semantic-test-current/previous workspaces add maintained test overlays
without changing production sources. Native focused race tests passed service,
middleware, DTO, repository and capability packages on both revisions. The first
composed handler run failed because its fixture still wrote schema-2 fixed host
facts. Updated that maintained test fixture to schema-3 dynamic facts; preserved
all behavior assertions. The rerun of TestEvalComposed passes on both revisions
(semantic-composed-final.log). It exercises four real-router evaluation flows,
normal native billing attribution/shared slots, external isolation and cancel.

Installer implementation now enumerates all Docker containers, resolves current
process Go metadata/source/version/layout/UID/GID, and invokes an isolated
matching-source Docker build. It keeps existing backup/journal/staged/process/
maintenance/recovery owners with dynamic paths and ownership. Source/build
failure still precedes protected target writes. This path is being tested.

scripts/universal_payload.py defines the exact functional source/license map.
scripts/package-sub2api.py now stages that source package rather than two fixed
server binaries. Installer TOOL_FACTS binds the plugin's own release/functional
source identity independently from discovered host releases. Docker build,
installer regression, lifecycle, final packaging and cleanup remain pending.


## Source-driven build and validation checkpoint

The required fresh current-source full unit-race run exited 1. Its log is
semantic-test-current/full-upstream-race.log. It records 30 race warnings and
103 top-level failures. These are current selected-source results. No complete
pristine A/B run establishes attribution. The required gate remains failed.
No unrelated upstream fixtures were changed.

The type engine now checks concrete pointer receiver identity, parameter/result
named types, interface signatures and critical fields. It checks native Forward,
post-slot directed recheck, RecordUsage, expired accounting creation, Key lookup,
delete/quota/rate repositories, runtime router/activation/maintenance/execute,
shared concurrency, account selection/profit gate and usage identity fields.
Compiler exports also compile generated call sites. The exact signature/field
checks pass for both actual revisions in types-signatures-current/previous.log.
Root adapter tests reject parameter/result/receiver/field drift, missing methods
and missing/ambiguous semantic hooks. UI registration remains an independent
exact text adapter. It is not described as Go AST adaptation.

Actual Docker failures and fixes:

1. First lifecycle: root readlink of /proc/1/exe failed under the container's
   ptrace policy. Discovery now reads it as the observed process UID.
2. Second lifecycle: corepack prepare lacked the pnpm shim. Builder enables pnpm.
3. Third lifecycle: enabling all corepack shims hit a dangling yarn symlink.
   Builder enables only pnpm.
4. Fourth lifecycle: native Vite had successfully written backend/internal/web/dist,
   but the builder tried to copy missing frontend/dist. Builder accepts the native
   embed output, optionally copies conventional frontend/dist and requires a real
   embedded index.html. No host-version branch chooses the output path.
5. Fifth lifecycle: matching-source frontend/native signature/embedded Linux arm64
   build and installation succeeded. Activation, repeat install and expiry-drift
   rejection passed. Recreation/reapplication and remaining cleanup checks are
   still running. Full acceptance is not claimed yet.

Raw logs: docker-universal-first/second/third/fourth/fifth.log and fixture-local
builder-*.log below universal-20261013. Fourth concrete failure is preserved in
codexeval-fixture-f7437f115d/builder-16d425b6.log.

The builder freezes independent Go/Node tool image digests and dependency locks.
Buildx uses a dedicated disposable builder, removed after success or failure.
No-buildx environments automatically use disposable standard Docker CLI
containers. That path uses the same functional builder, images and source profile.
It requires no user Go/Node/pnpm/Buildx installation. The dedicated no-buildx
lifecycle is running separately in docker-universal-no-buildx.log.

Installer checks include dynamic paths/UID/GID, both architectures, DockerHub,
GHCR/moving/local/renamed images, missing optional descriptors, wrong application,
source identity conflict, HTTPS URL credential/query/fragment rejection, simple
mode, failed-build cleanup and no target writes before a successful source build.
Functional source hashes are checked before building. A modified source payload
fails closed without a redundant checksum sidecar.

The repository check command passed (repository-final-check.log). It includes
helper tests, installer tests, shell syntax, root go vet, root Go race and UI
checks. Additional installation tests are recorded in installer-final.log.
README/operations/notices now describe the generic source package and unpublished
v0.3.0 entry. A local archive avoids only the release download. A fresh source
build still needs network access, disk and memory.

All new raw evidence and workspaces remain available for the independent checker.
Final cleanup waits for that review. The retained historical archive is unchanged.

Native GOOS/GOARCH build metadata now cross-checks the selected image platform.
The observed official image has vcs.modified=true from its release build. That
flag alone cannot prove divergent source. It is not an admission gate. Exact
revision/version/module plus selected-image binary parity remain authoritative.
A first macOS tar attempt included AppleDouble entries and failed the exact
allowlist check. It is not the final archive. Python tarfile produced a clean
63-file source-identical archive without platform metadata. Final identity will
be recorded after remaining installer validation.
The previous-source full race run also failed: 35 warnings and 126 top-level
failures. Both source-profile full-race gates remain visible and required.
Current-source native vet and integration passed.

## Independent workspace and broad-check findings

The preliminary hand-materialized semantic workspaces copied native Go product
bytes but omitted the UI text adaptation. The first full frontend runs caught
that omission in our private-URL test. Both independent workspaces now receive
the exact product UI recipe and functional frontend module. The maintained
materialization owner, scripts/universal_dev.py prepare, already invokes the
complete product adapt path. It now freezes production hashes before overlays
and rejects any overlay change to those bytes. production-byte-parity.json in
each semantic workspace verifies 3,264 current / 3,214 previous native product
files plus UI rules against the selected source. No test assertion was relaxed.

After correcting materialization, current frontend lint/typecheck pass. The full
run has 363 test files and 2,797 assertions passing, but exits 1 due to two
post-teardown requestAnimationFrame errors in KeyUsageView. Previous frontend
lint/typecheck pass, while its full test run exits 1 with six failures across
two files and one unhandled table-load error. These required gates stay failed.
Focused private-URL test outputs are frontend-private-url-final.log.

Current native integration exits 0. Previous integration exits 1 with:
TestAuthIdentityFoundationSchemas (native entc export reader rejects Go export
data v5 versus supported v4), TestOpenAIWSConnPool_AcquireQueueWaitMetrics
(observed queue 20.45ms below asserted 50ms), and
TestOpenCodeGoUsageRefreshSingleflightAndRunnerDeduplicateSharedGroup
(native 30-second manual refresh throttle returned 429). These test/owner paths
are not touched by adapter rules. A focused rerun of the two service failures
passes (native-integration-failed-cases-rerun.log). That rerun does not turn the
full integration gate green, and no pristine A/B attribution is claimed.
Neither the OpenCode quota refresh nor WS pool queue-metric fixture is the
evaluation billing/slot port. Maintained composed native routing/billing/shared
slot assertions passed on both revisions. No native fixture was modified.

## Earlier 63-file artifact checkpoint (superseded)

TOOL_FACTS release_tag is v0.3.0. Functional plugin SHA-256 is
5af4052ca6b3ff86959f45c58300dce66b9a0f03824447bdcad031d739a99c23.
The latest source-identical release staging is universal-20261013/release-final/.
Its clean Python-tarfile archive is
universal-20261013/sub2api-codex-eval-universal.tar.gz, 518,283 bytes,
63 regular files, SHA-256
20a511aa4ecc12b68a5ea01b5f1787b33b1e6f49c39976adebd8dd63aa7ded63.
Exact allowlist and every source/license byte were checked. Earlier package-*
folders and tar attempts are development evidence, not final deliverables.
The historical minimal-python archive remains at its original path and hash.

The latest repository check passed (repository-final-check-2.log). Installer
coverage now has 56 tests, including functional source alteration rejection
and standard Docker builder failure cleanup. Both revisions pass the focused
private-URL test (four assertions each). The current product-adapted backend
also cross-compiles with the actual built native frontend under embed for
Linux amd64 (adapter-proof-current/proof-embed-amd64.log). file identifies
proof-server-linux-amd64 as a statically linked ELF x86-64 binary. This is
architecture build evidence, not native amd64 loaded-process lifecycle proof.

At that earlier checkpoint buildx completed all 15 lifecycle assertions, including
candidate failure verified restoration, and is compiling the final reinstalled
fixture before terminal evidence.json emission. Standard Docker no-buildx has
completed 13 assertions including exact cleanup and original process restore,
and continues its uninitialized uninstall / candidate recovery cases. Both
run actual repeated matching-source builds, without a test binary fallback.

The maintained prepare owner was also actually executed against each selected
source using release-final. Both commands passed, including its new exact
production-byte check (materialization-owner-proof*.log). These fresh owners
include all Go AST and native UI adaptation before overlays.
Repository-final-check-2.log passes after all installer changes. No further
product source change is planned while the final lifecycle builds finish.

Buildx lifecycle has now exited 0 with all 15 checkpoints and final evidence.json
at codexeval-fixture-535416f93f/. Actual patched arm64 binary SHA is
9204ba16bd4d062e6e7dc1caf370d72a3d80a88946032b729fa2ecfbd1395d05.
Its fixture containers/network/volume were automatically removed. All task-owned
dedicated Buildx builders were removed after each build. Standard Docker
no-buildx subsequently exited 0 with all 15 cases (see current outcome).
A copy of that exact verified source-build binary and its non-secret installation
facts is retained under browser-proof-input for separate browser verification.
The browser fixture reuses this exact compiled candidate after checking its SHA,
source/module/platform/layout/UID/GID/image/plugin facts. This is development-only
reuse, not an installer fixed-binary fallback. Actual source-build lifecycle proof
comes from the independent full Buildx and no-buildx runs.

## Required gate matrix before independent review

| Gate | Current actual source | Previous actual source |
| --- | --- | --- |
| Same AST recipe / parse | pass, 73 edits | pass, 73 edits |
| Receiver/signature/field/compiler graph | pass | pass |
| Native preliminary compile | pass | pass |
| Native focused service/middleware/DTO/repository/evalroute race | pass | pass |
| Composed real-router four-flow/billing/shared-slot/cancel assertions | pass | pass |
| Native vet | pass | pass |
| Full upstream unit race | fail, 30 warnings / 103 failures | fail, 35 warnings / 126 failures |
| Native integration | pass | fail, three tests, two service cases rerun pass |
| Frontend frozen dependencies | pass | pass |
| Frontend lint / types | pass / pass | pass / pass |
| Full frontend test gate | exit 1, 2797 assertions pass, two unhandled errors | exit 1, six failures / one unhandled error |
| Private URL assertions | pass, four | pass, four |
| Actual maintained materialization and byte parity | pass | pass |

Root repository/helpers/installer/UI checks pass. Native amd64 embed build passes.
Actual native arm64 Buildx lifecycle passes all 15 assertions. Standard Docker
no-buildx terminal evidence remains pending at this checkpoint. No broad gate
was weakened, removed or replaced with a focused rerun. No native upstream
fixture repair, live deployment, paid request, commit, push, publication or
archive operation occurred. Publishing remains blocked by the required gates.

## Final license/package and runtime verification

The authoritative final archive is the 62-file path in Current outcome above.
Earlier 63-file/hash checkpoints are superseded. Final release staging is
release-license-final/. The only standalone ModelTrace license is now
LICENSE/ModelTrace-MIT.txt. The builder copies it into the native embed path
with identical bytes. Both Python functional-hash implementations include that
legal embed input. Maintained tests reject modifications to either functional
Go sources or this embedded license before building.

The final-source matching build used immutable Go/Node images, unchanged native
locks, native Vite embed output, complete receiver/signature/field/compiler
checks, and Linux arm64 embed compilation. final-source-install.log exited 0
after replacing and verifying the loaded process. Browser-final-source.log
passes the same controlled browser assertions on the latest runtime hash.
Browser Activate shows 评估已激活, and final-source-active-status.log exits 0
with patched/active state. No model request or paid verification occurred.

Last affected helper checks: helper-license-final.log passes 72 tests with one
historical pinned-checkout skip, installer-license-final.log passes 56 tests,
including source/license integrity and no-buildx cleanup. git diff --check passes.
Root adapter Go changes were already formatted/vetted/race-tested. The license
change does not alter native runtime or frontend behavior, so unrelated broad
gates were not rerun. Their required failures remain in the gate matrix.

Both full repeated-source-build lifecycle harnesses now have terminal exit 0,
15 checkpoints, identical prior-hash candidate SHA and automatic fixture cleanup.
The final-license build separately proves the latest package/build/runtime.
The browser fixture is codexeval-fixture-6ffca1c4e7 on loopback only. Its private
fixture credentials/storage stay at mode 600 and must not enter reports/payloads.
All task-owned generated workspaces/cache/fixtures wait for checker inspection
before authorized cleanup. Do not reuse old missing historical raw inputs.

Native administrator session/menu assertions pass on the final runtime.
The iframe carries only lang/theme, with no token/user/source URL fields.
The real host session refreshes user data through native auth APIs.
Docker buildx ls shows no task-owned builders left. The only retained task-owned
container set is the browser fixture for checker review.
Final report is below 32 KiB. Main owns spec/task updates and independent review.
