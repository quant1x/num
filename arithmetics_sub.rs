//! 与 arithmetics_sub.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::at_or_nan;

/// sub 逐元素减，长度口径同 add。
pub fn sub<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i] - at_or_nan(y, i)).collect()
}
