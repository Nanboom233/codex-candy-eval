# Selected source and live progress implementation

Status: implementation and final archived checks complete; original-agent independent review pending. No production deployment, paid request, commit, push, publication or evidence cleanup.

## Change

`install.py` resolves ordered hints with Git ls-remote and explicit depth-one `fetch --progress --depth=1 --no-tags`. Full SHA, advertised short commit and tag/default refs freeze one full commit. A short commit absent from advertised commits can be validated against a retrieved candidate tag's peeled commit; prefix mismatch or ambiguity fails with raw Git diagnostics. Explicit unresolved hints do not trigger implicit HEAD, full history, a fixed user version or an older checkout. Default HEAD is used only when original discovery supplied no hints. The 240-second source deadline is unchanged: unnecessary retrieval is removed rather than hidden with a longer timeout.

A helper-owned `mktemp -d /tmp/codexeval-source.XXXXXX` contains the selected tree and minimal shallow Git metadata. Source facts and tree are copied into the host `TemporaryDirectory` build workspace before strict bounded helper removal. Both builders consume this same prepared tree; no build-time clone or fetch remains. `adapter/build.py` is unchanged and needs only `git rev-parse HEAD`, which works on the copied depth-one checkout. Standard builder owns/chowns its copied /tmp/source, plugin and profile; no global Git safe.directory change. Actual input testing exposed Docker cp preserving the host UID and thereby Git's ownership rejection; this narrow chown fixes that real boundary.

`adapter/Dockerfile` copies adapter/core/LICENSE into /tmp/plugin, prepared source once into /tmp/source, profile into /tmp/profile.json and output under /tmp/output. Standard builder excludes source/profile from plugin-local copying and copies source separately once. It retains Git for build metadata/dependency tools, never for a second source retrieval. Host scratch uses tempfile/TMPDIR; application persistent backups/staging paths stay unchanged.

A shared live subprocess collector streams original stdout/stderr while preserving complete raw bytes, separate exception attributes and chronological logs. Pull, Git resolution, Buildx bootstrap/plain build and standard build use it. Existing numbered steps remain; only silent commands emit actual elapsed-time heartbeat. Normal failure does not replay already streamed bytes. Timeout/KeyboardInterrupt terminates the owned process group and drains pipes within a bounded cleanup budget. Detached pipe holders cannot cause an infinite wait; log-write failure disables further failed log writes while preserving primary error and teardown tracebacks. Original standard exceptions/argv/chains remain. Release download uses available-byte read1 where supported and reports measured bytes, with percent only for valid known Content-Length.

Changed owners: install.py, adapter/Dockerfile, docs/sub2api-operations.md, scripts/test-sub2api.sh routing, maintained tests/installer/test_source_progress.py, test_preflight.py, test_universal.py and test_docker_boundary.py. Core, native adapter logic, runtime ports, adapter/build.py and payload allowlist are unchanged. Functional package hash changes because Dockerfile is a functional build input; release tag remains unpublished v0.4.0.

## Real proof before product edits

System tempfile workroot: `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/codexeval-source-progress-zcbrhddf`. `selected-source.sh` and `raw-git-proof.json` retain original commands/streams. A real four-commit local bare Git remote proved fullSHA/tag/advertised short/HEAD cases: each rev-list count 1, .git/shallow present, exact expected tree ID, elapsed 0.101–0.122 seconds. Unadvertised unresolved short ID failed exit22 with original fatal in 0.048 seconds, no history or HEAD fallback. Later maintained tests run current product SOURCE_SCRIPT and cover short + matching tag with peeled advertisement withheld, plus mismatch failure.

## Actual builder input proof

`builder-proof.py`, `actual-builder-verified.log` and `actual-builder-input-proof.json`: exit0. Actual cached source helper retrieved depth1 and docker cp'd the selected local fixture tree through the product resolver. The fixture remote was mounted read-only only into that development source helper; neither builder can see it. Real native Linux/arm64 standard Docker builder and product Dockerfile Buildx (desktop-linux docker driver) consumed the same commit/tree, retained depth1 and all four /tmp paths, and produced identical fixture output. No upstream refetch could succeed. Product Dockerfile used cached golang:1.27.2-alpine/node:24-alpine bases; apk installed Python/Git build tools. Shared builder was used without changing its settings, not deleted or pruned.

The development-only adapter/build.py in this input proof validates paths/commit/tree/depth and writes fixture output. It is not the production adapter and is not a native server compilation or installation pass. Product adapter/build.py is unchanged. No paid/runtime build was substituted by this fixture. Current root gates separately cover maintained adapter/build contracts.

Initial proof failures retained: Docker socket was unavailable; one `open -a Docker` restored the existing local daemon. Cached alpine/git lacks git-daemon (server fixture exited; missing IP); host.docker.internal fixture port was refused (network cause not inferred). The final source-helper local read-only remote avoids those fixture assumptions. First standard build then failed Git dubious ownership after Docker cp; current owned-directory chown rerun passes. No failure is erased or counted as acceptance.

## Checks and errors

Commands use PYTHONDONTWRITEBYTECODE=1 / python3 -B. Full output logs live in workroot.

- Required `scripts/test-sub2api.sh check` passed exit0 (`root-check.log`) after initial change; `root-check-final.log` also exit0 during refinements. Final `scripts/test-sub2api.sh check` exit0 (`root-check-delivery.log`) after all product refinements; includes current 8 progress cases. These cover build helpers, installer transactions/universal/preflight/boundary/HTTP/progress, facts, shell syntax, root Go vet/race and UI. Historical failed full upstream race/frontend/integration gates are not replaced by these local passes.
- Current progress tests include real depth1 tree parity and hint ordering, standard host TMPDIR, non-TTY output observed before 1.5-second process exit, complete long raw streams/log/exception attributes on deadline, exact no-replay lines, actual KeyboardInterrupt group settlement, a detached pipe holder with bounded drain, primary log-write failure vs secondary close error and download bytes/valid length. Initial assertion counted command text in traceback as replay; corrected to count exact raw output lines. macOS returned EPERM for an already exited detached-holder process group; its original traceback is retained. Cleanup was corrected to continue wait/drain after that group error, removing ResourceWarning; tests explicitly settle their exact detached fixture PID.
- Existing fake-Docker/preflight mocks now target live collector and tree cp. Early preflight failure was a mock's recursive docker/command boundary and was fixed. Maintained Buildx/standard tests assert prepared-source paths, no source build args/refetch and dedicated resource removal. Initial/final logs remain distinct.

## Final canonical archive

Canonical archive 519831 bytes, 62 regular files.
SHA256 `773d9bfe3963e21f31b041992f1438fd1d9b5e1219ebfc13cf507374b0253a95`.
Installer SHA256 `c7d11a58d772892ca3493216d90542295cec388e603bc2263b939ab62b5a8844`.
Functional build/core source SHA256 `ec8fe6ac58b320f51cbc85d388f45f996df47123252287770030875b37ec88d5`.
Historical minimal archive SHA256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8` unchanged.

Fresh package stage/extraction passes exact 62-entry source parity and production allowlist; no tests/fixtures/tools/reports/Markdown/fixed host binary. After gates passed, candidate was copied byte-for-byte to canonical `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`. Canonical entries/bytes and current source parity were rechecked read-only. Extracted archive entry tests (`INSTALLER_TEST_PATH=<workroot>/archive-final/install.py`) all exit0: source/progress8 (`archive-progress.log`), preflight18 (`archive-preflight.log`), Docker boundary14 (`archive-boundary.log`), HTTP16 with PROBE_DOCKER_TESTS=1 and no skips (`archive-http.log`). No further tests/builds/repack after these passes. `git diff --check` exit0. The original core tree SHA256 remains `dd6d6eeea1f2bde8a200b993728bf4a784952f1ec7a4e24bdc0f8dac03b9a5e9`; new functional identity is due to the build Dockerfile, not core runtime changes.

## Limits and cleanup inventory

Unknown remote cause of user's240-second timeout, remote Git support for unadvertised full commits, iStoreOS/native amd64 loaded lifecycle, native billing/concurrency, production browser and published bootstrap remain unverified. Prefix failure is a bounded source-resolution failure, not proof of native incompatibility. No arbitrary future-fork compatibility claim.

Retain workroot scripts/raw proof/bare remote/logs/candidate/source workspaces/fixture inputs/outputs/extraction until checker. Exact attributable `/tmp/codexeval-source.*` host roots from early raw Git/test executions are listed in installer-source-progress-workroot.json; no broad /tmp glob cleanup is authorized. Required Go build/test caches and actual Docker builder cache remain. Fixture helper/build containers and host Git server process are removed by finally; preserve unrelated Docker resources/shared bases. No cleanup, prune or further unchanged runtime build forms part of this handoff.
