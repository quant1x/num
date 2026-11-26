"""逐元素加法（spec）。

口径：`y` 可以是标量或等长序列；两条序列长度不等抛 `ValueError`。
"""

from __future__ import annotations

from collections.abc import Sequence

Number = int | float


def _pairwise(x: Sequence[float], y: Number | Sequence[float]) -> list[tuple[float, float]]:
    if isinstance(y, (int, float)) and not isinstance(y, bool):
        return [(value, float(y)) for value in x]
    if len(x) != len(y):
        raise ValueError("两条序列长度不一致")
    return list(zip(x, y))


def add(x: Sequence[float], y: Number | Sequence[float]) -> list[float]:
    """返回逐元素和。"""
    return [left + right for left, right in _pairwise(x, y)]
