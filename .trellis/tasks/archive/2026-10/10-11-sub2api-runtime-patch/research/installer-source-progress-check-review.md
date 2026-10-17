# Installer Source and Progress Check Review

## Findings (fixed)

No local correctness or spec-compliance issue found in the reviewed source-resolution, live-output, temp-path, or builder-input changes. No product edits were needed.

## Findings (not fixed)

- File: `adapter/Dockerfile`
- Issue: The retained Buildx proof log reports two advisory `InvalidDefaultArgInFrom` warnings because `GO_IMAGE` and `NODE_IMAGE` intentionally have no defaults.
- Recommendation: Keep the values caller-supplied and validated by `install.py`; adding arbitrary image defaults would weaken source-driven tool selection. Both Buildx and standard builders passed with the explicit immutable image values.

No other code findings. The earlier user-reported 240-second retrieval delay was not reproduced against the live upstream, so its remote cause remains unknown. An unadvertised short commit may not be fetchable by every Git server; the current contract fails boundedly rather than fetching history or substituting HEAD.

## Verification

- Lint: No Python linter is configured. The actual Buildx proof completed with the two advisory warnings above.
- TypeCheck: No Python type-check command is configured. `py_compile` passed for the installer and installer tests. The repository check also passed Go vet and race tests.
- Tests: Pass. `scripts/test-sub2api.sh check` exited 0: 72 build-helper tests (1 historical skip), installer transaction/universal/preflight/boundary/HTTP/source-progress suites, shell syntax, Go vet/race and UI checks. Independent focused tests passed: source progress 8/8, preflight 18/18, universal 8/8 and Docker boundary 14/14. Extracted archive tests passed: source progress 8/8, preflight 18/18, boundary 14/14 and Docker-enabled clientless HTTP 16/16 with no skips.
- Package: The extracted 62-file package passed the repository payload validator and matches current source. `install.py` and `adapter/Dockerfile` are byte-identical to the extracted files. `git diff --check` passed.
- Spec: [universal-linux-delivery.md](../../../spec/backend/universal-linux-delivery.md) already records these source/progress contracts; no spec change was needed.

## Independent Evidence

The implementation report's local remote and both-builder proof was inspected independently. `actual-builder-input-proof.json` records the same depth-one revision `99ebdb01dcaee4e3de5a134a892734a58c26687b` and tree `8425893e917c909be9062addb9cb7ea62e7b326a` for standard Docker and Buildx. The read-only fixture remote was unavailable to both builders. The Buildx log shows `/tmp/source`, `/tmp/plugin`, `/tmp/output` and `/tmp/profile.json` inputs and a successful fixture verifier. This validates source reuse and Dockerfile input paths on native Linux/arm64; it is not a native Sub2API runtime compile or installation.

The independent progress tests used only a four-commit local Git fixture, not the upstream remote. They verified full SHA, tag, advertised short SHA, default HEAD when no hints exist, depth one, exact tree parity, bounded unresolved-prefix failure, output before child exit, complete streams/logs, timeout/cancellation, bounded detached-pipe drain and measured download bytes. The timeout/EPERM traceback emitted by the detached-pipe regression is its expected secondary teardown diagnostic; the test passed and settled its fixture process.

Canonical archive: `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`, 519831 bytes, 62 files, SHA256 `773d9bfe3963e21f31b041992f1438fd1d9b5e1219ebfc13cf507374b0253a95`.

Installer SHA256: `c7d11a58d772892ca3493216d90542295cec388e603bc2263b939ab62b5a8844`.

Functional source SHA256: `ec8fe6ac58b320f51cbc85d388f45f996df47123252287770030875b37ec88d5`. Core tree SHA256 remains `dd6d6eeea1f2bde8a200b993728bf4a784952f1ec7a4e24bdc0f8dac03b9a5e9`.

Earlier builder-fixture failures remain documented in the implementation evidence: the Docker daemon was initially unavailable, the cached helper lacked `git-daemon`, a host fixture port was refused, and Docker preserved a host UID that Git rejected. The final read-only helper fixture and local chown correction passed; these local cases do not explain the live upstream delay. My first read-only payload-validator invocation omitted `scripts` from `sys.path` and failed with `ModuleNotFoundError`; rerunning with the repository script path passed without source changes.

## Retained Cleanup Inventory

No cleanup, Docker prune, deployment, paid request or publication was performed. The implementation workroot `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/codexeval-source-progress-zcbrhddf` remains at 303412 KiB with proof scripts, builder workspaces, logs, fixtures and archive staging. The eight `/tmp/codexeval-source.*` roots listed in `installer-source-progress-workroot.json` were absent during this review; no removal was attempted. No source/builder fixture containers remain. The `desktop-linux` Buildx builder is running and retained; the separate `default` endpoint reports unavailable and was not changed or pruned. The shared Go build cache remains at 322324 KiB. Repository bytecode directories remain at `./__pycache__`, `./tests/installer/__pycache__` and `./.trellis/scripts/common/__pycache__`.
