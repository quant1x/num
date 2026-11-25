# tests/fixtures/elementwise —— N2 逐元素样本

本目录的 npz 按 `tests/contract/NAME-MAP.tsv` 的语义名组织，格式见 [../README.md](../README.md)。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `add-f64` / `add-f32` | `add` | f64 / f32 两条宽度 |
| `sub-f64` / `mul-f64` / `div-f64` | `sub` / `mul` / `div` | 正负值与零 |
| `div-zero-f64` | `div` | 除零得 `±inf`（含 `-0.0`） |
| `maximum-f64` / `minimum-f64` | `maximum` / `minimum` | 逐元素取大 / 取小 |
| `abs-f64` | `abs` | `-0.0` 归零 |
| `sqrt-f64` | `sqrt` | 0 / 1 / 2 / 4 |
| `pow-integer-f64` | `pow_integer` | 整数次幂，含负数 |
| `equal-f64` / `not-equal-f64` | `equal` / `not_equal` | NaN 与自身比较为 false / true |
| `greater-f64` | `greater` | 逐元素比较 |
| `logic-and` / `logic-or` / `logic-not` | `logic_*` | 布尔逐元素 |
| `where-f64` | `where` | 条件选择 |
| `add-empty-f64` | `add` | 零长度输入 |

样本只读：修改期望值必须同时改生成器与三套读取器，并在提交信息里说明依据。
