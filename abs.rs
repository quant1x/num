//! 与 abs.go 同语义的 Rust 实现。

use crate::float::Float;

/// abs 逐元素绝对值。
pub fn abs<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.abs()).collect()
}
