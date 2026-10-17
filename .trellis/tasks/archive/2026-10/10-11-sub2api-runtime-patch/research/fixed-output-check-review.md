# Fixed Output Revision Check Review

## Findings (fixed)

No implementation defects were found in the stable-output revision. The changed owners align on one raw binary path, one archive path, and one release asset basename.

## Findings (not fixed)

- The complete native unit-race gate remains failed with 153 top-level failures and 51 race warnings. This output-only revision did not rebuild or rerun that gate.
- Native N150/iStoreOS behavior, live billing, deployed browser acceptance, and production deployment remain unverified.

## Independent review

- `install.py` downloads `sub2api-codex-eval.tar.gz` under the configured release tag.
- `scripts/package-sub2api.py` defaults to `dist/build/sub2api-codex-eval` and `dist/build/sub2api-codex-eval.tar.gz`. Its no-argument mode uses temporary staging and removes it after atomic archive creation.
- CI stages, extracts, uploads, and publishes `dist/build/sub2api-codex-eval.tar.gz`. README and the operations guide name the same fixed outputs. Active delivery files contain no superseded dated outer artifact names.
- The canonical raw server is `125,632,672` bytes with SHA-256 `eed01f71c39048eef1a15fbf13652b5c44d9915fff1e34c59f3d58abf1c6c9e1`.
- The canonical archive is `39,016,368` bytes with SHA-256 `aed87fa204d56075d3547c2225eb333950e201b985c8a5a05d983396635eab96`.
- Raw `tarfile` inspection found exactly eight regular files (`install.py`, `cleanup.py`, the versioned v0.2.16 amd64 binary, and five `LICENSE/` texts) plus the `LICENSE/` directory. There are no duplicate, link, AppleDouble, or PAX metadata entries. Installer, cleanup, license, and binary bytes match their live sources.
- The actual no-argument packager completed successfully and reproduced the recorded canonical size and digest. No `sub2api-package-*` temporary staging directory remains.

## Verification

- Current installer suite: 99 cases ran, 95 passed and 4 explicit Docker cases skipped.
- Packaging tests: pass, including fixed defaults, automatic temporary staging cleanup, CI asset alignment, and raw archive parity.
- Canonical archive extraction: pass through `install.py`; extracted `install.py --help` and `cleanup.py --help` passed.
- Archived installer transactions: 52 cases passed with no skips.
- Python functional stamp parsing: pass.

No native/frontend rebuild, deployment, paid request, commit, push, or publication was performed for this path-only revision.
