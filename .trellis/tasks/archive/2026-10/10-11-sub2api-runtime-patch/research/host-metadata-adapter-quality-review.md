# Host Metadata Adapter Quality Review

- Date: 2026-10-12
- Scope: `eval_target_adapter.go`, its adjacent tests, patch 004, and a fresh
  replay from the pinned Sub2API source.
- Verdict: **PASS** for the metadata adapter quality gate.

## Findings (fixed)

- File: `research/host-metadata-adapter-delivery.md` and
  `spec/backend/host-target-adapter.md`.
- Issue: After this review passed, the records still said that the quality
  review was running or remained required.
- Fix: The main session updated both records to show that specification and
  quality reviews passed. It records the current patch hash and labels the
  initial-delivery hash as superseded historical data.

- File: `backend/internal/service/eval_target_adapter_test.go` in patch 004.
- Issue: A nil `*APIKeyService` stored in the private resolver interface was
  not exercised by the earlier missing-dependency test.
- Fix: The specification review added a regression test for that typed-nil
  path. The native `ResolveDirectedTargetGroup` method checks its nil receiver,
  and the adapter returns the fixed metadata error without exposing backend
  details.

## Findings (not fixed)

None within this quality review's scope.

No adapter code findings remain. This review does not cover routes, probes,
bootstrap, Wire, lifecycle, UI, installer, deployment, or live billing.

## Quality Review

Exact resolution uses the existing directed group resolver, reloads the
requested account, and verifies returned IDs, OpenAI identity, active status,
supported native type, and hydrated membership. The nil API-key service test
exercises the typed-nil interface call. Read errors stay sanitized, while
caller cancellation and deadlines preserve their standard error identities.

List projection keeps membership pairs by querying each group separately,
deduplicates account hydration into one batch, rejects foreign or stale rows,
sorts the output, and copies only approved scalar fields into explicit DTOs.
Unknown statuses become `unavailable`. The tests seed credential, token,
proxy, extra-data, note, and backend-error sentinels and assert that DTO JSON
does not contain them. Snapshot mutation does not alter returned values.

The adapter does not use stored schedulability to authorize a run. The normal
gateway remains responsible for live admission, policy, quota, concurrency,
and billing.

## Verification

- Lint: pass. `gofmt -l` reported no files in the maintained or replay trees.
  No separate Go linter is configured for the package.
- Type check: pass. Default and `unit` `go vet ./internal/service` passed in
  both trees. Root `go vet ./...` and replay `go vet ./internal/codexeval`
  passed.
- Tests: pass. Focused adapter race tests passed with default and `unit`
  tags in both trees. Replay `go test -race ./internal/codexeval` passed.
  Root `go test -race -count=1 ./...` passed for the evaluation core and
  original plugin packages.
- Replay: pass. A new directory was created from pinned commit
  `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. Patches 001 through 004 each
  passed sequential `git apply --check --whitespace=error` before application.
  The 28 evaluation-core files were staged before replay tests.
- Parity: pass. All 3,256 backend files, 931 frontend files, and 28 evaluation
  core files matched the maintained workspace. Both adapter files matched
  byte-for-byte. All three lockfiles matched the pinned source and maintained
  workspace.
- Fresh replay directory:
  `dist/sub2api-build/scratch-metadata-quality-vSsAGFLZ`.
- Patch 004 SHA-256:
  `850af726f80bb2f379e0cf9b3f395f28e9d8e89356106f0643c39ee5558f4088`.
- Frozen patch hashes remain unchanged:
- 001: `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`.
  - 002: `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`.
  - 003: `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.

The full upstream Sub2API service race suite was not run. Its pre-existing
fixture failures and the decision not to repair unrelated fixtures are
recorded in `research/upstream-race-baseline.md`. This review makes no claim
about that suite or live acceptance.
