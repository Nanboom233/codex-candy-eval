# Clientless Application Probe Review

## Findings (fixed)

- File: `install.py`
- Issue: Cached helper probe cleanup was non-strict and its exception was caught as a cache miss. A failed removal could therefore be followed by pull/reprobe and eventual success with an orphan. Source resolution also used non-strict cleanup and could return success after removal failed.
- Fix: Both paths now require bounded exact-name removal. Cleanup failure after otherwise successful work propagates. If probe/source work and cleanup both fail, the original work failure remains primary and the cleanup traceback is printed. Cleanup failure cannot enter cached-image fallback.
- File: `tests/installer/test_preflight.py`
- Issue: No focused regression covered cached/source success followed by cleanup failure or simultaneous primary and cleanup errors.
- Fix: Added those cases, including a no-pull assertion for cached helper cleanup failure.
- File: `tests/installer/test_docker_boundary.py`
- Issue: The cleanup-failure fixture rejected removal of every temporary tool, preventing tests from reaching the source resolver or builder cleanup under test.
- Fix: Scope the simulated denial to the source resolver or Buildx builder and update expected helper removals.

## Findings (not fixed)

No remaining code finding in this cleanup change. Live iStoreOS installation, native billing/concurrency, production browser behavior, and published bootstrap remain outside this review and are not established by local fixtures.

## Verification

- Lint: N/A; no Python linter is configured for the installer.
- TypeCheck: N/A; no Python type-check command is configured. `py_compile` passed for the installer and changed tests.
- Tests: Pass. Source `test_preflight.py` 18/18, `test_docker_boundary.py` 14/14, and `test_http_transport.py` 12 passed with 4 opt-in Docker cases skipped. Extracted archive preflight 18/18, boundary 14/14, and clientless HTTP 16/16 with Docker enabled and no skips. `git diff --check` passed.
- Spec: Existing [universal delivery contract](../../../spec/backend/universal-linux-delivery.md) already requires primary-error preservation and cleanup-failure reporting; no spec edit was needed.

## Independent Evidence

The implementation report's earlier archive checks predate this narrow cleanup correction. This review regenerated the canonical package, extracted it to a fresh directory, and reran the archived suites above. `install.py` matches the extracted archive byte-for-byte. The package validator accepted exactly 62 allowlisted regular files; functional runtime SHA256 remains `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759`.

Canonical archive: `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`, 518858 bytes, SHA256 `6ee8ffe55f43c564c22e764dd96be710548fbe7a03cf471283918b7b5ae5f823`.

Installer SHA256: `1275095a76e5eb8bbcc62d34aaa1a48e8e59344b68fea1bf1f9cb45ee8d8e316`.

An earlier clean-archive run before this repack had one intermittent fake-relay timeout in `test_actual_delayed_fixed_and_fragmented_chunked_requests`, subtest `(chunked=False, delay=0)`, at its 3-second request deadline. The exact archive command rerun and isolated case passed then; the current regenerated archive's full Docker-enabled 16-test matrix also passed. The initial timeout remains an observed intermittent result, not a claimed code fix.

## Retained Cleanup Inventory

No cleanup was performed. `/tmp/codex-application-clientless-fix-20261013` is retained at 9748 KiB with prior proof, logs, staging and extractions; the new review staging and extraction are also retained there. Docker has no matching `codexeval`/clientless/probe fixture containers. Retained images include `codexeval-clientless-fixture:20261013` (`sha256:c395b576ca7b62e5133a10942fdbb7b5c543222530de9510f580053c8a4ac74c`) and the cached `alpine/git:2.49.1` helper. Repository bytecode directories are `./__pycache__`, `./tests/installer/__pycache__`, and `./.trellis/scripts/common/__pycache__`. Shared Go build cache remains at 322280 KiB and was not cleared.
