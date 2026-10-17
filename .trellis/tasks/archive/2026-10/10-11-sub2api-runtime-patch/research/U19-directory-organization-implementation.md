# U19 Directory Organization Implementation

Date: 2026-10-18

## Current user correction: use Trash

The user rejected the task history directory because it does not follow the task layout.
The implementation owner moved the complete directory to macOS Trash with `/usr/bin/trash`.
The repository path `.trellis/tasks/10-11-sub2api-runtime-patch/history/` no longer exists.
This action includes the retired implementation files and older task snapshots.
No replacement archive was created. Current task artifacts and research remain in place.
Current task records and the work-commit plan now describe the Trash decision.
The original move report below records the rejected intermediate state. Its
relocation links no longer identify repository files and are retained only as
superseded evidence.

## Superseded intermediate attempt (before the Trash correction)

Organize the repository before work commits without changing runtime behavior.
Keep the active native integration, public commands, fixed artifacts, current fixture, and original CPA sources at their existing paths.
Move confirmed abandoned attempts into this task's `history/retired/` directory.
Do not stage, commit, rebuild, deploy, or run Docker acceptance in this step.

## Preservation moves

The [relocation inventory](../history/retired/README.md) records 42 path mappings.
Each directory mapping includes all descendants at their prior relative paths.
Direct `mv` commands preserved the file contents. No historical code or report was rewritten during relocation.

| Preserved directory | Source files, excluding Finder metadata |
| --- | ---: |
| `history/retired/adapter/` | 25 |
| `history/retired/internal/eval/` | 37 |
| `history/retired/patches/sub2api/v0.2.15/` | 18 |
| `history/retired/tests/host-overlay/` | 7 |
| `history/retired/tests/overlays/codexeval/` | 8 |

The other mappings preserve the compatibility profile, 11 helpers, 14 installer helpers, one browser helper, and the overlay module file.
Eight retired specifications moved from the active backend spec directory.
The full pre-U19 quality guide also moved into history before its focused replacement was added.
The total is 132 source and document files, excluding Finder metadata and the new relocation inventory.

`tests/go.mod` declares `codex-eval-development-overlays` and only isolates the retired overlay tests.
The current fixture and UI checks do not use it.
`seed_browser_fixture.py` uses the removed `installer.Installer` class and old fixture containers.
`sub2api-browser.js` uses `/codex-eval/ui` and four tabs. The current page uses `/codex-eval` and five tabs.
Neither helper has a current caller. Both moved to history without a replacement suite.

The [research inventory](./U19-directory-inventory.md#u19-implementation-classification-correction) now records these classification corrections.
Its earlier findings remain visible as historical research.

## Active document and ignore changes

- `.gitignore` ignores `__pycache__/`, Python bytecode, and root `output/`.
- `AGENTS.md` identifies current native owners, command scripts, fixed outputs, retained tests, and retired history.
- `AGENTS.md` separates current Sub2API checks from existing CPA-only build rules.
- `directory-structure.md` describes retained test paths and links to the relocation inventory.
- `quality-guidelines.md` removes retired mechanisms from active instructions and scopes retained CPA rules explicitly.
- The active quality guide retains the complete hybrid commit section without text changes.

Generated Python caches, browser output, and fixed build outputs remain on disk.
The repository hook rejected `rmdir` for untracked directories and directed use of macOS Trash.
The main session authorized that route.
The empty `compatibility/` and `tests/overlays/` directories went to Trash.
The residual `patches/` tree contained only Finder `.DS_Store` files and also went to Trash.
These residual directories contained no remaining source files.
Finder metadata inside the preserved directory moves remains with those historical copies.

## Checks and evidence

| Check | Result |
| --- | --- |
| Relocation destinations and absent former paths | Passed for all 42 mappings. The active quality guide is the documented replacement exception. |
| Relative Markdown targets in five changed documents | Passed. Includes current specs, inventory, relocation map, and `AGENTS.md`. |
| Hybrid commit section compared with preserved pre-U19 guide | Exact text match |
| Active-reference search outside task history and generated output | No retired owner references remain. Two matches were unrelated `storage-adapter` examples in Trellis skill documentation. |
| Python AST parse | Passed for root installer, two current build/package scripts, and three current fixture/version Python files. |
| Shell syntax | Passed for build wrapper, native fixture start script, and current test entry. |
| Node syntax | `node --check tests/ui/sub2api-dropdowns.cjs` passed. |
| Git ignore checks | Python caches, bytecode, root browser output, and fixed build output match ignore rules. |
| Unchanged active paths compared with the existing index | No new differences in native source, installer, four current scripts, current fixture/UI tests, original CPA owners, docs, README, or release workflow. |
| `git diff --check` | Passed |

The fixed executable remains 125599904 bytes.
The fixed archive remains 38700695 bytes.
This step did not calculate hashes or change either artifact.
The U18 build and browser evidence remains applicable because current product inputs did not change.

## Follow-up and limits

The main session owns task-manifest synchronization, independent review, final spec integration, and commit preparation.
Retired scripts retain old imports and root-relative paths. They are historical copies and need not run from their new location.
Historical reports retain former file paths. The relocation inventory resolves those paths without rewriting the reports.
Existing private and unrelated local changes remain untouched.
The pre-existing `install.sh` and `install.ps1` deletions were not part of U19.
No new product tests, broad suites, race checks, provider calls, Docker checks, or deployment ran.
This organization does not resolve or reinterpret the earlier history cancellation report.

## Main-session integration

Independent review passed in `U19-directory-organization-check.md` after its inventory correction.
Main synchronized the PRD, design, execution plan, context manifests, task status, and candidate work-commit plan.
The candidate plan contains 81 unique existing work files. The index still contains the earlier 80-file proposal.
No additional staging or commit occurred.
Task context validation and the final worktree whitespace check passed.
Context validation warns that the existing 34300-byte approval ledger exceeds the 32768-byte automatic injection limit.
The full ledger remains available for direct reading. Current task artifacts record the later decisions.
