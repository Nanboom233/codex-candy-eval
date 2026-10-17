# U18 Candy Header Status

The existing Candy account header now contains the next inspection time beside the account title and enabled-only control.
The same header retains active inspection feedback and progress.
The separate runtime card below the account list is removed.
Its card-visibility helper and calls are removed.

Actual runtime, evaluation, notification, and settings errors retain their existing error nodes inside the header.
Empty status and error nodes stay hidden.
No table column, card, explanatory caption, or replacement completion caption was added.
Settings-bottom inspection buttons, silent autosave, and independent inspection flow remain unchanged.

The Candy header can wrap at all widths.
At the existing narrow breakpoint, inspection status uses its own wrapped header line.
Long status and error text can wrap within the header.
The other three evaluation headers retain their existing structure and layout rules.
Actual browser layout checks belong to the main session after independent review and the final build.

Changed files:

- `native/backend/internal/codexeval/web/components.js`: Candy-only header status and error nodes.
- `native/backend/internal/codexeval/web/ui.html`: remove the lower runtime card.
- `native/backend/internal/codexeval/web/inspection.js`: remove the obsolete card helper and hide empty error nodes directly.
- `native/backend/internal/codexeval/web/style.css`: header status wrapping and actual error presentation.
- `tests/ui/sub2api-dropdowns.cjs`: update existing header placement and completion-state assertions.

Passed checks:

- `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`.
- `node --check` for changed JavaScript assets and the existing UI check.
- `git diff --check`.

Independent source review passed in `U18-header-status-check.md`.
The approved build and packaging commands then passed:

```sh
python3 scripts/build-sub2api.py --source /tmp/sub2api-codex-eval-v0217-check
python3 scripts/package-sub2api.py
```

The build reused existing dependencies and caches without `--check` or repeated test commands.
The required frontend build passed. Vite completed in 19.27 seconds.
The logged Go command sets the release version to `0.2.17`.
The build log is `/tmp/sub2api-U18-build.log`.

| Output | Size |
| --- | --- |
| `dist/build/sub2api-codex-eval` | 125599904 bytes |
| `dist/build/sub2api-codex-eval.tar.gz` | 38700695 bytes |

No backend work, additional tests, Docker, or browser operation ran in this implementation turn.
The main session owns retained-container browser acceptance.
