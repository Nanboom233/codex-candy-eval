# Robust preflight artifact cleanup

## Authorization and checkpoint

Status: cleanup complete, 2026-10-13. Independent review is complete without open current findings.
Main explicitly authorized cleanup of the implementation and independent-review generated roots.
The original 09:58 UTC inventory below is retained as the pre-review snapshot.
No product changes, tests or rebuilds are part of this cleanup.

Sizes, exact ownership, archive identities and maintained-source fingerprints were refreshed before effects.
The durable [execution checkpoint](robust-preflight-cleanup-checkpoint.json) records all nine exact Trash targets and 15 preservation fingerprints.
The checkpoint records actual effects and completed read-only preservation verification.

## Preserved items

Keep both archives at their original paths:

| Archive | Bytes | SHA-256 |
| --- | ---: | --- |
| `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz` | 525,565 | `f9980775d4aefe08906b9c0a480420f5454eb5e03164780b98643c4535176cfd` |
| `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz` | 73,379,445 | `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8` |

The current universal archive contains 62 regular functional/license files.
Its final functional hash is `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759`.
Release tag `v0.4.0` remains unpublished.
The historical minimal archive contains eight regular files.
These are the only current files observed under dist.

Preserve maintained product sources, tests/overlays, scripts, documentation, licenses, Trellis records and unrelated configuration.
Preserve all evidence until the checker finishes.
After approved cleanup, retained reports must distinguish historical canceled attempts from the successful final lifecycle.
The old 514,791-byte candidate was superseded by the checker's strict maintenance-parser fix.
The reviewed archive's 62 functional/license files match current maintained source bytes before cleanup.
The checker reran root checks and all 16 archived preflight tests before authorizing cleanup.

## Allocated space snapshot

Sizes use `du -sk`. They describe allocated directory space, without a permanent-reclamation claim.

| Root | KiB |
| --- | ---: |
| `/tmp/codex-robust-preflight-20261013` | 791,212 |
| `/Users/minami/Library/Caches/go-build` | 7,884,748 |
| Interrupted build temporary workspace | 1,960 |
| Default-source resolver temporary workspace | 4 |
| Repository Python cache roots below | 432 |

The first workroot reading was 791,944 KiB. Review activity makes the snapshot provisional.
APFS available space at the initial reading was 41,180,128 KiB.
Measure actual before/after available space during approved cleanup. Do not sum directory reductions into reclaimed space.

## Exact implementation workroot

The entire `/tmp/codex-robust-preflight-20261013` is implementation-generated evidence and staging.
It contains no maintained repository source or retained final archive.
Its final and baseline archive copies are redundant review evidence.
Move this whole root only after main confirms review completion.

Large groups observed during inventory:

| Generated group | KiB |
| --- | ---: |
| `current-source` | 94,740 |
| `previous-source` | 93,288 |
| `source-current` and `source-previous` unmaterialized clones | 16,812 each |
| `test-current`, `final-current`, `verified-current` | 74,688 each |
| `test-previous`, `final-previous`, `verified-previous` | 73,676 each |
| Successful fixture evidence `codexeval-fixture-d746c1cc83` | 113,908 |
| `archive-extracted` | 1,976 |
| `package-ready` | 1,896 |
| `package-initial`, `package-formatted`, `package-final` | 1,892 each |

Other children contain baseline installer copies, edit helpers, bytecode, raw logs, JSON evidence and interrupted fixture inputs.
Successful fixture evidence retains five real build logs/profiles/facts, loaded candidate and private disposable credentials.
The successful lifecycle completed 18 checkpoints with observed exit 0.
The original interrupted build evidence remains separate and is not a pass.

## Exact external temporary roots

Both targets were verified through this task's explicit recorded ownership:

- `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-build-7fdhzkrj`.
  Its build profile names `codexeval-fixture-0b435baa43-app`.
  It retains canceled pre-write build diagnostics and a duplicate functional context.
- `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-default-source-t115xe1s`.
  The recorded no-hints resolver JSON identifies this path.
  It retains only source resolver evidence.

No matching backup fragment, active build or universal-CI temporary root was observed in the inspected TMPDIR.
The inspected /tmp task-name families contained only the implementation workroot.
Do not delete a newly discovered path without checking task ownership and review completion.

## Exact repository bytecode roots

| Generated cache root | KiB |
| --- | ---: |
| `/Users/minami/Code/C/cpa-plugin-codex-candy-eval/__pycache__` | 80 |
| `/Users/minami/Code/C/cpa-plugin-codex-candy-eval/adapter/__pycache__` | 12 |
| `/Users/minami/Code/C/cpa-plugin-codex-candy-eval/scripts/__pycache__` | 80 |
| `/Users/minami/Code/C/cpa-plugin-codex-candy-eval/tests/installer/__pycache__` | 60 |
| `/Users/minami/Code/C/cpa-plugin-codex-candy-eval/.trellis/scripts/common/__pycache__` | 200 |

Only generated bytecode is in scope. Keep all source scripts and hook configuration.
Trellis hooks can recreate their bytecode during later tool calls.
Do not run tests after approved cleanup or loop indefinitely over regenerated harmless cache files.

## Docker and shared caches

No codexeval container, resolver, builder, fixture network or named data volume remains.
All five observed containers and two custom networks belong to unrelated kevbot work.
Buildx lists only the shared desktop builder and an unavailable default endpoint.
The unavailable default is an existing endpoint state, without a task-owned removal target.

Docker-wide read-only allocation snapshot:

| Type | Count | Size |
| --- | ---: | ---: |
| Images | 17 | 8.073 GB |
| Containers | 5 | 109.2 MB |
| Local volumes | 82 | 2.357 GB |
| Build cache entries | 102 | 4.477 GB |

Those totals are shared state, not task-owned cleanup targets.
Anonymous volumes cannot be safely attributed from current names alone and remain untouched.
Do not prune Docker, images, volumes, shared builders or global build cache.
Keep Go module cache and global npm/pnpm/Corepack stores.

After review completion and explicit cleanup dispatch, clear Go build/test cache only through `go clean -cache -testcache`.
Its current allocated size is 7,884,748 KiB. Do not clear it while independent checks are running.

## Planned execution boundary

Move only individually verified generated roots with `/usr/bin/trash` under the deletion-hook policy.
Do not use direct filesystem deletion, alternate delete APIs or global Trash emptying.
Trash contents remain recoverable and may retain their disk blocks.
Actual cleanup must record each effect, before/after sizes, free-space observations and source/archive preservation.
Those effects are now complete. See the execution and verification sections below.

## Refreshed execution inventory

| Target | Before KiB |
| --- | ---: |
| Implementation workroot | 795,524 |
| `/tmp/codex-robust-independent-review-20261013` | 160,148 |
| Interrupted build temporary workspace | 1,960 |
| Default-source temporary workspace | 4 |
| Five repository bytecode roots | 432 |
| Go build/test cache | 10,674,024 |
| Retained dist | 72,176 |

APFS available space at the execution checkpoint: 34,112,740 KiB.
The independent-review root contains only generated current/previous workspaces, adapter binary, package staging and archive extraction.
Main explicitly authorized that exact root together with the original inventory targets.
No task-owned Docker resource remains. No Docker cleanup effect is required.
Source, tests, scripts, documentation, active specifications and task records are fingerprinted before effects.
The two cleanup report/checkpoint files and generated bytecode/metadata are excluded from preservation fingerprints.

## Completed effects

Both commands completed with observed exit 0:

- `/usr/bin/trash` with the nine individually enumerated targets in the durable checkpoint.
- `go clean -cache -testcache`.

No direct filesystem deletion, alternate delete API, global Trash emptying or Docker cleanup occurred.
No module, npm, pnpm, Corepack or shared Docker cache was cleared.
No product source, maintained test, specification, task metadata or unrelated configuration was edited.
Only this cleanup report and its preservation checkpoint were written.
No tests, builds, deployment, paid calls, commit, push, publication or task archival occurred during cleanup.

## Read-only preservation verification

All 15 source/test/script/documentation/specification/task fingerprint groups match the pre-effect checkpoint exactly.
That includes 23 maintained test files, 24 adapter files, 19 specification files and 116 task files.
Generated bytecode, Finder metadata and the two cleanup files are excluded from that count and comparison.

Both archive sizes, regular-file counts and SHA-256 values remain unchanged from the reviewed identities above.
The universal archive remains 525,565 bytes with 62 files and SHA `f9980775d4aefe08906b9c0a480420f5454eb5e03164780b98643c4535176cfd`.
The historical minimal archive remains 73,379,445 bytes with eight files and its original SHA.
Dist contains exactly those two archives. No duplicate staging, logs, binaries or checksum sidecar remains there.

All nine original Trash-target paths are absent at verification time.
The implementation workroot, independent-review root, both exact TMPDIR roots and five bytecode roots are gone from their original locations.
Raw test/build/fixture credentials and evidence are now unavailable at those paths.
The preserved implementation and independent-review reports retain actual check results and their execution owners.
Historical canceled runs remain failures/interruption records. Cleanup does not convert them into successful acceptance.

All five unrelated kevbot containers remain. No task-owned Docker container, builder, network or named volume remains.
Anonymous volumes, every image, shared builders and Docker build cache remain untouched.

| Allocated directory space | Before KiB | After KiB |
| --- | ---: | ---: |
| Implementation workroot | 795,524 | absent, moved to Trash |
| Independent-review workroot | 160,148 | absent, moved to Trash |
| Go build/test cache | 10,674,024 | 8 |
| Dist | 72,176 | 72,176 |

APFS available space increased from 34,112,740 to 44,802,268 KiB.
The observed net increase is 10,689,528 KiB, approximately 10.19 GiB.
This is an actual free-space reading, not the sum of removed directory allocations.
Trash moves retain recoverable files and may retain their disk blocks. Concurrent filesystem activity can affect the net reading.

## Remaining boundaries

Cleanup is complete within the authorized scope and deletion-hook policy.
Permanent deletion of specific Trash items was not requested or performed. Global Trash remains untouched.
Trellis hooks or Finder may recreate harmless bytecode/metadata after verification.
No repeated cache-clearing loop or cache-generating tests followed cleanup.

Required broad upstream race/frontend/previous integration gates remain historical failures.
Native amd64 loaded lifecycle, published bootstrap and production deployment/billing/concurrency/browser acceptance remain unverified.
The task remains in progress. This cleanup does not grant publication readiness.
