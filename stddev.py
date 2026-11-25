"""总体标准差（spec）。

口径（与 Go 的 `Std` 一致）：两遍算法 —— 先求均值，再求「离差平方」的均值，最后开方；
分母是 `n`（总体口径，不是样本的 `n-1`）；空输入返回 `0.0`。
"""

from __future__ import annotations

import math
from collections.abc import Sequence

from .mean import mean


def std(x: Sequence[float]) -> float:
    """返回 x 的总体标准差。"""
    if len(x) == 0:
        return 0.0
    average = mean(x)
    total = 0.0
    for value in x:
        delta = value - average
        total += delta * delta
    return math.sqrt(total / len(x))
