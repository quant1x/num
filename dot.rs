//! 与 dot.go 同语义的 Rust 实现。

use crate::float::Float;

/// dot1d 一维点积，按第一个操作数长度对齐。
pub fn dot1d<T: Float>(x: &[T], y: &[T]) -> T {
    let mut total = T::zero();
    for (i, value) in x.iter().enumerate() {
        total += *value * y.get(i).copied().unwrap_or_else(T::nan);
    }
    total
}

/// dot2d 矩阵乘：`a` 的列数必须等于 `b` 的行数。
pub fn dot2d(a: &[Vec<f64>], b: &[Vec<f64>]) -> Vec<Vec<f64>> {
    assert!(!a.is_empty() && !b.is_empty(), "空矩阵");
    assert_eq!(a[0].len(), b.len(), "矩阵维度不匹配");
    let rows = a.len();
    let inner = b.len();
    let cols = b[0].len();
    let mut out = vec![vec![0.0; cols]; rows];
    for i in 0..rows {
        for k in 0..inner {
            let value = a[i][k];
            for j in 0..cols {
                out[i][j] += value * b[k][j];
            }
        }
    }
    out
}

/// dot2d1 矩阵乘向量。
pub fn dot2d1(a: &[Vec<f64>], v: &[f64]) -> Vec<f64> {
    assert!(!a.is_empty(), "空矩阵");
    assert_eq!(a[0].len(), v.len(), "矩阵与向量维度不匹配");
    a.iter()
        .map(|row| row.iter().zip(v.iter()).map(|(x, y)| x * y).sum())
        .collect()
}
