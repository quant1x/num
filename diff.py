"""一阶（或 n 阶）离散差分（spec）。

口径与 Go 一致：`diff[i] = x[i-n] - x[i]` —— 方向与 pandas / NumPy **相反**；
越界位为 `NaN`。
"""

from __future__ import annotations

from collections.abc import Sequence

from .shift import shift


def diff(x: Sequence[float], n: int = 1) -> list[float]:
    """返回方向为 `x[i-n] - x[i]` 的差分序列。"""
    shifted = shift(x, n)
    return [previous - value for previous, value in zip(shifted, x)]
