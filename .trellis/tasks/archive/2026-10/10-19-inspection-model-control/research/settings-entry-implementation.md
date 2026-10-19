# Settings-entry repair

## Cause and boundary

Patch `0006` adds `backend/internal/codexeval/web/app.js` and owns the page startup.
Its initial `switchTab()` restored the saved Settings tab and started `loadInspectionSettings()`.
The later `start()` called that function again while `inspectionSettingsLoading` was true.
That second call returned immediately, so `start()` loaded the catalog before the first settings response.

The catalog captured `selectedGroup = 0` before its group request.
The settings response then set the saved group to 202.
The catalog correctly rejected its response because the selected group had changed.
The page had no account inventory, and ordinary tab changes did not load that inventory.

The repair gives `switchTab()` an optional `loadSettings` parameter, with a default of true.
The initial tab restoration passes false.
`start()` now owns the first settings read and waits for it before loading the catalog.
Later Settings entries still read native settings through the existing path.
The repair changes three lines in `app.js`.

## Regression evidence

The existing harness now has an optional full-page initialization path.
This path executes the real event registration, saved-tab restoration, and authenticated `start()` call.
The case delays settings and group responses to expose the original startup order.
It checks the saved Settings tab, group options, selected group, and account rows across all four evaluation tabs.
It also checks one settings read at startup and a new read on later Settings entry.

The regression failed before the repair against `/tmp/sub2api-settings-entry-before/source`.
The group selector contained only `分组 202（不可用）` instead of A and B.
The final focused check passed:

```sh
node tests/ui/sub2api-dropdowns.cjs /tmp/sub2api-settings-entry-final/source
```

The complete existing UI harness passed, including automatic saves, shared selection, inventory failures, history, and model controls.
The effective patched tree is `/tmp/sub2api-settings-entry-final/source`.
The build entry's `prepare_tree()` created both trees from the original source and current ordered patches.
The original source is `/tmp/sub2api-codex-eval-v0217-check`.

## Actual browser evidence

Chrome exercised the actual native Vue sidebar, `/admin/codex-eval` route, and same-origin evaluation iframe.
The preview used the effective tree's complete page scripts and initialization.
Only API responses and the administrator session used synthetic fixture data.
The settings endpoint delayed its response by 600 milliseconds.
The group endpoint delayed its response by 150 milliseconds.

The browser started on native `/admin/announcements` with Settings saved as the evaluation tab.
The native `账户评估` sidebar link opened the Settings tab.
Switches to Candy, Fingerprint, ModelTrace, and Pelican showed group B and accounts `B only` and `shared`.
The group selector also contained group A.
The Candy page retained model `saved/literal` and four runs from browser preferences.

The browser then selected Settings, left through the native sidebar, and entered account evaluation again.
Settings remained selected on reentry.
A direct switch to Candy loaded the same group and accounts without another tab switch or manual refresh.
The screenshot showed the native sidebar, selected group B, saved model, four runs, and both account rows.
The fixture request log recorded settings before groups and group-filtered account reads on each entry.
It recorded no provider start request.

The preview remains at `http://127.0.0.1:60231/fixture-bootstrap`.
Its temporary entry is `/tmp/sub2api-settings-entry-preview.cjs`.
The read-only request log is `/fixture-evidence` on the same origin.
The retained Docker fixture was absent, so browser checks used synthetic APIs instead of a running native backend.
This check exercised the reported browser lifecycle without production data or provider quota.

## Delivery and checks

The normal delivery entry passed after the production repair:

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check
```

The build passed its existing frontend checks and produced the patched native frontend, server, and installer.
Both fixed outputs are static Linux/amd64 ELF executables:

- `dist/build/sub2api-codex-eval`
- `dist/build/sub2api-codex-eval-installer`

`git diff --check` and `git diff --cached --check` passed.
Backend and installer behavior did not change, so their completed focused checks remain valid.
The pending degraded-tag edits remain intact in the shared page patch and UI harness.
The main session owns independent review, specification synchronization, staging, and acceptance.
