//! 与 argmax.go 同语义的 Rust 实现。

use crate::float::Float;

/// argmax 最大值的首个下标；空输入返回 -1。
pub fn argmax<T: Float>(x: &[T]) -> i64 {
    let mut index = -1i64;
    let mut value = T::nan();
    for (i, item) in x.iter().enumerate() {
        if index < 0 || *item > value {
            value = *item;
            index = i as i64;
        }
    }
    index
}
