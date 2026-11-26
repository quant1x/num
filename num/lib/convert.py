"""类型转换：dtype、切片、bool / int / float。"""

from __future__ import annotations

import math
import struct
from collections.abc import Sequence

def _narrow_f32(value: float) -> float:
    return struct.unpack("f", struct.pack("f", value))[0]

def bool_to_int64(x: Sequence[bool]) -> list[int]:
    """返回 1 / 0 序列。"""
    return [1 if value else 0 for value in x]

def bool_to_float64(x: Sequence[bool]) -> list[float]:
    """返回 1.0 / 0.0 序列。"""
    return [1.0 if value else 0.0 for value in x]

def slice_to_bool(x: Sequence[float | int | bool]) -> list[bool]:
    """把数值切片转成布尔切片：非零为真。"""
    return [bool(value) for value in x]

_DEFAULTS = {"f4": math.nan, "f8": math.nan, "i1": 0, "i2": 0, "i4": 0, "i8": 0,
             "u1": 0, "u2": 0, "u4": 0, "u8": 0, "b1": False}

def dtype2int(x: Sequence[float]) -> list[int]:
    """返回向零截断后的整数序列。"""
    return [int(math.trunc(value)) for value in x]

def type_default(dtype: str):
    """返回该 dtype 的默认值。"""
    if dtype not in _DEFAULTS:
        raise ValueError(f"type_default: 未知 dtype {dtype!r}")
    return _DEFAULTS[dtype]

def slice_to_float32(x: Sequence[float | bool]) -> list[float]:
    """把数值切片转成 f32 精度（以 Python float 承载）。"""
    return [_narrow_f32(float(value)) for value in x]

def slice_to_int(x: Sequence[float | bool]) -> list[int]:
    """把数值切片转成整型切片（向零截断）。"""
    return [int(math.trunc(float(value))) for value in x]

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
    return [bool(value) or math.isnan(value) for value in x]

def to_int32(x: Sequence[float]) -> list[int]:
    """浮点转 int32：向零截断。"""
    return [int(math.trunc(value)) for value in x]

def to_int64(x: Sequence[float]) -> list[int]:
    """浮点转 int64：向零截断。"""
    return [int(math.trunc(value)) for value in x]

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
    return [int(math.trunc(value)) for value in x]

def slice_to_float64(x: Sequence[float | int | bool]) -> list[float]:
    """布尔 / 整型切片转浮点切片。"""
    return [float(value) for value in x]

def isnan(value: float) -> bool:
    """NaN 或 ±Inf 都算缺失（与 Go 的 `IsNaN` 一致）。"""
    return math.isnan(value) or math.isinf(value)
