"""保留小数位（spec）。

口径**照抄 Go 的 `Decimal`**（半程远离零，不是 Python 的银行家舍入）：

```text
trunc((value × 10^(d+1) + half) / 10) / 10^d      # half = ±5，符号随 value
```

- `digits` 缺省 2，负数按 0 处理；
- `NaN` 先按 0 处理；
- `IsNaN` 在 Go 里连 `±Inf` 一起算缺失 —— 本模块的 `is_nan` 保持一致。
"""

from __future__ import annotations

import math


def is_nan(value: float) -> bool:
    """NaN 或 ±Inf 都算缺失（与 Go 的 `IsNaN` 一致）。"""
    return math.isnan(value) or math.isinf(value)


def change_rate(base: float, current: float) -> float:
    """增长率 `current / base`。"""
    return current / base


def net_change_rate(base: float, current: float) -> float:
    """净增长率百分比：`(current/base - 1) * 100`。"""
    return (change_rate(base, current) - 1.0) * 100.0


def decimal(value: float, digits: int = 2) -> float:
    """把 value 保留 digits 位小数（半程远离零）。"""
    digits = max(digits, 0)
    if math.isnan(value):
        value = 0.0
    half = -5.0 if math.copysign(1.0, value) < 0 else 5.0
    n10 = math.pow(10.0, digits)
    nj1 = math.pow(10.0, digits + 1)
    return math.trunc((value * nj1 + half) / 10.0) / n10
