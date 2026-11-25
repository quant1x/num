# tests/fixtures/geometry —— 几何与多项式样本

本目录覆盖斜率、角度换算、三角斜边与多项式求值，格式见 [../README.md](../README.md)。
标量参数用 0 维数组存放；角度与三角函数的用例容差 `1e-12`。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `slope-f64` | `slope` | 两点斜率，x 为整数、y 为浮点 |
| `slope-to-degrees-f64` | `slope_to_degrees` | 斜率 1 → 45° |
| `degrees-to-slope-f64` | `degrees_to_slope` | 45° → 1 |
| `triangle-bevel-f64` | `triangle_bevel` | `slope*(xn-x1)+y1` |
| `polyval-f64` | `poly_val` | 系数最高次在前，Horner 法 |

`poly_fit`、`calculate_line_equation`、`fit_line`、`cross` 等结构体返回与 `find_percent` / `find_zscore`
未纳入本组：算法口径或数据表需先确认，另行建单。
