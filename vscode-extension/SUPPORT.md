# Support

## Check Server Status

Open a `.bt` file and select the btfmt language status item. Use **btfmt: Show Logs** for server startup and request diagnostics.

## Common Issues

### Server Does Not Start

Reset `btfmt.serverPath` to use the bundled binary, then run **btfmt: Restart Server**.

### Probe Completion Is Limited

btfmt does not require root. Portable catalog suggestions remain available without kernel access. Host-specific targets depend on kernel metadata readable by the current user.

### Remote Workspace Uses The Wrong Platform

Install btfmt in the Remote SSH, WSL, or Dev Container Extension Host rather than only in the local UI.

## Report An Issue

Open an issue at <https://github.com/fanyang89/bpftrace-formatter/issues> with:

- btfmt and VS Code versions
- Host or remote platform
- Relevant **btfmt: Show Logs** output
- A minimal `.bt` example when possible
