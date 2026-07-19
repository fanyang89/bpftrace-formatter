const CATALOG: &str = include_str!("tracepoint_catalog.txt");

pub(super) fn targets() -> impl Iterator<Item = &'static str> {
    CATALOG
        .lines()
        .filter(|line| !line.is_empty() && !line.starts_with('#'))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn catalog_is_sorted_unique_and_well_formed() {
        let targets = targets().collect::<Vec<_>>();
        assert!(targets.len() > 1_000);
        assert!(targets.windows(2).all(|pair| pair[0] < pair[1]));
        assert!(targets.iter().all(|target| {
            target
                .split_once(':')
                .is_some_and(|(system, event)| !system.is_empty() && !event.is_empty())
        }));
        assert!(targets.contains(&"sched:sched_switch"));
        assert!(targets.contains(&"syscalls:sys_enter_openat"));
    }
}
