"""最小值下标（spec）。

口径：返回**首个**最小值的下标；空输入返回 `-1`。
"""

from __future__ import annotations

from collections.abc import Sequence


def argmin(x: Sequence[float]) -> int:
    """返回 x 中首个最小值的下标。"""
    if len(x) == 0:
        return -1
    best = 0
    for index in range(1, len(x)):
        if x[index] < x[best]:
            best = index
    return best
