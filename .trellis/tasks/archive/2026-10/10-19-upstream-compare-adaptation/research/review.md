# Final Review

## Findings Fixed

`patches/0001-evaluation-core.patch` used `os.WriteFile`, which keeps an existing file's mode.
The loader now applies mode `0600` before writing. A permission error blocks migration and later saves.
The existing schema-1 fixture retries with a `0644` copy and checks the final mode.

## Findings Not Fixed

None. The approved variant, extraction, saved-label, migration, recovery, page, and documentation contracts match the effective candidate.

The runtime test checks failed extraction and its admitted variant directly.
It also checks saved success and cancellation through the same `runPelican` persistence path.
That path stores the complete evaluator result for every status.
The UI harness checks saved labels for failed and cancelled records.
A separate failed-run sequence would repeat the shared retry and persistence path.

## Verification

- Lint: Pass. Changed Go formatting and `go vet` passed.
- TypeCheck: Pass. The frontend build and Go server builds passed.
- Tests: Pass. The standard build/check entry passed. The affected runtime test passed after the permission fix.
- Browser: Pass. The main session checked normal and narrow layouts, both previews, animation, and detail close behavior. See [browser check](browser-check.md).
- Build evidence: `/tmp/pelican-variants-build.log`, `/tmp/pelican-variants-layout-build.log`, `/tmp/pelican-variants-copy-build.log`, and `/tmp/pelican-variants-permissions-build.log`.
- Repair evidence: `/tmp/pelican-variants-permissions-check.log`.
- `git diff --check`: Pass.
