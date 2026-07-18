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
  const fakeBinDir = fs.mkdtempSync(path.join(os.tmpdir(), 'btfmt-test-bin-'));
  const fakeBpftrace = path.join(fakeBinDir, 'bpftrace');
  fs.writeFileSync(
    fakeBpftrace,
    `#!/usr/bin/env node
const probes = [
  'tracepoint:sched:sched_switch',
  'tracepoint:syscalls:sys_enter_openat',
];
if (process.argv[2] === '-l') {
  for (const probe of probes) console.log(probe);
} else if (process.argv[2] === '-lv') {
  console.log('tracepoint:syscalls:sys_enter_openat');
  console.log('    int __syscall_nr');
  console.log('    const char * filename');
} else {
  process.exit(2);
}
`
  );
  fs.chmodSync(fakeBpftrace, 0o755);
  process.env.PATH = `${fakeBinDir}${path.delimiter}${process.env.PATH ?? ''}`;

  try {
    await runTests({
      extensionDevelopmentPath,
      extensionTestsPath,
      launchArgs: ['--disable-extensions'],
    });
  } finally {
    fs.rmSync(bundledBinary, { force: true });
    fs.rmdirSync(binDir);
    fs.rmSync(fakeBinDir, { recursive: true, force: true });
  }
}

void main().catch((error) => {
  console.error(error);
  process.exit(1);
});
