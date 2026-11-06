import * as fs from 'fs';
import * as path from 'path';

export function prepareWorkspace(workspace: string): void {
  fs.mkdirSync(path.join(workspace, '.vscode'), { recursive: true });
  fs.writeFileSync(path.join(workspace, '.vscode/settings.json'), JSON.stringify({
    '[bpftrace]': {
      'editor.defaultFormatter': 'fanyang89.btfmt',
      'editor.formatOnSave': true,
    },
  }));
  fs.writeFileSync(path.join(workspace, 'symbols.bt'),
    'BEGIN { @workspace_marker = count(); print(@workspace_marker); }\n');
}
