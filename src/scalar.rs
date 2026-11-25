//! N1 标量与类型边界。

use crate::float::Float;

/// DefaultValue 是各类型的"无效值"：浮点 NaN、整型 0、bool false、字符串 "NaN"。
pub trait DefaultValue {
    fn default_value() -> Self;
}

impl DefaultValue for f64 {
    fn default_value() -> Self {
        f64::NAN
    }
}

impl DefaultValue for f32 {
    fn default_value() -> Self {
        f32::NAN
    }
}

impl DefaultValue for i64 {
    fn default_value() -> Self {
        0
    }
}

impl DefaultValue for bool {
    fn default_value() -> Self {
        false
    }
}

impl DefaultValue for String {
    fn default_value() -> Self {
        String::from("NaN")
    }
}

/// type_default 返回类型的无效值。
pub fn type_default<T: DefaultValue>() -> T {
    T::default_value()
}

/// is_nan 按 Go 口径：NaN 与 ±Inf 都算无效值。
pub fn is_nan<T: Float>(value: T) -> bool {
    value.is_nan() || value.is_infinite()
}

/// bool_to_f64：true = 1，false = 0。
pub fn bool_to_f64(value: bool) -> f64 {
    if value { 1.0 } else { 0.0 }
}

/// bool_to_i64：true = 1，false = 0。
pub fn bool_to_i64(value: bool) -> i64 {
    i64::from(value)
}

/// ints_to_f64 把整型切片转成 f64 切片（SliceToFloat64 的 i64 分支）。
pub fn ints_to_f64(values: &[i64]) -> Vec<f64> {
    values.iter().map(|value| *value as f64).collect()
}

/// bools_to_f64 把布尔切片转成 f64 切片（SliceToFloat64 的 bool 分支）。
pub fn bools_to_f64(values: &[bool]) -> Vec<f64> {
    values.iter().map(|value| bool_to_f64(*value)).collect()
}

/// floats_to_int32 按向零截断转 int32（SliceFloatToInt32）。
pub fn floats_to_int32(values: &[f64]) -> Vec<i32> {
    values.iter().map(|value| *value as i32).collect()
}

/// floats_to_int64 按向零截断转 int64（SliceFloatToInt64）。
pub fn floats_to_int64(values: &[f64]) -> Vec<i64> {
    values.iter().map(|value| *value as i64).collect()
}
