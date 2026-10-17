# U17 Status Copy Removal

The page no longer shows saving or unsaved-change captions.
Successful saves still clear earlier save errors.
The account enable switch retains its disabled and busy states without a saving tooltip.
Autosave keeps its debounce, serialized writes, field revisions, and secret handling.

Candy no longer renders per-account inspection summaries.
The summary function, row parameter, markup, and CSS rules are removed.
Actual account evaluation and notification errors use the existing inspection error location.
The page renders error text with `textContent` and retains the account identity.
Existing Candy results and history remain unchanged.

The completed phase has no caption or stale account-progress text.
An empty inspection card is hidden.
Active progress, recovery progress, cancellation feedback, and the next scheduled time remain available.
Actual runtime, load, save, validation, and action errors remain visible.
The two inspection buttons remain at the bottom of settings, outside collapsed fields.

Changed product files:

- `native/backend/internal/codexeval/web/inspection.js`
- `native/backend/internal/codexeval/web/candy.js`
- `native/backend/internal/codexeval/web/components.js`
- `native/backend/internal/codexeval/web/ui.html`
- `native/backend/internal/codexeval/web/style.css`

The existing `tests/ui/sub2api-dropdowns.cjs` assertions now check silent saves, hidden empty completion state, and visible actual errors.
No backend code, inspection flow, or test suite was added.

Passed checks:

- `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`.
- `node --check` for the three changed JavaScript assets and the existing UI check.
- `git diff --check`.

Independent source review passed in `U17-status-copy-check.md`.
The approved build and package commands then passed:

```sh
python3 scripts/build-sub2api.py --source /tmp/sub2api-codex-eval-v0217-check
python3 scripts/package-sub2api.py
```

The build reused existing source and caches without `--check` or repeated test commands.
The required native frontend build completed. Vite finished in 15.50 seconds.
The logged Go command targets the static Linux/amd64 release with `main.Version=0.2.17`.
The log is `/tmp/sub2api-U17-build.log`.

| Output | Size |
| --- | --- |
| `dist/build/sub2api-codex-eval` | 125599904 bytes |
| `dist/build/sub2api-codex-eval.tar.gz` | 38701611 bytes |

No Docker or browser operation ran in this implementation turn.
The main session owns retained-container acceptance.
