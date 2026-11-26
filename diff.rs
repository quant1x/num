//! 与 diff.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::shift::shift;

/// diff 一阶离散差分：方向为 `s[i-n] - s[i]`（与 pandas 相反）。
pub fn diff<T: Float>(x: &[T], n: i64) -> Vec<T> {
    let shifted = shift(x, n);
    shifted
        .iter()
        .zip(x.iter())
        .map(|(left, right)| *left - *right)
        .collect()
}
