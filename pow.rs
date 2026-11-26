//! 与 pow.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::at_or_nan;

/// pow_integer 整数次幂。
pub fn pow_integer<T: Float>(x: &[T], n: i32) -> Vec<T> {
    x.iter().map(|value| value.powi(n)).collect()
}

/// pow_float 逐元素幂，按第一个操作数长度对齐。
pub fn pow_float<T: Float>(x: &[T], y: &[T]) -> Vec<T> {
    (0..x.len()).map(|i| x[i].powf(at_or_nan(y, i))).collect()
}
