"""指数与对数的「1±x」形态（spec）。

`expm1(x) = e^x - 1`、`log1p(x) = ln(1 + x)`；小值附近用专用实现保精度（Python 的 `math` 已提供）。
`log1p(-1)` 返回 `-Inf`，小于 -1 返回 `NaN`（IEEE 754）。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def expm1(x: Sequence[float]) -> list[float]:
    """返回逐元素 `e^x - 1`。"""
    return [math.expm1(value) for value in x]


def log1p(x: Sequence[float]) -> list[float]:
    """返回逐元素 `ln(1 + x)`。"""
    out: list[float] = []
    for value in x:
        if value == -1:
            out.append(-math.inf)
        elif value < -1:
            out.append(math.nan)
        else:
            out.append(math.log1p(value))
    return out
