# Research: Sub2API Host Metadata Adapter

- Query: Define the earliest host metadata adapter chunk for the passed Stage 4B run manager, including exact-target resolution, safe target-list projection, focused tests, and patch replay.
- Scope: internal
- Date: 2026-10-12

## Findings

### Stage 4B boundary

Stage 4B implementation and distinct independent specification and quality reviews passed ([run-state spec review](eval-run-state-spec-review.md#L1-L16), [run-state spec](../../../spec/backend/run-state.md#L1-L8)). The manager consumes `TargetResolver.Resolve(context.Context, Target) (TargetMetadata, error)` and treats resolver errors or mismatched target identities as unavailable explicit targets, without substituting another account ([runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L14-L23), [runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L234-L250), [run-state spec](../../../spec/backend/run-state.md#L38-L43)). Core `AccountMetadata` remains display-only and currently contains only `Name`; keep host-specific status and ownership fields out of it ([eval.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/eval.go#L13-L30)).

### First host-owned files

Limit this chunk to two new files in the pinned workspace:

- `backend/internal/service/eval_target_adapter.go`: `EvalTargetAdapter`, narrow dependency interfaces, `Resolve`, and `ListTargets` with service-owned target-list DTOs.
- `backend/internal/service/eval_target_adapter_test.go`: unit fixtures for exact resolution, list membership, status projection, and JSON redaction.

Do not edit Stage 4A/4B core files, frozen patches 001-003, handler/routes, probe/bootstrap, Wire, lifecycle, UI, installer, compatibility manifest, or build scripts. The patch collector sorts every `*.patch` in the required patch directory, so after implementation export the two new files as `patches/sub2api/v0.2.15/004-codex-eval-target-adapter.patch`; no build-script or manifest edit is indicated by the current workflow ([sub2api_build.py](../../../../scripts/sub2api_build.py#L214-L220)).

### Adapter contract and dependencies

Implement the existing core interface directly:

```go
func (a *EvalTargetAdapter) Resolve(
    ctx context.Context,
    target codexeval.Target,
) (codexeval.TargetMetadata, error)

func (a *EvalTargetAdapter) ListTargets(
    ctx context.Context,
) ([]EvalTargetGroupDTO, error)
```

Use small local interfaces so tests do not need the full admin service or repositories:

```go
type evalTargetAdminReader interface {
    GetAllGroupsByPlatform(context.Context, string) ([]service.Group, error)
    GetAccount(context.Context, int64) (*service.Account, error)
    GetAccountsByIDs(context.Context, []int64) ([]*service.Account, error)
}

type evalTargetMembershipReader interface {
    GetAccountIDsByGroupIDs(context.Context, []int64) ([]int64, error)
}

type evalDirectedGroupResolver interface {
    ResolveDirectedTargetGroup(context.Context, int64, int64) (*service.Group, error)
}
```

Keep these narrow interfaces private as adapter fields and test dependencies. The public production constructor takes existing exported `AdminService`, `GroupRepository`, and `*APIKeyService` dependencies, then assigns them to those fields. This avoids requiring the external `cmd/server` Wire graph to bind inaccessible private interface types. Store no credentials and defer Wire construction to the later integration chunk. `AdminService.GetGroup` already exists, but this adapter should not issue a redundant second group read: `ResolveDirectedTargetGroup` is the source-backed exact admission primitive for a selected target and itself reads the current group and membership. The list path starts from `GetAllGroupsByPlatform` instead ([admin_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_service.go#L34-L42), [admin_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_service.go#L73-L83), [eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go#L701-L724)).

### Exact calls and projections

For `Resolve`, reject nonpositive IDs; call `APIKeyService.ResolveDirectedTargetGroup(ctx, groupID, accountID)` and require the returned group ID to match, remain active, and be exactly `service.PlatformOpenAI`. Then call `AdminService.GetAccount(ctx, accountID)`. Require a nonnil result with the requested positive ID, exact OpenAI platform, active status, and one supported OpenAI Responses account type: `oauth`, `setup-token`, or `apikey`. Use native `AccountTypeOAuth`, `AccountTypeSetupToken`, and `AccountTypeAPIKey` constants. Their stored values are `oauth`, `setup-token`, and `apikey`. The existing account helpers define OAuth-like Codex types and API-key types, and the OpenAI scheduler recognizes the same two families ([account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/account.go#L1326-L1380), [openai_account_scheduler.go](../../../../dist/sub2api-build/workspace/backend/internal/service/openai_account_scheduler.go#L2993-L3005)). Cross-check that the reloaded account's hydrated `GroupIDs` still includes the requested group; a mismatch or lookup race returns an unavailable error. Return exactly `codexeval.TargetMetadata{Target: target, Account: codexeval.AccountMetadata{Name: account.Name}}`.

For `ListTargets`:

1. Call `AdminService.GetAllGroupsByPlatform(ctx, service.PlatformOpenAI)`. Keep only positive-ID groups whose platform is exactly OpenAI and status is active.
2. For each group, call `GroupRepository.GetAccountIDsByGroupIDs(ctx, []int64{group.ID})` with one group ID. That repository method returns a distinct union across all supplied groups; per-group calls preserve the group/account relation without adding a new repository query ([group_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/group_service.go#L17-L35), [group_repo.go](../../../../dist/sub2api-build/workspace/backend/internal/repository/group_repo.go#L1044-L1072)).
3. Deduplicate positive account IDs across groups, then call `AdminService.GetAccountsByIDs(ctx, ids)` once. Index returned accounts by ID. This existing batch method hydrates `GroupIDs`; the repository explicitly loads account-group bindings ([admin_account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_account.go#L55-L69), [account_repo.go](../../../../dist/sub2api-build/workspace/backend/internal/repository/account_repo.go#L281-L348)).
4. Build each group's account list only from IDs returned for that group, and require each hydrated account to have a matching ID, exact OpenAI platform, supported type, and the group ID in `account.GroupIDs`. Ignore missing rows caused by concurrent deletion/rebinding. Preserve active groups and supported-type member accounts even when an account is disabled or not schedulable so the list can show its stored state; `Resolve` rejects inactive accounts, and normal gateway scheduling handles current schedulability.

Use explicit DTOs rather than serializing `service.Group` or `service.Account`:

```go
type EvalTargetGroupDTO struct {
    ID       int64                  `json:"id"`
    Name     string                 `json:"name"`
    Status   string                 `json:"status"`
    Accounts []EvalTargetAccountDTO `json:"accounts"`
}

type EvalTargetAccountDTO struct {
    ID         int64  `json:"id"`
    Name       string `json:"name"`
    Type       string `json:"type"`
    Status     string `json:"status"`
    Schedulable bool   `json:"schedulable"`
}
```

Map group status to the fixed `active` value after checking it. Map account status through a small allowlist (`active`, `disabled`, `error`; otherwise `unavailable`) and emit the stored schedulable bool. Sort groups and accounts by ID for stable snapshots. Do not include `Credentials`, `Extra`, `Proxy` or proxy IDs/names, raw error messages, model configuration, token data, full service objects, or fields sourced from arbitrary maps. `GetAccountsByIDs` hydrates sensitive backend fields, so projection must be explicit even though the adapter reads only approved fields ([account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/account.go#L24-L68)).

The list is display data only. The group/account membership can change after it is returned; `RunManager.Start` calls `Resolve` again, and the subsequent normal gateway request rechecks directed target membership and continues through the existing API-key, group permission, subscription, balance, quota, model capability, scheduler eligibility, concurrency, and billing checks. The list's `schedulable` field is stored display state, not an eligibility guarantee: `Account.IsSchedulable()` also depends on live expiry, cooldown, overload, and quota state, and model admission is request-specific. Do not duplicate those request-time policies in the metadata adapter ([runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go#L234-L250), [account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/account.go#L182-L200), [api_key_auth_eval.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth_eval.go#L65-L89), [api_key_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth.go#L99-L187)).

### Focused test and patch gate

Use local function-backed stubs for the three narrow interfaces. No database, Redis, Wire graph, HTTP handler, network call, or live evaluation is needed. Cover:

- Resolver success returns exact positive IDs and name-only core metadata; mismatched group/account IDs, non-OpenAI/inactive group, nonmember, missing account, non-OpenAI account, unsupported type, inactive status, and membership drift all fail closed. A stored `schedulable` value does not override the gateway's live scheduler admission.
- Target-list request passes only `PlatformOpenAI`; per-group membership calls retain exact pairs; IDs are deduplicated before one batch hydration; nonpositive IDs, concurrent missing accounts, stale group binding, foreign platform/type, and non-active groups do not enter the DTO.
- DTO JSON includes only the declared fields and statuses. Seed credential, token, `Extra`, proxy, and raw backend error sentinels in service stubs and assert none occur in serialized output.
- Existing Stage 4B tests already prove that resolver errors become unavailable records for the originally requested target and do not start workers; retain those tests unchanged ([runs_test.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs_test.go#L880-L915)).

After implementation, run the focused service tests and touched-package vet. Export the change only as patch 004, apply/check patches 001-004 in order to a fresh pristine pinned `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` workspace, then compare the replayed adapter and test files byte-for-byte with the staged workspace. Confirm patches 001-003 remain byte-identical. Do not edit their source owners as part of this adapter chunk.

## Files Found

- [implement.md](../implement.md): passed Stage 4B spec-review checkpoint and ordered host-metadata chunk gate.
- [eval-run-state-spec-review.md](eval-run-state-spec-review.md): Stage 4B specification verdict and remaining quality/deployment caveats.
- [eval-run-state-plan.md](eval-run-state-plan.md): manager and state contracts.
- [host-runtime-wiring-plan.md](host-runtime-wiring-plan.md): earlier metadata boundary and host seam findings.
- [runs.go](../../../../dist/sub2api-build/workspace/backend/internal/codexeval/runs.go): `TargetResolver`, `TargetMetadata`, and resolver admission behavior.
- [admin_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_service.go), [admin_group.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_group.go), [admin_account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/admin_account.go): existing group/account list and hydration methods.
- [eval_key_runtime.go](../../../../dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go): exact active OpenAI group membership resolver.
- [group_service.go](../../../../dist/sub2api-build/workspace/backend/internal/service/group_service.go), [group_repo.go](../../../../dist/sub2api-build/workspace/backend/internal/repository/group_repo.go): membership interface and union semantics.
- [account.go](../../../../dist/sub2api-build/workspace/backend/internal/service/account.go), [account_repo.go](../../../../dist/sub2api-build/workspace/backend/internal/repository/account_repo.go): sensitive account shape, supported OpenAI type helpers, and batch membership hydration.
- [api_key_auth_eval.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth_eval.go), [api_key_auth.go](../../../../dist/sub2api-build/workspace/backend/internal/server/middleware/api_key_auth.go): runtime target membership and subsequent normal request policy.
- [sub2api_build.py](../../../../scripts/sub2api_build.py): patch directory is collected in sorted filename order.

## Related Specs

- [evaluation-core.md](../../../spec/backend/evaluation-core.md): numeric target identity and name-only core metadata.
- [run-state.md](../../../spec/backend/run-state.md): resolver revalidation, unavailable explicit targets, and no substitution.
- [error-handling.md](../../../spec/backend/error-handling.md): safe surfaced errors.
- [quality-guidelines.md](../../../spec/backend/quality-guidelines.md): focused checks and error visibility.

## External References

- None. Findings are from the pinned local Sub2API workspace and task specifications.

## Caveats / Not Found

- This plan covers only the metadata adapter and target-list projection. Its interfaces are source-backed proposals, not implemented or Wire-constructed types.
- Target-list and resolver reads are snapshots across separate service calls. They cannot authorize a later start or prove the live account remains schedulable; normal gateway admission remains authoritative.
- Both distinct Stage 4B local reviews passed. Later host integration and deployment acceptance remain separate gates.
- No source supports widening evaluation to non-OpenAI platforms in this chunk. No code, tests, patch export/replay, deployment, `.env` access, or live/model request was performed for this research.
