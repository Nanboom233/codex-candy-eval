#!/usr/bin/env python3
"""Build the standalone evaluation plugin for Sub2API v0.2.19."""
import argparse
import gzip
import os
from pathlib import Path
import runpy
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]
VERSION = "0.2.19"


def run(argv, cwd, env=None):
    print("+ " + " ".join(map(str, argv)), flush=True)
    subprocess.run(list(map(str, argv)), cwd=cwd, env=env, check=True)


def prepare_tree(source, workspace, check=False):
    tree = workspace / "source"
    shutil.copytree(source, tree, ignore=shutil.ignore_patterns("node_modules", ".git", ".entc"))
    for path in sorted((ROOT / "patches").glob("*.patch")):
        print(f"Apply {path.name}", flush=True)
        run(["git", "apply", "--binary", path], tree)
    return tree


def build(source, workspace, check=False):
    tree = prepare_tree(source, workspace, check)
    backend = tree / "backend"
    # Reuse the existing lockfile installation, then build the patched native
    # sidebar/view. An old embedded frontend cannot deliver the new route.
    dependencies = source / "frontend/node_modules"
    if not dependencies.is_dir():
        raise RuntimeError("Install the original frontend dependencies with pnpm install --frozen-lockfile first.")
    (tree / "frontend/node_modules").symlink_to(dependencies, target_is_directory=True)
    env = os.environ.copy()
    env.setdefault("COREPACK_ENABLE_DOWNLOAD_PROMPT", "0")
    env.setdefault("COREPACK_ENABLE_AUTO_PIN", "0")
    env.setdefault("COREPACK_HOME", str(ROOT / "dist/corepack"))
    env["PYTHONDONTWRITEBYTECODE"] = "1"
    if check:
        run(["corepack", "prepare", "pnpm@9.15.9", "--activate"], ROOT, env)
    run(["pnpm", "build"], tree / "frontend", env)
    env["GOFLAGS"] = "-mod=readonly"
    if check:
        host_env = dict(env, SUB2API_SOURCE=str(tree))
        for name in ("GOOS", "GOARCH", "CGO_ENABLED"):
            host_env.pop(name, None)
        support = runpy.run_path(str(ROOT / "tests/support.py"))
        support["test_entries"](ROOT, tree, host_env, run)
        run(["go", "vet", "./internal/server", "./internal/server/routes", "./internal/web"], backend, host_env)
    output = ROOT / "dist/build/sub2api-codex-eval"
    output.parent.mkdir(parents=True, exist_ok=True)
    env.update(GOOS="linux", GOARCH="amd64", CGO_ENABLED="0")
    run(["go", "build", "-v", "-tags=embed", "-trimpath", f"-ldflags=-s -w -X main.Version={VERSION} -X main.BuildType=release", "-o", output, "./cmd/server"], backend, env)
    return output


def prepare_installer(binary, workspace, source):
    staging = workspace / "installer"
    shutil.copytree(ROOT / "scripts/installer", staging)
    module_lines = (source / "backend/go.mod").read_text().splitlines()
    go_version = next(line for line in module_lines if line.startswith("go "))
    toolchain = [line for line in module_lines if line.startswith("toolchain ")]
    (staging / "go.mod").write_text("module sub2api-codex-eval-installer\n\n" + go_version + "\n" + "\n".join(toolchain) + "\n")
    with binary.open("rb") as original, gzip.open(staging / "payload.gz", "wb") as compressed:
        shutil.copyfileobj(original, compressed)
    license_dir = staging / "licenses/codex-candy-eval"
    license_dir.mkdir(parents=True)
    shutil.copyfile(ROOT / "LICENSE", license_dir / "LICENSE")
    shutil.copytree(ROOT / "licenses", license_dir, dirs_exist_ok=True)
    return staging


def package(binary, workspace, source, check=False):
    output = ROOT / "dist/build/sub2api-codex-eval-installer"
    if not binary.is_file():
        raise RuntimeError("Build dist/build/sub2api-codex-eval first.")
    staging = prepare_installer(binary, workspace, source)
    env = os.environ.copy()
    env["GOFLAGS"] = "-mod=readonly"
    env.update(GOOS="linux", GOARCH="amd64", CGO_ENABLED="0")
    run(["go", "build", "-trimpath", f"-ldflags=-s -w -X main.buildVersion={VERSION}", "-o", output, "."], staging, env)
    return output


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source", required=True, type=Path)
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    with tempfile.TemporaryDirectory(prefix="sub2api-codex-eval-") as temporary:
        print(f"Build workspace: {temporary}", flush=True)
        binary = build(args.source.resolve(), Path(temporary), args.check)
        print(binary, flush=True)
        print(package(binary, Path(temporary), args.source.resolve(), args.check), flush=True)


if __name__ == "__main__":
    main()
