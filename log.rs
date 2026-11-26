//! 与 log.go 同语义的 Rust 实现。

use crate::float::Float;

/// expm1 逐元素 `e^x - 1`。
pub fn expm1<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.exp_m1()).collect()
}

/// log1p 逐元素 `ln(1+x)`。
pub fn log1p<T: Float>(x: &[T]) -> Vec<T> {
    x.iter().map(|value| value.ln_1p()).collect()
}
