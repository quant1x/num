//! 与 logic_all.go 同语义的 Rust 实现。

/// all 全真；空输入为 true。
pub fn all(x: &[bool]) -> bool {
    x.iter().all(|value| *value)
}
