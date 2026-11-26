//! 与 decimals.go 同语义的 Rust 实现。

/// decimal 半程远离零；NaN 按 0 参与，负 digits 视为 0，默认 2 位。
pub fn decimal(value: f64, digits: Option<i32>) -> f64 {
    let mut digits = digits.unwrap_or(2);
    if digits < 0 {
        digits = 0;
    }
    let value = if value.is_nan() { 0.0 } else { value };
    let half = if value.is_sign_negative() { -5.0 } else { 5.0 };
    let scale = 10f64.powi(digits);
    let next_scale = 10f64.powi(digits + 1);
    ((value * next_scale + half) / 10.0).trunc() / scale
}

/// change_rate 返回 current / base。
pub fn change_rate(base: f64, current: f64) -> f64 {
    current / base
}

/// net_change_rate 返回 (change_rate - 1) * 100。
pub fn net_change_rate(base: f64, current: f64) -> f64 {
    (change_rate(base, current) - 1.0) * 100.0
}
