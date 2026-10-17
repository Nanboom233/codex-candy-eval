# Research: strict replay source and workspace preparation

- Query: Identify supported commands and inputs for a fresh independent
  production/test replay after `pristine/` was found not to be a valid Git
  checkout, while preserving the authoritative cleanup revision artifacts.
- Scope: internal
- Date: 2026-10-13

## Findings

The pinned local source is a real, clean Git checkout at
`dist/sub2api-v0.2.15`. It resolves to commit
`f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`, annotated tag object
`86a80c13dcba86f52f9ca815b5a471cecc236227`, and tag peeled commit
`f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. The supported command
`python3 scripts/sub2api_build.py verify-source compatibility/sub2api-v0.2.15.json dist/sub2api-v0.2.15`
returns `source ok`.

`dist/sub2api-build/installation-revision-20261013/pristine/` is only an
extracted source tree followed by `git init`; its `HEAD` has no commit and it
cannot satisfy `verify-source`. `scripts/sub2api_build.py:162-183` requires
`rev-parse HEAD`, the annotated tag object, the peeled commit, and an empty
`git status`. The production helper's archive extractor initializes an empty
Git repository after extraction (`scripts/build-sub2api.sh:69-95`), so that
repository is suitable as a build workspace but never as `verify-source`
input.

Use a new directory under the revision directory for every replay. The
supported production sequence is:

```sh
REV=dist/sub2api-build/installation-revision-20261013/replay-check-$(date -u +%Y%m%dT%H%M%SZ)
mkdir "$REV"
SUB2API_WORKSPACE="$REV/production" SUB2API_OUT_DIR="$REV/out" scripts/build-sub2api.sh prepare
SUB2API_WORKSPACE="$REV/production" SUB2API_OUT_DIR="$REV/out" scripts/build-sub2api.sh build --arch all
```

`scripts/build-sub2api.sh:9-14,59-96` requires the workspace and output
paths to remain under `dist/`; `prepare --force` preserves a mismatched
workspace, but a timestamped fresh path avoids reuse entirely. The build
checks the pinned source, patch inputs, lockfile baseline, staged core, Wire,
frontend, and both static binaries. It keeps the local pinned source fixed at
`dist/sub2api-v0.2.15`.

Create the independent test replay from that production workspace with the
build helper's supported CLI. The target must not exist before this command:

```sh
python3 scripts/sub2api_build.py create-test-workspace \
  compatibility/sub2api-v0.2.15.json . dist/sub2api-v0.2.15 \
  "$REV/production" "$REV/test"
```

`create-test-workspace` (`scripts/sub2api_build.py:591-620`) re-verifies the
real source checkout, archives exactly the pinned commit, applies the
production patch set, stages the allowlisted evaluation core, copies generated
Wire/frontend bytes from production, applies only the add-only test overlay,
and writes `$REV/test-evidence.json`. It enforces unchanged production bytes
before and after the overlay and checks the lockfile baseline. Run focused
native/core checks against `$REV/test`; do not point them at `pristine/` or at
an older `test-final`/`test-release` tree.

The current authoritative cleanup revision remains
`installation-revision-20261013/production-cleanup`,
`installation-revision-20261013/test-cleanup`, and
`installation-revision-20261013/out-cleanup`. Those paths are the final
delivery references, not inputs for a claimed fresh replay. A new replay must
write to a new task-owned path and record its own logs, evidence, source
identity, patch set, lockfile hashes, and output checksums.

If the local checkout is unavailable, obtain a fresh Git checkout of
`https://github.com/Wei-Shaw/sub2api`, fetch annotated tag `v0.2.15`, detach
at `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`, and verify the tag object and
clean status before using it as the `create-test-workspace` source. The
repository's `build-sub2api.sh` currently fixes production source at
`dist/sub2api-v0.2.15`; a custom source path therefore requires invoking the
helper functions directly or placing the fresh checkout at that expected
non-symlink path, without overwriting any existing task workspace.

## Files found

- `scripts/sub2api_build.py` — source identity checks, workspace safety,
  test-workspace construction, patch/core staging, and CLI dispatch.
- `scripts/build-sub2api.sh` — supported `check`, `prepare`, and `build`
  commands; fixed local source path and `SUB2API_WORKSPACE`/
  `SUB2API_OUT_DIR` overrides.
- `scripts/test-sub2api.sh` — integration runner; its default test target is
  timestamped under `dist/sub2api-build`, but direct `create-test-workspace`
  gives the checker an explicit task-owned target.
- `compatibility/sub2api-v0.2.15.json` — pinned commit/tag, patch directory,
  core allowlist, lockfiles, and toolchain contract.
- `dist/sub2api-v0.2.15` — verified local pinned Git checkout.
- `dist/sub2api-build/installation-revision-20261013/pristine` — invalid
  empty-HEAD extraction; not a source-verification input.
- `dist/sub2api-build/installation-revision-20261013/{production-cleanup,test-cleanup,out-cleanup}` — authoritative final cleanup artifacts.

## Code patterns

- `scripts/sub2api_build.py:162-183` — `verify_source` requires exact HEAD,
  annotated tag object, peeled commit, and clean status.
- `scripts/sub2api_build.py:444-472` — `prepare_workspace` rejects unmarked
  contents, refuses symlink workspaces, and preserves mismatched workspaces
  when forced.
- `scripts/sub2api_build.py:591-620` — test replay archives the verified
  pinned commit and checks production/test parity around the overlay.
- `scripts/build-sub2api.sh:9-14,59-96` — path overrides and safe prepare;
  `scripts/build-sub2api.sh:98-234` — deterministic patching, core staging,
  locked frontend, static builds, provenance, and checksums.
- `scripts/test-sub2api.sh:58-95` — integration test contract and why its
  generated default target should be replaced by an explicit revision path for
  the independent checker replay.

## External references

- Pinned upstream repository from `compatibility/sub2api-v0.2.15.json`:
  `https://github.com/Wei-Shaw/sub2api`.
- No external documentation was needed; the repository's manifest and helper
  implementation are authoritative.

## Related specs

- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/deployment-evidence.md`
- `.trellis/spec/backend/host-coordinator.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/installation-simplification-implementation.md`

## Caveats / Not Found

- No build, test gate, deployment, or product source was run or changed in
  this research pass. Only source identity inspection and read-only helper
  inspection were performed.
- The helper has no `--source` option for `build-sub2api.sh`; the fallback
  fresh-checkout path requires either the expected `dist/sub2api-v0.2.15`
  location or direct use of the helper APIs.
