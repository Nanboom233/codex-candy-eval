# Separate inspection verdicts from account priority

## Goal

Persist inspection verdicts independently from account priority. A confirmed degraded account needs one test when another round returns a wrong answer.
The automatic priority adjustment switch controls actual priority writes. It must not disable degradation detection or recovery.

## Background

- `patches/0005-inspection.patch:411-415` selects confirmation groups from actual priority.
- `patches/0005-inspection.patch:461` skips priority writes when `adjust_priority` is false.
- `patches/0005-inspection.patch:492-509` sends notifications only after successful priority writes.
- `tests/backend/internal/codexeval/inspection_test.go:314-361` covers adjustment off only with an already degraded actual priority.
  It misses an initially normal account that becomes degraded while actual priority stays unchanged.

The user approved task creation and the independent-verdict direction on 2026-10-18.
The code-level findings are in `research/verdict-state-contract.md`. The design and execution plan define the proposed repair.

## Requirements

- R1: Save each account's confirmed inspection verdict independently from its actual priority.
  The automatic priority adjustment switch must not disable verdict updates or determine the next round's test count.
  An account without a saved verdict starts as normal, regardless of its actual priority.
- R2: A confirmed degraded account that still answers incorrectly receives exactly one test in each later inspection round.
  A confirmed normal account that answers correctly also stops after one test.
  Possible degradation or recovery uses the existing confirmation count and recovery wait.
- R3: Apply the saved verdict to both scheduled and immediate inspection rounds.
  Keep manual Candy tests independent. Preserve account participation and the existing per-send opt-out checks.
- R4: Check the automatic priority adjustment switch only when a confirmed verdict changes.
  If enabled, apply the existing guarded priority update. If disabled, save the verdict without changing actual priority.
  A stable verdict does not trigger a priority update, including after the switch becomes enabled.
- R5: Use existing native storage and patch owners. Do not add a table, alternate scheduler, or separate status file.
  Preserve the previously removed error prompts and retained `下次巡检：` label/time.

## Acceptance Criteria

- AC1 (R1–R4): Start with a normal account and automatic priority adjustment off. Confirm wrong answers and save a degraded verdict.
  Its next wrong round makes one sender call and creates one result, while actual priority remains unchanged.
- AC2 (R1–R3): A degraded verdict recovers through the existing wait and confirmation sequence with adjustment off.
  The next correct round makes one sender call and creates one result.
- AC3 (R1, R3, R4): Failed, mixed, cancelled, incomplete, or intentionally opted-out confirmation work cannot change the saved verdict.
  Reuse existing confirmation, participation, and cancellation checks.
- AC4 (R3, R5): Native metadata saves preserve unrelated keys. Full account edits preserve the verdict key.
  Manual entry points and fixed counts remain usable. The page remains unchanged.
- AC5 (R4): A confirmed verdict change saves state in both adjustment modes. Only the enabled mode can write actual priority.
  Enabling the switch during a stable verdict does not trigger a priority write.
- AC6 (R1, R4): A failed verdict save prevents a priority write. A failed priority write retains the confirmed, saved verdict.
- AC7 (R1, R5): An account without a saved verdict starts as normal, even at the configured degraded or a custom priority.
  No priority value initializes a degraded verdict or exempts an account from confirmation.
- AC8: Run the established build/check entry for the final candidate. Reuse unchanged checks and add only missing regression assertions.

## Constraints

- Target Sub2API v0.2.17 Linux/amd64 with the existing seven production patches and tests outside those patches.
- Preserve the original upstream source, private configuration, real application data, and unrelated work.
- Use only the fixed `sub2api-codex-eval-test` fixture if runtime verification needs a container. Retain it for later use.
- The pending container-retention spec edit belongs to the preceding user instruction. Preserve it without adding product scope.
- Do not add race, coverage, full-suite, browser-engine, copy, or installer-restoration matrices.

## Out of Scope

- Redesign Candy controls, add status columns, or restore runtime/notification error prompts.
- Change manual test counts, provider sending, group filters, account participation, or global scheduling.
- Add account-copy rules, new public DTO fields, or a verdict reset interface.
- Infer verdicts from old priorities, migrate old priority-based state, or preserve its classification rules.
- Synchronize mismatched priorities during stable rounds, retry failed priority writes automatically, or change notification rules and text.
- Deploy to a production application or use provider quota.

## Artifact Status

The user narrowed the plan to independent verdicts and transition-only priority writes.
The user accepts the breaking change and removed compatibility with priority-based state.
The user approved implementation after the task scope and check plan audit.
