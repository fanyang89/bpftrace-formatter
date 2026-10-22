use tower_lsp::lsp_types::{
    CompletionItem, CompletionItemKind, Documentation, MarkupContent, MarkupKind,
};

pub(super) const CATALOG_VERSION: &str = "tree-sitter-bpftrace 0.3.2";

pub(super) type ContextMask = u8;
pub(super) const TOP_LEVEL: ContextMask = 1 << 0;
pub(super) const PROBE_PROVIDER: ContextMask = 1 << 1;
pub(super) const EXPRESSION: ContextMask = 1 << 2;
pub(super) const STATEMENT: ContextMask = 1 << 3;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub(super) enum CatalogKind {
    Provider,
    Function,
    Value,
    Keyword,
}

#[derive(Debug, Clone, Copy)]
pub(super) struct CatalogEntry {
    pub label: &'static str,
    pub kind: CatalogKind,
    pub contexts: ContextMask,
    pub detail: &'static str,
    pub documentation: &'static str,
    pub deprecated: bool,
}

impl CatalogEntry {
    pub(super) fn completion_item(self) -> CompletionItem {
        CompletionItem {
            label: self.label.to_string(),
            kind: Some(match self.kind {
                CatalogKind::Provider => CompletionItemKind::EVENT,
                CatalogKind::Function => CompletionItemKind::FUNCTION,
                CatalogKind::Value => CompletionItemKind::VARIABLE,
                CatalogKind::Keyword => CompletionItemKind::KEYWORD,
            }),
            detail: Some(self.detail.to_string()),
            documentation: Some(Documentation::MarkupContent(MarkupContent {
                kind: MarkupKind::Markdown,
                value: self.documentation.to_string(),
            })),
            deprecated: self.deprecated.then_some(true),
            ..CompletionItem::default()
        }
    }
}

pub(super) fn entries() -> &'static [CatalogEntry] {
    ENTRIES
}

pub(super) fn find(label: &str) -> Option<&'static CatalogEntry> {
    ENTRIES.iter().find(|entry| entry.label == label)
}

macro_rules! entry {
    ($label:literal, $kind:ident, $contexts:expr, $detail:literal, $doc:literal) => {
        CatalogEntry {
            label: $label,
            kind: CatalogKind::$kind,
            contexts: $contexts,
            detail: $detail,
            documentation: $doc,
            deprecated: false,
        }
    };
}

const ENTRIES: &[CatalogEntry] = &[
    entry!(
        "BEGIN",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Run once when bpftrace starts."
    ),
    entry!(
        "END",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Run once when bpftrace exits."
    ),
    entry!(
        "bench",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Benchmark probe provider."
    ),
    entry!(
        "fentry",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Kernel function entry probe using BTF."
    ),
    entry!(
        "fexit",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Kernel function exit probe using BTF."
    ),
    entry!(
        "hardware",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Hardware performance counter probe."
    ),
    entry!(
        "interval",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Periodic interval probe."
    ),
    entry!(
        "iter",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Kernel iterator probe."
    ),
    entry!(
        "kprobe",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Kernel function entry probe."
    ),
    entry!(
        "kretprobe",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Kernel function return probe."
    ),
    entry!(
        "profile",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Timed profile probe."
    ),
    entry!(
        "rawtracepoint",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Raw kernel tracepoint probe."
    ),
    entry!(
        "software",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Software performance counter probe."
    ),
    entry!(
        "test",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Test probe provider."
    ),
    entry!(
        "tracepoint",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Kernel tracepoint probe."
    ),
    entry!(
        "uprobe",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Userspace function entry probe."
    ),
    entry!(
        "uretprobe",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Userspace function return probe."
    ),
    entry!(
        "usdt",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Userspace statically defined tracepoint."
    ),
    entry!(
        "watchpoint",
        Provider,
        TOP_LEVEL | PROBE_PROVIDER,
        "probe provider",
        "Memory watchpoint probe."
    ),
    entry!(
        "avg",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Calculate an average aggregation."
    ),
    entry!(
        "cat",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Print the contents of a file."
    ),
    entry!(
        "clear",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Remove all values from a map."
    ),
    entry!(
        "count",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Count occurrences in a map aggregation."
    ),
    entry!(
        "delete",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Delete an element from a map."
    ),
    entry!(
        "exit",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Terminate tracing."
    ),
    entry!(
        "hist",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Build a power-of-two histogram."
    ),
    entry!(
        "join",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Join an array of strings."
    ),
    entry!(
        "kaddr",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Resolve a kernel symbol address."
    ),
    entry!(
        "kstack",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Capture a kernel stack trace."
    ),
    entry!(
        "ksym",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Resolve a kernel address to a symbol."
    ),
    entry!(
        "lhist",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Build a linear histogram."
    ),
    entry!(
        "max",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Calculate a maximum aggregation."
    ),
    entry!(
        "min",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Calculate a minimum aggregation."
    ),
    entry!(
        "print",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Print a value or map."
    ),
    entry!(
        "printf",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Print formatted output."
    ),
    entry!(
        "reg",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Read a register by name."
    ),
    entry!(
        "sprintf",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Format a string."
    ),
    entry!(
        "stats",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Calculate count, average, and total statistics."
    ),
    entry!(
        "str",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Read a string from a pointer."
    ),
    entry!(
        "strcmp",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Compare strings."
    ),
    entry!(
        "strlen",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Return string length."
    ),
    entry!(
        "sum",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Calculate a sum aggregation."
    ),
    entry!(
        "sym",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Resolve an address to a symbol."
    ),
    entry!(
        "system",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Execute a shell command asynchronously."
    ),
    entry!(
        "time",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Print the current time."
    ),
    entry!(
        "uaddr",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Resolve a userspace symbol address."
    ),
    entry!(
        "ustack",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Capture a userspace stack trace."
    ),
    entry!(
        "usym",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Resolve a userspace address to a symbol."
    ),
    entry!(
        "zero",
        Function,
        EXPRESSION | STATEMENT,
        "builtin function",
        "Set all map values to zero."
    ),
    entry!(
        "args",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Probe arguments exposed by BTF or tracepoint metadata."
    ),
    entry!(
        "comm",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Current task command name."
    ),
    entry!(
        "cpu",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Current CPU identifier."
    ),
    entry!(
        "curtask",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Pointer to the current task."
    ),
    entry!(
        "elapsed",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Elapsed nanoseconds since tracing started."
    ),
    entry!(
        "gid",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Current group identifier."
    ),
    entry!(
        "nsecs",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Current timestamp in nanoseconds."
    ),
    entry!(
        "pid",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Current process identifier."
    ),
    entry!(
        "retval",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Function return value."
    ),
    entry!(
        "tid",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Current thread identifier."
    ),
    entry!(
        "uid",
        Value,
        EXPRESSION | STATEMENT,
        "builtin value",
        "Current user identifier."
    ),
    entry!(
        "break",
        Keyword,
        STATEMENT,
        "keyword",
        "Exit the nearest loop."
    ),
    entry!(
        "config",
        Keyword,
        TOP_LEVEL,
        "keyword",
        "Start the root script configuration block."
    ),
    entry!(
        "continue",
        Keyword,
        STATEMENT,
        "keyword",
        "Continue the nearest loop."
    ),
    entry!("for", Keyword, STATEMENT, "keyword", "Start a loop."),
    entry!(
        "if",
        Keyword,
        EXPRESSION | STATEMENT,
        "keyword",
        "Start a conditional block."
    ),
    entry!(
        "import",
        Keyword,
        TOP_LEVEL,
        "keyword",
        "Import a bpftrace script or supported source file."
    ),
    entry!(
        "let",
        Keyword,
        TOP_LEVEL | STATEMENT,
        "keyword",
        "Declare a map or scratch variable."
    ),
    entry!(
        "macro",
        Keyword,
        TOP_LEVEL,
        "keyword",
        "Define a hygienic macro."
    ),
    entry!(
        "return",
        Keyword,
        STATEMENT,
        "keyword",
        "Return a value from a macro body."
    ),
    entry!(
        "while",
        Keyword,
        STATEMENT,
        "keyword",
        "Start a while loop."
    ),
];

#[derive(Debug)]
pub(super) struct Signature {
    pub parameters: &'static [(&'static str, &'static str)],
    pub variadic: bool,
}

// Signatures follow https://bpftrace.org/docs/release_024/stdlib.
// Optional parameters are bracketed; unlisted functions intentionally have no help.
pub(super) fn signature(name: &str) -> Option<Signature> {
    let parameters: &'static [(&'static str, &'static str)] = match name {
        "avg" | "max" | "min" | "stats" | "sum" => &[("int64 n", "Value to aggregate.")],
        "count" => &[],
        "cat" | "printf" | "system" => &[
            ("string fmt", "Constant format string."),
            ("...args", "Values substituted into the format string."),
        ],
        "clear" | "zero" => &[("map m", "Map to modify.")],
        "delete" => &[
            (
                "map m",
                "Map to modify; legacy indexed-map syntax is also supported.",
            ),
            ("[mapkey key]", "Key to delete; omit for scalar maps."),
        ],
        "exit" => &[("[int code]", "Optional process exit code.")],
        "hist" => &[
            ("int64 n", "Value to record in the histogram."),
            (
                "[int k]",
                "Log2 buckets per power of two, from 0 to 5; defaults to 0.",
            ),
        ],
        "lhist" => &[
            ("int64 n", "Value to record in the histogram."),
            ("int64 min", "Lower bound of the histogram."),
            ("int64 max", "Upper bound of the histogram."),
            ("int64 step", "Bucket width."),
        ],
        "join" => &[
            ("char *arr[]", "Array of strings to print."),
            ("[string sep]", "Separator; defaults to a space."),
        ],
        "kaddr" | "uaddr" => &[("string name", "Constant symbol name to resolve.")],
        "ksym" | "usym" => &[("uint64 addr", "Address to symbolize.")],
        "print" => &[
            ("value", "Value or map to print."),
            (
                "[int64 top]",
                "For maps only: print the highest top entries.",
            ),
            (
                "[int64 div]",
                "For maps only: divide values by this factor when printing.",
            ),
        ],
        "reg" => &[("string name", "Architecture-specific register name.")],
        "str" => &[
            ("char *data", "Pointer to a NUL-terminated string."),
            (
                "[uint32 length]",
                "Maximum length including the terminating NUL.",
            ),
        ],
        "time" => &[("string fmt", "strftime-style wall-clock format string.")],
        _ => return None,
    };
    Some(Signature {
        parameters,
        variadic: matches!(name, "cat" | "printf" | "system"),
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn signatures_reference_documented_builtin_functions() {
        let mut count = 0;
        for entry in entries() {
            if let Some(signature) = signature(entry.label) {
                count += 1;
                assert_eq!(entry.kind, CatalogKind::Function);
                assert!(!entry.documentation.is_empty());
                assert!(!signature.variadic || !signature.parameters.is_empty());
                for (label, documentation) in signature.parameters {
                    assert!(!label.is_empty());
                    assert!(!documentation.is_empty());
                }
            }
        }
        assert_eq!(count, 24);
    }

    #[test]
    fn catalog_labels_are_unique_and_versioned() {
        let mut labels: Vec<_> = entries().iter().map(|entry| entry.label).collect();
        labels.sort_unstable();
        labels.dedup();
        assert_eq!(labels.len(), entries().len());
        assert_eq!(CATALOG_VERSION, "tree-sitter-bpftrace 0.3.2");
    }
}
