"""字符串判据（spec，对应 Go 的 `type_strings.go` / `type_bool.go`）。

- `is_empty(text)`：去掉空白后为空串；
- `string_is_nan(text)`：去空白、转小写后等于 `nan`；
- `string_is_true` / `string_is_false`：**精确匹配**（Go 不做 trim），覆盖英文、数字与中文写法。

注意：`parse_*` 系列在调这两个判据前会先判空，但不会 trim，所以 `" true"` 不算真值。
"""

from __future__ import annotations

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
