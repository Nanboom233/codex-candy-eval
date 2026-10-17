# Consolidated Sub2API Check Review

Date: 2026-10-12
Scope: frozen patches 001-006, copied evaluation core, build tooling, and the
current Sub2API replay. Product implementation, deployment, paid requests,
task completion, and archival were outside this check.

## Verdict

Local patch, core, backend integration, frontend, and static build checks pass
except for the full unit-race suite. The migration is not ready for acceptance:
the backend has no evaluation API or production coordinator wiring, the
frontend has no evaluation page, the installer is absent, and the public
guidance remains CPA-specific.

The maintained workspace lacked its patch marker. Checks used the isolated
replay at dist/sub2api-build/audit-consolidated-20261012/replay-final/workspace.
All six frozen patches and the 28-file core matched their identities. Backend,
frontend, and core bytes matched the maintained source before generated build
outputs; all three lockfiles remain byte-identical to the frozen baseline.
The final core tree hash is
b8477eb71da52e832aefa690b7bdf4fa02c257a792f1669ff3cbb5126ec67af5.

## Findings (Fixed)

- scripts/sub2api_build.py: patch-path validation previously trusted Git
  diff headers and accepted whitespace errors. It now validates raw unified
  headers and applies patches with --whitespace=error-all.
- scripts/test_sub2api_build.py: added regressions for raw lockfile and
  traversal headers, quoted-octal paths, and whitespace rejection. Added
  isolated wrapper regressions for temporary Wire checksums and nested
  Corepack pnpm scripts.
- scripts/generate-wire.sh and scripts/build-sub2api.sh: Wire generation
  now uses a temporary modfile pair. This resolves the pinned
  github.com/google/subcommands v1.2.0 sum without changing frozen source
  lockfiles. The generated wire_gen.go remains byte-identical.
- scripts/pnpm, scripts/pnpm-support.sh, scripts/build-sub2api.sh, and
  scripts/test-sub2api.sh: the Corepack-only path now exposes the pinned
  pnpm wrapper to package scripts that invoke pnpm recursively.
- scripts/build-sub2api.sh: release builds disable automatic Go VCS metadata
  so binaries built beneath the CPA checkout do not identify themselves with
  the parent CPA revision.

Wire generator checksums and lockfile hashes are recorded in
dist/sub2api-build/audit-consolidated-20261012/wire-tool-dependencies.txt.
Both final binaries and source provenance are under
dist/sub2api-build/audit-consolidated-20261012/out/.

## Audit Artifacts

`out/provenance.json` is the build identity record. It binds the binaries to
Sub2API v0.2.15, source commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`,
the pinned image digests, Go/pnpm toolchain, six patch hashes, evaluation-core
tree hash, and three lockfile hashes. `out/SHA256SUMS` records both artifact
checksums. The Wire generator's temporary dependency checksums and retained
modfile pair are documented in `wire-tool-dependencies.txt`.

Saved command output is in
`dist/sub2api-build/audit-consolidated-20261012/`:

- `root-check-final.log` records the 65 helper tests, root vet/race, shell
  syntax, and UI harness.
- `build-check-final.log` records manifest, source identity, and toolchain
  checks.
- `backend-vet.log`, `backend-default-compile.log`,
  `backend-integration.log`, and `backend-unit-race.log` record replay backend
  results, including the full unit-race failure.
- `frontend-lint.log`, `frontend-typecheck.log`, `frontend-tests.log`, and
  `frontend-build-wrapper.log` record frontend checks.
- `build-linux-amd64-final.log`, `build-linux-arm64-final.log`, and
  `wire-generate-final.log` record final artifact and Wire commands.

## Verification

- Root checks: scripts/test-sub2api.sh check passed; 65 helper tests,
  go vet ./..., go test -race ./..., shell syntax, and UI harness passed.
- Build inputs: scripts/build-sub2api.sh check passed for the pinned source
  and toolchain.
- Backend vet: go vet ./... passed in the patched replay.
- Default-tag compile: go test -run '^$' ./... passed. Tests were compiled
  but not executed.
- Backend integration: go test -tags=integration ./... passed using the
  repository harness's temporary Postgres and Redis containers.
- Frontend lint: pnpm lint:check passed.
- Frontend typecheck: pnpm typecheck passed.
- Frontend tests: 361 files and 2,791 tests passed.
- Frontend production build: passed with the pinned frontend lockfile.
- Linux builds: static ELF x86-64 and arm64 binaries passed with Go 1.27.2,
  CGO_ENABLED=0, and VCS metadata disabled.
- Lockfiles: backend go.mod, backend go.sum, and frontend pnpm lockfile match
  the frozen hashes in out/provenance.json.
- Full backend unit-race suite: failed. It reported 18 race warnings and 58
  failed tests. See backend-unit-race.log; the declined pristine baseline
  was not rerun and unrelated fixtures were not changed.

The final artifact hashes are:

- amd64: 08ea6bf8881f87f567fb0f5fdbc0ee6fca04ef7dd63e066b342046858dc94dbb
- arm64: df217b31fb66820340b77e7d19e0a438871c5718b6b77c2cc72dcd84f3c3d1ac

See out/SHA256SUMS and out/provenance.json for checksums, patch hashes,
image identity, source commit, core tree, and lockfile hashes.

## Integration Boundary

backend/internal/server/router.go:126-134 registers the existing route
modules and no Codex evaluation route. backend/cmd/server/wire.go:34-61
does not include evaluation handlers or a topology runtime provider.
EvalTopologyRuntimeDeps declares the current-facts, evidence-writer,
listener-source, and local-probe ports, and its constructor rejects missing
dependencies. No production implementations or runtime construction were
found in the replay.

The Sub2API frontend/src tree has no evaluation page or API client. Root
README.md, install.sh, install.ps1, and .github/workflows/release.yml remain
CPA-specific. There is no Sub2API inspect/install/status/rollback/uninstall
entry point. These are migration gaps, not local review fixes.

## Requirement Matrix

| ID | Result | Evidence and remaining work |
| --- | --- | --- |
| R1 | PARTIAL | Four evaluation implementations and fixtures exist in the core. No host API or UI invokes them. |
| R2 | PARTIAL | Prompts, baselines, scoring, model handling, and controls remain in the unchanged core. End-to-end host behavior is absent. |
| R3 | PARTIAL | Directed Responses routing uses the normal gateway path. No evaluator producer assembles usage, billing, or shared concurrency. |
| R4 | PARTIAL | Protected token and directed-auth primitives exist. Admin run APIs and concurrent group selection are absent. |
| R5 | PARTIAL | Directed routing reuses target-group checks. Evaluation API authorization and complete integration are absent. |
| R6 | PARTIAL | Exact-account and no-fallback tests exist. Run-manager result handling is not connected to these failures. |
| R7 | PARTIAL | Lookup redaction and replay rejection primitives are tested. Admin evaluation routes, menu, and browser response checks are absent. |
| R8 | PARTIAL | Core persistence and retention tests exist. Host endpoints and refresh/restart UI flow are absent. |
| R9 | NOT IMPLEMENTED | No Sub2API binary replacement installer exists. |
| R10 | NOT IMPLEMENTED | Status, rollback, uninstall, transactional restore, and recreation detection are absent. |
| R11 | PARTIAL | Manifest checks and both static builds pass. No installer refuses incompatible targets before replacement. |
| R12 | FAIL | README, installers, and release workflow still describe or build the CPA plugin. |
| R13 | NOT IMPLEMENTED | No Sub2API evaluation UI or theme adaptation exists. |
| R14 | PARTIAL | Key and topology primitives exist, but runtime facts, listener, probe, readiness, revocation, and reapply wiring are absent. |

## Acceptance Matrix

| ID | Result | Evidence and remaining work |
| --- | --- | --- |
| A1 | PARTIAL | Core and transport fixtures pass; no integrated evaluation route or controlled host gateway flow. |
| A2 | LIVE ONLY | No live candy request, usage record, or charge was attempted. |
| A3 | PARTIAL | Directed scheduling primitives are tested; evaluator and ordinary traffic are not assembled together. |
| A4 | PARTIAL | Token and copied-group primitives exist; concurrent-group admin flow is absent. |
| A5 | PARTIAL | Authentication and authorization primitives are tested; no public evaluation API exercises them end to end. |
| A6 | PARTIAL | Normal gateway policy remains on the directed path; evaluation-level quota cases are not integrated. |
| A7 | PARTIAL | Directed failover tests preserve exact targets; manager history integration is absent. |
| A8 | PARTIAL | Lookup redaction/replay tests pass; admin routes, menu, and browser output are absent. |
| A9 | PARTIAL | State retention/recovery tests pass; host API/UI restart flow is absent. |
| A10 | NOT IMPLEMENTED | No installer verifies image identity or preserves deployment configuration. |
| A11 | NOT IMPLEMENTED | No installer idempotency, rollback restore, or uninstall behavior. |
| A12 | NOT IMPLEMENTED | No runtime-patch recreation inspection path. |
| A13 | PARTIAL | amd64/arm64 static builds pass; target compatibility refusal is unimplemented. |
| A14 | FAIL | Public installation and release guidance remains CPA-specific. |
| A15 | NOT IMPLEMENTED | The Sub2API evaluation UI and theme behavior are absent. |
| A16 | PARTIAL | Proof primitives have tests; real ingress inventory, listener capture, probes, and readiness wiring are absent. |
| A17 | PARTIAL | Owned-Key conflict checks have tests; installer revocation and reapplication are absent. |

## Key Boundary

| ID | Result | Evidence and remaining work |
| --- | --- | --- |
| K1 | PARTIAL | Raw token storage and separation primitives are tested. The assembled API/browser path does not exist. |
| K2 | PARTIAL | Lookup redaction and external replay rejection are tested. Full runtime and browser-output checks are absent. |
| K3 | PARTIAL | Directed gateway routing preserves normal admission checks; no integrated evaluation or live billing proof. |
| K4 | PARTIAL | Collector proof and anti-spoof tests exist. Production facts, listener, ingress, and probe adapters are absent. |
| K5 | NOT IMPLEMENTED | No planned rollback/uninstall path can stop workers and revoke only the owned Key. |
| K6 | ACCEPTED RESIDUAL | No stronger original-runtime lookup-secrecy or automatic-revocation claim is made. |
| K7 | PARTIAL | Owner/status/binding conflict primitives are tested; installer reapply behavior is absent. |

## Unresolved

The full unit-race run contains two handler/admin races outside the frozen
patches. TestGrokImportProbeFailureLogDoesNotIncludeErrorMessage reads a
bytes.Buffer while its asynchronous logger writes it. In
TestGrokOAuthHandlerQueryQuotaProbesUpstream, the test reads a stub map while
observed-model synchronization writes it through UpdateExtra. Other failures
include the documented pristine service fixtures and process-global test
state. These were left unchanged and the suite was not repeated.

No approved Sub2API live target was configured. No deployment, browser login,
paid model request, production charge, or application container lifecycle
test was run. Integration-test Postgres and Redis containers were disposable
and terminated by the test harness. No .env values were read.

The build artifacts establish compilation and provenance only. They do not
prove installed-container compatibility or close the missing runtime,
installer, UI, documentation, or live acceptance work. The task remains
in_progress.
