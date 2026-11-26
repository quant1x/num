//! 与 arithmetics_add.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::at_or_nan;

/// add 逐元素加：按第一个操作数长度对齐，第二个短补默认值（NaN）。
pub fn add<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i] + at_or_nan(y, i)).collect()
}
