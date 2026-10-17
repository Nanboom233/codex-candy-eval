# Source Findings

## Account Refresh

- `patches/0006-evaluation-ui.patch:228`: `load` runs for foreground refresh and background polling.
- `patches/0006-evaluation-ui.patch:232`: each load clears `accountListReady` before rendering.
- `patches/0006-evaluation-ui.patch:744`: `refreshCatalog` also clears readiness before the account request.
- `patches/0006-evaluation-ui.patch:783`: `manualCatalogReady` controls account admission.
- `patches/0006-evaluation-ui.patch:357`: `pageBlocker` changes each row action while readiness is false.
- `patches/0006-evaluation-ui.patch:266`: polling repeats every 2.5 seconds during work and every 20 seconds otherwise.

The delayed inventory probe confirmed a control flash: the test button changes from enabled to disabled, then back after success.
The row, selection, and expanded history remain intact. See `refresh-error-contract.md` for the command and observations.

## Inspection Counts

- `patches/0005-inspection.patch:386`: each account starts with one test.
- `patches/0005-inspection.patch:399`: a wrong normal account enters downgrade confirmation.
- `patches/0005-inspection.patch:402`: a correct degraded account enters recovery confirmation.
- `patches/0005-inspection.patch:421`: confirmation starts another Candy task with the maximum minus one.
- `tests/backend/internal/codexeval/inspection_test.go:92`: current tests require three calls for a transition and one for a correct normal account.

The user clarified the reported defect on 2026-10-18: a degraded account that still answers wrongly must receive only one test.
The existing test includes a correct degraded account, but it lacks a degraded account that remains wrong.
The stable degraded branch returned one call in the isolated automatic-round probe.
The same probe returned three calls when `List` reported normal priority and `Get` already reported degraded priority.
Confirmation selection uses the stale list priority despite the current-account read.
See `inspection-count-contract.md` for the executed cases, repair boundary, and production-evidence limit.
Keep the other three branches and manual fixed-count tests unchanged.

## Settings Switch

The user clarified on 2026-10-18 that the requested switch belongs to each account.
Add a Candy list column. Keep the existing global Settings switch unchanged.

- `patches/0006-evaluation-ui.patch:967`: the current account switch controls native `schedulable`, not inspection participation.
- `patches/0006-evaluation-ui.patch:377`: the existing account toggle writes `/accounts/:id/schedulable`.
- `patches/0005-inspection.patch:275`: `inspectionEligible` owns shared inspection account admission.
- `patches/0005-inspection.patch:786`: `inspectionAccount` projects native account state into inspection selection.
- Original `backend/internal/service/admin_service.go:94`: `UpdateAccountExtra` supports key-level account metadata updates.
- Original `backend/internal/service/admin_account.go:928`: the service delegates sanitized metadata updates to the repository.
- Original `backend/internal/repository/account_repo.go:2733`: `UpdateExtra` merges JSONB keys without replacing unrelated account metadata.
- Original `backend/internal/service/admin_account.go:668`: ordinary full `extra` updates replace that object after managed-key preservation.
- Original `backend/internal/handler/dto/mappers.go:257` and `:468`: full and lightweight account responses carry `extra`.

Use native account `extra` for participation. Add a focused administrator endpoint that merges only the participation key.
Preserve the owned key through ordinary account edits and retain existing schedulable behavior.
An absent participation key means enabled, preserving existing selection. The global timer remains separate.
The lightweight DTO carries `extra`, and current redaction removes only named unrelated keys.
See `account-participation-contract.md` for the concrete account endpoint, metadata merge, full-edit preservation, and worker admission gap.

## Redundant Error Summary

- `patches/0006-evaluation-ui.patch:1113`: the Candy account header includes `inspection-runtime-error` and `candy-settings-status`.
- `patches/0006-evaluation-ui.patch:1294`: `inspectionMessage` writes settings messages in both Settings and the Candy header.
- `patches/0006-evaluation-ui.patch:1352`: state rendering combines runtime, account, and notification failures into the header error container.
- `patches/0006-evaluation-ui.patch:1378`: a state-read failure also writes into the same header error container.
- `tests/ui/sub2api-dropdowns.cjs:266`: existing assertions require the runtime error container in the header.
- `tests/ui/sub2api-dropdowns.cjs:493`: existing assertions require repeated account and notification errors in that container.

The user requested removal of the redundant summary below `账户` and `仅启用` on 2026-10-18.
Update these references together. Keep original account result errors. Remove runtime and notification prompts instead of moving them into Settings.
The user explicitly requires retaining the `下次巡检：` label and its scheduled time.
The approved task requirement must supersede older header-location assertions and the corresponding spec wording during implementation.

## Validation Resources

The original source directory `/tmp/sub2api-codex-eval-v0217-check` exists.
Its installed `frontend/node_modules` directory is available.
The build entry copies the original tree, applies all seven patches, and attaches test files only for checks.
Its fixed outputs are the Linux/amd64 binary and archive under `dist/build/`.
The retained fixture is `sub2api-codex-eval-test`. Do not create another container.

This research pass ran isolated DOM and Go characterization probes on unchanged production patches.
It did not run a product build or actual browser acceptance check.
