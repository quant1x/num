//! N6 向量与容器。

use crate::float::Float;

/// zeros 全零。
pub fn zeros<T: Float>(shape: usize) -> Vec<T> {
    vec![T::zero(); shape]
}

/// ones 全一。
pub fn ones<T: Float>(shape: usize) -> Vec<T> {
    vec![T::one(); shape]
}

/// arange_f64 等间距序列，不含 end；step 非正时返回空。
pub fn arange_f64(start: f64, end: f64, step: f64) -> Vec<f64> {
    let mut out = Vec::new();
    if step <= 0.0 {
        return out;
    }
    let mut value = start;
    while value < end {
        out.push(value);
        value += step;
    }
    out
}

/// arange_i64 整数等间距序列。
pub fn arange_i64(start: i64, end: i64, step: i64) -> Vec<i64> {
    let mut out = Vec::new();
    if step <= 0 {
        return out;
    }
    let mut value = start;
    while value < end {
        out.push(value);
        value += step;
    }
    out
}

/// linspace 在 [start, end] 之间取 n 个等距点；endpoint=false 时不含 end。
pub fn linspace(start: f64, end: f64, n: usize, endpoint: bool) -> Vec<f64> {
    if n == 0 {
        return Vec::new();
    }
    if n == 1 {
        return vec![start];
    }
    let divisor = if endpoint { n - 1 } else { n } as f64;
    let step = (end - start) / divisor;
    (0..n).map(|i| start + i as f64 * step).collect()
}

/// range_i64 产生 0..n-1。
pub fn range_i64(n: usize) -> Vec<i64> {
    (0..n as i64).collect()
}

/// range_f64 产生 0..n-1。
pub fn range_f64(n: usize) -> Vec<f64> {
    (0..n).map(|i| i as f64).collect()
}

/// range_float 产生从 a 到 b（不含）步长 1 的序列。
pub fn range_float(a: f64, b: f64) -> Vec<f64> {
    arange_f64(a, b, 1.0)
}

/// repeat 同值重复 n 次。
pub fn repeat<T: Copy>(value: T, n: usize) -> Vec<T> {
    vec![value; n]
}

/// align 补默认值到 n：超长截断、不足补 value。
pub fn align<T: Float>(x: &[T], value: T, n: usize) -> Vec<T> {
    if x.len() == n {
        return x.to_vec();
    }
    if x.len() > n {
        return x[..n].to_vec();
    }
    let mut out = vec![value; n];
    out[..x.len()].copy_from_slice(x);
    out
}

/// shape 返回 (行, 列)；一维时列数为 0。
pub fn shape<T>(x: &[T]) -> (usize, usize) {
    (x.len(), 0)
}

/// shape_2d 返回二维矩阵的 (行, 列)。
pub fn shape_2d<T>(x: &[Vec<T>]) -> (usize, usize) {
    (x.len(), x.first().map_or(0, Vec::len))
}

/// concat1d 把三条等长序列逐元素拼成 n×3 矩阵（列式，不是按行）。
pub fn concat1d<T: Copy>(a: &[T], b: &[T], c: &[T]) -> Vec<Vec<T>> {
    (0..a.len()).map(|i| vec![a[i], b[i], c[i]]).collect()
}

/// transpose2d 二维转置。
pub fn transpose2d<T: Copy>(x: &[Vec<T>]) -> Vec<Vec<T>> {
    let rows = x.len();
    if rows == 0 {
        return Vec::new();
    }
    let cols = x[0].len();
    let mut out = vec![Vec::with_capacity(rows); cols];
    for row in x {
        for (j, value) in row.iter().take(cols).enumerate() {
            out[j].push(*value);
        }
    }
    out
}

/// fill 把 NaN 与 ±Inf 都替换成给定值（与 Go 的无效值口径一致）。
pub fn fill<T: Float>(x: &[T], value: T) -> Vec<T> {
    x.iter()
        .map(|item| {
            if item.is_nan() || item.is_infinite() {
                value
            } else {
                *item
            }
        })
        .collect()
}

/// fill_na 与 fill 同语义（Go 的两个入口只在签名上不同）。
pub fn fill_na<T: Float>(x: &[T], value: T) -> Vec<T> {
    fill(x, value)
}
