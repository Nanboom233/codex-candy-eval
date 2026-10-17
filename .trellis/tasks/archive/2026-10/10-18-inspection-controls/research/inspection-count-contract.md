# Research: Inspection Count Contract

- Query: Trace the effective patched automatic inspection path for a degraded account whose first answer is wrong. Identify count, priority, notification, settings, and reread behavior.
- Scope: internal
- Date: 2026-10-18

## Findings

### Static degraded-wrong branch already stops after one test

`patches/0005-inspection.patch` adds `backend/internal/codexeval/inspection.go`.
The coordinator reads accounts, filters eligibility and busy state, creates one
result, then calls the first phase with `runs=1` (`patches/0005-inspection.patch:359-386`).
It adds a wrong normal-priority account to downgrade confirmations and a correct
degraded-priority account to recovery confirmations. It does not add a wrong
degraded-priority account to either group (`patches/0005-inspection.patch:390-404`).
The confirmation loop only tests members of those groups and runs at most
`MaxTestsPerAccount - 1` more tests (`patches/0005-inspection.patch:407-429`).

Therefore, when the account-list snapshot has `Priority == DegradedPriority`,
and its first Candy result is a valid wrong answer, the coordinator makes one
inspection test in that round. The result remains `unchanged`; no target priority
is set (`patches/0005-inspection.patch:395-404`). Priority updates require a
`recommended` result and enabled adjustment (`patches/0005-inspection.patch:447-475`).
Notifications are sent only for successfully updated accounts (`patches/0005-inspection.patch:476-497`).
The degraded-wrong branch therefore keeps priority and sends no priority-change
notification, with adjustment either on or off.

### The fresh priority read is not used for classification

The account hook maps the native account's current `Priority` field into
`InspectionAccount` (`patches/0005-inspection.patch:786-790`). The list hook
uses that projection to make the round's `eligible` snapshot
(`patches/0005-inspection.patch:811-829`).

Before each first or confirmation phase task, `testPhase` calls `hooks.Get`.
It checks whether the returned account is still inspection-eligible, but it
does not compare or copy its current priority (`patches/0005-inspection.patch:515-535`).
After the first phase, group selection still uses the original `a.Priority`
from the list snapshot (`patches/0005-inspection.patch:392-404`). A concurrent
change from the configured normal priority to the configured degraded priority
after `List` can therefore leave the stale account in the normal-wrong group.
The coordinator can then run confirmations for an account that is already
degraded at the time of the first test. The later priority reread occurs only
after confirmation and prevents a stale-priority write; it cannot prevent those
extra tests (`patches/0005-inspection.patch:447-465`).

An isolated dynamic characterization confirmed this path on unchanged patches.
In `/tmp/inspection-controls-research-20261018`, the focused test
`TestResearchInspectionDegradedWrongCalls` passed with these observations:

- Stable `List=99`, `Get=99`, wrong answer: one sender call and `Tests == 1`
  with adjustment both on and off. Outcome was `unchanged`, with no target.
- Changed `List=10`, `Get=99`, wrong answer: three sender calls and
  `Tests == 3` with adjustment both on and off. With adjustment off, outcome
  was `recommended` for priority 99. With adjustment on, the final update
  reread rejected the stale priority with `账户状态或优先级已改变`. Both modes
  made zero priority updates and sent zero notifications.

Command, run from that isolated source tree:

```sh
GOWORK=off GOFLAGS=-mod=readonly go test ./internal/codexeval -run '^TestResearchInspectionDegradedWrongCalls$' -v -count=1
```

This confirms the code path on the unchanged patch candidate. It does not
establish that the user's production run had this exact timing or state change.
A stable degraded account in a fresh list does not take this path.

Preferred minimal fix: after the first test phase and before building
confirmation groups, reread each account's current priority and classify the
answer using that value, not the list snapshot. Keep the original priority for
the result record. Carry the fresh classified priority as the baseline for the
existing pre-update comparison, and verify it again before any confirmation
task. If the reread fails or the account is no longer eligible, do not launch
confirmations. This prevents stale-list priority from spending extra tests and
does not alter the Candy test count or sender path.

### Priority classification uses exact configured values

Default priorities are normal `10` and degraded `99`. Validation requires
normal to be positive and lower than degraded
(`patches/0005-inspection.patch:717-720,729-757`). The coordinator classifies
only exact equality with those configured values. An account at another
priority still gets the first test, but it enters neither confirmation group
(`patches/0005-inspection.patch:399-404`).

The account row's current priority and the configured values are distinct data
sources. A custom persisted priority pair is read by `SettingService` through
the native `SettingRepository.GetMultiple` path
(`patches/0005-inspection.patch:911-937`). Startup passes this method as the
coordinator's `Load` hook (`patches/0005-inspection.patch:808-810`); coordinator
initialization loads it once (`patches/0005-inspection.patch:82-87,134-149`).
The administrator GET handler reads the repository directly. The PUT handler
binds and validates settings, calls `UpdateCodexInspectionSettings`, and returns
the public values (`patches/0005-inspection.patch:846-869`). A successful
service update writes settings and calls `UpdateInspectionSettings`
(`patches/0005-inspection.patch:940-963`).

Each running round captures `i.settings` at start
(`patches/0005-inspection.patch:230-257`). Account priority classification
continues to use that snapshot for the whole round. A settings update during a
round changes future coordinator settings; it does not replace the active
round's local settings value.

### Candy execution does not retry a wrong answer

Each Candy task loops exactly `runs` times and records one `candyResult` for
each completed evaluation (`patches/0001-evaluation-core.patch:49-71`). Candy
evaluation calls `askCodex`, which calls `executeModel` once
(`patches/0001-evaluation-core.patch:73-85,284-291`). The retry helper is
`executeProbe`, documented for collection probes; its retryable statuses do not
retry invalid answers, and Candy does not call it
(`patches/0001-evaluation-core.patch:991-1046`). The `start` entry clamps Candy
run count to `1..candyMaxRuns`, while inspection passes `Runs: runs`
(`patches/0001-evaluation-core.patch:1972-1990`; `patches/0005-inspection.patch:532-535`).

`InspectionResult.Tests` increases by the number of Candy result rows matching
that task's run ID (`patches/0005-inspection.patch:593-620`). Thus for a
successful first-phase task, one coordinator test means one Candy result and
one direct sender invocation. A transport error also ends that Candy
evaluation; it does not retry the Candy request.

### Existing safeguards and test gap

One coordinator instance rejects an overlapping round under its mutex
(`patches/0005-inspection.patch:230-256`). The process-wide account `active`
map rejects a second manual or inspection Candy task for the same account
(`patches/0001-evaluation-core.patch:2031-2046`). Automatic and immediate
inspection both call this coordinator. Manual Candy routes use the public
`Start` wrapper, which sets `inspection=false`; preserve that separate request
path (`patches/0001-evaluation-core.patch:1972-1976`). A completed round creates
a new scheduled wait, so each later round is a separate inspection
(`patches/0005-inspection.patch:333-357`).

The current coordinator test counts three tests for a wrong normal account
and a correct degraded account, and one for a correct normal account. It does
not include a degraded account with a wrong first answer
(`tests/backend/internal/codexeval/inspection_test.go:25-33,92-99`). The
second run checks that adjustment-off does not update priorities or notify,
but does not assert its test counts (`tests/backend/internal/codexeval/inspection_test.go:111-123`).
The settings tests check loaded and saved priority values, but do not connect a
custom pair to account classification (`tests/backend/internal/service/codex_inspection_settings_test.go:77-161`).

## Regression Handoff

The check table in `../implement.md` defines the required cases.
Extend existing coordinator tests with stable degraded-wrong and stale-list/degraded-current assertions, with adjustment on and off.
Assert one sender call/result, unchanged priority, and no notification. Reuse existing confirmation and manual-count coverage.
Custom persisted priorities and repeated-round tests are not additional requirements for this repair.

## Files Found

- `patches/0005-inspection.patch` — inspection coordinator, native account hooks,
  settings handler, and settings service implementation.
- `patches/0001-evaluation-core.patch` — Candy run loop, evaluation entry,
  direct sender path, and native task admission.
- `patches/0007-sub2api-integration.patch` — native router initialization and
  administrator route registration.
- `tests/backend/internal/codexeval/inspection_test.go` — coordinator tests;
  lacks the degraded-wrong assertion.
- `tests/backend/internal/service/codex_inspection_settings_test.go` — settings
  persistence and priority validation tests.
- `.trellis/spec/backend/native-inspection.md` — current test count,
  confirmation, priority, and native settings contract.
- `.trellis/spec/backend/error-handling.md` — distinguishes wrong answers from
  request failures and documents separate retry behavior.

## Related Specs

- `.trellis/spec/backend/native-inspection.md`: sections 3, 4, and 6 define
  account tests, priority changes, and the existing focused checks.
- `.trellis/spec/backend/error-handling.md`: `executeProbe` retries transport
  failures for collection probes, not incorrect answers.
- `.trellis/spec/backend/quality-guidelines.md`: use the focused backend tests
  and existing patched source tree.

## External References

None. The behavior is defined by this repository's patch series and native
Sub2API handlers.

## Caveats / Not Found

- The dynamic probe used a temporary isolated copy and synthetic hook values.
  It characterizes unchanged code but does not prove the exact production
  account/config sequence behind the user's report.
- No stable-list code path in the inspection coordinator produces a second
  test for a wrong degraded account. The stale-priority mismatch is a confirmed
  code path and a plausible explanation when account priority changes after
  listing.
- Priority values outside the current configured pair are not treated as
  normal or degraded for transitions. They still receive one first test.
- This report traces the patch source. The original upstream tree remains
  unchanged; the build applies patches in `patches/series` order.
