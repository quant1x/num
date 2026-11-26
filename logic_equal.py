"""逐元素相等 / 不等（spec）。

口径与 Go 的 `Equal` / `NotEqual` 一致：按 `==` 比较，**NaN 与自身不相等**；
两条序列必须等长。
"""

from __future__ import annotations

from collections.abc import Sequence


def equal(x: Sequence[float], y: Sequence[float]) -> list[bool]:
    """返回逐元素是否相等。"""
    if len(x) != len(y):
        raise ValueError("equal: 两条序列长度不一致")
    return [left == right for left, right in zip(x, y)]


def not_equal(x: Sequence[float], y: Sequence[float]) -> list[bool]:
    """返回逐元素是否不等。"""
    if len(x) != len(y):
        raise ValueError("not_equal: 两条序列长度不一致")
    return [left != right for left, right in zip(x, y)]
