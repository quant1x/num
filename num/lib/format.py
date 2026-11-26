"""`.numb` 二进制格式编解码。"""

from __future__ import annotations

import struct
from dataclasses import dataclass

MAGIC = b"NUMB"

VERSION = 2

HEADER_SIZE = 9

WIDTHS = {"b1": 1, "i1": 1, "i2": 2, "i4": 4, "i8": 8,
          "u1": 1, "u2": 2, "u4": 4, "u8": 8, "f4": 4, "f8": 8}

class InvalidNumb(ValueError):
    """头部、版本、dtype 或长度不自洽。"""

@dataclass(frozen=True)
class Numb:
    dtype: str
    shape: tuple[int, ...]
    payload: bytes

    def count(self) -> int:
        total = 1
        for length in self.shape:
            total *= length
        return total

    def floats(self) -> list[float]:
        if self.dtype != "f8":
            raise InvalidNumb(f"dtype {self.dtype!r} 不是 f8")
        return list(struct.unpack(f"<{self.count()}d", self.payload))

    def floats32(self) -> list[float]:
        if self.dtype != "f4":
            raise InvalidNumb(f"dtype {self.dtype!r} 不是 f4")
        return list(struct.unpack(f"<{self.count()}f", self.payload))

    def int64s(self) -> list[int]:
        if self.dtype != "i8":
            raise InvalidNumb(f"dtype {self.dtype!r} 不是 i8")
        return list(struct.unpack(f"<{self.count()}q", self.payload))

    def bools(self) -> list[bool]:
        if self.dtype != "b1":
            raise InvalidNumb(f"dtype {self.dtype!r} 不是 b1")
        out = []
        for index, value in enumerate(self.payload):
            if value not in (0, 1):
                raise InvalidNumb(f"第 {index} 个布尔值非法（{value}）")
            out.append(bool(value))
        return out

def encode(numb: Numb) -> bytes:
    width = WIDTHS.get(numb.dtype)
    if width is None:
        raise InvalidNumb(f"未知 dtype 码 {numb.dtype!r}")
    expected = numb.count() * width
    if len(numb.payload) != expected:
        raise InvalidNumb(f"载荷 {len(numb.payload)} 字节，期望 {expected}")
    out = bytearray()
    out += MAGIC
    out += struct.pack("<H", VERSION)
    out += numb.dtype.encode("ascii")
    out += bytes([len(numb.shape)])
    for length in numb.shape:
        out += struct.pack("<Q", length)
    out += numb.payload
    return bytes(out)

def decode(data: bytes) -> Numb:
    if len(data) < HEADER_SIZE:
        raise InvalidNumb(f"文件过短（{len(data)} 字节）")
    if data[:4] != MAGIC:
        raise InvalidNumb("magic 不是 NUMB")
    version = struct.unpack("<H", data[4:6])[0]
    if version != VERSION:
        raise InvalidNumb(f"版本 {version} 不受支持（当前 {VERSION}）")
    dtype = data[6:8].decode("ascii")
    width = WIDTHS.get(dtype)
    if width is None:
        raise InvalidNumb(f"未知 dtype 码 {dtype!r}")
    ndim = data[8]
    header = HEADER_SIZE + 8 * ndim
    if len(data) < header:
        raise InvalidNumb(f"ndim={ndim} 但文件只有 {len(data)} 字节")
    shape = tuple(struct.unpack("<Q", data[HEADER_SIZE + 8 * i:HEADER_SIZE + 8 * i + 8])[0]
                  for i in range(ndim))
    count = 1
    for length in shape:
        count *= length
    want = header + count * width
    if len(data) != want:
        raise InvalidNumb(f"文件长度 {len(data)}，期望 {want}")
    return Numb(dtype=dtype, shape=shape, payload=data[header:])

def from_f64(shape: tuple[int, ...], values: list[float]) -> Numb:
    return Numb("f8", shape, struct.pack(f"<{len(values)}d", *values))
