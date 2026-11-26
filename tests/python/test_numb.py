"""`.numb` 读写的移植检验：读 6 个黄金样本、逐值核对，再重新编码要求逐字节一致。"""

from __future__ import annotations

import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT.parent) not in sys.path:
    sys.path.insert(0, str(ROOT.parent))

from num import numb  # noqa: E402

GOLDEN = ROOT / "tests" / "fixtures" / "numb"


class NumbPort(unittest.TestCase):
    def test_golden_values(self) -> None:
        vector = numb.decode((GOLDEN / "f8-vector.numb").read_bytes())
        self.assertEqual(vector.dtype, "f8")
        self.assertEqual(vector.shape, (5,))
        self.assertEqual(vector.floats(), [1.5, -2.0, 3.25, 0.0, 1e-3])

        matrix = numb.decode((GOLDEN / "f4-matrix.numb").read_bytes())
        self.assertEqual(matrix.shape, (2, 3))
        self.assertEqual(matrix.floats32(), [1.0, 2.0, 3.0, 4.0, 5.0, 6.0])

        integers = numb.decode((GOLDEN / "i8-vector.numb").read_bytes())
        self.assertEqual(integers.int64s(), [-1, 0, 7, 123456789])

        flags = numb.decode((GOLDEN / "b1-vector.numb").read_bytes())
        self.assertEqual(flags.bools(), [True, False, True])

        scalar = numb.decode((GOLDEN / "f8-scalar.numb").read_bytes())
        self.assertEqual(scalar.shape, ())
        self.assertEqual(scalar.floats(), [42.5])

        empty = numb.decode((GOLDEN / "f8-empty.numb").read_bytes())
        self.assertEqual(empty.shape, (0,))
        self.assertEqual(empty.payload, b"")

    def test_reencode_is_byte_identical(self) -> None:
        for path in sorted(GOLDEN.glob("*.numb")):
            original = path.read_bytes()
            self.assertEqual(numb.encode(numb.decode(original)), original, path.name)

    def test_rejects_malformed(self) -> None:
        valid = numb.encode(numb.from_f64((2,), [1.0, 2.0]))
        with self.assertRaises(numb.InvalidNumb):
            numb.decode(b"NUM")
        with self.assertRaises(numb.InvalidNumb):
            numb.decode(b"NUMX" + valid[4:])
        with self.assertRaises(numb.InvalidNumb):
            numb.decode(valid[:4] + b"\x01\x00" + valid[6:])
        with self.assertRaises(numb.InvalidNumb):
            numb.decode(valid[:-1])
        with self.assertRaises(numb.InvalidNumb):
            numb.decode(valid + b"\x00")
        self.assertEqual(len(numb.decode(valid).floats()), 2)


if __name__ == "__main__":
    unittest.main()
