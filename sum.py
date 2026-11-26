"""求和（spec）。

口径（与 Go 的 `Sum` / `__go_sum` 一致）：

- 顺序累加，不做分块或并行 —— 归约顺序敏感（CONTRIBUTING §5.6.3）；
- 空输入返回 `0.0`；
- `NaN` 参与累加后结果仍为 `NaN`（IEEE 754 传播）。
"""

from __future__ import annotations

from collections.abc import Sequence


def sum(x: Sequence[float]) -> float:  # noqa: A001 - 与合同名 SUM 保持一致
    """返回 x 的累加和。"""
    total = 0.0
    for value in x:
        total += value
    return total
