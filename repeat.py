"""常量序列（spec）。

口径与 Go 的 `Repeat` 一致：`n < 1` 返回空序列。
"""

from __future__ import annotations


def repeat(value: float, n: int) -> list[float]:
    """返回 n 个 value。"""
    if n < 1:
        return []
    return [value] * n
