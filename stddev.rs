//! 与 stddev.go 同语义的 Rust 实现。

use crate::covariance::variance;
use crate::float::Float;

/// std 总体标准差；空输入返回类型默认值。
pub fn std<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::nan();
    }
    variance(x).sqrt()
}
