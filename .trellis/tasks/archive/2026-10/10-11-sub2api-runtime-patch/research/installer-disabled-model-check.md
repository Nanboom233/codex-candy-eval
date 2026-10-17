# Independent Trellis check: U7–U9

Result: passed for the current scoped repair. No additional product edits were
needed. This review does not authorize a commit or deployment.

## Findings (fixed)

None. Reviewed the current PRD/design/implementation U7–U9 decisions, current
check context and implementation evidence, and the actual changed installer,
catalog, load/toggle flow and user documentation.

## Findings (not fixed)

No unresolved product findings in this repair. Main retains ownership of
spec/task synchronization and the new commit review. Historical checklist
entries and prior reviews do not introduce additional acceptance requirements.

## Verified behavior

- U7 changes readiness output only: startup progress omits transient transport
  errors, and success is printed after both native checks pass. Final failures
  preserve the last endpoint/error. Deadline, container state/address inspection,
  request bounds and restoration control flow remain unchanged.
- U8 creates default backups with `tempfile.mkdtemp(..., dir="/tmp")`, giving
  private directories and printing the exact retained executable path. Explicit
  `--backup-dir` retains its former timestamp layout. Backup completes before
  mutation, and failure/interrupt restoration still uses the saved executable.
  No cleanup, helper download or deletion of existing backups was added.
- U9 retains all accounts in `accountList` while querying models only for
  `schedulable !== false` accounts. Error names index that same filtered list.
  Loading succeeds only after discovery succeeds, preventing a subsequent poll
  from silently clearing enabled-account failures. Forced refresh resets the
  loaded flag; the existing native schedulable toggle invokes that refresh.
  Successful raw suggestions and direct model inputs remain available.

## Evidence reviewed

Evidence root:
`/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/sub2api-codex-eval-w8ax6xp8`.

Read `installer-progress-check.py` and both resulting logs. They use real local
TCP resets: one reset produces friendly progress followed by ready; persistent
resets produce bounded failure with the actual `/health` endpoint and reset.

Read `disabled-model-check.cjs` and `disabled-model-result.json`. Actual page
load/catalog/toggle functions execute with synthetic APIs: disabled discovery
is skipped, enabled failure survives polling, disabling clears its notice,
reenabling restores it, all rows and four direct model inputs survive. This is
focused JavaScript behavior evidence, not visual acceptance.

Read `installer-tmp-backup-check.py`, `installer-tmp-fixture/outcome.json` and
the success/failure-restoration logs. A real cached Docker fixture uses the
actual native executable for version/copy/replacement/restoration; Node supplies
reset/health/setup behavior. The successful run becomes ready after resets;
the truncated candidate times out and restores the original. Independently
checked both retained backup directories are mode0700 and their executable
bytes, plus the restored executable, match the final binary.

The retained three-second fixture failure is explicit: deliberate resets after
restart exceeded that fixture's restoration deadline. The eight-second fixture
rerun passed; product timeout/bounds were not increased. The earlier actual
native stopped-container restoration evidence remains applicable and was not
repeated.

Read the final build log: native frontend i18n/TypeScript/Vite and static native
Go build completed. Existing relevant Go tests/vet and actual Chrome
sidebar/origin/iframe/model-forwarding evidence are reused because those paths
were unchanged. No race, pristine rebuild, full suite or platform matrix was
added.

## Independent current checks

- JavaScript syntax: `node --check` on changed catalog passed.
- Python syntax: AST compilation of current installer passed.
- Changed-file whitespace: passed.
- Archive inventory: binary, installer and five necessary license texts only.
  Archived installer and binary match current source/output bytes.
- Binary SHA256: `23f2565f340b7c5ad9c542bb525cef5622fbebc8f50a1aa56e6fac2080290d4d`.
- Archive SHA256: `71ce4c4b0eaa802bcfd610b6bb2e9abf431f620c8185cea81f18bec9b3ba08c0`.

Hashes identify reviewed artifacts; they are not a new source/parity gate.
No product/artifact inputs changed during this check, so another rebuild is
unnecessary. No private environment, production/paid calls, CPA deployment,
pruning or business/history deletion occurred. N150/iStoreOS and production
HTTPS/proxy execution remain unperformed.
