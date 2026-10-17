# U7–U9 installer progress and disabled model discovery

Status: implementation and focused checks passed. Product/artifact inputs are
frozen for independent Trellis review. No commit or deployment was performed.

## Boundary and changes

The user confirmed installation succeeds after a transient startup reset. U7
therefore changes output only. `install.py` prints ordinary startup progress
with elapsed/deadline, container state, restart count and address. Both native
readiness checks must pass before `Application ready` appears. Final failures
retain the last actual request error. Existing state/address inspection,
deadlines, request bounds, replacement and rollback control flow are unchanged.

U8 changes the default host backup location. `tempfile.mkdtemp` creates a private
directory under explicit `/tmp`, with a container-ID prefix and random suffix.
The installer prints its exact `sub2api` path and retains it for restoration.
Explicit `--backup-dir` keeps its existing timestamped layout. No old backups,
caches or business/history data were deleted. No cleanup subsystem or download
path was added. Container-side same-directory replacement staging is unchanged.

U9 belongs in `native/backend/internal/codexeval/web/catalog.js`. Keep all accounts
in `accountList`, but query raw models only for `schedulable !== false` accounts.
Index failure notices against that filtered list. Only mark the catalog loaded
after discovery succeeds, so the next poll cannot silently clear an enabled
account's real error. A forced refresh also clears the loaded flag first.

Existing `app.js:setCredentialEnabled` already posts native schedulable and
forces `load({refresh:true})`. It refreshes suggestions/notices after a switch
change without another implementation. Account rows, enable controls, ordinary
account status/restrictions and editable model inputs remain intact.

Product files changed in this repair:

- `install.py`
- `native/backend/internal/codexeval/web/catalog.js`
- `README.md`
- `docs/sub2api-operations.md`

## Focused evidence

Evidence root:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8`.

`installer-progress-check.py`, session26109 exit0, uses a real local TCP reset.
One reset then healthy responses produced friendly progress and explicit ready,
without transport error text. Persistent resets retained the endpoint and raw
error in the final bounded failure. Logs: `installer-progress-ready.log` and
`installer-progress-failure.log`.

`disabled-model-check.cjs` exited0. It executes actual page load/catalog/toggle
functions in Node VM with synthetic API data. A disabled failing account was
never queried. An enabled failing account still produced its actual notice,
including after polling. All account rows remained. Disabling the failing
account cleared its obsolete notice through the existing forced refresh.
Reenabling restored the genuine error. All four direct model inputs retained
their entered unlisted name alongside successful raw suggestions. Result:
`disabled-model-result.json`. This is focused JavaScript behavior evidence, not
a new browser/layout claim. Prior actual Chrome origin/sidebar/iframe and model
forwarding evidence remains applicable and was not repeated.

`installer-tmp-backup-check.py`, final session33900 exit0, uses a cached local
Node Docker fixture with reset/health/setup responses and the final actual
native executable for version, copy, replacement and restoration. No image was
downloaded. The default installer created private mode0700 `/tmp` directories,
printed their exact paths and retained originals matching the executable bytes.
Normal startup became ready. A truncated candidate caused bounded failure;
restoration from the default `/tmp` backup succeeded and the restored executable
matched the original. Logs/result: `installer-tmp-fixture/`.

The first synthetic fixture used a three-second deadline while deliberately
resetting twice after each restart. Restoration readiness exceeded that fixture
deadline. Its logs remain as `short-fixture-*`. The fixture deadline became
eight seconds for the passing rerun. Product timeout/bounds were not changed.
Actual native stopped-container restoration remains covered by the prior
`native-browser-fixture/outcome.json` and was not repeated for this copy change.

Exact owned-container query after teardown:
`docker ps -a -q --filter label=codex-eval-check=w8ax6xp8` returned empty.
Only the owned disposable fixture container was removed. Backups remain intact.

Python compilation, changed JavaScript syntax and changed-file whitespace checks
passed. No Go source changed, so previously passed Go tests/vet were reused.

## Final build and package

Build session13620 exited0. Log: `installer-disabled-model-build.log`.
The ordinary builder reused original dependencies and Go caches, rebuilt its
native frontend with existing i18n/TypeScript/Vite steps and compiled the static
Linux/amd64 embed binary including the changed standalone catalog. The previous
final temporary source tree had been removed by its standard build lifecycle,
so no retained rebuilt native frontend tree was available for direct reuse.
No pristine baseline, race/full Go suite or browser matrix was added.

Package session70591 exited0. Fixed outputs:

- `dist/build/sub2api-codex-eval`
  SHA256 `23f2565f340b7c5ad9c542bb525cef5622fbebc8f50a1aa56e6fac2080290d4d`
- `dist/build/sub2api-codex-eval.tar.gz`
  SHA256 `71ce4c4b0eaa802bcfd610b6bb2e9abf431f620c8185cea81f18bec9b3ba08c0`

Archive inventory is the binary, current installer and five license texts.
The archived installer matches the current source. No test, log, development
source or cleanup payload was added.

N150/iStoreOS and production HTTPS/proxy execution remain unperformed. Local
Docker uses an ARM64 VM and amd64 execution support. No paid or production call,
private environment access or CPA deployment was performed.
