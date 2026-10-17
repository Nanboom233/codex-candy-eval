# Minimal Python Implementation

## Scope and current checkpoint

Read all updated role entries, PRD, design, plan, delivery research and active
specs before product edits. Main owns plans, task state and specs. This agent
owns the single installer, packager, affected developer/release inputs, tests,
public documentation and evidence. No child implementer/checker is dispatched.

Preserve the reviewed native runtime and its prior failed full gate. Use
`dist/sub2api-build/minimal-python-20261013/` for fresh package/evidence.
Retain all old workspaces, artifacts and unrelated workspace edits.
Do not read CPA private configuration, deploy, spend live quota, commit,
publish or archive the task.

The intended release has exactly eight regular files: install.py, both Linux
servers and five license/notice texts under LICENSE/. One Python process owns
bounded HTTPS download, safe extraction and Docker transactions. Status and
cleanup use embedded facts and journals without download. Packaging validates
the tracked literal facts against developer compatibility/provenance and
copies the installer bytes unchanged.

Implementation and affected checks are complete. Main will run the independent
full-scope review. Prior reviewed native binary hashes and race results are
inputs, not evidence for the new package behavior.

## Implementation milestone

The only public installation entry is `install.py`. Old install.sh,
install.ps1 and sub2api_install.py were removed under the user's authorization.
The local deletion hook required the macOS Trash, so those files were moved
there. No compatibility wrapper or second installer process remains.

One RELEASE_FACTS literal holds release/source/image/path/UID/GID/patch facts.
The developer helper AST-parses exactly one bounded literal without importing
the installer. It rejects duplicate fields, malformed assignments, dynamic
expressions and oversized source/constants. Build/package validation compares
it with development compatibility and production patch provenance. Runtime
status never supplies the expected patch identity.

The installer owns HTTPS-only redirects, a 192 MiB compressed download limit,
a 384 MiB expanded limit and exact archive paths/file types. Extraction accepts
only optional archive root/LICENSE directory entries and eight regular files.
It validates the full entry set before writing files. No package checksum list,
sidecar, manifest or provenance is read. Offline validation checks shape and
selected ELF architecture, then computes candidate/stage/loaded/backup identity
for the existing transaction. Status and cleanup do not require a bundle.

The same native journal/backup/immutable selection/drain/revocation/restore
paths remain. Maintenance additionally compares source/version/patch against
the embedded facts for every operation, including status and cleanup. UID/GID
checks match the pinned official runtime. Process messages number only actions
that run, without a checksum-preparation step.

Owners changed: install.py, scripts/package-sub2api.py,
scripts/sub2api_build.py, scripts/build-sub2api.sh, scripts/test-sub2api.sh,
scripts/test_sub2api_build.py, tests/installer/test_transactions.py,
tests/installer/docker_lifecycle.py, tests/installer/seed_browser_fixture.py,
.github/workflows/release.yml, README.md, docs/sub2api-operations.md,
THIRD_PARTY_NOTICES.md and licenses/THIRD_PARTY_NOTICES.txt. Native patches/core
and all original license texts remain unchanged.

`repo-check.log` passed 72 helper tests, 46 installer/package tests, shell
syntax, root vet/race and UI checks. New tests cover one-process stdin download,
offline discovery, download-free maintenance, exact nested extraction,
invalid archive paths/types/duplicates/size, fixed maintenance identity,
literal/development synchronization, byte-identical script/licenses and
archive-only release/tag guards. Installer fixtures do not depend on dist
artifacts in CI.

`input-reuse.log` verifies clean pinned source, original production source tree,
patch/core identities, source locks and both exact production binary hashes.
This justifies reusing out-cleanup without rebuilding or repeating the unchanged
full native/frontend gate. Its prior result remains failed: 16 race warnings,
nine top-level failures, no compiler failures/timeouts. No pristine-baseline
equivalence or unrelated fixture repair is claimed.

The first new arm64 lifecycle passed 14 checkpoints in lifecycle-arm64.log.
Final lifecycle adds actual packaged Python auto-discovery installation and
standalone download-free status, using the final minimal bundle. That run is
complete with 15 checkpoints in lifecycle-arm64-final.log and
codexeval-fixture-29e2b5090d/evidence.json. It verifies atomic expired activation,
peer denial, expiry conflicts, original-runtime paid ingress rejection and its
read-only boundary, PostgreSQL migration/cache invalidation/drift recovery,
recreation/reapplication, live drain/exact revocation/original restoration,
uninitialized uninstall and candidate-failure recovery.

Both arm64 fixtures were removed by the harness. Read-only Docker inventories
confirm no codexeval-fixture container, network or named data volume remains.
Unrelated Docker resources were untouched. Evidence directories and protected
fixture credentials remain in ignored dist; no credential content was printed.
No browser was launched for this follow-up.

## Commands and final artifacts

Relevant commands, from the repository root:

```sh
scripts/test-sub2api.sh check
python3 scripts/package-sub2api.py dist/sub2api-build/installation-revision-20261013/out-cleanup dist/sub2api-build/minimal-python-20261013/release-final
python3 tests/installer/docker_lifecycle.py dist/sub2api-build/minimal-python-20261013/release-final --platform linux/arm64
git diff --check
```

The input reuse verification checks pinned clean source at upstream commit
f2669c8cf62555cd92389b3f55920e9e6e7c6ff2, unchanged production source paths,
locks, core and patch inputs, and both production binary identities. Final
packaging does not rebuild the unchanged native runtime.

Final extracted package:
dist/sub2api-build/minimal-python-20261013/release-final/.
Final archive:
dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz.
payload-final.log records exact eight regular files, source-identical installer
bytes, and absence of JSON, Markdown, checksum artifacts and shell/PowerShell.
The final archive was built with Python tarfile from release-final and inspected
before completion. It has no checksum sidecar.

SHA-256 identities below are developer evidence only, outside the payload:

| Input/artifact | SHA-256 |
| --- | --- |
| Production patch/runtime stamp | 12a1865ff561c4c9f1c84df1c407fefb405db4093ec118d589bc651c065d9437 |
| Production staged tree | 4d28366dbdac4a597aa2f4c9bb509bc2feefa22ace8b98eb2c790251fd550fac |
| Linux amd64 binary | fd3109997b4d971e524a13ddd87b4b4c8f1cd8aad213036b0dae14ba33be213f |
| Linux arm64 binary | dd3760228f18cf63ab487ef6d3d433344d5c1b3d8ec89255f38a7deeb1bd129a |
| Final archive | 560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8 |
| Final tracked/bundled installer | c14d8c69277550ed501477edc40dfabce49fac9c6ce966c0d7b0fdb1ee5713a3 |

release-initial, initial/adapted logs, old archives and the previous independent
final-check-review-20261013.md remain historical evidence. They do not prove the
new eight-file delivery. The initial macOS path-comparison failure was fixed by
resolving test paths; no product change was required for that failure.

## Requirement mapping and remaining limits

R18/A24: one standard-library Python process handles online stdin bootstrap and
offline local discovery. All four maintenance operations use the same entry
without downloading or requiring bundle metadata. Tests establish HTTPS-only
redirects, byte limits, complete allowlist validation before extraction and
release-fact synchronization; the final arm64 lifecycle establishes actual
packaged installation and standalone status. The exact final archive establishes
the eight-file and license-only payload boundary. Online tests use controlled
responses; published online bootstrap remains unverified because no release was
published.

R9-R11: preserve official image/startup/mount/UID/GID and immutable-container
checks, journal-before-effects, original backup, candidate/staged/loaded-process
fingerprints, recovery and exact native cleanup. Unit regressions and the final
arm64 lifecycle pass. Native amd64 strict loaded-process acceptance remains open:
the earlier Rosetta run exposed a different /proc/1/exe and did not prove full
restoration success. Both production binaries retain their reviewed hashes.

R12/R15-R17: public docs and archive-only release workflow now use install.py;
tests/development manifests/provenance/evidence stay outside the payload. Tag
guards, mandatory choice even for one candidate, no-terminal/cancel/invalid
selection/ID drift and short actual-step output pass affected regressions.
R13 and K1-K7 native behavior retain the reviewed production bytes. This follow-up
does not claim a new browser or production-concurrency proof.

The required broad upstream race gate remains failed with 16 warnings and nine
top-level failures, no compilation failure or timeout. No pristine equivalence
is claimed. Unchanged broad native/frontend gates were not rerun, and unrelated
upstream fixtures were not repaired. Real paid billing, production concurrency,
deployed browser/payload, native amd64 lifecycle and published online bootstrap
remain external acceptance gaps. Windows installation is outside scope.

No CPA .env read, production deployment, paid use, commit, push, publication or
task archival occurred. Main owns AGENTS.md/spec/task synchronization and the
remaining independent review.
