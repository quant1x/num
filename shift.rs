//! 与 shift.go 同语义的 Rust 实现。

use crate::float::Float;

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
