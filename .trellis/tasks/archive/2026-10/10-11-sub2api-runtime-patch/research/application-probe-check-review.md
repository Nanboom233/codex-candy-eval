# Application Probe Check Review

Reviewed 2026-10-13 against the active application-probe contracts and the
canonical universal archive. No product source or spec changes were needed.

## Findings (fixed)

- None. The HTTP transport, listener ownership and health-address reuse paths
  match the active contract. No concrete implementation defect was found.

## Findings (not fixed)

- None.

## Verification

- Lint: pass. `scripts/test-sub2api.sh check` passed its shell syntax, Go vet,
  Go race, and UI checks.
- TypeCheck: pass. `go vet ./...` and `go test -race ./...` compiled and checked
  the Go packages; the installer has no separate configured Python typecheck.
- Tests: pass. The full check gate covered 72 foundation helpers (one
  historical skip), 50 transaction tests, 8 universal tests, 16 preflight
  tests, 14 Docker-boundary tests and 10 HTTP tests (two explicit Docker cases
  skipped in that gate).
- Working-tree transport: `PROBE_DOCKER_TESTS=1 python3
  tests/installer/test_http_transport.py` passed 10/10, including the actual
  cached BusyBox transport and packaged CLI fixture cases.
- Canonical-archive preflight: 16/16 passed against a fresh extraction.
- Canonical-archive Docker boundary: 14/14 passed against a fresh extraction.
- Canonical-archive HTTP transport: 10/10 passed with Docker cases enabled and
  no skips.
- Documentation/spec sync: `docs/sub2api-operations.md` describes bounded
  direct HTTP probes, app-owned listener discovery, safe loopback fallback,
  preservation of concrete transport errors and health-address reuse. These
  agree with the active runtime-installation contracts. No documentation
  update was needed.
- Package identity: the canonical archive is 518412 bytes with SHA256
  `73d4c8b6fde68309f49d2ca7f32a761ae01e3193d81b0aff786d8dea6b2fc4db`; its
  fresh extraction contains 62 files and its installer SHA256 is
  `8e04d4fd11760e220873be64cb5bc3b51fa00ea203a0539c1c1489024bdf04bf`.

The first direct archived Docker-boundary invocation omitted
`PYTHONDONTWRITEBYTECODE=1`. Python wrote `__pycache__` into that extracted
bundle, after which `local_bundle()` correctly rejected the extra file and the
CLI fell back to the download step. The canonical check script exports this
variable. A fresh extraction run with it set passed all archived suites; the
initial extracted tree was preserved for evidence and was not used for the
authoritative results.

## Retained Evidence and Cleanup Inventory

No cleanup was performed. The two evidence roots remain available until the
main session dispatches cleanup:

- `/tmp/codex-application-probe-fix-20261013`: 10844 KiB, including original
  implementation evidence and the reviewer's fresh archive extraction at
  `checker-archive-extracted`.
- `/tmp/codex-application-probe-research-20261013`: 316 KiB of raw transport
  observations and source snapshots.

The repository currently has generated Python bytecode at `__pycache__`
(104 KiB), `scripts/__pycache__` (60 KiB) and
`.trellis/scripts/common/__pycache__` (200 KiB). The shared Go build cache at
`/Users/minami/Library/Caches/go-build` is 322280 KiB. No matching application
probe Docker containers or local test processes remained after verification.
Shared Docker images/builders and language package caches were left untouched.

No deploy, paid request, production access, release publication, commit or
cleanup occurred. The Docker CLI fixture verified transport and no-write
behavior only; it does not establish a native Sub2API installation, billing,
production concurrency or deployed-browser result.
