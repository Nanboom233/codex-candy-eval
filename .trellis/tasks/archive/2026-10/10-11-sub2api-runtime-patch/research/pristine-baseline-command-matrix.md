# Pristine v0.2.16 Command Matrix — Upstream Evidence

Researcher reverified annotated tag v0.2.16 -> 01842576b3d5c1e2545b5a99ad555d36e28ed306.
Fresh official archive SHA256 eef7e0519b05b8254c981efbab7bfe74c5d89897c67053e81380987366daae89.
Research workroot: /var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/pristine-sub2api-v0.2.16.ifVBya.

| Gate | Actual upstream configuration | Execution distinction |
| --- | --- | --- |
| Backend unit | Go1.27.2, backend make test-unit -> go test -tags=unit ./... | Required upstream CI, no race flag |
| Backend integration | backend make test-integration -> go test -tags=integration ./... | CI=true forces unavailable Docker to fail instead of skip |
| Integration services | postgres:18.1-alpine3.23, redis:8.4-alpine through TestMain/testcontainers | Disposable owned services |
| Frontend required | Node20, pnpm9, frozen install, root make test-frontend | lint:check, typecheck and explicit critical Vitest list |
| Frontend release build | upstream package build and embed output | Actual full frontend build, preserved lock |
| Linux/amd64 runtime build | CGO_ENABLED=0, tags=embed, release ldflags, trimpath | Cross-compile from pristine source, verify version/architecture and disposable startup |
| Backend vet | go vet ./... | Supplemental diagnostic |
| Backend full race | go test -tags=unit -race ./... | Supplemental diagnostic and previous patched command comparison |
| Frontend full tests | pnpm run test:run | Supplemental beyond required critical CI list |

Current host is macOS arm64. Native source tests do not constitute Ubuntu Linux CI execution.
Any Linux/amd64 container execution uses emulation and is labeled accordingly.
Original patched153 failures/51 race warnings came from supplemental race invocation, not established as upstream required CI result.
Each row remains pending until actual command exit and complete diagnostics are recorded.
