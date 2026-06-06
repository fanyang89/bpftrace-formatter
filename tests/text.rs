use btfmt::text::{
    all_identifier_occurrences, full_range, identifier_at_position, offset_for_position,
    position_for_offset,
};
use tower_lsp::lsp_types::Position;

#[test]
fn offset_and_position_round_trip_ascii_and_unicode() {
    let text = "αBEGIN\n  printf(\"x\");\n";
    let offset = text.find("printf").unwrap();
    let pos = position_for_offset(text, offset);
    assert_eq!(offset_for_position(text, pos), offset);
    assert_eq!(position_for_offset(text, text.len()).line, 2);
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
