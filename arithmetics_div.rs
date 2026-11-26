//! 与 arithmetics_div.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::at_or_nan;

/// div 逐元素除，长度口径同 add。
pub fn div<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i] / at_or_nan(y, i)).collect()
}
