"""中位数（spec）。

口径**照抄 Go 的 `Median`，不排序**：

- 空输入 → `NaN`；
- 长度 1 → `0.0`（Go 的既有行为，属已知怪癖，待提案修正）；
- 偶数长度 → 位置 `n/2-1` 与 `n/2` 两个元素的平均；
- 奇数长度 → 位置 `n/2` 的元素。

因为不排序，"中位数"只是**按入参位置**取中间，调用方需自行保证有序。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def median(x: Sequence[float]) -> float:
    """返回按入参位置取到的中位数。"""
    if len(x) == 0:
        return math.nan
    if len(x) == 1:
        return 0.0
    if len(x) % 2 == 0:
        middle = len(x) // 2
        return (x[middle - 1] + x[middle]) / 2
    return x[len(x) // 2]
