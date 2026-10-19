# Go Installer Independent Check

## Findings Fixed

- The implementation's first revision registered `SIGTERM` in addition to Ctrl-C.
  The final source registers only `os.Interrupt`, matching the approved interruption behavior and the Python installer.
- The first revision rejected process arguments.
  The final source has no argument parameter, parser, or validation branch.
- The first revision extracted and validated the embedded payload twice in one entry test.
  The final test keeps the public entry transaction and the license-content check only.
- Go's default HTTP client followed redirects while the Python client treated a non-200 response as a failure.
  The final client disables redirects and the existing transaction table checks HTTP 302 restoration.
- The filtered selection fixture retained stale invalid `q` and textual `more` inputs.
  The implementer removed them and ran only `TestSelection/filtered`, which passed in 0.947 seconds.
  This test-only cleanup did not change the final installer or require another build.

## Findings Not Fixed

- None.

## Review

The Go entry reads no process arguments and extracts only its embedded server.
Its container menu filters names case-insensitively and uses a numbered `More` entry.
Selection resolves against the displayed list and retains the inspected container ID.

The install path preserves the v0.2.16/v0.2.17 source boundary, the v0.2.17 candidate requirement, and the exact 20-byte ELF comparison.
It validates and removes the staged candidate before backup, then retains a private `/tmp` backup before replacement.
Readiness re-inspects the selected ID and checks health and setup within the deadline.
Failure and interruption use a separate restoration context, including direct copy for a stopped container.
The extracted candidate has a separate temporary directory that is removed on exit.

The build stages the Go source, compressed server, and five license texts in a temporary module.
It runs installer tests for the host platform, then builds the static Linux/amd64 asset.
The release workflow uploads and publishes that installer only.
README, AGENTS.md, and the delivery specs describe the same entry and outputs.

The retained test container was not used for installer checks because `/app/sub2api` is its dependency launcher.
No production Docker installation or release publication was performed.

## Verification

- `gofmt -l scripts/installer/main.go tests/installer/installer_test.go`: pass.
- `git diff --check` for installer, build, release, user documentation, and delivery-spec paths: pass.
- `file` identifies both final outputs as static x86-64 Linux ELF executables.
- Focused installer tests: reused final implementation evidence. The host-platform Go test set passed in the build's temporary module after the final signal-isolation change.
- Build: reused final implementation evidence. The complete local build passed before that change; the installer was then checked and rebuilt against the completed server output.
- Model-control checks: reused existing handoff evidence because its patches and UI inputs are unchanged.

The corrected focused entry and redirect checks passed in the temporary host module:

```sh
go test . -run '^TestEntryAndEmbeddedPayload$|^TestInstallTransaction$/^redirect-rejected$'
```

The installer-only package step then rebuilt the Linux/amd64 installer from the completed raw server.
No server, UI, retained-container, production-install, or release-publication check was repeated.
