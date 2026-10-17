# Sub2API Administrator Evaluation Recovery Plan

## Goal

Restore the four original evaluations as an administrator account-testing feature
on Sub2API v0.2.16 Linux/amd64. Add the user-selected single-request endpoint
POST /api/v1/admin/accounts/:id/codex-eval using native account connectivity-test
sending logic. Deliver a prebuilt replacement for the user's N150 iStoreOS host.

Status: planning. No product implementation or task start is authorized.
Approved details: S01—S05 and C02/C03/C04/C06/C07. C08—C12 approved by user: 无误 继续. C13—C17 approved by user: 是. Next review: research/core-fourth-approval.md (C18—C22, pending).

## Evidence and decisions

Official source commit: 01842576b3d5c1e2545b5a99ad555d36e28ed306.
The untouched source compiled into a static Linux/amd64 binary with native frontend.
Completed checks include explicit platform differences, failures and skips;
see research/pristine-v0.2.16-baseline.md and research/pristine-baseline-review.md.
Reuse this evidence without rerunning the baseline or claiming all tests passed.
The old installed package remains unaccepted; compilation/installation do not prove functionality.
Race was explicitly removed from acceptance.

Earlier IDs/anchors remain in research/planning-before-minimal-replan-prd.md;
the immediately preceding Key proposal is preserved in research/before-admin-channel-prd.md.
The approval ledger owns current decisions and historical supersession.

| Earlier IDs | Current disposition |
| --- | --- |
| B1—B9, B5a | Completed pristine baseline evidence; race is historical only. |
| B10—B16 | Individual review, fixed version, breaking changes and no compatibility framework retained. |
| B17—B25 / F01—F06 | Approved standalone frontend behavior retained. |
| B26 / B01—B05 | Native admin authorization and received-error preservation retained; normal-user gateway assumptions superseded. |
| B27 / B06—B10 | User/Key/accounting identity design superseded; no exclusive administrator permission or automatic repair restored. |
| B28 / B11—B15 | Expanded gateway proposal rejected; subsequent shared-Key/target-header replacement also superseded. |
| C01/C05 | Required sending boundary covered/replaced by approved S02/S05. |
| C02/C03/C04/C06/C07 | Approved by user: 是 继续; see research/core-admin-channel-approval.md. |
| C08—C12 | Approved by user: 无误 继续; see research/core-second-approval.md. |
| C13—C17 | Approved by user: 是; see research/core-third-approval.md. |
| S01—S05 | S01—S05 administrator account endpoint details approved; remaining individual reviews pending. |

## Requirements

R1. Preserve the original four evaluations, question banks, algorithms, data and UI functions.
Review remaining core/history differences before changing them. No legacy migration
or implicit deletion of user records is authorized.

R2. Preserve approved F01—F06: /codex-eval through the native administrator custom menu,
zero native Vue edits, same-origin native administrator session and approved account/model APIs,
original toolbar/results/Pelican/filtering/polling/history/clear controls.

R3. Use the native administrator testing channel. Add the user-selected
POST /api/v1/admin/accounts/:id/codex-eval beside the existing /test endpoint.
The path's account ID is the sole sending target; there is no pool fallback.
Original /test behavior remains intact.

R4. Remove the internal-user, internal-Key, owner/accounting-row, expiry/renewal,
target-header and normal-user gateway integration requirements.
No Sub2API user billing/balance/quota/group-pricing machinery is added for evaluations.
Supplier requests still consume supplier quota; reported token usage is evaluation data.
The previous five-minute Key decision is superseded because that Key is no longer part of the design.

R5. Base sending on native account-test authentication, credential resolution,
model mapping, URL validation, proxy, TLS and native plugin transport.
Carry the complete original evaluation request and return complete answer/usage/error data.
Concrete full streaming request/response and native OpenAI Responses/CC boundaries are approved in S01—S03.

R6. Preserve original native API-key middleware, scheduler, sticky policies,
ordinary gateway billing/parsers, listening and shutdown.
Remove earlier evaluation-specific gateway, diagnostics, partial billing, cancellation,
private capability/pipe/producer and maintenance/activation additions during approved implementation.

R7. Preserve errors actually received by the evaluation layer under B05.
Account-test state side effects and task integration/cancellation follow approved S04—S05.
Do not infer user-billing drain guarantees for the new administrator testing channel.

R8. Build from pinned official source in the development environment using one
fixed-version ordinary patch and required evaluation resources. Remove multi-version
AST adaptation, compatibility configuration and obsolete migration/framework code.
Preserve upstream declared tool versions/locks and completed baseline evidence.

R9. Outputs: dist/build/sub2api-codex-eval and dist/build/sub2api-codex-eval.tar.gz.
Include install.py, standalone cleanup.py and necessary version/license files.
Installation never fetches source or compiles on the user's host. Cleanup targets
identifiable old evaluation build resources and preserves unrelated/business data.

R10. Record each approval with full explanations in batches of five.
S01—S05 and C02/C03/C04/C06/C07 are approved. Remaining core/history/delivery must converge before final implementation review.
The approved administrator batch does not authorize remaining core/history/delivery or product implementation.

## Acceptance criteria

A1—A6. Preserve completed baseline evidence with original failures/skips;
exact original statements remain in research/planning-before-minimal-replan-prd.md.
A7. Every changed scope item has current recorded approval before implementation.
A8. The new authenticated administrator endpoint sends only to account A; failures never send to B.
A9. Non-administrators cannot invoke the new endpoint. Original /test and ordinary gateway behavior remain intact.
A10. Evaluations need no internal Key/user/accounting setup and do not debit Sub2API users.
Actual supplier-reported usage and errors reach the evaluation without invented billing rows.
A11. Approved UI and original four evaluation behaviors remain available.
A12. The package contains the correct Linux/amd64 binary and standalone cleanup script.
A13. Package installation/rollback/cleanup are verified with disposable resources;
report build, tests, installation and actual feature results separately.
A14. Retired: internal-Key expiry/renewal acceptance is superseded by R4.
A15. Full request and completion/usage/error contracts are verified with local fixtures
for the protocols approved in S01—S03, without paid supplier requests.

Use native required and approved relevant checks; do not add race, paid production
requests or a new full-suite acceptance framework.

## Boundaries and open review

Authorized now: source research and planning artifact updates only.
No product code, canonical package, private configuration, external menu setting,
production request, database/key/history cleanup or deployment changes.
Use standard system temporary directories for development work.

Open: remaining core/history/delivery reviews. S01—S05 are approved;
C01/C05 host sending is covered by S02/S05. No Key-owner question remains.
Converge PRD/design/implementation plan and obtain final approval before task.py start.
