# Model control and Go installer design

## Boundary

Use the confirmed identity-whitelist scope in [the PRD](prd.md).
The storage investigation supplies the native schema and update evidence in [model-storage.md](research/model-storage.md).
The [native UI investigation](research/native-model-ui.md) supplies the control pattern and embedded-page integration evidence.
Keep scheduling, account availability, authentication, and the direct evaluation sender unchanged.

## Settings-entry data loading

Trace native navigation, restored tab selection, group refresh, and evaluation-tab account loading before editing.
Fix the existing page lifecycle owner. Reuse the shared group and account-loading paths.
Restoring Settings currently starts its settings read before `start()` owns initialization.
The later settings response changes the selected group and invalidates the concurrent group-catalog response.
Let `start()` own the first settings read before the first group refresh.
Keep settings refresh on later user tab entries. Keep obsolete-response protection.
Keep saved settings edits, filters, and native session behavior.
Do not add a refresh control, background retry mechanism, or new state store.
Keep production changes in patch `0006` and focused regression checks in the existing UI harness.
Use an actual browser to exercise Settings entry and a later evaluation-tab switch.

## Settings and form

| Field | Type | Default |
| --- | --- | --- |
| `control_models` | Boolean | `false` |
| `degraded_remove_models` | String array | `gpt-6.1-sol`, `gpt-6-sol`, `gpt-5.6-sol`, `gpt-6-astra` |

Use the existing inspection settings GET and PUT endpoints.
The settings service stores both fields through its existing JSON field persistence.
An absent list selects the defaults. An explicit empty list stays empty and removes no models.
Trim entries and remove duplicates while retaining their order.
Validate each entry with the existing model-name rules. Reject wildcard entries and the reserved placeholder name `blacklist`.
Compare names exactly after normalization.

Place `模型控制` directly after `优先级` as a sibling settings section.
Reuse the existing `.settings-card`, header, body, and responsive layout.
Use an accessible header switch with no separate visible switch caption or subtitle.
Show one field labelled `降级移除模型` while `control_models` is on.
Copy `ModelWhitelistSelector.vue` structure, dimensions, colors, and control placement into the existing HTML and JavaScript page.
Do not redesign the control or mount the Vue component.

Render selected models in the native two-column grid with 6 px gaps.
Use gray 12 px tags, 8 by 4 px padding, and the native 14 px model icons.
Each tag has a remove button with a model-specific accessible label.
Keep the native bottom separator and right-aligned dropdown arrow. Omit the model count.
Use a keyboard-accessible dropdown trigger without changing its visible layout. Do not nest buttons.
Give the trigger `aria-expanded` and `aria-controls`.
Use the native absolute dropdown below the field, with a 4 px gap, 8 px corners, and its shadow.
Keep the bordered search area and the 208 px maximum list height.
Show focusable option controls with the native selected state.
Use `catalogCache.models.ids` for current-group options. Do not add discovery requests.
Include selected IDs that the current catalog does not contain. Keep the selected array order.
Keep custom model entry available when the catalog lacks a name or contains no options.
Place custom entry below and outside the selection field, as the native component does.
Keep its `自定义模型名称` label, input, and `填入` text button.
Do not replace this row with an inline input or a compact plus button.
Trim input, ignore empty input, and prevent duplicate entries. Clear the input after a successful addition.
Apply the same normalization and validation to custom additions.

Support Tab navigation, Enter/Space activation, and Escape dismissal with focus returned to the trigger.
Close the dropdown on an outside click.
Do not add small explanatory text, counts, help paragraphs, tooltips, copy actions, or synchronization actions.
Keep successful automatic saves silent. Use the existing error area for actual validation or save errors.
Keep this switch independent from automatic inspection scheduling and priority adjustment.

Extend the existing form fill and serialization functions for the array field.
Use the existing debounce, serialized saves, field revisions, and read invalidation.
Keep the list when its section collapses. Preserve newer edits during reads and saves.

## Account data

Use only one new account metadata key:

```json
{"codex_eval_inspection_removed_models": ["gpt-6.1-sol", "gpt-6-sol"]}
```

The list contains only names that this feature actually removed.
Each name determines its original identity value, so no separate mapping snapshot is necessary.
Keep this list in `accounts.extra`, beside the existing inspection verdict.
Do not store credentials, a block flag, an original scheduling value, or a recovery phase.

Protect this metadata key during ordinary full account edits, as the integration already protects participation and verdict.
Treat an absent or cleared list as no pending restoration.
Use an empty list to clear pending restoration through the existing metadata format.

## Confirmed state owner

Extend `InspectionHooks.PersistVerdict` to receive the round's `InspectionSettings` with the confirmed verdict.
Keep a single persistence callback. Do not introduce a second optional callback or a fallback dispatch path.
After each account completes detection and existing confirmation rules establish its state, call this owner for that state.
Stable states need no extra confirmation. Apply each ready account without waiting for unrelated accounts or recovery confirmations.
Keep the existing confirmation-group order. This correction changes application after an established state, not confirmation scheduling.
Preserve cancellation, participation, eligibility, and mixed-result handling.
Use the round settings snapshot. Apply newer settings on the next detection.
The administrator adapter delegates account changes to a native service helper on `AccountTestService`.
That service already owns the account repository and direct evaluation integration.
No administrator route, repository interface, database table, or scheduler API is necessary.

The helper reads the current account and builds copies of mutable maps before it changes them.
Preserve all unrelated credential and metadata fields.
When model data changes, write the mapping, pending restoration list, and verdict with one `AccountRepository.Update` call.
The existing native update transaction saves these fields together and publishes the scheduler update.
Use the existing key-only `UpdateExtra` verdict path when only the verdict changes.
Skip persistence when the verdict, mapping, and pending removals already match.
Apply priority after successful persistence. Skip a matching priority.

Do not submit a partial credential object through `AdminService.UpdateAccount`.
The generic account merge retains sensitive keys only, so that path can delete absent ordinary credential fields.

## Degradation

After detection establishes a degraded state, apply model removal even when the saved verdict remains degraded.
Read the round's model settings rather than newer settings from another session.
If model control is off, preserve pending removals and save the verdict only when it changes.

For each configured name, remove the entry only if its stored key and value both equal that name.
Ignore absent names. Preserve other identity entries, aliases, wildcard entries, and compact mappings.
Do not change automatic-passthrough accounts because their model restrictions are outside the confirmed scope.
Do not create a mapping for an unrestricted account.
Merge actual removals into the pending restoration list without duplicate names.

If actual removal leaves no mapping entries, insert this placeholder:

```json
{"blacklist": "blacklist"}
```

The native whitelist stays nonempty and rejects ordinary real model names.
Keep `blacklist` reserved as the feature's placeholder. It is not a Sub2API deny keyword.
A literal request for `blacklist` still matches the native mapping rule.
Do not add a gateway rule solely for this placeholder name.

## Recovery

After detection establishes a normal state, restore pending removals even when the saved verdict remains normal.
Use the account's saved removed-model list, even if the global list changed or model control is now off.
The switch stops new removals. It must not strand a pending restoration.

Restore an absent name as `name -> name`.
Keep a currently present entry, including an administrator's replacement mapping.
Remove the reserved identity placeholder when restoring a pending removal list.
Clear that list and save the normal verdict in the same account write.
Do not replace the complete whitelist with an earlier snapshot.

The direct inspection sender bypasses ordinary gateway model admission.
The placeholder therefore leaves recovery inspection available without changing native `schedulable`.
Other native availability restrictions and account participation still apply.

## Failure and existing behavior

| Condition | Result |
| --- | --- |
| Account read or model update fails | Keep the previous saved transition data. Record the error. Do not write priority. |
| Optional priority update fails after transition persistence | Keep the saved verdict and model changes. Use the existing priority error behavior. |
| Confirmation is incomplete, mixed, cancelled, or opted out | Keep the current verdict and model data |
| Verdict stays stable | Apply model control and priority for that state. Skip matching writes. |
| Model control becomes enabled for an already degraded account | Apply removals after its next valid detection. |

Do not add automatic write retries, background synchronization, or new result fields.
Notify only after an actual priority write succeeds. Keep existing notification error handling.
Notification text must describe the confirmed state without claiming three tests for a stable-state correction.
Use the current result error owner for model persistence failures.

## Delivery and recovery

Production additions belong in patch `0005`. Form changes belong in patch `0006`.
The managed metadata preservation change belongs in patch `0007`.
Keep tests under `tests/` and preserve the pending degraded-tag edits in the shared UI files.
No schema migration or data backfill is necessary.

Older binaries do not restore this feature's pending removals.
This limitation belongs to a later runtime rollback decision. It adds no extra deployment gate.

## Self-extracting Go installer

The user added delivery migration to this same task. Preserve the completed model-control implementation and evidence above.
Build a static Linux/amd64 installer with the Go standard library and CGO disabled.
Embed a compressed payload containing the built server and required license texts.
Keep `dist/build/sub2api-codex-eval` as the raw server output.
Use `dist/build/sub2api-codex-eval-installer` as the single installation asset, replacing the archive/Python release entry.
Retain `scripts/build.py` as the developer build entry. The target host does not need Python.

The user downloads the installer, grants execute permission when needed, and runs it directly:

```sh
chmod +x sub2api-codex-eval-installer
./sub2api-codex-eval-installer
```

Docker CLI access and permissions remain necessary. The host needs no Go toolchain or external extraction tools.

Remove `install`, `--container`, `--binary`, `--backup-dir`, and `--timeout` from the installation interface.
Read running containers through the existing Docker listing and inspect path.
First show containers whose names contain `sub2api`, without case sensitivity. Match the name, not the image or version.
When this list omits running containers, append `More` with the next numeric index after the last container.
Selecting that index displays all running containers in one numbered list. Do not accept a separate `more` command.
If no name matches, show all running containers directly. If none are running, retain the existing error and exit behavior.
Show each container's name and shortened ID. Resolve a number against the currently displayed list.
Even a single match requires selection. Empty input has no default. Repeat the prompt after invalid input.
Omit the Cancel option and the `q` shortcut. Keep Ctrl-C interruption behavior.
Use the selected immutable container ID for all subsequent Docker operations.
Always use the embedded server and retain the original executable in a private directory under `/tmp`.
Use a fixed 120-second readiness deadline for installation and restoration.
Do not replace the removed arguments with environment settings, configuration files, or additional prompts.
Do not read or validate process arguments. The entry directly starts interactive installation.
Use one numeric selection prompt for containers and `More`. Keep four progress steps, backup output, startup progress, and final success/error messages.
Use Docker CLI with immutable selected container IDs.
Keep native readiness, version acceptance, executable checks, retained backups, stopped-container restoration, and interruption handling.

The [installer investigation](research/go-installer-parity.md) records the source-backed parity contract.
Keep Go installer source under `scripts/installer/` and tests under `tests/installer/`.
The build copies the source into a temporary module with the gzip server payload and five unchanged license texts.
Use standard-library imports only. Do not add a root Go module or tracked generated payloads.
Use the Go toolchain already selected from the upstream module.
For checks, copy the installer tests into that temporary module and run them for the build host.
Use separate build environment values for the final Linux/amd64 executable.
Embed the compressed server and required license texts. Do not add a license command or a general extraction command.
The default installation path retains its four visible steps. It adds no confirmation or extraction step.
Extract the embedded server into a private temporary directory and remove that directory on exit.
Retained executable backups use a separate directory and remain available.

Preserve the exact 20-byte ELF check and version boundaries in `scripts/install.py`.
Accept installed versions 0.2.16 and 0.2.17. Require candidate version 0.2.17 before backup and replacement.
Retain staged candidate cleanup, backup-before-replacement ordering, and direct Docker copy for stopped-container restoration.
Each readiness attempt reads current Docker facts and checks `/health` and `/setup/status`.
Keep the original direct HTTP 200 requirement. Do not follow redirects.
Retry restart and transport failures within the current deadline. Keep each HTTP request within three seconds and the remaining deadline.
After mutation, an installation error or interrupt restores the backup, restarts the container, and verifies the original application.
An interrupt must not cancel that restoration path. Preserve the existing exit behavior and printed backup location.

Replace archive packaging with static installer compilation in the existing build entry.
Publish only `sub2api-codex-eval-installer`. Keep the raw server as a local build output.
Remove the Python installer and obsolete archive entry after the Go replacement passes its focused checks.
Update README, repository instructions, delivery specs, and release references together.
Keep the sole third-party notice source under `licenses/`. Update its delivery paths while retaining attribution and license terms.

The retained container uses `/app/sub2api` as its fixture launcher and `/app/sub2api-bin` as its server.
Do not overwrite that launcher for installer checks.
Use a fake Docker command and local HTTP responses to check the changed Go installation transaction.
Reuse completed model-control, browser, and server evidence for unchanged inputs.
Container recreation still discards writable-layer replacement and requires installation again.
The user approved this installer scope after the task review. Implement it within this same task.
