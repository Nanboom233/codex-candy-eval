# Sub2API Binary Replacement Delivery

## 1. Scope / Trigger

Deliver the approved static Linux/amd64 v0.2.17 server for the standard `/app/sub2api` container layout.
One self-extracting Go installer contains the compressed server and required license texts.
Standalone cleanup was removed from delivery.

## 2. Signatures

User installation:

```sh
chmod +x sub2api-codex-eval-installer
./sub2api-codex-eval-installer
```

The installer does not read or validate command-line arguments. Select the container through its numbered menu.
Always use the embedded server, a private backup directory under `/tmp`, and a fixed 120-second readiness wait.
There are no installation options, license command, or external candidate override.

Developer build and package commands:

```sh
python3 scripts/build.py --source /path/to/original-sub2api --check
```

Fixed outputs:

```text
dist/build/sub2api-codex-eval
dist/build/sub2api-codex-eval-installer
```

## 3. Contracts

Reuse the existing original Sub2API source and dependency/build caches.
The build copies that tree into a standard system temporary directory and applies the production patches listed in `patches/series`.
Apply patches with `git apply --binary --unidiff-zero --whitespace=nowarn`.
The vendor patch uses a Git binary hunk to preserve `texmath.js` CRLF bytes.
There is no source overlay. Production patches contain no tests or test-only data.
With `--check`, copy `tests/backend/` into the temporary tree's `backend/` directory before running the existing package checks.
The original tree remains untouched.
Compile with Linux/amd64, CGO disabled and the native `embed` tag.
The v0.2.17 archive retains `0.2.16` in `backend/cmd/server/VERSION`.
Build with `-X main.Version=0.2.17` so archive builds report the actual target release.
No new Git clone, source pin, source manifest or repeated provenance check is required.
Reuse installed frozen-lockfile frontend dependencies, then run native `pnpm build` for the patched sidebar/view before Go embed compilation. Prebuilt original frontend output does not contain the new route.
The build entry and release workflow use pnpm 9.15.9.
The same build entry creates the installer after a successful server build.
The [installer source](../../../scripts/installer/main.go) uses only the Go standard library.
Assemble its module, gzip server payload, and license files in a temporary directory. Do not add a root Go module.
Build the installer for Linux/amd64 with CGO disabled. Keep the raw server as a local build output.
Publish only `sub2api-codex-eval-installer` as the installation asset.
Embed the project MIT license and the four texts under `licenses/`.
The project MIT license retains the embedded name `Project-MIT.txt`.
Build helpers, tests, fixture logs, and private files stay outside the installer payload.
Keep installation, manual rollback, and history instructions in `README.md`. There is no separate `docs/` directory.

Development checks reuse one Docker container named `sub2api-codex-eval-test`.
Its PostgreSQL and Redis test dependencies run inside the same container.
Reuse this container for later checks instead of creating random fixture containers.
Keep this container and its fixture data after checks and task completion.
Routine cleanup must retain it. Stop it when needed instead of deleting it.
If absent, create only this fixed-name container when the next runtime check needs it.
Keep this development layout outside the user package.
Remove only identified earlier task containers. Preserve unrelated containers, volumes, images, and caches.

The installer uses the Go standard library and Docker CLI. The target host needs no Python, compiler, or external extraction tool.
List running containers through quiet IDs and JSON inspect.
First show names containing `sub2api`, without case sensitivity.
If other running containers exist, append `More` with the next numeric index.
Selecting that index shows all running containers. If no names match, show all running containers directly.
Omit `More` when the list already includes all running containers.
Require a numeric selection even for one container. Resolve numbers against the currently displayed list.
Empty input has no default. Omit Cancel and the `q` shortcut. Keep Ctrl-C interruption.
Use the selected immutable container ID for all effects.
Accept the existing v0.2.16 installation for upgrade, or v0.2.17 for reinstallation. Check that the Linux/amd64 candidate reports v0.2.17 before replacement.
Retain the existing exact 20-byte ELF header check.
Extract the embedded server silently into a private temporary directory. Remove this directory on exit.
Keep its lifetime separate from the retained executable backup.

Copy the current executable to `/tmp/sub2api-codex-eval-<first-12-ID>-*/sub2api` in a private directory, and print its exact path.
Default backups remain available until the temporary directory is removed by the user or system.
Stage the replacement and rename it over `/app/sub2api`, then restart that container.
Wait for `/health` status `ok` and `/setup/status` with `needs_setup:false`.
Require HTTP 200 directly. Do not follow redirects for either readiness endpoint.
Use host HTTP access derived from inspected port bindings, host networking or container IP.
Inspect the selected container again on every readiness attempt so restarted addresses and bindings are current.
Retry restarting containers and connection resets within the deadline; limit each HTTP request and retry sleep to remaining time.
Stopped/dead/paused containers fail promptly with exit, OOM and runtime diagnostics.
Normal startup progress includes elapsed/deadline, current state, restart count and endpoint without transient transport error text.
Print `Application ready` only after both health and setup checks pass; terminal failures retain the last actual endpoint/error.
Print step progress and actual command failures.

After replacement fails or is interrupted, restore the backup and verify the restarted original.
For a stopped container, restore with Docker cp without requiring a running process inside it.
Keep the printed backup for manual restoration.
Container recreation discards writable-layer replacement and requires reapplication.

## 4. Validation & Error Matrix

| Condition | Outcome |
| --- | --- |
| No running target or interruption during selection | Exit before replacement |
| Blank, nonnumeric, or out-of-range selection | Repeat the selection prompt without a default |
| Build source, listed patch, or frontend dependencies missing | Fail the build and do not package a candidate |
| A production patch cannot apply | Stop with the Git error before compilation |
| Wrong release or executable architecture | Reject before replacement |
| Replacement, restart or readiness fails | Restore the exact backup and verify the original |
| Ctrl-C after replacement | Restore, verify and propagate interruption |
| Candidate stops the container | Restore through Docker cp, restart and verify |
| Connection reset or restart changes address | Retry using current inspect facts within the deadline |
| Transient reset followed by readiness | Show startup progress, then explicit ready, without presenting the reset as an install failure |
| Readiness deadline expires | Report last state/request error and restore the backup |
| Restoration fails | Report failure and retained backup location |

## 5. Good / Base / Bad Cases

Good: upgrade an existing v0.2.16 container to v0.2.17 and retain its exact executable backup.
Base: a container recreated from its original image needs the same package reapplied.
Bad: compile upstream sources or require nc/buildx/helper images on the user's host.

## 6. Tests Required

`--check` retains the focused server/UI checks and touched-package vet, and runs the installer checks in a temporary host-platform Go module.
Keep the host test environment separate from the Linux/amd64 release build environment.
Use [one focused installer test set](../../../tests/installer/installer_test.go) with fake Docker and local HTTP.
Check selection, version/ELF checks, candidate cleanup, replacement, readiness, and restoration within that set.
Include numbered `More`, no-match selection, stopped-container restoration, and Ctrl-C in that same test set.
Keep timeout control inside tests. The public deadline stays fixed at 120 seconds.
Build the final installer and check its executable format and release asset path.
The retained container's `/app/sub2api` is a fixture dependency launcher. Its server runs from `/app/sub2api-bin`.
Do not overwrite that launcher for installer tests. The transaction harness does not prove installation on an actual Docker target.
The [implementation](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/fresh-native-implementation.md) records actual Docker replacement and stopped-container restoration.
The [independent review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/fresh-native-check.md) records interruption handling and final rebuilt artifacts.
The [correction review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/fresh-native-corrections-check.md) records real TCP reset recovery, changed Docker port binding, bounded timeout and stopped-container restoration.
The [browser/model review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/native-browser-model-check.md) records rebuilt native frontend and actual Chrome checks of the final packaged artifact.
The [installer/discovery review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/installer-disabled-model-check.md) records startup output, retained terminal errors and actual Docker default `/tmp` backup/restoration.
The [v0.2.17 review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/native-v0217-integrated-check.md) records the final build and package after the settings-load repair.
The [v0.2.17 browser check](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/native-v0217-browser.md) records actual settings controls and the rebuilt candidate's retained-container smoke.
Reuse passed evidence when inputs and affected behavior remain unchanged.
No standalone cleanup tests, repeated pristine baseline, architecture matrix, or additional full-suite requirements apply.

## 7. Wrong vs Correct

Wrong: report a restored executable before checking the restarted application.
Correct: complete readiness checks before reporting successful restoration.
