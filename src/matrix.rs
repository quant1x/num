//! N6 线代：矩阵乘、矩阵乘向量、4×4 乘与求逆。

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

/// matmul 行主序 n×n 矩阵乘。
pub fn matmul(x: &[f64], y: &[f64], n: usize) -> Vec<f64> {
    assert_eq!(x.len(), n * n, "x 不是 n×n");
    assert_eq!(y.len(), n * n, "y 不是 n×n");
    let mut out = vec![0.0; n * n];
    for i in 0..n {
        for k in 0..n {
            let value = x[i * n + k];
            for j in 0..n {
                out[i * n + j] += value * y[k * n + j];
            }
        }
    }
    out
}

/// mat4mul 两个 4×4 行主序矩阵乘。
pub fn mat4mul(x: &[f64], y: &[f64]) -> Vec<f64> {
    assert_eq!(x.len(), 16, "x 不是 4×4");
    assert_eq!(y.len(), 16, "y 不是 4×4");
    matmul(x, y, 4)
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
