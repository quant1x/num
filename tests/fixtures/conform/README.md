# tests/fixtures/conform —— Go 实现对拍（M1 准入）

把 `tests/fixtures` 下每个 npz 交给 Go 实现计算，与 `expected` 逐值比对。这是"N1 ~ N3、N6 的 fixture 全绿"
的可执行判据，也是后续 C++ / Rust 对拍用例的模板。

## 文件

| 文件 | 内容 |
| --- | --- |
| `doc.go` | 包声明 |
| `conform_test.go` | 遍历分组 → 读 fixture → 按 `meta_semantic` 分派到 Go API → 比较 |

## 用法

```bash
go test ./tests/fixtures/conform/           # 全量
go test ./tests/fixtures/conform/ -run 'TestFixtures/container' -v
```

失败输出定位到"文件 + 第几位 + 实际值 / 期望值 / 容差"。

## 分派约定

- `meta_semantic` 决定调用哪个算子；dtype 决定走 f32 还是 f64 入口。
- 标量参数（`n`、`q`、`size`、`start` 等）以 0 元素数组存放，由 `scalarFloat` / `scalarInt` / `scalarBool` 读取。
- 二维结果（`concat1d` / `transpose2d`）按行主序展平后比较；`shape` 的期望是 `[行, 列]`。
- 不等长输出（`rolling`）走 `expected_flat` + `expected_lengths`：先逐个窗口比对长度，再比对展平值。
- 容差取 `meta_tolerance`：0 表示逐位相等，NaN 与 NaN 视为相等。

## 新增样本

1. 在 `gen_fixtures.py` 加样本；
2. 在 `compute` 加一个 `case`，复用现成的 `unary` / `binary` / `compare` / `reduceFloat` 等助手；
3. 重跑 `python tests/fixtures/gen/gen_fixtures.py` 与本用例；
4. 若该语义在 C++ 侧还没实现，同步更新 [../CPP-GAPS.md](../CPP-GAPS.md)。

## 边界

只验证 Go 一侧；跨语言一致性由三套读取器（C3）与后续各语言的 conform 用例承担。不修改 fixture，
不为了通过而放宽容差——容差只能由 `gen_fixtures.py` 的样本口径决定。
