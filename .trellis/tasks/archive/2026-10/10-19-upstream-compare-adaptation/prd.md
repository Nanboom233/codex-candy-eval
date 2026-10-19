# Adapt Latest Upstream Evaluation Changes

## Goal

Add the upstream Pelican variants to the native Sub2API evaluation page.
Users can run either drawing task and identify the selected version in saved results.

## Background

The user requested a new task and adaptation of this comparison:
https://github.com/Nanboom233/codex-candy-eval/compare/main...haowang02%3Acpa-plugin-codex-candy-eval%3Amain

The inspected local base is `619d144`. The shared upstream ancestor is `7a490680a2e75660250faa235ebeb521d37d8564`.
The comparison adds `e120868f1f5bb4cb10c175fcab9eec7738c4908d` and release commit `4f8367642cc1e629c4ee15a975e5e0fc53c4ac16`.
These commits add Pelican variants and mark old Pelican records. The previous synchronization explicitly deferred these changes.
The current product targets Sub2API v0.2.17. It uses seven production patches and separate tests.

## Requirements

- R1: Add a version control after Pelican reasoning effort, with `original` and `all-in-one` options.
- R2: Use the exact upstream prompts. Update the visible and copied question when the version changes.
- R3: Default to the original version and restore a valid browser preference. Save changes immediately without overwriting model or effort preferences.
- R4: Send the selected version for single-account and batch runs. Store it on successful, failed, and cancelled results.
- R5: Keep requests without `variant` usable as `original`. Reject unknown nonempty variants before account resolution or supplier requests.
- R6: Prefer the selected markup language and retain the upstream tolerant extraction rules. Wrap all-in-one output for the existing preview and detail frame.
- R7: Preserve HTML extraction behavior, bounded retries, usage totals, cancellation, and the 256 KiB code limit.
- R8: Show the saved version in latest-result, history, and detail metadata. Use picture labels where static or animated output can appear.
- R9: Upgrade native history to schema 2. Mark pre-variant Pelican records as `original` and preserve retained evaluation data.
- R10: Preserve schema-1 fingerprint history. Preserve the current unversioned retirement rule for obsolete fingerprint and ModelTrace data.
- R11: Keep native authentication, account selection, Settings, inspection, model control, and installer delivery behavior.
- R12: Update Chinese user instructions and affected English specifications. Keep new tests in the established focused checks.
- R13: Save the original history bytes beside the file before migration. If backup fails, report the error and block replacement and later saves.

## Acceptance Criteria

- AC1 (R1-R3): Both versions appear after effort. A version change updates the question, copy content, and saved preference.
  Version changes during asynchronous startup preserve the saved model and effort. Group changes preserve the selected version.
- AC2 (R4-R5): Single and batch starts send the selected version. Missing versions run the original prompt. Unknown versions make no supplier call.
- AC3 (R2, R6-R7): Both variants use their exact prompts, upstream extraction rules, code limit, and page preparation through the current retry owner.
- AC4 (R4, R8): Completed, failed, and cancelled records retain their version. Latest-result, history, and detail views show the saved label.
  Toolbar changes do not relabel or filter history. Both variants share the existing latest-20 retention limit.
- AC5 (R9-R10): Schema-1 migration preserves Candy, fingerprint, and Pelican data while adding only the Pelican version and schema marker.
- AC6 (R9-R10): Unversioned files follow the existing retirement rule. Invalid, null, and future-schema files retain overwrite protection.
- AC7 (R11-R12): The established build/check entry passes. The server and self-extracting installer build at their existing paths.
- AC8 (R1, R6, R8): Actual browser checks cover both versions, previews, details, and toolbar layout at normal and narrow widths.
- AC9 (R9-R10, R13): Migration saves an exact original-file copy before replacement. Backup failure preserves the file and blocks later writes.
  Migration-save failure retains the original file, recovery copy, and readable memory records. Schema-2 reload creates no migration copy.

## Out of Scope

- Changes to fingerprint attribution, Candy scoring, automatic inspection, transport, or installer transactions.
- Automated scoring of the all-in-one picture or its displayed Candy answer.
- Restoration of the retired CPA runtime or its migration framework.
- Remote publication, deployment, or a release.

## Planning Status

The source comparison and native owners have been inspected.
The design, execution plan, research notes, and context manifests accompany this PRD.
The user approved review items 1 through 4, including the original-file migration copy.
See [the detailed review](research/plan-review.md) for the decisions, source findings, and validation scope.
The user approved the final planning summary and requested a reread to remove redundant tests and added requirements.
The audit is recorded in [task audit](research/task-audit.md). It bounds validation and adds no product requirement.
The task is active. Implementation and browser evidence accompany the final review.
