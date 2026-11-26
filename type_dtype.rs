//! 与 type_dtype.go 同语义的 Rust 实现。

use crate::float::Float;
use crate::vek::to_int32;

/// dtype2int：与 to_int32 同义。
pub fn dtype2int(x: &[f64]) -> Vec<i32> {
    to_int32(x)
}

/// is_nan 按 Go 口径：NaN 与 ±Inf 都算无效值。
pub fn is_nan<T: Float>(value: T) -> bool {
    value.is_nan() || value.is_infinite()
}
