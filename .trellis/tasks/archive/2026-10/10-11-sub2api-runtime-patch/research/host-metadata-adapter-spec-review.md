# Host Metadata Adapter Specification Review

- Date: 2026-10-12
- Scope: the two target-adapter service files and patch 004 against
  `spec/backend/host-target-adapter.md`.
- Pinned source: Sub2API v0.2.15 at
  `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.
- Verdict: **PASS** for the metadata adapter specification boundary.

## Findings (fixed)

- File: `backend/internal/service/eval_target_adapter_test.go`, maintained by
  `patches/sub2api/v0.2.15/004-codex-eval-target-adapter.patch`.
- Issue: Passing a nil `*APIKeyService` to the constructor stores a typed-nil
  pointer in the private resolver interface. The existing missing-dependency
  test had a nil admin dependency, so short-circuiting prevented it from
  exercising that interface call.
- Fix: Added `TestEvalTargetAdapterNilAPIKeyDependency` with a nonnil admin
  dependency. It verifies that resolution returns the fixed unavailable error
  without panicking. The native `APIKeyService.ResolveDirectedTargetGroup`
  checks both a nil receiver and a missing group repository, so no production
  guard, reflection, or new abstraction is needed. Patch 004 was regenerated
  from the complete adapter files.

## Findings (not fixed)

- None within the metadata-adapter scope. The broader host-runtime plan remains
  **NEEDS-REFINEMENT** for provisioner-error recovery and observation-lock
  wording, as recorded in `host-runtime-plan-final-review.md`. Those runtime
  state-machine decisions are outside patch 004. This review does not cover
  routes, probes, bootstrap, Wire, lifecycle, UI, installer, deployment, or
  live billing.

## Specification Review

The implementation matches the adapter contract: exact positive target IDs,
fresh directed-group resolution, account identity/platform/type/status and
membership checks, and name-only core metadata. Listing preserves per-group
membership, hydrates deduplicated accounts in one batch, uses explicit
credential-free DTOs, projects status through an allowlist, sorts snapshots,
and leaves live schedulability and gateway admission to the normal request
path. Context cancellation is preserved, other read failures are sanitized,
and missing dependencies fail closed.

## Verification

- Lint: pass; `gofmt -l` reported no files in the maintained workspace or
  replay tree. The repository check contract specifies gofmt, vet, and race
  tests; no separate Go linter is configured for this package.
- TypeCheck: pass; default and `unit` `go vet ./internal/service` passed in
  both trees.
- Tests: pass; focused default and `unit` race tests passed in both the
  maintained workspace and replay tree.
- Replay: pass; sequential patches 001 through 004 were applied to pristine
  pinned source per the delivery record. This review confirmed byte equality
  for both adapter files between the maintained workspace and replay tree.
- Patch 004 SHA-256:
  `850af726f80bb2f379e0cf9b3f395f28e9d8e89356106f0643c39ee5558f4088`.
- Frozen patch hashes 001 through 003 remain unchanged.
