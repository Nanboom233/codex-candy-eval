# Implementation Result

The coordinator now uses the saved account verdict to select confirmations.
Absent verdicts start as normal, regardless of actual priority.
Confirmed transitions save the verdict before the optional priority update.
Stable verdicts do not write metadata or actual priority.
Existing notifications remain on the successful priority-update path.

## Changed Files

- `patches/0005-inspection.patch`: verdict type, coordinator persistence hook, native projection, and atomic key-only persistence.
- `patches/0007-sub2api-integration.patch`: preserve the verdict during full account edits.
- `tests/backend/internal/codexeval/inspection_test.go`: replace the old matrix with one adjustment-off lifecycle and extend existing boundary assertions.
  Add only the verdict-save and priority-write failure cases.
- `tests/backend/internal/handler/admin/codex_inspection_handler_test.go`: check the normal default, saved projection, and key-only metadata update.
- `tests/backend/internal/service/codex_inspection_settings_test.go`: extend the existing full-edit preservation assertion.

## Validation

Changed Go files passed `gofmt`. `git diff --check` passed.

The final candidate passed this entry with exit code 0:

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check
```

The entry passed frontend checks/build, installer version checks, existing UI assertions, focused Go tests, and vet.
It built and packaged the Linux/amd64 executable.

The adjustment-off lifecycle verified test counts of 3, 1, 3, and 1 across degradation, stable wrong, recovery, and stable correct.
It verified two verdict writes and unchanged actual priority.
Existing assertions cover mixed results, provider errors, cancellation, incomplete opt-out work, and manual fixed counts.
The failure cases verify save-before-priority ordering and durable verdicts after a failed priority write.

Log: `/tmp/inspection-verdict-state-build-check.log`.

Outputs:

- `dist/build/sub2api-codex-eval`
- `dist/build/sub2api-codex-eval.tar.gz`

No runtime container or provider quota was needed. The original upstream source remains unchanged.
The page remains unchanged. Prior browser evidence applies.

## Handoff

The candidate is ready for the Trellis check owner.
Reuse the completed build/check evidence for unchanged inputs.
No implementation blocker remains. No commit, push, or deployment was performed.
