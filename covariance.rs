//! 与 covariance.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::mean::mean;

/// variance 总体方差（ddof=0）；空输入返回 0（与 Go 一致）。
pub fn variance<T: Float>(x: &[T]) -> T {
    if x.is_empty() {
        return T::zero();
    }
    let average = mean(x);
    let mut total = T::zero();
    for value in x {
        let delta = *value - average;
        total += delta * delta;
    }
    total / T::from_f64(x.len() as f64)
}

/// covariance 总体协方差（ddof=0）；长度不等或为空返回 0。
pub fn covariance(x: &[f64], y: &[f64]) -> f64 {
    if x.len() != y.len() || x.is_empty() {
        return 0.0;
    }
    let mean_x = mean(x);
    let mean_y = mean(y);
    let mut total = 0.0;
    for (left, right) in x.iter().zip(y.iter()) {
        total += (left - mean_x) * (right - mean_y);
    }
    total / x.len() as f64
}
