"""零序列（spec）。

口径：返回长度为 `n` 的全零序列；`n < 1` 返回空序列。
"""

from __future__ import annotations


def ones(n: int) -> list[float]:
    """返回 n 个 1.0。"""
    if n < 1:
        return []
    return [1.0] * n


def zeros(n: int) -> list[float]:
    """返回 n 个 0.0。"""
    if n < 1:
        return []
    return [0.0] * n
