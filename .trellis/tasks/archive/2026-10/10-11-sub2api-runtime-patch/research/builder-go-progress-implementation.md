# Builder Go progress implementation

Status: implementation PAUSED by explicit user architecture replacement; candidate not promoted, original review not complete. See builder-go-progress-paused-checkpoint.json. No deployment, paid request, commit, push, publication or artifact cleanup.

## Changes

`adapter/build.py` prints flushed phase labels before native adapter compilation, frontend dependency installation, frontend build, native Go type checks and final backend compilation. Adapter/server `go build -v` exposes package activity. It preserves existing GOFLAGS and appends -mod=readonly before adaptation/typechecking and final build; no host-version branches, native contract relaxation or full-history/source refetch.

`adapter/typecheck.go` adds nativeTypeGraph(root): `go list -v -deps -export -json ./cmd/server` captures JSON stdout for its existing decoder, tees stderr live to os.Stderr plus an unbounded raw byte buffer, and preserves full stderr in the wrapped exec.ExitError. On command failure typecheck emits partial JSON stdout exactly once, then existing fail prints the wrapped exit cause. Previously streamed stderr is not reprinted. Errors.As retains the original exit cause. Import/export parsing and all method/field/receiver contracts remain unchanged. Native graph invocation also appends readonly even when called independently.

This revision does not yet alter the 1800-second build deadline, target CGO/GOARCH, trimpath/embed flags or cache behavior. The earlier log makes native typecheck the next likely phase after frontend, but no remote process trace proves the exact waiting command. Progress labels identify subsequent stages without equating heartbeat with compilation completion.

Changed owners: adapter/build.py, adapter/typecheck.go, maintained adapter/typecheck_progress_test.go and tests/installer/test_builder_progress.py, scripts/test-sub2api.sh routing, docs/sub2api-operations.md and installer functional hash. Native adapter semantics/core/runtime unchanged. Final package identity pending below; v0.4.0 remains unpublished.

## Evidence

Workroot `/var/folders/md/v0dh5x2n7t17w1j31yp4kkhh0000gn/T/codexeval-go-progress-8iek7wdx` was created with system tempfile.mkdtemp and recorded in builder-go-progress-workroot.json. No custom workroot or /host scratch. Local root Go cache uses current go env GOCACHE `/Volumes/User Data/caches/go/build`.

`PYTHONDONTWRITEBYTECODE=1 python3 -B tests/installer/test_builder_progress.py`: final barrier run 3/3 exit0 (`builder-progress-barrier.log`). Compiles actual current Go adapter from main.go/typecheck.go, then uses executable fake Go child with delayed .15-second stderr and a release-file barrier. Parent must observe complete original stderr while child is alive before releasing success JSON or exit17 with partial JSON. Exact rawstderr count1, complete partial stdout and exitcause line asserted. JSON-success branch reaches native module validation, demonstrating stderr did not enter the decoder. GO_CALL_LOG verifies list-v flags and retained GOFLAGS plus readonly. A real tiny standard-library Go module also executes go-list-v with export/json; stdout JSON decodes, its package compile appears in stderr and go.mod is unchanged. Full native host build is not claimed.

Actual build control flow test stubs adaptation commands only, invokes current builder.build, records flushed labels before each command and type/backend .1-second timed boundaries. Checks adapter/backend build-v, frontend/type/backend ordering and readonly GOFLAGS across commands. Mock commands do not count as actual compilation. Go unit `TestNativeTypeGraphPreservesJSONAndRawExitStreams` executes shell fakego and directly asserts returned JSON/partial bytes, errors.As exec.ExitError code17 and full rawstderr buffer without truncation.

Initial independent timing fixture waited only1.2s for a freshly compiled native binary and observed no early stderr, failing both subtests (`builder-progress.log`). Same product root gate then passed; changed fixture to synchronization barrier with bounded5s startup wait rather than relying on a fixed process-exit delay. No product deadline widening. Initial failures remain evidence, not acceptance.

`PYTHONDONTWRITEBYTECODE=1 scripts/test-sub2api.sh check`: exit0 (`root-check.log`), including 72 build-helper tests, installer regressions, new3 builder-progress tests, go vet/race and UI. This first root run predates final barrier/timing-only test refinements; final gate/package tests will be appended. gofmt applied to changed Go files; git diff --check exit0. Existing earlier full upstream race/frontend/integration failures remain historical open gates.

## Boundaries and cleanup inventory

Unknown user's exact waiting command/network/cache state, iStoreOS/native amd64 loaded lifecycle, native billing/concurrency, production/browser and published bootstrap remain unverified. No whole upstream native rebuild is necessary for this progress-only change; local tiny Go and actual adapter command boundaries provide focused evidence.

Retain workroot logs and later package staging/extraction until checker; temporary compiled adapter/fakego/tinyGo test fixtures clean their own tempfile roots. Current Go build/test cache from required checks and any new __pycache__ directories are generated outputs for later exact cleanup. No Docker fixture/build cache was created by this progress correction. Historical canonical archive remains preserved until verified replacement. No cleanup or global prune.

## User architecture replacement — paused checkpoint

The user explicitly rejected recompiling the complete host and requested a different method. All current edits/evidence are preserved; no further build/test/repack/promotion runs. No own live process remained at pause, so no process was killed. Canonical archive remains previous519831/62/SHA773d9...3a95; unpromoted candidate520385/62/SHA53fa...f538 is retained in workroot. Latest current branch includes --build-timeout positive seconds/default1800 and complete progress/readonly edits, but these are not represented as a solution for the30-minute stall or accepted final architecture. Go environment/cache/compilation flags were not changed beyond readonly/verbose diagnostics.

Final current source root check passed; archived builder5 and fresh isolated boundary14 passed. The separately named archive universal/transactions logs actually exercised current source because those maintained suites ignore INSTALLER_TEST_PATH; do not claim archived coverage for them. One shared extraction later lacked all regular files and caused initial boundary failures; its cause is unproven. The earlier suspicion of transaction tamper is withdrawn (that suite owns separate temporary bundles). Fresh independent extraction passed. Candidate has not been promoted. Precise identities, maintained-file hashes, limits and no-pending-process inventory are in paused checkpoint. New no-host-recompile architecture research is owned by root; no cleanup until review.
