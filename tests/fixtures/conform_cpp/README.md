# tests/fixtures/conform_cpp —— C++ 对拍（C3 之后）

读取 `tests/fixtures` 的 npz，把 `num_core` 已有算子的结果与 `expected` 逐值比较。能算的必须全绿；
C++ 尚未实现的语义跳过并打印名字，**跳过不算通过**，对应缺口见 [../CPP-GAPS.md](../CPP-GAPS.md) §1。

## 构建与运行

```bash
cmake --build --preset release-static --target num_fixture_conform
ctest --preset release-static -R fixtures_conform_cpp
# 或直接跑
./build/release-static/Release/num_fixture_conform.exe tests/fixtures
```

输出示例：

```text
conformed 94 fixtures, skipped 35 (27 semantics)
  skip ceil
  skip cos
  ...
```

## 覆盖范围

- 类型：`type_default`、`nan_check`、`slice_to_float64`、`slice_float_to_int`
- 逐元素：add / sub / mul / div / maximum / minimum / abs / sqrt / pow_integer / expm1 / log1p / decimal /
  change_rate / net_change_rate
- 比较与逻辑：equal / not_equal / gt / gte / lt / lte / logic_and / logic_or / logic_not / where
- 归约：sum / mean / std / variance / min / max / argmin / argmax / median / count / all / any / cumsum /
  dot1d / covariance
- 构造与整形：zeros / ones / arange / linspace / range / repeat / align / shape / concat1d / transpose2d /
  fill / fill_na
- 窗口：shift / diff

## 未覆盖（跳过）

滚动族（rolling / rolling_apply）、prod / norm / manhattan_norm / quantile / none、neg / inv / ceil / floor /
round、sin / cos / exp / log / log2 / log10、select / find / gather / scatter、logic_xor、range_float、
pow_float、bool_to_numeric 与转换族。补齐后从跳过列表移除，并在 `CPP-GAPS.md` 销账。

## 口径

- 数值比较用 fixture 的 `meta_tolerance`；NaN 与 NaN 视为相等。
- C++ 的比较 / 逻辑按"缺位补 0"对齐，Go 的 `equal` / `not_equal` 按第一个操作数长度补 NaN；
  当前样本长度一致，不受影响，差异记录在 `CPP-GAPS.md`。
- 只读 fixture，不写文件、不联网。
