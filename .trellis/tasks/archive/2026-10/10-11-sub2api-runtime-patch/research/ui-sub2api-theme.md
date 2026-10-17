# Research: Sub2API Theme Adaptation for the Embedded Evaluation UI

- **Query**: Map the existing CPA evaluation page (`internal/plugin/web/*`) to the pinned Sub2API v0.2.15 design tokens and theme mechanism, keeping structure unchanged.
- **Scope**: internal (pinned upstream `dist/sub2api-v0.2.15` @ `f2669c8cf62555cd92389b3f55920e9e6e7c6ff2` + local `internal/plugin/web/`). No external searches needed; local source is authoritative.
- **Date**: 2026-10-11

## 1. Verified Upstream Tokens and Styles

All upstream anchors are under `dist/sub2api-v0.2.15/frontend/`.

### Palette (`tailwind.config.js`)

| Role | Scale | Key values |
|---|---|---|
| `primary` (accent) | Tailwind **teal** | 50 `#f0fdfa`, 100 `#ccfbf1`, 400 `#2dd4bf`, 500 `#14b8a6`, 600 `#0d9488`, 700 `#0f766e`, 900 `#134e4a`, 950 `#042f2e` (`tailwind.config.js:9-21`) |
| `accent` and `dark` | Tailwind **slate**, identical scales | 50 `#f8fafc`, 100 `#f1f5f9`, 200 `#e2e8f0`, 300 `#cbd5e1`, 400 `#94a3b8`, 500 `#64748b`, 600 `#475569`, 700 `#334155`, 800 `#1e293b`, 900 `#0f172a`, 950 `#020617` (`tailwind.config.js:23-49`) |
| `gray` | default Tailwind gray | 50 `#f9fafb`, 100 `#f3f4f6`, 200 `#e5e7eb`, 300 `#d1d5db`, 400 `#9ca3af`, 600 `#4b5563`, 700 `#374151`, 900 `#111827` |
| `darkMode: 'class'` | `tailwind.config.js:4` — `.dark` on `<html>` | |

### Base and semantic usage (`src/style.css`)

| Surface | Light | Dark | Anchor |
|---|---|---|---|
| body bg / text | `bg-gray-50` `#f9fafb` / `text-gray-900` `#111827` | `bg-dark-950` `#020617` / `text-gray-100` `#f3f4f6` | `style.css:20` |
| default border | `border-gray-200` `#e5e7eb` | `border-dark-700` `#334155` | `style.css:7` |
| card | `bg-white`, `border-gray-100`, `rounded-2xl` (1rem), `shadow-card` = `0 1px 3px rgba(0,0,0,.04),0 1px 2px rgba(0,0,0,.06)` | `bg-dark-800/50`, `border-dark-700/50` | `style.css:228-234`, `tailwind.config.js:72` |
| primary button | gradient `primary-500→600`, white text, `shadow-primary-500/25`; hover `600→700` | `shadow-primary-500/20` | `style.css:82-87` |
| secondary/ghost button | `bg-white border-gray-200 text-gray-700`, hover `bg-gray-50 border-gray-300` | `bg-dark-800 border-dark-600 text-gray-200`, hover `bg-dark-700 border-dark-500` | `style.css:89-100` |
| danger/success/warning buttons | gradients `red-500→600`, `emerald-500→600`, `amber-500→600` | shadow `*/20` | `style.css:102-120` |
| focus | `focus:ring-2 ring-primary-500/50 ring-offset-2`; inputs `focus:border-primary-500 ring-primary-500/30` | same | `style.css:77,174` |
| input | `bg-white border-gray-200`, `rounded-xl`, placeholder `text-gray-400` | `bg-dark-800 border-dark-600`, placeholder `dark-400`, disabled `bg-dark-900` | `style.css:167-176` |
| table | th `text-gray-600 bg-gray-50 border-gray-200`; td `text-gray-700 border-gray-100`; row hover `bg-gray-50`; container `border-gray-200 rounded-xl` | th `text-dark-300 bg-dark-800/50 border-dark-700`; td `text-gray-300 border-dark-800`; hover `bg-dark-800/30` | `style.css:311-339` |
| badges | `primary-100/primary-700`, `emerald-100/700`, `amber-100/700`, `red-100/700`, `gray-100/700`, `purple-100/700` | `*-900/30` bg + `*-400` text (`gray` → `dark-700/dark-300`) | `style.css:342-369` |
| tabs | container `bg-gray-100`, tab `text-gray-600`, active `bg-white text-gray-900 shadow-sm` | container `bg-dark-800`, tab `text-dark-400`, active `bg-dark-700 text-white` | `style.css:638-654` |
| status pills (PluginsView) | enabled `emerald-100/700`, error `red-100/700`, starting `amber-100/700`, disabled `gray-100/600` | `emerald-900/30`+`emerald-300`, `red-900/30`+`red-300`, `amber-900/30`+`amber-300`, `dark-700`+`gray-300` | `PluginsView.vue:705-723` |
| progress/switch | progress track `bg-gray-200`, bar gradient `primary-500→400`; switch off `bg-gray-300`, on `bg-primary-500`, thumb `bg-white` | track `bg-dark-700`; switch off `bg-dark-600` | `style.css:657-685` |
| code / pre | inline `bg-gray-100 text-primary-600`; block `bg-gray-900 text-gray-100` (stays dark both themes) | `bg-dark-800 text-primary-400`; block unchanged | `style.css:688-699` |
| modal/dropdown | overlay `bg-black/50` (dialog `bg-black/40`); content `bg-white border-gray-200 rounded-2xl` | overlay `bg-black/60`; `bg-dark-800 border-dark-700` | `style.css:391-428,484-510` |
| tooltip-adjacent (toast) | `bg-white border-l-4` (emerald/red/amber/primary) | `bg-dark-800` | `style.css:512-537` |
| selection | `bg-primary-500/20`, `text-primary-900` | `text-primary-100` | `style.css:64-66` |
| radii / shadows | controls `rounded-xl` (0.75rem), cards `rounded-2xl` (1rem); `shadow-card`, `shadow-card-hover` `0 10px 40px rgba(0,0,0,.08)`, `shadow-glass` `0 8px 32px rgba(0,0,0,.08)` | same | `tailwind.config.js:67-75` |
| link color (markdown content) | `text-primary-500` hover `primary-600` | same | `CustomPageView.vue:556` |

Font stack: `system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif` (`tailwind.config.js:52-65`). Mono: `ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace`.

## 2. Token Mapping: Local Variables → Sub2API

Local components already consume only CSS variables; every color flows through `theme.css`/`style.css` tokens (`grep` confirms zero hex colors in JS; the only JS color work is `data-theme`). Keep the token architecture; replace the token *values* and the host blocks. **Mapping fidelity note**: values marked with a `style.css` anchor are source-exact for that semantic role; rows labeled *design inference* map local constructs that have no upstream analog (tinted button hovers, tooltips, bordered warning tags, soft-danger buttons) — they borrow the nearest upstream token rather than claim exactness.

### Replace host blocks

`theme.css:19-130` defines two host blocks (`:root:not([data-cpamp-plugin-host])` for CPAMC, `:root[data-cpamp-plugin-host]` for CPAMP). Both must be replaced by **one** Sub2API token block keyed on `:root` + `:root[data-theme="dark"]` — keeping the existing `data-theme` attribute is the recommended canonical mechanism (see §3); it minimizes selector churn in `style.css`/`theme.css` and needs no parent CSS-class convention inside a separate document. Delete the CPAMP blob animation, `--glass` backdrop filter, `body::before` top band, and the `header{display:none}` CPAMP rule — Sub2API's `CustomPageView` renders no plugin chrome, so the eval page's own header should stay visible but can be slimmed to match `page-title`/`page-description` styling (`style.css:591-601`).

### Variable mapping (recommended values)

| Local token | Sub2API light | Sub2API dark | Source |
|---|---|---|---|
| `--bg` | `#f9fafb` (gray-50) | `#020617` (dark-950) | `style.css:20` |
| `--bg-image` | `none` | `none` | upstream body has no image |
| `--surface` (cards) | `#ffffff` | `rgb(30 41 59 / 50%)` (dark-800/50) | `style.css:229` |
| `--sunken` (insets, th, pre) | `#f9fafb` (gray-50) | `rgb(30 41 59 / 50%)` | `style.css:322` |
| `--field` | `#ffffff` | `#1e293b` (dark-800) | `style.css:169` |
| `--line` / `--card-line` | `#e5e7eb` (gray-200) / `#f3f4f6` (gray-100) | `#334155` (dark-700) / `rgb(51 65 85 / 50%)` | `style.css:7,231` |
| `--ink` | `#111827` (gray-900) | `#f3f4f6` (gray-100) | `style.css:20` |
| `--muted` | `#6b7280`-family gray-500/600 (`#6b7280`) | `#94a3b8` (dark-400) | `style.css:91,295,321` |
| `--faint` | `#9ca3af` (gray-400) | `#64748b` (dark-500) | `style.css:587` |
| `--control-bg` / `--control-line` | `#ffffff` / `#e5e7eb` | `#1e293b` / `#475569` — **contextual**: form controls and secondary buttons use `dark-600` (`#475569`), NOT the `dark-700` (`#334155`) used for card/table borders | `style.css:89-95,170` |
| `--control-hover` / `--control-hover-line` | `#f9fafb` / `#d1d5db` | `#334155` (dark-700) / `#64748b` (dark-500) | `style.css:93-94` |
| `--accent` | `#14b8a6` (primary-500) | `#2dd4bf` (primary-400) | `tailwind.config.js:9-21` |
| `--accent-ink` | `#0d9488` (primary-600) | `#2dd4bf` (primary-400) | `style.css:275,576` |
| `--accent-soft` / `--accent-soft-hover` | `#ccfbf1` (primary-100) / `#99f6e4` (primary-200) | `rgb(19 78 74 / 20%)` / `rgb(19 78 74 / 30%)` — exact for selected/hover via `sidebar-link-active`; **design inference** for `.btn` and `.tab` hover reuse, upstream has no tinted-button equivalent | `style.css:348,574-577` |
| `--action-line` | `rgb(20 184 166 / 30%)` | `rgb(45 212 191 / 30%)` | `ring-primary-500/30` |
| `--success` / `--success-ink` / `--success-soft` | `#10b981` / `#047857` (emerald-700) / `#d1fae5` (emerald-100) | `#34d399` (emerald-400) / `#34d399` / `rgb(6 78 59 / 30%)` (emerald-900/30) — exact badge semantics; verdict-pill reuse is design inference | `style.css:351,707` |
| `--danger` / `--danger-ink` / `--danger-soft` | `#ef4444` / `#b91c1c` (red-700) / `#fee2e2` (red-100) | `#f87171` (red-400) / `#f87171` / `rgb(127 29 29 / 30%)` (red-900/30) — exact badge semantics; `.btn.danger` uses a red **gradient** upstream (`style.css:102-106`), so the local soft-danger button keeps soft styling as design inference | `style.css:360,709` |
| `--warning-ink` / `--warning-soft` / `--warning-line` | `#b45309` (amber-700) / `#fef3c7` (amber-100) / `rgb(245 158 11 / 30%)` | `#fbbf24` (amber-400) / `rgb(120 53 15 / 30%)` (amber-900/30) / `rgb(251 191 36 / 30%)` — exact badge semantics; `--warning-line` is a local invention (upstream badges are borderless), keep for `.tag.warn`/flash borders | `style.css:355,711` |
| `--tag-neutral-*` | `bg-gray-100 text-gray-700` + `border-gray-200` | `bg-dark-700 text-dark-300` + `border-dark-600` | `badge-gray` `style.css:363-365` |
| `--history-bg` / `--history-head-bg` | `var(--bg)` / `#f9fafb` | `var(--surface)` / `rgb(30 41 59 / 50%)` | table header pattern `style.css:319-324` |
| `--radius` | `12px` keep, or `16px` for cards + `12px` controls (upstream: `rounded-2xl` cards, `rounded-xl` controls) | same | `style.css:221-234` |
| `--shadow` | `0 1px 3px rgba(0,0,0,.04), 0 1px 2px rgba(0,0,0,.06)` | same (upstream does not vary card shadow) | `tailwind.config.js:72` |
| `--float` | `0 10px 15px -3px rgb(0 0 0 / 10%)` or `shadow-2xl` for modals | same | `style.css:377,401` |
| `--glass` | `none` (delete; Sub2API uses `backdrop-blur` only on `.glass` opt-in classes) | `none` | — |
| `--tip-bg` / `--tip-ink` / `--tip-line` | `#1e293b` / `#f3f4f6` / `#334155` — **design inference**: upstream has no plain tooltip token; chosen as an inverted `dark-800` surface for contrast in both themes | same (upstream dialogs/toasts reuse `bg-dark-800`, no dark-variant tip exists) | nearest analog: `style.css:372-380` |
| `--popover-bg` / `--dialog-bg` | `#ffffff` | `#1e293b` (dark-800) | `style.css:374,400` |
| `dialog::backdrop` (`style.css:284`) | `rgb(0 0 0 / 40%)` | `rgb(0 0 0 / 60%)` | `style.css:486` |
| `--page-top`, `--page-inline`, `--section-gap`, `--page-width` | keep (layout, not theme); `--page-width` may stay 1280 (close to upstream container) | same | — |

### Hardcoded / special cases

| Location | Item | Decision |
|---|---|---|
| `theme.css:4-17` `--premium-*`, `--plus-*`, `--go-*`, `--business-*` + `style.css:171-183` `.plan-*` badges | CPA plan-tier badge materials (gold/platinum gradients) with no Sub2API analog. **Recommendation (conservative, not approved)**: drop the metallic `--premium-*`/`--plus-*`/`--go-*`/`--business-*` and `.plan-elite` materials entirely; render whatever account metadata `/targets` actually returns (account name, platform, group name) through the neutral native badge tokens (`badge-gray`: `bg-gray-100 text-gray-700` light, `bg-dark-700 text-dark-300` dark; `badge-primary` teal for the selected/target group if a label is needed). Do NOT invent group→plan semantics or arbitrary per-group colors. |
| `style.css:62` switch thumb `#fff` | keep — upstream thumb is `bg-white` (`style.css:677-683`). |
| `style.css:274` `.preview iframe { background:#fff }` | keep — pelican srcdoc pages are opaque-origin and render their own page; upstream plugin iframe also pins `bg-white dark:bg-dark-900` on the frame element only (`PluginsView.vue:321`). |
| `style.css:12` `:where(body …) { all: revert }` CPAMP base-style reset | **remove** — it exists only to undo CPAMP's injected base styles; under Sub2API nothing injects, and `all: revert` would fight the new token sheet. |
| `--marks-width`, `--list-columns`, `--toolbar-columns`, media-query row reflow (`style.css:360-410`) | keep unchanged — structure/layout, not color. |
| `.markdown` rules (`style.css:228-239`) | **A single `--sunken` cannot represent both table/inset surfaces and upstream's fixed-dark `pre`.** Upstream `pre`/`code-block` stays `bg-gray-900 text-gray-100` in BOTH themes (`style.css:695-699`, `CustomPageView.vue:563`); local `.markdown pre`, `.problem-body`, `.fp-stat`, `.mt-summary`, `.preview` all share `--sunken`. Recommended split (no layout change): keep `--sunken` for insets/th/td-header; add `--code-soft`/`--code-ink` for inline code (light `gray-100`/`primary-600`, dark `dark-800`/`primary-400`) and `--pre-bg`/`--pre-ink` for code blocks (`#111827`/`#f3f4f6`, identical in both themes). Apply `--pre-*` to `.markdown pre`, `.problem-body pre`, `.mt-sample pre`; leave `.fp-stat`/`.mt-summary`/`.preview` on `--sunken`. |
| `ui.html:18-24` `#login` CPA password form | remove/replace — Sub2API embeds provide the admin session; auth comes from `auth_token` headers per patch-contracts §7, no login form. |

## 3. Theme Ownership and Synchronization in the Iframe

### Verified upstream behavior

1. **Storage key**: `localStorage['theme']` ∈ `'light' | 'dark'`, written by `AppSidebar.vue:866-870` (`toggleTheme`); missing → `prefers-color-scheme` fallback (`main.ts:25-31`, `AppSidebar.vue:935-943`). No `'white'`/`'light'`/'paper' variants — only two themes.
2. **Class**: `documentElement.classList` toggles `dark` (`main.ts:30`, `AppSidebar.vue:868`). Tailwind `darkMode:'class'`.
3. **Iframe URL**: `buildEmbeddedUrl` appends `theme=dark|light` (+ `user_id`, `token`, `lang`, `ui_mode=embedded`, `src_host`, `src_url`) (`utils/embedded-url.ts:16-46`). `detectTheme()` reads the parent's `html.dark` (`embedded-url.ts:48-51`).
4. **Runtime sync**: `CustomPageView.vue:424-435` mounts a `MutationObserver` on `documentElement` `class` attribute → `pageTheme` ref updates → `embeddedUrl` recomputes → **iframe `src` changes and the iframe reloads**. There is no postMessage theme push; sync is reload-based. Observer disconnected in `onUnmounted` (`CustomPageView.vue:446-452`).
5. **Planned patch** (`patch-contracts.md §7`, `source-audit.md:88-94`): narrow `embedded-url.ts` exception for the same-origin eval resource omits `token`, `user_id`, `src_host`, `src_url` but should **keep `theme`** (and optionally `lang`) so the reload-based sync still works.

### Recommended eval-page theme resolution

The eval page runs same-origin, non-sandboxed inside `CustomPageView`'s iframe. Precedence (first wins):

1. `?theme=` URL param (`dark|light`) — authoritative for this load; set by parent.
2. `localStorage['theme']` — same key as parent (shared origin).
3. `matchMedia('(prefers-color-scheme: dark)')` — final fallback, matching `main.ts:25-31`.

Apply by setting `document.documentElement.dataset.theme = 'dark' | 'light'` — **recommendation (conservative)**: retain the existing local `[data-theme]` selectors as canonical. The iframe document need not adopt the parent's `.dark` class convention; keeping `data-theme` minimizes CSS churn and the two mechanisms are equivalent inside a separate document.

Runtime update reality — a `storage` listener alone does NOT prevent reload: `CustomPageView` watches `documentElement` class and recomputes `embeddedUrl` reactively, so a theme toggle changes the iframe `src` and reloads it regardless of anything the child does. Reload-free sync would require an additional parent-side patch (e.g., a postMessage theme push in `CustomPageView`, or omitting `theme` from the rebuilt URL) — beyond a listener and beyond the planned minimal `embedded-url.ts` exception.

**Recommendation (minimal scope)**: accept the native reload behavior. The page resolves theme at load time only, and restores transient UI from server state after reload (run progress, histories, selected tab via existing `PREF_STORE`).

Resolution semantics:
- **Load time (authoritative)**: `?theme=` param (always present when embedded, maintained by the parent) → `localStorage['theme']` (direct-open scenario) → `prefers-color-scheme` (last resort).
- **Live events**: `storage` and `matchMedia` listeners only matter for the direct-open scenario or a missing param; when embedded, the parent's reload supersedes them. Keep the existing listener pattern (`app.js:459-460`) only for that fallback case. The current Node harness excludes initialization, so this behavior needs new tests.
- **Cleanup**: the page is a single-document app with no SPA unmount; module-level listeners are fine.
- **Do NOT write** `localStorage['theme']` from the eval page — the sidebar owns it.

### What to delete from `app.js`

- `panelValue()` + the `enc::v1::` XOR obfuscation (`app.js:25-39`) — CPA-panel specific; Sub2API stores plain values.
- `panelKey()` reading `cli-proxy-auth` (`app.js:47-50`) — replaced by Sub2API `auth_token` read for API headers only.
- Theme whitelist `["light","white","dark"]` → `["light","dark"]`; `applyTheme` rewrite per §3 (`app.js:41-46,458-460`).
- `sessionStorage` CPA key bootstrap (`app.js:483-488,523-525`) — replaced by parent-provided admin token in headers; **no token in URL** (hard requirement, `patch-contracts.md §7`).

### Two upstream precedents to reuse

- `handler/admin/plugin_handler.go:250-254`: a route-local response overriding `X-Frame-Options: SAMEORIGIN` and setting a strict CSP with `frame-ancestors 'self'` — the exact pattern the eval UI route needs, since the global policy is `X-Frame-Options: DENY` + `frame-ancestors 'none'` (`security_headers.go:131`, `config.go:35`). Handler-set headers post-middleware is proven to work.
- Native plugin-UI bridge (`PluginsView.vue:315-326,600-703` + `plugin_manager.go`): opaque-origin `sandbox="allow-scripts"` iframe + postMessage `bridge_token`, asset CSP `connect-src 'none'`. **Not usable as-is for the eval page**: its CSP forbids all `fetch`, its only channels are `config.load/save/test`, `plugin.status`, `ui.resize`, `ui.notify`, and it receives no theme at all. Extending the bridge would be a much larger frontend patch; the custom-menu + `embedded-url` exception in the design is the smaller correct path. Documented here so the option isn't rediscovered blindly.

## 4. Minimal File Change Map

| File | Change | Size |
|---|---|---|
| `internal/plugin/web/theme.css` (or a renamed `tokens` block copied into `internal/eval/web/`) | Rewrite: delete both host blocks, CPAMP blobs/`--glass`/`body::before`, plan-badge materials; emit one `:root` + `:root[data-theme="dark"]` token set per §2 | rewrite ~130→~80 lines |
| `internal/plugin/web/style.css` | Remove `:where(…){all:revert}` line 12; retarget any literals per §2 special cases; keep all layout/grid/media rules | ~10-line edit |
| `internal/plugin/web/app.js` | Replace `panelValue/panelKey/applyTheme` + init lines 458-460,523-525,483-488 with Sub2API theme resolution + `auth_token` API helper; switch API envelope to `{code,message,data}` (`patch-contracts §6-7`) | localized rewrites in 5 spots |
| `internal/plugin/web/ui.html` | Remove `#login` card + footer CPA links/version placeholder semantics; keep everything else (tabs, toolbars, lists, dialogs, tip) | small |
| `internal/plugin/web/components.js`, `candy.js`, `fingerprint.js`, `modeltrace.js`, `pelican.js`, `credentials.js`, `catalog.js` | No color/theme code (verified: no hex/rgb in JS). Only data-shape changes for new API contract — out of this report's scope | 0 for theme |
| `pelican.js:30` `sandbox="allow-scripts"` srcdoc | keep as-is — already opaque-origin, matches upstream plugin-iframe precedent | 0 |
| upstream `frontend/src/utils/embedded-url.ts` | patch: same-origin eval resource → `theme` (+ `lang`) only, no `token`/`user_id`/`src_*` | existing patch contract |
| upstream eval route handler | replicate `plugin_handler.go:250-254` pattern: `X-Frame-Options: SAMEORIGIN`, route CSP with `frame-ancestors 'self'`. **Recommendation**: use the host CSP nonce (`GetNonceFromContext`, `security_headers.go:98-106`) injected into the assembled HTML's `<script>`/`<style>` tags, NOT `'unsafe-inline'` — the page embeds 11 inline scripts and Temml CSS, and nonce preserves the global CSP strength the design requires | existing patch contract |

`plugin.go:34-38,85` embeds `theme.css`+`style.css` behind `/*APP_STYLES*/` — the assembly point survives the move to `internal/eval/` unchanged.

## 5. Test Matrix

| # | Scenario | Method | Assert |
|---|---|---|---|
| T1 | Light theme iframe load | menu → eval page | `data-theme="light"`; `--bg` `#f9fafb`; no white-on-white text |
| T2 | Dark theme iframe load | set localStorage `theme=dark`, reload menu | `data-theme="dark"`; `--bg` `#020617`; cards `dark-800/50` |
| T3 | Theme toggle while iframe open | sidebar toggle | parent's MutationObserver rebuilds `embeddedUrl` → iframe **reloads** with new `theme=` param (unavoidable without a parent patch); page re-resolves theme at load and restores state from server (`/state`, `PREF_STORE` tab) |
| T4 | Fallbacks | clear `theme` key, toggle OS scheme | page follows `prefers-color-scheme`; `theme=` param overrides OS |
| T5 | Contrast | computed style or WCAG check | body text `#111827`/`#f3f4f6` on `#f9fafb`/`#020617` ≥ 7:1; muted `#6b7280`/`#94a3b8` ≥ 4.5:1 on their backgrounds; badge ink vs soft bg ≥ 4.5:1 both themes |
| T6 | Button states | hover/active/disabled/focus-visible | primary gradient visible both themes; disabled `opacity-50` + `cursor-not-allowed`; focus ring `primary-500/50` visible on both bg |
| T7 | Row states | select, hover, open | selected `accent-soft`; open row `--sunken`; checkbox `accent-color` teal; switch on `primary-500` |
| T8 | Busy/error/verdict | run eval, force failure | `verdict.ok` emerald, `.bad` red, `.warn` amber, `.err/.skip` bordered neutral; spinner `border-current` inherits; `.global-flash[data-tone]` readable both themes |
| T9 | History cards | render + expand | `--history-bg/head-bg` separation visible in dark mode (currently identical to card → verify border compensates) |
| T10 | Markdown/math | expanded candy answer | Inline code uses `--code-*`; blocks use `--pre-*`; tables scroll; Temml remains readable; blockquote borders remain visible in both themes |
| T11 | Pelican | preview + dialog | srcdoc sandbox still renders (`background:#fff` intentional); `.pelican-stage` scales; preview border `--control-line` |
| T12 | Mobile ≤768px | narrow viewport | existing media queries still reflow rows/toolbars/dialogs; no horizontal overflow beyond `.markdown table`/`.fp-detail-scroll` intentional scroll |
| T13 | No credential leak | inspect URL/iframe src, devtools network, logs | iframe URL has `theme`/`lang` only — no `token`, `user_id`, `src_*`, JWT, or eval key; API calls carry `Authorization` header only |
| T14 | Embedding headers | curl/route test | eval UI returns `X-Frame-Options: SAMEORIGIN` and CSP `frame-ancestors 'self'`; all other routes keep `DENY`/`'none'` |
| T15 | `node scripts/test-ui.cjs` equivalent | extend theme/auth environment tests | Test initialization and envelope unwrap explicitly. Add `location.search`, `documentElement`, `matchMedia`, and listener stubs because current initialization is excluded |
| T16 | iFrame reload preserves state | toggle theme mid-run | run state re-fetched from `/state` after reload; persisted history + selected tab intact (per R8); open dialog/scroll loss is accepted native behavior |

## 6. Unresolved Risks and Decisions

1. **`frame-ancestors 'none'` + `X-Frame-Options: DENY` are global** (`config.go:35`, `security_headers.go:131`). The eval route must override per-response; `plugin_handler.go:250-254` proves the pattern. **Recommendation**: route-local `SAMEORIGIN` + `frame-ancestors 'self'`, and use `GetNonceFromContext` for the inline scripts/styles rather than `'unsafe-inline'`. Residual work: the `/*APP_STYLES*/`/`script` assembly in `plugin.go:79-99` must inject the per-request nonce at serve time — the HTML is currently pre-rendered once, so serving becomes per-request (or nonce applied by string substitution on the cached byte string).
2. **Plan badges — recommended resolution**: delete the CPA plan materials; render returned account metadata (name/platform/group) via neutral `badge-gray`/`badge-primary` tokens only. Residual risk: if `credentials.js` tag rendering depends on plan fields, the `/targets` payload must supply safe label strings or the tags render empty — verify at implementation.
3. **Iframe reload on theme toggle is unavoidable without a parent patch** — a storage listener in the child cannot stop `embeddedUrl` recomputation. **Recommendation**: accept reload (native upstream behavior) + server-state restoration; note that persisted tab selection already survives via `PREF_STORE`. Residual loss: open dialogs, scroll position, expanded rows. Reload-free sync is documented as requiring an extra parent-side postMessage/bridge patch — explicitly out of minimal scope.
4. **`data-theme` vs `dark` class — recommended resolution**: keep local `[data-theme="dark"]` as canonical inside the iframe document; the parent's `.dark` class need not be adopted. No residual risk.
5. **Tooltip tokens are invented**: Sub2API has no plain-tooltip component token; `--tip-*` values in §2 are chosen from the nearest analog (dropdown/modal). Flagged as nearest-fit, not source-exact.
6. **Plugin-UI bridge is a dead end for this page** (connect-src 'none', no theme channel) — confirmed so the alternative isn't re-litigated; custom-menu path remains correct per design.
7. **`auth_token` storage verified**: `stores/auth.ts:17` uses plain `localStorage['auth_token']` (`AUTH_TOKEN_KEY`), with `refresh_token` and `token_expires_at` siblings — readable directly from the same-origin iframe. No obfuscation. The page should re-read it per request (rotation) and stop on 401 per `patch-contracts §7`. Resolved — no longer a risk.
8. **Temml/vendor CSS order**: `/*APP_STYLES*/` concatenates temml → themes → styles; keep that order so token overrides win.

## Related Specs

- `.trellis/tasks/10-11-sub2api-runtime-patch/research/patch-contracts.md` §7 — admin page contracts (envelope, embedding, CSP, menu).
- `.trellis/tasks/10-11-sub2api-runtime-patch/research/source-audit.md:88-98` — embedded-url token leak + security headers findings this report extends.
- `.trellis/spec/backend/quality-guidelines.md` — applies to Go handler changes, not CSS.
