"""一维点积（spec）。

口径（与 Go 的 `Dot1D` 一致）：按 `a` 的长度顺序累加 `a[i] * b[i]`；
长度不等抛 `ValueError`（Go 侧是 panic 或静默截断，Python 选择显式失败）。
"""

from __future__ import annotations

from collections.abc import Sequence


def dot1d(a: Sequence[float], b: Sequence[float]) -> float:
    """返回 a 与 b 的内积。"""
    if len(a) != len(b):
        raise ValueError("dot1d: 两条序列长度不一致")
    total = 0.0
    for left, right in zip(a, b):
        total += left * right
    return total
