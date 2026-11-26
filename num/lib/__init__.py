"""NumPy 风格的 lib 工具命名空间。"""

from .array import align
from .convert import (
    _DEFAULTS,
    bool_to_float64,
    bool_to_int64,
    dtype2int,
    from_bool,
    from_float,
    from_int32,
    from_int64,
    isnan,
    slice_float_to_int,
    slice_to_bool,
    slice_to_float32,
    slice_to_float64,
    slice_to_int,
    to_bool,
    to_float32,
    to_float64,
    to_int32,
    to_int64,
    type_default,
)
from .fill import fill, fill_na
from .format import InvalidNumb, Numb, decode, encode, from_f64
from .index import compress, find, scatter, take
from .parse import (
    is_empty,
    parse_bool,
    parse_float32,
    parse_float64,
    parse_int32,
    parse_int64,
    string_is_false,
    string_is_nan,
    string_is_true,
)

__all__ = [
    "align", "bool_to_float64", "bool_to_int64", "compress", "dtype2int", "encode", "decode",
    "fill", "fill_na", "find", "from_bool", "from_f64", "from_float", "from_int32", "from_int64",
    "InvalidNumb", "is_empty", "isnan", "Numb", "parse_bool", "parse_float32", "parse_float64",
    "parse_int32", "parse_int64", "scatter", "slice_float_to_int", "slice_to_bool",
    "slice_to_float32", "slice_to_float64", "slice_to_int", "string_is_false", "string_is_nan",
    "string_is_true", "take", "to_bool", "to_float32", "to_float64", "to_int32", "to_int64",
    "type_default",
]
