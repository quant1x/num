//! 与 minimum.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::at_or_nan;

/// minimum 按最长对齐取小，NaN 口径同 maximum。
pub fn minimum<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| {
            let left = at_or_nan(x, i);
            let right = at_or_nan(y, i);
            if right < left { right } else { left }
        })
        .collect()
}
