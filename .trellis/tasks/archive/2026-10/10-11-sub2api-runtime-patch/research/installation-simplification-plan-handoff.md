# Research: Installation Simplification and Runtime Wiring Handoff

- Query: Finish the one-command installation and first-activation design, settle
  in-process Sub2API dispatch, protected live-runtime cleanup, and legacy native
  accounting migration.
- Scope: mixed (pinned Sub2API source, existing patch contracts, product plan)
- Date: 2026-10-12

## Findings

### In-process dispatch

The audited server already exposes the composition points needed for a genuine
in-process request. `server/http.go:30-92` returns the configured `*gin.Engine`;
`http.go:114-174` applies the normal body limit and then wraps the final handler
with `EvalHostRuntime.Guard`. The real listener is retained before serving and
the evaluation manager drains before HTTP shutdown (`cmd/server/main.go:172-213`).
The core request and stream limits, cancellation, and usage parsing remain in
`internal/eval/http.go:27-141`, `model.go:20-67`, and `sse.go:52-216`.

Proposed ownership and signatures:

```go
// internal/service/eval_host_runtime.go
func (r *EvalHostRuntime) BindRouter(router http.Handler)
type evalServeHTTPTransport struct{ handler http.Handler }
func newEvalServeHTTPTransport(router http.Handler) http.RoundTripper
func (t *evalServeHTTPTransport) RoundTrip(*http.Request) (*http.Response, error)
```

Call `BindRouter` from `ProvideEvalHTTPServer` after `ProvideHTTPServer` has
received the router and before the outer `Guard` is installed. The transport
invokes that retained router directly, so it cannot recurse through the outer
guard or the TCP listener. Its request is exactly `POST /v1/responses`, with a
bounded body, no authorization or target headers, no query/fragment/redirect,
and a context capability carrying the target and runtime identity. Its response
writer must stream through `io.Pipe`, publish status/headers on first write,
implement `http.Flusher`, and close the pipe when the request context is
cancelled. The existing SSE reader continues to enforce event/output limits.

`openai_gateway_handler.go:643-850` already performs directed admission,
post-slot recheck, forwarding, and usage callbacks. The generic Responses
handler deliberately rejects directed context (`gateway_handler_responses.go:184-187`),
so the transport must select the existing OpenAI route only.

### Capability and native admission

Keep the capability type in the internal routing package beside the existing
target context (`internal/pkg/evalroute/evalroute.go:15-101`), with unexported
fields and an opaque context key. The runtime is the only producer in product
code; middleware is the only consumer. The capability contains the exact API
key row ID, owner ID, target, fixed expected expiry, and a process-opaque
runtime nonce. It has no HTTP serialization.

At the start of `APIKeyAuthMiddleware` (`api_key_auth.go:35-90`), check the
capability before public credential extraction. The private branch loads the
exact row by ID, verifies owner, `StatusActive`, nil external group binding,
IP policy, lookup identity, and the expected fixed expiry. It skips only the
ordinary external credential-expiry check; native group/account admission,
subscriptions, quotas, billing, and concurrency remain in force. Copy the key
and resolved group per request, then install the existing directed context used
by `api_key_auth_eval.go:65-90`. Headers and URLs cannot mint this capability;
tests must cover forged headers, public-context isolation, and route recursion.

The private transport must set `Request.RemoteAddr` explicitly to the
server-owned loopback source (`127.0.0.1:0`; use the canonical `::1` form only
when an exact legacy policy requires it) and send no forwarded-IP headers.
Native resolution then continues through `ip.GetSecurityClientIP`/Gin's
`ClientIP` fallback (`pkg/ip/ip.go:53-86,153-171`) and the existing API-key ACL
check (`api_key_auth.go:139-153`). A fresh accounting row uses canonical
`127.0.0.1` in its immutable whitelist. An old row with any other whitelist is
a migration conflict; there is no IP bypass or second private-path exemption.

### First-use administrator entry

The current UI has a real first-use gap. `AppSidebar.vue:35-100,775-859`
constructs built-in admin items separately from `adminSettingsStore.customMenuItems`;
the latter is loaded from the settings API (`stores/adminSettings.ts:54-79`).
`CustomPageView.vue:222-260` resolves its iframe target only from public/admin
custom-menu settings, while `/custom/:id` is the existing authenticated route
(`router/index.ts:389-397`). The backend page handler also treats configured
custom-menu visibility as authoritative (`page_handler.go:204-247`).

Use a narrow static admin-only sidebar item with a reserved ID/path such as
`/custom/codex-eval`, and add a synthetic reserved `CustomMenuItem` fallback in
`CustomPageView` that points to same-origin `/codex-eval/ui` only for an admin.
This reuses the existing iframe, theme reload, and exact URL credential
exclusion in `frontend/src/utils/embedded-url.ts`; it does not require a DB
menu mutation or manual URL entry. A later authenticated activation may merge
the same owned menu item for persistence. Migration may remove only an exact
owned prior entry; unrelated menu entries and conflicting edits remain intact
and conflicts are reported. Suppress a duplicate configured item by reserved
ID so the sidebar never shows two entries.

### Installer and release journey

The current installer requires explicit container, origin, administrator token,
and ingress revision (`sub2api_install.py:47-94,326-328`) and performs menu/API
mutation during installation (`sub2api_install.py:262-280,413-433`). Replace
the public shell/PowerShell entry with tagged release download, checksum
verification, safe extraction, and invocation of the existing Python
transaction. Container selection is always explicit: show an interactive list
with name, immutable ID, image, status, and compatibility, and require the
operator to choose even when exactly one candidate exists. Do not treat Enter
or a blank response as a default. In unattended mode, `--container` is the
explicit choice; without a controlling terminal or that flag, fail before any
mutation. Read interactive selection from the controlling terminal rather than
bootstrap script stdin (so `curl | sh` cannot answer implicitly), re-inspect
the selected immutable ID immediately before mutation, and never switch targets
automatically. Remove JWT-file, origin, and ingress-revision inputs from the
default install command. Installation verifies image/source/platform/binary
identity, stages atomically, restarts only the selected container, and reports
`installed_uninitialized` until browser setup completes.

The administrator logs into the native panel, opens the built-in entry, and
activates through the existing admin middleware. The revised production
contract has no mandatory proxy attestation, external runner, topology
inventory, or topology-readiness gate. Business fingerprint probe data remains
runtime-owned evaluation data. The in-process dispatch removes the operator
probe executable and topology artifacts from the default server payload; keep
test overlays, fixtures, and validation scripts in the separate test workspace.

### Protected live-runtime cleanup

Add a runtime-owned Unix socket under the evaluation data root. Suggested owner
and API:

```go
// internal/service/eval_maintenance.go
func (r *EvalHostRuntime) ServeMaintenance(ctx context.Context) error
// cmd/server/maintenance_client.go (same server binary, before Wire startup)
func runMaintenanceClient(ctx context.Context, socketPath string) error
```

The daemon creates a no-symlink, private socket and serves one bounded
maintenance request. `cmd/server` dispatches the client mode before
`initializeApplication`/Wire/server startup; it does not launch a second
server. The daemon validates Unix peer credentials and the protected identity
journal, then calls the existing `EvalHostRuntime.Deactivate`, which already
drains the manager, revokes the exact owned row, and removes protected state
(`eval_host_runtime.go:243-283`). The same server binary is transiently
launched through `docker exec` as the runtime UID, sends no secret in argv,
environment, URL, or SQL, and owns no runtime state. The UDS is for maintenance
only; the live patched runtime owns drain and revocation. If the socket is
absent because the container is stopped or still runs the original binary,
cleanup is incomplete: retain the journal and backup, refuse restore-as-success,
and require reapplication of the supported patch before retrying.

### Legacy native-row migration

`EvalProvisionIntent` is currently version 1 and has no expiry field
(`001-codex-eval-auth.patch:634-680`; `evalIntentVersion = 1`). The audited
legacy create path leaves expiry nil when `ExpiresInDays` is nil, and the
repository persists that nil with `SetNillableExpiresAt`
(`api_key_service.go:213-224,584-588`; `api_key_repo.go:45-58`). Bump the
intent schema to version 2 and add an accounting-identity version plus the
fixed past `ExpiresAt`. A planned migration must first load the exact row by ID
and verify lookup, owner, active status, nil group binding, IP policy, and
`ExpiresAt == nil`; any legacy expiry value is a conflict.

The DB update and intent file cannot be one atomic transaction. Persist a
versioned `migration_pending` intent before the effect, recording the row ID,
owner/lookup tuple, old nil expiry, and fixed target expiry. Update the row's
expiry exactly once through the native repository, invalidate native API-key
caches, and finalize the v2 intent only after rereading the row. After an
interruption, recovery permits only the explicitly recorded old/nil or new
fixed expiry: nil retries the planned update, the fixed value finalizes the
intent, and any other value is a repair conflict. Preserve row ID,
counters/history, and restart the local connection if the ORM/cache requires
it. Never create a replacement or silently repair an external edit. Existing
WebSocket sessions may outlive later expiry or status changes and remain an
accepted residual.

### Legacy menu cleanup boundary

Historical rollback ownership in `sub2api_install.py:262-280` matched the
complete generated menu entry, including the absolute `origin +
"/codex-eval/ui"` URL, before deleting it. The protected legacy journal kept
only `menu_owned=true`, so that flag is necessary but insufficient for safe
cleanup. A legacy entry may be deleted only when its full non-origin tuple
matches and its absolute URL is provably the current same-origin evaluation
URL. The server-local maintenance path has no browser origin: an
unverifiable or foreign origin, changed URL, duplicate ID, or any other
conflict must leave the entry untouched. Report cleanup as incomplete and
retain the journal and backup; do not claim rollback or uninstall success.
Do not add an origin prompt or other user-input requirement, and leave
unrelated menu entries untouched.

## Caveats / Not Found

- The parent host/console repositories are not checked out beside this
  workspace; conclusions use the audited v0.2.15 copy under `dist/` and the
  pinned patch series.
- No live target, browser session, build, test, deployment, or billing path was
  exercised. Required implementation checks are transport cancellation/flushing,
  no recursive admin dispatch, reserved-entry admin isolation/duplicate
  suppression, UDS peer/journal checks, cleanup-incomplete recovery, and
  v1-to-v2 row-preserving migration.
- No further user-owned decision is currently required; the remaining choices
  are implementation details covered by the contracts above.

## Related specs

- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/deployment-evidence.md`
- `.trellis/spec/backend/host-coordinator.md`
- `research/host-runtime-wiring-plan.md`
- `research/installation-simplification-options.md`
- `research/patch-native-accounting-identity.md`
