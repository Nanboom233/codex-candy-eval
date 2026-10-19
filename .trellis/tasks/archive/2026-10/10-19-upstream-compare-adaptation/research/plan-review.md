# Detailed Adaptation Review

## Review Status

The user requested a detailed review of each adaptation item on 2026-10-20.
The user approved the complete plan after all four behavior reviews.
This file preserves the discussion and source findings. `task-audit.md` defines verification ownership and removes duplicate check obligations.

## Item 1: Version Selection, Prompt Updates, and Preferences

### Source Findings

Upstream `web/pelican.js:89` initializes the version select from the shared variant definitions.
Upstream `web/components.js:43` selects the first option when the saved value does not match an available choice.
The first upstream variant is `original`. The control has no empty placeholder or separate confirmation action.
Upstream `web/app.js:538` saves toolbar changes through a shared per-control preference handler.
That handler differs from the native product's grouped preference object.

Native `patches/0006-evaluation-ui.patch:1869` saves current model and effort values together through `plSavePrefs`.
Native `initializeModelPickers` creates an empty model control before asynchronous startup finishes.
Native `start()` shows the page, awaits inspection settings, and then calls `fillModels` to restore the model.
Adding the variant control to the existing grouped saver can overwrite a saved model during that asynchronous interval.
This follows from the current call order. No runtime reproduction ran during this planning review.

Native `initializeCopyButtons` reads the adjacent code element's current `textContent` on every click.
It does not capture the initial prompt. Updating that element reuses the current copy action without another handler.
Native narrow layouts already place toolbar fields in two columns and the Pelican batch button across both columns.

### Approved Behavior

1. Keep field order: group, model, effort, version, and test action.
2. Offer `original` as `原版` and `all-in-one` as `多合一版`.
3. Default a first visit or an invalid saved choice to `original` through the existing select helper.
4. Restore a valid saved version for the current browser origin.
5. Update the visible prompt and save the version immediately when the version changes.
6. Reuse the current copy action so it copies the visible prompt.
7. Keep the current model and effort choices when the version changes.
8. Merge only the changed version field into stored Pelican preferences during a version change.
9. Keep the selected version across group changes. Version choice remains a browser preference, like model and effort.
10. Apply a changed version to later starts. An admitted run retains the request version it already received.
11. Keep the same native single-account and batch entry points. Both read the current version when submitted.
12. Keep narrow field pairs as group/model and effort/version, followed by the full-width test action.

Use the Go variant definitions as the authoritative source for IDs, labels, and prompts.
Inject that data during page assembly before initialization.
Keep the exact upstream prompt text and use `textContent` to show it.
Do not add another server settings write or an extra save action for this browser preference.

### Acceptance Evidence Needed

- A fresh preference object selects the original version and shows its prompt.
- A saved all-in-one version restores its exact prompt after a page reload.
- Version changes update copied text and stored version before a run starts.
- A version change before asynchronous settings resolve preserves the saved literal model and effort.
- Single-account and batch requests contain the current selected version.
- Polling and group changes preserve the version and prompt.
- An admitted request keeps its original version after the user changes the control.
- Actual browser inspection verifies the five-control layout at normal and narrow widths.

### Discussion State

The user approved this item with `无误` on 2026-10-20.
The PRD and design now include these decisions.

## Item 2: Request Validation and Output Extraction

### Source Findings

Pinned upstream source: `4f8367642cc1e629c4ee15a975e5e0fc53c4ac16:internal/plugin/pelican.go`.
Its request validator requires an explicit known variant before starting runs.
Native existing `StartRequest` callers omit this field, including the current runtime tests.
The native start branch can normalize the new field before account resolution and worker creation.

Upstream `extractCode` first finds a code block whose language matches the variant, ignoring case.
It returns that block's trimmed content without checking its root element or markup validity.
If no matching block exists, it takes the first code block whose trimmed content starts with `<`.
If no code block exists, it accepts the whole trimmed answer only when it starts with `<`.
Explanatory prose followed by unfenced markup does not pass that last rule.
An empty preferred-language block returns empty content. The extractor does not continue to another block.

Upstream `pelicanVariant.page` rejects empty extracted code and code larger than 262,144 bytes.
It applies the limit before adding the SVG display wrapper.
The original variant returns extracted code unchanged. The all-in-one variant adds the upstream HTML and viewport styles.
Both variants retain the native `html` result field for the existing preview and detail consumers.

The upstream test reuses an HTML-labeled code block for an all-in-one result containing an SVG element.
This confirms that the requested language is a preference, not a strict acceptance requirement.
The upstream code has no XML parser, root-element check, image-quality check, or Candy-answer scoring.

Native `askCodex` already owns three-attempt retries, usage accumulation, per-attempt deadlines, and cancellation.
Native Pelican already uses an extraction validator through that owner.
Native failed and cancelled results remain saved, whereas upstream drops a run when its context is cancelled.
Keep the native result lifecycle and add the variant before the supplier call.

### Approved Behavior

1. Send the selected variant explicitly from both native page start paths.
2. Treat omitted or empty native request variants as `original` to preserve current callers.
3. Reject any other unknown variant with a native HTTP 400 response before account lookup or supplier calls.
4. Resolve the exact prompt and expected language from the shared variant definition.
5. Keep the upstream extraction order and tolerant markup fallback rules.
6. Do not require an SVG root or add an XML parser for all-in-one answers.
7. Apply the existing 256 KiB limit to extracted bytes before adding the display wrapper.
8. Save the prepared page once through the current result path, with the selected variant on every recorded outcome.
9. Keep existing retry limits, delays, deadlines, usage accounting, and explicit stop behavior.
10. Interpret completed generation as accepted extracted output. Do not assign correctness to the depicted scene or Candy answer.

The selected route preserves upstream acceptance and avoids extra retries caused only by code-block labels.

### Acceptance Evidence Needed

- Missing variants run the original prompt. Both known explicit variants select their exact prompt.
- Unknown variants create no worker, account lookup, or supplier request.
- Preferred-language blocks win over earlier markup blocks, ignoring language case.
- The existing markup fallback, bare-markup rule, and empty-block behavior remain unchanged.
- All-in-one markup receives viewport styles. Original markup remains unchanged.
- Code at the byte limit is accepted. Code above the limit fails through the existing validator.
- Rejected output uses the existing bounded retries and accumulates observed usage.
- Failed and cancelled saved results retain the selected version and the native status.

### Discussion State

The user approved this item with `是` on 2026-10-20.
The PRD and design now include these decisions.

## Item 3: Saved Version Labels, History, Details, and Previews

### Source Findings

Upstream `web/pelican.js` adds version metadata to history cards and detail headers through `pelicanVariantMeta`.
Its latest-result row still calls the shared `latestResult` helper without version metadata.
Thus an unexpanded account row does not identify which version produced its latest result.
The toolbar selection belongs to the next request and can differ from that saved result's version.

Native `recordMeta` and `historyCard` already accept optional extra metadata.
Native `latestResult` currently shows verdict, model, effort, and usage, without an extra-metadata argument.
The native Pelican summary clears only HTML. A new variant field remains available in polling summaries.
Version labels therefore need no extra detail request.

Native preview cache keys contain account ID and record ID.
Each saved result has its own ID, so mixed versions do not need a new cache-key scheme.
Native detail fetching already rejects obsolete responses after another record opens or the dialog closes.
Retain those response guards and the existing `{type, pelican}` response wrapper.

The preview uses an 800 by 600 viewport in a 4:3 container.
List previews use scale 0.15. History previews use scale 0.3. Detail scaling follows available stage width.
The preview freeze script pauses browser animations, SVG animations, and timers after the initial display interval.
The detail frame omits that script, and the close handler clears its content.
The existing iframe sandbox permits scripts without giving model content the parent origin.

### Approved Behavior

1. Show `原版` or `多合一版` beside the saved model in the latest-result metadata.
2. Show the same version metadata in every history card and detail header.
3. Read labels from the saved record's version, not the toolbar's current selection.
4. Keep both versions in the same account history, ordered newest first, with the existing shared retention limit.
5. Keep history visible when the toolbar version changes. Do not add a version filter or another history section.
6. Retain version labels on failed and cancelled records and keep the current error presentation.
7. Use `画面`, `鹈鹕画面`, `放大查看`, and related picture labels where either static or animated output can appear.
8. Keep the current thumbnail, enlarged detail, preview cache, response guards, and iframe isolation.
9. Freeze thumbnail animations through the current script. Allow animations in the detail frame and clear them on close.
10. Keep the 4:3 preview viewport and existing responsive detail stage. Do not rewrite generated SVG coordinates or add image export controls.

The approved latest-result version label extends upstream display by one small metadata item.
It identifies mixed-version results before the user expands history.
Keeping the version only in history and details follows upstream more closely but requires expansion to identify a recent result.

### Acceptance Evidence Needed

- A latest all-in-one result shows its saved version while the toolbar selects original, and vice versa.
- Mixed-version history preserves ordering and shows each saved version in its card and detail header.
- Failed and cancelled records retain their version and visible error outcome.
- Polling summaries include the version while omitting HTML, without another request to resolve the label.
- Original HTML animation and all-in-one SVG both display in thumbnails and details.
- Actual browser inspection verifies thumbnail freezing, detail playback, close behavior, and narrow-width metadata wrapping.

### Discussion State

The user approved this item with `是` on 2026-10-20, including the latest-result version label.
The PRD, design, and execution plan now include these decisions.

## Item 4: History Migration and Executable Rollback

### Source Findings

Native `state.go` in patch `0001` accepts only explicit schema 1 or an absent marker.
It decodes current fingerprint records only after the explicit marker identifies the supported format.
The absent-marker route retains Candy and Pelican records and discards obsolete fingerprint and ModelTrace payloads.
All retained fields decode before loaded maps replace current memory maps.
The saver writes a same-directory temporary file with mode `0600`, then renames it over `history.json`.
Unreadable, malformed, null, and unsupported-version snapshots block later overwrites.

Upstream `migration.go` advances schema 1 to 2 and sets old Pelican variants to `original`.
Upstream `state.go:loadState` also writes the original bytes to a version-named backup before saving the migrated snapshot.
A backup-write failure blocks later state writes. Successful decoding has already made the retained records available in memory.
This backup behavior was not included in the initial native adaptation proposal.

The native installer backs up the executable. It does not restore evaluation history.
After a native schema-2 save, restoring the current old executable leaves history unreadable and saving blocked.
Keeping marker 1 would let the old typed saver discard new version fields during its next write.
Marker 2 therefore prevents silent version loss during an executable downgrade.
Restoring a pre-upgrade history snapshot also omits records created after that snapshot.

### Approved Behavior

1. Keep the history path, account keys, result identities, native map names, and latest-20 retention rule.
2. Read schemas 1 and 2. Keep the existing absent-marker rule for older retired data.
3. On schema-1 migration, retain Candy, current fingerprint, and Pelican records and mark old Pelican variants as `original`.
4. On absent-marker migration, retain Candy and Pelican records, mark Pelican variants, and keep the existing retired-data removal rule.
5. Save schema 2 after successful decoding and migration. A schema-2 file needs no migration write.
6. Keep malformed, null, and future-version overwrite protection and the current atomic write path.
7. Before replacing an old-format file, save one copy of its original bytes beside the history file.
8. Use the pre-upgrade copy only for migration recovery. Do not add periodic backups or installer changes.
9. If the original copy cannot be saved, retain readable memory records, report the error, and block file replacement and later saves.
10. If the migrated snapshot cannot be saved, retain memory records and the original disk file and report the save error.
11. Explain that an executable downgrade requires the matching old history snapshot and loses later records from the active history.

The approved decision changes the initial no-backup proposal after inspection of the upstream state saver and native rollback boundary.
It adds a single migration copy through the existing history owner, without a new storage service or migration framework.
The alternative keeps the current no-backup policy and makes a manual pre-upgrade history copy necessary for later downgrade recovery.
The user approved the migration-copy scope.

### Acceptance Evidence Needed

- Schema-1 migration retains current fingerprint data and old Pelican HTML, IDs, metadata, and status fields.
- Both older routes label retained Pelican records as original and save schema 2.
- The absent-marker route preserves its existing retired-data removal rule.
- Schema-2 reload retains both variants without another migration write.
- Original-copy bytes match the input bytes, before any migration or serialization.
- Copy failure protects the original file and blocks later writes with a visible error.
- Migration-save failure keeps the original file, copied snapshot, and loaded memory records.
- Invalid and future-schema files remain unchanged after attempted saves and clears.

### Discussion State

The user approved this item with `是` on 2026-10-20.
The PRD, design, and execution plan now require the original-file migration copy.

## Item 5: Validation and Delivery

### Source Findings

`scripts/build.py --check` applies all seven patches to a copy of the original Sub2API source.
It runs the UI harness, focused backend tests, affected-package vet, frontend compilation, server build, and focused installer tests.
It packages the existing Linux/amd64 server into the self-extracting installer.
The focused selector includes `Codex` and `ExtractHTML`. New variant and migration cases must remain selected.

The existing release workflow builds against Sub2API v0.2.17 and uploads the installer at its fixed path.
The comparison's release-only commit changes the CPA plugin version to 0.4.2.
It supplies no native Sub2API version or installer behavior change.
The native product version remains Sub2API v0.2.17.

### Validation Scope

- Extend the current Pelican and runtime tests for approved extraction, variant routing, statuses, and migration-copy behavior.
- Extend the current UI harness for preference preservation, question copying, submitted variants, mixed history, and saved version labels.
- Keep synthetic model responses and temporary history files. No supplier quota or real user history is needed.
- Run the existing build/check entry and inspect the actual page at normal and narrow widths.
- Exercise both variants, thumbnail freezing, detail playback, close behavior, and saved preference restoration.
- Reuse valid unchanged transport and inspection evidence. Do not add full-suite, race, coverage, or architecture matrices.

### Delivery Scope

- Keep production changes in patches `0001` and `0006`, with tests under `tests/`.
- Update the focused test selector only if renamed or added tests require it.
- Update Chinese usage and history recovery guidance and the affected English specifications.
- Deliver `dist/build/sub2api-codex-eval` and `dist/build/sub2api-codex-eval-installer`.
- Keep installation behavior, native host version, and release workflow unchanged.
- No remote publication, production deployment, or release creation forms part of this task.

### Final Planning Gate

Items 1 through 4 have explicit user approval. No unresolved product choice remains.
The user approved the final summary with `无误` and requested removal of redundant tests and added requirements.
The main session completed that audit in `task-audit.md`. Task activation and implementation are authorized.
