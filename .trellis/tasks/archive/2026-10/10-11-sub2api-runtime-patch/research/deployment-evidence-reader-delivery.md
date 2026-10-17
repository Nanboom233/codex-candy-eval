# Protected Deployment Evidence Reader Delivery

Date: 2026-10-12. Scope: isolated Sub2API patch 005.

## Context and Ownership

The implementer used native role context from the complete saved SubagentStart hook artifact.
The tool's visible hook output was truncated. The implementer read the saved artifact in bounded segments.
The frozen reader plan, current PRD, design, execution plan, and injected backend specifications governed this work.
The main confirmed explicit presence of all four `EvalProbeCase` tuple members.

Changed files:

- `dist/sub2api-build/workspace/backend/internal/service/eval_deployment_evidence.go`.
- `dist/sub2api-build/workspace/backend/internal/service/eval_deployment_evidence_test.go`.
- `patches/sub2api/v0.2.15/005-codex-eval-deployment-evidence.patch`.
- This delivery report.

The pristine source, patches 001-004, core files, lockfiles, build helpers, routes, middleware, Wire, bootstrap, lifecycle, UI, and installer remain unchanged.
No `.env`, CPA deployment, live target, paid request, container operation, commit, or push was used.

## Delivered Interface

The service exports the frozen types:

- `EvalDeploymentCurrentFacts`.
- `EvalDeploymentInventory`.
- `EvalDeploymentIngress`.
- `EvalDeploymentProxyHop`.
- `EvalDeploymentRunnerLocation`.
- `EvalDeploymentCaseDestination`.
- `EvalDeploymentEvidence`.

The method remains:

```go
func (s *EvalKeyStore) LoadDeploymentEvidence(
    cfg *config.Config,
    current EvalDeploymentCurrentFacts,
    expectedArtifactID string,
    expectedSHA256 string,
) (EvalDeploymentEvidence, error)
```

All fields and JSON names match the frozen reader plan. The loader returns the decoded inventory and lowercase exact-byte SHA-256.
It requires both artifact expectations. It rejects unavailable or different current facts and simple-mode configuration.
The reader derives the security fingerprint from live config. It compares every other current binding independently.

Filesystem checks use anchored `openat`, `O_NOFOLLOW`, and descriptor `fstat` from `/` through the fixed inventory path.
Private directories require effective UID ownership and mode `0700`. The file requires owned regular type and `0400` or `0600`.
Special permission bits fail. Nonblocking file opens prevent FIFO reads from waiting for a writer.
The loader closes every descriptor, creates nothing, and changes no file bytes.

The reader bounds exact bytes at 1 MiB and derived required cases at 256.
It checks the case limit before collector cross-product allocation and verifies the final collector case count.
Strict schema validation rejects duplicate keys at every depth, null, missing members, unknown fields, and trailing values.
The typed decoder also rejects wrong scalar types. The digest hashes accepted bytes before decoding.
Creation uses Unix seconds with five minutes of future tolerance. Expiry must exceed creation and the current time.
No inventory lifetime maximum was added.

The collector owns required case derivation. Local `direct_allowed` cases have no external destinations.
Additional valid, unique runner descriptors are permitted. Every case destination must reference a declared valid runner.
All declared runner descriptors are validated, including those with no destination references.
External coverage requires exact case tuples, recognized live headers, known runners, ingress identity, and ordered hop IDs.
The final hop source must match its ingress source. Ingress origins match the public/TLS hostname and TLS port.
Origins require HTTPS and an empty or root path. Userinfo, query, fragment, unsupported ports, loopback, and local names fail.
Ambiguous numeric IPv4 shorthand and non-IPv6 bracket hosts also fail.
The loader makes no network request and returns freshly decoded, detached slices on every load.

## Producer and Consumer Continuity

The later privileged operator/installer produces `deployment-evidence/inventory.json` under `EvalKeyStore.Dir()`.
It must inspect every public listener, proxy hop, rewrite, peer address, and deployment identity independently.
Protection and audit references do not prove operator authority or deployment completeness.

Every JSON member is required. Every `case` includes `kind`, `ingress_id`, `header`, and `source`.
Inapplicable tuple fields must be explicit empty strings. Ordinary `json.Marshal(EvalProbeCase)` omits them and is insufficient.
The tests include this producer serialization regression. Patch 005 adds no producer or writer.
Producers may include valid unused runner descriptors. Consumers must not require every declared runner to be referenced.

A later trusted provider supplies every `EvalDeploymentCurrentFacts` field from independent current sources.
Pinned `BuildInfo` supplies only version and build type. Missing source, patch, deployment, image, binary, or ingress facts remain unavailable.
Neither this inventory nor an admin request supplies expected current facts. The live security fingerprint comes only from config.

The coordinator later consumes `Inventory`, `SHA256`, and the same current-facts contract.
It can reconstruct peers from ingress ID/source and derive required cases through `NewEvalTopologyCollector`.
It must reload inventory/current facts at Begin, completion, evidence-only retry, and restart readiness checks.
Runtime session evidence, one-use observations, listener binding, readiness, admin auth, and host execution remain separate work.
A readable inventory neither establishes topology proof nor activates evaluation readiness.

## Validation

Formatting: `gofmt -l` reports no changed files in maintained and final replay trees.

Local root checks passed:

```sh
go vet ./...
go test -race ./...
```

Maintained copied-core checks passed:

```sh
go vet ./internal/codexeval
go test -race ./internal/codexeval
```

Final maintained and replay reader/service checks passed:

```sh
go test -tags unit -race ./internal/service -run '^TestEvalDeployment' -count=1
go test -race ./internal/service -run '^TestEvalDeployment' -count=1
go vet ./internal/service
go vet -tags unit ./internal/service
```

Reader tests use the frozen `unit` tag. The default race invocation validates compilation and reports `[no tests to run]`.
The final replay also passed `go vet ./internal/codexeval` and `go test -race ./internal/codexeval`.
These results follow the unused-runner specification correction and its new pristine replay.
Root and maintained-core checks above were completed during the initial delivery. Their sources remain unchanged.

Eleven top-level reader tests cover valid protected fixtures, exact-byte hashes, detached output, all current bindings, and semantic failures.
They also cover recursive missing/null members, duplicates at every struct depth, escaped duplicate keys, wrong types, bounds, and timestamps.
Filesystem tests cover missing storage, root/ancestor/evidence/file symlinks, non-directory parents, bad modes, FIFO, and socket paths.
Owner/type checks use the production stat validator without privileged `chown` or device creation.
Repeated successful and failed reads verify descriptor counts remain stable.
The FIFO/socket rejection test has a bounded wait. Socket fixtures use a short canonical temporary root for macOS path limits.
Case tests include exact 256-case acceptance, overflow rejection, live-header drift, direct-only topology, and mismatched coverage.
A positive unused-runner regression also confirms invalid unused descriptors still fail validation.
Origin tests include HTTPS restrictions, explicit approved ports, wrong ingress mapping, loopback variants, and malformed names.

Linux `amd64` and `arm64` service unit-test binaries compiled from final source using `CGO_ENABLED=0`.
`file` identified static ELF x86-64 and ARM aarch64 executables. These are compile checks, not Linux test execution or final release binaries.

## Final Patch and Replay

Patch 005 SHA-256:

`3402198f3f75e73ebf46e748ef4a289001a13d0632692c3bb493a630055650cb`.

Final pristine replay:

`dist/sub2api-build/scratch-deployment-reader-specfix-1wk5jtxn`.

The replay started from a new archive of `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`.
An isolated Git repository prevents nested-repository new-file skips.
Patches 001-005 passed sequential `git apply --check --whitespace=error` and strict application.
Both new reader files were asserted present. The 28 core files were copied unchanged before checks.
Full byte comparisons passed for all 3,258 backend files, 931 frontend files, and 28 core files.
`backend/go.mod`, `backend/go.sum`, and `frontend/pnpm-lock.yaml` match maintained and pristine source.
The earlier `scratch-deployment-reader-csfn52pb` and `scratch-deployment-reader-final-vnhun2qq` are superseded intermediate replays.
The final replay follows the specification correction permitting valid unused runner descriptors.
The prior patch hash `991d323e24d381c75b1754a84d09b55a791393f30e039b20713d693fd90f4dcc` is superseded.
Replay identity details remain in `deployment-reader-replay-evidence.json` at the final replay root.

Frozen identities remain:

| Input | SHA-256 |
| --- | --- |
| 001 | `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b` |
| 002 | `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5` |
| 003 | `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4` |
| 004 | `850af726f80bb2f379e0cf9b3f395f28e9d8e89356106f0643c39ee5558f4088` |
| Core 28-file historical identity | `3eabe2af794582dc8e41dfd637e57fc55940b9e57bc13e9606c1d748f9a8a7ba` |

The historical core hash concatenates sorted relative UTF-8 paths, NUL, and each raw 32-byte file SHA-256.
The build helper instead uses hex-text file digests and yields `b8477eb71da52e832aefa690b7bdf4fa02c257a792f1669ff3cbb5126ec67af5`.
These are different hash algorithms over identical file bytes, not source drift.

## Limits and Next Review

The full upstream service race suite remains blocked by documented pristine fixture failures in `research/upstream-race-baseline.md`.
The user declined unrelated fixture repair. Focused tests do not replace that suite or establish a full-service pass.
No Docker lifecycle, authenticated browser, real candy request, billing, or live deployment acceptance occurred.
No trusted current-facts provider or operator inventory installation exists yet. Missing facts keep the reader unavailable.

The corrected delivery passed final specification review and distinct independent quality review.
See [specification review](deployment-evidence-reader-spec-review.md) and [quality review](deployment-evidence-reader-quality-review.md).
Quality independently replayed the unchanged final patch at `dist/sub2api-build/scratch-deployment-quality-uTdpnl`.
Main captured the reader contract in `.trellis/spec/backend/deployment-evidence.md`.
No coordinator implementation was started. Later product chunks retain their existing planning and recovery contracts.
