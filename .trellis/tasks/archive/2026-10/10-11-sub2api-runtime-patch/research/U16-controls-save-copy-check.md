# U16 Inspection Controls and Save Status Check

## Findings fixed

No additional product fix was required during this independent review.

The reviewed implementation satisfies both U16 corrections:

- `ui.html` has one `立即巡检` button and one `停止当前巡检` button. Both
  are at the end of the settings form and outside the collapsible automatic
  inspection and priority fields.
- `inspection.js` clears the transient status after a save with no newer edit.
  The existing `有未保存的修改` status remains when a newer edit exists.
  Validation, write, and action errors still use the existing visible status
  elements.
- Candy keeps inspection progress, results, history, and runtime errors.

## Findings not fixed

None within the U16 scope.

No backend, persistence, scheduling, container, build, browser, or deployment
work was performed. U15 backend and timer evidence remains applicable.

## Verification

- UI check: pass
  - `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`
  - Covers control placement, collapsed-field behavior, manual start/stop,
    successful-save status clearing, newer-edit preservation, and save errors.
- JavaScript syntax: pass
  - `node --check native/backend/internal/codexeval/web/inspection.js`
  - `node --check tests/ui/sub2api-dropdowns.cjs`
- Diff whitespace: pass (`git diff --check`).

No additional test suite or acceptance gate was added.
