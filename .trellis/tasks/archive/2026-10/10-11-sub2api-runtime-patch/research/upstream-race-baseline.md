# Research: upstream `-race` baseline blocker (pristine sub2api v0.2.15 `internal/service`)

- **Query**: Identify actual race stack pairs from the pristine `go test -tags=unit -race ./internal/service -count=1` run (exit 1, 10 race warnings, 54 top-level fails), classify fixture-only vs production, and scope the minimal test-fixture repair to let a full race gate pass.
- **Scope**: internal (dist/sub2api-v0.2.15/backend @ f2669c8, git-clean) + captured log
- **Date**: 2026-10-11

## Confirmed facts

Baseline command on the pristine tree: `cd dist/sub2api-v0.2.15/backend && go test -tags=unit -race ./internal/service -count=1` → `FAIL ... 261.677s`, exit 1. Captured log contains **10 `WARNING: DATA RACE` blocks** which resolve to **5 distinct race sources**, all rooted in `_test.go` fixtures. No production (`*.go` non-test) symbol writes unsynchronized data of its own; every race pair has at least one endpoint in test code, and each production endpoint is either (a) reading a test-stub-owned field, or (b) a gin v1.9.1 library global written only from tests.

### Race 1 — `openAICaptureHandler` fields (4 warnings: log lines 438/472/506/540)

- **Write/Write & Write/Read pairs**: `(*openAICaptureHandler).ServeHTTP` writes `lastHeaders` (:72), `lastPath` (:73), `lastBody` (:77), reads/writes `status` (:79-80) in `internal/service/channel_monitor_checker_body_test.go`.
- **Owning test**: `TestRunCheck_QuotaProbeAttachesSnapshotToPrimaryRowOnly` (`channel_monitor_quota_mode_test.go:161`). Monitor has `PrimaryModel` + `ExtraModels`; production `runChecksConcurrent` (`channel_monitor_service.go:713-722`) deliberately fans out one `runCheckForModel` per model via `errgroup`, so two concurrent `httptest` server goroutines hit the same handler. **Production concurrency is intended** (comment at :696 "对 primary + extra 模型并发执行检测"); the handler fields are plain fields — **fixture-only race**.
- Same struct also used in `channel_monitor_checker_body_test.go` tests with `ExtraModels`/multi-request paths; `captureHandler` (:27) has the identical pattern and should be hardened at the same time (same file, no observed warning but same exposure).

### Race 2 — `openAIGrokFreeQuotaGateCache` global reassignment (1 warning: log line 3358)

- **Write**: `TestGrokFreeQuotaGateIsSchedulerOnlyAdminPathUnfiltered` `grok_free_quota_gate_test.go:241` executes `openaiGrokFreeQuotaGateCache = sync.Map{}` — assigns the package-level variable itself.
- **Read**: leftover background refresh goroutine `scheduleGrokFreeQuotaStatsRefresh.func1` (`grok_free_quota_gate.go:264`, `sync/atomic.LoadUint32` inside `sync.Map` internals), spawned by `TestOpenAIAccountSchedulerLoadBalanceAppliesGrokFreeQuotaGate` (`grok_free_quota_gate_test.go:222`, `t.Parallel` at :188) via `filterGrokFreeQuotaAccountsForOpenAI` → `filterGrokFreeQuotaAccountsCore(..., &openaiGrokFreeQuotaGateCache, ...)` (`openai_gateway_scheduling.go:1698`). The background goroutine dereferences `&global` after the test reassigns the variable → race on the variable slot. **Test-global ownership bug, fixture-side fix**; production's `*sync.Map` use is internally consistent.
- Identical resets at `grok_free_quota_gate_test.go:71, 99, 139, 201` (5 total). `freeQuotaRefreshInFlight` (`grok_free_quota_gate.go:131`) is keyed by the `*sync.Map` pointer, so in-place clearing keeps coalescing semantics.

### Race 3 — `dashboardRepoStub.recomputeCalls` (2 warnings: log lines 12932/12979)

- **Write**: `(*dashboardRepoStub).RecomputeRange` `recomputeCalls++` (`usage_cleanup_service_test.go:68`), called from production `DashboardAggregationService.TriggerRecomputeRange` background goroutine (`dashboard_aggregation_service.go:169,173,212` — async by design).
- **Read**: `require.Eventually` closures at `usage_cleanup_service_test.go:583` (`TestUsageCleanupServiceExecuteTaskDashboardRecomputeError`) and `:611` (`...RecomputeSuccess`). **Fixture-only**: the stub lacks the `mu` that its sibling `cleanupRepoStub` already uses.

### Race 4 — `emailBindCacheStub.data.Attempts` (1 warning: log line 13072)

- **Write**: `(*emailBindCacheStub).IncrVerificationCodeAttempts` `s.data.Attempts++` (`auth_service_email_bind_test.go:1173`).
- **Read**: production `verifyCodeWithAttempts` `data.Attempts >= maxVerifyCodeAttempts` (`email_service.go:410`) on the `*VerificationCodeData` pointer returned by the stub's `get` — the stub returns the same shared pointer to both goroutines.
- **Owning test**: `TestAuthServiceBindEmailIdentity_AllowsOnlyOneConcurrentAliasVariant` (`auth_service_email_bind_test.go:291-296`) intentionally races two `BindEmailIdentity` calls on the alias-variant path; production `AuthService.VerifyOAuthEmailCode` → `EmailService.VerifyCode` chain is legitimate concurrent use. **Fixture-only** (a real `EmailCache` returns per-call copies); the stub needs a mutex and must return a copy from `get`.

### Race 5 — gin v1.9.1 globals `ginMode`/`modeName` (2 warnings: log lines 13349/13407)

- **Writes**: `gin.SetMode` → `mode.go:72` (`ginMode`) and `mode.go:77` (`modeName`) — plain non-atomic globals in gin v1.9.1.
- **Reads**: `gin.IsDebugging` (`debug.go:20`) via `debugPrint` ← `gin.New` ← `gin.CreateTestContext`.
- **Owning tests**: `TestForwardAsAnthropic_StoresStreamingResponseIDWithoutUsage` (`openai_compat_model_test.go:1486`, `t.Parallel` :1485) racing `TestForwardAsAnthropic_DisablesAPIKeyContinuationWhenUpstreamRequiresWebSocketV2` (:937 SetMode, :969 CreateTestContext). **Any** two parallel tests in the package can trip this: `gin.SetMode` appears ~700+ times across ~179 `internal/service/*_test.go` files (verified by grep), always `gin.SetMode(gin.TestMode)`. **Fixture-only** (library-global written exclusively from tests under `go test`).

### Why 54 top-level fails from 10 warnings

Confirmed by log pattern: once the race detector fires, (a) the report is emitted asynchronously, so `testing.go:1867 "race detected during execution"` attaches to whichever test happens to be running at flush time — this is how pure functions (`TestUnionFloat`, `TestPcParseFloat`, `TestExtractUpstreamModelIDs`, `TestInferGoogleOneTier`, `TestJoinTypes`, `TestValidateProviderRequest`, etc. — all 0.00s) "fail with race"; (b) leaked background goroutines (grok refresh goroutine, dashboard recompute goroutine, httptest conn goroutines) outlive their parent test, so the *other* half of a pair lands inside an unrelated parallel test; (c) parallel tests in the same package share process-wide globals (gin mode), so one racy pair contaminates both participants and any observer. Therefore fail counts (baseline 54 vs implementer's ~96/~120 varying runs) are **not comparable run-to-run** and individual "race detected" fails are not evidence the named test is itself racy — only the 10 stack pairs are authoritative.

## Minimal candidate repair scope (test-fixture-only; no production changes, no dropped tests, no `-race` removal, no suite serialization, no timing hacks)

All edits are in `dist/sub2api-v0.2.15/backend/internal/service/` `_test.go` files only. The implementer should land this as a **new patch** (e.g. `patches/sub2api/v0.2.15/000-test-fixture-race-fixes.patch`) applied to the pristine tree — see approval note below.

| # | File | Symbols / lines | Candidate fix |
|---|---|---|---|
| 1 | `channel_monitor_checker_body_test.go` | `openAICaptureHandler` (:62-70), `ServeHTTP` (:71-95+), all field reads in tests; `captureHandler` (:27-54) same pattern | Add `sync.Mutex` (or `atomic` fields) guarding `lastBody/lastHeaders/lastPath/status/rawResponse/responsesLeadingReasoning`; expose locked getters or lock around assertions. Preserve the Chinese doc comment at :26 and the field semantics — do not change `answerFromOpenAIRequest` or response logic. |
| 2 | `grok_free_quota_gate_test.go` | lines 71, 99, 139, 201, 241 (`openaiGrokFreeQuotaGateCache = sync.Map{}`) | Replace variable reassignment with in-place clear: `openaiGrokFreeQuotaGateCache.Clear()` (Go 1.23+, go.mod says `go 1.27.2`) **or** Range+Delete. Do NOT keep `= sync.Map{}` — production passes `&openaiGrokFreeQuotaGateCache`. Optionally also drain stale `freeQuotaRefreshInFlight` entries; keep the comments at :190-194 (scheduler-only gate documentation) and :229-231 (admin-path rationale) intact. |
| 3 | `usage_cleanup_service_test.go` | `dashboardRepoStub` (:57-69), readers :583, :611 | Make `recomputeCalls` `atomic.Int64` (or add `mu sync.Mutex` like sibling `cleanupRepoStub` already has); read via `.Load()` in the `require.Eventually` closures. Assertions stay `== 1`. |
| 4 | `auth_service_email_bind_test.go` | `emailBindCacheStub` (`IncrVerificationCodeAttempts` :1172-1178 and the `get` used by `VerifyCode`) | Add `mu sync.Mutex` to the stub; lock all accessors mutating `s.data`; have the getter return a **copy** of `*VerificationCodeData` so production `email_service.go:410` never reads a field another goroutine increments. Preserve the concurrency-test intent (the `start` channel + two goroutines at :287-296) — do not serialize it. |
| 5 | ~179 `*_test.go` files | every `gin.SetMode(gin.TestMode)` call | Remove all per-test `gin.SetMode(gin.TestMode)` calls and add **one** `func init() { gin.SetMode(gin.TestMode) }` (or `TestMain`) in a single new `_test.go` file so the only write happens before any test goroutine exists; afterwards all accesses are read-only → no race. Do not substitute a mutex around `SetMode` — `IsDebugging`/`CreateTestContext` reads are not under it, so a lock cannot fix this. (gin defaults `modeName` via `init()` before testing flags register, so relying on `test.v` auto-detection is not guaranteed — the explicit one-shot init is required.) This is a large mechanical but semantically no-op diff. |

Targeted verification commands for the implementer (run per-race, not the full suite):

```sh
cd dist/sub2api-v0.2.15/backend
go test -tags=unit -race ./internal/service -count=1 -run 'TestRunCheck_QuotaProbeAttachesSnapshotToPrimaryRowOnly|TestRunCheck_QuotaProbeQuotaFailureKeepsProbeStatus'          # race 1
go test -tags=unit -race ./internal/service -count=1 -run 'TestGrokFreeQuotaGate|TestOpenAIAccountSchedulerLoadBalanceAppliesGrokFreeQuotaGate|TestFilterGrokFreeQuotaAccounts'  # race 2 (run repeatedly, e.g. -count=5, timing-dependent)
go test -tags=unit -race ./internal/service -count=1 -run 'TestUsageCleanupServiceExecuteTaskDashboardRecompute'                                                              # race 3
go test -tags=unit -race ./internal/service -count=1 -run 'TestAuthServiceBindEmailIdentity_AllowsOnlyOneConcurrentAliasVariant'                                              # race 4
go test -tags=unit -race ./internal/service -count=1 -run 'TestForwardAsAnthropic_'                                                                                            # race 5
```

Expected assertions (must remain unchanged): quota snapshot on primary row only (`results[1].Quota == nil`, history[0].Quota set); admin path unfiltered (`filterGrokFreeQuotaAccounts` returns len 0 for over-gate FREE account via Eventually); `dashboardRepo.recomputeCalls == 1` via Eventually in both cleanup tests; exactly one of the two concurrent `BindEmailIdentity` calls succeeds; streaming response ID stored without usage / WebSocket-required continuation disabled.

## Inferences (not directly proven by the log)

- Warning count 10 = 4 (race 1) + 1 (race 2) + 2 (race 3) + 1 (race 4) + 2 (race 5); all confirmed by stack pairs above.
- Additional `sync.Map{}`-style global resets or shared-pointer stubs may exist elsewhere in the package (log shows "further late usage/timing/testglobals" noise); the implementer should grep `= sync.Map{}` and stub counters package-wide after fixing these five and re-run the gate once — no full-suite rerun was performed for this report.
- The Gin mode repair may touch many test files. Main-session search found 1,109 `gin.SetMode` matches in service test sources, so the report's earlier approximate call count is an underestimate.
- Do not remove intentional mode-switching assertions indiscriminately. Inspect non-TestMode uses and package initialization before choosing the mechanical repair.
- A Gin upgrade is a dependency change outside the current approved repair scope. Do not use it to bypass the verification blocker.

## Confirmed Scope Decision

The user selected `Keep the original patch scope` after reviewing the verified baseline blocker.

Do not create `000-test-fixture-race-fixes.patch`, alter unrelated upstream fixtures, upgrade Gin, disable race detection, or skip tests.

Continue the approved feature implementation. The full upstream race acceptance remains blocked and must not be reported as passed.

The repair table above is research evidence only. It does not authorize those changes.

## Caveats / Not Found

- Log read selectively; the ~54 top-level "race detected" fails are attribution noise, not 54 distinct races — only the 10 stack pairs are real.
- Did not run the suite again (per instructions); per-race `-run` verification is left to the implementer.
- `captureHandler` (`channel_monitor_checker_body_test.go:27`) shares the race-1 pattern but produced no warning in this run — fix it together with `openAICaptureHandler` since the same production fan-out applies.
