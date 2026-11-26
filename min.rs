//! 与 min.go 同语义的 Rust 实现。

use crate::float::Float;

/// min 最小值；空输入返回 NaN。
pub fn min<T: Float>(x: &[T]) -> T {
    let mut value = T::nan();
    let mut first = true;
    for item in x {
        if first {
            value = *item;
            first = false;
        } else if *item < value {
            value = *item;
        }
    }
    value
}
