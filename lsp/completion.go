package lsp

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/fanyang89/bpftrace-formatter/parser"
)

// CompletionItemKind constants for convenience
var (
	kindFunction = protocol.CompletionItemKindFunction
	kindVariable = protocol.CompletionItemKindVariable
	kindKeyword  = protocol.CompletionItemKindKeyword
	kindConstant = protocol.CompletionItemKindConstant
	kindEvent    = protocol.CompletionItemKindEvent
)

// builtinFunctions defines bpftrace built-in functions with documentation
var builtinFunctions = []struct {
	name   string
	detail string
	doc    string
}{
	{"printf", "printf(fmt, ...)", "Print formatted output"},
	{"print", "print(expr)", "Print value"},
	{"str", "str(ptr[, len])", "Convert pointer to string"},
	{"ksym", "ksym(addr)", "Kernel symbol name for address"},
	{"usym", "usym(addr)", "User symbol name for address"},
	{"kaddr", "kaddr(name)", "Kernel symbol address"},
	{"uaddr", "uaddr(name)", "User symbol address"},
	{"reg", "reg(name)", "Get register value"},
	{"system", "system(cmd)", "Execute shell command"},
	{"exit", "exit()", "Terminate bpftrace"},
	{"cgroupid", "cgroupid(path)", "Get cgroup ID"},
	{"ntop", "ntop([af,] addr)", "Convert IP address to string"},
	{"pton", "pton(str)", "Convert string to IP address"},
	{"kstack", "kstack([limit])", "Kernel stack trace"},
	{"ustack", "ustack([limit])", "User stack trace"},
	{"cat", "cat(filename)", "Print file contents"},
	{"signal", "signal(sig)", "Send signal to process"},
	{"strncmp", "strncmp(s1, s2, n)", "Compare strings"},
	{"override", "override(rc)", "Override return value (kprobes)"},
	{"buf", "buf(ptr, len)", "Get buffer as hex string"},
	{"sizeof", "sizeof(type)", "Size of type"},
	{"strftime", "strftime(fmt, nsecs)", "Format timestamp"},
	{"join", "join(arr[, sep])", "Join array elements"},
	{"time", "time(fmt)", "Print formatted time"},
	{"kptr", "kptr(addr)", "Cast to kernel pointer"},
	{"uptr", "uptr(addr)", "Cast to user pointer"},
	{"macaddr", "macaddr(addr)", "Format MAC address"},
	{"bswap", "bswap(n)", "Byte swap"},
	{"skboutput", "skboutput(path, skb, len, offset)", "Output skb to pcap file"},
	{"path", "path(struct path *)", "Get filesystem path"},
	{"unwatch", "unwatch(addr)", "Remove watchpoint"},
	{"nsecs", "nsecs", "Current timestamp in nanoseconds"},
	{"elapsed", "elapsed", "Nanoseconds since bpftrace start"},
	{"pid", "pid", "Process ID"},
	{"tid", "tid", "Thread ID"},
	{"uid", "uid", "User ID"},
	{"gid", "gid", "Group ID"},
	{"cgroup", "cgroup", "Cgroup ID"},
	{"cpu", "cpu", "CPU ID"},
	{"comm", "comm", "Process name"},
	{"curtask", "curtask", "Current task struct"},
	{"rand", "rand", "Random number"},
	{"ctx", "ctx", "Probe context"},
	{"args", "args", "Probe arguments struct"},
	{"retval", "retval", "Return value (kretprobe/uretprobe)"},
	{"func", "func", "Current function name"},
	{"probe", "probe", "Current probe name"},
	{"username", "username", "Current username"},
	{"arg0", "arg0", "First argument"},
	{"arg1", "arg1", "Second argument"},
	{"arg2", "arg2", "Third argument"},
	{"arg3", "arg3", "Fourth argument"},
	{"arg4", "arg4", "Fifth argument"},
	{"arg5", "arg5", "Sixth argument"},
	{"arg6", "arg6", "Seventh argument"},
	{"arg7", "arg7", "Eighth argument"},
	{"arg8", "arg8", "Ninth argument"},
	{"arg9", "arg9", "Tenth argument"},
	{"sarg0", "sarg0", "First stack argument"},
	{"sarg1", "sarg1", "Second stack argument"},
	{"sarg2", "sarg2", "Third stack argument"},
	{"sarg3", "sarg3", "Fourth stack argument"},
	{"sarg4", "sarg4", "Fifth stack argument"},
}

// mapFunctions defines map aggregation functions
var mapFunctions = []struct {
	name   string
	detail string
	doc    string
}{
	{"count", "count()", "Count occurrences"},
	{"sum", "sum(n)", "Sum values"},
	{"avg", "avg(n)", "Average value"},
	{"min", "min(n)", "Minimum value"},
	{"max", "max(n)", "Maximum value"},
	{"stats", "stats(n)", "Statistics (count, avg, total)"},
	{"hist", "hist(n)", "Power-of-2 histogram"},
	{"lhist", "lhist(n, min, max, step)", "Linear histogram"},
	{"delete", "delete(@map[key])", "Delete map entry"},
	{"clear", "clear(@map)", "Clear all map entries"},
	{"zero", "zero(@map)", "Zero all map values"},
	{"len", "len(@map)", "Number of map entries"},
}

// probeTypes defines bpftrace probe types
var probeTypes = []struct {
	name    string
	detail  string
	doc     string
	example string
}{
	{"BEGIN", "BEGIN", "Run once at start", "BEGIN { }"},
	{"END", "END", "Run once at end", "END { }"},
	{"kprobe", "kprobe:function", "Kernel function entry", "kprobe:vfs_read { }"},
	{"kretprobe", "kretprobe:function", "Kernel function return", "kretprobe:vfs_read { }"},
	{"uprobe", "uprobe:binary:function", "User function entry", "uprobe:/bin/bash:readline { }"},
	{"uretprobe", "uretprobe:binary:function", "User function return", "uretprobe:/bin/bash:readline { }"},
	{"tracepoint", "tracepoint:category:name", "Kernel tracepoint", "tracepoint:syscalls:sys_enter_read { }"},
	{"usdt", "usdt:binary:probe", "User static probe", "usdt:/usr/lib/libc.so.6:probe { }"},
	{"profile", "profile:hz:rate", "CPU profiling", "profile:hz:99 { }"},
	{"interval", "interval:s:duration", "Timed intervals", "interval:s:1 { }"},
	{"software", "software:event:count", "Software events", "software:page-faults:100 { }"},
	{"hardware", "hardware:event:count", "Hardware events", "hardware:cache-misses:1000000 { }"},
	{"watchpoint", "watchpoint:addr:len:mode", "Memory watchpoint", "watchpoint:0x1234:4:rw { }"},
	{"asyncwatchpoint", "asyncwatchpoint:addr:len:mode", "Asynchronous memory watchpoint", "asyncwatchpoint:0x1234:4:rw { }"},
}

// keywords defines bpftrace keywords
var keywords = []string{
	"if", "else", "while", "for", "return", "unroll", "sizeof",
}

// CompletionForPosition returns completion items for a document position.
func CompletionForPosition(doc *Document, pos protocol.Position) []protocol.CompletionItem {
	if doc == nil {
		return defaultCompletions()
	}

	// Determine context based on cursor position
	context := determineCompletionContext(doc, pos)

	switch context.kind {
	case contextProbeStart:
		return probeTypeCompletions()
	case contextMapName:
		return mapCompletions(doc, context.prefix)
	case contextVariable:
		return variableCompletions(doc, context.prefix)
	case contextFunctionCall:
		return functionCompletions(context.prefix)
	case contextMapFunction:
		return mapFunctionCompletions(context.prefix)
	case contextStatement:
		return statementCompletions(context.prefix)
	default:
		return defaultCompletions()
	}
}

type completionContextKind int

const (
	contextUnknown completionContextKind = iota
	contextProbeStart
	contextMapName
	contextVariable
	contextFunctionCall
	contextMapFunction
	contextStatement
)

type completionContext struct {
	kind   completionContextKind
	prefix string
}

func determineCompletionContext(doc *Document, pos protocol.Position) completionContext {
	if doc == nil || doc.Text == "" {
		return completionContext{kind: contextProbeStart}
	}

	// Get the text before cursor
	offset := offsetForPosition(doc.Text, pos)
	if offset <= 0 {
		return completionContext{kind: contextProbeStart}
	}

	textBefore := doc.Text[:offset]
	lines := strings.Split(textBefore, "\n")
	if len(lines) == 0 {
		return completionContext{kind: contextProbeStart}
	}

	currentLine := lines[len(lines)-1]
	trimmed := strings.TrimLeft(currentLine, " \t")

	// Check if inside a block first to get proper context
	insideBlock := isInsideBlock(textBefore)

	// Check for @ prefix (map) - only if cursor is right after @ or typing map name
	if prefix, ok := markerPrefixInCode(trimmed, '@'); ok {
		return completionContext{kind: contextMapName, prefix: prefix}
	}

	// Check for $ prefix (variable) - only if cursor is right after $ or typing variable name
	if prefix, ok := markerPrefixInCode(trimmed, '$'); ok {
		return completionContext{kind: contextVariable, prefix: prefix}
	}

	// Check if we're right after = in a map assignment (for map function completion)
	// Pattern: @map = <cursor> or @map[key] = <cursor>
	if insideBlock {
		if prefix, ok := getMapAssignmentPrefix(trimmed); ok {
			return completionContext{kind: contextMapFunction, prefix: prefix}
		}
	}

	// Check if we're at the start of a probe definition
	if isProbeContext(doc, pos, textBefore) {
		return completionContext{kind: contextProbeStart, prefix: trimmed}
	}

	// Check if inside a block (for statements/function calls)
	if insideBlock {
		lastWord := extractLastWord(trimmed)
		return completionContext{kind: contextStatement, prefix: lastWord}
	}

	if trimmed != "" {
		lastWord := extractLastWord(trimmed)
		return completionContext{kind: contextStatement, prefix: lastWord}
	}

	return completionContext{kind: contextUnknown}
}

// getMapAssignmentPrefix checks if cursor is right after = in a map assignment
// and returns the prefix being typed (empty string if just after =).
// Returns (prefix, true) if in map assignment context, ("", false) otherwise.
func getMapAssignmentPrefix(line string) (string, bool) {
	if isInStringOrComment(line) {
		return "", false
	}

	// Find the last standalone assignment operator (=), not comparison/compound operators.
	eqIdx, ok := findLastAssignmentOperator(line)
	if !ok {
		return "", false
	}

	// Check the assignment LHS is a map lvalue (@name or @name[...]).
	beforeEq := line[:eqIdx]
	if !hasMapLValueBeforeAssignment(beforeEq) {
		return "", false
	}

	// Get what's after the =
	afterEq := strings.TrimLeft(line[eqIdx+1:], " \t")

	// If empty, user just typed = and we should show all map functions
	if afterEq == "" {
		return "", true
	}

	// Check if user is typing a simple identifier (potential map function name)
	// Should not contain operators, parentheses with content, etc.
	for _, r := range afterEq {
		if !isWordChar(r) {
			return "", false
		}
	}

	if !hasMapFunctionPrefix(afterEq) {
		return "", false
	}

	return afterEq, true
}

func hasMapFunctionPrefix(prefix string) bool {
	for _, f := range mapFunctions {
		if strings.HasPrefix(f.name, prefix) {
			return true
		}
	}

	return false
}

func isInStringOrComment(line string) bool {
	inString := false
	stringDelimiter := byte(0)
	escaped := false

	for i := 0; i < len(line); i++ {
		c := line[i]

		if inString {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == stringDelimiter {
				inString = false
				stringDelimiter = 0
			}
			continue
		}

		if c == '/' && i+1 < len(line) && line[i+1] == '/' {
			return true
		}

		if c == '"' || c == '\'' {
			inString = true
			stringDelimiter = c
		}
	}

	return inString
}

func findLastAssignmentOperator(line string) (int, bool) {
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] != '=' {
			continue
		}

		if i > 0 {
			switch line[i-1] {
			case '=', '!', '<', '>', '+', '-', '*', '/', '%':
				continue
			}
		}

		if i+1 < len(line) && line[i+1] == '=' {
			continue
		}

		return i, true
	}

	return 0, false
}

func hasMapLValueBeforeAssignment(beforeEq string) bool {
	trimmed := strings.TrimRight(beforeEq, " \t")
	for i := 0; i < len(trimmed); i++ {
		if trimmed[i] != '@' {
			continue
		}
		if isMapLValue(trimmed[i:]) {
			return true
		}
	}

	return false
}

func isMapLValue(candidate string) bool {
	if candidate == "" || candidate[0] != '@' {
		return false
	}

	i := 1
	for i < len(candidate) && isWordChar(rune(candidate[i])) {
		i++
	}
	hasName := i > 1
	i = skipInlineSpaces(candidate, i)

	if !hasName {
		if i == len(candidate) {
			return true
		}
		if candidate[i] != '[' {
			return false
		}
	}

	for i < len(candidate) {
		i = skipInlineSpaces(candidate, i)
		if i == len(candidate) {
			return true
		}
		if candidate[i] != '[' {
			return false
		}

		depth := 1
		i++
		for i < len(candidate) && depth > 0 {
			switch candidate[i] {
			case '[':
				depth++
			case ']':
				depth--
			}
			i++
		}
		if depth != 0 {
			return false
		}
	}

	return true
}

func skipInlineSpaces(s string, i int) int {
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return i
}

func isProbeContext(doc *Document, _ protocol.Position, textBefore string) bool {
	if braceDepth(textBefore) > 0 {
		return false
	}

	lines := strings.Split(textBefore, "\n")
	if len(lines) == 0 {
		return true
	}

	currentLine := strings.TrimLeft(lines[len(lines)-1], " \t")
	if currentLine == "" {
		return true
	}

	if strings.Contains(currentLine, "{") {
		return false
	}

	if hasPredicateStart(currentLine) {
		return false
	}

	tokenEnd := strings.IndexAny(currentLine, ": \t")
	if tokenEnd < 0 {
		tokenEnd = len(currentLine)
	}
	token := currentLine[:tokenEnd]
	if hasProbeTypePrefix(token) {
		return true
	}

	// Also check using AST if available
	if doc.ParseResult != nil && doc.ParseResult.Tree != nil {
		// Outside top-level declaration sites, do not force probe-type completion.
		return false
	}

	return false
}

func hasProbeTypePrefix(token string) bool {
	for _, p := range probeTypes {
		if strings.HasPrefix(p.name, token) {
			return true
		}
	}

	return false
}

func hasPredicateStart(line string) bool {
	for i := 1; i < len(line); i++ {
		if line[i] != '/' {
			continue
		}
		if line[i-1] == ' ' || line[i-1] == '\t' {
			return true
		}
	}

	return false
}

func markerPrefixInCode(line string, marker byte) (string, bool) {
	last := -1
	inString := false
	stringDelimiter := byte(0)
	inLineComment := false
	escaped := false

	for i := 0; i < len(line); i++ {
		c := line[i]

		if inLineComment {
			break
		}

		if inString {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == stringDelimiter {
				inString = false
				stringDelimiter = 0
			}
			continue
		}

		if c == '/' && i+1 < len(line) && line[i+1] == '/' {
			inLineComment = true
			i++
			continue
		}

		if c == '"' || c == '\'' {
			inString = true
			stringDelimiter = c
			continue
		}

		if c == marker {
			last = i
		}
	}

	if last < 0 {
		return "", false
	}

	prefix := line[last+1:]
	for _, r := range prefix {
		if !isWordChar(r) {
			return "", false
		}
	}

	return prefix, true
}

func isInsideBlock(textBefore string) bool {
	return braceDepth(textBefore) > 0
}

func braceDepth(text string) int {
	depth := 0
	inString := false
	stringDelimiter := byte(0)
	inLineComment := false
	escaped := false

	for i := 0; i < len(text); i++ {
		c := text[i]

		if inLineComment {
			if c == '\n' {
				inLineComment = false
			}
			continue
		}

		if inString {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == stringDelimiter {
				inString = false
				stringDelimiter = 0
			}
			continue
		}

		if c == '/' && i+1 < len(text) && text[i+1] == '/' {
			inLineComment = true
			i++
			continue
		}

		if c == '"' || c == '\'' {
			inString = true
			stringDelimiter = c
			continue
		}

		switch c {
		case '{':
			depth++
		case '}':
			depth--
		}
	}

	return depth
}

func extractLastWord(line string) string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return ""
	}
	lastWord := words[len(words)-1]

	end := len(lastWord)
	for end > 0 && !isWordChar(rune(lastWord[end-1])) {
		end--
	}
	if end == 0 {
		return ""
	}

	start := end - 1
	for start >= 0 && isWordChar(rune(lastWord[start])) {
		start--
	}

	return lastWord[start+1 : end]
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

func offsetForPosition(text string, pos protocol.Position) int {
	var line uint32
	var offset int
	for i, r := range text {
		if line == pos.Line {
			// Count UTF-16 code units to match pos.Character
			var col uint32
			for j, r2 := range text[i:] {
				if col >= pos.Character {
					return i + j
				}
				if r2 == '\n' {
					return i + j
				}
				col += utf16ColumnWidth(r2)
			}
			return len(text)
		}
		if r == '\n' {
			line++
		}
		offset = i
	}
	_ = offset
	return len(text)
}

func probeTypeCompletions() []protocol.CompletionItem {
	items := make([]protocol.CompletionItem, 0, len(probeTypes))
	for _, p := range probeTypes {
		doc := p.doc + "\n\nExample:\n```bpftrace\n" + p.example + "\n```"
		items = append(items, protocol.CompletionItem{
			Label:         p.name,
			Kind:          &kindEvent,
			Detail:        &p.detail,
			Documentation: doc,
		})
	}
	return items
}

func mapCompletions(doc *Document, prefix string) []protocol.CompletionItem {
	// Collect existing map names from the document
	maps := collectMapNames(doc)
	items := make([]protocol.CompletionItem, 0, len(maps))

	for _, name := range maps {
		if prefix == "" || strings.HasPrefix(name, prefix) {
			detail := "map"
			insertText := name // Insert only the name, not the @ prefix
			items = append(items, protocol.CompletionItem{
				Label:      "@" + name,
				Kind:       &kindVariable,
				Detail:     &detail,
				InsertText: &insertText,
			})
		}
	}

	return items
}

func variableCompletions(doc *Document, prefix string) []protocol.CompletionItem {
	// Collect existing variable names from the document
	vars := collectVariableNames(doc)
	items := make([]protocol.CompletionItem, 0, len(vars))

	// Add user-defined variables
	for _, name := range vars {
		if prefix == "" || strings.HasPrefix(name, prefix) {
			detail := "variable"
			insertText := name // Insert only the name, not the $ prefix
			items = append(items, protocol.CompletionItem{
				Label:      "$" + name,
				Kind:       &kindVariable,
				Detail:     &detail,
				InsertText: &insertText,
			})
		}
	}

	return items
}

func constantCompletions(prefix string) []protocol.CompletionItem {
	items := make([]protocol.CompletionItem, 0, len(builtinFunctions))

	for _, f := range builtinFunctions {
		if strings.Contains(f.detail, "(") {
			continue
		}
		if prefix == "" || strings.HasPrefix(f.name, prefix) {
			items = append(items, protocol.CompletionItem{
				Label:         f.name,
				Kind:          &kindConstant,
				Detail:        &f.detail,
				Documentation: f.doc,
			})
		}
	}

	return items
}

func functionCompletions(prefix string) []protocol.CompletionItem {
	items := make([]protocol.CompletionItem, 0, len(builtinFunctions))

	for _, f := range builtinFunctions {
		// Only add functions (have parentheses)
		if strings.Contains(f.detail, "(") {
			if prefix == "" || strings.HasPrefix(f.name, prefix) {
				items = append(items, protocol.CompletionItem{
					Label:         f.name,
					Kind:          &kindFunction,
					Detail:        &f.detail,
					Documentation: f.doc,
					InsertText:    &f.name,
				})
			}
		}
	}

	return items
}

func mapFunctionCompletions(prefix string) []protocol.CompletionItem {
	items := make([]protocol.CompletionItem, 0, len(mapFunctions))

	for _, f := range mapFunctions {
		if prefix == "" || strings.HasPrefix(f.name, prefix) {
			items = append(items, protocol.CompletionItem{
				Label:         f.name,
				Kind:          &kindFunction,
				Detail:        &f.detail,
				Documentation: f.doc,
			})
		}
	}

	return items
}

func statementCompletions(prefix string) []protocol.CompletionItem {
	items := make([]protocol.CompletionItem, 0, len(keywords)+len(builtinFunctions))

	// Add keywords
	for _, kw := range keywords {
		if prefix == "" || strings.HasPrefix(kw, prefix) {
			items = append(items, protocol.CompletionItem{
				Label: kw,
				Kind:  &kindKeyword,
			})
		}
	}

	// Add functions
	items = append(items, functionCompletions(prefix)...)

	// Add map functions
	items = append(items, mapFunctionCompletions(prefix)...)

	// Add built-in constants (e.g. pid/tid/uid)
	items = append(items, constantCompletions(prefix)...)

	return items
}

func defaultCompletions() []protocol.CompletionItem {
	items := make([]protocol.CompletionItem, 0, len(probeTypes)+len(builtinFunctions))

	// Probe types
	items = append(items, probeTypeCompletions()...)

	// Functions
	for _, f := range builtinFunctions {
		if strings.Contains(f.detail, "(") {
			items = append(items, protocol.CompletionItem{
				Label:         f.name,
				Kind:          &kindFunction,
				Detail:        &f.detail,
				Documentation: f.doc,
			})
		}
	}

	// Constants
	items = append(items, constantCompletions("")...)

	return items
}

// collectMapNames walks the AST to find all map names used in the document
func collectMapNames(doc *Document) []string {
	if doc == nil || doc.ParseResult == nil || doc.ParseResult.Tree == nil {
		return nil
	}

	names := make(map[string]struct{})
	var walk func(node antlr.Tree)
	walk = func(node antlr.Tree) {
		if node == nil {
			return
		}

		if term, ok := node.(antlr.TerminalNode); ok {
			text := term.GetText()
			if name, found := strings.CutPrefix(text, "@"); found {
				// Remove any trailing brackets or content
				if idx := strings.Index(name, "["); idx > 0 {
					name = name[:idx]
				}
				if name != "" {
					names[name] = struct{}{}
				}
			}
		}

		for i := 0; i < node.GetChildCount(); i++ {
			walk(node.GetChild(i))
		}
	}

	walk(doc.ParseResult.Tree)

	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	return result
}

// collectVariableNames walks the AST to find all variable names used in the document
func collectVariableNames(doc *Document) []string {
	if doc == nil || doc.ParseResult == nil || doc.ParseResult.Tree == nil {
		return nil
	}

	names := make(map[string]struct{})
	var walk func(node antlr.Tree)
	walk = func(node antlr.Tree) {
		if node == nil {
			return
		}

		// Check for variable context
		if ctx, ok := node.(*parser.VariableContext); ok {
			if ctx.VARIABLE() != nil {
				text := ctx.VARIABLE().GetText()
				if name, found := strings.CutPrefix(text, "$"); found && name != "" {
					names[name] = struct{}{}
				}
			}
		}

		// Also check terminal nodes directly
		if term, ok := node.(antlr.TerminalNode); ok {
			text := term.GetText()
			if name, found := strings.CutPrefix(text, "$"); found && name != "" && isValidIdentifier(name) {
				names[name] = struct{}{}
			}
		}

		for i := 0; i < node.GetChildCount(); i++ {
			walk(node.GetChild(i))
		}
	}

	walk(doc.ParseResult.Tree)

	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	return result
}

func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_') {
				return false
			}
		} else {
			if !isWordChar(r) {
				return false
			}
		}
	}
	return true
}
