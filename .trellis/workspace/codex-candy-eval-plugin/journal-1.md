# Journal - codex-candy-eval-plugin (Part 1)

> AI development session journal
> Started: 2026-10-11

---



## Session 1: Initialize Trellis and portable AI integration
<!-- trellis-session: v=2 fp=6d77a5b80789d5fb -->

**Date**: 2026-10-11
**Task**: Initialize Trellis and portable AI integration
**Branch**: `main`

### Summary

Initialize source-backed plugin specs and portable Trellis AI integrations. Keep task artifacts in the native Trellis lifecycle.

### Main Changes

- Define JSON storage, errors, concurrency, and host ABI guidance.
- Track Trellis 0.6.17 and portable Claude, Codex, and Devin entry points.
- Use git-commit for ordinary work and native Trellis for artifact commits.
- Leave machine-only settings and the parallel task untouched.

### Git Commits

| Hash | Message |
|------|---------|
| `2e01ddb` | docs(trellis): Define source-backed plugin guidelines |
| `f84347d` | chore(trellis): Track portable runtime and AI integrations |

### Testing

- [OK] go vet ./... passed.
- [OK] go test -race -count=1 ./... passed.
- [OK] node scripts/test-ui.cjs passed.
- [OK] Python syntax, JSON/TOML parsing, scope, and spec-link checks passed.
- [OK] Work and archive commit signatures passed local GPG verification.

### Status

[OK] **Completed**

### Next Steps

- Keep the retained machine-only settings outside unrelated commits.
- Continue parallel tasks in their own Trellis context.


## Session 2: Complete the native Sub2API v0.2.17 migration
<!-- trellis-session: v=2 fp=6321206dc6aaf2a1 -->

**Date**: 2026-10-18
**Task**: Complete the native Sub2API v0.2.17 migration
**Branch**: `feat/sub2api-runtime-patch`

### Summary

Commit the complete native evaluation migration, finish the approved patch-series layout, and archive the task.

### Main Changes

- Port four evaluations, native administrator sending, sidebar and iframe, history, Candy inspection, database settings, priority adjustment, and notifications.
- Use seven production patches, separate tests, two Python command entries, and README instructions. Remove obsolete CPA contents and redundant source and documentation trees.
- Create signed work commit 6cac3a2 and archive commit ecf1e37. Recover the native archive auto-commit from its nonexistent untracked source path using the archived task path only.

### Git Commits

| Hash | Message |
|------|---------|
| `6cac3a2` | feat(sub2api): port the complete evaluation plugin to v0.2.17 |

### Testing

- [OK] Reuse completed native Go, frontend, retained-container, and browser evidence. U22 patch composition, byte preservation, focused Go/UI/installer checks, and independent review passed.
- [OK] Repack the seven-entry installation archive with the unchanged Linux/amd64 executable. No U22 product rebuild or Docker operation.

### Status

[OK] **Completed**

### Next Steps

- The historical context canceled report remains unresolved. Keep its evidence limits in the archived task. Local configuration changes remain outside these commits. No remote push ran.


## Session 3: Fix account refresh and inspection controls
<!-- trellis-session: v=2 fp=0eebcdd215c203e6 -->

**Date**: 2026-10-18
**Task**: Fix account refresh and inspection controls
**Branch**: `fix/inspection-controls`

### Summary

Fixed refresh readiness and stale-priority inspection counts. Added persistent per-account inspection switches and per-send opt-out checks. Removed runtime and notification prompts while retaining the next inspection time. Corrected narrow Candy switch alignment after browser feedback. Existing build/check and affected layout checks passed. Browser acceptance and signed work commit completed. Finished the archive commit with its destination path after the native auto-commit referenced an untracked source path.

### Git Commits

| Hash | Message |
|------|---------|
| `75f5a3f` | fix(codex-eval): stabilize refresh and account inspection controls |

### Status

[OK] **Completed**


## Session 4: Independent inspection verdict state
<!-- trellis-session: v=2 fp=a2ff233e2213cfde -->

**Date**: 2026-10-18
**Task**: Independent inspection verdict state
**Branch**: `fix/inspection-controls`

### Summary

Separate saved inspection verdicts from actual priority and retain the single reusable test container.

### Main Changes

- Default absent verdicts to normal. Save confirmed transitions before optional priority writes.
- Keep stable rounds to one test without compatibility, synchronization, retries, or notification changes.
- Retain sub2api-codex-eval-test and its fixture data after checks.
- Resolve the archive pathspec failure with the generated archive message and actual archived-directory scope.

### Git Commits

| Hash | Message |
|------|---------|
| `8c22806` | fix(codex-eval)!: separate inspection verdicts from priority |
| `aecf699` | docs: retain the reusable Sub2API test container |

### Testing

- [OK] The established build/check and final Trellis review passed. The adjustment-off lifecycle verified 3/1/3/1 sender calls.
- [OK] Both work commits and the task archive commit passed signature verification.

### Status

[OK] **Completed**


## Session 5: Complete model control and Go installer delivery
<!-- trellis-session: v=2 fp=8002efb7cb55a7e5 -->

**Date**: 2026-10-19
**Task**: Complete model control and Go installer delivery
**Branch**: `feat/inspection-model-control`

### Summary

Complete confirmed-state model control, Go installer delivery, Settings entry, and the Chinese user guide.

### Main Changes

- Apply model control and priority after each confirmed account detection, including stable states.
- Deliver the self-extracting Go installer and publish it through the existing release workflow.
- Fix Settings-first group loading and document user operations with original project credits.
- Create signed work commit c4cc8e8 and archive only the completed model-control task.
- Recover the archive commit with its destination path after Trellis referenced an untracked source directory.

### Git Commits

| Hash | Message |
|------|---------|
| `c4cc8e8` | feat(codex-eval)!: control degraded models and add a Go installer |

### Testing

- [OK] Reuse passing focused backend, installer, UI, actual-browser, build, and independent-review evidence.
- [OK] Pass the final staged diff check and verify the work and archive signatures.
- [OK] Keep production installation and GitHub Actions execution outside this local handoff.

### Status

[OK] **Completed**

### Next Steps

- Fast-forward local main to the completed task commits and preserve unrelated working changes.


## Session 6: Commit remaining account tags and repository guidance
<!-- trellis-session: v=2 fp=8d132af266d75c34 -->

**Date**: 2026-10-19
**Task**: Commit remaining account tags and repository guidance
**Branch**: `main`

### Summary

Commit the saved-verdict account tags and the remaining repository guidance after final user approval.

### Main Changes

- Show degraded account tags from saved inspection verdicts in all four evaluation lists.
- Commit the existing recovery-refresh assertions and account-display contract.
- Keep the Trellis entry guide and record the return to local main rule.
- Remove the old degraded-tag task directory and clear its active pointer at the user request.

### Git Commits

| Hash | Message |
|------|---------|
| `47e8406` | feat(codex-eval): show saved degraded account verdicts |
| `a2c212c` | docs: streamline repository instructions and finish rules |

### Testing

- [OK] Reuse passing UI, build, and normal-width and narrow-width visual evidence for unchanged inputs.
- [OK] Pass the staged diff check and verify both work commit signatures.

### Status

[OK] **Completed**
