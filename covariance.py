"""方差与协方差（spec）。

口径（与 Go 的 `Variance` / `Covariance` 一致）：总体口径，分母为 `n`；
长度不等或空输入返回 `0.0`（不抛异常）。
"""

from __future__ import annotations

from collections.abc import Sequence

from .mean import mean


def variance(x: Sequence[float]) -> float:
    """返回 x 的总体方差。"""
    if len(x) == 0:
        return 0.0
    average = mean(x)
    total = 0.0
    for value in x:
        delta = value - average
        total += delta * delta
    return total / len(x)


def covariance(x: Sequence[float], y: Sequence[float]) -> float:
    """返回 x 与 y 的总体协方差。"""
    if len(x) != len(y) or len(x) == 0:
        return 0.0
    mean_x = mean(x)
    mean_y = mean(y)
    total = 0.0
    for left, right in zip(x, y):
        total += (left - mean_x) * (right - mean_y)
    return total / len(x)
