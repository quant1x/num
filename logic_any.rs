//! 与 logic_any.go 同语义的 Rust 实现。

/// any 任一为真；空输入为 false。
pub fn any(x: &[bool]) -> bool {
    x.iter().any(|value| *value)
}
