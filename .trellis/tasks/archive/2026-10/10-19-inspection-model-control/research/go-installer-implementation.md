# Go installer implementation

## Result

The static Linux/amd64 installer contains the compressed server and five license texts.
Installation uses numbered container selection and the existing four steps. The installer does not read process arguments.

The first list contains case-insensitive `sub2api` name matches.
A numbered `More` entry expands the list to all running containers.
Without matches, the first list contains all running containers.
Empty input and invalid values select nothing. No Cancel option or `q` shortcut exists.
Every operation after selection uses the inspected immutable container ID.

The installer preserves installed v0.2.16/v0.2.17 acceptance and requires candidate v0.2.17.
It checks the original 20-byte ELF header and removes the staged version-check candidate.
It retains a private `/tmp` backup before replacement and prints the exact backup path.
It uses a fixed 120-second deadline for each installation or restoration readiness wait.
Each attempt reads fresh Docker facts and checks native health and setup endpoints.
HTTP requests use a three-second maximum within the remaining deadline.

Failure or interruption after replacement restores the backup, restarts the container, and checks readiness again.
Stopped-container restoration uses direct Docker copy without an executable inside the container.
Restoration uses a separate context. Docker commands use separate process groups so terminal signals do not cancel restoration.
The temporary extracted server is separate from retained backups and is removed when installation exits.

## Files

- `scripts/installer/main.go`: embedded payload, selection, Docker commands, validation, replacement, readiness, interruption, and restoration.
- `scripts/build.py`: temporary installer module, upstream Go version, gzip payload, embedded licenses, host tests, and static release build.
- `tests/installer/installer_test.go`: one fake-Docker and local HTTP test set.
- `.github/workflows/release.yml`: upload and publish only the installer asset.
- `README.md`: direct installation, selection, fixed defaults, recovery, build, and output instructions.
- `AGENTS.md`: current owners, outputs, host checks, and retained-fixture boundary.
- `licenses/THIRD_PARTY_NOTICES.txt`: embedded license paths and unchanged attribution terms.
- Removed `scripts/install.py` and `tests/installer/test_native_version.py` after their replacement passed.

Project specs remain under the main session's ownership.
Product patches, model-control tests, and the staged index were not changed by this implementation.

## Checks

The final host-platform installer test set passed through `scripts/build.py`'s `package(..., check=True)` function.
That function copied tests into the temporary module and ran `go test .` before Linux/amd64 compilation.
The final result was `ok sub2api-codex-eval-installer 8.545s`.

The same test set covers filtered selection, numbered More, no-match display, no blank default, and immutable target IDs.
It covers upgrade, reinstallation, candidate version/ELF rejection, backup ELF rejection, and candidate cleanup.
It covers success, restart/address changes, incomplete-setup timeout, interruption, stopped-container restoration, and restoration failure diagnostics.
It checks backup-before-replacement ordering, retained private backups, embedded licenses, and extracted payload cleanup.
The parameter-free public entry runs successfully through the same fake Docker command and local HTTP server.

The normal build passed:

```sh
python3 -B scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check
```

It applied the seven production patches in a temporary copy, built the native frontend, and compiled both outputs.
The build completed with existing upstream frontend warnings about chunk size, imports, and Browserslist data.
After the final signal-isolation change, `package(..., check=True)` checked and rebuilt only the installer using the completed server output.
No unchanged service or UI acceptance checks were repeated.

`gofmt` formatted both changed Go files.
Python AST parsing accepted `scripts/build.py`.
`git diff --check` passed.
Current delivery references contain no retired Python installer, tar archive, or installation options.
The release workflow uploads and publishes the same installer filename.
`file` identified both outputs as static x86-64 Linux ELF executables.

## Outputs and limits

- `dist/build/sub2api-codex-eval`: raw server, approximately 120 MiB.
- `dist/build/sub2api-codex-eval-installer`: single installation asset, approximately 44 MiB.

The obsolete local tar archive was removed from the output directory.
The target host needs Docker CLI access. It needs no Python, Go toolchain, or extraction command.

The host-platform tests simulate Docker and use real local HTTP responses.
The Linux installer was compiled and inspected, not installed into an actual Docker target during this task.
The retained `sub2api-codex-eval-test` container was not changed because `/app/sub2api` is its dependency launcher.
No production installation or release publication occurred.
The candidate is ready for independent delivery review.

## Argument-check correction and scope audit

The user removed the added argument rejection check after the first implementation.
`run` now receives only the context and input/output streams. `main` does not read or forward process arguments.
The existing entry test no longer checks argument rejection. No replacement parser, configuration, or check was added.

The audit compared installer code, build assembly, and tests with the Python implementation and the approved plan.
No removed-interface test remains.
Version, ELF, backup, readiness, candidate cleanup, and restoration checks follow the inherited installation contract.
Container name filtering and numbered More follow the user's selection decisions.
Temporary-module assembly, payload extraction, and license checks support the single-file delivery.
Resource cleanup and separate restoration contexts preserve the required interruption behavior.
The build retains one installer test set. It adds no race, coverage, architecture matrix, or second runtime gate.
The entry test no longer extracts and validates the same payload twice. Its public-entry transaction covers extraction, ELF validation, and cleanup.
License content checks remain in that existing test.

The audit found that Go followed HTTP redirects, while the Python installer rejected non-200 responses without following redirects.
The HTTP client now returns redirect responses to the existing non-200 check.
One case in the existing transaction harness checks `HTTP 302` rejection and restoration.
The checker also found added SIGTERM handling. Signal registration now covers only Ctrl-C, as the Python installer did.
Docker process-group isolation and the separate restoration context remain in place.

The corrected entry and HTTP transaction checks passed in the temporary host module:

```sh
go test . -run '^TestEntryAndEmbeddedPayload$|^TestInstallTransaction$/^redirect-rejected$'
```

The result was `ok sub2api-codex-eval-installer 3.607s`.
The current `package(..., check=False)` helper then rebuilt only the Linux/amd64 installer from the completed raw server.
This targeted evidence supersedes the earlier entry evidence. The unchanged transaction evidence remains valid.
`gofmt` and `git diff --check` passed for the correction.
The raw server, retained Docker fixture, and product/UI checks were not changed or repeated.

The final selection cleanup removed stale `q` and textual `more` inputs from the filtered selection case.
That case retains blank input, an out-of-range number, and valid numeric selection.
Only the affected case ran in the temporary host module:

```sh
go test . -run '^TestSelection$/^filtered$'
```

The result was `ok sub2api-codex-eval-installer 0.947s`.
Production source and the completed installer were unchanged. No installer rebuild was necessary.
