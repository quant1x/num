# tests/fixtures/time —— 频率与时间序列样本

本目录用 **UTC 纳秒整数**表示时间，避免 npz 存时间类型或 ISO 解析差异。格式见 [../README.md](../README.md)。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `parse-freq` | `parse_freq` | D / 2H / 30min / 500ms → 纳秒 |
| `date-range` | `date_range` | 2026-01-01T00:00:00Z 起、4 个点、12H |

约定：
- `parse_freq` 的 `freq` 是 `|S8` 字符串数组，`expected` 是对应的 int64 纳秒。
- `date_range` 的 `start` 是 int64 纳秒、`periods` 是 int64、`freq` 是 `|S8` 标量，`expected` 是 int64 序列。
- conform 以 int64 精确比较，不经过 float64（`1.7e18` 纳秒会丢精度）。

频率单位与 Go 一致、大小写敏感：`N`/`ns`、`U`/`us`/`µs`、`L`/`ms`、`S`、`T`/`min`、`H`、`D`。
