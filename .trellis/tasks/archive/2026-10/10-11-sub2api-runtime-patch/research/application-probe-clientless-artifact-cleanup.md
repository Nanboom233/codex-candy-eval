# Clientless probe artifact cleanup

Status: COMPLETE within the authorized exact scope. Independent review preceded effects. No tests, build, repack, deployment, paid request, commit or publication ran during cleanup.

[Durable checkpoint](application-probe-clientless-cleanup-checkpoint.json) records per-file and group hashes, exact paths, archive identities, resource inventories and command exits before/after. All 12 maintained groups match exactly: installer, 25 maintained tests, 24 adapter files, core, legacy, scripts, docs/licenses, 19 specs, 132 prior task records, workspace and configuration. The new cleanup report/checkpoint are excluded from their own preservation comparison.

| Exact generated root | Before KiB | Result |
| --- | ---: | --- |
| `/tmp/codex-application-clientless-fix-20261013` | 9748 | Moved to Trash |
| Repository `__pycache__` | 112 | Moved to Trash |
| `tests/installer/__pycache__` | 172 | Moved to Trash |
| `.trellis/scripts/common/__pycache__` | 200 | Moved to Trash |

`/usr/bin/trash` received only those four absolute paths and exited 0. All original paths are absent. Evidence remains recoverable in Trash; raw proof/logs/staging/extractions are no longer at their original paths. No global Trash emptying or direct filesystem deletion occurred.

The owned Docker fixture `codexeval-clientless-fixture:20261013` was checked twice as `sha256:c395b576ca7b62e5133a10942fdbb7b5c543222530de9510f580053c8a4ac74c`, with its sole tag and no referencing container. `docker image rm --no-prune <exact-tag>` exited 0, removing that ID only. Before/after Docker image-ID inventories match after excluding that single ID. Shared Node/alpine-git images and unrelated resources remain. All five unrelated containers match exactly; no matching fixture/helper container or test/build process remains.

`go clean -cache -testcache` exited 0. Go build cache fell from 322280 KiB to 8 KiB. Retained dist remains 72168 KiB. No global prune, shared dependency cache cleanup, tests or cleanup loop ran.

Both canonical archives retain exact bytes, hashes and regular-file counts:

- Universal `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`: 518858 bytes, 62 files, SHA256 `6ee8ffe55f43c564c22e764dd96be710548fbe7a03cf471283918b7b5ae5f823`.
- Historical `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz`: 73379445 bytes, 8 files, SHA256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8`.

Latest installer SHA256 remains `1275095a76e5eb8bbcc62d34aaa1a48e8e59344b68fea1bf1f9cb45ee8d8e316`; unchanged functional runtime identity remains `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759`. The review's source/archive parity remains valid because preserved source and archive bytes are unchanged.

APFS available space before 46509156 KiB, after 46831616 KiB (observed net 322460 KiB). This is a filesystem observation affected by concurrent activity and recoverable Trash, not a permanent-reclamation claim. Hooks may recreate harmless bytecode after this checkpoint.

Current independent review remains authoritative for the latest cleanup fix and intermittent transport-test observation. Historical failed broad integration gates, iStoreOS/native amd64 lifecycle, real billing/concurrency, production browser and published-bootstrap limits remain. Cleanup does not establish release readiness or archive the task. Main performs final read-only verification and task/spec bookkeeping.
