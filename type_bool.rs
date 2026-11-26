//! 与 type_bool.go 同语义的 Rust 实现。

/// bool_to_f64：true = 1，false = 0。
pub fn bool_to_f64(value: bool) -> f64 {
    if value { 1.0 } else { 0.0 }
}

/// bool_to_i64：true = 1，false = 0。
pub fn bool_to_i64(value: bool) -> i64 {
    i64::from(value)
}

/// bools_to_f64 把布尔切片转成 f64 切片（SliceToFloat64 的 bool 分支）。
pub fn bools_to_f64(values: &[bool]) -> Vec<f64> {
    values.iter().map(|value| bool_to_f64(*value)).collect()
}
