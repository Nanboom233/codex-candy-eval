# Research: Source Resolver Helper Platform and Tool Acquisition

- Query: Determine whether the current Docker helper is cacheable/runnable across platforms, whether local images provide Git, and what the resolver must validate before source resolution.
- Scope: mixed (internal code plus local Docker observations)
- Date: 2026-10-13

## Findings

### Current code contract

- `install.py:53` fixes the resolver helper to `alpine:3.23`.
- `install.py:345-369` unconditionally runs `docker pull alpine:3.23`, then starts the helper with the daemon default platform and installs Git at runtime using `apk add --no-cache git`.
- `install.py:180-186` drops Docker stdout/stderr and maps every nonzero Docker command to one generic `InstallError`, so helper acquisition cannot currently distinguish registry, platform, or daemon failures.
- `install.py:367-389` writes resolver output only after `docker run` returns and attempts a best-effort `docker rm -f`; a pull failure occurs before the resolver name and cleanup path are established.
- `.trellis/spec/backend/universal-linux-delivery.md:61-66,91-94,116-117` requires source-hint resolution before target writes, allows only optional discovery evidence, and says users provide Python and Docker while the package obtains build tools automatically. It does not authorize a target-architecture-specific helper or fixed host fallback.

### Local Docker observations

The available Docker CLI/server are 29.9.0/29.8.2. The daemon reports Linux `aarch64`; Buildx's `desktop-linux` worker advertises `linux/arm64`, `linux/amd64`, and several other emulated platforms. The image cache contains `alpine:3.23` as Linux arm64 (`Architecture=arm64`, `Variant=v8`, `Os=linux`) with a multi-platform OCI descriptor.

- `docker image inspect alpine:3.23 --platform linux/amd64` returns a remote amd64 descriptor, but `docker run --pull=never --platform linux/amd64 alpine:3.23 ...` fails because the amd64 platform variant is not locally present. Inspecting a manifest is therefore not proof that the requested platform is runnable from cache.
- `docker run --rm --pull=never alpine:3.23 ...`, `golang:1.27.2-alpine`, `node:24-alpine`, and `redis:8.4-alpine` all expose `sh` and `apk` but no `git`. The existing builder image is not a Git-capable resolver fallback.
- With `--network none`, `apk add --no-cache git` fails against Alpine repository DNS and reports `git (no such package)`. The current resolver therefore has an additional package-repository network dependency after the unconditional registry pull.
- The Buildx worker's advertised amd64 emulation does not make an amd64 variant available to ordinary `docker run`; source resolution must either use the daemon's runnable local platform or explicitly pull a helper variant for the required platform.

### Design implications

1. Cache-first admission must inspect and run the exact helper reference with `--pull=never` before attempting registry access. A descriptor-only inspection is insufficient; a bounded probe that checks `/bin/sh` and `git --version` (or a fixed preinstalled-tool contract) is needed.
2. Runtime `apk add git` should not be the normal path. It introduces a second network and repository failure mode and may succeed with a mutable package version unrelated to the package release. A small immutable helper image with Git already installed is preferable, provided its reference and platform acquisition are handled with bounded diagnostics.
3. Helper platform is a tool-execution concern, not target compatibility evidence. The implementation may select a runnable Linux helper for the Docker daemon or request the daemon's intended platform explicitly; it must not compare helper architecture to the selected Sub2API image as a source-compatibility gate.
4. If a pull is necessary, preserve the command's exit status and a bounded redacted stderr excerpt. Retry only classified transient network/registry failures. A cached image that cannot execute the required shell/tool should be reported as unsuitable and may be replaced only through the same bounded acquisition path.
5. The resolver still needs to run before `protect_root`, journals, backups, staging, or target writes. Helper pull/cache changes are local Docker state only and must not alter the selected container transaction.

## Files found

- `install.py:53,180-186,345-389`: fixed helper reference, Docker error mapping, source-resolution command and cleanup.
- `.trellis/spec/backend/universal-linux-delivery.md:61-66,91-94,116-117`: current universal source/build and no-write contracts.
- `adapter/Dockerfile:1-13`: later build image and source clone; confirms the existing builder also does not provide a resolver Git binary by contract.
- `tests/installer/test_preflight.py:180-198`: current resolver tests mock a successful helper pull and do not cover cache/platform/tool acquisition or pull stderr.

## External references

- Docker CLI/server behavior was observed locally with Docker 29.9.0/29.8.2 (`docker image inspect --platform`, `docker run --pull=never`, `docker buildx inspect --bootstrap`). No external web source was required.

## Caveats / Not Found

- No remote registry pull or real source clone was attempted in this check; network and credential causes for the user's failure remain unknown.
- No Git-capable helper reference was selected or added to the product. The local image cache did not contain `alpine/git` or `bitnami/git`.
- Buildx emulation support is local daemon configuration and cannot be assumed on a deployment host. The resolver must report the actual helper-platform failure rather than infer capability from Buildx metadata.
- This research does not recommend host-Git resolution; the separate source-resolver failure report records its credential and host-network boundary risks.
