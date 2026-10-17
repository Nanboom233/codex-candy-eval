# Research: Builder Timeout Analysis

- Query: Measure the cold and warm native Go graph/build stages for the user's source hint and establish what the 1,800-second deadline covers.
- Scope: mixed
- Date: 2026-10-14

## Findings

### Files found

- `install.py` — 1,800-second default build deadline, live-process timeout, builder dispatch and CLI timeout setting.
- `adapter/typecheck.go` — `go list -v -deps -export -json ./cmd/server`, JSON stdout capture and live stderr tee.
- `adapter/build.py` — stage order, readonly modules, type graph and final backend build flags.
- `adapter/Dockerfile` — one BuildKit `RUN` encloses the entire Python build pipeline.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md` — user's post-frontend timeout report and required stage visibility.
- `.trellis/spec/backend/universal-linux-delivery.md` — streaming, original diagnostics, deadline and Go-stage requirements.
- `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/builder-go-timeout.GseZIC/workroot/` — retained exact-source test inputs, isolated Go caches, logs and Linux/arm64 binary; 2.7 GB, retained for review.

### Exact-source Go measurements

The `018425` short hint resolved via the public GitHub commit API to `01842576b3d5c1e2545b5a99ad555d36e28ed306`. The source archive came from the corresponding codeload URL; `backend/go.mod` declares `go 1.27.2`. The source/hash is a reproduction input only, not a release pin. No Git command was run.

The host has Go 1.27.2 (`darwin/arm64`). A fresh isolated `GOCACHE` and `GOMODCACHE` under system `TMPDIR` were used. The command was:

```sh
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-mod=readonly \
go list -v -deps -export -json ./cmd/server
```

Cold graph/export completed successfully in 152 seconds. It emitted 141 module-download notices, 1,122 stderr lines total, and 4,828,616 JSON stdout bytes. The isolated module cache reached 502 MB and Go build cache 1.0 GB. The identical warm graph completed in under one second with no stderr output. This establishes that the Go type graph can include substantial cold module acquisition and export compilation; it does not establish the user's proxy latency.

With that cache, a vanilla upstream backend build used `GOOS=linux`, `GOARCH=arm64`, `CGO_ENABLED=0`, `GOFLAGS=-mod=readonly`, `-buildvcs=false`, `-trimpath`, and `-tags=embed`. A one-file `backend/internal/web/dist/index.html` fixture satisfied the Go embed pattern; the frontend was not built. The compile succeeded in 50 seconds, emitted 981 package lines, and produced a 148,173,378-byte static AArch64 ELF. An identical warm build completed in under one second with no package output. The cache reached 2.0 GB after the compile.

The graph measurement used the final CGO/OS/architecture settings but did not use `-tags=embed` or `-trimpath`; the following server build emitted hundreds of package compilation lines. This demonstrates that graph exports do not fully warm the final build under these flags. The run did not isolate which flag caused each miss, nor include the adapter-injected code or actual frontend assets.

### Deadline scope and interpretation

- `install.py:102, 810-845, 883-885, 946-968`: `BUILD_TIMEOUT` defaults to 1,800 seconds and `live_process()` applies it to the complete Buildx subprocess or fallback `docker exec` build subprocess. Its heartbeat reports elapsed command time; silence is not a Go-stage completion signal.
- The Buildx command covers all Dockerfile stages and its single `RUN python3 .../build.py`, including APK setup, frontend dependency/build work, type graph, and server compilation. The fallback timer begins at `docker exec`; its shell runs APK setup before the same complete `build.py` pipeline.
- Tool-image acquisition and Buildx bootstrap use separate command timeouts; source resolution uses `SOURCE_TIMEOUT=240`. The build deadline is not one timer per internal phase.
- Current `install.py:1263-1269` accepts a positive-seconds `--build-timeout` with the existing 1,800-second default. Keep the limit finite and configurable; do not remove the deadline or increase it solely to hide a slow stage.

### Related code patterns and specs

- `adapter/typecheck.go:99-110`: type graph preserves JSON stdout, enables `go list -v`, tees stderr live, and retains captured stderr on exit failure.
- `adapter/build.py:71-91`: frontend, native typecheck and backend build have flushed phase labels. The final Go build uses `-v`; both Go stages keep module mode readonly.
- `prd.md:196-201` and `universal-linux-delivery.md:103-108`: require internal Go package activity, separate JSON stdout/live stderr, and complete failures without duplicate output.

### External references

- GitHub REST commit lookup and codeload archive for the exact reproduction source hash above. Go module downloads used the public Go module proxy (`https://proxy.golang.org,direct`).
- Local Go command documentation: `go env` reported toolchain 1.27.2; `go doc` / `go help list` were used for local command behavior.

## Caveats / Not Found

- The measured Go stages total about 202 seconds locally without frontend build or adapter adaptation; they do not explain the reported 1,800-second timeout. No target Docker/BuildKit run, Linux emulation, CPU/memory profile, raw user builder log, or proxy timing trace was available. Do not assign the remaining time to network, CPU, OOM, or any single stage.
- For repeatability, the graph command explicitly set `CGO_ENABLED=0` to match the final backend build. The production typegraph code currently inherits the builder container's CGO default and only sets readonly `GOFLAGS`; that default was not measured. The compile-count comparison is evidence of cache divergence under the omitted `-tags=embed` / `-trimpath` flags in this controlled run, not a complete production cache-key attribution.
- The cold run did download many dependencies from the public proxy, but none timed out. It cannot distinguish whether the user's slow interval was sequential downloads, compilation, resource contention, or Docker/architecture behavior.
- The source archive omitted `.git`; the focused commands did not call production `adapter/build.py` or mutate the repository. The embed fixture and all generated caches/binary remain only in the standard-temp workroot listed above. The temporary research symlink used to preserve the running process path was moved to macOS Trash by the local filesystem hook; the workroot itself remains intact. No live production target or paid request was used.
