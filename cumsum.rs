//! 与 cumsum.go 同语义的 Rust 实现。

use crate::float::Float;

/// cumsum 前缀和。
pub fn cumsum<T: Float>(x: &[T]) -> Vec<T> {
    let mut out = Vec::with_capacity(x.len());
    let mut total = T::zero();
    for value in x {
        total += *value;
        out.push(total);
    }
    out
}
