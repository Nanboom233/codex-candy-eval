# Research: Stage 4B Run Manager and Durable State

- Query: Define a minimal Stage 4B manager and file-state contract against the implemented Stage 4A `internal/eval` APIs.
- Scope: internal
- Date: 2026-10-12

## Findings

### Stable Stage 4A boundary

`internal/eval` already exposes the complete Stage 4B execution boundary:

- `Target` is the routing identity and contains positive numeric `GroupID` and
  `AccountID` ([eval.go:13-16](../../../../internal/eval/eval.go#L13)).
- `AccountMetadata` is display-only and currently contains only `Name`
  ([eval.go:18-21](../../../../internal/eval/eval.go#L18)).
- `EvaluationRequest` carries type, exact target, display metadata, run ID,
  requested model/effort, mode, sample concurrency, and `OnProgress`
  ([eval.go:45-57](../../../../internal/eval/eval.go#L45)).
- `RunOne(ctx, ModelExecutor, EvaluationRequest)` performs one evaluation and
  leaves cancellation, batching, persistence, and account exclusion to its
  caller ([eval.go:81-82](../../../../internal/eval/eval.go#L81)).
- `Progress` is a copied value. The core invokes callbacks serially outside its
  data locks. Its final callback is emitted after collection workers settle
  ([eval.go:33-43](../../../../internal/eval/eval.go#L33),
  [README.md:10-14](../../../../internal/eval/README.md#L10)).
- `ModelExecutor.Execute` takes only a caller context and a `ResponsesRequest`
  containing the exact target, plain model, effort, body, and protocol headers
  ([model.go:20-28](../../../../internal/eval/model.go#L20)). `HTTPExecutor` adds
  the numeric target headers and never accepts credentials from the manager
  ([http.go:94-110](../../../../internal/eval/http.go#L94)).
- Every typed result embeds `ResultInfo`, so the selected group/account,
  requested/effective model and effort, run ID, and terminal status survive in
  JSON ([eval.go:59-79](../../../../internal/eval/eval.go#L59)). Candy, fingerprint,
  ModelTrace, and pelican have exported result structs and detail summaries
  ([candy.go:17-28](../../../../internal/eval/candy.go#L17),
  [fingerprint.go:37-53](../../../../internal/eval/fingerprint.go#L37),
  [modeltrace.go:56-74](../../../../internal/eval/modeltrace.go#L56),
  [pelican.go:19-30](../../../../internal/eval/pelican.go#L19)).

No Stage 4A core field needs to be added for Stage 4B. Group names and owner
identity belong in a later host metadata DTO. The result already preserves the
numeric group ID. A manager-level `StorageError` reports persistence failure;
an `ArtifactStatus`, `GroupName`, or owner field in `ResultInfo` would create a
new JSON contract without solving a 4B execution need.

### Narrow manager API

Keep the manager in `internal/eval/runs.go` and keep the store in
`internal/eval/state.go`. The following exported types are sufficient for the
later administrator adapter while avoiding Sub2API handler or repository
types:

```go
type TargetMetadata struct {
    Target  Target          `json:"target"`
    Account AccountMetadata `json:"account"`
}

type TargetResolver interface {
    Resolve(context.Context, Target) (TargetMetadata, error)
}

type StartRequest struct {
    Type        EvaluationType `json:"type"`
    Targets     []Target       `json:"targets"`
    Model       string         `json:"model"`
    Effort      string         `json:"effort"`
    Mode        string         `json:"mode,omitempty"`
    Concurrency int            `json:"concurrency,omitempty"`
    Runs        int            `json:"runs,omitempty"` // candy only; zero means one
}

type RunSnapshot struct {
    ID          string         `json:"id"`
    Type        EvaluationType `json:"type"`
    Target      Target         `json:"target"`
    Account     AccountMetadata `json:"account"`
    Model       string         `json:"model"`
    Effort      string         `json:"effort"`
    Mode        string         `json:"mode,omitempty"`
    Runs        int            `json:"runs"`
    Progress    Progress       `json:"progress"`
    Status      string         `json:"status"` // queued, running, cancelling, completed, failed, cancelled
    Error       string         `json:"error,omitempty"`
    StartedAt   time.Time      `json:"started_at"`
    FinishedAt  time.Time      `json:"finished_at,omitempty"`
}

type StartSummary struct {
    Started     []RunSnapshot `json:"started"`
    Unavailable []RunSnapshot `json:"unavailable"` // one terminal typed failure per rejected target/batch
    Busy        []Target      `json:"busy"`
}

type StateStore interface {
    Load() (PersistedState, error)
    Save(PersistedState) error
}

type HistoryScope string

const HistoryAll HistoryScope = "all"

type ManagerState struct {
    Active       []RunSnapshot              `json:"active"`
    Results      map[int64][]CandyResult     `json:"results"`
    Fingerprints map[int64][]FingerprintResult `json:"fingerprints"`
    ModelTraces  map[int64][]TraceResult     `json:"modeltraces"`
    Pelicans     map[int64][]PelicanResult   `json:"pelicans"`
    StorageError string                      `json:"storage_error,omitempty"`
    Quiescing    bool                        `json:"quiescing"`
}

func NewRunManager(parent context.Context, executor ModelExecutor,
    resolver TargetResolver, store StateStore) (*RunManager, error)
func (m *RunManager) Start(context.Context, StartRequest) (StartSummary, error)
func (m *RunManager) State() ManagerState
func (m *RunManager) Record(EvaluationType, int64, string) (EvaluationResult, bool)
func (m *RunManager) ClearHistory(context.Context, HistoryScope) error
func (m *RunManager) Flush(context.Context) error
func (m *RunManager) Cancel(string) bool
func (m *RunManager) Shutdown(context.Context) error

func NewFileStateStore(path string) (*FileStateStore, error)
```

`StateStore` is a small injectable interface over `PersistedState` with
`Load` and `Save` methods. `FileStateStore` is the production implementation.
The manager owns one store instance and one manager instance per configured
file. It does not add a file-locking service or a second scheduler.

`TargetResolver.Resolve` receives the exact positive numeric target chosen by
the administrator. It must return the same target or an error. A mismatch is
an unavailable-target failure, never a reason to query another group or
account. Host code remains responsible for administrator authorization and
normal group/account admission. Gateway execution remains responsible for
subscriptions, balance, quotas, concurrency, model capability, and billing.

### Admission and batch behavior

`Start` first validates the request and positive IDs, checks manager quiescing
and storage health, then resolves each target with the HTTP request context.
Resolver calls happen outside manager locks. A resolver error becomes one
explicit terminal typed failure under the requested account ID. It does not
start a worker and does not select a replacement target.

Admission uses `busy[accountID]`, not `(groupID, accountID)`. A running candy,
fingerprint, ModelTrace, or pelican therefore blocks the same account selected
through another group. Duplicate account IDs in one request are admitted in
request order; later duplicates are returned in `Busy`. The account slot is
released only after the worker has settled and the terminal state transition
has been attempted.

The manager creates one logical run per accepted target. Candy runs are capped
at `CandyMaxRuns` (10) and execute serially in that worker. Zero means one,
matching the old runner's default; non-candy types accept only one. Each candy
iteration gets a unique `RunID` derived from the logical run ID and its
sequence number, so typed history records remain distinct. The outer
`RunSnapshot.Progress` aggregates candy `Done`, `Valid`, and `Errors`; for the
other types it mirrors the latest core progress counts.

Admission is a transaction. Under `saveMu` followed by `stateMu`, build the
candidate snapshot containing all active intents and unavailable failures,
then call `Store.Save`. Commit the in-memory reservations only after `Save`
succeeds. Add workers to the wait group before releasing `stateMu`, but launch
them only after the save and locks complete. A failed startup save rolls back
all reservations and starts no executor call. A request context cancelled
before this transaction aborts without a record. Once the durable transaction
commits, the worker uses the manager lifecycle context even if the HTTP
request is cancelled or the acknowledgement is lost.

### Durable schema and file safety

Use one versioned JSON envelope. Keep the existing typed map names and key
histories by numeric account ID. JSON encodes integer map keys as decimal
strings, while the Go API remains `map[int64]...`:

```go
type PersistedState struct {
    SchemaVersion int `json:"schema_version"`
    Results       map[int64][]CandyResult       `json:"results"`
    Fingerprints  map[int64][]FingerprintResult `json:"fingerprints"`
    ModelTraces   map[int64][]TraceResult       `json:"modeltraces"`
    Pelicans      map[int64][]PelicanResult     `json:"pelicans"`
    Active        map[string]ActiveIntent       `json:"active,omitempty"`
}

type ActiveIntent struct {
    ID string `json:"id"`
    Target
    AccountMetadata
    Type        EvaluationType `json:"type"`
    Model       string         `json:"model"`
    Effort      string         `json:"effort"`
    Mode        string         `json:"mode,omitempty"`
    Concurrency int            `json:"concurrency,omitempty"`
    Runs        int            `json:"runs"`
    Completed   int            `json:"completed"`
    StartedAt   time.Time      `json:"started_at"`
}
```

`schema_version` must be exactly `1`. Decode through a pointer so a top-level
JSON `null` is rejected. Reject unknown schema versions, malformed JSON,
negative or zero IDs, invalid active counters, and invalid evaluation types.
Absent optional maps normalize to empty maps. A decode or validation failure
sets `StorageError` and blocks paid starts; it never replaces the original
file with an empty snapshot.

`FileStateStore` must reject a symlink at the state path, at the temporary
path, or in any existing parent component. The destination and temporary
file are in the same directory. Write a fresh temporary file with mode
`0600`, flush and close it, then rename it over the destination. Remove only
the temporary file after a failed write. A failed save leaves the previous
destination untouched and exposes the error through `ManagerState`; no new
paid start is admitted while the error is present. There is intentionally no
cross-process lock. The one-manager-per-file rule owns serialization.

Retention trims every typed map to its latest 20 records per account after
load and after appending. The account key is `AccountID`; each record's
`ResultInfo.Target.GroupID` preserves the selected group, including when an
account has history from multiple groups. Existing CPA state files are never
opened or migrated.

### Recovery and terminal persistence

Manager construction loads the envelope before accepting starts. For each
`ActiveIntent`, create one synthetic typed result with the exact target,
account name, requested model/effort, run ID suffix `:interrupted`, status
`interrupted`, and a safe error such as `run interrupted during restart`.
For a candy batch this is one task-level interruption record; already persisted
completed iterations remain intact. Delete active intents in a candidate
snapshot and save that candidate before exposing the manager as recovered. Do
not call the resolver or executor during recovery.

If the recovery save fails, retain the original loaded snapshot in memory and
on disk, set `StorageError`, and block all new starts. This prevents duplicate
interruption records on a later process restart and makes the failure visible.
The manager may expose the original active intents as a recovery-blocked state;
it must not claim they were durably interrupted.

Workers derive `runCtx` from the manager lifecycle context, never from
`context.Background()` or the HTTP request. `Cancel` marks the snapshot
`cancelling` and invokes its explicit cancel function without network work.
The callback passed to `RunOne` only copies `Progress` into in-memory state
under `stateMu`; it does not resolve metadata, call the executor, or save a
file. A cancellation before output still produces a typed `cancelled` result
with observed usage when the core provides it. Cancellation covers all four
evaluation types.

Each terminal result is appended to its typed account map, active intent
progress is advanced, and a candidate snapshot is saved while `saveMu` and
`stateMu` serialize transitions. Incorrect candy answers remain completed
records because `RunOne` returns them without an error. A transport or
evaluation error retains the typed result and error. If the terminal save
fails, keep the completed result and terminal in-memory snapshot, set
`StorageError`, release the account slot, and block later paid starts. Never
silently discard the result or report persistence success.

`Flush(ctx)` is the only retry path. It is explicit and synchronous: under the
same lock order, save the complete current in-memory snapshot once. A
successful write clears a transient write `StorageError` and reopens starts.
It must refuse to overwrite a file after a load/decode error or a
recovery-save failure because those states deliberately preserve the original
bytes. There is no background retry and no automatic paid retry. If `Flush`
fails, retain the error and the start block.

`ClearHistory` takes `saveMu` then `stateMu`, builds a candidate with only the
selected typed maps emptied, and saves it. Active intents, progress, cancel
functions, and busy-account reservations are copied unchanged. On save
failure restore the previous history maps and return the error. A successful
clear does not cancel a worker; a worker that completes afterward appends and
saves its result normally.

`Shutdown(ctx)` sets `quiescing`, cancels every run, releases `stateMu`, and
waits for the worker group with the caller's deadline. It must not hold
`stateMu` while waiting. A deadline returns `context.DeadlineExceeded` while
workers remain tracked and the manager remains quiescing; the local gateway
must stay alive until a later shutdown call observes worker completion. A
successful shutdown has no active workers and has attempted terminal saves. If
any terminal save failed, `Shutdown` returns that storage error after the
workers drain rather than claiming persistence success. The later host wiring
closes the local HTTP executor/gateway only after this method succeeds.

### Exact file ownership for one 4B implementer

- `internal/eval/state.go`: `PersistedState`, `ActiveIntent`, `StateStore`,
  `FileStateStore`, schema validation, safe load/save, retention, and history
  clear candidate logic.
- `internal/eval/runs.go`: `TargetResolver`, start request/response types,
  `RunManager`, admission, batch workers, progress snapshots, cancellation,
  recovery orchestration, and shutdown.
- `internal/eval/state_test.go`: JSON schema, null/unknown/corrupt input,
  symlink and mode checks, same-directory atomic replacement, retention, and
  clear rollback tests.
- `internal/eval/runs_test.go`: admission, metadata mismatch, startup
  acknowledgement, all four workers, progress, cancellation, recovery,
  final-save failure, and shutdown race tests.

Do not edit Stage 4A algorithm or transport files for this step. Do not add
administrator routes, Wire bindings, readiness/bootstrap, gateway shutdown,
or a second persistence system in the 4B-owned files.

### Required focused assertions

1. **Startup transaction and acknowledgement:** block a spy store's first
   `Save`; assert no executor call until the save succeeds and the worker is
   launched. Make that save fail; assert the request errors, no worker call
   occurs, no `busy` reservation remains, and the destination bytes are
   unchanged. Cancel the HTTP request after a successful commit; assert the
   worker still completes under the lifecycle context.
2. **Exact target and admission:** a resolver returning another group or
   account creates a failed typed record carrying the requested IDs and zero
   executor calls. Start account `11` in group `7`, then try account `11` in
   group `9` through another evaluation type; assert `Busy` and no second
   request. After terminal completion, assert the slot is reusable.
3. **Candy cap and retention:** request 20 candy runs and assert exactly ten
   serial `RunOne` calls. Produce 21 completed records for one account and
   assert only the latest 20 remain while each record retains its numeric group
   ID and requested/effective identity.
4. **Unavailable and recovery:** an unavailable selected target appears as a
   failed typed record. Seed an active intent, construct a new manager, assert
   one interrupted record and zero executor/resolver calls. A recovery save
   fault must leave the original file bytes and active intent untouched and
   reject a new paid start.
5. **Persistence faults and clear rollback:** fail a terminal save; assert the
   completed typed result remains in `State`, `StorageError` is non-empty, the
   busy account is released, and a new start is blocked. Restore the fake store
   and call `Flush`; assert the full in-memory result becomes durable and new
   starts can proceed. Fail `ClearHistory` while a blocking worker runs; assert
   histories and active/busy state remain unchanged. On success, only history
   is cleared and the worker continues.
6. **Progress and races:** feed callback updates while `State` and `Cancel`
   run concurrently. Assert callbacks perform no store writes or resolver
   calls and no data race occurs. Verify final progress counts match the
   corresponding `TraceResult` and `FingerprintResult` `Done/Total/Valid/Errors`
   fields; Stage 4A tests establish those count contracts
   ([eval_test.go:487-520](../../../../internal/eval/eval_test.go#L487)).
7. **Cancellation and shutdown:** cancel candy, pelican, fingerprint, and
   ModelTrace during executor calls, assert `cancelled` terminal records and
   observed usage, then assert account slots release exactly once. Use an
   executor that ignores cancellation to force a shutdown deadline; assert the
   method returns the deadline, workers remain tracked, and a later call after
   release drains them before gateway close.
8. **File safety:** assert mode `0600`, same-directory temporary writes,
   temporary cleanup, symlink rejection, top-level `null` rejection, unknown
   schema rejection, and preservation of corrupt/original bytes. Run all
   manager tests with `go test -race`.

## Caveats / Not Found

- No source-backed requirement calls for a group name, owner ID, or artifact
  status in Stage 4A result JSON. If the Stage 5 browser DTO needs a group
  label, add it to the host response assembled from the resolver, not to
  `ResultInfo` or the persisted typed records.
- The exact Sub2API account/group repository DTOs and authorization calls are
  intentionally not imported here. Stage 5 must implement `TargetResolver`
  against the pinned host and keep normal authorization, quota, subscription,
  balance, and model restrictions in their existing host/gateway paths.
- The manager does not prove local gateway readiness or topology. Bootstrap
  and shutdown wiring remain later owners. No live container, billing account,
  or paid request was available for this research.
- JSON rename provides atomic replacement at the file boundary. It is not an
  `fsync` or directory-durability guarantee, matching the repository
  persistence guide
  ([database-guidelines.md:33-53](../../../../.trellis/spec/backend/database-guidelines.md#L33)).
