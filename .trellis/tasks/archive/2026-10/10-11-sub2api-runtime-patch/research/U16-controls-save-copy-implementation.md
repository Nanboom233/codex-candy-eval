# U16 Inspection Controls and Save Status

The two inspection buttons now belong to the settings tab again.
They remain outside the collapsed automatic-inspection fields.
Candy retains account progress, results, history, and runtime feedback.
The settings tab still has no account list or save button.

Successful automatic saves clear their transient status.
They do not display `已保存`.
Pending newer edits still display their existing unsaved status.
Validation, write, and action failures retain their existing error display.
Persistence, field revisions, scheduling, and inspection flow are unchanged.

Changed files:

- `native/backend/internal/codexeval/web/ui.html`: restore the two settings buttons.
- `native/backend/internal/codexeval/web/inspection.js`: clear the successful-save status.
- `tests/ui/sub2api-dropdowns.cjs`: update the existing control-location and save-status assertions.

Passed checks:

- `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`.
- `node --check native/backend/internal/codexeval/web/inspection.js`.
- `node --check tests/ui/sub2api-dropdowns.cjs`.
- `git diff --check`.

No new tests or suites were added.
Independent review passed in `U16-controls-save-copy-check.md`.

The approved build and packaging commands then passed:

```sh
python3 scripts/build-sub2api.py --source /tmp/sub2api-codex-eval-v0217-check
python3 scripts/package-sub2api.py
```

The build reused existing source and caches without `--check` or repeated backend suites.
Its required frontend build passed. Vite completed in 18.97 seconds.
The build log is `/tmp/sub2api-U16-build.log`.

| Output | Result |
| --- | --- |
| `dist/build/sub2api-codex-eval` | 125599904 bytes. Static, stripped ELF 64-bit x86-64. |
| `dist/build/sub2api-codex-eval.tar.gz` | 38701694 bytes. Seven expected entries. |

The build command sets the native version to `0.2.17`.
The archive contains the executable, installer, and five license texts.
The main session owns final browser verification in the retained container.
