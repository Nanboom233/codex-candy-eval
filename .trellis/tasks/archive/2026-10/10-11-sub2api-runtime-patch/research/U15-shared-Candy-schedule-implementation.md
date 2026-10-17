# U15 Shared Candy Account Selection

## Approved contract

The user clarified that only account selection is shared.
Inspection retains its independent model, effort, maximum count, and confirmation flow.
Manual Candy retains fixed runs for row, selected-account, and all-account tests.

The existing native inspection settings add four shared fields:

| Field | Default | Owner |
| --- | --- | --- |
| `group_id` | `0` until a group is selected | Current page group |
| `account_type` | `all` | Candy account-type filter |
| `account_plan` | `all` | Candy subscription filter, including `Pro ALL` |
| `enabled_only` | `false` | Candy enabled-account filter |

The same native `SettingService` and `SettingRepository` persist these fields.
The coordinator resolves current group membership before each round.
Availability and busy checks still apply.
Inspection keeps native `Account.IsSchedulable()` checks, including expiration and quota restrictions.
Shared account selection does not change manual administrator admission or remove native inspection safeguards.
The separate `skip_free` override is removed.
An unselected group cannot start a round or enable the timer.

The settings tab has no account list or save button.
It retains inspection model, effort, maximum count, wait bounds, priority fields, and notifications.
Its model suggestions use the current page group catalog.
The page does not fetch an independent all-account catalog.

Candy owns account progress, records, history, inspection outcomes, and inspection stop controls.
Inspection retains the first test and conditional confirmations up to its maximum count.
Recovery retains its 240-second wait.
Priority changes, notifications, and cancellation retain their existing owners.

## Source changes

`inspection_settings.go` adds the shared fields and removes `skip_free`.
The existing settings service persists the new fields without another repository or configuration store.
`codex_inspection_handler.go` queries the saved group and exposes current membership, type, and plan to the coordinator.
The coordinator reapplies shared filters before tests and priority updates.
It uses the selected group as result origin for accounts with multiple groups.
Subscription labels and Pro aliases match the existing Candy page definitions.

Independent review found that shadow accounts missed their native parent subscription.
Both the list and the single-account recheck now use the same inspection conversion.
This conversion reuses native `accountListResponseFromService` and `enrichShadowParents`.
An account's own plan takes precedence, as it does in Candy.
Otherwise, the native parent plan supplies the subscription filter value.
Each recheck resolves the current parent plan without copying credentials into the shadow.

The coordinator retains its pending wait for notification, model, or filter edits.
Changes to scheduling bounds or the enabled switch draw a new wait.
Each completed round still draws a new random wait.

`inspection.js` saves valid edits after a 400-millisecond delay.
Only one settings write can run at a time.
Field revisions preserve edits made during a read or save.
Successful saves clear submitted secrets from the form.
Blank secrets still retain their stored values through the native settings service.
Validation and repository failures remain visible and retain unsaved values.
A start action saves pending valid edits before it starts the round.

`app.js` loads native settings before initial account selection.
Saved filters stay selected when the current inventory has no matching accounts.
An absent saved group cannot silently select another group.
The other evaluation tabs retain their account, model, and run controls.

`catalog.js` no longer has an independent inspection catalog or model request path.
`candy.js` and `components.js` no longer have inspection-only read-only list adapters.
`ui.html` and `style.css` remove the settings account list and save button.
Inspection progress, outcomes, errors, and start/stop controls now use Candy.
Existing Candy records, history details, and answer controls remain in use.

The local fixture helper now creates its groups before saving settings.
It uses the actual U13 Group A ID and explicit shared filter defaults.
This change does not execute the fixture or modify a running container.

## Focused checks

Passed:

- `SUB2API_SOURCE=/tmp/sub2api-codex-eval-v0217-check node tests/ui/sub2api-dropdowns.cjs`.
- The initial Go command without `-tags unit` passed the inspection checks but excluded the tagged native lifetime test.
- After the shadow repair: `go test -tags unit -overlay /tmp/sub2api-U15-go-overlay.json ./internal/codexeval ./internal/service ./internal/handler/admin -run 'TestCodexInspection|TestCodexEvalBackgroundNativeStreamLifetime' -count=1`.
- `go vet -overlay /tmp/sub2api-U15-go-overlay.json ./internal/codexeval ./internal/service ./internal/handler/admin`.
- After the shadow repair: `go vet -tags unit -overlay /tmp/sub2api-U15-go-overlay.json ./internal/handler/admin`.
- `gofmt` on changed Go files.
- Python AST parsing for `tests/installer/native_inspection_fixture.py`.
- `git diff --check`.

Go checks used `/tmp/sub2api-U11-check.CqGH5I/backend` and current source through the overlay.
The UI check covers debounce, serialized writes, newer edits, secrets, validation, and save failures.
It also covers missing-group/filter preservation, current membership, inventory failures, manual entry points, and Candy history.
The existing inspection checks retain original conditional counts, recovery, priority updates, notifications, cancellation, and restart behavior.
The new shadow regression checks batch parent resolution, own-plan precedence, and current-plan resolution for the single-account recheck.
The UI regression checks the same shadow under `Pro ALL` and a specific plan.
The tagged native lifetime check now passed with the inspection and shadow checks.
It covers Responses and Chat conversion completion and explicit cancellation through the native sender.
The same UI command passed again after the shadow repair.

## Delivery state and limits

Independent source review passed after the shadow account repair.
The approved build and packaging commands completed without `--check` or repeated test suites.

Commands:

```sh
python3 scripts/build-sub2api.py --source /tmp/sub2api-codex-eval-v0217-check
python3 scripts/package-sub2api.py
```

The build used existing dependencies and Go caches.
The native frontend build completed, including its required locale and type checks.
Vite completed in 11.71 seconds.
The build log is `/tmp/sub2api-U15-build.log`.

| Output | Result |
| --- | --- |
| `dist/build/sub2api-codex-eval` | 125599904 bytes. ELF 64-bit x86-64, statically linked and stripped. |
| `dist/build/sub2api-codex-eval.tar.gz` | 38701721 bytes. Seven expected entries. |

Go build metadata confirms `GOOS=linux`, `GOARCH=amd64`, and `CGO_ENABLED=0`.
The logged build command sets `main.Version=0.2.17` and `main.BuildType=release`.
The main session will verify the runtime version in the retained Linux container.
The archive contains the executable, `install.py`, and five license texts.
It contains no development or test files.

Docker, browser, production, and paid-provider checks have not run in this implementation turn.
The main session owns subsequent runtime, browser, and timer acceptance checks.
The earlier history `context canceled` report remains unresolved without its triggering operation.
No cancellation result was changed into success.
