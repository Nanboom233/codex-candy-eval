# U18 Final Browser Check

The final v0.2.17 binary passed browser checks in the retained `sub2api-codex-eval-test` container.
The fixture helper reused container `bd6c9583c3e6` at `http://127.0.0.1:60198`.

The browser enabled automatic inspection through the settings switch and its native automatic save.
Candy showed `下次巡检：10-17 23:18` inside the account header, beside `账户` and `仅启用`.
The separate lower runtime card was absent from the DOM.
Other evaluation headers contained no inspection status nodes.

The desktop screenshot showed the timestamp on the account header line.
At a 430-pixel viewport, the timestamp wrapped inside the same header.
The embedded document had equal client and scroll widths, 381 pixels.
Both screenshots were inspected:

- `output/playwright/U18-header-desktop.png`
- `output/playwright/U18-header-narrow.png`

The browser disabled automatic inspection after the check.
The next-run text cleared and its empty status node became hidden.
Both inspection action buttons remained on settings. Save-status text stayed empty.
No evaluation request was needed for this layout correction.

Independent source review and existing UI checks passed.
Unchanged U15 backend/timer and U17 error/completion evidence remains applicable.
Build log: `/tmp/sub2api-U18-build.log`.
Binary: `dist/build/sub2api-codex-eval`, 125599904 bytes.
Archive: `dist/build/sub2api-codex-eval.tar.gz`, 38700695 bytes.
