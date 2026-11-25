"""直线族（spec，对应 Go 的 `slope.go` 里 Line / FitLine / CalculateLineEquation）。

直线用**斜截式** `y = slope·x + intercept` 表达；`equation()` 给出 `Ax + By + C = 0`
的 `(A, B, C) = (slope, -1, intercept)`，`degrees()` 是 `atan(slope)` 的角度。

`fit_line` 在样本不足或 x 全等（垂直线）时返回 `success=False`，此时 Line 的取值按 Go 的
退化策略：单点给 `slope=0, intercept=y[0]`，其余给全 0。
"""

from __future__ import annotations

import math
from dataclasses import dataclass
from collections.abc import Sequence


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
