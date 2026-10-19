# Implement model control and Go installer delivery

## Current phase

The initial model-control implementation passed its checks. Reuse evidence for unchanged UI and metadata paths.
The user added the Go installer to this same task and requested planning before migration.
The installer parity research and expanded design are recorded in this task.
The corrected Go installer implementation and independent scope review are complete.
See [the implementation report](research/go-installer-implementation.md) and [the delivery check](research/go-installer-check.md).
The correction removed argument inspection, SIGTERM handling, and duplicate payload checks. HTTP readiness no longer follows redirects.
The affected entry and redirect checks passed. The installer was rebuilt from the unchanged server output.
The user then requested model control and priority after each completed, confirmed account detection.
The correction passed focused tests, affected-package vet, the final server/installer build, and independent full-scope review.
See [the implementation report](research/inspection-timing-implementation.md) and [the independent review](research/inspection-timing-check.md).
The user approved the work commit. Signed commit `c4cc8e8` contains the approved work, and its signature verified successfully.
The implementation and checks are complete. Trellis archival, session recording, and return to local `main` complete the handoff.

## Strict README rewrite

The README is a Chinese user guide for installation, evaluation, inspection, recovery, and practical builds.
Build instructions retain dependencies, commands, and output paths.
The usage and acknowledgments use the original README at `cdd4fb9:README.md` as their source.
The guide adapts original evaluation steps to current Sub2API groups and accounts.
It includes fingerprint request counts, record controls, request costs, timeouts, and original acknowledgments.
The guide omits version and upgrade descriptions, product decisions, and internal implementation details.
The author and independent reviewer passed focused documentation checks. No finding remains.
The review confirmed original usage counts, costs, timeouts, record controls, credits, and current commands.
The final Chinese README is in work commit `c4cc8e8`, including the user's final staged edits.
No product check or build ran for this documentation change.
Earlier checks remain in [the rewrite evidence](research/readme-strict-implementation.md) and [the review](research/readme-strict-check.md).
Existing product evidence remains valid. Final work-commit approval is complete.

1. Read the complete README, the strict `output` skill, and direct product/build contracts.
2. Rewrite the complete document with short active sentences and numbered procedures.
   Use Chinese as the user requested. Supplement usage and acknowledgments from the original source README.
   Keep necessary commands and user behavior accurate. Preserve practical build instructions.
   Omit version descriptions, upgrade descriptions, product decisions, and internal implementation details.
3. Review the complete rewrite for meaning, strict style, and local references.
   Use documentation checks only. Reuse all valid product evidence.
4. Update the pending work-commit proposal after independent documentation review. Do not commit without final approval.

## Settings-entry repair

The three-line startup repair passed the existing UI harness and the actual native sidebar/iframe flow.
The normal build regenerated both delivery outputs. Independent review passed with no remaining findings.
See [the repair evidence](research/settings-entry-implementation.md) and [the independent review](research/settings-entry-check.md).

1. Assign the page lifecycle investigation, patch correction, and existing UI regression case to one implementer.
2. Reproduce native entry with Settings as the saved tab. Trace the missing group or account load.
3. Fix the existing owner and run the affected UI harness. Exercise the reported flow in an actual browser.
4. Rebuild the raw server and installer with the existing build entry. Reuse unchanged backend and installer checks.
5. Obtain independent review of this correction. Update the affected source-backed spec and task evidence.

## Confirmed-state correction

1. Assign the coordinator, persistence helper, and affected existing tests to the original model-control implementer.
2. Apply controls when each account's detection establishes its state under the existing confirmation rules.
   Include stable states. Keep confirmation counts, group order, recovery waits, cancellation, participation, and eligibility rules.
3. Reuse one account-application owner. Save model data and verdict before optional priority changes.
   Skip unchanged writes. Merge actual removals without duplicates. Keep restoration behavior after switch-off.
4. Update existing focused lifecycle and service tests for stable-state application and unchanged writes.
   Do not add a second test suite or extra runtime fixtures.
5. Run affected focused Go tests and vet. Rebuild the raw server and installer through the existing build entry.
   Reuse unchanged UI and installer checks.
6. Assign independent review to the existing model-control checker. Synchronize affected specs and commit proposal.

## Completed work and evidence

The settings, transition helper, native model selector, and metadata protection are complete in patches `0005`, `0006`, and `0007`.
Their requirements remain in [the PRD](prd.md) and [the design](design.md).
The [implementation handoff](research/implementation-handoff.md) and [independent check](research/check-handoff.md) record the completed checks.
Reuse the focused Go, UI, build, and actual-render evidence for unchanged inputs.
Only a later change to those inputs requires the affected check again.

Begin each dispatch prompt with `Active task: .trellis/tasks/10-19-inspection-model-control`.
Use native context injection, with child-side context loading as the fallback.
Each agent must execute its assigned role directly without dispatching another implementation or check agent.

## Handoff and limits

The main session owns planning, shared-file coordination, specification updates, and final acceptance.
The implementer owns the coherent product change and its focused checks.
The checker owns independent verification and necessary in-scope corrections.
Planning approval and task activation are complete. Follow the remaining implementation and check steps below.

The model-control callback and compound account-write contracts are already recorded in the relevant specs.
Keep that callback and the key-only verdict merge when only the verdict changes.
The confirmed-state correction also requires the affected inspection and persistence specs to match its verified behavior.

The added delivery scope produces `dist/build/sub2api-codex-eval-installer` alongside the raw server build output.
It replaces the archive/Python installation entry and updates the build, release workflow, README, and affected delivery specs.
The installer execution scope below records completed delivery work. Repeat only checks affected by the active confirmed-state correction.
Actual production installation, publishing a release, and data migration remain outside this task.
Use the current quality guidelines for work-commit approval and Trellis-managed archival and session records.

## Go installer execution after plan approval

1. Resume this task after approval of the expanded plan. Do not create another task or child task.
2. Assign the Go installer, build integration, focused checks, and delivery references to one `trellis-implement` owner.
   Preserve the staged model-control work and unrelated degraded-tag edits. Keep tests outside production source.
3. Port `scripts/install.py` into `scripts/installer/` with standard-library Go code and the parameter-free installation contract.
   Remove installation argument parsing and overrides. Keep numbered container selection, embedded payload, private `/tmp` backups, and a fixed 120-second deadline.
   Do not inspect process arguments or reject them at entry.
   First show case-insensitive `sub2api` name matches. Append a numbered `More` option when other running containers exist.
   Show the full list directly when no name matches. Keep explicit selection even for one match.
   Remove the Cancel option and `q` shortcut. Use numeric selection for `More`, with Ctrl-C for interruption.
   Add embedded payload cleanup and keep required license texts in the installer. Do not add another command interface.
4. Replace archive packaging in `scripts/build.py` with temporary-module assembly and static Linux/amd64 installer compilation.
   Keep the developer build entry and raw server output. Change release upload and publication to the single installer asset.
5. Replace the Python-import installer check with focused Go checks under `tests/installer/`.
   Cover parameter-free entry, numbered container selection, version and ELF rejection, candidate cleanup, and selected immutable container IDs.
   Extend the same selection cases for name filtering, numbered `More`, and no-match display. Check selection against the currently displayed list.
   Do not retain tests for removed argument parsing, external candidates, custom backup directories, or configurable timeouts.
   Use one fake-Docker transaction harness with local HTTP responses for success, timeout, restart changes, and failure or interrupt restoration.
   Include the stopped-container restoration path in that harness. Do not create another Docker fixture or add a runtime test flag.
   Keep timeout control inside tests. The public readiness deadline remains fixed at 120 seconds.
6. Run `gofmt`, the focused installer checks, and the existing build entry to produce the final installer.
   Reuse completed server and UI checks where their inputs remain unchanged. Do not add race, coverage, or full-suite gates.
   Inspect the installer format and release asset path. Check embedded payload extraction through the same focused harness.
7. Remove the replaced Python installer and obsolete archive references.
   Update README, AGENTS.md, and affected directory, delivery, and quality references to the verified Go entry and check command.
   Preserve unrelated changes in shared specifications.
8. Assign the completed delivery candidate and valid evidence to one `trellis-check` owner.
   Resolve in-scope findings, record the combined task result, then follow the existing commit-approval and finish workflow.

Do not run the installer against `sub2api-codex-eval-test` because `/app/sub2api` is its dependency launcher.
Task-reference checks passed during planning. The implementation requires the focused installer checks and final delivery build below.

The build's `--check` path runs `go test .` inside the temporary installer module, using the host platform.
Run those focused checks once for the final installer source. Keep the Linux/amd64 release build environment separate.
The local delivery build uses the existing command:

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check
```

CI retains `python3 scripts/build.py --source upstream --check` with the new installer checks integrated.
Do not repeat the unchanged local server/UI checks merely to exercise installer packaging.
Update the paths in `licenses/THIRD_PARTY_NOTICES.txt` alongside the delivery documentation.

Installer checks use one test set for the changed Go entry and transaction. Do not add a second acceptance suite.
Do not add race, coverage, full-suite, provider-quota, architecture-matrix, hash, or repeated browser gates.
Record the focused checks and final output paths in the implementation handoff.
