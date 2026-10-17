# U20 CPA Cleanup Implementation

Date: 2026-10-18

## Scope and dependency findings

The user requested removal of old CPA contents.
U20 replaces the earlier requirement to retain those sources.
Keep the completed native integration, current commands, fixtures, legal attribution, private files, and reusable caches.
Do not create an archive or rebuild the unchanged binary.

`scripts/build-sub2api.py` applies `native/sub2api.patch` and copies `native/backend/` and `native/frontend/` into an existing upstream tree.
It reads the upstream Go module and frontend dependencies. It does not read `internal/plugin/`.
The release workflow selects `upstream/backend/go.mod`.
Current fixture and UI files read `native/`, the root installer, and the fixed native binary.
No current command or test imports the removed root module or old source graph.

The root Go module declared `cpa-codex-candy-eval` and served only the old shared-library build.
Its `go.sum` contained that module's dependency records.
The placeholder `.env.example` contained only CPA deployment settings.
The active installer does not read dotenv files.

The screenshot `docs/images/example.png` shows the old CPAMC credential page.
No current guide or command uses it.

## Actual removals

The following ten targets moved to macOS Trash through `/usr/bin/trash`.
These targets contain 54 tracked files in total.
Directory moves include all descendants and any Finder metadata within those directories.

| Removed target | Contents and reason |
| --- | --- |
| `cmd/` | Two files under `cmd/cpa-codex-candy-eval/`. Only the old C ABI entry point and Finder metadata remained. |
| `internal/` | 44 files under `internal/plugin/`. Old handlers, tests, runners, data, page assets, and renderer copies. |
| `examples/` | Old `inspection.py` using CPA management and plugin routes. |
| `docs/inspection.md` | Old external CPA inspection instructions. |
| `docs/images/` | Old `example.png` screenshot. |
| `scripts/deploy.sh` | Old shared-library build, SSH upload, and CPA restart path. |
| `scripts/test-ui.cjs` | Old UI harness reading `internal/plugin/web/`. |
| `go.mod` | Old root plugin module. |
| `go.sum` | Old root plugin dependency records. |
| `.env.example` | Old CPA deployment placeholder. |

`install.sh` and `install.ps1` were already absent before U20.
Read-only `git show HEAD:<path>` confirmed both are exclusively old CPA shared-library installers.
They install `.so`, `.dylib`, or `.dll` files under the old host's `plugins/<os>/<arch>` paths.
Their existing deletions belong in the U20 cleanup commit plan.
This turn did not restore, stage, or remove them again.
The complete cleanup plan therefore contains 56 tracked deletions, including those two prior deletions.

The `dist/` scan found no CPA-named artifact or shared-library candidate outside reusable caches.
No additional `dist/` path moved to Trash.
Current native outputs and reusable dependency/build directories remain in place.
Private `.env`, machine configuration, sibling repositories, and unrelated output remain untouched.

## Current document changes

- `AGENTS.md` now describes only current Sub2API owners, commands, installation, and focused checks.
- Its managed Trellis block is unchanged. General private-data and commit rules remain.
- `directory-structure.md` now maps native tasks, persistence, inspection, transport, page assets, and adjacent tests.
- `database-guidelines.md` retains native settings and JSON history contracts and removes the old CPA storage section.
- `error-handling.md` and `logging-guidelines.md` now cite native owners and their actual response/error channels.
- `index.md` and `quality-guidelines.md` remove original-plugin source requirements and retained CPA rules.
- The hybrid commit section remains unchanged.
- `runtime-installation.md` describes native theme variables without an old host reference.
- README development instructions now describe only `native/` and the upstream Go module.
- `.gitignore` no longer has an exception for the removed `.env.example`.

Current documentation links do not point to removed sources.
Historical task and workspace evidence remains unchanged by this cleanup.

## License and package decision

Keep root `LICENSE` and all required third-party license texts.
Keep native footer source URLs, copyright attribution, embedded renderer notices, and ModelTrace attribution.
Do not rename the repository, Git remote, or local checkout.

The package script now calls the root MIT license `LICENSE/Project-MIT.txt` instead of `LICENSE/CPA-MIT.txt`.
Both third-party notice copies use that filename.
They also identify the actual v0.2.17 upstream release and source URL.
The native delivery spec and README record the neutral filename.
License text and copyright holders remain unchanged.

Ran only the package command:

```sh
python3 scripts/package-sub2api.py
```

It reused the existing binary and regenerated `dist/build/sub2api-codex-eval.tar.gz`.
The binary remains 125599904 bytes with its prior modification time.
The archive is 38700696 bytes.
No binary build, Docker operation, deployment, or hash calculation ran.

The archive contains exactly these seven files:

```text
sub2api-codex-eval/sub2api-codex-eval
sub2api-codex-eval/install.py
sub2api-codex-eval/LICENSE/Project-MIT.txt
sub2api-codex-eval/LICENSE/ModelTrace-MIT.txt
sub2api-codex-eval/LICENSE/Renderers-MIT.txt
sub2api-codex-eval/LICENSE/Sub2API-LGPL-3.0.txt
sub2api-codex-eval/LICENSE/THIRD_PARTY_NOTICES.txt
```

## Checks

| Check | Actual result |
| --- | --- |
| Old-source imports and direct reads before removal | No active command, workflow, fixture, or UI harness reads the old graph. |
| Removed paths | All ten selected targets are absent. The two older installer deletions remain absent. |
| Hidden tracked-project reference scan | Checked 230 existing text files. Excluded private/local settings and historical task/workspace evidence. |
| Remaining old host names | Only unchanged native footer attribution and five nonexecuting CSS/JavaScript comments remain. They are not dependencies. |
| Current old-source references | No removed source imports, reads, commands, package-license names, or current document pointers remain. |
| Current local Markdown targets | Passed for 16 documents, including backend specs and shared guides. |
| Active Python AST syntax | Passed for six installer, build, package, and fixture files. |
| Shell syntax | Passed for the build wrapper, native fixture start script, and current test entry. |
| Node syntax | `node --check tests/ui/sub2api-dropdowns.cjs` passed. |
| Archive inventory | Exact seven-file set. No development, old CPA, or extra payload files. |
| Archived installer and license contents | Match their current source files byte for byte. |
| Notice copies | Root Markdown and packaged text copies match byte for byte. |
| Hybrid commit section and managed AGENTS block | Match their prior staged text. |
| Current native/runtime/test inputs | No added difference against the existing index. |
| `git diff --check` | Passed |

No runtime behavior or current native asset changed.
Existing product tests and U18 binary/browser evidence remain applicable.
No new suite, race gate, provider call, or repeated product build was necessary.

## Remaining tree and handoff

Active product paths are `native/`, root `install.py`, four current scripts, `tests/`, current docs, and required license files.
The repository keeps its release workflow, Trellis files, ignored outputs, and machine configuration.
It has no root Go module, old shared-library entry point, external inspection script, or old CPA screenshot tree.

The main session owns task metadata, independent review, commit-plan synchronization, and final approval.
This turn did not stage or commit any file.
The earlier history cancellation report remains unchanged and unresolved by this source cleanup.
