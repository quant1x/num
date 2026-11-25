"""字符串解析（spec，对应 Go 的 `Parse*` 系列）。

Go 侧 `IgnoreParseExceptions = true`，解析失败不抛异常，而是回落到「空值」：

| 输入 | bool | int32 / int64 | float32 / float64 |
| --- | --- | --- | --- |
| 空串 / 全空白 | False | 0 | **NaN** |
| 真值文本（`true` / `yes` / `1` / `真` …） | True | 1 | 1.0 |
| 假值文本（`false` / `no` / `0` / `假` …） | False | 0 | 0.0 |
| `nan` | False | 0 | NaN |
| 合法数字 | False | 十进制，其次十六进制 | 十进制浮点 |
| 其它 | False | 0 | **NaN** |

细节（都由 Go 的 `strconv` 行为决定，Python 内置函数更宽松，必须收紧）：

- **不 trim**：`" 7 "` 解析失败（整型 0、浮点 NaN）；
- 整型十六进制不认 `0x` 前缀（`"0x1f"` 失败，`"1f"` = 31）；
- `parse_bool` 只认 `strconv.ParseBool` 的写法（`1` / `t` / `true` 及其大写变体、`0` / `f` / `false`），
  不认 `yes` / `on` / 中文 —— 与 `string_is_true` 的范围不同。
"""

from __future__ import annotations

import math
import re
import struct

from .type_strings import is_empty, string_is_false, string_is_nan, string_is_true

_BOOL_TRUE = frozenset({"1", "t", "T", "TRUE", "true", "True"})
_DECIMAL = re.compile(r"^[+-]?[0-9]+$")
_HEX = re.compile(r"^[+-]?[0-9a-fA-F]+$")
_FLOAT = re.compile(r"^[+-]?([0-9]+(\.[0-9]*)?|\.[0-9]+)([eE][+-]?[0-9]+)?$")
INT32_MIN, INT32_MAX = -(2**31), 2**31 - 1


def _narrow_f32(value: float) -> float:
    return struct.unpack("f", struct.pack("f", value))[0]


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
