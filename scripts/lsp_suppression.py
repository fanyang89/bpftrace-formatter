#!/usr/bin/env python3
"""Check CLI/LSP suppression parity without changing the general smoke suite."""

import os
import pathlib
import subprocess
import tempfile

from lsp_smoke import LSPClient, apply_text_edits


def main():
    binary = os.environ["BTFMT_PATH"]
    source = (
        "BEGIN{@before=1;\n"
        " \t// btfmt: off\r\n"
        " @x=  2;  \r\n\r\n"
        "// btfmt: on\r\n"
        "@after=3;}\n"
        "// btfmt: off\nEND{ exit( ); } \t"
    )
    with tempfile.TemporaryDirectory(prefix="btfmt-suppression-") as directory:
        root = pathlib.Path(directory)
        (root / ".btfmt.json").write_text("{}", encoding="utf-8")
        expected = subprocess.run(
            [binary, "--config", str(root / ".btfmt.json"), "-"],
            input=source.encode(),
            stdout=subprocess.PIPE,
            check=True,
        ).stdout.decode()
        client = LSPClient([binary, "lsp"], os.environ.copy())
        try:
            response, _ = client.wait_for_response(
                client.request(
                    "initialize",
                    {
                        "processId": None,
                        "rootUri": root.as_uri(),
                        "capabilities": {},
                    },
                )
            )
            assert "error" not in response, response
            client.notify("initialized", {})
            # Separate documents avoid relying on notification processing order
            # between didChange and formatting for the idempotence check.
            for index, text in enumerate([source, expected]):
                uri = (root / f"suppression-{index}.bt").as_uri()
                client.notify(
                    "textDocument/didOpen",
                    {
                        "textDocument": {
                            "uri": uri,
                            "languageId": "bpftrace",
                            "version": 1,
                            "text": text,
                        }
                    },
                )
                client.wait_for(
                    lambda message: (
                        message.get("method") == "textDocument/publishDiagnostics"
                        and message.get("params", {}).get("uri") == uri
                    )
                )
                response, _ = client.wait_for_response(
                    client.request(
                        "textDocument/formatting",
                        {
                            "textDocument": {"uri": uri},
                            "options": {
                                "tabSize": 4,
                                "insertSpaces": True,
                            },
                        },
                    )
                )
                assert "error" not in response, response
                assert apply_text_edits(text, response["result"] or []) == expected
            response, _ = client.wait_for_response(client.request("shutdown", None))
            assert "error" not in response, response
            client.notify("exit", None)
        finally:
            client.close()


if __name__ == "__main__":
    main()
