"""逐元素取大（spec）。

口径与 Go 一致：**只在 `y > x` 时覆盖**，因此 `NaN` 不传播（`x` 为 `NaN` 且 `y` 不是时结果仍是 `NaN`）。
两条序列必须等长，否则抛 `ValueError`。
"""

from __future__ import annotations

from collections.abc import Sequence


def maximum(x: Sequence[float], y: Sequence[float]) -> list[float]:
    """返回逐元素较大者。"""
    if len(x) != len(y):
        raise ValueError("maximum: 两条序列长度不一致")
    return [max(left, right) for left, right in zip(x, y)]
