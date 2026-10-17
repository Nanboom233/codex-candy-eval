# Application probe transport failure

## Scope and checkpoint

Read-only research handoff, 2026-10-13. No product edit, test/source build, repack, production access, paid request or target application write occurred. The previous reviewed universal archive remains526230 bytes /62 regular files /SHA256 `b39fe125b06da0363f100709b89a0247a7c630e07526af7a6c523b0dee0eaab2` until a separately reviewed fix.
Original researcher observed empty output/PIPE_RC0 on a closed localhost port and a successful immediate Node response, then became unavailable. Original implementer continued actual transport observations in the existing owned `codex-research-http-probe` fixture (`node:24-alpine`, BusyBox1.37.0), without a new pull/dependency. A disposable additional TCP process ran only inside that fixture, with code passed directly through argv/stdin, no files written. `docker rm -f codex-research-http-probe` completed exit0 after observations. No other Docker resource was touched.
Raw evidence: `/tmp/codex-application-probe-research-20261013/{pipeline.json,direct-nc.json,framed-nc.json,framing-disconnect.json,tools.json,source-config.json}`. Source reference snapshots are in that research root. Keep it until independent review.

## Proven fault

Current owner: `install.py:352-384`, `application_capabilities()`. The selected-container shell runs `printf ... | nc -w 3 127.0.0.1 PORT | head -c 8193`, then host Python parses stdout with HTTPResponse. The pipeline has two independent transport defects:

1. BusyBox nc exits on stdin EOF before an asynchronous response arrives. The producer printf immediately closes stdin. This is a response race, not proof of absent application capabilities.
2. Shell pipeline status is head's status. Closed-port nc exit1 or missing nc exit127 becomes shell exit0. Python then sees empty stdout and raises RemoteDisconnected; final catch synthesizes “lacks native health and setup-status capabilities”, replacing the transport failure.

The actual Python literal and shell printf were tested, not inferred. Bytes are exactly `GET /health HTTP/1.1\r\nHost: 127.0.0.1\r\nConnection: close\r\n\r\n` (hex starts474554 and ends0d0a0d0a). CRLF escaping is correct.

| Actual executable case | Observed result |
| --- | --- |
| Current pipeline, immediate response on127.0.0.1:18110 | exit0,178 raw bytes, valid HTTP200/chunked; success can race |
| Current pipeline, normal1.5-second delayed response | exits0 in0.045 seconds, stdout/stderr empty |
| Current pipeline, normal4-second delayed response | exits0 in0.046 seconds, stdout/stderr empty |
| Current pipeline, closed18101 | exits0 in0.039 seconds, stdout/stderr empty |
| Current pipeline with nonexistent nc name | exits0, empty stdout, raw `missing_nc: not found` stderr |
| Current pipeline with nonexistent head name | exit127, raw missing_head stderr |
| Direct `docker exec -i ... nc -w3 ...`, communicate(request) | closes stdin too early; even immediate response can exit0 empty |
| Same direct nc, stdin kept open | receives immediate/delayed complete response; process/pipes can remain open until caller closes stdin, so waiting only for process exit is insufficient |
| Host bounded framing collection, immediate chunked | complete178 bytes,0.226 seconds, close stdin then exit0 |
| Host bounded framing collection,1.5-second chunked | complete177 bytes,1.579 seconds, close stdin then exit0 |
| Host bounded framing collection,4-second delay with3-second deadline | no response by3.043 seconds; explicit deadline outcome, never capability rejection |
| Direct nc, closed18101 | exit1, empty stderr preserved; no head to mask status |
| Direct nc, fixed Content-Length15 body | exact full200 header/body retained, exit0 |
| Direct nc, peer disconnect before response | exit0/empty; classify actual RemoteDisconnected, not missing capability |
| Direct nc, malformed status line | raw `NOT HTTP` retained, parser failure required |
| Direct nc, declared Content-Length99 with15-byte body | raw truncated reply retained; parser IncompleteRead or transport deadline required, never successful typed capability |

Keeping stdin open without response framing is not by itself a fix: BusyBox nc may stay alive after receiving the complete response because stdin remains open. The research framed collector parses only supported complete Content-Length/chunked responses before closing stdin and draining/waiting. Source constant JSON handlers ordinarily use framed HTTP responses; close-delimited response behavior needs an explicit bounded failure contract or a proven client path, not an idle-delay heuristic presented as success.

## Available HTTP tools and portability

The fixture has `/usr/bin/nc` and `/usr/bin/wget`, no curl. BusyBox nc help supports ordinary `nc HOST PORT` and `-w SEC`; the host can enforce deadlines independently and omit optional -w to avoid incompatibility with other nc applets. Do not silently assume nc exists; preserve direct launch/shell `command -v` failures and choose a fallback only after actual capability observation.
BusyBox wget supports `-S -O - -T SEC -Y off`, but follows redirects by default and does not advertise a no-redirect flag. Actual `wget --max-redirect=0 ...` exited1 with its full unrecognized-option usage text. Therefore a blanket wget replacement would violate the existing no-redirect contract. GNU-specific flags cannot be assumed available. No dependency or helper-image installation is required for the proven direct-nc approach.

## Address and port discovery gap

Current code probes only SERVER_PORT (or8080),8080 and Docker Config.ExposedPorts, at most8 ports, always127.0.0.1. It ignores SERVER_HOST, YAML config, actual listening sockets and healthcheck URL hints. These are separate possible causes of a wrong endpoint; the user's actual deployment endpoint remains unknown and no production access was attempted.
Source-backed reference, v0.2.15 `backend/internal/config/config.go`:

- Lines1830-1836: YAML config plus Viper AutomaticEnv with dot→underscore mapping.
- Lines2018-2031: explicit CONFIG_FILE, otherwise DATA_DIR, /app/data, working directory, ./config and /etc/sub2api search.
- Lines2040-2041: defaults0.0.0.0:8080.
- Lines3827-3847: GetServerAddress documents and implements environment > config file > defaults using server.host/server.port.

Thus a valid server can run on a YAML-defined port without SERVER_PORT or updated EXPOSE, or bind a specific address/IPv6 rather than127.0.0.1. Dockerfile lines152-157 expose8080 but its healthcheck uses localhost:${SERVER_PORT:-8080}; this is a hint, not authoritative actual configuration. Compose publishes host `${SERVER_PORT:-8080}:8080` while setting container SERVER_PORT8080 (lines29,49-50); published host ports must not be mistaken for in-container listener ports.
The fixture's actual `/proc/net/tcp` LISTEN row was `0100007F:46BE` (127.0.0.1:18110), with no SERVER_PORT or exposed-port hint. `/proc/net/tcp6` had no listener. Read-only listener discovery can identify real local candidates without reading the entire credential-bearing config.yaml or executing arbitrary healthcheck commands. Exact native health/setup responses, not a socket or image name, establish functionality.
No installed-binary checksum/metadata, version or source/image matching gate should be added to solve this transport/address fault.
Known source URLs inspected directly: https://raw.githubusercontent.com/Wei-Shaw/sub2api/v0.2.15/backend/internal/config/config.go , /Dockerfile and /deploy/docker-compose.yml. The guessed backend/internal/server/server.go path returned404; no conclusion relies on it. Only relevant source excerpts were reported; no private configuration was read/dumped.

## Narrow recommended implementation

Choose host Python controlling direct `docker exec -i <selected immutable ID> nc <local address> <port>` with stdout/stderr separate and stdin kept open until the HTTP reply is complete or a deadline/process failure occurs. Avoid shell printf/head pipelines and unsupported nc options. Existing host Python supplies polling/selectors, framing parsing and bounded collection; no target file writes, new package or helper pull is needed.

1. Check selected immutable ID before probe; discover executable transport capability without hiding failures. Send exact static GET /health and /setup/status requests; use no proxy, no redirect and no external hostname.
2. Bound the entire client execution and raw response to8192 bytes plus a1-byte overflow sentinel. Stream stdout with host-side limits, and preserve complete original stderr independently. Do not let a malicious large response become an unbounded capture_output allocation. Keep decoded body ≤4096 bytes and exact typed HTTP200/native JSON admission checks.
3. Recognize complete Content-Length and chunked framing incrementally. On completion close stdin, drain within deadline, wait/terminate the owned Docker exec client as necessary. Empty/disconnected/malformed/incomplete replies preserve actual parser/transport exception; nonzero process exit uses CalledProcessError with original argv/stdout/stderr; deadline uses TimeoutExpired with captured streams. Clean local process resources without replacing the primary failure.
4. Collect bounded validated in-container candidate addresses/ports from SERVER_HOST/SERVER_PORT and optional actual /proc/net/tcp/tcp6 LISTEN rows, then safe Docker exposed/healthcheck URL hints/defaults. Do not execute healthcheck strings or dump whole YAML/env. Prefer actual listening candidates, handle wildcard bind through appropriate local loopback, validate same-container specific bind addresses. Keep finite overall candidates/time budget and require both exact native endpoint contracts.
5. Distinguish actual valid HTTP with missing/wrong native fields from inability to reach/read HTTP. When all attempts fail transport, re-raise concrete underlying exception (preserving prior raw attempts/chains) instead of fabricated “lacks capability”. Only observed incompatible native responses can justify a functional-capability validation error. Keep ordinary progress and the user's complete raw stdout/stderr/traceback rule; no Cause/Recovery/Next interpretation block.

## Required verification and limits

Add real disposable socket/subprocess tests, beyond the current mocks that return preassembled HTTP bytes: immediate/1.5-second response, fixed/chunked/fragmented replies, disconnect, slow timeout, malformed/truncated/oversize response, nonzero/missingtool exit, stdin lifecycle, process cleanup, redirects rejected, and configured nondefault/bound listener discovery. Keep production payload free of fixtures/tests and assert no target writes. Exercise final packaged CLI through its actual transport boundary before independent checker review/repack.
Research proves this BusyBox1.37 EOF/pipeline failure and a working host-controlled framed response path, not the user's remote daemon/tool version or exact configuration. Native IPv6/specific-bind execution was not tested in this research; source and listener discovery make those explicit implementation/check cases. No end-to-end installation, native amd64/iStoreOS, deployment, billing or publication acceptance is claimed. Report is ready for main review; product edits wait for its dispatch.
