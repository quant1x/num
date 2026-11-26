"""缺失值填充。"""

from __future__ import annotations

import math
from collections.abc import Sequence

def _is_missing(value: float) -> bool:
    return math.isnan(value) or math.isinf(value)

def fill(x: Sequence[float], value: float) -> list[float]:
    """把 NaN / ±Inf 换成 value。"""
    return [value if _is_missing(item) else item for item in x]

def fill_na(x: Sequence[float], value: float) -> list[float]:
    """把 NaN / ±Inf 换成 value（与 `fill` 同口径，入口名对齐 Go）。"""
    return fill(x, value)
