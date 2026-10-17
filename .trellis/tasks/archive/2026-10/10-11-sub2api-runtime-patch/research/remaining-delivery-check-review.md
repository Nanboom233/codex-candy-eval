# Remaining Delivery Check Review

Date: 2026-10-12
Scope: final Sub2API runtime, UI, installer, build, release, and test-isolation delivery.

## Verdict

The local delivery is substantially verified. Focused race suites, affected vet,
integration, frontend, artifact, payload, Docker lifecycle, and browser checks
passed. The required full upstream unit-race gate remains failed. No approved
live target exists, so live usage charges, deployed behavior, and complete
external ingress topology remain unverified. The Trellis task must remain
in progress.

One scoped production validation patch and its separate test-overlay
regression were fixed during review. Private split-DNS/VPN destinations remain
allowed when the named runner can reach them, consistent with the accepted
non-loopback origin contract. This does not replace the operator's duty to
attest to the complete ingress inventory and runner source.

## Findings (Fixed)

- File: `production/005-codex-eval-deployment-evidence.patch` and
  `test-overlay/005-codex-eval-deployment-evidence.patch`.
  Issue: the inventory reader rejected mapped loopback but checked
  `IsUnspecified` before unmapping IPv4-mapped IPv6 literals. That allowed
  `[::ffff:0.0.0.0]` through inventory validation, contrary to the frozen
  unspecified-literal rule.
  Fix: normalize the parsed address with `Unmap` before loopback and unspecified
  checks while retaining zone rejection. Added the mapped-unspecified origin
  regression to the separate test overlay.

## Findings (Not Fixed)

- Full unit-race gate: the final `go test -tags=unit -race ./...` failed with
  23 race warnings and 112 top-level failures, with zero compiler errors. The
  failures were confined to native `internal/handler/admin` and
  `internal/service` packages. No maintained overlay test name intersected the
  failed-test set. The prior source classification recorded 27 warnings and
  105 failures before this validator fix, while the historical run reported
  18 warnings and 58 failures. These counts are not interchangeable. No
  pristine baseline was repeated and no unrelated native fixtures were
  repaired. Keep the gate failed in reports. Evidence:
  `dist/sub2api-build/mapped-unspecified-full-unit-race.log` and the prior
  `research/full-race-source-classification.md`.
- Live acceptance: no approved Sub2API target was available. Real candy billing
  attribution and charges, production concurrency, deployed browser behavior,
  full independently inspected ingress inventory, and deployed payload
  inspection remain open. Local fixtures do not establish these facts.
- Windows host execution: packaged PowerShell `inspect` passed through local
  PowerShell, but native Windows execution was unavailable.

## Verification

- Repository checks: `scripts/test-sub2api.sh check` passed, including 68
  helper tests, five installer unit tests, shell syntax, root `go vet ./...`,
  root `go test -race ./...`, and the UI harness.
- Pinned replay and production/test separation: fresh replay passed for
  Sub2API v0.2.15 at `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`. Production
  and test-overlay workspaces had identical hashes for all 2,630 production
  paths. Evidence: `dist/sub2api-build/test-mapped-unspecified-20261012-evidence.json`.
  Test overlay identity is
  `6ab04a1d2c595b7721dd049b55c7f83cca66f0076864db913e9615a9dc17ab43`;
  augmented-tree identity is
  `006f964a847a347a0923f2e316fe253d762dfe11dcd41e0c20bd12f1331cb083`.
- Focused race: service, server, operator probe, DTO, and evaluation-core
  packages passed, including the mapped-unspecified regression. Evidence:
  `dist/sub2api-build/mapped-unspecified-focused-race.log`.
- Gateway/core race: handler, server middleware, evaluation route, and core
  packages passed. Evidence:
  `dist/sub2api-build/test-composed-final/gateway-core-race.log`.
- Affected vet: runtime package vet and gateway/core `unit` vet passed. Full
  backend `go vet ./...` also passed in the prior final replay. Evidence:
  `dist/sub2api-build/mapped-unspecified-affected-vet.log` and the earlier
  `dist/sub2api-build/test-review-final/vet.log`.
- Backend integration: final `go test -tags=integration ./...` passed.
- Frontend: pinned lint and typecheck passed; 362 files and 2,795 tests passed;
  production build passed.
- Wire generation: output was byte-identical, SHA-256
  `f5322d8332e75271401557194e003921a458b24b53aad7cd4318c64b1f6d81f5`.
- Artifacts: static stripped amd64 and arm64 server and operator binaries built
  after the fix. Payload allowlist and binary checksums passed. The fresh
  payload is `dist/sub2api-build/review-release-20261012-mapped-unspecified/`;
  it contains no test or development-verification inputs.
- Docker lifecycle: the prior full review passed 12 checks for each of arm64
  and amd64, including install, idempotence, conflict preservation, rollback,
  uninstall, recreation/reapply, and candidate-failure restoration. This
  validator-only fix did not change installer code; lifecycle tests were not
  repeated with the rebuilt server binary.
- Browser: authenticated local fixture passed four evaluation tabs, themes,
  history, API header authentication, Markdown/math, pelican sandbox, and CSP
  checks. Its records were synthetic.
- Packaged entry points: POSIX status and PowerShell inspect passed locally.
- Address policy: the accepted contract requires a non-loopback destination
  and separately rejects loopback, unspecified, link-local, zone-qualified,
  and local-hostname destinations. It does not require public CIDR
  routability. The named approved runner may use private DNS/VPN reachability;
  origin binding, TLS hostname verification, redirects-disabled probing, and
  independently complete ingress inventory remain required.

## Requirement Matrix

| Requirement | Review result |
| --- | --- |
| R1–R2 / A1 | Four evaluations and preserved controls passed in the assembled local gateway fixture. No live inference was performed. |
| R3 / A2–A3 | Local native accounting and shared concurrency fixture passed. Real charges and production concurrency remain unverified. |
| R4–R6 / A4–A7 | Exact Key, group, and account routing; normal admission; no fallback; and cancellation fixtures passed. |
| R7 / A8 | Admin authorization, private-token exclusion, lookup redaction, and replay checks passed locally. |
| R8 / A9 | Retention, recovery, cancellation, history, and restart checks passed locally. |
| R9 / A10 | Both Docker image identities and original deployment configuration remained unchanged during lifecycle checks. |
| R10 / A11–A12 | Idempotence, restore, rollback, uninstall, and recreation/reapply checks passed on both architectures. |
| R11 / A13 | Pinned amd64/arm64 builds and compatibility rejection checks passed locally. |
| R12 / A14 | Sub2API installation, operations, and release guidance replaced CPA-specific public paths. |
| R13 / A15 | Local browser checks passed for layout, theme changes, details, history, and sandbox behavior. Deployed browser behavior remains open. |
| R14 / A16–A17 | Local listener, spoof, proof, conflict, and scoped revocation checks passed. Real ingress completeness remains an operator/live obligation. |
| R15 / A18 | Separate production and overlay inputs, path/hash parity, release allowlist, and artifact inspection passed. |

## Key Boundary Matrix

| Contract | Review result |
| --- | --- |
| K1 | Raw token stays in protected server files; local checks found no browser, log, URL, process-argument, or artifact exposure. |
| K2 | Exact owned-row redaction and external lookup replay rejection passed local regression checks. |
| K3 | Normal Key permissions, subscriptions, quota, balance, pricing, and concurrency remain in the gateway path. Local fixtures passed. |
| K4 | Local topology and spoof fixtures fail closed. Named-runner private DNS/VPN reachability is allowed by the precise non-loopback contract; actual ingress completeness remains unverified. |
| K5 | Planned rollback/uninstall revoked only the owned Key and preserved unrelated Keys in both Docker runs. |
| K6 | Accepted original-runtime lookup visibility/use residual remains documented; automatic revocation after unplanned recreation is not claimed. |
| K7 | Owner/status/binding conflicts were detected without silent repair or reenabling; missing-row recovery remains explicit operator repair. |

## Artifact Identities

- Production patch set: `c437e1eff6792bf87dcf9c1baeeceeb4c1ec8c96d281b199885413218af32169`.
- Production core/resource tree: `71cebfd65c9167eacaff687d677e51a032ed12503c876185a0d6236ff2f7bd58`.
- Production staged tree: `4f6385b0a243e8cdea82a1db32ac24c020e006b350662293da0510e159088ba0`.
- Development overlay patch set: `6ab04a1d2c595b7721dd049b55c7f83cca66f0076864db913e9615a9dc17ab43`.
- Development augmented tree: `006f964a847a347a0923f2e316fe253d762dfe11dcd41e0c20bd12f1331cb083`.
- amd64 server: `a068211e3279a84694482da680fe209c6934fd10817f17b54a2b290c9f49c0fe`.
- amd64 operator probe: `58eaf2ebfe3f3b18ea87f490060202702baf2bfe72de7f3d78eeb7338ac5fde1`.
- arm64 server: `377cc310972be3492181896a0dc5f1b458f5c7fb76ea83e972afc099abe95929`.
- arm64 operator probe: `3733f8b2720d2e573e6a6f274df70e6445b0ec5f9ab9d798bcfe76ba1a92b88f`.

All evidence paths are relative to `dist/sub2api-build/` unless otherwise
stated. The saved failed race gate is a required blocker, not an exception to
completion. No deployment, paid request, CPA `.env` access, commit, or archive
was performed.
