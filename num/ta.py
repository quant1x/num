"""num 专有的技术指标与滚动算子。"""

from __future__ import annotations

import math
from collections.abc import Callable, Sequence
from dataclasses import dataclass

def change_rate(base: float, current: float) -> float:
    """增长率 `current / base`。"""
    return current / base

def net_change_rate(base: float, current: float) -> float:
    """净增长率百分比：`(current/base - 1) * 100`。"""
    return (change_rate(base, current) - 1.0) * 100.0

def decimal(value: float, digits: int = 2) -> float:
    """把 value 保留 digits 位小数（半程远离零）。"""
    if digits < 0:
        digits = 0
    if math.isnan(value):
        value = 0.0
    half = -5.0 if math.copysign(1.0, value) < 0 else 5.0
    n10 = math.pow(10.0, digits)
    nj1 = math.pow(10.0, digits + 1)
    return math.trunc((value * nj1 + half) / 10.0) / n10

def _window_at(window: int | Sequence[int], index: int) -> int:
    if isinstance(window, int):
        return window
    return int(window[index])

def rolling(
    x: Sequence[float],
    window: int | Sequence[int],
    min_periods: int | None = None,
) -> list[list[float]]:
    """返回每个位置的滚动窗口。"""
    out: list[list[float]] = []
    for index in range(len(x)):
        size = _window_at(window, index)
        if size <= 0:
            out.append([])
            continue
        available = min(index + 1, size)
        threshold = size if min_periods is None else min_periods
        if available < threshold:
            out.append([])
            continue
        out.append(list(x[index + 1 - available:index + 1]))
    return out

def rolling_apply(
    x: Sequence[float],
    window: int | Sequence[int],
    apply: Callable[[Sequence[float]], float],
    min_periods: int | None = None,
) -> list[float]:
    """对每个窗口应用归约函数，空窗口为 NaN。"""
    out: list[float] = []
    for block in rolling(x, window, min_periods):
        out.append(math.nan if not block else float(apply(block)))
    return out

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

@dataclass(frozen=True)
class Line:
    """斜截式直线。"""

    slope: float = 0.0
    intercept: float = 0.0
    offset: float = 0.0

    def equation(self) -> tuple[float, float, float]:
        """返回 Ax + By + C = 0 的 (A, B, C)。"""
        return (self.slope, -1.0, self.intercept)

    def degrees(self) -> float:
        """斜率对应的角度（度）。"""
        return math.degrees(math.atan(self.slope))

    def y(self, x: float) -> float:
        """给定 x 求 y。"""
        return self.slope * x + self.intercept

def calculate_line_equation(x1: float, y1: float, x2: float, y2: float) -> Line:
    """两点确定直线。"""
    slope = (y2 - y1) / (x2 - x1)
    return Line(slope=slope, intercept=y1 - slope * x1, offset=x1)

def fit_line(x: Sequence[float], y: Sequence[float]) -> tuple[Line, bool]:
    """最小二乘拟合直线，返回 (Line, 是否成功)。"""
    count = len(x)
    if count == 0 or len(y) != count:
        return Line(), False
    if count < 2:
        return Line(slope=0.0, intercept=y[0], offset=x[0]), False
    sum_x = sum(x)
    sum_y = sum(y)
    sum_xy = sum(left * right for left, right in zip(x, y))
    sum_xx = sum(value * value for value in x)
    denominator = count * sum_xx - sum_x * sum_x
    if abs(denominator) < 1e-10:
        return Line(), False
    slope = (count * sum_xy - sum_x * sum_y) / denominator
    intercept = (sum_y - slope * sum_x) / count
    return Line(slope=slope, intercept=intercept, offset=x[0]), True

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
