"""整型切片转换（spec，对应 Go 的 `SliceToInts`）。

口径：数值向零截断（1.9 → 1、-2.9 → -2），布尔按 1 / 0 处理。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def slice_to_int(x: Sequence[float | bool]) -> list[int]:
    """把数值切片转成整型切片（向零截断）。"""
    return [math.trunc(float(value)) for value in x]
