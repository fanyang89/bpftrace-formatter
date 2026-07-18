use std::collections::HashMap;
use tree_sitter::{Node, Tree};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub(super) enum AccessKind {
    Read,
    Write,
    Declaration,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub(super) struct ByteRange {
    pub start: usize,
    pub end: usize,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub(super) struct RenameResult {
    pub replacement: String,
    pub ranges: Vec<ByteRange>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
enum SymbolKind {
    Scratch,
    Map,
}

impl SymbolKind {
    fn sigil(self) -> char {
        match self {
            Self::Scratch => '$',
            Self::Map => '@',
        }
    }
}

#[derive(Debug)]
struct Scope {
    parent: Option<usize>,
    kind: ScopeKind,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum ScopeKind {
    Global,
    Lexical,
    Macro,
}

#[derive(Debug)]
struct Binding {
    kind: SymbolKind,
    name: String,
    scope: usize,
    occurrences: Vec<usize>,
    definition: Option<usize>,
    definition_is_declaration: bool,
}

#[derive(Debug)]
struct Occurrence {
    binding: usize,
    range: ByteRange,
    access: AccessKind,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
enum NamedKind {
    Macro,
    Parameter,
}

#[derive(Debug)]
struct NamedBinding {
    kind: NamedKind,
    name: String,
    scope: usize,
    occurrences: Vec<usize>,
    definitions: Vec<usize>,
}

#[derive(Debug)]
struct NamedOccurrence {
    binding: usize,
    range: ByteRange,
    access: AccessKind,
}

#[derive(Debug)]
pub(super) struct SymbolIndex {
    scopes: Vec<Scope>,
    bindings: Vec<Binding>,
    scope_bindings: Vec<HashMap<(SymbolKind, String), usize>>,
    occurrences: Vec<Occurrence>,
    named_bindings: Vec<NamedBinding>,
    named_scope_bindings: Vec<HashMap<(NamedKind, String), usize>>,
    named_occurrences: Vec<NamedOccurrence>,
}

impl SymbolIndex {
    pub(super) fn new(source: &str, tree: &Tree) -> Option<Self> {
        if tree.root_node().has_error() {
            return None;
        }
        let mut index = Self {
            scopes: vec![Scope {
                parent: None,
                kind: ScopeKind::Global,
            }],
            bindings: Vec::new(),
            scope_bindings: vec![HashMap::new()],
            occurrences: Vec::new(),
            named_bindings: Vec::new(),
            named_scope_bindings: vec![HashMap::new()],
            named_occurrences: Vec::new(),
        };
        index.predeclare_macros(source, tree.root_node());
        index.walk_node(source, tree.root_node(), 0);
        Some(index)
    }

    pub(super) fn prepare_range(&self, offset: usize) -> Option<ByteRange> {
        if let Some(occurrence) = self.occurrence_at(offset) {
            return self.bindings[occurrence.binding]
                .definition
                .is_some()
                .then_some(occurrence.range);
        }
        let occurrence = self.named_occurrence_at(offset)?;
        (!self.named_bindings[occurrence.binding]
            .definitions
            .is_empty())
        .then_some(occurrence.range)
    }

    #[cfg(test)]
    fn definition_at(&self, offset: usize) -> Option<ByteRange> {
        self.definitions_at(offset).into_iter().next()
    }

    pub(super) fn definitions_at(&self, offset: usize) -> Vec<ByteRange> {
        if let Some(occurrence) = self.occurrence_at(offset) {
            let binding = &self.bindings[occurrence.binding];
            return binding
                .definition
                .map(|idx| vec![self.occurrences[idx].range])
                .unwrap_or_default();
        }
        let Some(occurrence) = self.named_occurrence_at(offset) else {
            return Vec::new();
        };
        self.named_bindings[occurrence.binding]
            .definitions
            .iter()
            .map(|idx| self.named_occurrences[*idx].range)
            .collect()
    }

    pub(super) fn references_at(&self, offset: usize, include_declaration: bool) -> Vec<ByteRange> {
        if let Some(occurrence) = self.occurrence_at(offset) {
            let binding = &self.bindings[occurrence.binding];
            return binding
                .occurrences
                .iter()
                .filter(|idx| include_declaration || Some(**idx) != binding.definition)
                .map(|idx| self.occurrences[*idx].range)
                .collect();
        }
        let Some(occurrence) = self.named_occurrence_at(offset) else {
            return Vec::new();
        };
        let binding = &self.named_bindings[occurrence.binding];
        binding
            .occurrences
            .iter()
            .filter(|idx| include_declaration || !binding.definitions.contains(idx))
            .map(|idx| self.named_occurrences[*idx].range)
            .collect()
    }

    pub(super) fn highlights_at(&self, offset: usize) -> Vec<(ByteRange, AccessKind)> {
        if let Some(occurrence) = self.occurrence_at(offset) {
            return self.bindings[occurrence.binding]
                .occurrences
                .iter()
                .map(|idx| {
                    let occurrence = &self.occurrences[*idx];
                    (occurrence.range, occurrence.access)
                })
                .collect();
        }
        let Some(occurrence) = self.named_occurrence_at(offset) else {
            return Vec::new();
        };
        self.named_bindings[occurrence.binding]
            .occurrences
            .iter()
            .map(|idx| {
                let occurrence = &self.named_occurrences[*idx];
                (occurrence.range, occurrence.access)
            })
            .collect()
    }

    pub(super) fn rename_at(
        &self,
        offset: usize,
        new_name: &str,
    ) -> Result<Option<RenameResult>, String> {
        let Some(occurrence) = self.occurrence_at(offset) else {
            return self.rename_named_at(offset, new_name);
        };
        let binding = &self.bindings[occurrence.binding];
        let name = normalize_name(new_name, binding.kind)?;
        if self.bindings.iter().enumerate().any(|(idx, candidate)| {
            idx != occurrence.binding
                && candidate.kind == binding.kind
                && (candidate.scope == binding.scope
                    || (binding.kind == SymbolKind::Scratch
                        && self.scopes_related(candidate.scope, binding.scope)))
                && candidate.name == name
        }) {
            return Err(format!(
                "rename target {}{} already exists in this scope",
                binding.kind.sigil(),
                name
            ));
        }

        Ok(Some(RenameResult {
            replacement: format!("{}{}", binding.kind.sigil(), name),
            ranges: binding
                .occurrences
                .iter()
                .map(|idx| self.occurrences[*idx].range)
                .collect(),
        }))
    }

    fn occurrence_at(&self, offset: usize) -> Option<&Occurrence> {
        self.occurrences
            .iter()
            .filter(|occurrence| occurrence.range.start <= offset && offset <= occurrence.range.end)
            .min_by_key(|occurrence| occurrence.range.end - occurrence.range.start)
    }

    fn named_occurrence_at(&self, offset: usize) -> Option<&NamedOccurrence> {
        self.named_occurrences
            .iter()
            .filter(|occurrence| occurrence.range.start <= offset && offset <= occurrence.range.end)
            .min_by_key(|occurrence| occurrence.range.end - occurrence.range.start)
    }

    fn rename_named_at(
        &self,
        offset: usize,
        new_name: &str,
    ) -> Result<Option<RenameResult>, String> {
        let Some(occurrence) = self.named_occurrence_at(offset) else {
            return Ok(None);
        };
        let binding = &self.named_bindings[occurrence.binding];
        if binding.definitions.is_empty() {
            return Ok(None);
        }
        let name = new_name.trim();
        if !valid_identifier(name) {
            return Err(format!("invalid rename target {new_name:?}"));
        }
        if self
            .named_bindings
            .iter()
            .enumerate()
            .any(|(idx, candidate)| {
                idx != occurrence.binding
                    && candidate.kind == binding.kind
                    && candidate.scope == binding.scope
                    && candidate.name == name
            })
        {
            return Err(format!("rename target {name} already exists in this scope"));
        }
        Ok(Some(RenameResult {
            replacement: name.to_string(),
            ranges: binding
                .occurrences
                .iter()
                .map(|idx| self.named_occurrences[*idx].range)
                .collect(),
        }))
    }

    fn walk_node(&mut self, source: &str, node: Node<'_>, scope: usize) {
        match node.kind() {
            "macro_definition" => {
                let macro_scope = self.new_macro_scope(scope);
                if let Some(parameters) = node.child_by_field_name("parameters") {
                    let mut cursor = parameters.walk();
                    for child in parameters.named_children(&mut cursor) {
                        if is_variable(child) {
                            self.record_variable(
                                source,
                                child,
                                macro_scope,
                                Some(AccessKind::Declaration),
                            );
                        } else if child.kind() == "identifier" {
                            self.record_named(
                                source,
                                child,
                                macro_scope,
                                NamedKind::Parameter,
                                AccessKind::Declaration,
                            );
                        }
                    }
                }
                if let Some(body) = node.child_by_field_name("body") {
                    self.walk_children(source, body, macro_scope);
                }
                return;
            }
            "action" => {
                let action_scope = self.new_scope(scope);
                self.walk_children(source, node, action_scope);
                return;
            }
            "block" | "block_expression" => {
                let block_scope = self.new_scope(scope);
                self.walk_children(source, node, block_scope);
                return;
            }
            "for_statement" => {
                self.walk_for_statement(source, node, scope);
                return;
            }
            "scratch_variable" | "map_variable" => {
                self.record_variable(source, node, scope, None);
            }
            "identifier" => self.record_identifier(source, node, scope),
            _ => {}
        }
        self.walk_children(source, node, scope);
    }

    fn walk_children(&mut self, source: &str, node: Node<'_>, scope: usize) {
        let mut cursor = node.walk();
        for child in node.named_children(&mut cursor) {
            self.walk_node(source, child, scope);
        }
    }

    fn walk_for_statement(&mut self, source: &str, node: Node<'_>, scope: usize) {
        let Some(body) = node.child_by_field_name("body") else {
            self.walk_children(source, node, scope);
            return;
        };
        let body_scope = self.new_scope(scope);
        let mut cursor = node.walk();
        for child in node.named_children(&mut cursor) {
            if same_node(child, body) {
                continue;
            }
            if child.kind() == "scratch_variable" {
                self.record_variable(source, child, body_scope, Some(AccessKind::Declaration));
            } else {
                self.walk_node(source, child, scope);
            }
        }
        self.walk_children(source, body, body_scope);
    }

    fn predeclare_macros(&mut self, source: &str, root: Node<'_>) {
        let mut cursor = root.walk();
        for node in root.named_children(&mut cursor) {
            if node.kind() != "macro_definition" {
                continue;
            }
            let Some(name) = node.child_by_field_name("name") else {
                continue;
            };
            self.record_named(source, name, 0, NamedKind::Macro, AccessKind::Declaration);
        }
    }

    fn record_identifier(&mut self, source: &str, node: Node<'_>, scope: usize) {
        if node
            .parent()
            .is_some_and(|parent| parent.kind() == "macro_definition")
        {
            return;
        }
        let Some(name) = source.get(node.byte_range()) else {
            return;
        };
        if let Some(binding) = self.resolve_named_binding(scope, NamedKind::Parameter, name) {
            self.push_named_occurrence(binding, node, AccessKind::Read);
            return;
        }
        let is_call = node.parent().is_some_and(|parent| {
            parent.kind() == "call_expression" && field_matches(parent, "function", node)
        });
        let is_expression = node.parent().is_some_and(|parent| {
            !matches!(
                parent.kind(),
                "config_assignment" | "probe" | "struct_type" | "field_expression"
            )
        });
        if (is_call || is_expression) && in_executable_context(node) {
            if let Some(binding) = self.named_scope_bindings[0]
                .get(&(NamedKind::Macro, name.to_string()))
                .copied()
            {
                self.push_named_occurrence(binding, node, AccessKind::Read);
            }
        }
    }

    fn record_named(
        &mut self,
        source: &str,
        node: Node<'_>,
        scope: usize,
        kind: NamedKind,
        access: AccessKind,
    ) {
        let Some(name) = source.get(node.byte_range()) else {
            return;
        };
        if !valid_identifier(name) {
            return;
        }
        let binding = self.named_binding_in_scope(scope, kind, name);
        self.push_named_occurrence(binding, node, access);
    }

    fn push_named_occurrence(&mut self, binding: usize, node: Node<'_>, access: AccessKind) {
        let occurrence = self.named_occurrences.len();
        self.named_occurrences.push(NamedOccurrence {
            binding,
            range: ByteRange {
                start: node.start_byte(),
                end: node.end_byte(),
            },
            access,
        });
        let binding = &mut self.named_bindings[binding];
        binding.occurrences.push(occurrence);
        if access == AccessKind::Declaration {
            binding.definitions.push(occurrence);
        }
    }

    fn record_variable(
        &mut self,
        source: &str,
        node: Node<'_>,
        scope: usize,
        forced_access: Option<AccessKind>,
    ) {
        let Some((kind, name, range)) = variable_identity(source, node) else {
            return;
        };
        let (access, plain_assignment) = forced_access
            .map(|access| (access, false))
            .unwrap_or_else(|| classify_access(source, node, kind));
        let binding = match kind {
            SymbolKind::Scratch => {
                if access == AccessKind::Declaration {
                    self.binding_in_scope(scope, kind, &name)
                } else {
                    self.resolve_binding(scope, kind, &name)
                        .unwrap_or_else(|| self.binding_in_scope(scope, kind, &name))
                }
            }
            SymbolKind::Map => {
                if access == AccessKind::Declaration && scope != 0 {
                    self.binding_in_scope(scope, kind, &name)
                } else {
                    self.resolve_non_global_binding(scope, kind, &name)
                        .unwrap_or_else(|| {
                            self.macro_scope(scope)
                                .map(|scope| self.binding_in_scope(scope, kind, &name))
                                .unwrap_or_else(|| self.binding_in_scope(0, kind, &name))
                        })
                }
            }
        };
        let occurrence_idx = self.occurrences.len();
        self.occurrences.push(Occurrence {
            binding,
            range,
            access,
        });
        let binding = &mut self.bindings[binding];
        binding.occurrences.push(occurrence_idx);

        let is_definition = access == AccessKind::Declaration
            || (access == AccessKind::Write && (kind == SymbolKind::Map || plain_assignment));
        if access == AccessKind::Declaration {
            if !binding.definition_is_declaration {
                binding.definition = Some(occurrence_idx);
                binding.definition_is_declaration = true;
            }
        } else if is_definition && binding.definition.is_none() {
            binding.definition = Some(occurrence_idx);
        }
    }

    fn new_scope(&mut self, parent: usize) -> usize {
        let idx = self.scopes.len();
        self.scopes.push(Scope {
            parent: Some(parent),
            kind: ScopeKind::Lexical,
        });
        self.scope_bindings.push(HashMap::new());
        self.named_scope_bindings.push(HashMap::new());
        idx
    }

    fn new_macro_scope(&mut self, parent: usize) -> usize {
        let idx = self.scopes.len();
        self.scopes.push(Scope {
            parent: Some(parent),
            kind: ScopeKind::Macro,
        });
        self.scope_bindings.push(HashMap::new());
        self.named_scope_bindings.push(HashMap::new());
        idx
    }

    fn binding_in_scope(&mut self, scope: usize, kind: SymbolKind, name: &str) -> usize {
        let key = (kind, name.to_string());
        if let Some(binding) = self.scope_bindings[scope].get(&key) {
            return *binding;
        }
        let binding = self.bindings.len();
        self.bindings.push(Binding {
            kind,
            name: name.to_string(),
            scope,
            occurrences: Vec::new(),
            definition: None,
            definition_is_declaration: false,
        });
        self.scope_bindings[scope].insert(key, binding);
        binding
    }

    fn named_binding_in_scope(&mut self, scope: usize, kind: NamedKind, name: &str) -> usize {
        let key = (kind, name.to_string());
        if let Some(binding) = self.named_scope_bindings[scope].get(&key) {
            return *binding;
        }
        let binding = self.named_bindings.len();
        self.named_bindings.push(NamedBinding {
            kind,
            name: name.to_string(),
            scope,
            occurrences: Vec::new(),
            definitions: Vec::new(),
        });
        self.named_scope_bindings[scope].insert(key, binding);
        binding
    }

    fn resolve_named_binding(
        &self,
        mut scope: usize,
        kind: NamedKind,
        name: &str,
    ) -> Option<usize> {
        loop {
            if let Some(binding) = self.named_scope_bindings[scope].get(&(kind, name.to_string())) {
                return Some(*binding);
            }
            scope = self.scopes[scope].parent?;
        }
    }

    fn resolve_binding(&self, mut scope: usize, kind: SymbolKind, name: &str) -> Option<usize> {
        loop {
            if let Some(binding) = self.scope_bindings[scope].get(&(kind, name.to_string())) {
                return Some(*binding);
            }
            scope = self.scopes[scope].parent?;
        }
    }

    fn resolve_non_global_binding(
        &self,
        mut scope: usize,
        kind: SymbolKind,
        name: &str,
    ) -> Option<usize> {
        while scope != 0 {
            if let Some(binding) = self.scope_bindings[scope].get(&(kind, name.to_string())) {
                return Some(*binding);
            }
            scope = self.scopes[scope].parent?;
        }
        None
    }

    fn macro_scope(&self, mut scope: usize) -> Option<usize> {
        loop {
            if self.scopes[scope].kind == ScopeKind::Macro {
                return Some(scope);
            }
            scope = self.scopes[scope].parent?;
        }
    }

    fn scopes_related(&self, left: usize, right: usize) -> bool {
        self.is_ancestor(left, right) || self.is_ancestor(right, left)
    }

    fn is_ancestor(&self, ancestor: usize, mut scope: usize) -> bool {
        loop {
            if ancestor == scope {
                return true;
            }
            let Some(parent) = self.scopes[scope].parent else {
                return false;
            };
            scope = parent;
        }
    }
}

fn variable_identity(source: &str, node: Node<'_>) -> Option<(SymbolKind, String, ByteRange)> {
    let kind = match node.kind() {
        "scratch_variable" => SymbolKind::Scratch,
        "map_variable" => SymbolKind::Map,
        _ => return None,
    };
    let start = node.start_byte();
    let mut end = node.end_byte();
    if kind == SymbolKind::Map {
        let text = source.get(start..end)?;
        end = start + text.find('[').unwrap_or(text.len());
    }
    let text = source.get(start..end)?;
    let mut chars = text.chars();
    if chars.next()? != kind.sigil() {
        return None;
    }
    let name = chars.as_str();
    if name.is_empty() && kind == SymbolKind::Map {
        return Some((kind, String::new(), ByteRange { start, end }));
    }
    if !valid_identifier(name) {
        return None;
    }
    Some((kind, name.to_string(), ByteRange { start, end }))
}

fn classify_access(source: &str, node: Node<'_>, kind: SymbolKind) -> (AccessKind, bool) {
    let Some(parent) = node.parent() else {
        return (AccessKind::Read, false);
    };
    if parent.kind() == "declaration_statement" && field_matches(parent, "name", node)
        || parent.kind() == "map_declaration"
    {
        return (AccessKind::Declaration, false);
    }
    if parent.kind() == "assignment_statement" && field_matches(parent, "left", node) {
        let plain_assignment = parent
            .child_by_field_name("operator")
            .and_then(|operator| source.get(operator.byte_range()))
            == Some("=");
        return (AccessKind::Write, plain_assignment);
    }
    if parent.kind() == "update_expression" {
        return (AccessKind::Write, false);
    }
    if kind == SymbolKind::Map && is_mutating_map_call(source, node) {
        return (AccessKind::Write, false);
    }
    (AccessKind::Read, false)
}

fn is_mutating_map_call(source: &str, node: Node<'_>) -> bool {
    let mut current = node;
    let arguments = loop {
        let Some(parent) = current.parent() else {
            return false;
        };
        match parent.kind() {
            "arguments" => break parent,
            "parenthesized_expression" => current = parent,
            _ => return false,
        }
    };
    let Some(call) = arguments
        .parent()
        .filter(|parent| parent.kind() == "call_expression")
    else {
        return false;
    };
    call.child_by_field_name("function")
        .and_then(|function| source.get(function.byte_range()))
        .is_some_and(|name| matches!(name, "clear" | "delete" | "zero"))
}

fn field_matches(parent: Node<'_>, field: &str, child: Node<'_>) -> bool {
    parent
        .child_by_field_name(field)
        .is_some_and(|candidate| same_node(candidate, child))
}

fn same_node(left: Node<'_>, right: Node<'_>) -> bool {
    left.kind() == right.kind()
        && left.start_byte() == right.start_byte()
        && left.end_byte() == right.end_byte()
}

fn in_executable_context(mut node: Node<'_>) -> bool {
    while let Some(parent) = node.parent() {
        if matches!(
            parent.kind(),
            "action" | "block" | "block_expression" | "macro_definition"
        ) {
            return true;
        }
        node = parent;
    }
    false
}

fn is_variable(node: Node<'_>) -> bool {
    matches!(node.kind(), "scratch_variable" | "map_variable")
}

fn normalize_name(new_name: &str, kind: SymbolKind) -> Result<String, String> {
    let mut name = new_name.trim();
    if let Some(sigil) = name.chars().next().filter(|ch| matches!(ch, '$' | '@')) {
        if sigil != kind.sigil() {
            return Err(format!("rename target must use the {} sigil", kind.sigil()));
        }
        name = &name[sigil.len_utf8()..];
    }
    if name.is_empty() && kind == SymbolKind::Map && new_name.trim() == "@" {
        return Ok(String::new());
    }
    if !valid_identifier(name) {
        return Err(format!("invalid rename target {new_name:?}"));
    }
    Ok(name.to_string())
}

fn valid_identifier(name: &str) -> bool {
    let mut bytes = name.bytes();
    bytes
        .next()
        .is_some_and(|byte| byte.is_ascii_alphabetic() || byte == b'_')
        && bytes.all(|byte| byte.is_ascii_alphanumeric() || byte == b'_')
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::parse;

    fn index(source: &str) -> SymbolIndex {
        let tree = parse::ensure_valid(source).unwrap();
        SymbolIndex::new(source, &tree).unwrap()
    }

    #[test]
    fn scratch_variables_follow_action_and_block_scopes() {
        let source = concat!(
            "BEGIN { $x = 1; if (1) { print($x); let $x = 2; print($x); } print($x); }\n",
            "END { $x = 3; print($x); }\n",
        );
        let index = index(source);
        let first_outer_read = source.find("print($x)").unwrap() + "print(".len();
        let inner_read = source[first_outer_read + 1..].find("print($x)").unwrap()
            + first_outer_read
            + 1
            + "print(".len();
        let final_read = source.rfind("print($x)").unwrap() + "print(".len();

        assert_eq!(index.references_at(first_outer_read, true).len(), 3);
        assert_eq!(index.references_at(inner_read, true).len(), 2);
        assert_eq!(index.references_at(final_read, true).len(), 2);
        assert_ne!(
            index.definition_at(first_outer_read),
            index.definition_at(final_read)
        );
    }

    #[test]
    fn maps_are_global_except_macro_parameters() {
        let source = concat!(
            "macro use_map(@m) { @m = count(); print(@m); }\n",
            "BEGIN { @m = count(); print(@m); }\n",
            "END { print(@m); }\n",
        );
        let index = index(source);
        let macro_read = source.find("print(@m)").unwrap() + "print(".len();
        let global_read = source.rfind("print(@m)").unwrap() + "print(".len();

        assert_eq!(index.references_at(macro_read, true).len(), 3);
        assert_eq!(index.references_at(global_read, true).len(), 3);
        assert_ne!(
            index.definition_at(macro_read),
            index.definition_at(global_read)
        );
    }

    #[test]
    fn for_variables_are_declared_in_the_loop_body() {
        let source = "BEGIN { for ($i : 0..2) { print($i); } }\n";
        let index = index(source);
        let read = source.find("print($i)").unwrap() + "print(".len();

        assert_eq!(index.references_at(read, true).len(), 2);
        assert_eq!(
            index.definition_at(read),
            Some(ByteRange {
                start: source.find("$i").unwrap(),
                end: source.find("$i").unwrap() + 2,
            })
        );
    }

    #[test]
    fn map_indexes_and_unicode_keep_exact_symbol_ranges() {
        let source = "BEGIN { printf(\"你好\"); @m[\"[\"] = 1; print(@m[\"[\"]); }\n";
        let index = index(source);
        let read = source.find("print(@m").unwrap() + "print(".len();

        let references = index.references_at(read, true);
        assert_eq!(references.len(), 2);
        assert!(references
            .iter()
            .all(|range| &source[range.start..range.end] == "@m"));
        assert_eq!(index.references_at(read, false).len(), 1);
        assert_eq!(
            index
                .highlights_at(read)
                .into_iter()
                .map(|(_, access)| access)
                .collect::<Vec<_>>(),
            vec![AccessKind::Write, AccessKind::Read]
        );
    }

    #[test]
    fn block_expressions_have_lexical_scope() {
        let source = "BEGIN { $value = { let $x = 1; $x }; $x = 2; print($x); }\n";
        let index = index(source);
        let inner_read = source.find("$x }").unwrap();
        let outer_read = source.find("print($x)").unwrap() + "print(".len();

        assert_eq!(index.references_at(inner_read, true).len(), 2);
        assert_eq!(index.references_at(outer_read, true).len(), 2);
        assert_ne!(
            index.definition_at(inner_read),
            index.definition_at(outer_read)
        );
    }

    #[test]
    fn rename_rejects_capture_across_related_scopes() {
        let source = "BEGIN { $x = 1; if (1) { let $y = 2; print($x); } }\n";
        let index = index(source);
        let outer = source.find("$x").unwrap();

        assert!(index.rename_at(outer, "y").is_err());
    }

    #[test]
    fn undeclared_macro_maps_do_not_resolve_to_globals() {
        let source = concat!("macro bump() { @m++; }\n", "BEGIN { @m = 1; print(@m); }\n",);
        let index = index(source);
        let macro_map = source.find("@m").unwrap();
        let global_map = source.rfind("@m").unwrap();

        assert_eq!(index.references_at(macro_map, true).len(), 1);
        assert_eq!(index.references_at(global_map, true).len(), 2);
        assert_ne!(
            index.definition_at(macro_map),
            index.definition_at(global_map)
        );
    }

    #[test]
    fn anonymous_maps_can_be_navigated_and_renamed() {
        let source = "BEGIN { @ = count(); print(@); }\n";
        let anonymous_index = index(source);
        let read = source.rfind('@').unwrap();

        assert_eq!(anonymous_index.references_at(read, true).len(), 2);
        let rename = anonymous_index.rename_at(read, "counts").unwrap().unwrap();
        assert_eq!(rename.replacement, "@counts");
        assert_eq!(rename.ranges.len(), 2);

        let named_source = "BEGIN { @counts = count(); print(@counts); }\n";
        let named_index = index(named_source);
        let named_read = named_source.rfind("@counts").unwrap();
        let rename = named_index.rename_at(named_read, "@").unwrap().unwrap();
        assert_eq!(rename.replacement, "@");
        assert_eq!(rename.ranges.len(), 2);
    }

    #[test]
    fn mutating_map_calls_are_writes() {
        let source = "BEGIN { @m = 1; clear((@m)); zero((@m)); delete((@m)); }\n";
        let index = index(source);
        let read = source.rfind("@m").unwrap();
        let accesses: Vec<_> = index
            .highlights_at(read)
            .into_iter()
            .map(|(_, access)| access)
            .collect();

        assert_eq!(accesses, vec![AccessKind::Write; 4]);
    }

    #[test]
    fn macro_families_and_expression_parameters_are_indexed() {
        let source = concat!(
            "macro choose(value) { value }\n",
            "macro choose($value) { $value }\n",
            "BEGIN { $x = 1; print(choose(1)); print(choose($x)); }\n",
        );
        let index = index(source);
        let first_call = source.find("choose(1)").unwrap();
        let parameter_read = source.find("value }").unwrap();

        assert_eq!(index.definitions_at(first_call).len(), 2);
        assert_eq!(index.references_at(first_call, true).len(), 4);
        assert_eq!(index.references_at(parameter_read, true).len(), 2);

        let rename = index.rename_at(first_call, "select").unwrap().unwrap();
        assert_eq!(rename.replacement, "select");
        assert_eq!(rename.ranges.len(), 4);

        let builtin = source.find("print(").unwrap();
        assert!(index.prepare_range(builtin).is_none());
        assert!(index.rename_at(builtin, "display").unwrap().is_none());
    }

    #[test]
    fn rename_validates_sigils_collisions_and_ranges() {
        let source = "BEGIN { $x = 1; $other = 2; print($x); }\n";
        let index = index(source);
        let read = source.find("print($x)").unwrap() + "print(".len();

        let rename = index.rename_at(read, "renamed").unwrap().unwrap();
        assert_eq!(rename.replacement, "$renamed");
        assert_eq!(rename.ranges.len(), 2);
        assert!(index.rename_at(read, "@wrong").is_err());
        assert!(index.rename_at(read, "1bad").is_err());
        assert!(index.rename_at(read, "$other").is_err());
    }

    #[test]
    fn erroneous_trees_do_not_produce_an_index() {
        let source = "BEGIN { $x =";
        let parsed = parse::parse(source).unwrap();
        assert!(SymbolIndex::new(source, &parsed.tree).is_none());
    }
}
