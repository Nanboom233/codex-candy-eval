# Research: Operator Probe Origin Boundary

- Query: Determine whether accepted topology contracts require globally public external probe origins, assess RFC1918/ULA behavior, DNS pinning and runner dispatch, and identify focused regressions.
- Scope: internal
- Date: 2026-10-12

## Findings

### Accepted boundary

The approved design repeatedly describes direct-public and declared-ingress destinations as “public origins” and assigns them to a named external runner: `host-runtime-wiring-plan.md:232-235,483-498,510-517,523-530`; PRD A16 requires public-ingress fixtures (`prd.md:76`). The protected inventory also calls the ingress field `public_hostname` and rejects local DNS suffixes (`production/005-codex-eval-deployment-evidence.patch:314-344,411`).

The operative address rule is narrower, however. The host plan says “non-loopback destination” and rejects literal loopback addresses, while `deployment-evidence.md:98-107,174` rejects loopback/unspecified literals and local hostnames. Neither rule prohibits RFC1918 IPv4 or RFC4193 IPv6 ULA. No accepted contract says private origins are required, but the detailed non-loopback rule currently permits them when the named runner can reach them. Thus a private-only origin is not affirmatively approved as a use case, yet rejecting it would narrow the literal accepted rule. Generic “public origin” wording alone does not settle that CIDR policy. Treat an explicit RFC1918/ULA prohibition as optional public-only hardening pending a product decision; do not present it as a required fix under technical implementation authority.

The specific concern is real if public-only means globally reachable: a declared HTTPS origin resolving to `10/8`, `172.16/12`, `192.168/16`, or `fc00::/7` can receive the one-use probe ticket from an external runner with private routing. TLS hostname verification, disabled redirects, and a named runner do not prove the destination is public. This is a compatibility-policy question under the current non-loopback rule, not a demonstrated bypass of accepted case/source/hop checks. If the product decision is to enforce public-only destinations, scope the check to externally executed probe origins in both inventory validation and the `client(true)` worker. Do not apply it to the admin API client (`client(false)`), private proxy-hop peers, or other internal topology fields. `net.IP.IsPrivate` covers exactly RFC1918 and RFC4193; it is not a complete globally-routable-address policy. Broader IANA special-purpose filtering would be a separate scope decision.

### Current implementation and paths

- `patches/sub2api/v0.2.15/production/005-codex-eval-deployment-evidence.patch:314-344,347-375,411`: literal origins and ingress `public_hostname`/`tls_hostname` use `evalDeploymentHostname`; it excludes loopback, unspecified, scoped literals, ambiguous numeric forms, and selected local suffixes, but permits RFC1918/ULA literals.
- `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:63-90,98-170`: the probe transport sets `Proxy: nil`, resolves each hostname, rejects an answer set if any answer is loopback, unspecified, link-local unicast, or scoped, then dials each already-resolved IP literal. It does not reject RFC1918/ULA.
- `dist/sub2api-build/workspace/backend/cmd/codex-eval-probe/main.go:60-93,98-170`: staged-source equivalent of the production worker; `origin` checks URL syntax, while the public destination policy is enforced in `client(true)` during `runTickets`.
- `patches/sub2api/v0.2.15/test-overlay/007-codex-eval-runtime.patch:18-50` and `tests/host-overlay/backend/cmd/codex-eval-probe/main_test.go:15-48`: current origin/client tests cover malformed origins, loopback rejection, and redirects; they do not cover private literals, DNS answer sets, or private-range answers.
- `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch:269-307` and staged source `main.go:269-307`: tickets are grouped by named runner. The local runner calls `runTickets`; other runner groups are handed as JSON stdin to an operator-configured `runner-commands` argv. The dispatcher is not required by code to run this binary's `--ticket-input` mode. The payload contains tickets and inventory but not the administrator JWT. This is an operator trust/attestation boundary; no exploit of the accepted proof flow is established without an untrusted or misconfigured dispatcher.

### DNS and IP-family behavior

For a worker that runs `runTickets`, all DNS answers are checked before any dial attempt. Under the current policy, any loopback, unspecified, link-local-unicast, or scoped answer rejects the full set; RFC1918 and ULA answers are permitted. If public-only CIDR policy is approved, `IsPrivate` must likewise reject a mixed public/private set before dialing. The dial loop uses each checked answer as a numeric address (`net.JoinHostPort(ip.IP.String(), port)`), so the dialer does not resolve the hostname a second time between validation and connection. Every new connection repeats lookup and checking. Default HTTP TLS verification still uses the request hostname. This avoids resolution-to-dial rebinding in this transport.

The existing `net.IP` predicates handle IPv4-mapped forms: Go 1.27.2 `IsLoopback`, `IsPrivate`, and `IsLinkLocalUnicast` normalize through `To4`; `IsUnspecified` uses `Equal`, which treats mapped IPv4 as equal to its IPv4 form (`/opt/homebrew/Cellar/go/1.27.2/libexec/src/net/ip.go:119-153,181-185,392-405`). Thus adding `ip.IP.IsPrivate()` would also reject mapped RFC1918 addresses; current unspecified and loopback dial checks already reject mapped `0.0.0.0` and `127/8` answers. Link-local zones are separately rejected.

There is a narrower reader inconsistency under the accepted literal rule: `evalDeploymentHostname` calls `addr.Unmap().IsLoopback()` but calls `addr.IsUnspecified()` without unmapping. In Go 1.27.2, `netip.Addr.IsUnspecified` compares only native `0.0.0.0` and `::`, whereas `Unmap` explicitly converts IPv4-mapped IPv6 (`/opt/homebrew/Cellar/go/1.27.2/libexec/src/net/netip/netip.go:477-485,666-672`). Consequently `[::ffff:0.0.0.0]` can pass inventory hostname validation even though it denotes the unspecified IPv4 address; `client(true)` later rejects it because `net.IP.IsUnspecified` normalizes mapped addresses. This is a fail-later validator bug, not evidence that RFC1918/ULA rejection is already part of the contract. Normalize once (`addr = addr.Unmap()`) before the accepted loopback/unspecified checks and add a reader regression for mapped unspecified.

Do not substitute `IsGlobalUnicast` for a private-range check: Go documents that it returns true for private IPv4 and local IPv6 unicast (`net/ip.go:188-203`). Do not expand the current scope to every special-purpose/reserved range without an explicit policy.

### Recommendation and regressions

1. Under the currently accepted contract, preserve RFC1918/ULA reachability; record public-only CIDR rejection as optional hardening requiring an explicit decision. If the product owner confirms these cases must prove an Internet-public ingress, reject private literal origins in the inventory validator and reject `IsPrivate()` on every DNS answer in the public worker, before dialing. Keep `client(false)` and internal hop address handling unchanged.
2. Fix the definite mapped-unspecified reader mismatch by unmapping parsed IP literals before testing loopback/unspecified; add cases for `[::ffff:0.0.0.0]` and `[::ffff:127.0.0.1]`. Keep existing IPv4/IPv6 unspecified, loopback, link-local and zone checks.
3. Add deterministic resolver/dialer coverage: check every answer before dialing; under current policy a public-plus-loopback/unspecified/link-local/scoped set must fail without a dial; the dial target must be the checked numeric answer. A public-plus-private set must fail without a dial only if public-only policy is separately approved; mapped private follows that same conditional rule. Keep tests local and avoid public-network probes.
4. Keep the runner-command trust explicit in operator acceptance: the configured nonlocal command is expected to invoke the approved versioned `--ticket-input` worker on the named runner, whose transport applies the same origin checks. If the implementation claims this as a code-enforced guarantee, arbitrary `runner-commands` do not establish it and need a separate design decision. Preserve case/origin/ticket matching and do not pass the admin JWT to child commands.

## Files Found

- `.trellis/tasks/10-11-sub2api-runtime-patch/research/host-runtime-wiring-plan.md`: approved named-runner, public-origin, destination, and probe acceptance plan.
- `.trellis/spec/backend/deployment-evidence.md`: accepted inventory reader contract; literal address rule is loopback/unspecified, with DNS policy delegated to the runner.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md`: K4/A16 topology and public-ingress acceptance boundary.
- `patches/sub2api/v0.2.15/production/005-codex-eval-deployment-evidence.patch`: origin and ingress-host validation implementation.
- `patches/sub2api/v0.2.15/production/007-codex-eval-runtime.patch`: operator client, DNS resolver/dialer, ticket worker, and runner dispatch.
- `patches/sub2api/v0.2.15/test-overlay/007-codex-eval-runtime.patch`: operator-client test additions.
- `tests/host-overlay/backend/cmd/codex-eval-probe/main_test.go`: maintained operator-client test source.
- `dist/sub2api-build/workspace/backend/cmd/codex-eval-probe/main.go`: staged copy of the operator-client source used for line-level confirmation.
- `/opt/homebrew/Cellar/go/1.27.2/libexec/src/net/ip.go` and `net/netip/netip.go`: local Go 1.27.2 address predicate implementations.

## Code Patterns

- Inventory hostname validation accepts literal IPs by `netip.ParseAddr`, then checks `Zone`, `Unmap().IsLoopback()`, and raw `IsUnspecified` (`production/005...patch:314-320`). It is the correct place for the accepted literal-address checks; public-only private rejection would be a separate policy addition.
- The worker checks every `LookupIPAddr` result before the dial loop and dials numeric address strings (`production/007...patch:70-82`). Add any public-only policy at this point so every answer, including mapped answers, is handled before a connection.
- `runTickets` verifies exact inventory bytes, runner-to-case coverage, ticket/origin match, and approved TLS port before posting `/v1/responses` (`production/007...patch:98-170`). Keep this binding intact.
- The dispatcher uses `exec.CommandContext` with configured argv and the ticket-bearing work body on stdin (`production/007...patch:269-307`). It is not itself a DNS-checked transport.
- Current tests assert loopback cannot be reached by `client(true)` and that redirects are not followed (`main_test.go:26-48`); no RFC1918/ULA or DNS-answer test currently resolves this ambiguity.

## External References

- No external references were fetched. Local standard-library docs/source are Go 1.27.2 (`go doc net.IP.IsPrivate`, `go doc net/netip.Addr.IsPrivate`; source paths listed above). They define `IsPrivate` as RFC1918 IPv4 plus RFC4193 ULA and document that it is not a complete security policy.

## Related Specs

- `.trellis/spec/backend/deployment-evidence.md:98-107,174`: current normative reader and runner split.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/deployment-evidence-reader-plan.md:174`: earlier exact-origin implementation plan.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/host-runtime-wiring-plan.md:232-235,483-530,1018-1047`: named external runner, public-origin intent, narrower non-loopback rule, and acceptance tests.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md:50,76`: K4 and A16.

## Caveats / Not Found

- The accepted text mixes “public origin” intent with an operative non-loopback rule and does not resolve RFC1918/ULA destinations. Private origins are not explicitly required or approved, but the detailed contract currently permits them; public-only rejection is therefore optional until clarified.
- The arbitrary nonlocal runner command makes the checked-worker guarantee depend on operator configuration; no test proves dispatch invokes the approved `--ticket-input` worker. This is a trust-boundary caveat, not a demonstrated bypass under the accepted operator-controlled configuration.
- No tests, builds, network probes, code edits, spec edits, or implementation changes were performed for this research.
