#!/usr/bin/env python3
"""安装包冒烟（D3）：装到临时前缀，再用 `find_package(num)` 编一个消费方并运行。

用法（仓库根）：

    python tests/install/smoke.py

前置：已经按 BUILDING.md §6.2 构建过 `build/release-static`（安装要装那份产物）。
脚本会验证两件事 —— 安装树里有 `include/num/*.h` 与 `lib/cmake/num/numConfig.cmake`，
且消费方只用安装面就能编过并算出正确结果。
"""

from __future__ import annotations

import os
import shutil
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
BUILD = ROOT / "build" / "release-static"
PREFIX = ROOT / "temp" / "install"
CONSUMER_BUILD = ROOT / "temp" / "install-consumer"
EXPECTED = "10 2.5"


def run(command: list[str], what: str) -> None:
    result = subprocess.run(command, cwd=ROOT, capture_output=True, text=True,
                            encoding="utf-8", errors="replace")
    if result.returncode != 0:
        print(f"FAIL {what}\n{result.stdout}\n{result.stderr}")
        sys.exit(1)


def main() -> None:
    if not (BUILD / "CMakeCache.txt").exists():
        print(f"缺少构建目录：{BUILD}（先跑 cmake --build --preset release-static）")
        sys.exit(1)

    xtensor_dir = BUILD / "generators"
    if not (xtensor_dir / "xtensor-config.cmake").exists():
        print(f"找不到 xtensor 的 CMake 包：{xtensor_dir}")
        sys.exit(1)

    for path in (PREFIX, CONSUMER_BUILD):
        shutil.rmtree(path, ignore_errors=True)

    run(["cmake", "--install", str(BUILD), "--prefix", str(PREFIX)], "安装")
    for relative in ("include/num/sum.h", "lib/cmake/num/numConfig.cmake"):
        if not (PREFIX / relative).exists():
            print(f"FAIL 安装树缺少 {relative}")
            sys.exit(1)

    run([
        "cmake", "-S", str(ROOT / "tests" / "install" / "consumer"), "-B", str(CONSUMER_BUILD),
        "-DCMAKE_BUILD_TYPE=Release",
        # 消费方需要同时找到 num 与它的依赖 xtensor（conan 生成的包在构建目录的 generators 下）
        f"-DCMAKE_PREFIX_PATH={PREFIX}{os.pathsep}{xtensor_dir}",
    ], "配置消费方")
    run(["cmake", "--build", str(CONSUMER_BUILD), "--config", "Release"], "构建消费方")

    exe = CONSUMER_BUILD / "Release" / "consumer.exe"
    if not exe.exists():
        exe = CONSUMER_BUILD / "consumer"
    output = subprocess.run([str(exe)], capture_output=True, text=True,
                            encoding="utf-8", errors="replace").stdout.strip()
    if output != EXPECTED:
        print(f"FAIL 消费方输出 {output!r}，期望 {EXPECTED!r}")
        sys.exit(1)
    print(f"PASS 安装包冒烟（{output}）")


if __name__ == "__main__":
    main()
