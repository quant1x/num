"""读取 `tests/fixtures` 的 npz，用 Python spec 重算并逐值对拍。

运行（仓库根）：

    python -m unittest discover -s tests/python -v

包布局见 `docs/proposals/0003-python-numpy-package-layout.md`：仓库根直接 `import num`。
"""

from __future__ import annotations

import glob
import math
import sys
import unittest
from pathlib import Path

import numpy as np

ROOT = Path(__file__).resolve().parents[2]
if str(ROOT) not in sys.path:
    sys.path.insert(0, str(ROOT))

import num  # noqa: E402  (路径注入后再导入)
from num import core, linalg  # noqa: E402
from num.lib import convert, index  # noqa: E402


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


def strings_of(data: np.lib.npyio.NpzFile, name: str) -> list[str]:
    array = np.asarray(data[name])
    return [item.decode("ascii") if isinstance(item, bytes) else str(item)
            for item in array.reshape(-1).tolist()]


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
    "variance": lambda data: num.var(values_of(data, "x")),
    "dot1d": lambda data: num.dot(values_of(data, "x"), values_of(data, "y")),
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
    "sub": lambda data: num.subtract(values_of(data, "x"), values_of(data, "y")),
    "mul": lambda data: num.multiply(values_of(data, "x"), values_of(data, "y")),
    "div": lambda data: num.divide(values_of(data, "x"), values_of(data, "y")),
    "equal": lambda data: num.equal(values_of(data, "x"), values_of(data, "y")),
    "not_equal": lambda data: num.not_equal(values_of(data, "x"), values_of(data, "y")),
    "greater": lambda data: num.greater(values_of(data, "x"), values_of(data, "y")),
    "greater_equal": lambda data: num.greater_equal(values_of(data, "x"), values_of(data, "y")),
    "less": lambda data: num.less(values_of(data, "x"), values_of(data, "y")),
    "less_equal": lambda data: num.less_equal(values_of(data, "x"), values_of(data, "y")),
    "logic_and": lambda data: num.logical_and(values_of(data, "x"), values_of(data, "y")),
    "logic_or": lambda data: num.logical_or(values_of(data, "x"), values_of(data, "y")),
    "logic_not": lambda data: num.logical_not(values_of(data, "x")),
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
    "nan_check": lambda data: [num.isnan(value) for value in values_of(data, "x")],
    "dtype2int": lambda data: num.dtype2int(values_of(data, "x")),
    "bool_to_numeric": lambda data: (
        num.bool_to_float64([bool(value) for value in values_of(data, "x")])
        if np.asarray(data["expected"]).dtype.kind == "f"
        else num.bool_to_int64([bool(value) for value in values_of(data, "x")])
    ),
    "expm1": lambda data: num.expm1(values_of(data, "x")),
    "log1p": lambda data: num.log1p(values_of(data, "x")),
    "neg": lambda data: num.negative(values_of(data, "x")),
    "inv": lambda data: num.reciprocal(values_of(data, "x")),
    "prod": lambda data: num.prod(values_of(data, "x")),
    "cumprod": lambda data: num.cumprod(values_of(data, "x")),
    "norm": lambda data: num.norm(values_of(data, "x")),
    "manhattan_norm": lambda data: linalg.manhattan_norm(values_of(data, "x")),
    "distance": lambda data: linalg.distance(values_of(data, "x"), values_of(data, "y")),
    "manhattan_distance": lambda data: linalg.manhattan_distance(
        values_of(data, "x"), values_of(data, "y")
    ),
    "cosine_similarity": lambda data: linalg.cosine_similarity(
        values_of(data, "x"), values_of(data, "y")
    ),
    "logic_xor": lambda data: num.logical_xor(
        [bool(value) for value in values_of(data, "x")],
        [bool(value) for value in values_of(data, "y")],
    ),
    "none": lambda data: num.none([bool(value) for value in values_of(data, "x")]),
    "quantile": lambda data: num.quantile(
        values_of(data, "x"), float(np.asarray(data["q"]).reshape(-1)[0])
    ),
    "shape": lambda data: num.shape(matrix_of(data, "x")),
    "concat1d": lambda data: linalg.concat1d(
        values_of(data, "a"), values_of(data, "b"), values_of(data, "c")
    ),
    "transpose2d": lambda data: linalg.transpose(matrix_of(data, "x")),
    "dot2d": lambda data: linalg.dot2d(matrix_of(data, "a"), matrix_of(data, "b")),
    "dot2d1": lambda data: linalg.dot2d1(matrix_of(data, "a"), values_of(data, "v")),
    "matmul": lambda data: linalg.matmul(matrix_of(data, "x"), matrix_of(data, "y")),
    "inverse": lambda data: linalg.inv(matrix_of(data, "a")),
    "mat4mul": lambda data: linalg.mat4mul(matrix_of(data, "x"), matrix_of(data, "y")),
    "arange": lambda data: num.arange(
        float(np.asarray(data["start"]).reshape(-1)[0]),
        float(np.asarray(data["end"]).reshape(-1)[0]),
        float(np.asarray(data["step"]).reshape(-1)[0]),
    ),
    "linspace": lambda data: num.linspace(
        float(np.asarray(data["start"]).reshape(-1)[0]),
        float(np.asarray(data["end"]).reshape(-1)[0]),
        int(np.asarray(data["n"]).reshape(-1)[0]),
        bool(np.asarray(data["endpoint"]).reshape(-1)[0]),
    ),
    "range": lambda data: core.range(int(np.asarray(data["n"]).reshape(-1)[0])),
    "rolling": lambda data: num.rolling(
        values_of(data, "x"), window_of(data, "window"), optional_int(data, "min_periods")
    ),
    "rolling_apply": lambda data: num.rolling_apply(
        values_of(data, "x"), window_of(data, "window"), num.mean
    ),
    "from_bool": lambda data: convert.from_bool([bool(value) for value in values_of(data, "x")]),
    "from_float": lambda data: convert.from_float(values_of(data, "x"), expected_code(data)),
    "from_int32": lambda data: convert.from_int32([int(value) for value in values_of(data, "x")]),
    "from_int64": lambda data: convert.from_int64([int(value) for value in values_of(data, "x")]),
    "to_bool": lambda data: convert.to_bool(values_of(data, "x")),
    "to_float32": lambda data: convert.to_float32(values_of(data, "x")),
    "to_float64": lambda data: convert.to_float64(values_of(data, "x")),
    "to_int32": lambda data: convert.to_int32(values_of(data, "x")),
    "to_int64": lambda data: convert.to_int64(values_of(data, "x")),
    "slice_float_to_int": lambda data: convert.slice_float_to_int(
        values_of(data, "x"), expected_code(data)
    ),
    "slice_to_float64": lambda data: convert.slice_to_float64(values_of(data, "x")),
    "type_default": lambda data: [num.type_default(expected_code(data))],
    "fill": lambda data: num.fill(values_of(data, "x"), values_of(data, "value")[0]),
    "fill_na": lambda data: num.fill_na(values_of(data, "x"), values_of(data, "value")[0]),
    "poly_fit": lambda data: num.polyfit(
        values_of(data, "x"), values_of(data, "y"), int(np.asarray(data["deg"]).reshape(-1)[0])
    ),
    "poly_val": lambda data: num.polyval(values_of(data, "p"), values_of(data, "x")),
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
    "ones": lambda data: num.ones(int(np.asarray(data["n"]).reshape(-1)[0])),
    "range_float": lambda data: num.range_float(
        values_of(data, "a")[0], values_of(data, "b")[0]
    ),
    "ceil": lambda data: num.ceil(values_of(data, "x")),
    "floor": lambda data: num.floor(values_of(data, "x")),
    "round": lambda data: num.round(values_of(data, "x")),
    "sin": lambda data: convert.to_float32(num.sin(values_of(data, "x"))),
    "cos": lambda data: convert.to_float32(num.cos(values_of(data, "x"))),
    "exp": lambda data: convert.to_float32(num.exp(values_of(data, "x"))),
    "log": lambda data: convert.to_float32(num.log(values_of(data, "x"))),
    "log2": lambda data: convert.to_float32(num.log2(values_of(data, "x"))),
    "log10": lambda data: convert.to_float32(num.log10(values_of(data, "x"))),
    "pow_float": lambda data: num.power(values_of(data, "x"), values_of(data, "y")),
    "pow_integer": lambda data: num.power(values_of(data, "x"), 2),
    "find": lambda data: [index.find(values_of(data, "x"), values_of(data, "value")[0])],
    "gather": lambda data: index.take(
        values_of(data, "x"), [int(value) for value in values_of(data, "idx")]
    ),
    "select": lambda data: index.compress(
        values_of(data, "x"), [bool(value) for value in values_of(data, "mask")]
    ),
    "scatter": lambda data: index.scatter(
        values_of(data, "x"),
        [int(value) for value in values_of(data, "idx")],
        int(np.asarray(data["size"]).reshape(-1)[0]),
    ),
    "change_rate": lambda data: [
        num.change_rate(values_of(data, "base")[0], values_of(data, "current")[0])
    ],
    "net_change_rate": lambda data: [
        num.net_change_rate(values_of(data, "base")[0], values_of(data, "current")[0])
    ],
    "parse_freq": lambda data: [num.parse_freq(item) for item in strings_of(data, "freq")],
    "date_range": lambda data: num.date_range(
        int(np.asarray(data["start"]).reshape(-1)[0]),
        int(np.asarray(data["periods"]).reshape(-1)[0]),
        strings_of(data, "freq")[0],
    ),
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
        self.assertGreaterEqual(covered, 160, f"只覆盖了 {covered} 个样本")


class EmptyInputs(unittest.TestCase):
    def test_reductions(self) -> None:
        self.assertEqual(num.sum([]), 0.0)
        self.assertTrue(math.isnan(num.mean([])))
        self.assertEqual(num.std([]), 0.0)
        self.assertEqual(num.var([]), 0.0)
        self.assertEqual(num.covariance([], []), 0.0)

    def test_elementwise(self) -> None:
        self.assertEqual(num.abs([]), [])
        self.assertEqual(num.cumsum([]), [])
        self.assertEqual(num.dot([], []), 0.0)

    def test_length_mismatch(self) -> None:
        with self.assertRaises(ValueError):
            num.dot([1.0, 2.0], [1.0])
        with self.assertRaises(ValueError):
            num.maximum([1.0, 2.0], [1.0])
        with self.assertRaises(ValueError):
            num.minimum([1.0, 2.0], [1.0])


if __name__ == "__main__":
    unittest.main()
