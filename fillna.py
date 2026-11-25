"""缺失值填充（spec）。

口径与 Go 一致：`NaN` 与 `±Inf` 都算缺失（`Float64IsNaN`）。

Go 留了两个入口 —— `Fill`（字符串 / f32 / f64）与 `FillNa`（额外覆盖 int64）。
浮点序列上两者行为相同，Python 各保留一个函数以便按名对位。
"""

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
