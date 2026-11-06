# Publishing The VS Code Extension

## Publisher

Use the permanent Marketplace publisher ID `fanyang89`. The extension ID is `fanyang89.btfmt`.

## Pre-release

1. Keep `Cargo.toml` and `package.json` on the same `major.minor.patch` version.
2. Run the GitHub Actions `Release` workflow with `pre_release` enabled.
3. Download the `release-files` artifact after all test and platform build jobs pass.
4. Upload each platform-specific VSIX to the Visual Studio Marketplace publisher management page as a pre-release.
5. Install the Marketplace pre-release on at least one local and one Remote workspace before promotion.

The initial Marketplace release uses manual upload so the repository does not require a long-lived Azure DevOps Personal Access Token.

## Stable Release

1. Promote a tested pre-release by updating to the next stable version.
2. Create a matching `v<version>` tag.
3. Confirm the release workflow publishes the GitHub Release and all four platform-specific VSIX files.
4. Upload the VSIX files to the Marketplace without the pre-release flag.

## Automation

Each release matrix job installs its finished VSIX into an isolated VS Code
profile and runs the editor smoke suite before uploading the extension artifact.
Only an empty test harness runs in development mode; the btfmt client and server
come from the installed package. A failing `btfmt` executable placed first on
PATH prevents accidental use of a developer's global server installation.
Both editor test launchers enable VS Code's smoke-test driver on a temporary
local debugging port. Before exercising the editor, the suite awaits
`driver.whenWorkbenchRestored()`, including workbench contributions such as save
participants; extension-host startup alone does not guarantee those are ready.

Run the same check locally after packaging:

```bash
# Linux (requires xvfb)
xvfb-run -a npm run test:vsix -- btfmt-linux-x64-0.4.0.vsix
# macOS / Windows: use the matching platform package
npm run test:vsix -- <platform-specific.vsix>
```

The suite checks activation, completion, semantic highlighting, workspace search,
unsaved edits, and restart recovery. It also opens a real `.bt` file, enables
format-on-save, executes the editor's Save command, verifies disk contents,
closes and reopens the editor, formats selections, checks suppression boundaries,
and verifies diagnostics and formatting recover after correcting invalid syntax.

Before promotion, also perform a Remote SSH smoke check on a test host:

1. Install the host-platform VSIX in the remote Extension Host with no custom
   `btfmt.serverPath` setting.
2. Open a remote `.bt` file, enable format-on-save, edit and save it, then reopen
   it and verify the remote disk contents.
3. Format a selection next to an `off`/`on` region and verify surrounding bytes
   and the disabled region remain unchanged.
4. Introduce and correct a syntax error; verify diagnostics clear and saving
   formats again. Check completion, signature help, and navigation.
5. Reload the remote window and restart the language server; repeat a save and
   check the btfmt output log for the bundled server path and unexpected errors.

Hosted-runner tests do not validate Remote SSH transport or a real remote install.

Keep the `vsce.baseContentUrl` and `vsce.baseImagesUrl` settings in `package.json` pointed at the extension subdirectory. The Marketplace README uses these URLs for relative links and images; `vsce` otherwise assumes the repository root.

Future Marketplace automation must use Microsoft Entra ID workload identity federation. Do not add a long-lived `VSCE_PAT`; global Azure DevOps PATs retire on December 1, 2026.
