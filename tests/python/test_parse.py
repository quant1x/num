"""字符串判据与解析的移植检验。

期望值取自 Go 侧同一组输入（`num.ParseBool` / `ParseInt32` / `ParseInt64` /
`ParseFloat32` / `ParseFloat64` 与 `IsEmpty` / `StringIsTrue` / `StringIsFalse` / `StringIsNaN`）。
"""

from __future__ import annotations

import math
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))

from num.lib import parse  # noqa: E402

# (输入, bool, int32, int64, float32, float64)
CASES = [
    ("true", True, 1, 1, 1.0, 1.0),
    ("TRUE", True, 1, 1, 1.0, 1.0),
    ("1", True, 1, 1, 1.0, 1.0),
    ("t", True, 0, 0, math.nan, math.nan),
    ("yes", False, 1, 1, 1.0, 1.0),
    ("on", False, 1, 1, 1.0, 1.0),
    ("真", False, 1, 1, 1.0, 1.0),
    ("false", False, 0, 0, 0.0, 0.0),
    ("0", False, 0, 0, 0.0, 0.0),
    ("f", False, 15, 15, math.nan, math.nan),
    ("no", False, 0, 0, 0.0, 0.0),
    ("", False, 0, 0, math.nan, math.nan),
    (" ", False, 0, 0, math.nan, math.nan),
    ("nan", False, 0, 0, math.nan, math.nan),
    ("NaN", False, 0, 0, math.nan, math.nan),
    ("1.5", False, 0, 0, 1.5, 1.5),
    ("-2", False, -2, -2, -2.0, -2.0),
    ("0x1f", False, 0, 0, math.nan, math.nan),
    (" 7 ", False, 0, 0, math.nan, math.nan),
    ("abc", False, 2748, 2748, math.nan, math.nan),
    ("1f", False, 31, 31, math.nan, math.nan),
    ("-1f", False, -31, -31, math.nan, math.nan),
    ("1e3", False, 483, 483, 1000.0, 1000.0),
    ("inf", False, 0, 0, math.inf, math.inf),
]


class ParsePort(unittest.TestCase):
    def test_parse_table(self) -> None:
        for text, expected_bool, expected_i32, expected_i64, expected_f32, expected_f64 in CASES:
            label = repr(text)
            self.assertEqual(parse.parse_bool(text), expected_bool, f"{label}: bool")
            self.assertEqual(parse.parse_int32(text), expected_i32, f"{label}: int32")
            self.assertEqual(parse.parse_int64(text), expected_i64, f"{label}: int64")
            for got, want, name in (
                (parse.parse_float32(text), expected_f32, "float32"),
                (parse.parse_float64(text), expected_f64, "float64"),
            ):
                if isinstance(want, float) and math.isnan(want):
                    self.assertTrue(math.isnan(got), f"{label}: {name} 应为 NaN，实际 {got}")
                elif math.isinf(want):
                    self.assertTrue(math.isinf(got) and got == want, f"{label}: {name} 应为 {want}")
                else:
                    self.assertEqual(got, want, f"{label}: {name}")

    def test_string_predicates(self) -> None:
        self.assertTrue(parse.is_empty("   "))
        self.assertFalse(parse.is_empty(" 7 "))
        self.assertTrue(parse.string_is_nan(" NaN "))
        self.assertFalse(parse.string_is_nan("nanx"))
        self.assertTrue(parse.string_is_true("yes"))
        self.assertTrue(parse.string_is_true("真"))
        self.assertFalse(parse.string_is_true(" true"))  # 精确匹配，不 trim
        self.assertTrue(parse.string_is_false("off"))
        self.assertTrue(parse.string_is_false("假"))


if __name__ == "__main__":
    unittest.main()
