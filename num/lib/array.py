"""数组辅助：对齐。"""

from __future__ import annotations

from collections.abc import Sequence

def align(x: Sequence[float], value: float, n: int) -> list[float]:
    """把 x 对齐到长度 n。"""
    out = list(x[:n])
    if len(out) < n:
        out.extend([value] * (n - len(out)))
    return out
