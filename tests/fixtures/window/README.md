# tests/fixtures/window —— N4 滚动族样本

本目录覆盖滑动窗口、移动统计与位移/差分，格式见 [../README.md](../README.md)。

嵌套结果（`Rolling` 返回 `[][]T`）用两个数组编码：`expected_flat`（按窗口顺序展平）与
`expected_lengths`（每个窗口的长度）。空窗口用长度 0 表示，不补 NaN。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `rolling-basic-f64` | `rolling` | 常量窗口 3，前两个位置为空 |
| `rolling-min1-f64` | `rolling` | `min_periods=1`，首位置即出窗口 |
| `rolling-min2-f64` | `rolling` | `min_periods=2`，partial window |
| `rolling-min-window-f64` | `rolling` | `min_periods == window`，与不传 min 一致 |
| `rolling-min-gt-window-f64` | `rolling` | `min_periods > window`，全部为空 |
| `rolling-window-array-f64` | `rolling` | 逐位置窗口 + min_periods |
| `rolling-apply-mean-f64` | `rolling_apply` | 移动均值，窗口不足处为 NaN |
| `shift-right-f64` / `shift-left-f64` / `shift-zero-f64` | `shift` | 正 / 负 / 零位移，越界补默认值 |
| `diff-n1-f64` / `diff-n2-f64` | `diff` | 方向为 `s[i-n]-s[i]`，前 n 位为 NaN |
| `cumprod-f64` | `cumprod` | 前缀积 |

`min_periods` 的取值口径：

- 不传或 `<= 0`：要求完整窗口，窗口不足的位置返回空切片；
- `1 <= min_periods <= window`：可用观测数达到 `min_periods` 即计算，窗口不足时用已有观测；
- `min_periods > window`：不产生任何窗口（与 pandas 的非法参数不同，这里选择静默跳过，已在 README 写明）。
