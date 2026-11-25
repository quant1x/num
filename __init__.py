"""quant1x num 的 Python spec 包。

包目录就是仓库根：`.py` 与同名 `.go` / `.h` / `.rs` 平铺在同一目录（CONTRIBUTING §3.2）。
从源码运行时把仓库的**父目录**加入 `sys.path`，再 `import num`；安装后直接 `import num`。

当前只落地首批标量与归约算子，语义以本目录的 `.py` 为准，其余仍以 Go 实现为临时基准。
"""

from .abs import abs
from .argmax import argmax
from .argmin import argmin
from .arithmetics_add import add
from .arithmetics_div import div
from .arithmetics_mul import mul
from .arithmetics_sub import sub
from .count import count
from .covariance import covariance, variance
from .cumsum import cumsum
from .diff import diff
from .dot import dot1d
from .logic import gt, gte, logic_and, logic_not, logic_or, lt, lte
from .logic_all import all
from .logic_any import any
from .logic_equal import equal, not_equal
from .max import max
from .maximum import maximum
from .mean import mean
from .median import median
from .min import min
from .minimum import minimum
from .shift import shift
from .sqrt import sqrt
from .stddev import std
from .sum import sum

__all__ = [
    "abs", "add", "all", "any", "argmax", "argmin", "count", "covariance", "cumsum", "diff",
    "div", "dot1d", "equal", "gt", "gte", "logic_and", "logic_not", "logic_or", "lt", "lte",
    "max", "maximum", "mean", "median", "min", "minimum", "mul", "not_equal", "shift", "sqrt",
    "std", "sub", "sum", "variance",
]
