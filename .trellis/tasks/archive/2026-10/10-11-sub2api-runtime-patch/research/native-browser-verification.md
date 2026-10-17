# Actual browser verification for U5-U6

Main session drove real headless Chrome through Playwright CLI against the final rebuilt artifact in disposable Docker. This replaces the earlier Node-VM-only acceptance for the affected sidebar/iframe/model paths.

## Actual user flow

1. Opened `http://127.0.0.1:55182/login` and filled the disposable native administrator login form.
2. Closed the native first-login tour, then clicked the actual sidebar `账户评估` link.
3. Native route `/admin/codex-eval` rendered the administrator layout, native beaker icon and evaluation iframe.
4. Filled the iframe model combobox with `browser-unconfigured` and clicked the account's actual `测试` button.
5. Candy completed successfully; its row/history displayed `browser-unconfigured(low)` and `1/1` correct.

No DOM mocks, API mocks, injected session token or programmatic replacement of the start action was used. Upstream was the existing local fake supplier, so this consumed no provider quota.

## Browser and model observations

Actual browser DOM evaluation after the flow returned:

```json
{
  "parentOrigin": "http://127.0.0.1:55182",
  "frameOrigin": "http://127.0.0.1:55182",
  "framePath": "/codex-eval",
  "framePort": "55182",
  "frameLoaded": true,
  "sidebarIcon": true,
  "iconPaths": 1,
  "rawSuggestions": ["gpt-5.5", "browser-unconfigured"],
  "selectedModel": "browser-unconfigured",
  "theme": "light"
}
```

Chrome console reported `Total messages: 0 (Errors: 0, Warnings: 0)` after iframe loading and the evaluation. The unchanged parent CSP enforced `frame-src 'self'`; the child document was readable from the parent and its controls were operated in the real frame.

Authenticated comparison in that same browser returned:

```json
{
  "normalModels": [],
  "evaluationModels": ["gpt-5.5", "browser-unconfigured"]
}
```

The normal path was `/api/v1/admin/accounts/1/models`; the evaluation path was `/api/v1/admin/accounts/1/codex-eval/models`. Stored mappings deliberately conflicted with both actual supplier model names. The normal native projection remained active while the evaluation suggestions exposed the supplier names.

Actual local upstream logs after the UI action included:

```json
{"path":"/v1/chat/completions","model":"browser-unconfigured","stream":true}
```

The fixture's preceding native capability probe separately used `provider-remapped` with `stream:false`; that ordinary behavior was retained.

## Visual evidence

Screenshot: `output/playwright/sub2api-native-port.png` in the repository workspace. Main viewed it: native sidebar beaker icon/selected entry, native header, loaded iframe, selected model, successful account result and history are visible. CLI snapshots for the native sidebar and real iframe are retained under `.playwright-cli/` as local artifacts.

Build/runtime ownership and teardown are recorded by the implementer in `native-browser-model-repair.md`; independent review is recorded in `native-browser-model-check.md`.

This is actual Chrome and Docker evidence on the development machine, including a nonstandard HTTP port. N150/iStoreOS execution and the user's production HTTPS origin were not exercised.
