"""滚动窗口（spec）。

口径与 Go 的 `Rolling` / `RollingApply` 一致：

- `window` 可以是整数或逐位置序列；第 i 位的窗口长度是 `min(i+1, window[i])`；
- `min_periods` 缺省等于**该位置的窗口长度**（即只算完整窗口）；
  `min_periods > 窗口长度` 时该位置给空窗口；
- `rolling` 返回每个位置的窗口（列表的列表），空位是空列表；
- `rolling_apply` 对每个窗口调用归约函数，空窗口返回 `NaN`。
"""

from __future__ import annotations

import math
from collections.abc import Callable, Sequence


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
