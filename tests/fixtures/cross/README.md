# tests/fixtures/cross —— 穿越样本

本目录覆盖 `cross`：a 相对 b 的上穿 / 下穿，格式见 [../README.md](../README.md)。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `cross-f64` | `cross` | a=[1,3,1,3]、b=[2,2,2,2]，穿越下标 [1,2,3]、状态 [1,-1,1] |

约定：fixture 用 `expected_x`（下标）与 `expected_state`（上穿 1、下穿 -1）两个数组描述结果；
conform 按顺序逐个比较。判据是**前一根与当前根都严格分居两侧**，相等不计。
