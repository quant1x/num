"""逐元素乘法（spec）。

口径：`y` 可以是标量或等长序列；长度不等抛 `ValueError`。
"""

from __future__ import annotations

from collections.abc import Sequence

from .arithmetics_add import _pairwise


def mul(x: Sequence[float], y: float | Sequence[float]) -> list[float]:
    """返回逐元素积。"""
    return [left * right for left, right in _pairwise(x, y)]
