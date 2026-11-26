"""比较与逻辑族（spec）。

口径：`x` 为标量或等长序列；比较对 `NaN` 一律为 `False`；逻辑按「非零为真」判定。
`and` / `or` / `not` 是 Python 关键字，故落成 `logic_and` / `logic_or` / `logic_not`。
"""

from __future__ import annotations

from collections.abc import Sequence

from .arithmetics_add import _pairwise


def gt(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x > y`。"""
    return [left > right for left, right in _pairwise(x, y)]


def gte(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x >= y`。"""
    return [left >= right for left, right in _pairwise(x, y)]


def lt(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x < y`。"""
    return [left < right for left, right in _pairwise(x, y)]


def lte(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x <= y`。"""
    return [left <= right for left, right in _pairwise(x, y)]


def logic_and(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素逻辑与（两侧都非零为真）。"""
    return [bool(left) and bool(right) for left, right in _pairwise(x, y)]


def logic_or(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素逻辑或（任一侧非零为真）。"""
    return [bool(left) or bool(right) for left, right in _pairwise(x, y)]


def logic_not(x: Sequence[float]) -> list[bool]:
    """逐元素逻辑非。"""
    return [not bool(value) for value in x]
