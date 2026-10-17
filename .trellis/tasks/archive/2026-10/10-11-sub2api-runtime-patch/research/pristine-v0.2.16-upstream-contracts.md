# Research: pristine-v0.2.16-upstream-contracts

- Query: Identify the exact upstream CI, dependency-lock, test harness, frontend embed, and Linux/amd64 build contracts for pristine Sub2API v0.2.16.
- Scope: internal
- Date: 2026-10-14

## Findings

The verified annotated tag `v0.2.16` resolves to commit `01842576b3d5c1e2545b5a99ad555d36e28ed306`. The official codeload archive is `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/pristine-sub2api-v0.2.16.ifVBya/source-v0.2.16.tar.gz`, SHA-256 `eef7e0519b05b8254c981efbab7bfe74c5d89897c67053e81380987366daae89`. The extracted source is under `.../extracted/sub2api-01842576b3d5c1e2545b5a99ad555d36e28ed306`; the archive has no Git metadata. `source-files-before.sha256` and `lockfiles-before.sha256` in the workroot preserve file fingerprints. The frozen dependency fingerprints are `backend/go.mod` `809e6806dc248bd44cea4198ddc21fd6ed04f80003300c3668367b91bae319a8`, `backend/go.sum` `b17668c6300c04027bd397650e60e28e514dc4abb9ab113e3cea634fdd708109`, `frontend/package.json` `d1fd3d91ccf9bbca6cf0a79428343109a16efa6b3675fd01869d37db5d63083e`, and `frontend/pnpm-lock.yaml` `1d8e54cc8e1f1bd3e7593805f607d9314ff81dab98ab5a9a6396475e148b26ce`.

Upstream CI is defined by `.github/workflows/backend-ci.yml`. It requires Go `1.27.2`; backend `make test-unit` invokes `go test -tags=unit ./...` (no race flag), and `make test-integration` invokes `go test -tags=integration ./...`. The workflow also runs `golangci-lint` v`2.14.0` with `--timeout=30m`. Frontend CI installs pnpm major `9`, Node `20`, and uses `pnpm install --frozen-lockfile`, then root `make test-frontend`. Root `Makefile` runs `frontend` `lint:check`, `typecheck`, and the explicit `test-frontend-critical` Vitest file list. Treat `go vet ./...`, `go test -tags=unit -race ./...`, and `pnpm run test:run` as supplemental diagnostics, separate from those exact CI gates. For a claim covering the complete workflow, note its separate `shell` job and `release-helpers` Python/shell job as well.

Run the repository integration command with `CI=true`: `backend/internal/repository/integration_harness_test.go:53-60` exits successfully without running tests when Docker is unavailable and `CI` is empty, but fails when any nonempty `CI` value is set. This package harness uses PostgreSQL `18.1-alpine3.23` and Redis `8.4-alpine`, applies migrations, and shares those services across repository integration tests (`:32-35`, `:63-100`). It calls `os.Exit(m.Run())` at `:134`; deferred container termination at `:76` and `:86` is therefore not executed on normal completion. Testcontainers may use Ryuk, but verify disposable containers are gone after the command. Separate integration packages also start Redis testcontainers; `internal/server/routes/auth_rate_limit_integration_test.go:73-79` can still skip based on Docker socket detection regardless of `CI`, so preserve actual test/skip counts rather than inferring complete coverage from command success. `internal/integration/e2e_gateway_test.go:70-86` has a separate `TestMain`; individual live tests may skip when API keys are absent or the target account is unavailable.

The frontend release output is `backend/internal/web/dist` (`frontend/vite.config.ts:106-109`; release workflow `.github/workflows/release.yml:98-105`). The real build is `cd frontend && pnpm run build`, which executes locale tests, `vue-tsc -b`, and Vite (`frontend/package.json:8-10`). The release binary uses `-tags=embed`; `backend/internal/web/embed_on.go:29` embeds `all:dist`. Build the frontend before compiling that variant. The lock format is pnpm lockfile v9 (`frontend/pnpm-lock.yaml:1`); retain frozen install and do not run `go mod tidy` or dependency-updating commands during the baseline.

Release compile settings are in `.goreleaser.yaml:9-32`: backend `./cmd/server`, `-tags=embed`, `CGO_ENABLED=0`, stripped symbols, and `main.Commit`, `main.Date`, `main.BuildType=release`. The release workflow targets Linux/amd64 and uses Go `1.27.2` (`.github/workflows/release.yml:108-165`). The backend Makefile's ordinary `build` target instead uses `CGO_ENABLED=0`, `-trimpath`, and `main.Version`, but omits `-tags=embed` (`backend/Makefile:3-7`). A direct lock-preserving Linux/amd64 release-equivalent compile should therefore set `GOOS=linux GOARCH=amd64 CGO_ENABLED=0`, pass `-tags=embed`, and set release ldflags explicitly. Do not call a direct compiler build a GoReleaser release/archive. The codeload tree has `backend/cmd/server/VERSION` = `0.2.15`; `release_matrix.py` normally rewrites that file from the tag before building (`.github/release-tools/release_matrix.py:53-73`). Because the pristine evidence tree must remain unchanged and has no Git metadata, explicitly set `-X main.Version=0.2.16` in the direct-build command (alongside commit/date/build type) so `-version` reports the tested tag. `main.go:29-50,59-75` documents that version fallback and CLI output.

For loaded runtime checks, `-version` is a CLI-only startup path (`backend/cmd/server/main.go:59-67`). Setup mode starts a setup-only server and registers setup routes (`main.go:69-90,97-132`); its status route is `/setup/status` (`internal/setup/handler.go:22-36`). Normal configured mode registers `/health` and `/setup/status` (`internal/server/routes/common.go:9-31`). Thus setup-mode status and normal-mode health are distinct checks. Keep the server and any PostgreSQL/Redis services disposable; use only the upstream binary built from this extracted tree.

### Files Found

- `backend-ci.yml` — upstream Go, frontend, linter, shell, and release-helper CI gates.
- `Makefile`, `backend/Makefile` — critical frontend list, test tags, and ordinary backend build target.
- `backend/go.mod`, `backend/go.sum` — Go toolchain declaration and dependency lock set.
- `frontend/package.json`, `frontend/pnpm-lock.yaml`, `frontend/vite.config.ts` — frontend scripts, frozen package lock, and embedded build output.
- `backend/internal/repository/integration_harness_test.go` — integration `TestMain`, Docker handling, service images, migrations, and cleanup structure.
- `backend/internal/server/routes/auth_rate_limit_integration_test.go` — separately skipped Redis testcontainer coverage.
- `backend/internal/integration/e2e_gateway_test.go` — E2E `TestMain` and external-key-dependent test setup.
- `.goreleaser.yaml`, `.github/release-tools/release_matrix.py`, `.github/workflows/release.yml` — release compiler flags, version-file mutation, and target workflow.
- `backend/cmd/server/main.go`, `backend/internal/setup/handler.go`, `backend/internal/server/routes/common.go` — version/setup/normal-server smoke-test behavior.

### Related Specs

- `.trellis/workflow.md` — research is a planning artifact; preserve the source tree and record findings in the task research directory.
- `.trellis/spec/backend/index.md` — requires task/spec context and separation of upstream source checks from local patched-source checks.
- `.trellis/spec/backend/build-guards.md` — useful local lockfile and reproducibility guidance; its patch-specific requirements are not upstream Sub2API CI gates.
- `.trellis/spec/backend/minimal-linux-delivery.md` — local delivery/runtime verification context; do not substitute it for the pristine upstream baseline.

### External References

No external documentation was needed; all version and command contracts above come from the exact tagged source and workflow files.

## Caveats / Not Found

No build or test command was run in this research pass. The current host is macOS arm64; local native test results do not prove Ubuntu/Linux execution equivalence, and Docker's Linux VM is aarch64. Keep those platform facts alongside every result. The direct Linux/amd64 compiler command and actual smoke-test outcomes belong in the execution report/checkpoint, not this source-contract report. Do not infer pristine upstream test failures from the earlier patched runtime's race failures.
