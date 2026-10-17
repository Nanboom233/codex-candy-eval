# Source resolver and raw installation errors

## Current checkpoint

Raw-output revision implemented and verified. Full root check, archived preflight and archived boundary/CLI checks all completed with observed exit0. Final read-only archive/source parity and `git diff --check` pass. Ready for the original independent checker; no cleanup has run.
The latest user instruction supersedes earlier bounded/redacted diagnostics. Failed commands relay complete original stdout/stderr, retain original argv/output/stderr in standard subprocess exceptions, and print full Python traceback with chains. No diagnostic truncation, redaction or Cause/Recovery/Next error explanation remains.
Main owns task/specs. Original implementer performed the changes without child agents. No production deployment, paid request, commit, push, publication, task archive or evidence cleanup occurred. No .env was read. Evidence remains at `/tmp/codex-source-resolver-fix-20261013` for the original checker.

## Concrete implementation

- Use verified `alpine/git:2.49.1`, with Git already installed. Resolver does not run apk.
- Probe cached Linux shell/Git with `--pull=never --network none`; helper execution platform is independent of the selected target platform.
- Acquire unavailable tools automatically, with at most two attempts for transient DNS/timeout/reset/rate-limit/server failures. Auth/denied/certificate/missing-manifest/platform errors take precedence and stop promptly.
- Keep ordered source hints. Resolve full/short commits and tags to a concrete checkout. Failed hints never fall back to HEAD or an older source.
- Probe Go/Node executables before registry access; accept multiple valid RepoDigests and local image IDs.
- Prepare a dedicated Buildx builder. Absent/unusable preparation selects standard Docker before compilation; a failed compilation never changes builders.
- Keep Docker/source/build deadlines of 120/240/1800 seconds. Build cancellation/deadline terminates the local process group, then cleans owned Docker resources.
- `command()` and resolver nonzero exit raise `CalledProcessError` with original argv and both original streams. Timeouts retain `TimeoutExpired` and captured partial streams. Startup, download, extraction and JSON failures preserve original exceptions.
- Build output spools stdout/stderr separately without pipe deadlock; the complete combined `builder.log` is retained, and failed streams relay on their original channels without truncation. Build exceptions retain separate complete stdout/stderr.
- Main prints standard complete traceback and returns nonzero. Contextual validation errors remain ordinary exceptions. Normal progress/status messages remain.
- Cleanup retries do not replace the primary error or stop cleanup of other owned resources. Its final traceback includes the generated resource in the original command, without a custom explanation block.
- Post-replacement failures restore and verify the original, save `failed_restored`, then re-raise the identical primary exception. Secondary restoration/stage cleanup errors print their traceback and cannot claim successful restoration. Stopped-container copy/start also restores executable UID:GID and mode.

## Final raw-output verification

Observed exit 0:

| Command | Result / retained log |
| --- | --- |
| `scripts/test-sub2api.sh check` | `repository-check-raw-final.log`: 72 helpers (one historical skip), 49 transactions, 8 universal, 16 preflight, 14 Docker boundary; Go vet/race, shell and UI checks pass |
| Root focused transactions | 49 pass, `raw-transactions-final.log` |
| Root focused universal / preflight | 8 / 16 pass, `raw-universal-final.log`, `raw-preflight-final.log` |
| Fresh packaging and real extraction | `raw-final-package-verification.json/.log`: strict allowlist, 62 regular files, all source/license bytes identical; functional identity unchanged |
| Archived real PTY / preflight | 16 pass, `archived-raw-preflight.log`, inherited `PYTHONDONTWRITEBYTECODE=1` |
| Archived boundary / full CLI | 14 pass, observed exit0, `archived-raw-boundary.log`; extracted installer loaded and full CLI invoked |

Executable fake Docker is launched on PATH through actual subprocesses, rather than only mocked Docker methods. Cases cover warm/cold cache, transient/permanent pulls, mixed-error precedence, final executable failure, clone/ref/files/checkout/Go declarations, Docker exit125, timeouts, cleanup, multiple digests, local IDs, Buildx preparation/fallback and standard container cleanup.
Raw tests assert exact stdout/stderr longer than240 characters (including Unicode and ANSI), a build stdout over10000 characters plus complete stderr, original command/exit/exception stream attributes, nested exception traceback and absence of custom error blocks. Actual Python build subprocesses verify separate complete streams/log bytes and SIGINT process-group termination. Full packaged CLI tests cover acquisition denial, clone/cleanup failure and cancellation without any selected-container write.
Transaction tests inject KeyboardInterrupt, concrete validation and `CalledProcessError` after replacement; they check identical primary exception, restoration order and journal phases. A secondary restoration exception preserves its entire nested traceback and never records `failed_restored`.
The initial raw preflight run failed one stale expectation: malformed JSON now correctly raises original `JSONDecodeError` instead of the former translated `InstallError`. The assertion was updated; final root and archive runs pass. That failed run is retained in `raw-preflight-initial.log` and is not counted as a pass.

## Prior actual Docker evidence, retained within its scope

Before the final raw-output-only revision, actual helper manifest/pulls and network-disabled Git probes passed. Real multiarch helper index digest: `c0280cf9572316299b08544065d3bf35db65043d5e3963982ec50647d2746e26`. Native arm64 reported aarch64/Git2.49.1; emulated amd64 reported x86_64/Git2.49.1 on an arm64 daemon.
Seven real resolver runs passed: full/short commits and pure tags on both execution platforms, plus explicit no-hint HEAD on arm64. Commit/tag resolved `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`; HEAD resolved `4d9ac3e482106d2da3a39fe3e546d8d6801f1cc1`. Warm traces had no pull or target writes. Retained logs/results: `real-resolver*`, `real-extra-refs*`.
Actual arm64 Buildx and standard source builds both exited0, with logs/profile/output/facts in `real-buildx.log`, `real-standard.log`, `real-build-buildx/`, `real-build-standard/`. Actual stopped restore fixture exited0 and confirmed running plus `1000:1000:755` in `stopped-restore-final.log` and `stopped-restore-result.json`.
These successful source/helper/build/ownership checks were not repeated for the final raw-output revision. Runtime, core and adapter/build inputs remain unchanged. Emulated amd64 is not native amd64 lifecycle acceptance.

## Installation-chain evidence matrix

| Stage | Evidence and limits |
| --- | --- |
| Selection | Final real PTY/preflight and transaction tests cover piped stdin, explicit choice, cancel and ID drift; CLI traces prove no mutation on failures |
| Release/download | Focused tests cover HTTPS, redirect/size/extraction/allowlist/source parity; published online bootstrap not exercised |
| Application discovery | Typed bounded native health/setup probes and source-hint tests pass; unrelated applications fail before mutation |
| Helper acquisition | Actual native/emulated tools and warm probes retained; final executable fake CLI covers registry failures, retries and raw diagnostics |
| Source resolution | Seven actual ref/platform runs retained; final fake CLI covers raw clone/ref/Docker/timeout failures and owned cleanup |
| Go/Node acquisition | Actual build profiles retained; final fake CLI covers cached runnable tools, digests and local IDs |
| Buildx / standard | Both actual arm64 builds retained; final subprocess tests cover preparation fallback, full logs, exit/timeout/cancellation and cleanup continuation |
| Backup/private journal | Final49 transaction tests cover ownership/journal-before-effects; no new broad native suite |
| Replacement / recovery | Final injected primary/secondary failures preserve exact exceptions and phases; actual stopped-copy fixture confirms ownership/mode and running only |
| Restart/status/rollback/uninstall | Focused final tests plus unchanged earlier18-point arm64 lifecycle evidence; full native lifecycle not rerun for this installer-only correction |
| Admin/accounting/billing/concurrency | Runtime/core/adapter unchanged; earlier native evidence stays scoped to its owner. No new production charges/concurrency claim |

## Final payload and files

Archive: `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`, release `v0.4.0`, unpublished.
Size515570 bytes; 62 regular functional/license files.
SHA256: `c0eeb0ec01121399ba19f4b7d3f897bf0a7c27494f16a07a73a3b6c930283353`.
Installer SHA256: `c66dd0f4b5c4fc2716318e7df922a6468909b39ac7c09ad519d8fc2a14831714`.
Functional SHA256: `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759`.
Every extracted file matches current sources and fresh staging. Historical eight-file archive remains at its original path, SHA256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8`.
The prior bounded-output b422 archive is superseded; only a workroot evidence copy `before-raw-output-override.tar.gz` remains. Earlier reviewed archive is also retained as review evidence, not a delivery.

Changed files: `install.py`, `README.md`, `docs/sub2api-operations.md`; maintained tests `test_transactions.py`, `test_preflight.py`, `test_universal.py`, new `test_docker_boundary.py`; check-suite routing `scripts/test-sub2api.sh`; this report. Main's independent task/spec/config changes were preserved. No adapter/runtime/core functional bytes changed.

## Earlier failed attempts and remaining limits

Earlier method mocks needed the new acquisition boundary. Initial fake CLI startup timeouts under concurrent compilers were too short; isolated traces proved two retry calls, so normal test deadlines were widened while timeout cases retain explicit bounds. Initial stopped fixture lacked its target directory; corrected disposable fixture uses existing `/tmp` paths.
Earlier parallel archived PTY imports created extracted-release bytecode; strict allowlist rejected it and three CLI cases attempted online discovery before any target write. Only that exact generated bytecode root was moved to Trash. Final archive tests inherit disabled bytecode; no payload validation was weakened. Failed logs remain separate from passes.
Required broad native race/frontend/previous integration gates remain historical failures, without a pristine A/B attribution. No iStoreOS acceptance, native amd64 loaded lifecycle, real billing/concurrency, deployed browser flow, production deployment or published online bootstrap is established. No arbitrary-future-fork compatibility guarantee is made. Original checker follows this handoff; generated evidence remains until review completes.
