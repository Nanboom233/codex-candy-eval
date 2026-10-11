# JSON Persistence Guidelines

> This plugin stores history in a JSON file. It does not use a database or ORM.

---

## Overview

<!--
Document your project's database conventions here.

Questions to answer:
- What ORM/query library do you use?
- How are migrations managed?
- What are the naming conventions for tables/columns?
- How do you handle transactions?
-->

The legacy filename remains for spec-link compatibility. The storage contract lives in [state.go](../../../internal/plugin/state.go), not in a database layer.

`statePath` defaults to `plugins/cpa-codex-candy-eval-state.json`, relative to the host working directory. `persistedState` stores four credential-ID maps: candy, fingerprint, ModelTrace, and pelican history.

Running tasks, configured credential metadata, and inspection tracking are not fields of this snapshot. Do not assume every in-memory map survives restart.

Evidence: `persistedState`, `currentState`, and `loadState` in [state.go](../../../internal/plugin/state.go).

---

## Snapshot Access and Writes

<!-- How should queries be written? Batch operations? -->

Hold the package mutex `mu` when accessing shared histories during normal execution. Call `saveStateLocked` only with that lock held.

`loadState` acquires its own lock and reads once per process. A missing file means no saved history. Other read or decode errors set `stateLoadError` and `storageError`.

`saveStateLocked` serializes all four histories together. It writes `statePath + ".tmp"` with mode `0600`, then renames that file over `statePath`.

The temporary file sits beside the destination. Preserve this placement and the read-error guard. Atomic replacement is not a database transaction or an fsync durability guarantee.

Existing runner pattern from [runCandy](../../../internal/plugin/candy.go):

```go
mu.Lock()
candyResults[id] = appendRecent(candyResults[id], r)
candyRunning[id].Done++
_ = saveStateLocked()
mu.Unlock()
```

Runner save failures update `storageError`. They do not undo completed evaluations. History removal must restore the previous maps when saving fails.

Evidence: `saveStateLocked` and `clearHistoryResponse` in [state.go](../../../internal/plugin/state.go).

---

## Snapshot Compatibility

<!-- How to create and run migrations -->

There is no migration runner. JSON decoding accepts older snapshots that omit newer history maps. Keep existing JSON field names compatible unless a separate migration is planned.

`loadState` replaces only maps present in the snapshot, then trims each history. A JSON `null` snapshot is invalid and must not trigger replacement of the original file.

Do not treat the plugin registration `schemaVersion` as a persisted-file migration number. It appears in host registration, not in `persistedState`.

Evidence:

- [Load and save tests](../../../internal/plugin/plugin_test.go): `TestStatePersistence` and `TestUnreadableStateIsPreserved`.
- [Host registration](../../../internal/plugin/plugin.go): `schemaVersion`.

---

## Field Names and History Limits

<!-- Table names, column names, index names -->

Keep the JSON keys `results`, `fingerprints`, `modeltraces`, and `pelicans`. Each map uses the CPA credential ID as its key.

Use `appendRecent` to append records. Use `keepRecent` when you restore snapshots. `historyLimit` retains the latest 20 records per credential and test type.

Do not key saved results by display name, email, or a masked API-key preview. Credential identity and display metadata have separate roles.

Evidence: [history helpers](../../../internal/plugin/state.go), [credential types](../../../internal/plugin/credentials.go), and `TestRunAllKeepsRecentHistory` in [plugin tests](../../../internal/plugin/plugin_test.go).

---

## Common Mistakes

<!-- Database-related mistakes your team has made -->

- Do not overwrite malformed or unreadable history with a new empty snapshot. `stateLoadError` deliberately blocks saves.
- Do not clear running-task maps when removing saved history. Existing tests require in-flight work to remain intact.
- Do not report history removal as successful when persistence failed. Restore the previous histories. Return an error.
- Do not introduce SQL, an ORM, or a migration dependency for this file-backed store without a separate task.
- Do not create deployment state files or change production history during tests. Use `setupTest` and its temporary `statePath`.

Verification: `TestStatePersistence`, `TestUnreadableStateIsPreserved`, and `TestFingerprintValidationAndStorageFailure` in [plugin tests](../../../internal/plugin/plugin_test.go) and [fingerprint tests](../../../internal/plugin/fingerprint_test.go).
