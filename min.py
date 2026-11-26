"""最小值（spec）。

口径：顺序扫描，返回首个最小值的**值**（不是下标）；空输入返回 `0.0`（与 Go 的 `Min` 一致）。
"""

from __future__ import annotations

from collections.abc import Sequence


def min(x: Sequence[float]) -> float:  # noqa: A001 - 与合同名 MIN 保持一致
    """返回 x 的最小值。"""
    if len(x) == 0:
        return 0.0
    best = x[0]
    for value in x[1:]:
        if value < best:  # noqa: PLR1730 - 调用内建 min 会递归到本函数
            best = value
    return best
