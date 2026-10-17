# Deployment Evidence Reader Quality Review

Date: 2026-10-12. Scope: final quality review of patch 005.

## Context and Verdict

Context mode: native Trellis hook injection. The visible hook output was
truncated, so the complete saved hook artifact was read in bounded sections.
The frozen plan, delivery, specification PASS, PRD, design, execution plan,
and directory/logging guides were read before review.

Verdict: PASS. No product files required correction. No code findings remain
unfixed for this reader scope.

## Findings (fixed)

None.

## Findings (not fixed)

None in the assigned reader scope. The full upstream service race suite was
not rerun because the task records reproducible pristine fixture failures and
the user declined unrelated fixture repair. This review does not claim that
suite passed.

Linux test execution, release binaries, Docker lifecycle checks, live topology,
authenticated browser behavior, and real billing remain outside this local
reader gate. No approved live Sub2API target or trusted current-facts provider
is available. The reader remains unavailable until the later provider supplies
independent current facts.

## Review Evidence

The reader matches the frozen exported API and remains a read-only consumer.
It walks the absolute store path with anchored `openat` calls and
`O_NOFOLLOW`, checks opened descriptors with `fstat`, enforces effective-UID
ownership and exact private modes, rejects special files, and bounds reads.
The strict schema pass rejects duplicate, unknown, missing, null, and trailing
JSON data before typed decoding. Exact accepted bytes determine the digest.
Current deployment facts are compared independently, while the security
fingerprint comes from live configuration. Required external cases come from
the existing collector, and every declared runner is validated while valid
unused runner descriptors remain allowed. HTTPS origin and ingress/hop
bindings are checked without network requests. Every successful load returns
newly decoded slices.

The reader adds no writer, model or network call, readiness change, route,
ticket, provider, coordinator, or lifecycle behavior.

The final interface remains `EvalDeploymentCurrentFacts`,
`EvalDeploymentInventory`, `EvalDeploymentIngress`,
`EvalDeploymentProxyHop`, `EvalDeploymentRunnerLocation`,
`EvalDeploymentCaseDestination`, `EvalDeploymentEvidence`, and
`(*EvalKeyStore).LoadDeploymentEvidence(cfg, current, expectedArtifactID,
expectedSHA256)`.

The later operator/installer must produce the protected inventory and encode
all four case tuple fields explicitly, including empty inapplicable values.
The trusted current-facts provider must source version, build, deployment,
source, ordered patch, image, binary, and ingress facts independently. The
coordinator must reload the inventory and facts at its required lifecycle
checks. A valid inventory alone does not establish topology proof or readiness.

## Verification

- Lint: pass. `gofmt -l` returned no changed files in maintained and replay
  trees.
- TypeCheck: pass. Default and `unit` service vet passed in both trees; root
  `go vet ./...` and replay copied-core vet passed.
- Tests: pass. Focused `TestEvalDeployment` race tests passed in maintained and
  replay trees. Root `go test -race ./...` and replay copied-core race tests
  passed. Default service race commands compiled successfully and reported
  `[no tests to run]`; they are compile checks, not test execution.

A new isolated replay was created at
`dist/sub2api-build/scratch-deployment-quality-uTdpnl` from pristine commit
`f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. Patches 001 through 005 passed
sequential `git apply --check --whitespace=error` and application. Both new
reader files were asserted present and compared byte-for-byte with the
maintained tree. The full backend (3,258 files), frontend (931 files), and
copied evaluation core (28 files) matched. All three lockfiles matched
maintained and pristine source. The 28 core files were staged in the isolated
replay.

Final patch 005 SHA-256:
`3402198f3f75e73ebf46e748ef4a289001a13d0632692c3bb493a630055650cb`.

Frozen patch hashes remained unchanged:

- 001: `1274cf97605c7eb134c3825fdbe16d8786dba04ba5acbda93f75d8bfa9df922b`
- 002: `f03dafc2d21c7e49dcda567ac2e5ff51a835d5feb4e69b3210dc3892bf041fd5`
- 003: `6a7602528be3c68244856a3ee3b57467bf24a4ad829436afeeb021a8f643bed4`
- 004: `850af726f80bb2f379e0cf9b3f395f28e9d8e89356106f0643c39ee5558f4088`

The existing full upstream service race baseline remains a documented blocker;
it was not represented as a pass. No Linux test execution, release build,
container lifecycle, deployment, paid evaluation, or live acceptance was
performed.

## Handoff

The main session may run `trellis-update-spec` after this quality PASS, then
synchronize the task checkpoint and manifests before the next chunk. The next
product work remains the scoped coordinator chunk; this review does not
authorize combined runtime integration.
