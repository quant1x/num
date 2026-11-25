"""quant1x num 的 Python spec 包。

包目录就是仓库根：`.py` 与同名 `.go` / `.h` / `.rs` 平铺在同一目录（CONTRIBUTING §3.2）。
从源码运行时把仓库的**父目录**加入 `sys.path`，再 `import num`；安装后直接 `import num`。

当前只落地首批标量与归约算子，语义以本目录的 `.py` 为准，其余仍以 Go 实现为临时基准。
"""

from .abs import abs
from .align import align
from .argmax import argmax
from .argmin import argmin
from .arithmetics_add import add
from .arithmetics_div import div
from .arithmetics_mul import mul
from .arithmetics_sub import sub
from .count import count
from .cross import cross
from .covariance import covariance, variance
from .cumsum import cumsum
from .decimals import decimal, is_nan
from .diff import diff
from .fillna import fill, fill_na
from .dot import dot1d
from .line import Line, calculate_line_equation, fit_line
from .logic import gt, gte, logic_and, logic_not, logic_or, lt, lte
from .logic_all import all
from .logic_any import any
from .log import expm1, log1p
from .logic_equal import equal, not_equal
from .max import max
from .maximum import maximum
from .mean import mean
from .median import median
from .min import min
from .minimum import minimum
from .repeat import repeat
from .rolling import rolling, rolling_apply
from .polynomial import poly_fit, poly_val
from .shift import shift
from .slope import degrees_to_slope, slope, slope_to_degrees, triangle_bevel
from .sqrt import sqrt
from .stddev import std
from .sum import sum
from .type_bool import bool_to_float64, bool_to_int64
from .type_dtype import dtype2int, type_default
from .where import where
from .zeros import zeros

__all__ = [
    "abs", "add", "align", "all", "any", "argmax", "argmin", "bool_to_float64", "bool_to_int64",
    "count", "covariance", "cross", "cumsum", "decimal", "degrees_to_slope", "diff", "div",
    "dot1d", "dtype2int", "equal",
    "expm1", "fill", "fill_na", "gt", "gte", "is_nan", "log1p", "logic_and", "logic_not", "logic_or", "lt", "lte", "Line", "calculate_line_equation", "fit_line",
    "max", "maximum", "mean", "median", "min", "minimum", "mul", "not_equal", "poly_fit", "poly_val", "repeat", "rolling",
    "rolling_apply", "shift", "slope", "slope_to_degrees",
    "sqrt", "std", "sub", "sum", "triangle_bevel", "type_default", "variance", "where", "zeros",
]
