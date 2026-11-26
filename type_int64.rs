//! 与 type_int64.go 同语义的 Rust 实现。

/// floats_to_int64 按向零截断转 int64（SliceFloatToInt64）。
pub fn floats_to_int64(values: &[f64]) -> Vec<i64> {
    values.iter().map(|value| *value as i64).collect()
}
