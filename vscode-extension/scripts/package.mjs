import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';

const [, , target, binaryArgument, channel = 'release'] = process.argv;
const hostTarget = new Map([
  ['linux:x64', 'linux-x64'],
  ['darwin:x64', 'darwin-x64'],
  ['darwin:arm64', 'darwin-arm64'],
  ['win32:x64', 'win32-x64'],
]).get(`${process.platform}:${process.arch}`);

if (!target || !binaryArgument) {
  throw new Error('usage: node scripts/package.mjs <target> <btfmt-binary> [release|pre-release]');
}
if (target !== hostTarget) {
  throw new Error(`target ${target} does not match the current host ${hostTarget ?? 'unsupported'}`);
}
if (channel !== 'release' && channel !== 'pre-release') {
  throw new Error(`unsupported channel: ${channel}`);
}

const binary = path.resolve(binaryArgument);
if (!fs.statSync(binary).isFile()) {
  throw new Error(`btfmt binary not found: ${binary}`);
}
execFileSync(binary, ['--version'], { stdio: 'inherit' });

const binDir = path.resolve('bin');
if (fs.existsSync(binDir)) {
  throw new Error(`refusing to replace existing staging directory: ${binDir}`);
}
fs.mkdirSync(binDir);

const bundledName = process.platform === 'win32' ? 'btfmt.exe' : 'btfmt';
const bundledBinary = path.join(binDir, bundledName);
const vsce = path.resolve('node_modules', '@vscode', 'vsce', 'vsce');

try {
  fs.copyFileSync(binary, bundledBinary);
  fs.chmodSync(bundledBinary, 0o755);

  const files = execFileSync(process.execPath, [vsce, 'ls', '--no-dependencies'], {
    encoding: 'utf8',
  });
  for (const required of [
    `bin/${bundledName}`,
    'images/icon.png',
    'images/completion.png',
    'CHANGELOG.md',
    'SUPPORT.md',
  ]) {
    if (!files.includes(required)) {
      throw new Error(`package is missing required file: ${required}`);
    }
  }
  process.stdout.write(files);

  const args = ['package', '--no-dependencies', '--target', target];
  if (channel === 'pre-release') {
    args.push('--pre-release');
  }
  execFileSync(process.execPath, [vsce, ...args], { stdio: 'inherit' });
} finally {
  fs.rmSync(binDir, { recursive: true, force: true });
}
