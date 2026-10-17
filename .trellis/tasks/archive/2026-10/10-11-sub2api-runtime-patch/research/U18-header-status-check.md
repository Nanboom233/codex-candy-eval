# U18 Candy Header Status Check

## Findings fixed

No additional product fix was required during this independent review.

The reviewed source satisfies U18:

- The Candy account header owns the inspection status and next-run timestamp
  beside the account title and `仅启用` control.
- The other three evaluation headers do not receive inspection status nodes.
- The separate lower runtime card and its visibility adapter are removed.
- Active progress, cancellation feedback, next-run time, and actual runtime,
  evaluation, notification, and settings errors remain available.
- Empty status and error nodes remain hidden. Completed state does not add a
  caption or stale progress.
- No new table column, card, explanatory caption, or replacement completion
  copy was added.
- Settings-bottom start/stop controls and U17 silent-save behavior remain
  unchanged.
- The Candy header can wrap, and the narrow breakpoint gives the inspection
  status its own full-width line. Long status and error text can wrap.

## Findings not fixed

None within the U18 scope.

No backend, persistence, scheduling, build, package, Docker, browser, or
deployment work was performed. U15 timer/backend and U17 source evidence remain
valid.

## Verification

- UI check: pass
  - `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`
  - Covers Candy-header placement, removal of the lower card, active progress,
    next-run status, completion cleanup, error forwarding, shared history, and
    settings controls.
- JavaScript syntax: pass
  - `node --check native/backend/internal/codexeval/web/inspection.js`
  - `node --check native/backend/internal/codexeval/web/components.js`
  - `node --check native/backend/internal/codexeval/web/candy.js`
  - `node --check tests/ui/sub2api-dropdowns.cjs`
- Product structure scan: pass; inspection status/error nodes are generated
  only for Candy, and the lower runtime-card identifier is absent from product
  assets.
- Diff whitespace: pass (`git diff --check`).

No additional test suite or acceptance gate was added.
