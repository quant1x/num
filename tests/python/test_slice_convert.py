"""切片转换的移植检验（期望值取自 Go 的 SliceToBool / SliceToInts / SliceToFloat32）。"""

from __future__ import annotations

import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))

from num.lib import convert  # noqa: E402


class SliceConvert(unittest.TestCase):
    def test_float_source(self) -> None:
        source = [1.9, -2.9, 0.0, 0.1]
        self.assertEqual(convert.slice_to_bool(source), [True, True, False, True])
        self.assertEqual(convert.slice_to_int(source), [1, -2, 0, 0])
        narrowed = convert.slice_to_float32(source)
        self.assertEqual(narrowed[0], convert.slice_to_float32([1.9])[0])
        self.assertEqual(narrowed[3], 0.10000000149011612)

    def test_int_and_bool_sources(self) -> None:
        self.assertEqual(convert.slice_to_bool([-2, 0, 3]), [True, False, True])
        self.assertEqual(convert.slice_to_int([-2, 0, 3]), [-2, 0, 3])
        self.assertEqual(convert.slice_to_float32([-2, 0, 3]), [-2.0, 0.0, 3.0])
        self.assertEqual(convert.slice_to_int([True, False]), [1, 0])
        self.assertEqual(convert.slice_to_float32([True, False]), [1.0, 0.0])


if __name__ == "__main__":
    unittest.main()
