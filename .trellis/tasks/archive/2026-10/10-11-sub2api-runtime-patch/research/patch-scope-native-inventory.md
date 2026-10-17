# Research: Current Native Adapter Scope Inventory

- Query: Inventory the current structural rules, UI rules, runtime templates and adapter build path against pristine Sub2API v0.2.16; identify independent behaviors, upstream anchors, provenance, dependencies, deletion effects and pending retain/remove recommendations. Assess the maintenance startup and shutdown/drain changes explicitly.
- Scope: internal
- Date: 2026-10-14

## Findings

The audit is anchored to the verified pristine Sub2API tag `v0.2.16`, revision `01842576b3d5c1e2545b5a99ad555d36e28ed306`, archive SHA-256 `eef7e0519b05b8254c981efbab7bfe74c5d89897c67053e81380987366daae89`. The extracted tree is recorded in the machine-readable inventory. No pristine source was changed.

The inventory covers all 73 entries in `adapter/rules.json`, its six declaration groups containing 13 function definitions, all 15 injected import files, all 8 UI edits, all 16 runtime templates, `adapter/build.py`, and the release payload allowlist. Each review unit has an `N01`-`N19` ID and remains `pending` in [patch-scope-native-inventory.json](patch-scope-native-inventory.json). Rule ordinals in the prose below are one-based; the JSON uses exact zero-based JSON pointers. Host source paths in the JSON are relative to its recorded pristine-tree base.

| ID | Behavior | Recommendation | Main dependency or effect |
| --- | --- | --- | --- |
| N01 | Same-binary maintenance capability/client modes | Retain conditionally | Installer status and cleanup currently invoke these modes. |
| N02 | Maintenance goroutine plus HTTP listener rewrite | Split; remove listener rewrite | The in-process transport binds the Gin handler directly; the maintenance UDS does not use the TCP listener. |
| N03 | Evaluation admission stop and manager/producer drain | Retain conditionally; review timeout loop | Native server shutdown does not drain the evaluation manager; repeated new timeouts currently have no total bound. |
| N04 | Application, handler and service Wire graph | Retain only required edges; consider native Wire regeneration | Generated `wire_gen.go` rewrites are build glue, not separate product behavior. |
| N05 | Real-router binding, server wrapper and external-header guard | Retain conditionally | Needed for native handler/middleware execution and external capability isolation. |
| N06 | Administrator evaluation API and static shell route | Retain conditionally | Enables status, activation, targeting, runs, history and deactivation in-panel. |
| N07 | Process-bound private capability and public credential rejection | Retain conditionally | No native context capability exists; this is the isolation boundary. |
| N08 | Fixed-expiry native accounting Key and conditional lifecycle | Retain conditionally; review migration/helper pieces separately | Native usage rows require the API Key foreign key and ID. |
| N09 | Redact the owned Key lookup in native DTO output | Retain while DTOs expose the row | The upstream mapper copies `APIKey.Key`. |
| N10 | Exact target resolution and safe metadata list | Retain conditionally | Native group/account readers exist, but not exact evaluator target semantics. |
| N11 | Restrict capability to OpenAI Responses and pin target through handler admission | Retain conditionally | Prevents endpoint crossover, target switching and partial-usage loss. |
| N12 | Exact scheduler selection and fallback/sticky exclusions | Retain conditionally | Native scheduler is intentionally multi-account and can select alternates. |
| N13 | Directed cancellation and partial usage behavior, including WebSocket reads | Retain conditionally; review subpaths separately | Native requests detach some streams so work can outlive client cancellation. |
| N14 | Protected Unix-socket status/deactivation | Retain conditionally; review informational response fields separately | Installer verifies live drain and exact revocation before restore. |
| N15 | Reserved-menu collision check without mutation | Retain conditionally with built-in menu | Prevents activation over a conflicting user menu and avoids adopting legacy entries. |
| N16 | Built-in sidebar/custom-page integration | Retain conditionally | There is no native static extension hook for this evaluator route. |
| N17 | Embedded SPA bypass and same-origin theme URL handling | Retain with the route, conditional on UI | Needed for the explicit shell handler to win over SPA fallback. |
| N18 | Protected deployment ID plus informational version/source/hash fields | Keep ID validation conditionally; source/hash output is optional | Current installer checks deployment ID/state/drained/revoked; it ignores source/hash response fields. |
| N19 | Source adaptation, runtime copy, typecheck and Linux build pipeline | Retain only while source-build install remains selected | All retained feature rules/templates are build inputs; allowlist must track removals. |

The lifecycle assessment is specific. In pristine [main.go](/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/pristine-sub2api-v0.2.16.ifVBya/extracted/sub2api-01842576b3d5c1e2545b5a99ad555d36e28ed306/backend/cmd/server/main.go:134), Sub2API starts `app.Server.ListenAndServe()` and later calls `app.Server.Shutdown(ctx)` at lines 171-194. The evaluation transport calls the Gin handler directly ([eval_http_transport.go.in](../../../../adapter/runtime/backend/internal/service/eval_http_transport.go.in:47)); `ProvideEvalHTTPServer` binds that handler into `EvalHostRuntime` ([adapter/rules.json](../../../../adapter/rules.json), declaration under `/Declarations/internal/server/http.go/0`). Neither path needs the explicit TCP listener added by rule 2. The maintenance service creates its own Unix socket. Rule 2 can therefore be split: keep its maintenance goroutine only if N01/N14 remain, and remove the `net.Listen`/`http.Server.Serve` replacement.

Drain has a different justification. Pristine Sub2API has no evaluation manager; its `Application.Cleanup` later starts parallel cleanup steps, while the evaluation manager's shutdown cancels runs and waits on its drain channel. Rule 3 provides the only explicit ordering before `Server.Shutdown` and dependency cleanup. It is technically distinct from the listener rewrite. Its current loop replaces an expired five-second context and retries without a total deadline, so retaining the drain behavior does not imply retaining that unbounded retry policy.

The maintenance CLI/socket is not redundant with native HTTP status. The current installer probes `--codex-eval-maintenance-capabilities`, invokes same-binary status/deactivate, and requires `drained` plus `owned_key_revoked` before rollback/uninstall ([install.py](../../../../install.py:813)). The native server has no such mode. This makes N01/N14 conditional on retaining those installer workflows, while leaving each pending for the user's individual decision.

Requirement provenance is deliberately limited. The active task requires a complete inventory and individual approval in `prd.md` B10-B12; it explicitly says earlier specs/implementation are not current approval. Earlier product PRDs/designs explain why native accounting, exact target selection, UI, maintenance and cancellation were introduced, but they do not settle any `N` item. The JSON records both provenance tiers and sets all items to `pending`.

### Files Found

- `adapter/rules.json` — 73 AST rules, six declaration groups (13 definitions), and 15 import additions used by the source adapter.
- `adapter/ui-rules.json` — eight guarded frontend edits across sidebar, custom page and embedded URL helper.
- `adapter/runtime/` — 15 backend Go templates plus one frontend menu helper copied into the selected host source.
- `adapter/build.py` — validates source identity/locks, applies structural and UI edits, copies runtime/core, builds frontend, type-checks native contracts and compiles the Linux server.
- `scripts/universal_payload.py` — explicit allowlist includes all current adapter rules, build logic and runtime templates.
- pristine Sub2API `backend/cmd/server/main.go`, `backend/internal/server/http.go`, `backend/internal/server/router.go`, `backend/internal/server/routes/admin.go`, auth middleware, API-key service/repository, gateway handlers/services, generated Wire graph and frontend sources — host integration points cited by item.
- `.trellis/spec/backend/host-target-adapter.md`, `patch-native-installation.md`, `runtime-installation.md`, `evaluation-core.md` — related existing contracts; not current retention approvals.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md` and `design.md` — current audit scope and individual-approval boundary.

### Code Patterns

- `backend/cmd/server/main.go:55-94` parses native flags before entering normal server startup; lines 171-194 use `ListenAndServe` then `Shutdown`.
- `backend/internal/server/http.go:29-45,112-163` constructs the Gin router and native `http.Server`; the custom provider wraps these existing pieces.
- `adapter/runtime/backend/internal/service/eval_host_runtime.go.in:83-89` binds a router handler; `eval_http_transport.go.in:47-106` invokes it directly, with a pipe-backed response body and producer tracking.
- `backend/cmd/server/wire.go:27-64` is the declarative Wire application graph; `wire_gen.go:35-365` is generated graph output patched by many structural rules.
- `backend/internal/handler/dto/mappers.go:80-125` copies the API Key into DTOs; rule 29 clears only the owned row's lookup.
- `backend/internal/service/api_key_service.go:494-565` is native Key creation; adapter declarations add a fixed-expiry private path and exact repository expiry/revoke predicates.
- `backend/internal/server/middleware/api_key_auth.go:35-110` is the native credential ingress extended for context capability; Google auth and exact credential markers are separately guarded.
- `backend/internal/handler/openai_gateway_handler.go:389-420,640-790,850-1005` contains the Responses parse, slot/admission, retry and usage flow modified for directed requests.
- `backend/internal/service/openai_account_scheduler.go:2163-2335` and `openai_gateway_scheduling.go:905-930,1138-1165` cover native scheduler and alternate selectors guarded for directed requests.
- `backend/internal/service/openai_ws_forwarder_v2.go:390-415,526-605` detaches ordinary upstream reads; directed evaluation changes keep cancellation tied to the manager.
- `adapter/runtime/backend/internal/service/eval_maintenance.go.in:46-206` validates bounded maintenance requests, deployment identity and peer access; the Linux helper uses `SO_PEERCRED`.
- `install.py:813-834,1039-1070` verifies maintenance status and deactivation before transaction restoration.
- `adapter/build.py:14-51,53-104` is the current host adaptation/build path; `scripts/universal_payload.py:5` lists functional build inputs.

### External References

None. Findings use the exact tagged source and current workspace files.

### Related Specs

The host adapter, native installation, runtime installation, evaluation core, current task PRD/design, and pristine upstream contract report are listed in the JSON inventory. Existing spec statements remain technical provenance, not user approval.

## Caveats / Not Found

- No product changes or approval decisions were made. Every item remains `pending`.
- This is an inventory, not a claim that the current implementation works on the user's deployment. No build or test was run.
- The audit covers the current adapter build inputs requested here; it does not inventory the evaluation core, installer, package contents, historical patch files, or retired production patches as independent product items.
- No `adapter/build.py` regeneration, upstream source mutation, deployment, paid request, or artifact cleanup was performed. One read-only `git status` query was inadvertently run while checking the output scope; no Git mutation occurred.
- The active task's no-batch rule applies: recommendations must be presented to the user one concrete item at a time before any product change.
