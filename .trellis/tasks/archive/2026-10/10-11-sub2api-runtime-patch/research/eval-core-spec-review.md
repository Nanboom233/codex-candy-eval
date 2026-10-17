# Stage 4A Evaluation Core Review

- Date: 2026-10-12
- Scope: `internal/eval` transport, SSE parsing, model normalization, four evaluation operations, progress callbacks, and staged Sub2API copy.
- Verdict: Stage 4A specification PASS; independent quality review pending.

## Findings fixed

### Short-value truncation

`truncate` now trims leading and trailing whitespace for values below the limit as well as bounded values. The regression test covers both paths.

### Fingerprint effort validation

Fingerprint rejects an explicit `none` effort in the model suffix or effort field. General model normalization still accepts `none` by omitting the wire effort. The validation regression covers both forms.

### Model suffix slicing

During review, a suffix-extraction edit referenced an index outside its scope. The reviewer fixed this mechanical compile issue by deriving the index in `NormalizeModel`. It was introduced during the review correction, not present in the original implementation that had passed vet.

## Evidence

- The core has no CPA imports, host callbacks, or credential types.
- `RunOne` preserves target IDs, requested and effective model values, and caller-owned cancellation.
- HTTP execution owns the token, fixed routing headers, loopback origin, proxy-disabled transport, and redirect rejection.
- SSE tests cover chunk boundaries, terminal failures, usage retention, output-item reconstruction, bounds, and cancellation.
- Evaluation tests cover candy, fingerprint, ModelTrace, pelican, retry identity, sample limits, progress serialization, and cancellation settlement.
- Maintained and staged trees contain the same files byte-for-byte.
- Data, license, and ModelTrace golden inputs match their source copies.

## Verification

These checks provide specification-review evidence. They do not replace the independent quality review.

- `gofmt -w internal/eval`: pass
- `go vet ./...`: pass
- `go test -race -count=1 ./...`: pass
- `python3 scripts/sub2api_build.py stage-eval-core compatibility/sub2api-v0.2.15.json . dist/sub2api-build/workspace`: pass
- Pinned backend `go vet ./internal/codexeval`: pass
- Pinned backend `go test -race -count=1 ./internal/codexeval`: pass
- `diff -rq internal/eval dist/sub2api-build/workspace/backend/internal/codexeval`: pass
- `git diff --check`: pass
- Data, license, and golden byte comparisons: pass
- Frozen patch hashes: 001 `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`; 002 `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`; 003 `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`
- Latest fresh root and copied-package race test runs: 2026-10-12; exact wall-clock time was not captured.

## Not fixed or not applicable

- The full upstream race baseline remains blocked by known failures on the unmodified pinned tree. No unrelated fixtures were changed.
- Docker lifecycle, deployment, authenticated browser checks, live billing, and live model calls were not run because no approved target or available Docker daemon exists.
- Stage 4B durable run state, administrator wiring, topology bootstrap, UI integration, patched binaries, and installer work remain outside this review.
