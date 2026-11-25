#!/usr/bin/env python3
"""生成跨语言 fixture（C2）。

输入与期望输出由 NumPy 产出，落盘为 ``tests/fixtures/<分组>/<用例>.npz``。npz 内只放数据数组，
外加两个元数据标量：

``meta_semantic``
    语义名（``|S32`` 标量），对应 ``tests/contract/NAME-MAP.tsv`` 的 ``semantic`` 列。
``meta_tolerance``
    浮点比较容差（``f8`` 标量）；0 表示逐位相等。
``meta_checksum``
    所有数据数组按名字排序后串联的 FNV-1a 64 校验和（``u8`` 标量），供各语言读取器自检。

以 ``meta_`` 开头的条目是元数据，其余条目按名字对应算子的输入与期望输出（期望输出统一叫 ``expected``）。

zip 里的时间戳与条目顺序固定，重复执行产物逐字节一致。
"""

from __future__ import annotations

import argparse
import io
import math
import os
import zipfile

import numpy as np
import numpy.lib.format as npy_format

FIXED_ZIP_TIME = (1980, 1, 1, 0, 0, 0)
FNV_OFFSET_BASIS = 0xCBF29CE484222325
FNV_PRIME = 0x100000001B3
MASK64 = (1 << 64) - 1


def fnv1a64(payload: bytes) -> int:
    digest = FNV_OFFSET_BASIS
    for byte in payload:
        digest ^= byte
        digest = (digest * FNV_PRIME) & MASK64
    return digest


def write_npz(path: str, arrays: dict[str, np.ndarray]) -> None:
    payload = b"".join(np.ascontiguousarray(arrays[name]).tobytes()
                       for name in sorted(k for k in arrays if not k.startswith("meta_")))
    checksum = fnv1a64(payload)

    entries = dict(arrays)
    entries["meta_checksum"] = np.array(checksum, dtype=np.uint64)

    with zipfile.ZipFile(path, "w", compression=zipfile.ZIP_STORED) as archive:
        for name in sorted(entries):
            buffer = io.BytesIO()
            npy_format.write_array(buffer, np.asarray(entries[name]), allow_pickle=False)
            info = zipfile.ZipInfo(f"{name}.npy", date_time=FIXED_ZIP_TIME)
            info.compress_type = zipfile.ZIP_STORED
            info.external_attr = 0o600 << 16
            archive.writestr(info, buffer.getvalue())


def save(out_dir: str, group: str, case: str, semantic: str, arrays: dict[str, np.ndarray],
         tolerance: float = 0.0) -> None:
    target_dir = os.path.join(out_dir, group)
    os.makedirs(target_dir, exist_ok=True)
    path = os.path.join(target_dir, f"{case}.npz")
    payload = dict(arrays)
    payload["meta_semantic"] = np.array(semantic.encode("ascii"), dtype="S32")
    payload["meta_tolerance"] = np.array(tolerance, dtype=np.float64)
    write_npz(path, payload)


def f64(values) -> np.ndarray:
    return np.asarray(values, dtype=np.float64)


def f32(values) -> np.ndarray:
    return np.asarray(values, dtype=np.float32)


def elementwise(out_dir: str) -> None:
    group = "elementwise"
    x = f64([1.5, -2.0, 0.0, 3.25])
    y = f64([0.5, 4.0, -1.0, 2.0])
    save(out_dir, group, "add-f64", "add", {"x": x, "y": y, "expected": np.add(x, y)})
    save(out_dir, group, "sub-f64", "sub", {"x": x, "y": y, "expected": np.subtract(x, y)})
    save(out_dir, group, "mul-f64", "mul", {"x": x, "y": y, "expected": np.multiply(x, y)})
    save(out_dir, group, "div-f64", "div", {"x": x, "y": y, "expected": np.divide(x, y)})

    x32, y32 = f32([1.5, -2.0, 0.25, 3.25]), f32([0.5, 4.0, -1.0, 2.0])
    save(out_dir, group, "add-f32", "add", {"x": x32, "y": y32, "expected": np.add(x32, y32)})

    with np.errstate(divide="ignore", invalid="ignore"):
        zero = f64([0.0, -0.0, 1.0, -1.0])
        expected = np.divide(f64([1.0, -1.0, 1.0, 1.0]), zero)
    save(out_dir, group, "div-zero-f64", "div", {"x": f64([1.0, -1.0, 1.0, 1.0]), "y": zero,
                                                "expected": expected})

    maximum = f64([1.0, -5.0, 3.5, 0.0])
    minimum = f64([2.0, -1.0, 3.5, -0.5])
    save(out_dir, group, "maximum-f64", "maximum",
         {"x": maximum, "y": minimum, "expected": np.maximum(maximum, minimum)})
    save(out_dir, group, "minimum-f64", "minimum",
         {"x": maximum, "y": minimum, "expected": np.minimum(maximum, minimum)})

    save(out_dir, group, "abs-f64", "abs",
         {"x": f64([-2.5, -0.0, 0.0, 7.25]), "expected": f64([2.5, 0.0, 0.0, 7.25])})
    save(out_dir, group, "sqrt-f64", "sqrt",
         {"x": f64([0.0, 1.0, 2.0, 4.0]), "expected": np.sqrt(f64([0.0, 1.0, 2.0, 4.0]))})
    save(out_dir, group, "pow-integer-f64", "pow_integer",
         {"x": f64([-2.0, -1.0, 0.0, 1.0, 2.0, 3.0]),
          "expected": np.power(f64([-2.0, -1.0, 0.0, 1.0, 2.0, 3.0]), 2)})

    left = f64([1.0, 2.0, float("nan"), 4.0])
    right = f64([1.0, 3.0, float("nan"), 4.0])
    save(out_dir, group, "equal-f64", "equal", {"x": left, "y": right, "expected": np.equal(left, right)})
    save(out_dir, group, "not-equal-f64", "not_equal",
         {"x": left, "y": right, "expected": np.not_equal(left, right)})
    save(out_dir, group, "greater-f64", "greater",
         {"x": f64([1.0, 2.0, 3.0]), "y": f64([2.0, 2.0, 2.0]),
          "expected": np.greater(f64([1.0, 2.0, 3.0]), f64([2.0, 2.0, 2.0]))})

    a = np.array([True, True, False, False])
    b = np.array([True, False, True, False])
    save(out_dir, group, "logic-and", "logic_and", {"x": a, "y": b, "expected": np.logical_and(a, b)})
    save(out_dir, group, "logic-or", "logic_or", {"x": a, "y": b, "expected": np.logical_or(a, b)})
    save(out_dir, group, "logic-not", "logic_not", {"x": a, "expected": np.logical_not(a)})

    condition = np.array([True, False, True, False])
    save(out_dir, group, "where-f64", "where",
         {"condition": condition, "x": f64([1.0, 2.0, 3.0, 4.0]), "y": f64([10.0, 20.0, 30.0, 40.0]),
          "expected": np.where(condition, f64([1.0, 2.0, 3.0, 4.0]), f64([10.0, 20.0, 30.0, 40.0]))})

    save(out_dir, group, "add-empty-f64", "add",
         {"x": f64([]), "y": f64([]), "expected": f64([])})


def reduction(out_dir: str) -> None:
    group = "reduction"
    x = f64([1.0, 2.0, 3.0, 4.0])
    save(out_dir, group, "sum-basic-f64", "sum", {"x": x, "expected": np.array(np.sum(x))})
    save(out_dir, group, "mean-basic-f64", "mean", {"x": x, "expected": np.array(np.mean(x))})
    save(out_dir, group, "std-basic-f64", "std",
         {"x": x, "expected": np.array(np.std(x))}, tolerance=1e-12)
    save(out_dir, group, "var-basic-f64", "variance",
         {"x": x, "expected": np.array(np.var(x))}, tolerance=1e-12)
    save(out_dir, group, "min-basic-f64", "min", {"x": x, "expected": np.array(np.min(x))})
    save(out_dir, group, "max-basic-f64", "max", {"x": x, "expected": np.array(np.max(x))})
    save(out_dir, group, "argmin-basic-f64", "argmin",
         {"x": f64([3.0, 1.0, 4.0, 2.0]), "expected": np.array(np.argmin(f64([3.0, 1.0, 4.0, 2.0])),
                                                              dtype=np.int64)})
    save(out_dir, group, "argmax-basic-f64", "argmax",
         {"x": f64([3.0, 1.0, 4.0, 2.0]), "expected": np.array(np.argmax(f64([3.0, 1.0, 4.0, 2.0])),
                                                              dtype=np.int64)})

    outliers = f64([0.0, 1.5, -2.0, 0.0, 3.5])
    save(out_dir, group, "count-numeric-f64", "count",
         {"x": outliers, "expected": np.array(np.count_nonzero(outliers), dtype=np.int64)})
    flags = np.array([True, False, True, True])
    save(out_dir, group, "count-bool", "count",
         {"x": flags, "expected": np.array(np.count_nonzero(flags), dtype=np.int64)})

    save(out_dir, group, "cumsum-basic-f64", "cumsum", {"x": x, "expected": np.cumsum(x)})
    save(out_dir, group, "prod-basic-f64", "prod", {"x": f64([1.5, 2.0, -3.0]),
                                                    "expected": np.array(np.prod(f64([1.5, 2.0, -3.0])))})
    a, b = f64([1.0, 2.0, 3.0]), f64([4.0, 5.0, 6.0])
    save(out_dir, group, "dot1d-basic-f64", "dot1d", {"x": a, "y": b, "expected": np.array(np.dot(a, b))})

    # 向量对归约：欧氏距离 / 曼哈顿距离 / 余弦相似度（Go 侧 x32 / x64 的 Distance 族）。
    save(out_dir, group, "distance-basic-f64", "distance",
         {"x": a, "y": b, "expected": np.array(np.sqrt(np.sum((a - b) ** 2)))}, tolerance=1e-12)
    save(out_dir, group, "manhattan_distance-basic-f64", "manhattan_distance",
         {"x": a, "y": b, "expected": np.array(np.sum(np.abs(a - b)))}, tolerance=1e-12)
    save(out_dir, group, "cosine_similarity-basic-f64", "cosine_similarity",
         {"x": a, "y": b,
          "expected": np.array(np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b)))}, tolerance=1e-12)

    with np.errstate(invalid="ignore"):
        nan_sum = np.sum(f64([1.0, float("nan"), 2.0]))
    save(out_dir, group, "sum-nan-f64", "sum",
         {"x": f64([1.0, float("nan"), 2.0]), "expected": np.array(nan_sum)}, tolerance=1e-12)

    long_sum = np.array([0.1] * 100, dtype=np.float64)
    save(out_dir, group, "sum-long-f64", "sum", {"x": long_sum, "expected": np.array(np.sum(long_sum))},
         tolerance=1e-12)
    long_sum32 = np.array([0.1] * 100, dtype=np.float32)
    save(out_dir, group, "sum-long-f32", "sum",
         {"x": long_sum32, "expected": np.array(np.sum(long_sum32, dtype=np.float32))}, tolerance=1e-6)

    # Median 按入参顺序取中间元素、不排序（NAME-MAP 的 median 行），期望值按同一口径手算。
    unsorted = f64([5.0, 1.0, 3.0])
    save(out_dir, group, "median-unsorted-f64", "median",
         {"x": unsorted, "expected": np.array(unsorted[len(unsorted) // 2])})


def scalar_layer(out_dir: str) -> None:
    group = "scalar"
    save(out_dir, group, "type-default-f64", "type_default",
         {"expected": np.array(np.nan, dtype=np.float64)})
    save(out_dir, group, "type-default-f32", "type_default",
         {"expected": np.array(np.nan, dtype=np.float32)})
    save(out_dir, group, "type-default-i64", "type_default",
         {"expected": np.array(0, dtype=np.int64)})
    save(out_dir, group, "type-default-bool", "type_default",
         {"expected": np.array(False)})

    probe = f64([np.nan, np.inf, -np.inf, -0.0, 1.0])
    save(out_dir, group, "nan-check-f64", "nan_check", {"x": probe, "expected": np.isnan(probe) | np.isinf(probe)})
    probe32 = f32([np.nan, np.inf, -np.inf, -0.0, 1.0])
    save(out_dir, group, "nan-check-f32", "nan_check",
         {"x": probe32, "expected": np.isnan(probe32) | np.isinf(probe32)})

    flags = np.array([True, False, True])
    save(out_dir, group, "bool-to-float64", "bool_to_numeric",
         {"x": flags, "expected": np.array([1.0, 0.0, 1.0], dtype=np.float64)})
    save(out_dir, group, "bool-to-int64", "bool_to_numeric",
         {"x": flags, "expected": np.array([1, 0, 1], dtype=np.int64)})

    ints = np.array([1, -2, 3], dtype=np.int64)
    save(out_dir, group, "slice-to-float64-int64", "slice_to_float64",
         {"x": ints, "expected": ints.astype(np.float64)})
    save(out_dir, group, "slice-to-float64-bool", "slice_to_float64",
         {"x": flags, "expected": flags.astype(np.float64)})

    trunc = f64([1.9, -2.9, 0.0, 7.0])
    save(out_dir, group, "slice-float-to-int32", "slice_float_to_int",
         {"x": trunc, "expected": np.trunc(trunc).astype(np.int32)})
    save(out_dir, group, "slice-float-to-int64", "slice_float_to_int",
         {"x": trunc, "expected": np.trunc(trunc).astype(np.int64)})


def decimal_expected(value: float, digits: int = 2) -> float:
    """按契约口径复算 Decimal：NaN 归零、半数远离零、负 digits 视为 0。"""
    if math.isnan(value):
        value = 0.0
    digits = max(0, digits)
    half = -5.0 if math.copysign(1.0, value) < 0 else 5.0
    scale = 10.0 ** digits
    next_scale = 10.0 ** (digits + 1)
    return math.trunc((value * next_scale + half) / 10.0) / scale


def elementwise_extras(out_dir: str) -> None:
    group = "elementwise"
    x32a, y32 = f32([1.5, -2.0, 0.25, 3.25]), f32([0.5, 4.0, -1.0, 2.0])
    save(out_dir, group, "sub-f32", "sub", {"x": x32a, "y": y32, "expected": np.subtract(x32a, y32)})
    save(out_dir, group, "mul-f32", "mul", {"x": x32a, "y": y32, "expected": np.multiply(x32a, y32)})
    save(out_dir, group, "div-f32", "div", {"x": x32a, "y": y32, "expected": np.divide(x32a, y32)})
    save(out_dir, group, "maximum-f32", "maximum",
         {"x": x32a, "y": y32, "expected": np.maximum(x32a, y32)})
    save(out_dir, group, "minimum-f32", "minimum",
         {"x": x32a, "y": y32, "expected": np.minimum(x32a, y32)})
    save(out_dir, group, "abs-f32", "abs",
         {"x": f32([-2.5, -0.0, 0.0, 7.25]), "expected": np.abs(f32([-2.5, -0.0, 0.0, 7.25]))})
    save(out_dir, group, "sqrt-f32", "sqrt",
         {"x": f32([0.0, 1.0, 2.0, 4.0]), "expected": np.sqrt(f32([0.0, 1.0, 2.0, 4.0]))},
         tolerance=1e-6)

    values = f64([1.234, -1.235, 2.5, -2.5])
    save(out_dir, group, "decimal-default-f64", "decimal",
         {"x": values, "expected": f64([decimal_expected(v, 2) for v in values])})
    save(out_dir, group, "decimal-digits0-f64", "decimal",
         {"x": values, "digits": np.array(0, dtype=np.int64),
          "expected": f64([decimal_expected(v, 0) for v in values])})

    save(out_dir, group, "change-rate-f64", "change_rate",
         {"base": np.array(2.0), "current": np.array(3.0), "expected": np.array(1.5)})
    save(out_dir, group, "net-change-rate-f64", "net_change_rate",
         {"base": np.array(2.0), "current": np.array(3.0), "expected": np.array(50.0)})

    px, py = f64([2.0, 3.0, 4.0]), f64([0.5, 2.0, 3.0])
    save(out_dir, group, "pow-float-f64", "pow_float",
         {"x": px, "y": py, "expected": np.power(px, py)})

    ex = f64([-0.5, 0.0, 0.5, 1.0])
    save(out_dir, group, "expm1-f64", "expm1", {"x": ex, "expected": np.expm1(ex)}, tolerance=1e-12)
    lg = f64([-0.5, 0.0, 0.5, 1.0])
    save(out_dir, group, "log1p-f64", "log1p", {"x": lg, "expected": np.log1p(lg)}, tolerance=1e-12)

    rounding = f64([1.4, 1.6, -1.4, -1.6, 2.25])
    save(out_dir, group, "ceil-f64", "ceil", {"x": rounding, "expected": np.ceil(rounding)})
    save(out_dir, group, "floor-f64", "floor", {"x": rounding, "expected": np.floor(rounding)})
    save(out_dir, group, "round-f64", "round",
         {"x": rounding, "expected": f64([math.floor(v + 0.5) if v >= 0 else math.ceil(v - 0.5) for v in rounding])})

    save(out_dir, group, "neg-f64", "neg", {"x": x32a.astype(np.float64), "expected": -x32a.astype(np.float64)})
    with np.errstate(divide="ignore"):
        inv_x = f64([2.0, -4.0, 0.5, 8.0])
        save(out_dir, group, "inv-f64", "inv", {"x": inv_x, "expected": np.divide(1.0, inv_x)})

    trig = f32([0.0, 0.5, 1.0, -1.0])
    save(out_dir, group, "sin-f32", "sin", {"x": trig, "expected": np.sin(trig)}, tolerance=1e-6)
    save(out_dir, group, "cos-f32", "cos", {"x": trig, "expected": np.cos(trig)}, tolerance=1e-6)
    pos = f32([0.5, 1.0, 2.0, 4.0])
    save(out_dir, group, "exp-f32", "exp", {"x": pos, "expected": np.exp(pos)}, tolerance=1e-5)
    save(out_dir, group, "log-f32", "log", {"x": pos, "expected": np.log(pos)}, tolerance=1e-6)
    save(out_dir, group, "log2-f32", "log2", {"x": pos, "expected": np.log2(pos)}, tolerance=1e-6)
    save(out_dir, group, "log10-f32", "log10", {"x": pos, "expected": np.log10(pos)}, tolerance=1e-6)

    left, right = f64([1.0, 2.0, 3.0]), f64([2.0, 2.0, 2.0])
    save(out_dir, group, "gte-f64", "greater_equal", {"x": left, "y": right, "expected": np.greater_equal(left, right)})
    save(out_dir, group, "lt-f64", "less", {"x": left, "y": right, "expected": np.less(left, right)})
    save(out_dir, group, "lte-f64", "less_equal", {"x": left, "y": right, "expected": np.less_equal(left, right)})

    a, b = np.array([True, True, False, False]), np.array([True, False, True, False])
    save(out_dir, group, "logic-xor", "logic_xor", {"x": a, "y": b, "expected": np.logical_xor(a, b)})

    mask = np.array([True, False, True, False])
    save(out_dir, group, "select-f64", "select",
         {"x": f64([1.0, 2.0, 3.0, 4.0]), "mask": mask,
          "expected": f64([1.0, 3.0])})

    haystack = f64([10.0, 20.0, 30.0])
    save(out_dir, group, "find-hit-f64", "find",
         {"x": haystack, "value": np.array(20.0), "expected": np.array(1, dtype=np.int64)})
    save(out_dir, group, "find-miss-f64", "find",
         {"x": haystack, "value": np.array(99.0), "expected": np.array(-1, dtype=np.int64)})

    source = f64([10.0, 20.0, 30.0, 40.0])
    idx = np.array([2, 0, 3], dtype=np.int64)
    save(out_dir, group, "gather-f64", "gather", {"x": source, "idx": idx, "expected": np.take(source, idx)})
    scatter_expected = np.zeros(6, dtype=np.float64)
    scatter_expected[[1, 3, 4]] = source[:3]
    save(out_dir, group, "scatter-f64", "scatter",
         {"x": f64([10.0, 20.0, 30.0]), "idx": np.array([1, 3, 4], dtype=np.int64),
          "size": np.array(6, dtype=np.int64), "expected": scatter_expected})


def reduction_extras(out_dir: str) -> None:
    group = "reduction"
    x32a = f32([1.0, 2.0, 3.0, 4.0])
    save(out_dir, group, "mean-f32", "mean", {"x": x32a, "expected": np.array(np.mean(x32a, dtype=np.float32))},
         tolerance=1e-6)
    save(out_dir, group, "std-f32", "std", {"x": x32a, "expected": np.array(np.std(x32a, dtype=np.float32))},
         tolerance=1e-5)
    save(out_dir, group, "min-f32", "min", {"x": x32a, "expected": np.array(np.min(x32a))})
    save(out_dir, group, "max-f32", "max", {"x": x32a, "expected": np.array(np.max(x32a))})
    save(out_dir, group, "prod-f32", "prod", {"x": f32([1.5, 2.0, -3.0]),
                                              "expected": np.array(np.prod(f32([1.5, 2.0, -3.0])))})

    even = f64([4.0, 1.0, 3.0, 2.0])
    save(out_dir, group, "median-even-f64", "median",
         {"x": even, "expected": np.array((even[1] + even[2]) / 2)})

    flags = np.array([True, True, False])
    save(out_dir, group, "all-bool", "all", {"x": flags, "expected": np.array(np.all(flags))})
    save(out_dir, group, "any-bool", "any", {"x": flags, "expected": np.array(np.any(flags))})
    save(out_dir, group, "none-bool", "none", {"x": flags, "expected": np.array(not np.any(flags))})

    save(out_dir, group, "norm-f64", "norm", {"x": f64([3.0, 4.0]), "expected": np.array(5.0)})
    save(out_dir, group, "manhattan-norm-f64", "manhattan_norm",
         {"x": f64([-3.0, 4.0, -5.0]), "expected": np.array(12.0)})
    quant = f64([1.0, 2.0, 3.0, 4.0])
    save(out_dir, group, "quantile-q025-f64", "quantile",
         {"x": quant, "q": np.array(0.25), "expected": np.array(np.quantile(quant, 0.25))})

    cx, cy = f64([1.0, 2.0, 3.0]), f64([2.0, 4.0, 7.0])
    save(out_dir, group, "covariance-f64", "covariance",
         {"x": cx, "y": cy, "expected": np.array(np.cov(cx, cy, ddof=0)[0, 1])}, tolerance=1e-12)

    save(out_dir, group, "sum-inf-f64", "sum",
         {"x": f64([1.0, np.inf, 2.0]), "expected": np.array(np.inf)})


def container(out_dir: str) -> None:
    group = "container"
    zeros = np.zeros(4, dtype=np.float64)
    save(out_dir, group, "zeros-f64", "zeros", {"n": np.array(4, dtype=np.int64), "expected": zeros})
    ones = np.ones(5, dtype=np.float64)
    save(out_dir, group, "ones-f64", "ones", {"n": np.array(5, dtype=np.int64), "expected": ones})

    save(out_dir, group, "arange-int", "arange",
         {"start": np.array(2, dtype=np.int64), "end": np.array(12, dtype=np.int64),
          "step": np.array(3, dtype=np.int64), "expected": np.arange(2, 12, 3, dtype=np.int64)})
    save(out_dir, group, "arange-step-f64", "arange",
         {"start": np.array(0.0), "end": np.array(1.0), "step": np.array(0.25),
          "expected": np.arange(0.0, 1.0, 0.25, dtype=np.float64)})

    save(out_dir, group, "linspace-f64", "linspace",
         {"start": np.array(0.0), "end": np.array(1.0), "n": np.array(5, dtype=np.int64),
          "endpoint": np.array(True), "expected": np.linspace(0.0, 1.0, 5)})
    save(out_dir, group, "linspace-endpoint-false", "linspace",
         {"start": np.array(0.0), "end": np.array(1.0), "n": np.array(5, dtype=np.int64),
          "endpoint": np.array(False), "expected": np.linspace(0.0, 1.0, 5, endpoint=False)})

    save(out_dir, group, "range-int", "range",
         {"n": np.array(5, dtype=np.int64), "expected": np.arange(5, dtype=np.int64)})
    save(out_dir, group, "range-f64", "range",
         {"n": np.array(4, dtype=np.int64), "expected": np.arange(4, dtype=np.float64)})
    save(out_dir, group, "range-float-f32", "range_float",
         {"a": np.array(1.0, dtype=np.float32), "b": np.array(4.0, dtype=np.float32),
          "expected": np.arange(1.0, 4.0, dtype=np.float32)})

    save(out_dir, group, "repeat-f64", "repeat",
         {"value": np.array(2.5), "n": np.array(4, dtype=np.int64),
          "expected": np.full(4, 2.5, dtype=np.float64)})
    save(out_dir, group, "repeat-int", "repeat",
         {"value": np.array(7, dtype=np.int64), "n": np.array(3, dtype=np.int64),
          "expected": np.full(3, 7, dtype=np.int64)})

    base = f64([1.0, 2.0, 3.0])
    save(out_dir, group, "align-expand-f64", "align",
         {"x": base, "value": np.array(-1.0), "n": np.array(5, dtype=np.int64),
          "expected": np.concatenate([base, np.full(2, -1.0)])})
    save(out_dir, group, "align-truncate-f64", "align",
         {"x": base, "value": np.array(-1.0), "n": np.array(2, dtype=np.int64),
          "expected": base[:2]})

    save(out_dir, group, "shape-1d", "shape",
         {"x": f64([1.0, 2.0, 3.0]), "expected": np.array([3, 0], dtype=np.int64)})
    matrix = np.array([[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]], dtype=np.float64)
    save(out_dir, group, "shape-2d", "shape",
         {"x": matrix, "expected": np.array([2, 3], dtype=np.int64)})

    a = np.array([1, 2], dtype=np.int64)
    b = np.array([3, 4], dtype=np.int64)
    c = np.array([5, 6], dtype=np.int64)
    save(out_dir, group, "concat1d-int", "concat1d",
         {"a": a, "b": b, "c": c, "expected": np.stack([a, b, c], axis=1)})

    source = np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int64)
    save(out_dir, group, "transpose2d-int", "transpose2d",
         {"x": source, "expected": np.ascontiguousarray(np.transpose(source))})

    with np.errstate(invalid="ignore"):
        holes = f64([1.0, np.nan, 3.0, np.nan])
    save(out_dir, group, "fill-f64", "fill",
         {"x": holes, "value": np.array(0.0),
          "expected": np.where(np.isnan(holes), 0.0, holes)})
    save(out_dir, group, "fill-na-f64", "fill_na",
         {"x": holes, "value": np.array(-9.0),
          "expected": np.where(np.isnan(holes), -9.0, holes)})


def rolling_windows(values, window, min_periods=None):
    """按 Go 的 v1At 口径生成窗口：min 未设置要求完整窗口，设置后用 partial window。"""
    out = []
    for i in range(len(values)):
        size = window[i] if isinstance(window, list) else window
        available = i + 1
        if min_periods is None or min_periods <= 0:
            if size > available:
                out.append([])
                continue
        else:
            if size > available:
                if available < min_periods:
                    out.append([])
                    continue
                size = available
            elif size < min_periods:
                out.append([])
                continue
        out.append(list(values[available - size:available]))

    lengths = [len(block) for block in out]
    flat = [value for block in out for value in block]
    return lengths, flat


def shift_expected(values, n):
    out = []
    for i, _ in enumerate(values):
        j = i - n
        out.append(values[j] if 0 <= j < len(values) else float("nan"))
    return out


def window_layer(out_dir: str) -> None:
    group = "window"
    series = f64([1.0, 2.0, 3.0, 4.0, 5.0])

    def save_rolling(case, window, min_periods=None):
        lengths, flat = rolling_windows(list(series), window, min_periods)
        payload = {
            "x": series,
            "window": np.array(window, dtype=np.int64),
            "expected_lengths": np.array(lengths, dtype=np.int64),
            "expected_flat": f64(flat),
        }
        if min_periods is not None:
            payload["min_periods"] = np.array(min_periods, dtype=np.int64)
        save(out_dir, group, case, "rolling", payload)

    save_rolling("rolling-basic-f64", 3)
    save_rolling("rolling-min1-f64", 3, 1)
    save_rolling("rolling-min2-f64", 3, 2)
    save_rolling("rolling-min-window-f64", 3, 3)
    save_rolling("rolling-min-gt-window-f64", 3, 4)
    save_rolling("rolling-window-array-f64", [2, 3, 3, 3, 3], 2)

    means = [float("nan") if i < 2 else float(np.mean(series[i - 2:i + 1])) for i in range(len(series))]
    save(out_dir, group, "rolling-apply-mean-f64", "rolling_apply",
         {"x": series, "window": np.array(3, dtype=np.int64), "expected": f64(means)})

    save(out_dir, group, "shift-right-f64", "shift",
         {"x": series, "n": np.array(2, dtype=np.int64), "expected": f64(shift_expected(list(series), 2))})
    save(out_dir, group, "shift-left-f64", "shift",
         {"x": series, "n": np.array(-1, dtype=np.int64), "expected": f64(shift_expected(list(series), -1))})
    save(out_dir, group, "shift-zero-f64", "shift",
         {"x": series, "n": np.array(0, dtype=np.int64), "expected": series})

    for n in (1, 2):
        shifted = shift_expected(list(series), n)
        expected = [float("nan") if math.isnan(a) else a - b for a, b in zip(shifted, series)]
        save(out_dir, group, f"diff-n{n}-f64", "diff",
             {"x": series, "n": np.array(n, dtype=np.int64), "expected": f64(expected)})

    save(out_dir, group, "cumprod-f64", "cumprod",
         {"x": f64([1.5, 2.0, -3.0]), "expected": np.cumprod(f64([1.5, 2.0, -3.0]))})



def linalg_layer(out_dir: str) -> None:
    group = "linalg"
    a = f64([[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]])
    b = f64([[7.0, 8.0], [9.0, 10.0], [11.0, 12.0]])
    save(out_dir, group, "dot2d-f64", "dot2d",
         {"a": a, "b": b, "expected": np.matmul(a, b)}, tolerance=1e-12)

    v = f64([1.0, 2.0, 3.0])
    save(out_dir, group, "dot2d1-f64", "dot2d1",
         {"a": a, "v": v, "expected": np.matmul(a, v)}, tolerance=1e-12)

    x = f64([[1.0, 2.0], [3.0, 4.0]])
    y = f64([[5.0, 6.0], [7.0, 8.0]])
    save(out_dir, group, "matmul-f64", "matmul",
         {"x": x, "y": y, "n": np.array(2, dtype=np.int64), "expected": np.matmul(x, y)},
         tolerance=1e-12)

    m4 = f64([[1.0, 0.0, 0.0, 0.0], [0.0, 2.0, 0.0, 0.0],
              [0.0, 0.0, 3.0, 0.0], [0.0, 0.0, 0.0, 4.0]])
    n4 = f64([[1.0, 2.0, 3.0, 4.0], [5.0, 6.0, 7.0, 8.0],
              [9.0, 10.0, 11.0, 12.0], [13.0, 14.0, 15.0, 16.0]])
    save(out_dir, group, "mat4mul-f64", "mat4mul",
         {"x": m4, "y": n4, "expected": np.matmul(m4, n4)}, tolerance=1e-12)

    inv = f64([[4.0, 7.0], [2.0, 6.0]])
    save(out_dir, group, "inverse-f64", "inverse",
         {"a": inv, "expected": np.linalg.inv(inv)}, tolerance=1e-12)


def geometry_layer(out_dir: str) -> None:
    group = "geometry"
    save(out_dir, group, "slope-f64", "slope",
         {"x1": np.array(1, dtype=np.int64), "y1": np.array(1.0), "x2": np.array(4, dtype=np.int64),
          "y2": np.array(10.0), "expected": np.array(3.0)})
    save(out_dir, group, "slope-to-degrees-f64", "slope_to_degrees",
         {"x": np.array(1.0), "expected": np.array(45.0)}, tolerance=1e-12)
    save(out_dir, group, "degrees-to-slope-f64", "degrees_to_slope",
         {"x": np.array(45.0), "expected": np.array(1.0)}, tolerance=1e-12)
    save(out_dir, group, "triangle-bevel-f64", "triangle_bevel",
         {"slope": np.array(2.5), "x1": np.array(1, dtype=np.int64), "y1": np.array(10.0),
          "xn": np.array(5, dtype=np.int64), "expected": np.array(20.0)})
    fit_x = f64([0.0, 1.0, 2.0, 3.0])
    fit_y = f64([1.0, 3.0, 7.0, 13.0])
    save(out_dir, group, "polyfit-f64", "poly_fit",
         {"x": fit_x, "y": fit_y, "deg": np.array(2, dtype=np.int64),
          "expected": np.polyfit(fit_x, fit_y, 2)}, tolerance=1e-10)

    coefficients = f64([2.0, -3.0, 1.0])
    xs = f64([0.0, 1.0, 2.0, 3.0])
    save(out_dir, group, "polyval-f64", "poly_val",
         {"p": coefficients, "x": xs, "expected": np.polyval(coefficients, xs)})


def conversion_layer(out_dir: str) -> None:
    group = "conversion"
    flags = np.array([True, False, True])
    save(out_dir, group, "from-bool-f64", "from_bool",
         {"x": flags, "expected": f64([1.0, 0.0, 1.0])})
    i32 = np.array([-2, 0, 3], dtype=np.int32)
    save(out_dir, group, "from-int32-f64", "from_int32",
         {"x": i32, "expected": i32.astype(np.float64)})
    i64 = np.array([-2, 0, 3], dtype=np.int64)
    save(out_dir, group, "from-int64-f64", "from_int64",
         {"x": i64, "expected": i64.astype(np.float64)})
    wide = f64([1.5, -2.25])
    save(out_dir, group, "from-float64-f32", "from_float",
         {"x": wide, "expected": wide.astype(np.float32)})
    narrow = f32([1.5, -2.25])
    save(out_dir, group, "from-float32-f64", "from_float",
         {"x": narrow, "expected": narrow.astype(np.float64)})
    mixed = f64([0.0, -0.0, 2.0, np.nan])
    with np.errstate(invalid="ignore"):
        to_bool = np.not_equal(mixed, 0.0)
    save(out_dir, group, "to-bool-f64", "to_bool", {"x": mixed, "expected": to_bool})
    trunc = f64([1.9, -2.9, 0.0])
    save(out_dir, group, "to-int32-f64", "to_int32",
         {"x": trunc, "expected": np.trunc(trunc).astype(np.int32)})
    save(out_dir, group, "to-int64-f64", "to_int64",
         {"x": trunc, "expected": np.trunc(trunc).astype(np.int64)})
    save(out_dir, group, "to-float32-f64", "to_float32",
         {"x": wide, "expected": wide.astype(np.float32)})
    save(out_dir, group, "to-float64-f32", "to_float64",
         {"x": narrow, "expected": narrow.astype(np.float64)})
    save(out_dir, group, "dtype2int-f64", "dtype2int",
         {"x": trunc, "expected": np.trunc(trunc).astype(np.int32)})


def line_layer(out_dir: str) -> None:
    group = "line"
    slope = 3.0
    intercept = -1.0
    probe = f64([0.0, 2.0, 5.0])
    degrees = math.degrees(math.atan(slope))
    save(out_dir, group, "line-equation-f64", "calculate_line_equation",
         {"x1": np.array(1.0), "y1": np.array(2.0), "x2": np.array(4.0), "y2": np.array(11.0),
          "probe_x": probe, "expected_y": f64([slope * x + intercept for x in probe]),
          "expected_a": np.array(slope), "expected_b": np.array(-1.0), "expected_c": np.array(intercept),
          "expected_degrees": np.array(degrees), "expected_success": np.array(1.0)})

    fit_x = f64([0.0, 1.0, 2.0, 3.0])
    fit_y = f64([1.0, 3.0, 5.0, 7.0])
    save(out_dir, group, "fit-line-f64", "fit_line",
         {"x": fit_x, "y": fit_y, "probe_x": fit_x,
          "expected_y": f64([2.0 * x + 1.0 for x in fit_x]),
          "expected_a": np.array(2.0), "expected_b": np.array(-1.0), "expected_c": np.array(1.0),
          "expected_degrees": np.array(math.degrees(math.atan(2.0))), "expected_success": np.array(1.0)})

    single_x = f64([2.0])
    single_probe = f64([0.0, 1.0, 2.0])
    save(out_dir, group, "fit-line-single-f64", "fit_line",
         {"x": single_x, "y": f64([5.0]), "probe_x": single_probe,
          "expected_y": f64([5.0, 5.0, 5.0]),
          "expected_a": np.array(0.0), "expected_b": np.array(-1.0), "expected_c": np.array(5.0),
          "expected_degrees": np.array(0.0), "expected_success": np.array(0.0)})


def cross_layer(out_dir: str) -> None:
    group = "cross"
    a = f64([1.0, 3.0, 1.0, 3.0])
    b = f64([2.0, 2.0, 2.0, 2.0])
    save(out_dir, group, "cross-f64", "cross",
         {"a": a, "b": b,
          "expected_x": np.array([1, 2, 3], dtype=np.int64),
          "expected_state": np.array([1, -1, 1], dtype=np.int64)})


def time_layer(out_dir: str) -> None:
    group = "time"
    freqs = np.array([b"D", b"2H", b"30min", b"500ms"], dtype="S8")
    nanos = np.array([86400 * 10**9, 2 * 3600 * 10**9, 30 * 60 * 10**9, 500 * 10**6],
                     dtype=np.int64)
    save(out_dir, group, "parse-freq", "parse_freq", {"freq": freqs, "expected": nanos})

    start = np.array(1767225600 * 10**9, dtype=np.int64)  # 2026-01-01T00:00:00Z
    periods = np.array(4, dtype=np.int64)
    step = 12 * 3600 * 10**9
    expected = np.array([int(start) + i * step for i in range(4)], dtype=np.int64)
    save(out_dir, group, "date-range", "date_range",
         {"start": start, "periods": periods,
          "freq": np.array(b"12H", dtype="S8"), "expected": expected})

def main() -> None:
    parser = argparse.ArgumentParser(description="生成跨语言 fixture")
    parser.add_argument("--out", default=os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."),
                        help="fixture 根目录，默认为 tests/fixtures")
    args = parser.parse_args()

    out_dir = os.path.abspath(args.out)
    scalar_layer(out_dir)
    elementwise(out_dir)
    elementwise_extras(out_dir)
    reduction(out_dir)
    reduction_extras(out_dir)
    container(out_dir)
    window_layer(out_dir)
    linalg_layer(out_dir)
    geometry_layer(out_dir)
    conversion_layer(out_dir)
    line_layer(out_dir)
    cross_layer(out_dir)
    time_layer(out_dir)
    print(f"fixtures written to {out_dir}")


if __name__ == "__main__":
    main()
