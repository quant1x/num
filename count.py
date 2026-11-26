"""非零个数（spec）。

口径：统计非零元素（布尔序列即 `True` 的个数）；空输入为 0。
"""

from __future__ import annotations

from collections.abc import Sequence


def count(x: Sequence[float] | Sequence[bool]) -> int:
    """返回 x 中非零元素的个数。"""
    total = 0
    for value in x:
        if value:
            total += 1
    return total
