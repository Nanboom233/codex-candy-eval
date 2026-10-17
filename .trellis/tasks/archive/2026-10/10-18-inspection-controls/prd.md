# Account Refresh and Automatic Inspection Controls

## Goal

Keep the account page stable during refresh. Test a degraded account once when its first answer remains wrong.
Let administrators control each account's inspection participation from a dedicated Candy page column.
Remove inspection runtime and notification error prompts. Keep the `下次巡检：` timestamp.

## Requirements

- R1: Keep account rows, selections, expanded history, and action controls stable during successful background refresh.
  Continue account membership and result updates. Block starts after an actual inventory failure or group change.
  Evidence: `patches/0006-evaluation-ui.patch:232`, `:744`, and `:358` reset readiness and change row controls during polling.
  A delayed DOM probe confirmed the unchanged row's test button becomes disabled during a successful refresh.
- R2: Test each eligible account once at the start of an inspection round.
  If a degraded account still answers incorrectly, end its work for that round after exactly one test.
  Keep its priority unchanged. Do not add confirmation tests or send a priority-change notification.
  A correct normal account also stops after one test.
  Only a wrong normal account or a correct degraded account can enter confirmation tests.
  The configured maximum includes the initial test. Keep the existing recovery wait before recovery confirmations.
  Errors and mixed outcomes must still prevent priority changes. Keep manual test counts independent.
  Evidence: `patches/0005-inspection.patch:386`, `:399`, `:402`, and `:421` own initial tests and confirmation selection.
  Evidence: `tests/backend/internal/codexeval/inspection_test.go:25` lacks an explicit degraded account that remains wrong.
  The user identified the degraded-wrong branch as the reported count defect on 2026-10-18.
  A synthetic automatic-round probe confirmed extra tests when the list priority is normal but the current priority is already degraded.
- R3: Add an `启用自动巡检` column to the Candy account list, with one independent switch for each account.
  Save participation by native account ID. The same account retains its value across groups, reloads, and server restarts.
  Turning the switch off excludes that account from inspection. It does not disable the account or change normal traffic scheduling.
  Turning it on permits inspection when the global schedule and existing group, availability, and filter rules allow it.
  Keep manual Candy tests usable regardless of this participation switch.
  Preserve the existing global inspection switch and immediate start/stop controls in Settings.
  Keep existing accounts included when no participation value exists, preserving current behavior.
  A request already in flight can finish. Recheck participation before later inspection tests and priority updates.
  Intentional opt-out ends remaining inspection work without a test-error message. Preserve completed test records.
  Show account switch save errors beside the affected control. Preserve other account fields.
  Align the account and inspection switches on the same baseline in narrow rows.
  The user clarified that this switch is account-level on 2026-10-18.
- R4: Remove the redundant automatic inspection error summary below `账户` and `仅启用`.
  Keep account test errors in their existing result records. Do not repeat them in a header-wide error list.
  Remove inspection runtime and notification error prompts from the page. Do not move them to Settings or another page area.
  Keep the `下次巡检：` label and its scheduled time in the Candy account header.
  Show switch validation and save failures beside the affected control.
  Evidence: `patches/0006-evaluation-ui.patch:1113` creates the header error container.
  Evidence: `patches/0006-evaluation-ui.patch:1352` combines runtime, account, and notification errors in that container.

## Constraints

- Target Sub2API v0.2.17 on Linux/amd64.
- Keep production changes in the existing seven patches.
- Keep tests outside production patches.
- Preserve the original upstream tree, private configuration, and real application data.
- Reuse the named test container and installed frontend dependencies.

## Acceptance Criteria

- AC1: A delayed successful background refresh preserves row nodes, selection, expansion, and unchanged button state.
- AC2: Membership changes remove stale accounts. Failed inventory reads block starts until a successful refresh.
- AC3: A degraded account with a wrong first answer makes exactly one test call and creates one result per inspection round.
  It retains its priority and sends no priority-change notification. Verify automatic priority adjustment both on and off.
  Include a normal list snapshot followed by a degraded current-priority read before the first test.
- AC4: A correct normal account stops after one test. Transition candidates retain the configured confirmation limit and recovery wait.
  Failed, mixed, and cancelled runs do not produce false confirmations.
- AC5: Manual fixed-count tests retain their current count and entry points.
- AC6: Each Candy account row has an independent inspection switch. Changing account A leaves account B and the global switch unchanged.
  Participation persists across page reload, group changes, and server restart. Absent values retain existing account participation.
  Both row switches align at narrow widths, including rows with different history lengths.
- AC7: An opted-out account receives no new inspection test or priority update after participation is rechecked.
  Its account enabled state, traffic scheduling, and manual Candy tests remain unchanged.
  Failed saves retain the last saved state, show the account error, and permit retry without losing other account metadata.
  If a multi-test worker is active, its next request checks participation again. Opt-out does not create an evaluation failure.
- AC8: The existing build/check entry and one browser pass verify the effective patched candidate.
- AC9: Account test failures do not produce a repeated error summary below the Candy account header.
  Their original result errors remain available. Neither Candy nor Settings displays inspection runtime or notification error prompts.
  The `下次巡检：` label and scheduled time remain visible when a next run exists.
  Switch save failures remain visible beside the control. The header leaves no empty error row after removal.

## Validation Scope

Use `implement.md` as the check plan. Add assertions only for changed behavior without existing coverage.
Reuse existing confirmation, cancellation, manual-count, and failure checks for AC2, AC4, and AC5.
Native database storage provides persistence across groups and restarts. Do not add separate restart, copy, import/export, or DTO test matrices.
Run the established build/check entry once for the final candidate. Use one browser pass for layout and switch interaction.
Reuse valid results for unchanged inputs. Repeat a check only after a relevant change or failure.

## Out of Scope

- Change inspection selection beyond account participation, or change priority, wait, notification, or manual evaluation rules.
- Redesign task progress or add cumulative progress fields as part of the reported test-count repair.
- Deploy to a real application or use provider quota.
- Change upstream account management pages without evidence of an affected owner.

## Planning Status

Code-level research and isolated characterization probes are complete. The design and execution plan include their findings.
The user approved this plan on 2026-10-18. Implementation, bounded checks, independent review, and browser acceptance are complete.
See `research/implementation-evidence.md`, `research/check-evidence.md`, and `research/browser-evidence.md` for evidence.
The corrected executable and archive are available in `dist/build/`.
The user approved the exact work commit. Signed commit `75f5a3f` contains the verified work, and its signature passed verification.
