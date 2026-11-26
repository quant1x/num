//! 与 maximum.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::at_or_nan;

/// maximum 按最长对齐取大；NaN 只在 `y > x` 成立时被覆盖，不传播。
pub fn maximum<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| {
            let left = at_or_nan(x, i);
            let right = at_or_nan(y, i);
            if right > left { right } else { left }
        })
        .collect()
}
