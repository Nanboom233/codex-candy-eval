# Sub2API Binary Replacement Delivery

## 1. Scope / Trigger

Deliver the approved static Linux/amd64 v0.2.17 server for the standard `/app/sub2api` container layout.
The package contains the binary, installer and necessary license texts.
Standalone cleanup was removed from delivery.

## 2. Signatures

User installation after unpacking the archive:

```sh
python3 install.py install --container sub2api
```

The install operation is the default.
Options are `--container`, `--binary`, `--backup-dir` and positive `--timeout`.
Defaults use the adjacent binary, a private backup directory under `/tmp` and a 120-second readiness wait.

Developer build and package commands:

```sh
python3 scripts/build.py --source /path/to/original-sub2api --check
```

Fixed outputs:

```text
dist/build/sub2api-codex-eval
dist/build/sub2api-codex-eval.tar.gz
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
The same build entry packages the executable after a successful build.

The archive has an outer `sub2api-codex-eval/` directory.
Its files are `sub2api-codex-eval`, `install.py` and five `LICENSE/` texts.
The root project MIT license uses the archive name `LICENSE/Project-MIT.txt`.
Build helpers, tests and fixture logs stay outside the package.
The repository installer lives at `scripts/install.py`. Its archive entry remains `install.py` beside the executable.
Keep installation, manual rollback, and history instructions in `README.md`. There is no separate `docs/` directory.

Development checks reuse one Docker container named `sub2api-codex-eval-test`.
Its PostgreSQL and Redis test dependencies run inside the same container.
Reuse this container for later checks instead of creating random fixture containers.
Keep this development layout outside the user package.
Remove only identified earlier task containers. Preserve unrelated containers, volumes, images, and caches.

The installer uses Python's standard library and Docker CLI.
Without `--container`, list running containers and require a number; empty input has no default.
Use the selected immutable container ID for all effects.
Read Docker facts through quiet IDs and JSON inspect.
Accept the existing v0.2.16 installation for upgrade, or v0.2.17 for reinstallation. Check that the Linux/amd64 candidate reports v0.2.17 before replacement.

Copy the current executable to a private `/tmp/sub2api-codex-eval-<container-id>-*/sub2api` directory created with `tempfile.mkdtemp(dir="/tmp")`, and print its exact path.
An explicit `--backup-dir` retains the chosen directory's timestamped backup layout.
Default backups remain available until the temporary directory is removed by the user or system.
Stage the replacement and rename it over `/app/sub2api`, then restart that container.
Wait for `/health` status `ok` and `/setup/status` with `needs_setup:false`.
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
| No running target, cancelled selection or missing binary | Exit before replacement |
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

`--check` runs the relevant evaluation/sender/admin tests and touched-package vet before the static build.
Check the changed replacement and failure-restoration paths using local disposable resources.
The [implementation](../../tasks/10-11-sub2api-runtime-patch/research/fresh-native-implementation.md) records actual Docker replacement and stopped-container restoration.
The [independent review](../../tasks/10-11-sub2api-runtime-patch/research/fresh-native-check.md) records interruption handling and final rebuilt artifacts.
The [correction review](../../tasks/10-11-sub2api-runtime-patch/research/fresh-native-corrections-check.md) records real TCP reset recovery, changed Docker port binding, bounded timeout and stopped-container restoration.
The [browser/model review](../../tasks/10-11-sub2api-runtime-patch/research/native-browser-model-check.md) records rebuilt native frontend and actual Chrome checks of the final packaged artifact.
The [installer/discovery review](../../tasks/10-11-sub2api-runtime-patch/research/installer-disabled-model-check.md) records startup output, retained terminal errors and actual Docker default `/tmp` backup/restoration.
The [v0.2.17 review](../../tasks/10-11-sub2api-runtime-patch/research/native-v0217-integrated-check.md) records the final build and package after the settings-load repair.
The [v0.2.17 browser check](../../tasks/10-11-sub2api-runtime-patch/research/native-v0217-browser.md) records actual settings controls and the rebuilt candidate's retained-container smoke.
Reuse passed evidence when inputs and affected behavior remain unchanged.
No standalone cleanup tests, repeated pristine baseline or additional full-suite requirements apply.

## 7. Wrong vs Correct

Wrong: report a restored executable before checking the restarted application.
Correct: complete readiness checks before reporting successful restoration.
