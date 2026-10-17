# Artifact Cleanup Check

Review date: 2026-10-13. This was a read-only cleanup verification. No build,
test, deployment, paid request, Docker prune, image removal, or source cleanup
was performed by this check.

## Findings (fixed)

- The retained archive is present at its original path:
  `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz`.
- It is `73,379,445` bytes with SHA-256
  `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8`.
- Tar inspection reports exactly eight regular files and no nonregular members:
  `install.py`, two Linux binaries, and five `LICENSE/` texts.
- Archive and current tracked `install.py` bytes are equal at 39,082 bytes.
  Both hashes are
  `c14d8c69277550ed501477edc40dfabce49fac9c6ce966c0d7b0fdb1ee5713a3`.
- The functional roots `internal/eval/`, `patches/`, `compatibility/`, and
  `licenses/` retain their recorded counts and full tree hashes:

  ```text
  internal/eval/  38  0feb47f776d35d0ac9f98386894f66e5b28438adb178521e718375203c959bf0
  patches/       18  f92e2e878dfa4a2ddc665c992c339bf3331c0b540e7ae0320f2b9a74870d634b
  compatibility/  1  3ae13c68f7b784d046339f0b38436e5d0214a8d0caa03b644f029cf04dd795f2
  licenses/       4  335646704d751f999da92db44e584fdbaa20dd6ab97da91b6e5b8d18fe75c3f6
  ```

- Maintained sources remain present: `scripts/test-sub2api.sh`,
  `scripts/test_sub2api_build.py`, `scripts/test-ui.cjs`, `tests/`,
  `internal/eval/`, `patches/`, `compatibility/`, and `licenses/`.
- A source-only audit excluding `.DS_Store`, `__pycache__`, `.pyc`, and `.pyo`
  reports these current inventories: `tests/` 20 files,
  `internal/eval/` 37 files, `patches/` 18 files, `compatibility/` 1 file, and
  `licenses/` 4 files. Their current source-only tree hashes are, respectively,
  `000b4a21e835758d8645842ea3533fcb8ea7ccdae204e880d3f41f5b1dad1748`,
  `dd6d6eeea1f2bde8a200b993728bf4a784952f1ec7a4e24bdc0f8dac03b9a5e9`,
  `f92e2e878dfa4a2ddc665c992c339bf3331c0b540e7ae0320f2b9a74870d634b`,
  `3ae13c68f7b784d046339f0b38436e5d0214a8d0caa03b644f029cf04dd795f2`, and
  `335646704d751f999da92db44e584fdbaa20dd6ab97da91b6e5b8d18fe75c3f6`.
- Go build/test cache is reduced to 8 KiB and two files at
  `/Users/minami/Library/Caches/go-build`.
- `/tmp/sub2api-patch-check2`, `/tmp/sub2api-replay2`,
  `/tmp/sub2api-installer-check.log`, and the enumerated task-specific
  `/tmp`/`TMPDIR` fragment candidates are absent at their original paths.
- Read-only Docker inventories show no task container, network, or volume.
  The official `ghcr.io/wei-shaw/sub2api` image remains, as image removal was
  explicitly excluded.
- `git diff --check` passes.

## Findings (not fixed)

- Finder/tool activity regenerated `.DS_Store` metadata during verification,
  including `dist/.DS_Store` and `dist/sub2api-build/.DS_Store`. The retained
  archive is unchanged. The recorded `tests/` full-tree fingerprint differs
  only because `tests/.DS_Store` and `tests/host-overlay/.DS_Store` were
  regenerated; functional test sources remain intact. These metadata files
  were not deleted or chased in a loop.
- The macOS Trash directory cannot be inspected from this session (`du` is
  denied with `Operation not permitted`). The cleanup record states that
  inventoried artifacts were moved there with `/usr/bin/trash`, Trash was not
  emptied, and their disk blocks remain unreclaimed. No permanent deletion was
  attempted.

## Verification commands

The check used only read-only commands, including `find dist -type f`, `stat`,
`shasum -a 256`, `tar -tvzf`, a `PYTHONDONTWRITEBYTECODE=1` archive/source byte
comparison and functional-root fingerprint calculation, Docker `ps`/`network
ls`/`volume ls`/`image ls`, `find` for task temporary paths, and `git diff
--check`. No test or build command was run. The implementation inventory and
cleanup actions are recorded in `research/artifact-cleanup.md`.

