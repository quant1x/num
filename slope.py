"""斜率与角度互转（spec，对应 Go 的 `slope.go`）。

- `slope(x1, y1, x2, y2)`：两点斜率 `(y2-y1)/(x2-x1)`；
- `slope_to_degrees(slope)`：`atan(slope)` 转角度；
- `degrees_to_slope(degrees)`：`tan` 转斜率；
- `triangle_bevel(slope, x1, y1, xn)`：`slope*(xn-x1)+y1`。
"""

from __future__ import annotations

import math


def slope(x1: int, y1: float, x2: int, y2: float) -> float:
    """两点斜率。"""
    return (y2 - y1) / (x2 - x1)


def slope_to_degrees(value: float) -> float:
    """斜率转角度。"""
    return math.degrees(math.atan(value))


def degrees_to_slope(degrees: float) -> float:
    """角度转斜率。"""
    return math.tan(math.radians(degrees))


def triangle_bevel(value: float, x1: int, y1: float, xn: int) -> float:
    """按斜率延拓到 xn 处的值。"""
    return value * (xn - x1) + y1
