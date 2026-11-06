#!/usr/bin/env python3
"""Exercise selection formatting over the real stdio LSP transport."""

import os
import pathlib
import tempfile

from lsp_smoke import LSPClient, offset_to_position


def main():
    client = LSPClient([os.environ["BTFMT_PATH"], "lsp"], os.environ.copy())
    try:
        with tempfile.TemporaryDirectory(prefix="btfmt-range-") as directory:
            root = pathlib.Path(directory)
            (root / ".btfmt.json").write_text('{"indent":{"size":2}}')
            uri = (root / "selection.bt").as_uri()

            def request(method, params):
                response, _ = client.wait_for_response(client.request(method, params))
                assert "error" not in response, response
                return response["result"]

            initialized = request("initialize", {"capabilities": {}})
            assert (
                initialized["capabilities"]["documentRangeFormattingProvider"] is True
            )
            client.notify("initialized", {})
            version = 0

            def document(text):
                nonlocal version
                version += 1
                if version == 1:
                    client.notify(
                        "textDocument/didOpen",
                        {
                            "textDocument": {
                                "uri": uri,
                                "languageId": "bpftrace",
                                "version": version,
                                "text": text,
                            }
                        },
                    )
                else:
                    client.notify(
                        "textDocument/didChange",
                        {
                            "textDocument": {"uri": uri, "version": version},
                            "contentChanges": [{"text": text}],
                        },
                    )
                client.wait_for(
                    lambda message: (
                        message.get("method") == "textDocument/publishDiagnostics"
                        and message["params"].get("version") == version
                    )
                )

            def params(text, selection):
                start = text.index(selection)
                return {
                    "textDocument": {"uri": uri},
                    "range": {
                        "start": offset_to_position(text, start),
                        "end": offset_to_position(text, start + len(selection)),
                    },
                    # The same repository config as full formatting wins over
                    # client options, even when the editor asks for tabs.
                    "options": {"tabSize": 8, "insertSpaces": False},
                }

            def format_selection(text, selection):
                document(text)
                return request("textDocument/rangeFormatting", params(text, selection))

            def offset(text, position):
                lines = text.splitlines(keepends=True)
                line = position["line"]
                if line == len(lines):
                    return len(text)
                prefix = lines[line].encode("utf-16-le")[: position["character"] * 2]
                return sum(map(len, lines[:line])) + len(prefix.decode("utf-16-le"))

            source = '// 中文😀\r\nBEGIN { printf("😀"); $x=1+2; // 注释😀\r\n $y= 3;\r\n}\r\n'
            edits = format_selection(source, "1+")
            assert len(edits) == 1, edits
            edit = edits[0]
            assert edit["range"] == {
                "start": {"line": 1, "character": 21},
                "end": {"line": 2, "character": 0},
            }, edit
            assert edit["newText"] == "\n  $x = 1 + 2; // 注释😀\n", edit
            start = offset(source, edit["range"]["start"])
            end = offset(source, edit["range"]["end"])
            updated = source[:start] + edit["newText"] + source[end:]
            assert updated.startswith('// 中文😀\r\nBEGIN { printf("😀");\n')
            assert updated.endswith(" $y= 3;\r\n}\r\n")
            assert format_selection(updated, "$x = 1 + 2;") == []

            source = "BEGIN{if(pid){$x=1;}else{$x=2;}}\n"
            edits = format_selection(source, source)
            full = request(
                "textDocument/formatting",
                {
                    "textDocument": {"uri": uri},
                    "options": {"tabSize": 8, "insertSpaces": False},
                },
            )
            assert edits == full

            for text, selection in [
                ("BEGIN{$x=1;}\n", ""),
                ("BEGIN{$x=;}\n", "$x"),
                ("// comment\nBEGIN{}\n", "comment"),
                ("#define X 1\nBEGIN{$x=X;}\n", "#define"),
                ("BEGIN{\n// btfmt: off\n$x=1;\n// btfmt: on\n}\n", "$x"),
            ]:
                assert format_selection(text, selection) == [], (text, selection)

            text = 'BEGIN{printf("😀");}\n'
            document(text)
            invalid = params(text, "😀")
            invalid["range"]["start"]["character"] += 1
            response, _ = client.wait_for_response(
                client.request("textDocument/rangeFormatting", invalid)
            )
            assert response["error"]["code"] == -32602, response
            request("shutdown", None)
            client.notify("exit", None)
            client.p.stdin.close()
            assert client.p.wait(timeout=5) == 0
    finally:
        client.close()


if __name__ == "__main__":
    main()
