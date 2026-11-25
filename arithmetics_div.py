"""逐元素除法（spec）。

口径：`y` 可以是标量或等长序列；长度不等抛 `ValueError`。
按 IEEE 754 处理除零（Python 的 `/` 会抛 `ZeroDivisionError`，这里显式实现）：

- 被除数为 `NaN` → `NaN`；
- `0/0` → `NaN`；
- 非零 / `±0` → `±Inf`，符号由两侧符号相乘决定。
"""

from __future__ import annotations

import math
from collections.abc import Sequence

from .arithmetics_add import _pairwise


def _divide(numerator: float, denominator: float) -> float:
    if math.isnan(numerator):
        return math.nan
    if denominator == 0:
        if numerator == 0:
            return math.nan
        return math.copysign(math.inf, numerator) * math.copysign(1.0, denominator)
    return numerator / denominator


def div(x: Sequence[float], y: float | Sequence[float]) -> list[float]:
    """返回逐元素商。"""
    return [_divide(left, right) for left, right in _pairwise(x, y)]
