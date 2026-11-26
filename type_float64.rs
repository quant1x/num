//! 与 type_float64.go 同语义的 Rust 实现。

/// ints_to_f64 把整型切片转成 f64 切片（SliceToFloat64 的 i64 分支）。
pub fn ints_to_f64(values: &[i64]) -> Vec<f64> {
    values.iter().map(|value| *value as f64).collect()
}
