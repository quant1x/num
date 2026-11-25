# tests/fixtures/linalg —— N6 线代样本

本目录覆盖矩阵乘、矩阵乘向量、4×4 乘与求逆，格式见 [../README.md](../README.md)。
矩阵按行主序存放为二维数组，`expected` 也是二维；浮点容差 `1e-12`。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `dot2d-f64` | `dot2d` | 2×3 乘 3×2 |
| `dot2d1-f64` | `dot2d1` | 2×3 乘 3 维向量 |
| `matmul-f64` | `matmul` | 2×2 行主序，`n=2` |
| `mat4mul-f64` | `mat4mul` | 4×4 行主序 |
| `inverse-f64` | `inverse` | 2×2 可逆矩阵，不做主元搜索 |

`polyfit` / `polyval` 与几何、时间语义未纳入本组：算法口径需先确认，另行建单。
