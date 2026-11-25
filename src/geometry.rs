//! 几何与多项式求值：斜率、角度换算、三角斜边、Horner 求值。

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
