//! 与 matrix.go 同语义的 Rust 实现。

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

/// inverse 2×2 及以上方阵求逆；与 Go 一致，不做主元搜索。
pub fn inverse(a: &[Vec<f64>]) -> Vec<Vec<f64>> {
    let n = a.len();
    assert!(n > 0, "空矩阵");
    assert!(a.iter().all(|row| row.len() == n), "不是方阵");
    let mut aug = vec![vec![0.0; 2 * n]; n];
    for i in 0..n {
        for j in 0..n {
            aug[i][j] = a[i][j];
        }
        aug[i][n + i] = 1.0;
    }
    for col in 0..n {
        let pivot = aug[col][col];
        for value in aug[col].iter_mut() {
            *value /= pivot;
        }
        let pivot_row = aug[col].clone();
        for (row, values) in aug.iter_mut().enumerate() {
            if row == col {
                continue;
            }
            let factor = values[col];
            if factor == 0.0 {
                continue;
            }
            for (dst, src) in values.iter_mut().zip(pivot_row.iter()) {
                *dst -= factor * *src;
            }
        }
    }
    aug.iter().map(|row| row[n..].to_vec()).collect()
}
