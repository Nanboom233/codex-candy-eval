# Sub2API Binary Replacement Delivery

## 1. Scope / Trigger

Deliver standalone programs from the selected upstream tag for Linux amd64/arm64, Darwin amd64/arm64, and Windows amd64.
The verified integration baseline is v0.2.19.
Each self-contained Go installer carries its native program, required Linux container program, and original license texts.
Support official Docker/Compose, native executable, and Apple container installations through their existing startup owners.
Standalone cleanup was removed from delivery.

## 2. Signatures

User installation:

```sh
chmod +x sub2api-plugins-installer-v0.2.19-linux-amd64
./sub2api-plugins-installer-v0.2.19-linux-amd64
```

The installer does not read or validate command-line arguments. Select the installation form, then its target.
Use embedded programs, a retained private backup, and the fixed 120-second readiness wait.
Native programs use same-filesystem staging. Other temporary files use the host's temporary directory.
There is no license command, external candidate override, or persistent installer configuration.

Developer build and package commands:

```sh
python3 scripts/build.py
```

Actions exposes a required string input named `upstream_tag` for manual runs, with default `v0.2.19`. The user can replace it before the run starts.
Pass the requested value through `UPSTREAM_TAG` in the process environment. Do not interpolate user input into shell commands.
Local builds use the same entry and environment input, for example `UPSTREAM_TAG=v0.2.19 python3 scripts/build.py`.
Push and pull-request builds resolve the latest official release through the public upstream release API.
Local builds without `UPSTREAM_TAG` use the same resolver. Resolve once before source preparation and reuse the resulting actual `tag_name`.
Use no user credentials for this public lookup. Query failure or an unusable response fails the build without a version fallback.

Fixed outputs:

```text
dist/build/<os>-<arch>/sub2api-plugins-<upstream-tag>-<os>-<arch>
dist/build/<os>-<arch>/sub2api-plugins-installer-<upstream-tag>-<os>-<arch>
```

## 3. Contracts

Download `https://github.com/Wei-Shaw/sub2api.git` at the selected tag with depth one into a temporary workspace.
Reuse normal dependency/build caches and apply top-level `patches/*.patch` files in lexical filename order.
Apply standard contextual patches with `git apply --binary`. Stop on failure without weaker-context retries.
Use the command runner as the single patch-progress owner. Do not print another `Apply ...` line for the same command.
The vendor patch uses a Git binary hunk to preserve `texmath.js` CRLF bytes.
There is no source overlay. Production patches contain no tests or test-only data.
Discover Python, Go, CJS, and native Vue test entries directly from flat `tests/` on every build.
Use `// sub2api-package:` annotations to stage native suites in their temporary package and run the discovered cases with `go test`.
Keep all assertions and execution support in `tests/`. Use `assets/` for data only.
Local source checkouts remain untouched. The entry needs no source-path, check, or test-skipping argument.
Prepare source, dependencies, and the frontend once. Run host tests once, then compile the five supported targets.
Use CGO disabled and the native `embed` tag. Add `.exe` to Windows programs.
The v0.2.19 source retains `0.2.18` in `backend/cmd/server/VERSION`.
Resolve the actual selected tag once. Remove only its leading `v` for canonical server and installer version metadata.
Use that same selection for source download, test context, both linker values, artifact names, validation, and upload paths.
Do not read a hardcoded version separately during output validation.
Install frontend dependencies with the frozen lockfile and keep Go modules readonly.
Run native `pnpm build` before Go embed compilation. Prebuilt original frontend output does not contain the new route.
Use the same selected stable pnpm for dependency installation, frontend compilation, and native Vue tests.
Translate upstream overrides and effective build settings into temporary pnpm configuration without changing the lockfile or override values.
The same build entry creates the installer after a successful server build.
The [installer source](../../../scripts/installer/) uses standard-library installation operations and PTerm v0.12.84 for terminal interaction.
Prepare dependencies, compressed programs, and license files in its temporary module. Do not add a root Go module.
Compile each installer for its host with CGO disabled. Keep each raw native program as a local output.
Linux installers reuse their native payload for Docker. Darwin installers additionally embed the matching Linux architecture for containers.
Darwin arm64 supports Apple container with Linux arm64. Windows amd64 embeds Linux amd64 for Docker.
Use the selected upstream tag, including `v`, in output names. Windows names add `.exe`.
Upload each executable directly with `archive: false`. Do not create tar files or another product archive.
This mode ignores `name` and derives the artifact name from the file's basename.
Five platform pairs produce ten direct product artifacts. Unix downloads need `chmod +x` before execution.
Keep executable outputs under `dist/build/<os>-<arch>/`. Do not retain duplicate legacy output copies.
Stage the complete new output set on the destination filesystem before publication.
Swap directories through rename. If publication fails, restore the old directory at its original path.
Do not copy partial new outputs into the live directory or nest the old directory during restoration.
Collect root `LICENSE` and third-party `licenses/` inputs under installer resources `licenses/` with separate namespaces.
Preserve original license text and its copyright notices without a fixed file count.
Use original license files only. Do not add a separate `THIRD_PARTY_NOTICES.txt` summary.
Build helpers, tests, fixture logs, and private files stay outside the installer payload.
Keep installation, manual rollback, and history instructions in `README.md`. There is no separate `docs/` directory.
The approved exact README restoration retains historical names and arguments. Use the current entry and filenames in this spec.
Use `LGPL-3.0-or-later` for this adaptation. Preserve inherited MIT and dependency notices and required GPL/LGPL texts.

The complete standalone kit contains the builder, installer, aligned common tests, and `.github/workflows/build.yml`.
Local and CI builds use the same argument-free entry without a parent checkout or remote build Action.
Use push, pull request, and manual triggers on `ubuntu-latest` with read-only repository permissions.
Reference official Actions by exact release tags such as `@v7.0.1`, without commit hashes or major-only tags.
Disable checkout credential persistence. Publish Actions artifacts only, not Releases or container images.
Use verified current stable tools. Preserve failures through `pipefail` when capturing build output.
Before upload, reject missing, empty, or non-executable outputs with an error annotation and an explicit job-summary message.
Treat output paths as data in shell messages. Do not evaluate them through command substitution.
Upload available build logs as diagnostics only when the build step fails.
Keep summaries to status, failed step, and available diagnostic locations. Do not copy build-log bodies or tails into summaries.
Link a diagnostic artifact only after its upload succeeds. Otherwise, direct the reader to the native step log.
Report output validation and upload failures through their own step logs and job-summary messages.

Manual container fixture helpers live under flat `tests/native_*` and stay outside automatic discovery.
No retained container or manual browser session is required for this build migration.
Preserve unrelated containers, volumes, images, and caches.

The target host needs no Python, compiler, or external extraction tool. Container modes use the selected runtime CLI.
Native mode controls the effective official systemd service or asks the user to stop/start through the original launch method.
It does not create a service or reconstruct manual arguments, environment, or working directory.
Resolve native symlinks and preserve ownership and permissions during same-filesystem replacement.
Ask for the native address with default `127.0.0.1:8080` and Linux program path default `/opt/sub2api/sub2api`.

Use PTerm inline selection, text input, confirmations, and timed spinners only with actual supported terminal stdin/stdout.
Injected or redirected streams and `TERM=dumb` use plain numbered prompts and sequential lines without cursor controls.
Honor `NO_COLOR`. Stop and join spinner rendering before prompts, diagnostics, and final results.
Show elapsed time for unknown-duration stages without an invented percentage.
Keep animated stage text short and width-aware. PTerm's spinner alone cannot clear automatic terminal wrapping.
Use the selected component's area rendering for wrapped frames. Do not build another terminal redraw engine.
Show one live elapsed timer. Do not embed sampled elapsed time, container state, or endpoint data in animated frame text.

Override every PTerm prompt's interrupt handler. Cancel the operation context and check it before accepting any returned choice.
UI callbacks do not exit directly or own rollback. Keep command diagnostics with one final error owner.

Resolve language once from the first nonempty `LC_ALL`, `LC_MESSAGES`, or `LANG` value.
Chinese locales use Simplified Chinese. Other, absent, `C`, and `POSIX` locales use English.
Accept locale case, `_`/`-` separators, and encoding suffixes. Higher-priority `C` overrides lower-priority Chinese.

Localize installer-owned messages. Keep paths, commands, versions, identifiers, and external diagnostics literal.
Do not add a language menu, argument, configuration file, OS-language lookup, or locale framework.

Docker target selection retains the following behavior.
List running containers through quiet IDs and JSON inspect.
First show names containing `sub2api`, without case sensitivity.
If other running containers exist, append `More` with the next numeric index.
Selecting that index shows all running containers. If no names match, show all running containers directly.
Omit `More` when the list already includes all running containers.
Require explicit selection even for one container. Plain prompts require a number from the currently displayed list.
Empty input has no default. Omit Cancel and the `q` shortcut. Keep Ctrl-C interruption.
Use the selected immutable container ID for all effects.
The Docker helper buffers child stdout and stderr. Forward buffered output only for successful non-quiet commands.
On command failure, return captured detail with operation context for the final error owner to print once.
Quiet failures retain their captured detail. Preserve cancellation causes and separate restoration outcomes.
Read installed and staged executable versions from bounded `-version` probes in their target environment, using both output streams.
Normalize only an optional leading `v` and compare complete application version tokens with the build version.
Image tags such as `latest` and labels do not determine executable versions.
Mismatch or unresolved probes produce one combined warning and `Continue installation? [y/N]` before backup or replacement.
Only explicit `y` or `yes` continues. Empty input, refusal, EOF, or interruption cancels with candidate cleanup.
Matching versions need no extra version prompt. Keep operational and executable-format failures separate from version warnings.
Use Go's ELF, Mach-O, and PE readers to check the selected runtime format and architecture before replacement.
Extract the embedded server silently into a private temporary directory. Remove this directory on exit.
Keep its lifetime separate from the retained executable backup.

Copy the current executable into a private backup and print its exact path.
Default backups remain available until the temporary directory is removed by the user or system.
Stage the replacement and rename it over `/app/sub2api`, then restart that container.
Wait for `/health` status `ok` and `/setup/status` with `needs_setup:false`.
Require HTTP 200 directly. Do not follow redirects for either readiness endpoint.
Use host HTTP access derived from inspected port bindings, host networking or container IP.
Inspect the selected container again on every readiness attempt so restarted addresses and bindings are current.
Retry restarting containers and connection resets within the deadline; limit each HTTP request and retry sleep to remaining time.
Stopped/dead/paused containers fail promptly with exit, OOM and runtime diagnostics.
Plain startup progress includes elapsed/deadline, current state, restart count, and endpoint without transient transport error text.
Dynamic progress uses the short readiness stage and live component timer. Final failures retain the actual endpoint and error.
Print `Application ready` only after both health and setup checks pass; terminal failures retain the last actual endpoint/error.
Print step progress and actual command failures.

Before any application or startup mutation, cancellation removes the stage without restarting the unchanged application.
After program replacement or a managed stop begins, recovery restores the original program and startup state, then checks readiness.
Recovery commands and manual prompts use a separate restoration context, rather than the canceled installation context.
Unavailable recovery input retains the backup and reports the pending manual action without claiming success.
For a stopped container, restore with Docker cp without requiring a running process inside it.
Keep the printed backup for manual restoration.
Container recreation discards writable-layer replacement and requires reapplication.

Native systemd installation stops the effective service, replaces the program, starts that service, and checks readiness.
Native manual installation asks the user to stop, replaces the program, asks for original-method startup, and checks readiness.
Recovery follows the same startup owner. Preserve data, service configuration, arguments, environment, and working directory.
Program rollback does not reverse database migrations.

Apple installation targets `/app/storage/runtime/sub2api` and preserves `/app/storage/runtime/base-image-id`.
Copy the stage into the running application container, preserve ownership, and set executable access before rename.
Keep the backup on the host. Stop/start only that application container, leaving PostgreSQL and Redis running.
Apple container copy requires a running container. If stopped, start the same container before restoration.
If its runtime or volume cannot start, retain the host backup and report incomplete recovery.
Normal recreation retains the runtime program while the base image ID is unchanged. Image changes require reinstallation.
Data-only backups do not include that runtime program.

## 4. Validation & Error Matrix

| Condition | Outcome |
| --- | --- |
| No running target or interruption during selection | Exit before replacement |
| Blank, nonnumeric, or out-of-range selection | Repeat the selection prompt without a default |
| Source download, dependency, patch, or test failure | Fail the build and preserve previous outputs |
| A production patch cannot apply | Stop with the Git error before compilation |
| Expected CI output is missing, empty, or non-executable | Fail validation, name the output in an annotation and summary, and skip product upload |
| CI build step fails | Fail the job and retain its available build log as diagnostics |
| Automatic upstream release lookup fails | Report the lookup failure and preserve previous outputs |
| Mismatched or unresolved release | Warn and require explicit confirmation before replacement |
| Wrong executable architecture | Reject before replacement |
| Replacement, restart or readiness fails | Restore the exact backup and verify the original |
| Ctrl-C after replacement | Restore, verify and propagate interruption |
| Candidate stops the container | Restore through Docker cp, restart and verify |
| Connection reset or restart changes address | Retry using current inspect facts within the deadline |
| Transient reset followed by readiness | Show startup progress, then explicit ready, without presenting the reset as an install failure |
| Readiness deadline expires | Report last state/request error and restore the backup |
| Restoration fails | Report failure and retained backup location |
| Native startup has no service owner | Ask the user to stop/start through the original launcher |
| Apple container is stopped during recovery | Start the same container before copying the host backup |
| Recovery input or Apple runtime remains unavailable | Retain backup and report incomplete recovery |

## 5. Good / Base / Bad Cases

Good: install v0.2.19 with an exact backup and an explicit user decision for version warnings.
Base: a container recreated from its original image needs the same package reapplied.
Bad: compile upstream sources or require nc/buildx/helper images on the user's host.

## 6. Tests Required

The default build runs discovered tests and focused vet, then compiles the standalone server and installer.
Keep host tests separate from cross-compilation environments.
Keep aligned `tests/test_build.py`, `tests/test_installer_test.go`, and `tests/support.py` copies in each complete build kit.
Standalone builds run common and business cases once. Combined builds run plugin business cases and the parent's common cases once.
Use the same fake-Docker and local HTTP installer cases for transaction, version, readiness, and restoration coverage.
Capture expected failure output in the existing build fixture. Reveal it when an assertion unexpectedly fails.
Extend the existing Docker-helper case for single error output. Do not add a separate output suite or mirror progress strings.
The public readiness deadline stays fixed at 120 seconds. Standalone builds require no parent checkout.
Extend existing cases for changed native/Apple control and recovery, host/payload selection, locale precedence, and plain/cancellation behavior.
Do not multiply shared cases by installation form, language, rendering mode, or target.
Inspect actual terminal selection, spinner-to-prompt transitions, and cancellation without a snapshot suite.
For progress rendering changes, include a long readiness wait in a narrow terminal. Short prompts do not establish redraw correctness.
Build final platform pairs and retain their platform paths. Cross-compilation does not prove platform runtime installation.
Historical runtime reviews below record earlier container and browser evidence. They do not add current migration acceptance gates.
The [implementation](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/fresh-native-implementation.md) records actual Docker replacement and stopped-container restoration.
The [independent review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/fresh-native-check.md) records interruption handling and final rebuilt artifacts.
The [correction review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/fresh-native-corrections-check.md) records real TCP reset recovery, changed Docker port binding, bounded timeout and stopped-container restoration.
The [browser/model review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/native-browser-model-check.md) records rebuilt native frontend and actual Chrome checks of the final packaged artifact.
The [installer/discovery review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/installer-disabled-model-check.md) records startup output, retained terminal errors and actual Docker default `/tmp` backup/restoration.
The [v0.2.17 review](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/native-v0217-integrated-check.md) records the final build and package after the settings-load repair.
The [v0.2.17 browser check](../../tasks/archive/2026-10/10-11-sub2api-runtime-patch/research/native-v0217-browser.md) records actual settings controls and the rebuilt candidate's retained-container smoke.
Reuse passed evidence when inputs and affected behavior remain unchanged.
No standalone cleanup tests, repeated pristine baseline, architecture matrix, or additional full-suite requirements apply.
Do not add naming compatibility cases, generated-repository fixtures, separate workflow-lint gates, or deliberate remote failure runs.

## 7. Wrong vs Correct

Wrong: report a restored executable before checking the restarted application.
Correct: complete readiness checks before reporting successful restoration.
