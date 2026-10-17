# Research: Interactive Selection and Source Preflight

- Query: Diagnose the terminal and native source-identity failures, then map identity gates that must be removed while preserving robust build, install, readiness, and recovery behavior.
- Scope: mixed
- Date: 2026-10-13

## Findings

### Interactive selection

`install.py:198-228` prints every container and then opens `/dev/tty` as `r+`. The parent agent reproduced this in a real controlling PTY: Python raises `io.UnsupportedOperation` for `open('/dev/tty', 'r+')`; it is caught through `OSError` and reported as “A terminal or --container is required.” The existing test at `tests/installer/test_universal.py:41-46` substitutes a seekable `StringIO`, so it cannot reproduce this PTY failure. Python defines `r+` as read/write update mode; separate `/dev/tty` read and write handles avoid requiring one bidirectional buffered stream. This PTY cause is confirmed by the parent agent's experiment, not independently rerun here.

Keep explicit selection and cancellation semantics. Open `/dev/tty` separately for output and input, flush the prompt before reading, require a valid number or explicit cancel, and treat EOF as cancellation. Do not consume `sys.stdin` as the selection channel because the installer may itself be launched from piped bootstrap source. Without a controlling terminal, require `--container` and fail before any target mutation. Add a real `pty.fork` regression, alongside piped-stdin/non-TTY no-mutation, explicit-container, invalid choice, EOF, and cancel cases.

### Source identity error

`install.py:265-295` copies the running executable, requires inline Go metadata, requires a parseable `Sub2API ... (commit: ..., ...)` response, and rejects exact conflicts among `vcs.revision`, OCI revision, version-output commit, and OCI version. The reported error can come from checks at lines 276, 281, or 284; without the user's actual image label and binary values, the triggering pair is unknown. Do not claim a specific packaging defect from the error alone.

The official Sub2API v0.2.15 `.goreleaser.yaml` uses the same `.Commit` value for `main.Commit` and `org.opencontainers.image.revision`; its GHCR entries also set `org.opencontainers.image.source` and version. The tagged `main.go` prints `Version`, `Commit`, and `Date` for `-version`. Therefore a same-checkout official build is expected to agree, but this does not justify making those hints mandatory for local, forked, copied, or custom images. Go defines `vcs.revision` as the source checkout's revision, and `-buildvcs=false` omits VCS stamping. OCI defines `source` and `revision` as descriptive metadata; those labels do not prove which executable bytes are running.

### Bounded source resolution

Use a finite best-effort hint order, then make the exact resolved checkout the build input:

1. Parse `-version` if supported. Treat version, commit, and build date as hints; tolerate missing output and nonstandard formats.
2. Read OCI `source`, `version`, and `revision` labels independently when present. Ignore malformed or disagreeing labels as identity evidence; retain a usable repository URL or version as a hint.
3. Read Go module path and VCS fields only when inline metadata is available. Treat them as optional repository/revision hints; never require Go build metadata or compare it against labels or `-version`.
4. Derive a repository from a usable OCI source URL, then a usable module path. Resolve revision in order from a valid full commit, a short commit resolved by Git, then a version tag/ref resolved by Git. For a valid checkout, freeze its full `git rev-parse HEAD` and pass that commit to the builder. The builder's checkout-vs-requested-commit check (`adapter/build.py:52-62`) is build-input reproducibility, not installed-binary attestation, and may remain.
5. Take the Go toolchain version from the resolved source's `backend/go.mod` or the existing builder contract, not from binary metadata. Do not insert a fixed host-version fallback or ask the operator to supply per-machine facts. If hints cannot resolve a repository and commit, or checkout/build fails, return the concrete source-resolution, download, or build error before changing the target.

Never silently substitute the pinned compatibility image/source as a fallback for an unresolved target. Do not require a binary SHA, `/proc` executable identity, source/version equality, exact OCI label match, or image/config digest equality to continue, activate, or recover.

### Mandatory binary gates and operational behavior

- `install.py:241-304`: remove `/proc/1/exe` as a mandatory source of binary path/identity; remove `go_build_info`, `-version`, platform metadata, source/version/OCI strong-match, and `binary_sha()` rejection paths. Retain actual selected-container identity, supported Linux image architecture for the new build, process UID/GID, billing mode, persistent writable data location, and safe regular-file/parent/symlink/writeability checks needed for file replacement. Resolve the replacement path from available process/container startup facts; if no safe target path can be found, report that concrete operational failure.
- `install.py:440-494`: remove the `--binary-sha256` maintenance argument, response hash/version/source comparisons, on-disk and `/proc/1/exe` SHA gates, and restart source/version equality. Keep container-running/health checks, localhost `/health`, startup retries, maintenance capability/state, deployment ownership, and repair-state handling.
- `install.py:496-520, 548-588, 602-608, 664-671, 684-718`: remove original/patched/current/staged binary hash comparisons and image-runtime/config digest gates that block idempotent install, status, retry, rollback, or uninstall. Preserve journal phase transitions, protected original backup, staged file replacement, atomic rename, permissions, selected-container checks, drain/revoke ordering, retained backup on failure, and actionable recovery. Keep hashes used to validate the installer release payload separate from target-binary identity.
- `adapter/runtime/backend/internal/service/eval_runtime_ports.go.in:14-95`: remove SHA format/current-binary hashing via `/proc/self/exe`, plus hard readiness comparisons on `EvalHostVersion`, source revision, host image/runtime IDs, image config digest, and target binary SHA. Those checks currently feed `Ready` and `Activate` via `eval_host_runtime.go.in:104-121,151-185`, which can leave a successfully replaced binary impossible to activate. Keep protected installation ownership/deployment ID and functional runtime/store/account readiness.
- `adapter/runtime/backend/internal/service/eval_maintenance.go.in:21-39,48-62,176-208` and `adapter/runtime/backend/cmd/server/maintenance_client.go.in:14-29`: remove binary SHA from the local maintenance request/response and equality gates, plus source/version/image identity gates reached through current runtime facts. Keep peer credentials, socket path/permissions, bounded protocol, deployment ID ownership, `deactivate`, drain completion, and owned-Key revocation. These gates currently can block status and make rollback/uninstall unable to drain or revoke.
- `adapter/build.py:63-66,84-91`: stop embedding image runtime ID, image config digest, and detected host version as facts that readiness must match. Data directory remains runtime configuration. Source revision may be reported as build provenance, but must not attest the deployed executable or gate runtime readiness.

Hash fields in existing journals/installer facts should become optional migration-era data, not requirements. Status must not infer `binary_conflict` from a missing or differing target SHA; use journal phase plus maintenance capability/readiness and provide a recoverable state when the runtime is unavailable. Restore from the protected backup based on owned journal, backup presence/type/permissions, selected target path, and startup/health outcome; do not make digest equality a restoration prerequisite. Image/config digests and source labels may be displayed as hints but cannot block a retry or restoration.

### Files and patterns found

- `install.py`: selection, source resolution, transactional install, health, status, rollback, and binary gates.
- `adapter/build.py`, `adapter/rules.json`, `adapter/runtime/backend/`: source adaptation, early maintenance dispatch, readiness, and maintenance gates.
- `tests/installer/test_universal.py`, `tests/host-overlay/backend/`, `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch`: terminal/identity expectations and retained fixed-host checks.
- `.trellis/spec/backend/universal-linux-delivery.md`, `runtime-installation.md`, `patch-native-installation.md`: superseded binary-attestation contracts.

### External references

- Sub2API v0.2.15 [GoReleaser config](https://github.com/Wei-Shaw/sub2api/blob/v0.2.15/.goreleaser.yaml), lines 658-666 and 786-822: `main.Commit` and OCI revision use `.Commit`; OCI source/version are also set.
- Sub2API v0.2.15 [server main](https://github.com/Wei-Shaw/sub2api/blob/v0.2.15/backend/cmd/server/main.go), lines 630-640, 675-686: version fields and `-version` output.
- Go 1.27.2 [runtime/debug build settings](https://go.dev/src/runtime/debug/mod.go), lines 193-196: `vcs.revision`, `vcs.time`, and `vcs.modified` describe the source tree.
- Go [build command documentation](https://go.dev/src/cmd/go/internal/work/build.go), lines 232-239: VCS stamping behavior and `-buildvcs=false`.
- OCI Image Specification [pre-defined annotations](https://github.com/opencontainers/image-spec/blob/main/annotations.md): source and revision are descriptive annotation values.
- Python 3.11 [built-in `open`](https://docs.python.org/3.11/library/functions.html#open), lines 583-598: `r+` means read/write update mode.

## Caveats / Not Found

- The actual runtime metadata and labels from the user's selected container were not available, so the exact value pair causing the reported source-conflict branch is unknown.
- Grok Search returned HTTP 403 authentication errors for the Sub2API and terminal searches after its configuration check passed. Those two searches were completed with the configured web-search fallback; official Sub2API tag sources, official Go docs, OCI spec, and Python docs were retrieved. No private deployment configuration, Docker target, or `.env` was inspected.
- The PTY failure reproduction was supplied by the parent agent; this researcher did not rerun it. The source-resolution order above is an implementation recommendation and requires tests for commit/tag resolution, ambiguous source hints, download failures, and no-mutation on unresolved source.

### Maintenance Dispatch Safety (Follow-up)

The official v0.2.15 `main()` calls `flag.Parse()` before checking `-version`, setup mode, and `runMainServer()` (`main.go:669-737`). Go flag parsing stops at the first non-flag argument, so an unpatched binary invoked as `sub2api codex-eval-maintenance status ...` ignores later flags, then may enter setup/autosetup or normal server startup (`initializeApplication`, Wire/DB). The adapter rule inserts the maintenance dispatch before `logger.InitBootstrap()` (`adapter/rules.json:4-11`), so that command is safe only through the known adapted client path. Never probe an unknown/original server binary with the maintenance subcommand.

The `control.sock` path, owner, and mode alone do not prove a live process. Require a bounded `status` response with the expected protocol and matching deployment ID; this proves a live capability, not binary identity. Keep peer-credential and deployment-owner checks; on connect/protocol failure, retain the journal/backup and report retryable cleanup unavailability. For deactivation, require `drained` and `owned_key_revoked`. Use a known adapter or dedicated protocol client; never fall back to starting the server binary.

Sources: [Sub2API v0.2.15 `main.go`](https://github.com/Wei-Shaw/sub2api/blob/v0.2.15/backend/cmd/server/main.go), [Go flag syntax](https://pkg.go.dev/flag#hdr-Command_line_flag_syntax); local dispatch `adapter/rules.json:4-11`, socket `eval_maintenance.go.in:85-124,160-208`.

### Native HTTP Application-Capability Probe

Compared `/tmp/codex-robust-preflight-20261013/current-source/backend` and `/tmp/codex-robust-preflight-20261013/previous-source/backend`.

Both supplied source trees register the same unauthenticated root endpoints in `backend/internal/server/routes/common.go`:

- `GET /health` returns HTTP 200 JSON `{"status":"ok"}` (`:12-14`). This is useful as a liveness check but is too generic to identify the application by itself.
- `GET /setup/status` returns HTTP 200 JSON `{"code":0,"data":{"needs_setup":false,"step":"completed"}}` (`:21-31`). It is a stronger functional marker for the normal Sub2API router.

`backend/internal/server/router.go:58-93,116-130` applies request/context/logging, CORS, security-header, timing, and optional frontend middleware globally, then registers the common routes directly on the root Gin engine. Authentication is added to later API groups, not these two handlers. The embedded frontend middleware in both revisions bypasses `/setup/` and `/health` (`backend/internal/web/embed_on.go:355-363`), so it does not replace their responses. The handlers return constant literals and do not query configuration/database state, disclose private data, or mutate application state. Normal request logging can still record these GETs.

Safe capability check: issue only bounded, unauthenticated GETs to `/health` and `/setup/status`; require successful HTTP status and parse JSON fields/types rather than matching body substrings. Require `status == "ok"`, `code == 0`, `data.needs_setup == false`, and `data.step == "completed"`. `/health` alone must not count as application identification. Do not use the telemetry POST endpoint. The anonymous `/api/v1/settings/public` endpoint is unnecessary and less suitable: its route documents a per-request database query and returns configurable public settings (`backend/internal/server/routes/auth.go:241-247`; `backend/internal/handler/setting_handler.go:37-62`).

The two supplied `common.go` and `router.go` files compare identical. `embed_on.go` differs in other route bypasses, but both revisions explicitly bypass these probe paths. This identifies a functional Sub2API-like application contract, not executable provenance; another service can imitate the response, which is consistent with R11's removal of source/version/image/binary identity gates.

The current installer probe at `install.py:285-303` checks these same JSON values with loopback `wget` and a 4096-byte response bound. Its `wget | head` shell pipeline does not independently retain the HTTP status, so the implementation should verify a 2xx status in addition to validating the JSON contract.
