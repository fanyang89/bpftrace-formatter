import * as fs from 'fs';
import * as os from 'os';
import * as path from 'path';
import { runTests } from '@vscode/test-electron';

async function main(): Promise<void> {
  const extensionDevelopmentPath = path.resolve(__dirname, '../..');
  const extensionTestsPath = path.resolve(__dirname, './suite/index');
  const source = process.env.BTFMT_PATH;
  if (!source) {
    throw new Error('BTFMT_PATH must point to a built btfmt binary');
  }

  const binDir = path.join(extensionDevelopmentPath, 'bin');
  const binaryName = process.platform === 'win32' ? 'btfmt.exe' : 'btfmt';
  const bundledBinary = path.join(binDir, binaryName);
  if (fs.existsSync(bundledBinary)) {
    throw new Error(`test binary already exists: ${bundledBinary}`);
  }

  fs.mkdirSync(binDir, { recursive: true });
  fs.copyFileSync(path.resolve(source), bundledBinary);
  fs.chmodSync(bundledBinary, 0o755);

  const workspacePath = fs.realpathSync.native(fs.mkdtempSync(path.join(os.tmpdir(), 'btfmt-extension-')));
  try {
    fs.writeFileSync(
      path.join(workspacePath, 'symbols.bt'),
      'BEGIN { @workspace_marker = count(); print(@workspace_marker); }\n'
    );
    await runTests({
      extensionDevelopmentPath,
      extensionTestsPath,
      launchArgs: [workspacePath, '--disable-extensions', '--disable-workspace-trust'],
    });
  } finally {
    fs.rmSync(bundledBinary, { force: true });
    fs.rmdirSync(binDir);
    fs.rmSync(workspacePath, { recursive: true, force: true });
  }
}

void main().catch((error) => {
  console.error(error);
  process.exit(1);
});
