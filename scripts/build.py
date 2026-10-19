#!/usr/bin/env python3
"""Apply the production patches, check changes, build, and package Sub2API v0.2.17."""
import argparse
import gzip
import os
from pathlib import Path
import shutil
import subprocess
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
        run(["node", ROOT / "tests/ui/sub2api-dropdowns.cjs", tree], ROOT, env)
        run(["go", "test", "-tags=unit", "./internal/codexeval", "./internal/service", "./internal/handler/admin", "-run", "Codex|FingerprintNormalization|FingerprintStatistics|FingerprintAttribution|ModelTraceUpstreamParity|ModelTraceParseAndThreshold|ExtractHTML|HasStandalone21"], backend, env)
        run(["go", "vet", "./internal/codexeval", "./internal/service", "./internal/handler/admin", "./internal/server", "./internal/server/routes", "./internal/web"], backend, env)
    output = ROOT / "dist/build/sub2api-codex-eval"
    output.parent.mkdir(parents=True, exist_ok=True)
    env.update(GOOS="linux", GOARCH="amd64", CGO_ENABLED="0")
    run(["go", "build", "-v", "-tags=embed", "-trimpath", "-ldflags=-s -w -X main.Version=0.2.17 -X main.BuildType=release", "-o", output, "./cmd/server"], backend, env)
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
    license_dir = staging / "LICENSE"
    license_dir.mkdir()
    shutil.copyfile(ROOT / "LICENSE", license_dir / "Project-MIT.txt")
    for text in sorted((ROOT / "licenses").glob("*.txt")):
        shutil.copyfile(text, license_dir / text.name)
    return staging


def package(binary, workspace, source, check=False):
    output = ROOT / "dist/build/sub2api-codex-eval-installer"
    if not binary.is_file():
        raise RuntimeError("Build dist/build/sub2api-codex-eval first.")
    staging = prepare_installer(binary, workspace, source)
    env = os.environ.copy()
    env["GOFLAGS"] = "-mod=readonly"
    if check:
        for test in sorted((ROOT / "tests/installer").glob("*_test.go")):
            shutil.copyfile(test, staging / test.name)
        host_env = env.copy()
        for name in ("GOOS", "GOARCH"):
            host_env.pop(name, None)
        run(["go", "test", "."], staging, host_env)
    env.update(GOOS="linux", GOARCH="amd64", CGO_ENABLED="0")
    run(["go", "build", "-trimpath", "-ldflags=-s -w", "-o", output, "."], staging, env)
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
