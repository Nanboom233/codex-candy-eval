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
