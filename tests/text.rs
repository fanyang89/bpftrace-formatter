use btfmt::text::{
    all_identifier_occurrences, full_range, identifier_at_position, offset_for_position,
    position_for_offset, LineIndex,
};
use tower_lsp::lsp_types::Position;

#[test]
fn offset_and_position_round_trip_ascii_and_unicode() {
    let text = "α😀BEGIN\n  printf(\"x\");\n";
    let offset = text.find("printf").unwrap();
    let pos = position_for_offset(text, offset);
    assert_eq!(offset_for_position(text, pos), offset);
    assert_eq!(position_for_offset(text, text.len()).line, 2);
    assert_eq!(position_for_offset(text, 1), Position::new(0, 0));

    let index = LineIndex::new(text);
    assert_eq!(index.position_for_offset(text, offset), pos);
    assert_eq!(index.offset_for_position(text, pos), offset);
    assert_eq!(index.full_range(text), full_range(text));
}

#[test]
fn identifier_lookup_supports_variables_and_maps() {
    let text = "BEGIN { $x = 1; @m[$x] = count(); }";
    let pos = position_for_offset(text, text.find("$x").unwrap() + 1);
    let (name, range) = identifier_at_position(text, pos).unwrap();
    assert_eq!(name, "$x");
    assert_eq!(
        &text[offset_for_position(text, range.start)..offset_for_position(text, range.end)],
        "$x"
    );
}

#[test]
fn identifier_lookup_handles_end_of_file_and_unicode() {
    let text = "$value";
    let (name, _) = identifier_at_position(text, position_for_offset(text, text.len())).unwrap();
    assert_eq!(name, "$value");

    let text = "BEGIN { printf(\"😀\"); } 😀";
    assert!(identifier_at_position(text, position_for_offset(text, text.len())).is_none());
}

#[test]
fn identifier_occurrences_match_whole_tokens() {
    let text = "$x $xx @x $x";
    let ranges = all_identifier_occurrences(text, "$x");
    assert_eq!(ranges.len(), 2);
}

#[test]
fn full_range_spans_entire_document() {
    let text = "BEGIN\n{}";
    let range = full_range(text);
    assert_eq!(range.start, Position::new(0, 0));
    assert_eq!(offset_for_position(text, range.end), text.len());
}
