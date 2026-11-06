#!/usr/bin/env python3
"""Full-document semantic tokens protocol tests, including Restricted Mode."""

import os
import pathlib
import tempfile
import threading

from lsp_smoke import LSPClient


TOKEN_TYPES = [
    "variable",
    "map",
    "parameter",
    "macro",
    "function",
    "property",
    "type",
    "event",
    "comment",
    "string",
    "number",
    "keyword",
]


def request(client, method, params):
    response, _ = client.wait_for_response(client.request(method, params))
    assert "error" not in response, response
    return response["result"]


def tokens(client, uri, text):
    result = request(
        client,
        "textDocument/semanticTokens/full",
        {
            "textDocument": {"uri": uri},
        },
    )
    data = result["data"]
    assert len(data) % 5 == 0, data
    lines = text.split("\n")
    line = column = 0
    previous_end = (0, 0)
    decoded = []
    for offset in range(0, len(data), 5):
        delta_line, delta_start, length, kind, modifiers = data[offset : offset + 5]
        assert all(
            isinstance(value, int) and value >= 0 for value in data[offset : offset + 5]
        )
        line += delta_line
        column = delta_start if delta_line else column + delta_start
        assert (line, column) >= previous_end
        assert length > 0 and 0 <= kind < len(TOKEN_TYPES) and modifiers < 4
        raw_line = lines[line].removesuffix("\r").encode("utf-16-le")
        assert 2 * (column + length) <= len(raw_line)
        value = raw_line[column * 2 : (column + length) * 2].decode("utf-16-le")
        assert "\n" not in value and "\r" not in value
        decoded.append((value, TOKEN_TYPES[kind], modifiers))
        previous_end = (line, column + length)
    return decoded


def sync_document(client, uri, text, version, opening=False):
    if opening:
        client.notify(
            "textDocument/didOpen",
            {
                "textDocument": {
                    "uri": uri,
                    "languageId": "bpftrace",
                    "version": version,
                    "text": text,
                },
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
            and message["params"].get("uri") == uri
            and message["params"].get("version") == version
        )
    )


def check_session(root, trusted):
    marker = root / "bpftrace-called"
    fake_bpftrace = root / "bpftrace"
    fake_bpftrace.write_text(f"#!/bin/sh\ntouch '{marker}'\n")
    fake_bpftrace.chmod(0o755)
    (root / "other.bt").write_text(
        "macro external() { 1 }\nBEGIN { @disk = count(); }\n"
    )
    env = dict(os.environ, PATH=str(root) + os.pathsep + os.environ.get("PATH", ""))
    client = LSPClient([os.environ["BTFMT_PATH"], "lsp"], env)
    watchdog = threading.Timer(30, client.close)
    watchdog.start()
    try:
        result = request(
            client,
            "initialize",
            {
                "processId": None,
                "rootUri": root.as_uri(),
                "capabilities": {
                    "textDocument": {
                        "semanticTokens": {
                            "requests": {"full": True},
                            "tokenTypes": TOKEN_TYPES,
                            "tokenModifiers": ["declaration", "defaultLibrary"],
                            "formats": ["relative"],
                            "multilineTokenSupport": False,
                            "overlappingTokenSupport": False,
                        }
                    }
                },
                "initializationOptions": {"btfmt": {"trusted": trusted}},
            },
        )
        capability = result["capabilities"]["semanticTokensProvider"]
        assert capability["full"] is True and capability.get("range") is False
        assert capability["legend"] == {
            "tokenTypes": TOKEN_TYPES,
            "tokenModifiers": ["declaration", "defaultLibrary"],
        }
        client.notify("initialized", {})
        uri = (root / "document.bt").as_uri()
        assert tokens(client, uri, "") == []
        text = (
            "/*😀\r\n終*/\r\nmacro add(value) { value + 1 }\r\n"
            'BEGIN { let $x = pid; printf("😀中"); @total[$x] = count(); add($x); external(); }\r\n'
        )
        sync_document(client, uri, text, 1, opening=True)
        first = tokens(client, uri, text)
        for expected in [
            ("add", "macro", 1),
            ("add", "macro", 0),
            ("value", "parameter", 1),
            ("value", "parameter", 0),
            ("$x", "variable", 1),
            ("$x", "variable", 0),
            ("@total", "map", 0),
            ("pid", "variable", 2),
            ("printf", "function", 2),
            ("count", "function", 2),
            ("/*😀", "comment", 0),
            ("終*/", "comment", 0),
            ('"😀中"', "string", 0),
        ]:
            assert expected in first, (expected, first)
        assert not any(value in ("@disk", "external") for value, _, _ in first)
        assert tokens(client, uri, text) == first

        incomplete = 'BEGIN { @fresh = count(); printf("😀"); $unfinished ='
        sync_document(client, uri, incomplete, 2)
        changed = tokens(client, uri, incomplete)
        assert ("@fresh", "map", 0) in changed
        assert not any(
            value in ("@total", "$x", "value", "add") for value, _, _ in changed
        )

        sync_document(client, uri, "", 3)
        assert tokens(client, uri, "") == []
        sync_document(client, uri, "BEGIN { @reopened = 1; }", 4)
        assert ("@reopened", "map", 0) in tokens(
            client, uri, "BEGIN { @reopened = 1; }"
        )
        client.notify("textDocument/didClose", {"textDocument": {"uri": uri}})
        client.wait_for(
            lambda message: (
                message.get("method") == "textDocument/publishDiagnostics"
                and message["params"].get("uri") == uri
            )
        )
        assert tokens(client, uri, "") == []
        sync_document(client, uri, "BEGIN { @again = 1; }", 1, opening=True)
        assert ("@again", "map", 0) in tokens(client, uri, "BEGIN { @again = 1; }")
        assert not marker.exists(), "semantic highlighting executed bpftrace"
        request(client, "shutdown", None)
        client.notify("exit", None)
        client.p.stdin.close()
        assert client.p.wait(timeout=5) == 0
    finally:
        watchdog.cancel()
        client.close()


def main():
    base = (
        pathlib.Path(os.environ.get("TMPDIR", str(pathlib.Path.home() / "tmp")))
        / "opencode"
    )
    base.mkdir(parents=True, exist_ok=True)
    for trusted in (False, True):
        with tempfile.TemporaryDirectory(
            prefix="btfmt-semantic-", dir=base
        ) as directory:
            check_session(pathlib.Path(directory), trusted)
    print("Semantic tokens protocol tests passed (trusted and Restricted Mode).")


if __name__ == "__main__":
    main()
