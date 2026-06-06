use tower_lsp::lsp_types::{Position, Range};

pub fn line_starts(text: &str) -> Vec<usize> {
    let mut starts = vec![0];
    for (idx, byte) in text.bytes().enumerate() {
        if byte == b'\n' {
            starts.push(idx + 1);
        }
    }
    starts
}

pub fn position_for_offset(text: &str, offset: usize) -> Position {
    let starts = line_starts(text);
    let offset = offset.min(text.len());
    let line_idx = match starts.binary_search(&offset) {
        Ok(idx) => idx,
        Err(idx) => idx.saturating_sub(1),
    };
    let line_start = starts[line_idx];
    let character = text[line_start..offset]
        .chars()
        .map(|ch| ch.len_utf16() as u32)
        .sum();
    Position::new(line_idx as u32, character)
}

pub fn offset_for_position(text: &str, position: Position) -> usize {
    let starts = line_starts(text);
    let line = position.line as usize;
    if line >= starts.len() {
        return text.len();
    }

    let start = starts[line];
    let end = starts.get(line + 1).copied().unwrap_or(text.len());
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

pub fn range_for_offsets(text: &str, start: usize, end: usize) -> Range {
    Range::new(
        position_for_offset(text, start),
        position_for_offset(text, end),
    )
}

pub fn full_range(text: &str) -> Range {
    range_for_offsets(text, 0, text.len())
}

pub fn identifier_at_position(text: &str, position: Position) -> Option<(String, Range)> {
    let offset = offset_for_position(text, position);
    let bytes = text.as_bytes();
    if bytes.is_empty() {
        return None;
    }

    let mut start = offset.min(bytes.len());
    if start == bytes.len() && start > 0 {
        start -= 1;
    }
    while start > 0 && is_ident_byte(bytes[start - 1]) {
        start -= 1;
    }

    let mut end = offset.min(bytes.len());
    while end < bytes.len() && is_ident_byte(bytes[end]) {
        end += 1;
    }

    if start >= end {
        return None;
    }
    let value = text[start..end].to_string();
    Some((value, range_for_offsets(text, start, end)))
}

pub fn all_identifier_occurrences(text: &str, needle: &str) -> Vec<Range> {
    if needle.is_empty() {
        return Vec::new();
    }

    let mut ranges = Vec::new();
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
            ranges.push(range_for_offsets(text, start, end));
        }
        search_start = end;
    }
    ranges
}

fn is_ident_byte(byte: u8) -> bool {
    byte.is_ascii_alphanumeric() || matches!(byte, b'_' | b'$' | b'@')
}
