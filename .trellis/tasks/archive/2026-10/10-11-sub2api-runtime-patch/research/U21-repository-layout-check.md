# U21 Repository Layout Check

Review date: 2026-10-18.

Main-session follow-up: the index now matches the current 138-path candidate.
Stale staged `native/` entries and the duplicate root notice were removed from the index.
Current source paths, current caller changes, and original CPA removals are staged. No commit occurred.

## Findings (fixed)

- File: `dist/sub2api-build/.DS_Store`
- Issue: A Finder metadata file remained under the preserved build-cache tree,
  although U21 requires Finder metadata cleanup from `dist/`.
- Fix: Moved the file to macOS Trash with `/usr/bin/trash`. Build-cache
  contents and fixed delivery outputs were not changed.

## Findings (not fixed)

- The Git index still contains the prior staged `native/` additions and the
  root `THIRD_PARTY_NOTICES.md` addition, while the working tree contains the
  new `src/`, `patches/`, and sole notice paths. This is the expected
  no-staging handoff for the layout correction. The main session owns refreshing
  the work-commit candidate; this review did not stage or alter the index.

## Source and caller review

The 54 moved source inputs match their former staged contents and executable
modes:

- `native/backend/` -> `src/backend/` (52 files)
- `native/frontend/` -> `src/frontend/` (1 file)
- `native/sub2api.patch` -> `patches/sub2api.patch` (1 file)

The old `native/` tree and root `THIRD_PARTY_NOTICES.md` are absent. The build
script reads `patches/sub2api.patch` and overlays `src/backend/` and
`src/frontend/`. The UI harness reads the moved embedded page. Current README,
AGENTS instructions, backend specs, operations guide, and notices use the new
paths. No current command, workflow, fixture, or UI harness reads an old
`native/` path or the removed notice file.

The fixed binary remains at `dist/build/sub2api-codex-eval` with size
`125599904` bytes. The regenerated archive remains at the fixed path and has
exactly seven entries: the binary, `install.py`, and five required license or
notice files. Archive installer and license entries match their source files.
The sole notice source is `licenses/THIRD_PARTY_NOTICES.txt`.

Current Finder metadata and Python bytecode scans are clean outside Git and
task evidence. The reusable `dist/` build/dependency cache remains available;
only generated metadata was removed. Private files, agent/tool directories,
Git data, task research, and unrelated caches remain untouched.

## Verification

- Lint: N/A; U21 changes repository layout and generated files only, and no
  configured linter applies to this scope.
- TypeCheck: N/A; no typed source changed and the upstream Go module is outside
  this repository.
- Tests: existing owner evidence reused. The U21 implementation report records
  the existing UI harness and two installer tests passing; no duplicate product
  run was needed for a path-only review.
- Moved content and mode comparison: pass for all 54 files against the old
  staged paths.
- Current Python, Shell, and Node syntax: pass.
- Current Markdown links: pass, 76 relative links.
- Archive inventory and installer/license source bytes: pass, 7 entries.
- Finder metadata and bytecode cleanup: pass outside excluded Git/task paths.
- `git diff --check`: pass.
- No Docker operation, binary rebuild, deployment, staging, commit, or hash
  calculation was performed by this review.

The main session owns the staged-path refresh, task metadata, and final commit
approval.
