#!/usr/bin/env python3
"""Apply the production patches, check changes, build, and package Sub2API v0.2.17."""
import argparse
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parents[1]


def run(argv, cwd, env=None):
    print("+ " + " ".join(map(str, argv)), flush=True)
    subprocess.run(list(map(str, argv)), cwd=cwd, env=env, check=True)


def prepare_tree(source, workspace, check=False):
    tree = workspace / "source"
    shutil.copytree(source, tree, ignore=shutil.ignore_patterns("node_modules", ".git", ".entc"))
    for name in (ROOT / "patches/series").read_text().splitlines():
        if not name.strip() or name.startswith("#"):
            continue
        print(f"Apply {name}", flush=True)
        with (ROOT / "patches" / name).open("rb") as patch:
            subprocess.run(["git", "apply", "--binary", "--unidiff-zero", "--whitespace=nowarn", "-"], cwd=tree, stdin=patch, check=True)
    if check:
        shutil.copytree(ROOT / "tests/backend", tree / "backend", dirs_exist_ok=True)
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
        run([sys.executable, ROOT / "tests/installer/test_native_version.py"], ROOT, env)
        run(["node", ROOT / "tests/ui/sub2api-dropdowns.cjs", tree], ROOT, env)
        run(["go", "test", "-tags=unit", "./internal/codexeval", "./internal/service", "./internal/handler/admin", "-run", "Codex|FingerprintNormalization|FingerprintStatistics|FingerprintAttribution|ModelTraceUpstreamParity|ModelTraceParseAndThreshold|ExtractHTML|HasStandalone21"], backend, env)
        run(["go", "vet", "./internal/codexeval", "./internal/service", "./internal/handler/admin", "./internal/server", "./internal/server/routes", "./internal/web"], backend, env)
    output = ROOT / "dist/build/sub2api-codex-eval"
    output.parent.mkdir(parents=True, exist_ok=True)
    env.update(GOOS="linux", GOARCH="amd64", CGO_ENABLED="0")
    run(["go", "build", "-v", "-tags=embed", "-trimpath", "-ldflags=-s -w -X main.Version=0.2.17 -X main.BuildType=release", "-o", output, "./cmd/server"], backend, env)
    return output


def package(binary):
    output = ROOT / "dist/build/sub2api-codex-eval.tar.gz"
    if not binary.is_file():
        raise RuntimeError("Build dist/build/sub2api-codex-eval first.")
    with tarfile.open(output, "w:gz") as archive:
        for source, target in [(binary, "sub2api-codex-eval"), (ROOT / "scripts/install.py", "install.py"),
                               (ROOT / "LICENSE", "LICENSE/Project-MIT.txt")]:
            archive.add(source, arcname="sub2api-codex-eval/" + target)
        for source in sorted((ROOT / "licenses").glob("*.txt")):
            archive.add(source, arcname="sub2api-codex-eval/LICENSE/" + source.name)
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
        print(package(binary), flush=True)


if __name__ == "__main__":
    main()
