# Research: Source Resolver Docker Failure

- Query: Attribute the step-6 generic Docker failure, trace source-resolver bootstrap/network/cleanup, and define bounded automatic remediation and validation.
- Scope: internal
- Date: 2026-10-13

## Findings

### Failure attribution

For the inspected `install.py`, the reported `Cause: Docker could not complete the selected operation.` at `Resolve source hints through Git` maps to a nonzero exit from the first command in `resolve_source`: `docker pull alpine:3.23` (`install.py:345-350`). It runs before the resolver container is named, before the `try/finally`, and before any resolver log is created. No other nonzero `self.docker` call escapes during this step: later `docker rm -f` errors are caught and ignored (`:384-389`). `docker()` maps nonzero exits to this generic message and discards captured stdout/stderr; missing Docker or a 120-second timeout instead produces `Docker did not finish the selected operation.` (`:180-186`). Thus the output identifies the failing command, not whether the daemon saw DNS, registry auth/rate limit, TLS, daemon availability, or another pull error.

An in-memory subprocess reproduction returned status 1 with fake stderr `denied: registry DNS failure`; `resolve_source()` emitted exactly the generic cause and recorded only `['docker', 'pull', 'alpine:3.23']`. It issued no later command. The current repository `install.py` and the matching resolver flow in this task use the same command order.

### Resolver dependency and network chain

`docker pull alpine:3.23` is unconditional; this path does not inspect or admit a locally cached helper image first. A cached tag therefore cannot avoid this registry operation. The subsequent `docker run --name <random> --entrypoint sh alpine:3.23 ...` has no `--platform` or `--network` option, so the helper uses the Docker daemon's default execution platform/network, independent of the selected Sub2API target architecture (`install.py:349-369`).

The helper installs Git at runtime with `apk add --no-cache git`, then clones the validated HTTPS source using `git clone --filter=blob:none --no-checkout`. `apk` package access and the clone occur inside the helper container and can have different egress/DNS behavior from host-side release download or daemon-side registry pull. Package-install failures, checkout failures, and Docker run failures other than the explicitly mapped exits use `Source resolution could not complete in Docker.`; clone failures map to the source-download message (exit 21), unresolved refs to exit 22, and missing `backend/go.mod` or `frontend/package.json` to exit 23 (`:351-383`). The build later clones the selected source again and again installs Git in its Go builder (`adapter/Dockerfile:1-13`; standard Docker fallback at `install.py:442-465`).

On this research machine, Docker CLI/server are 29.9.0/29.8.2 and the daemon reports Linux/aarch64. The cached `alpine:3.23` image is Linux/arm64. An owned `docker run --rm --pull=never --network none alpine:3.23 ...` confirmed Alpine 3.23.6 has `/bin/sh` and `apk`, but no Git. A disposable network-disabled `apk add git` emitted an Alpine repository DNS failure. This proves the current helper has a second outbound-network dependency; it does not establish the user's daemon platform or network cause. The resolver pull/run do not receive the target's `amd64` platform, so the reported target architecture alone does not explain the step-6 pull error.

### Error and transaction behavior

When `docker run` returns, `stdout+stderr` are written to `source-resolution.log`; the run uses a 240-second timeout. The `finally` block attempts `docker rm -f <random-name>`, but swallows cleanup errors, preserving the primary error while possibly leaving a stopped helper container. A timeout or process-launch error occurs before the log write. `install()` places the build workspace in `TemporaryDirectory`, so any source log is removed when installation exits (`install.py:367-389,611-615`).

Source resolution/build occurs before `protect_root`, journal save, backup, staging, or binary replacement (`install.py:632-660`). On this failure path `mutated` remains false: no selected-container filesystem, persistent data, journal, or application binary write occurs. The attempted helper image pull may update Docker's local image cache. `main()` reports the step, generic cause, “No changes were made,” and a generic retry instruction, but no Docker stderr (`:805-835`).

A second in-memory reproduction let the helper pull succeed, returned an `apk`-like status 1 from `docker run`, and made cleanup fail. The cause became `Source resolution could not complete in Docker.`; captured resolver stderr existed only in the temporary log path, and cleanup's own error was not surfaced. The first failure was separately reproduced at the pull boundary above.

### Recommended bounded fix

Use a generic Git-capable resolver tool image and make helper admission separate from target compatibility: inspect the exact helper reference once; accept a cached image only when it is Linux, runnable by the Docker daemon, and has the required shell/Git tools. If absent or unsuitable, pull that same generic helper reference with a short bounded retry for transient failures, then re-inspect it. Do not require the helper architecture to equal the selected target architecture, require target image digests/descriptors, or add host-specific tool facts. The current Alpine helper can remain only if Git is already included; relying on runtime `apk add` preserves a second package-repository dependency.

Report registry acquisition separately from in-container source resolution. Preserve the Docker exit code and a bounded, control-character-stripped, credential-redacted stderr excerpt in the user-visible failure; do not print raw environment/configuration or an unbounded log. Retain a private diagnostic file only when it is intentionally kept after failure; the current temporary workspace is deleted. Retry only a small fixed number of transient registry/network failures (for example temporary DNS, timeout, 429/5xx); fail promptly on authentication, missing reference, or unsupported helper platform. Do not assume a registry mirror or force host networking. Never fall back to a fixed v0.2.15 source/binary or reinterpret a failed hinted ref as proof of compatibility. If resolver cleanup fails, keep the primary cause and report the secondary cleanup status and generated helper name; retry cleanup automatically within a bound.

### Selected policy and host-Git alternative

The current fix should keep source resolution inside Docker: admit a cached Git-capable helper first, acquire that same helper with bounded diagnostics only when it is unavailable, and preserve the existing container network boundary. Do not install Git on the host or silently switch to host networking when the helper cannot be acquired. The helper must provide the shell and Git tools in its image; runtime `apk add git` remains a second network dependency and is not a reliable acquisition strategy.

A host-Git fast path is technically feasible but remains an optional future alternative, not part of this fix. If a later revision adopts it, use `shutil.which("git")` without installing anything, clone with argv-only commands into a disposable private workspace, and isolate it with a fresh mode-`0700` `HOME`, `GIT_CONFIG_NOSYSTEM=1`, empty global/system config paths, an empty `credential.helper`, `GIT_TERMINAL_PROMPT=0`, and a failing `GIT_ASKPASS`; pass no credentials. Bound clone, ref verification, checkout, and metadata commands, redact control characters/credential-like text in bounded diagnostics, and remove the checkout on success or failure. A present-but-failing host Git process must report its failure rather than retry through Docker or silently choose `HEAD`; `HEAD` remains an explicit candidate only when `source_refs()` returns `['HEAD']` because no valid hints exist (`install.py:65-71,347-348`). The arbitrary HTTPS source label would still grant the host process host-network reach, so this alternative requires an explicit trust-boundary decision.

The selected Docker path and any future host path must freeze the same full `git rev-parse HEAD` and read the same `backend/go.mod` module/toolchain values. Discovery is separate from the actual native checkout and compile, which still occur in Docker from `SOURCE_URL` plus the frozen `SOURCE_REVISION` (`install.py:470-498`; `adapter/Dockerfile:7-12`); `adapter/build.py:52-57` rejects a checkout whose full revision differs. Source resolution remains before protected target writes.

### Adjacent builder pulls and hidden failure surfaces

Resolving source with host Git does not remove later network/tool dependencies. `verify_bundle()` unconditionally pulls `golang:<version>-alpine` and `node:24-alpine` with the selected target platform (`install.py:482-490`) before Buildx or the standard builder starts. Those pulls use the same `docker()` wrapper, so registry/auth/DNS/platform errors still lose stderr and become the generic Docker message. Buildx then pulls the declared images as needed for `adapter/Dockerfile:1-5`; its builder installs `python3 git` with `apk` (`:5`) and clones the source again (`:7-12`). The standard-Docker fallback performs the same `apk add --no-cache python3 git` and source clone inside the native container (`install.py:442-465`, especially `:458-460`). Frontend `corepack prepare` and frozen `pnpm install` in `adapter/build.py:64-72` add further dependency-network requirements.

The changed diagnostic/cleanup scope must cover these adjacent stages as well as the first resolver pull:

- `docker()` (`install.py:180-186`) suppresses stderr for every pull/create/start/exec/remove failure and maps command timeouts to one generic message.
- `resolve_source()` bounds its `docker run` to 240 seconds but deletes its only resolver log with the enclosing temporary workspace; cleanup errors are swallowed (`:367-389`).
- `build_process()` (`:434-440`) has no overall timeout and can wait indefinitely while printing periodic progress. A failed Buildx/standard build can therefore outlive the intended bound.
- Buildx cleanup (`:507-511`) and standard-builder cleanup (`:466-468`) are not isolated from the primary build result; a `buildx rm`/`docker rm` failure can replace the actual source/build error. Resolver cleanup intentionally does the opposite and silently drops its secondary error. The implementation should preserve the primary failure and report bounded secondary cleanup status in both cases.

These are failure-reporting and bounded-cleanup corrections around the existing source-driven build, not permission to add a fixed source, host binary, or target-specific fallback.

### Validation required

- Subprocess-level fake-Docker test: make `docker pull <helper>` return nonzero with recognizable stderr; assert the CLI identifies the helper-pull stage, emits only bounded/redacted diagnostics, does not invoke resolver `run`/`rm`, and performs no target writes. Cover cached-image admission without a registry call, one transient retry, and permanent auth/platform failure without unbounded retry.
- Host-Git path test if adopted: Git present resolves only the declared refs into a disposable workspace without invoking the helper; Git absent automatically uses the cached/pulled helper; a present but failing Git process reports its own bounded error without HEAD/default-source retry. Test with ambient credential helpers/config and malicious local/redirect URLs to prove the host path cannot read credentials or reach host-local services.
- Resolver-process tests: exercise package/tool startup failure, clone exit 21, unresolved hint exit 22, missing native files exit 23, Docker-run status 125, timeout, and cleanup failure. Assert source diagnostics survive long enough to explain the cause and cleanup cannot replace it.
- Real Linux/amd64 Docker source-resolver check: use a helper that already has Git and verify default-branch plus commit/tag resolution against the exact selected source URL; assert the full 40-character revision and Go/toolchain values. Exercise container-network failure separately, verify bounded cleanup, and assert no target writes. The available research daemon is arm64, so no real amd64 source-resolver run was claimed here.

### Files found

- `install.py:180-189,345-389,611-660,805-835`: Docker wrapper, source resolver, pre-write install order, temporary workspaces, and CLI errors.
- `adapter/Dockerfile:1-13`: later builder installs Git and clones the source again.
- `tests/installer/test_preflight.py:180-198`: resolver-success and exit-21/22/23 tests; the mocked Docker pull always succeeds, so pull-wrapper stderr is not covered.
- `tests/installer/test_universal.py:142-170`: failed source build must precede target state and standard builder containers are removed; no resolver pull failure is covered.
- `.trellis/spec/backend/universal-linux-delivery.md:61-95,120-130`: source/build and pre-write capability contracts; no fixed source downgrade.

### External references

None required; command and behavior findings are from local source/tests and the observed Docker CLI/server versions above.

## Caveats / Not Found

- The user's one-line stderr cannot reveal the underlying registry or Docker-daemon cause. The command attribution assumes the executed installer matches the inspected current `install.py`; it does not diagnose the remote host's registry, auth, DNS, proxy, platform, or Docker configuration.
- The local network-disabled `apk` check is a disposable reproduction of the package-manager network dependency, not the user's failed source-resolution run. No deployment, paid request, target-container mutation, source downgrade, or real amd64 resolver execution was performed.
