# Work Commit Proposal

## Staged Files

```text
.trellis/spec/backend/database-guidelines.md
.trellis/spec/backend/logging-guidelines.md
.trellis/spec/backend/native-inspection.md
.trellis/spec/backend/runtime-installation.md
patches/0001-evaluation-core.patch
patches/0005-inspection.patch
patches/0006-evaluation-ui.patch
patches/0007-sub2api-integration.patch
tests/backend/internal/codexeval/inspection_test.go
tests/backend/internal/handler/admin/codex_inspection_handler_test.go
tests/backend/internal/service/codex_inspection_settings_test.go
tests/ui/sub2api-dropdowns.cjs
```

## Exact Message

```text
fix(codex-eval): stabilize refresh and account inspection controls

Keep degraded wrong-answer inspections to one test per round.
Add per-account inspection switches with atomic metadata saves.
Remove runtime and notification prompts. Retain the next run time.
Align Candy account switches on narrow screens.
```

## Signing and Scope

The effective `commit.gpgSign` value is true. Use the configured signing policy without an override.
The complete staged diff has been reviewed. The work commit includes only the twelve files above.
Task artifacts remain under the Trellis lifecycle workflow. Generated binaries and browser files stay outside the work commit.
The user approved this exact work commit. Commit `75f5a3f` used the configured signing policy, and signature verification passed.
No push ran.

## Validation

The established build/check entry passed UI, focused Go, vet, build, and package checks.
After the narrow CSS correction, the affected UI check and standard build passed.
Browser acceptance confirmed independent saved switches, the retained next inspection time, and aligned narrow-screen controls.
See `implementation-evidence.md`, `check-evidence.md`, and `browser-evidence.md`.
