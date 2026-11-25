"""float32 切片转换（spec，对应 Go 的 `SliceToFloat32`）。

口径：数值收窄到 f32 精度后回落 f64 表示，布尔按 1.0 / 0.0。
"""

from __future__ import annotations

import struct
from collections.abc import Sequence


def _narrow_f32(value: float) -> float:
    return struct.unpack("f", struct.pack("f", value))[0]


def slice_to_float32(x: Sequence[float | bool]) -> list[float]:
    """把数值切片转成 f32 精度（以 Python float 承载）。"""
    return [_narrow_f32(float(value)) for value in x]
