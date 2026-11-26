"""条件选择（spec）。

口径与 NumPy 的 `where` 一致：`condition` 为真取 `x`，否则取 `y`；三个序列等长。
Go 的 `Where` 要求条件与取值同类型，Python 侧按真值判断。
"""

from __future__ import annotations

from collections.abc import Sequence


def where(condition: Sequence[bool], x: Sequence[float], y: Sequence[float]) -> list[float]:
    """按条件从 x / y 中逐位置取值。"""
    if len(condition) != len(x) or len(x) != len(y):
        raise ValueError("where: 条件与取值长度不一致")
    return [left if flag else right for flag, left, right in zip(condition, x, y)]
