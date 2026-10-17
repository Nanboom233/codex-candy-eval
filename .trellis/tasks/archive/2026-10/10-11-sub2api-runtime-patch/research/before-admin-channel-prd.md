# Sub2API Minimal Recovery Plan

Latest steering: pause C01—C05 and investigate extending native account connectivity tests.
See research/native-connectivity-dispatch.md. The Key/local-gateway requirements below
are the prior planning branch under review, not authorization to proceed. The administrator-test
alternative changes user billing and scheduler semantics and awaits user choice.

## Goal

Restore the four original evaluations on Sub2API v0.2.16 Linux/amd64.
Add only internal-key authentication and hard account selection to the native gateway.
Deliver a prebuilt binary package for replacement on the user's N150 iStoreOS host.

Status: planning. No product implementation or task start is authorized.
The full current proposal is research/minimal-replan.md.

## Evidence and earlier decisions

The official source commit is 01842576b3d5c1e2545b5a99ad555d36e28ed306.
The untouched source compiled into a static Linux/amd64 binary with the native frontend.
Required-check results include explicit platform differences, failures and skips.
See research/pristine-v0.2.16-baseline.md and research/pristine-baseline-review.md.
Reuse this completed evidence. Do not claim every test passed or repeat baseline work without a new reason.

The old installed package remains unaccepted by the user.
Its installer and build evidence does not prove evaluation functionality.
The user withdrew the assistant-added race requirement.

Earlier requirement IDs and anchors remain losslessly preserved in research/planning-before-minimal-replan-prd.md.
Current dispositions follow the table below. The approval ledger owns exact authorization records.

| Earlier IDs | Current disposition |
| --- | --- |
| B1—B9, B5a | Completed pristine research and original failure evidence. Race is excluded from acceptance. |
| B10—B12 | Review actual modifications individually, with full explanations in batches of five. |
| B13—B16 | Fixed version, breaking changes, no compatibility framework, frontend decisions reviewed separately. |
| B17—B25 | F01—F06 approved. See research/frontend-first-approval.md. |
| B26 | B01—B05 approved. See research/backend-first-approval.md. |
| B27 | B06—B10 historically approved. See research/backend-identity-approval.md. Per-user parts of B06/B08/B09 and the expiry exemption are superseded below. |
| B28 | The expanded B11—B15 proposal is rejected. Internal key plus selection header replaces it. |

The older B14 breaking-change list is evidence, not the current design.
The latest proposal and user instructions override conflicting historical specs.

## Current requirements

R1. Preserve the original four evaluations, question banks, algorithms, data, and UI functions.
Review remaining core and history differences before changing them.
Do not migrate legacy formats or implicitly delete user records.

R2. Preserve approved frontend decisions F01—F06.
Use /codex-eval through the native administrator custom menu without editing native Vue source.
Use the native same-origin administrator session and approved account/model APIs.
Retain approved toolbar, result, Pelican, filtering, polling, history and clear controls.

R3. Use one shared internal Key and native accounting row for evaluations.
Native administrator login still controls start/history/management operations.
Usage belongs to that Key's native owner, not each initiating administrator.
Preserve request-local group binding and normal native limits.
Do not repair disabled/deleted Key data automatically.
Latest steering supersedes the per-user parts of B06/B08/B09.

R4. Reuse the shared Key with a five-minute sliding expiry through native ExpiresAt.
Before each evaluation HTTP call, the trusted evaluation module resets ExpiresAt to now + 5 minutes.
The gateway validates expiry but never renews it merely because a bearer request arrives.
Remove per-administrator mappings, per-request keys, consumption tables and body/header binding.
The next authorized evaluation call may renew the same expired Key through the native update path.
Disabled/deleted Keys are not repaired; balances and quotas are not reset.
See research/shared-key-expiry.md.

R5. Proposed header: X-Codex-Eval-Target: <group_id>:<account_id>.
Honor it only after internal-key validation and only on POST /v1/responses.
Reject ordinary-key attempts, malformed targets and missing required targets.
Forward only the target account. Failed or excluded targets never fall back to another account.

R6. Use the native local HTTP listener for evaluation requests.
Remove the custom pipe transport and producer completion framework.
Native eligibility, concurrency, same-account retries, protocol conversion, errors, billing and disconnect handling remain authoritative.
Restore native listening and shutdown. Remove evaluation-specific diagnostics, partial billing and cancellation changes.

R7. Preserve B05 errors received by the evaluation layer.
Do not intercept or reconstruct upstream errors already transformed by the native gateway.
Evaluation cancellation stops further work. Native upstream draining retains its original semantics.

R8. Build from the pinned official source in the development environment.
Use one fixed-version ordinary patch and required evaluation resources.
Remove multi-version AST adaptation, compatibility configuration and obsolete migration/framework code.
Keep source locks and original baseline evidence intact.

R9. Final outputs are dist/build/sub2api-codex-eval and dist/build/sub2api-codex-eval.tar.gz.
The package includes install.py, standalone cleanup.py, and necessary version/license files.
Installation never retrieves source or compiles on the user's host.
Cleanup targets identifiable old evaluation build resources and preserves unrelated resources and business data.

R10. Keep changes individually reviewable and recorded.
Complete remaining core/history/delivery decisions before final implementation review.
Do not infer approval for those items from the new backend simplification.
The next proposed batch C01—C05 is fully described in research/core-first-approval.md; it is pending user approval.

## Acceptance criteria

A1—A6. Preserve the completed baseline source/build/check evidence and independent review.
Their exact original statements remain in research/planning-before-minimal-replan-prd.md.
A7. Every changed scope item has current recorded approval before implementation.
A8. Valid internal key and target select only account A, including failure and retry paths.
A9. Ordinary keys cannot activate directed forwarding. Ordinary native requests retain their behavior.
A10. Accounting belongs to the shared Key owner with native group, balance and quota checks.
A14. Each trusted evaluation call resets the same Key to five minutes; gateway requests cannot renew it.
After five minutes without a trusted renewal, new requests are rejected. An admitted request may finish.
Renewal failure prevents dispatch; native authentication cache invalidation accompanies renewal.
A11. The approved evaluation UI and original four evaluation behaviors remain available.
A12. The package contains the correct Linux/amd64 binary and the standalone cleanup script.
A13. Package installation, rollback and cleanup are verified with disposable resources.
Report build, tests, installation and feature results separately.

Use native required checks and already approved relevant checks.
Do not add race, paid production evaluation or new full-suite acceptance frameworks.

## Boundaries and open review

Authorized now: source research and planning artifact updates.
Existing product code, canonical package, external menu settings and user data remain untouched.
No CPA deployment, private configuration access, production requests or actual cleanup is authorized here.
Use standard system temporary directories for development work.

The user has rejected per-administrator and single-use request keys.
The shared Key uses the user-approved five-minute sliding expiry and existing target header.
Concrete shared-owner preparation remains for review.
Remaining core/history/delivery changes retain their individual review gate.
Await final planning approval before task.py start.
