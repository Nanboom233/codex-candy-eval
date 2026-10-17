# Installer source/progress artifact cleanup

Status: COMPLETE after accepted independent review. Cleanup ran no tests, builds, repack, deployment, paid request, commit or publication.

[Checkpoint](installer-source-progress-cleanup-checkpoint.json) preserves exact paths, per-file/group hashes, both archive identities, commands/exits and resource inventories. All 12 maintained groups match exactly before/after, including 26 test files, 24 adapter files, 19 specs, 138 existing task records, installer/core/scripts/docs/licenses/config/workspace. Only this new report/checkpoint are excluded from their own comparison.

| Exact generated root group | Before KiB | Result |
| --- | ---: | --- |
| System-temp `codexeval-source-progress-zcbrhddf` workroot | 303412 | Moved to Trash |
| Repository `__pycache__` | 116 | Moved to Trash |
| `tests/installer/__pycache__` | 148 | Moved to Trash |
| `.trellis/scripts/common/__pycache__` | 92 | Moved to Trash |
| Eight individually listed `/tmp/codexeval-source.*` roots | 896 | Moved to Trash |

The eight source roots were reported absent during review but were present at this cleanup. Each directory's Git config/advertised refs tied it to this workroot or its disposable Git tests; all eight approved exact paths are recorded in the checkpoint. No unlisted /tmp paths were touched. `/usr/bin/trash` received exactly the 12 refreshed paths and exited0; every original path is absent. Raw logs/proofs/staging/source/extracted copies remain recoverable in Trash. No direct delete, broad glob, Trash emptying or cleanup loop occurred.

`go env GOCACHE` resolved to `/Volumes/User Data/caches/go/build` for this session. `go clean -cache -testcache` exited0 and reduced that cache from322324 KiB to8 KiB. This uses the actual configured cache, not the historical default path. Dist remains72168 KiB.

No Docker deletion/prune ran. Five unrelated containers, all image IDs, network IDs and volume names match exactly. `desktop-linux` builder/settings and shared cache remain unchanged; task cache entries cannot be proven exclusively owned without risking shared builder state, so they are retained. Separate unavailable default endpoint is unchanged. No fixture/helper container or test/build process remains.

Both canonical archives remain byte-identical:

- Universal `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`:519831 bytes,62 regular files,SHA256 `773d9bfe3963e21f31b041992f1438fd1d9b5e1219ebfc13cf507374b0253a95`.
- Historical `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz`:73379445 bytes,8 regular files,SHA256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8`.

Installer SHA256 remains `c7d11a58d772892ca3493216d90542295cec388e603bc2263b939ab62b5a8844`; functional build/core identity remains `ec8fe6ac58b320f51cbc85d388f45f996df47123252287770030875b37ec88d5`; core tree remains `dd6d6eeea1f2bde8a200b993728bf4a784952f1ec7a4e24bdc0f8dac03b9a5e9`. Source/package parity established by review remains valid because source and archive bytes are unchanged.

APFS available space before86320916 KiB, after86265136 KiB; observed net change-55780 KiB. Recoverable Trash and concurrent filesystem activity affect this observation; no permanent-reclamation claim. Hooks may regenerate harmless bytecode after this checkpoint.

Independent review and implementation reports retain original execution limits and failures. Real native billing/concurrency, iStoreOS/native amd64 loaded lifecycle, production browser, live upstream timeout cause and published-bootstrap acceptance remain unverified. Cleanup does not certify publication readiness or archive the task. Main continues final read-only checks/bookkeeping.
