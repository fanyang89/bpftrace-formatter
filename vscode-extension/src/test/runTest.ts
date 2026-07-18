import * as fs from 'fs';
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

  try {
    await runTests({
      extensionDevelopmentPath,
      extensionTestsPath,
      launchArgs: ['--disable-extensions'],
    });
  } finally {
    fs.rmSync(bundledBinary, { force: true });
    fs.rmdirSync(binDir);
  }
}

void main().catch((error) => {
  console.error(error);
  process.exit(1);
});
