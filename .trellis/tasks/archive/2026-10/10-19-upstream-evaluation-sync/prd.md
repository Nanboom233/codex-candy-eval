# Sync upstream evaluation changes

## Goal

Give native Sub2API users the upstream fingerprint detector, retry behavior, and current request context.
Keep the existing installation and account inspection experience.

## Background

The user requested a new task to analyze and synchronize upstream changes.
Local `main` is `2aedffe`. The shared base is `6b46d4f`.
The comparison found 21 local commits and 11 upstream commits after that base.
Both `upstream/main` and the fork's `main` pointed to `7a49068` during research.
The upstream source is `haowang02/cpa-plugin-codex-candy-eval`.
The user's fork is `Nanboom233/codex-candy-eval`.
Local product source lives in seven ordered patches. Upstream product source uses `internal/plugin/`.
The delivery target remains Sub2API v0.2.17 on Linux/amd64.

## Technical Evidence

- The detector has 36 challenges and 60 enrolled candidates.
  It adds no Go module dependency. Its resources total about 4.55 MB and decode into about 8.2 MiB of floats.
- Native account IDs and history fields differ from the CPA state format.
  The CPA migration cannot preserve the native history contract.
- Native effort controls already omit the tooltips removed upstream.
- A read-only merge-tree probe reports 43 conflict entries, mostly in CPA files removed by the native port.

Sources:

- [Evaluation changes](research/upstream-evaluations.md)
- [UI and saved history](research/upstream-ui-history.md)
- [Commit inventory and Git integration](research/git-integration.md)
- [Detector code contracts](research/code-detector.md)
- [Retry and inspection code contracts](research/code-retry-inspection.md)
- [Native UI and history code contracts](research/code-ui-history.md)
- [Patch, resource, and check ownership](research/code-delivery.md)

## Requirements

- R1: Account for each upstream commit after the shared base using actual source changes and tests.
  Record which behavior already exists locally and which behavior needs adaptation.
- R2: Plan the integration through the current production patches, native module, and existing test owners.
- R3: Preserve native administrator access, account selection, transport, cancellation, and supplier-quota behavior.
  Preserve existing Candy and Pelican history data and formats.
- R4: Preserve confirmed-state model control, priority adjustment, recovery inspection, and degraded account tags.
  Keep automatic inspection on Candy. Upstream's multi-test inspection was already outside the native port at the shared base.
- R5: Preserve the parameter-free Go installer, numbered container selection, backups, and failure recovery.
- R6: Account for changes to test workflows and history without applying the CPA state migration to native data.
- R7: Keep the Chinese user README and practical build instructions current with approved behavior.
- R8: Create an integration branch from local `main` and normally merge `upstream/main`.
  Resolve conflicts through native source ownership. Retain both histories without rebasing or force-pushing.
  After approved commits and task finish, fast-forward local `main` and push normally to `origin/main`.
  Check remote changes before publication and preserve any new remote work.
- R9: After plan approval, implement the selected upstream behavior and verify the resulting native delivery.
  Complete ordinary work commits, Trellis finish, and the previously requested normal push to the user's fork.
- R10: Replace the old manual fingerprint test and ModelTrace with the upstream LM fingerprint detector.
  Use the three upstream test pages plus the native Settings tab. Remove the old sampling-mode controls.
  Keep the existing native layout and account controls. Add no explanatory small print or tooltips.
  Select three challenges per run. Use upstream ranking, aliases, concurrency, and partial-result behavior.
  After the sequence check passes, three valid answers produce probabilities. One or two produce rankings without probabilities.
  A failed sequence check produces `not_model` without rankings or probabilities.
- R11: Do not migrate old fingerprint or ModelTrace data. This is an accepted breaking change.
  Start the new fingerprint history empty. Keep Candy and Pelican history data and formats unchanged.
  Detect an unversioned native file and immediately save its Candy and Pelican records in the new versioned snapshot.
  Do not add backup, conversion, or legacy-history compatibility flows for the retired tests.
- R12: Apply the upstream retry behavior, with at most three total attempts per logical request.
  Wait two seconds before the first retry and four before the second.
  Retry transient request failures and unusable answers. Do not retry valid wrong Candy text, HTTP 401, or HTTP 429.
  Cancellation stops pending waits and further sends. Count attempts and total usage across attempts.
  Each logical Candy test produces one final result, regardless of its internal attempts.
  Each completed result counts once toward confirmation. Request errors do not count as confirmation answers.
  Preserve inspection participation checks before each supplier send, including retries.
  An opt-out during retry backoff stops later sends and leaves the account outcome unchanged without a synthetic failed result.
- R13: Adopt the upstream Codex request identity and context update through the native account sender.
  Change the default identity from `codex-tui/0.160.0` to `codex-tui/0.163.0` with the upstream platform fields.
  Use the updated context, `history_initialization: new`, and the fingerprint developer instructions in the corresponding context field.
  Preserve account-specific User-Agent overrides, authentication, and proxy rules.

## Acceptance Criteria

- AC1 (R1): A source-backed inventory accounts for every upstream commit and its local destination or exclusion reason.
- AC2 (R2, R3): The design identifies actual native owners and preserves current Sub2API integration contracts.
- AC3 (R4, R5): The plan explicitly retains local inspection, model control, priority, tag, and installer behavior.
- AC4 (R6): The plan explains changes to evaluations, user controls, saved history, dependencies, and resource costs.
- AC5 (R7): The plan includes affected Chinese user instructions and correct project attribution.
- AC6 (R8): The Git plan retains local and upstream commits and names the final publication branch.
- AC7: Verification follows changed behavior and reuses valid evidence. No race, coverage, or full-suite gate is added.
- AC8 (R9): The selected upstream behavior works through the native administrator page and account sender.
  The checked server and Go installer build successfully through the existing entry point.
- AC9 (R8, R9): Final local and fork history retain the local work and the upstream commits without a forced update.
- AC10 (R10): Manual evaluation offers Candy, fingerprint detection, and Pelican. No old fingerprint or ModelTrace page remains.
  The new detector follows upstream challenge, attribution, concurrency, and partial-result behavior.
  Settings-first entry still loads shared groups. Native model selectors and account controls remain usable without extra explanatory text.
- AC11 (R3, R11): Candy and Pelican retain existing records and formats.
  Loading an unversioned native file writes the new snapshot through the existing atomic saver.
  Old fingerprint and ModelTrace records do not enter the new detector history or a migration/backup flow.
- AC12 (R12): Retries follow the approved limit, delays, final statuses, and cancellation rules.
  Valid wrong Candy text ends after one attempt. Internal retries do not increase logical inspection or confirmation counts.
  New results show actual attempts and accumulated usage without rewriting old Candy or Pelican records.
- AC13 (R3, R13): Native evaluation requests use the approved default identity and context fields.
  Account-specific User-Agent overrides, authentication, and proxy selection retain their current behavior.

## Scope Limits

- Planning and read-only research come first. Product edits require approval of the final plan.
- Do not replace native Sub2API delivery with the original CPA plugin installation or authentication flow.
- Do not add compatibility platforms, new installer parameters, production account tests, or extra infrastructure.
- Do not publish a release or overwrite remote history as part of analysis.

## Decision Record

The user approved all five items through separate discussion.

| Item | Requirement |
| --- | --- |
| 1: Test replacement | R4, R10 |
| 2: Breaking history change | R3, R6, R11 |
| 3: Retry rules | R12 |
| 4: Request identity and context | R13 |
| 5: Normal merge and fork push | R8, R9 |

## Execution Authorization

No product or scope decision remains open.
The requested deeper code research is complete.
The code reports define the native adapters, plugin no-replay boundary, inspection wait limit, and focused check owners.
The user approved implementation with one final adjustment: migrate retained Candy and Pelican histories when detecting the old format.
The task audit keeps existing native test owners and adds only the missing shared-runner test owner.
