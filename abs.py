"""逐元素绝对值（spec）。

口径：逐元素 `abs`，保持输入长度；`NaN` 仍为 `NaN`，`-0.0` 归一为 `0.0`。
"""

from __future__ import annotations

import builtins
from collections.abc import Sequence


def abs(x: Sequence[float]) -> list[float]:  # noqa: A001 - 与合同名 ABS 保持一致
    """返回 x 的逐元素绝对值。"""
    return [builtins.abs(value) for value in x]
