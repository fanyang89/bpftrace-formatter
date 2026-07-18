#!/usr/bin/env python3
"""LSP smoke test for btfmt."""

import json
import os
import pathlib
import subprocess
import sys
import tempfile
import threading
import time
import traceback


def _now_ms():
    return int(time.time() * 1000)


def _frame(payload_bytes: bytes) -> bytes:
    return b"Content-Length: %d\r\n\r\n" % len(payload_bytes) + payload_bytes


class LSPClient:
    def __init__(self, argv, env):
        self.p = subprocess.Popen(
            argv,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=env,
            bufsize=0,
        )
        self._next_id = 1
        self._stderr_buf = bytearray()
        self._start_stderr_drain()

    def _start_stderr_drain(self):
        # Best-effort: capture stderr for debugging on failure.
        def run():
            try:
                while True:
                    b = self.p.stderr.read(4096)
                    if not b:
                        break
                    self._stderr_buf.extend(b)
            except Exception:
                pass

        t = threading.Thread(target=run, daemon=True)
        t.start()

    def stderr_text(self):
        try:
            return self._stderr_buf.decode("utf-8", errors="replace")
        except Exception:
            return ""

    def send(self, msg: dict):
        raw = json.dumps(msg, separators=(",", ":"), ensure_ascii=True).encode("utf-8")
        self.p.stdin.write(_frame(raw))
        self.p.stdin.flush()

    def request(self, method: str, params):
        rid = self._next_id
        self._next_id += 1
        msg = {"jsonrpc": "2.0", "id": rid, "method": method}
        if params is not None:
            msg["params"] = params
        self.send(msg)
        return rid

    def notify(self, method: str, params):
        msg = {"jsonrpc": "2.0", "method": method}
        if params is not None:
            msg["params"] = params
        self.send(msg)

    def _read_exact(self, n: int) -> bytes:
        out = bytearray()
        while len(out) < n:
            chunk = self.p.stdout.read(n - len(out))
            if not chunk:
                raise EOFError("unexpected EOF")
            out.extend(chunk)
        return bytes(out)

    def read_message(self, timeout_s: float = 10.0) -> dict:
        # Minimal stdio framing reader: parse headers until CRLFCRLF then read JSON body.
        deadline = time.time() + timeout_s
        header = bytearray()
        while True:
            if time.time() > deadline:
                raise TimeoutError("timeout waiting for LSP headers")
            b = self.p.stdout.read(1)
            if not b:
                raise EOFError("unexpected EOF while reading headers")
            header.extend(b)
            if header.endswith(b"\r\n\r\n"):
                break
            if len(header) > 64 * 1024:
                raise ValueError("header too large")

        headers = header.decode("ascii", errors="replace").split("\r\n")
        content_length = None
        for line in headers:
            if line.lower().startswith("content-length:"):
                content_length = int(line.split(":", 1)[1].strip())
                break
        if content_length is None:
            raise ValueError("missing Content-Length header")
        body = self._read_exact(content_length)
        return json.loads(body.decode("utf-8"))

    def wait_for_response(self, rid: int, timeout_s: float = 10.0):
        deadline = time.time() + timeout_s
        notifications = []
        while True:
            remaining = max(0.01, deadline - time.time())
            msg = self.read_message(timeout_s=remaining)
            if isinstance(msg, dict) and msg.get("id") == rid:
                return msg, notifications
            notifications.append(msg)

    def wait_for(self, predicate, timeout_s: float = 10.0):
        deadline = time.time() + timeout_s
        while True:
            remaining = max(0.01, deadline - time.time())
            msg = self.read_message(timeout_s=remaining)
            if predicate(msg):
                return msg

    def close(self):
        try:
            if self.p.poll() is None:
                self.p.terminate()
                self.p.wait(timeout=2)
        except Exception:
            try:
                self.p.kill()
            except Exception:
                pass


def apply_text_edits(text: str, edits):
    # LSP positions are UTF-16; this smoke uses ASCII-only fixtures.
    lines = text.splitlines(True)

    def pos_to_offset(line: int, character: int) -> int:
        if line < 0:
            return 0
        if line >= len(lines):
            return sum(len(x) for x in lines)
        return sum(len(lines[i]) for i in range(line)) + min(
            character, len(lines[line])
        )

    def edit_key(e):
        r = e.get("range") or {}
        s = r.get("start") or {}
        return (s.get("line", 0), s.get("character", 0))

    # Apply in reverse order.
    out = text
    for e in sorted(edits, key=edit_key, reverse=True):
        r = e.get("range") or {}
        s = r.get("start") or {}
        en = r.get("end") or {}
        start = pos_to_offset(s.get("line", 0), s.get("character", 0))
        end = pos_to_offset(en.get("line", 0), en.get("character", 0))
        new_text = e.get("newText", "")
        out = out[:start] + new_text + out[end:]
    return out


def offset_to_position(text: str, offset: int):
    before = text[:offset]
    return {
        "line": before.count("\n"),
        "character": len(before.rsplit("\n", 1)[-1].encode("utf-16-le")) // 2,
    }


def fail(summary, step: str, msg: str, **extra):
    summary["ok"] = False
    summary["failed_step"] = step
    summary["error"] = msg
    if extra:
        summary["error_extra"] = extra
    print(json.dumps(summary, indent=2, sort_keys=True))
    sys.exit(1)


def record(summary, step: str, started_ms: int, **data):
    entry = {"step": step, "ms": _now_ms() - started_ms}
    entry.update(data)
    summary["steps"].append(entry)


def main():
    summary = {
        "ok": True,
        "steps": [],
        "started_ms": _now_ms(),
    }

    lsp = None
    try:
        t0 = _now_ms()
        with tempfile.TemporaryDirectory(prefix="btfmt-lsp-smoke-") as td:
            td_path = pathlib.Path(td)
            doc_path = td_path / "smoke.bt"
            helper_path = td_path / "helper.bt"
            config_path = td_path / ".btfmt.json"
            second_root = td_path / "second"
            second_root.mkdir()
            second_doc_path = second_root / "smoke.bt"
            second_config_path = second_root / ".btfmt.json"
            env = os.environ.copy()
            # Make smoke deterministic: avoid picking up user/global configs.
            env["HOME"] = str(td_path)
            env["XDG_CONFIG_HOME"] = str(td_path)
            env["XDG_DATA_HOME"] = str(td_path)
            env["XDG_CACHE_HOME"] = str(td_path)
            fake_bin = td_path / "bin"
            fake_bin.mkdir()
            fake_bpftrace = fake_bin / "bpftrace"
            fake_bpftrace.write_text(
                """#!/usr/bin/env python3
import fnmatch
import sys

probes = [
    "tracepoint:sched:sched_switch",
    "tracepoint:syscalls:sys_enter_openat",
    "tracepoint:syscalls:sys_enter_openat2",
]
if len(sys.argv) == 3 and sys.argv[1] == "-l":
    for probe in probes:
        if fnmatch.fnmatch(probe, sys.argv[2]):
            print(probe)
elif len(sys.argv) == 3 and sys.argv[1] == "-lv":
    for probe in probes:
        if fnmatch.fnmatch(probe, sys.argv[2]):
            print(probe)
            print("    int __syscall_nr")
            print("    const char * filename")
            print("    int flags")
else:
    raise SystemExit(2)
""",
                encoding="utf-8",
            )
            fake_bpftrace.chmod(0o755)
            env["PATH"] = str(fake_bin) + os.pathsep + env.get("PATH", "")
            btfmt_path = os.environ.get("BTFMT_PATH") or str(
                (pathlib.Path.cwd() / "btfmt").resolve()
            )
            config_path.write_text('{"indent":{"size":2}}\n', encoding="utf-8")
            second_config_path.write_text('{"indent":{"size":6}}\n', encoding="utf-8")

            # Intentionally unformatted input.
            doc_text_v1 = (
                'import "helper.bt";\n'
                'tracepoint:syscalls:sys_enter_openat{$target=1;@opens[$target]=count();printf("%d", imported(1));printf("openat: %s\\n",str(args.filename));}\n'
                'tracepoint:syscalls:sys_enter_openat2{printf("openat2: %s\\n",str(args->filename));}\n'
            )
            helper_text = "macro imported(value) { value }\n"
            doc_path.write_text(doc_text_v1, encoding="utf-8")
            helper_path.write_text(helper_text, encoding="utf-8")
            doc_uri = doc_path.absolute().as_uri()
            helper_uri = helper_path.absolute().as_uri()
            second_doc_text = "BEGIN{exit();}\n"
            second_doc_path.write_text(second_doc_text, encoding="utf-8")
            second_doc_uri = second_doc_path.absolute().as_uri()

            cli = subprocess.run(
                [btfmt_path, str(doc_path)],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                env=env,
                cwd=str(td_path),
            )
            if cli.returncode != 0:
                fail(
                    summary,
                    "cli_format",
                    "CLI formatter returned non-zero",
                    returncode=cli.returncode,
                    stderr=cli.stderr[-4000:],
                )
            expected_formatted = cli.stdout
            record(
                summary, "cli_format", t0, bytes=len(expected_formatted.encode("utf-8"))
            )

            lsp = LSPClient([btfmt_path, "lsp"], env=env)
            record(summary, "spawn", t0, pid=lsp.p.pid)

            init_params = {
                "processId": None,
                "rootUri": td_path.absolute().as_uri(),
                "capabilities": {
                    "workspace": {"workspaceEdit": {"documentChanges": True}}
                },
                "clientInfo": {"name": "task-lsp-smoke", "version": "0"},
                "initializationOptions": {"btfmt": {"configPath": ".btfmt.json"}},
                "workspaceFolders": [
                    {"uri": td_path.absolute().as_uri(), "name": "btfmt-lsp-smoke"},
                    {"uri": second_root.absolute().as_uri(), "name": "second"},
                ],
            }
            rid = lsp.request("initialize", init_params)
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(summary, "initialize", "initialize returned error", resp=resp)
            capabilities = (resp.get("result") or {}).get("capabilities") or {}
            for supported in (
                "definitionProvider",
                "referencesProvider",
                "documentHighlightProvider",
                "renameProvider",
            ):
                if not capabilities.get(supported):
                    fail(
                        summary,
                        "initialize",
                        f"expected capability is not enabled: {supported}",
                        resp=resp,
                    )
            trigger_characters = (capabilities.get("completionProvider") or {}).get(
                "triggerCharacters"
            ) or []
            if not {":", "."}.issubset(trigger_characters):
                fail(
                    summary,
                    "initialize",
                    "dynamic completion triggers are missing",
                    resp=resp,
                )
            record(summary, "initialize", t0)

            lsp.notify("initialized", {})
            record(summary, "initialized", t0)

            lsp.notify(
                "textDocument/didOpen",
                {
                    "textDocument": {
                        "uri": doc_uri,
                        "languageId": "bpftrace",
                        "version": 1,
                        "text": doc_text_v1,
                    }
                },
            )
            record(summary, "didOpen", t0)

            rid = lsp.request(
                "textDocument/formatting",
                {
                    "textDocument": {"uri": doc_uri},
                    "options": {"tabSize": 4, "insertSpaces": True},
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(summary, "formatting", "formatting returned error", resp=resp)
            edits = resp.get("result")
            if not isinstance(edits, list):
                fail(
                    summary, "formatting", "formatting result is not a list", resp=resp
                )
            lsp_formatted = apply_text_edits(doc_text_v1, edits)
            if lsp_formatted != expected_formatted:
                fail(
                    summary,
                    "formatting",
                    "LSP formatting output mismatch vs CLI",
                    expected_prefix=expected_formatted[:200],
                    got_prefix=lsp_formatted[:200],
                )
            record(
                summary,
                "formatting",
                t0,
                edits=len(edits),
                bytes=len(lsp_formatted.encode("utf-8")),
            )

            lsp.notify(
                "textDocument/didOpen",
                {
                    "textDocument": {
                        "uri": second_doc_uri,
                        "languageId": "bpftrace",
                        "version": 1,
                        "text": second_doc_text,
                    }
                },
            )
            rid = lsp.request(
                "textDocument/formatting",
                {
                    "textDocument": {"uri": second_doc_uri},
                    "options": {"tabSize": 4, "insertSpaces": True},
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(
                    summary,
                    "multi_root_formatting",
                    "formatting returned error",
                    resp=resp,
                )
            second_formatted = apply_text_edits(
                second_doc_text, resp.get("result") or []
            )
            if "\n      exit();" not in second_formatted:
                fail(
                    summary,
                    "multi_root_formatting",
                    "second workspace config was not applied",
                    formatted=second_formatted,
                )
            record(summary, "multi_root_formatting", t0)

            rid = lsp.request(
                "textDocument/hover",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": {"line": 0, "character": 1},
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(summary, "hover", "hover returned error", resp=resp)
            record(summary, "hover", t0, has_result=resp.get("result") is not None)

            rid = lsp.request(
                "textDocument/documentSymbol",
                {"textDocument": {"uri": doc_uri}},
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(
                    summary,
                    "documentSymbol",
                    "documentSymbol returned error",
                    resp=resp,
                )
            syms = resp.get("result")
            if not isinstance(syms, list):
                fail(
                    summary,
                    "documentSymbol",
                    "documentSymbol result is not a list",
                    resp=resp,
                )
            record(summary, "documentSymbol", t0, count=len(syms))

            rid = lsp.request(
                "textDocument/completion",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": {"line": 0, "character": 0},
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(summary, "completion", "completion returned error", resp=resp)
            completion = resp.get("result")
            if isinstance(completion, dict):
                completion_count = len(completion.get("items") or [])
            elif isinstance(completion, list):
                completion_count = len(completion)
            else:
                completion_count = 0
            if completion_count == 0:
                fail(summary, "completion", "expected completion items", resp=resp)
            record(summary, "completion", t0, count=completion_count)

            rid = lsp.request(
                "textDocument/completion",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": offset_to_position(
                        doc_text_v1, doc_text_v1.index("imported(1)") + 5
                    ),
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            imported_completion = resp.get("result") or {}
            imported_items = (
                imported_completion.get("items")
                if isinstance(imported_completion, dict)
                else imported_completion
            ) or []
            if "imported" not in {item.get("label") for item in imported_items}:
                fail(
                    summary,
                    "context_completion",
                    "imported macro is missing from completion",
                    resp=resp,
                )

            rid = lsp.request(
                "textDocument/completion",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": offset_to_position(
                        doc_text_v1, doc_text_v1.index("$target") + 4
                    ),
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            variable_completion = resp.get("result") or {}
            variable_items = (
                variable_completion.get("items")
                if isinstance(variable_completion, dict)
                else variable_completion
            ) or []
            if [item.get("label") for item in variable_items] != ["$target"]:
                fail(
                    summary,
                    "context_completion",
                    "scratch completion is not scope-filtered",
                    resp=resp,
                )

            rid = lsp.request(
                "textDocument/completion",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": offset_to_position(
                        doc_text_v1, doc_text_v1.index("tracepoint:syscalls") + 15
                    ),
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            target_completion = resp.get("result") or {}
            target_items = (
                target_completion.get("items")
                if isinstance(target_completion, dict)
                else target_completion
            ) or []
            if "syscalls:sys_enter_openat" not in {
                item.get("label") for item in target_items
            }:
                fail(
                    summary,
                    "context_completion",
                    "probe target completion is missing",
                    resp=resp,
                )

            rid = lsp.request(
                "textDocument/completion",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": offset_to_position(
                        doc_text_v1, doc_text_v1.index("args.filename") + 5
                    ),
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            field_completion = resp.get("result") or {}
            field_items = (
                field_completion.get("items")
                if isinstance(field_completion, dict)
                else field_completion
            ) or []
            if "filename" not in {item.get("label") for item in field_items}:
                fail(
                    summary,
                    "context_completion",
                    "args field completion is missing",
                    resp=resp,
                )
            record(summary, "context_completion", t0)

            imported_position = offset_to_position(
                doc_text_v1, doc_text_v1.index("imported(1)") + 1
            )
            rid = lsp.request(
                "textDocument/definition",
                {"textDocument": {"uri": doc_uri}, "position": imported_position},
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            imported_definition = resp.get("result") or {}
            if imported_definition.get("uri") != helper_uri:
                fail(
                    summary,
                    "cross_file_definition",
                    "definition did not resolve to imported helper",
                    resp=resp,
                )
            record(summary, "cross_file_definition", t0)

            rid = lsp.request(
                "textDocument/references",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": imported_position,
                    "context": {"includeDeclaration": True},
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            imported_references = resp.get("result") or []
            if len(imported_references) != 2 or {
                item.get("uri") for item in imported_references
            } != {doc_uri, helper_uri}:
                fail(
                    summary,
                    "cross_file_references",
                    "expected imported definition and call",
                    resp=resp,
                )
            record(summary, "cross_file_references", t0)

            rid = lsp.request(
                "textDocument/rename",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": imported_position,
                    "newName": "renamed_import",
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            imported_changes = (resp.get("result") or {}).get("documentChanges") or []
            if len(imported_changes) != 2:
                fail(
                    summary,
                    "cross_file_rename",
                    "expected edits for open root and closed helper",
                    resp=resp,
                )
            versions = {
                item.get("textDocument", {}).get("uri"): item.get(
                    "textDocument", {}
                ).get("version")
                for item in imported_changes
            }
            if versions != {doc_uri: 1, helper_uri: None}:
                fail(
                    summary,
                    "cross_file_rename",
                    "unexpected document versions",
                    versions=versions,
                )
            if any(
                edit.get("newText") != "renamed_import"
                for item in imported_changes
                for edit in item.get("edits") or []
            ):
                fail(
                    summary,
                    "cross_file_rename",
                    "unexpected replacement",
                    resp=resp,
                )
            record(summary, "cross_file_rename", t0)

            symbol_position = offset_to_position(
                doc_text_v1, doc_text_v1.index("$target") + 1
            )
            rid = lsp.request(
                "textDocument/definition",
                {"textDocument": {"uri": doc_uri}, "position": symbol_position},
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            definition = resp.get("result")
            if "error" in resp or not isinstance(definition, dict):
                fail(summary, "definition", "expected definition location", resp=resp)
            record(summary, "definition", t0)

            rid = lsp.request(
                "textDocument/references",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": symbol_position,
                    "context": {"includeDeclaration": True},
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            references = resp.get("result")
            if (
                "error" in resp
                or not isinstance(references, list)
                or len(references) != 2
            ):
                fail(summary, "references", "expected two references", resp=resp)
            record(summary, "references", t0, count=len(references))

            rid = lsp.request(
                "textDocument/documentHighlight",
                {"textDocument": {"uri": doc_uri}, "position": symbol_position},
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            highlights = resp.get("result")
            if (
                "error" in resp
                or not isinstance(highlights, list)
                or len(highlights) != 2
            ):
                fail(
                    summary,
                    "documentHighlight",
                    "expected two highlights",
                    resp=resp,
                )
            if {highlight.get("kind") for highlight in highlights} != {2, 3}:
                fail(
                    summary,
                    "documentHighlight",
                    "expected read and write highlights",
                    resp=resp,
                )
            record(summary, "documentHighlight", t0, count=len(highlights))

            rid = lsp.request(
                "textDocument/prepareRename",
                {"textDocument": {"uri": doc_uri}, "position": symbol_position},
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp or not isinstance(resp.get("result"), dict):
                fail(summary, "prepareRename", "expected rename range", resp=resp)
            record(summary, "prepareRename", t0)

            rid = lsp.request(
                "textDocument/rename",
                {
                    "textDocument": {"uri": doc_uri},
                    "position": symbol_position,
                    "newName": "renamed",
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            changes = (resp.get("result") or {}).get("documentChanges") or []
            if "error" in resp or len(changes) != 1:
                fail(
                    summary, "rename", "expected one versioned document edit", resp=resp
                )
            document_edit = changes[0]
            if document_edit.get("textDocument", {}).get("version") != 1:
                fail(summary, "rename", "rename edit is not versioned", resp=resp)
            rename_edits = document_edit.get("edits") or []
            if len(rename_edits) != 2 or any(
                edit.get("newText") != "$renamed" for edit in rename_edits
            ):
                fail(summary, "rename", "unexpected rename edits", resp=resp)
            record(summary, "rename", t0, edits=len(rename_edits))

            # Trigger diagnostics with invalid syntax.
            doc_text_v2 = 'tracepoint:syscalls:sys_enter_openat { printf("x");\n'
            lsp.notify(
                "textDocument/didChange",
                {
                    "textDocument": {"uri": doc_uri, "version": 2},
                    "contentChanges": [{"text": doc_text_v2}],
                },
            )

            def is_publish(msg):
                return (
                    isinstance(msg, dict)
                    and msg.get("method") == "textDocument/publishDiagnostics"
                    and isinstance(msg.get("params"), dict)
                    and msg["params"].get("uri") == doc_uri
                )

            msg = lsp.wait_for(is_publish, timeout_s=10.0)
            if msg.get("params", {}).get("version") != 2:
                fail(
                    summary,
                    "diagnostics_publish",
                    "diagnostics version does not match document",
                    response=msg,
                )
            diags = msg.get("params", {}).get("diagnostics")
            if not isinstance(diags, list) or len(diags) == 0:
                fail(
                    summary,
                    "diagnostics_publish",
                    "expected non-empty diagnostics",
                    msg=msg,
                )
            record(summary, "diagnostics_publish", t0, count=len(diags))

            # Fix diagnostics.
            lsp.notify(
                "textDocument/didChange",
                {
                    "textDocument": {"uri": doc_uri, "version": 3},
                    "contentChanges": [{"text": doc_text_v1}],
                },
            )
            msg = lsp.wait_for(is_publish, timeout_s=10.0)
            if msg.get("params", {}).get("version") != 3:
                fail(
                    summary,
                    "diagnostics_clear",
                    "diagnostics version does not match document",
                    response=msg,
                )
            diags = msg.get("params", {}).get("diagnostics")
            # Accept both null and empty list as "no diagnostics"
            if diags is not None and (not isinstance(diags, list) or len(diags) != 0):
                fail(
                    summary,
                    "diagnostics_clear",
                    "expected empty diagnostics",
                    response=msg,
                )
            record(summary, "diagnostics_clear", t0)

            # A stale full-document update must not replace the latest snapshot.
            lsp.notify(
                "textDocument/didChange",
                {
                    "textDocument": {"uri": doc_uri, "version": 2},
                    "contentChanges": [{"text": doc_text_v2}],
                },
            )
            rid = lsp.request(
                "textDocument/formatting",
                {
                    "textDocument": {"uri": doc_uri},
                    "options": {"tabSize": 4, "insertSpaces": True},
                },
            )
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(summary, "stale_change", "formatting returned error", resp=resp)
            stale_result = apply_text_edits(doc_text_v1, resp.get("result") or [])
            if stale_result != expected_formatted:
                fail(
                    summary,
                    "stale_change",
                    "stale change replaced latest document snapshot",
                    formatted=stale_result,
                )
            record(summary, "stale_change", t0)

            rid = lsp.request("shutdown", None)
            resp, _ = lsp.wait_for_response(rid, timeout_s=10.0)
            if "error" in resp:
                fail(summary, "shutdown", "shutdown returned error", resp=resp)
            record(summary, "shutdown", t0)

            lsp.notify("exit", None)
            record(summary, "exit", t0)

            try:
                lsp.p.wait(timeout=2.0)
            except Exception:
                pass

        summary["finished_ms"] = _now_ms()
        print(json.dumps(summary, indent=2, sort_keys=True))
    except SystemExit:
        raise
    except Exception as e:
        if lsp is not None:
            stderr = lsp.stderr_text()
        else:
            stderr = ""
        summary["ok"] = False
        summary["failed_step"] = "exception"
        summary["error"] = str(e)
        summary["traceback"] = traceback.format_exc()
        if stderr:
            summary["server_stderr"] = stderr[-4000:]
        print(json.dumps(summary, indent=2, sort_keys=True))
        sys.exit(1)
    finally:
        if lsp is not None:
            lsp.close()


if __name__ == "__main__":
    main()
