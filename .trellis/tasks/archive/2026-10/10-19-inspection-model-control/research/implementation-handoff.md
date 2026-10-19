# Implementation record

## Change boundary

Add confirmed-transition model control to the existing inspection settings, native account update, and embedded form.
Patch `0005` owns settings and the account helper. Patch `0006` owns the widget and automatic save integration.
Patch `0007` preserves the pending removal list during ordinary account edits.
Extend existing tests under `tests/` for the changed contracts.
Keep scheduling, direct sending, native availability, unrelated mappings, and pending degraded-tag changes unchanged.
No new transport, scheduler state, repository interface, or migration is necessary.

## Evidence

The implementation adds settings defaults, model normalization, transition settings, native compound updates, and pending-list preservation.
The page uses the native selector layout instead of the earlier inline-input design.

- `python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check` passed.
  The entry passed frontend checks, the UI harness, focused Go tests, vet, and the Linux build.
- The final UI revision changed only patch `0006` after that check.
  `node tests/ui/sub2api-dropdowns.cjs /tmp/sub2api-model-control-final` passed for that revision.
- `python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check` passed for the final UI revision.
  It refreshed `dist/build/sub2api-codex-eval` and `dist/build/sub2api-codex-eval.tar.gz`.
- `gofmt -l` returned no changed production Go files. `git diff --check` passed.

## Native UI comparison

The reference mounts the original Vue component with four synthetic model names.
Its source is `frontend/src/components/account/ModelWhitelistSelector.vue` in the unchanged upstream tree.

The final widget copies the two-column grid, gray tags, icon sizes, remove controls, border, padding, dropdown placement, and custom-entry row.
The custom input sits below the selector. Its visible button uses the native Chinese label `填入`.
The dropdown overlays later content. Its search area and 208-pixel option limit match the native component.
Counts, explanations, copy controls, and quick actions are absent.

Visual evidence:

- Native reference: `.playwright-cli/page-2026-10-18T19-36-11-008Z.png`.
- Final desktop: `.playwright-cli/page-2026-10-18T19-34-12-091Z.png`.
- Final narrow dropdown: `.playwright-cli/page-2026-10-18T19-39-04-484Z.png`.

The browser inspection exercised the switch, dropdown, search, selection, removal, custom entry, and automatic saves.
Tab reached the option and custom input. Enter opened the dropdown and added a model. Space selected an option.
Escape returned focus to the trigger. An outside click closed the dropdown.
At 390 pixels, tags retained two columns and the document had no horizontal overflow.
The synthetic settings endpoint confirmed automatic persistence after selection and removal.

## Handoff

The effective patched tree is `/tmp/sub2api-model-control-final`.
The synthetic preview is `http://127.0.0.1:60219`.
Run `node /tmp/sub2api-model-control-preview.cjs` if that preview stops.
The preview assembles actual page files in memory. It uses synthetic settings and sends no provider requests.

The native Vue reference is `http://127.0.0.1:60220/native-reference.html`.
Its temporary root is `/tmp/sub2api-model-control-implementation/frontend`.
Run `./node_modules/.bin/vite --config native-reference.config.mjs` from that root if the reference stops.
Both preview services remain available for review.

The main session owns specification synchronization and final acceptance.
An independent check remains the next workflow step.
No installation, release, commit, or real-account operation occurred.

## Native brand icon correction

The independent check found that non-OpenAI models used letters instead of native brand icons.
Patch `0006` now copies all native `ModelIcon.vue` SVG data and model-name classification rules.
The page retains the native fallback only for model names without a known brand.
This correction adds no dependencies or new assets.

The existing UI harness passed against `/tmp/sub2api-model-control-final` after this correction.
The normal build entry passed and refreshed the final executable and archive.
Unchanged Go check evidence remains valid.

The browser inspection added Claude, Gemini, GLM, Qwen, and DeepSeek names through the custom-entry control.
The tags displayed the native shapes and colors without fallback letters.
Screenshot: `.playwright-cli/page-2026-10-18T19-55-48-517Z.png`.
The effective tree and synthetic preview above contain this final correction.
