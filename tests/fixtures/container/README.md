# tests/fixtures/container —— N6 向量与容器样本

本目录覆盖构造、整形与索引类算子，格式见 [../README.md](../README.md)。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `zeros-f64` / `ones-f64` | `zeros` / `ones` | 全零 / 全一 |
| `arange-int` / `arange-step-f64` | `arange` | 整型步长、浮点步长；不含 end |
| `linspace-f64` / `linspace-endpoint-false` | `linspace` | endpoint 两种取值 |
| `range-int` / `range-f64` / `range-float-f32` | `range` / `range_float` | 0..n-1 与浮点区间 |
| `repeat-f64` / `repeat-int` | `repeat` | 同值重复 |
| `align-expand-f64` / `align-truncate-f64` | `align` | 补默认值 / 截断 |
| `shape-1d` / `shape-2d` | `shape` | 返回 (行, 列) |
| `concat1d-int` | `concat1d` | **n×3 列式拼接**（不是按行） |
| `transpose2d-int` | `transpose2d` | 二维转置 |
| `fill-f64` / `fill-na-f64` | `fill` / `fill_na` | 缺失值替换；样本只含 NaN，不含 ±Inf |

注意：

- `Ones` 在 Go 按入参切片长度生成，不接受形状整数；样本只存 `n`，由对拍用例构造长度为 n 的切片。
- `fill` / `fill_na` 会把 ±Inf 也当作缺失（NumPy 只填 NaN）；样本不覆盖 Inf，差异登记在
  [../CPP-GAPS.md](../CPP-GAPS.md) §3。
- `concat1d` 的期望值是 `np.stack([a, b, c], axis=1)`，与 NumPy `concatenate` 的默认轴不同。
