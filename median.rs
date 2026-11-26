//! 与 median.go 同语义的 Rust 实现。

use crate::float::Float;

/// median 按入参顺序取中间元素、不排序：空输入 NaN、单元素 0、偶数位取中间两值平均。
pub fn median<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::nan();
    }
    if x.len() == 1 {
        return T::zero();
    }
    if x.len().is_multiple_of(2) {
        let i = x.len() / 2;
        return (x[i - 1] + x[i]) / T::from_f64(2.0);
    }
    x[x.len() / 2]
}
