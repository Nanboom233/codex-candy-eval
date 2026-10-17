# Clientless application probe implementation

Status: implementation complete; original-agent independent review pending. No deployment, paid request, commit, push, publication or cleanup of retained evidence occurred.

## Change and scope

`install.py` now uses uniquely named `docker run -i --pull=never --rm --name <owned> --network container:<selected immutable ID> --entrypoint nc <validated immutable helper image ID> <address> <port>` for every HTTP probe. No target nc/curl/wget is executed or installed. Helper launch has no mounts, published ports, privilege, PID namespace sharing or Docker socket. Shell/Git/nc capability acquisition is cache-first, once per Installer, before the 48-second candidate budget; the same ID feeds source resolution. Installation, startup/recovery and rollback/uninstall health share this transport through `http_get`.

The existing framed collector keeps stdin open until a bounded complete response, retains fixed/chunked/interim/trailer/EOF parsing, IPv4/IPv6 authorities and original stdout/stderr. Each exact helper name is explicitly removed with two bounded five-second attempts even after success, launch failure, timeout or KeyboardInterrupt. Removal failure cannot turn into success; a primary exception remains primary and cleanup prints its complete traceback. Stream teardown likewise preserves the primary exception. Complete subprocess argv/streams, TimeoutExpired and CalledProcessError attributes and traceback chains remain intact. HTTP wire/body bounds alone remain 8192/4096 bytes.

Owned listening endpoints are no longer capped before ordering. Valid SERVER_PORT is promoted within owned endpoints before the eight-candidate cap. Other owned endpoints remain candidates. This does not attribute the user's 6185 listener or prove that 6185 is wrong; ownership remains required for shared-network or non-loopback probes.

Files changed for this correction: `install.py`, `docs/sub2api-operations.md`, `tests/installer/test_http_transport.py`, `test_preflight.py`, `test_universal.py`, `test_docker_boundary.py`, and this report. Runtime/adapter/core and release mappings are unchanged. Main owns task/spec metadata.

## Actual raw proof before product edits

Workroot: `/tmp/codex-application-clientless-fix-20261013`. `raw-helper-proof.py` and `raw-helper-proof.json` preserve exact commands/results. Cached `node:24-alpine` derived fixture removed all nc/curl/wget aliases, ran with network none, and served both routes after 1.5 seconds. `command -v` confirmed all three absent. Cached helper ID `sha256:c0280cf9572316299b08544065d3bf35db65043d5e3963982ec50647d2746e26` produced valid Content-Length health and setup replies in 1.707/1.657 seconds, exit 0. Target Docker diff was empty before/after; target and exact helpers were removed in finally. No new image pull.

Initial fixture proof omitted Docker `-i`, returned empty EOF and failed an assertion. Corrected proof passed before product edits. That failed attempt is not acceptance.

## Verification and outcomes

All commands below used `PYTHONDONTWRITEBYTECODE=1` and `python3 -B` where applicable. Full logs remain in workroot. Standard root check skips optional Docker cases; explicit Docker run executes all.

- `scripts/test-sub2api.sh check`: exit 0 (`root-check.log`); repeated after final teardown/product refinement, exit 0 (`root-check-final.log`). Foundation, transactions, universal, preflight, boundary, transport, installer facts, shell syntax, root Go vet/race and UI checks pass. Go race second run used valid caches. This is the local root gate, not the historical failed full upstream integration gate.
- `python3 -B -m unittest discover -s tests/installer -p test_preflight.py`: 16 tests, exit 0 (`preflight.log`). `test_universal.py`: exit 0 (`universal.log`). First transport fake run hit an accidental test-only NameError and was fixed. First boundary run overapplied its cleanup-failure fixture to HTTP helpers and stopped before source stage; scoped the fixture to source/build cleanup and corrected acquisition-stage expectations. Final maintained suites pass.
- `PROBE_DOCKER_TESTS=1 python3 -B tests/installer/test_http_transport.py`: 16 tests, no skips, exit 0 (`real-http-tests-final.log`). Actual isolated clientless fixture covers fixed/chunked delayed 1.5-second replies, wildcard/specific IPv4, wrapper-child IPv6 ownership, restart health, empty disconnect, deadline, KeyboardInterrupt, failed helper namespace launch and exact removal. Actual owned 6185 plus configured 8080 resolves 8080 and logs only health/setup. Target diffs remain empty. Packaged CLI passes delayed native pair then intentionally stops at absent persistent data mount before any target mutation.
- First real failure fixture expected RemoteDisconnected but BusyBox nc retained the half-closed socket and reached its original three-second TimeoutExpired. Test now accepts those two original transport outcomes and asserts exact removal; no product deadline widening. Subsequent actual matrix passes.
- Extracted canonical archive: `INSTALLER_TEST_PATH=<workroot>/archive-extracted/install.py PROBE_DOCKER_TESTS=1 python3 -B tests/installer/test_http_transport.py`: 16 tests, exit 0 (`archive-http.log`). Same extracted installer preflight 16 tests exit 0 (`archive-preflight.log`), Docker boundary 14 tests exit 0 (`archive-boundary.log`). Maintained clientless Docker tests now create their own uniquely tagged cached-base fixture and remove that exact image after the class; final self-contained rerun also passed all 16 tests, exit 0 (`archive-http-final.log`).

Simulation/socket cases also assert cache acquisition once/before budget, helper reuse in source resolution, >8 owned-port prioritization, no target tools or privileged mounts/ports, malformed/truncated/redirect/bounds and full stderr, strict cleanup failure after success and primary preservation, cancellation and restart address retention. Tests and fixture tooling remain outside the final payload.

## Archive identity

Canonical `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`, v0.4.0 unpublished.
Size 518738 bytes, 62 regular files.
SHA256 `9f82539dd720781c5865581b0cb064d7f9f447b380e4dbd91197d5d9f463a919`.
Installer SHA256 `36d7b344c416fa91ac8975f019d40541c04d44c3fabd5d224bd714a7b5315992`.
Functional runtime SHA256 `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759` unchanged.
Historical minimal archive SHA256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8` unchanged.

`package.py` stages a fresh allowlisted source bundle, writes deterministic gzip/tar, checks exact regular entries, extracts and validates source parity. No tests, fixtures, audit tools, logs, Markdown or fixed host binary are shipped. `package-facts.json` retains identities. Canonical previous archive was replaced only by the verified current 62-file archive.

## Remaining boundaries and retained outputs

iStoreOS/native amd64 full installation, native loaded lifecycle, real billing/concurrency, production/browser and published bootstrap remain unverified. No cause is assigned to remote registry or the user's exact listener. Historical full upstream unit-race/frontend gates remain failed with no pristine A/B attribution; this narrow correction does not rerun unchanged runtime builds or repair those failures.

Retained generated outputs: workroot raw proof script/JSON, derived fixture Dockerfile/build log, temporary test editing scripts, test/check logs, fresh staged bundle, extracted archive, package script/facts; the raw-proof derived image `codexeval-clientless-fixture:20261013`; root Go build/test cache generated by required checks. Newly self-contained tests remove only their own target/helper containers and uniquely tagged image. Evidence/root-cache cleanup is deferred until independent review and must preserve both final archives, maintained files and other agents' roots/images.
