//! 与 mean.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::sum::sum;

/// mean 均值；空输入返回 NaN。
pub fn mean<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::nan();
    }
    sum(x) / T::from_f64(x.len() as f64)
}
