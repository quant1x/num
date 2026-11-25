#!/usr/bin/env python3
"""D2 跨语言 CSV 互读互写验证。

Python 写出带缺失值、无穷、引号与 UTF-8 的 CSV；Go / C++ / Rust 各自读入并按 D1 重新写出；
再用 Python 读回三份产物，与原文逐值比较，并检查每门语言重复写出的文件逐字节一致。
"""

from __future__ import annotations

import csv
import math
import os
import subprocess
import sys
from pathlib import Path

SCHEMA = [("date", "s"), ("close", "f"), ("volume", "i"), ("flag", "b"), ("note", "s")]
SCHEMA_SPEC = ",".join(f"{name}:{code}" for name, code in SCHEMA)

ROWS = [
    {"date": "2026-01-05", "close": 10.5, "volume": 100, "flag": True, "note": "沪市"},
    {"date": None, "close": math.inf, "volume": None, "flag": False, "note": 'a,b"c'},
    {"date": "2026-01-07", "close": -math.inf, "volume": -7, "flag": True, "note": None},
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


def parse_cell(kind: str, text: str):
    if text == "":
        return None
    if kind == "f":
        if text == "inf":
            return math.inf
        if text == "-inf":
            return -math.inf
        return float(text)
    if kind == "i":
        return int(text)
    if kind == "b":
        if text not in ("true", "false"):
            raise ValueError(f"非法布尔值 {text!r}")
        return text == "true"
    return text


def read_csv(path: Path):
    data = path.read_bytes()
    if data.startswith(b"\xef\xbb\xbf"):
        raise ValueError(f"{path}: 含 BOM")
    if b"\r" in data:
        raise ValueError(f"{path}: 含 CR，写出应为 LF")
    with path.open("r", encoding="utf-8", newline="") as handle:
        reader = csv.reader(handle)
        header = next(reader)
        if header != [name for name, _ in SCHEMA]:
            raise ValueError(f"{path}: 表头 {header}")
        rows = []
        for line, record in enumerate(reader, start=2):
            if len(record) != len(SCHEMA):
                raise ValueError(f"{path}: 第 {line} 行字段数 {len(record)}")
            rows.append({name: parse_cell(kind, record[i]) for i, (name, kind) in enumerate(SCHEMA)})
        return rows


def compare(label: str, got, want) -> None:
    if len(got) != len(want):
        raise AssertionError(f"{label}: 行数 {len(got)} != {len(want)}")
    for index, (actual, expected) in enumerate(zip(got, want), start=1):
        for name, kind in SCHEMA:
            a, b = actual[name], expected[name]
            if kind == "f" and isinstance(a, float) and isinstance(b, float):
                if math.isnan(a) and math.isnan(b):
                    continue
                if a == b:
                    continue
                raise AssertionError(f"{label}: 第 {index} 行 {name}: {a} != {b}")
            if a != b:
                raise AssertionError(f"{label}: 第 {index} 行 {name}: {a!r} != {b!r}")


def run(command, cwd: Path) -> None:
    result = subprocess.run(command, cwd=cwd, text=True, capture_output=True)
    if result.returncode != 0:
        raise RuntimeError(f"{command[0]} 失败: {result.stderr.strip() or result.stdout.strip()}")


def main() -> None:
    root = repo_root()
    work = root / "temp" / "tests" / "io" / "roundtrip"
    work.mkdir(parents=True, exist_ok=True)

    source = work / "python.csv"
    write_csv(source, ROWS)

    cpp_binary = Path(os.environ.get("NUM_CPP_CSV_BIN",
                                     root / "build" / "release-static" / "Release" / "num_csv_roundtrip.exe"))
    if not cpp_binary.is_file():
        cpp_binary = root / "build" / "release-static" / "num_csv_roundtrip"
    if not cpp_binary.is_file():
        raise SystemExit("找不到 num_csv_roundtrip，可先 cmake --build --preset release-static")

    builders = {
        "go": lambda in_path, out_path: run(
            ["go", "run", "./tests/io/go_roundtrip", SCHEMA_SPEC, str(in_path), str(out_path)], root),
        "cpp": lambda in_path, out_path: run(
            [str(cpp_binary), SCHEMA_SPEC, str(in_path), str(out_path)], root),
        "rust": lambda in_path, out_path: run(
            ["cargo", "run", "--quiet", "--release", "--manifest-path", "tests/io/rs/Cargo.toml",
             "--", SCHEMA_SPEC, str(in_path), str(out_path)], root),
    }

    outputs = {}
    for language, build in builders.items():
        out_path = work / f"{language}.csv"
        build(source, out_path)
        compare(f"{language} 读 Python", read_csv(out_path), ROWS)
        outputs[language] = out_path
        print(f"PASS {language}: 读 Python 产物并重写")

    for language, out_path in outputs.items():
        second = work / f"{language}.second.csv"
        builders[language](out_path, second)
        if second.read_bytes() != out_path.read_bytes():
            raise AssertionError(f"{language}: 重复写出不一致")
        print(f"PASS {language}: 重复写出逐字节一致")

    go_path = outputs["go"]
    compare("Python 读 Go", read_csv(go_path), ROWS)
    print("PASS Python: 读 Go 产物")
    print("ALL PASS")


if __name__ == "__main__":
    try:
        main()
    except Exception as error:  # noqa: BLE001 - 命令行入口统一报错
        print(f"FAIL: {error}", file=sys.stderr)
        sys.exit(1)
