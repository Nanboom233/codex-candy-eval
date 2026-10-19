# Implementation Handoff

## Effective Source

The retained patched source is `/tmp/pelican-variants-native/source`.
Its frontend dependencies link to the retained original source installation.
The standard build log is `/tmp/pelican-variants-build.log`.

## Changes

- `patches/0001-evaluation-core.patch` adds shared variants, request normalization, saved versions, markup extraction, and schema 2 migration.
- `patches/0006-evaluation-ui.patch` adds shared page configuration, version selection, prompt updates, saved labels, and picture labels.
- `tests/backend/internal/codexeval/pelican_test.go` extends the existing extraction and code-size cases.
- `tests/backend/internal/codexeval/runtime_test.go` extends runtime, version, migration, and recovery-copy cases.
- `tests/ui/sub2api-dropdowns.cjs` extends startup, preference, prompt, payload, and saved-label cases.

The existing focused test selector includes all changed tests.
The build script does not need a change.

## Migration

The loader accepts unversioned, schema-1, and schema-2 files.
Schema 1 retains current fingerprint records. Unversioned files retain the existing retirement rule.
Older Pelican records receive `variant: original`.
Migration copies the original bytes to `history.v0.json` or `history.v1.json` beside `history.json` before replacement.
The copy enforces mode `0600` before writing, including when the version-based path already exists.
Open, permission, write, or close failure retains decoded records and blocks migration and later saves.
The existing atomic saver writes schema 2. A failed migration save retains memory records, the original file, and its copy.
An executable downgrade also needs its compatible pre-upgrade history file.

## Validation

The focused page harness passed against the retained tree.
The standard build/check entry passed with exit status 0.
It passed frontend checks/build, the UI harness, focused Go tests, vet, installer tests, and both Linux/amd64 builds.
The final runtime-fixture adjustment also passed `go test -tags=unit ./internal/codexeval -run '^TestCodexEvalRuntime$'`.
`git diff --check` and changed-Go formatting checks passed.
The fixed outputs are `dist/build/sub2api-codex-eval` and `dist/build/sub2api-codex-eval-installer`.

The browser found toolbar clipping in a 780-pixel native iframe.
The Pelican toolbar now uses two field columns and a full-width test button at widths up to 900 pixels.
Pelican latest-result metadata can wrap. The short version label does not shrink.
The affected UI harness and `git diff --check` passed after this repair.
The output rebuild log is `/tmp/pelican-variants-layout-build.log`.
Vite rebuilt the retained frontend. The Linux/amd64 server and installer rebuilt successfully with the final CSS.
The rebuild reused the standard build flags and package helper. Unchanged tests did not run again.

The main session confirmed the repaired toolbar and complete version labels at normal, 780-pixel, and narrow widths.
Generic Pelican interface labels now use `画面`, `鹈鹕画面`, and `放大查看` as approved.
The exact upstream prompts remain unchanged. The retained source contains the same label changes.
The final embedded-output rebuild log is `/tmp/pelican-variants-copy-build.log`.
The Linux/amd64 server and installer rebuilt successfully after the label correction. `git diff --check` passed.

Independent review found that `os.WriteFile` retained the mode of an existing migration-copy path.
The copy now applies `Chmod(0600)` before writing. Permission errors use the existing backup-failure path.
The existing schema-1 fixture changes the copy mode to `0644` before its retry and retains its final `0600` assertion.
The affected check log is `/tmp/pelican-variants-permissions-check.log`.
The output rebuild log is `/tmp/pelican-variants-permissions-build.log`.
`TestCodexEvalRuntime`, formatting checks, and `git diff --check` passed after the permission repair.
The Linux/amd64 server and installer rebuilt successfully. The retained source contains the final repair.

## Main-Owned Acceptance

Inspect both previews and detail animation at normal and narrow widths.
The main session completed the toolbar and version-label browser checks.
Update the Chinese README and affected native specifications with these effective contracts.
Keep the recovery copy paths in the existing downgrade instructions.
