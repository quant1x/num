"""任一为真（spec）。

口径：任一元素非零为 `True`；空输入为 `False`。
"""

from __future__ import annotations

import builtins
from collections.abc import Sequence


def any(x: Sequence[float]) -> bool:  # noqa: A001 - 与合同名 ANY 保持一致
    """返回是否存在为真的元素。"""
    return builtins.any(bool(value) for value in x)
