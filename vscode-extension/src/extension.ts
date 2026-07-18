import * as fs from 'fs';
import * as path from 'path';
import * as vscode from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  RevealOutputChannelOn,
  ServerOptions,
  State,
} from 'vscode-languageclient/node';

let client: LanguageClient | undefined;
let stateSubscription: vscode.Disposable | undefined;
let lifecycle: Promise<void> = Promise.resolve();

function getBundledBinaryPath(context: vscode.ExtensionContext): string | undefined {
  const binaryName = process.platform === 'win32' ? 'btfmt.exe' : 'btfmt';
  const binaryPath = path.join(context.extensionPath, 'bin', binaryName);
  return fs.existsSync(binaryPath) ? binaryPath : undefined;
}

function getServerPath(context: vscode.ExtensionContext): string {
  const configuredPath = vscode.workspace
    .getConfiguration('btfmt')
    .get<string>('serverPath', '')
    .trim();
  if (configuredPath) {
    return configuredPath;
  }
  return getBundledBinaryPath(context) ?? 'btfmt';
}

export function activate(context: vscode.ExtensionContext): void {
  const log = vscode.window.createOutputChannel('btfmt', { log: true });
  const status = vscode.languages.createLanguageStatusItem('btfmt.server', {
    language: 'bpftrace',
  });
  status.name = 'btfmt';
  const fileWatcher = vscode.workspace.createFileSystemWatcher('**/*.bt');
  context.subscriptions.push(log, status, fileWatcher);

  context.subscriptions.push(
    vscode.commands.registerCommand('btfmt.restartLsp', () =>
      enqueue(() => restartClient(context, log, status, fileWatcher)).catch(() => undefined)
    ),
    vscode.commands.registerCommand('btfmt.showLogs', () => log.show(true)),
    vscode.commands.registerCommand('btfmt.openSettings', () =>
      vscode.commands.executeCommand(
        'workbench.action.openSettings',
        `@ext:${context.extension.id}`
      )
    ),
    vscode.workspace.onDidGrantWorkspaceTrust(() =>
      enqueue(() => restartClient(context, log, status, fileWatcher)).catch(() => undefined)
    ),
    vscode.workspace.onDidChangeConfiguration((event) => {
      if (event.affectsConfiguration('btfmt.serverPath')) {
        void enqueue(() => restartClient(context, log, status, fileWatcher)).catch(
          () => undefined
        );
        return;
      }
      if (!client || !event.affectsConfiguration('btfmt.configPath')) {
        return;
      }
      void client.sendNotification('workspace/didChangeConfiguration', {
        settings: { btfmt: buildSettings() },
      });
    })
  );

  void enqueue(() => startClient(context, log, status, fileWatcher)).catch(() => undefined);
}

async function startClient(
  context: vscode.ExtensionContext,
  log: vscode.LogOutputChannel,
  status: vscode.LanguageStatusItem,
  fileWatcher: vscode.FileSystemWatcher
): Promise<void> {
  const serverPath = getServerPath(context);
  log.info(`Using server: ${serverPath}`);
  updateStatus(status, 'starting');

  const serverOptions: ServerOptions = {
    command: serverPath,
    args: ['lsp'],
    options: { env: { ...process.env } },
  };
  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ language: 'bpftrace' }],
    initializationOptions: { btfmt: buildSettings() },
    synchronize: { configurationSection: 'btfmt', fileEvents: fileWatcher },
    outputChannel: log,
    revealOutputChannelOn: RevealOutputChannelOn.Never,
    middleware: {
      provideDocumentFormattingEdits: async (document, options, token, next) => {
        log.debug(`Formatting ${document.uri.toString()}`);
        try {
          const result = await withTimeout(next(document, options, token), 35_000);
          const editCount = Array.isArray(result) ? result.length : 0;
          log.debug(`Formatting completed with ${editCount} edits`);
          return result;
        } catch (error) {
          log.error(error instanceof Error ? error : String(error));
          throw error;
        }
      },
    },
  };

  const nextClient = new LanguageClient('btfmt', 'btfmt', serverOptions, clientOptions);
  client = nextClient;
  stateSubscription?.dispose();
  stateSubscription = nextClient.onDidChangeState((event) => {
    log.debug(`Server ${formatState(event.oldState)} -> ${formatState(event.newState)}`);
    if (client !== nextClient) {
      return;
    }
    if (event.newState === State.Starting) {
      updateStatus(status, 'starting');
    } else if (event.newState === State.Running) {
      updateStatus(status, vscode.workspace.isTrusted ? 'ready' : 'limited');
    } else {
      updateStatus(status, 'stopped');
    }
  });

  try {
    await nextClient.start();
    if (client === nextClient) {
      updateStatus(status, vscode.workspace.isTrusted ? 'ready' : 'limited');
      log.info('Language server ready');
    }
  } catch (error) {
    if (client === nextClient) {
      client = undefined;
      updateStatus(status, 'error');
    }
    log.error(error instanceof Error ? error : String(error));
    throw error;
  }
}

async function restartClient(
  context: vscode.ExtensionContext,
  log: vscode.LogOutputChannel,
  status: vscode.LanguageStatusItem,
  fileWatcher: vscode.FileSystemWatcher
): Promise<void> {
  log.info('Restarting language server');
  await stopClient();
  await startClient(context, log, status, fileWatcher);
}

async function stopClient(): Promise<void> {
  const previous = client;
  client = undefined;
  stateSubscription?.dispose();
  stateSubscription = undefined;
  if (previous) {
    await previous.stop();
  }
}

function enqueue(operation: () => Promise<void>): Promise<void> {
  const next = lifecycle.then(operation, operation);
  lifecycle = next.catch(() => undefined);
  return next;
}

export function deactivate(): Thenable<void> {
  return enqueue(stopClient);
}

function buildSettings(): Record<string, unknown> {
  const config = vscode.workspace.getConfiguration('btfmt');
  return {
    configPath: config.get<string>('configPath', ''),
    trusted: vscode.workspace.isTrusted,
  };
}

function updateStatus(
  status: vscode.LanguageStatusItem,
  phase: 'starting' | 'ready' | 'limited' | 'stopped' | 'error'
): void {
  status.busy = phase === 'starting';
  switch (phase) {
    case 'starting':
      status.severity = vscode.LanguageStatusSeverity.Information;
      status.text = '$(sync~spin) btfmt';
      status.detail = 'Starting language server';
      break;
    case 'ready':
      status.severity = vscode.LanguageStatusSeverity.Information;
      status.text = '$(check) btfmt';
      status.detail = 'Language features ready';
      break;
    case 'limited':
      status.severity = vscode.LanguageStatusSeverity.Warning;
      status.text = '$(shield) btfmt';
      status.detail = 'Limited in Restricted Mode';
      break;
    case 'stopped':
      status.severity = vscode.LanguageStatusSeverity.Warning;
      status.text = '$(circle-slash) btfmt';
      status.detail = 'Language server stopped';
      break;
    case 'error':
      status.severity = vscode.LanguageStatusSeverity.Error;
      status.text = '$(error) btfmt';
      status.detail = 'Language server failed to start';
      break;
  }
  const canRestart = phase === 'error' || phase === 'stopped';
  status.command = {
    command: canRestart ? 'btfmt.restartLsp' : 'btfmt.showLogs',
    title: canRestart ? 'Restart Server' : 'Show Logs',
  };
}

function withTimeout<T>(
  result: vscode.ProviderResult<T>,
  ms: number
): Promise<T | null | undefined> {
  let timer: NodeJS.Timeout | undefined;
  return Promise.race([
    Promise.resolve(result),
    new Promise<never>((_, reject) => {
      timer = setTimeout(
        () => reject(new Error(`formatting timed out after ${ms}ms`)),
        ms
      );
    }),
  ]).finally(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });
}

function formatState(state: State): string {
  switch (state) {
    case State.Starting:
      return 'starting';
    case State.Running:
      return 'running';
    case State.Stopped:
      return 'stopped';
    default:
      return 'unknown';
  }
}
