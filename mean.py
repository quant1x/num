"""均值（spec）。

口径：`sum(x) / len(x)`，顺序累加；空输入为 `NaN`（与 Go 的 `Mean` 在空切片上 0/0 的结果一致）。
"""

from __future__ import annotations

import math
from collections.abc import Sequence

from .sum import sum


def mean(x: Sequence[float]) -> float:
    """返回 x 的算术平均。"""
    if len(x) == 0:
        return math.nan
    return sum(x) / len(x)
