# Workspace symbols

The LSP server advertises `workspaceSymbolProvider` and handles `workspace/symbol`.
In VS Code, use **Go to Symbol in Workspace** to search indexed bpftrace scripts.

## Search and navigation

- The query is trimmed and matched as a case-insensitive substring of the displayed
  name. An empty or whitespace-only query returns all available symbols. There is
  no fuzzy matching, ranking, or result cap separate from the workspace index.
- Global maps use their `@name` spelling (including the unnamed `@` map) and LSP
  `Variable` kind. Each file contributes one entry per map, at its first indexed
  definition (declaration, assignment, or update), not at every reference.
- Macro definitions use the bare macro name and `Function` kind. Each file
  contributes one entry per macro name, at its first definition. Calls, unresolved
  names, built-ins, parameters, scratch variables, and macro-local maps are excluded.
- Probe/action blocks use the whitespace-normalized probe list and `Event` kind,
  matching document symbols. Distinct blocks remain separate, even when their
  probe lists have the same spelling.
- Results are ordered by name, URI, and source range. Each result includes a file
  URI and a UTF-16 selection range for direct navigation; the container is the
  file's URI path. Same-name definitions in different files remain distinct.

## Workspace scope and lifecycle

Search uses the existing workspace index and immutable document snapshots. It does
not scan directories, load files, resolve imports, or reparse documents on a query.
All active workspace folders are searched: the LSP request has no originating
document or root selector. Files from unrelated or removed roots are excluded,
including open buffers. Overlapping roots do not duplicate entries, and symbols
with the same name in separate roots are not merged.

Unsaved open documents override their disk snapshots. File-backed new buffers
inside a workspace are included even before saving. Non-file/untitled documents
are excluded. Closing a buffer restores the latest indexed disk snapshot, or
removes the result if the file was deleted. Save and watched-file create/change/
delete notifications use the same index lifecycle. Folder changes rebuild that
index while preserving eligible open-buffer contents.

The existing `.bt` file filter, hidden-file filtering, `.ignore`, and Git ignore
rules apply to disk files and open-buffer results. Git rules require a repository,
as with the existing index walker. Symlink files and directories are excluded.
Notifications for `.ignore`/`.gitignore` changes or deletion rebuild the existing
index. Clients must send watched-file notifications to reflect external changes;
the server does not run its own file watcher.

Restricted Mode (`initializationOptions.btfmt.trusted: false`) does not scan roots
and always returns an empty workspace-symbol list, even if cached snapshots are
present. Current-document features remain available. Trust is fixed for a server
session; the VS Code client restarts the server when workspace trust is granted.

## Limits

The workspace scan retains its existing limits: 5,000 files, 2 MiB per file, and
50 MiB total. Results cover the available index, not necessarily every file in a
workspace that exceeds these limits. The legacy `SymbolInformation[]` response has
no completeness flag. Files with parser errors are omitted until a valid snapshot
is available, rather than returning stale definitions. Imported files outside the
active roots are not included, and workspace search does not evaluate macros or
infer runtime symbol identities across scripts.
