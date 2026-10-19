# Confirmed account application

## Change boundary

Apply model control and priority after each account establishes a valid state through the existing confirmation rules.
Include stable states and process completed accounts before unrelated confirmations finish.
Keep runtime sending, confirmation counts, recovery waits, cancellation, and participation unchanged.
Patch `0005` owns completion collection, account application, unchanged writes, and removal-list deduplication.
Extend existing coordinator and account-helper tests for these changes.
Reuse UI, installer, and unchanged transport evidence.

## Evidence

Patch `0005` now collects each completed account before the remaining accounts finish.
Stable results use the same account application owner as confirmed transitions.
The owner saves model changes and verdict before optional priority changes.
It retains the existing confirmation counts, recovery wait, participation checks, cancellation, and result fields.
The service skips matching writes and merges actual removals without duplicate names.
Notifications describe the confirmed state without claiming three tests for stable results.

The existing coordinator tests cover ready-account application, stable-state correction, matching priorities, and adjustment-off behavior.
The existing service tests cover repeated removals, unique saved names, stable recovery, and unchanged account writes.

The following checks passed against `/tmp/sub2api-confirmed-account-control/backend`:

```sh
GOFLAGS=-mod=readonly go test -tags=unit ./internal/codexeval ./internal/service ./internal/handler/admin -run CodexInspection
GOFLAGS=-mod=readonly go vet ./internal/codexeval ./internal/service ./internal/handler/admin
```

`gofmt -l` returned no changed production or test Go files. `git diff --check` passed.
The full production patch series applied to the unchanged original source.

The existing build entry passed:

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check
```

It rebuilt the frontend, Linux/amd64 server, and embedded installer.
The build retained nonblocking upstream dependency and chunk-size warnings.
The final outputs are `dist/build/sub2api-codex-eval` and `dist/build/sub2api-codex-eval-installer`.
Reuse prior UI and installer checks because their inputs did not change in this correction.

## Handoff

The effective tree is `/tmp/sub2api-confirmed-account-control`.
The source patch and affected tests are ready for independent review.
The main session owns specification synchronization and final acceptance.
This correction used synthetic tests and no provider requests.
It did not install a server, publish a release, stage files, or create a commit.
