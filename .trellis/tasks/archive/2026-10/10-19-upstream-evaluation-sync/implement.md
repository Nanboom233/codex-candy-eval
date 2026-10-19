# Implementation Plan

## Entry Condition

The user approved the complete plan and requested immediate execution after the task audit.
The final adjustment migrates retained Candy and Pelican records when an unversioned native file loads.
Use this task throughout integration and delivery.

## Execution Order

1. Activate this task with `task.py start` after approval.
2. Read the task artifacts and curated implementation context.
3. Inspect current changes and refresh upstream and fork refs.
4. Create `sync/upstream-main-2026-10-19` from local `main`.
5. Reconcile any newly fetched commits, then merge the reviewed upstream tip with `--no-commit`.
6. Resolve conflicts through the ownership table in `design.md` and the upstream commit inventory.
7. Retain `mustDecode` in `codex.go`, then replace the old fingerprint and ModelTrace owners in patches `0001` and `0002`.
8. Adapt the request context, typed transport boundary, and shared retry runner through the native sender.
9. Preserve Candy confirmation counts and check participation before retry sends.
10. Update task dispatch, state, details, history, and clear paths for the three test types.
11. Migrate unversioned Candy/Pelican history through the existing atomic saver and discard retired test records.
12. Adapt the current page, preferences, controls, and result rendering without explanatory small print.
13. Replace retired algorithm tests and test-only data under `tests/backend/`.
14. Extend existing inspection, runtime, request, and UI tests. Add only the missing native runner test owner.
15. Update the existing focused Go test selection in `scripts/build.py`.
16. Update Chinese README usage, Credits, detector license, and third-party notice.
17. Format changed Go code and run the bounded validation below.
18. Run a Trellis check agent against the complete candidate and fix in-scope findings.
19. Update affected specs from the checked implementation.
20. Present exact work commit files, message, and effective signing behavior for final approval.
21. Create approved work commits, then archive and record the session through Trellis.
22. Fast-forward local `main`, check the current fork tip, and push normally to `origin/main`.

Use one `trellis-implement` agent for the coupled retry, detector, history, and page contracts.
Use a separate `trellis-check` agent after implementation reaches a stable candidate.
Do not split this integration into new product tasks.

### Code Sequence Within the Native Adaptation

1. Retain the shared decoder helper before deleting `fingerprint_algorithm.go`.
2. Add the detector decoder and algorithm with its runtime assets in patches `0001` and `0002`.
3. Mark actual transport-attempt errors without marking host, proxy, client, or account setup errors.
4. Preserve observed response status before any error-body read or retry decision.
5. Move `askCodex` into `model.go` and adapt Candy, fingerprint, and Pelican callers together.
6. Remove `executeProbe` after both retired callers disappear.
7. Pass the existing inspection callback through the Candy worker before every attempt.
8. Handle intentional opt-out before record creation and progress increments.
9. Replace fingerprint progress and result types through `startChecked`, `State`, and detail handlers.
10. Add the raw history read boundary and migrate retained old records through the existing saver after successful decoding.
11. Remove ModelTrace dispatch, maps, detail and clear branches, retention, and save rollback fields.
12. Update native page state, fingerprint controls, result details, and embedded license.
13. Remove retired frontend state, preferences, script declarations, markup, and resource consumers.

Keep native numeric account IDs and active-run ownership throughout these steps.
Do not copy CPA routes, catalog admission, common-record rewrites, or migration helpers.

## Focused Validation

| Existing owner | Required changed behavior |
| --- | --- |
| `fingerprint_test.go` | Adapt only upstream `TestFingerprintUpstreamParity` and `TestFingerprintParseAndChallenges`, with their golden fixture |
| New `model_test.go`, `TestCodexEvalRetries` | Retry limit and waits, identical requests, cumulative usage, final statuses, validators, and backoff cancellation |
| Existing `codex_test.go` | Updated default identity/context and developer-instruction placement |
| Existing sender tests | Setup versus transport marker, retained status after error-body failure, and existing 401/429 account effects |
| Existing inspection tests | Wrong Candy text stays one answer, retries preserve confirmation counts, and opt-out prevents later sends and synthetic results |
| Existing runtime tests | Three dispatch paths, native cancellation, detector collection, minimal legacy migration, reload, and retained Candy/Pelican payloads |
| Existing UI harness | Three pages, detector controls/results, retired tab fallback, native group/model controls, and Settings-first entry |
| Actual browser | Changed page render, detector controls/details, and Settings-to-test group loading with synthetic data |
| Build entry | Patched frontend/server, focused Go tests and vet, and packaged Go installer with current licenses |

Reuse unchanged native sender, delayed-stream lifetime, installer, and inspection evidence.
Extend current owners for changed behavior. The native shared runner has no current focused owner.
Do not copy the CPA runner, collector, Candy, Pelican, or state-migration suites.
Use synthetic supplier responses and temporary history data.
No production account or provider quota is part of acceptance.

Replace the existing Go test selector with this expression:

```text
Codex|FingerprintUpstreamParity|FingerprintParseAndChallenges|ExtractHTML|HasStandalone21
```

Remove only the retired algorithm names and their test-only data.
The `Codex` fragment selects existing native runtime, sender, request, inspection, and new shared-runner tests.
Keep the existing UI harness, focused vet command, and host-platform installer tests.

Use each test owner for its own contract.
Algorithm tests own scores, sequence thresholds, aliases, and null probabilities.
Runtime tests own collected samples, status precedence, summary/detail persistence, and cancellation.
The UI harness owns rendered null probabilities, `not_model`, retired-tab fallback, and Settings-first navigation.
Do not repeat algorithm golden assertions through the page or service tests.
Extend the current history fixture for migration. Do not add a separate migration suite or future-format test matrix.

Run the established check once the candidate is ready:

```sh
python3 scripts/build.py --source /tmp/sub2api-codex-eval-v0217-check --check
```

The entry applies patches, builds the frontend, runs the UI harness, focused Go tests, and vet, then builds both executables.
Its installer checks remain in the host-platform temporary module.
Reuse the retained `sub2api-codex-eval-test` container only if the native browser check needs it.
Do not run the installer against that container's `/app/sub2api` dependency launcher.
Do not add another installer transaction, full suite, race, coverage, benchmark, architecture matrix, or hash gate.
Repeat checks only for changed inputs, actual failures, or unresolved findings.

## Sensitive Boundaries and Recovery

- Review `0001` history decoding before loading a retired fingerprint payload into the new record type.
- Keep Candy and Pelican JSON fields and native account IDs unchanged.
- Keep sender setup errors distinct from transient failures. Keep HTTP 401 and 429 final.
- Keep opt-out distinct from a failed logical Candy result.
- Keep native Settings initialization and account inventory ownership intact.
- Keep ModelTrace removals consistent across code, tests, resources, licenses, README, and build test selection.
- Keep the Git merge pending until the candidate passes checks and receives final commit approval.
- Preserve unrelated changes, ignored private files, the original Sub2API tree, and retained fixture data.

## Completion Evidence

Record actual checks and material limits with the final candidate.
Confirm that every upstream commit has a native destination or an exclusion reason.
Confirm that local `main` includes the checked integration, archive, and session commits.
Confirm that the normal fork push succeeds and local `main` tracks `origin/main`.
Deliver the result, build output paths, and any unresolved acceptance issue.

## Execution Checkpoint

Implementation, bounded checks, browser inspection, independent review, and affected spec updates are complete.
See `research/implementation-handoff.md`, `research/browser-check.md`, and `research/independent-check.md` for the candidate and actual evidence.
The staged work set contains 26 paths. It excludes task and workspace bookkeeping.
The pending merge retains local `2aedffe` and approved upstream `7a49068` as its parents.
Git signing is enabled through the existing configuration.
The next action is final approval of the exact signed merge commit.
After that commit, finish the task, record the session, fast-forward local `main`, and push normally to the fork.
