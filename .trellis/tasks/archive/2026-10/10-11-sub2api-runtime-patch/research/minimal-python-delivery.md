# Research: minimal Python delivery

- Query: Collapse the Linux release and installer to one Python entry point,
  remove ineffective package checks and redundant manifests/docs, preserve the
  meaningful container transaction checks, and define one online/offline flow.
- Scope: internal
- Date: 2026-10-13

## Findings

### Current release boundary

`scripts/package-sub2api.py:10-24,61-81` currently emits a 15-file payload:
two Bash/PowerShell entry points, `sub2api_install.py`, a runtime JSON
manifest, README/operations Markdown, four top-level notice files,
`provenance.json`, `SHA256SUMS`, and two Linux binaries. `install.sh` and
`install.ps1` contain duplicate Python bootstraps (`install.sh:10-96`), while
the real transaction is in `sub2api_install.py`. The release workflow also
publishes an archive `.sha256` sidecar (`.github/workflows/release.yml:41-74`).

The package self-checks are not an independent trust anchor: the bundled
checksum manifest and provenance are shipped beside the bytes they verify.
Remove the package `SHA256SUMS`/provenance checks, the archive sidecar
checksum, and their user-facing explanations. Keep strict archive shape
validation (regular files, no symlinks/links, no traversal, duplicate or
unknown members, bounded size) because that protects extraction behavior even
without a package digest.

### Minimal Linux release tree

The release archive should contain exactly:

```text
install.py
sub2api-v0.2.15-linux-amd64
sub2api-v0.2.15-linux-arm64
LICENSE/
  CPA-MIT.txt
  Sub2API-LGPL-3.0.txt
  ModelTrace-MIT.txt
  Renderers-MIT.txt
  THIRD_PARTY_NOTICES.txt
```

`install.py` is the renamed, single public installer. Remove `install.sh`,
`install.ps1`, and the old `sub2api_install.py` name rather than retaining a
compatibility wrapper. Omit top-level README/operations Markdown and the
runtime `sub2api-v0.2.15.json`, `provenance.json`, and `SHA256SUMS`. Keep the
repository's development documentation and Trellis files; this boundary only
changes the release payload.

The current legal files are `LICENSE`,
`licenses/Sub2API-LGPL-3.0.txt`, `licenses/ModelTrace-MIT.txt`, and
`licenses/Renderers-MIT.txt`. Move their unchanged text under `LICENSE/`,
renaming the project license to `CPA-MIT.txt`. Convert the current
`THIRD_PARTY_NOTICES.md` attribution/URL text to
`LICENSE/THIRD_PARTY_NOTICES.txt`; it records Sub2API LGPL, ModelTrace MIT,
the renderer notices, and source URLs. These are legal notices, not ordinary
product documentation, and must remain with the distributed binaries.

### Embedded compatibility facts

The tagged `install.py` should contain one bounded literal `RELEASE_FACTS`
block, copied byte-for-byte into the archive. It replaces the runtime JSON
manifest and carries only installer-required facts:

```text
release_tag = v0.2.15
expected_version = 0.2.15
upstream_commit = f2669c8cf62555cd92389b3f55920e9e6e7c6ff2
image_index_digest = sha256:da0fd94daebf890419e21e956821fd6c2008fc666f9617f2eb917f4492245da3
linux/amd64: manifest sha256:71d627c5f0d383be3f9098957ac20e0adc2c65ffbc758d7a7b1f7c0084bc5745; config sha256:6c51fac82f6348083b2dad597177185e64ebf7c734592116edcb017014c2edc6
linux/arm64: manifest sha256:8b81de4ebe8fae8773e88b4f452851c75dd206ada98c6559a9efa8322397ec21; config sha256:b7cf64a3ed962aa518565c09197ad16ca64df4a20ef1edca76b1d3c79d90c539
binary_path = /app/sub2api
data_dir = /app/data
entrypoint = /app/docker-entrypoint.sh
runtime_uid = 1000
runtime_gid = 1000
runtime_patch_sha256 = <production patch-set identity for this tagged build>
```

The source literal is the online bootstrap's authority. The build/package
check must parse only that designated block with `ast.literal_eval` (bounded
source span, exactly one block, no imports) and compare its source/image/
identity values with the development compatibility manifest and its
`runtime_patch_sha256` with development `out/provenance.json`'s
`productionPatchSetSha256`. The packager then copies the exact `install.py`
bytes; it must not generate a second installer from provenance. This keeps the
online stdin script and offline bundled script synchronized without a JSON
artifact, a second process, or trusting the runtime-observed patch stamp.

Build/test manifests, lock hashes, patch/core identities, provenance, and
checksums may remain in `dist/` and CI evidence. They are not runtime inputs
and must not be copied into the release tree.

### Checks that remain meaningful

Do not weaken the transaction while removing package self-checks. Preserve:

- pinned image index/platform manifest/config identity, source revision,
  version, architecture, startup command, writable data mount, runtime UID,
  and the independently exported OCI manifest/config hash checks;
- candidate binary architecture/ELF validation and the transaction's computed
  `patched_sha256`;
- the protected journal's deployment ID, image/runtime/platform facts,
  `original_sha256`, `patched_sha256`, phase, and runtime patch identity;
- immutable original-backup hash, staged-file hash, installed binary hash,
  and loaded `/proc/1/exe` hash comparisons;
- maintenance invocation through the same server binary with runtime UID,
  deployment ID, patched binary hash, response size/schema, source revision,
  version, and embedded runtime patch identity;
- drain and exact owned-record revocation before restoring the original,
  image-drift detection, stopped/conflicting runtime failure, retained journal,
  retained backup, and incomplete-cleanup reporting.

The installer may continue using SHA-256 internally for these identity checks;
the removed checks are only the self-authenticating package/archive manifests.

### Single-process online and offline flow

Online bootstrap:

```sh
curl -fsSL https://raw.githubusercontent.com/haowang02/cpa-plugin-codex-candy-eval/v0.2.15/install.py \
  | python3 - install --container sub2api
```

`install.py` must parse the operation and explicit container choice, check the
Docker CLI, and use `/dev/tty` (or `CONIN$` only if Windows support still
exists during migration) for interactive selection. With no `--bundle`, the
same process downloads the tagged archive over HTTPS with the standard
library, extracts its bounded allowlist into a temporary directory, and uses
the embedded `RELEASE_FACTS`; no shell wrapper, JSON manifest, checksum
sidecar, or subprocess installer is involved. A no-TTY invocation without
`--container` fails before mutation. Status/inspect/rollback/uninstall can use
the journal and embedded facts without downloading binaries; install needs the
archive only when no local bundle is supplied.

Offline/bundled operation:

```sh
python3 /path/release/install.py install \
  --container sub2api --bundle /path/release
```

The bundle path is validated against the same minimal tree and the selected
architecture's binary is used. Running `python3 install.py` from a release
directory may infer its parent as the bundle, but an explicit `--bundle` is
the unambiguous documented form. The online archive may include the same
`install.py` for offline use; the online process continues executing its
tagged source bytes and does not load a second installer.

### Packaging and test changes required

Update `scripts/package-sub2api.py` to stage the minimal tree, aggregate the
license files, validate the literal fact block against development inputs, and
write no package checksum/provenance files. Update `.github/workflows/release.yml`
to publish only the archive, remove the archive `.sha256` asset, and remove
PowerShell acceptance. Rewrite installer tests that currently import
`sub2api_install.py`, parse `install.sh`, build the old `PAYLOAD`, or require
`SHA256SUMS` (`tests/installer/test_transactions.py` and
`tests/installer/docker_lifecycle.py`). Add coverage for strict minimal-tree
allowlisting, safe archive extraction, literal/provenance synchronization,
source/package installer byte equality, online stdin execution using the
embedded patch stamp, explicit/no-TTY selection, and preserved rollback,
maintenance, journal, backup, loaded-binary, and cleanup-incomplete behavior.
Keep build/test manifests and evidence outside the package.

## Files found

- `scripts/package-sub2api.py` — current release allowlist, runtime metadata
  projection, provenance/checksum copying, and payload validation.
- `install.sh` and `install.ps1` — duplicate Python bootstrap bodies for online
  download, archive extraction, and handoff to the transaction script.
- `sub2api_install.py` — actual Docker selection, identity checks, staged
  replacement, maintenance, journal/backup recovery, and rollback transaction.
- `compatibility/sub2api-v0.2.15.json` — development-only release/image/source
  facts and build inputs to validate against the embedded literal.
- `licenses/*`, `LICENSE`, `THIRD_PARTY_NOTICES.md` — legal texts to aggregate
  under release `LICENSE/`.
- `.github/workflows/release.yml` — current archive/checksum publication and
  installer/lifecycle gates.
- `tests/installer/test_transactions.py`,
  `tests/installer/docker_lifecycle.py` — old payload/bootstrap assumptions and
  meaningful lifecycle coverage to preserve while adapting entry names.
- `dist/sub2api-build/installation-revision-20261013/release-cleanup/` — current
  superseded 15-file payload showing the files to remove or aggregate.

## Code patterns

- `sub2api_install.py:119-192` — image/platform/source/startup/mount/UID
  identity checks; preserve with embedded facts.
- `sub2api_install.py:194-214,286-347` — binary and loaded-process hashes,
  journal validation, maintenance identity, and original restore; preserve.
- `sub2api_install.py:236-264` — package allowlist, SHA256SUMS, and provenance
  checks; remove package digest/provenance dependency while retaining ELF and
  runtime identity checks.
- `sub2api_install.py:349-432,481-531` — staged install, failure recovery,
  drain/revoke-before-restore, and cleanup-incomplete behavior; preserve.
- `install.sh:39-87` — online/offline bootstrap behavior to move into the same
  Python process.
- `scripts/package-sub2api.py:36-80` — build-only allowlist/provenance checks
  to replace with minimal-tree staging and literal-fact synchronization.

## External references

- Pinned release source from the compatibility file:
  `https://github.com/Wei-Shaw/sub2api`.
- Tagged online installer/archive origin currently uses
  `https://raw.githubusercontent.com/haowang02/cpa-plugin-codex-candy-eval/v0.2.15/`
  and the matching GitHub release archive URL.

## Related specs

- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/deployment-evidence.md`
- `.trellis/spec/backend/host-coordinator.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/installation-simplification-plan-handoff.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/installation-simplification-implementation.md`

## Caveats / Not Found

- No product files, packager, release, build, test, deployment, or live
  resources were changed or run in this research pass.
- Existing lifecycle evidence and current cleanup payload are superseded
  inputs; the implementer must adapt tests and package inspection to the new
  `install.py`/`LICENSE/` tree rather than preserving obsolete payload names.
- The exact `runtime_patch_sha256` literal must be updated with the release
  patch set and checked against development provenance before packaging; it
  must never be learned from a runtime status response.
