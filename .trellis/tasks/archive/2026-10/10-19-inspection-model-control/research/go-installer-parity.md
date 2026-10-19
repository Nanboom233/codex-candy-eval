# Research: Go self-extracting installer parity

- Query: Plan a single self-extracting Linux/amd64 Go installer while preserving the current installation logic and user experience.
- Scope: internal
- Date: 2026-10-19

## Findings

### Planning decision after research

The user removed `install`, `--container`, `--binary`, `--backup-dir`, and `--timeout` from the new installation interface.
Use the parameter-free entry in task R9 and `design.md`.
Always select a numbered container, use the embedded server, retain private `/tmp` backups, and use a fixed 120-second readiness deadline.
First show case-insensitive `sub2api` name matches. Append a numbered `More` option when the filtered list omits running containers.
Remove the Cancel option and `q` shortcut. Use Ctrl-C for interruption.
If no name matches, show all running containers directly. A single match still requires selection.
The Python CLI details below document the old implementation. They are not Go compatibility requirements.
Do not retain argument-specific tests or replace removed arguments with other configuration inputs.
Keep required license texts embedded. The plan review removed the optional license command as unnecessary scope.

### Current public contract

The current installer is `scripts/install.py`. It uses only Python's standard library and the Docker CLI.
The host must still provide Docker access, but it does not need Sub2API source code or a build tool.

- `scripts/install.py:196-216` defines the public interface.
- The optional operation is `install` and is the default.
- The options are `--container`, `--binary`, `--backup-dir`, and positive `--timeout`.
- The readiness timeout defaults to `120` seconds.
- Runtime failures return exit code `1` and argument parsing failures use argparse exit code `2`.

The normal user-visible sequence has four numbered steps.

1. `scripts/install.py:140-158` checks the existing release, candidate ELF header, staged candidate version, and staged-file cleanup.
2. `scripts/install.py:159-172` creates and prints the retained backup path.
3. `scripts/install.py:175-180` replaces the selected executable and restarts the selected container.
4. `scripts/install.py:180-193` waits for readiness and prints the existing final success or restoration messages.

Use numeric selection for containers and `More`. Preserve the four step labels, backup output, diagnostics, and final message.
Payload extraction must remain silent so it does not create a fifth install step.

### Container selection and immutable target

`select_container` in `scripts/install.py:31-50` has two modes.

- `--container` inspects that name directly.
- Without it, `docker ps -q` lists running IDs, each ID is inspected, and the user must enter a number.
- Empty input has no default. `q` cancels. Invalid numbers repeat the prompt.

After selection, `scripts/install.py:134-139` stores the inspected immutable ID and rejects a stopped target.
Every later Docker operation uses that ID, not the display name.
The new Go entry retains only numbered interactive selection for installation.
Do not replace this with implicit single-container selection or an auto-switch after a failed inspect.

### Version, ELF, staging, and cleanup

The existing installation accepts an installed Sub2API version containing `0.2.16` or `0.2.17`.
The candidate must report `0.2.17` (`scripts/install.py:140-156`).
`scripts/build.py:59-60` injects `main.Version=0.2.17` into the server build.

The local candidate check compares exactly the first 20 ELF bytes.
It requires ELF class 64, little-endian data, current ELF version, and machine `x86-64`.
The current error says static, but the code only checks these 20 bytes.
The Go port must retain the same exact comparison instead of adding a new static-link inspection.

The candidate is copied to `/app/sub2api.codex-eval-new`, made executable, and executed with `-version`.
The staged path is removed in a `finally` block (`scripts/install.py:151-158`).
Use a deferred cleanup path in Go after a successful staged copy.

### Backup and replacement transaction

The default backup uses `tempfile.mkdtemp` under `/tmp` with prefix
`sub2api-codex-eval-<first-12-ID>-` (`scripts/install.py:159-167`).
The retained executable is `<backup-dir>/sub2api`, and its exact path is printed.
The backup is intentionally not deleted after a successful install.

`--backup-dir` creates a timestamped child named `<first-12-ID>-<UTC timestamp>`.
The new Go entry removes this override. It retains only the private default backup and never deletes that backup automatically.

The backup is complete before any mutation (`scripts/install.py:167-177`).
For a running container, `replace` copies to the staged path, then runs
`chmod 755 ... && mv -f ... /app/sub2api` (`scripts/install.py:123-131`).
For a stopped container, it copies directly to `/app/sub2api` without `docker exec`.
The public install rejects a stopped selected target, but the stopped branch is required for restoration after a candidate stops the container.

The transaction marks itself changed before replacement.
Any replacement, restart, readiness, or interrupt error restores the exact backup, restarts the container, and runs readiness again (`scripts/install.py:173-191`).
If restoration fails, the error retains the backup path.
The Go signal handler must preserve this path and must not cancel restoration.

### Readiness and endpoint behavior

`endpoint` in `scripts/install.py:52-65` derives the address from fresh Docker inspect facts.

1. Read `SERVER_PORT`, defaulting to `8080`.
2. Prefer the first published TCP binding.
3. Map wildcard host bindings to `127.0.0.1`.
4. Use `127.0.0.1:SERVER_PORT` for host networking.
5. Otherwise use the first Docker network IP.

`readiness` in `scripts/install.py:68-120` re-inspects the selected ID on every attempt.
It fails promptly for stopped, dead, and paused states. It retries restarting states within the deadline.
Terminal diagnostics retain the state and last request details, with exit code, OOM, and runtime errors where applicable.
It requests `/health` and `/setup/status` in that order.
The required responses are HTTP 200, `status=ok`, and `data.needs_setup=false`.

Each HTTP request has a maximum three-second timeout, bounded by the remaining overall deadline.
Connection, JSON, HTTP, and readiness failures retry until the deadline.
Progress prints elapsed time, deadline, state, restart count, and current endpoint.
Transient request text is retained for the terminal error, but is not printed in ordinary startup progress.
Success prints `Application ready` only after both endpoints pass.

The Go port should use `net/http` with a per-request context deadline and a monotonic overall deadline.
It must close every response body and connection and must preserve the current retry and progress cadence.
Do not add a new backoff policy, health endpoint, or readiness state.

### Self-extracting payload boundary

The recommended release output is one static executable:

```text
dist/build/sub2api-codex-eval-installer
```

Keep `dist/build/sub2api-codex-eval` as the raw server build input.
Do not ship that raw server as a second release asset.

Use a small tracked Go installer source directory with standard-library imports only.
The existing repository has no root Go module (`AGENTS.md:59-60`).
Build the installer in a temporary staging module so the repository does not gain a root module or generated payload files.

The staging module should contain:

- `main.go` copied from tracked installer source and a minimal staging-only `go.mod`.
- A gzip-compressed copy of the raw server binary.
- The project license and four files under `licenses/`.

Use `embed` for the compressed server and license bytes.
On a normal install, write the server to a private `os.MkdirTemp` directory, set mode `0755`, and remove only that extraction directory on exit.
The extracted candidate is separate from the retained Docker backup.
The new installation entry always uses the embedded candidate. It has no external candidate override.

The normal install path must not need Python, `tar`, `unzip`, `gzip`, or a helper image.
Go standard-library gzip extraction meets the self-extracting requirement without adding a host dependency.

The old archive carried five license texts (`scripts/build.py:64-74`).
Embed their unchanged bytes in the installer. Do not add a separate license operation.
Do not add a general extraction mode or a second CLI framework.
The default install command and output must remain unchanged.

### Build and release wiring

`scripts/build.py:34-61` currently builds the raw server after applying the patch series and frontend build.
`scripts/build.py:64-86` then creates the tar archive.
Replace only the packaging stage with:

1. Build the raw server using the existing upstream Go module and flags.
2. Stage the Go installer source, compressed server, and license bytes in a temporary directory.
3. Build the static Linux/amd64 installer with `CGO_ENABLED=0`, `-trimpath`, and stripped release flags.
4. Emit the installer at the fixed output path above.

The developer command remains `python3 scripts/build.py --source ... [--check]`.
The Python build entry is a developer/CI tool, not a target-host dependency.

`.github/workflows/release.yml:22-39` already provisions Go from the upstream module and runs the build entry.
Change the upload artifact to the installer executable.
`release.yml:41-55` should publish that one file instead of the tar archive.
Do not add a second archive, download step, source build, or helper image.

### Existing test and fixture boundaries

`tests/installer/test_native_version.py:18-53` currently imports the Python installer and mocks Docker.
Replace its Python-module dependency with one focused changed-boundary check against the Go installer or its standard-library test package.
It must retain the two installed-version cases and the candidate-version rejection before backup or restart.

The current fixed fixture is not a direct installer target.
`tests/installer/native_inspection_fixture.py:76-92` copies the server to `/app/sub2api-bin` and installs
`native_fixture_start.sh` as `/app/sub2api`.
The wrapper starts PostgreSQL, Redis, the supplier, and then `/app/sub2api-bin` (`tests/installer/native_fixture_start.sh:1-41`).
The fixture therefore uses `/app/sub2api` as a launcher, not as the native server executable.

Running the new installer against this container would replace the launcher and make the next restart skip the fixture dependencies.
That would produce a false installer failure and would violate the fixed-container preservation rule.

Use the fixed container only for the existing server/frontend smoke checks.
Use a fake Docker command plus a local HTTP server for the new Go transaction boundary.
The fake must model inspect, `ps -q`, `cp`, `exec`, `restart`, staged candidate version, and readiness address changes.
This check can prove selected-ID use, candidate rejection, backup-before-mutation, restart/readiness, and restoration without altering the retained fixture.
Do not create a random fixture container or add a test-only production flag.

Reuse completed installer readiness and restoration evidence when the Go behavior is covered by focused fake-Docker tests.
Do not add race, coverage, full-suite, architecture-matrix, or UI checks for this packaging change.

### Smallest acceptable validation set

Use only checks that cover changed behavior:

- One focused Go test set for parameter-free selection, version and ELF checks, candidate cleanup, installation, readiness, and restoration.
  Use the same fake-Docker and local HTTP harness for transaction paths.
- The existing build entry compiles the release installer. Check its Linux/amd64 format and the single release asset path.
- Reuse completed server, UI, and Go evidence for unchanged product inputs. CI retains the existing `--check` entry.

Do not add a second runtime fixture, a browser matrix, a full Go suite, or a new deployment gate.

## Files found

- `scripts/install.py` — current Docker selection, version/ELF checks, replacement, readiness, backup, and restoration transaction.
- `scripts/build.py` — upstream patch application, static server build, and current archive packaging entry.
- `.github/workflows/release.yml` — Go setup, build/check invocation, artifact upload, and release publication.
- `tests/installer/test_native_version.py` — existing focused installer release-boundary test.
- `tests/installer/native_inspection_fixture.py` — fixed-container preparation and `/app/sub2api-bin` server placement.
- `tests/installer/native_fixture_start.sh` — fixture startup wrapper installed at `/app/sub2api`.
- `.trellis/spec/backend/minimal-linux-delivery.md` — current archive, installer, Docker, backup, readiness, and validation contracts.
- `.trellis/spec/backend/quality-guidelines.md` — focused-check and no-extra-matrix rules.
- `.trellis/spec/backend/directory-structure.md` — no-root-module, package ownership, and test placement rules.
- `AGENTS.md` — repository layout, fixed outputs, fixture reuse, and host dependency rules.
- `README.md` — current user installation, rollback, build, package, and fixed-output documentation.

## Code patterns

- `scripts/install.py:31-50` — explicit or numbered container selection with no blank default.
- `scripts/install.py:52-65` — endpoint selection from current inspect data.
- `scripts/install.py:68-120` — bounded readiness loop and progress output.
- `scripts/install.py:123-131` — running/stopped replacement paths.
- `scripts/install.py:134-193` — version checks, backup, transaction, restoration, and user messages.
- `scripts/install.py:196-216` — CLI defaults and exit behavior.
- `scripts/build.py:34-61` — raw server build boundary.
- `scripts/build.py:64-74` — current binary and five-license archive inputs.
- `tests/installer/test_native_version.py:18-53` — minimal current version-boundary assertions.
- `tests/installer/native_inspection_fixture.py:76-152` — retained fixture lifecycle and server/wrapper split.
- `tests/installer/native_fixture_start.sh:1-41` — fixture dependency startup and `/app/sub2api-bin` execution.

## External references

- Go toolchain version is supplied by `upstream/backend/go.mod` (`go 1.27.2`) and `.github/workflows/release.yml:22-25`.
- The proposed installer uses Go standard-library packages such as `embed`, `compress/gzip`, `os`, `os/exec`, `net/http`, `encoding/json`, and `os/signal`.
- No external runtime library or target-host download is required.

## Related specs

- `.trellis/spec/backend/minimal-linux-delivery.md:5-19,43-89,91-125` — current package, installer, backup, readiness, and validation contract.
- `.trellis/spec/backend/quality-guidelines.md:12-29,37-40` — focused checks and excluded test matrices.
- `.trellis/spec/backend/directory-structure.md:9-14,20-41` — patch ownership, no root module, and package boundaries.
- `.trellis/spec/backend/index.md:42-60` — quality gate and spec-only reference checks.
- `AGENTS.md:31-62,72-78` — patch/test layout, fixed container, host dependencies, and installer verification.

## Caveats / Not Found

- No installer implementation, build, release, or runtime check was executed during this research pass.
- The fixed `sub2api-codex-eval-test` container was not used because its `/app/sub2api` path is a fixture launcher. The server runs from `/app/sub2api-bin`.
- The planned transaction checks simulate Docker and HTTP. They do not prove a new installation on an actual Docker target.
  Do not reconfigure or recreate the retained fixture for this task.
- The final plan retains embedded license texts and removes the optional license operation. No general extraction command is required.
- Parent planning artifacts must update the current minimal-linux delivery, directory, README, build, release, and focused-check references. This report does not edit specs or implementation files.
