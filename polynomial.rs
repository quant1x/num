//! 与 polynomial.go 同语义的 Rust 实现。

use crate::matrix::inverse;

/// poly_val 多项式求值：系数最高次在前，Horner 法。
pub fn poly_val(p: &[f64], x: &[f64]) -> Vec<f64> {
    let mut y = vec![0.0; x.len()];
    for coefficient in p {
        for (value, argument) in y.iter_mut().zip(x.iter()) {
            *value = *value * *argument + *coefficient;
        }
    }
    y
}

/// poly_fit 最小二乘拟合，正规方程 XTX⁻¹XTy；系数按幂次降序（p[0] 最高次）。
pub fn poly_fit(x: &[f64], y: &[f64], deg: i64) -> Vec<f64> {
    assert!(!x.is_empty(), "空输入");
    assert_eq!(x.len(), y.len(), "x/y 长度不一致");
    assert!(deg >= 0, "deg 不能为负");
    let n = deg as usize + 1;
    let mut vandermonde = vec![vec![0.0; n]; x.len()];
    for (i, value) in x.iter().enumerate() {
        for (j, slot) in vandermonde[i].iter_mut().enumerate() {
            let power = deg as usize - j;
            *slot = value.powi(power as i32);
        }
    }
    let mut xtx = vec![vec![0.0; n]; n];
    let mut xty = vec![0.0; n];
    for j in 0..n {
        for k in 0..n {
            xtx[j][k] = vandermonde.iter().map(|row| row[j] * row[k]).sum();
        }
        xty[j] = vandermonde
            .iter()
            .zip(y.iter())
            .map(|(row, target)| row[j] * target)
            .sum();
    }
    let inverse_matrix = inverse(&xtx);
    (0..n)
        .map(|j| (0..n).map(|k| inverse_matrix[j][k] * xty[k]).sum())
        .collect()
}
