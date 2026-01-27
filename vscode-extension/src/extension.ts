import * as vscode from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} from 'vscode-languageclient/node';

let client: LanguageClient | undefined;

export function activate(context: vscode.ExtensionContext): void {
  const config = vscode.workspace.getConfiguration('btfmt');
  const serverPath = config.get<string>('serverPath') || 'btfmt';

  const serverOptions: ServerOptions = {
    command: serverPath,
    args: ['lsp'],
    options: { env: { ...process.env } },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: 'file', language: 'bpftrace' }],
    initializationOptions: { btfmt: buildSettings() },
    synchronize: { configurationSection: 'btfmt' },
    outputChannelName: 'btfmt LSP',
  };

  client = new LanguageClient('btfmt', 'btfmt LSP', serverOptions, clientOptions);
  void client.start();
  context.subscriptions.push(client);

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

function buildSettings(): Record<string, unknown> {
	const config = vscode.workspace.getConfiguration('btfmt');
	const configPath = config.get<string>('configPath');

	const settings: Record<string, unknown> = {};
	if (configPath) {
		settings.configPath = configPath;
	}
	return settings;
}
