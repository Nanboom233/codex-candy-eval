# Docker Format Compatibility Review

## Findings (fixed)

- None. The reviewed installer and cleanup paths no longer pass `--format` to Docker, and no scoped defect required a code fix.

## Findings (not fixed)

- The user's iStoreOS Docker CLI and Python 3.9 environment remain unverified. The old-CLI subprocess cases are simulated compatibility checks, not host acceptance; this review ran under Python 3.14.
- The previously recorded full native race failure and live N150/iStoreOS acceptance gaps remain outside this Docker CLI review. No deployment or broader native run was performed.
- `tests/installer/test_http_transport.py` retains a `docker ps --format` use for opt-in fixture cleanup. It is test-only and excluded from the package; production `install.py` and `cleanup.py` have no such invocation.

## Verification

- Lint: N/A; no Python lint configuration or installed linter is present.
- TypeCheck: N/A; no Python type-check configuration or installed checker is present.
- Tests: PASS. Compatibility subprocess suite: 9 passed. Cleanup regressions: 10 passed. Current installer suite: 108 run, 104 passed, 4 explicit Docker skips.
- Archived checks: both `--help` entry points passed; all 9 compatibility tests passed against the extracted archived scripts. Packager raw-tar validation confirmed exactly eight regular payload files plus `LICENSE`, with source, license, and binary parity.
- Read-only Docker checks: `docker ps -a --quiet --no-trunc` succeeded. `docker buildx ls` succeeded and reported the unavailable default-builder daemon on stderr; this diagnostic remains visible.
- Archive: `dist/build/sub2api-codex-eval.tar.gz`, 39,016,308 bytes, SHA-256 `39c876d274a6499479bf0fbcbcf8d93c46c52053b01b2dbbc6ee3b497554f6fc`.
- Runtime: `dist/build/sub2api-codex-eval`, 125,632,672 bytes, SHA-256 `eed01f71c39048eef1a15fbf13652b5c44d9915fff1e34c59f3d58abf1c6c9e1`.

No runtime build, resource mutation, deployment, paid request, commit, or publication was performed.
