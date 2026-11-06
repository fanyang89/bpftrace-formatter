import * as assert from 'assert';
import * as vscode from 'vscode';

export async function testEditingLifecycle(): Promise<void> {
  const root = vscode.workspace.workspaceFolders?.[0].uri;
  assert.ok(root);
  const uri = vscode.Uri.joinPath(root, 'editing lifecycle 中文.bt');
  await vscode.workspace.fs.writeFile(uri, Buffer.from('BEGIN{}\n'));
  let document = await vscode.workspace.openTextDocument(uri);
  assert.strictEqual(document.languageId, 'bpftrace');
  assert.strictEqual(vscode.workspace.getConfiguration('editor', document).get('formatOnSave'), true);
  let editor = await vscode.window.showTextDocument(document, { preview: false });
  const protectedRegion = '// btfmt: off\nEND{  @keep=1; }  \t\n// btfmt: on\n';
  const source = 'BEGIN{printf("中文😀");}\n' + protectedRegion;
  const expected = 'BEGIN\n{\n    printf("中文😀");\n}\n' + protectedRegion;

  await replace(document, source);
  assert.ok(document.isDirty);
  await save(document);
  assert.strictEqual(document.getText(), expected, 'save did not apply formatting');
  assert.strictEqual(await diskText(uri), expected, 'formatted content was not written to disk');

  // Make another real edit, then save; the output must remain stable.
  await replace(document, ' ' + expected);
  await save(document);
  assert.strictEqual(await diskText(uri), expected, 'second save changed formatted output');
  await vscode.commands.executeCommand('workbench.action.closeActiveEditor');
  document = await vscode.workspace.openTextDocument(uri);
  editor = await vscode.window.showTextDocument(document, { preview: false });
  assert.strictEqual(document.getText(), expected, 'reopened document differs from disk');

  const selectionSource = 'BEGIN{exit();}\n' + protectedRegion + 'interval:s:1{exit();}\n';
  await replace(document, selectionSource);
  editor.selection = new vscode.Selection(0, 0, 0, 'BEGIN{exit();}'.length);
  await vscode.commands.executeCommand('editor.action.formatSelection');
  const selected = 'BEGIN\n{\n    exit();\n}\n' + protectedRegion + 'interval:s:1{exit();}\n';
  assert.strictEqual(document.getText(), selected, 'selection changed surrounding content');
  const start = document.positionAt(selected.indexOf('@keep'));
  editor.selection = new vscode.Selection(start, start.translate(0, '@keep'.length));
  await vscode.commands.executeCommand('editor.action.formatSelection');
  assert.strictEqual(document.getText(), selected, 'selection changed a suppressed region');

  await replace(document, 'BEGIN { $x=; }\n');
  await waitForDiagnostics(uri, true);
  await replace(document, 'BEGIN{$x=1;}\n');
  await waitForDiagnostics(uri, false);
  await save(document);
  assert.strictEqual(await diskText(uri), 'BEGIN\n{\n    $x = 1;\n}\n',
    'formatting did not recover after a syntax error');

  await vscode.commands.executeCommand('btfmt.restartLsp');
  await replace(document, 'BEGIN{$x=2;}\n');
  await save(document);
  assert.strictEqual(await diskText(uri), 'BEGIN\n{\n    $x = 2;\n}\n',
    'save formatting did not recover after restarting the server');
  await vscode.commands.executeCommand('workbench.action.closeActiveEditor');
  console.log('Editing lifecycle passed: save, reopen, selection, suppression, diagnostics, restart.');
}

async function replace(document: vscode.TextDocument, text: string): Promise<void> {
  const edit = new vscode.WorkspaceEdit();
  edit.replace(document.uri, new vscode.Range(0, 0, document.lineCount, 0), text);
  assert.ok(await vscode.workspace.applyEdit(edit));
}

async function save(document: vscode.TextDocument): Promise<void> {
  await vscode.commands.executeCommand('workbench.action.files.save');
  assert.strictEqual(document.isDirty, false, 'save did not finish');
}

async function diskText(uri: vscode.Uri): Promise<string> {
  return Buffer.from(await vscode.workspace.fs.readFile(uri)).toString('utf8');
}

async function waitForDiagnostics(uri: vscode.Uri, expectedError: boolean): Promise<void> {
  const deadline = Date.now() + 15_000;
  while (Date.now() < deadline) {
    const errors = vscode.languages.getDiagnostics(uri)
      .filter((diagnostic) => diagnostic.severity === vscode.DiagnosticSeverity.Error);
    if ((errors.length > 0) === expectedError) {
      return;
    }
    await new Promise((resolve) => setTimeout(resolve, 50));
  }
  throw new Error(`diagnostics did not ${expectedError ? 'appear' : 'clear'} for ${uri}`);
}
