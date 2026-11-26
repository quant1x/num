//! 与 rolling.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::window::Window;

/// rolling 返回每个位置的窗口切片；窗口不足时为空的 Vec。
///
/// `min_periods` 与 Go 的 v1At 同口径：可用观测数达到阈值即用已有观测（partial window），
/// 大于窗口长度时不产生窗口。
pub fn rolling<T: Float>(x: &[T], window: &Window, min_periods: Option<i64>) -> Vec<Vec<T>> {
    let minimum = min_periods.unwrap_or(0);
    let mut blocks = Vec::with_capacity(x.len());
    for index in 0..x.len() {
        let size = window.at(index);
        let available = (index + 1) as i64;
        let length = if size < 0 {
            None
        } else if minimum <= 0 {
            if size > available { None } else { Some(size) }
        } else if size > available {
            if available < minimum {
                None
            } else {
                Some(available)
            }
        } else if size < minimum {
            None
        } else {
            Some(size)
        };
        match length {
            Some(length) => {
                let start = index + 1 - length as usize;
                blocks.push(x[start..=index].to_vec());
            }
            None => blocks.push(Vec::new()),
        }
    }
    blocks
}

/// rolling_apply 对每个完整窗口求值；窗口不足的位置填默认值（NaN）。
pub fn rolling_apply<T: Float, F>(x: &[T], window: &Window, apply: F) -> Vec<T>
where
    F: Fn(f64, &[T]) -> T,
{
    let mut out = Vec::with_capacity(x.len());
    for index in 0..x.len() {
        let size = window.at(index);
        if size <= 0 || size as usize > index + 1 {
            out.push(T::nan());
            continue;
        }
        let start = index + 1 - size as usize;
        out.push(apply(size as f64, &x[start..=index]));
    }
    out
}
