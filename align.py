"""长度对齐（spec）。

口径与 Go 的 `Align` 一致：长度不足用 `value` 补齐，超出则截断，结果长度恰为 `n`。
"""

from __future__ import annotations

from collections.abc import Sequence


def align(x: Sequence[float], value: float, n: int) -> list[float]:
    """把 x 对齐到长度 n。"""
    out = list(x[:n])
    if len(out) < n:
        out.extend([value] * (n - len(out)))
    return out
