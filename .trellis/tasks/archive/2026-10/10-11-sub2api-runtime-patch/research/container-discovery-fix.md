# Research: universal container discovery and runtime compatibility

- Query: Explain the reported Docker Hub/latest installation failures, and assess a universal single-Python package that discovers a selected Sub2API environment without embedding its version, source commit, image digest, or repository name while preserving native targeting, billing, concurrency, activation, and cleanup.
- Scope: internal
- Date: 2026-10-13

## Findings

### Reported failures and source-confirmed causes

The supplied reproduction is `docker ps` showing container `sub2api` (`9dfd087c7f3a`) with image `weishaw/sub2api:latest`. The first install reports `No Sub2API container is available.` The next attempt with `--container sub2api` reaches selection and reports `Selected container platform evidence unavailable` before staging or replacing a binary.

Both failures follow directly from the current installer:

- `install.py:174-195` lists candidates only when `Config.Image` contains the literal `wei-shaw/sub2api`. It does not recognize the supplied `weishaw/sub2api:latest` string, so the container is discarded before interactive selection.
- `install.py:216-240` uses fixed v0.2.15 OCI facts. When Docker omits `ImageManifestDescriptor`, line 231 accepts the selected image only if its `.Image` equals the one embedded v0.2.15 config digest. The reported latest image has a different config identity, so it fails with the observed error. This is an intentional fail-closed check under the old pin, not evidence that the selected container is not Sub2API.

The runtime pin is broader than those two installer checks. `install.py:19-32` embeds release/version, upstream commit, image index/platform/config digests, paths, UID/GID, and patch-set SHA. `install.py:233-240` then checks repo digests, source and version labels. The installer package name and binary lookup also use `expected_version` (`install.py:337`). Removing only the initial substring filter exposes the second rejection; removing only the second rejection would allow the installer to replace the newer whole server with the v0.2.15 binary.

### Generic discovery contract

Discovery and compatibility need separate results. Discovery may identify a likely application container from runtime-observed metadata and read-only capability probes, then list it for the required user selection. It must not use a repository-name substring as the identity test, and it must not treat a moving tag such as `latest` as an immutable version. After selection, the installer must re-inspect the same container ID and complete compatibility/build preflight before any protected write, backup, or binary replacement.

For a generic detector, useful read-only hints include the configured image reference, image ID, OCI title/source/version/revision labels, repo tags/digests, OS/architecture, process executable/build metadata, and application health or capability responses where available. These hints can narrow the candidate list; labels and tags are self-asserted or mutable, so they do not authenticate an image or prove that a patch applies. Generic enumeration also needs to avoid treating a database or cache container as a target merely because it shares a stack or has a familiar name. A candidate should pass an application capability/profile check, and the user still chooses the target.

The reported `weishaw/sub2api:latest` reference can be used as a regression fixture, but should not become another baked-in repository exception: the latest user direction explicitly rejects hardcoded repository names. An image reference parser may normalize familiar registry syntax for reporting and diagnostics, but a normalized name is a discovery hint only. Do not infer registry/path identity with substring matching. Do not infer a platform manifest digest from the tag or the container name.

### Docker identity and platform evidence

The current fallback assumes a particular Docker inspection shape. `ImageManifestDescriptor` is optional in Docker's inspect response. `.Image` is commonly the local image/config ID; it is not necessarily the OCI index or platform-manifest digest. `.Os` and `.Architecture` identify the local image platform, but do not identify all layer bytes or the source revision. `RepoDigests` may be absent for locally loaded/tagged images and may represent an index or a platform manifest depending on the image store and pull path. A Docker Hub tag is mutable.

Consequently, a version-flexible installer cannot replace the fixed digest checks with a claim that the image is official based only on `.Image`, labels, or `latest`. It needs a documented identity source per target: for example, retrieve the immutable OCI index/manifest for the selected source reference and platform, verify the manifest's config digest against the local image ID, and verify source provenance for the exact commit used to build. Handle both Docker API responses with a platform descriptor and older responses without it. Handle the image-save representation actually returned by the supported Docker engine; the current `image_config()` path expects `blobs/sha256/...` OCI entries (`install.py:263-289`) and does not establish that every classic Docker save archive has that layout.

If the local image has no repo digest, trustworthy source/build attestation, or recoverable source revision, the installer can still report it as a candidate. It cannot claim a verified match for a source build and should stop before mutation. An explicit user selection is not a replacement for source/build provenance when the action installs a whole executable into that container.

### Why one current server binary cannot be universal

The delivery currently replaces `/app/sub2api` with a newly built Sub2API server; it is not a separate runtime plugin. The build command compiles `./cmd/server` from a pinned source checkout and links the evaluation module into it (`scripts/build-sub2api.sh:195-213`). The seven production patches live under `patches/sub2api/v0.2.15/production/`; they change native authentication/repository logic, account scheduling/Responses forwarding, cancellation, target adaptation, server wiring/router/runtime, and frontend integration. They depend on the host's Go types, service interfaces, middleware, generated Wire graph, scheduler, and UI files.

The native runtime facts provider in patch 007 has additional literal guards: `EvalSourceRevision` must equal the v0.2.15 commit, `Build.Version` must equal `0.2.15`, and the image index/platform/config digests must equal the v0.2.15 values (`patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:1518-1541`). Replacing the installer constants alone therefore leaves the new binary unready. Removing the installer checks while continuing to ship the same old binary would silently downgrade any newer target and break the explicit current requirement.

The compatibility/build pipeline reinforces this binding: `compatibility/sub2api-v0.2.15.json` has one release, source commit, image identity, build recipe, and patch directory; `scripts/sub2api_build.py:77-152,163-170,332-341,679-699` validates those inputs and synchronizes the installer facts with that patch set; `scripts/build-sub2api.sh:10-12,34-37,195-213` prepares and builds the pinned source. The patch context is not a generic source transformer. No local rule can establish semantic compatibility with an arbitrary future upstream commit merely by applying these patches with fuzz.

### Native billing and integration alternatives

The current native path is the reason a static host build was selected. The runtime creates an opaque process-context capability containing the exact group/account target, then calls the real router (`patch 007:728-745,766-770,956-966,1044-1065`). The target adapter uses native admin/account/group repositories (`patch 004:22-39,57-75`). The scheduler patch forces the requested account and rejects fallback; the normal gateway pipeline then performs authorization, account eligibility, slot admission, forwarding, and usage recording. This matches PRD R3/R5 and A2/A3.

The earlier source-port audit found the existing Sub2API transport plugin protocol sees an `account_id` already selected by the host; it does not expose request-scoped account selection (`research/sub2api-port.md:15-31`). A transport plugin can run a request through the selected account's native forwarder, but cannot force the exact account requested by evaluation. Calling the ordinary public API can preserve native billing and normal shared concurrency, but its scheduler chooses the account; it is not equivalent to exact group/account targeting. Calling provider/upstream APIs directly loses Sub2API usage rows, charges, and shared native concurrency and remains outside the accepted requirement.

The local implementation does not show a dynamic in-process extension ABI for administrator routes plus pre-selection target dispatch. A first-class upstream extension ABI is the clean long-term route to a genuinely version-independent extension: Sub2API would need to expose a stable, versioned hook for admin activation/lifecycle, native target validation, request-scoped exact-account selection, normal auth/quota/billing/slots, cancellation, and history. Existing transport-plugin capability is insufficient. A sidecar or HTTP proxy cannot manufacture the current context-only capability, and an external header/loopback marker would violate the security contract.

Go `plugin` or `LD_PRELOAD` is not a drop-in way to alter this server. The current server is compiled as one Go executable (build uses `CGO_ENABLED=0`); a Go plugin requires the host to be built to load it and compatible Go package/build settings, and it cannot add hooks to a process that never calls `plugin.Open`. Native scheduler/service methods are not exported as a stable C ABI. These mechanisms do not avoid the need for an upstream extension point or a matching server build.

### Feasible architecture choices

1. **Universal installer with matching-source build at install time.** Keep one Python installer and the user selection flow. Discover the selected runtime's source revision and platform from verified metadata/provenance; fetch that exact source; run a pinned builder; apply a patch implementation known to fit that source; compile the replacement from the same upstream commit; stamp it with the discovered version, commit, patch-set identity, and observed image facts; and bind those values to protected per-install facts/journal. This avoids shipping one old server binary for every target and can preserve native accounting when the build is genuinely based on the running version. It needs network access, trustworthy source/provenance, a Go/Node/pnpm toolchain or a builder image, disk/RAM/CPU, per-target dependency locks, and a patch preflight. The existing v0.2.15 patches do not currently satisfy the generic-patch step. Unknown upstream changes must fail before mutation until ported and fully checked; “latest” is not automatically supported merely because it is latest.

2. **First-class upstream extension contract.** Add a stable host API inside Sub2API that lets a separately delivered extension register admin lifecycle/UI and dispatch an evaluation through native billing with a validated exact account target. This is the only route to one extension binary/package that does not replace the whole server, assuming upstream adopts and maintains the ABI. It requires upstream changes and server releases that implement the contract; current transport plugins cannot provide it.

3. **Ordinary public API sidecar.** This is feasible without replacing Sub2API if group-level rather than exact-account targeting is acceptable. Normal API-key requests keep native billing and scheduling, but the selected account is not guaranteed. It therefore changes product behavior and cannot be substituted silently.

### Source-driven adapter and AST-injection feasibility

The matching-source option can be made more general than a version matrix, but it still needs a semantic adapter for each host shape. A practical build engine could inspect the selected source with `go/packages`/`go/ast`, identify capability patterns (the router constructor, the native OpenAI Responses handler, account scheduler/service, usage/billing service, Wire provider set, build-info provider, and admin route registry), and generate a small adapter package plus narrowly scoped edits. The generated adapter would carry the evaluation core, activation/UI handlers, private capability type, exact target resolver, and maintenance endpoint. It would bind to the host's discovered service interfaces instead of comparing a release string.

That approach is a build-time capability check, not a runtime trick. The engine must require exactly one unambiguous match for every required hook and prove the resulting program still has the native account, group, quota, billing, slot, cancellation, and router paths. It should emit a machine-readable adapter plan containing source tree hash, matched symbols, generated edits, patch identity, toolchain/lock hashes, and invariant results. Any missing, duplicate, or semantically changed hook must stop before the container is modified. A textual regex or `go/ast` insertion around a similarly named function is insufficient: it can compile while bypassing a billing stage, re-entering failover, using stale account pointers, or exposing an external capability.

The minimum native surface is larger than one injected handler. The current implementation changes at least six backend concerns (Key/auth/repository identity, Responses scheduling and no-fallback semantics, cancellation/stream accounting, exact target adaptation, runtime/router/Wire/maintenance, and UI), touching about 3,822 production patch lines before the UI patch. A generic generator may reduce the *host-specific diff* to a generated adapter and a few hook edits, but it cannot make these invariants disappear. The scheduler and usage call sites are especially semantic: the directed target must remain fixed through every retry/wait branch, the post-slot account must be the one forwarded and billed, and directed WebSocket cancellation must not return a partially used connection to the shared pool. These rules require executable tests against the selected source, not only AST shape checks.

An AST engine also needs a stable way to identify host data types across releases. If the upstream project changes package layout, generic names, receiver signatures, generated Wire ordering, Ent schema fields, or middleware flow, the adapter needs a new strategy. `go:linkname`, Go `plugin`, binary detours, or `LD_PRELOAD` do not solve this: the current server is a statically built `CGO_ENABLED=0` executable and does not load external Go plugins, and unexported symbols/linker details are not a supported ABI. A source-driven engine is therefore implementable only as a guarded source build with explicit semantic capability tests and an expanding adapter library. It cannot truthfully promise every future `latest` image on the first release.

The least disruptive path that preserves current native guarantees is a universal **installer** plus a matching-source build with a fail-closed adapter engine. The user experience changes: the first install may need to download the exact source, start or use a builder image, install/verify Go and frontend dependencies, compile one binary for the selected target architecture, and wait minutes; an image with no trustworthy source provenance or an unsupported host shape is reported as unsupported without mutation. A prebuilt eight-file archive remains useful only for targets whose exact source/image/patch identity matches that archive. A truly universal *binary/plugin* package requires the upstream extension contract described above. Main should revise the PRD/design/implementation plan to select this source-build package or explicitly accept the weaker ordinary-API sidecar behavior before implementation; the old v0.2.15-only product restriction is superseded.

### Recommended universal package contract

The implementable universal deliverable is one Linux Python package whose installer and build engine contain no target release/version/commit/image-digest/repository allowlist. It may pin its own builder toolchain and adapter-engine version, but all host identity comes from the selected container and the build it creates. The release payload can keep `install.py` and aggregated `LICENSE/` texts, plus the generic evaluation core/adapter templates; host-specific server binaries are build outputs, not fixed package members. An optional cached binary is accepted only when its generated profile exactly matches the selected target and is never used as a downgrade fallback.

The process is:

1. Enumerate containers without an image-name substring filter and show safe, sanitized metadata. Keep explicit user selection even when one candidate exists.
2. Build a capability profile for the selected immutable container ID: executable path and hash, OS/architecture/UID/GID, writable data path, process/build metadata, OCI source/revision/version labels, immutable image/config/manifest evidence when available, and native API/capability probes. Cross-check independent facts; a label or tag alone is only a hint.
3. Resolve exact source provenance from OCI attestations/labels or an equivalent verifiable source record. If source URL/commit or the host's required capability profile is unavailable, stop with a concrete unsupported reason before backup or write. The installer never replaces a target with an older cached binary.
4. Run the source-driven structural adapter. Match required native contracts (router construction, admin route registration, OpenAI Responses admission, scheduler target selection, usage/billing, cancellation, Wire/build info, and persistent schema access), generate the host adapter and dynamic build facts, and reject missing/ambiguous/semantically changed anchors.
5. Build in an isolated builder, run compile and focused native invariants, then bind the actual source commit, source tree hash, adapter hash, toolchain locks, generated runtime stamp, image/platform/config facts, binary path, data path, and UID/GID to protected installation state. Stage and verify the candidate/backup/loaded process exactly as today.

This gives automatic environment identification without trusting a repository name and without silently downgrading. It still reports unsupported hosts when provenance or semantic capabilities are absent; “universal” means one package and one adaptation protocol, not a guarantee that every arbitrary future fork can be patched without a new adapter rule.

### Minimum executable validation before product implementation

The next planning gate can be demonstrated without touching a user's container:

- Run a read-only fake-Docker/profile fixture containing `weishaw/sub2api:latest`, `wei-shaw/sub2api:0.2.15`, a same-name malicious image, Redis, and Postgres. Assert all containers are enumerated, the user selects by immutable ID, and only the selected profile proceeds to capability checks; assert no name substring grants compatibility.
- Restore the pinned source or use an isolated copy, run a prototype `go/packages`/`go/ast` matcher, and require exactly one match for each native contract. Generate a no-op adapter and compile the unmodified host plus generated package. Record matched symbols and source-tree hash; deliberately rename/remove one anchor and assert a specific pre-mutation failure.
- Generate the dynamic facts/runtime stamp from the source commit, adapter/core hash, toolchain lock hash, and target platform. Verify that the runtime reads protected facts rather than a v0.2.15 literal, and that a changed source/image/binary identity fails reapplication.
- Build one disposable patched server from a known source profile and run focused tests for exact account pinning through retries/waits, native usage/billing and shared slots, cancellation/connection ownership, admin activation, and rollback. The existing lifecycle tests remain the oracle for backup, loaded-process hash, cleanup ordering, and history.
- Only after the profile/adapter/build checks pass should the package and install transaction be exercised on disposable amd64/arm64 targets. A host with no source provenance or failed semantic matcher must produce no mutation and no archive/binary fallback.

This experiment proves the architecture and its fail-closed behavior; it does not claim support for arbitrary `latest` until a real source profile passes the same checks.

### Persisted schemas and identity after removing compile-time pins

The current source contains several independent schemas that should not be conflated with Sub2API's release version:

- Installer journal schema is 2 and the reader accepts 1 or 2 (`install.py:297-311`).
- Installer facts are schema 2 (`install.py:505-507`); patch 007's runtime facts reader accepts installer-facts schema 1 or 2 (`patch 007:1490-1503,1522-1541`).
- Native maintenance request/result protocol is schema 1 (`patch 007:1141-1173`).
- The evaluation history file has `stateSchemaVersion = 1` (`internal/eval/state.go:14,160,196`).
- The native ownership intent introduced by patch 001 uses intent version 2 and has an explicit migration from intent version 1 (`patch 001:605-607,948-957`).

These tool/runtime protocols can remain versioned independently. They do not prove that a newer Sub2API database schema or service API is compatible. A matching-source build must run the selected upstream source's regular database migrations and must revalidate the Key fields, Ent/repository calls, gateway usage foreign keys, and cache/slot contracts. Protected facts should record the actual target/build identity created for that installation and compare it on restart/reapplication; compile-time literals for a single target version/image can be removed, but exact per-install identity checks and binary/process hashes should remain.

### Optional Docker diagnostic boundary

The following read-only commands expose only the container identity/reference/state, local image identity/platform, repo tags/digests, three selected OCI labels, and Docker client/server versions. They avoid environment variables, full labels, mounted paths, command contents, configuration, and tokens:

```sh
docker inspect --format '{{.Id}}|{{.Name}}|{{.Config.Image}}|{{.Image}}|{{.State.Status}}' sub2api
image_id=$(docker inspect --format '{{.Image}}' sub2api)
docker image inspect --format '{{.Id}}|{{.Os}}|{{.Architecture}}|{{json .RepoTags}}|{{json .RepoDigests}}|{{index .Config.Labels "org.opencontainers.image.version"}}|{{index .Config.Labels "org.opencontainers.image.revision"}}|{{index .Config.Labels "org.opencontainers.image.source"}}' "$image_id"
docker version --format 'client={{.Client.Version}} server={{.Server.Version}} server_api={{.Server.APIVersion}}'
```

These facts can explain the missing descriptor path and are useful for a local regression fixture, but they are not a prerequisite for the universal design and do not authorize removing compatibility checks or installing a different server build. The installer must obtain equivalent facts automatically from the selected container or fail closed.

### Packaging and cleanup handoff

An installer-only discovery correction can reuse the retained two-architecture archive if the exact source/image/patch identities still match its recorded provenance. The ignored archive is present at `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz`, size 73,379,445 bytes, SHA-256 `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8`. It contains exactly the eight-file payload (`install.py`, two `v0.2.15` binaries, and five `LICENSE/` texts). The older implementation report records a different historical archive identity (`fde738f...`) under a superseded directory; do not substitute that hash for the retained archive. The retained binary names and embedded installer facts are still v0.2.15, so this archive cannot solve arbitrary `latest` compatibility; it is only usable for a verified matching target or for installer-only test/repackaging work.

The universal source-build design changes native runtime/patch code and cannot be delivered by re-tarring the old v0.2.15 binaries. It requires a new runtime build from the actual selected source revision, or a future upstream ABI implementation. Put any new build output in a fresh, separate `dist/` task directory and record source commit, patch/core identities, builder/toolchain, binary hash, and package hash. Preserve the retained archive identity as historical baseline, while allowing the task's approved cleanup to remove superseded raw logs/workspaces; do not imply that every prior artifact must remain. No live server command, deployment, paid request, commit, publication, or filesystem cleanup was run in this research.

## Files found

- `install.py` — current single-process installer, fixed release facts, candidate filter, selected-image checks, and transaction journal.
- `compatibility/sub2api-v0.2.15.json` — single pinned upstream commit, GHCR image/index/platform/config digests, build recipe, and patch directory.
- `scripts/sub2api_build.py` — validates pinned source, patch inputs, build facts, and installer literal synchronization.
- `scripts/build-sub2api.sh` — prepares and builds a full replacement `cmd/server` binary from the selected manifest's pinned source.
- `patches/sub2api/v0.2.15/production/001-codex-eval-auth.patch` — native Key/auth/repository changes and versioned ownership intent.
- `patches/sub2api/v0.2.15/production/002-codex-eval-scheduling.patch` — directed account selection in native gateway scheduling.
- `patches/sub2api/v0.2.15/production/004-codex-eval-target-adapter.patch` — exact native group/account resolution.
- `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch` — static server wiring, private router capability, maintenance, and hardcoded release/runtime checks.
- `internal/eval/state.go` — evaluation history schema version and persistence reader.
- `tests/installer/test_transactions.py` — current container selection and installer identity fixtures, presently using only the GHCR v0.2.15 reference.
- `tests/installer/docker_lifecycle.py` — disposable lifecycle fixture pinned to the current GHCR image index.
- `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz` — retained ignored eight-file v0.2.15 baseline archive; it is not a universal target binary.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/sub2api-port.md` — prior source audit of native transport-plugin capabilities and exact-account routing limitations.
- `.trellis/spec/backend/runtime-installation.md` and `minimal-linux-delivery.md` — current pinned install and delivery contracts, now superseded by the user's universal-package request and requiring planning updates.
- `.trellis/workflow.md` — current Trellis planning/execution gates; complex scope requires revised PRD/design/implementation plan before implementation.

## Code patterns

- `install.py:174-210` — current candidate enumeration, target prompt, explicit target, and immutable container ID capture.
- `install.py:216-240,263-289` — current fixed image digest/source/version checks, absent descriptor fallback, and OCI blob verification.
- `install.py:452-490` — installer validates the target and extracts the original binary before writing the journal or replacing it; preserve this no-downgrade/no-mutation boundary.
- `scripts/sub2api_build.py:163-170,332-341,679-699` — verifies exact upstream source and expected facts against the one manifest/patch set.
- `scripts/build-sub2api.sh:10-12,34-37,195-213` — pins the source path and compiles the replacement server with version, commit, and patch-set build stamps.
- `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:728-745,766-770,956-966,1044-1065` — binds a private context capability to an exact group/account target and executes the real router.
- `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:1518-1541` — fixed version, commit, image index/platform/config gates in native readiness facts.
- `patches/sub2api/v0.2.15/production/004-codex-eval-target-adapter.patch:22-39,57-75` — target resolution uses native repositories and group membership.
- `tests/installer/test_transactions.py:25-69` — selection tests currently seed only `ghcr.io/wei-shaw/sub2api:0.2.15`, so neither reported failure is covered.
- `internal/eval/state.go:14,160,196` — independent evaluation history schema version 1.

## External references

- No new external documentation was needed. The Docker reference and reported error strings came from the user; the runtime/build findings are based on local installer, manifest, patches, and prior source-audit artifacts.

## Related specs

- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/minimal-linux-delivery.md`
- `.trellis/spec/backend/patch-native-installation.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/sub2api-port.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/patch-native-accounting-identity.md`

## Caveats / Not Found

- The user-reported commands and errors were not replayed against their Docker daemon. The inspected local code contains both exact failure conditions and the provided reproduction follows them.
- The local pinned upstream source checkout disappeared during the concurrent cleanup. The retained archive remains under an ignored `dist/` path as recorded above; its payload and v0.2.15 identities are available for inspection, but current latest-image labels, repo digests, registry manifest, Docker engine version, source attestation, and exact config digest are not available in this workspace.
- No arbitrary upstream version was compiled or tested. The current patches only establish a v0.2.15 integration. A capability-probed universal installer does not make those native code patches compatible with arbitrary versions.
- No product files, specs, tests, binaries, containers, server state, or release assets were changed. No Docker command, build, deployment, paid request, publication, commit, or filesystem cleanup was run.
