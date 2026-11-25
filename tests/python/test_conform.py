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
from num import matrix, range as sequence, shape as shape_module  # noqa: E402
from num.x32 import vek  # noqa: E402


def semantic_of(data: np.lib.npyio.NpzFile) -> str:
    return bytes(data["meta_semantic"]).rstrip(b"\x00").decode("ascii")


def matrix_of(data: np.lib.npyio.NpzFile, name: str):
    array = np.asarray(data[name])
    return array.tolist()


def window_of(data: np.lib.npyio.NpzFile, name: str):
    array = np.asarray(data[name])
    if array.ndim == 0:
        return int(array)
    return [int(value) for value in array.reshape(-1).tolist()]


def optional_int(data: np.lib.npyio.NpzFile, name: str):
    if name not in data.files:
        return None
    return int(np.asarray(data[name]).reshape(-1)[0])


DTYPE_CODES = {"float64": "f8", "float32": "f4", "int64": "i8", "int32": "i4", "bool": "b1"}


def expected_code(data: np.lib.npyio.NpzFile) -> str:
    array = np.asarray(data["expected"] if "expected" in data.files else data["expected_flat"])
    return DTYPE_CODES.get(str(array.dtype), str(array.dtype))


def flatten(value) -> list:
    if isinstance(value, (list, tuple)):
        out: list = []
        for item in value:
            out.extend(flatten(item))
        return out
    return [value]


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
    "where": lambda data: num.where(
        [bool(value) for value in values_of(data, "condition")],
        values_of(data, "x"),
        values_of(data, "y"),
    ),
    "align": lambda data: num.align(
        values_of(data, "x"),
        values_of(data, "value")[0],
        int(np.asarray(data["n"]).reshape(-1)[0]),
    ),
    "repeat": lambda data: num.repeat(
        values_of(data, "value")[0], int(np.asarray(data["n"]).reshape(-1)[0])
    ),
    "zeros": lambda data: num.zeros(int(np.asarray(data["n"]).reshape(-1)[0])),
    "decimal": lambda data: [
        num.decimal(value, int(np.asarray(data["digits"]).reshape(-1)[0]))
        if "digits" in data.files else num.decimal(value)
        for value in values_of(data, "x")
    ],
    "nan_check": lambda data: [num.is_nan(value) for value in values_of(data, "x")],
    "dtype2int": lambda data: num.dtype2int(values_of(data, "x")),
    "bool_to_numeric": lambda data: (
        num.bool_to_float64([bool(value) for value in values_of(data, "x")])
        if np.asarray(data["expected"]).dtype.kind == "f"
        else num.bool_to_int64([bool(value) for value in values_of(data, "x")])
    ),
    "expm1": lambda data: num.expm1(values_of(data, "x")),
    "log1p": lambda data: num.log1p(values_of(data, "x")),
    "neg": lambda data: vek.neg(values_of(data, "x")),
    "inv": lambda data: vek.inv(values_of(data, "x")),
    "prod": lambda data: vek.prod(values_of(data, "x")),
    "cumprod": lambda data: vek.cumprod(values_of(data, "x")),
    "norm": lambda data: vek.norm(values_of(data, "x")),
    "manhattan_norm": lambda data: vek.manhattan_norm(values_of(data, "x")),
    "distance": lambda data: vek.distance(values_of(data, "x"), values_of(data, "y")),
    "manhattan_distance": lambda data: vek.manhattan_distance(
        values_of(data, "x"), values_of(data, "y")
    ),
    "cosine_similarity": lambda data: vek.cosine_similarity(
        values_of(data, "x"), values_of(data, "y")
    ),
    "logic_xor": lambda data: vek.logic_xor(
        [bool(value) for value in values_of(data, "x")],
        [bool(value) for value in values_of(data, "y")],
    ),
    "none": lambda data: vek.none([bool(value) for value in values_of(data, "x")]),
    "quantile": lambda data: vek.quantile(
        values_of(data, "x"), float(np.asarray(data["q"]).reshape(-1)[0])
    ),
    "shape": lambda data: shape_module.shape(matrix_of(data, "x")),
    "concat1d": lambda data: matrix.concat1d(
        values_of(data, "a"), values_of(data, "b"), values_of(data, "c")
    ),
    "transpose2d": lambda data: matrix.transpose2d(matrix_of(data, "x")),
    "dot2d": lambda data: matrix.dot2d(matrix_of(data, "a"), matrix_of(data, "b")),
    "dot2d1": lambda data: matrix.dot2d1(matrix_of(data, "a"), values_of(data, "v")),
    "matmul": lambda data: matrix.matmul(matrix_of(data, "x"), matrix_of(data, "y")),
    "inverse": lambda data: matrix.inverse(matrix_of(data, "a")),
    "mat4mul": lambda data: matrix.mat4mul(matrix_of(data, "x"), matrix_of(data, "y")),
    "arange": lambda data: sequence.arange(
        float(np.asarray(data["start"]).reshape(-1)[0]),
        float(np.asarray(data["end"]).reshape(-1)[0]),
        float(np.asarray(data["step"]).reshape(-1)[0]),
    ),
    "linspace": lambda data: sequence.linspace(
        float(np.asarray(data["start"]).reshape(-1)[0]),
        float(np.asarray(data["end"]).reshape(-1)[0]),
        int(np.asarray(data["n"]).reshape(-1)[0]),
        bool(np.asarray(data["endpoint"]).reshape(-1)[0]),
    ),
    "range": lambda data: sequence.range(int(np.asarray(data["n"]).reshape(-1)[0])),
    "rolling": lambda data: num.rolling(
        values_of(data, "x"), window_of(data, "window"), optional_int(data, "min_periods")
    ),
    "rolling_apply": lambda data: num.rolling_apply(
        values_of(data, "x"), window_of(data, "window"), num.mean
    ),
    "from_bool": lambda data: vek.from_bool([bool(value) for value in values_of(data, "x")]),
    "from_float": lambda data: vek.from_float(values_of(data, "x"), expected_code(data)),
    "from_int32": lambda data: vek.from_int32([int(value) for value in values_of(data, "x")]),
    "from_int64": lambda data: vek.from_int64([int(value) for value in values_of(data, "x")]),
    "to_bool": lambda data: vek.to_bool(values_of(data, "x")),
    "to_float32": lambda data: vek.to_float32(values_of(data, "x")),
    "to_float64": lambda data: vek.to_float64(values_of(data, "x")),
    "to_int32": lambda data: vek.to_int32(values_of(data, "x")),
    "to_int64": lambda data: vek.to_int64(values_of(data, "x")),
    "slice_float_to_int": lambda data: vek.slice_float_to_int(
        values_of(data, "x"), expected_code(data)
    ),
    "slice_to_float64": lambda data: vek.slice_to_float64(values_of(data, "x")),
    "type_default": lambda data: [num.type_default(expected_code(data))],
    "fill": lambda data: num.fill(values_of(data, "x"), values_of(data, "value")[0]),
    "fill_na": lambda data: num.fill_na(values_of(data, "x"), values_of(data, "value")[0]),
    "poly_fit": lambda data: num.poly_fit(
        values_of(data, "x"), values_of(data, "y"), int(np.asarray(data["deg"]).reshape(-1)[0])
    ),
    "poly_val": lambda data: num.poly_val(values_of(data, "p"), values_of(data, "x")),
    "slope": lambda data: [
        num.slope(
            int(np.asarray(data["x1"]).reshape(-1)[0]),
            values_of(data, "y1")[0],
            int(np.asarray(data["x2"]).reshape(-1)[0]),
            values_of(data, "y2")[0],
        )
    ],
    "slope_to_degrees": lambda data: [num.slope_to_degrees(values_of(data, "x")[0])],
    "degrees_to_slope": lambda data: [num.degrees_to_slope(values_of(data, "x")[0])],
    "triangle_bevel": lambda data: [
        num.triangle_bevel(
            values_of(data, "slope")[0],
            int(np.asarray(data["x1"]).reshape(-1)[0]),
            values_of(data, "y1")[0],
            int(np.asarray(data["xn"]).reshape(-1)[0]),
        )
    ],
    "cross": lambda data: num.cross(values_of(data, "a"), values_of(data, "b")),
    "fit_line": lambda data: num.fit_line(values_of(data, "x"), values_of(data, "y")),
    "calculate_line_equation": lambda data: (
        num.calculate_line_equation(
            values_of(data, "x1")[0],
            values_of(data, "y1")[0],
            values_of(data, "x2")[0],
            values_of(data, "y2")[0],
        ),
        True,
    ),
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
                raw = DISPATCH[semantic](data)
                if semantic == "cross":
                    expected_x = [int(value) for value in values_of(data, "expected_x")]
                    expected_state = [int(value) for value in values_of(data, "expected_state")]
                    self.assertEqual([item[0] for item in raw], expected_x, f"{Path(path).name}: 穿点下标")
                    self.assertEqual([item[1] for item in raw], expected_state, f"{Path(path).name}: 穿点方向")
                    covered += 1
                    continue
                if semantic in {"fit_line", "calculate_line_equation"}:
                    line, success = raw
                    expected_a = values_of(data, "expected_a")[0]
                    expected_b = values_of(data, "expected_b")[0]
                    expected_c = values_of(data, "expected_c")[0]
                    a, b, c = line.equation()
                    self.assertTrue(close(a, expected_a, tolerance), f"{Path(path).name}: a {a} != {expected_a}")
                    self.assertTrue(close(b, expected_b, tolerance), f"{Path(path).name}: b {b} != {expected_b}")
                    self.assertTrue(close(c, expected_c, tolerance), f"{Path(path).name}: c {c} != {expected_c}")
                    expected_degrees = values_of(data, "expected_degrees")[0]
                    self.assertTrue(
                        close(line.degrees(), expected_degrees, tolerance),
                        f"{Path(path).name}: degrees {line.degrees()} != {expected_degrees}",
                    )
                    expected_success = values_of(data, "expected_success")[0]
                    self.assertEqual(1 if success else 0, int(expected_success), f"{Path(path).name}: success")
                    probe = values_of(data, "probe_x")
                    expected_y = values_of(data, "expected_y")
                    self.assertEqual([line.y(value) for value in probe], expected_y, f"{Path(path).name}: y")
                    covered += 1
                    continue
                if semantic == "rolling":
                    lengths = [int(value) for value in values_of(data, "expected_lengths")]
                    self.assertEqual(
                        [len(block) for block in raw], lengths,
                        f"{Path(path).name}: 每个位置的窗口长度不对",
                    )
                expected_key = "expected" if "expected" in data.files else "expected_flat"
                expected = np.asarray(data[expected_key]).reshape(-1).tolist()
                actual = flatten(raw)
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
        self.assertGreaterEqual(covered, 138, f"只覆盖了 {covered} 个样本")


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
