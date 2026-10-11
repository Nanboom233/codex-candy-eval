import ast
import hashlib
import json
import os
import re
import subprocess
import tomllib
from pathlib import Path


def git(*args):
    return subprocess.check_output(["git", *args], text=True)


paths = git("diff", "--cached", "--name-only", "--diff-filter=ACM").splitlines()
assert paths, "No installation files are staged"
manifest = Path(__file__).with_name("commit-files.txt").read_text().splitlines()
assert paths == manifest, "The staged file list differs from the approved manifest"
hashes = json.loads(Path(".trellis/.template-hashes.json").read_text())["hashes"]
counts = {"python": 0, "json": 0, "toml": 0, "template_match": 0}
changed, unregistered = [], []
tracked = set(git("ls-files").splitlines())
private = {".devin/config.local.json", ".devin/hooks.v1.json"}
secret = re.compile(r"(?m)^-----BEGIN (?:RSA |EC |OPENSSH |PGP )?PRIVATE KEY|\b(?:sk-(?:proj-)?[A-Za-z0-9_-]{25,}|gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{40,})\b")
for name in paths:
    assert name not in private, f"Machine-only configuration staged: {name}"
    assert not name.startswith(".trellis/tasks/"), f"Task bookkeeping staged: {name}"
    assert not name.startswith(".trellis/workspace/") or name == ".trellis/workspace/index.md", f"Session bookkeeping staged: {name}"
    assert "/__pycache__/" not in name and not name.endswith(".pyc"), f"Cache staged: {name}"
    data = subprocess.check_output(["git", "show", f":{name}"])
    text = data.decode("utf-8")
    assert not secret.search(text), f"Possible private key or token in {name}"
    if name.endswith(".py"):
        ast.parse(text, filename=name)
        counts["python"] += 1
    elif name.endswith(".json"):
        json.loads(text)
        counts["json"] += 1
    elif name.endswith(".toml"):
        config = tomllib.loads(text)
        counts["toml"] += 1
        if name == ".codex/config.toml":
            assert "mcp_servers" not in config, "Machine-only MCP settings staged"
            assert config["agents"]["max_depth"] == 1, "Codex recursion guard missing"
    if name in hashes:
        if hashlib.sha256(data).hexdigest() == hashes[name]:
            counts["template_match"] += 1
        else:
            changed.append(name)
    else:
        unregistered.append(name)
    if name in {".trellis/spec/backend/index.md", ".trellis/spec/backend/quality-guidelines.md"}:
        for link in re.findall(r"\]\(([^)]+)\)", text):
            target = os.path.normpath(str(Path(name).parent / link.split("#", 1)[0]))
            assert target in tracked, f"Untracked spec link: {name} -> {target}"
print("Staged installation checks passed:", len(paths), "files", counts)
print("Template differences for review:", changed)
print("Files without template hashes:", unregistered)
