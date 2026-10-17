# Research: Installer Source and Progress Checkpoint

- Query: Identify a bounded source retrieval and reuse path, live progress gaps, and focused verification for the current installer request.
- Scope: internal
- Date: 2026-10-14

## Findings

### Files found

- `install.py` — installer entry point, Git source resolver, Docker builders, release download and temporary-workspace ownership.
- `adapter/Dockerfile` — Buildx source acquisition and build steps.
- `adapter/build.py` — adapts a checked-out upstream tree and builds the native package.
- `tests/installer/test_preflight.py` — current mocked source-resolution tests.
- `tests/installer/test_universal.py` — source-hint and fallback-container builder tests.
- `tests/installer/test_docker_boundary.py` — Docker subprocess, timeout, cancellation and cleanup tests.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md` — current progress and source-retrieval requirements.
- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md` — universal package and installer behavior contract.
- `.trellis/spec/backend/universal-linux-delivery.md` — source helper, image acquisition, build and error-propagation rules.

### Code patterns

- `install.py:681-706`: `resolve_source()` runs `git clone --filter=blob:none --no-checkout` without `--depth`, in `/host`, then checks ordered hints and checks out one. The helper subprocess uses `capture_output=True`, so Git's progress is not visible until it exits. Timeout output is relayed only after timeout; failure output is relayed after a nonzero exit.
- `install.py:830-840`: the fallback Docker builder clones the source a second time and uses `/host`; this duplicates source transfer after resolution. It also installs `git` with `apk` for that clone.
- `adapter/Dockerfile:1-12`: the Buildx path independently clones `$SOURCE_URL` and checks out `$SOURCE_REVISION`, so it also duplicates retrieval.
- `install.py:779-818`: `build_process()` redirects both child streams to temporary files, emits only a generic wait message every 30 seconds, and relays captured content after failure. This does preserve complete output for `CalledProcessError` / timeout, but not live build progress.
- `install.py:296-353` and `install.py:258-270`: image pulls and Docker commands use `subprocess.run(..., capture_output=True)`, so Docker download output is buffered. `install.py:849-883` invokes Buildx with `--progress=plain`, but `build_process()` still hides it until completion/failure.
- `install.py:194-207`: `download_release()` reads bounded chunks into the archive but does not report transferred bytes. The user-facing release step therefore has no in-flight download progress.
- `install.py:849-875`, `install.py:997-1015`: host build and release scratch are already allocated with `tempfile.TemporaryDirectory`, which honors the system temp directory/TMPDIR. Container Git scratch is the `/host` exception called out by the PRD.
- `install.py:686-727`: the resolver has unique helper naming and exact-name cleanup in `finally`; its source facts are parsed from `SOURCE_*` lines after command completion. Any streaming implementation must continue collecting those lines and attach complete captured stdout/stderr and tracebacks on failures.

### Handoff direction

- Keep hint order and resolve one selected full commit. Use a bounded shallow retrieval strategy for an advertised branch/tag or supported exact commit, then check out and validate that one tree. If a short commit cannot be resolved without full history, fail within the source deadline with the underlying Git diagnostics; do not infer a network cause or fall back to `HEAD`/an older revision.
- Materialize the immutable selected tree once under container `/tmp`, copy it into the host `build_workspace`, and make both Buildx and fallback Docker builds consume that same tree as builder input. Remove the build-time Git clone; the fallback builder then no longer needs Git for source acquisition.
- Stream subprocess output while retaining full stdout/stderr for logs and exceptions. Apply the same mechanism to `docker pull`, source resolution, Buildx/fallback build, and chunked release download. Report byte counts or command output, not guessed percentages.
- Continue using `TemporaryDirectory` for host artifacts and `/tmp` for disposable container paths. Keep source resolution/build failures before target writes and preserve current exact-helper cleanup semantics.

### Focused verification targets

- Extend resolver tests for full SHA, short SHA, tag, no hints, bounded unresolved short-SHA failure, source facts parsing, live output before process exit, timeout output retention and cleanup.
- Verify both Buildx and fallback builds consume the prepared source snapshot and make no second upstream Git request; assert no source scratch path is `/host` and fallback does not install Git solely to clone source.
- Verify output from a slow fake pull/build and a chunked download becomes observable before process completion, while complete stdout/stderr, log files, timeout/cancellation behavior and primary errors remain intact.
- Run focused installer tests directly as the repository does: `python3 tests/installer/test_preflight.py`, `python3 tests/installer/test_universal.py`, and `python3 tests/installer/test_docker_boundary.py`. A live remote source resolution is a separate integration proof and must not be represented by these mocks.

### Related specs

- PRD current requirements: `prd.md:182-194` — live progress, selected source only, no full history or duplicate clone, standard temp paths, preserve failures.
- Design: `design.md:28-35, 61-69` — build/progress messages, exact source/tool identities, hint handling.
- Universal delivery: `universal-linux-delivery.md:115-129` — helper preflight, transient retry boundary, complete original diagnostics, process-group timeout and cleanup.

### External references

- None consulted. The shallow-fetch plan is a handoff direction, not a claim that every Git server supports resolving arbitrary unadvertised short object IDs.

## Caveats / Not Found

- The previously summarized `installer-progress.md` was not present in the shared task directory, so this report was grounded in the current PRD, specs, installer and tests.
- No production target, `.env`, external Docker daemon or remote Git retrieval was accessed. The cause of the reported 240-second timeout remains unknown.
- A short SHA may identify an object that is not directly advertised by the remote. The implementation must fail boundedly when that object cannot be resolved without fetching history; do not silently substitute `HEAD`.
- The focused test command above has not been run in this research task.
