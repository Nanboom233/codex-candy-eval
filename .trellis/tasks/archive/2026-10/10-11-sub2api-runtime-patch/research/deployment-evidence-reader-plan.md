# Research: Protected Deployment Evidence Reader

- Query: Define the isolated Sub2API patch 005 reader for the protected ingress inventory, including its schema, current-fact inputs, filesystem checks, and focused tests.
- Scope: internal
- Date: 2026-10-12

## Findings

### Source boundary

Patch 004 is the metadata adapter. The next isolated host change should add only `backend/internal/service/eval_deployment_evidence.go` and `backend/internal/service/eval_deployment_evidence_test.go`, exported as `patches/sub2api/v0.2.15/005-codex-eval-deployment-evidence.patch`. The implementation reads the fixed `filepath.Join(store.Dir(), "deployment-evidence", "inventory.json")` path. The deployment operator produces and installs this inventory; patch 005 is its read-only consumer. A later trusted runtime/deployment provider owns collection and validation of the current facts passed to the reader. Patch 005 does not create the directory, write evidence, cache mutable data, change readiness, or add routes, tickets, Wire providers, bootstrap, runner, or lifecycle behavior.

`EvalKeyStore.Dir()` is available at `eval_key_runtime.go:80-109`. The existing `ensureDir` creates the path and checks only the final component's directory type and `0700` permission bits. It does not check ownership or parent components. `fileModeOK` at `:112-120` rejects directories and symlinks and checks permission bits, but it does not require a regular file or verify ownership. Do not treat either helper as a complete protected-inventory check, and do not call `ensureDir` from this read-only loader.

The live security configuration can be compared today: `EvalSecuritySnapshotFromConfig` captures forwarded-IP trust, configured headers, trusted proxies, and the configured flag, and `Fingerprint()` hashes the normalized snapshot (`eval_topology.go:45-85`). The collector supplies the exact required case set from the validated loopback sources, ingress peers, and recognized headers (`:292-325`, `:464-515`). `EvalIngressTopology.Validate` already canonicalizes loopback sources and binds ingress IDs to the sealed proof (`eval_key_runtime.go:336-362`).

Other required current facts are not exposed by the pinned runtime. `service.BuildInfo` and the handler build info carry only `Version` and `BuildType` (`service/wire.go:30-34`, `handler/handler.go:73-77`, `cmd/server/main.go:146-149`). No current field supplies the running source commit, runtime patch-set identity, deployment ID, host image index/platform/config digest, binary checksum, or ingress-configuration revision. The reader must take these expected values as explicit current facts from a later trusted provider and fail closed when any required value is unavailable. Do not treat an artifact's self-declared value as the current value, infer an image digest from the version, or invent a signing key. The security fingerprint is the only listed revision fact directly computable from the current config.

### Frozen exported inventory contract

Expose exactly these exported Go types in `package service`:

```go
type EvalDeploymentCurrentFacts struct {
	Version                  string
	BuildType                string
	DeploymentID             string
	SourceRevision           string
	RuntimePatchSHA256       string
	ImageIndexDigest         string
	PlatformManifestDigest   string
	ImageConfigDigest        string
	BinarySHA256              string
	IngressRevision          string
}

type EvalDeploymentInventory struct {
	SchemaVersion          int                                  `json:"schema_version"`
	ArtifactID             string                               `json:"artifact_id"`
	Version                string                               `json:"version"`
	BuildType              string                               `json:"build_type"`
	DeploymentID           string                               `json:"deployment_id"`
	SourceRevision         string                               `json:"source_revision"`
	RuntimePatchSHA256     string                               `json:"runtime_patch_sha256"`
	ImageIndexDigest       string                               `json:"image_index_digest"`
	PlatformManifestDigest string                               `json:"platform_manifest_digest"`
	ImageConfigDigest      string                               `json:"image_config_digest"`
	BinarySHA256           string                               `json:"binary_sha256"`
	SecurityFingerprint    string                               `json:"security_fingerprint"`
	IngressRevision        string                               `json:"ingress_revision"`
	CreatedAt              int64                                `json:"created_at"`
	ExpiresAt              int64                                `json:"expires_at"`
	OperatorAuditReference string                               `json:"operator_audit_reference"`
	LoopbackSources        []string                             `json:"loopback_sources"`
	Ingresses              []EvalDeploymentIngress              `json:"ingresses"`
	RunnerLocations        []EvalDeploymentRunnerLocation       `json:"runner_locations"`
	CaseDestinations       []EvalDeploymentCaseDestination      `json:"case_destinations"`
}

type EvalDeploymentIngress struct {
	ID                 string                    `json:"id"`
	Source             string                    `json:"source"`
	PublicHostname     string                    `json:"public_hostname"`
	ListenerHost       string                    `json:"listener_host"`
	ListenerPort       int                       `json:"listener_port"`
	CanonicalPath      string                    `json:"canonical_path"`
	RewritePath        string                    `json:"rewrite_path"`
	TLSHostname        string                    `json:"tls_hostname"`
	TLSPort            int                       `json:"tls_port"`
	TLSPolicy          string                    `json:"tls_policy"`
	Hops               []EvalDeploymentProxyHop `json:"hops"`
}

type EvalDeploymentProxyHop struct {
	ID                 string `json:"id"`
	Kind               string `json:"kind"`
	Listener           string `json:"listener"`
	NextHop            string `json:"next_hop"`
	ObservedPeerSource string `json:"observed_peer_source"`
	TLS                bool   `json:"tls"`
}

type EvalDeploymentRunnerLocation struct {
	ID                  string   `json:"id"`
	ExpectedSourceCIDRs []string `json:"expected_source_cidrs"`
	ClientVersion       string   `json:"client_version"`
	AuditReference      string   `json:"audit_reference"`
}

type EvalDeploymentCaseDestination struct {
	Case             EvalProbeCase `json:"case"`
	Origin           string        `json:"origin"`
	RunnerLocationID string        `json:"runner_location_id"`
	ApprovedTLSPort  int           `json:"approved_tls_port"`
	HopIDs           []string      `json:"hop_ids"`
}

type EvalDeploymentEvidence struct {
	Inventory EvalDeploymentInventory
	SHA256    string
}

func (s *EvalKeyStore) LoadDeploymentEvidence(
	cfg *config.Config,
	current EvalDeploymentCurrentFacts,
	expectedArtifactID string,
	expectedSHA256 string,
) (EvalDeploymentEvidence, error)
```

These exported names, fields, and the `LoadDeploymentEvidence` signature are
the concrete patch-005 interface. Every JSON member is required, including
`approved_tls_port`; zero means only port 443 is permitted. Arrays may be empty
only where the topology allows it. `source_revision`,
`runtime_patch_sha256`, image-index/platform/config
digests, and `binary_sha256` are separate identities; the deployment review
warns against treating them as interchangeable (`research/patch-contracts.md`,
Runtime Installer Transactions). The inventory contains no self-digest.
`LoadDeploymentEvidence` hashes the exact bytes and returns the lowercase hex
digest separately, compares it with `expectedSHA256`, and compares
`expectedArtifactID` with the file's ID. Both expectations are required.

The later trusted runtime/deployment provider constructs
`EvalDeploymentCurrentFacts` from current sources and owns their provenance and
validation. The reader compares every field in that value with the
corresponding inventory field on every load. It computes
`SecurityFingerprint` from the passed live config and compares it with the
inventory; it never accepts a caller-supplied current fingerprint. Missing
current facts fail closed. `Version` and `BuildType` are available from the
current host, but the other facts require a later trusted deployment/runtime
provider. The inventory file, admin request, and caller-supplied artifact
metadata are never sources for expected current facts. The `expectedArtifactID`
and `expectedSHA256` arguments bind this read to the operator/runner's selected
artifact; both must be non-empty, and the SHA-256 expectation must be lowercase
hex for exactly 32 bytes.

The topology portion contains canonical `loopback_sources`, the ingress
inventory, runner locations, and one case-destination record per external
probe. Each ingress's `ID` and `Source` project to
`EvalIngressPeer{ID, Source}`. The listener, host/path rewrite, TLS policy, and
ordered hop list provide the detail needed for independent deployment review.
The final hop's observed peer source must equal the ingress `Source`. Each
external case binds the full `EvalProbeCase` tuple to one HTTPS origin, one
runner-location ID, and its ordered hop IDs. A nonzero `approved_tls_port`
permits only that explicit port for the case. Runtime-derived direct-loopback
origins remain outside the file.

The artifact producer must emit all four case-tuple members (`kind`,
`ingress_id`, `header`, and `source`), including explicit empty strings when a
member does not apply. The host's `EvalProbeCase` JSON tags use `omitempty` for
some fields, so `json.Marshal(EvalProbeCase)` alone can omit required empty
members. The later operator/installer producer must use a dedicated artifact
encoder/DTO or explicit map that writes the complete tuple; the consumer keeps
the required-member validation unchanged.

After decoding and canonicalizing sources/peers with existing helpers, create
`EvalTopologyCollector` from the live config and inventory topology. Compare
`CaseDestinations` with `RequiredCases()`: `direct_allowed` cases are local-only
and have no external destination; `direct_public`, `spoofed_loopback`,
`ingress_allowed`, and `ingress_spoofed` each require exactly one destination.
Reject duplicate, missing, extra, or mismatched cases. This derives coverage
from current recognized headers and the artifact's ingress IDs.
Each destination must reference a known, valid runner location; runner IDs must
be unique, but unused valid runner descriptors are allowed and need not be
referenced by a case.

### File and decode validation

- Require an initialized store and an absolute, non-empty store root. Walk every component from `/` through the store and `deployment-evidence` using directory file descriptors and anchored `openat` calls with `O_DIRECTORY|O_NOFOLLOW`; reject symlinks and non-directories at every component. Require the private store root and evidence directory to be owned by the effective runtime UID and have mode `0700`; system ancestors require only safe directory/no-symlink checks. Open the fixed inventory name relative to the evidence-directory descriptor with `O_NOFOLLOW`; require the opened file to be owned by that UID, regular, and mode `0400` or `0600`. Use `fstat` on the opened descriptors for these checks. Do not create missing parents or permit production symlink exceptions.
- Read at most 1 MiB using `io.LimitReader(limit+1)` and reject oversize input before parsing. Freeze the required topology-case count at 256; reject a live `RequiredCases()` set above that bound and require the artifact's external destinations to match the derived set exactly. These are fixed safety limits, not measurements of one deployment.
- Hash the exact accepted byte slice with SHA-256 before JSON parsing. Do not normalize, re-marshal, or embed a digest in the JSON. A whitespace-only change must produce a different digest.
- First validate the raw JSON structure against the schema: require a top-level object; reject duplicate object members at every nesting level, every `null` value (including array elements), and every missing required member; then decode with `json.Decoder.DisallowUnknownFields` and require exactly one JSON value with no trailing bytes. Validate required semantic values and schema version, reject duplicate case entries and all extra data. `DisallowUnknownFields` alone is insufficient because Go accepts duplicate members with last-value-wins behavior and does not enforce required-member presence.
- Interpret `created_at` and `expires_at` as Unix seconds. Permit creation at most five minutes in the future; require expiry strictly later than both the current time and creation time. Do not impose a maximum inventory lifetime. This is separate from the topology-proof's existing 24-hour maximum validity.
- Validate every external origin as an exact origin, not a URL prefix. Require HTTPS, a non-empty hostname, no userinfo, query, fragment, opaque form, or non-root path, and only port 443 unless the inventory explicitly declares the approved TLS port. Reject literal loopback destinations and local hostnames. The one-shot runner later owns DNS resolution checks, trusted TLS roots/hostname verification, redirects-off, ambient-proxy-off, and destination credential isolation; the reader performs no network requests.
- Return detached inventory values. Copy slices, nested hop lists, and maps before returning. Avoid a cache; if later added, return a deep clone on every read.

The protected path is an installation boundary, not cryptographic provenance. A correctly owned file and audit-reference string do not prove that the operator's source export was complete or authorized. The deployment operator must independently inspect every public listener, DNS/CDN/WAF/load-balancer/reverse-proxy hop, rewrite, and peer address before installing it. The collector cannot discover omitted hops or independently distinguish public routes that arrive from the same peer.

### Focused patch and test gates

Keep the source ownership to the reader and its adjacent test. Use the existing `unit` build-tag and `testify/require` style from `eval_key_runtime_test.go`; create temporary protected roots and explicitly set directory/file modes. Canonicalize temporary test roots (for example with `filepath.EvalSymlinks`) before applying production checks, since macOS temporary paths can include `/var` aliases. Production must still reject symlinks in every absolute path component. No test may call a public host or run a model request.

- Fixed-path behavior: missing root, evidence directory, or file fails; the reader never accepts a caller-provided path or creates missing storage.
- Filesystem safety: reject a symlink at every absolute root component, the evidence parent, or inventory file; reject non-directory parents, wrong owner/mode, FIFO/device/socket destinations, and non-regular files. Exercise owner mismatch through a testable owner-check helper so the test does not require `chown` privileges.
- Bounded strict JSON: test the 1 MiB limit, invalid/truncated JSON, top-level and nested `null`, unknown nested fields, missing required members, trailing values/bytes, duplicate members at every nesting level, invalid schema version, missing/zero required bindings, malformed timestamps, creation beyond five-minute skew, expiry not later than now or creation, and acceptance of long validity periods.
- Current facts: change each source/runtime/patch/image/binary/deployment/ingress/config binding independently and assert rejection. Empty current facts must fail closed. Verify the live `EvalSecuritySnapshot` fingerprint changes when forwarded-IP trust, headers, trusted proxies, or the configured flag changes.
- Case coverage: derive requirements from `EvalTopologyCollector.RequiredCases()` and test exact success, duplicate, missing, extra, wrong ingress ID/source, wrong spoof header, and a destination attached to a local-only direct case.
- Destination safety: reject HTTP, userinfo, query, fragment, path prefix, unsupported port, loopback literal, localhost, malformed host, and a case mapped to the wrong ingress/runner. Assert the reader performs no dial.
- Digest and ownership of outputs: assert exact-byte SHA-256, then add whitespace and assert the digest changes while decoded facts remain equal. Mutate every returned nested slice/map and verify a fresh read is unchanged.
- Replay `001` through `005` onto pristine pinned Sub2API `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2`; compare the new source/test files with the maintained workspace. The patch must not include a session-evidence writer, tickets, routes, middleware, Wire, bootstrap, runner, or lifecycle files.

## Files Found

- `dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime.go`: fixed store root access, current mode checks, canonicalization helpers, provisioner intent recovery, and topology identity validation.
- `dist/sub2api-build/workspace/backend/internal/service/eval_topology.go`: live security fingerprint, exact required cases, and topology collector.
- `dist/sub2api-build/workspace/backend/internal/service/wire.go`: build metadata contains only version and build type.
- `dist/sub2api-build/workspace/backend/internal/service/plugin_package.go`: existing bounded-read and strict-decoder pattern (`:223-230`, `:382-395`).
- `dist/sub2api-build/workspace/backend/internal/service/eval_key_runtime_test.go` and `eval_topology_test.go`: adjacent `unit` test helpers, temporary storage setup, and config/topology fixtures.
- `research/host-runtime-plan-final-review.md`: next-smallest-chunk recommendation and external deployment-evidence caveat.
- `research/patch-contracts.md`: identity distinctions for source, image, binary, and patch artifacts.

## Related Specs

- `.trellis/spec/backend/quality-guidelines.md`: focused adjacent tests, persistence visibility, and verification reporting.
- `.trellis/spec/backend/run-state.md`: protected storage and fail-closed corruption/recovery principles; the new inventory is a separate file owner and must not reuse run-state storage.
- `research/host-runtime-wiring-plan.md`: parent host workflow, probe cases, marker/session boundaries, and later runtime integration.

## External References

- No external documentation was needed. The reader and decoder patterns are available in the pinned Sub2API source.

## Caveats / Not Found

- The pinned host does not expose enough current runtime facts to accept the proposed artifact yet. The later deployment/runtime adapter must supply and validate source commit, ordered patch identity, deployment identity, image index/platform/config digests, running binary checksum, and ingress revision. Missing values must keep inventory status stale/unavailable.
- The source does expose the live security fingerprint and `EvalKeyStore.Dir()`, but its existing path helpers do not establish ownership, all parent safety, or regular-file type for this inventory.
- The fixed limits are 1 MiB for the inventory and 256 derived required cases. Inventory expiry has no maximum lifetime; only the five-minute future-creation skew and strict expiry-after-now-and-creation rules apply. These limits are independent of the topology proof's 24-hour lifetime.
- The inventory reader cannot establish operator authorization or completeness. Docker is available through `desktop-linux` (Docker Desktop 29.8.2, linux/arm64), but no approved live target or deployment inventory was used for this read-only research.
