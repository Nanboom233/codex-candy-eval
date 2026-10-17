# U17 Final Browser Check

The final v0.2.17 candidate passed browser checks in the retained `sub2api-codex-eval-test` container.
The fixture helper reused container `bd6c9583c3e6` at `http://127.0.0.1:60198`.
No new container was created.

The browser changed the notification username and received a successful native settings PUT.
Both save-status locations stayed empty. The username was restored to empty through automatic persistence.
An invalid notification server showed the actual validation error `请输入网址。`.
Restoring the valid address saved automatically and cleared the error.

Both inspection buttons remained at the bottom of settings with automatic inspection off.
The browser started one local fixture round and switched to Candy.
Active progress showed `账户进度：0/2`.
After the round, phase and progress text were empty. The empty runtime card was hidden.
No row summary or removed status caption appeared. Existing account results and history remained visible.

The final screenshots were inspected:

- `output/playwright/U17-candy-no-summary.png`
- `output/playwright/U17-settings-silent-save.png`

The existing focused UI checks cover actual evaluation and notification errors in the shared runtime error location.
U15 backend and timer evidence remains applicable. No backend logic changed in U17.

Build log: `/tmp/sub2api-U17-build.log`.
Binary: `dist/build/sub2api-codex-eval`, 125599904 bytes.
Archive: `dist/build/sub2api-codex-eval.tar.gz`, 38701611 bytes.
The retained runtime reported Sub2API 0.2.17.
