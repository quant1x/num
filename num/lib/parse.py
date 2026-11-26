"""字符串解析与判定。"""

from __future__ import annotations

import math
import re

from .convert import _narrow_f32

_TRUE_TEXTS = frozenset({"true", "TRUE", "True", "y", "Y", "yes", "Yes", "YES", "on", "On", "ON",
                         "1", "真", "对", "好"})

_FALSE_TEXTS = frozenset({"false", "FALSE", "False", "n", "N", "no", "No", "NO", "off", "Off", "OFF",
                          "0", "假", "错", "坏"})

def is_empty(text: str) -> bool:
    """去掉首尾空白后是否为空。"""
    return text.strip() == ""

def string_is_nan(text: str) -> bool:
    """去空白、转小写后是否为 `nan`。"""
    return text.strip().lower() == "nan"

def string_is_true(text: str) -> bool:
    """是否属于真值文本（精确匹配）。"""
    return text in _TRUE_TEXTS

def string_is_false(text: str) -> bool:
    """是否属于假值文本（精确匹配）。"""
    return text in _FALSE_TEXTS

_BOOL_TRUE = frozenset({"1", "t", "T", "TRUE", "true", "True"})

_DECIMAL = re.compile(r"^[+-]?[0-9]+$")

_HEX = re.compile(r"^[+-]?[0-9a-fA-F]+$")

_FLOAT = re.compile(r"^[+-]?([0-9]+(\.[0-9]*)?|\.[0-9]+)([eE][+-]?[0-9]+)?$")

INT32_MIN, INT32_MAX = -(2**31), 2**31 - 1

INT32_MIN, INT32_MAX = -(2**31), 2**31 - 1

def parse_bool(text: str) -> bool:
    """按 `strconv.ParseBool` 的口径解析；其余输入为 False。"""
    return text in _BOOL_TRUE

def _parse_int(text: str, minimum: int, maximum: int) -> int:
    if is_empty(text):
        return 0
    if string_is_true(text):
        return 1
    if string_is_false(text) or string_is_nan(text):
        return 0
    for pattern, base in ((_DECIMAL, 10), (_HEX, 16)):
        if not pattern.match(text):
            continue
        value = int(text, base)
        if minimum <= value <= maximum:
            return value
    return 0

def parse_int32(text: str) -> int:
    """解析成 int32；失败为 0。"""
    return _parse_int(text, INT32_MIN, INT32_MAX)

def parse_int64(text: str) -> int:
    """解析成 int64；失败为 0。"""
    return _parse_int(text, -(2**63), 2**63 - 1)

def parse_float64(text: str) -> float:
    """解析成 float64；空串或失败为 NaN。"""
    if is_empty(text):
        return math.nan
    if string_is_true(text):
        return 1.0
    if string_is_false(text):
        return 0.0
    lowered = text.lower()
    if lowered in {"inf", "+inf", "infinity", "+infinity"}:
        return math.inf
    if lowered in {"-inf", "-infinity"}:
        return -math.inf
    if lowered == "nan":
        return math.nan
    if not _FLOAT.match(text):
        return math.nan
    return float(text)

def parse_float32(text: str) -> float:
    """解析成 float32（返回收窄后的值）；空串或失败为 NaN。"""
    return _narrow_f32(parse_float64(text))
