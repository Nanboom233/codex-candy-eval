# Research: Test and Release Isolation

- Query: Separate production functionality and release inputs from test, fixture, mock, audit, log, and development-validation material while preserving the four evaluations and their security/runtime behavior.
- Scope: internal
- Date: 2026-10-12

## Findings

### Current overlap

The existing patch series is mixed: all six patches add production changes and `_test.go` files. The current additions total 13 Go test files: patch 001 adds six; 002 adds two; 003 adds two; 004, 005, and 006 add one each. Examples are `patches/sub2api/v0.2.15/001-codex-eval-auth.patch:2,200,475,1244,2152,3551`, `002-codex-eval-scheduling.patch:26,1527`, and `006-codex-eval-topology-runtime.patch:772`.

`internal/eval` contains 28 paths: seven `_test.go` files and `testdata/modeltrace-golden.json` are test-only. `scripts/sub2api_build.py:472-486` currently calls unrestricted `shutil.copytree(src, dest)`, so all 28 paths enter the staged host source tree. Go's server binary build does not compile `_test.go`, but their source and the golden test fixture still enter that build workspace.

The runtime uses explicit business data, which must remain: Codex context, fingerprint probes and baselines, and the ModelTrace bank and its license are embedded at `internal/eval/codex.go:36-46`, `fingerprint_algorithm.go:38-45`, and `modeltrace_algorithm.go:16-20`. `ModelExecutor`, `TargetResolver`, and `StateStore` are production contracts at `internal/eval/eval.go:81-86`, `runs.go:19-23`, and `state.go:40-47`; they are useful test seams but also required runtime abstractions, not mock-only code.

The test harnesses are development tools: `scripts/test_sub2api.sh:39-53,56-87` runs helper, local, patched-backend, and frontend checks; its integration mode currently validates and tests the shared prepared workspace. `scripts/test-ui.cjs:1-50` uses a mocked DOM and `node:vm`. The compatibility helper tests live in `scripts/test_sub2api_build.py`.

The current `.github/workflows/release.yml:11-24` checks the CPA plugin, builds platform `.so`/`.dylib`/`.dll` outputs, and the publish job uploads `assets/*` (`:63-99`). It is not yet a Sub2API release workflow. A future release must create an explicit Sub2API allowlist package and must not upload the development workspace, test overlays, check logs, or test reports.

### Recommended separation

Use two version-specific patch identities. Put production changes under `patches/sub2api/v0.2.15/production/`; extract test additions from legacy patches 001-006 into `patches/sub2api/v0.2.15/test-overlay/`. Keep planned production patches 007 runtime and 008 UI within the production series; their tests go only in the overlay. Record old 001-006 digests as historical audit evidence, not as a constraint on the new split.

Move maintained evaluation tests to a separate overlay tree, such as `tests/overlays/codexeval/`, preserving their original package-relative paths and filenames. The test-workspace materializer copies those files into the same `backend/internal/codexeval` package and applies the Sub2API test overlay there, so same-package tests retain access to unexported implementation details. Test overlays may add `_test.go`, test fixtures, and frontend test files only; they must not change production source, generated runtime files, compatibility behavior, or lockfiles. Use fakes in `_test.go` files instead of test-only branches or mock implementations in production files.

Stage release and test workspaces independently from the same pinned source commit and production patch set. The production workspace receives only production patches and an explicit evaluation-core file allowlist. The disposable test workspace receives the same production inputs, then the test overlays. Compare path/hash manifests for every production file between the two workspaces before test execution and after overlay application; reject overlays that touch production paths. Run patched product Go/frontend integration tests only in the disposable workspace. Repository-owned build-helper tests and mocked UI harnesses may run in the development checkout or CI, but their source, caches, logs, fixtures, reports, and test binaries stay out of release staging and deployment.

The core allowlist must include production Go files, embedded product UI assets when added, and runtime data/license resources. Exclude `_test.go`, `testdata/`, golden fixtures, mocks, audit samples, and dev-only documentation from the copied runtime package. Preserve the existing evaluation algorithms, prompts, history, normal Responses transport, auth/Key ownership, topology evidence, billing, concurrency, quota and model checks, exact-target/no-fallback behavior, and UI. Keep installer identity/checksum/rollback checks and the operator topology-evidence client as product operations; they are not unit-test fixtures. The operator client runs from the administrator's environment before activation, not inside the deployed server container.

### Identity and packaging rules

Maintain distinct identities:

- `productionPatchSetSha256`: ordered production patch names and hashes; included in production provenance.
- `testOverlayPatchSetSha256`: test-only overlay identity retained in development check evidence, never in the release package or deployed server.
- `productionCoreSourceTreeSha256`: allowlisted maintained production core and required runtime resources.
- `productionStagedTreeSha256`: complete staged production package, recalculated after UI or other production resources are added. This is a new full-tree digest and is distinct from the historical 28-path tree digest.
- `testAugmentedTreeSha256`: complete ephemeral test workspace digest, recorded only in development evidence.
- Per-file production path/hash manifests: compare production paths in the production and test workspaces to prove both derive from identical production bytes.

Keep historical patch/core digests with the existing audit records. Do not claim the old complete-tree digest still identifies a new tree containing additional UI assets. A test-overlay change creates a new test overlay identity and augmented-tree digest without changing the production release identities when production bytes are unchanged.

Put a packaging allowlist check in build/release tooling, outside the product binary and outside the deployed runtime. Construct the release directory from explicit production outputs, then reject any unexpected path or test/dev category before upload. Allow only the required server binaries, the operator evidence client if distributed, installer entry points, compatibility manifest, production provenance/checksums, public operational documentation, and license notices. Exclude patch sources, Go/JS/Python test files, testdata, fixtures, mocks, test harnesses, audits, logs, test reports, temporary workspaces, and caches. Deploy only the server binary and required runtime state/resources; do not install test scripts or reports into the container.

The operator evidence client and installer preflight checks remain product operations: they establish the protected deployment facts required before paid evaluation activation and safe installation. Developer unit/integration/browser suites, mocks, synthetic gateway data, and local audit commands are development verification and do not ship or execute in the deployed runtime. Normal candy, fingerprint, ModelTrace, and pelican runs remain the product's evaluation feature.

## Files Found

- `patches/sub2api/v0.2.15/001-006-*.patch`: production changes are interleaved with 13 added Go test files across the six patches.
- `internal/eval/`: 28 paths currently combine runtime Go/data/license assets with seven Go test files and one golden fixture.
- `compatibility/sub2api-v0.2.15.json:60-68`: one flat patch directory and one whole-directory core source/destination are currently declared.
- `scripts/sub2api_build.py:319-334,472-486`: validates source presence and symlinks, then copies every file without a production allowlist.
- `scripts/build-sub2api.sh:149-224`: applies one patch set, stages the core, regenerates Wire, builds frontend/server, and writes provenance/checksums.
- `scripts/test-sub2api.sh:39-87` and `scripts/test-ui.cjs:1-50`: development checks currently share the prepared build workspace or use a mocked Node DOM.
- `.github/workflows/release.yml:11-99`: current workflow tests and releases CPA plugin binaries with a broad `assets/*` upload.
- `.trellis/spec/backend/{build-guards,evaluation-core,run-state,quality-guidelines}.md`: current build, runtime-data, persistence, and development-check contracts.

## Code Patterns

- Runtime embedded assets are declared by production `go:embed` directives, while golden fixtures are opened only by tests (`internal/eval/fingerprint_algorithm.go:38-45`, `codex.go:36-46`, `modeltrace_algorithm.go:16-20`, `modeltrace_test.go:13-18`). This provides an explicit allowlist boundary.
- Production core interfaces permit test doubles without adding mock branches to business code (`internal/eval/eval.go:81-86`, `runs.go:19-23`, `state.go:40-47`). Keep these interfaces and place fakes in the test overlay.
- Current stage copying has no exclusions (`scripts/sub2api_build.py:472-486`); provenance fingerprints the full source tree (`:513-537`). Extend the production manifest/helper to hash the allowlisted staged tree, and compute the overlay identity separately.
- The existing developer workflow already separates test and build commands but not source workspaces (`scripts/test-sub2api.sh:56-87`); the release replacement should build from a fresh production-only workspace and run checks in a second disposable workspace.
- Existing release packaging uses a wildcard upload (`.github/workflows/release.yml:89-99`); use an explicit output allowlist for Sub2API artifacts.

## External References

- No external documentation was fetched. Source pin remains Wei-Shaw/Sub2API v0.2.15, commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`; toolchain pins remain Go 1.27.2 and pnpm 9.15.9.

## Related Specs

- `.trellis/spec/backend/build-guards.md`: strict patching, pinned build inputs, source/core provenance, and separation of helper checks from runtime acceptance.
- `.trellis/spec/backend/evaluation-core.md`: production algorithms, embedded datasets/licenses, and copied-package expectations.
- `.trellis/spec/backend/run-state.md`: paid-run lifecycle, durable state, and test requirements.
- `.trellis/spec/backend/quality-guidelines.md`: Go/UI development checks and the current CPA-specific release requirement that must be replaced in the migration.
- `.trellis/workflow.md`: one implementation handoff followed by a full quality check; this physical input split does not create review batches.

## Caveats / Not Found

- No code, release packaging, or workflow changes were made, and no build or test was run for this research.
- The current release workflow still builds and publishes the CPA plugin; replacing it remains part of the planned implementation.
- The known full unit-race failures and live-target acceptance gaps remain unchanged.
- Exact final allowlist filenames and new hashes require implementation-time artifacts; historical digests are evidence only and do not constrain the user-requested production/test split.
