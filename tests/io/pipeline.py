#!/usr/bin/env python3
"""M6 数据管道验证：按主键覆盖 + 幂等。

Python 写出基线与增量；Go / C++ / Rust 各自把增量合并进基线；Python 读回三份结果逐值比较，
再重复合并一次，确认产物逐字节一致。
"""

from __future__ import annotations

import csv
import math
import os
import shutil
import subprocess
import sys
from pathlib import Path

SCHEMA = [("account", "s"), ("close", "f"), ("volume", "i")]
SCHEMA_SPEC = ",".join(f"{name}:{code}" for name, code in SCHEMA)
KEY = "account"

BASE = [
    {"account": "a", "close": 1.5, "volume": 10},
    {"account": "b", "close": 2.5, "volume": 20},
    {"account": "c", "close": 3.5, "volume": 30},
]
DELTA = [
    {"account": "b", "close": 9.5, "volume": 99},
    {"account": "d", "close": 4.5, "volume": 40},
]
EXPECTED = [
    {"account": "a", "close": 1.5, "volume": 10},
    {"account": "b", "close": 9.5, "volume": 99},
    {"account": "c", "close": 3.5, "volume": 30},
    {"account": "d", "close": 4.5, "volume": 40},
]


def repo_root() -> Path:
    return Path(__file__).resolve().parents[2]


def cell_text(kind: str, value) -> str:
    if value is None:
        return ""
    if kind == "f":
        if math.isinf(value):
            return "inf" if value > 0 else "-inf"
        return repr(float(value))
    if kind == "b":
        return "true" if value else "false"
    return str(value)


def write_csv(path: Path, rows) -> None:
    with path.open("w", encoding="utf-8", newline="") as handle:
        writer = csv.writer(handle, lineterminator="\n")
        writer.writerow([name for name, _ in SCHEMA])
        for row in rows:
            writer.writerow([cell_text(kind, row[name]) for name, kind in SCHEMA])


def read_csv(path: Path):
    data = path.read_bytes()
    if data.startswith(b"\xef\xbb\xbf"):
        raise ValueError(f"{path}: 含 BOM")
    if b"\r" in data:
        raise ValueError(f"{path}: 含 CR")
    with path.open("r", encoding="utf-8", newline="") as handle:
        reader = csv.reader(handle)
        header = next(reader)
        if header != [name for name, _ in SCHEMA]:
            raise ValueError(f"{path}: 表头 {header}")
        rows = []
        for line, record in enumerate(reader, start=2):
            if len(record) != len(SCHEMA):
                raise ValueError(f"{path}: 第 {line} 行字段数 {len(record)}")
            row = {}
            for (name, kind), text in zip(SCHEMA, record):
                if text == "":
                    row[name] = None
                elif kind == "f":
                    row[name] = math.inf if text == "inf" else -math.inf if text == "-inf" else float(text)
                elif kind == "i":
                    row[name] = int(text)
                else:
                    row[name] = text
            rows.append(row)
        return rows


def compare(label: str, got, want) -> None:
    if len(got) != len(want):
        raise AssertionError(f"{label}: 行数 {len(got)} != {len(want)}")
    for index, (actual, expected) in enumerate(zip(got, want), start=1):
        if actual != expected:
            raise AssertionError(f"{label}: 第 {index} 行 {actual!r} != {expected!r}")


def run(command, cwd: Path) -> None:
    result = subprocess.run(command, cwd=cwd, text=True, capture_output=True, check=False)
    if result.returncode != 0:
        raise RuntimeError(f"{command[0]} 失败: {result.stderr.strip() or result.stdout.strip()}")


def main() -> None:
    root = repo_root()
    work = root / "temp" / "tests" / "io" / "pipeline"
    work.mkdir(parents=True, exist_ok=True)

    base_path = work / "base.csv"
    delta_path = work / "delta.csv"
    write_csv(base_path, BASE)
    write_csv(delta_path, DELTA)

    cpp_binary = Path(os.environ.get("NUM_CPP_MERGE_BIN",
                                     root / "build" / "release-static" / "Release" / "num_csv_merge.exe"))
    if not cpp_binary.is_file():
        cpp_binary = root / "build" / "release-static" / "num_csv_merge"
    if not cpp_binary.is_file():
        raise SystemExit("找不到 num_csv_merge，可先 cmake --build --preset release-static")

    tools = {
        "go": lambda target, delta: run(
            ["go", "run", "./tests/io/go_merge", SCHEMA_SPEC, str(target), str(delta), KEY], root),
        "cpp": lambda target, delta: run([str(cpp_binary), SCHEMA_SPEC, str(target), str(delta), KEY], root),
        "rust": lambda target, delta: run(
            ["cargo", "run", "--quiet", "--release", "--manifest-path", "tests/io/rs/Cargo.toml",
             "--bin", "num-csv-merge", "--", SCHEMA_SPEC, str(target), str(delta), KEY], root),
    }

    for language, merge in tools.items():
        target = work / f"{language}.csv"
        shutil.copyfile(base_path, target)
        merge(target, delta_path)
        compare(f"{language} 合并结果", read_csv(target), EXPECTED)
        first = target.read_bytes()
        merge(target, delta_path)
        if target.read_bytes() != first:
            raise AssertionError(f"{language}: 重复合并不幂等")
        print(f"PASS {language}: 合并更新 + 追加 + 幂等")

    print("ALL PASS")


if __name__ == "__main__":
    try:
        main()
    except Exception as error:  # noqa: BLE001 - 命令行入口统一报错
        print(f"FAIL: {error}", file=sys.stderr)
        sys.exit(1)
