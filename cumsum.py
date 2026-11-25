"""前缀和（spec）。

口径：逐位置累加，返回与输入等长的新序列；空输入返回空序列。
"""

from __future__ import annotations

from collections.abc import Sequence


def cumsum(x: Sequence[float]) -> list[float]:
    """返回 x 的前缀和。"""
    out: list[float] = []
    total = 0.0
    for value in x:
        total += value
        out.append(total)
    return out
