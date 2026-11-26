"""dtype 相关的标量转换（spec）。

- `dtype2int`：向零截断（1.9 → 1、-2.9 → -2），不是四舍五入；
- `is_nan` 见 `decimals.py`（NaN 与 ±Inf 都算缺失）。
"""

from __future__ import annotations

import math
from collections.abc import Sequence


def dtype2int(x: Sequence[float]) -> list[int]:
    """返回向零截断后的整数序列。"""
    return [math.trunc(value) for value in x]


# dtype 码 → 默认值：浮点为 NaN，整型为 0，布尔为 False。
_DEFAULTS = {"f4": math.nan, "f8": math.nan, "i1": 0, "i2": 0, "i4": 0, "i8": 0,
             "u1": 0, "u2": 0, "u4": 0, "u8": 0, "b1": False}


def type_default(dtype: str):
    """返回该 dtype 的默认值。"""
    if dtype not in _DEFAULTS:
        raise ValueError(f"type_default: 未知 dtype {dtype!r}")
    return _DEFAULTS[dtype]
