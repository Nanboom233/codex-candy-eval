const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const { createRequire } = require("node:module");
const source = process.argv[2] || process.env.SUB2API_SOURCE;
if (!source) throw new Error("Set SUB2API_SOURCE to the patched native source tree.");
const { JSDOM } = createRequire(path.resolve(source, "frontend/package.json"))("jsdom");
const web = path.resolve(source, "backend/internal/codexeval/web");
const read = (name) => fs.readFileSync(path.join(web, name), "utf8");
const settings = () => ({ enabled: false, group_id: 101, account_type: "all", account_plan: "all", enabled_only: false,
  model: "saved/literal", effort: "low", max_tests_per_account: 3,
  min_wait_minutes: 15, max_wait_minutes: 35,
  adjust_priority: false, normal_priority: 10, degraded_priority: 99,
  ntfy_server: "https://ntfy.sh", ntfy_topic: "", ntfy_username: "", ntfy_password: "", bark_url: "",
  ntfy_password_configured: true, bark_url_configured: false });
const idle = { phase: "idle", enabled: false, running: false, results: [] };
const emptyState = () => ({ active: [], candy: {}, fingerprint: {}, modeltrace: {}, pelican: {} });
const deferred = () => { let resolve; const promise = new Promise((done) => { resolve = done; }); return { promise, resolve }; };
const tick = () => new Promise((resolve) => setImmediate(resolve));
function harness(preferences = {}) {
  const dom = new JSDOM(read("ui.html"), { url: "http://127.0.0.1:60198/codex-eval", runScripts: "outside-only" });
  const w = dom.window;
  w.localStorage.setItem("sub2api-codex-eval.prefs", JSON.stringify({ model: "saved/literal" }));
  for (const [key, value] of Object.entries(preferences)) w.localStorage.setItem(key, value);
  w.eval([read("vendor/markdown-it.min.js"), read("vendor/markdown-it-cjk-friendly.min.js"),
    read("vendor/texmath.js"), read("vendor/temml.min.js"), read("components.js"), read("catalog.js"),
    read("candy.js"), read("fingerprint.js").replace('/*FINGERPRINT_CONFIG*/{}', JSON.stringify({ modes: [
      { id: "quick", name: "快速", cells: 1, samples_per_cell: 1 }], default_concurrency: 1, max_concurrency: 32 })),
    read("modeltrace.js").replace('/*MODELTRACE_CONFIG*/{}', JSON.stringify({ requests: 1, default_concurrency: 1, max_concurrency: 32 })),
    read("pelican.js"), read("inspection.js"), read("app.js").split("// Events and initialization.")[0], `
    initializeLayout(); initializeModelPickers(); initializeCandy(); initializeFingerprint(); initializeModelTrace(); initializePelican(); initializeInspection();
    window.test = {fillModels,fillInspectionSettings,loadInspectionSettings,saveInspectionSettings,markInspectionEdit,
      updateModelPicker,renderInspectionState,validateInspectionFields,loadInspectionState,
      refreshCatalog,resetCatalogCache,stateCredentials,visibleCredentials,targetsFor,fillCredentialTypes,fillPlanFilter,renderCandyRow,renderList,render,
      batchCredentials,runCandy,runFingerprint,runModelTrace,runPelican,candyTip,showLogin,load,
      setAccountInspection, selectCandy(id){candySelected.add(id);render()}, setUpdate(fn){update=fn},
      setCredentials(data){credentials=stateCredentials(data)},
      setAPI(fn){api=fn}, setModels(ids){catalogCache.models={ids}}, setGroup(id){selectedGroup=id},
      saved(){return inspectionSettingsLoaded}, dirty(){return inspectionDirtyFields.size}};
  `].join("\n"));
  return { dom, w, h: w.test, $: (id) => w.document.getElementById(id), close: () => w.close() };
}
function choose(ctx, id, value) {
  const select = ctx.$(id + "-select");
  select.value = value;
  select.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
}
const fixtureAccount = (id, name, group_ids, extra = {}) => ({ id, name, group_ids, platform: "openai",
  type: "apikey", status: "active", schedulable: true, priority: 10, credentials: {}, ...extra });
function groupFixture() {
  const accounts = [fixtureAccount(11, "A only", [101]), fixtureAccount(22, "B only", [202]),
    fixtureAccount(33, "shared", [101, 202], { type: "oauth", credentials: { plan_type: "pro" } }),
    fixtureAccount(44, "ungrouped", []), fixtureAccount(55, "disabled", [202], { schedulable: false }),
    fixtureAccount(66, "restricted", [202], { status: "error", error_message: "native restriction" }),
    fixtureAccount(77, "busy", [202])];
  const calls = [];
  const data = emptyState();
  data.active.push({ id: "busy-candy", type: "candy", target: { account_id: 77, group_id: 202 }, progress: { done: 0, total: 1 } });
  return { accounts, calls, data, async api(url) {
    calls.push(url);
    if (url.includes("/groups/all")) return [{ id: 101, name: "A" }, { id: 202, name: "B" }];
    if (url.includes("/accounts?")) {
      const query = new URL(url, "http://fixture").searchParams;
      assert.equal(query.has("group_id"), false);
      const group = query.get("group");
      const selected = group ? accounts.filter((a) => a.group_ids.includes(Number(group))) : accounts;
      const page = Number(query.get("page"));
      return { pages: Math.max(1, selected.length), items: JSON.parse(JSON.stringify(selected.slice(page - 1, page))) };
    }
    if (url.endsWith("/inspection/state")) return { ...idle, phase: "completed", results: [{ account_id: 33,
      account_name: "stale response name", priority: 10, target_priority: 99, tests: 3, outcome: "updated",
      run_ids: ["inspection-run"], error: "evaluation error", notification_error: "notification error" }] };
    if (url.endsWith("/state")) return data;
    return [{ id: "provider/raw" }];
  } };
}

(async () => {
  let ctx = harness();
  try {
    ctx.h.setModels(["provider/raw", "gpt-6.1-sol"]);
    ctx.h.fillModels();
    assert.equal(ctx.$("model-select").value, "saved/literal");
    for (const id of ["model", "fp-model", "mt-model", "pl-model"]) {
      ctx.h.setModels(["provider/raw", "gpt-6.1-sol"]);
      ctx.h.fillModels();
      choose(ctx, id, "provider/raw");
      assert.equal(ctx.$(id).value, "provider/raw");
      const select = ctx.$(id + "-select");
      select.selectedIndex = select.options.length - 1;
      select.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
      assert.equal(ctx.$(id).hidden, false, id);
      ctx.$(id).value = "custom/literal-" + id;
      ctx.$(id).dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
      ctx.h.setModels([]);
      ctx.h.fillModels();
      assert.equal(ctx.$(id).value, "custom/literal-" + id);
    }
    ctx.$("model").value = "";
    ctx.h.fillModels();
    assert.equal(ctx.$("model").value, "", "A refresh must retain an unfinished custom value");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    ctx.h.fillModels();
    const ids = ["model", "fp-model", "mt-model", "pl-model"];
    for (const id of ids) {
      const select = ctx.$(id + "-select");
      select.selectedIndex = select.options.length - 1;
      select.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
      ctx.$(id).value = "typed/" + id;
      ctx.$(id).dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    }
    for (const [id, value] of [["runs", "7"], ["fp-concurrency", "3"], ["mt-concurrency", "2"], ["effort", "high"], ["pl-effort", "max"]]) {
      ctx.$(id).value = value;
      ctx.$(id).dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    }
    const preferences = Object.fromEntries(Object.entries(ctx.w.localStorage));
    const reloaded = harness(preferences);
    try {
      reloaded.h.fillModels();
      for (const id of ids) assert.equal(reloaded.$(id).value, "typed/" + id, "Reload before blur lost " + id);
      for (const [id, value] of [["runs", "7"], ["fp-concurrency", "3"], ["mt-concurrency", "2"], ["effort", "high"], ["pl-effort", "max"]]) {
        assert.equal(reloaded.$(id).value, value, "Reload before blur lost " + id);
      }
    } finally { reloaded.close(); }
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const saveResponse = deferred();
    let persisted = settings();
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.includes("/accounts?")) return { pages: 1, items: [] };
      if (url.endsWith("/inspection/state")) return idle;
      if (url.endsWith("/state")) return emptyState();
      if (options.method === "PUT") {
        persisted = { ...options.body, ntfy_password: "", bark_url: "", ntfy_password_configured: true };
        return saveResponse.promise;
      }
      return { ...persisted };
    });
    await ctx.h.loadInspectionSettings();
    const form = ctx.$("inspection-settings");
    for (const [name, value] of [["ntfy_topic", "saved-topic"], ["ntfy_password", "fixture-secret"]]) {
      const field = form.elements.namedItem(name);
      field.value = value;
      field.dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    }
    form.dispatchEvent(new ctx.w.Event("submit", { bubbles: true, cancelable: true }));
    await tick();
    const wait = form.elements.namedItem("min_wait_minutes");
    wait.value = "22";
    wait.dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    saveResponse.resolve({ ...persisted });
    await tick();
    assert.equal(wait.value, "22", "A save response must preserve a later edit to another field");
    assert.equal(form.elements.namedItem("ntfy_password").value, "", "A saved secret must clear even when another field changes");
    assert.equal(ctx.$("inspection-save-status").textContent, "");
    persisted.ntfy_topic = "saved-in-another-session";
    await ctx.h.loadInspectionSettings();
    assert.equal(form.elements.namedItem("ntfy_topic").value, persisted.ntfy_topic, "A submitted field must no longer remain dirty");
    assert.equal(wait.value, "22");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const oldRead = deferred();
    let pendingRead = false;
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.includes("/accounts?")) return { pages: 1, items: [] };
      if (url.endsWith("/inspection/state")) return idle;
      if (url.endsWith("/state")) return emptyState();
      if (options.method === "PUT") return { ...options.body };
      return pendingRead ? oldRead.promise : settings();
    });
    await ctx.h.loadInspectionSettings();
    pendingRead = true;
    const loading = ctx.h.loadInspectionSettings();
    await tick();
    ctx.$("inspection-model").value = "newly/saved";
    ctx.$("inspection-model").dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    ctx.$("inspection-settings").dispatchEvent(new ctx.w.Event("submit", { bubbles: true, cancelable: true }));
    await tick();
    oldRead.resolve(settings());
    await loading;
    assert.equal(ctx.$("inspection-model").value, "newly/saved", "An older GET must not replace a completed save");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    let persisted = settings();
    ctx.h.setAPI(async (url) => {
      if (url.includes("/accounts?")) return { pages: 1, items: [] };
      if (url.endsWith("/inspection/state")) return idle;
      if (url.endsWith("/state")) return emptyState();
      return { ...persisted };
    });
    await ctx.h.loadInspectionSettings();
    persisted = { ...persisted, model: "saved/other-session", min_wait_minutes: 21, max_wait_minutes: 32 };
    await ctx.h.loadInspectionSettings();
    assert.equal(ctx.$("inspection-model").value, persisted.model, "Reopening settings must read saved native values");
    assert.equal(ctx.$("inspection-settings").elements.namedItem("min_wait_minutes").value, "21");
    ctx.$("inspection-model").value = "newer/unsaved";
    ctx.$("inspection-model").dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    persisted = { ...persisted, model: "saved/newer-session", max_wait_minutes: 33 };
    await ctx.h.loadInspectionSettings();
    assert.equal(ctx.$("inspection-model").value, "newer/unsaved", "A reread must preserve unsaved fields");
    assert.equal(ctx.$("inspection-settings").elements.namedItem("max_wait_minutes").value, "33");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const invalid = { ...settings(), degraded_priority: 5 };
    let writes = 0;
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.endsWith("/inspection/state")) return idle;
      if (options.method === "PUT") { writes++; throw new Error("priority pair must be corrected"); }
      return invalid;
    });
    await ctx.h.loadInspectionSettings();
    const form = ctx.$("inspection-settings");
    const adjustment = form.elements.namedItem("adjust_priority");
    const degraded = form.elements.namedItem("degraded_priority");
    assert.equal(ctx.$("inspection-priority-fields").hidden, true);
    ctx.h.validateInspectionFields();
    assert.notEqual(degraded.validationMessage, "", "Autosave must reject an invalid collapsed priority pair");
    adjustment.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
    form.dispatchEvent(new ctx.w.Event("submit", { bubbles: true, cancelable: true }));
    await tick();
    assert.equal(writes, 0, "An invalid pair must not be saved");
    assert.equal(ctx.$("inspection-save-status").textContent, "降级优先级须大于正常优先级");
    adjustment.checked = true;
    adjustment.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
    assert.equal(ctx.$("inspection-priority-fields").hidden, false);
    assert.notEqual(degraded.validationMessage, "", "enabling adjustment must validate the visible priority pair");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    ctx.h.fillInspectionSettings(settings());
    const form = ctx.$("inspection-settings");
    const adjustment = form.elements.namedItem("adjust_priority");
    const normal = form.elements.namedItem("normal_priority");
    assert.equal(ctx.$("inspection-priority-fields").hidden, true);
    assert.equal(normal.value, "10");
    adjustment.checked = true;
    adjustment.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
    assert.equal(ctx.$("inspection-priority-fields").hidden, false);
    normal.value = "0";
    assert.equal(normal.checkValidity(), false);
    normal.value = "1";
    assert.equal(normal.checkValidity(), true);
    adjustment.checked = false;
    adjustment.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
    assert.equal(ctx.$("inspection-priority-fields").hidden, true);
    assert.equal(normal.value, "1");
    assert.equal(form.elements.namedItem("include_patterns"), null);
    assert.equal(form.elements.namedItem("exclude_patterns"), null);
    assert.equal(ctx.$("inspection-ntfy-configured"), null);
    assert.equal(ctx.$("inspection-bark-configured"), null);
    const header = ctx.$("candy-credentials").querySelector(".credential-head");
    assert.equal(ctx.$("inspection-runtime"), null);
    assert.equal(ctx.$("inspection-status").parentElement, header);
    assert.equal(ctx.$("inspection-runtime-error"), null);
    assert.equal(ctx.$("candy-settings-status"), null);
    for (const prefix of ["fp-", "mt-", "pl-"]) assert.equal(ctx.$(prefix + "rows").parentElement.querySelector(".inspection-status"), null);
    ctx.h.renderInspectionState(idle);
    assert.equal(ctx.$("inspection-status").hidden, true);
    assert.equal(ctx.$("inspection-phase").textContent, "");
    assert.equal(ctx.$("inspection-next").textContent, "");
    ctx.h.renderInspectionState({ ...idle, enabled: true, next_run_at: "2026-10-17T12:00:00Z" });
    assert.equal(ctx.$("inspection-status").hidden, false);
    assert.ok(ctx.$("inspection-next").textContent.startsWith("下次巡检："));
    ctx.h.renderInspectionState({ ...idle, running: true, phase: "testing", started_at: "2026-10-17T12:00:00Z", done: 1, total: 2 });
    assert.equal(ctx.$("inspection-phase").textContent, "正在测试");
    assert.equal(ctx.$("inspection-progress").textContent, "账户进度：1/2");
    ctx.h.renderInspectionState({ ...idle, phase: "failed", error: "actual fixture failure" });
    assert.equal(ctx.$("inspection-runtime-error"), null);
    assert.equal(ctx.$("inspection-phase").textContent, "");
    assert.equal(ctx.$("settings-panel").textContent.includes("actual fixture failure"), false);
    assert.equal(header.hidden, false);
    ctx.h.renderInspectionState({ ...idle, phase: "completed", started_at: "2026-10-17T12:00:00Z", done: 2, total: 2 });
    assert.equal(ctx.$("inspection-phase").textContent, "");
    assert.equal(ctx.$("inspection-progress").textContent, "");
    assert.equal(ctx.$("inspection-status").hidden, true);
    assert.equal(header.hidden, false);
    ctx.h.renderInspectionState({ ...idle, phase: "completed", next_run_at: "2026-10-17T12:30:00Z" });
    assert.equal(ctx.$("inspection-status").hidden, false);
    assert.ok(ctx.$("inspection-next").textContent.startsWith("下次巡检："));
  } finally { ctx.close(); }

  ctx = harness();
  try {
    let state = idle;
    let saved;
    const actions = [];
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.includes("/accounts?")) return { pages: 1, items: [] };
      if (url.endsWith("/inspection/state")) return state;
      if (url.endsWith("/state")) return emptyState();
      if (options.method === "PUT") { saved = options.body; return saved; }
      if (options.method === "POST") {
        actions.push(url);
        state = url.endsWith("/start") ? { ...idle, running: true, phase: "testing",
          started_at: "2026-10-17T12:00:00Z", done: 1, total: 2,
          results: [{ account_name: "fixture", priority: 10, tests: 1, outcome: "testing" }] }
          : { ...state, running: false, phase: "cancelled" };
        return state;
      }
      return settings();
    });
    await ctx.h.loadInspectionSettings();
    const form = ctx.$("inspection-settings");
    const enabled = ctx.$("inspection-enabled");
    const fields = ctx.$("inspection-auto-fields");
    assert.equal(fields.hidden, true);
    for (const input of [enabled, form.elements.namedItem("adjust_priority")]) {
      assert.ok(input.getAttribute("aria-label"));
      assert.equal(input.closest("label").textContent.trim(), "", "Header switches must have no visible captions");
    }
    for (const [id, panel] of [["inspection-start", "settings"], ["inspection-stop", "settings"], ["inspection-status", "candy"]]) {
      assert.equal(fields.contains(ctx.$(id)), false, id + " must remain outside collapsed configuration");
      assert.equal(ctx.$(panel + "-panel").contains(ctx.$(id)), true, id + " belongs to " + panel);
    }
    enabled.checked = true;
    enabled.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
    assert.equal(fields.hidden, false);
    ctx.h.updateModelPicker("inspection-model", [], "custom/collapsed-model");
    form.elements.namedItem("min_wait_minutes").value = "22";
    form.elements.namedItem("max_wait_minutes").value = "31";
    form.elements.namedItem("max_tests_per_account").value = "2";
    enabled.checked = false;
    enabled.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
    assert.equal(fields.hidden, true);
    form.dispatchEvent(new ctx.w.Event("submit", { bubbles: true, cancelable: true }));
    await tick();
    assert.equal(saved.enabled, false);
    assert.equal(saved.model, "custom/collapsed-model");
    assert.equal(saved.min_wait_minutes, 22);
    assert.equal(saved.max_wait_minutes, 31);
    assert.equal(saved.max_tests_per_account, 2);
    assert.equal(fields.hidden, true);
    assert.equal(ctx.$("inspection-start").disabled, false);
    ctx.$("inspection-start").click();
    await tick();
    assert.equal(ctx.$("inspection-stop").disabled, false);
    assert.equal(ctx.$("inspection-phase").textContent, "正在测试");
    assert.equal(ctx.$("inspection-progress").textContent, "账户进度：1/2");
    assert.equal(ctx.$("inspection-results"), null);
    assert.equal(fields.hidden, true);
    ctx.$("inspection-stop").click();
    await tick();
    assert.equal(ctx.$("inspection-phase").textContent, "巡检已停止");
    assert.deepEqual(actions.map((url) => url.split("/").pop()), ["start", "stop"]);
    enabled.checked = true;
    enabled.dispatchEvent(new ctx.w.Event("change", { bubbles: true }));
    assert.equal(fields.hidden, false);
    assert.equal(form.elements.namedItem("min_wait_minutes").value, "22");
    assert.equal(ctx.$("inspection-model").value, "custom/collapsed-model");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const firstLoad = deferred();
    const saveResponse = deferred();
    let saved;
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.includes("/accounts?")) return { pages: 1, items: [{ id: 1, name: "first" }] };
      if (url.endsWith("/models")) throw new Error("active catalog error");
      if (url.endsWith("/inspection/state")) return idle;
      if (url.endsWith("/state")) return emptyState();
      if (options.method === "PUT") { saved = options.body; return saveResponse.promise; }
      return firstLoad.promise;
    });
    const loading = ctx.h.loadInspectionSettings();
    ctx.$("inspection-model").value = "edited/before-load";
    ctx.$("inspection-model").dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    firstLoad.resolve(settings());
    await loading;
    assert.equal(ctx.$("inspection-model").value, "edited/before-load");
    assert.equal(ctx.$("inspection-catalog-error"), null);
    assert.equal(ctx.h.saved(), true, "Catalog failure must not block a valid saved/custom model");
    ctx.h.updateModelPicker("inspection-model", [], "custom/exact-model");
    ctx.$("inspection-settings").dispatchEvent(new ctx.w.Event("submit", { bubbles: true, cancelable: true }));
    await tick();
    assert.equal(saved.model, "custom/exact-model");
    assert.ok(!("include_patterns" in saved) && !("exclude_patterns" in saved));
    assert.equal(saved.ntfy_password, "");
    assert.equal(saved.bark_url, "");
    ctx.$("inspection-model").value = "newer/unsaved-edit";
    ctx.$("inspection-model").dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    saveResponse.resolve({ ...saved, ntfy_password_configured: true });
    await tick();
    assert.equal(ctx.$("inspection-model").value, "newer/unsaved-edit");
    assert.equal(ctx.$("inspection-save-status").textContent, "");
    ctx.h.fillInspectionSettings(saved);
    assert.equal(ctx.$("inspection-model-select").value, "custom/exact-model");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const fixture = groupFixture();
    const prefixes = ["", "fp-", "mt-", "pl-"];
    const submissions = [];
    ctx.h.setAPI(fixture.api);
    ctx.h.setUpdate((_path, options) => submissions.push(options.body));
    ctx.h.setGroup(101);
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(fixture.data);
    ctx.h.fillCredentialTypes();
    ctx.h.render();
    const ids = (prefix) => [...ctx.$(prefix + "rows").querySelectorAll("[data-row]")].map((row) => Number(row.dataset.row)).sort((a, b) => a - b);
    for (const prefix of prefixes) assert.deepEqual(ids(prefix), [11, 33], prefix);
    for (const [prefix, run] of prefixes.map((prefix, i) => [prefix, [ctx.h.runCandy, ctx.h.runFingerprint, ctx.h.runModelTrace, ctx.h.runPelican][i]])) {
      const targets = ctx.h.targetsFor(ctx.h.batchCredentials(prefix, new Set()).map((a) => a.id));
      await run({ targets });
    }
    assert.deepEqual(submissions.map((body) => [...body.targets].map((target) => target.account_id).sort((a, b) => a - b)),
      [[11, 33], [11, 33], [11, 33], [11, 33]]);
    assert.ok(submissions.every((body) => body.targets.every((target) => target.group_id === 101)));
    const modelRequests = fixture.calls.filter((url) => url.endsWith("/models")).length;
    await ctx.h.refreshCatalog();
    assert.equal(fixture.calls.filter((url) => url.endsWith("/models")).length, modelRequests, "Unchanged inventory must reuse raw models");
    fixture.accounts.find((a) => a.id === 11).group_ids = [202];
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(fixture.data);
    ctx.h.render();
    for (const prefix of prefixes) assert.deepEqual(ids(prefix), [33], "An ordinary poll must remove moved accounts: " + prefix);
    assert.deepEqual([...ctx.h.targetsFor(["11", "33"])].map((target) => target.account_id), [33]);
    fixture.accounts.find((a) => a.id === 11).group_ids = [101];
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(fixture.data);
    await ctx.h.loadInspectionState();
    ctx.h.render();
    for (const prefix of prefixes) assert.deepEqual(ids(prefix), [11, 33], "Inspection polling must not replace inventory: " + prefix);
    ctx.h.setGroup(202);
    assert.equal(ctx.h.targetsFor(["11"]).length, 0, "A stale row cannot produce a target for a new group");
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(fixture.data);
    ctx.h.fillCredentialTypes();
    ctx.h.render();
    for (const prefix of prefixes) {
      assert.deepEqual(ids(prefix), [22, 33, 55, 66, 77]);
      assert.deepEqual([...ctx.h.batchCredentials(prefix, new Set())].map((a) => Number(a.id)).sort((a, b) => a - b), [22, 33]);
    }
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const fixture = groupFixture();
    const stale = deferred();
    let held = false;
    ctx.h.setAPI(async (url) => {
      const query = new URL(url, "http://fixture").searchParams;
      if (url.includes("/accounts?") && query.get("group") === "101" && query.get("page") === "2") {
        held = true;
        return stale.promise;
      }
      return fixture.api(url);
    });
    ctx.h.setGroup(101);
    const old = ctx.h.refreshCatalog();
    await tick();
    assert.equal(held, true);
    ctx.h.setGroup(202);
    await ctx.h.refreshCatalog();
    stale.resolve({ pages: 2, items: fixture.accounts.filter((a) => a.id === 33) });
    await old;
    ctx.h.setCredentials(fixture.data);
    ctx.h.fillCredentialTypes();
    for (const prefix of ["", "fp-", "mt-", "pl-"]) {
      ctx.h.fillPlanFilter(prefix);
      assert.ok([...ctx.h.visibleCredentials(prefix)].every((a) => a.group_ids.includes(202)), "Late group A responses cannot commit into B");
    }
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const fixture = groupFixture();
    fixture.data.candy["33"] = [{ run_id: "inspection-run:1", time: "2026-10-17T12:00:00Z", model: "saved/literal",
      effort: "low", ok: true, answer: "**shared Candy answer**", duration_ms: 20, input_tokens: 10, output_tokens: 4 }];
    ctx.h.setAPI(fixture.api);
    ctx.h.setGroup(101);
    ctx.h.fillInspectionSettings(settings());
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(fixture.data);
    ctx.h.render();
    await ctx.h.loadInspectionState();
    const list = ctx.$("rows");
    const row = list.querySelector('[data-row="33"]').closest(".row");
    assert.equal(row.querySelector(".name").textContent, "shared", "Native identity must win over the result snapshot name");
    assert.equal(row.querySelector(".row-summary"), null);
    assert.equal(ctx.$("inspection-runtime-error"), null);
    assert.equal(ctx.$("candy-panel").textContent.includes("notification error"), false);
    assert.equal(ctx.$("settings-panel").textContent.includes("notification error"), false);
    assert.ok(list.querySelector("[data-candy-run]"));
    assert.equal(ctx.$("inspection-clear"), null);
    assert.equal(ctx.$("inspection-select-all"), null);
    const mark = row.querySelector("[data-auth]");
    assert.ok(ctx.h.candyTip(mark).includes("<strong>shared Candy answer</strong>"));
    row.querySelector("[data-row]").click();
    const openRow = list.querySelector('[data-row="33"]').closest(".row");
    assert.ok(openRow.querySelector(".history-card .answer-content").innerHTML.includes("<strong>shared Candy answer</strong>"));
    const answer = openRow.querySelector("[data-answer]");
    answer.click();
    assert.equal(answer.getAttribute("aria-expanded"), "true");
    const node = openRow.querySelector("[data-row]");
    ctx.h.renderInspectionState({ ...idle, results: [{ account_id: 33, priority: 10, tests: 2, outcome: "testing" }] });
    assert.equal(list.querySelector('[data-row="33"]'), node, "Shared list patching must keep the existing row node");
    ctx.h.fillCredentialTypes();
    ctx.$("credential-type").value = "oauth";
    ctx.h.render();
    assert.deepEqual([...list.querySelectorAll("[data-row]")].map((entry) => Number(entry.dataset.row)), [33]);
    ctx.h.showLogin("fixture session ended");
    assert.equal(ctx.$("inspection-results"), null);
    assert.equal(ctx.h.candyTip(mark), "", "Session reset must discard inspection records");
  } finally { ctx.close(); }
  ctx = harness();
  try {
    const fixture = groupFixture();
    fixture.data.candy["33"] = [{ time: "2026-10-17T12:00:00Z", model: "saved/literal", error: "original test error" }];
    const delayedInventory = deferred();
    const save = deferred();
    let delay = false, failSave = false;
    const writes = [];
    ctx.h.setGroup(101);
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.endsWith("/codex-eval-inspection") && options.method === "PUT") {
        writes.push(options.body);
        if (failSave) throw new Error("fixture save failure");
        await save.promise;
        fixture.accounts.find((a) => a.id === 33).extra = { codex_eval_inspection_enabled: options.body.enabled };
        return { account_id: 33, enabled: options.body.enabled };
      }
      if (delay && url.includes("/accounts?")) await delayedInventory.promise;
      return fixture.api(url);
    });
    await ctx.h.load();
    ctx.h.selectCandy("33");
    ctx.$("rows").querySelector('[data-row="33"]').click();
    const row = ctx.$("rows").querySelector('[data-row="33"]');
    const button = row.querySelector("[data-candy-run]");
    const before = button.outerHTML;
    assert.equal(row.querySelector("[data-candy-select]").checked, true);
    assert.ok(row.closest(".row").querySelector(".answer-content").textContent.includes("original test error"));
    delay = true;
    const refreshing = ctx.h.load();
    await tick();
    assert.equal(ctx.$("rows").querySelector('[data-row="33"]'), row);
    assert.equal(row.querySelector("[data-candy-run]").outerHTML, before);
    assert.equal(row.querySelector("[data-candy-select]").checked, true);
    assert.ok(row.closest(".row").classList.contains("open"));
    delay = false;
    delayedInventory.resolve();
    await refreshing;
    assert.equal(ctx.$("rows").querySelector('[data-row="33"]'), row);
    assert.equal(row.querySelector("[data-candy-run]").outerHTML, before);
    assert.equal(row.querySelector("[data-candy-select]").checked, true);
    assert.ok(row.closest(".row").classList.contains("open"));
    assert.ok(ctx.$("candy-credentials").querySelector(".list-head").textContent.includes("启用自动巡检"));
    for (const prefix of ["fp-", "mt-", "pl-"]) assert.equal(ctx.$(prefix + "rows").querySelector("[data-inspection-enable]"), null);
    assert.equal(row.querySelector("[data-enable]").getAttribute("aria-label"), "启用此账户");
    const saving = ctx.h.setAccountInspection("33", false);
    await tick();
    await ctx.h.load();
    assert.equal(row.querySelector("[data-inspection-enable]").checked, false);
    assert.equal(row.querySelector("[data-inspection-enable]").disabled, true);
    assert.equal(row.querySelector("[data-enable]").checked, true);
    assert.equal(ctx.$("rows").querySelector('[data-inspection-enable="11"]').checked, true);
    assert.equal(ctx.$("inspection-enabled").checked, false);
    save.resolve();
    await saving;
    assert.deepEqual(JSON.parse(JSON.stringify(writes)), [{ enabled: false }]);
    assert.equal(row.querySelector("[data-inspection-enable]").checked, false);
    assert.equal(row.querySelector("[data-inspection-enable]").disabled, false);
    failSave = true;
    await ctx.h.setAccountInspection("33", true);
    assert.equal(row.querySelector("[data-inspection-enable]").checked, false);
    assert.equal(row.querySelector("[data-inspection-enable]").disabled, false);
    assert.ok(row.querySelector(".inspection-cell [role=alert]").textContent.includes("fixture save failure"));
    assert.equal(ctx.$("inspection-save-status").textContent, "");
    assert.equal(ctx.$("inspection-runtime-error"), null);
    ctx.h.setAPI(async () => { throw new Error("hidden runtime read failure"); });
    await ctx.h.loadInspectionState();
    assert.equal(ctx.$("candy-panel").textContent.includes("hidden runtime read failure"), false);
    assert.equal(ctx.$("settings-panel").textContent.includes("hidden runtime read failure"), false);
  } finally { ctx.close(); }
  ctx = harness();
  try {
    const fixture = groupFixture();
    let inventoryFailed = false;
    let modelFailed = false;
    const starts = [];
    ctx.h.setGroup(101);
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.includes("/accounts?") && inventoryFailed) throw new Error("fixture inventory offline");
      if (url.endsWith("/models") && modelFailed) throw new Error("fixture model error");
      if (url.endsWith("/start")) { starts.push(options.body); return { started: [], busy: [], unavailable: [] }; }
      return fixture.api(url);
    });
    await ctx.h.load();
    const oldTargets = ctx.h.targetsFor(["11"]);
    assert.equal(oldTargets.length, 1);
    fixture.accounts.find((a) => a.id === 11).group_ids = [202];
    inventoryFailed = true;
    await ctx.h.load();
    assert.equal(ctx.$("catalog-error").hidden, false);
    assert.ok(ctx.$("catalog-error").textContent.includes("fixture inventory offline"));
    assert.equal(ctx.h.targetsFor(["11", "33"]).length, 0);
    for (const prefix of ["", "fp-", "mt-", "pl-"]) {
      assert.equal(ctx.$(prefix + "run-batch").disabled, true);
      assert.equal(ctx.$(prefix + "rows").querySelector("[data-row] > button").disabled, true);
    }
    for (const run of [ctx.h.runCandy, ctx.h.runFingerprint, ctx.h.runModelTrace, ctx.h.runPelican]) await run({ targets: oldTargets });
    assert.equal(starts.length, 0, "Neither row nor batch helpers may submit targets from a failed inventory refresh");
    inventoryFailed = false;
    await ctx.h.load();
    assert.equal(ctx.$("catalog-error").hidden, true);
    assert.deepEqual([...ctx.h.targetsFor(["11", "33"])].map((target) => target.account_id), [33]);
    assert.equal(ctx.$("run-batch").disabled, false);
    ctx.$("rows").querySelector("[data-candy-run]").click();
    await tick();
    assert.equal(starts.length, 1, "A successful inventory refresh must restore row starts");
    modelFailed = true;
    await ctx.h.load({ refresh: true });
    assert.equal(ctx.$("catalog-error").hidden, false);
    assert.ok(ctx.$("catalog-error").textContent.includes("fixture model error"));
    assert.equal(ctx.$("run-batch").disabled, false, "Raw-model discovery failure must not invalidate fresh account membership");
    await ctx.h.runCandy({ targets: ctx.h.targetsFor(["33"]) });
    assert.equal(starts.length, 2, "A literal model remains usable with a fresh inventory");
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const fixture = groupFixture();
    fixture.data.candy["33"] = [{ run_id: "retained:1", time: "2026-10-17T12:00:00Z", model: "saved/literal",
      effort: "low", ok: true, answer: "retained Candy record", duration_ms: 20, input_tokens: 10, output_tokens: 4 }];
    let stateFailed = false;
    const starts = [];
    ctx.h.setGroup(101);
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.endsWith("/state") && stateFailed) throw new Error("fixture evaluation state offline");
      if (url.endsWith("/start")) { starts.push(options.body); return { started: [], busy: [], unavailable: [] }; }
      return fixture.api(url);
    });
    await ctx.h.load();
    fixture.accounts.find((a) => a.id === 11).group_ids = [202];
    Object.assign(fixture.accounts.find((a) => a.id === 33), { name: "shared updated", priority: 99 });
    stateFailed = true;
    await ctx.h.load();
    assert.equal(ctx.$("load-error").hidden, false);
    assert.ok(ctx.$("load-error").textContent.includes("fixture evaluation state offline"));
    for (const prefix of ["", "fp-", "mt-", "pl-"]) {
      const list = ctx.$(prefix + "rows");
      assert.deepEqual([...list.querySelectorAll("[data-row]")].map((row) => Number(row.dataset.row)), [33]);
      assert.equal(list.querySelector(".name").textContent, "shared updated");
      assert.ok(list.textContent.includes("P99"));
      assert.deepEqual([...ctx.h.batchCredentials(prefix, new Set())].map((a) => Number(a.id)), [33]);
    }
    assert.deepEqual([...ctx.h.targetsFor(["11", "33"])].map((target) => target.account_id), [33]);
    ctx.$("rows").querySelector(".toggle").click();
    assert.ok(ctx.$("rows").querySelector(".answer-content").textContent.includes("retained Candy record"));
    ctx.$("rows").querySelector("[data-candy-run]").click();
    await tick();
    for (const [prefix, run] of [["", ctx.h.runCandy], ["fp-", ctx.h.runFingerprint], ["mt-", ctx.h.runModelTrace], ["pl-", ctx.h.runPelican]]) {
      await run({ targets: ctx.h.targetsFor(ctx.h.batchCredentials(prefix, new Set()).map((a) => a.id)) });
    }
    assert.equal(starts.length, 5);
    assert.ok(starts.every((body) => body.targets.length === 1 && body.targets[0].account_id === 33 && body.targets[0].group_id === 101));
    assert.equal(ctx.$("load-error").hidden, false, "Starting a valid current account must not hide the record-fetch error");
    stateFailed = false;
    await ctx.h.load();
    assert.equal(ctx.$("load-error").hidden, true);
  } finally { ctx.close(); }

  ctx = harness();
  try {
    const fixture = groupFixture();
    const writes = [];
    const firstSave = deferred();
    let activeWrites = 0;
    let maximumWrites = 0;
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.endsWith("/inspection/state")) return idle;
      if (options.method === "PUT") {
        activeWrites++;
        maximumWrites = Math.max(maximumWrites, activeWrites);
        writes.push(options.body);
        if (writes.length === 1) await firstSave.promise;
        activeWrites--;
        return { ...options.body, ntfy_password: "", bark_url: "" };
      }
      if (url.includes("/settings/")) return { ...settings(), account_type: "oauth", account_plan: "Pro ALL", enabled_only: true };
      return fixture.api(url);
    });
    await ctx.h.loadInspectionSettings();
    assert.equal(ctx.$("inspection-save"), null);
    assert.equal(ctx.$("inspection-results"), null);
    assert.equal(ctx.$("inspection-settings").elements.namedItem("skip_free"), null);
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(fixture.data);
    ctx.h.render();
    assert.deepEqual([...ctx.h.batchCredentials("", new Set())].map((a) => Number(a.id)), [33]);
    assert.equal(ctx.$("credential-plan").value, "Pro ALL", "A saved filter must survive a group with one Pro tier");
    const before = fixture.calls.length;
    await ctx.h.loadInspectionState();
    assert.equal(fixture.calls.length, before, "Inspection state must not fetch another account/model catalog");
    const topic = ctx.$("inspection-settings").elements.namedItem("ntfy_topic");
    for (const value of ["a", "ab", "abc"]) {
      topic.value = value;
      topic.dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    }
    await new Promise((resolve) => setTimeout(resolve, 450));
    assert.equal(writes.length, 1, "Input changes must debounce without Submit");
    assert.equal(ctx.$("inspection-save-status").textContent, "");
    assert.equal(writes[0].group_id, 101);
    assert.equal(writes[0].account_type, "oauth");
    assert.equal(writes[0].account_plan, "Pro ALL");
    assert.equal(writes[0].enabled_only, true);
    topic.value = "newer";
    topic.dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    const pendingSave = ctx.h.saveInspectionSettings();
    firstSave.resolve();
    await pendingSave;
    assert.equal(maximumWrites, 1, "Native settings writes must be serialized");
    assert.equal(writes.length, 2);
    assert.equal(writes[1].ntfy_topic, "newer");
    assert.equal(ctx.h.dirty(), 0);
    assert.equal(ctx.$("inspection-save-status").textContent, "");
    assert.equal(ctx.$("candy-settings-status"), null);
  } finally { ctx.close(); }

  ctx = harness();
  try {
    let writes = 0;
    let failing = true;
    let lastBody;
    const fixture = groupFixture();
    ctx.h.setAPI(async (url, options = {}) => {
      if (url.endsWith("/inspection/state")) return idle;
      if (options.method === "PUT") {
        writes++;
        lastBody = options.body;
        if (failing) throw new Error("fixture database failure");
        return { ...options.body };
      }
      if (url.includes("/settings/")) return { ...settings(), group_id: 999, account_type: "oauth", account_plan: "Free" };
      return fixture.api(url);
    });
    await ctx.h.loadInspectionSettings();
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(emptyState());
    ctx.h.render();
    assert.equal(ctx.$("group").value, "999", "An absent saved group must not switch to a broader group");
    assert.equal(ctx.$("credential-plan").value, "Free", "An empty inventory must not reset a saved subscription filter");
    assert.equal(ctx.h.batchCredentials("", new Set()).length, 0);
    const topic = ctx.$("inspection-settings").elements.namedItem("ntfy_topic");
    topic.value = "retain-after-failure";
    topic.dispatchEvent(new ctx.w.Event("input", { bubbles: true }));
    assert.equal(await ctx.h.saveInspectionSettings(), false);
    assert.equal(ctx.$("inspection-save-status").textContent, "fixture database failure");
    assert.equal(ctx.h.dirty(), 1);
    assert.equal(topic.value, "retain-after-failure");
    failing = false;
    await ctx.h.saveInspectionSettings();
    assert.equal(writes, 2);
    assert.equal(lastBody.ntfy_topic, "retain-after-failure");
    assert.equal(ctx.h.dirty(), 0);
  } finally { ctx.close(); }
  ctx = harness();
  try {
    const fixture = groupFixture();
    fixture.accounts.push(fixtureAccount(88, "shadow", [101], { type: "oauth", parent_plan_type: "pro" }));
    ctx.h.setAPI(async (url) => url.includes("/settings/")
      ? { ...settings(), account_type: "oauth", account_plan: "Pro ALL" } : fixture.api(url));
    await ctx.h.loadInspectionSettings();
    await ctx.h.refreshCatalog();
    ctx.h.setCredentials(fixture.data);
    ctx.h.render();
    assert.deepEqual([...ctx.h.batchCredentials("", new Set())].map((a) => Number(a.id)).sort((a, b) => a - b), [33, 88],
      "Candy Pro ALL must include the shadow's native parent plan");
    ctx.$("credential-plan").value = "Pro 200";
    assert.deepEqual([...ctx.h.batchCredentials("", new Set())].map((a) => Number(a.id)).sort((a, b) => a - b), [33, 88],
      "Candy's specific plan filter must include the same shadow account");
  } finally { ctx.close(); }
  console.log("Settings and account checks passed: autosave, native shared selection, manual actions, inventory failure protections, Candy history, and independent inspection controls.");
})().catch((error) => { console.error(error); process.exitCode = 1; });
