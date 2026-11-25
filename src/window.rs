//! N4 窗口与滚动族。

use crate::float::Float;

/// Window 是滚动窗口参数：常量或逐位置长度。
#[derive(Clone)]
pub enum Window {
    Constant(i64),
    PerPosition(Vec<i64>),
}

impl Window {
    fn at(&self, index: usize) -> i64 {
        match self {
            Window::Constant(value) => *value,
            Window::PerPosition(values) => values.get(index).copied().unwrap_or(0),
        }
    }
}

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

/// shift 索引平移：n > 0 向右、n < 0 向左，越界位填默认值。
pub fn shift<T: Float>(x: &[T], n: i64) -> Vec<T> {
    let length = x.len();
    let mut out = vec![T::nan(); length];
    if n == 0 {
        return x.to_vec();
    }
    let offset = n.unsigned_abs() as usize;
    if offset >= length {
        return out;
    }
    if n > 0 {
        out[offset..].copy_from_slice(&x[..length - offset]);
    } else {
        out[..length - offset].copy_from_slice(&x[offset..]);
    }
    out
}

/// diff 一阶离散差分：方向为 `s[i-n] - s[i]`（与 pandas 相反）。
pub fn diff<T: Float>(x: &[T], n: i64) -> Vec<T> {
    let shifted = shift(x, n);
    shifted
        .iter()
        .zip(x.iter())
        .map(|(left, right)| *left - *right)
        .collect()
}

/// cumprod 前缀积。
pub fn cumprod<T: Float>(x: &[T]) -> Vec<T> {
    let mut out = Vec::with_capacity(x.len());
    let mut total = T::one();
    for value in x {
        total *= *value;
        out.push(total);
    }
    out
}
