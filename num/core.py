"""NumPy 顶层命名空间：逐元素、归约、创建、比较、逻辑。"""

from __future__ import annotations

import builtins
import math
from collections.abc import Sequence

Number = int | float

def _pairwise(x: Sequence[float], y: Number | Sequence[float]) -> list[tuple[float, float]]:
    if isinstance(y, (int, float)) and not isinstance(y, bool):
        return [(value, float(y)) for value in x]
    if len(x) != len(y):
        raise ValueError("两条序列长度不一致")
    return list(zip(x, y))

def _divide(numerator: float, denominator: float) -> float:
    if math.isnan(numerator):
        return math.nan
    if denominator == 0:
        if numerator == 0:
            return math.nan
        return math.copysign(math.inf, numerator) * math.copysign(1.0, denominator)
    return numerator / denominator

def abs(x: Sequence[float]) -> list[float]:  # noqa: A001 - 与合同名 ABS 保持一致
    """返回 x 的逐元素绝对值。"""
    return [builtins.abs(value) for value in x]

def add(x: Sequence[float], y: Number | Sequence[float]) -> list[float]:
    """返回逐元素和。"""
    return [left + right for left, right in _pairwise(x, y)]

def subtract(x: Sequence[float], y: float | Sequence[float]) -> list[float]:
    """返回逐元素差。"""
    return [left - right for left, right in _pairwise(x, y)]

def multiply(x: Sequence[float], y: float | Sequence[float]) -> list[float]:
    """返回逐元素积。"""
    return [left * right for left, right in _pairwise(x, y)]

def divide(x: Sequence[float], y: float | Sequence[float]) -> list[float]:
    """返回逐元素商。"""
    return [_divide(left, right) for left, right in _pairwise(x, y)]

def maximum(x: Sequence[float], y: Sequence[float]) -> list[float]:
    """返回逐元素较大者。"""
    if len(x) != len(y):
        raise ValueError("maximum: 两条序列长度不一致")
    return [right if right > left else left for left, right in zip(x, y)]

def minimum(x: Sequence[float], y: Sequence[float]) -> list[float]:
    """返回逐元素较小者。"""
    if len(x) != len(y):
        raise ValueError("minimum: 两条序列长度不一致")
    return [right if right < left else left for left, right in zip(x, y)]

def sqrt(x: Sequence[float]) -> list[float]:
    """返回 x 的逐元素平方根。"""
    out: list[float] = []
    for value in x:
        if math.isnan(value) or value < 0:
            out.append(math.nan)
        else:
            out.append(math.sqrt(value))
    return out

def expm1(x: Sequence[float]) -> list[float]:
    """返回逐元素 `e^x - 1`。"""
    return [math.expm1(value) for value in x]

def log1p(x: Sequence[float]) -> list[float]:
    """返回逐元素 `ln(1 + x)`。"""
    out: list[float] = []
    for value in x:
        if value == -1:
            out.append(-math.inf)
        elif value < -1:
            out.append(math.nan)
        else:
            out.append(math.log1p(value))
    return out

def greater(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x > y`。"""
    return [left > right for left, right in _pairwise(x, y)]

def greater_equal(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x >= y`。"""
    return [left >= right for left, right in _pairwise(x, y)]

def less(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x < y`。"""
    return [left < right for left, right in _pairwise(x, y)]

def less_equal(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素 `x <= y`。"""
    return [left <= right for left, right in _pairwise(x, y)]

def logical_and(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素逻辑与（两侧都非零为真）。"""
    return [bool(left) and bool(right) for left, right in _pairwise(x, y)]

def logical_or(x: Sequence[float], y: float | Sequence[float]) -> list[bool]:
    """逐元素逻辑或（任一侧非零为真）。"""
    return [bool(left) or bool(right) for left, right in _pairwise(x, y)]

def logical_not(x: Sequence[float]) -> list[bool]:
    """逐元素逻辑非。"""
    return [not bool(value) for value in x]

def equal(x: Sequence[float], y: Sequence[float]) -> list[bool]:
    """返回逐元素是否相等。"""
    if len(x) != len(y):
        raise ValueError("equal: 两条序列长度不一致")
    return [left == right for left, right in zip(x, y)]

def not_equal(x: Sequence[float], y: Sequence[float]) -> list[bool]:
    """返回逐元素是否不等。"""
    if len(x) != len(y):
        raise ValueError("not_equal: 两条序列长度不一致")
    return [left != right for left, right in zip(x, y)]

def where(condition: Sequence[bool], x: Sequence[float], y: Sequence[float]) -> list[float]:
    """按条件从 x / y 中逐位置取值。"""
    if len(condition) != len(x) or len(x) != len(y):
        raise ValueError("where: 条件与取值长度不一致")
    return [left if flag else right for flag, left, right in zip(condition, x, y)]

def sum(x: Sequence[float]) -> float:
    """返回 x 的累加和。"""
    total = 0.0
    for value in x:
        total += value
    return total

def mean(x: Sequence[float]) -> float:
    """返回 x 的算术平均。"""
    if len(x) == 0:
        return math.nan
    return sum(x) / len(x)

def var(x: Sequence[float]) -> float:
    """返回 x 的总体方差。"""
    if len(x) == 0:
        return 0.0
    average = mean(x)
    total = 0.0
    for value in x:
        delta = value - average
        total += delta * delta
    return total / len(x)

def covariance(x: Sequence[float], y: Sequence[float]) -> float:
    """返回 x 与 y 的总体协方差。"""
    if len(x) != len(y) or len(x) == 0:
        return 0.0
    mean_x = mean(x)
    mean_y = mean(y)
    total = 0.0
    for left, right in zip(x, y):
        total += (left - mean_x) * (right - mean_y)
    return total / len(x)

def std(x: Sequence[float]) -> float:
    """返回 x 的总体标准差。"""
    if len(x) == 0:
        return 0.0
    average = mean(x)
    total = 0.0
    for value in x:
        delta = value - average
        total += delta * delta
    return math.sqrt(total / len(x))

def min(x: Sequence[float]) -> float:  # noqa: A001 - 与合同名 MIN 保持一致
    """返回 x 的最小值。"""
    if len(x) == 0:
        return 0.0
    best = x[0]
    for value in x[1:]:
        if value < best:
            best = value
    return best

def max(x: Sequence[float]) -> float:  # noqa: A001 - 与合同名 MAX 保持一致
    """返回 x 的最大值。"""
    if len(x) == 0:
        return 0.0
    best = x[0]
    for value in x[1:]:
        if value > best:
            best = value
    return best

def argmin(x: Sequence[float]) -> int:
    """返回 x 中首个最小值的下标。"""
    if len(x) == 0:
        return -1
    best = 0
    for index in builtins.range(1, len(x)):
        if x[index] < x[best]:
            best = index
    return best

def argmax(x: Sequence[float]) -> int:
    """返回 x 中首个最大值的下标。"""
    if len(x) == 0:
        return -1
    best = 0
    for index in builtins.range(1, len(x)):
        if x[index] > x[best]:
            best = index
    return best

def median(x: Sequence[float]) -> float:
    """返回按入参位置取到的中位数。"""
    if len(x) == 0:
        return math.nan
    if len(x) == 1:
        return 0.0
    if len(x) % 2 == 0:
        middle = len(x) // 2
        return (x[middle - 1] + x[middle]) / 2
    return x[len(x) // 2]

def count(x: Sequence[float] | Sequence[bool]) -> int:
    """返回 x 中非零元素的个数。"""
    total = 0
    for value in x:
        if value:
            total += 1
    return total

def all(x: Sequence[float]) -> bool:  # noqa: A001 - 与合同名 ALL 保持一致
    """返回是否所有元素都为真。"""
    return builtins.all(bool(value) for value in x)

def any(x: Sequence[float]) -> bool:  # noqa: A001 - 与合同名 ANY 保持一致
    """返回是否存在为真的元素。"""
    return builtins.any(bool(value) for value in x)

def cumsum(x: Sequence[float]) -> list[float]:
    """返回 x 的前缀和。"""
    out: list[float] = []
    total = 0.0
    for value in x:
        total += value
        out.append(total)
    return out

def ones(n: int) -> list[float]:
    """返回 n 个 1.0。"""
    if n < 1:
        return []
    return [1.0] * n

def zeros(n: int) -> list[float]:
    """返回 n 个 0.0。"""
    if n < 1:
        return []
    return [0.0] * n

def arange(start: float, end: float, step: float = 1.0) -> list[float]:
    """返回 [start, end) 上的等距序列。"""
    out: list[float] = []
    value = start
    while value < end:
        out.append(value)
        value += step
    return out

def linspace(start: float, end: float, n: int, endpoint: bool = True) -> list[float]:
    """返回 n 个等距点。"""
    if n <= 0:
        return []
    if n == 1:
        return [start]
    segments = n - 1 if endpoint else n
    step = (end - start) / segments
    return [start + index * step for index in builtins.range(n)]

def range_float(a: float, b: float) -> list[float]:
    """返回 [a, b) 上的浮点序列（步长 1）。"""
    return arange(a, b, 1.0)

def range(n: int) -> list[int]:  # noqa: A001 - 与合同名 RANGE 保持一致
    """返回 0 .. n-1。"""
    if n <= 0:
        return []
    return list(builtins.range(n))

def repeat(value: float, n: int) -> list[float]:
    """返回 n 个 value。"""
    if n < 1:
        return []
    return [value] * n

def shape(x) -> tuple[int, int]:
    """返回 (行, 列)；一维的第二维为 0。"""
    if len(x) == 0:
        return (0, 0)
    first = x[0]
    if isinstance(first, (list, tuple)):
        return (len(x), len(first))
    return (len(x), 0)

def diff(x: Sequence[float], n: int = 1) -> list[float]:
    """返回方向为 `x[i-n] - x[i]` 的差分序列。"""
    shifted = shift(x, n)
    return [previous - value for previous, value in zip(shifted, x)]

def shift(x: Sequence[float], n: int) -> list[float]:
    """返回平移 n 位后的序列。"""
    out = [math.nan] * len(x)
    for index in builtins.range(len(x)):
        source = index - n
        if 0 <= source < len(x):
            out[index] = x[source]
    return out

def negative(x: Sequence[float]) -> list[float]:
    """逐元素取相反数。"""
    return [-value for value in x]

def reciprocal(x: Sequence[float]) -> list[float]:
    """逐元素取倒数（按 IEEE 754 处理除零）。"""
    return [_divide(1.0, value) for value in x]

def ceil(x: Sequence[float]) -> list[float]:
    """逐元素向上取整。"""
    return [math.ceil(value) for value in x]

def floor(x: Sequence[float]) -> list[float]:
    """逐元素向下取整。"""
    return [math.floor(value) for value in x]

def _round_half_away(value: float) -> float:
    if math.isnan(value) or math.isinf(value):
        return value
    return math.floor(value + 0.5) if value >= 0 else math.ceil(value - 0.5)

def round(x: Sequence[float]) -> list[float]:  # noqa: A001 - 与合同名 ROUND 保持一致
    """逐元素四舍五入（半程远离零，与 Go 的 `math.Round` 一致）。"""
    return [_round_half_away(value) for value in x]

def sin(x: Sequence[float]) -> list[float]:
    """逐元素正弦。"""
    return [math.sin(value) for value in x]

def cos(x: Sequence[float]) -> list[float]:
    """逐元素余弦。"""
    return [math.cos(value) for value in x]

def exp(x: Sequence[float]) -> list[float]:
    """逐元素指数。"""
    return [math.exp(value) for value in x]

def _log_of(value: float, func) -> float:
    if math.isnan(value):
        return math.nan
    if value < 0:
        return math.nan
    if value == 0:
        return -math.inf
    return func(value)

def log(x: Sequence[float]) -> list[float]:
    """逐元素自然对数（非正数给 -Inf / NaN，与 IEEE 754 一致）。"""
    return [_log_of(value, math.log) for value in x]

def log2(x: Sequence[float]) -> list[float]:
    """逐元素以 2 为底的对数。"""
    return [_log_of(value, math.log2) for value in x]

def log10(x: Sequence[float]) -> list[float]:
    """逐元素以 10 为底的对数。"""
    return [_log_of(value, math.log10) for value in x]

def logical_xor(x: Sequence[bool], y: Sequence[bool]) -> list[bool]:
    """逐元素逻辑异或。"""
    if len(x) != len(y):
        raise ValueError("logic_xor: 两条序列长度不一致")
    return [bool(left) != bool(right) for left, right in zip(x, y)]

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

def none(x: Sequence[bool]) -> bool:
    """是否全部为假；空输入为 True。"""
    return not builtins.any(bool(value) for value in x)

def pow_float(x: Sequence[float], y: Sequence[float]) -> list[float]:
    """逐元素浮点幂。"""
    if len(x) != len(y):
        raise ValueError("pow_float: 两条序列长度不一致")
    return [math.pow(base, exponent) for base, exponent in zip(x, y)]

def pow_integer(x: Sequence[float], n: int) -> list[float]:
    """逐元素整数次幂。"""
    return [math.pow(value, n) for value in x]


def power(x: Sequence[float], y: int | float | Sequence[float]) -> list[float]:
    """逐元素幂：整数指数走 pow_integer，其余走 pow_float。"""
    if isinstance(y, int) and not isinstance(y, bool):
        return pow_integer(x, y)
    return pow_float(x, y)
