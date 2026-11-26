//! 与 align.go 同语义的 Rust 实现。

use crate::float::Float;

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
