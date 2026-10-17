# U21 Repository Layout Implementation

Date: 2026-10-18

## Source moves

Use standard source and patch directories. Keep public commands and fixed output names unchanged.
Direct `mv` commands preserved source contents and modes.

| Former path | Current path |
| --- | --- |
| `native/backend/` | `src/backend/` |
| `native/frontend/` | `src/frontend/` |
| `native/sub2api.patch` | `patches/sub2api.patch` |

The backend move contains 52 files. The frontend move contains one file.
The patch is the 54th file.
All 54 files match their former staged contents and executable modes.
The residual `native/` directory contained only Finder metadata and went to macOS Trash.
No compatibility wrapper or replacement archive directory was added.

## Commands and current references

`scripts/build-sub2api.py` reads `patches/sub2api.patch` and copies `src/backend/` and `src/frontend/` into upstream directories.
The existing source copy, patch order, frontend dependency reuse, checks, and compilation commands remain unchanged.
The build still writes `dist/build/sub2api-codex-eval`.
The release workflow calls the same public scripts and needs no path change.

The existing UI harness now reads `src/backend/internal/codexeval/web/`.
The native Docker fixture reads the fixed binary and its adjacent fixture helpers, so its paths remain unchanged.
The installer test sets `sys.dont_write_bytecode = True` before importing root `install.py`.
Its normal command no longer creates root `__pycache__/`.

README, AGENTS, and the directory spec show the actual root layout.
Current backend specs and source links use `src/` and `patches/sub2api.patch`.
Permanent directory and quality specs no longer describe intermediate cleanup history.
Historical research retains its original source paths. The mapping above records the move.
The main session owns current task metadata and context updates.

## Generated-file cleanup

All removals used `/usr/bin/trash`.
The following generated directories went to Trash:

- Root `__pycache__/`, `output/`, and `.playwright-cli/`.
- `tests/installer/__pycache__/`.
- `.codex/hooks/__pycache__/`.
- `scripts/__pycache__/`.
- `.trellis/scripts/common/__pycache__/`.

Finder `.DS_Store` files went to Trash from these directories:

```text
.
dist/
tests/
.agents/
.agents/skills/
.claude/
.claude/skills/
.codex/
.github/
.trellis/
.trellis/tasks/
.trellis/workspace/
.trellis/spec/
.trellis/.runtime/
.trellis/scripts/
.devin/
.devin/skills/
```

The residual `native/` move also removed its Finder metadata from the repository.
The ignore rule now uses the correct `.DS_Store` capitalization.
Existing Python cache, bytecode, browser-output, and build-output ignore rules remain.
No selected generated file remains outside preserved Git and dependency/build-cache paths at the final check.

Tool integration directories and their skills remain intact.
Private configuration, task research, workspace records, Git data, fixed outputs, and reusable dependency/build caches remain intact.
No broad cache purge ran.

## Single notice source and package

The duplicate root `THIRD_PARTY_NOTICES.md` went to Trash.
`licenses/THIRD_PARTY_NOTICES.txt` is the sole notice source.
README links to that file. Required licenses, source attribution, and native renderer notices remain intact.

The notice now identifies `src/` and `patches/sub2api.patch` as the current source locations.
The main session authorized an archive-only update because this packaged notice changed.
Ran `python3 scripts/package-sub2api.py` with the existing binary.
No binary build ran.

The archive still has exactly seven files: the binary, installer, and five license/notice texts.
The installer and all license/notice entries match their source files byte for byte.
The binary remains 125599904 bytes with its prior modification time.
The regenerated archive is 38700698 bytes.
No hashes were calculated.

## Checks

| Check | Result |
| --- | --- |
| Moved file contents and executable modes | All 54 files match the former staged paths. |
| Build overlay path resolution | New source owners resolve to upstream `backend/` and `frontend/`. Existing frontend dependencies are available. |
| Patch dry run | Passed all nine file edits against `/tmp/sub2api-codex-eval-v0217-check`. |
| Existing UI harness | Passed with `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check`. |
| Existing installer tests | `python3 tests/installer/test_native_version.py` passed both tests. Root `__pycache__/` remained absent. |
| Current local Markdown links | Passed for 15 documents. |
| Current obsolete-path scan | No `native/`, duplicate notice filename, or incorrect Finder ignore spelling remains outside historical evidence. |
| Python, shell, and Node syntax | Passed for current command and fixture files. |
| Ignore behavior | Finder metadata, Python caches, browser output, and Playwright output match ignore rules. |
| Archive inventory and source bytes | Exact seven-entry package and matching installer/license texts. |
| `git diff --check` | Passed |

## Final root and limits

```text
src/ patches/ scripts/ tests/ docs/ licenses/ dist/
install.py README.md AGENTS.md LICENSE
.github/ .trellis/ .agents/ .codex/ .claude/ .devin/ .git/
.gitignore .gitattributes
```

`dist/build/` contains the fixed executable and installation archive.
Other existing `dist/` dependency/build directories remain available.
Current source bytes and runtime behavior did not change.
No Docker operation, deployment, broad suite, race check, staging, or commit ran.
The earlier history cancellation report remains unchanged.
