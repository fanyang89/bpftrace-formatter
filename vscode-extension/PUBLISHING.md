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

Future Marketplace automation must use Microsoft Entra ID workload identity federation. Do not add a long-lived `VSCE_PAT`; global Azure DevOps PATs retire on December 1, 2026.
