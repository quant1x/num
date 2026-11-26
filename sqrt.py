"""逐元素平方根（spec）。

口径：负数返回 `NaN`（不是抛异常），`NaN` 仍为 `NaN`。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def sqrt(x: Sequence[float]) -> list[float]:
    """返回 x 的逐元素平方根。"""
    out: list[float] = []
    for value in x:
        if math.isnan(value) or value < 0:
            out.append(math.nan)
        else:
            out.append(math.sqrt(value))
    return out
