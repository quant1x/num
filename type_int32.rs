//! 与 type_int32.go 同语义的 Rust 实现。

/// floats_to_int32 按向零截断转 int32（SliceFloatToInt32）。
pub fn floats_to_int32(values: &[f64]) -> Vec<i32> {
    values.iter().map(|value| *value as i32).collect()
}
