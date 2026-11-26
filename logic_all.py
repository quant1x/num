"""全真判定（spec）。

口径：所有元素非零为 `True`；空输入为 `True`（空真，与 Go 的 `All` 一致）。
"""

from __future__ import annotations

import builtins
from collections.abc import Sequence


def all(x: Sequence[float]) -> bool:  # noqa: A001 - 与合同名 ALL 保持一致
    """返回是否所有元素都为真。"""
    return builtins.all(bool(value) for value in x)
