"""quant1x num 的 Python spec：按 NumPy 公开命名空间组织。

顶层是 ufunc / 归约 / 创建 / 逻辑；`num.linalg`、`num.polynomial`、`num.lib` 对应 NumPy 的子命名空间，
`num.ta`、`num.ztable`、`num.freq` 是 num 专有能力。`import num` 不加载 `ztable` 的 428 KB 数据表。
"""

from .core import (
    abs,
    add,
    all,
    any,
    arange,
    argmax,
    argmin,
    ceil,
    cos,
    count,
    covariance,
    cumsum,
    cumprod,
    diff,
    divide,
    equal,
    exp,
    expm1,
    floor,
    greater,
    greater_equal,
    less,
    less_equal,
    linspace,
    log,
    log10,
    log1p,
    log2,
    logical_and,
    logical_not,
    logical_or,
    logical_xor,
    max,
    maximum,
    mean,
    median,
    min,
    minimum,
    multiply,
    negative,
    none,
    not_equal,
    ones,
    power,
    prod,
    quantile,
    range_float,
    reciprocal,
    repeat,
    round,
    shape,
    shift,
    sin,
    sqrt,
    std,
    subtract,
    sum,
    var,
    where,
    zeros,
)
from .lib import (
    align,
    bool_to_float64,
    compress,
    take,
    bool_to_int64,
    dtype2int,
    fill,
    fill_na,
    isnan,
    type_default,
)
from .linalg import dot, inv as _matrix_inv, matmul, norm, transpose
from .polynomial import polyfit, polyval
from .ta import (
    BREAK_THROUGH,
    FALL_DRASTICALLY,
    Line,
    calculate_line_equation,
    change_rate,
    cross,
    decimal,
    degrees_to_slope,
    fit_line,
    net_change_rate,
    rolling,
    rolling_apply,
    slope,
    slope_to_degrees,
    triangle_bevel,
)
from .freq import date_range, parse_freq

__all__ = [
    "abs", "add", "all", "any", "arange", "argmax", "argmin", "ceil", "cos", "count",
    "cumsum", "cumprod", "covariance", "diff", "divide", "dot", "equal", "exp", "expm1", "floor",
    "greater", "greater_equal", "less", "less_equal", "linspace", "log", "log10", "log1p",
    "log2", "logical_and", "logical_not", "logical_or", "logical_xor", "matmul", "max",
    "maximum", "mean", "median", "min", "minimum", "multiply", "negative", "norm",
    "not_equal", "ones", "polyfit", "polyval", "power", "prod", "quantile", "range_float", "reciprocal",
    "repeat", "round", "shape", "shift", "sin", "sqrt", "std", "subtract", "sum",
    "transpose", "var", "where", "zeros",
    "align", "bool_to_float64", "bool_to_int64", "dtype2int", "fill", "fill_na", "isnan",
    "type_default", "date_range", "parse_freq",
    "BREAK_THROUGH", "FALL_DRASTICALLY", "Line", "calculate_line_equation", "change_rate",
    "cross", "decimal", "degrees_to_slope", "fit_line", "net_change_rate", "rolling",
    "rolling_apply", "slope", "slope_to_degrees", "triangle_bevel",
]

# 一版兼容别名：旧名不进 __all__，旧模块路径不保留。
variance = var
gt = greater
gte = greater_equal
lt = less
lte = less_equal
logic_and = logical_and
logic_or = logical_or
logic_not = logical_not
logic_xor = logical_xor
neg = negative
inv = reciprocal
inverse = _matrix_inv
sub = subtract
mul = multiply
div = divide
dot1d = dot
transpose2d = transpose
poly_fit = polyfit
poly_val = polyval
gather = take
select = compress
is_nan = isnan
