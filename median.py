"""中位数（spec）。

口径与 `np.median` 一致：

- 空输入 → `NaN`；
- 任一元素为 `NaN` → `NaN`（NaN 传播）；
- 其余先排序：偶数长度取中间两值均值，奇数长度取中点，单元素返回其本身。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def median(x: Sequence[float]) -> float:
    """返回排序后的中位数；偶数长度取中间两值均值，任一 NaN 传播。"""
    if len(x) == 0:
        return math.nan
    for value in x:
        if value != value:  # noqa: PLR0124 - NaN 自比较即 NaN 判定
            return math.nan
    ordered = sorted(x)
    if len(ordered) % 2 == 0:
        middle = len(ordered) // 2
        return (ordered[middle - 1] + ordered[middle]) / 2
    return float(ordered[len(ordered) // 2])
