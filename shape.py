"""形状查询（spec）。

口径与 Go 的 `Shape` 一致：一维返回 `(len, 0)`，二维返回 `(行, 列)`。
"""

from __future__ import annotations


def shape(x) -> tuple[int, int]:
    """返回 (行, 列)；一维的第二维为 0。"""
    if len(x) == 0:
        return (0, 0)
    first = x[0]
    if isinstance(first, (list, tuple)):
        return (len(x), len(first))
    return (len(x), 0)
