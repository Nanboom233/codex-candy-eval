# Source Resolver Check Review

## Findings (fixed)

- File: `install.py`
- Issue: The initial immutable-image cache-probe edit during review placed image-ID validation and execution under an early-return branch. A valid cached helper therefore raised `UnboundLocalError`. Focused regressions caught this review-time editing error; it was corrected before final checks and packaging. It was not the user-reported helper-pull failure.
- Fix: Move image-ID validation and the no-pull probe onto the valid-metadata path. `resolve_source()` now retains the acquired image ID and uses that same immutable ID for its networked source-resolution run, closing the tag-rebinding gap between probe and use.

- File: `install.py`
- Issue: The atomic `mv; sync` command can replace the target and then fail during `sync`. Marking replacement only after the command returns would leave this uncertain-effect failure on the non-restoration path.
- Fix: Set the replacement-attempt flag before running `mv; sync`, so any command failure attempts original-binary restoration. `tests/installer/test_transactions.py::test_uncertain_atomic_replace_failure_restores_original` injects this failure and asserts restoration and journal ordering.

- File: `tests/installer/test_docker_boundary.py`
- Issue: The warm-cache regression checked only the probe image ID, not the image used for actual source resolution.
- Fix: Assert that both commands use the same ID.

- File: `tests/installer/test_preflight.py`
- Issue: Two resolver fakes returned empty metadata and did not model the acquired-image contract.
- Fix: Return a valid immutable image ID from those fakes.

- File: `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`
- Issue: The existing archive did not contain the final reviewed installer.
- Fix: Preserve its prior bytes at `/tmp/codex-source-resolver-fix-20261013/archive-before-checker-fix.tar.gz`, then rebuild and extract the current 62-file allowlisted package. Staging and extracted files compare byte-for-byte.

## Findings (not fixed)

- No unresolved local defect was found in the reviewed resolver, raw-diagnostic, cleanup, or replacement-attempt paths.
- Native amd64 lifecycle, iStoreOS acceptance, full live installation, deployed browser flow, paid billing/concurrency, and published online bootstrap remain unverified. Earlier arm64 Buildx/standard build and 18-point lifecycle evidence predates this final installer revision; it is not counted as a fresh pass.
- Required broad historical gates remain failed without a pristine A/B baseline for attribution. The local repository check below passes; that does not resolve the historical failures.

## Verification

- Lint: pass (`git diff --check`, shell syntax check).
- TypeCheck: pass (`go vet ./...`; included in the integrated check).
- Tests: pass. `scripts/test-sub2api.sh check` completed with 72 foundation helpers (1 historical skip), 50 transaction tests, 8 universal tests, 16 preflight tests, and 14 Docker-boundary tests; Go race tests and UI checks also passed. The focused transaction and boundary suites passed after the final edits.
- Archive tests: extracted-package preflight 16/16 and Docker-boundary 14/14 passed. Staging and extraction matched exactly; 62 regular files, 526230 bytes, SHA256 `b39fe125b06da0363f100709b89a0247a7c630e07526af7a6c523b0dee0eaab2`.
- Actual resolver: Docker 29.9 client / 29.8.2 arm64 daemon; cached helper ID `sha256:c0280cf9572316299b08544065d3bf35db65043d5e3963982ec50647d2746e26`. Current `resolve_source()` completed against `refs/tags/v0.2.15`, resolving `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` and Go `1.27.2`. Log: `/tmp/codex-source-resolver-fix-20261013/actual-resolver-arm64-final/source-resolution.log`.
- Earlier actual helper/build evidence is preserved in `source-resolver-implementation.md` and is reported only within its original scope. No deployment, paid request, publication, commit, push, or evidence cleanup occurred.

## Retained Evidence

The complete checker workroot is `/tmp/codex-source-resolver-fix-20261013`. It contains the original source-resolver evidence, the pre-checker archive copy, the checked release staging/extraction, and the fresh actual resolver log. Nothing in that workroot was removed.
