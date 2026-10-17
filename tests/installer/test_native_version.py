"""Check the native replacement release boundaries."""
import contextlib
import importlib.util
import io
from pathlib import Path
import sys
import tempfile
from types import SimpleNamespace
import unittest
from unittest.mock import patch

sys.dont_write_bytecode = True
spec = importlib.util.spec_from_file_location("native_installer", Path(__file__).resolve().parents[2] / "scripts/install.py")
installer = importlib.util.module_from_spec(spec)
spec.loader.exec_module(installer)


class NativeReleaseChecks(unittest.TestCase):
    def run_install(self, installed, candidate):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            binary = root / "sub2api-codex-eval"
            header = b"\x7fELF\x02\x01\x01\x00" + b"\x00" * 8 + b"\x02\x00\x3e\x00"
            binary.write_bytes(header)
            calls = []

            def docker(*args, **kwargs):
                calls.append(args)
                if args[:3] == ("exec", "fixture", installer.BINARY):
                    return "Sub2API " + installed
                if args[:3] == ("exec", "fixture", installer.STAGED):
                    return "Sub2API " + candidate
                if args[:2] == ("cp", "fixture:" + installer.BINARY):
                    Path(args[2]).write_bytes(header)
                return ""

            args = SimpleNamespace(container="fixture", binary=binary, backup_dir=root / "backups", timeout=5)
            item = {"Id": "fixture", "Name": "/fixture", "State": {"Running": True}}
            with patch.object(installer, "select_container", return_value=item), patch.object(installer, "docker", side_effect=docker), \
                    patch.object(installer, "replace"), patch.object(installer, "readiness"), contextlib.redirect_stdout(io.StringIO()):
                installer.install(args)
            return calls

    def test_upgrade_and_reinstall(self):
        for installed in ("0.2.16", "0.2.17"):
            with self.subTest(installed=installed):
                calls = self.run_install(installed, "0.2.17")
                self.assertIn(("exec", "fixture", installer.STAGED, "-version"), calls)
                self.assertIn(("restart", "fixture"), calls)

    def test_wrong_candidate_fails_before_backup_or_restart(self):
        with self.assertRaisesRegex(RuntimeError, "Expected candidate Sub2API v0.2.17"):
            self.run_install("0.2.16", "0.2.16")


if __name__ == "__main__":
    unittest.main()
