# U19 Directory Organization Check

## Current Result: Trash Layout

The complete `.trellis/tasks/10-11-sub2api-runtime-patch/history/` directory is
absent from the repository after the user-authorized `/usr/bin/trash` action.
No replacement archive was created. Current task artifacts and research remain
under the task directory.

Current PRD, design, implementation plan, task metadata, AGENTS instructions,
directory spec, and quality guide describe the Trash decision. A current-path
scan found no dangling `history/retired` or relocation-inventory pointer in
those active records. The old links remain only inside the explicitly marked
superseded intermediate section of the implementation report and in this
report's preserved review below.

## Current Findings (fixed)

- File: `.trellis/tasks/10-11-sub2api-runtime-patch/research/U19-directory-organization-implementation.md`
- Issue: The rejected intermediate move report had a normal current-looking
  heading and contained links to the removed `history/retired` directory.
- Fix: Marked that section as superseded evidence and stated that its links no
  longer identify repository files.

## Current Findings (not fixed)

None. The discarded archive contents are in macOS Trash by user instruction;
no repository replacement is required.

## Current Verification

- `test ! -e .trellis/tasks/10-11-sub2api-runtime-patch/history`: pass.
- Active docs and task metadata contain no current `history/retired` or
  relocation-inventory path.
- Active native paths, root installer, current build/test/package scripts,
  retained fixture files, original CPA sources, fixed outputs, and runtime
  behavior remain unchanged.
- No product build, broad test suite, race check, Docker operation, staging,
  commit, or deployment was run.

## Superseded Intermediate Review (before the Trash correction)

- File: `.trellis/tasks/10-11-sub2api-runtime-patch/history/retired/README.md`
  and `research/U19-directory-inventory.md`
- Issue: The final classification added `tests/go.mod`, the retired browser
  helper, the obsolete fixture seeder, and the pre-U19 quality guide after the
  initial inventory. The inventory still described only the earlier 13
  installer files, and the quality-guide row was not standard code-formatted.
- Fix: Normalized the quality-guide mapping and added a clear final-relocation
  note. The final README now parses as 42 mappings, and the inventory explains
  why the initial counts differ from the final map.

## Superseded Findings (not fixed)

None within U19. Historical copies retain their old imports and root-relative
paths by design. They are preserved evidence and are not active tools.

## Superseded Verification

- All 42 relocation mappings have existing destinations; all former active
  paths are absent.
- Preserved content counts match the implementation inventory: 25 adapter
  files, 37 evaluation files, 18 v0.2.15 patch files, 7 host-overlay files,
  8 overlay test files, 28 individual helper/test/module files, and 9 retired
  specifications. This is 132 source/document files, excluding Finder
  metadata and the new inventory README.
- Active tree search found no current build, test, installer, release, or CPA
  entry point that references the retired owners. Remaining matches are the
  intended relocation links and current native checks.
- Relative Markdown link check: pass, 35 links checked.
- Python AST, shell syntax, Node syntax, ignore rules, and `git diff --check`
  results from the U19 implementation remain applicable.
- Active native paths, root `install.py`, current build/test/package scripts,
  retained fixture files, original CPA sources, fixed output paths, and runtime
  behavior were not changed by U19.
- No product build, broad test suite, race check, Docker operation, staging,
  commit, or deployment was run for this organization review.
