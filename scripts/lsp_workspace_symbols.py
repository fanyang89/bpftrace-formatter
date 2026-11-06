#!/usr/bin/env python3
"""Workspace symbol protocol and lifecycle coverage using the shared stdio client."""

import os
import pathlib
import tempfile
import threading

from lsp_smoke import LSPClient, offset_to_position


def request(client, method, params):
    response, _ = client.wait_for_response(client.request(method, params))
    assert "error" not in response, response
    return response["result"]


def query(client, text=""):
    return request(client, "workspace/symbol", {"query": text})


def folder(path):
    return {"uri": path.as_uri(), "name": path.name}


def open_document(client, path, text, version=1):
    client.notify(
        "textDocument/didOpen",
        {
            "textDocument": {
                "uri": path.as_uri(),
                "languageId": "bpftrace",
                "version": version,
                "text": text,
            }
        },
    )
    diagnostics(client, path, version)


def diagnostics(client, path, version):
    client.wait_for(
        lambda message: (
            message.get("method") == "textDocument/publishDiagnostics"
            and message["params"]["uri"] == path.as_uri()
            and message["params"].get("version") == version
        )
    )


def watched(client, path, kind):
    client.notify(
        "workspace/didChangeWatchedFiles",
        {
            "changes": [
                {"uri": path.as_uri(), "type": kind},
            ]
        },
    )


def run(binary, root, trusted):
    first, second, outside = (root / name for name in ("first", "second", "outside"))
    for directory in (first, second, outside):
        directory.mkdir()
    main = first / "main.bt"
    helper = first / "helper.bt"
    other = second / "other.bt"
    external = outside / "external.bt"
    main_text = (
        'import "helper.bt";\r\n'
        'BEGIN { printf("😀"); @Count = Measure(1); @Count++; Measure(2); }\r\n'
    )
    main.write_text(main_text)
    helper.write_text("macro Measure(value) { value + 1 }\n")
    other.write_text("BEGIN { @Count = 2; }\n")
    external.write_text("BEGIN { @external = 1; }\n")
    (first / ".ignore").write_text("ignored*.bt\n")
    ignored = first / "ignored.bt"
    ignored.write_text("BEGIN { @ignored = 1; }\n")

    client = LSPClient([binary, "lsp"], dict(os.environ))
    watchdog = threading.Timer(30, client.close)
    watchdog.start()
    try:
        initialized = request(
            client,
            "initialize",
            {
                "processId": os.getpid(),
                "capabilities": {},
                "workspaceFolders": [folder(first), folder(second)],
                "initializationOptions": {"btfmt": {"trusted": trusted}},
            },
        )
        assert initialized["capabilities"]["workspaceSymbolProvider"] is True
        client.notify("initialized", {})
        open_document(client, main, main_text)
        open_document(client, external, external.read_text())
        open_document(client, ignored, ignored.read_text())

        if not trusted:
            assert query(client) == []
            assert request(
                client,
                "textDocument/documentSymbol",
                {
                    "textDocument": {"uri": main.as_uri()},
                },
            )
            watched(client, helper, 2)
            client.notify(
                "workspace/didChangeWorkspaceFolders",
                {
                    "event": {
                        "added": [folder(outside)],
                        "removed": [],
                    }
                },
            )
            assert query(client, "Count") == []
        else:
            all_symbols = query(client)
            assert all_symbols == query(client, "   ")
            assert [symbol["name"] for symbol in all_symbols] == [
                "@Count",
                "@Count",
                "BEGIN",
                "BEGIN",
                "Measure",
            ], all_symbols
            assert query(client, "NoSuchSymbol") == []
            count = query(client, "  cOuNt ")
            assert len(count) == 2 and all(symbol["kind"] == 13 for symbol in count)
            first_count = count[0]
            start = main_text.index("@Count")
            assert first_count["location"] == {
                "uri": main.as_uri(),
                "range": {
                    "start": offset_to_position(main_text, start),
                    "end": offset_to_position(main_text, start + len("@Count")),
                },
            }, first_count
            macro = query(client, "MEAS")
            assert len(macro) == 1 and macro[0]["kind"] == 12
            assert macro[0]["location"]["uri"] == helper.as_uri()

            # Full-document changes override disk data; stale versions cannot restore it.
            client.notify(
                "textDocument/didChange",
                {
                    "textDocument": {"uri": main.as_uri(), "version": 2},
                    "contentChanges": [{"text": "BEGIN { @unsaved = 1; }"}],
                },
            )
            diagnostics(client, main, 2)
            client.notify(
                "textDocument/didChange",
                {
                    "textDocument": {"uri": main.as_uri(), "version": 1},
                    "contentChanges": [{"text": main_text}],
                },
            )
            assert len(query(client, "@unsaved")) == 1
            assert len(query(client, "@Count")) == 1
            main.write_text("BEGIN { @saved = 1; }")
            client.notify(
                "textDocument/didSave", {"textDocument": {"uri": main.as_uri()}}
            )
            assert query(client, "@saved") == []
            client.notify(
                "textDocument/didClose", {"textDocument": {"uri": main.as_uri()}}
            )
            diagnostics(client, main, None)
            assert len(query(client, "@saved")) == 1
            assert query(client, "@unsaved") == []

            created = first / "created.bt"
            open_document(client, created, "BEGIN { @new_buffer = 1; }")
            assert len(query(client, "@new_buffer")) == 1
            created.write_text("BEGIN { @created = 1; }")
            watched(client, created, 1)
            assert query(client, "@created") == []
            created.unlink()
            watched(client, created, 3)
            assert len(query(client, "@new_buffer")) == 1
            client.notify(
                "textDocument/didClose", {"textDocument": {"uri": created.as_uri()}}
            )
            diagnostics(client, created, None)
            assert query(client, "@new_buffer") == []

            created.write_text("BEGIN { @created = 1; }")
            watched(client, created, 1)
            assert len(query(client, "@created")) == 1
            created.write_text("BEGIN { @changed = 1; }")
            watched(client, created, 2)
            assert query(client, "@created") == []
            assert len(query(client, "@changed")) == 1
            created.unlink()
            watched(client, created, 3)
            assert query(client, "@changed") == []
            watched(client, ignored, 2)
            assert query(client, "@ignored") == []
            assert query(client, "@external") == []

            open_document(client, other, "BEGIN { @second_buffer = 1; }")
            client.notify(
                "workspace/didChangeWorkspaceFolders",
                {
                    "event": {
                        "added": [],
                        "removed": [folder(second)],
                    }
                },
            )
            assert query(client, "@second_buffer") == []
            assert all(
                symbol["location"]["uri"].startswith(first.as_uri() + "/")
                for symbol in query(client)
            )
            client.notify(
                "workspace/didChangeWorkspaceFolders",
                {
                    "event": {
                        "added": [folder(second)],
                        "removed": [],
                    }
                },
            )
            assert len(query(client, "@second_buffer")) == 1

        response, _ = client.wait_for_response(
            client.request("workspace/symbol", {"query": 42})
        )
        assert response["error"]["code"] == -32602, response
        request(client, "shutdown", None)
        client.notify("exit", None)
        client.p.stdin.close()
        assert client.p.wait(timeout=5) == 0
    finally:
        watchdog.cancel()
        client.close()


if __name__ == "__main__":
    binary = os.environ["BTFMT_PATH"]
    for trusted in (True, False):
        with tempfile.TemporaryDirectory(
            prefix="btfmt-workspace-symbols-"
        ) as directory:
            run(binary, pathlib.Path(directory), trusted)
    print("workspace symbol protocol tests passed")
