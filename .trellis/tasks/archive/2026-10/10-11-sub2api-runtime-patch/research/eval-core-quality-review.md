# Stage 4A Evaluation Core Quality Review

- Date: 2026-10-12
- Scope: `internal/eval` transport, SSE parsing, four evaluation operations, progress callbacks, and staged Sub2API copy.
- Verdict: **QUALITY PASS**

## Findings (fixed)

### Empty query delimiter accepted by the HTTP executor

`NewHTTPExecutor` rejected non-empty `RawQuery` values but did not reject
`url.URL.ForceQuery`. An origin such as `http://127.0.0.1?` therefore passed
validation and produced an endpoint with a trailing query delimiter.

The constructor now rejects `ForceQuery`, and the origin-validation test covers
the empty-query case. The executor remains bound to a literal loopback origin.

## Findings (not fixed)

- `.trellis/spec/backend/evaluation-core.md` still says the independent review is
  pending. The main session owns spec files, so it should update that status after
  recording this verdict.
- The full pinned upstream service race suite remains blocked by known failures
  on the pristine source. No unrelated fixture changes or race suppression were
  made.
- Docker lifecycle and live gateway checks remain unavailable because the Docker
  daemon is unreachable and no approved target is configured. These are later
  migration gates, not Stage 4A local quality failures.

No other Stage 4A code defects were found.

## Verification

- Formatting: `gofmt -l internal/eval` clean.
- Lint/type checks: no separate lint command is configured; root and copied-package `go vet` passed.
- Root tests: `go test -race -count=1 ./...` passed. Evaluation core: 10.303s. Existing plugin: 16.985s. Wall time: 17.853s.
- Root vet: `go vet ./...` passed.
- Copied package: `go vet ./internal/codexeval` passed in the pinned backend module.
- Copied package tests: `go test -race -count=1 ./internal/codexeval` passed in 10.225s.
- Staging: `python3 scripts/sub2api_build.py stage-eval-core compatibility/sub2api-v0.2.15.json . dist/sub2api-build/workspace` passed.
- Copy identity: all 24 files under `internal/eval` match `backend/internal/codexeval` byte-for-byte.
- Algorithm parity: fingerprint and ModelTrace algorithm files match the originals apart from the package declaration.
- Asset parity: Codex context, fingerprint probes/baselines, ModelTrace bank, license, README, and golden fixture match their originals byte-for-byte.
- Patch integrity: frozen patch SHA-256 values remain 001 `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`, 002 `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`, and 003 `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`.
- `git diff --check` passed.

The Stage 4A local quality gate is complete. Deployment, live billing, browser
integration, and the pristine upstream full-suite race gate remain incomplete.
