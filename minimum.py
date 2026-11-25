"""逐元素取小（spec）。

口径与 Go 一致：**只在 `y < x` 时覆盖**，`NaN` 不传播；两条序列必须等长。
"""

from __future__ import annotations

from collections.abc import Sequence


def minimum(x: Sequence[float], y: Sequence[float]) -> list[float]:
    """返回逐元素较小者。"""
    if len(x) != len(y):
        raise ValueError("minimum: 两条序列长度不一致")
    return [right if right < left else left for left, right in zip(x, y)]
