const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');

const source = process.argv[2];
if (!source) throw new Error('Specify the patched research source.');
const testPath = path.resolve(__dirname, '../../../../tests/ui/sub2api-dropdowns.cjs');
const prefix = fs.readFileSync(testPath, 'utf8').split('(async () => {')[0]
  .replace('setUpdate(fn){update=fn},', 'setUpdate(fn){update=fn}, setCandySelection(ids){candySelected.clear();ids.forEach(id=>candySelected.add(id))}, setCandyExpanded(ids){ids.forEach(id=>candyExpanded.add(id))},');
const sandbox = { require, process: { argv: ['node', testPath, source], env: process.env }, setImmediate, URL };
vm.runInNewContext(prefix + '\nthis.helpers={harness,groupFixture,deferred,tick};', sandbox, { filename: testPath });
const { harness, groupFixture, deferred, tick } = sandbox.helpers;

(async () => {
  const ctx = harness();
  const fixture = groupFixture();
  try {
    ctx.h.setGroup(101);
    ctx.h.setAPI(fixture.api);
    await ctx.h.load();
    assert.ok(ctx.$('rows').querySelector('[data-row="11"]'), JSON.stringify({ rows: ctx.$('rows').textContent,
      catalogError: ctx.$('catalog-error')?.textContent, loadError: ctx.$('load-error')?.textContent, calls: fixture.calls }));
    ctx.h.setCandySelection(['11']);
    ctx.h.setCandyExpanded(['11']);
    ctx.h.render();
    const row = ctx.$('rows').querySelector('[data-row="11"]').parentElement;
    const button = row.querySelector('[data-candy-run="11"]');
    const checkbox = row.querySelector('[data-candy-select="11"]');
    checkbox.focus();
    const before = { disabled: button.disabled, selected: checkbox.checked, expanded: row.classList.contains('open') };
    const gate = deferred();
    let gated = false;
    ctx.h.setAPI(async (url, options) => {
      if (!gated && url.includes('/accounts?')) {
        gated = true;
        await gate.promise;
      }
      return fixture.api(url, options);
    });
    const refresh = ctx.h.load();
    await tick();
    const during = { disabled: button.disabled, selected: checkbox.checked, expanded: row.classList.contains('open'),
      sameRow: row === ctx.$('rows').querySelector('[data-row="11"]').parentElement,
      reason: button.getAttribute('data-tip') };
    gate.resolve();
    await refresh;
    const after = { disabled: button.disabled, selected: checkbox.checked, expanded: row.classList.contains('open'),
      sameRow: row === ctx.$('rows').querySelector('[data-row="11"]').parentElement };
    assert.equal(before.disabled, false);
    assert.equal(during.disabled, true);
    assert.equal(after.disabled, false);
    assert.equal(during.sameRow && after.sameRow, true);
    assert.equal(before.selected && during.selected && after.selected, true);
    assert.equal(before.expanded && during.expanded && after.expanded, true);
    ctx.h.renderInspectionState({ phase: 'failed', running: false, error: 'runtime failure', results: [
      { account_id: 11, account_name: 'A only', error: 'record failure', notification_error: 'notification failure' },
    ] });
    const summary = ctx.$('inspection-runtime-error');
    const errorSummary = { inHeader: summary.parentElement.classList.contains('credential-head'), text: summary.textContent };
    assert.equal(errorSummary.inHeader, true);
    assert.equal(errorSummary.text.includes('record failure'), true);
    console.log(JSON.stringify({ before, during, after, errorSummary }, null, 2));
  } finally {
    ctx.close();
  }
})().catch((error) => { console.error(error); process.exitCode = 1; });
