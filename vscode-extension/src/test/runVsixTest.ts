import * as fs from 'fs';
import * as os from 'os';
import * as path from 'path';
import { runTests, runVSCodeCommand } from '@vscode/test-electron';

async function main(): Promise<void> {
  const argument = process.argv[2];
  if (!argument || process.argv.length !== 3) {
    throw new Error('usage: npm run test:vsix -- <platform-specific.vsix>');
  }
  const vsix = path.resolve(argument);
  if (!fs.statSync(vsix).isFile()) {
    throw new Error(`VSIX not found: ${vsix}`);
  }
  const root = fs.mkdtempSync(path.join(os.tmpdir(), 'btfmt-vsix-'));
  try {
    const workspace = path.join(root, 'workspace');
    const harness = path.join(root, 'harness');
    const extensions = path.join(root, 'extensions');
    const profile = path.join(root, 'profile');
    const poisonPath = path.join(root, 'path');
    fs.mkdirSync(workspace);
    fs.mkdirSync(harness);
    fs.mkdirSync(poisonPath);
    // Fail if the installed extension falls back to a developer's PATH binary.
    fs.writeFileSync(path.join(poisonPath, process.platform === 'win32' ? 'btfmt.exe' : 'btfmt'),
      '#!/bin/sh\nexit 99\n', { mode: 0o755 });
    const pathKey = Object.keys(process.env).find((key) => key.toLowerCase() === 'path') ?? 'PATH';
    fs.writeFileSync(path.join(workspace, 'symbols.bt'),
      'BEGIN { @workspace_marker = count(); print(@workspace_marker); }\n');
    // Only this empty harness is loaded in development mode. btfmt must come
    // from the installed VSIX, including its compiled client and native server.
    fs.writeFileSync(path.join(harness, 'package.json'), JSON.stringify({
      name: 'btfmt-test-harness', publisher: 'btfmt-tests', version: '0.0.0',
      engines: { vscode: '^1.85.0' },
    }));
    const profileArgs = ['--user-data-dir', profile, '--extensions-dir', extensions];
    const installed = await runVSCodeCommand([
      ...profileArgs, '--install-extension', vsix, '--force',
    ]);
    process.stdout.write(installed.stdout);
    await runTests({
      extensionDevelopmentPath: harness,
      extensionTestsPath: path.resolve(__dirname, 'suite/index'),
      extensionTestsEnv: {
        BTFMT_INSTALLED_EXTENSIONS: extensions,
        [pathKey]: poisonPath + path.delimiter + (process.env[pathKey] ?? ''),
      },
      launchArgs: [workspace, ...profileArgs, '--disable-workspace-trust', '--skip-welcome'],
    });
  } finally {
    fs.rmSync(root, { recursive: true, force: true, maxRetries: 5, retryDelay: 200 });
  }
}

void main().catch((error) => {
  console.error(error);
  process.exit(1);
});
