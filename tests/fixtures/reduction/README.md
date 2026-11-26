# tests/fixtures/reduction —— N3 归约样本

本目录的 npz 按 `tests/contract/NAME-MAP.tsv` 的语义名组织，格式见 [../README.md](../README.md)。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `sum-basic-f64` / `mean-basic-f64` | `sum` / `mean` | 基础归约 |
| `sum-long-f64` / `sum-long-f32` | `sum` | 100 个 `0.1`；容差 1e-12 / 1e-6，归约顺序敏感 |
| `sum-nan-f64` | `sum` | NaN 传播 |
| `std-basic-f64` / `var-basic-f64` | `std` / `variance` | ddof=0，容差 1e-12 |
| `min-basic-f64` / `max-basic-f64` | `min` / `max` | 极值 |
| `argmin-basic-f64` / `argmax-basic-f64` | `argmin` / `argmax` | 返回下标（i64） |
| `count-numeric-f64` / `count-bool` | `count` | 非零计数 / true 计数 |
| `cumsum-basic-f64` | `cumsum` | 前缀和 |
| `prod-basic-f64` | `prod` | 含负数 |
| `dot1d-basic-f64` | `dot1d` | 一维点积 |
| `median-basic-f64` / `median-even-f64` | `median` | 排序后取中间值：奇长取中点，偶长取中间两值平均 |
| `median-single-f64` / `median-nan-f64` | `median` | 单元素返回其本身；NaN 传播 |
