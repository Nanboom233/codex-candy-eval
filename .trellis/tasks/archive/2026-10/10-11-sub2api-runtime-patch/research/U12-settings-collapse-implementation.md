# U12 Settings Collapse Implementation

## Scope

The automatic inspection switch now controls the visibility of its configuration body.
The body starts hidden and follows loaded settings and local switch changes.
Closing the body retains all input values and includes them in explicit saves.

Both card header switches have accessible names without visible captions.
Manual start and stop now share the existing settings action row.
Runtime progress and errors appear above that row, outside the collapsed configuration.
The existing results section remains outside the configuration.
No explanatory card, paragraph, tooltip, or new backend setting was added.

## Preserved behavior

The priority switch still controls both priority inputs.
Normal priority still defaults to 10 and has minimum 1.
All five model dropdowns retain suggestions and custom literal values.
Free filtering, native account availability, and independent model discovery remain unchanged.
Secret redaction and blank-secret saves remain unchanged.
Removed account patterns, idle copy, and notification configured-status copy remain absent.

The checker fixes for invalid loaded priority pairs and hidden-field validation remain intact.
This change does not modify backend scheduling, persistence, or account tests.
Reuse the recorded U11 backend checks and independent review evidence.

## Changed files

- `native/backend/internal/codexeval/web/ui.html`: added the collapsible configuration body and accessible switch names. Moved runtime feedback and manual actions.
- `native/backend/internal/codexeval/web/inspection.js`: updates both configuration bodies after settings loads and form changes.
- `native/backend/internal/codexeval/web/style.css`: removed the obsolete manual-action row selectors.
- `tests/ui/sub2api-dropdowns.cjs`: checks collapse, retained values, save payloads, accessible names, and manual start/stop with scheduling off.

## Checks

- Focused settings and model UI checks passed.
- JavaScript syntax checks passed.
- Diff whitespace checks passed.
- The native frontend passed its three locale tests, TypeScript checks, and Vite build.
- The native Linux/amd64 build passed through `scripts/build-sub2api.py`.
- Packaging passed through `scripts/package-sub2api.py`.

The executable is a stripped, static Linux ELF for amd64.

| Output | Size |
| --- | ---: |
| `dist/build/sub2api-codex-eval` | 125587616 bytes |
| `dist/build/sub2api-codex-eval.tar.gz` | 38700463 bytes |

The archive has seven entries: the executable, the installer, and license files.
The archive contains no development or test files.

Build source: `/tmp/sub2api-codex-eval-v0217-check`.
Build log: `/tmp/sub2api-U12-build.log`.
The main session owns container loading, actual browser verification, and independent review.
These final U12 checks remain pending at this handoff.
No Docker, production, paid-request, staging, commit, or cleanup action ran.
