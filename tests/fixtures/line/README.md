# tests/fixtures/line —— 直线与拟合样本

本目录覆盖两点直线方程与最小二乘拟合，格式见 [../README.md](../README.md)。Line 是多字段返回：
fixture 用 `expected_a` / `expected_b` / `expected_c` / `expected_degrees` / `expected_success` 与
`probe_x` → `expected_y` 描述，conform 逐字段比较。

| 用例 | 语义 | 覆盖点 |
| --- | --- | --- |
| `line-equation-f64` | `calculate_line_equation` | 两点 (1,2)/(4,11)，斜率 3、截距 -1 |
| `fit-line-f64` | `fit_line` | 完美线性 x=[0..3]、y=2x+1，成功位 true |
| `fit-line-single-f64` | `fit_line` | 单点：水平线 y=5、成功位 false |

字段口径：A = slope、B = -1、C = intercept（Go `Line.Equation()`、C++ `Line::equation()`）；
`expected_degrees` 是倾斜角；`expected_success` 为 1/0。
