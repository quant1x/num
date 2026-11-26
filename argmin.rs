//! 与 argmin.go 同语义的 Rust 实现。

use crate::float::Float;

/// argmin 最小值的首个下标；空输入返回 -1。
pub fn argmin<T: Float>(x: &[T]) -> i64 {
    let mut index = -1i64;
    let mut value = T::nan();
    for (i, item) in x.iter().enumerate() {
        if index < 0 || *item < value {
            value = *item;
            index = i as i64;
        }
    }
    index
}
