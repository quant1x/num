"""NumPy linalg：点积、矩阵乘、求逆与范数。"""

from __future__ import annotations

import math
from collections.abc import Sequence

def dot(a: Sequence[float], b: Sequence[float]) -> float:
    """返回 a 与 b 的内积。"""
    if len(a) != len(b):
        raise ValueError("dot1d: 两条序列长度不一致")
    total = 0.0
    for left, right in zip(a, b):
        total += left * right
    return total

def concat1d(a, b, c) -> list[list[float]]:
    """按列拼成 n×3。"""
    if len(a) != len(b) or len(a) != len(c):
        raise ValueError("concat1d: 三条序列长度不一致")
    return [[left, middle, right] for left, middle, right in zip(a, b, c)]

def transpose(x) -> list[list[float]]:
    """二维转置。"""
    if not x:
        return []
    width = len(x[0])
    return [[row[column] for row in x] for column in range(width)]

def dot2d(a, b) -> list[list[float]]:
    """矩阵乘 a·b。"""
    return matmul(a, b)

def dot2d1(a, v) -> list[float]:
    """矩阵乘向量。"""
    if a and len(a[0]) != len(v):
        raise ValueError("dot2d1: 列数与向量长度不一致")
    return [sum(left * right for left, right in zip(row, v)) for row in a]

def matmul(a, b) -> list[list[float]]:
    """矩阵乘 a·b。"""
    if not a or not b:
        return []
    if len(a[0]) != len(b):
        raise ValueError("matmul: 维度不匹配")
    columns = list(zip(*b))
    return [[sum(left * right for left, right in zip(row, column)) for column in columns]
            for row in a]

def mat4mul(a, b) -> list[list[float]]:
    """4×4 矩阵乘。"""
    if len(a) != 4 or len(b) != 4:
        raise ValueError("mat4mul: 只接受 4×4")
    return matmul(a, b)

def inv(a) -> list[list[float]]:
    """Gauss-Jordan 求逆。"""
    size = len(a)
    if size == 0 or any(len(row) != size for row in a):
        raise ValueError("inverse: 需要方阵")
    augmented = [list(map(float, row)) + [1.0 if i == j else 0.0 for j in range(size)]
                 for i, row in enumerate(a)]
    for column in range(size):
        pivot = max(range(column, size), key=lambda r: abs(augmented[r][column]))
        if augmented[pivot][column] == 0:
            raise ValueError("inverse: 矩阵奇异")
        augmented[column], augmented[pivot] = augmented[pivot], augmented[column]
        scale = augmented[column][column]
        augmented[column] = [value / scale for value in augmented[column]]
        for row in range(size):
            if row == column:
                continue
            factor = augmented[row][column]
            if factor != 0:
                augmented[row] = [current - factor * base
                                  for current, base in zip(augmented[row], augmented[column])]
    return [row[size:] for row in augmented]

def norm(x: Sequence[float]) -> float:
    """欧氏范数。"""
    total = 0.0
    for value in x:
        total += value * value
    return math.sqrt(total)

def manhattan_norm(x: Sequence[float]) -> float:
    """曼哈顿范数。"""
    total = 0.0
    for value in x:
        total += abs(value)
    return total

def distance(x: Sequence[float], y: Sequence[float]) -> float:
    """欧氏距离。"""
    if len(x) != len(y):
        raise ValueError("distance: 两条序列长度不一致")
    total = 0.0
    for left, right in zip(x, y):
        delta = left - right
        total += delta * delta
    return math.sqrt(total)

def manhattan_distance(x: Sequence[float], y: Sequence[float]) -> float:
    """曼哈顿距离。"""
    if len(x) != len(y):
        raise ValueError("manhattan_distance: 两条序列长度不一致")
    total = 0.0
    for left, right in zip(x, y):
        total += abs(left - right)
    return total

def cosine_similarity(x: Sequence[float], y: Sequence[float]) -> float:
    """余弦相似度。"""
    if len(x) != len(y):
        raise ValueError("cosine_similarity: 两条序列长度不一致")
    dot = 0.0
    norm_x = 0.0
    norm_y = 0.0
    for left, right in zip(x, y):
        dot += left * right
        norm_x += left * left
        norm_y += right * right
    return dot / math.sqrt(norm_x * norm_y)
