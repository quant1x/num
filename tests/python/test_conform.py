"""读取 `tests/fixtures` 的 npz，用 Python spec 重算并逐值对拍。

运行（仓库根）：

    python -m unittest discover -s tests/python -v

包布局见 `docs/proposals/0002-python-package-layout.md`：仓库根就是包，因此这里把仓库的
**父目录**放进 `sys.path` 之后再 `import num`。
"""

from __future__ import annotations

import glob
import math
import sys
import unittest
from pathlib import Path

import numpy as np

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT.parent) not in sys.path:
    sys.path.insert(0, str(ROOT.parent))

import num  # noqa: E402  (路径注入后再导入)


def semantic_of(data: np.lib.npyio.NpzFile) -> str:
    return bytes(data["meta_semantic"]).rstrip(b"\x00").decode("ascii")


def values_of(data: np.lib.npyio.NpzFile, name: str) -> list[float]:
    return np.asarray(data[name]).reshape(-1).tolist()


def close(actual: float, expected: float, tolerance: float) -> bool:
    if math.isnan(actual) and math.isnan(expected):
        return True
    if actual == expected:
        return True
    return abs(actual - expected) <= max(tolerance, 1e-12) * max(1.0, abs(expected))


DISPATCH = {
    "sum": lambda data: num.sum(values_of(data, "x")),
    "mean": lambda data: num.mean(values_of(data, "x")),
    "abs": lambda data: num.abs(values_of(data, "x")),
    "cumsum": lambda data: num.cumsum(values_of(data, "x")),
    "std": lambda data: num.std(values_of(data, "x")),
    "variance": lambda data: num.variance(values_of(data, "x")),
    "dot1d": lambda data: num.dot1d(values_of(data, "x"), values_of(data, "y")),
    "covariance": lambda data: num.covariance(values_of(data, "x"), values_of(data, "y")),
    "sqrt": lambda data: num.sqrt(values_of(data, "x")),
    "maximum": lambda data: num.maximum(values_of(data, "x"), values_of(data, "y")),
    "minimum": lambda data: num.minimum(values_of(data, "x"), values_of(data, "y")),
    "max": lambda data: num.max(values_of(data, "x")),
    "min": lambda data: num.min(values_of(data, "x")),
    "argmax": lambda data: num.argmax(values_of(data, "x")),
    "argmin": lambda data: num.argmin(values_of(data, "x")),
    "count": lambda data: num.count(values_of(data, "x")),
    "median": lambda data: num.median(values_of(data, "x")),
    "shift": lambda data: num.shift(values_of(data, "x"), int(np.asarray(data["n"]).reshape(-1)[0])),
    "diff": lambda data: num.diff(values_of(data, "x"), int(np.asarray(data["n"]).reshape(-1)[0])),
    "add": lambda data: num.add(values_of(data, "x"), values_of(data, "y")),
    "sub": lambda data: num.sub(values_of(data, "x"), values_of(data, "y")),
    "mul": lambda data: num.mul(values_of(data, "x"), values_of(data, "y")),
    "div": lambda data: num.div(values_of(data, "x"), values_of(data, "y")),
    "equal": lambda data: num.equal(values_of(data, "x"), values_of(data, "y")),
    "not_equal": lambda data: num.not_equal(values_of(data, "x"), values_of(data, "y")),
    "greater": lambda data: num.gt(values_of(data, "x"), values_of(data, "y")),
    "greater_equal": lambda data: num.gte(values_of(data, "x"), values_of(data, "y")),
    "less": lambda data: num.lt(values_of(data, "x"), values_of(data, "y")),
    "less_equal": lambda data: num.lte(values_of(data, "x"), values_of(data, "y")),
    "logic_and": lambda data: num.logic_and(values_of(data, "x"), values_of(data, "y")),
    "logic_or": lambda data: num.logic_or(values_of(data, "x"), values_of(data, "y")),
    "logic_not": lambda data: num.logic_not(values_of(data, "x")),
    "all": lambda data: num.all(values_of(data, "x")),
    "any": lambda data: num.any(values_of(data, "x")),
}


class FixtureConformance(unittest.TestCase):
    def test_fixtures(self) -> None:
        covered = 0
        for path in sorted(glob.glob(str(ROOT / "tests/fixtures/**/*.npz"), recursive=True)):
            with np.load(path) as data:
                semantic = semantic_of(data)
                if semantic not in DISPATCH:
                    continue
                tolerance = float(np.asarray(data["meta_tolerance"]))
                expected = np.asarray(data["expected"]).reshape(-1).tolist()
                actual = DISPATCH[semantic](data)
                if not isinstance(actual, list):
                    actual = [actual]
                self.assertEqual(
                    len(actual), len(expected),
                    f"{Path(path).name}: 长度 {len(actual)} != {len(expected)}",
                )
                for index, (got, want) in enumerate(zip(actual, expected)):
                    self.assertTrue(
                        close(float(got), float(want), tolerance),
                        f"{Path(path).name} {semantic}: 第 {index} 位 {got} != {want}",
                    )
                covered += 1
        self.assertGreaterEqual(covered, 57, f"只覆盖了 {covered} 个样本")


class EmptyInputs(unittest.TestCase):
    def test_reductions(self) -> None:
        self.assertEqual(num.sum([]), 0.0)
        self.assertTrue(math.isnan(num.mean([])))
        self.assertEqual(num.std([]), 0.0)
        self.assertEqual(num.variance([]), 0.0)
        self.assertEqual(num.covariance([], []), 0.0)

    def test_elementwise(self) -> None:
        self.assertEqual(num.abs([]), [])
        self.assertEqual(num.cumsum([]), [])
        self.assertEqual(num.dot1d([], []), 0.0)

    def test_length_mismatch(self) -> None:
        with self.assertRaises(ValueError):
            num.dot1d([1.0, 2.0], [1.0])
        with self.assertRaises(ValueError):
            num.maximum([1.0, 2.0], [1.0])
        with self.assertRaises(ValueError):
            num.minimum([1.0, 2.0], [1.0])


if __name__ == "__main__":
    unittest.main()
