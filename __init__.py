"""quant1x num 的 Python spec 包。

包目录就是仓库根：`.py` 与同名 `.go` / `.h` / `.rs` 平铺在同一目录（CONTRIBUTING §3.2）。
从源码运行时把仓库的**父目录**加入 `sys.path`，再 `import num`；安装后直接 `import num`。

当前只落地首批标量与归约算子，语义以本目录的 `.py` 为准，其余仍以 Go 实现为临时基准。
"""

from .abs import abs
from .argmax import argmax
from .argmin import argmin
from .count import count
from .covariance import covariance, variance
from .cumsum import cumsum
from .diff import diff
from .dot import dot1d
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
    "abs", "argmax", "argmin", "count", "covariance", "cumsum", "diff", "dot1d", "max",
    "maximum", "mean", "median", "min", "minimum", "shift", "sqrt", "std", "sum", "variance",
]
