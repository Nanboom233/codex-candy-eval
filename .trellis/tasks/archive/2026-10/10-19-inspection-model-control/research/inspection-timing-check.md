# Confirmed-state timing review

## Decision

No material timing issue found. The implementation meets AC5 and the
clarified confirmation boundary.

`inspector.run` applies stable accounts from the first-test completion
callback. It queues only accounts whose saved verdict may transition. It
keeps the existing degraded-then-normal confirmation-group order and the
recovery wait. During a confirmation phase, `testPhase` applies each
account from its own completion callback; it does not wait for the other
accounts in that phase to finish. `applyAccount` persists the verdict and
model changes before the optional priority update.

The existing coordinator test checks that a confirmed degraded account is
updated before an unrelated recovery confirmation starts. The focused
stable-state test checks that `max_tests_per_account=1` needs no extra
confirmation and that stable verdicts still apply with priority adjustment
off. These tests and the code path cover the clarified timing. No additional
test matrix is needed.

## Verification

The implementation evidence records passing focused Go tests, `go vet`,
formatting, diff checks, and the final build against the reviewed source
tree. The implementation inputs are unchanged, so those checks remain
valid. No source or test change was needed in this review.

## Full-scope review

The model-control helper and its focused service tests cover the complete
account-data path:

- Stable degraded results remove only exact identity entries, preserve
  aliases and wildcards, keep `schedulable`, and save actual removals with
  the degraded verdict.
- Stable normal results restore the saved removal list even when the switch
  is off or the global list changed. They preserve an administrator's
  current mapping, remove the `blacklist` placeholder, and clear pending
  removals in one native update.
- Matching model data and verdicts skip writes. Verdict-only changes use the
  existing extra merge. Actual removals use a de-duplicated saved list.
- Compound model failures leave credentials, pending removals, and verdict
  unchanged. The coordinator does not update priority after that failure.
- A successful verdict/model write precedes priority adjustment. Matching
  priorities are skipped. Priority failures retain the saved account state,
  and notifications are sent only after an actual priority write. A
  notification failure remains on the result without reversing the update.
- Existing cancellation and participation checks prevent later sends and
  account application. Manual fixed-count tests remain available. When an
  account opts out before a later send, the asynchronous pre-send callback
  marks it `unchanged` and returns before recording an error. The runtime
  does not propagate that callback result as a `startChecked` start error,
  and incomplete collection preserves `unchanged`. The timing test covers
  ready accounts applying before unrelated confirmations.

The final UI and Go-installer evidence is reusable for this correction:
the UI inputs and installer files did not change. The recorded UI harness and
browser comparison cover native selector structure, selection, custom entry,
keyboard dismissal, narrow layout, icons, and automatic saves. The recorded
installer checks cover parameter-free entry, filtered selection, immutable
container IDs, candidate and backup validation, readiness, rollback, signal
handling, redirect rejection, embedded payload, and license texts. No live
Docker installation or provider request was used.

## Findings

No material findings remain. Source tracing confirms that the opt-out
pre-send callback preserves `unchanged` without adding an evaluation error.
The model-control, persistence, priority, notification, UI, and installer
checks align with the active task and recorded evidence.
