# Independent U12 Settings Collapse Check

## Source review

- Automatic inspection collapse: `#inspection-enabled` controls
  `#inspection-auto-fields`. `updateInspectionVisibility()` runs after settings
  load and on form changes. With automatic inspection off, the fields are
  hidden; enabling it reveals them.
- Switch names: the two header switches have no visible caption. Their
  accessible names are `自动巡检` and `自动调整优先级`, provided by `aria-label`.
  Each also points to its controlled field region with `aria-controls`.
- Saved values: collapsing does not clear or disable fields. The submit handler
  serializes every named `form.elements` entry, including hidden fields. The UI
  test verifies that the model, wait range, and test count survive collapse and
  appear in the save payload and remain after the successful save response
  refreshes the form.
- Manual controls and feedback: start/stop buttons, runtime phase, next-run
  time, progress, runtime error, and results are outside
  `#inspection-auto-fields`. They remain available when automatic inspection
  is disabled. The UI test starts and stops a run with the switch off and checks
  that progress and completion remain visible while the configuration stays
  collapsed.
- U11 priority fixes remain present: normal priority defaults to 10 and has
  `min="1"`; priority fields follow the independent adjustment switch. Browser
  pair validation runs only when adjustment is enabled, so hidden invalid
  legacy values reach the API and any API validation error appears in the save
  status. The settings service still normalizes a stored nonpositive normal
  priority to 10 without writing it, preserves the stored degraded value, and
  validates initial settings before the coordinator can schedule or start.
  U11's focused service, coordinator, and UI regression checks cover these
  paths.

No U12 defect was found in the reviewed source.

## Findings (fixed)

No product code change was needed in this independent review.

## Findings (not fixed)

No U12 finding remains open.

## Work record

Removed three empty task-owned Docker networks:
`codexeval-inspection-0246d0c8`, `codexeval-inspection-b1f09e84`, and
`codexeval-inspection-c1fc2bc4`. Preserved the sole reusable task container
`sub2api-codex-eval-test` and its active network. No unrelated Docker resources
were changed.

## Verification

- UI checks: pass. `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node
  tests/ui/sub2api-dropdowns.cjs`
- JavaScript syntax: pass. `node --check
  native/backend/internal/codexeval/web/inspection.js` and the UI test.
- Whitespace: pass. `git diff --check`.
- Fixed delivery artifacts exist at `dist/build/sub2api-codex-eval` and
  `dist/build/sub2api-codex-eval.tar.gz`.
- Docker inventory after cleanup: one task-owned container,
  `sub2api-codex-eval-test`; one task-owned network in use by that container.

