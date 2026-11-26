//! 与 sqrt.go 同语义的 Rust 实现。

use crate::float::Float;

/// sqrt 逐元素平方根。
pub fn sqrt<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.sqrt()).collect()
}
