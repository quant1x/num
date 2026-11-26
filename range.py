"""等距序列（spec）。

- `arange(start, end, step=1)`：按累加生成，`v < end` 为止（浮点累加与 Go 一致，不做乘法重算）；
- `linspace(start, end, n, endpoint=True)`：`step = (end-start)/(n-1)`（`endpoint=False` 时分母为 `n`），
  第 i 个元素为 `start + i*step`；
- `range(n)`：`0 .. n-1`。
"""

from __future__ import annotations

import builtins


def arange(start: float, end: float, step: float = 1.0) -> list[float]:
    """返回 [start, end) 上的等距序列。"""
    out: list[float] = []
    value = start
    while value < end:
        out.append(value)
        value += step
    return out


def linspace(start: float, end: float, n: int, endpoint: bool = True) -> list[float]:
    """返回 n 个等距点。"""
    if n <= 0:
        return []
    if n == 1:
        return [start]
    segments = n - 1 if endpoint else n
    step = (end - start) / segments
    return [start + index * step for index in builtins.range(n)]


def range_float(a: float, b: float) -> list[float]:
    """返回 [a, b) 上的浮点序列（步长 1）。"""
    return arange(a, b, 1.0)


def range(n: int) -> list[int]:  # noqa: A001 - 与合同名 RANGE 保持一致
    """返回 0 .. n-1。"""
    if n <= 0:
        return []
    return list(builtins.range(n))
