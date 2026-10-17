# Browser Acceptance

## Candidate and Method

The retained `sub2api-codex-eval-test` container runs the packaged v0.2.17 candidate at `http://127.0.0.1:60198`.
The existing fixture uses synthetic accounts and a local supplier.
Playwright used the fixture's private storage-state file without printing its contents.

## Account Controls and Refresh

The Candy list shows separate `启用` and `启用自动巡检` columns.
Turning off inspection for `u13-a-only` persisted after page reload.
The account's native enabled switch and the second account's inspection switch remained on.
Manual test controls remained usable for the opted-out account.
Ordinary background refresh retained the visible rows and controls.
The existing DOM harness owns delayed-response and error-injection assertions.

## Error Removal and Next Run

The removed header elements `inspection-runtime-error` and `candy-settings-status` are absent.
The header has no empty error row.
Temporarily enabling the global schedule showed `下次巡检：10-18 04:20` in the Candy header.
The timestamp remained visible at a 390-pixel viewport without horizontal page overflow.
The fixture's global schedule was then disabled and the account switch restored.

## Alignment Correction

The first narrow screenshot showed different vertical positions for the two row switches.
The user identified this defect before final acceptance.
The implementation owner changed only Candy's narrow-screen control alignment in the existing CSS.
The existing UI harness passed. The standard build without `--check` refreshed the executable and archive.
Unchanged backend checks and vet results remain valid from `implementation-build.log`.

The fixture then loaded the corrected package.
At a 390-pixel viewport, both visible account rows have identical native and inspection switch top positions:

| Row | Native switch top | Inspection switch top | Offset |
| --- | --- | --- | --- |
| `u13-a-only` | 760.5703125 px | 760.5703125 px | 0 px |
| `u13-shared` | 979.6015625 px | 979.6015625 px | 0 px |

The corrected screenshot confirms aligned controls despite different history lengths.
Screenshots are `/tmp/inspection-controls-wide.png` and `/tmp/inspection-controls-narrow-aligned.png`.
The final candidate retains the existing desktop grid. The correction applies only at narrow Candy widths.

## Result

The bounded browser acceptance is complete. No provider quota or production application was used.
Backend sender counts and synthetic failure paths use the already passed focused tests.
