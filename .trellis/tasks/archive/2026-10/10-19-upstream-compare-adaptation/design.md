# Native Pelican Variant Design

## Change Boundary

Adapt upstream behavior to the existing native module. Keep production edits in patches `0001` and `0006`.
Keep tests under `tests/`. Change the focused build selector only if test names require it.
Keep Sub2API version `0.2.17` and the local plugin version separate from upstream release metadata.
The release-only CPA version update does not change the native host version.

## Shared Variant Definition

Define the two upstream variants in the native `pelican.go` owner.
Each definition contains its ID, label, exact prompt, and private markup language.
Use the definition for request validation, prompt routing, output extraction, and page configuration.
Inject the public variant definitions through `web.go` into the existing Pelican script.
Do not maintain a second prompt table in JavaScript.

## Request and Result Flow

Add optional JSON field `variant` to `StartRequest`.
In the Pelican admission branch, map an absent or empty variant to `original` and reject unknown IDs with HTTP 400.
Validate before account resolution and goroutine creation.
Pass the resolved variant to the native runner and evaluator.
Initialize the result version before sending so every recorded outcome retains it.
Keep the existing native target, result identity, state maps, detail wrappers, and task lifecycle.
Candy and fingerprint requests ignore the Pelican-only field.

The shared retry runner calls the selected variant's output validator.
Prefer the requested code-block language, then keep the current markup fallback rules.
Match language labels without case sensitivity. Do not add XML parsing or root-element validation.
If no preferred block exists, use the first block whose trimmed content starts with `<`.
Accept an unfenced answer only when it contains no fenced block and its trimmed text starts with `<`.
An empty preferred block remains rejected. Do not continue to a later block in that case.
Limit extracted code to 256 KiB before adding the SVG display wrapper.
Measure extracted bytes before the wrapper. Keep original output unchanged and wrap all-in-one output with the upstream styles.
Wrap SVG with the upstream viewport styles and retain the native HTML result field.
Keep existing iframe sandbox, CSP, preview freeze, page cache, and detail playback.
The all-in-one image has no automatic answer or image-quality score.
Use the shared retry owner without changing delays, deadlines, usage accounting, or cancellation.
Keep failed and cancelled native records, including their version. Do not copy upstream cancellation's record-discard behavior.

## Native History Upgrade

Use a named schema constant with value 2 in native `state.go`.
Accept the existing unversioned format, schema 1, and schema 2.
Decode fingerprints for schemas 1 and 2. Keep the existing unversioned retirement rule.
For pre-2 files, set saved Pelican versions to `original` after retained fields decode successfully.
Before replacing an older snapshot, save its original bytes in the same directory with mode `0600`.
Use the upstream version-based backup name, such as `history.v1.json` or `history.v0.json`.
Keep this migration copy outside the atomic saver's temporary path. Do not add periodic backups.
If copying fails, retain decoded records in memory, expose the error, and block migration replacement and later saves.
Use the current atomic saver once to persist the upgraded snapshot.
Do not rename native `results`, alter record identities, or synthesize status and attempt fields.
Preserve unreadable-file protection, storage errors, retention, and failed-clear rollback.
Keep future schemas invalid. A failed upgrade save retains loaded data and reports the existing storage error.
The previous binary cannot read schema 2. An executable rollback therefore also needs a compatible pre-upgrade history file.
Restoring an older history snapshot excludes later records from the active history.
Keep backup and migration in the current history owner. Do not change the installer transaction.
Do not migrate real user data during development.

## Native Page Adaptation

Use native IDs `pl-variant` and `pl-prompt` with the existing `pl-` controls.
Add the version control after `pl-effort` and extend the existing Pelican preference object.
Default missing or obsolete saved choices to `original` through the normal select helper.
Restore a valid saved choice for the current browser origin and retain it across group changes.
On a version change, merge only the version field into stored preferences. Preserve the saved model and effort during asynchronous startup.
Update and save the version immediately. Keep model and effort controls unchanged.
Use `textContent` for the displayed prompt so the existing copy action reads the current question.
Use each admitted request's version for its full lifetime, even when the user changes the control.
Add saved version metadata to latest-result rows, history cards, and detail headers.
Read the saved version from each record. Do not derive it from the current toolbar selection.
Keep mixed versions together in newest-first history with the current shared latest-20 limit.
Toolbar version changes do not filter or relabel saved records.
Show the version on failed and cancelled records through the current outcome and error presentation.
Keep summary versions while omitting HTML. Labels need no extra detail request.
Use picture labels for shared previews and dialog controls.
Keep thumbnail animation freezing, full detail playback, and frame removal when the detail dialog closes.
Retain the 800 by 600 viewport, 4:3 preview container, record cache keys, and stale-detail response guards.
Retain the native group field, model picker, eligible targets, switches, and Settings initialization.
Adjust only the Pelican toolbar tracks and necessary responsive rules.
Keep field order as group, model, effort, version, and test action.
At narrow widths, pair group/model and effort/version. Keep the test action on its own full-width row.
Keep the previously removed fixed request-count caption absent.

## Verification and Delivery

Extend native Pelican, runtime/persistence, and UI tests for the changed contracts.
Use synthetic responses and temporary histories. Do not call a real model for validation.
Run the established build/check entry and inspect the effective browser page.
Reuse unchanged transport, installer, and inspection evidence where the source and inputs remain valid.
Update the native runtime, persistence, and error-handling specifications after implementation.
Update the Chinese Pelican instructions and existing recovery guidance for the original-file migration copy.
Keep both existing Linux/amd64 output paths and the current installer package contents.
The existing build/check entry covers focused backend tests, UI assertions, vet, frontend build, and installer checks.
Inspect the changed visual result at normal and narrow widths. Do not repeat assertions already covered by the UI harness.
Keep remote publication, release creation, and production installation outside this task.
Use [task audit](research/task-audit.md) as the verification ownership map. Reuse existing checks for unchanged behavior.

## Source Evidence

See [upstream comparison research](research/upstream-comparison.md) for commit identities and native owner anchors.
