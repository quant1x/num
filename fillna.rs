//! 与 fillna.go 同语义的 Rust 实现。

use crate::float::Float;

/// fill 把 NaN 与 ±Inf 都替换成给定值（与 Go 的无效值口径一致）。
pub fn fill<T: Float>(x: &[T], value: T) -> Vec<T> {
    x.iter()
        .map(|item| {
            if item.is_nan() || item.is_infinite() {
                value
            } else {
                *item
            }
        })
        .collect()
}

/// fill_na 与 fill 同语义（Go 的两个入口只在签名上不同）。
pub fn fill_na<T: Float>(x: &[T], value: T) -> Vec<T> {
    fill(x, value)
}
