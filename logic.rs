//! 与 logic.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::support::{at_or_false, at_or_nan};

/// greater 逐元素 `>`，按最长对齐；NaN 参与比较恒为 false。
pub fn greater<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) > at_or_nan(y, i))
        .collect()
}

/// greater_equal 逐元素 `>=`。
pub fn greater_equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) >= at_or_nan(y, i))
        .collect()
}

/// less 逐元素 `<`。
pub fn less<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) < at_or_nan(y, i))
        .collect()
}

/// less_equal 逐元素 `<=`。
pub fn less_equal<T: Float>(x: &[T], y: &[T]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_nan(x, i) <= at_or_nan(y, i))
        .collect()
}

/// logic_and 逐元素逻辑与，按最长对齐，缺位按 false。
pub fn logic_and(x: &[bool], y: &[bool]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_false(x, i) && at_or_false(y, i))
        .collect()
}

/// logic_or 逐元素逻辑或，按最长对齐。
pub fn logic_or(x: &[bool], y: &[bool]) -> Vec<bool> {
    let size = x.len().max(y.len());
    (0..size)
        .map(|i| at_or_false(x, i) || at_or_false(y, i))
        .collect()
}

/// logic_not 逐元素逻辑非。
pub fn logic_not(x: &[bool]) -> Vec<bool> {
    x.iter().map(|value| !value).collect()
}
