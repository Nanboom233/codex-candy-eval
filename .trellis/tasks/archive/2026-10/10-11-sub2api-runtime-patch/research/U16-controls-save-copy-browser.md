# U16 Final Browser Check

The final v0.2.17 candidate passed the changed UI checks in the retained `sub2api-codex-eval-test` container.
The fixture helper reused the same container and origin, `http://127.0.0.1:60198`.
The Linux runtime reported Sub2API 0.2.17.

Both inspection buttons appear at the end of the settings form.
They remain visible with automatic inspection off and its configuration collapsed.
Candy contains neither button. Its account results and history remain in place.

The browser edited the notification username and confirmed the value through the native settings API.
After the successful automatic save, the status text was empty.
An invalid notification server showed a validation error.
Restoring valid input saved automatically and cleared the error without displaying `已保存`.
The fixture username was restored to empty. Automatic inspection remains off.

Screenshot: `output/playwright/U16-settings-buttons.png`.

Independent review: `U16-controls-save-copy-check.md`.
Build log: `/tmp/sub2api-U16-build.log`.
Binary: `dist/build/sub2api-codex-eval`, 125599904 bytes.
Archive: `dist/build/sub2api-codex-eval.tar.gz`, 38701694 bytes, with seven runtime/license entries.

The U15 backend, group-selection, and timer evidence remains applicable.
No backend logic changed in U16.
