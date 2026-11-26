"""频率解析与时间序列。"""

from __future__ import annotations

_NANOS = {
    "N": 1, "ns": 1,
    "U": 1_000, "us": 1_000, "\u00b5s": 1_000,
    "L": 1_000_000, "ms": 1_000_000,
    "S": 1_000_000_000,
    "T": 60 * 1_000_000_000, "min": 60 * 1_000_000_000,
    "H": 3_600 * 1_000_000_000,
    "D": 86_400 * 1_000_000_000,
}

def parse_freq(text: str) -> int:
    """把频率字符串解析成纳秒。"""
    value = text.strip()
    if not value:
        raise ValueError("parse_freq: 空的频率字符串")
    index = 0
    while index < len(value) and value[index].isdigit():
        index += 1
    count = int(value[:index]) if index else 1
    unit = value[index:]
    if not unit:
        raise ValueError("parse_freq: 缺少单位")
    if unit not in _NANOS:
        raise ValueError(f"parse_freq: 不支持的单位 {unit!r}")
    return count * _NANOS[unit]

def date_range(start: int, periods: int, freq: str) -> list[int]:
    """等间隔生成纳秒时间戳。"""
    step = parse_freq(freq)
    return [start + index * step for index in range(periods)]
