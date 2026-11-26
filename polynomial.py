"""多项式拟合与求值（spec，对应 Go 的 `polynomial.go`）。

- `poly_fit(x, y, degree)`：最小二乘，系数**按次数降序**（`c[0]` 是最高次项）；
- `poly_val(p, x)`：按降序系数用 Horner 法求值。

degree 超过 `len(x)-1` 或方程奇异时抛 `ValueError`。
"""

from __future__ import annotations

from collections.abc import Sequence


def _solve(matrix: list[list[float]], rhs: list[float]) -> list[float]:
    size = len(matrix)
    augmented = [row[:] + [rhs[i]] for i, row in enumerate(matrix)]
    for column in range(size):
        pivot = max(range(column, size), key=lambda r: abs(augmented[r][column]))
        if augmented[pivot][column] == 0:
            raise ValueError("poly_fit: 方程组奇异")
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
    return [row[-1] for row in augmented]


def poly_fit(x: Sequence[float], y: Sequence[float], degree: int) -> list[float]:
    """最小二乘拟合，返回按次数降序的系数。"""
    if len(x) != len(y):
        raise ValueError("poly_fit: x 与 y 长度不一致")
    if degree < 0 or degree >= len(x):
        raise ValueError("poly_fit: degree 超出样本数")
    size = degree + 1
    # 正规方程：Σ x^(i+j) · c_j = Σ y·x^i
    matrix = [[sum(value ** (i + j) for value in x) for j in range(size)] for i in range(size)]
    rhs = [sum(y[k] * x[k] ** i for k in range(len(x))) for i in range(size)]
    ascending = _solve(matrix, rhs)
    return list(reversed(ascending))


def poly_val(p: Sequence[float], x: Sequence[float]) -> list[float]:
    """按降序系数求值。"""
    out: list[float] = []
    for value in x:
        total = 0.0
        for coefficient in p:
            total = total * value + coefficient
        out.append(total)
    return out
