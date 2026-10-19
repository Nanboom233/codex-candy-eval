# Independent Check

## Findings (Fixed)

- File: `patches/0006-evaluation-ui.patch` (`web/ui.html`)
- Issue: The Pelican page said each account sent one request. The shared runner can send up to three attempts for one logical test.
- Fix: The implementation owner removed the fixed request-count caption. Candy and Pelican pages have no other fixed request-count caption. No replacement copy or test was added.

## Findings (Not Fixed)

None.

## Review

The implementation matches the approved native contracts:

- Only an absent `schema_version` selects history migration. Loading retains Candy and Pelican records, drops retired histories, and immediately saves schema version 1 through the existing atomic saver. Explicit `null` and unsupported markers protect the source file.
- The shared runner builds each request once, makes at most three attempts, waits two then four seconds, and accumulates attempts and observed usage. Permanent statuses and unmarked setup or plugin errors remain final.
- A known HTTP status remains available after response-body read failure. Only the actual native transport boundary marks a pre-response error as retryable.
- Inspection checks participation before each send, including after backoff. One logical Candy result counts once. Intentional opt-out adds no synthetic result.
- The detector, native Settings page, account controls, and four-map-to-three-map state change follow the approved design. Tests use the existing focused owners, with one new owner for the shared runner.
- The merge has no unresolved paths and uses the approved `7a49068` upstream parent. Later upstream Pelican changes remain excluded.

The final history-marker correction is covered by the existing runtime persistence fixture. The task audit's test boundary remains intact. No full-suite, race, coverage, benchmark, architecture, or hash gate was added.

## Verification

- Lint: Pass. The final `--check` run passed `go vet` for the changed backend packages and affected server packages.
- Type check and build: Pass. The final `--check` run compiled the patched frontend and Linux/amd64 server. After the caption removal, the existing UI harness passed and the build entry regenerated both server and installer outputs.
- Tests: Pass. The final `--check` run passed the focused Go selector, UI harness, and host-platform installer tests. The final markup-only repair reused the unchanged Go and installer results.
- Source review: `git diff --check` passed for patch 0006 after the repair. The checked candidate is `/tmp/upstream-sync-native/source`.

Evidence: `/tmp/upstream-sync-native-build-check-final.log`, `/tmp/upstream-sync-native-build-ui-final.log`, and `research/implementation-handoff.md`.
