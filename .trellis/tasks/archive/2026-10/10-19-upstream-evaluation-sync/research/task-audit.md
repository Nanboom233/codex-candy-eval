# Task Audit Before Execution

The user approved execution with minimal migration of retained history.
The PRD, design, implementation plan, and history research now record that decision.

## Requirement Boundary

Detect unversioned native history and save Candy/Pelican records in the schema-1 snapshot during load.
Retain their existing payload fields and native account IDs.
Discard retired fingerprint and ModelTrace records without backup, conversion, or a legacy page.
Reuse the existing atomic saver and error display.
Keep every other approved behavior unchanged.

## Check Boundary

Adapt only the two upstream detector algorithm tests and their golden fixture.
Add one native shared-runner test owner because no current test directly owns that behavior.
Extend existing request, sender, runtime, inspection, and UI owners for their changed contracts.
Put minimal migration assertions in the current persistence fixture.
Do not add another migration suite, future-format matrix, or duplicate algorithm checks across layers.
Reuse unchanged installer and stream-lifetime evidence.
Keep the existing build entry and bounded test selector.
Do not add full-suite, race, coverage, benchmark, architecture, production-account, or hash gates.

## Authorization

The latest user message approves execution after this audit.
No further planning approval is required for this explicit adjustment.
Ordinary work commits still require the existing final files/message/signing review.
