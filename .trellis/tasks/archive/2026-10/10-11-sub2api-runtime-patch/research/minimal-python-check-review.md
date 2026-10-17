# Minimal Python Linux Delivery Check Review

Review date: 2026-10-13. Scope is the authorized minimal Linux follow-up for
`.trellis/tasks/10-11-sub2api-runtime-patch`. This review covers the current
worktree, the eight-file package, the existing reviewed native runtime, and the
disposable arm64 lifecycle evidence. No CPA `.env`, deployment, paid request,
commit, push, publication, or task archival was performed.

## Findings (fixed)

No mechanical implementation defect was found in this review. No product code
or spec file was changed. `sub2api_install.py` is absent from the worktree and
from tracked files. `install.sh` and `install.ps1` are deleted as authorized.

## Findings (not fixed)

- The required broad upstream patched unit/race gate remains failed: 16 race
  warnings and nine top-level failures, with no compiler failure or timeout.
  No pristine A/B equivalence or attribution claim is made. The unrelated
  upstream fixtures were not repaired.
- Native amd64 loaded-process attestation remains blocked by Docker Desktop
  Rosetta. The arm64 lifecycle passed. Real paid billing, production
  concurrency, a deployed browser/payload inspection, and production browser
  acceptance remain unverified.
- The public online bootstrap is unverified because the tagged release was not
  published. Controlled HTTPS download tests passed, but the public URL and
  published archive were not exercised. Windows installation is explicitly out
  of scope.

These are acceptance or environment gates, not local mechanical issues. The
minimal delivery spec keeps them visible and main owns any later spec/task
synchronization.

## Implementation reviewed

`install.py` contains the one bounded `RELEASE_FACTS` literal and the complete
standard-library installer. `local_bundle`, `validate_bundle`,
`download_release`, and `extract_release` implement offline discovery,
HTTPS-only redirects, 192 MiB compressed and 384 MiB expanded limits, exact
allowlisting, path/type/link/duplicate checks, and pre-write archive
validation. `Installer.selection` reads `/dev/tty`, requires a choice even for
one candidate, accepts explicit `--container`, rechecks immutable ID, and keeps
selection separate from piped installer stdin. `Installer.inspect`,
`install`, `status`, `rollback`, `uninstall`, `maintenance`, and
`restore_binary` retain image/source/platform/startup/UID/GID checks, journal
and original-backup identity, candidate/staged/file/loaded-process hashes,
live drain, exact native revocation, and verified restoration.

The native runtime owners remain:

- `backend/internal/service/eval_host_runtime.go`:
  `EvalHostRuntime.BindRouter`, `Activate`, `Deactivate`, `ServeMaintenance`.
- `backend/internal/service/eval_runtime_ports.go`, `eval_identity.go`,
  `eval_key_runtime.go`, `eval_http_transport.go`, `eval_maintenance.go`,
  and `eval_target_adapter.go`.
- `backend/internal/service/openai_eval_directed.go` for canonical
  `POST /v1/responses` dispatch.
- `backend/internal/server/http.go` and the API-key middleware files for the
  public credential boundary.
- `backend/internal/handler/admin_codex_eval.go` and
  `backend/internal/server/routes/admin.go` for native administrator APIs.
- `backend/cmd/server/main.go` and `backend/cmd/server/maintenance_client.go`
  for live runtime and maintenance-client separation.

The maintenance socket is `/app/data/codex-eval/control.sock`. Evaluation APIs
are under `/api/v1/admin/codex-eval/`; the embedded shell is `/codex-eval/ui`.
`POST /v1/responses` is the only directed gateway route.

`scripts/package-sub2api.py` owns `validate_build`, `validate_payload`, and
`package`. It AST-validates one literal facts assignment, synchronizes it with
the development compatibility manifest and production provenance, checks the
production core tree and both binary identities, copies the tracked installer
and licenses byte-for-byte, and stages only the exact eight-file tree.

## Reused inputs and identities

`dist/sub2api-build/minimal-python-20261013/input-reuse.log` verifies the clean
pinned source checkout, production patch/core inputs, lockfiles, and both
reviewed native output hashes. Reuse is limited to unchanged native inputs and
outputs; package validation still checks the facts and provenance before copy.

| Identity | SHA-256 or value |
| --- | --- |
| Upstream source | `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` |
| Production patch/runtime stamp | `12a1865ff561c4c9f1c84df1c407fefb405db4093ec118d589bc651c065d9437` |
| Production staged tree | `4d28366dbdac4a597aa2f4c9bb509bc2feefa22ace8b98eb2c790251fd550fac` |
| Linux amd64 binary | `fd3109997b4d971e524a13ddd87b4b4c8f1cd8aad213036b0dae14ba33be213f` |
| Linux arm64 binary | `dd3760228f18cf63ab487ef6d3d433344d5c1b3d8ec89255f38a7deeb1bd129a` |

## Final package

Directory: `dist/sub2api-build/minimal-python-20261013/release-final/`.
Archive: `dist/sub2api-build/minimal-python-20261013/sub2api-codex-eval.tar.gz`.
The archive inspection reports eight regular files and no nonregular entries:

```text
install.py
sub2api-v0.2.15-linux-amd64
sub2api-v0.2.15-linux-arm64
LICENSE/CPA-MIT.txt
LICENSE/Sub2API-LGPL-3.0.txt
LICENSE/ModelTrace-MIT.txt
LICENSE/Renderers-MIT.txt
LICENSE/THIRD_PARTY_NOTICES.txt
```

The archive contains no Markdown, JSON, provenance, `SHA256SUMS`, checksum
sidecar, shell/PowerShell installer, test, fixture, wrapper, or development
manifest. The installer in the package is byte-identical to the tracked
`install.py`.

| Artifact | SHA-256 |
| --- | --- |
| Archive | `560cd57d30074ef9f3fbbe73640583fcc01dca024bf0fe92fc1e3dacd76d9eb8` |
| Tracked and bundled `install.py` | `c14d8c69277550ed501477edc40dfabce49fac9c6ce966c0d7b0fdb1ee5713a3` |
| Bundled Linux amd64 server | `fd3109997b4d971e524a13ddd87b4b4c8f1cd8aad213036b0dae14ba33be213f` |
| Bundled Linux arm64 server | `dd3760228f18cf63ab487ef6d3d433344d5c1b3d8ec89255f38a7deeb1bd129a` |

## Verification

Commands run during this review:

```sh
scripts/test-sub2api.sh check
python3 scripts/test_sub2api_build.py
python3 -m unittest discover -s tests/installer -p 'test*.py'
python3 -m py_compile install.py scripts/package-sub2api.py scripts/sub2api_build.py tests/installer/test_transactions.py tests/installer/docker_lifecycle.py
bash -n scripts/build-sub2api.sh scripts/test-sub2api.sh scripts/generate-wire.sh scripts/pnpm-support.sh scripts/pnpm
go vet ./...
go test -race ./...
node scripts/test-ui.cjs
git diff --check
PYTHONDONTWRITEBYTECODE=1 python3 scripts/package-sub2api.py dist/sub2api-build/installation-revision-20261013/out-cleanup dist/sub2api-build/minimal-python-20261013/release-check-2
```

Evidence paths:

- `dist/sub2api-build/minimal-python-20261013/independent-check.log`:
  72 helper tests, 46 installer/package tests, installer-fact validation,
  shell syntax, root vet/race, and UI checks passed.
- `dist/sub2api-build/minimal-python-20261013/independent-installer.log`:
  46 installer/package tests passed.
- `dist/sub2api-build/minimal-python-20261013/independent-package.log`:
  fresh package construction passed.
- `dist/sub2api-build/minimal-python-20261013/independent-syntax-diff.log`:
  empty output means `git diff --check` passed after Python compilation and
  shell syntax checks.
- `dist/sub2api-build/minimal-python-20261013/payload-final.log`:
  exact final payload, byte equality, and forbidden-artifact checks.
- `dist/sub2api-build/minimal-python-20261013/lifecycle-arm64-final.log` and
  `codexeval-fixture-29e2b5090d/evidence.json`: 15 disposable arm64
  checkpoints covering packaged auto-discovery installation, download-free
  status, uninitialized and active states, peer denial, expiry conflicts,
  original-runtime limits, migration/cache invalidation, recreation and
  reapplication, drain/revocation/original restore, uninstall, and candidate
  recovery. Docker inventory after the run has no codexeval fixture container,
  network, or named data volume.
- `dist/sub2api-build/minimal-python-20261013/input-reuse.log`: source, lock,
  patch/core, staged-tree, and binary reuse identities.
- `dist/sub2api-build/minimal-python-20261013/repo-check.log`: implementation
  checkpoint with 72 helper and 46 installer/package tests.

The required broad patched upstream race result remains the known failure
listed above. It is separate from the passing root `go test -race ./...` in the
local check command. No complete pristine A/B replay was run or claimed.

## Requirement mapping

### R1-R18

| Requirement | Review result |
| --- | --- |
| R1 | Four evaluation cores, native routing, and embedded UI remain in the reviewed production inputs; no new live billing claim. |
| R2 | Prompts, Codex context, baselines, scoring, model/effort, and batch controls remain in `internal/eval` and production overlays. |
| R3 | Directed transport enters native `/v1/responses`; native accounting and shared slots remain source-owned. Real charges remain open. |
| R4 | Server-owned accounting identity and fixed past expiry are retained; arm64 lifecycle verifies atomic activation and identity recovery. |
| R5 | Internal capability and exact owned-row checks remain in the native patches; expiry and peer conflicts are exercised. |
| R6 | Target adapter preserves exact requested targets and no fallback path. |
| R7 | Native admin middleware and administrator menu/API remain source-owned; no browser deployment check in this follow-up. |
| R8 | Durable typed history and restart/recovery behavior remain; lifecycle refresh/restart paths pass. |
| R9 | Official image, mount, startup, volume, UID/GID, journal, backup, and reversible replacement checks remain. |
| R10 | Install, inspect, status, rollback, uninstall, reapplication, drain, exact revocation, and honest incomplete cleanup are covered. |
| R11 | Embedded facts bind source/image/platforms and both static Linux binaries; native amd64 loaded attestation remains blocked. |
| R12 | README, operations docs, build helpers, and archive-only workflow describe Sub2API and `install.py`. |
| R13 | Four-tab native-theme UI and controls are retained; root UI checks pass, deployed browser behavior remains open. |
| R14 | Private UDS maintenance and canonical real-router dispatch are retained; lifecycle peer denial passes. |
| R15 | Production/test isolation, explicit allowlists, and payload exclusion checks pass. |
| R16 | One Python process owns HTTPS/download/extraction/Docker work; mandatory selection and offline/online tests pass. |
| R17 | Numbered actual-step output, target/result/next-action messages, stderr failures, and secret/raw-dump exclusions pass focused tests. |
| R18 | Exact eight-file Linux package with aggregated `LICENSE/` texts and no sidecar/checksum/runtime metadata passes. |

### A1-A24

| Acceptance | Evidence/status |
| --- | --- |
| A1-A8 | Preserved native evaluation, target, admin, capability, and billing paths are covered by production overlays and prior focused evidence; no new live-charge claim. |
| A9 | Durable history/restart behavior is covered by the run-state tests and arm64 lifecycle. |
| A10-A12 | Arm64 lifecycle verifies official image preservation, idempotence, candidate recovery, recreation detection, and reapplication. |
| A13 | Both static ELF outputs and embedded architecture facts pass; native amd64 loaded-process proof is open. |
| A14-A15 | Documentation/workflow and root UI checks pass; no deployed browser check was run. |
| A16-A17 | Maintenance peer/build identity, expiry conflicts, drain, exact revocation, restore ordering, and incomplete cleanup pass arm64 lifecycle/tests. |
| A18 | Production/test parity and exact payload validation pass; broad patched race gate remains failed. |
| A19 | Same-process online stdin, offline discovery, bounded HTTPS/archive validation, mandatory selection, cancellation, invalid input, and ID drift tests pass. |
| A20-A22 | Packaged activation, interrupted/reapplied identity, migration, PostgreSQL cache invalidation, and conflict recovery pass in the final arm64 lifecycle. |
| A23 | Output and failure-state tests pass; lifecycle records actual step/result messages. |
| A24 | The final archive and download-free maintenance path contain and use only the required standalone Python entry plus Linux binaries/licenses. |

### K1-K7

| Boundary | Review result |
| --- | --- |
| K1 | Protected state and UDS ownership checks remain in runtime/installer; no credential content was printed or packaged. |
| K2 | External credential/header/peer paths remain rejected by the native patches; maintenance identity is fixed to embedded facts and journal. |
| K3 | Only the owned row's immutable past expiry is exempted; native status, owner, binding, limits, quota, balance, subscription, and concurrency remain authoritative. |
| K4 | Readiness uses protected build/identity/dispatcher facts; retired public topology/probe inventories are not package prerequisites. |
| K5 | Cleanup drains and verifies exact revocation before restore; failed cleanup retains journal/backup and does not claim success. |
| K6 | Original pinned normal runtime paid rejection/read-only/session limits are demonstrated in the arm64 lifecycle; simple mode and deliberate edits remain outside the guarantee. |
| K7 | Recreation/reapplication detects missing or conflicting identity and does not silently rebind, rewrite expiry, or replace a conflicting row. |

## Final gate

Local minimal-package, installer, helper, root, and arm64 disposable checks pass.
The review does not claim the failed broad patched race gate, pristine A/B
equivalence, native amd64 loaded-process attestation, live billing,
production concurrency, deployed browser/payload acceptance, or published
online bootstrap. No release was published and no external target was changed.
