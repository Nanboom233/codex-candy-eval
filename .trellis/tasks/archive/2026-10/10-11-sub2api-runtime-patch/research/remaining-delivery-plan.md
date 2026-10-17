# Research: Remaining Sub2API Delivery Plan

- Query: Define one cohesive remaining implementation and final-check handoff from the completed consolidated audit.
- Scope: internal
- Date: 2026-10-12

## Findings

### Baseline and boundary

The pinned upstream is Sub2API v0.2.15 at `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. The current 001-006 patch digests and the existing 28-path `internal/eval` tree digest remain historical audit evidence; the patches mix production changes with 13 test files, and the core tree includes seven `_test.go` files plus a golden fixture. Preserve prior source bytes and digests where applicable, but split the new production patch/core inputs from their test overlays. The complete 28-path digest is not the production staged-tree identity. Record new production patch-set, production-core, production-staged-tree, and separate development-only test-overlay identities. Existing component reviews remain in `implement.md` and prior delivery reports. The consolidated audit passed helper tests, root checks, focused backend tests, backend integration, frontend lint/typecheck/2,791 tests/build, Wire generation, and static Linux amd64/arm64 builds. Its full unit-race suite failed with 18 race warnings and 58 test failures. The pristine baseline was not rerun and unrelated fixture repair was declined. Audit binaries/provenance are under `dist/sub2api-build/audit-consolidated-20261012/out/` and must be rebuilt after implementation.

Production evaluation routes/Wire/bootstrap/readiness, real probe adapters, evaluation UI, installer lifecycle, and Sub2API public guidance are still absent. The current-facts provider, evidence writer, listener source, and local probe are required nonnil coordinator ports, not optional enhancements. No live Sub2API target is approved/configured.

Keep R1-R14, A1-A17, and K1-K7 unchanged in the PRD; include R15 and A18 for production/test isolation. Continue to preserve normal billing, usage, subscription/quota/model rules, user/account concurrency, exact OpenAI target eligibility, and no fallback. No sidecar, schema migration, Compose/image change, business rebinding, automatic upgrade, or provider expansion is in scope.

### Contract signatures and boundaries

Patch 006 defines the current runtime API and DTOs. Wire concrete production implementations into these contracts. Repartition 001-006 as required to isolate production code from test overlays, preserving their accepted runtime behavior and security contracts rather than treating historical patch hashes as a freeze.

```go
type EvalDeploymentCurrentFactsProvider interface {
    CurrentEvalDeploymentFacts(context.Context) (EvalDeploymentCurrentFacts, error)
}
type EvalSessionEvidenceWriter interface {
    CheckAvailable(context.Context) error
    Save(context.Context, EvalSessionEvidence) error
}
type EvalBoundListenerSource interface {
    ActualEvalListener() (EvalBoundListener, error)
}
type EvalTopologyLocalProbe interface {
    ProbeLoopback(context.Context, EvalBoundListener, []EvalProbeLocalTicket) error
}
func NewEvalTopologyRuntime(EvalTopologyRuntimeDeps) (*EvalTopologyRuntime, error)
func (*EvalTopologyRuntime) Begin(context.Context, int64, EvalProbeBeginRequest) (EvalProbeBeginResult, error)
func (*EvalTopologyRuntime) Observe(*gin.Context) (bool, error)
func (*EvalTopologyRuntime) Complete(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Abort(context.Context, int64, string) (EvalProbeSessionStatus, error)
func (*EvalTopologyRuntime) Status(context.Context, int64, string) (EvalProbeSessionStatus, error)
```

`EvalDeploymentCurrentFacts` binds version/build type, deployment ID, source revision, runtime patch digest, image index/platform/config digests, actual binary digest, and ingress revision. BuildInfo supplies only version/build type. A trusted runtime/installer provenance provider must independently source the remaining facts; neither inventory JSON nor admin input may supply expected-current values. Any unavailable or mismatched fact keeps readiness false.

The operator independently inspects every public listener, proxy hop, rewrite, and expected peer, then installs `deployment-evidence/inventory.json` under the protected `EvalKeyStore.Dir()` tree. Patch 005 is the strict reader. The collector proves declared paths only; it cannot discover omitted proxies. External topology completeness remains an operator/live acceptance obligation.

Write separate runtime-owned session evidence under that protected tree, bound to the inventory's exact-byte SHA-256, current facts/security fingerprint, session/owner, exact case set, expiry, and SHA-256 of the sealed proof bytes. Use a fixed path, strict versioned schema, private owner/mode, no-symlink validation, atomic replacement, and idempotent Save for identical session/inventory/proof identity; conflicting content fails. `CheckAvailable` must not create files. Store no raw token, ticket, ticket hash, or raw proof in session evidence. On restart, revalidate the exact inventory, current facts, proof, expiry, and owned Key. Missing completion evidence requires a fresh same-owner probe unless a valid distinct durable completion binding was written.

### One coherent implementation sequence

Assign one `trellis-implement` the full remaining delivery. This order is internal dependency order, not separate implementation/review batches.

1. **Provenance and durable evidence.** Add a trusted `EvalDeploymentCurrentFactsProvider` using independent runtime/install provenance and actual binary/image facts. Add the protected session-evidence writer/reader and restart validation beside the existing deployment inventory reader. Load the inventory and facts at Begin, completion, evidence-only retry, and startup readiness. Make failures statusable and fail closed; ordinary Sub2API startup continues.

   Preserve coordinator recovery transitions: Seal failure invalidates ephemeral state and requires a fresh probe. Every `Initialize` error may be partially committed, so publish owner-bound `provisioning_recovery_required`, clear tickets, and allow only same-owner explicit recovery through the existing exact intent. The current `row_pending` intent with no row has no automatic resume path; keep it blocked for explicit operator repair. `completion_pending` is only successful Initialize plus evidence-save failure; retry must revalidate facts/proof/Key and write evidence only, never call Initialize. Abort clears ephemeral state only and preserves persisted intent/Key. Under the coordinator mutex consume a ticket/increment its observation lease, call `collector.Observe` unlocked, then publish only for the same generation; abort waits for leases.

2. **Real probes and ingress boundary.** Change server startup to bind and retain the actual `net.Listener`, expose it through `ActualEvalListener`, and serve that listener. Build local direct-loopback probes from its actual port and verify a self-request reaches this server. Add the one-shot standard-library operator client at `backend/cmd/codex-eval-probe/main.go`: verify local inventory bytes/digest and exact case-to-origin coverage, call authenticated Begin/Complete as the same administrator, keep tickets in memory, and probe public/declared ingress paths from the named external network. Use verified HTTPS, exact `/v1/responses`, redirects and ambient proxies disabled, and send no admin/private/gateway credential to probe origins.

   Add ticket lookup/consumption at the start of the `/v1` chain. Wrap the final composed `http.Server.Handler` in `internal/server/http.go`, after body-limit handling and outside Gin. If the reserved `X-Codex-Eval-Probe-Ticket` marker is absent, pass through unchanged. If present, reject unless method/path are exactly `POST /v1/responses`; do so before Gin redirect/no-route/method handling. A valid ticket is one-use and terminates before API-key lookup, forwarding, account selection, usage, or billing. Unmarked ordinary requests retain current behavior.

3. **Host runtime, API, Wire, and shutdown.** Compose the existing `EvalKeyStore`, `EvalKeyProvisioner`, protected inventory reader, evidence writer, facts provider, listener source, probe, coordinator, target adapter, and run manager through service/handler Wire providers. Bootstrap only as a recovery/guard operation. Keep status/admin APIs and ordinary host traffic available when evaluation is not ready. Construct one retained nonnil lazy executor for `NewRunManager`; each execution resolves readiness and the current delegate, returning typed `ModelError{Status: 503, Code: "not_ready"}` when evidence drifts or is absent. The runtime alone owns raw token and real `HTTPExecutor`; use the literal loopback origin and normal `POST /v1/responses` route.

   Add admin-only operations under `/api/v1/admin/codex-eval/`: `status`, `probes/begin`, `probes/complete`, `probes/abort`, `targets`, `state`, `start`, `cancel`, `records/:type/:account_id/:run_id`, `clear-history`, and `flush`. Derive actor ID from existing admin auth context, never request data. Use `{code,message,data}` and set no-store on probe responses. A Begin response may return external one-use tickets only to the authenticated same-owner Begin/operator flow; keep them in memory in the operator client. Never disclose tickets through Status, history, logs, audits, or URLs. Never expose the raw evaluation token or Key lookup value in any response, URL, log, or audit. `/targets` is a snapshot only; Start re-resolves exact group/account membership. Existing gateway checks remain authoritative.

   Add one service/runtime owner in the current DI graph, one admin handler in `AdminHandlers`/`RegisterAdminRoutes`, and regenerate `cmd/server/wire_gen.go` with `scripts/generate-wire.sh`; do not hand-edit generated wiring. Extend `Application` and `cmd/server/main.go` so manager drain/cancellation precedes HTTP shutdown, while the local gateway and dependencies remain available until workers and terminal saves settle. Only then close the executor and existing dependencies. Evaluation storage/load failures remain inspectable and do not create an empty fallback or make ordinary host startup fatal.

4. **Structure-preserving UI.** Maintain a plain-JavaScript embedded page from the existing evaluation UI, with its assets/licenses owned by the maintained evaluation package and copied by the existing core build path. Keep four tabs, account rows, batch controls, details, result/history layout, Markdown/math renderers, and pelican sandbox. Add group/account selection where the Sub2API API requires it; do not redesign the page. Use the verified light/dark tokens in `ui-sub2api-theme.md` and preserve local `data-theme`.

   Serve the static shell at `/codex-eval/ui`, with per-response host CSP nonce and route-local `X-Frame-Options: SAMEORIGIN`/`frame-ancestors 'self'`; leave global security headers unchanged and do not use `unsafe-inline`. Patch only the exact same-origin `frontend/src/utils/embedded-url.ts` menu exception. Pass only `theme`/optional `lang`, never JWT or source-page parameters. The iframe reads the current same-origin `auth_token` for each API request, sends it only in headers, and stops on expiry. Preserve unrelated custom-menu settings. Accept native theme-toggle reload and restore server state/history/preferences.

5. **Installer and recovery transaction.** Replace CPA installer/public entry behavior with Sub2API `inspect`, `install`, `status`, `rollback`, and `uninstall`, retaining only the supported host targets. Before mutation, compare actual release/source commit, image index/platform/config digest, architecture, binary path, startup settings, writable replacement support, and checksums to `compatibility/sub2api-v0.2.15.json`; reject mismatches without upgrade or replacement. Build/verify the matching static binary, back up `/app/sub2api` and checksum, persist an install journal, stage beside the target, and replace atomically without overwriting the running inode. Restart only the selected disposable/target container; preserve official image, Compose, volumes, and existing business bindings.

   Make status inspect actual binary/image identity so container recreation is detected before reapplication. Reinstall is idempotent and reuses the exact owned Key/menu; it never rebinds or reenables externally changed rows. Planned rollback/uninstall drains evaluations, revokes only the verified owned Key via `RevokeOwnedKey`, removes only the owned menu/provisioning resources, restores the original binary only under its verified image/checksum contract, and retains history and unrelated Keys. Keep the transaction journal and report incomplete if revocation/restore fails. Do not claim automatic revocation after unplanned recreation (accepted K6).

6. **Build/test isolation, release, and public guidance.** Keep the three upstream lockfiles unchanged and retain the pinned build guards: raw/Git patch path validation, strict whitespace application, temporary Wire modfile/checksum, pinned pnpm shim for nested scripts, and `-buildvcs=false`. Replace CPA README/install/release steps and ensure the release workflow never builds or deploys the old `.so` plugin.

   Split the mixed 001-006 material into production-only patches and a separate test overlay; keep 007 runtime and 008 UI patches production-only, with their tests in the overlay. Stage production source from an explicit core/resource allowlist that excludes `_test.go`, `testdata`, fixtures, mocks, audits, logs, and development-validation payload. Preserve the runtime Codex context, fingerprint probes/baselines, ModelTrace bank/license, embedded UI, normal gateway/auth/billing behavior, and required installer/security checks. Apply test overlays only to a separate disposable workspace derived from the same pinned source and production inputs. Restrict overlays to test files/data, retaining same-package access for Go tests by placing `_test.go` at its package path only inside that workspace. Compare production path/hash manifests between production and test workspaces before and after overlay application. Run patched product Go/frontend integration tests there; run repository helper and mocked UI harness checks only from the development checkout or CI. Never include test/validation payload or output in the final release or deployed server runtime. The one-shot operator evidence client and installer compatibility/rollback checks remain product operations; keep the client on the operator host, outside the server container. Produce static Linux amd64/arm64 server artifacts, the operator evidence client artifact, production provenance/checksums/license notices.

   Give the production patch set, test overlay, production core source, complete production staged tree, and test-augmented tree distinct hashes. Preserve historical 001-006 and 28-path digests as evidence. The complete staged-tree hash must be recalculated when UI/product assets are added. Store test-overlay and test-workspace identities only in development evidence, not production provenance or the release package. Add an allowlist packaging check outside product code that creates the release directory from approved server/probe binaries, installers, compatibility manifest, production provenance/checksums, public operations documentation, and license notices, and rejects any unexpected test/dev path. Build production binaries from a fresh production-only workspace, not the test-augmented workspace.

### Owners, producers, and consumers

| Data or surface | Producer/owner | Consumer or guard |
| --- | --- | --- |
| Protected deployment inventory | Deployment operator/installer after independent topology inspection | Existing patch-005 reader; coordinator reloads by artifact ID and exact-byte digest |
| Current deployment facts | New trusted runtime/installer provenance provider | Inventory reader, coordinator, bootstrap, readiness and status |
| Session evidence | Coordinator after successful probe and Key initialization; new protected atomic writer | Startup readiness/recovery and status; evidence-only retry after strict revalidation |
| Actual listener identity | `cmd/server/main.go` retained `net.Listener`, surfaced through server/runtime provider | Local probe origin and lazy executor |
| Probe observation/tickets | Coordinator plus pre-Gin guard and gateway-chain observation hook | Collector; never general gateway forwarding/accounting |
| Evaluation target metadata | Existing exact `EvalTargetAdapter` / `TargetResolver` | Admin `/targets` and manager admission; Start resolves again |
| Evaluation run state/results | Existing `codexeval.RunManager` and `FileStateStore` | Admin state/record/history APIs and UI |
| Page/menu/theme/auth | Maintained embedded evaluation assets plus Sub2API frontend exact-URL patch | Same-origin admin iframe and authenticated APIs |
| Production/test source split | Production patch/core allowlists and test-only overlays | Production build uses only production inputs; disposable checks apply overlays separately |
| Binary transaction | Sub2API installer and journal | `inspect/status/rollback/uninstall`, actual binary/image checks and owned-Key lifecycle |

### Expected patch and artifact set

- Production patches under `patches/sub2api/v0.2.15/production/`, including production-only 007 runtime and 008 UI changes. Repartition historical 001-006 changes into production patches and test overlays; the split outputs are ordered build inputs, not separate review gates.
- Test-only Go/frontend patch overlays under `patches/sub2api/v0.2.15/test-overlay/` and maintained core tests/fixtures under a separate tree such as `tests/overlays/codexeval/`. Apply these only to a disposable test workspace; preserve same-package Go test paths there.
- Updated `compatibility/sub2api-v0.2.15.json` and build helpers with separate production and test-overlay identities, an explicit production core/resource allowlist, workspace separation, and production-path byte comparisons. Keep `scripts/test-sub2api.sh`, `scripts/test-ui.cjs`, helper unit tests, and disposable Docker checks as development-only tools.
- Sub2API root install entry points, transaction/journal implementation, README/installation/operations docs, license notices, and Sub2API release workflow. Keep Docker/build outputs under ignored `dist/`. Release packaging uses an external explicit allowlist and does not contain test overlays, fixture/mocks, audit/logs, test reports, caches, or development validation scripts.
- Final `dist/sub2api-build/out/` static amd64/arm64 server binaries, operator evidence client binary, installers if distributed, compatibility manifest, production `provenance.json`, license notices, and `SHA256SUMS`. Keep separate `productionPatchSetSha256`, `productionCoreSourceTreeSha256`, and new complete `productionStagedTreeSha256`; preserve historical 001-006 and 28-path digests only as evidence. Record `testOverlayPatchSetSha256` and `testAugmentedTreeSha256` only in development check evidence. Compare production file/hash manifests across both workspaces. The operator client and installer preflight are product security operations, not developer test payloads, and the probe runs on the operator host rather than in the deployed server container.

These production patches and test overlays are ordered inputs, not separate implementation or review gates. `implement.md` and the role context should carry the final production hashes, development-only overlay identity, and owned file lists.

### Validation and acceptance status

The implementer prepares one cohesive delivery for one full-scope `trellis-check`; production/test source isolation is an input boundary, not a separate implementation or approval batch. Patched product Go/frontend integration tests run only in disposable test workspaces. Repository-owned build-helper/UI harness checks may run from the development checkout or CI. No test/validation payload or output ships or executes in the deployed server runtime. The checker should verify production-only build inputs and artifact allowlisting, compare production path/hash manifests with the test-augmented workspace, re-run locked source/replay/build gates against the production patch series, backend vet/default compile, integration tests, focused new unit-race tests, frontend lint/typecheck/tests/build, Wire regeneration checks, both static artifacts/provenance, and disposable Docker install/failure-restore/rollback/uninstall/recreation tests. Test the probe guard through a real bound listener, not only `gin.Engine.ServeHTTP`. Check normal routing/billing/concurrency and ensure invalid probes create no auth lookup, upstream call, usage, or charge.

Retain the full backend unit-race gate and report its failures; do not skip or suppress it. The consolidated audit's 18 warnings/58 failures are the known pre-integration result. Fix races/failures introduced by this delivery, but do not repair unrelated upstream fixtures that the user declined to change. The completed helper/backend-integration/frontend/build checks are useful baseline evidence and must not be described as migration acceptance.

| Criteria | Local delivery evidence | Still external/live |
| --- | --- | --- |
| R1-R8, A1, A3-A9, K1-K3 | Four evaluation flows through admin API and normal gateway; exact target/no fallback; persistence/restart; auth/redaction; concurrency, quota/billing fixtures. | A2 actual candy usage/charge and production-account concurrency remain live-only. |
| R14, A16-A17, K4-K7 | Protected facts/inventory/evidence, spoof/ingress fixtures, same-actor session/recovery, Key conflict/revoke, installer reapply and container recreation tests. | Operator must independently attest the real proxy/ingress inventory; K6 remains accepted residual. |
| R13, A15 and R7/R8, A8/A9 | Existing structure/theme/API/menu tests plus local browser verification, nonce and no-token URL/header checks. | Authenticated browser flow on a real approved deployment remains live-only. |
| R9-R11, A10-A13 | Disposable container identity, atomic install/restore, idempotence, status/recreation, amd64/arm64 build and provenance. | No production image/container change without an explicitly approved target. |
| R12, A14 | README/install/release guidance and built artifact describe Sub2API; old CPA deployment is not invoked. | No external requirement beyond published artifact review. |
| R15, A18 | Production patch/core/resource and release allowlists; test overlays restricted to separate disposable workspaces; independent pre/post production path/hash parity and artifact inspection prove no test or development-validation content enters production artifacts or staged runtime. | Inspect the deployed payload against the runtime allowlist when an approved target is available; no live target is currently configured. |

After the single check, main updates source-backed specs and Phase 3.4 records per `.trellis/workflow.md`. Keep the task `in_progress` and do not commit/archive/finish while a required local gate fails or live-only acceptance remains unverified. With no approved live target, report the remaining live criteria as blocked rather than claiming completion.

### Files found and code patterns

- `internal/eval/runs.go:98-138`: `NewRunManager` preserves an inspectable blocked manager when load or recovery persistence fails, and turns unfinished persisted runs into interrupted records without executor calls. `:211-240` validates exact positive group/account targets before target resolution.
- `internal/eval/state.go:40-85,88-156`: one-owner versioned state storage rejects symlink/nonregular paths and replaces snapshots through a private same-directory temporary file. The runtime session-evidence writer is a separate owner and needs its own schema and lifecycle checks.
- `internal/eval/http.go:27-64`: `HTTPExecutor` requires a literal loopback HTTP origin, owns a transport with `Proxy: nil`, refuses redirects, and targets `/v1/responses`. It cannot establish which listener is actually bound; the host must provide that fact separately.
- `internal/eval/runs.go:612-688`: `Flush`, cancellation, and `Shutdown` retain manager ownership of work and wait for tracked workers. The host must call manager shutdown before HTTP shutdown so those workers can still use the local gateway.
- `scripts/sub2api_build.py:263-307,513-537`: patch path validation includes raw unified headers; provenance binds upstream/image/toolchain, patch and lockfile hashes, core tree hash, and artifacts. `scripts/build-sub2api.sh:104-224` sequences source validation, patching, Wire regeneration, locked frontend build, static amd64/arm64 builds, provenance, and checksums.
- `scripts/sub2api_build.py:319-334,472-486`: current `stage_eval_core` copies the entire core directory, including tests and `testdata`, into the host workspace. The isolation plan replaces this with explicit production allowlisting and a separately materialized test workspace.
- `patches/sub2api/v0.2.15/001-codex-eval-auth.patch:2,200,475,1244,2152,3551` and `002-codex-eval-scheduling.patch:26,1527`: examples of production patches adding same-package `_test.go` files; all six current patches contain test additions.
- `.github/workflows/release.yml:11-24,63-99`: current workflow tests the CPA plugin, packages old shared libraries, and uploads a wildcard asset directory. Replace with production-only Sub2API artifact allowlisting; CI tests stay outside release assets.
- `scripts/test-sub2api.sh:56-87`: integration tests currently run in the same prepared build workspace. Use a second disposable workspace with test overlays so final binaries are built from production-only inputs.
- `internal/eval/modeltrace_test.go:13-18` and `internal/eval/testdata/modeltrace-golden.json`: test-only golden fixture; contrast with production embedded datasets and license at `codex.go:36-46`, `fingerprint_algorithm.go:38-45`, and `modeltrace_algorithm.go:16-20`.
- `dist/sub2api-build/workspace/backend/cmd/server/main.go:171-194` (audited v0.2.15 workspace): startup currently calls `ListenAndServe`, and shutdown occurs before the existing process cleanup. `dist/sub2api-build/workspace/backend/internal/server/http.go:112-163` composes the Gin handler, body limit, and final `http.Server.Handler`; the probe marker guard belongs outside that final chain while preserving the body limit.
- `dist/sub2api-build/workspace/backend/internal/handler/wire.go:180-296` and `dist/sub2api-build/workspace/backend/internal/server/routes/admin.go:13-33` (audited v0.2.15 workspace): handlers are supplied through Wire and admin routes inherit the shared auth, panel limiting, audit, and compliance middleware. The new API should join these owners and use generated `wire_gen.go`.
- `dist/sub2api-build/workspace/frontend/src/utils/embedded-url.ts:1-46` (audited v0.2.15 workspace): the shared iframe URL builder currently adds user ID, JWT, theme, language, and source-page parameters. The migration must narrowly exclude those credentials/source values only for the exact same-origin evaluation URL; broader menu behavior is shared.
- `internal/plugin/plugin.go:31-103` and `internal/plugin/web/ui.html:26-40,50-87` show the maintained embedded page's asset assembly and four-tab/account-list structure. `research/ui-sub2api-theme.md:46-56,141-175` records the Sub2API token map and focused UI/browser assertions.
- `research/host-runtime-refinement-response.md:9-20` records the corrected coordinator error states, exact-intent recovery limit, evidence-only retry, and observation-lease rule. In particular, an existing `row_pending` intent with no Key row has no automatic resume path.

### Related research/specs and external source

- Current audit: `research/consolidated-check-review.md` (full local result and R/A/K matrices).
- Test/release isolation: `research/test-release-isolation.md` (mixed inputs, production allowlist, overlay workspaces, identities, and package boundary).
- Runtime: `research/host-runtime-wiring-plan.md`, `research/host-runtime-refinement-response.md`, `research/host-runtime-plan-refined-review.md`, `research/host-coordinator-delivery.md`.
- UI and lifecycle: `research/ui-sub2api-theme.md`, `research/private-key-lifecycle.md`, `design.md`, and `research/patch-contracts.md`.
- Specs: `.trellis/spec/backend/{evaluation-core,run-state,host-target-adapter,deployment-evidence,host-coordinator,build-guards,quality-guidelines}.md`.
- External source pin: Wei-Shaw/Sub2API v0.2.15 at commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`; pinned toolchain Go 1.27.2 and pnpm 9.15.9. No external documentation was fetched.
- The `dist/sub2api-build/workspace/` references above are ignored local audit/build material, not maintained source. Their version-specific anchors are corroborated by `research/source-audit.md` and the pinned commit; regenerate the workspace before relying on them if it is absent or stale.
- The current build helper copies test sources into its staged tree and the current release workflow is still CPA-specific. Production/test workspace separation, explicit staged-tree hashes, and release allowlisting remain implementation work. Development checks remain required for engineering review but are never release or deployed-runtime content.

## Caveats / Not Found

- No live Sub2API target, real ingress inventory, authenticated browser session, or paid evaluation is available. Keep these acceptance items pending until the target is provided and approved.
- The runtime cannot infer omitted reverse proxies. Code and disposable fixtures do not substitute for operator topology attestation.
- Current facts must have an independent provenance source for each field. If actual image/deployment/ingress facts cannot be sourced, report unavailable and keep readiness false.
- The current provisioner cannot automatically resume a `row_pending` intent when its Key row is missing. Surface an explicit repair-required status; never create a replacement, rebind, or reenable silently.
- Full unit-race failures remain a known quality blocker; the user declined unrelated upstream fixture repairs.
- This continuation only completed and source-anchored the delivery plan; it did not change product code or rerun tests/builds.
