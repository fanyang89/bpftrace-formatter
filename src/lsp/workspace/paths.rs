use ignore::gitignore::Gitignore;
use std::fs;
use std::path::{Component, Path, PathBuf};
use tower_lsp::lsp_types::Url;

pub(super) fn uri_key(uri: &Url) -> Url {
    uri.to_file_path()
        .ok()
        .and_then(|path| Url::from_file_path(path).ok())
        .unwrap_or_else(|| uri.clone())
}

pub(super) fn canonicalize(path: &Path) -> Option<PathBuf> {
    let path = fs::canonicalize(path).ok()?;
    // Windows canonical paths use a verbatim prefix (\\?\), while LSP file
    // URIs decode to ordinary paths. Compare both in the URI representation.
    Url::from_file_path(path).ok()?.to_file_path().ok()
}

// Apply the walker's standard path filters to individual lifecycle updates, including
// file-backed buffers that do not exist on disk yet. This does not walk directories.
pub(super) fn is_eligible(path: &Path, roots: &[PathBuf]) -> bool {
    if path.extension().and_then(|extension| extension.to_str()) != Some("bt") {
        return false;
    }
    let Some(root) = roots
        .iter()
        .filter(|root| path.starts_with(root))
        .max_by_key(|root| root.components().count())
    else {
        return false;
    };
    let Some(parent) = path.parent() else {
        return false;
    };
    let ancestors: Vec<_> = parent.ancestors().collect();
    let repository = ancestors.iter().position(|dir| dir.join(".git").exists());
    let ignores: Vec<_> = ancestors
        .iter()
        .map(|dir| Gitignore::new(dir.join(".ignore")).0)
        .collect();
    let gitignores: Vec<_> = repository
        .map(|end| {
            ancestors[..=end]
                .iter()
                .map(|dir| Gitignore::new(dir.join(".gitignore")).0)
                .collect()
        })
        .unwrap_or_default();
    let exclude = repository.map(|end| {
        let mut builder = ignore::gitignore::GitignoreBuilder::new(ancestors[end]);
        builder.add(ancestors[end].join(".git/info/exclude"));
        builder.build().unwrap_or_else(|_| Gitignore::empty())
    });
    let global = repository.map(|_| {
        ignore::gitignore::GitignoreBuilder::new(root)
            .build_global()
            .0
    });
    let mut current = root.to_path_buf();
    for component in path.strip_prefix(root).unwrap().components() {
        let Component::Normal(name) = component else {
            return false;
        };
        current.push(name);
        if fs::symlink_metadata(&current).is_ok_and(|metadata| metadata.file_type().is_symlink()) {
            return false;
        }
        let is_dir = current != path;
        let matched = ignores
            .iter()
            .chain(gitignores.iter())
            .chain(exclude.iter())
            .chain(global.iter())
            .filter(|matcher| current.starts_with(matcher.path()))
            .map(|matcher| matcher.matched(&current, is_dir))
            .find(|matched| !matched.is_none());
        match matched {
            Some(matched) if matched.is_ignore() => return false,
            Some(_) => {}
            None if name.to_string_lossy().starts_with('.') => return false,
            None => {}
        }
    }
    true
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn canonical_roots_match_lsp_paths_for_saved_and_unsaved_files() {
        let directory = tempfile::tempdir().unwrap();
        let canonical = fs::canonicalize(directory.path()).unwrap();
        let root = canonicalize(directory.path()).unwrap();
        for name in ["saved.bt", "unsaved.bt"] {
            let path = canonical.join(name);
            if name == "saved.bt" {
                fs::write(&path, "BEGIN {}").unwrap();
            }
            let uri_path = Url::from_file_path(&path).unwrap().to_file_path().unwrap();
            assert!(is_eligible(&uri_path, std::slice::from_ref(&root)));
        }
    }
}
