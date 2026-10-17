# Saved full-race source classification

The saved final command remains failed: **27 race warnings and 105 top-level failing tests**. No compilation failure remains. This analysis reads saved logs and pinned sources only. It does not run a pristine baseline or repair unrelated tests.

Evidence: `dist/sub2api-build/full-race-source-classification.json`. Each race entry records every repository frame, exact source line and pinned-byte comparison. Each failing test records its source definition, native-byte comparison, immediate failure output and historical-name match.

All 105 top-level failures report only `race detected during execution of test` in their immediate indented failure output. No independent assertion mismatch or panic is reported there. All 105 test sources match pinned upstream bytes. None belong to the maintained overlay function-name set. This checks actual function names, rather than the misleading `TestEval` prefix, which also matches native scheduling-threshold tests.

| Direct race group | Warnings | Source and observed conflict |
| --- | --- | --- |
| Administrator Grok import probe | 4 | `grok_import_probe_test.go:277` reads a plain `bytes.Buffer` while a native scheduler worker logs to it through `grok_import_probe.go:144`. |
| Administrator Grok quota/model refresh | 2 | `grok_oauth_handler_test.go:42,162` writes/reads the unsynchronized transport call slice while native `grok_observed_models.go:147` refreshes models asynchronously. |
| Channel monitor capture | 4 | `channel_monitor_checker_body_test.go:72–80` concurrently writes fields on its reused HTTP capture handler. |
| Content moderation runtime snapshot | 1 | `content_moderation.go:1559` reads the runtime cache TTL while `content_moderation_cyber_test.go:264` writes it during background refresh. |
| Grok free quota cache | 1 | `grok_free_quota_gate.go:264` accesses the cache while `grok_free_quota_gate_test.go:241` replaces the global `sync.Map`. |
| Dashboard aggregation cleanup | 2 | `usage_cleanup_service_test.go:68` increments an unsynchronized counter while tests at lines 583/611 inspect it during native dashboard background recomputation. |
| Email verification cache | 2 | `auth_service_email_bind_test.go:1173–1174` increments/reads an unsynchronized verification-attempt counter in concurrent alias binding tests. |
| OpenAI compatibility native tests | 9 | Parallel `openai_compat_model_test.go` tests call global `gin.SetMode` while other tests construct contexts or set the same global. |
| Responses buffered/streaming native tests | 2 | Parallel `gateway_forward_as_responses_test.go:502,533` tests call global `gin.SetMode` concurrently. |

Direct race-frame sources are byte-identical to the pinned source except for two scheduler owner files. Their three participating functions are also byte-identical: `selectByLoadBalance`, `listSchedulableAccounts` and `filterGrokFreeQuotaAccountsForOpenAI`. The saved stack traverses ordinary Grok scheduling into the unchanged native soft-quota cache. The evaluation target path is not present in this race stack.

The historical reviewed log reported 18 warnings and 58 failed tests. Forty-one names in the final 105 also occur in that saved historical failure set. Different scheduling and race-contamination windows can fail unrelated parallel tests. Source attribution supports the native ownership classification, but does not prove identical execution or causality for each additional failed test. The count change is reported explicitly. The full-race gate remains a blocker and must not be suppressed or described as passing.

Eight direct race groups have frame-file evidence in the saved historical log, including native Grok and global Gin conflicts. The final buffered/streaming Responses Gin race frames are absent from that historical log, so their historical recurrence is not claimed. Their test sources still match the pinned source. Local focused delivery tests pass independently. The full-scope checker should review this classification and the retained failed gate before handoff.
