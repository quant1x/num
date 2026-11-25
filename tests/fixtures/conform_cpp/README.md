# tests/fixtures/conform_cpp —— C++ 对拍（C3 之后）

读取 `tests/fixtures` 的 npz，把 `num_core` 已有算子的结果与 `expected` 逐值比较。当前 **160 个样本全部通过**；
C++ 尚未实现的语义会跳过并打印名字，跳过不算通过，缺口见 [../CPP-GAPS.md](../CPP-GAPS.md)。

## 构建与运行

```bash
cmake --build --preset release-static --target num_fixture_conform
ctest --preset release-static -R fixtures_conform_cpp
# 或直接跑
./build/release-static/Release/num_fixture_conform.exe tests/fixtures
```

输出示例：

```text
conformed 160 fixtures, skipped 0 (0 semantics)
  skip ceil
  skip cos
  ...
```

## 覆盖范围

- 类型：`type_default`、`nan_check`、`slice_to_float64`、`slice_float_to_int`
- 逐元素：add / sub / mul / div / maximum / minimum / abs / sqrt / pow_integer / pow_float / expm1 / log1p /
  decimal / change_rate / net_change_rate；neg / inv / ceil / floor / round / sin / cos / exp / log / log2 / log10
- 比较与逻辑：equal / not_equal / gt / gte / lt / lte / logic_and / logic_or / logic_not / logic_xor / where
- 索引：select / find / gather / scatter
- 归约：sum / mean / std / variance / min / max / argmin / argmax / median / count / all / any / none /
  prod / norm / manhattan_norm / quantile / cumsum / dot1d / covariance
- 构造与整形：zeros / ones / arange / linspace / range / range(a,b) 浮点区间 / repeat / align / shape /
  concat1d / transpose2d / fill / fill_na
- 窗口：rolling（含 min_periods）/ rolling_apply / shift / diff / cumprod
- 线代：dot2d / dot2d1 / matmul / mat4mul / inverse
- 几何与多项式：slope / slope_to_degrees / degrees_to_slope / triangle_bevel / poly_fit / poly_val
- 转换族：from_bool / from_int32 / from_int64 / from_float / to_bool / to_int32 / to_int64 / to_float32 / to_float64 / dtype2int
- 直线：calculate_line_equation / fit_line（A/B/C、角度、探针 Y、成功位）
- 穿越：cross（下标 + 状态）
- 时间：parse_freq / date_range（纳秒整数）

## 未覆盖（跳过）

无：160 个样本全部覆盖。没有 fixture 的接口（mat4mul、转换族、`any` 入口）不在对拍范围，
清单见 `CPP-GAPS.md` §1。

## 口径

- 数值比较用 fixture 的 `meta_tolerance`；NaN 与 NaN 视为相等。
- C++ 的比较 / 逻辑按"缺位补 0"对齐，Go 的 `equal` / `not_equal` 按第一个操作数长度补 NaN；
  当前样本长度一致，不受影响，差异记录在 `CPP-GAPS.md`。
- 只读 fixture，不写文件、不联网。
