"""布尔数值转换（spec）。

口径：`True → 1`、`False → 0`；分别提供 int64 与 float64 两种落点。
"""

from __future__ import annotations

from collections.abc import Sequence


def bool_to_int64(x: Sequence[bool]) -> list[int]:
    """返回 1 / 0 序列。"""
    return [1 if value else 0 for value in x]


def bool_to_float64(x: Sequence[bool]) -> list[float]:
    """返回 1.0 / 0.0 序列。"""
    return [1.0 if value else 0.0 for value in x]


def slice_to_bool(x: Sequence[float | int | bool]) -> list[bool]:
    """把数值切片转成布尔切片：非零为真。"""
    return [bool(value) for value in x]
