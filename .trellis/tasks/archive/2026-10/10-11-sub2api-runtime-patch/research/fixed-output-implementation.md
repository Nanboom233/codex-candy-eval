# Fixed Output Implementation

## Canonical files

- Archive: `dist/build/sub2api-codex-eval.tar.gz`, 39,016,368 bytes.
  SHA256: `aed87fa204d56075d3547c2225eb333950e201b985c8a5a05d983396635eab96`.
- Raw patched server: `dist/build/sub2api-codex-eval`, 125,632,672 bytes.
  SHA256: `eed01f71c39048eef1a15fbf13652b5c44d9915fff1e34c59f3d58abf1c6c9e1`.

The raw server is copied from the previously reviewed v0.2.16 amd64 output.
No native or frontend compilation ran. The compiled functional stamp remains unchanged.

## Changed owners

install.py downloads the stable asset basename `sub2api-codex-eval.tar.gz` under the existing release tag.
scripts/package-sub2api.py defaults to the fixed raw binary and archive paths.
Running `python3 -B scripts/package-sub2api.py` packages and writes the canonical archive through automatic system-temp staging.
Explicit fresh staging and --binary/--archive overrides remain available for development/test callers.
The same explicit Python archive writer validates raw members before atomic replacement.

Release CI stages its server at the fixed raw path, runs the default packager, extracts the fixed archive for its disposable lifecycle, and uploads/publishes the stable basename.
README and operations instructions identify the fixed output files.
Affected transaction/packaging tests verify the URL, CI asset alignment, fixed defaults and automatic staging cleanup.

The archive's internal payload remains eight files: install.py, cleanup.py, the versioned v0.2.16 amd64 server, and five LICENSE texts.
Only LICENSE is an archive directory. There are no AppleDouble entries, links, extended metadata, source or test files.
All installer/cleanup/license bytes match current sources. Archived server bytes match the canonical raw server.

## Executed checks

- Current Python suite: 99 cases total, 95 passed and four explicit Docker cases skipped.
- Extracted canonical archive transaction suite: 52 cases passed, no skips.
- Raw tar-member allowlist, duplicates/types/metadata, source/license/binary parity and installer extraction passed.
- Extracted install.py --help and cleanup.py --help passed.
- Installer functional stamp validation passed.
- Active owner/reference search found no prior dated output directory or versioned outer archive name in changed active files.
- Automatic packaging/test temporary staging is absent after completion.

Validation identities and raw payload paths are recorded in research/fixed-output-validation.json.
Raw check logs remain at /tmp/fixed-output-current-suite.log and /tmp/fixed-output-archive-transactions.log.

## Preservation and limits

The previous prebuilt archive remains unchanged at SHA256 `24f463eb70cb20a4127984bfbad84aeaf9eeab1ba5d06db7e6927d643d08cc01`.
Historical universal and minimal archives retain their previously recorded hashes.
No historical report or unrelated work was removed.
No broad Go, native, frontend or Docker lifecycle checks were repeated for this output-name change.
The earlier complete native unit-race failure remains unresolved: 153 top-level failures and 51 race warnings.
Native N150/iStoreOS, real billing and deployed browser acceptance remain unverified.
No production deployment, paid request, commit, push or publication occurred.
