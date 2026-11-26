//! 与 where.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::{at_or_false, at_or_nan};

/// where_select 条件为真取 x，否则取 y；按最长对齐，缺位为默认值（NaN）。
pub fn where_select<T: Float>(condition: &[bool], x: &[T], y: &[T]) -> Vec<T> {
    let size = condition.len().max(x.len()).max(y.len());
    (0..size)
        .map(|i| {
            if at_or_false(condition, i) {
                at_or_nan(x, i)
            } else {
                at_or_nan(y, i)
            }
        })
        .collect()
}
