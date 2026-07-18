import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import {
  LanguageClient,
  LanguageClientOptions,
  RevealOutputChannelOn,
  ServerOptions,
  State,
} from 'vscode-languageclient/node';
import { execFile } from 'child_process';

let client: LanguageClient | undefined;

function getBundledBinaryPath(context: vscode.ExtensionContext): string | undefined {
  const binaryName = process.platform === 'win32' ? 'btfmt.exe' : 'btfmt';
  const binaryPath = path.join(context.extensionPath, 'bin', binaryName);

  if (fs.existsSync(binaryPath)) {
    return binaryPath;
  }
  return undefined;
}

function getServerPath(context: vscode.ExtensionContext): string {
  const config = vscode.workspace.getConfiguration('btfmt');
  const configuredPath = config.get<string>('serverPath');

  // If user explicitly configured a path, use it
  if (configuredPath && configuredPath !== 'btfmt') {
    return configuredPath;
  }

  // Try bundled binary first
  const bundledPath = getBundledBinaryPath(context);
  if (bundledPath) {
    return bundledPath;
  }

  // Fall back to PATH
  return 'btfmt';
}

export function activate(context: vscode.ExtensionContext): void {
  const outputChannel = vscode.window.createOutputChannel('btfmt LSP');
  const fileWatcher = vscode.workspace.createFileSystemWatcher('**/*.bt');
  context.subscriptions.push(outputChannel, fileWatcher);
  startClient(context, outputChannel, fileWatcher);

  context.subscriptions.push(
    vscode.commands.registerCommand('btfmt.restartLsp', () =>
      restartLsp(context, outputChannel, fileWatcher)
    )
  );

  context.subscriptions.push(
    vscode.workspace.onDidChangeConfiguration((event) => {
      if (event.affectsConfiguration('btfmt.serverPath')) {
        void restartLsp(context, outputChannel, fileWatcher);
        return;
      }
      if (!client || !event.affectsConfiguration('btfmt.configPath')) {
        return;
      }
      client.sendNotification('workspace/didChangeConfiguration', {
        settings: { btfmt: buildSettings() },
      });
    })
  );
}

function startClient(
  context: vscode.ExtensionContext,
  outputChannel: vscode.OutputChannel,
  fileWatcher: vscode.FileSystemWatcher
): void {
  const serverPath = getServerPath(context);
  outputChannel.appendLine(`[Info ] Using server path: ${serverPath}`);

  execFile(serverPath, ['--help'], { timeout: 5000 }, (err, _stdout, _stderr) => {
    if (err) {
      outputChannel.appendLine(`[Error] Cannot run ${serverPath}: ${err.message}`);
    } else {
      outputChannel.appendLine(`[Info ] Server binary OK`);
    }
  });

  const serverOptions: ServerOptions = {
    command: serverPath,
    args: ['lsp'],
    options: { env: { ...process.env } },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ language: 'bpftrace' }],
    initializationOptions: { btfmt: buildSettings() },
    synchronize: { configurationSection: 'btfmt', fileEvents: fileWatcher },
    outputChannel,
    revealOutputChannelOn: RevealOutputChannelOn.Error,
    middleware: {
      provideDocumentFormattingEdits: async (document, options, token, next) => {
        outputChannel.appendLine(`[Format] request ${document.uri.toString()}`);
        try {
          const result = await withTimeout(next(document, options, token), 35_000);
          const count = Array.isArray(result) ? result.length : 0;
          outputChannel.appendLine(`[Format] response edits=${count}`);
          return result;
        } catch (err) {
          const message = err instanceof Error ? err.stack ?? err.message : String(err);
          outputChannel.appendLine(`[Format] error ${message}`);
          throw err;
        }
      },
    },
  };

  const nextClient = new LanguageClient('btfmt', 'btfmt LSP', serverOptions, clientOptions);
  client = nextClient;
  outputChannel.appendLine('[Info ] btfmt LSP activated');
  context.subscriptions.push(
    nextClient.onDidChangeState((event) => {
      outputChannel.appendLine(
        `[State] ${formatState(event.oldState)} -> ${formatState(event.newState)}`
      );
    }),
    nextClient
  );
  void nextClient.start().catch((err) => {
    const message = err instanceof Error ? err.stack ?? err.message : String(err);
    outputChannel.appendLine(`[Error] failed to start: ${message}`);
    outputChannel.show(true);
  });
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}

async function restartLsp(
  context: vscode.ExtensionContext,
  outputChannel: vscode.OutputChannel,
  fileWatcher: vscode.FileSystemWatcher
): Promise<void> {
  const previous = client;
  client = undefined;
  if (previous) {
    await previous.stop();
  }
  startClient(context, outputChannel, fileWatcher);
}

function buildSettings(): Record<string, unknown> {
  const config = vscode.workspace.getConfiguration('btfmt');
  return { configPath: config.get<string>('configPath', '') };
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
