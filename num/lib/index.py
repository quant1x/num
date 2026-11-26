"""索引与掩码：take / compress / put / find。"""

from __future__ import annotations

from collections.abc import Sequence

def find(x: Sequence[float], value: float) -> int:
    """返回首次出现的下标；未命中返回 -1。"""
    for index, item in enumerate(x):
        if item == value:
            return index
    return -1

def take(x: Sequence[float], idx: Sequence[int]) -> list[float]:
    """按索引取值。"""
    return [x[index] for index in idx]

def compress(x: Sequence[float], mask: Sequence[bool]) -> list[float]:
    """按掩码过滤。"""
    if len(x) != len(mask):
        raise ValueError("select: 长度不一致")
    return [value for value, keep in zip(x, mask) if keep]

def scatter(x: Sequence[float], idx: Sequence[int], size: int) -> list[float]:
    """把 x 按索引散落到长度为 size 的零序列。"""
    if len(x) != len(idx):
        raise ValueError("scatter: 值与索引长度不一致")
    out = [0.0] * size
    for value, position in zip(x, idx):
        out[position] = value
    return out
