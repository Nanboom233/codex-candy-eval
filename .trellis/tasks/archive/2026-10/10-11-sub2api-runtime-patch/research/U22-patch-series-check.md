# U22 Patch Series Check

Review date: 2026-10-18.

## Findings (fixed)

None. No in-scope mechanical defect required a repair.

## Findings (not fixed)

- The index and commit candidate still reflect the prior U21 paths until the
  main session refreshes them. The main session owns that handoff. This review
  did not stage files or change the index.

## Candidate review

The current candidate uses seven production patches in the order listed by
`patches/series`. The series covers 50 production paths: 41 additions and 9
upstream edits. None of the production patches contains a Go test or test-only
data path. All 17 test files live under `tests/`, including the 12 Go test and
golden-data inputs under `tests/backend/`.

Only `scripts/build.py` and `scripts/install.py` remain in `scripts/`. The
build entry reads the series, applies each patch, attaches backend tests only
for `--check`, runs the existing UI and installer checks, builds, and packages.
The installer test imports its current `scripts/install.py` path. The UI
harness uses the effective patched upstream tree. The repository has no
production `src/` overlay or `docs/` directory.

README contains installation, manual rollback, backup lifetime, and history
instructions. The release workflow calls the combined build/check/package
entry. The legal notice identifies the patch series, Sub2API v0.2.17,
ModelTrace, and renderer licenses. The package keeps the installer's archive
name as `install.py` beside the executable.

## Verification

- Lint: N/A; no configured linter applies to this layout change.
- TypeCheck: N/A; the review did not change typed source. The upstream Go
  module remains the build owner.
- Tests: the U22 implementation report records successful focused Go tests,
  UI harness, installer tests, and patch composition in both production and
  check trees. This review reused that evidence. It independently ran the two
  existing installer tests after their import-path move.
- Entry help: `scripts/build.py --help` and `scripts/install.py --help` pass.
- Syntax: Python AST, Node syntax, fixture shell syntax, and `git diff --check`
  pass.
- Current Markdown links: pass, 70 relative links across README, AGENTS, and
  active backend specs.
- Patch series: pass, 7 listed patches and 50 production paths; no test or
  test-data paths appear in production patches.
- Layout: pass, exactly 2 script entries, 17 test files under `tests/`, and no
  `src/` or `docs/` directory.
- Package: pass, exactly 7 expected entries. The executable, installer, and
  five license/notice entries match their source bytes.
- No server build, Docker operation, deployment, staging, commit, or hash
  calculation was performed by this review.

The U22 implementation report provides the patch-application, package-test,
and UI-harness evidence. The main session owns specification synchronization,
the final work-commit candidate, and commit approval.
