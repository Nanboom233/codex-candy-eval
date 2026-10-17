# Research: Builder Go Progress

- Query: Locate the post-frontend quiet phase and identify how Go progress can be shown without corrupting the typecheck JSON stream.
- Scope: internal
- Date: 2026-10-14

## Findings

- `adapter/build.py:11-12` runs child commands with inherited stdout/stderr and `check=True`, so the Python wrapper itself does not buffer their output.
- Build order in `adapter/build.py:69-85`: pnpm install/build; copy and validate embedded frontend; run the adapter `types` operation at line 78; then compile the Linux server with `go build` at lines 80-85. There is no phase label before the typecheck and the final `go build` lacks `-v`.
- The strongest candidate for a quiet wait after frontend completion is the typecheck at `adapter/build.py:78`. It invokes `adapter/typecheck.go:99-105`, which runs `go list -deps -export -json ./cmd/server` and calls `cmd.Output()`. JSON stdout must remain captured for the decoder at lines 107-120. Successful-command stderr is currently not connected to the parent process, so it is discarded.
- Local `go doc os/exec.Cmd.Output` confirms `Output()` returns stdout and populates `ExitError.Stderr` only when `Cmd.Stderr` is nil. Set `Cmd.Stderr` to a tee that writes live to `os.Stderr` and retains raw bytes; `Output()` can continue capturing stdout for JSON. Add `-v` to `go list` so package activity appears on stderr. On failure, preserve the returned partial stdout and captured stderr in the error/log path while avoiding a second print of stderr bytes already streamed live.
- Add a flushed phase message immediately before `engine types` (for example, `Check native Go type contracts.`) so a quiet toolchain still identifies the current stage. Add `-v` to the final `go build` at line 84 to expose package compilation progress. Keep the existing `GOFLAGS=-mod=readonly` at line 80; `-v` is an explicit Go command flag and does not relax module immutability. Use similarly short labels before adapter compilation and final server compilation if the installer should identify all Go stages.
- Existing task contract: PRD `prd.md:182-186` requires progress before command exit and complete original diagnostics; design `design.md:161` and universal delivery spec `universal-linux-delivery.md:103-106` require live chunk teeing, separate stdout/stderr, no duplicate failure replay, and measured download bytes.

## Caveats / Not Found

- This is a code-path diagnosis, not a confirmed runtime trace. The typecheck is the first Go stage after frontend completion, but no current log proves it had started during the reported wait.
- No full build, network operation, or production target was used for this research.
