# Control degraded models and deliver a Go installer

## Goal

Remove selected whitelist entries from degraded accounts.
Restore those entries when inspection confirms that the account returns to normal.
Deliver the resulting server through one self-extracting Go installer with no installation arguments.
Keep the existing interactive selection, installation steps, and recovery behavior.

## Requirements

- R1: Add a `模型控制` section beside `优先级` in the Settings tab.
  Use an independent switch that defaults to off. Show the model list only when the switch is on.
- R2: Label the editable list `降级移除模型`.
  Its default entries are `gpt-6.1-sol`, `gpt-6-sol`, `gpt-5.6-sol`, and `gpt-6-astra`.
  Copy the native Sub2API whitelist control's structure and styling for removable tags, search, selection, and custom entry.
  Keep custom model entry available when a name is absent from the catalog.
  Do not add small explanatory text to the section or its controls.
- R3: Target only identity entries in `accounts.credentials.model_mapping`, with the format `model -> model`.
  When inspection confirms degradation, remove entries whose model names appear in the configured removal list.
  Do not add absent models. Keep other entries unchanged.
- R4: When inspection confirms recovery, restore the models that model control removed from that account.
  Restore them as their original identity entries.
  Use saved removals even if the global list changes or the switch becomes off.
  Preserve a currently present entry when an administrator replaces its mapping before recovery.
- R5: Use the saved inspection verdict. Model control must work when automatic priority adjustment is off.
  After each account completes detection and the existing confirmation rules establish its state, apply model control and priority for that state.
  Include stable states. Skip writes when the account already matches the target.
- R6: Use the existing settings persistence and account update paths.
  Keep unrelated account data and existing manual evaluation behavior unchanged.
- R7: If model removal empties the complete `model_mapping` object, stop ordinary real-model requests for that account.
  Keep recovery inspection available. Restore ordinary model access after confirmed recovery restores the removed entries.
  Use the `"blacklist": "blacklist"` placeholder instead of an extra scheduler block state.
- R8: Deliver one self-extracting Linux/amd64 Go executable containing the built server and required license texts.
  The target host must not need Python or external extraction tools.
- R9: Port the current installer logic to Go with a parameter-free installation entry.
  Remove `install`, `--container`, `--binary`, `--backup-dir`, and `--timeout` from the new installation interface.
  Do not read, parse, or validate command-line arguments.
  Keep numbered interactive container selection, the embedded candidate, private backups under `/tmp`, and a fixed 120-second readiness deadline.
  First show running containers whose names contain `sub2api`, without case sensitivity.
  Offer `More` as the next numbered option to show all running containers when the first list omits any.
  If no name matches, show all running containers directly.
  Preserve four progress steps, readiness checks, backup retention, and failure or interruption restoration.
- R10: Reuse `scripts/build.py` and the release workflow to generate and publish the single installation asset.
  Installation must not download sources, compile, or fetch helper images.
  Keep Docker CLI access as the existing host requirement.
- R11: When native navigation opens account evaluation on the saved Settings tab, refresh the shared group data.
  Switching to an evaluation tab must load its selected group's accounts without a manual refresh.
  Preserve saved settings, filters, tab selection, and existing native authentication.
- R12: Rewrite the complete README in strict `output` style.
  Use Chinese as the user requested. Apply the skill's plain-language principles to Chinese prose.
  Write for users. Keep installation, use, recovery, and practical build instructions.
  Remove product decisions and internal implementation details. Do not move them into another document.
  Omit version numbers, version checks, and upgrade descriptions from the README.
  Keep necessary commands, paths, feature behavior, and recovery conditions accurate.
  Use the original CPA Codex Candy Eval README's usage and acknowledgments sections to supplement the guide.
  Adapt usage to Sub2API and retain verified original project and contributor links.
  Use short active sentences, consistent terms, and numbered procedures.

## Acceptance Criteria

- AC1 (R1): Settings shows separate priority and model control switches at the same section level.
  The model list follows only the model control switch.
- AC2 (R2): The first settings read supplies the four exact defaults. Edits persist through the existing automatic save flow.
  Users can search, select, remove, and add model IDs through the native-style controls.
  Saved IDs remain visible when the catalog lacks them. An empty catalog still permits custom additions.
  Keyboard controls work, and the native two-column tags fit narrow widths. No small explanatory text appears.
- AC3 (R3): Confirmed degradation removes only matching identity entries. An absent model causes no model-data change.
  Preserve aliases, wildcard entries, and unrelated identity entries.
- AC4 (R4): Confirmed recovery restores only the identity entries that model control removed.
  Edited global settings do not change the restoration set. Recovery preserves an administrator's replacement entry.
- AC5 (R5): The removal and recovery lifecycle works with priority adjustment off.
  Stable degraded accounts apply the current removal settings after a valid detection.
  Stable normal accounts restore saved removals after a valid detection.
  Priority adjustment uses the confirmed state, including stable states, and skips matching priorities.
  Once detection and required confirmations establish its state, each account applies its result without waiting for unrelated accounts.
  Keep the existing confirmation-group order and recovery wait.
- AC6 (R6): Existing priority controls, participation, confirmation rules, manual tests, and unrelated account fields retain their behavior.
- AC7 (R7): Removing the final mapping entry blocks ordinary real-model requests but permits recovery inspection.
  The stored mapping contains `"blacklist": "blacklist"` instead of an empty object.
  Confirmed recovery restores the removed entries and removes the reserved identity placeholder.
  A remaining alias or wildcard keeps its native admission behavior.
- AC8 (R8): A single executable contains the server and required license texts and runs without Python on the target host.
- AC9 (R9): Running the executable without arguments starts numbered container selection, followed by the existing four installation steps.
  The first list contains name matches for `sub2api`. Selecting the numbered `More` option expands it to all running containers.
  If there are no matches, the first list contains all running containers. Omit `More` when the list already includes all containers.
  Empty input has no default. Omit the Cancel option and the `q` shortcut. Users can interrupt with Ctrl-C.
  No installation argument or replacement configuration prompt is available.
  A single match still requires explicit selection. Numbers select entries from the currently displayed list.
  The installer uses its embedded server, prints the retained backup path, and waits up to 120 seconds for readiness.
  All effects use the selected immutable container ID.
- AC10 (R9): Version, architecture, readiness, timeout, stopped-container, and interruption behavior match the current installer contract.
- AC11 (R10): The build produces the raw server and self-extracting installer. The release publishes the installer as its installation asset.
- AC12 (R11): Enter from Sub2API with Settings as the saved tab, then switch to an evaluation tab.
  The group selector and account rows load through the existing data path.
  Reentry must not require another tab switch or manual refresh.
- AC13 (R12): The complete README follows strict `output` rules and gives users practical installation, use, recovery, and build instructions.
  The guide uses Chinese and incorporates applicable original usage and acknowledgments from the actual source README.
  Build instructions retain prerequisites, commands, and output paths. Product decisions and internal implementation details are absent.
  It contains no version or upgrade descriptions. Literal model IDs remain unchanged.
  Documentation checks cover sentence lengths, references, commands, and changed claims. No product rebuild is necessary.

## Scope Limits

- Limit this task to inspection model control, Go installer delivery, the Settings-entry repair, and the strict README rewrite.
  Use focused checks for changed behavior.
- Do not extend this task to unrestricted accounts, wildcard exclusions, alias targets, automatic passthrough, or compact-only mappings.
- Keep server changes in the existing ordered patches. Keep installer/build changes in their existing delivery area and tests under `tests/`.
- Preserve the pending degraded-tag task and its unrelated working changes.
- Do not use production accounts or provider quota for local checks.
- Keep this delivery limited to Linux/amd64. Do not add source downloads, a general extraction command, or another installation interface.

## Evidence

- `patches/0005-inspection.patch:25`: `InspectionAccount` contains a saved verdict independent from priority.
- `patches/0005-inspection.patch:487`: confirmed transitions save the verdict before the optional priority change.
- `patches/0006-evaluation-ui.patch:2344`: the existing priority section provides the adjacent settings pattern.
- `.trellis/spec/backend/native-inspection.md`: settings use native storage, automatic saves, and existing confirmation rules.
- [Account model storage investigation](research/model-storage.md) traces the schema, frontend serialization, account updates, runtime admission, catalogs, and shadow accounts.
- [Native model UI investigation](research/native-model-ui.md) traces the whitelist selector, native setting rows, and embedded settings integration.
- `scripts/install.py` and `.trellis/spec/backend/minimal-linux-delivery.md` define the current installation mechanism.
  R9 replaces their argument interface with the user's parameter-free installation decision.
- [Go installer investigation](research/go-installer-parity.md) records the installation contract, build boundary, and focused verification scope.
- Account model restrictions use PostgreSQL `accounts.credentials.model_mapping`.
  Whitelist entries use identity mappings. Explicit aliases use request-side keys and upstream target values.
- An empty mapping selects native default admission. Wildcards can still match a removed exact key.
  Automatic passthrough bypasses the mapping restriction.

## Technical Constraints

- Preserve other credential fields. The generic account credential update does not merge absent ordinary fields.
- Keep only one new account metadata list for actual removed identity model names.
  Do not store a complete whitelist, authentication credentials, a block flag, or an original scheduling value.
- Save changed model data, restoration data, and the verdict together before the optional priority update.
  Failed model persistence must not leave a saved transition without its model changes.
- Preserve the existing recovery-test entry point, which bypasses ordinary gateway scheduling.
- Keep the existing semantics for mappings outside the confirmed identity-entry scope.
- Keep native `schedulable` unchanged. Use the existing normal/degraded verdict and existing scheduler admission.
- `blacklist` remains an ordinary exact key, not a native deny keyword. A literal request for that name still matches.
- See [the design](design.md) for fields, confirmed-state ownership, failure behavior, and recovery.
- See [the execution plan](implement.md) for ordered work, ownership, and bounded acceptance checks.

## Status

The model-control implementation and independent checks are complete. Its native UI and spec synchronization are complete.
See [implementation evidence](research/implementation-handoff.md) and [independent review](research/check-handoff.md).
The user added the Go installer to this same task and requires planning before that migration.
The Go installer correction and implementation scope audit are complete.
The entry does not read process arguments. Added SIGTERM behavior and duplicate payload checks are removed.
HTTP readiness retains the original direct-response behavior without redirects.
See [installer implementation evidence](research/go-installer-implementation.md) and [independent delivery review](research/go-installer-check.md).
The new installer was not installed into an actual Docker target. The approved transaction checks used fake Docker and local HTTP.
The affected checks and installer-only rebuild passed.
The confirmed-state timing correction passed affected checks and the delivery rebuild.
See [the correction evidence](research/inspection-timing-implementation.md) and [the independent review](research/inspection-timing-check.md).
Work-commit approval and Trellis finish remain.
The Settings-entry group-refresh repair passed focused UI, actual browser, and delivery build checks in this same task.
See [the repair evidence](research/settings-entry-implementation.md) and [the passing independent review](research/settings-entry-check.md).
