# Implementation Handoff

## Current Candidate

The effective source tree is `/tmp/upstream-sync-native/source`.
The page is stable and ready for the main session's synthetic browser inspection.
The updated `tests/ui/sub2api-dropdowns.cjs` passed against that tree on 2026-10-19.

The detector, shared retry runner, native sender, runtime, history, and page changes are present in that tree.
Unversioned native history immediately saves retained Candy and Pelican records through the existing atomic saver.
The page removes ModelTrace and fingerprint modes while retaining Settings and native account controls.

Focused native evaluation, service, and administrator tests passed against the candidate.
The main session resolved the Candy test conflict and retained the unchanged native test.
The implementation adds the missing shared-runner owner and extends current runtime, sender, inspection, and UI owners.
Production patches 0001, 0002, 0004, 0005, and 0006 now contain the candidate.
The first full build entry passed before the stricter history-marker correction.
Command: `python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check`.
Complete log: `/tmp/upstream-sync-native-build-check.log`.
That run passed patch application, frontend checks, UI, focused Go tests, vet, server build, and installer test/build.
The main session then identified a history-marker boundary mismatch.
The reader now distinguishes an absent marker from an explicit null, zero, or unsupported marker.
Only an absent marker selects migration. Any present marker must decode to version 1.
The existing persistence test now checks explicit-null protection without a future-format matrix.
The final build rerun uses the corrected marker boundary and records `/tmp/upstream-sync-native-build-check-final.log`.
It passed frontend checks, UI, focused Go tests, vet, and the Linux server build.
The installer tests passed: `ok sub2api-codex-eval-installer 11.134s`.
The installer build passed. The final build entry exited with status 0.
The main session owns README, AGENTS, task/spec updates, and browser acceptance.

The focused transport subtest passed through the real native HTTP repository.
Account lookup failures stay unmarked. An actual interrupted network request carries the new marker and preserves its underlying error.
The native runtime test also passed with a parse retry and fingerprint cancellation.

The preview fixture's synthetic ranking rows should include a numeric `score` when checking score display.
Real detector rows always provide that field.

## Check Results

`node tests/ui/sub2api-dropdowns.cjs /tmp/upstream-sync-native/source` passed after the final page edits.
It reported: `Settings and account checks passed`.

The focused Go selector passed for evaluation, service, and administrator packages.
The evaluation package reported `ok .../internal/codexeval 0.761s` before its added cancellation assertion.
The service and administrator packages reported `ok .../internal/service 8.143s` and `ok .../internal/handler/admin 1.226s`.
The changed runtime test then passed separately: `ok .../internal/codexeval 1.937s`.
The native transport classification subtest passed: `ok .../internal/service 2.660s`.
The final corrected candidate reports `ok .../internal/codexeval 2.041s`, `ok .../internal/service 8.092s`, and `ok .../internal/handler/admin 1.218s`.
Changed Go files passed `gofmt`.
The affected source, tests, and notices passed `git diff --check`.

The first persistence assertion reused a decoded map and retained earlier fixture entries.
Resetting that local test snapshot fixed the assertion. No product change was needed for that failure.

## Changed Files

- Production patches: 0001, 0002, 0004, 0005, and 0006.
- Build: `scripts/build.py` changes only the focused test selector.
- Tests: current Codex request, fingerprint, inspection, runtime, sender, stream-lifetime, and UI owners.
- New test owner: `tests/backend/internal/codexeval/model_test.go`.
- Test data: new LM-detector golden fixture replaces ModelTrace golden data.
- Notices: `licenses/LM-detector-MIT.txt` and `licenses/THIRD_PARTY_NOTICES.txt`.
- Retired algorithm test: `modeltrace_test.go` is removed.

README and AGENTS changes belong to the main session.
No task/spec/workspace path was staged by this implementer.
No commit, branch switch, remote operation, installer transaction, or release ran.

## Final Delivery

Final build log: `/tmp/upstream-sync-native-build-check-final.log`.
Effective source for review and browser work: `/tmp/upstream-sync-native/source`.
Server: `dist/build/sub2api-codex-eval`.
Installer: `dist/build/sub2api-codex-eval-installer`.
Both outputs are stripped static Linux/amd64 ELF executables.

The final product checks are complete.
The main session owns the independent Trellis review, spec synchronization, final commit approval, and approved Git completion.
No unresolved implementation issue remains at this handoff.
Browser acceptance remains with the main session and uses synthetic data.

## Review Repair: Request Count Copy

The independent review found a stale Pelican caption: `每个账户 1 次请求`.
The new retry runner can send multiple attempts for one logical test.
The caption is removed from `web/ui.html` and patch 0006.
The Candy and Pelican page modules contain no other fixed request-count caption.
No replacement explanation or small print was added.

The existing UI harness passed against the changed page. No new test was added.
The prior backend checks remain valid because this repair changes only markup.
The build regenerated the server and installer with the final markup and exited with status 0.
Command: `python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check`.
Log: `/tmp/upstream-sync-native-build-ui-final.log`.
This build reuses the valid backend and installer test results instead of repeating them.
Both existing output paths now contain the repaired page.
`git diff --check` passed for patch 0006.

## Merge Choices

Retired CPA runtime, docs, examples, and UI test tooling remain deleted.
The golden detector fixture remains under the native test package.
Runtime detector resources live only in the effective production tree and patch 0002.
The root detector license reuses the prior identical MIT text under `licenses/LM-detector-MIT.txt`.
No installer, release workflow, vendor patch, or native integration patch change is planned.

No unmerged entry remains after the main session resolved the Candy test conflict.

Patch 0002 places its generated Git binary hunk before text assets to keep the ordinary repository whitespace check clean.
Patch 0004 uses zero-context hunks for the two original-file changes, consistent with the build's existing `--unidiff-zero` apply mode.
These patch-format choices do not change the effective source.

## Tool Boundary

The checkout hook rejected `git checkout --ours` because the merge has local changes.
The main session took ownership of those three text conflicts.
The edit hook requires whole-file removals through macOS Trash.
The implementation used that reversible operation for the retired temporary source and algorithm tests.
