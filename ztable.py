"""正态分布表查询（spec，对应 Go 的 `ztable.go`）。

数据在 `ztable_data.py`（两张 10001 项的表，由 Go 侧数据导出），本模块只做索引与二分：

- `find_percent(z)`：表里第一个 ≥ z 的位置 / 10000；
- `find_zscore(percent)`：`table[int(percent*10000) % 10000]`；
- `confidence_interval_to_zscore(interval)`：同上，查另一张表；
- `zscore_to_confidence_interval(z)`：Go 的近似二分（比较器用 `int((z-a[i])*10000)` 的符号），
  未命中时取前一位，再 / 10000。Go 把 `±Inf` / 溢出转成 minInt64，这里用 `_go_int64` 显式模拟。

数据资产 `ztable_data.py` 有 2 万项，故本模块**不随 `import num` 加载**，用 `from num import ztable` 取用。
"""

from __future__ import annotations

import bisect
import math

from .ztable_data import percentToZscore, z_table

SCALE = 10000


def _go_int64(value: float) -> int:
    """模拟 Go 的 float64 → int 转换：溢出、±Inf、NaN 都落到 minInt64（amd64 语义）。"""
    if math.isnan(value) or value >= 2**63 or value <= -(2**63):
        return -(2**63)
    return int(value)


def find_percent(zscore: float) -> float:
    """返回 z 分对应的累计概率。"""
    index = bisect.bisect_left(percentToZscore, zscore)
    return index / SCALE


def find_zscore(percent: float) -> float:
    """按累计概率查 z 分。"""
    return percentToZscore[int(percent * SCALE) % SCALE]


def confidence_interval_to_zscore(interval: float) -> float:
    """按置信区间查 z 分。"""
    return z_table[int(interval * SCALE) % SCALE]


def zscore_to_confidence_interval(zscore: float) -> float:
    """按 z 分查置信区间（近似二分，与 Go 的 `__SearchFloat64s` 同口径）。"""
    low, high = 0, len(z_table)
    while low < high:
        middle = (low + high) // 2
        if _go_int64((zscore - z_table[middle]) * SCALE) < 0:
            high = middle
        else:
            low = middle + 1
    found = low < len(z_table) and _go_int64((zscore - z_table[low]) * SCALE) == 0
    index = low if found else low - 1
    return index / SCALE
