//! 与 repeat.go 同语义的 Rust 实现。

/// repeat 同值重复 n 次。
pub fn repeat<T: Copy>(value: T, n: usize) -> Vec<T> {
    vec![value; n]
}
