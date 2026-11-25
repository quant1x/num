"""dtype 相关的标量转换（spec）。

- `dtype2int`：向零截断（1.9 → 1、-2.9 → -2），不是四舍五入；
- `is_nan` 见 `decimals.py`（NaN 与 ±Inf 都算缺失）。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def dtype2int(x: Sequence[float]) -> list[int]:
    """返回向零截断后的整数序列。"""
    return [int(math.trunc(value)) for value in x]
