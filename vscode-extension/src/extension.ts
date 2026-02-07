import * as vscode from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  RevealOutputChannelOn,
  ServerOptions,
  State,
  Trace,
} from 'vscode-languageclient/node';
import { execFile } from 'child_process';

let client: LanguageClient | undefined;

export function activate(context: vscode.ExtensionContext): void {
  const config = vscode.workspace.getConfiguration('btfmt');
  const serverPath = config.get<string>('serverPath') || 'btfmt';

  const outputChannel = vscode.window.createOutputChannel('btfmt LSP');
  const traceChannel = vscode.window.createOutputChannel('btfmt LSP Trace');

  outputChannel.appendLine(`[Info ] Using server path: ${serverPath}`);

  // Test if the server is accessible (execFile auto-cleans up)
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
    synchronize: { configurationSection: 'btfmt' },
    outputChannel,
    traceOutputChannel: traceChannel,
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

  client = new LanguageClient('btfmt', 'btfmt LSP', serverOptions, clientOptions);
  outputChannel.appendLine('[Info ] btfmt LSP activated');
  client.onDidChangeState((event) => {
    outputChannel.appendLine(
      `[State] ${formatState(event.oldState)} -> ${formatState(event.newState)}`
    );
  });
  void client.setTrace(Trace.Verbose);
  void client.start().catch((err) => {
    const message = err instanceof Error ? err.stack ?? err.message : String(err);
    outputChannel.appendLine(`[Error] failed to start: ${message}`);
    outputChannel.show(true);
  });

  context.subscriptions.push(client, outputChannel, traceChannel);

  context.subscriptions.push(
    vscode.commands.registerCommand('btfmt.restartLsp', () => {
      void restartLsp();
    })
  );

  context.subscriptions.push(
    vscode.workspace.onDidChangeConfiguration((event) => {
      if (!client) {
        return;
      }
      if (!event.affectsConfiguration('btfmt')) {
        return;
      }
      client.sendNotification('workspace/didChangeConfiguration', {
        settings: { btfmt: buildSettings() },
      });
    })
  );
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}

async function restartLsp(): Promise<void> {
	if (!client) {
		return;
	}
	await client.restart();
}

function buildSettings(): Record<string, unknown> {
	const config = vscode.workspace.getConfiguration('btfmt');
	const configPath = config.get<string>('configPath');

	const settings: Record<string, unknown> = {};
	if (configPath) {
		settings.configPath = configPath;
	}
	return settings;
}

function withTimeout<T>(promise: Thenable<T>, ms: number): Promise<T> {
  return Promise.race([
    promise as Promise<T>,
    new Promise<T>((_, reject) =>
      setTimeout(() => reject(new Error(`formatting timed out after ${ms}ms`)), ms)
    ),
  ]);
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
