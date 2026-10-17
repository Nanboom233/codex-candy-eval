# Host Metadata Adapter Delivery

- Date: 2026-10-12
- Scope: exact-target metadata resolution and safe target-list snapshots.
- Status: distinct specification and quality reviews passed after one regression-test addition.

## Current Artifact After Specification Review

The [specification review](host-metadata-adapter-spec-review.md) added `TestEvalTargetAdapterNilAPIKeyDependency`.
Production code did not change. The current patch 004 SHA-256 is
`850af726f80bb2f379e0cf9b3f395f28e9d8e89356106f0643c39ee5558f4088`.
The checks and hash below describe initial delivery before that test addition.
The [distinct quality review](host-metadata-adapter-quality-review.md) passed fresh pristine replay and final checks of the corrected artifact.
Its replay tree is `dist/sub2api-build/scratch-metadata-quality-vSsAGFLZ`.
All 3256 backend, 931 frontend, and 28 core files matched the maintained workspace.

## Files Changed

- `dist/sub2api-build/workspace/backend/internal/service/eval_target_adapter.go`: new 197-line service adapter.
- `dist/sub2api-build/workspace/backend/internal/service/eval_target_adapter_test.go`: new 380-line test file.
- `patches/sub2api/v0.2.15/004-codex-eval-target-adapter.patch`: maintained two-file patch.
- This delivery report.

Stage 4A/4B sources, frozen patches 001-003, and upstream dependency lockfiles remain unchanged.
No handler, route, Wire, lifecycle, probe, bootstrap, UI, installer, manifest, or build helper was edited.

## Implemented Contracts

`NewEvalTargetAdapter(AdminService, GroupRepository, *APIKeyService)` follows native constructor conventions.
It assigns existing exported dependencies to private narrow reader interfaces.
A compile-time constructor signature assertion covers later external construction without private Wire interface bindings.
The adapter implements `codexeval.TargetResolver` directly.

`Resolve` rejects nonpositive IDs before any reads.
It resolves the exact directed group, then reloads the requested account without alternate-target reads.
It checks exact IDs, active OpenAI group/account status, supported account families, and hydrated membership.
Success returns the requested target and account name only.
Stored schedulability and expiry do not duplicate gateway admission policy.

`ListTargets` reads active exact-OpenAI groups and preserves group/account relations through one-group membership queries.
Positive account IDs are deduplicated before one batch hydration.
Hydrated platform, type, and membership checks exclude missing, rebound, unsupported, and unsolicited rows.
Groups and accounts sort by numeric ID, and empty collections encode as JSON arrays.
Inactive member accounts remain visible with their stored schedulability.

The DTO projection contains only IDs, names, native type, allowlisted status, stored schedulability, and nested account arrays.
Unknown account statuses become `unavailable`.
Credentials, Extra, proxy metadata, notes, model configuration, and raw backend errors never enter DTOs.
Read failures return a fixed safe metadata error, while cancellation and deadline errors preserve standard context identity.
Caller context is passed directly to every dependency and checked between reads.

These snapshots do not authorize later starts or activate evaluations.
Normal gateway checks remain authoritative for permissions, subscriptions, models, quota, concurrency, and billing.

## Test Evidence

Fixtures were written before implementation.
The initial test invocation failed compilation because the adapter and constructor did not exist.
After implementation, the first JSON fixture failed because it used prose spellings `api_key` and `setup_token`.
The native constants are `apikey` and `setup-token` in `internal/domain/constants.go:65-67`.
Only fixture expectations changed for that mismatch, and production code retained native constants and helpers.
That fixture mistake was not a pre-existing product defect.

Nine test functions cover constructor dependencies, all three supported native account types, exact lookup order,
invalid identities, group/account mismatch, status/platform/type rejection, missing rows, nonmembership, and membership drift.
List fixtures cover per-group pairs, one deduplicated batch, stable sorting, missing/rebound rows, explicit JSON shape,
all status projections, private sentinels, empty arrays, and independent snapshots.
An account bound to both groups appears only under its queried relations.
After a displayed account loses membership, a new `Resolve` rejects that target.
Error and cancellation fixtures cover every read stage, early cancellation, cancellation between reads, and absent dependencies.

Fresh final checks on the implementation workspace:

- `gofmt -l` for both adapter files: clean.
- `go test -race -count=1 ./internal/service -run '^TestEvalTargetAdapter'`: pass, 3.612s.
- `go test -tags=unit -race -count=1 ./internal/service -run '^TestEvalTargetAdapter'`: pass, 2.326s.
- `go vet ./internal/service`: pass.
- `go vet -tags=unit ./internal/service`: pass.
- Root `go vet ./...`: pass.
- Root `go test -race -count=1 ./...`: pass, evaluation core 10.779s and original plugin 16.645s.
- `scripts/build-sub2api.sh check`: pass, pinned source identity and toolchain verified.
- `python3 scripts/test_sub2api_build.py`: all 60 foundation tests pass.
- Build-input discovery: all four patches found in deterministic filename order.
- `git diff --check`: pass.

## Patch Replay and Parity

Initial delivery patch 004 SHA-256, superseded by the specification-review correction above:
`b656a4acf1676a38341dfe8a1dee6177cbafa6f894f3685421c9fc72a5884940`.

The export uses established bare unified headers with `a/backend/` and `b/backend/` paths.
An earlier `diff --git`/index export was silently skipped by Git in the nested ignored workspace despite exit status zero.
The full-tree comparison caught both missing adapter files.
The corrective export changed patch representation only and did not modify the Git helper.

The successful fresh replay workspace is:
`dist/sub2api-build/scratch-target-adapter-1791775754367192000`.
The failed-format scratch workspace remains separately preserved under the same scratch prefix.

The fresh tree came from read-only `git archive` of pinned commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.
Each patch passed sequential `git apply --check --whitespace=error` before actual strict application.
All 28 maintained evaluation-core files were safely staged into the replay before service compilation.
The existing implementation copy already matched those 28 files, so it needed no replacement or restaging.

Complete replay parity passed for all 3256 backend files and 931 frontend files.
This includes both new adapter files and every frozen-patch source owner.
All three declared lockfiles match pristine source and their recorded workspace baselines.
The pristine source identity and clean working tree were verified again after replay.

Frozen patch SHA-256 values remain:

- 001: `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`.
- 002: `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`.
- 003: `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.

Fresh replay checks passed after core staging:

- `go test -tags=unit -race -count=1 ./internal/service -run '^TestEvalTargetAdapter'`: pass, 3.685s.
- `go vet ./internal/service`: pass.
- `go vet -tags=unit ./internal/service`: pass.

## Remaining Gates

Distinct specification and quality reviews passed for the corrected artifact. Later host integration remains separately gated.
The broader host probe/runtime plan has separate refinements and review gates.
The full upstream service race suite remains baseline-blocked, and its fixture repair was declined.
Only focused service tests were run, so this report does not claim a full upstream service-suite pass.
Docker lifecycle, live topology, deployment, browser integration, live billing, and paid model calls were not run.
Docker is available after the latest resume. Container lifecycle checks remain pending, and no approved live Sub2API target exists.
No `.env` access, CPA deployment, commit, or push occurred.
