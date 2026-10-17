# Research: Installation Simplification Options

- Query: How can the Sub2API patch provide one-command installation followed by first activation in the existing administrator web panel, without a user-created JWT file, hand-authored JSON manifest, or separately run probe command, while retaining K1-K7 and topology evidence requirements?
- Scope: mixed (internal source/spec evidence with product-level option analysis)
- Date: 2026-10-12

## Findings

### 1. The current command combines independent responsibilities

The installer currently validates an explicitly named container, reads an optional administrator JWT file, accepts an administrator origin, and requires an independently supplied ingress revision before `install` can proceed (`sub2api_install.py:49-70, 326-331, 436-443`). It also performs authenticated post-restart checks and mutates the custom menu (`sub2api_install.py:234-280, 377-380`). The protected transaction journal and installer facts are written inside the container (`sub2api_install.py:176-220, 358-374`); those are internal state and do not need to remain user-authored artifacts.

The binary transaction itself has a smaller trust boundary. Image/source/platform identity, pristine-binary comparison, immutable backup, staged atomic replacement, restart, and restoration are host/container operations (`sub2api_install.py:96-140, 346-393`; `research/patch-contracts.md:360-376`). They can run without administrator API credentials. A simplified install can therefore:

1. Discover exactly one running container matching the pinned Sub2API identity, or fail with an actionable list when zero or multiple candidates exist.
2. Verify architecture, image/source digests, startup path, writable data mount, UID and original binary checksum.
3. Stage, atomically replace, restart, and verify local process/health/hash state.
4. Generate protected installation facts internally and print/open the stable evaluation URL.

The command should defer menu merge, private-Key creation, ingress revision, and topology readiness to an authenticated browser setup flow. This removes `--admin-token-file`, `--origin`, and `--ingress-revision` from the default install journey without weakening R9-R11. A failed candidate must still restore the original binary, and image drift must still block old-backup reuse.

### 2. Existing browser authentication can supply activation authorization

Native Sub2API stores the current access token in `localStorage` and the API client attaches it as a same-origin `Authorization` header (`dist/sub2api-build/workspace/frontend/src/stores/auth.ts:104-125`; `dist/sub2api-build/workspace/frontend/src/api/client.ts:41-47`). The evaluation page contract already requires a static same-origin shell whose APIs use the existing administrator login; it must not receive a second secret or put JWTs in URLs (`research/patch-contracts.md:302-324`).

This supports a browser-first activation path:

- The host command installs only the runtime and emits `/codex-eval/ui` (or opens it after printing it).
- The administrator logs into the normal panel and opens that route, even before a menu entry exists.
- The page calls authenticated setup/status APIs using the parent panel's current session.
- The server creates/recovers its private evaluation token and internal JSON state, then performs a narrow, conflict-aware menu merge from the authenticated request. The browser never sees the raw token or stored lookup.

The current installer-side settings PUT and menu ownership checks are therefore implementation machinery that can move behind the authenticated setup endpoint (`sub2api_install.py:262-280`). The endpoint must preserve unrelated menu entries and report conflicts rather than repairing them silently.

### 3. Topology evidence cannot be inferred from Docker or the browser

K4 requires source-IP and proxy validation before activation. The topology collector can validate only the declared cases; it cannot discover an omitted CDN, WAF, load balancer, or reverse proxy. The research contract explicitly makes complete ingress inventory an operator/deployment acceptance obligation (`research/host-runtime-wiring-plan.md:170-203, 258-272`). The protected inventory is intentionally an operator-installed attestation, and the runtime must not accept inventory or expected-current facts from an administrator request (`research/host-runtime-wiring-plan.md:206-217`; `.trellis/spec/backend/runtime-installation.md:46-58`).

The external runner is also a real security boundary. Direct public and declared-ingress cases must run from the approved external network location; a Docker-host loopback, server-side hairpin, or container request cannot prove the external source/peer path (`research/host-runtime-wiring-plan.md:470-481`). The client additionally enforces exact origins, TLS/DNS policy, no redirects, and no ambient proxy (`research/host-runtime-wiring-plan.md:483-527`). A browser-only replacement is unsound under the current contract: the custom ticket header normally triggers a cross-origin CORS preflight, and browser networking cannot enforce the runner's direct transport policy (`research/host-runtime-wiring-plan.md:528-537`).

The runtime's own bootstrap behavior reinforces this split. Missing inventory/proof leaves normal Sub2API traffic available while evaluation remains not-ready (`research/host-runtime-wiring-plan.md:160-168`; `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:1071-1091`). `Begin` loads the protected artifact and can ensure private state, but it still requires the artifact and current facts before issuing tickets (`patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:1145-1166`).

### 4. Options for the no-manifest/no-probe-command journey

**Option A: binary-only install plus browser activation (recommended baseline).** The one command performs container discovery and binary replacement only. The browser setup page owns menu merge, Key initialization, status, and any later probe orchestration. Internal files remain generated and protected; users never edit JSON or copy a JWT. This is implementable without changing K1-K7, but evaluation remains `not_ready` until topology evidence exists.

**Option B: integrated one-command setup with a preconfigured external verifier.** The installer/browser flow starts a short-lived verifier automatically. The verifier can be a deployment-provided runner or a separately managed service whose origin, runner identity, and trust policy are already configured. The user sees a single activation action and no JWT file, JSON manifest, or probe command. This requires a new trust/deployment contract for the verifier and still needs an authoritative source of complete ingress coverage; it cannot be made universal from Docker inspection alone.

**Option C: guided generation of internal inventory plus one-shot probe.** The setup UI or command collects the public origins/hops/runners, generates the protected inventory internally, and invokes the probe client as part of the same flow. This removes hand-edited JSON and a separately typed probe command, but it still requires truthful operator input (or preconfigured deployment metadata) for every ingress and runner. It must not claim completeness from generated defaults.

**Option D: browser-only or container-only probing.** Reject for the current acceptance boundary. Browser CORS/preflight and transport controls do not provide the required runner guarantee; container/Docker-host requests cannot prove public source and proxy behavior. Accepting either would weaken K4 and contradict the existing topology evidence contract.

### 5. Rollback/uninstall remains a separate authentication design point

The current rollback/uninstall sequence drains evaluations, calls authenticated deactivation, verifies exact owned-Key revocation, removes owned menu state, and only then restores the original binary (`sub2api_install.py:413-433`; `research/private-key-lifecycle.md:135-151`). Removing JWT files from installation does not automatically provide an authenticated destructive flow. A browser-owned “deactivate and prepare restore” operation could persist an owner-bound completion marker, after which a host command performs only the image-bound restore; that is a new protocol and needs explicit design/testing. Until then, retaining a credentialed rollback path is safer than silently treating local binary restoration as proof of K5 cleanup.

## Proposed simplest feasible journey

1. Run one installer command. Auto-select only one supported container; fail closed on ambiguity. Perform version/platform/image/binary checks, atomic replacement, restart and local readiness checks. Generate protected transaction facts internally.
2. Print/open the direct same-origin evaluation URL. The administrator uses the existing Sub2API login; no JWT file or URL token is involved.
3. The page performs authenticated activation and menu merge. It creates/reuses internal private-Key state only through normal authorization and reports conflicts without repair.
4. If a deployment has a preconfigured external verifier/runner, activation invokes it automatically. Otherwise the page reports that topology evidence is missing and leaves evaluation disabled; it must not synthesize a ready state from local Docker facts.
5. Keep rollback/uninstall as a separately designed browser-authorized deactivation plus host restore transaction, preserving exact owned-Key revocation and image/checksum checks.

## Critical caveats / follow-up decisions

- A one-command installer can remove explicit container selection only when exactly one compatible target is discoverable; multiple supported containers still require an explicit choice.
- Public origin and complete ingress coverage cannot be safely guessed from container metadata. They must come from a trusted deployment integration, preconfigured verifier, or operator input collected inside the guided flow.
- Browser execution cannot replace the external runner under K4 without changing the contract.
- The highest-value product decision is whether to provide/approve a deployment-managed external verifier (Option B) or accept a guided operator-provided topology step (Option C). Until that decision, Option A should install successfully but keep evaluation not-ready.
- No live target, external runner, complete ingress inventory, production browser session, or real billing path was available for this research. No tests/builds/deployments were run.

## External references

- Sub2API v0.2.15, pinned source commit `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` as recorded in `prd.md` (the repository's compatibility target; no external network lookup was required for this report).

## Related specs

- `.trellis/spec/backend/runtime-installation.md`
- `.trellis/spec/backend/deployment-evidence.md`
- `.trellis/spec/backend/host-coordinator.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/host-runtime-wiring-plan.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/patch-contracts.md`
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/private-key-lifecycle.md`
