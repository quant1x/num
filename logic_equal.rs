//! 与 logic_equal.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::at_or_nan;

/// equal 逐元素相等：按第一个操作数长度对齐，NaN 与自身比较为 false。
pub fn equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    (0..x.len()).map(|i| x[i] == at_or_nan(y, i)).collect()
}

/// not_equal 逐元素不等：NaN 与自身比较为 true。
pub fn not_equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    (0..x.len()).map(|i| x[i] != at_or_nan(y, i)).collect()
}
