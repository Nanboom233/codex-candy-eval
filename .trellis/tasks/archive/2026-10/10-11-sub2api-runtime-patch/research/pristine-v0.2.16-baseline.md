# Pristine Sub2API v0.2.16 Baseline

## Outcome and boundaries

The untouched official source compiled successfully into a static Linux/amd64 server with the real embedded frontend.
The newly compiled binary loaded through emulation and passed real setup, frontend and configured-health checks.
Frontend required checks, full frontend tests, native integration, vet, upstream linter, shell and release-helper checks passed.
An isolated Linux ARM64 unit run passed after correcting DNS and temporary-write prerequisites.
Native macOS full unit runs retained DNS failures. Supplemental native full race failed and has two unmodified minimal reproductions.
This is a planning baseline. It does not accept the plugin package or activate product implementation.

Every original command record remains under the workroot `commands/` and `runs/` directories.
Complete, separate original stdout and stderr remain under `logs/`. No failure output was replaced with a summary.
The structured [outcomes](pristine-v0.2.16-baseline-outcomes.json) contain exact argv, environments, timestamps, exits, durations, counts and failure/skip names.
The durable [checkpoint](pristine-baseline-checkpoint.json) records all prerequisite attempts and final preservation checks.

## Source, tools and preservation

Workroot:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/pristine-sub2api-v0.2.16.ifVBya`

Official annotated tag object: `c3cb8588ce6a5aef4dc29fa4f71ac91f3c0bdf3d`.
Exact commit: `01842576b3d5c1e2545b5a99ad555d36e28ed306`.
Official archive SHA256: `eef7e0519b05b8254c981efbab7bfe74c5d89897c67053e81380987366daae89`.
The researcher reverified the official tag and inspected the [upstream contracts](pristine-v0.2.16-upstream-contracts.md).

All 4,322 original archive files matched byte-for-byte before and after commands.
Both the original extracted tree and separate Linux scratch tree preserve those files.
Upstream-generated frontend dist, dependencies, tsbuildinfo and temporary test outputs are separate additions.
No plugin, adapter, core, overlay, installer, assertion or fixture modification entered either tree.
`source-files-before.sha256`, `lockfiles-before.sha256` and `purity-final.json` retain the identities.
`go mod verify` passed. Dependency locks and source VERSION remain unchanged.

| Frozen file | SHA256 |
| --- | --- |
| backend/go.mod | 809e6806dc248bd44cea4198ddc21fd6ed04f80003300c3668367b91bae319a8 |
| backend/go.sum | b17668c6300c04027bd397650e60e28e514dc4abb9ab113e3cea634fdd708109 |
| frontend/package.json | d1fd3d91ccf9bbca6cf0a79428343109a16efa6b3675fd01869d37db5d63083e |
| frontend/pnpm-lock.yaml | 1d8e54cc8e1f1bd3e7593805f607d9314ff81dab98ab5a9a6396475e148b26ce |

Native host: macOS27.0.1 ARM64, Go1.27.2 darwin/arm64.
Frontend: official Node20.20.2 darwin/arm64 and explicitly pinned pnpm9.15.9.
Node tarball SHA256 matches official SHASUMS: `466e05f3477c20dfb723054dfebffe55bc74660ee77f612166fca121dacb65b6`.
Upstream golangci-lint2.14.0 binary reports its official build tool Go1.27.0.
Release helper tests use isolated Python3.12 and required PyYAML6.0.3.
Linux diagnostic: official `golang:1.27.2-alpine`, Linux/arm64, default CGO_ENABLED=0.
Image ID: `sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d`.
These environments differ from upstream Ubuntu CI runners. No complete Ubuntu CI claim is made.

All 165 maintained product snapshot files remain unchanged.
Canonical plugin archive SHA256 remains `39c876d274a6499479bf0fbcbcf8d93c46c52053b01b2dbbc6ee3b497554f6fc`.
Canonical plugin binary SHA256 remains `eed01f71c39048eef1a15fbf13652b5c44d9915fff1e34c59f3d58abf1c6c9e1`.
The pure baseline artifact was never copied into `dist/build`.

## Actual checks

Log filenames below are relative to the workroot. Each stem has `.stdout` and `.stderr` files.
Backend Makefile commands keep the official tags and no race flag for required unit/integration checks.
GOFLAGS adds readonly module mode and verbose telemetry. It does not change assertions or selected test tags.

| Check and actual command | Environment | Exit / duration | Result / log stem |
| --- | --- | --- | --- |
| `make test-unit` | native Mac, initial download proxy | 2 /302.419s | 3 top-level CAPTCHA failures; `backend-unit` |
| `make test-unit` | native Mac, proxy removed | 2 /208.221s | 1 DNS-dependent matrix failure; `backend-unit-clean-env` |
| `make test-unit` | native Mac, subsequent DNS recovery attempt | 2 /202.569s | 1 different DNS-dependent matrix subtest; `backend-unit-dns-recovered` |
| `make test-integration` | native Mac, proxy removed, CI=true, actual Docker | 0 /143.773s | 7,127 top-level passes,15 top-level skips; `backend-integration-clean-env` |
| `go test -tags=unit ./...` | Linux ARM64, read-only source, first DNS prerequisite | 1 /215.836s | `.entc` write prerequisite and second hostname DNS failure; `linux-unit-full` |
| `go test -tags=unit ./...` | Linux ARM64, byte-identical writable scratch, two actual DNS hosts | 0 /195.427s | 11,806 top-level passes,18 skips; `linux-unit-writable-all-dns` |
| `make test-frontend` | Node20/pnpm9 native Mac | 0 /67.713s | lint, typecheck,573 tests/30 critical files; `frontend-required` |
| `pnpm run test:run` | Node20/pnpm9 native Mac | 0 /99.706s | 2,869 tests/372 files; `frontend-full` |
| `pnpm run build` | Node20/pnpm9 native Mac | 0 /43.439s | locale test, vue-tsc and embedded Vite output; `frontend-build` |
| `go vet ./...` | native Mac | 0 /8.337s | passed; `backend-vet-clean-env` |
| `golangci-lint run --timeout=30m` | v2.14.0 native Mac | 0 /43.353s | 0 issues; `backend-lint-clean-env` |
| `go test -tags=unit -race -json ./...` | native Mac, proxy removed | 1 /412.096s | 140 top-level failures,28 race warnings; `backend-race-clean-env` |
| release matrix unittest discovery | Python3.12 native Mac | 0 /1.701s | 10 passed; `release-helper-tests` |
| all eight applicable shell commands | native Mac | all0 | `upstream-shell-1` through `upstream-shell-8` |
| `go mod verify` | isolated verified module cache | 0 /1.784s | all modules verified; `go-mod-verify` |

The eight shell commands are exact upstream deployment syntax/tests, Compose simple-mode test and release shell syntax check.
They appear individually in structured outcomes. These are helper tests, not production deployment.

Linux final unit emitted 22,907 passing test/subtest events and18 skip events.
Its 56 cached package results replay earlier passing package runs. Failed packages reran after prerequisite corrections.
Integration emitted14,061 passing test/subtest events and16 skip events, including one nested skip.
Counts distinguish top-level tests from subtests and package events.
Raw skip reasons and exact names remain in `runs/*-test-counts.json` and structured outcomes.
Skipped coverage includes explicitly enabled live/network checks, prompt-audit environment tests and upstream optional plugin-runtime tests.
No skip count was presented as passing coverage.

## Prerequisite failures and unmodified reproductions

Initial system resolution failed for Node, Go proxy, npm and GitHub.
Official Node curl exited6. macOS socket lookups, router DNS and public DNS queries timed out.
A disposable Docker DNS probe also timed out. Exact logs and failed command records remain.
A TLS-verified direct connection fetched official Node20. A workroot-only CONNECT proxy used authenticated DNS-over-HTTPS for dependency downloads.
Destination TLS validation remained enabled. Host DNS configuration and hosts files were not modified.
An inherited Corepack shim initially selected pnpm12 and frozen install failed.
Explicit pnpm9.15.9 preparation and a direct pinned executable corrected that prerequisite without changing package files or locks.
The final frozen install passed. All earlier prerequisite outputs remain retained.

Initial backend unit and integration runs contained exactly these three repository failures:

- `TestAliyunCaptchaVerifier_VerifySuccess`: expected JSON success, observed `invalid character '<' looking for beginning of value`.
- `TestAliyunCaptchaVerifier_VerifyResultFalse`: same JSON error.
- `TestAliyunCaptchaVerifier_APIErrorNormalized`: expected `AliyunCaptchaAPIError`, observed `json.SyntaxError`.

The unchanged four-case focused test exited1 with the download proxy and0 after clearing only proxy environment variables.
Proxy logs recorded the HTTP POST501 responses. Local CAPTCHA fixtures remained unmodified.
HTTPS_PROXY-only plus localhost NO_PROXY also passed all four cases.
Original initial failure records remain distinct from the later clean integration pass.

Native clean unit failed at `TestValidateCreateParams_CheckModeMatrix/quota_probe_requires_primary_model`.
It expected `CHANNEL_MONITOR_MISSING_PRIMARY_MODEL` but observed `CHANNEL_MONITOR_ENDPOINT_UNREACHABLE`.
The original validator resolves `api.kimi.com` before checking model presence.
The unchanged minimal test initially failed, then passed with HTTPS-only and without any proxy.
A system recheck showed this hostname had begun resolving while other hostnames still failed.
Those results do not prove the HTTPS proxy repaired Go's direct DNS resolver.
The subsequent full native unit run instead failed at `probe_requires_api_key` while resolving `api.openai.com`.
Its expected `CHANNEL_MONITOR_MISSING_API_KEY` was preceded by the same unresolved-endpoint error.
No repeated blanket native reruns were used to hide these results.

An isolated Linux environment used current TLS DoH-derived A records for both exact hostnames via Docker `--add-host`.
It had no external network. Local httptest traffic remained available.
The first read-only source mount also prevented ent's unchanged `TestAuthIdentityFoundationSchemas` from creating its `.entc` temporary workspace.
A separate byte-identical writable scratch copy corrected that environmental requirement.
The combined unchanged matrix/schema minimal check passed, followed by the full Linux unit pass.
Every original source byte in both copies remained identical afterward.
Linux Alpine CGO0/ARM64 and actual DNS-host entries are explicit differences from Ubuntu CI and native Mac runs.

## Supplemental race result

The pristine native full race run exited1 with28 `WARNING: DATA RACE` blocks.
It recorded140 top-level failing tests,149 failures including subtests, and two failed packages:
`internal/handler/admin` and `internal/service`.
Exact names, complete race stacks and parsed blocks remain in `runs/backend-race-test-counts.json` and `race-blocks.json`.
Many fail events can follow concurrent races. Their count is not a count of independent root causes.

Unmodified focused race reproductions also failed:

- Admin: `TestGrokImportProbeFailureLogDoesNotIncludeErrorMessage` and `TestGrokOAuthHandlerQueryQuotaProbesUpstream`, exit1,7 race warnings.
- Service: `TestRunCheck_QuotaProbeAttachesSnapshotToPrimaryRowOnly`, exit1,4 race warnings.

The admin trace includes the test's `bytes.Buffer.Bytes` read at `grok_import_probe_test.go:277` concurrent with scheduler logging at `grok_import_probe.go:144`.
The service trace records concurrent writes to the original capture handler at `channel_monitor_checker_body_test.go:72`.
These are measured pristine-source traces. No test synchronization, race suppression or assertion changes were applied.
This evidence does not individually explain the previous patched153 failures/51 warnings.
Required upstream unit CI omits `-race`. The supplemental race result remains failed and requires review before broader integration conclusions.

## Real compile and loaded runtime

Actual direct release-equivalent command used Go1.27.2, `GOOS=linux GOARCH=amd64 CGO_ENABLED=0`, `-tags=embed`, `-trimpath`, and upstream release ldflags.
The exact command is in `commands/linux-amd64-build.json`; stdout/stderr and24.013s exit0 are retained.
Ldflags set Version0.2.16, the exact source commit, build date and BuildType=release.
The original VERSION file says0.2.15 and remained unchanged.
Upstream release tooling normally changes that file, so the direct compiler version override is documented rather than editing pristine evidence.
This direct compiler output is not a GoReleaser release or published archive.

Pure artifact:
`artifacts/sub2api-pristine-v0.2.16-linux-amd64`

Size:123,682,976 bytes.
SHA256: `9bc90c95a5d8811e1e93aff55175ab0be7d999740fd4e5cb632d48d586a5a2bf`.
`file` reports a stripped, statically linked x86-64 ELF.
The new artifact is distinct from the canonical plugin binary.

A Docker `-version` invocation loaded that exact ELF and reported0.2.16/exactcommit.
Execution used binfmt/Rosetta through an ARM64 Alpine base on the ARM64 Docker VM.
The initial forced amd64-base command exited125 because the cached Alpine image was ARM64. That original output remains.
No native N150/iStoreOS execution is claimed.

The same pure binary passed five actual HTTP startup checks:
setup `/setup/status` needs_setup=true, embedded frontend root, configured `/health` status=ok,
configured `/setup/status` needs_setup=false, and configured embedded frontend root.
Normal startup used fresh owned PostgreSQL18.1-alpine3.23 and Redis8.4-alpine on an internal network.
No installer, plugin, imported business database or provider key was involved.
The final smoke evidence is `smoke-final/outcome.json`, `smoke-final/steps.json`, raw HTTP outputs and per-command logs.

Two earlier smoke harness prerequisites remain recorded:
internal-network published ports were unavailable, and initial PostgreSQL UNIX-socket readiness preceded final TCP startup.
The corrected harness probes process localhost with BusyBox wget and waits for PostgreSQL TCP readiness.
No server source or test assertion changed.
All exact owned smoke containers/networks were removed after logs/data preservation, with zero teardown errors.
Final Docker inventory matches the five original unrelated container IDs exactly. Testcontainers left no remaining session containers.
The workroot, source copies, dependency/build caches, baseline artifact, databases, scripts and full diagnostics remain for review.
The scoped download proxy process stopped after verification. Its script and logs remain.

No product implementation, task start, private .env use, deployment, paid request, commit, publication or evidence cleanup occurred.
Main must review these results before any new plugin plan or implementation authorization.
