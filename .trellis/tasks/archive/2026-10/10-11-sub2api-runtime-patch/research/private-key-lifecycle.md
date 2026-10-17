# Research: Private Evaluation Key Lifecycle

- **Query**: Can the private evaluation token / stored `s2eval-<digest>` lifecycle satisfy normal billing/auth, no browser-usable credential, and unchanged official image/schema more strongly after unplanned container recreation?
- **Scope**: internal source review of pinned `dist/sub2api-v0.2.15` at `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`
- **Date**: 2026-10-11
- **Status**: the user confirmed conditional K1-K7 after this research. No patch was applied and no live validation was executed.
- **Authority**: the current PRD owns the accepted boundary. Research alternatives below remain evidence, not approved implementation exceptions.

## Bottom Line

**Conditional recommendation: keep a normal `api_keys` row, retain raw-token mapping, and add a verified loopback-only IP restriction through normal `Create`.**

Normal authentication enforces this restriction before user/group/subscription/billing checks (`api_key_auth.go:132-146`, Google variant `:92-105`).

However, an IP whitelist alone does not prove remote clients cannot authenticate. Forwarded-IP trust and reverse-proxy topology must also be verified.

Use only tested container-local source addresses. Do not widen the whitelist or change Compose, schema, business bindings, or host security policy.

This does **not** make the digest cryptographically unusable: a process able to originate from an allowed loopback address can still use it on the original executable. It is a privacy/exposure boundary, not revocation after unplanned recreation. Any claim of stronger post-recreation denial requires an original-binary-visible credential path that does not exist under R9.

A server-generated control-character marker is **not** a clean server-only credential. The source has a Google query-key path that can present values ordinary HTTP headers cannot (`api_key_auth_google.go:221-258`), while PostgreSQL cannot store NUL and the 128-byte schema/API ceiling cannot be exceeded (`ent/schema/api_key.go:37-40`, `api_key_service.go:49,739-741`). A CR/LF marker can be persisted by bypassing public `CustomKey` validation, but remains presentable through the Google query path; trailing-whitespace variants would rely on the brittle normalization behavior explicitly excluded by the request.

## Relevant Source Contract

| File / anchor | Finding |
|---|---|
| `backend/ent/schema/api_key.go:37-40` | `api_keys.key` is non-empty, unique, `MaxLen(128)`. |
| `internal/service/api_key_service.go:49` | `MaxAPIKeyCredentialBytes = 128`. |
| `internal/service/api_key_service.go:402-419` | Public `ValidateCustomKey` allows only `[A-Za-z0-9_-]`; control characters cannot use public `Create`. |
| `internal/service/api_key_service.go:492-598` | `Create` performs user lookup, optional group binding permission, duplicate checks, creation limits, ownership, status, IP ACL, quota/rate fields, then invalidates auth cache. `GroupID == nil` is supported. |
| `internal/repository/api_key_repo.go:45-75` | Repository creation writes the supplied `Key` through ent; a narrow server method could bypass public charset validation, but uniqueness/DB constraints remain. |
| `internal/server/middleware/api_key_auth.go:49-96` | Main extraction accepts Bearer, `x-api-key`, and `x-goog-api-key`; `key`/`api_key` query parameters are rejected before lookup. |
| `internal/server/middleware/api_key_auth_google.go:221-258` | Google extraction accepts `x-goog-api-key`, Bearer, `x-api-key`, and query `key` on `/v1beta` and `/antigravity/v1beta`; values are trimmed. |
| `internal/service/api_key_service.go:739-790` | `GetByKey` does exact cached authentication and then sets `apiKey.Key` to the presented lookup argument. |
| `internal/service/api_key_auth_cache_impl.go:194-197,269-331,443-452` | Auth cache key is SHA-256 of lookup input; materialized `APIKey.Key` is also the lookup input. Raw token must therefore be translated before `GetByKey`. |
| `internal/server/middleware/api_key_auth.go:132-146` | Existing IP whitelist/blacklist enforcement runs before user/group checks in normal auth. |
| `internal/handler/dto/mappers.go:80-111` | DTO returns `APIKey.Key` verbatim, including user-facing list/detail and nested structures. |
| `internal/handler/api_key_handler.go:180-227,229-318` | User create/update/delete routes expose ordinary management over an owned row; `custom_key` creation is idempotency-wrapped. |
| `internal/handler/admin/apikey_handler.go:25-76` | Admin can rebind/unbind the row and reset rate usage; response returns full DTO key. |
| `internal/server/middleware/audit_log.go:189-227` | Request bodies are audited after `service.RedactAuditBody`; `custom_key`, `key`, authorization-like fields, token/secret/apikey names are redacted (`service/audit_log.go:122-159`). |
| `internal/handler/ops_error_logger.go:1260,1380,1568` and `openai_gateway_handler.go:4323,4398` | Ops/error metadata records an 8-character key prefix, not the full key. Marker design must keep this prefix non-sensitive and free of controls. |
| `internal/service/gateway_usage_billing.go:382-385` | Quota exhaustion invalidates auth cache using `p.APIKey.Key`; this must remain the stored lookup value, not the raw private token. |
| `internal/handler/openai_gateway_handler.go:2129-2208` | Long-lived OpenAI WS turns can call `GetByKey(conn.Key)` to refresh billing group snapshots. `conn.Key` must therefore remain the stored lookup value. |

## Credential Extraction Trace

1. **Main gateway auth**: query `key` and `api_key` are rejected (`api_key_auth.go:49-56`). Bearer parses `Authorization` by a literal space and trims the token (`:59-67`); `x-api-key` is read untrimmed (`:70-73`); `x-goog-api-key` is a later fallback (`:81-84`). All are limited to 128 bytes (`:18,74-79,290-296`).
2. **Google/Gemini auth**: extraction is `x-goog-api-key` → Bearer → `x-api-key` → query `key` for allowed `/v1beta` paths (`api_key_auth_google.go:221-258`). Header values and query values are trimmed (`:226-249`).
3. **Service calls**: only the two middlewares call `APIKeyService.GetByKey` for ingress in production source. A separate WS-turn refresh reuses `conn.Key` established by authentication; it is not a new public credential presentation path (`openai_gateway_handler.go:2129-2197`).
4. **Consequence**: a marker containing CR or LF cannot be placed in ordinary HTTP header values, but can be represented in an allowed Google query parameter. NUL would be a stronger marker but conflicts with PostgreSQL text storage and cannot be used without a forbidden storage/schema change. Invalid UTF-8 is also unsafe because `varchar` storage must satisfy database encoding. Do not use trailing spaces/control trim tricks as the security boundary.

## Alternatives Evaluated

### A. Accept digest residual + privacy boundary + rollback rotation

**Contract**: raw token stays only in `/app/data/codex-eval/`; `s2eval-<digest>` is an ordinary custom key; patched auth maps only the exact configured raw token to the digest; patched ingress rejects external digest replay; uninstall/rollback rotates/revokes when planned.

- **Preserves**: normal `api_keys` accounting, quota, rate windows, owner concurrency, subscriptions, cache invalidation, group checks, usage rows, and R9.
- **Weakness**: after unplanned recreation the original binary displays the digest and accepts it as an unrestricted key wherever normal Key auth applies. If unbound, inference depends on `IsUngroupedKeySchedulingAllowed`; source review cannot assume it is false. Planned rollback can revoke because the patch/installer is still running; unplanned recreation cannot guarantee revocation.
- **Verdict**: feasible but weakest. It should not be the final answer without at least IP scoping.

### B. Server-only marker in the existing row

**Contract considered**: a narrow internal method creates/looks up a `key` that cannot pass public `ValidateCustomKey`, while reusing owner/user/group/subscription/quota/status/repository logic. The raw private token remains outside DB; server-side mapping resolves it to the marker.

- **What works**: an internal method can persist printable-invalid-to-public-`CustomKey` strings because repository `Create` writes `SetKey` directly (`api_key_repo.go:45-75`). Existing auth then performs exact lookup, cache, billing, and account attribution unchanged.
- **What fails**: “never accepted by original HTTP auth” is not guaranteed. CR/LF blocks normal header presentation but remains representable through Google query `key` (`api_key_auth_google.go:246-250`). NUL cannot be stored in PostgreSQL text. Exceeding 128 bytes is forbidden by schema and service checks. Trailing CR/LF would defeat query `TrimSpace`, but that is precisely a whitespace-normalization trick and would remain brittle across extraction changes.
- **Additional cost/risk**: bypassing public charset validation creates a policy exception, DTO/logging/UI escape concerns, and a non-printable unique key that ordinary update/delete/cache code may carry into audit or diagnostics. Prefix-wide protections would regress ordinary keys; protection must be exact row identity.
- **Verdict**: **not recommended** as the primary contract. It is defensible only if stakeholders explicitly approve an internal non-HTTP marker and accept residual Google-query/other transport handling plus strict DTO/log redaction. It does not solve original-binary misuse more reliably than the loopback-scoped digest.

### C. Dedicated internal registration / existing metadata or scoped normal row

#### C1. Loopback-scoped ordinary digest row — recommended

Create the same `s2eval-<digest>` row through normal `APIKeyService.Create`, with `IPWhitelist` restricted to the evaluated container-local source addresses. This preserves all normal checks and gives the original binary a real, source-enforced barrier after recreation.

- **R4/R5**: row remains a normal Key; directed capability still requires exact configured row ID, owner, lookup digest, patched mapping, and routing headers. Group selection is request-scoped.
- **R7**: browsers receive no raw token. The original DTO can still expose the digest. Remote usability depends on verified source-IP resolution and proxy topology.
- Local code at an allowed source can use the digest while the original executable runs. This residual requires explicit acceptance.
- **R9**: no schema, image, Compose, controller, or business binding change.
- **Approval caveat**: the loopback address list must be verified on the actual deployment (`127.0.0.1` vs `::1`; no assumed Docker network translation). This is deployment evidence required before implementation, not a silently chosen global setting.

#### Source-IP Safety Gate Added by Main Review

`config/config.go:2115` defaults legacy forwarded-IP trust to true. `pkg/ip/ip.go:164-171` then uses raw forwarding-header resolution for security checks.

Disabling that compatibility mode uses Gin's trusted-proxy chain. The task must not change either setting automatically.

A public request might resolve to loopback through forwarded headers or a local reverse proxy. Therefore, a loopback whitelist is not an unconditional isolation proof.

Before initialization reports readiness, verify native source-IP behavior for direct clients and every deployed ingress path. Reject an unsupported or unverified topology.

Tests must include public traffic with conflicting forwarding headers and a local reverse proxy. Verify both main and Google authentication paths.

Keep raw-token privacy, browser-visible lookup values, ordinary inference, and directed capability as separate acceptance claims.

#### C2. Separate evaluation auth store / other metadata tables

A separate store keyed by a hash could map the private token to a normal row ID, but public gateway authentication still needs a presentable stored `key` unless patched auth bypasses or rewrites lookup. Existing user attributes/settings are not consulted by `GetByKey`; storing metadata there does not remove the ordinary credential property of the row. A disabled/expired row would avoid external use but requires the patch to bypass a normal status/expiry check, weakening Key lifecycle semantics and violating “preserve normal policy” unless explicitly approved.

- **Verdict**: not recommended. It adds lifecycle complexity without a stronger invariant than C1.

## Meaning of “Normally Registered Evaluation Key”

The phrase should **not** require the public HTTP `POST /api/v1/keys` endpoint itself. A narrow server-side call to `APIKeyService.Create` can still be a normally registered Key because it reuses the same ownership, group permission, uniqueness, creation-limit, quota, status, repository, and cache logic.

However, “normally registered” should be clarified to state whether it permits:

1. **Normal service create with normal public input constraints** — recommended for the loopback digest row; no charset exception.
2. **Internal service create reusing checks but bypassing `ValidateCustomKey`** — required only for control-character markers; this is a policy exception that must be explicitly approved and should not be silently called “normal.”
3. **Repository-only create** — not recommended because it can skip service-level ownership/group/creation-limit policy unless those checks are duplicated exactly.

## Required Identity and Ownership Guards

Use exact configured identity, never a reserved prefix:

- Persist a provisioning intent containing immutable row ID once created, owner user ID, expected lookup digest, storage phase, and patch version/checksum.
- Directed capability requires raw-token mapping + authenticated row ID + configured owner ID. A name or `s2eval-` prefix grants nothing.
- On `ExistsByKey` conflict, recover by exact lookup identity and owner/status/journal; never retry blindly because `ExistsByKey` conflict returns `ErrAPIKeyExists` and increments custom-key error count (`api_key_service.go:544-551`).
- Do not reserve or reject all `s2eval-*` ordinary keys.
- Patch user `Update`, `Delete`, and admin `UpdateGroup`/rate-reset paths only by exact owned row ID. The original binary cannot enforce these guards after recreation; that residual management exposure must be acknowledged.
- Never silently re-enable or rebind an externally disabled/rebound owned row. Status should report conflict and require explicit administrator repair/revocation.
- DTO redaction should apply by exact owned row ID at every output boundary, not by prefix. Do not mutate `APIKey.Key` before billing/cache paths that need the stored lookup value.

## Hidden Raw-Token Risks Checked

- `GetByKey` and cache materialization assign the **lookup argument** to `APIKey.Key`. Mapping must occur before `GetByKey`, so context and billing see the digest/marker rather than the raw token.
- Auth caches key by SHA-256 of the lookup argument; using the digest consistently preserves invalidation and prevents raw token cache residency. If a server-only marker is used, the marker—not raw token—must be the lookup argument.
- Detached billing uses `p.APIKey.Key` for quota-exhaustion invalidation; do not overwrite it with raw token.
- Long-lived WS billing refresh calls `GetByKey(conn.Key)`; `conn.Key` must remain the stored lookup value.
- DTO mappers return the stored lookup value verbatim; browser safety depends on redaction plus loopback scoping.
- Audit JSON bodies redact `custom_key`, `key`, authorization/cookie, token/secret/apikey fields. New evaluation endpoints that accept or generate credentials should omit bodies entirely or verify redaction tests.
- Ops metadata records only an 8-byte prefix. Keep a safe fixed marker prefix if using digest; avoid control prefixes.

## Lifecycle Failure Modes

- File write succeeds but DB create fails → retryable intent, no row to revoke.
- DB row exists but row ID/file write incomplete → recover by exact lookup digest, then verify owner/status before adopting.
- Patch removed after row/file persist → original sees a normal row; C1 limits use to loopback but cannot revoke without a running privileged actor.
- Planned rollback/uninstall → revoke/rotate owned row and delete token state after transaction checkpoint.
- Owner deletes/rebinds/disables row while original binary runs → patched status detects exact ID mismatch/disabled state; do not resurrect automatically.
- Partial cache invalidation → update/delete paths already invalidate by stored `apiKey.Key`; custom internal operations must do the same and publish L1 invalidation.
- Audit/log paths → ensure raw token never enters request bodies, URL query, process args, or error strings. Marker values with controls add log-injection/UI escaping risk.

## Planning Questions and Confirmed Resolution

1. K6 accepts original-runtime lookup visibility and local-source usability. K4 still requires source-IP/proxy validation before activation.
2. Normal service-level registration retains existing public validation. No control-marker or validation exception is approved.
3. Original administrators retain native management behavior. K7 requires conflict detection on reapplication, not silent reversal of external changes.
4. Exact loopback addresses remain deployment evidence, not an unresolved risk preference. Do not broaden the whitelist to bypass a failed check.
5. Original administrator DTO lookup visibility is accepted by K6. K1 still excludes the raw private token under every runtime state.

This confirmation resolves the Key risk decision. It does not approve product implementation, deployment, or any unexecuted security claim.

## Minimal Owner-File Map for a Future Patch

| Owner | Responsibility |
|---|---|
| `internal/service/api_key_service.go` | Optional narrow provisioning/recovery helpers; keep normal `Create` semantics unless marker exception approved. |
| `internal/server/middleware/api_key_auth.go` | Map raw token to stored lookup only after extraction and before `GetByKey`; reject external stored-lookup replay; preserve normal checks. |
| `internal/server/middleware/api_key_auth_google.go` | Equivalent external marker/digest replay guard if non-HTTP marker is approved; required for query-path completeness. |
| `internal/handler/api_key_handler.go` | Exact owned-row update/delete guard without prefix-wide regressions. |
| `internal/handler/admin/apikey_handler.go` | Exact owned-row rebind/reset guard. |
| `internal/handler/dto/mappers.go` | Exact owned-row key redaction in all DTO paths. |
| `internal/server/middleware/audit_log.go` / `internal/service/audit_log.go` | Add route omission/redaction only if new eval APIs carry secret-bearing bodies. |
| `internal/service/gateway_usage_billing.go` | No functional change expected; regression-test invalidation sees stored lookup key. |
| `internal/handler/openai_gateway_handler.go` | Preserve `conn.Key` as stored lookup for WS-turn refresh. |
| installer/status implementation | Verify row identity, token file mode, binary hash, and post-recreation patch absence before reapplying. |

## Regression Tests Required

| Test | Required assertion |
|---|---|
| Raw-token mapping | External raw token maps to digest only for configured exact row; `APIKey.Key` in context is digest, never raw. |
| Digest replay while patched | External digest is rejected before `GetByKey`; local worker path still authenticates. |
| Original residual contract | Original accepts digest only from configured loopback ACL; remote-source request fails before user/group checks. Fixture must cover IPv4 and any approved IPv6 source. |
| Google alias path | Marker/digest checks cover Bearer, `x-api-key`, `x-goog-api-key`, and allowed query `key`; deprecated main query remains rejected. |
| DTO/browser privacy | User list/detail, admin user/group key lists, nested usage, errors, and audit output contain no raw token; owned stored value is redacted where required. |
| Normal lifecycle | Owner status, quotas, rate windows, subscriptions, concurrency, group permission, and usage attribution remain identical to an ordinary Key. |
| Cache consistency | Create/update/delete/quota-exhaustion invalidation use the stored lookup value; raw token never becomes an auth cache key. |
| Recovery | Existing row + missing intent, intent + conflict, externally disabled row, owner mismatch, and duplicate `ExistsByKey` all fail closed with repair status—not silent re-enable/recreate. |
| Identity guards | Same-prefix ordinary keys can be updated/deleted normally; only exact owned row ID is protected. |
| Marker feasibility (if revisited) | DB round-trip, DTO escaping, audit redaction, log prefix safety, main auth, Google query auth, and WS refresh all explicitly enumerated; no claim of “unpresentable” without every path passing. |
| Concurrency | Parallel evaluations across groups preserve immutable Key/Group copies, and no shared cached row is rebound. |

## Final Recommendation

C1 is a conditional improvement, not a complete privacy guarantee. Verify source-IP/proxy behavior before relying on loopback restrictions.

Keep patched replay rejection and DTO redaction. The user explicitly accepted original-runtime lookup visibility and local-source usability as K6.

Planned rollback/uninstall can revoke owned resources. Unplanned recreation cannot guarantee revocation without another persistent enforcement mechanism.

Reject control-character and oversized markers as general solutions. They conflict with query extraction, existing storage, or stable public credential constraints.

Stronger original-runtime derivative secrecy is outside the confirmed K6 boundary. The normal-row scheme does not claim that stronger guarantee.
