# Pelican Browser Check

## Candidate and Fixture

The effective source is `/tmp/pelican-variants-native/source`.
The temporary fixture is `/tmp/pelican-variants-preview.cjs` on port 60232.
Chrome rendered the native Vue layout and evaluation iframe with synthetic authentication, groups, accounts, settings, and two Pelican records.
The fixture used short prompt placeholders because the UI harness owns prompt and copy assertions.
No real account, history file, or supplier request entered this check.

## Visual Results

Both records appeared in the same account history with their saved version labels.
Their 4:3 previews rendered the supplied pages. The thumbnails held their animation position.
Opening the original record showed its version and animated page.
Successive screenshots showed the moving circle at different positions in the detail frame.
Opening the all-in-one record showed its own version and page.
Closing each detail returned to the list. The existing DOM assertion owns removal of the detail frame.

The initial layout clipped the batch button at a 1100-pixel parent viewport, with an approximately 780-pixel evaluation iframe.
The initial latest-result version also lost readable text through truncation.
The implementation owner fixed only the Pelican toolbar breakpoint and latest-result wrapping.

The repaired 1100-pixel parent view showed group/model and effort/version pairs above a full-width batch button.
The complete saved version remained readable in the latest result.
At a 600-pixel parent viewport, the toolbar retained the same pairs and full-width action.
The detail metadata and 4:3 page remained usable at that narrow width.
The normal 1920-pixel parent viewport retained the single-row toolbar and readable version metadata.
The temporary viewport override was reset after inspection.

## Verification Boundary

The existing UI harness owns preferences, startup restoration, prompt/copy updates, submitted payloads, and record-derived labels.
This check inspected visual layout, previews, animation, and detail interaction only.
The final generic text correction uses the approved `画面` labels. It does not change layout or behavior.
