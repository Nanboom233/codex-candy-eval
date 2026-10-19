# Execution Plan

## Planning Gate

- [x] Create the requested task and inspect the supplied comparison.
- [x] Inspect native runtime, persistence, page, tests, and delivery contracts.
- [x] Write the PRD, design, research notes, and context manifests.
- [x] Present the final planning summary and obtain the required subsequent approval.
- [x] Reread all task files and remove redundant checks and added requirements.
- [x] Run `python3 .trellis/scripts/task.py start .trellis/tasks/10-19-upstream-compare-adaptation` after that approval.

## Implementation

Dispatch one `trellis-implement` owner for the coupled runtime, history, page, and tests.
Start its prompt with the active task path. Use native context injection and the child-side loading fallback.
The main session owns coordination, browser acceptance, specification updates, and final delivery.

- [x] Add the shared variant definitions, request normalization, result version, and markup extraction in patch `0001`.
- [x] Add the schema-2 native migration without changing schema-1 fingerprint meaning or unversioned retirement behavior.
- [x] Save original history bytes before migration and block replacement and later saves if the copy fails.
- [x] Add page configuration, version selection, prompt updates, preferences, and record metadata in patch `0006`.
- [x] Show each saved version in latest-result rows as well as history and detail metadata.
- [x] Preserve saved model and effort when the version changes before asynchronous model restoration finishes.
- [x] Adapt picture labels and toolbar layout while preserving native controls and sandbox behavior.
- [x] Extend existing native Pelican, runtime/persistence, and UI checks only for gaps listed in `research/task-audit.md`.
- [x] Keep new test names within the existing focused selector, or update that selector with the tests.
- [x] Update the Chinese README Pelican instructions and existing history recovery guidance.

## Validation

Use the retained original source at `/tmp/sub2api-codex-eval-v0217-check` if its original files and dependencies remain usable.
This retained source exists and has frontend dependencies. It has no `.git` directory.
Do not depend on `git describe` to identify this source tree.
The last task's effective patched tree is `/tmp/upstream-sync-native/source` and is reference evidence only.

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check
git diff --check
```

Use the verification ownership map in `research/task-audit.md`.
Run the standard build/check entry once for the final candidate.
Rerun only a failed or affected check after a repair. Do not add another baseline or full-suite run.
The successful standard build supplies the server and installer evidence. Do not add a separate artifact test.
Use actual browser inspection for visual and animation behavior that the DOM harness cannot establish.

## Review and Finish

- [x] Dispatch `trellis-check` against the final candidate and valid verification evidence.
- [x] Have the reviewer reuse passing evidence and check scope without repeating the standard build.
- [x] Return in-scope repairs to the implementation owner and resolve acceptance failures.
- [x] Update runtime, persistence, and error-handling specs with the effective contracts and source anchors.
- [x] Inspect the complete diff against every acceptance criterion.
- [x] Follow the existing work-commit approval and signing rules if a work commit is requested.
- [x] Finish and archive through the native Trellis workflow after the required work commit.

## Recovery

Keep changes limited to the named patches, tests, build selector if needed, documentation, and task artifacts.
Development checks use temporary data. No production installation or real history upgrade runs in this task.
If the history marker advances in a later deployment, the old executable needs a compatible pre-upgrade history file.
The approved migration saves an original-file copy before replacement. Keep its path available in the existing recovery guidance.
An executable downgrade also requires its matching old-format history snapshot. Restoring that snapshot excludes later active-history records.

## Evidence

The standard build/check entry passed. Its log is `/tmp/pelican-variants-build.log`.
The final runtime-fixture adjustment passed its affected runtime test.
Browser inspection found a clipped batch action and truncated saved version near the toolbar breakpoint.
The implementation owner repaired only the affected Pelican layout and rebuilt the embedded server and installer.
See [implementation handoff](research/implementation-handoff.md) and [browser check](research/browser-check.md) for evidence and source paths.
The final generic label correction retains the approved `画面` wording.
Independent review found that an existing migration copy could retain its previous permissions.
The implementation now applies `0600` before writing and blocks saves if this step fails.
The existing runtime fixture and affected output rebuild passed after that repair.
The reviewer accepted the combined source and test evidence without another failed-run sequence or standard build.
See [final review](research/review.md).
The user approved the exact nine-file commit, message, and configured signing behavior.
Work commit `b67452e` contains those files. Its signature is valid.
