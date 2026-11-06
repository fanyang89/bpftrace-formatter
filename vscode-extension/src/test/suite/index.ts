import * as assert from 'assert';
import * as path from 'path';
import * as vscode from 'vscode';
import { testEditingLifecycle } from './lifecycle';
import { waitForWorkbench } from '../workbench';

const extensionId = 'fanyang89.btfmt';

export async function run(): Promise<void> {
  await waitForWorkbench();
  const extension = vscode.extensions.getExtension(extensionId);
  assert.ok(extension, `extension not found: ${extensionId}`);
  const installedRoot = process.env.BTFMT_INSTALLED_EXTENSIONS;
  if (installedRoot) {
    assert.strictEqual(path.dirname(path.relative(installedRoot, extension.extensionPath)), '.');
    assert.strictEqual(vscode.workspace.getConfiguration('btfmt').get('serverPath'), '');
    console.log(`Testing installed extension: ${extension.extensionPath}`);
  }
  await extension.activate();
  assert.strictEqual(extension.isActive, true);

  const document = await vscode.workspace.openTextDocument({
    language: 'bpftrace',
    content: 'BEGIN{exit();}\n',
  });

  const edits = await formattingEdits(document.uri);
  assert.ok(edits.length > 0, 'formatting provider returned no edits');
  const expected = 'BEGIN\n{\n    exit();\n}\n';
  assert.strictEqual(applyTextEdits(document, edits), expected);
  const formattedDocument = await vscode.workspace.openTextDocument({
    language: 'bpftrace',
    content: expected,
  });
  const unchangedEdits = await vscode.commands.executeCommand<vscode.TextEdit[]>(
    'vscode.executeFormatDocumentProvider',
    formattedDocument.uri,
    { tabSize: 4, insertSpaces: true }
  );
  assert.ok(!unchangedEdits || unchangedEdits.length === 0);

  const selectionSource = 'BEGIN{exit();}\nEND{exit();}\n';
  const selectionDocument = await vscode.workspace.openTextDocument({
    language: 'bpftrace',
    content: selectionSource,
  });
  const selectionEdits = await vscode.commands.executeCommand<vscode.TextEdit[]>(
    'vscode.executeFormatRangeProvider',
    selectionDocument.uri,
    new vscode.Range(0, 0, 0, 'BEGIN{exit();}'.length),
    { tabSize: 4, insertSpaces: true }
  );
  assert.ok(selectionEdits?.length, 'range formatting provider returned no edits');
  assert.strictEqual(
    applyTextEdits(selectionDocument, selectionEdits),
    expected + 'END{exit();}\n'
  );

  const suppressedSource = '// btfmt: off\nBEGIN{exit();}\n// btfmt: on\n';
  const suppressedDocument = await vscode.workspace.openTextDocument({
    language: 'bpftrace',
    content: suppressedSource + 'END{exit();}\n',
  });
  const suppressedEdits = await formattingEdits(suppressedDocument.uri);
  assert.strictEqual(
    applyTextEdits(suppressedDocument, suppressedEdits),
    suppressedSource + 'END\n{\n    exit();\n}\n'
  );

  const semanticDocument = await vscode.workspace.openTextDocument({
    language: 'bpftrace',
    content: 'BEGIN { $value = 1; @count = $value; print(@count); }\n',
  });
  const legend = await vscode.commands.executeCommand<vscode.SemanticTokensLegend>(
    'vscode.provideDocumentSemanticTokensLegend',
    semanticDocument.uri
  );
  assert.ok(legend, 'semantic tokens provider did not register a legend');
  const tokens = await vscode.commands.executeCommand<vscode.SemanticTokens>(
    'vscode.provideDocumentSemanticTokens',
    semanticDocument.uri
  );
  assert.ok(tokens?.data.length, 'semantic tokens provider returned no tokens');
  assert.strictEqual(tokens.data.length % 5, 0);
  let line = 0;
  let character = 0;
  const highlighted = new Map<string, string>();
  for (let index = 0; index < tokens.data.length; index += 5) {
    const [deltaLine, deltaStart, length, tokenType] = tokens.data.slice(index, index + 5);
    line += deltaLine;
    character = deltaLine === 0 ? character + deltaStart : deltaStart;
    highlighted.set(
      semanticDocument.getText(new vscode.Range(line, character, line, character + length)),
      legend.tokenTypes[tokenType]
    );
  }
  assert.strictEqual(highlighted.get('$value'), 'variable');
  assert.ok(highlighted.has('@count'), 'maps were not semantically highlighted');
  assert.strictEqual(highlighted.get('print'), 'function');

  const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
  assert.ok(workspaceFolder, 'test workspace was not opened');
  const workspaceUri = vscode.Uri.joinPath(workspaceFolder.uri, 'symbols.bt');
  const symbols = await vscode.commands.executeCommand<vscode.SymbolInformation[]>(
    'vscode.executeWorkspaceSymbolProvider',
    'workspace_marker'
  );
  assert.strictEqual(symbols?.length, 1, 'workspace symbols were missing or duplicated');
  assert.strictEqual(symbols[0].name, '@workspace_marker');
  assert.strictEqual(symbols[0].location.uri.toString(), workspaceUri.toString());

  const workspaceDocument = await vscode.workspace.openTextDocument(workspaceUri);
  const change = new vscode.WorkspaceEdit();
  change.replace(
    workspaceUri,
    new vscode.Range(0, 0, workspaceDocument.lineCount, 0),
    'BEGIN { @unsaved_marker = count(); }\n'
  );
  assert.ok(await vscode.workspace.applyEdit(change));
  const updatedSymbols = await vscode.commands.executeCommand<vscode.SymbolInformation[]>(
    'vscode.executeWorkspaceSymbolProvider',
    'unsaved_marker'
  );
  assert.strictEqual(updatedSymbols?.length, 1, 'unsaved workspace symbols were not updated');
  assert.strictEqual(updatedSymbols[0].name, '@unsaved_marker');
  const oldSymbols = await vscode.commands.executeCommand<vscode.SymbolInformation[]>(
    'vscode.executeWorkspaceSymbolProvider',
    'workspace_marker'
  );
  assert.strictEqual(oldSymbols?.length, 0, 'workspace symbol search returned stale data');

  const completions = await completionItems(document.uri, new vscode.Position(0, 8));
  assert.ok(
    completions.items.some((item) => item.label === 'exit'),
    'completion provider did not return the matching builtin'
  );

  const dynamicSource =
    'tracepoint:syscalls:sys_enter_openat { print(args.filename); }\n';
  const dynamicDocument = await vscode.workspace.openTextDocument({
    language: 'bpftrace',
    content: dynamicSource,
  });
  const targets = await completionItems(
    dynamicDocument.uri,
    dynamicDocument.positionAt(dynamicSource.indexOf('tracepoint:sysc') + 'tracepoint:sysc'.length)
  );
  assert.ok(
    targets.items.some((item) => item.label === 'syscalls:sys_enter_openat'),
    'completion provider did not return the probe target'
  );
  assert.ok(
    targets.items.some((item) => item.label === 'syscalls:sys_enter_access'),
    'completion provider did not return the portable probe target'
  );
  const fields = await completionItems(
    dynamicDocument.uri,
    dynamicDocument.positionAt(dynamicSource.indexOf('args.') + 'args.'.length)
  );
  assert.ok(
    fields.items.some((item) => item.label === 'filename'),
    'completion provider did not return the args field'
  );

  const commands = await vscode.commands.getCommands(true);
  assert.ok(commands.includes('btfmt.showLogs'));
  assert.ok(commands.includes('btfmt.openSettings'));
  await Promise.all([
    vscode.commands.executeCommand('btfmt.restartLsp'),
    vscode.commands.executeCommand('btfmt.restartLsp'),
  ]);
  const restartedEdits = await formattingEdits(document.uri);
  assert.ok(restartedEdits.length > 0, 'formatting failed after LSP restart');
  assert.strictEqual(applyTextEdits(document, restartedEdits), expected);
  await testEditingLifecycle();
}

function applyTextEdits(
  document: vscode.TextDocument,
  edits: readonly vscode.TextEdit[]
): string {
  let text = document.getText();
  const ordered = [...edits].sort(
    (left, right) =>
      document.offsetAt(right.range.start) - document.offsetAt(left.range.start)
  );
  for (const edit of ordered) {
    const start = document.offsetAt(edit.range.start);
    const end = document.offsetAt(edit.range.end);
    text = text.slice(0, start) + edit.newText + text.slice(end);
  }
  return text;
}

async function formattingEdits(uri: vscode.Uri): Promise<vscode.TextEdit[]> {
  const deadline = Date.now() + 15_000;
  while (Date.now() < deadline) {
    const edits = await vscode.commands.executeCommand<vscode.TextEdit[]>(
      'vscode.executeFormatDocumentProvider',
      uri,
      { tabSize: 4, insertSpaces: true }
    );
    if (edits?.length) {
      return edits;
    }
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  throw new Error('timed out waiting for btfmt formatting provider');
}

async function completionItems(
  uri: vscode.Uri,
  position: vscode.Position
): Promise<vscode.CompletionList> {
  const deadline = Date.now() + 15_000;
  while (Date.now() < deadline) {
    const completions = await vscode.commands.executeCommand<vscode.CompletionList>(
      'vscode.executeCompletionItemProvider',
      uri,
      position
    );
    if (completions?.items.length) {
      return completions;
    }
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  throw new Error('timed out waiting for btfmt completion provider');
}
