# Robust preflight implementation

## Current checkpoint

Status: product implementation and final packaging complete, 2026-10-13. Fresh two-source focused native/type checks passed with observed exit 0. Final disposable arm64 lifecycle passed all 18 checkpoints with observed exit 0. The original implementer loaded current implement.jsonl, its contexts, PRD, design, plan, quality rules and research. Latest user instruction removes target binary attestation throughout discovery, transactions, readiness and maintenance. Historical strict-attestation statements are superseded.

Owners: install.py, adapter/build.py and runtime templates, maintained tests, docs and package. Main owns specs and task metadata. No children are spawned. Preserve unrelated work.

Confirmed research: Python buffered r+ requires seekability on /dev/tty. Separate terminal readers/writers and terminal-only stdin fallback retain piped-bootstrap selection. Source/version/OCI values become hints, not identity gates. Source checkout is resolved automatically through Git, frozen for this build, with toolchain declarations from source. No installed-binary bytes are parsed.

Remove host/backup/staged/loaded SHA and ELF gates, mandatory build-info/parseable version, source/version/image matching, readiness executable hashing and maintenance binary fields. Keep protected ownership, immutable selected container, billing mode, data mount/UID, private no-symlink storage, original backup, journal-before-effects, atomic replacement, bounded health retries, live drain and exact accounting-Key revocation. Legacy attestation fields are optional ignored migration data. Functional package/source integrity remains.

New generated workspaces and release staging live under /tmp/codex-robust-preflight-20261013 outside root Go discovery. Evidence stays until independent review. Prior raw evidence was cleaned and cannot be reused as new proof. Baseline final archive: v0.4.0, 514343 bytes, 62 files, SHA ff096a9de4a6a7d6b69dd6f76eb107e34aab31c65549e53aa35711bed0f9c8b1. Historical minimal archive 560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8 remains untouched.

Required new verification: real PTY including code through stdin, source short/full/tag/no-hint resolution and download failure, no-write regressions, root checks, fresh two-source type/native/composed checks and newly built disposable lifecycle. Required broad gates remain failed pending actual applicable runs. No unrelated upstream fixture fixes. No live deployment, paid request, commit, push, tag or publication.

## Execution evidence

The retained baseline installer reproduced the PTY failure using actual controlling PTY and piped Python code (pty-before.log, expected nonzero). The new test uses real descriptors and passes all nine preflight cases (preflight-new.log), including terminal-only fallback, EOF/cancel/invalid input and no selection without tty. Earlier StringIO cases remain supplemental tests.

Installer now reads separate tty channels, resolves safe executable paths through process or startup facts, treats native -version/OCI labels as optional hints and never reads host executable bytes. A disposable Docker Git resolver freezes the selected full checkout and reads backend/go.mod tool declarations. Absent hints discover default-project HEAD, while unresolved supplied refs produce a concrete error.

Installer target/backup/staged/process fingerprint and ELF checks are removed. Journals retain phases/path/UID ownership and ignore legacy hashes. Native readiness accepts protected schema-3 deployment ownership without executable/source/version/image attestation. Maintenance sends no binary field and checks deployment ownership. Source/package stamps are provenance only.

Safe maintenance uses new leading --codex-eval-maintenance-capabilities and --codex-eval-maintenance flags handled before native startup. Original native flag.Parse rejects those unknown flags before Wire/DB. The lifecycle adds original/stale-socket/recreation no-server/SQL-effect cases.

Maintained overlays now originate from tests/host-overlay, with current schema-3 ownership-only fixtures and explicit ignored legacy metadata regression. Exported add-only test patches are synchronized and filter generated metadata. The old fixed-host production exporter is explicitly retired so it cannot reintroduce old runtime gates. Historical patch bytes remain development evidence, outside active build/package paths.

Earlier root scripts/test-sub2api.sh check passed (repository-check.log), including 72 helpers/one historical skip, 48 transaction, eight universal, then-nine preflight, root vet/race and UI. Final focused/type checks are complete. Two attempted git checkout commands on unmaterialized fresh --no-checkout clones were rejected by PreToolUse as dirty trees. No bypass occurred. Fresh separate git clone --branch workspaces succeeded.

No product release/paid/deployment action occurred. Final archive has been replaced and verified as source-identical. Final lifecycle passes. Independent review remains pending. No cleanup of new evidence has started.

## Final verification checkpoint

Evidence root: `/tmp/codex-robust-preflight-20261013/`. These successful commands have observed exit 0.

| Check | Actual result and evidence |
| --- | --- |
| Final preflight | 15 tests pass, `preflight-final.log` |
| Extracted final installer | Same 15 preflight tests pass, `preflight-packaged.log` |
| Transactions | 48 tests pass, `transactions-final.log` |
| Universal discovery | Eight tests pass, `universal-final.log` |
| Build/package helpers | 72 tests, one historical skip, `helpers-final.log` |
| Fresh preparation | 73 rules and production/test byte parity pass, `prepare-verified-current.log` and `prepare-verified-previous.log` |
| Native type contracts | Both sources pass, `types-verified-current.log` and `types-verified-previous.log` |
| Focused native race | Both sources pass, `native-focused-verified-current.log` and `native-focused-verified-previous.log` |
| DTO/repository race | Meaningful accounting/redaction tests execute and pass, `native-dto-repository-current.log` and `native-dto-repository-previous.log` |
| Real no-app container | Alpine refused without installer writes or PID-start change, `no-app-preflight.log` |
| Actual absent-hint source | Docker Git resolver succeeds before writes, `source-default-real.log` and `.json` |
| Final archive | Real installer extraction/validation and source/license bytes pass, `package-final-verification.log` and `.json` |

HTTP preflight now uses bounded raw loopback HTTP, retaining status through parsing.
It checks typed JSON, completed setup, redirects/errors, chunked/truncated responses and 4,096/8,192-byte limits.
The integer code field cannot accept false. Health alone cannot identify the application.
No credentials, proxies, redirects or target temporary files are used.

Fresh native commands run in `verified-current/backend` and `verified-previous/backend`:

```sh
go test -tags=unit -race ./internal/service ./internal/server/middleware ./internal/handler/dto ./internal/repository ./internal/handler -run 'TestEval(Native|RuntimeDeployment|Maintenance|Composed)' -count=1
go test -tags=unit -race ./internal/handler/dto ./internal/repository -run 'TestEval(Accounting|OwnedRow)' -count=1
```

The first filtered command has no matching DTO/repository tests. The second explicitly executes their relevant tests.
An earlier fixture wrongly used create-only writes for replacement. Protected atomic replacement fixed it without weakening the metadata-drift assertion.
Fresh verified workspaces were prepared after corrected test-overlay export. Earlier failed/debug workspaces are not final acceptance.

Source profiles remain actual v0.2.15 `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` and v0.2.14 `0363b8cdba8cec3e2ba4b2dbd49c4481143fa55d`.
The independent no-hints resolver discovered current default HEAD `4d9ac3e482106d2da3a39fe3e546d8d6801f1cc1` and declared Go 1.27.2.
Affected Python compilation and diff whitespace checks pass. Shell syntax passed independently of an obsolete temporary-log read.

### Final artifact

Release `v0.4.0` remains unpublished. Functional hash is `aff326691ea26c57609827b3325922838657fbdb44cf93c5d116d3d9c5c6d759`.
Archive: `dist/sub2api-build/universal-20261013/sub2api-codex-eval-universal.tar.gz`.
Size: 514,791 bytes. Exact payload: 62 regular functional/license files.
SHA-256: `155bbd03abb7fe140a94b00350843b2a4938da3ade363874a4d072c2ad030b13`.
Fresh staging `package-ready/` and real-installer extraction `archive-extracted/` match all maintained source/license bytes.
Historical minimal archive retains SHA `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8` at its original path.
Previous universal archive remains temporarily in the evidence root for review.

### Lifecycle checkpoint

Final actual command completed with observed exit 0:

```sh
python3 tests/installer/docker_lifecycle.py /tmp/codex-robust-preflight-20261013/package-ready --platform linux/arm64
```

Log: `lifecycle-ready.log`. Fixture evidence: `codexeval-fixture-d746c1cc83/`.
First actual matching-source arm64 build/export/install passed with final functional hash.
Real Redis/PostgreSQL failed preflight before writes and retained the same accounting row count.
Initial activation, repeat install, public single-Python installer/status, ownership and expiry drift checks completed.
Original-runtime recreation and stale-socket cleanup preserved PID-start and accounting row count, without a second server.
Actual matching-source reapplication, legacy migration, counter/cache preservation and drift refusal completed.
Live drain, exact owned-Key revocation and health-verified original restoration completed.
The third real build and uninitialized uninstall passed. The fourth build passed candidate-failure original recovery with `failed_restored` journal.
The fifth actual build passed retry installation after recovery. Overall command passed all 18 checkpoints, exit 0.
All five retained build profiles match actual source, arm64 platform and final functional hash.
Supplementary capture confirms loaded maintenance provenance carries final functional hash. The candidate and safe reply remain available for review.
Final evidence: `codexeval-fixture-d746c1cc83/evidence.json`. Automatic harness cleanup removed only its containers, network and named data volume.
Raw diagnostics and all final workspaces remain. Developer binary digests in evidence are informational, without product admission checks.
Builder/resolver logs, profiles and build facts are copied before temporary workspace removal.

Earlier lifecycle runs are not passes. `lifecycle-first.log` used a debug hash and ended with explicit interrupt exit 130.
`lifecycle-final.log` lost its host process. Original exit was unobserved. Builder reports `CANCELED` and `context canceled`.
Its target had no installation journal and remained before writes. Preserved canceled evidence is under `interrupted-build/`.
Only exact orphan resources from that owned attempt were removed. No global prune, unrelated resources or cache deletion occurred.

### Remaining boundaries

Independent review and new-artifact cleanup remain pending. Final lifecycle completion does not resolve external acceptance or historical broad failures.
Required broad native race/frontend/previous integration gates remain historical failures. They were not rerun or repaired here.
Fresh focused passes do not replace those gates. First-build frontend checks cover the added menu module and production build.
Native amd64 loaded lifecycle, production billing/concurrency/deployed browser acceptance and published bootstrap remain unverified.
No arbitrary future-fork compatibility or publication readiness is claimed. No CPA deployment or private configuration read occurred.
