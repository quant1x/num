# tests/fixtures/conversion —— N1 转换族样本

本目录覆盖 bool / 整型 / 浮点之间的序列转换与 `dtype2int`，格式见 [../README.md](../README.md)。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `from-bool-f64` | `from_bool` | true=1 / false=0 |
| `from-int32-f64` / `from-int64-f64` | `from_int32` / `from_int64` | 整型转 f64 |
| `from-float64-f32` / `from-float32-f64` | `from_float` | 浮点宽度互转 |
| `to-bool-f64` | `to_bool` | 非零为真 |
| `to-int32-f64` / `to-int64-f64` | `to_int32` / `to_int64` | 向零截断 |
| `to-float32-f64` / `to-float64-f32` | `to_float32` / `to_float64` | 浮点宽度互转 |
| `dtype2int-f64` | `dtype2int` | 与 `to_int32` 同义 |

注意：`to_bool` 的样本不含 NaN。Go 的 AVX2 路径把 NaN 判为 false，而 Go 参考路径 `x != 0` 判为 true，
两侧不一致，已登记在 `num.md` §5；fixture 暂不覆盖该格。
