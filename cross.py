"""上下穿（spec，对应 Go 的 `cross.go`）。

返回 `(索引, 状态)` 列表：上穿 `1`、下穿 `-1`；比较是严格的（相等不算穿）。
"""

from __future__ import annotations

from collections.abc import Sequence

BREAK_THROUGH = 1
FALL_DRASTICALLY = -1


def cross(a: Sequence[float], b: Sequence[float]) -> list[tuple[int, int]]:
    """返回 a 相对 b 的穿点与方向。"""
    if len(a) != len(b):
        raise ValueError("cross: 两条序列长度不一致")
    out: list[tuple[int, int]] = []
    for index in range(1, len(a)):
        front = index - 1
        if a[front] < b[front] and a[index] > b[index]:
            out.append((index, BREAK_THROUGH))
        elif a[front] > b[front] and a[index] < b[index]:
            out.append((index, FALL_DRASTICALLY))
    return out
