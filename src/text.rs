use tower_lsp::lsp_types::{Position, Range};

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct LineIndex {
    starts: Vec<usize>,
}

impl LineIndex {
    pub fn new(text: &str) -> Self {
        let mut starts = vec![0];
        for (idx, byte) in text.bytes().enumerate() {
            if byte == b'\n' {
                starts.push(idx + 1);
            }
        }
        Self { starts }
    }

    pub fn position_for_offset(&self, text: &str, offset: usize) -> Position {
        let offset = floor_char_boundary(text, offset.min(text.len()));
        let line_idx = match self.starts.binary_search(&offset) {
            Ok(idx) => idx,
            Err(idx) => idx.saturating_sub(1),
        };
        let line_start = self.starts[line_idx];
        let character = text[line_start..offset]
            .chars()
            .map(|ch| ch.len_utf16() as u32)
            .sum();
        Position::new(line_idx as u32, character)
    }

    pub fn offset_for_position(&self, text: &str, position: Position) -> usize {
        let line = position.line as usize;
        if line >= self.starts.len() {
            return text.len();
        }

        let start = self.starts[line];
        let end = self.starts.get(line + 1).copied().unwrap_or(text.len());
        let mut utf16_count = 0u32;
        for (rel, ch) in text[start..end].char_indices() {
            if utf16_count >= position.character {
                return start + rel;
            }
            utf16_count += ch.len_utf16() as u32;
            if utf16_count > position.character {
                return start + rel;
            }
        }
        end
    }

    pub fn range_for_offsets(&self, text: &str, start: usize, end: usize) -> Range {
        Range::new(
            self.position_for_offset(text, start),
            self.position_for_offset(text, end),
        )
    }

    pub fn full_range(&self, text: &str) -> Range {
        self.range_for_offsets(text, 0, text.len())
    }
}

pub fn line_starts(text: &str) -> Vec<usize> {
    LineIndex::new(text).starts
}

pub fn position_for_offset(text: &str, offset: usize) -> Position {
    LineIndex::new(text).position_for_offset(text, offset)
}

pub fn offset_for_position(text: &str, position: Position) -> usize {
    LineIndex::new(text).offset_for_position(text, position)
}

pub fn range_for_offsets(text: &str, start: usize, end: usize) -> Range {
    LineIndex::new(text).range_for_offsets(text, start, end)
}

pub fn full_range(text: &str) -> Range {
    LineIndex::new(text).full_range(text)
}

pub fn identifier_at_position(text: &str, position: Position) -> Option<(String, Range)> {
    let index = LineIndex::new(text);
    identifier_at_position_with_index(text, &index, position)
}

pub fn identifier_at_position_with_index(
    text: &str,
    index: &LineIndex,
    position: Position,
) -> Option<(String, Range)> {
    let offset = index.offset_for_position(text, position);
    let bytes = text.as_bytes();
    if bytes.is_empty() {
        return None;
    }

    let anchor = if bytes.get(offset).is_some_and(|byte| is_ident_byte(*byte)) {
        offset
    } else if offset > 0 && is_ident_byte(bytes[offset - 1]) {
        offset - 1
    } else {
        return None;
    };

    let mut start = anchor;
    while start > 0 && is_ident_byte(bytes[start - 1]) {
        start -= 1;
    }

    let mut end = anchor + 1;
    while end < bytes.len() && is_ident_byte(bytes[end]) {
        end += 1;
    }

    if start >= end {
        return None;
    }
    let value = text[start..end].to_string();
    Some((value, index.range_for_offsets(text, start, end)))
}

fn floor_char_boundary(text: &str, mut offset: usize) -> usize {
    while !text.is_char_boundary(offset) {
        offset -= 1;
    }
    offset
}

pub fn all_identifier_occurrences(text: &str, needle: &str) -> Vec<Range> {
    if needle.is_empty() {
        return Vec::new();
    }

    let mut ranges = Vec::new();
    let index = LineIndex::new(text);
    let mut search_start = 0;
    while let Some(rel) = text[search_start..].find(needle) {
        let start = search_start + rel;
        let end = start + needle.len();
        let before = start
            .checked_sub(1)
            .and_then(|idx| text.as_bytes().get(idx))
            .copied();
        let after = text.as_bytes().get(end).copied();
        if before.is_none_or(|b| !is_ident_byte(b)) && after.is_none_or(|b| !is_ident_byte(b)) {
            ranges.push(index.range_for_offsets(text, start, end));
        }
        search_start = end;
    }
    ranges
}

fn is_ident_byte(byte: u8) -> bool {
    byte.is_ascii_alphanumeric() || matches!(byte, b'_' | b'$' | b'@')
}
