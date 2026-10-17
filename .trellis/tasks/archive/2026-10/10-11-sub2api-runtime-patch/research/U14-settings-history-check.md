# U14 Settings and History Check

## Findings (fixed)

None. The source review found no confirmed code defect that needs a local fix.

## Findings (not fixed)

- The source of the user's earlier `context canceled` history record remains unknown. The new integration test proves that ending the start and state requests does not cancel a background Responses or Chat Completions evaluation. It also proves that an explicit stop cancels the supplier request and records one cancelled result. The test does not identify the cause of an earlier record from a retained or production run. Do not report that historical event as resolved without its run context.

## Review

- The four evaluation pages save existing preferences on `input`. Their `change` handlers still normalize controls and refresh the view.
- Each settings-tab entry reads native settings again. Dirty fields remain unchanged during a read.
- A save advances the read generation. A late GET cannot replace that save. A successful response clears only submitted fields that were not edited again. Blank secret fields retain the native save contract.
- The native settings service still reads and writes the existing settings repository. It keeps secret preservation and validates the priority values. The current service files match the reviewed native source.
- Manual starts use a background context after the start handler returns. Inspection starts retain the coordinator context. The native upstream transport owns each attempt until its response body closes. The added test uses the actual handlers, runtime, account-test sender, native HTTP upstream, and a delayed local SSE supplier.
- The integration test covers both Responses and Chat Completions conversion. It preserves successful history and observed token counts. Explicit stop remains a cancelled result.

## Verification

- Lint: no U14-specific linter was run. `node --check` passed for all five changed JavaScript files.
- Type check: passed through the focused Go package test compilation. `go vet` passed for `./internal/service`.
- Tests: passed. The focused UI regression passed:

  ```sh
  SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs
  ```

  The selected Go settings and request-lifetime tests passed in `/tmp/sub2api-U11-check.CqGH5I/backend` with `/tmp/sub2api-U14-go-overlay.json`:

  ```sh
  go test -v -tags=unit -overlay /tmp/sub2api-U14-go-overlay.json ./internal/service -run '^TestCodex(InspectionSettings|EvalBackgroundNativeStreamLifetime)' -count=1 -timeout=60s
  go vet -tags=unit -overlay /tmp/sub2api-U14-go-overlay.json ./internal/service
  ```

- Formatting: `gofmt -d` reported no change for the new Go test. `git diff --check` passed for the reviewed UI files.
- Docker, browser, and delivery build: not run in this source review. The main session owns the retained-container and final artifact checks.
