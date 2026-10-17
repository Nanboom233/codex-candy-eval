# U20 CPA Cleanup Check

Review date: 2026-10-18.

## Findings (fixed)

None. The U20 implementation is internally consistent and required no local
code or documentation repair.

## Findings (not fixed)

None within the authorized U20 cleanup scope.

The native page still contains a few non-executing historical host names in
CSS or JavaScript comments, plus the existing footer attribution. They do not
reference a removed file, affect runtime behavior, or add a build input. U20
explicitly preserves the native source and fixed binary, so this check does
not edit those comments or trigger a rebuild.

## Scope review

The 56 planned deletions are present:

- `cmd/`
- `internal/`
- `examples/`
- `docs/inspection.md`
- `docs/images/`
- `scripts/deploy.sh`
- `scripts/test-ui.cjs`
- `go.mod`
- `go.sum`
- `.env.example`
- the already-authorized `install.sh` and `install.ps1` deletions

The active Sub2API paths remain present: `native/`, `install.py`, the four
current build/test/package scripts, retained installer and UI fixtures, the
current operations guide, and all required license files. The task
`history/` directory is absent as required by U19. No active command,
workflow, fixture, or UI harness references the removed CPA source graph.

The fixed delivery outputs remain in `dist/build/`. The archive contains only
the executable, `install.py`, and the five required license/notice files. The
package uses `LICENSE/Project-MIT.txt`; no current file refers to the removed
`CPA-MIT.txt` name. Sub2API v0.2.17, ModelTrace, and renderer attribution is
present in the current notices and license inventory.

The retained Docker fixture checks for the single named
`sub2api-codex-eval-test` container and fails when it is absent. It uses
`docker inspect`, `docker cp`, `docker start`, and `docker exec` only; it does
not run `docker create` or `docker run`, so later checks reuse the same
container. The task record documents removal of the 11 identified earlier
task containers while preserving unrelated Docker resources.

## Verification

- Lint: N/A; no configured Python/Go linter is part of this documentation and
  deletion-only review.
- TypeCheck: N/A; U20 does not change typed source, and the root Go module was
  intentionally removed.
- Tests: existing product and browser evidence reused because native inputs and
  the fixed binary are unchanged.
- Python AST: pass for `install.py`, current build/package scripts, and the
  retained fixture files (6 files).
- Shell syntax: pass for the current build wrapper, fixture start script, and
  test entry point.
- Node syntax: pass for `tests/ui/sub2api-dropdowns.cjs`.
- Markdown links: pass for 76 relative links in current README, guide, docs,
  and spec files.
- Package inventory: pass; `dist/build/sub2api-codex-eval.tar.gz` has exactly
  7 expected entries.
- `git diff --check`: pass.
- No Docker operation, rebuild, deployment, staging, commit, or hash
  calculation was performed by this review.

The main session owns task metadata synchronization, the work-commit
candidate, staging, and final approval.
