const assert = require('node:assert/strict');
const { test } = require('node:test');
require('vscode-languageserver-protocol/node');
const { Delayer } = require('vscode-languageclient/lib/common/utils/async');

test('pending change delivery can retrigger without losing its task', async () => {
  const delayer = new Delayer(0);
  let retry;
  await delayer.trigger(() => {
    retry = delayer.trigger(() => 'delivered');
  });
  assert.equal(await retry, 'delivered');
});

test('forced delivery preserves a re-entrant task', async () => {
  const delayer = new Delayer(0);
  let retry;
  void delayer.trigger(() => {
    retry = delayer.trigger(() => 'delivered');
  });
  delayer.forceDelivery();
  assert.equal(await retry, 'delivered');
});
