# User-authorized artifact cleanup

## Scope and inventory

The user requested generated-artifact cleanup and explicitly retained maintained
test sources. This follow-up removes generated files only, without rebuilding,
running tests, deployment, paid requests, commit, publication or task archival.
Main owns task/spec changes. Other agents' source/configuration work is retained.

Before cleanup, allocated space (`du -sk`) is:

| Root | KiB |
| --- | ---: |
| dist/ | 49,320,744 |
| /Users/minami/Library/Caches/go-build | 4,162,716 |
| /tmp/sub2api-patch-check2 | 47,560 |
| /tmp/sub2api-replay2 | 64,224 |
| .playwright-cli/ | 36 |

The APFS filesystem initially reports 9,874,512 KiB available (`df -k .`).
Disk accounting may differ from the sum of directory allocations because APFS
shares/clones blocks and other applications may write concurrently.

## Retained payload

Keep only this generated release artifact at its original path:

`dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz`

Expected SHA-256:
`560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8`.
It contains exactly eight regular files: install.py, both Linux server binaries
and five LICENSE/ texts. The extracted release-final copy is redundant and is
authorized for removal.

## Removed or moved removal roots

- Every immediate dist entry except sub2api-build: .DS_Store, corepack,
  sub2api-v0.2.15, sub2api-codex-eval-v0.2.15-linux-amd64-arm64,
  sub2api-codex-eval-v0.2.15-linux-amd64-arm64.tar.gz and its .sha256 sidecar.
- Every sub2api-build entry except minimal-python-20261013. These inspected
  roots are this migration's pinned/replay/production/test workspaces, previous
  workspaces, duplicate release directories, build/test binaries, official-image
  export, temporary fixtures/credentials, raw logs and browser snapshots.
- Every minimal-python-20261013 entry except the retained archive, including
  release-initial/release-final, both fixtures and raw evidence logs.
- Repository .playwright-cli snapshots (11 generated files dated 2026-10-12).
- Generated __pycache__ roots: repository root, tests/installer, scripts,
  .codex/hooks and .trellis/scripts/common. Remove standalone generated pyc/pyo
  files if found outside those roots, without removing Python sources.
- /tmp/sub2api-patch-check2 and /tmp/sub2api-replay2 (inspected backend/replay
  workspace trees), /tmp/sub2api-installer-check.log, task-specific codex-eval
  check logs and superseded 001/003 patch snapshots listed below.
- TMPDIR/sub2api-backup-part-* inspected fragment files: 52 files, 221 logical
  bytes. These task-named test fragments are not shared temporary directories.
- Go build/test cache through `go clean -cache -testcache`; preserve
  /Users/minami/go/pkg/mod and global npm/pnpm/corepack stores.

No task-owned Docker containers remain; only unrelated kevbot containers are
present. No Docker prune or image removal is authorized or performed.
Unrelated /tmp/installbuilder_installer.log remains untouched.

## Preserved source fingerprints before cleanup

For each root, SHA-256 hashes the sorted `path + NUL + file SHA + newline` rows,
excluding generated __pycache__/pyc/pyo. These fingerprints will be compared
again after cleanup without tests/builds.

| Root | Files | Tree SHA-256 |
| --- | ---: | --- |
| tests/ (including test-overlay/) | 22 | 73d294085121a14df1f9412b20a541b86b9a5d516b140205f267b9b61c136b43 |
| internal/eval/ | 38 | 0feb47f776d35d0ac9f98386894f66e5b28438adb178521e718375203c959bf0 |
| patches/ | 18 | f92e2e878dfa4a2ddc665c992c339bf3331c0b540e7ae0320f2b9a74870d634b |
| compatibility/ | 1 | 3ae13c68f7b784d046339f0b38436e5d0214a8d0caa03b644f029cf04dd795f2 |
| licenses/ | 4 | 335646704d751f999da92db44e584fdbaa20dd6ab97da91b6e5b8d18fe75c3f6 |

Retain scripts/test*, functional sources, root license, docs/images/example.png,
all Trellis plans/research/history/journals/session records and unrelated
.devin/.codex configuration. Removing .codex/hooks/__pycache__ removes bytecode,
not hook scripts/configuration.

## Evidence availability

After cleanup, old raw dist and /tmp evidence paths are historical and no longer
available. Their preserved Trellis reports record the actual prior results,
including the failed upstream race gate and platform/live limitations. Do not
claim that raw logs, binary outputs, provenance or disposable fixture credentials
remain available. The final archive remains the only retained dist artifact.

## Actual completion and deletion-hook limit

The permitted deletion command completed successfully (exit 0):

```sh
find dist/sub2api-build -mindepth 1 -maxdepth 1 ! -name minimal-python-20261013 -exec rm -rf -- {} +
go clean -cache -testcache
```

The separate direct `rm -rf` command for the inventoried remaining dist roots,
caches and /tmp workspaces was rejected by PreToolUse. Its stated reason was
that targets lacked Git records, contained uncommitted changes, or could not
have their Git state verified. It explicitly required `/usr/bin/trash` instead.
The implementer stopped direct deletion and informed main. Main instructed the
implementer to follow the hook, without bypassing it or emptying all of Trash.

Every remaining inventoried target above was then moved with `/usr/bin/trash`,
using individually enumerated absolute paths (no broad shared-/tmp wildcard).
Both Trash commands completed with exit 0. These items remain recoverable in
macOS Trash and their occupied disk blocks are not reclaimed yet. No permanent
Trash deletion, alternate delete API or global Trash emptying was attempted.

The enumerated /tmp files include:

- 001-codex-eval-auth.patch.before-comment-strip, .before-quality,
  .before-revoke-lock and .before-settings-fingerprint.
- 003-codex-eval-cancellation.patch and 003-orig-backup.patch.
- codex-eval-core-race.log and codex-eval-core-vet.log.
- codex-eval-maintained-check-0.log through -5.log and
  codex-eval-replay-check-0.log through -5.log.
- codex-eval-root-race-final.log, -root-race.log, -root-vet-final.log,
  -root-vet.log, -runtime-race-final.log, -runtime-race.log,
  -runtime-tests.log, -service-unit-vet.log and -service-vet.log.

Post-cleanup allocated space:

| Root | Before KiB | After KiB |
| --- | ---: | ---: |
| dist/ | 49,320,744 | 71,680 |
| Go build/test cache | 4,162,716 | 8 |

APFS available space increased from 9,874,512 to 48,138,924 KiB: a real observed
net gain of 38,264,412 KiB (36.49 GiB). Directory reductions are not identical to
physical reclaimed blocks; moved Trash items and APFS accounting explain why a
47 GiB dist reduction is not a 47 GiB free-space claim.

Post-cleanup checks (without rebuilding/testing):

- All five source-tree fingerprints/counts match the before table exactly.
  scripts/test-sub2api.sh, scripts/test_sub2api_build.py and scripts/test-ui.cjs
  remain. Trellis reports/plans/history and unrelated configuration remain.
- The retained archive remains 73,379,445 bytes with the expected SHA-256 and
  exactly eight regular files. Its install.py bytes match current source
  (installer SHA c14d8c69277550ed501477edc40dfabce49fac9c6ce966c0d7b0fdb1ee5713a3).
- No inspected task-specific /tmp/TMPDIR candidates remain. The 52 fragment
  files and both replay workspace roots are absent from their original paths.
- Task-owned Docker resources were already absent; unrelated Docker containers,
  networks, volumes and official images were not removed.
- `git diff --check` passed. No product/source change was made for cleanup.

Finder recreated two .DS_Store files during inspection, and tool-run Trellis
hooks recreated .trellis/scripts/common/__pycache__. These are harmless generated
metadata/bytecode, not additional payloads or maintained sources. They are moved
again in the final cleanup call, but may regenerate on subsequent Finder/tool
activity; do not rerun tests or loop indefinitely to chase automatic caches.

The workspace cleanup is complete within the hook's deletion policy. Remaining
permanent reclamation of the specifically moved Trash items is a hook-limited
follow-up. Main owns independent verification and final user reporting.
