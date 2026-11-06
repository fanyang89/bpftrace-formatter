# Language client backport

`vscode-languageclient+9.0.1.patch` backports the pending-change cleanup and
re-entrant Delayer fix from microsoft/vscode-languageserver-node commit
`d11d518e5a61c64dd2e54f227b3bd5a43290a7eb`:

https://github.com/microsoft/vscode-languageserver-node/commit/d11d518e5a61c64dd2e54f227b3bd5a43290a7eb

Stopping a client with a full-text change pending otherwise leaves a timer that
retries delivery after shutdown and clears its own re-entrant task. This produced
`Delivering pending changes failed: TypeError: this.task is not a function`
in the installed-VSIX restart tests. The fix is not in npm release 11.0.0 either.

`npm ci` applies the patch through `postinstall`. Remove the patch and
`patch-package` when upgrading to a release containing the upstream fix.
`npm run test:client` exercises re-entrant delivery with the installed dependency.
