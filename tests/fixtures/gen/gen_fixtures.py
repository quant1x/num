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


def main() -> None:
    parser = argparse.ArgumentParser(description="生成跨语言 fixture")
    parser.add_argument("--out", default=os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."),
                        help="fixture 根目录，默认为 tests/fixtures")
    args = parser.parse_args()

    out_dir = os.path.abspath(args.out)
    elementwise(out_dir)
    reduction(out_dir)
    print(f"fixtures written to {out_dir}")


if __name__ == "__main__":
    main()
