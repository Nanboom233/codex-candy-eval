# Research: Core C18-C22 Evidence

- Query: Compare Candy batch execution and state/history persistence in the current `internal/eval` package with the original `internal/plugin` implementation, and identify which active-intent, validation, file-store, snapshot, detail, clear, and flush additions are required by the approved administrator channel.
- Scope: internal source
- Date: 2026-10-15

## Findings

### C18: Candy batches, settlement, and result persistence

The current manager validates Candy batches in `internal/eval/runs.go:141-162`, defaults `Runs` to one, and caps it at `CandyMaxRuns` (10). `:412-503` executes each iteration serially for one account, creates a distinct run ID (`intent.ID:iteration`), passes the manager-owned context into `RunOne`, appends each result, updates cumulative progress, and releases the account reservation only after the worker exits (`:413-420`). The current `internal/eval/candy.go:30-42` preserves the original one-request result, numeric usage, duration, 4,000-rune answer display cap, standalone-21 verdict, and error/cancellation result status.

The original `internal/plugin/candy.go:54-72` clamps the requested batch to one through ten and starts one goroutine per selected credential. `:75-98` runs iterations serially, checks `quiescing` between requests, appends each completed result, increments the in-memory progress counter, and calls `_ = saveStateLocked()` after each result. `:100-111` preserves usage and error text even when the request fails; a successful wrong answer remains a normal result with `OK=false` rather than a transport failure.

The material persistence difference is failure handling. Original workers ignore the return value of `saveStateLocked`; `internal/plugin/state.go:87-106` records a human-readable `storageError` on failure and leaves the in-memory result/history present. The worker does not create an active intent, claim that the failed write was durable, or change the result into a synthetic record. The current worker treats `m.storageErr != nil` as a terminal condition at `internal/eval/runs.go:472-500`, which is a new manager policy and can stop later Candy iterations after a save failure.

Recommendation: retain the cap, per-account serial order, distinct per-iteration identifiers, error/usage/wrong-answer semantics, and final busy cleanup. Persist each finished result as one history snapshot and expose save failures explicitly; never report a failed `Save` as durable success. Do not add a second batch bookkeeping store or claim that the original plugin had durable active/completed intents. Whether a storage failure stops subsequent paid iterations is a separate runtime policy; it is not inherited from the original worker, which continued in memory while exposing `storageError`.

### C19: Active intents, restart recovery, and fabricated interruption history

The current state schema adds `ActiveIntent` and an optional `active` map at `internal/eval/state.go:17-38`. `internal/eval/runs.go:100-138` loads the map on startup, appends a synthetic `interrupted` result for every active intent through `syntheticResult`, deletes the active entries, and saves the recovery snapshot before returning. `:268-331` writes the full active intent before launching any worker, and `:516-521` exposes untracked persisted intents as `recovery_blocked` snapshots. This is a durable task/recovery subsystem, not a history feature.

The original `internal/plugin/state.go:13-18,27-66` persisted only the four completed-history maps. There is no active-run schema, pre-ack intent write, auto-resume, restart reconciliation, or fabricated `interrupted` record. A process restart simply loses in-flight memory state; results already saved after prior Candy iterations remain in history. The original `internal/plugin/candy.go:75-98` confirms that worker lifecycle and persisted history were separate.

Recommendation: keep runtime tasks and progress in memory and remove the persisted `ActiveIntent`/restart-recovery path. Do not auto-resume work and do not create an `interrupted` history row for work that never produced a model result. A restart may discard an in-flight task; already completed records remain subject to the normal history snapshot. This is a deliberate breaking change: do not add a migration layer, silently import CPA state, or delete/repair existing files. If a future product requirement asks for crash recovery, it must be a new, explicitly approved state contract rather than an assumption derived from history retention.

### C20: State validation, corruption protection, and retention

The current `internal/eval/state.go:14-15,31-38` introduces `schema_version=1`, typed maps keyed by numeric account IDs, and an optional active map. `:163-193` fills nil maps and trims every history to the latest 20 records. `:195-287` additionally validates positive/matching account and group IDs, nonempty run IDs, a fixed status enum, model/effort normalization, collection counters, active-account uniqueness, active counters, and full `StartRequest` options. Any failure blocks the manager and prevents writes.

The original `internal/plugin/state.go:27-66` only distinguishes missing, unreadable, malformed, or JSON-null state. It preserves the source bytes by setting `stateLoadError`/`storageError`, refuses future saves while that latch is set (`:87-96`), accepts absent history maps, and trims each map to `historyLimit=20` (`:68-79`). It does not validate historical model names, current request options, status enums, counter relationships, group identity, or an active-run schema. Thus strict validation can reject otherwise readable historical records when the evaluation result/status vocabulary evolves, and it couples storage recovery to live request validation.

Recommendation: retain the observable guarantees that matter to history: malformed/unreadable/null input never overwrites the original file, storage remains visibly blocked until repaired, missing maps normalize to empty collections, and each account/type keeps the latest 20 records. Use one fixed new format for the new administrator channel and fail closed on an unsupported format; do not build compatibility or migration branches. Remove live `NormalizeModel`/`StartRequest` validation from historical records and remove all active-intent validators together with C19. Any additional identity checks must be justified by a concrete handler contract, not by the existence of the history file alone. Never delete or rewrite an existing source file as part of loading or format transition.

### C21: Atomic history file store and path hardening

The current `internal/eval/state.go:46-85` performs component-by-component `Lstat` checks and rejects symlinks, non-directory parents, and nonregular destinations. `:113-156` validates state, creates missing parent directories with mode `0700`, applies `0700` to the final parent, creates a unique same-directory temporary file, applies mode `0600`, writes and closes it, rechecks the temporary/destination paths, and atomically renames the temporary file. The code explicitly does not coordinate multiple processes and does not fsync the directory.

The original `internal/plugin/state.go:86-107` serializes all histories under one mutex, writes JSON to `statePath + ".tmp"` with mode `0600`, removes the temporary path on failure, and renames it over the destination. The host already supplies the `plugins/` directory, so it does not create or chmod parent directories and has no symlink-component defense. Both implementations avoid exposing a partially written JSON snapshot; neither proves power-loss durability or multi-process ownership.

Recommendation: retain one owner, one same-directory temporary write, mode `0600`, close-before-rename, cleanup of only the temporary file owned by that save, and atomic replacement. The native runtime should choose and create its dedicated private state directory; do not make the evaluation core a general filesystem policy engine. Component symlink/regular-file checks and parent `0700` enforcement are defensible only when the runtime cannot otherwise guarantee ownership; they are not required by the original history semantics and should not be used to introduce alternate paths, shared locks, fsync claims, or process coordination. A failed validation or rename must leave the destination bytes untouched and must be reported as a save failure.

### C22: Snapshots, detail lookup, clear, and flush

The current manager deep-copies pointer/slice-rich result data in `internal/eval/state.go:289-341`. `internal/eval/runs.go:505-525` returns a cloned state snapshot and active progress; `:527-564` retrieves a record by account and ID/run ID for all four result types; `:566-610` clears one type or all histories by saving a candidate snapshot before publishing it; and `:612-635` adds an explicit `Flush` operation that retries the full snapshot and clears a recoverable storage error. The web page depends on these user-facing behaviors: `internal/eval/web/app.js:51-69,459-469` maps detail/clear requests, `internal/eval/web/components.js:199-209` loads full Fingerprint/ModelTrace detail records, and `internal/eval/web/pelican.js:10-23,75-80` loads Pelican HTML separately from list summaries. Candy has no detail fetch: `internal/eval/web/candy.js:25-55,64-68` renders the complete result from `/state` and expands the answer inline.

The original `internal/plugin/state.go:109-131` provides detail lookup for Fingerprint, ModelTrace, and Pelican under the package mutex and returns a fixed not-found error; it has no Candy detail branch. `:134-158` clears a selected type or all types, saves the complete history snapshot, and restores the previous in-memory maps if saving fails. Candy history is already present in the normal state response; it does not need a new record endpoint or `ID||RunID` lookup alias. There is no original `Flush` maintenance endpoint; ordinary workers attempt `saveStateLocked` and expose `storageError` when persistence fails.

Recommendation: retain the actual page controls: immutable state snapshots for safe reads, Fingerprint/ModelTrace/Pelican detail lookup, per-type and all-history clear, not-found behavior, and rollback/no-visible-clear on failed persistence. Candy remains state-list-only with inline answer expansion; remove its redundant record lookup unless a separate API requirement is approved. Deep copies are useful for manager-owned mutable slices and should remain only at the handler boundary that needs them. Keep a visible storage error and never clear it merely because an in-memory operation succeeded. Remove `Flush`, `recovery_blocked`, and other maintenance-only state-machine surface unless a concrete native lifecycle caller is approved; no production caller for `Flush` was found. Normal task completion and shutdown do not imply that a disk write succeeded. Clear and detail operations must preserve unrelated histories and must not delete existing files outside the selected history snapshot.

## Files Found

- `internal/plugin/candy.go` — original Candy validation, serial batches, result/error fields, progress, and per-result save calls.
- `internal/plugin/state.go` — original history-only persistence, corruption latch, latest-20 retention, atomic temp replacement, detail lookup, and clear rollback.
- `internal/plugin/plugin.go` — original state load/registration and management route ownership.
- `internal/eval/candy.go` — current Candy algorithm/result behavior.
- `internal/eval/runs.go` — current batch worker, active intent persistence, recovery, snapshots, record lookup, clear, flush, and shutdown state.
- `internal/eval/state.go` — current schema, strict validators, retention, file store, cloning, and atomic save implementation.
- `internal/eval/web/app.js`, `components.js`, and `pelican.js` — current detail, clear-history, summary, polling, and full HTML consumers.
- `tests/overlays/codexeval/runs_test.go` — current tests for pre-ack active intents, save failures, Candy batches, unavailable synthetic records, record lookup, clear, and flush.
- `tests/overlays/codexeval/state_test.go` — current tests for symlink/path policy, strict schema/identity validation, corruption protection, recovery records, retention, and temporary ownership.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/patch-scope-core-tooling-inventory.json` — inventory entries C18-C22 and their current/new behavior comparison.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/sub2api-current-modifications.md` — concise list of current C18-C22 additions.

## Related Specs

- `.trellis/spec/backend/run-state.md` — current durable-manager contract; its active-intent, recovery, and strict validation requirements are the additions under review here, not evidence that the original plugin had them.
- `.trellis/spec/backend/evaluation-core.md` — Candy result, cancellation, progress, and persistence boundary.
- `.trellis/spec/backend/minimal-linux-delivery.md` — runtime-owned data directory and delivery boundaries relevant to C21.
- `.trellis/spec/backend/patch-native-installation.md` — native administrator lifecycle and file ownership boundaries.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md`, `design.md`, and `implement.md` — current administrator channel, breaking-change, no-migration, and no-product-edit planning constraints.

## Caveats / Not Found

- The current test overlay intentionally asserts active-intent/recovery and strict validators; those tests document the current implementation, not user approval. Removing C19/C20 behavior requires revising those tests during the authorized implementation phase.
- No product code, build, test, deployment, cleanup, or live provider request was run for this evidence.
- The original save path uses a fixed `.tmp` name under a host-owned plugin directory; the native runtime may still need a private directory and unique temporary name for concurrent/owned writes. That is an ownership decision, not evidence for a full filesystem security framework.
- This research does not decide C23 manager-wide cancellation/shutdown; it only identifies which persistence and maintenance claims C18-C22 should not make. Existing in-memory cancellation remains covered by the earlier approved core/runtime decisions.
- Breaking-change scope forbids automatic migration and implicit deletion, but the final implementation still needs one explicit destination path for new history. Existing unrelated CPA files should remain untouched.
