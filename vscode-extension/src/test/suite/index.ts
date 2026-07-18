import * as assert from 'assert';
import * as vscode from 'vscode';

const extensionId = 'local.btfmt-lsp';

export async function run(): Promise<void> {
  const extension = vscode.extensions.getExtension(extensionId);
  assert.ok(extension, `extension not found: ${extensionId}`);
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

  await vscode.commands.executeCommand('btfmt.restartLsp');
  const restartedEdits = await formattingEdits(document.uri);
  assert.ok(restartedEdits.length > 0, 'formatting failed after LSP restart');
  assert.strictEqual(applyTextEdits(document, restartedEdits), expected);
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
