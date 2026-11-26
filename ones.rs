//! 与 ones.go 同语义的 Rust 实现。

use crate::float::Float;

/// ones 全一。
pub fn ones<T: Float>(shape: usize) -> Vec<T> {
    vec![T::one(); shape]
}
