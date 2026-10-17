# Code-Level Patch Contracts

## Status and Evidence

The initial sections preserve approved planning contracts. Later checkpoints record applied patches and their scoped review evidence.
Sections 14 and 15 identify the withdrawn native candidate and final Stage 2C scope.
No local review proves that a patched deployment works.

The user approved task activation and implementation. The source audit records checks against commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.

Source paths beginning with `backend/` or `frontend/` refer to that upstream commit. Existing local evaluation sources remain under `internal/plugin/`.

## 1. Patch Packaging

Keep maintained Go evaluation code in this repository. Copy it into `backend/internal/codexeval/` during a reproducible build.

Keep upstream changes as version-specific patches under proposed `patches/sub2api/v0.2.15/`. Keep upstream source in ignored `dist/`, not in Git.

Proposed local owners:

| Owner | Responsibility |
| --- | --- |
| `internal/eval/` | Four evaluation algorithms, worker state, SSE transport, history, and embedded page |
| `patches/sub2api/v0.2.15/` | Authentication, scheduler, cancellation, DTO, routes, wiring, and frontend source patches |
| `scripts/build-sub2api.sh` | Exact source validation, patch application, frontend build, Wire generation, static Linux builds |
| `scripts/test-sub2api.sh` | Focused and full patched-host checks with exact build tags |
| `scripts/test-runtime-patch.sh` | Disposable Docker lifecycle fixtures |
| `compatibility/sub2api-v0.2.15.json` | Source, index, platform manifest, image config, architecture, and artifact identities |
| `install.sh`, `install.ps1` | Inspect, install, status, rollback, and uninstall entry points |

These paths are proposed. Do not create product files during the planning phase.

Do not package a CPA shared library. Preserve data, prompts, algorithms, renderer licenses, and ModelTrace attribution when moving evaluation code.

Build with Go 1.27.2, `CGO_ENABLED=0`, `-trimpath`, and `-tags embed`. Build the pinned frontend before Go compilation.

Use the existing Wire dependency. Update provider declarations and regenerate `cmd/server/wire_gen.go` rather than hand-editing generated wiring.

Provide upstream license notices, exact patch sources, and reproducible build inputs with distributed binaries. Upstream uses LGPL-3.0.

## 2. Trusted Routing Context

Add a small dependency-free package at proposed `backend/internal/pkg/evalroute/`. It must not import handlers or services.

Proposed contracts:

```go
type Target struct {
    GroupID int64
    AccountID int64
    RunID string
}

func WithTarget(ctx context.Context, target Target) context.Context
func FromContext(ctx context.Context) (Target, bool)
```

Use a private typed context key. Raw HTTP headers are not trusted context.

The authentication middleware is the sole public ingress that creates this context after normal Key, IP, and active-user authentication.

Parse exactly one positive int64 group header and account header. Reject missing pairs, repeated values, zero, negative values, overflow, and unsupported paths.

Accept directed requests only on `POST /v1/responses`. Reject aliases, subpaths, WebSocket, images, other providers, and composite groups in v1.

Reject reserved headers from ordinary Keys. Reject a private evaluation token without its required target headers.

Delete reserved routing headers before forwarding. Never send the private token or routing control headers to an upstream provider.

### Authentication source owners

| File and anchor | Patch |
| --- | --- |
| `server/middleware/api_key_auth.go:148-165` | Resolve target after active-user authentication, before group availability and permission checks |
| `service/api_key_service.go:739-790` | Keep normal cached lookup semantics. Pass the stored lookup value only after exact ingress-token recognition |
| `service/api_key_service.go:481-489` | Reuse normal binding permission through a narrow service method |
| `server/routes/gateway.go:213-227` | Keep existing platform dispatch, and reject trusted directed context in unsupported dispatch paths |
| `handler/gateway_handler_responses.go:183-224` | Fail closed before generic account selection when directed context reaches this handler |

### Request-copy contract

Copy the API Key and target Group. Set `GroupID` to a fresh request-local value and `Group` to the separate copy.

Do not modify the stored binding, authentication snapshots, Group maps, User, account groups, or account priority.

Keep `SetOpsFallbackAPIKey`, normal group context, subscription context, and pricing context aligned with the request copy.

Call the existing group checks and subscription/limit checks after selection. Reject simple mode because it bypasses billing and group membership semantics.

Do not clear a shared Group's fallback fields. OpenAI-only validation and explicit handler guards must prevent unsupported fallback behavior.

## 3. Private Token Storage and Key Lifecycle

Plaintext registration plus DTO redaction cannot preserve token privacy after container recreation. The original program returns the database Key value.

### Options

1. Plaintext row plus DTO redaction: smallest patch, but privacy fails when the original binary runs.
2. Private random token plus deterministic lookup value: selected under confirmed K1-K7. Keep normal Key status, quotas, owner, subscriptions, and accounting.
3. Separate evaluation authentication store: larger change that duplicates Key authentication and risks bypassing normal checks.

For option 2, generate one 32-byte random private token. Store it only in protected `/app/data/codex-eval/` configuration.

Derive a namespaced SHA-256 lookup value, such as `s2eval-<hex-digest>`. Register that value through normal server-side `APIKeyService.Create`.

The lookup value is not a second random token. No schema migration is required.

At the main HTTP ingress, map only the configured private token to the lookup value before normal authentication lookup.

Keep `APIKeyService.GetByKey` and internal refresh/recovery lookups unchanged. Preserve the stored lookup value in caches and invalidation calls.

Reject externally presented owned lookup values in both main and Google ingress middleware. Reject the private raw token on every unsupported endpoint.

Internal lookup calls are not public credential replay. Do not reject them globally or break normal cache, recovery, or WebSocket billing refresh.

Grant directed capability only when the raw token digest, immutable Key ID, and configured owner match. A name or prefix match grants no capability.

Reject the lookup value itself as a gateway credential while the patch is active. It must not authorize directed requests.

The original binary cannot authenticate the private token because that token is absent from the database. Reapplication restores only the token-to-lookup mapping.

The user confirmed the conditional Key boundary K1-K7. This scheme still needs dedicated tests and is not a verified privacy guarantee.

The raw-token claim does not apply to the deterministic stored lookup value. The original binary can authenticate that value as a normal Key.

Unbound inference is disabled only when the host's ungrouped-scheduling setting rejects it. Do not assume that setting or modify it for this feature.

The [private-Key lifecycle research](private-key-lifecycle.md) evaluates stronger boundaries. It recommends considering a scoped normal row only after source-IP/proxy validation.

Loopback restrictions are not an unconditional proof because native security checks may trust forwarding headers or observe a local proxy address.

Do not alter host trust settings. Fail readiness on an unverified or unsupported topology, and review native lookup-value exposure explicitly.

K6 explicitly accepts original-runtime lookup visibility and possible normal-Key use from approved local sources. No automatic revocation is promised after unplanned recreation.

Keep raw-token secrecy across every runtime state. No control-marker, relaxed validation, or schema change has been authorized.

On planned rollback/uninstall, stop evaluations and verify owned-Key revocation. Incomplete revocation remains a transaction failure, not successful cleanup.

### Provisioning and recovery

- Initialize only through authenticated `POST /api/v1/admin/codex-eval/initialize`.
- Fix the charging owner to the authenticated administrator who initializes the Key.
- Do not choose an arbitrary first administrator, add funds, or grant subscriptions.
- Provision the owned row without changing existing business bindings.
- Persist a provisioning intent before creating the row. Recover creation/file-write interruptions by exact lookup identity.
- Protect the directory with mode `0700` and files with mode `0600` under the runtime UID.
- Restore the private identity before accepting requests. Missing or corrupt state disables evaluations, not ordinary traffic.
- Make repeated initialization idempotent. Never reenable an externally disabled Key silently.
- Redact the owned lookup value in `handler/dto/mappers.go:80-111`, including nested usage DTOs.
- Keep other Keys unchanged. Test list, detail, update, administrator user-Key lists, and nested usage results.
- Revoke only task-owned resources during uninstall. Preserve history by default.

Never echo the raw token through browser initialization, API errors, audit request bodies, process arguments, or logs.

## 4. Forced Account Selection

Intercept at `service/openai_account_scheduler.go:2298-2331`, after quota/privacy/profit context setup and before advanced/legacy dispatch or affinity lookup.

This point covers advanced scheduling on and off. Adding a field only to `OpenAIAccountScheduleRequest` would leave the legacy path unpatched.

Proposed helper signature:

```go
func (s *OpenAIGatewayService) selectDirectedEvaluationAccount(
    ctx context.Context,
    req OpenAIAccountScheduleRequest,
    target evalroute.Target,
) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error)
```

The helper must perform these steps:

1. Confirm the target group equals the request group and the platform is OpenAI.
2. Refuse an excluded target. Never reset exclusions or substitute a candidate.
3. Load the exact target through the existing schedulable-account path.
4. Confirm direct membership using `openAIStickyAccountMatchesGroup`, not a simple-mode shortcut.
5. Apply `checkChannelPricingRestriction` for the request model.
6. Reuse `resolveFreshSchedulableOpenAIAccount` and `recheckSelectedOpenAIAccountFromDB`.
7. Reuse `isAccountRequestCompatibleReason` and transport/capability checks.
8. Acquire normal account concurrency or return an existing `AccountWaitPlan` for that same account.
9. Attach the normal profit gate and emit a directed schedule decision.
10. Return an unavailable-target error without entering normal load balancing.

Relevant eligibility owners:

- `openai_account_scheduler.go:1775-1831`: model, privacy, quota pause, runtime block, proxy quarantine, shadow parent, channel, capability, and profit.
- `openai_gateway_scheduling.go:1539-1657`: current snapshot, fresh database state, scheduling thresholds, group membership, and profit.
- `openai_account_scheduler.go:574-610`: existing account-slot and wait-plan pattern.

Do not set `StickyAccountID` alone. Sticky escape and load-balance fallthrough can still choose another account.

Do not exclude every other account by scanning the pool. That depends on a mutable account snapshot and misses affinity paths.

Do not call `Forward` directly after `GetByID`. That bypasses normal admission and handler accounting.

### Handler guards

Retain `openai_gateway_handler.go:599-630` for user concurrency, eligibility recheck, and inflight balance reservation.

Retain `768-780` for account concurrency and final profit admission. Retain `817-859` for normal usage, including partial-result accounting.

For directed requests, terminate on profit veto, unsupported continuation, or failover instead of entering account reselection.

Preserve upstream health reporting. Do not let the loop switch accounts or groups after an error.

Keep normal forwarding normalization on the same account. Preserve worker-level probe retries, but each retry must keep its original group and account.

### Corrected review boundaries

The [gateway research review](gateway-contract-review.md) narrows the common-interception claim to the HTTP Responses selector.

The public `SelectAccountWithLoadAwareness` can bypass `selectAccountWithSchedulerOnce`. Reject unexpected directed context there and in the generic handler.

The quarantine wrapper at `openai_account_scheduler.go:2178-2194` can rerun selection with a bypass context. Return before this retry for directed requests.

Also reject an already bypass-marked directed context. Do not weaken the normal quarantine policy or change ordinary fail-open behavior.

Suppress directed affinity writes at `openai_profit_control.go:344-348` and `357-375`. Ignoring affinity reads alone does not prevent sticky-cache contamination.

After slot acquisition, reread the exact target and recheck membership, runtime/model/capability eligibility, and profit admission before forwarding.

`ProfitControlVetoLatest` only checks the profit rate gate. It does not replace a full post-wait account-eligibility check.

Propagate the admitted account into the handler's local forward/accounting variable. Currently it captures `account` before the slot helper refreshes `selection.Account`.

Guard the actual Responses `continue` sites: `762` continuation mismatch, `776` profit veto, `933` handler retry, and `965` account switch.

The error response at `969-971` and profit-exhausted helper at `2241-2249` are terminal, not additional reselection loops.

Group-specific multiplier calculation lives in `openai_gateway_usage.go:202-203`. The same request Key copy must reach this path.

## 5. Cancellation and Shutdown

`gateway_usage_billing.go:573-587` removes cancellation with `context.WithoutCancel`. The current worker timeout alone cannot stop a gateway upstream.

For trusted evaluation context only, preserve the incoming cancellation in both upstream-context helpers. Keep ordinary detached behavior unchanged.

Do not change `detachedBillingContext`. Already observed partial usage must still reach accounting after cancellation.

Test normal and passthrough HTTP forwarding. Confirm the upstream request context ends and account/user slots release.

HTTP Responses can also select upstream WebSocket transport at `openai_gateway_forward.go:887`. POST-only ingress does not exclude that path.

`openai_ws_forwarder_v2.go:402` directly detaches read cancellation. Cover trusted evaluation cancellation there, including the current leased connection's ownership.

Interrupt that lease without returning incomplete turn state to the pool. Keep the shared pool and ordinary draining behavior unchanged.

Preserve already observed usage in the partial forward result and normal billing task. Do not fabricate final token totals when cancellation precedes authoritative usage.

`openai_ws_http_bridge.go:493` is the opposite bridge: WebSocket client ingress to an HTTP upstream. It is not evidence for HTTP-to-WebSocket coverage.

Change `askCodex` to accept its run context. It currently creates `context.Background()` at local `codex.go:189`.

Give candy and pelican cancellable run contexts too. Current cancellation handlers only cover fingerprint and ModelTrace.

Stop evaluations before `cmd/server/main.go:190` starts HTTP shutdown. Cancel workers first, then drain their tasks while the local gateway remains available.

Persist terminal interruption state even when graceful shutdown reaches its deadline. Do not treat a partial SSE answer as success.

## 6. Evaluation API and Persistence

Use Sub2API's `{code,message,data}` administrator response envelope from `internal/pkg/response/response.go:14-56`.

Do not return the CPA management envelope from new administrator handlers. Update the browser helper to unwrap `data` and handle `message` and `reason`.

Proposed routes beneath `/api/v1/admin/codex-eval`:

| Method | Suffix | Contract |
| --- | --- | --- |
| GET | `/status` | Sanitized patch provenance, readiness, owner ID, and storage state |
| POST | `/initialize` | Idempotent private-Key and menu initialization, no token response |
| GET | `/targets` | Permitted OpenAI groups, member-account metadata, and model availability |
| GET | `/state` | Progress and latest histories |
| POST | `/run` | Evaluation type, group, accounts, model, effort, and existing batch options |
| POST | `/cancel` | Cancel selected run IDs across all four evaluation types |
| GET | `/records/:id` | Full record scoped to its account and type |
| DELETE | `/results` | Clear selected history scope, with write rollback on failure |

All operations inherit admin authentication, compliance guard, audit middleware, and panel limits from `server/routes/admin.go:26-33`.

The core must expose a narrow account metadata provider and HTTP model executor. It must never return upstream credentials through `/targets`.

Use numeric Sub2API IDs at the API boundary. Keep a stable account namespace internally, and preserve per-account exclusion across evaluation types.

Capture group ID, group name, account ID, owner ID, requested model, effort, and run ID once at run start.

A selected unavailable account produces a failed record. Do not silently filter it out as current `selectedCredentials` does.

Persist a versioned active-run record before acknowledging worker startup. Reconcile unfinished entries as interrupted at startup, without sending paid retries automatically.

Preserve atomic snapshots, corrupt-file protection, visible storage errors, and the latest 20 records per account and evaluation type.

Preserve the current probe retry policy and all four scoring algorithms. Keep group selection out of mutable global run settings.

The HTTP SSE adapter must require a terminal Responses event, parse chunk boundaries incrementally, and bound response buffers.

Handle non-2xx responses, SSE error events, incomplete responses, tool-only output, large lines, cancellation, and item-only output reconstruction.

Normalize CPA-only model effort suffixes before sending an HTTP model name. Record the requested and effective model/effort without changing evaluation prompts.

## 7. Administrator Page

Keep the page at `/codex-eval/ui`. Add an exact bypass to `web/embed_on.go:355-385` so the SPA does not intercept it.

Serve only a static shell without account or result data. All data comes from authenticated administrator APIs.

The browser reads the current same-origin `auth_token` for each API request. It does not store another management secret or refresh token.

On token expiry, stop polling and request parent-console login or refresh. Do not run an independent refresh loop that races the parent's token rotation.

Patch `frontend/src/utils/embedded-url.ts` only for the exact same-origin evaluation resource. Do not attach JWT or source-page URL parameters there.

Keep other custom-menu URLs unchanged. Include tests for same-origin matching and cross-origin lookalikes.

Reuse the existing CSP nonce. Allow only same-origin embedding on this resource, without weakening global security headers.

Keep pelican output inside its existing opaque-origin sandbox. Preserve vendored Markdown/math renderers.

Use an administrator-only stable menu ID. Merge only that entry and preserve unrelated JSON entries.

Use a narrow settings update with concurrency protection. Do not read and rewrite all settings through a stale administrator form.

The installer supplies and verifies the public origin. Never construct the menu URL from an untrusted request Host header.

### User-confirmed visual boundary

Keep the four evaluation tabs, account rows, batch controls, detail views, and results/history structure. Do not replace the page with a new dashboard.

Adapt only host-dependent controls and visual primitives: colors, surface layers, borders, text, status states, focus, and light/dark behavior.

Use exact pinned Sub2API design values. Retain the plain-JavaScript implementation and renderer assets rather than introducing Vue into the copied evaluation page.

The current Node UI harness excludes initialization. Add explicit parent-theme, storage-event, mutation-observer, and cleanup tests instead of relying on the existing harness alone.

The [UI research report](ui-sub2api-theme.md) provides the source-backed token map and theme ownership contract.

Keep local `data-theme` selectors. Use URL `theme`, shared storage, then OS preference at load time. The parent owns theme selection.

Keep only theme/language parameters in the evaluation iframe URL. Accept native iframe reload and restore server progress/history plus saved preferences.

Do not claim a child storage listener prevents parent-triggered reload. Reload-free synchronization would require another parent-side patch and is not the minimal proposal.

Separate fixed-dark code-block colors from table/inset tokens. Keep safe metadata labels and use neutral native badges instead of CPA metallic styles.

Substitute the host CSP nonce at response time. Keep global headers unchanged and avoid `unsafe-inline` script policy.

Verify actual light/dark contrast and initialization behavior in tests. The research report is source evidence, not a browser verification result.

## 8. Runtime Installer Transactions

Compare separate source, image-index, platform-manifest, image-config, binary, and patch identities. They are not interchangeable digests.

Preflight must inspect only sanitized container fields. Never print container environment values or raw deployment configuration.

Check the exact target container ID, image ID, architecture, startup command, writable root layer, runtime UID, data volume, and binary mount state.

Refuse read-only roots, binary bind mounts, unsupported source/image identities, unexpected launchers, and normal-mode incompatibility.

Persist original and patched binary checksums, image/config identity, container identity, owned resource IDs, and transaction phase in the data volume.

Keep the original backup immutable. Never back up a previously patched binary as the official original.

Stage the candidate beside `/app/sub2api`, on the same filesystem. Verify mode, ownership, architecture, provenance, and checksum before atomic rename.

Never overwrite the running executable in place. Restart only the selected container, without changing Compose, image, or account/group bindings.

After restart, check ordinary health, patch status, admin authentication, and gateway readiness. A health response alone cannot prove patch installation.

Recover failed initialization and replacement from the transaction journal. Preserve histories and unrelated resources.

Rollback restores the recorded original only when its image/source identity still matches the target. Uninstall additionally revokes owned resources.

Container recreation changes writable-layer state. `status` must compare the actual binary and patch endpoint, not trust persisted installation metadata.

Do not automatically apply an old backup or patch after an official image upgrade. Require a compatible artifact for the new identity.

## 9. Regression Tests Before Implementation

| Test | Required assertion |
| --- | --- |
| Ordinary Key with reserved headers | Rejected before any upstream call |
| Missing, duplicate, or overflowing target header | Rejected without trusted context |
| Two groups with one private token | Independent Group/subscription/pricing objects and correct usage attribution |
| Cache snapshot immutability | Original Key/Group binding remains unchanged under `-race` |
| Private token storage | Database row, DTOs, nested usage, errors, and logs never contain the raw token |
| Lookup-value replay | Cannot grant directed capability or authenticate on the patched gateway |
| Recreation privacy | Original binary cannot authenticate or reveal the private token |
| Advanced scheduler on/off | Both select exactly the requested member account |
| Affinity points elsewhere | No call to the other account, and no unrelated sticky-cache mutation |
| Account removed between selection and admission | Failure, not substitution |
| Model/privacy/quota/profit/cooldown failures | Preserve each normal veto with zero alternate-account calls |
| Upstream 429/520 and first-output timeout | Directed call never reaches account B. Ordinary requests still fail over |
| Generic Responses dispatch | Trusted evaluation context fails closed before selection |
| Shared user/account slots | Ordinary and evaluation traffic use the same counters and limits |
| Cancellation | Upstream context ends, slots release, and observed partial usage persists |
| Four evaluation fixtures | Same scoring, prompts, Codex fields, retries, and renderer behavior |
| Start/save interruption | Recover unfinished runs as interrupted, without an automatic paid request |
| Private initialization interruption | Reinstall does not create duplicate Key rows or menu entries |
| API/UI contract | Sub2API envelopes, auth expiry, same-origin embedding, and nonce behavior |
| Runtime lifecycle | No image/Compose changes, idempotence, fault rollback, uninstall retention, recreation detection |

Adapt `handler/openai_responses_failover_cancel_test.go` to assert exact upstream account-ID call sequences.

Its current fixtures use simple mode. Add normal-mode billing and shared-concurrency fixtures rather than using them as accounting proof.

## 10. Review Gates

1. Review the private lookup scheme and UI source exceptions before activation.
2. Implement routing tests and normal-traffic regressions before connecting paid workers.
3. Verify cancellation and partial accounting before enabling batch evaluations.
4. Verify installer fault injection and original-binary rollback in disposable containers before live installation.
5. Run one real candy evaluation only against an approved Sub2API target.

This machine has no Docker CLI and no configured Sub2API deployment target. Container tests and live acceptance remain blocked, not passed.

## 11. K4 Topology Evidence Contract (Stage 2A implementation, current)

Implemented in `internal/service/eval_topology.go`, consumed by `eval_key_runtime.go`,
`eval_identity.go`, and `server/middleware/api_key_auth_eval.go`.

### Evidence model

- `EvalSecuritySnapshot` binds live `ForwardedClientIPSettings()` (trust flag + header list)
  and `Server.TrustedProxies`/`TrustedProxiesConfigured`. Its SHA-256 fingerprint is embedded
  in the sealed proof; any drift in forwarded-IP trust, accepted headers, or trusted-proxy list
  invalidates the proof at validation time.
- `EvalIngressPeer{ID, Source}` is the server/installer-owned deployment ingress plan entry:
  a declared ingress path plus its expected peer source address (canonical `netip` string).
  `NewEvalTopologyCollector(cfg, loopbackSources, []EvalIngressPeer)` canonicalizes and sorts
  peers; duplicate IDs, invalid IDs, and unparseable sources fail closed.
- `EvalTopologyCollector.Observe(*gin.Context, EvalProbeCase)` records only server-observed
  fields: actual method/path, `req.RemoteAddr` (port stripped), the resolved security client IP
  via `ip.GetSecurityClientIP` with the live trust flag, and recognized forwarded-header values.
  The probe label only selects which plan-derived case the request must satisfy; it is not
  accepted as evidence. Duplicate, unknown, or mismatched cases fail. The method must be
  `POST` and path must be `/v1/responses`.
- Required cases: each whitelisted loopback source (direct, remote=source, resolved=source,
  allowed), one direct public source (resolved=remote, denied), one spoofed-loopback probe per
  recognized forwarded header (public remote, loopback in that header, denied), and per
  declared ingress an allowed probe (remote == declared peer source, resolved loopback
  allowed) plus a spoofed probe per header (remote == declared peer, forwarded loopback,
  denied). Caller-supplied labels cannot erase ingress IDs: the plan ingress set is sealed
  into the proof and `EvalIngressTopology.Validate` cross-checks claim vs sealed inventory.
- Honest limit: the server cannot invent the deployment inventory. If a proxy hop exists but
  the installer does not declare it in `[]EvalIngressPeer`, only a direct-path proof can be
  produced; detection of undeclared proxies is an installer/inventory responsibility, not
  something this primitive fabricates.

### Sealed proof and freshness

- `EvalTopologyProof{Version, Fingerprint, Whitelist, Ingresses, Observations, SealedAt,
  ExpiresAt, MAC}`. The HMAC-SHA256 (keyed by SHA-256 of domain+private token) covers all
  fields including timestamps. `SealedAt`/`ExpiresAt` are Unix seconds with a fixed bound
  `evalTopologyProofMaxAge = 24h` and future-skew bound `5m`.
- Validation rejects: absent/malformed proof, version mismatch, `SealedAt<=0`,
  `ExpiresAt<=SealedAt`, validity window > 24h, `SealedAt` more than 5m in the future, and
  `now > ExpiresAt`. Fingerprint drift, MAC failure, incomplete/inconsistent/tampered
  observations all fail closed identically (`ErrEvalTopologyUnverified`).
- `Seal` additionally refuses collectors older than the max age or created in the future.
- `InstallEvalRuntimeWithProof` now fully revalidates (fingerprint + MAC + freshness +
  observations) and stores a deep clone of the proof, so a caller cannot enable eval with an
  unvalidated proof or mutate installed evidence. `SecuritySnapshot()` returns cloned slices.
- Every gate funnels through `validateEvalTopologyProofAgainstSnapshot`: `Initialize`,
  `verifyLiveTopologyProof` (adopt/Load), `BootstrapRuntime` (expired proof installs
  identity-only guards so replay rejection and DTO redaction survive), and per-request
  `EvalRuntimeProofValid` in the API-key middleware. `evalProofStructureValid` also rejects
  persisted intents with non-positive or inverted timestamps as corrupt.

### Fresh command outputs (this session)

```text
gofmt -l internal/service internal/server/middleware        -> clean
go vet -tags=unit ./internal/service ./internal/server/middleware \
  ./internal/pkg/evalroute ./internal/handler/dto           -> clean
go test -tags=unit -race ./internal/service -run 'TestEval' -count=1
  -> ok  github.com/Wei-Shaw/sub2api/internal/service  2.345s
go test -tags=unit -race ./internal/server/middleware ./internal/pkg/evalroute \
  ./internal/handler/dto -count=1
  -> ok middleware 3.385s / evalroute 2.705s / dto 3.342s
git apply --check on pristine f2669c8 -> OK
scratch apply + diff -r vs workspace -> REPLAY_IDENTICAL
```

Full upstream `go test -tags=unit -race ./internal/service` remains baseline-blocked
(10 race reports / ~54 failures on the unmodified pinned tree; user declined fixture repair).

### Patch export method

`patches/sub2api/v0.2.15/001-codex-eval-auth.patch` regenerated deterministically with
Python `difflib.unified_diff` over the 14 sorted touched paths (3 modified, 11 new),
`a/backend/...`/`b/backend/...` labels, epoch timestamp for new files, mtime for modified.
14 files, +3863 lines.

### Spec-compliance check (trellis-check, Stage 2A re-review)

Finding fixed in place: `EvalKeyProvisioner.Initialize` dropped the freshly validated
`topo.Proof` whenever a persisted intent existed. `adoptExistingRow` re-validated only
`intent.Proof`, so once the sealed proof aged past `ExpiresAt` (24h bound) the deployment
could never renew — re-initialization failed closed permanently unless
`provisioning.json` was deleted by hand. Renewal is the designed lifecycle for a 24h
proof, so this was a soundness defect, not deferred wiring.

Fix (`eval_key_runtime.go`, intent!=nil branch): canonical-compare the fresh proof's
`[]EvalIngressPeer{ID,Source}` inventory against the persisted proof (fail closed on
drift with `ErrEvalTopologyUnverified`), adopt the fresh proof into the intent, persist
it via `SaveIntent`, then adopt the row. Whitelist drift stays bound by
`verifyLiveTopologyProof` (proof whitelist == intent whitelist) and row-IP binding by
`verifyEvalOwnedRow`. Two regressions added: renewal after expiry restores `Load`
without creating a second row; an ingress-inventory swap on renew is rejected.

Re-exported the 14-file patch with the same difflib method (epoch for new files, mtime
for modified). `git apply --check` on pristine `f2669c8` passes; scratch apply +
`diff -r` against the workspace backend is identical.

Fresh commands (this review):

```text
gofmt -l internal/service internal/server/middleware internal/pkg/evalroute \
  internal/handler/dto                                          -> clean
go vet -tags=unit ./internal/service ./internal/server/middleware \
  ./internal/pkg/evalroute ./internal/handler/dto               -> clean
go test -tags=unit -race ./internal/service -run 'TestEval' -count=1
  -> ok 2.272s (incl. TestEvalProvisionerRenewsExpiredProofWithFreshEvidence,
         TestEvalProvisionerRejectsInventoryDriftOnRenew)
go test -tags=unit -race ./internal/server/middleware \
  ./internal/pkg/evalroute ./internal/handler/dto -count=1
  -> ok middleware 3.573s / evalroute 1.384s / dto 3.527s
git apply --check (pristine)                                    -> OK
scratch apply + diff -r backend vs workspace                    -> REPLAY_IDENTICAL
```

Confirmed unchanged scope: ordinary-Key reserved-header 403 before lookup, private token
confined to POST /v1/responses with directed headers, stored-lookup replay rejected on
main and Google ingresses (private token also rejected on Google), owned-row identity
bound by ID+UserID+constant-time lookup compare, target-group resolution applied to a
Key/Group copy before normal availability/permission/subscription/billing checks,
bootstrap/revoke failure paths install identity-only guards, per-request
`EvalRuntimeProofValid` revalidates MAC/freshness/live-config fingerprint. Deferred and
tracked, not waived: probe-producer/admin-API and bootstrap wiring (Stage 5+), undeclared
proxy inventory (installer stage), full pristine service suite (baseline-blocked).

### Quality re-review (trellis-check, Stage 2A third pass)

Two primitive-level defects found and fixed in place; patch re-exported.

1. Spoof observations did not bind the resolved client (`eval_topology.go`,
   `checkObservation`). `EvalProbeSpoofedLoopback` and `EvalProbeIngressSpoofed`
   required only a loopback token in the arrived header plus `!allowed`
   resolution. A resolver that honored a forwarded loopback outside the
   whitelist (e.g. whitelist `{::1}`, resolved `127.0.0.1`) passed the case,
   certifying the exact mechanism an attacker would use with a whitelisted
   value. Both branches now also require `!resolved.IsLoopback()`, so any
   forwarded-to-loopback resolution fails regardless of whitelist contents.
   Regression `TestEvalTopologySpoofRejectsResolvedLoopbackBeyondWhitelist`
   covers direct and ingress cases plus the safe public-tail pass case.

2. No mid-collection config-drift guard. The collector froze
   `EvalSecuritySnapshot` at construction, but `Observe` resolves via
   `ip.GetSecurityClientIP`, which honors the request-scoped
   `forwardedIPSettings` override set by `SessionBindingContext` from live
   config. A transient A->B->A drift during collection could seal evidence
   gathered under different resolver semantics than the sealed fingerprint.
   The collector now keeps `cfg` and rejects `Observe`/`Seal` whenever the
   live fingerprint differs from the sealed snapshot. Regression
   `TestEvalTopologyCollectorRejectsConfigDriftDuringCollection` covers a
   rejected drift observation, incomplete-seal refusal, and empty-coverage
   seal refusal under drift.

Re-exported `patches/sub2api/v0.2.15/001-codex-eval-auth.patch` with the same
difflib method (sorted 14 paths, a/b labels, epoch for new files, file mtime
otherwise); the method byte-reproduces the other 12 file sections.

Fresh commands (this review):

```text
gofmt -l internal/service/eval_topology.go internal/service/eval_topology_test.go -> clean
go vet -tags=unit ./internal/service ./internal/server/middleware \
  ./internal/pkg/evalroute ./internal/handler/dto                              -> clean
go test -race -tags=unit -run 'TestEval|TestResolveDirected' ./internal/service
  -> ok 3.734s
go test -race -tags=unit -count=1 ./internal/pkg/evalroute ./internal/handler/dto
  -> ok evalroute 1.257s / dto 1.960s
go test -race -tags=unit -count=1 -run 'TestEval' ./internal/server/middleware
  -> ok 2.016s
git apply --check (pristine f2669c8)                                          -> OK
git archive f2669c8 -> scratch apply -> diff -rq backend vs workspace          -> REPLAY_IDENTICAL
diff -rq baseline backend vs workspace backend                                 -> exactly the 14 patch paths
```

Strict-boundary notes recorded for later stages: probe payload adequacy
(which spoof value is sent per header), path-level ingress identity beyond
`remoteAddr == peer.Source`, admin/probe wiring, and bootstrap/renewal
callers remain deferred — `checkObservation` binds evidence to declared
peers and resolved outcomes, and the runtime keeps per-request
MAC/freshness/fingerprint revalidation, but no production caller exists yet.
Loopback-sourced ingress peers cannot seal (ingress-spoof cases would always
resolve whitelisted loopback); this fails closed and is acceptable.

## 12. Stage 2B Directed Scheduling — implemented contracts

Implemented in `dist/sub2api-build/workspace`, exported as
`patches/sub2api/v0.2.15/002-codex-eval-scheduling.patch` (9 paths,
relative to pristine+001). Scratch trees preserved under
`dist/sub2api-build/scratch-*/`.

### Data flow (no-fallback guarantee)

Middleware sets `evalroute.WithTarget` on the request ctx only for the
owned eval Key on `POST /v1/responses`. `selectAccountWithSchedulerOnce`
intercepts the directed ctx after quota/privacy/profit ctx installation
and before advanced/legacy dispatch, sticky lookup, and guardian
affinity — the candidate set is exactly `target.AccountID`. Selection
reuses `tryAcquireAccountSlot` + `AccountWaitPlan` admission, so user/
account concurrency and inflight reservation are unchanged. The handler
revalidates via `RecheckDirectedEvaluationAccount` after slot/wait
admission and propagates the fresh `latest` into both `Forward` and
`RecordUsage` inputs — closing the documented stale-local-capture gap.

All escape paths fail closed: excluded target, group/platform mismatch,
bypass-marked quarantine ctx, legacy `selectAccountWithLoadAwareness`,
`selectAccountForModelWithExclusionsStickyHit`,
`selectAccountByPreviousResponseIDForCapability`, and the generic
provider Responses handler all return
`ErrEvalDirectedRoutingUnsupported` (or wrapped `ErrNoAvailableAccounts`
for target unavailability). Continuation mismatch, profit veto, and
account-switch failover terminate; same-account `RetryableOnSameAccount`
continues only because re-entry re-runs the pinned directed selection
with full fresh checks. Sticky writes are suppressed at
`setStickySessionAccountID`, `bindOpenAIStickySessionDuringSelection`,
and `BindStickySessionAfterProfitAdmission`; foreign affinity is never
read for directed ctx. Shared Group/fallback fields, cache keys, and
scheduling settings are untouched.

### Replay method

Bare-difflib unified diff identical to 001 (sorted paths, `a/`+`b/`
labels, `1970-01-01 08:00:00` for new files, file mtime otherwise). Note:
`diff --git` headers must NOT be added — git 2.54 silently skips such
sections without `index` lines (verified: "Skipped patch" + 0 files
changed while `--check` exits 0).

```text
git apply --check --whitespace=error 001; 002   -> OK (fresh scratch tree)
apply 001; apply 002; diff -rq backend vs workspace -> REPLAY_IDENTICAL
go test -tags=unit -race -run 'TestDirected' ./internal/service -> ok 3.547s
focused service+handler regressions -race -> ok 91.168s / 17.625s
```

Deferred: cancellation/WS-drain stage (forwarder_v2.go:402 detach,
lease MarkBroken, partial usage), eval core, admin/probe APIs,
installer. Full pristine service suite still baseline-blocked.

### Stage 2B independent spec review addendum (check agent)

Source review confirmed the nine-path patch implements the contract as
designed; no code defect required self-fix. Verified against workspace
source, not just the diff:

- `openai_account_scheduler.go` ~2328: directed intercept after
  quota/privacy/profit ctx installation and platform normalization,
  before guardian affinity and advanced/legacy dispatch.
- `openai_eval_directed.go`: `freshDirectedEvaluationAccount` rejects
  quarantine-bypass ctx, image capability, group mismatch, non-positive
  account; then channel-pricing → `getSchedulableAccount` → normalized
  platform guard (the exact `account.Platform != platform ||
  !IsOpenAICompatible()` check inside
  `openAICompatibleAccountEligibilityFailureReasonBeforeProfit` is what
  actually refuses normalized non-OpenAI providers) →
  `openAIStickyAccountMatchesGroup` (normal-mode direct membership,
  independent of `RunMode`) → `resolveFreshSchedulableOpenAIAccount`
  (model/rate-limit/quota/parent/runtime-block/threshold/quarantine +
  profit veto) → `recheckSelectedOpenAIAccountFromDB`
  (`accountRepo.GetByID` when `schedulerSnapshot != nil`, fresh object
  re-checked for membership/privacy/eligibility) → capability,
  transport, `isAccountRequestCompatibleReason` (composite-ownership,
  privacy, runtime, quarantine, quota-pause, shadow-parent health,
  model, channel restriction, capability, profit veto).
- Slot admission reuses `tryAcquireAccountSlot` /
  `AccountWaitPlan{FallbackWaitTimeout,FallbackMaxWaiting}` — identical
  parameters to the ordinary fallback path
  (`openai_gateway_scheduling.go` ~1496).
- Handler: `evalroute.FromContext` read once; profit-veto and
  continuation-mismatch branches terminate; post-slot
  `RecheckDirectedEvaluationAccount` under
  `ContextWithSelectionProfitGate` replaces `account`/`selection.Account`
  before `Forward`, `submitResponsesUsage`/`RecordUsage`, and
  `resolveOpenAIUpstreamEndpoint`; release func runs on recheck failure;
  account-switch failover terminates after
  `ReportOpenAIAccountScheduleResult`; same-account
  `RetryableOnSameAccount` re-enters selection with the target still
  eligible. `apiKey.GroupID`/`Group` are the middleware request-scoped
  copies set by `applyEvalDirectedTarget`, so the directed group is what
  reaches selection and billing.
- Auth-side enforcement: `ResolveDirectedTargetGroup` requires
  `group.Platform == PlatformOpenAI` plus `GetAccountIDsByGroupIDs`
  membership; simple mode, non-POST and non-`/v1/responses` paths are
  rejected before the target ctx is ever installed, so Messages /
  WebSocket handlers cannot carry a directed ctx.
- Fresh scratch `dist/sub2api-build/scratch-check2b`: sequential
  `git apply --check --whitespace=error` then apply of 001+002 on
  pristine `f2669c8` succeeded; `diff -rq` vs workspace is clean
  (REPLAY_IDENTICAL). 001 patch file unchanged
  (sha256 1274cf97…). Focused `-race` runs: `TestDirected` ok 3.523s;
  handler focused set ok 18.611s; service focused set ok 90.357s.

Review finding — coverage gap (returned BLOCKED to implementer): zero
directed coverage in `internal/handler` tests; the post-slot test calls
`RecheckDirectedEvaluationAccount` against static fixtures instead of
mutating the target during a real slot wait, and nothing proves the
promised account-ID upstream sequence, sticky-write suppression at the
handler layer, slot release on recheck failure, billing attribution, or
the directed-vs-ordinary 429/520 contrast. Ordinary regression passes
alone do not discharge the stage-3 gate ("exact upstream account-ID
call sequences, normal failover regressions, slot release, and billed
partial usage"). Minor, non-blocking: implementer-added comments exceed
the preserve-comments instruction (left in place per review rules), and
the `s.concurrencyService == nil` fallback in
`selectDirectedEvaluationAccount` is unreachable
(`tryAcquireAccountSlot` already returns an acquired no-op result when
the service is nil).

### Section 12 addendum — handler-level coverage (implementer)

`internal/handler/openai_eval_directed_handler_test.go` now exercises
the full `Responses` handler with real WaitPlan waits: mid-wait
mutations (membership, model_mapping, deletion, privacy, account quota,
runtime cooldown, shadow-parent health, channel restriction) all
terminate 503 post-recheck with released slots and no upstream traffic;
fresh-object propagation is asserted at the upstream `Authorization`
header and the persisted `UsageLog`; directed 429/502/520 terminates
while ordinary requests fail over; pool-mode same-account retry stays
pinned; profit veto and continuation-mismatch terminate post-slot;
sticky reads/writes are suppressed for directed requests. Replay of the
re-exported 10-path 002 patch on pristine+001 is byte-identical to the
workspace (`scratch-replay3`).

### Section 12 addendum 2 — final rework (implementer)

All Stage2B-added comments were reverted (verified: patch contains no
`+` comment lines other than `//go:build unit`; zero baseline lines
removed in any modified file). Added
`TestEvalDirectedResponsesRetryDelayInvalidationTerminates` (target
made ineligible inside the blocked first upstream call → same-account
retry re-selection fails closed: `[501]` call trace, balanced slots) and
`TestEvalDirectedGenericResponsesRejectsDirectedContext` (real
`GatewayHandler.Responses`: directed ctx → 403, zero selector snapshot
calls; ordinary control invokes the selector). Service:
`TestDirectedRejectsUnsupportedResponsesCapability` covers the nonempty
`requiredCapability` gate via
`extra.openai_responses_mode=force_chat_completions` (plain-text
Responses uses `requiredCapability=""`, so this gate is not reachable
from the handler path). Final replay: `scratch-replay5`, pristine
`f2669c8` + 001 + re-exported 002 (10 paths), `--whitespace=error`
clean, `diff -rq` byte-identical to workspace. 001 SHA-256 unchanged.
Gate pending re-review.

### Stage 2B re-review addendum (check agent, post-follow-up)

The implementer follow-up delivered the previously-missing handler
coverage in `internal/handler/openai_eval_directed_handler_test.go`
(1208 lines, `//go:build unit`, 11 `TestEvalDirected*` cases) and
reverted all unauthorized new comments; 002 is now 10 paths / ~2007
lines. Verified against workspace source and real handler execution:

- Real `Responses` handler invoked through a gin test context with the
  eval-route target installed on the request ctx; upstream calls
  recorded per-account by a stub `HTTPUpstream`.
- Directed 429/502/520 produce exactly one upstream call on the target;
  the ordinary control proves alternate-account failover still works,
  so fixtures do not mask switching. `pool_mode` same-account retry
  yields `[501,501]` (eligible retry re-runs pinned selection); a
  mid-delay `TempUnschedulableUntil` mutation terminates with a single
  call — eligible and ineligible retry are distinct cases.
- Wait-plan tests enter a real `AcquireAccountSlotWithWaitTimeout`
  queue (deny → `IncrementAccountWaitCount` signal → mutate repo →
  grant). Eight post-slot mutations (membership, model mapping,
  deletion, privacy flag, quota, runtime cooldown, shadow parent,
  channel restriction) each terminate with 503, zero upstream calls,
  and balanced acquire/release counts for both account and user slots.
- `WaitPlanForwardsFreshAccount` proves `latest` propagation: mutated
  `api_key` credential reaches the upstream `Authorization` header and
  mutated `RateMultiplier` reaches `UsageLog.AccountRateMultiplier`
  with correct `AccountID`/`GroupID`.
- Profit veto mid-wait, OAuth continuation mismatch (400), cross-group
  target (pre-selection reject, zero acquires), and the generic
  provider Responses 403 before any selector call are all covered.

Reviewer self-fix (bounded, test-only): the delivered
`stickySetKeys()` skipped `openai:`-containing keys — i.e. every real
sticky write — making `require.Empty` vacuous, and the foreign binding
was seeded under the raw hash while the read path queries the
`openai:`-prefixed key, making the ignore-foreign-binding assertion
vacuous. Fixed both; the filter now collects `openai:`-prefixed keys
excluding legitimate `openai:response:` /
`openai:http-response-owner:` continuation records, and the seed uses
the production cache key. 002 was re-exported (test-hunk only change)
and sequential 001→002 replay on pristine `f2669c8` is byte-identical
to the workspace; 001 sha256 unchanged
(`1274cf97…`). Scratch trees preserved under
`dist/sub2api-build/scratch-check2b-r*/`.

```text
TestEvalDirected ./internal/handler -race  -> ok 5.575s
TestDirected     ./internal/service  -race -> ok 3.951s
focused handler regressions -race          -> ok 18.335s
gofmt/vet touched files                    -> clean
```

Spec verdict: PASS for the Stage 2B contract. Independent QUALITY
review remains the pending scheduler gate; pristine full-suite race
baseline, Docker daemon, and live-target acceptance unchanged.

### Stage 2B independent quality review addendum (check agent)

Independent quality re-review of the final 10-path 002 patch, verified
against workspace source line-by-line (not the diff alone) and fresh
test runs:

- Exhaustive `continue`-site audit of `OpenAIGatewayHandler.Responses`:
  directed requests can only ever re-enter selection via the
  same-account `RetryableOnSameAccount` retry (:967), which re-runs the
  pinned directed selection with full fresh eligibility — verified by
  both the eligible `[501,501]` case and the mid-retry-delay
  invalidation `[501]`-only case. Every other loop escape (continuation
  mismatch :765, profit veto :777, account-switch :970 including
  first-output-timeout and OAuth-429 failover paths that fall through to
  it) terminates before `failedAccountIDs`/`switchCount`/
  `RecordOpenAIAccountSwitch` mutation.
- Post-slot `RecheckDirectedEvaluationAccount` replaces both the local
  `account` and `selection.Account`; `Forward`, upstream endpoint
  resolution, health reporting, and `RecordUsage` all consume the fresh
  DB object. No stale scheduler snapshot reaches accounting.
- Fixture soundness re-audited for residual masking: repo embeds are nil
  (unexpected calls panic, not silently pass); `schedulerSnapshot` is
  non-nil so `recheckSelectedOpenAIAccountFromDB` exercises the real
  `accountRepo.GetByID` path; wait-plan tests synchronize through the
  real `IncrementAccountWaitCount` signal; usage submission runs
  synchronously (nil worker pool → inline fallback) so UsageLog
  assertions are deterministic; `stickySetKeys` now collects production
  `openai:`-prefixed keys and excludes only documented
  `openai:response:`/`openai:http-response-owner:` records (prefixes
  confirmed at openai_ws_state_store.go:15-17); the foreign binding is
  seeded at the production key shape and `session_id` is a recognized
  session header.
- Ordinary controls present for every masked-behavior claim: alternate-
  account failover, generic-handler selector invocation, sticky eager
  write under ordinary ctx.
- Non-blocking notes: unreachable `concurrencyService == nil` fallback
  in `selectDirectedEvaluationAccount` (cosmetic dead code);
  `refreshStickySessionTTL` is unguarded but unreachable for directed
  ctx today (all callers sit behind directed-rejected selectors);
  normalized OpenAI-compatible request platforms enter directed
  selection but the target account itself must still pass the exact
  OpenAI platform check — no substitution escape.

```text
gofmt -l (10 files) / go vet -tags=unit service+handler -> clean
TestEvalDirected handler -race   -> ok 5.525s
TestDirected service -race       -> ok 2.711s
handler focused regressions -race -> ok 17.120s
service focused regressions -race -> ok 92.031s
git apply --check --whitespace=error 002 on scratch-pristine001 -> OK
git apply -R --check 002 on workspace -> OK (exact match)
002 adds zero comment lines besides //go:build unit
001 sha256: 1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b (unchanged)
```

**Quality verdict: PASS.** Stage 2B scheduler gate cleared on both
spec and quality. Inherited blockers unchanged: full pristine service
race suite baseline-blocked (fixture repair declined by user), Docker
daemon unreachable, no approved live target, cancellation/WS-drain
stage deferred as designed.

## 13. Stage 2C Directed Cancellation — implemented contracts

Implements section 5 (Cancellation and Shutdown) for the gateway
patch surface.

### Context policy (ownership)

- `upstreamRequestContext(ctx)` /
  `streamUpstreamRequestContext(ctx, stream)` in
  `gateway_usage_billing.go` are the single decision point: trusted
  directed target (`evalroute.FromContext`) → request context passes
  through unchanged; otherwise → previous detached behavior
  (`context.WithoutCancel` semantics) bit-for-bit.
- Call sites converted: `openai_gateway_forward.go` transport loop,
  `openai_gateway_passthrough.go`, `openai_gateway_cc_pipeline.go`,
  `gateway_forward_as_responses.go`,
  `gateway_forward_as_chat_completions.go`.
- `detachedBillingContext` and the billing-side
  `detachUpstreamContext` call inside `finalizePostUsageBilling` are
  deliberately not routed through the new helpers — usage recording
  must never observe request cancellation.

### WebSocket lease contract (HTTP → upstream WS v2)

- Directed: `upstreamReadCtx = context.WithCancel(ctx)`;
  `markClientDisconnected` fires the cancel; the detached-drain
  `continue` requires `!directedEval`. Read failures follow the
  existing `lease.MarkBroken()` path, the deferred `Release()` drops
  the broken lease, and the shared pool remains open. A later turn
  dials a new connection — partial turn state is never reused.
- Ordinary: unchanged `context.WithoutCancel` drain to terminal.
- `openai_ws_http_bridge.go` (client-WS → HTTP) is out of scope per
  section 5 and untouched.

### Partial result and billing contract

- Directed stream/passthrough errors propagate a partial
  `OpenAIForwardResult` (`openAIDirectedStreamPartialResult`,
  `openAIDirectedPassthrough{Stream,NonStream}PartialResult`, WS tail
  return in `Forward`) so `submitResponsesUsage` commits exactly one
  record carrying authoritatively observed usage. Progressive
  pre-terminal usage merges through `mergeOpenAIUsageNonZero`; nothing
  is synthesized when no usage was observed. Interrupted results are
  never marked scheduling-successful.
- Directed requests remain bound to the admin-selected
  group/account: scheduling (patch 002) plus `handleFailoverExhausted`
  keep cancellation from ever switching account or group.

### Evidence

- Focused unit: service ok 2.615s, handler ok 4.290s; race ok
  2.453s/4.138s; ordinary WS drain regression ok (race 1.982s).
- `go vet` (default and `-tags unit`) clean; `gofmt -l` clean; patch
  adds zero comment lines besides `//go:build unit`.
- Replay: `git apply --check --whitespace=error` 001→002→003 on
  pristine `f2669c8` → APPLY_OK; the initial 9-file export was later
  superseded by the final 11-file export recorded below.
- Historical initial Patch 003 SHA-256:
  `7a0da381b05b4a47e354f6a767f3a0db1b3775d2cbb16b5cb0d0a52456933801`.
- Pending: independent spec and quality review gates.

### Section 13 addendum — independent spec review (check agent)

Source-verified against `dist/sub2api-build/workspace` (not the
implementer report):

- Historical directed-reachability claim (superseded by Section 14):
  `api_key_auth_eval.go:51` restricts the
  trusted target context to `POST /v1/responses` only, so every other
  `detachUpstreamContext`/`detachStreamUpstreamContext` call site
  (chat completions, messages, images, anthropic-native, grok, WS
  ingress bridge) is unreachable for directed requests. The five
  converted call sites were previously treated as covering the full
  reachable surface; the two
  `gateway_forward_as_*` conversions are harmless belt-and-suspenders
  behind the generic-handler directed rejection.
- `handleScanErr`/`finalizeStream` return `resultWithUsage()` on
  `context.Canceled`, so `openAIDirectedStreamPartialResult` sees real
  observed usage; passthrough mirrors its own success-path result
  fields (BillingModel empty unless imageCount, matching baseline).
- Progressive usage on non-terminal events is a documented baseline
  extension (`parseSSEUsageBytesWithType`: "retaining progressive
  usage from compatible upstreams on any event type"), so the
  delta-event usage fixtures exercise the real accepted contract, not
  an impossible shape.
- WS path: acquire/prewarm/write use `ctx` (cancel propagates through
  dial, queue, acquire, write); `upstreamReadCtx=WithCancel(ctx)`;
  `markClientDisconnected` cancels instead of detaching; read errors
  hit the existing `lease.MarkBroken()` + deferred `Release()`; the
  `!directedEval` gate on the drain `continue` and the line-796
  `clientDisconnected && terminalEventCount == 0` partial-result
  return are verified in source. Pool stays open (test: second turn
  dials a fresh conn, dial count 2).
- Billing ctx stays detached: `detachedBillingContext`,
  `wrapUsageRecordTaskContext`, and the
  `finalizePostUsageBilling` detach are untouched — canceling a
  directed request cannot lose an already-committed usage record.

Reviewer self-fix (bounded, patch re-exported, replay-verified):

- `openai_gateway_handler.go` `UpstreamFailoverError` branch dropped
  `result` on every terminal path (`handleFailoverExhausted`,
  directed-terminate, same-account retry `continue`). Baseline
  invariant "failover errors always carry result=nil" no longer holds
  for directed: `openAIDirectedStreamPartialResult` and the WS tail
  can return `result != nil` together with a failover-typed error.
  Fix: one `submitResponsesUsage(result)` line added as the first
  statement inside `errors.As(err, &failoverErr)`. No-op for ordinary
  traffic (result is always nil there); each directed attempt's
  observed usage now bills exactly once on every handler branch,
  satisfying the section-13 contract.

Fresh verification this review:

```text
gofmt -l openai_gateway_handler.go                  -> clean
go vet ./internal/service ./internal/handler        -> clean (default tags)
go vet -tags=unit ./internal/handler                -> clean
go test -tags=unit -race -run 'TestEvalDirected' ./internal/handler
  -> ok 5.590s
go test -tags=unit -race -run 'TestEvalDirected|TestDirected' ./internal/service
  -> ok 2.908s
go test -tags=unit -race -run 'TestOpenAI|TestResponses|TestFailover|TestProfit|TestGateway' ./internal/handler
  -> ok 19.979s
go test -tags=unit -race -run 'TestOpenAI|TestSelectAccount|TestEval|TestProfit|TestSticky|TestGateway' ./internal/service
  -> ok 93.545s
go test -race -run 'TestForwardOpenAIWSV2' ./internal/service
  -> ok 3.879s (ordinary detached-drain regression, untagged)
Re-exported 003 via diff -ruN a/backend b/backend (patchgen-003);
git apply --check --whitespace=error on scratch-check2c
(pristine+001+002, preserved) -> APPLY_CHECK_OK; apply +
diff -rq backend vs workspace -> byte-identical.
New patch: 10 paths (added openai_gateway_handler.go), SHA-256
31803abfd098f6c48a336c1002df82b49577ae4eaed79831947a832601e62402;
superseded file preserved at /tmp/003-orig-backup.patch.
Patch contains zero added comment lines besides //go:build unit.
001/002 SHA-256 unchanged.
```

Coverage notes (not blocking): cancel during WaitPlan/slot-acquire
relies on the shared `ctx` path (code-reviewed, no dedicated fixture);
cancel-during-WS-acquire is covered implicitly by acquireCtx deriving
from request ctx. Full pristine service suite remains
baseline-blocked; Docker/live acceptance unchanged.

## 14. Withdrawn Native Anthropic Candidate

An intermediate review misclassified native Anthropic conversion as reachable
for directed evaluation. Its Kimi fixture called `Forward` directly and bypassed
the exact OpenAI selector. The native additions were withdrawn.

The superseded 11-path candidate had SHA-256
`32a5f532115b111d1fe643fa4b49e95deae8213ccdb175b7e528da2ad3b31cda`.
Its native-only test results do not count toward the final directed scope gate.
The final correction below and [independent source review](native-cancellation-final-review.md)
record the retained scope and final patch identity.

## 15. Stage 2C Final Scope Correction

Trusted directed evaluation is constrained twice to exact
`PlatformOpenAI`: `ResolveDirectedTargetGroup` requires an OpenAI group, and
`freshDirectedEvaluationAccount` fixes the scheduler platform and rejects
provider platform IDs such as Kimi. A `PlatformOpenAI` account is not a
multi-protocol account and therefore does not select `APIProtocolAnthropic`.
The native Responses conversion is reachable for ordinary OpenAI gateway
traffic, but not for the current directed evaluation selector.

The native production edits and native-only fixtures were withdrawn. The
final cancellation patch retains the real directed HTTP, passthrough, CC,
WS lease, and failover partial-billing paths:

- Patch 003: 10 paths, SHA-256
  `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.
- Patch 001 SHA-256 remains
  `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`.
- Patch 002 SHA-256 remains
  `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`.
- Strict patch application to the fresh `scratch-review-003-r3`
  passed, and its backend tree is byte-identical to the workspace.

Final Stage 2C directed scope verdict: **PASS**. Native-only tests are not
part of the final gate. The pristine full service race baseline, Docker
daemon, and approved live target remain blocked.

### Section 13 quality re-review addendum

The reviewer corrected the stale handler comment claiming every failover error
has a nil result, renamed the zero-usage test to describe its zero-token record,
and added `TestEvalDirectedResponses_FailoverPartialUsageBilledOnce`. The new
handler regression verifies one usage row with the observed 3/5 tokens, the
directed account ID, and balanced slots for a failover-typed partial result.

Final patch 003 SHA-256 is
`6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.
Strict application on pristine+001+002 matched the workspace. Directed service
and handler race tests, the ordinary WebSocket drain regression, unit vet, and
formatting checks pass. No behavior or scope changed beyond the reviewed
failover partial-billing fix and its regression coverage.
