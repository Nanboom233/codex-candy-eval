# Sub2API Runtime and Installation

The [universal Linux delivery](./universal-linux-delivery.md) defines actual host discovery, source builds and release payload.
Preserve native runtime, identity, maintenance and transaction behavior.

## 1. Scope / Trigger

Read this contract before changing administrator APIs, readiness, UI embedding, maintenance, shutdown or binary transactions.
The [patch-native contract](./patch-native-installation.md) defines the complete active architecture.
Historical topology and operator contracts remain in [deployment evidence](./deployment-evidence.md) and [coordinator](./host-coordinator.md).
Their proof sessions and network probes are retired from production.

Runtime owners are [runtime additions](../../../adapter/runtime/backend/internal/service/eval_host_runtime.go.in) and [runtime ports](../../../adapter/runtime/backend/internal/service/eval_runtime_ports.go.in).
Accounting and conditional repository wiring use [structural rules](../../../adapter/rules.json) and [Key runtime](../../../adapter/runtime/backend/internal/service/eval_key_runtime.go.in).
Native UI owners are [UI rules](../../../adapter/ui-rules.json) and [menu integration](../../../adapter/runtime/frontend/src/utils/codex-eval-menu.ts).
Embedded resources remain under [internal/eval](../../../internal/eval/ui.go).
The installer is [install.py](../../../install.py).

The [native runtime review](../../tasks/10-11-sub2api-runtime-patch/research/final-check-review-20261013.md) records unchanged runtime evidence and unresolved gates.
The minimal Python review records historical fixed-package checks.
See [universal implementation](../../tasks/10-11-sub2api-runtime-patch/research/universal-implementation.md) for current evidence.
The [independent universal review](../../tasks/10-11-sub2api-runtime-patch/research/universal-check-review.md) is complete.
Required broad gates and external acceptance remain unresolved.

## 2. Signatures

```go
func NewEvalHostRuntime(cfg *config.Config, keys *APIKeyService, admin AdminService, groups GroupRepository, settings SettingRepository, build BuildInfo) *EvalHostRuntime
func (r *EvalHostRuntime) BindRouter(router http.Handler)
func (r *EvalHostRuntime) Ready(ctx context.Context) error
func (r *EvalHostRuntime) State(ctx context.Context) string
func (r *EvalHostRuntime) Activate(ctx context.Context, actor int64) error
func (r *EvalHostRuntime) Shutdown(ctx context.Context) error
func (r *EvalHostRuntime) Deactivate(ctx context.Context, actor int64) error
func (r *EvalHostRuntime) ServeMaintenance(ctx context.Context) error
func EvalMaintenanceClient(ctx context.Context, request EvalMaintenanceRequest) (EvalMaintenanceResult, error)
```

```sh
python3 install.py
python3 install.py install --container <explicit-container> --bundle <release-directory>
python3 install.py inspect --container <explicit-container>
python3 install.py status --container <explicit-container>
python3 install.py rollback --container <explicit-container>
python3 install.py uninstall --container <explicit-container>
```

Install is the default operation. Only install downloads when no local bundle exists.
Maintenance uses embedded facts and protected journals without a bundle or download.
Without `--container`, require an explicit selection through the controlling terminal.
Never select a sole candidate automatically. Empty input does not select a container.
No administrator-token file, origin, inventory or ingress revision is required.

## 3. Contracts

### Native administrator and internal gateway

Register operations under `/api/v1/admin/codex-eval/` with native administrator middleware.
Use native `{code,message,data}` envelopes. Derive actors from authentication.
Activation accepts no caller-supplied owner, identity or installation facts.
Keep state, targets, start, cancel, records, clear-history, flush and deactivate under native authentication.

Add the static administrator entry `/custom/codex-eval` before activation.
Serve `/codex-eval/ui` as a non-sensitive same-origin shell with nonce styles/scripts and route-local framing.
Only authenticated administrator APIs return evaluation data or perform operations.
Read the current administrator token only for API headers.
Preserve four tabs, native theme reload, preferences, history and the pelican `allow-scripts` sandbox.

Bind the real Gin router before installing the external guard.
Use the private capability transport only for canonical `POST /v1/responses`.
Internal requests carry no credential or target header.
Read the current exact owned row and preserve native authorization, quota, billing and shared slots.
Exempt only its recorded fixed credential expiry.
Preserve [run-state](./run-state.md), directed scheduling and cancellation contracts.
Invalidate admission, cancel producers and drain manager workers before closing dependencies.

### Protected state and recovery

Keep private state under `<actual-data-dir>/codex-eval/`.
Store installer-observed facts in `installer-facts.json` with schema version 3.
Require a valid deployment ID. Treat version/source/plugin stamps and legacy binary/image fields as optional provenance.
Legacy schema-1/2 installation facts are not automatically admitted by universal readiness.
Do not compare binary checksums, loaded executables or source/version/image stamps with protected facts.
Use strict JSON, private modes, no-symlink access and atomic synced writes.

Persist fixed past expiry before the first native accounting-row insert.
Version-2 identity intents record initialization, migration and revocation phases.
Recover only exact recorded identities. Missing or drifted rows remain conflicts.
Legacy migration preserves row ID, counters and history.
Conditional expiry/revocation uses native repositories and cache invalidation.
Only cleanup verification may inspect the exact soft-deleted revoked row.
Ordinary bootstrap and admission retain native soft-delete filtering.

### Installation and planned cleanup

Download the plugin-tagged archive through bounded HTTPS. Validate the complete functional tar entry set before extraction.
The same Python process owns download, extraction and installation.
Use the functional package beside a file-based installer and automatically build the matching source in Docker.
Offline builds require matching source, tools and dependencies locally. No fixed host server is bundled.
List candidates and require user selection even with one compatible candidate.
Reinspect the selected immutable container ID before effects.
Reject unresolved source, missing native capabilities, unsupported platforms or unsafe replacement paths before mutation.

Keep immutable `original-sub2api` and private `installation.json` under the state root.
Journal before staging and atomic replacement. Restart only the selected container.
Verify protected backup file type, permissions and owner. Verify startup health and functional maintenance readiness.
A running-container restore must restart even during candidate recovery.
A stopped-container restore must start the selected container.
Verify selected-container startup and local health before reporting original restoration success.

Use the protected Unix socket `<actual-data-dir>/codex-eval/control.sock` for maintenance only.
The same server binary runs a bounded client mode before Wire/server startup.
Validate peer UID, socket ownership, bounded protocol and expected deployment ownership.
Probe explicit maintenance flags. Original or unavailable runtimes must never start a duplicate server.
The live runtime owns drain and exact native Key revocation.
Do not launch a duplicate server or use SQL cleanup to claim drain.

Restore after verified cleanup success. Retain history, journal and original backup.
Unavailable or conflicted cleanup remains incomplete.
Legacy `menu_owned` alone cannot authorize deletion.
Activation may preserve an exact legacy tuple without adopting it.
Maintenance without same-origin ownership proof retains absolute menu entries and reports incomplete cleanup.

Print short process messages with the selected target, current step, actual result and next action.
Distinguish `installed_uninitialized` from `active`.
Installer failures use nonzero exit status and preserve complete raw command stdout/stderr and Python tracebacks. Do not append explanatory summaries or unrelated configuration dumps.

## 4. Validation & Error Matrix

| Condition | Outcome |
| --- | --- |
| Missing router, facts or admission | `not_ready`, no model execution |
| No owned intent | `installed_uninitialized`, ordinary traffic remains available |
| Different activation owner | `session_owner_conflict` |
| Deleted, disabled, rebound or expiry-edited row | `owned_key_repair_required`, no repair or replacement |
| Unsafe ancestry or corrupt state | Block evaluation and retain ordinary host availability |
| Public credential or forged evaluation header | Reject before forwarding or billing |
| No selection, cancellation or selected-ID drift | No container mutation |
| Unresolved source, unsupported platform or missing native capability | Refuse before replacement |
| Candidate failure | Restore and verify the original, retain failure journal |
| Missing live cleanup owner or menu conflict | Incomplete cleanup, retain journal and backup |
| Interrupted exact revocation | Resume `revoke_pending` using native soft-deleted-row verification |
| Missing or drifting binary/image metadata | Continue through functional discovery and readiness |

## 5. Good / Base / Bad Cases

Good: install, open the native administrator entry, activate, then use the real normally billed gateway.
Base: uninitialized installation leaves ordinary traffic available and evaluation unavailable.
Bad: restore the original without verified revocation and report complete cleanup.

## 6. Tests Required

Keep tests in [host overlays](../../../tests/host-overlay/) and [core overlays](../../../tests/overlays/codexeval/).
Run them only in fresh disposable test workspaces.
Assert atomic expiry, exact-row conflicts, row-preserving migration and native cache invalidation.
Test real PostgreSQL JSON argument binding and Ent soft-delete verification.
Join native admin activation, internal routing, all four evaluations, usage, slots and cancellation with a controlled upstream.
Run [Docker lifecycle checks](../../../tests/installer/docker_lifecycle.py) and [browser checks](../../../tests/ui/sub2api-browser.js).
Assert mandatory selection, original restore, recreation, durable cleanup retry, theme reload, nonce rendering and sandbox isolation.
Verify independent production/test parity, exact payload paths and source-identical installer/license bytes.
Assert missing/drifting binary metadata never gates installer, activation, status, maintenance or rollback.
Assert original/stale-socket probes start no server. Verify deployment ownership and drain/revocation results.
Test real terminal input with bootstrap code on stdin, EOF, cancellation and terminal-only fallback.
Keep full upstream race failures visible. Local fixtures do not prove real charges or deployed behavior.
Native amd64 lifecycle still needs fresh functional install, health, activation and cleanup evidence.

## 7. Wrong vs Correct

Wrong: accept network headers or loopback addresses as internal authorization.
Correct: require the private capability and current exact accounting identity.

Wrong: verify a soft-deleted row through an ordinary filtered query.
Correct: use `mixins.SkipSoftDelete` only for exact cleanup verification and retain normal admission filtering.

Wrong: pass a Go whitelist slice directly to the SQL driver.
Correct: serialize JSON before binding the conditional mutation argument.

Wrong: package tests or repair disabled Keys to make installation appear successful.
Correct: release only allowlisted product outputs and report exact ownership conflicts.
