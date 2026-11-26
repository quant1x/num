//! 与 sum.go 同语义的 Rust 实现。

use crate::float::Float;

/// sum 求和。
pub fn sum<T: Float>(x: &[T]) -> T {
    let mut total = T::zero();
    for value in x {
        total += *value;
    }
    total
}
