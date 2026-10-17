# Research: No-Host-Build Delivery

- Query: Compare release-time prebuilt runtime artifacts with binary deltas for avoiding installation-time compilation; identify compatibility and rollback constraints.
- Scope: internal
- Date: 2026-10-14

## Findings

### Files found

- `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz` — retained historical 73,379,445-byte archive (about 70 MiB), SHA-256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8`; contents include `install.py`, five license/notice files, and `sub2api-v0.2.15-linux-amd64` (125,321,376 bytes) plus `sub2api-v0.2.15-linux-arm64` (116,195,488 bytes).
- `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz` — current source-package snapshot, 508 KiB, SHA-256 `773d9bfe3963e21f31b041992f1438fd1d9b5e1219ebfc13cf507374b0253a95`; contains functional installer, adapter/runtime and evaluation sources, and licenses, with no host server binary.
- `.github/workflows/release.yml:31-65` — runs source/lifecycle checks, packages the universal source payload, and publishes only that archive.
- `adapter/build.py:14-51,53-94` — structurally adapts the selected host source, copies runtime additions, embeds the frontend, type-checks and compiles the server for the selected architecture. It writes `EvalDataDir` from the per-install profile into generated Go source (`:63-65`).
- `tests/installer/universal_ci.py:17-44` — validates the adapter against the two latest actual upstream releases and records their source revisions as test evidence; it does not build or publish reusable host binaries.
- `install.py:462-503,733-756` — discovers architecture and data directory, resolves source hints to an exact Git revision, and records source/version/toolchain facts used for the per-install build.
- `adapter/build.py:57-65,84-94` — checks that the checkout matches the selected revision and emits build facts from the selected profile, including architecture, source revision and dependency-lock hashes.
- `scripts/package-sub2api.py:13-32` — enforces an exact source/license allowlist and stages no host binary.
- `install.py:917-976,1102-1159` — builds the resolved source before target writes, then journals, backs up the original, stages and atomically replaces the binary, restarts and verifies; on candidate failure it restores and verifies the original while retaining the journal.
- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md:22-35` and `prd.md:31-41,103-107` — require a universal source package, exclude fixed-version host binaries, authorize an automatic selected-source Docker build, and remove binary-identity gates.
- `.trellis/spec/backend/universal-linux-delivery.md:46-60,98-125,145-165` and `runtime-installation.md:94-119` — define source/capability discovery, host-build contract, payload contents, and recoverable replacement behavior.

### Comparison and recommendation

If the requirement is **no compilation on the installation host**, complete prebuilt patched server binaries are clearer than binary deltas. CI can compile full binaries for each supported architecture and source profile; installation then selects a matching artifact. This still requires compilation in CI and a finite supported-profile catalog. It conflicts with the active contract's universal source package, exclusion of fixed-version host servers, automatic source discovery without a version allowlist, and support for differing compatible source revisions. A reusable prebuilt also needs a solution for `EvalDataDir`, which is currently compiled into generated source per installation.

A binary delta can reduce download size when the original server is an exact known base, but it does not remove release-time compilation. Applying it safely requires identifying/verifying that exact base and refusing mismatches. No delta generation or application mechanism was found. Exact-base gating conflicts with the current removal of binary checksums, ELF checks and host/source/version/image identity gates; omitting the gate risks corrupting an unknown base. Deltas therefore add compatibility and recovery complexity without satisfying a stronger requirement than complete prebuilt binaries.

### Concrete prebuilt flow and unsupported profiles

A prebuilt design would have release CI compile one complete patched server per declared `(source_revision, architecture)` profile, publish those artifacts, and let the installer resolve the observed source hints plus architecture to exactly one artifact. `build-facts.json` already records dynamic build-profile facts during installation, but CI currently records profile revisions only as validation evidence and publishes no binaries. The installer would need a finite, explicit profile-to-artifact lookup. An unknown or unresolvable profile must fail before backup/replacement; there must be no old-binary fallback and no local compilation fallback. This lookup turns the active universal capability-based contract into a finite source-support catalog and requires an explicit contract change.

Each complete artifact can be reused across installations only after removing the compiled per-install `EvalDataDir` assumption. The runtime must derive the actual protected data path safely at execution time or the artifact still has to be built per target. No such generic runtime path mechanism was found. Embedding installation-specific data paths as an artifact selector would multiply the profile dimensions and defeat reusable release-time artifacts.

If “no compilation anywhere” includes release production, neither option fits the current architecture: the host hooks and runtime are compiled into the server. No dynamic extension/ABI path for avoiding that compilation was found. A separately adopted upstream extension interface would be a different architecture.

### Transaction boundary to preserve

Changing where the candidate binary comes from must preserve the existing installation transaction. Complete container selection and functional/source compatibility checks before target mutation. Keep the protected original backup and journal, stage the selected candidate, atomically replace only the selected container's binary, restart it, and verify startup health plus functional maintenance readiness. On failure after replacement, restore and verify the original and retain failure evidence. Candidate preparation failure must leave the target unchanged. Archive checksums may validate downloaded release payloads, but must not be repurposed as removed host-binary identity gates.

### Related specs and external references

The controlling contracts are `universal-linux-delivery.md`, `runtime-installation.md`, `patch-native-installation.md`, the active task PRD, and its approved design. No external references were needed for this repository-grounded comparison.

## Caveats / Not Found

- The historical fixed-binary archive is a local retained artifact, not the active release format or evidence of arbitrary host compatibility.
- No prebuilt catalog, binary-delta tooling, exact-base support contract, or generic runtime-data-path resolution was found.
- CI's two-source validation is not itself a supported prebuilt catalog; it neither publishes binaries nor defines installer behavior for an unlisted source revision.
- No builds, binary inspection/disassembly, release publication, or renewed native API/plugin survey were performed.
