# tests/fixtures/scalar —— N1 标量与类型边界样本

本目录覆盖类型默认值、NaN / ±Inf 判据与标量/切片转换，格式见 [../README.md](../README.md)。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `type-default-f64` / `-f32` / `-i64` / `-bool` | `type_default` | 浮点 NaN、整型 0、bool false |
| `nan-check-f64` / `-f32` | `nan_check` | `NaN` / `±Inf` 判为真，`-0.0` 与普通值判为假 |
| `bool-to-float64` / `bool-to-int64` | `bool_to_numeric` | true=1 / false=0 |
| `slice-to-float64-int64` / `-bool` | `slice_to_float64` | 整型与布尔切片转 f64 |
| `slice-float-to-int32` / `-int64` | `slice_float_to_int` | 向零截断 |

注意：

- `Float32IsNaN` / `Float64IsNaN` **把 ±Inf 也算作无效值**，与 NumPy 的 `isnan` 不同；样本按 Go 口径断言。
- `BoolToFloat64` / `BoolToInt64` 是标量入口，样本按数组逐元素调用；不是切片 API。
