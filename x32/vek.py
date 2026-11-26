"""向量算子（spec，对应 Go 的 `x32/vek.go`）。

一份实现同时承担 `x32` 与 `x64` 两个宽度：Python 没有 f32 / f64 之分，语义按 f64 计算，
f32 样本按 fixture 的容差比对。

口径：

- `neg` / `inv` 逐元素；`inv` 按 IEEE 754 处理除零（`1/±0 → ±Inf`、`1/NaN → NaN`）；
- `prod` 顺序累乘，空输入为 `1.0`；`cumprod` 为前缀积；
- `quantile` 先在副本上排序，再按 `f = (n-1)q`、`i = floor(f)` 线性插值；`q` 超出 [0,1] 抛 `ValueError`；
- `norm` / `manhattan_norm` / `distance` / `manhattan_distance` / `cosine_similarity` 顺序累加；
- `logic_xor` / `none` 按布尔处理（`none` 空输入为 `True`）。
"""

from __future__ import annotations

import builtins
import math
import struct as _struct
from collections.abc import Sequence

from ..arithmetics_div import _divide

_math = math


def neg(x: Sequence[float]) -> list[float]:
    """逐元素取相反数。"""
    return [-value for value in x]


def inv(x: Sequence[float]) -> list[float]:
    """逐元素取倒数（按 IEEE 754 处理除零）。"""
    return [_divide(1.0, value) for value in x]


def prod(x: Sequence[float]) -> float:
    """返回所有元素的乘积；空输入为 1.0。"""
    total = 1.0
    for value in x:
        total *= value
    return total


def cumprod(x: Sequence[float]) -> list[float]:
    """返回前缀积。"""
    out: list[float] = []
    total = 1.0
    for value in x:
        total *= value
        out.append(total)
    return out


def quantile(x: Sequence[float], q: float) -> float:
    """返回 q 分位（线性插值）。"""
    if not 0.0 <= q <= 1.0:
        raise ValueError("quantile: q 必须在 [0, 1] 之间")
    if len(x) == 0:
        raise ValueError("quantile: 序列为空")
    ordered = sorted(x)
    if len(ordered) == 1:
        return ordered[0]
    if q == 0.0:
        return ordered[0]
    if q == 1.0:
        return ordered[-1]
    position = (len(ordered) - 1) * q
    lower = math.floor(position)
    weight = position - lower
    first = ordered[lower]
    second = ordered[lower + 1]
    return first + (second - first) * weight


def norm(x: Sequence[float]) -> float:
    """欧氏范数。"""
    total = 0.0
    for value in x:
        total += value * value
    return math.sqrt(total)


def manhattan_norm(x: Sequence[float]) -> float:
    """曼哈顿范数。"""
    total = 0.0
    for value in x:
        total += abs(value)
    return total


def distance(x: Sequence[float], y: Sequence[float]) -> float:
    """欧氏距离。"""
    if len(x) != len(y):
        raise ValueError("distance: 两条序列长度不一致")
    total = 0.0
    for left, right in zip(x, y):
        delta = left - right
        total += delta * delta
    return math.sqrt(total)


def manhattan_distance(x: Sequence[float], y: Sequence[float]) -> float:
    """曼哈顿距离。"""
    if len(x) != len(y):
        raise ValueError("manhattan_distance: 两条序列长度不一致")
    total = 0.0
    for left, right in zip(x, y):
        total += abs(left - right)
    return total


def cosine_similarity(x: Sequence[float], y: Sequence[float]) -> float:
    """余弦相似度。"""
    if len(x) != len(y):
        raise ValueError("cosine_similarity: 两条序列长度不一致")
    dot = 0.0
    norm_x = 0.0
    norm_y = 0.0
    for left, right in zip(x, y):
        dot += left * right
        norm_x += left * left
        norm_y += right * right
    return dot / math.sqrt(norm_x * norm_y)


def logic_xor(x: Sequence[bool], y: Sequence[bool]) -> list[bool]:
    """逐元素逻辑异或。"""
    if len(x) != len(y):
        raise ValueError("logic_xor: 两条序列长度不一致")
    return [bool(left) != bool(right) for left, right in zip(x, y)]


def _narrow_f32(value: float) -> float:
    """把 f64 收窄到 f32 再回落 f64，用于模拟 f32 精度。"""
    return _struct.unpack("f", _struct.pack("f", value))[0]


def from_bool(x: Sequence[bool]) -> list[float]:
    """布尔转浮点：True → 1.0。"""
    return [1.0 if value else 0.0 for value in x]


def from_int32(x: Sequence[int]) -> list[float]:
    """int32 转浮点。"""
    return [float(value) for value in x]


def from_int64(x: Sequence[int]) -> list[float]:
    """int64 转浮点。"""
    return [float(value) for value in x]


def from_float(x: Sequence[float], target: str = "f4") -> list[float]:
    """浮点宽度转换：`target="f4"` 收窄到 f32 精度，`"f8"` 保持。"""
    if target == "f8":
        return [float(value) for value in x]
    return [_narrow_f32(value) for value in x]


def to_bool(x: Sequence[float]) -> list[bool]:
    """浮点转布尔：非零为真，**NaN 也算真**（与 Go 的汇编口径一致）。"""
    return [bool(value) or _math.isnan(value) for value in x]


def to_int32(x: Sequence[float]) -> list[int]:
    """浮点转 int32：向零截断。"""
    return [int(_math.trunc(value)) for value in x]


def to_int64(x: Sequence[float]) -> list[int]:
    """浮点转 int64：向零截断。"""
    return [int(_math.trunc(value)) for value in x]


def to_float32(x: Sequence[float]) -> list[float]:
    """收窄到 f32 精度后回落 f64。"""
    return [_narrow_f32(value) for value in x]


def to_float64(x: Sequence[float]) -> list[float]:
    """加宽到 f64（Python 里即原样）。"""
    return [float(value) for value in x]


def slice_float_to_int(x: Sequence[float], target: str = "i8") -> list[int]:
    """浮点切片转整型切片：向零截断；`target` 只影响调用方对宽度的预期。"""
    if target not in {"i4", "i8"}:
        raise ValueError(f"slice_float_to_int: 不支持的 target {target!r}")
    return [int(_math.trunc(value)) for value in x]


def slice_to_float64(x: Sequence[float | int | bool]) -> list[float]:
    """布尔 / 整型切片转浮点切片。"""
    return [float(value) for value in x]


def ceil(x: Sequence[float]) -> list[float]:
    """逐元素向上取整。"""
    return [_math.ceil(value) for value in x]


def floor(x: Sequence[float]) -> list[float]:
    """逐元素向下取整。"""
    return [_math.floor(value) for value in x]


def round(x: Sequence[float]) -> list[float]:  # noqa: A001 - 与合同名 ROUND 保持一致
    """逐元素四舍五入（半程远离零，与 Go 的 `math.Round` 一致）。"""
    return [_round_half_away(value) for value in x]


def _round_half_away(value: float) -> float:
    if _math.isnan(value) or _math.isinf(value):
        return value
    return _math.floor(value + 0.5) if value >= 0 else _math.ceil(value - 0.5)


def sin(x: Sequence[float]) -> list[float]:
    """逐元素正弦。"""
    return [_math.sin(value) for value in x]


def cos(x: Sequence[float]) -> list[float]:
    """逐元素余弦。"""
    return [_math.cos(value) for value in x]


def exp(x: Sequence[float]) -> list[float]:
    """逐元素指数。"""
    return [_math.exp(value) for value in x]


def log(x: Sequence[float]) -> list[float]:
    """逐元素自然对数（非正数给 -Inf / NaN，与 IEEE 754 一致）。"""
    return [_log_of(value, _math.log) for value in x]


def log2(x: Sequence[float]) -> list[float]:
    """逐元素以 2 为底的对数。"""
    return [_log_of(value, _math.log2) for value in x]


def log10(x: Sequence[float]) -> list[float]:
    """逐元素以 10 为底的对数。"""
    return [_log_of(value, _math.log10) for value in x]


def _log_of(value: float, func) -> float:
    if _math.isnan(value):
        return _math.nan
    if value < 0:
        return _math.nan
    if value == 0:
        return -_math.inf
    return func(value)


def pow_float(x: Sequence[float], y: Sequence[float]) -> list[float]:
    """逐元素浮点幂。"""
    if len(x) != len(y):
        raise ValueError("pow_float: 两条序列长度不一致")
    return [_math.pow(base, exponent) for base, exponent in zip(x, y)]


def pow_integer(x: Sequence[float], n: int) -> list[float]:
    """逐元素整数次幂。"""
    return [_math.pow(value, n) for value in x]


def find(x: Sequence[float], value: float) -> int:
    """返回首次出现的下标；未命中返回 -1。"""
    for index, item in enumerate(x):
        if item == value:
            return index
    return -1


def gather(x: Sequence[float], idx: Sequence[int]) -> list[float]:
    """按索引取值。"""
    return [x[index] for index in idx]


def select(x: Sequence[float], mask: Sequence[bool]) -> list[float]:
    """按掩码过滤。"""
    if len(x) != len(mask):
        raise ValueError("select: 长度不一致")
    return [value for value, keep in zip(x, mask) if keep]


def scatter(x: Sequence[float], idx: Sequence[int], size: int) -> list[float]:
    """把 x 按索引散落到长度为 size 的零序列。"""
    if len(x) != len(idx):
        raise ValueError("scatter: 值与索引长度不一致")
    out = [0.0] * size
    for value, position in zip(x, idx):
        out[position] = value
    return out


def none(x: Sequence[bool]) -> bool:
    """是否全部为假；空输入为 True。"""
    return not builtins.any(bool(value) for value in x)
