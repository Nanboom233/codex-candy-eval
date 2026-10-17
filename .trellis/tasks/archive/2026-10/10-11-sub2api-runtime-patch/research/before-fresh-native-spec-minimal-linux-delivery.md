# Minimal Linux Binary Replacement Delivery

## 1. Scope / Trigger

Latest user explicitly restores binary replacement and requests a patched Sub2API v0.2.16 for N150 Linux/amd64.
This supersedes install-time full builds, universal source payload, proposed catalog review and historical v0.2.15 gates.
Build the actual plugin-enabled runtime here once. Installation performs no source, frontend, dependency or Go compilation.
Include user-requested independent cleanup.py. Native dispatch/accounting/admin/history contracts remain authoritative.

## 2. Signatures

```sh
python3 install.py
python3 install.py install --container sub2api
python3 install.py status --container sub2api
python3 install.py rollback --container sub2api
python3 install.py uninstall --container sub2api
python3 cleanup.py
python3 cleanup.py --dry-run
python3 cleanup.py --select <exact-listed-id>
```

Keep one Linux Python installer plus the separately requested cleanup utility. Explicit container selection is mandatory.
Keep tests/build helpers outside the release. Standard library Python and Docker CLI remain user prerequisites.
Use fixed delivery outputs: dist/build/sub2api-codex-eval.tar.gz and dist/build/sub2api-codex-eval.
Keep the versioned inner binary name for the v0.2.16 payload. Active installer asset lookup, packaging, CI and docs use the fixed archive basename.
Reuse the reviewed runtime when changing output paths. Historical archives and reports remain evidence.

Developer packaging uses the fixed raw binary and writes the fixed archive:

```sh
python3 scripts/package-sub2api.py
```

No-argument packaging uses automatic system-temp staging. Explicit staging and binary/archive overrides remain developer options.

## 3. Contracts

Build source is verified v0.2.16 tag, full commit01842576b3d5c1e2545b5a99ad555d36e28ed306, Go1.27.2.
Deliver Linux/amd64 patched runtime for the user's standard /app/data layout. Do not reuse the vanilla arm64 timing binary.
Retain source-version information and native plugin integration. Build frontend and backend with actual adapters/core, no fixtures in payload.
Prebuilt installation uses local candidate loading and current functional layout/readiness/transaction checks. No Git/build fallback.
Do not reintroduce host binary checksum/ELF/build-metadata/image allowlists. A known other host release or conflicting data layout fails before mutation.
Use immutable selected container ID, protected journal/backups, atomic stage/replace, restart, native maintenance and healthy restoration.
Keep clientless target probing through owned network helper; preserve complete raw errors and traceback chains.
Minimal payload has eight regular files: install.py, cleanup.py, sub2api-v0.2.16-linux-amd64 and five required LICENSE/ texts.
No Markdown, test/dev/source/build material, manifests or redundant SHA sidecars.
Archive download/extraction still bounds size and rejects traversal, duplicate entries, symlinks and special files.

cleanup.py inventories old installer-owned stopped tool/source/http containers, unique codexeval-builder names and stale standard-temp workspaces.
Exact prefix, current-use and ancestry checks are required. Never remove business containers, application data, databases, Redis or volumes.
Anonymous fallback containers and shared unused Go/Node images require explicit exact-ID selection and refreshed reference checks.
Retain probe helper image needed by the new installer, shared builders and any referenced images. No global Docker/image/builder prune.
Missing optional Buildx must not prevent ordinary container/image/temp inventory on the user's Docker installation.
Enumerate containers with docker ps -a --quiet --no-trunc and read facts through JSON inspect. Avoid Docker --format dependencies.
Detect optional Buildx from successful Docker help before calling it. Parse plain builder listings conservatively and retain ambiguous entries.
Recognize digest-created tool candidates through inspected image tags. Require explicit selection and recheck references before removal.
Print concise plans/results with exact removed and retained items. Keep raw Docker errors on failures.
Default cleanup inventories resources and requires explicit selection. --select accepts repeatable exact listed IDs for unattended cleanup.
--dry-run always performs inventory without deletion. Without a terminal or --select, keep every resource.

## 4. Validation & Error Matrix

| Condition | Required outcome |
| --- | --- |
| Missing/canceled container selection | Fail before mutation |
| Target lacks HTTP clients | Owned namespace helper probes actual native routes |
| Known wrong release/architecture/data layout | Fail before backup/replacement, no downgrade |
| Invalid/oversized archive or unexpected file | Reject before extraction effects |
| Candidate/readiness fails after replacement | Restore original and verify health, preserve failure/journal |
| Native cleanup ownership/revocation conflict | Keep journal/backup/history and report incomplete cleanup |
| Old build resource active/ambiguous/shared | Retain and show reason |
| Explicit selected unused image/container changes references | Refuse removal |
| Default packaging invoked repeatedly | Atomically replace the fixed archive and remove temporary staging |
| Docker CLI lacks formatting flags or optional Buildx | Quiet IDs and JSON inspect work without formatting; absent/ambiguous builders remain |
| Docker help/list/inspect has genuine failure | Preserve original output and exception, no deletion |

## 5. Good / Base / Bad Cases

Good: selected v0.2.16 standard-layout container receives reviewed prebuilt amd64 plugin runtime and is activated in native panel.
Base: cleanup dry run lists confirmed leftovers and unused tool candidates without effects.
Bad: call old v0.2.15 or vanilla arm64 binary a v0.2.16 patched runtime.
Bad: make install silently compile or run Docker global prune for cleanup.

## 6. Tests Required

Verify annotated tag/commit, full actual patched compilation, embedded frontend and developer artifact architecture/version.
Run focused native contract/type tests and actual disposable lifecycle where executable on local Docker; distinguish emulated amd64 from native.
Run root gofmt/vet/race/UI and current installer transaction/rollback/clientless regressions, package-source/license parity and meaningful cleanup mixed-resource tests.
Current regressions live in tests/installer/test_prebuilt.py, test_cleanup.py and test_transactions.py.
Developer build and test materialization use scripts/prebuilt_dev.py and scripts/prebuilt_test.py outside the release.
scripts/package-sub2api.py enforces the eight-file allowlist and source/license parity.
Validate raw tar members with Python tarfile. Reject AppleDouble metadata, unexpected directories, links and duplicate entries.
Create archives from explicit allowlisted paths. Native macOS tar listings may hide AppleDouble members and cannot prove payload parity.
Exercise archived installer and cleanup script, ensure no installer source/Git/build work or test fixtures enter payload.
Retain known historical broad failures and unverified iStoreOS/live billing limits explicitly. Do not substitute static build for loaded lifecycle.

Current v0.2.16 evidence is recorded in the [implementation](../../tasks/10-11-sub2api-runtime-patch/research/v0.2.16-prebuilt-implementation.md) and [independent review](../../tasks/10-11-sub2api-runtime-patch/research/v0.2.16-prebuilt-check-review.md).
Actual patched compilation, frontend lint/2869 tests, native vet and integration passed.
The prebuilt baseline installer suite ran 98 cases with 94 passes and four explicit Docker skips.
Baseline independent archived cleanup and transactions passed ten and 51 cases. Implementer-owned actual clientless transport passed 16 cases.
The 18-checkpoint actual Docker lifecycle passed under amd64 emulation on macOS arm64.
Complete native unit race failed with 51 warnings and 153 top-level failures. Attribution lacks a pristine comparison.
Native N150/iStoreOS, live billing and direct Python 3.9 execution remain unverified. Publication remains blocked by required gate failures.

The fixed-output revision uses the same binary and records its checks in [fixed-output implementation](../../tasks/10-11-sub2api-runtime-patch/research/fixed-output-implementation.md).
The [independent fixed-output review](../../tasks/10-11-sub2api-runtime-patch/research/fixed-output-check-review.md) confirms default packaging and canonical source/archive parity.
The revision ran 99 installer cases with 95 passes and four skips. Archived installer transactions passed 52 cases.
Check fixed defaults, download asset and CI upload/publish names together. Verify the actual canonical archive through raw tarfile and installer extraction.
Output-only revisions need affected Python/package checks and unchanged runtime identity. Do not rebuild the runtime for path changes.

The [Docker formatting review](../../tasks/10-11-sub2api-runtime-patch/research/docker-format-compatibility-check-review.md) covers the latest CLI compatibility fix.
Its current suite ran 108 cases with 104 passes and four skips. Nine real-subprocess compatibility cases passed against source and archived scripts.
Current Docker read-only listings passed. Simulated legacy CLI checks do not prove real iStoreOS acceptance.

## 7. Wrong vs Correct

Wrong: extend timeout and compile the entire server on the user's host again.
Correct: build requested runtime here and ship a complete binary replacement package.

Wrong: delete generic Go/Node images by repository name or remove shared builders/volumes.
Correct: remove only confirmed owned leftovers or explicit unused exact IDs after reference checks.
