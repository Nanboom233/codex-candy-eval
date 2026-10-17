# Fresh native implementation evidence

Status: implementation, final build, fixed package and disposable runtime
verification completed. Main owns independent check,
task/spec updates and commit. No production deployment or paid provider request.

## Workspaces

- Original already-built source: `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/pristine-sub2api-v0.2.16.ifVBya/extracted/sub2api-01842576b3d5c1e2545b5a99ad555d36e28ed306`.
- Initial materialized implementation: `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8/source`.
- Workroot: `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8`.
- Old delivery files/artifacts preserved under `workroot/previous-delivery/`.
- Final build entry log: `workroot/build-entry.log`; session 60259 exited 0.
- Disposable fixture script: `workroot/native-smoke-complete.py`.
- Initial final fixture evidence: `workroot/smoke-complete/`; session 80949 exited 1
  on background Candy fixture assertion and removed all exact owned containers.
- Completed corrected fixture: `workroot/native-smoke-final-complete.py`,
  `workroot/smoke-final-complete/`; session 47947 exited 0. Runtime and actual
  rollback checks passed; all four owned containers and the owned network were
  removed. Final `docker ps -a -q --filter label=codex-eval-check=w8ax6xp8`
  returned no containers. Unrelated containers were untouched.
- Actual standalone JavaScript/API check: `workroot/native-ui-runtime-check.cjs`,
  evidence `workroot/smoke-final-complete/ui-runtime-checks.json`.

## Implemented boundaries

Fresh `native/backend/internal/codexeval` ports the current original four
evaluation runners, data, renderers and history, with memory-only task state,
native target IDs and direct native sender. Native `AccountTestService` supplies
exact-account sending and native Responses/CC bridge. Native admin handler
supplies task APIs, account sender API and standalone page.

Ordinary `native/sub2api.patch` modifies exactly three original files:

- `backend/internal/server/router.go`
- `backend/internal/server/routes/admin.go`
- `backend/internal/web/embed_on.go`

There are no native frontend source changes, new user/internal keys, billing,
scheduler, sticky logic or compatibility/overlay dependencies. `cleanup.py`
was removed from delivery; no cleanup/prune was executed.

Fresh root `install.py`, `scripts/build-sub2api.py`, wrapper and package script
deliver a static Linux amd64 binary with embedded original Vue frontend and
standalone evaluation page. Installer replaces only `/app/sub2api`, stores a
host backup, uses inspect JSON and native health/setup endpoints, and restores
stopped containers through Docker cp. No nc/buildx/helper image/source fetch.

## Completed evidence

- Earlier focused tests, touched-package vet, native compile and static build
  passed (`workroot/source/{focused-test,vet,compile,patched-build}.log`).
- Ordinary patch dry run against untouched original passed.
- Final fresh build entry (`--check`) passed after latest state-error and
  stream-retry refinements: evaluation/service/admin relevant tests, touched
  package vet, actual static Linux amd64 build, all exit 0.
- Fixed binary: `dist/build/sub2api-codex-eval`.
- Final package generation and archive inventory passed after final build;
  package session 70921 exited 0. Archive contains only binary, installer and
  five license texts; no cleanup, dev/build helpers, tests or fixture logs.
- JavaScript syntax checks passed for standalone page scripts.
- Actual disposable Docker installer replacement, `/codex-eval` page/CSP,
  native unauthenticated rejection, authenticated state and direct sender
  request to local fake upstream passed.
- All four background tasks completed through the actual native sender with
  the local fake upstream; Candy returned the expected 21, fingerprint and
  ModelTrace completed collection/record retrieval, and Pelican saved/retrieved
  its generated HTML. Synthetic answers verify execution, not model quality.
- Standalone JavaScript was loaded from the actual served page into the existing
  Node VM fake DOM pattern and called actual native HTTP APIs. It loaded native
  group/account/models, rendered account/history rows, exercised all four tab
  handlers, checked preserved Markdown/CJK/math behavior and HTML escaping,
  changed the real synthetic account active/inactive with native PUT, then ran
  and retrieved Pelican/ModelTrace/fingerprint records. No mocked API response.
- Actual installer failure was induced with a 20-byte malformed ELF. Container
  stopped; installer restored the prior executable using Docker cp, restarted
  it, passed original readiness endpoints and served the evaluation page again.
  `installer-failure.stderr` contains `Original executable restored.`
- Python AST syntax, shell syntax and `git diff --check` passed.

## Final artifact hashes

These identify the implementation handoff artifacts before independent review.
The checker subsequently fixed error-stream/status/cancellation/installation
issues and regenerated both fixed outputs; current verification is recorded in
`fresh-native-check.md`. These historical hashes do not identify the latest files.

- `dist/build/sub2api-codex-eval`: SHA256
  `7481db5904a23445db59e5218897f919009618d898c8b00a7a7b1b8fc271a299`.
- `dist/build/sub2api-codex-eval.tar.gz`: SHA256
  `35f9b56f814faac1785f103f3d4923fb1ac86c8d7a073e2f212d9e1a1b7835bf`.
- `file` reports stripped static x86-64 ELF.

## Fixture corrections and limits

Earlier fixture failures are retained as evidence. Alpine fixture needed UTC
because it lacked timezone data; published host ports changed after restart;
native compliance had to be acknowledged on the disposable administrator;
native group creation required rate_multiplier=1. Native background Responses
capability probe correctly classified the incomplete always-SSE fake upstream
as CC-only. The fixture now emits proper CC and Responses events. Its first
prompt matcher inspected the entire captured Codex context and emitted HTML
for Candy; corrected matcher reads only the last user prompt. These required
fixture changes, not authentication/transport bypass or extra product logic.

## Handoff and remaining limits

Main should dispatch the independent Trellis check against the approved task.
CUA browser inventory failed native pipe startup and IAB was unavailable, so
the Node VM check is functional JavaScript/API evidence, not a visual browser
or CSS layout claim. The fixture runs actual amd64 ELF through binfmt/Rosetta
on ARM64 Docker VM; this is not execution on N150. No production deployment,
real provider request, cache cleanup or business/history deletion occurred.
Existing original baseline evidence was reused without pristine/full/race/
coverage gates. Old unreferenced experimental directories remain historical
workspace evidence and are not inputs to the fresh build, installer or archive.
