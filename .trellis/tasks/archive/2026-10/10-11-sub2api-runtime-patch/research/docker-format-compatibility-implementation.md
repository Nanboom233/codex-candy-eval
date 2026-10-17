# Docker CLI Formatting Compatibility

## Scope and cause boundary

The reported output is `unknown flag: --format`, followed by `See docker --help.`
The exact remote command and Docker version remain unknown.
Both product entry points used formatted `docker ps`. Cleanup also used formatted Buildx listing.
A real executable fake Docker reproduces the exact output and exit125 for both formatted command classes.
This establishes subprocess compatibility behavior. It does not establish acceptance on the user's iStoreOS host.

## Implementation

`install.py` and `cleanup.py` enumerate full container IDs with `docker ps -a --quiet --no-trunc`.
Existing JSON inspection supplies container facts. Explicit selection and immutable selected-ID rechecks remain active.
Neither production entry point invokes `--format`.

Cleanup detects the optional `buildx` command from successful `docker --help` output.
It recognizes the ordinary command listing and the starred CLI-plugin listing.
Absent Buildx prints a retain-builders result. Container, image and temporary-workspace inventory continue without invoking the missing command.

Available Buildx uses `docker buildx ls` without formatting options.
Only exact owned builder names under the known table header and `docker-container` driver can be eligible.
Unknown headers/drivers, duplicate names, indented ambiguous rows, missing node containers and active builders stay retained.
The existing process-activity, exact selection and refreshed inventory guards still apply before removal.
Shared builders and business containers are excluded. No global prune was added.

Command failures remain fatal. Raw stdout/stderr are relayed once and original exceptions retain their byte fields.
CLI entry points preserve full tracebacks. Invalid inspect JSON remains fatal.
Cleanup additionally relays stderr from successful commands, preserving Buildx diagnostics even when Buildx returns zero.

## Changed files

- `install.py`: quiet full-ID enumeration.
- `cleanup.py`: quiet full-ID enumeration, optional Buildx detection, conservative table parsing and successful stderr relay.
- `tests/installer/test_cleanup.py`: ordinary Buildx table and absent-plugin fixtures.
- `tests/installer/test_docker_compatibility.py`: nine real-subprocess compatibility tests.
- `tests/installer/current_suite.py`: includes the new compatibility cases.
- Canonical `dist/build/sub2api-codex-eval.tar.gz`: repackaged reviewed runtime and current Python entry points.
- This report and `docker-format-compatibility-validation.json`: delivery evidence.

## Validation

`python3 -B -m unittest discover -s tests/installer -p test_docker_compatibility.py -v`: nine passed.
Coverage includes exact error reproduction, explicit immutable selection, missing choice, identity drift, optional missing Buildx,
legacy tables, ambiguous/running builders, successful stderr diagnostics, fatal daemon/permission/help/Buildx errors and invalid JSON.
Failure cases emit complete 12KB stdout/stderr sentinels exactly once plus the traceback.
Fake executable logs establish inventory-only operations and no product `--format` calls.

`python3 -B tests/installer/current_suite.py`: 108 cases, 104 passed, four explicit Docker skips, exit0.
Log: `/tmp/docker-format-current-suite.log`.

`python3 -B scripts/package-sub2api.py`: default canonical packaging passed.
The maintained Python writer stages in system temp and atomically replaces the archive.
It preserves historical archives and does not rebuild the runtime.

The actual canonical archive passed installer extraction and packager validation.
Raw members contain exactly eight regular files and the `LICENSE` directory.
No duplicate paths, links, extended metadata, AppleDouble files, sources, tests or sidecars are present.
Extracted installer, cleanup and all license bytes match maintained source.
Extracted runtime matches the reviewed canonical raw binary.
Both archived `--help` entry points passed. Archived real-subprocess compatibility suite passed nine cases.
Log: `/tmp/docker-format-archive-compatibility.log`.

Current real Docker quiet-ID listing and unformatted Buildx table succeeded read-only.
The existing unrelated desktop builder was retained.
Buildx returned zero while reporting an unavailable default-builder daemon endpoint on stderr.
The new stderr behavior preserves this diagnostic. No real Docker resources were mutated.

## Artifact identity and remaining limits

Canonical archive: `dist/build/sub2api-codex-eval.tar.gz`, 39,016,308 bytes.
SHA256: `39c876d274a6499479bf0fbcbcf8d93c46c52053b01b2dbbc6ee3b497554f6fc`.

Canonical runtime: `dist/build/sub2api-codex-eval`, 125,632,672 bytes.
Unchanged SHA256: `eed01f71c39048eef1a15fbf13652b5c44d9915fff1e34c59f3d58abf1c6c9e1`.
The internal runtime filename remains `sub2api-v0.2.16-linux-amd64`.

No new runtime/frontend build, broad native tests, deployment, paid requests, commit, publication or cache cleanup occurred.
Required original-checker review remains next. Main owns task/spec synchronization.
The real user host, native N150/iStoreOS and direct Python3.9 execution remain unverified.
Historical broad native unit-race failures still block publication. This Python fix does not alter those results.
