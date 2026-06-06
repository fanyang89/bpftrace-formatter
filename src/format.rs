use crate::config::{BraceStyle, Config};
use crate::parse;
use anyhow::Result;

#[derive(Debug, Clone, PartialEq, Eq)]
enum TokenKind {
    Word,
    String,
    Comment,
    Shebang,
    Preprocessor,
    Newline,
    Symbol,
    Operator,
}

#[derive(Debug, Clone, PartialEq, Eq)]
struct Token {
    kind: TokenKind,
    text: String,
}

pub fn format_source(source: &str, config: &Config) -> Result<String> {
    parse::ensure_valid(source)?;
    Ok(format_tokens(&tokenize(source), config))
}

fn format_tokens(tokens: &[Token], config: &Config) -> String {
    let mut writer = Writer::new(config);
    let mut top_level_block_closed = false;

    for token in tokens {
        match token.kind {
            TokenKind::Newline => {}
            TokenKind::Shebang => {
                writer.write_line_token(&token.text);
                writer.blank_lines(config.line_breaks.empty_lines_after_shebang);
            }
            TokenKind::Preprocessor => writer.write_line_token(&token.text),
            TokenKind::Comment => writer.write_comment(&token.text),
            TokenKind::Symbol => match token.text.as_str() {
                "{" => {
                    writer.open_block();
                    top_level_block_closed = false;
                }
                "}" => {
                    let was_top_level = writer.indent == 1;
                    writer.close_block();
                    if was_top_level {
                        top_level_block_closed = true;
                    }
                }
                ";" => writer.semicolon(),
                "," => writer.comma(),
                "(" => writer.open_paren(),
                ")" => writer.close_paren(),
                "[" => writer.open_bracket(),
                "]" => writer.close_bracket(),
                ":" | "." | "->" => writer.write_raw(&token.text),
                _ => writer.write_word(&token.text),
            },
            TokenKind::Operator => writer.operator(&token.text),
            TokenKind::Word | TokenKind::String => {
                if top_level_block_closed && writer.at_line_start() {
                    writer.blank_lines(config.line_breaks.empty_lines_between_probes);
                    top_level_block_closed = false;
                }
                writer.write_word(&token.text);
            }
        }
    }

    writer.finish()
}

struct Writer<'a> {
    config: &'a Config,
    out: String,
    indent: usize,
    line_has_content: bool,
    pending_space: bool,
}

impl<'a> Writer<'a> {
    fn new(config: &'a Config) -> Self {
        Self {
            config,
            out: String::new(),
            indent: 0,
            line_has_content: false,
            pending_space: false,
        }
    }

    fn finish(mut self) -> String {
        while self.out.ends_with([' ', '\t', '\n']) {
            self.out.pop();
        }
        self.out.push('\n');
        self.out
    }

    fn at_line_start(&self) -> bool {
        !self.line_has_content
    }

    fn write_indent(&mut self) {
        if !self.line_has_content {
            if self.config.indent.use_spaces {
                self.out
                    .push_str(&" ".repeat(self.indent * self.config.indent.size));
            } else {
                self.out.push_str(&"\t".repeat(self.indent));
            }
            self.line_has_content = true;
        }
    }

    fn maybe_space(&mut self) {
        if self.pending_space
            && self.line_has_content
            && !self
                .out
                .ends_with([' ', '\n', '\t', '(', '[', '{', ':', '.'])
        {
            self.out.push(' ');
        }
        self.pending_space = false;
    }

    fn write_raw(&mut self, text: &str) {
        self.write_indent();
        self.out.push_str(text);
    }

    fn write_word(&mut self, text: &str) {
        self.write_indent();
        self.maybe_space();
        if needs_space_before_word(&self.out) {
            self.out.push(' ');
        }
        self.out.push_str(text);
        self.pending_space = true;
    }

    fn write_line_token(&mut self, text: &str) {
        if self.line_has_content {
            self.newline();
        }
        self.write_indent();
        self.out.push_str(text.trim_end());
        self.newline();
    }

    fn write_comment(&mut self, text: &str) {
        if self.line_has_content && self.config.comments.preserve_inline {
            if !self.out.ends_with(' ') {
                self.out.push(' ');
            }
            self.out.push_str(text.trim_end());
            self.newline();
            return;
        }
        self.write_line_token(text);
    }

    fn operator(&mut self, op: &str) {
        if matches!(op, "->" | ".." | "/") || (op == "*" && self.out.ends_with(['_', ':', '.'])) {
            trim_trailing_space(&mut self.out);
            self.write_raw(op);
            self.pending_space = false;
            return;
        }

        if self.config.spacing.around_operators {
            self.pending_space = true;
            self.maybe_space();
        }
        self.write_raw(op);
        self.pending_space = self.config.spacing.around_operators;
    }

    fn comma(&mut self) {
        trim_trailing_space(&mut self.out);
        self.write_raw(",");
        self.pending_space = self.config.spacing.around_commas;
    }

    fn semicolon(&mut self) {
        trim_trailing_space(&mut self.out);
        self.write_raw(";");
        self.newline();
    }

    fn open_paren(&mut self) {
        trim_trailing_space(&mut self.out);
        self.write_raw("(");
        self.pending_space = self.config.spacing.around_parentheses;
    }

    fn close_paren(&mut self) {
        trim_trailing_space(&mut self.out);
        self.write_raw(")");
        self.pending_space = true;
    }

    fn open_bracket(&mut self) {
        trim_trailing_space(&mut self.out);
        self.write_raw("[");
        self.pending_space = self.config.spacing.around_brackets;
    }

    fn close_bracket(&mut self) {
        trim_trailing_space(&mut self.out);
        self.write_raw("]");
        self.pending_space = true;
    }

    fn open_block(&mut self) {
        trim_trailing_space(&mut self.out);
        match self.config.blocks.brace_style {
            BraceStyle::SameLine => {
                if self.config.spacing.before_block_start && self.line_has_content {
                    self.out.push(' ');
                }
            }
            BraceStyle::NextLine | BraceStyle::Gnu => {
                if self.line_has_content {
                    self.newline();
                }
                if matches!(self.config.blocks.brace_style, BraceStyle::Gnu) {
                    self.indent += 1;
                }
            }
        }
        self.write_raw("{");
        self.newline();
        if self.config.blocks.indent_statements {
            self.indent += 1;
        }
    }

    fn close_block(&mut self) {
        if self.config.blocks.indent_statements {
            self.indent = self.indent.saturating_sub(1);
        }
        if matches!(self.config.blocks.brace_style, BraceStyle::Gnu) {
            self.indent = self.indent.saturating_sub(1);
        }
        if self.line_has_content {
            self.newline();
        }
        self.write_indent();
        self.out.push('}');
        self.newline();
    }

    fn blank_lines(&mut self, count: usize) {
        if self.line_has_content {
            self.newline();
        }
        for _ in 0..count {
            if !self.out.ends_with("\n\n") {
                self.out.push('\n');
            }
        }
    }

    fn newline(&mut self) {
        trim_trailing_space(&mut self.out);
        if !self.out.ends_with('\n') {
            self.out.push('\n');
        }
        self.line_has_content = false;
        self.pending_space = false;
    }
}

fn needs_space_before_word(out: &str) -> bool {
    out.chars()
        .last()
        .is_some_and(|ch| ch.is_ascii_alphanumeric() || matches!(ch, '_' | '$' | '@' | ']' | ')'))
}

fn trim_trailing_space(out: &mut String) {
    while out.ends_with([' ', '\t']) {
        out.pop();
    }
}

fn tokenize(source: &str) -> Vec<Token> {
    let bytes = source.as_bytes();
    let mut tokens = Vec::new();
    let mut idx = 0;
    let mut at_line_start = true;

    while idx < bytes.len() {
        let byte = bytes[idx];
        if byte == b'\n' {
            tokens.push(Token::new(TokenKind::Newline, "\n"));
            idx += 1;
            at_line_start = true;
            continue;
        }
        if byte.is_ascii_whitespace() {
            idx += 1;
            continue;
        }
        if at_line_start && source[idx..].starts_with("#!") {
            let end = read_to_line_end(source, idx);
            tokens.push(Token::new(TokenKind::Shebang, &source[idx..end]));
            idx = end;
            at_line_start = false;
            continue;
        }
        if at_line_start && byte == b'#' {
            let end = read_to_line_end(source, idx);
            tokens.push(Token::new(TokenKind::Preprocessor, &source[idx..end]));
            idx = end;
            at_line_start = false;
            continue;
        }
        if source[idx..].starts_with("//") {
            let end = read_to_line_end(source, idx);
            tokens.push(Token::new(TokenKind::Comment, &source[idx..end]));
            idx = end;
            at_line_start = false;
            continue;
        }
        if matches!(byte, b'\'' | b'\"') {
            let end = read_string(source, idx, byte);
            tokens.push(Token::new(TokenKind::String, &source[idx..end]));
            idx = end;
            at_line_start = false;
            continue;
        }
        if let Some(op) = read_multi_operator(&source[idx..]) {
            tokens.push(Token::new(TokenKind::Operator, op));
            idx += op.len();
            at_line_start = false;
            continue;
        }
        if is_word_start(byte) {
            let end = read_word(source, idx);
            tokens.push(Token::new(TokenKind::Word, &source[idx..end]));
            idx = end;
            at_line_start = false;
            continue;
        }
        let ch = source[idx..].chars().next().unwrap();
        let kind = if is_operator_char(ch) {
            TokenKind::Operator
        } else {
            TokenKind::Symbol
        };
        tokens.push(Token::new(kind, &source[idx..idx + ch.len_utf8()]));
        idx += ch.len_utf8();
        at_line_start = false;
    }

    tokens
}

impl Token {
    fn new(kind: TokenKind, text: &str) -> Self {
        Self {
            kind,
            text: text.to_string(),
        }
    }
}

fn read_to_line_end(source: &str, start: usize) -> usize {
    source[start..]
        .find('\n')
        .map(|rel| start + rel)
        .unwrap_or(source.len())
}

fn read_string(source: &str, start: usize, quote: u8) -> usize {
    let bytes = source.as_bytes();
    let mut idx = start + 1;
    while idx < bytes.len() {
        if bytes[idx] == b'\\' {
            idx += 2;
            continue;
        }
        if bytes[idx] == quote {
            return idx + 1;
        }
        idx += 1;
    }
    source.len()
}

fn read_multi_operator(source: &str) -> Option<&'static str> {
    const OPS: &[&str] = &[
        "<<=", ">>=", "==", "!=", "<=", ">=", "&&", "||", "<<", ">>", "->", "+=", "-=", "*=", "/=",
        "%=", "&=", "|=", "^=", "++", "--", "..",
    ];
    OPS.iter().copied().find(|op| source.starts_with(op))
}

fn read_word(source: &str, start: usize) -> usize {
    let bytes = source.as_bytes();
    let mut idx = start;
    while idx < bytes.len() {
        let byte = bytes[idx];
        if byte.is_ascii_alphanumeric() || matches!(byte, b'_' | b'$' | b'@') {
            idx += 1;
        } else {
            break;
        }
    }
    idx
}

fn is_word_start(byte: u8) -> bool {
    byte.is_ascii_alphanumeric() || matches!(byte, b'_' | b'$' | b'@')
}

fn is_operator_char(ch: char) -> bool {
    matches!(
        ch,
        '=' | '+' | '-' | '*' | '/' | '%' | '&' | '|' | '^' | '<' | '>' | '!' | '~' | '?'
    )
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn formats_basic_probe() {
        let got = format_source("BEGIN{printf(\"x\",1);}", &Config::default()).unwrap();
        assert!(got.contains("BEGIN"));
        assert!(got.contains("printf"));
        assert!(got.ends_with('\n'));
    }
}
