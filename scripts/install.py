#!/usr/bin/env python3
"""Replace a standard v0.2.17 Linux/amd64 Sub2API container binary."""
import argparse
import datetime
import http.client
import json
from pathlib import Path
import re
import subprocess
import sys
import tempfile
import time

BINARY = "/app/sub2api"
STAGED = "/app/sub2api.codex-eval-new"


def docker(*args, quiet=False):
    result = subprocess.run(["docker", *args], stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE, text=True)
    if not quiet:
        for output, stream in ((result.stdout, sys.stdout), (result.stderr, sys.stderr)):
            if output:
                print(output, end="", file=stream, flush=True)
    if result.returncode:
        message = result.stderr.strip() or result.stdout.strip()
        raise RuntimeError(f"docker {' '.join(args[:2])} failed ({result.returncode}): {message}")
    return result.stdout


def inspect(container):
    return json.loads(docker("inspect", container, quiet=True))[0]


def select_container(name):
    if name:
        return inspect(name)
    ids = docker("ps", "-q", quiet=True).split()
    containers = [inspect(cid) for cid in ids]
    if not containers:
        raise RuntimeError("No running containers.")
    for i, item in enumerate(containers, 1):
        print(f"[{i}] {item['Name'].lstrip('/')} ({item['Id'][:12]})")
    while True:
        value = input("Select a container number. No default. Enter q to cancel: ").strip()
        if value.lower() == "q":
            raise RuntimeError("Cancelled.")
        if value.isdigit() and 1 <= int(value) <= len(containers):
            return containers[int(value) - 1]


def endpoint(item):
    env = dict(entry.split("=", 1) for entry in item["Config"].get("Env", []) if "=" in entry)
    port = int(env.get("SERVER_PORT", "8080"))
    bindings = item["NetworkSettings"].get("Ports", {}).get(f"{port}/tcp") or []
    if bindings:
        binding = bindings[0]
        host = binding["HostIp"]
        return ("127.0.0.1" if host in ("", "0.0.0.0", "::") else host, int(binding["HostPort"]))
    if item["HostConfig"].get("NetworkMode") == "host":
        return "127.0.0.1", port
    for network in item["NetworkSettings"].get("Networks", {}).values():
        if network.get("IPAddress"):
            return network["IPAddress"], port
    raise RuntimeError("No reachable application address in Docker inspect.")


def readiness(cid, timeout):
    started = time.monotonic()
    deadline = started + timeout
    last_error = "not ready"
    next_progress = 0.0
    previous_state = None
    while time.monotonic() < deadline:
        item = inspect(cid)
        state = item["State"]
        status = state.get("Status", "unknown")
        elapsed = time.monotonic() - started
        if not state.get("Running") and not state.get("Restarting"):
            detail = state.get("Error") or "no runtime error reported"
            raise RuntimeError(f"Container {status} after {elapsed:.1f}s: exit={state.get('ExitCode')}, OOMKilled={state.get('OOMKilled', False)}, error={detail}; last request: {last_error}")
        if state.get("Paused"):
            raise RuntimeError(f"Container paused after {elapsed:.1f}s; last request: {last_error}")
        address = "no address yet"
        try:
            if state.get("Restarting"):
                raise RuntimeError("Container is restarting")
            host, port = endpoint(item)
            address = f"{host}:{port}"
            for path in ("/health", "/setup/status"):
                remaining = deadline - time.monotonic()
                if remaining <= 0:
                    raise RuntimeError("Readiness deadline reached")
                connection = http.client.HTTPConnection(host, port, timeout=min(3, remaining))
                try:
                    connection.request("GET", path)
                    response = connection.getresponse()
                    data = json.loads(response.read())
                    if response.status != 200:
                        raise RuntimeError(f"{path}: HTTP {response.status}")
                    if path == "/health" and data.get("status") != "ok":
                        raise RuntimeError(f"{path}: {data}")
                    if path == "/setup/status" and data.get("data", {}).get("needs_setup") is not False:
                        raise RuntimeError("Application setup is incomplete.")
                except (OSError, ValueError, http.client.HTTPException, RuntimeError) as error:
                    raise RuntimeError(f"{path} at {address}: {error}") from error
                finally:
                    connection.close()
            print(f"Application ready after {time.monotonic() - started:.1f}s: health OK, setup complete.", flush=True)
            return
        except (OSError, ValueError, http.client.HTTPException, RuntimeError) as error:
            last_error = str(error)
        current_state = (status, item.get("RestartCount", 0), address)
        if time.monotonic() >= next_progress or current_state != previous_state:
            elapsed = time.monotonic() - started
            print(f"Application is starting; waiting for health and setup [{elapsed:.1f}/{timeout}s, state={status}, restarts={current_state[1]}, address={address}].", flush=True)
            next_progress = time.monotonic() + 5
            previous_state = current_state
        time.sleep(min(1, max(0, deadline - time.monotonic())))
    raise RuntimeError(f"Readiness timed out after {time.monotonic() - started:.1f}s: state={status}, restarts={item.get('RestartCount', 0)}, last request: {last_error}")


def replace(cid, source):
    if not inspect(cid)["State"]["Running"]:
        # A bad replacement may have stopped the container. Docker cp works
        # on stopped containers, so restoration needs no executable inside it.
        docker("cp", str(source), f"{cid}:{BINARY}")
        return
    docker("cp", str(source), f"{cid}:{STAGED}")
    docker("exec", "-u", "0", cid, "sh", "-c",
           f"chmod 755 {STAGED} && mv -f {STAGED} {BINARY}")


def install(args):
    item = select_container(args.container)
    cid = item["Id"]
    print(f"Container: {item['Name'].lstrip('/')} ({cid[:12]})", flush=True)
    if not item["State"]["Running"]:
        raise RuntimeError("Select a running container.")
    print("[1/4] Check v0.2.17 and the local binary.", flush=True)
    version = docker("exec", cid, BINARY, "-version", quiet=True)
    if not any(re.search(r"\b" + re.escape(release) + r"\b", version) for release in ("0.2.16", "0.2.17")):
        raise RuntimeError(f"Expected Sub2API v0.2.16 or v0.2.17, received: {version.strip()}")
    source = args.binary.resolve()
    if not source.is_file():
        raise RuntimeError(f"Binary not found: {source}")
    with source.open("rb") as stream:
        header = stream.read(20)
    if header != b"\x7fELF\x02\x01\x01\x00" + b"\x00" * 8 + b"\x02\x00\x3e\x00":
        raise RuntimeError("Expected a static Linux/amd64 ELF executable.")
    docker("cp", str(source), f"{cid}:{STAGED}")
    try:
        docker("exec", "-u", "0", cid, "chmod", "755", STAGED)
        candidate_version = docker("exec", cid, STAGED, "-version", quiet=True)
        if not re.search(r"\b0\.2\.17\b", candidate_version):
            raise RuntimeError(f"Expected candidate Sub2API v0.2.17, received: {candidate_version.strip()}")
    finally:
        docker("exec", "-u", "0", cid, "rm", "-f", STAGED)
    print("[2/4] Back up the current executable.", flush=True)
    if args.backup_dir is None:
        backup_dir = Path(tempfile.mkdtemp(prefix=f"sub2api-codex-eval-{cid[:12]}-", dir="/tmp"))
    else:
        stamp = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%S%fZ")
        backup_dir = args.backup_dir.resolve() / f"{cid[:12]}-{stamp}"
        backup_dir.mkdir(parents=True)
    backup = backup_dir / "sub2api"
    docker("cp", f"{cid}:{BINARY}", str(backup))
    with backup.open("rb") as stream:
        original_header = stream.read(20)
    if original_header != header:
        raise RuntimeError("The existing executable must also be Linux/amd64.")
    print(f"Backup: {backup}", flush=True)
    changed = False
    try:
        print("[3/4] Replace the binary and restart the selected container.", flush=True)
        # The original backup is complete before any container mutation.
        changed = True
        replace(cid, source)
        docker("restart", cid)
        print("[4/4] Wait for native health and setup status.", flush=True)
        readiness(cid, args.timeout)
    except (Exception, KeyboardInterrupt):
        if changed:
            print("Installation failed; restoring the original executable.", file=sys.stderr, flush=True)
            try:
                replace(cid, backup)
                docker("restart", cid)
                readiness(cid, args.timeout)
                print("Original executable restored.", file=sys.stderr, flush=True)
            except Exception as rollback_error:
                raise RuntimeError(f"Restoration failed: {rollback_error}. Backup: {backup}") from rollback_error
        raise
    print("Installed. Open 账户评估 in the native administrator sidebar.", flush=True)


def main():
    here = Path(__file__).resolve().parent
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("operation", nargs="?", choices=["install"], default="install")
    parser.add_argument("--container")
    parser.add_argument("--binary", type=Path, default=here / "sub2api-codex-eval")
    parser.add_argument("--backup-dir", type=Path, help="Backup directory (default: a private temporary directory under /tmp)")
    parser.add_argument("--timeout", type=int, default=120)
    args = parser.parse_args()
    if args.timeout <= 0:
        parser.error("--timeout must be positive")
    try:
        install(args)
    except (RuntimeError, OSError, EOFError, KeyboardInterrupt) as error:
        print(f"Error: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
