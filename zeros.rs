//! 与 zeros.go 同语义的 Rust 实现。

use crate::float::Float;

/// zeros 全零。
pub fn zeros<T: Float>(shape: usize) -> Vec<T> {
    vec![T::zero(); shape]
}
