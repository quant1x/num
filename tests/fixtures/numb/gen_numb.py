#!/usr/bin/env python3
"""生成 `.numb` 黄金样本（D1 §3.2，版本 2）。

格式：magic `NUMB` | version u16 LE | dtype 2 字节 ASCII | ndim u8 | dims u64 LE × ndim | data。
三端（Go / C++ / Rust）读取这些文件并对读出的值与形状做断言，再把值重新编码，要求与文件逐字节一致。
"""

from __future__ import annotations

import os
import struct

HEADER_VERSION = 2


def encode(dtype: str, shape: list[int], payload: bytes) -> bytes:
    assert len(dtype) == 2, dtype
    out = bytearray()
    out += b"NUMB"
    out += struct.pack("<H", HEADER_VERSION)
    out += dtype.encode("ascii")
    out += struct.pack("<B", len(shape))
    for length in shape:
        out += struct.pack("<Q", length)
    out += payload
    return bytes(out)


def main() -> None:
    here = os.path.dirname(os.path.abspath(__file__))
    cases = {
        "f8-vector.numb": encode("f8", [5], struct.pack("<5d", 1.5, -2.0, 3.25, 0.0, 1e-3)),
        "f4-matrix.numb": encode("f4", [2, 3], struct.pack("<6f", *[float(v) for v in range(1, 7)])),
        "i8-vector.numb": encode("i8", [4], struct.pack("<4q", -1, 0, 7, 123456789)),
        "b1-vector.numb": encode("b1", [3], bytes([1, 0, 1])),
        "f8-scalar.numb": encode("f8", [], struct.pack("<d", 42.5)),
        "f8-empty.numb": encode("f8", [0], b""),
    }
    for name, data in cases.items():
        with open(os.path.join(here, name), "wb") as handle:
            handle.write(data)
    print(f"wrote {len(cases)} files to {here}")


if __name__ == "__main__":
    main()
