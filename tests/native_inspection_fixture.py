#!/usr/bin/env python3
"""Prepare the fixed native test container for local browser checks."""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile
import time
import urllib.request

ROOT = Path(__file__).resolve().parents[1]
CONTAINER = "sub2api-codex-eval-test"


def docker(*args):
    process = subprocess.run(["docker", *map(str, args)], capture_output=True)
    if process.returncode:
        raise RuntimeError(process.stderr.decode(errors="replace").strip())
    return process.stdout


def inspect(name):
    process = subprocess.run(["docker", "inspect", name], capture_output=True)
    return json.loads(process.stdout)[0] if process.returncode == 0 else None


def private_file(path, value):
    path.write_text(value)
    path.chmod(0o600)


def write_browser_state(work, origin, token, user):
    user_id = user.get("id", 1) if isinstance(user, dict) else 1
    state = {"cookies": [], "origins": [{"origin": origin, "localStorage": [
        {"name": "auth_token", "value": token},
        {"name": "auth_user", "value": json.dumps(user)},
        {"name": f"admin_guide_{user_id}_admin_v4_interactive", "value": "true"},
    ]}]}
    path = work / "browser-storage-state.json"
    private_file(path, json.dumps(state))
    return path


def resolve_built_binary(arch="amd64"):
    arch = "arm64" if arch in ("arm64", "aarch64") else "amd64"
    other = "amd64" if arch == "arm64" else "arm64"
    candidates = [
        ROOT / f"dist/build/linux-{arch}/sub2api-plugins-v0.2.20-linux-{arch}",
        ROOT.parent.parent / f"dist/build/linux-{arch}/sub2api-plugins-v0.2.20-linux-{arch}",
        ROOT / f"dist/build/linux-{other}/sub2api-plugins-v0.2.20-linux-{other}",
        ROOT.parent.parent / f"dist/build/linux-{other}/sub2api-plugins-v0.2.20-linux-{other}",
        ROOT / "dist/build/sub2api-plugins",
        ROOT.parent.parent / "dist/build/sub2api-plugins",
    ]
    for path in candidates:
        if path.is_file():
            return path
    for base in (ROOT, ROOT.parent.parent):
        matches = [p for p in sorted(base.glob(f"dist/build/linux-{arch}/sub2api-plugins*")) if p.is_file() and "installer" not in p.name]
        if matches:
            return matches[0]
    return ROOT / "dist/build/sub2api-plugins"


def stage_container_file(current, source, destination):
    source = Path(source)
    for mount in current.get("Mounts", []):
        if mount.get("Destination") == destination and mount.get("Type") == "bind":
            target = Path(mount["Source"])
            if target.is_dir():
                shutil.rmtree(target)
            target.parent.mkdir(parents=True, exist_ok=True)
            if target.resolve() != source.resolve():
                shutil.copy2(source, target)
            target.chmod(0o755)
            return
    docker("cp", source, CONTAINER + ":" + destination)


def seed_group_fixture(api, origin, work):
    groups = {row["name"]: row for row in api("/api/v1/admin/groups/all")}
    selected = {}
    for key, name, platform in (
            ("a", "U13 Group A", "openai"),
            ("b", "U13 Group B", "openai"),
            ("claude", "U13 Claude Group", "anthropic"),
            ("gemini", "U13 Gemini Group", "gemini"),
            ("antigravity", "U13 Antigravity Group", "antigravity")):
        selected[key] = groups.get(name) or api("/api/v1/admin/groups", "POST", {
            "name": name, "platform": platform, "rate_multiplier": 1, "subscription_type": "standard"})
    accounts = {}
    page = 1
    while True:
        response = api(f"/api/v1/admin/accounts?lite=1&page={page}&page_size=100")
        accounts.update((row["name"], row) for row in response["items"])
        if page >= response["pages"] or not response["items"]:
            break
        page += 1
    entries = {}
    specs = (
        ("a_only", "u13-a-only", "openai", "apikey", [selected["a"]["id"]],
         {"api_key": "fixture-a-only", "base_url": "http://127.0.0.1:18081"},
         {"email": "a-only@example.com"}),
        ("b_only", "u13-b-only", "openai", "apikey", [selected["b"]["id"]],
         {"api_key": "fixture-b-only", "base_url": "http://127.0.0.1:18081"},
         {"email": "b-only@example.com"}),
        ("shared", "u13-shared", "openai", "apikey", [selected["a"]["id"], selected["b"]["id"]],
         {"api_key": "fixture-shared", "base_url": "http://127.0.0.1:18081"},
         {"email": "shared@company.com"}),
        ("oauth", "u13-oauth", "openai", "oauth", [selected["b"]["id"]],
         {"access_token": "fixture-oauth-token", "id_token": "eyJ.fixture.idtoken", "plan_type": "pro", "base_url": "http://127.0.0.1:18081"},
         {"email": "oauth@company.com"}),
        ("claude_key", "u13-claude-key", "anthropic", "apikey", [selected["claude"]["id"]],
         {"api_key": "fixture-claude-key", "base_url": "http://127.0.0.1:18081"},
         {"email": "claude@company.com"}),
        ("gemini_key", "u13-gemini-key", "gemini", "apikey", [selected["gemini"]["id"]],
         {"api_key": "fixture-gemini-key", "base_url": "http://127.0.0.1:18081"},
         {"email": "gemini@company.com"}),
        ("antigravity_key", "u13-antigravity-key", "antigravity", "apikey", [selected["antigravity"]["id"]],
         {"api_key": "fixture-antigravity-key", "base_url": "http://127.0.0.1:18081"},
         {"email": "antigravity@company.com"}),
    )
    for key, name, platform, auth_type, membership, creds, extra in specs:
        body = {"name": name, "platform": platform, "type": auth_type, "group_ids": membership, "priority": 10,
                "status": "active", "error_message": "", "schedulable": True,
                "credentials": creds, "extra": extra,
                "upstream_billing_probe_enabled": False}
        if name in accounts:
            account = api(f"/api/v1/admin/accounts/{accounts[name]['id']}", "PUT", body)
        else:
            account = api("/api/v1/admin/accounts", "POST", {**body, "concurrency": 1})
        if account.get("schedulable") is False:
            account = api(f"/api/v1/admin/accounts/{account['id']}/schedulable", "POST", {"schedulable": True})
        entries[key] = {"id": account["id"], "name": name, "platform": platform, "type": auth_type, "group_ids": membership}
    manifest = {"origin": origin, "groups": {
        key: {"id": group["id"], "name": group["name"]} for key, group in selected.items()}, "accounts": entries}
    path = work / "u13-groups.json"
    private_file(path, json.dumps(manifest, indent=2))
    return path


def main():
    current = inspect(CONTAINER)
    if current is None:
        raise RuntimeError("The retained test container is absent. This tool does not create another container.")
    environment = dict(item.split("=", 1) for item in current["Config"]["Env"] if "=" in item)
    bindings = current["HostConfig"]["PortBindings"]["8080/tcp"]
    origin = "http://127.0.0.1:" + bindings[0]["HostPort"]
    work = Path(tempfile.mkdtemp(prefix="sub2api-inspection-fixture-", dir="/tmp"))
    for mount in current["Mounts"]:
        if mount["Destination"] == "/app/data" and mount["Type"] == "bind":
            Path(mount["Source"]).mkdir(mode=0o700, parents=True, exist_ok=True)
    if current["State"]["Running"]:
        docker("stop", "--time", "30", CONTAINER)
    image_arch = docker("image", "inspect", current["Image"], "--format", "{{.Architecture}}").decode().strip()
    stage_container_file(current, resolve_built_binary(image_arch), "/app/sub2api-bin")
    stage_container_file(current, ROOT / "tests/native_fixture_start.sh", "/app/sub2api")
    stage_container_file(current, ROOT / "tests/native_fixture_supplier.py", "/app/fixture-supplier.py")
    docker("start", CONTAINER)
    deadline = time.monotonic() + 240
    while time.monotonic() < deadline:
        process = subprocess.run(["docker", "exec", CONTAINER, "psql", "-h", "127.0.0.1", "-U", "sub2api", "-d", "sub2api", "-tAc", "SELECT 1"], capture_output=True)
        if process.returncode == 0:
            break
        time.sleep(1)
    else:
        raise RuntimeError("PostgreSQL startup timed out. Inspect the retained container logs.")
    docker("exec", CONTAINER, "touch", "/app/fixture.ready")
    token = ""

    def api(path, method="GET", body=None, auth=True):
        headers = {"Content-Type": "application/json"}
        if auth:
            headers["Authorization"] = "Bearer " + token
        request = urllib.request.Request(origin + path, method=method, headers=headers,
                                         data=None if body is None else json.dumps(body).encode())
        with urllib.request.urlopen(request, timeout=35) as response:
            payload = json.load(response)
        return payload.get("data", payload)

    deadline = time.monotonic() + 120
    while time.monotonic() < deadline:
        try:
            if api("/health", auth=False).get("status") == "ok" and not api("/setup/status", auth=False).get("needs_setup"):
                break
        except (OSError, ValueError):
            pass
        time.sleep(1)
    else:
        raise RuntimeError("Application startup timed out. Inspect the retained container logs.")
    email, password = environment["ADMIN_EMAIL"], environment["ADMIN_PASSWORD"]
    token = api("/api/v1/auth/login", "POST", {"email": email, "password": password}, auth=False)["access_token"]
    private_file(work / "browser-login.json", json.dumps({"origin": origin, "email": email, "password": password, "token": token, "container": CONTAINER}))
    storage_state = write_browser_state(work, origin, token, api("/api/v1/auth/me"))
    compliance = api("/api/v1/admin/compliance")
    if compliance["required"]:
        api("/api/v1/admin/compliance/accept", "POST", {"phrase": compliance["ack_phrase_en"], "language": "en"})
    accounts = api("/api/v1/admin/accounts?page_size=100")
    rows = accounts.get("items", []) if isinstance(accounts, dict) else accounts
    for account in rows:
        if account["name"].startswith("inspection-"):
            credentials = dict(account.get("credentials", {}))
            credentials["base_url"] = "http://127.0.0.1:18081"
            api(f"/api/v1/admin/accounts/{account['id']}", "PUT", {"credentials": credentials})
    group_manifest = seed_group_fixture(api, origin, work)
    group_id = json.loads(group_manifest.read_text())["groups"]["a"]["id"]
    settings_path = "/api/v1/admin/settings/codex-eval-inspection"
    settings = api(settings_path)
    settings.update(enabled=False, group_id=group_id, account_type="all", account_plan="all", enabled_only=False,
                    model="inspection-literal", adjust_priority=False,
                    ntfy_server="http://127.0.0.1:18081/ntfy",
                    bark_url="http://127.0.0.1:18081/bark/fixture-private-key")
    api(settings_path, "PUT", settings)
    version = docker("exec", CONTAINER, "/app/sub2api", "-version").decode().strip()
    print("Reusable fixture ready: " + origin)
    print("Browser credentials: " + str(work / "browser-login.json"))
    print("Browser storage state: " + str(storage_state))
    print("Group fixture manifest: " + str(group_manifest))
    print("Runtime version: " + version)


if __name__ == "__main__":
    main()
