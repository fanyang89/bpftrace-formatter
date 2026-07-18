use anyhow::{Context, Result};
use std::fs::{self, OpenOptions};
use std::io::Write;
use std::path::Path;
use tempfile::NamedTempFile;

pub(crate) fn write_new(path: &Path, contents: &[u8]) -> Result<()> {
    let mut file = OpenOptions::new()
        .write(true)
        .create_new(true)
        .open(path)
        .with_context(|| format!("creating {}", path.display()))?;
    let result = file
        .write_all(contents)
        .with_context(|| format!("writing {}", path.display()))
        .and_then(|_| {
            file.sync_all()
                .with_context(|| format!("syncing {}", path.display()))
        });
    if let Err(error) = result {
        drop(file);
        let _ = fs::remove_file(path);
        return Err(error);
    }
    Ok(())
}

pub(crate) fn write_atomic(path: &Path, contents: &[u8]) -> Result<()> {
    let target = if fs::symlink_metadata(path)?.file_type().is_symlink() {
        fs::canonicalize(path).with_context(|| format!("resolving {}", path.display()))?
    } else {
        path.to_path_buf()
    };
    let metadata = fs::metadata(&target)?;
    let parent = target
        .parent()
        .filter(|parent| !parent.as_os_str().is_empty())
        .unwrap_or_else(|| Path::new("."));
    let mut temp = NamedTempFile::new_in(parent)
        .with_context(|| format!("creating temporary file in {}", parent.display()))?;
    temp.as_file()
        .set_permissions(metadata.permissions())
        .with_context(|| format!("preserving permissions for {}", target.display()))?;
    temp.write_all(contents)
        .with_context(|| format!("writing temporary file for {}", target.display()))?;
    temp.as_file()
        .sync_all()
        .with_context(|| format!("syncing temporary file for {}", target.display()))?;
    temp.persist(&target)
        .map_err(|err| err.error)
        .with_context(|| format!("replacing {}", target.display()))?;
    Ok(())
}
