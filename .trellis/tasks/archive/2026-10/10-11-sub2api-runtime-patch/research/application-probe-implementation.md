# Application probe implementation

## Final checkpoint

Implementation and affected verification complete; ready for the original independent checker. Main owns current specs/task state. No child agent, production deployment, paid request, commit/push/tag/publication/task archive or evidence cleanup occurred.
Research: [actual transport diagnosis](application-probe-transport-failure.md). BusyBox1.37 stdin EOF races the response and shell head hides nc exit; correct CRLF bytes were proven. The fix replaces that pipeline with host-controlled direct nc and owned application-listener discovery.
Evidence workroot `/tmp/codex-application-probe-fix-20261013`; research root `/tmp/codex-application-probe-research-20261013`. Both remain for review.

## Concrete code changes

`install.py` owns all changes:

- `http_get()` directly launches `docker exec -i -u0 <selected-ID> nc <numeric-local-address> <port>`. No printf/head, optional nc -w flag, proxy/redirect, new dependency/helper image or target temp file is used.
- Host sends exact GET /health or /setup/status bytes and keeps stdin open until response framing completes, client failure/EOF or deadline. Concurrent bounded stdout collection plus disk-spooled complete stderr avoids pipe deadlock and the BusyBox EOF race.
- `http_body()` caps raw HTTP at8192 bytes (one-byte overflow sentinel), decoded body4096, checks HTTP200 and identity encoding, supports Content-Length and chunked framing. It requires the declared fixed-body length, every chunk boundary and terminating zero-chunk/trailer blank line. Partial packets cannot become successful short reads; complete malformed/status/oversize errors fail promptly. Interim100 Continue is supported. Close-delimited replies require actual client EOF; no idle-delay success heuristic.
- Each probe has3-second full deadline and up to8 endpoint candidates share48 seconds, plus bounded5-second optional listener discovery/short teardown. Owned local processes close stdin, wait briefly, then TERM/KILL their process group on failure. Full captured original HTTP/stderr relay before original exception propagates. CalledProcessError/TimeoutExpired retain original argv/exit/streams; parser errors remain parser errors. Only locally initiated TERM after a complete reply is accepted; externally terminated client is a failure.
- `listening_endpoints()` reads only socket inodes of the already resolved executable under PID1 or its bounded descendant tree (≤32 processes, depth3), matched to /proc/1/net/tcp/tcp6 LISTEN rows. It reads no environment/config file contents. Only app-owned endpoints become candidates; wildcard binds map to127.0.0.1/::1, specific binds remain exact; addresses/ports are validated/deduplicated and capped8.
- Owned sockets override stale env/exposed-port hints. Without ownership, a private container namespace may use validated numeric loopback/localhost SERVER_HOST, SERVER_PORT and exposed-port hints/default8080. Non-loopback hints alone never contact an external address. Host/container-shared namespaces require owned listeners and cannot scan other namespace services. Arbitrary healthcheck strings are never executed.
- Both native endpoints must return exact typed JSON fields. Empty, disconnected, process, timeout and malformed replies keep their original concrete exceptions. A valid incompatible native response can fail its actual field contract. After multiple failed candidates, the first preferred endpoint exception is re-raised; every attempted failure traceback stays visible. No fabricated “lacks native capabilities” conclusion replaces transport failures.
- `inspect()` retains health_host/health_port in actual facts. Restart/original health verification uses the same transport/address rather than the old127.0.0.1 wget command, so specific binds and IPv6 survive the entire lifecycle.

No binary checksum/metadata/version/source/image admission was added. Runtime/core/adapter/build functional inputs are unchanged. Application discovery still precedes protected files, backup, replacement and restart. No fixture/helper code enters product payloads.

## Meaningful tests and actual results

All final commands below have observed exit0:

| Command / boundary | Evidence |
| --- | --- |
| `scripts/test-sub2api.sh check` | `repository-final.log`:72 helpers (one historical skip),50 transactions,8 universal,16 preflight,14 Docker boundary,10 HTTP tests (2 explicit Docker checks skipped); Go vet/race, shell syntax and UI pass |
| Explicit cached Docker suite `PROBE_DOCKER_TESTS=1 ... test_http_transport.py` | `http-docker-final.log`:10 tests pass, no skips; actual BusyBox nc/server/owned socket/wrapper/IPv6 plus full packaged CLI |
| Final extracted archive preflight | `archived-preflight-reviewed-candidate.log`:16 pass |
| Final extracted archive subprocess/CLI | `archived-boundary-reviewed-candidate.log`:14 pass |
| Final extracted archive HTTP with explicit Docker | `archived-http-reviewed-candidate.log`:10 pass, no skips, including final archived installer actual delayed BusyBox CLI |
| Payload extraction and source parity | `final-package-verification.json`:62 exact files, strict installer validate_bundle; final read-only recheck after archived suites passes |
| Whitespace check | `git diff --check` exit0 |

New maintained `tests/installer/test_http_transport.py` uses real local sockets and an executable subprocess relay independently modeling nc EOF behavior. Actual socket tests cover immediate and1.5-second delayed fixed/fragmented chunked responses, exact request bytes, trailers/interim framing, disconnect, malformed status, truncated Content-Length, oversized wire/body,302 rejection without extra requests, deadline, missing nc exit127, nonzero exit19/full >240-char stdout/stderr, preferred original exception, real IPv6 health and restart address, shared-network no-scan and unverified external bind rejection. Child-process PIDs must be gone afterward.
The actual cached `node:24-alpine` fixtures use BusyBox1.37 nc. Matrix covers wildcard IPv4, specific127.0.0.1, a shell-wrapper child listening on::1, nondefault18280 and1.5-second health/setup bodies using fixed/chunked framing. Owned listener discovery selects them without SERVER_PORT/exposed hints; no nc process remains. Every exact fixture is removed in finally.
Final packaged CLI also operates on an actual1.5-second delayed BusyBox fixture at18381, passes both native endpoint contracts, then intentionally stops at the fixture's missing persistent application mount before any installation effect. Before/after `docker diff` is identical. This is real probe/CLI transport acceptance, not a native Sub2API installation or billing claim.
Existing preflight mocks now test parsed transport/body boundaries instead of assuming the old shell pipeline. The fake-Docker subprocess suite handles nc stdin/response explicitly. All tests stay outside the62-file package.

## Initial failure and corrections

The first socket suite failed one assertion: first configured endpoint correctly produced RemoteDisconnected, but later default8080 refusal became the final exception. The final code retains the first preferred concrete exception while preserving tracebacks of every attempt. Initial failed `http-initial.log` is retained, not counted as a pass.
Root and initial archive checks passed before the final non-loopback fallback guard. That guard prevents unproven container-external addresses from being contacted. Only the final root/archived logs above and final archive identity below are authoritative. Superseded logs/staging/archive copies remain in the workroot for review, not in dist as deliveries.

## Final archive

Canonical `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`, releasev0.4.0 unpublished.
Size518412 bytes; exactly62 regular production/license files.
SHA256 `73d4c8b6fde68309f49d2ca7f32a761ae01e3193d81b0aff786d8dea6b2fc4db`.
Installer SHA256 `8e04d4fd11760e220873be64cb5bc3b51fa00ea203a0539c1c1489024bdf04bf`.
Functional SHA256 remains `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759`.
Every current source/license byte matches fresh staging and final extracted archive after all archived tests. No test, fixture, report, source snapshot, bytecode, host binary or separate runtime manifest is delivered.
Historical minimal archive remains73379445 bytes/8 files/SHA256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8` at its original path. Before-fix reviewed b39fe archive is only a workroot evidence copy.

Changed maintained files: install.py; docs/sub2api-operations.md; scripts/test-sub2api.sh; tests/installer/test_preflight.py, test_universal.py, test_docker_boundary.py and new test_http_transport.py; this report. Main task/spec and all unrelated source/config edits were preserved.

## Generated cleanup inventory and limits

Inventory only; no cleanup authorized/executed yet. Entire implementation root8832 KiB and research root316 KiB are task-owned logs/fixture scripts/source excerpts/staging/extracted packages/superseded archive copies. Final canonical archives remain in dist. These exact roots can be refreshed after checker completion and moved via the established /usr/bin/trash policy only on main's explicit cleanup dispatch.
Generated repository Python bytecode roots and global Go build/test cache (currently322280 KiB) require refreshed inspection after review; do not clear while checker runs. Never clear modules/npm/pnpm/Corepack/shared Docker caches, images/builders or anonymous/unrelated volumes. No codexeval-http or codex-research-http-probe fixture remains after final checks. The unproven prior /tmp/codex-universal-archive-list.txt remains untouched.
Full source builds/native suites were not rerun because runtime/core/adapter inputs remain byte-identical. Historical broad race/frontend/previous integration failures retain their previous execution scope and unresolved attribution. No iStoreOS/native amd64 loaded lifecycle, arbitrary-fork guarantee, real billing/concurrency, production deployment, deployed browser or published online bootstrap acceptance is established.
Process ownership discovery requires readable /proc, readlink/awk and an identifiable resolved executable within the bounded PID tree. Wrappers deeper than3 or >32 descendants and unavailable socket ownership can fall back only to safe private-namespace loopback hints; shared namespaces/specific non-loopback binds fail explicitly. Link-local IPv6 zones and arbitrary DNS bind names are not inferred. Unknown nc implementations that do not finish on stdin close need the bounded local-client teardown; actual no-orphan evidence here is BusyBox1.37 and the independent relay, not every nc variant. Missing nc preserves its actual error without a guessed wget fallback or new dependency. These are explicit delivery limits, not binary/version gates.
