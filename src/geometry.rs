//! 几何与多项式：斜率、角度换算、三角斜边、拟合与求值。

/// slope 两点斜率：`(y2 - y1) / (x2 - x1)`。
pub fn slope(x1: i64, y1: f64, x2: i64, y2: f64) -> f64 {
    (y2 - y1) / (x2 - x1) as f64
}

/// triangle_bevel 三角斜边：`slope * (xn - x1) + y1`。
pub fn triangle_bevel(slope: f64, x1: i64, y1: f64, xn: i64) -> f64 {
    slope * (xn - x1) as f64 + y1
}

/// degrees_to_slope 角度转斜率。
pub fn degrees_to_slope(degrees: f64) -> f64 {
    degrees.to_radians().tan()
}

/// slope_to_degrees 斜率转角度。
pub fn slope_to_degrees(m: f64) -> f64 {
    m.atan().to_degrees()
}

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
    let inverse_matrix = crate::matrix::inverse(&xtx);
    (0..n)
        .map(|j| (0..n).map(|k| inverse_matrix[j][k] * xty[k]).sum())
        .collect()
}
