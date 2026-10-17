# Research: application probe without target clients

- Query: Determine whether the installer can validate Sub2API when the selected image has no `nc`, `curl`, or `wget`; compare target clients, a daemon-side helper, and host access through published ports. Assess the reported 6185 versus configured 8080 discrepancy without attributing the listener beyond current evidence.
- Scope: mixed
- Date: 2026-10-13

## Findings

### Files found

- `install.py` — selected-container inspection, HTTP transport, owned-listener discovery, and image acquisition.
- `.trellis/tasks/10-11-sub2api-runtime-patch/prd.md` — reports the user’s exit-126 failure, claimed Compose `SERVER_PORT=8080`, and the no-client acceptance requirement.
- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md` — transport, bounded endpoint discovery, raw-error, and no-target-write design contract.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/application-probe-transport-failure.md` — prior BusyBox framing evidence and Sub2API v0.2.15 source findings.
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/interactive-source-preflight-fix.md` — source-version-scoped Sub2API configuration and native route contracts.
- `.trellis/spec/backend/universal-linux-delivery.md` — installer probe limits, immutable helper-image handling, cleanup, and required no-client acceptance.
- `.trellis/spec/backend/patch-native-installation.md` — selected-container identity, no-mutation-before-preflight, and raw diagnostic requirements.
- `.trellis/spec/backend/runtime-installation.md` — runtime and installation constraints for the selected Linux target.
- `.trellis/workflow.md` — Trellis research persistence and task-artifact workflow.

### Current behavior and endpoint selection

`install.py:420-424` starts `docker exec -i -u 0 <selected-id> nc <address> <port>`. The HTTP socket collector runs on the installer host, but the network client itself must exist in the target container. A missing executable therefore prevents HTTP inspection before response framing or JSON validation can run. The current limits are 8,192 wire bytes, 4,096 body bytes, 3 seconds per request, and 48 seconds across candidates (`install.py:59-62,420-448,551-597`).

`listening_endpoints()` reads `/proc/net/tcp` and `tcp6` and retains listening socket inodes held by the selected application executable (`install.py:490-549`). It converts wildcard binds to loopback and truncates the result to eight. When any owned endpoints exist, `application_capabilities()` tries that list as-is; only if it is empty does it try the environment/default/exposed-port hints (`install.py:551-578`). Consequently an owned listener candidate on 6185 can be attempted ahead of `SERVER_PORT=8080`. Ownership establishes a process-held listening socket, not that the port serves Sub2API’s HTTP routes.

Recommended order: validate `SERVER_PORT` as a numeric port, promote it only when it appears among discovered owned endpoints, then try the remaining bounded owned endpoints. Keep the cap after this prioritization so early unrelated owned sockets cannot crowd out the configured port. Require the same successful HTTP 200 and typed JSON pair on a candidate: `/health` with `status == "ok"`, and `/setup/status` with integer `code == 0`, boolean `data.needs_setup == false`, and `data.step == "completed"`. A socket, port hint, health-only response, or image label is not sufficient. The two endpoint contracts are source-verified for Sub2API v0.2.15 in the existing source-preflight report; they do not prove the user’s exact running version.

### Transport choices

1. **Detected target client.** Reusing a client already in the selected image avoids helper startup and image acquisition. Presence alone is not enough: validate supported flags, retain stderr/status, disable proxies and redirects, apply absolute deadlines and byte bounds, and reject non-200 responses. The earlier BusyBox research found that `wget` follows redirects by default and does not accept `--max-redirect=0`; behavior varies by implementation. The current direct `nc` path is already a target-client dependency, and the user’s reported image lacks it.

2. **One-shot helper in the target network namespace.** This works without installing tools or files in the target. The proof below used the Docker Engine’s `--network container:<target-id>` mode; a helper’s `127.0.0.1` reached a loopback-only server in the target. Use the immutable helper image ID from the existing `acquire_image()` cache-first flow (`install.py:283-320`), with `--pull=never` for execution. The existing source resolver already obtains `alpine/git:2.49.1` and uses its immutable ID for execution (`install.py:629-677`); in this local cache it contains `sh`, Git, `nc`, and `wget`. The application probe currently precedes that acquisition (`install.py:377-390`, then bundle source resolution at `install.py:795`), so a fallback must acquire lazily and pass the same inspected image ID forward, or move the shared acquisition earlier without repeating it.

   A direct raw-socket form is `docker run --pull=never --rm --network container:<selected-id> --entrypoint nc <helper-image-id> <validated-address> <port>`, with the existing host-side bounded HTTP framing/parser. Give the helper a unique owned name as well as `--rm`; if the host Docker client times out or is terminated, explicitly `docker rm -f` that exact name in cleanup because the daemon-side container may outlive its CLI process. Preserve the original Docker stdout/stderr, HTTP bytes, timeout/parser exception, and any cleanup failure. Do not use target mounts, host binds, published ports, `--privileged`, a shared PID namespace, or a Docker-socket mount. The helper only needs the target’s network namespace and must issue static, bounded GETs to discovered owned endpoints.

3. **Installer-host access through published ports.** This is conditional on an exact published-port mapping and a route from the installer host to the Docker daemon host. With a remote Docker context, the client’s loopback is not the daemon’s loopback; ports published only on the daemon or on its loopback may be unreachable from the installer host. A host-side request may also reach a different listener if the mapping or host address is assumed. It is less portable than namespace sharing and does not work for unpublished container ports.

### Disposable no-client proof

On 2026-10-13, Docker CLI 29.9.0 / Engine 29.8.2, a disposable image was built from the locally cached `node:24-alpine` image with `docker build --pull=false`; the image’s `/bin` and `/usr/bin` client links for `nc`, `curl`, and `wget` were removed. An in-image `command -v` check reported all three absent. The target ran a Node HTTP server bound only to `127.0.0.1:18385`, on `--network none`, with the two expected JSON routes. No production or user application container was used.

The helper command used the fixture’s immutable image ID and `docker run --rm --pull=never --network container:<fixture-target-id> --entrypoint node <fixture-image-id> -e <bounded-client-script>`. The script made direct Node HTTP requests with an absolute 2.5-second timeout, a 4,096-byte body bound, no proxy support, and no redirect-following behavior. It returned:

- `/health`: HTTP 200, `{"status":"ok"}`.
- `/setup/status`: HTTP 200, `{"code":0,"data":{"needs_setup":false,"step":"completed"}}`.

`docker diff <target-id>` was empty both before and after the helper. The named helper was absent after `--rm`; a shell trap removed the exact target fixture and derived image. The helper used no pull, mount, published port, or target filesystem command. The test establishes network-namespace reachability and target immutability for this Docker Engine setup; it does not establish the user’s daemon behavior or target endpoint configuration.

### Related specs and external references

- `.trellis/spec/backend/universal-linux-delivery.md:73-88,103-123,186-190` — exact HTTP contract, byte/time limits, cache-first immutable tools, cleanup, and real no-client test requirement.
- `.trellis/spec/backend/patch-native-installation.md:94-126,138-140` — preflight before replacement, selected immutable ID, and complete raw diagnostics.
- `.trellis/tasks/10-11-sub2api-runtime-patch/design.md:39-55` — transport correction, listener ownership, error preservation, and helper acquisition constraints.
- Docker Engine networking, container network mode: https://docs.docker.com/engine/network/#container-networks
- Docker CLI `docker container run` reference (`--network`, `--pull`, and `--rm`): https://docs.docker.com/reference/cli/docker/container/run/
- Local version evidence: `docker version --format 'client={{.Client.Version}} server={{.Server.Version}}'` returned `client=29.9.0 server=29.8.2`; `docker run --help` documents `--network`, `--pull=never`, and `--rm` behavior.

## Caveats / Not Found

- The user’s exact running image, `docker inspect` output, target executable list, runtime environment, and HTTP responses were not accessed. The supplied Compose value and reported 6185 attempt are not proof of the live container configuration or route behavior.
- No evidence identifies the 6185 listener beyond “owned listener candidate”; the current selection order can try it before environment hints. Whether 6185 serves either native route remains unknown.
- The namespace test used a local Docker Engine and an isolated Node fixture. Remote daemon, rootless Engine, host-network targets, and the user’s exact image/helper architecture were not exercised. Validate `container:<selected-id>` support and helper cleanup against the target Docker Engine before claiming those environments accepted.
- The successful fixture helper was Node-based. The cached `alpine/git:2.49.1` helper was separately inspected for `nc`, but was not used in the no-client network request itself. No source or production image was pulled.
- No product code, test, spec, production container, deployment, or paid request was changed or run. Full installer regression and packaged-entry validation belong to implementation/check after the transport design is applied.
