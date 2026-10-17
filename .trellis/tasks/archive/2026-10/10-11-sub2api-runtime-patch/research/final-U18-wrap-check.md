# Final U18 Wrap Check

## Findings (fixed)

None. The reviewed v0.2.17 candidate matches the approved native evaluation,
inspection, persistence, delivery, and UI contracts. No mechanical source,
syntax, or whitespace defect required a product edit.

The single-container contract is implemented by
`tests/installer/native_inspection_fixture.py`. It inspects the fixed
`sub2api-codex-eval-test` name and fails when it is absent. It does not create
another fixture container. The README, AGENTS instructions, delivery spec,
and task R9 state the same rule.

## Findings (not fixed)

- The N150/iStoreOS host and production HTTPS/proxy path were not exercised.
  The local Linux/amd64 build and retained-container checks do not prove those
  environments.
- The historical `context canceled` report has no identified triggering
  operation. Existing lifetime evidence covers normal request completion and
  explicit stop cancellation, but does not resolve that historical event.

These are recorded scope limits, not defects introduced by U18.

## Verification

- Native focused Go tests and touched-package `go vet`: pass in the integrated
  v0.2.17 check report.
- Frontend build, focused UI check, JavaScript syntax checks, and package
  generation: pass in the U18 and integrated reports.
- `git diff --check`: pass for the reviewed work paths.
- Fixed artifacts remain `dist/build/sub2api-codex-eval` and
  `dist/build/sub2api-codex-eval.tar.gz`; the archive contains only the
  executable, installer, and required license texts.
- Browser verification passed in the retained `sub2api-codex-eval-test`
  container at the nonstandard port, including native sidebar/iframe,
  settings, inspection start/stop, and Candy header status.
- The main session's read-only Docker inventory found exactly one task-owned
  container, `sub2api-codex-eval-test` (`bd6c9583c3e6`), running. Eleven
  earlier task containers were removed. Unrelated containers, images,
  volumes, caches, and application data were preserved.

No race, coverage, full-suite, cross-platform, paid-provider, production, or
additional cleanup gate was added.
