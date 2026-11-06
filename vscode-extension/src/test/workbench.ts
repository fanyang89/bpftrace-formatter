import * as net from 'net';

export async function debugPort(): Promise<number> {
  const server = net.createServer();
  await new Promise<void>((resolve, reject) => {
    server.once('error', reject);
    server.listen(0, '127.0.0.1', resolve);
  });
  const port = (server.address() as net.AddressInfo).port;
  await new Promise<void>((resolve, reject) => server.close((error) => error ? reject(error) : resolve()));
  return port;
}

export async function waitForWorkbench(): Promise<void> {
  const port = Number(process.env.BTFMT_TEST_DEBUG_PORT);
  if (!Number.isInteger(port) || port <= 0) {
    throw new Error('BTFMT_TEST_DEBUG_PORT is required');
  }
  const response = await fetch(`http://127.0.0.1:${port}/json/list`, {
    signal: AbortSignal.timeout(10_000),
  });
  if (!response.ok) {
    throw new Error(`Cannot inspect test workbench: HTTP ${response.status}`);
  }
  const pages = await response.json() as { type: string; url: string; webSocketDebuggerUrl: string }[];
  const page = pages.find((page) => page.type === 'page' && page.url.includes('workbench'));
  if (!page) {
    throw new Error('Test workbench page not found');
  }
  // VS Code's own UI tests await both restored lifecycle and contributions.
  // Extension-host startup alone does not guarantee save participants exist.
  const socket = new WebSocket(page.webSocketDebuggerUrl);
  try {
    await new Promise<void>((resolve, reject) => {
      const timer = setTimeout(() => reject(new Error('Workbench restoration timed out')), 30_000);
      const finish = (error?: Error) => {
        clearTimeout(timer);
        error ? reject(error) : resolve();
      };
      socket.addEventListener('error', () => finish(new Error('Workbench debugger connection failed')));
      socket.addEventListener('close', () => finish(new Error('Workbench debugger closed before restoration')));
      socket.addEventListener('open', () => socket.send(JSON.stringify({
        id: 1,
        method: 'Runtime.evaluate',
        params: {
          expression: '(async () => { await globalThis.driver.whenWorkbenchRestored(); return true; })()',
          awaitPromise: true,
          returnByValue: true,
        },
      })));
      socket.addEventListener('message', (event) => {
        const message = JSON.parse(String(event.data));
        if (message.id !== 1) {
          return;
        }
        if (message.error || message.result?.exceptionDetails || message.result?.result?.value !== true) {
          finish(new Error(`Workbench restoration failed: ${JSON.stringify(message)}`));
        } else {
          finish();
        }
      });
    });
  } finally {
    socket.close();
  }
  console.log('Workbench restored, including save participants.');
}
