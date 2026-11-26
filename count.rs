//! 与 count.go 同语义的 Rust 实现。

use crate::float::Float;

/// count 统计非零元素。
pub fn count<T: Float>(x: &[T]) -> usize {
    x.iter().filter(|value| **value != T::zero()).count()
}

/// count_true 统计布尔序列中 true 的个数。
pub fn count_true(x: &[bool]) -> usize {
    x.iter().filter(|value| **value).count()
}
