//! 与 median.go 同语义的 Rust 实现。

use crate::float::Float;

/// median 先排序再取中位数：空输入 NaN、单元素返回其本身、偶数位取中间两值平均、
/// 任一 NaN 传播。
pub fn median<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::nan();
    }
    if x.iter().any(|value| value.is_nan()) {
        return T::nan();
    }
    let mut ordered = x.to_vec();
    ordered.sort_unstable_by(|left, right| left.partial_cmp(right).expect("NaN 已在上方拦截"));
    let n = ordered.len();
    if n.is_multiple_of(2) {
        let i = n / 2;
        return (ordered[i - 1] + ordered[i]) / T::from_f64(2.0);
    }
    ordered[n / 2]
}
