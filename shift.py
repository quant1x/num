"""索引平移（spec）。

口径与 Go 一致：`n > 0` 向右（取 `x[i-n]`），`n < 0` 向左（取 `x[i-n]`），越界位填 `NaN`；
`n == 0` 返回原序列的副本。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def shift(x: Sequence[float], n: int) -> list[float]:
    """返回平移 n 位后的序列。"""
    out = [math.nan] * len(x)
    for index in range(len(x)):
        source = index - n
        if 0 <= source < len(x):
            out[index] = x[source]
    return out
