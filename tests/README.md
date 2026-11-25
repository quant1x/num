# tests —— 跨语言一致性验收

本目录收**与语言无关的验收资产**：符号契约表、fixture 生成器与读取器。它不是某个 Go 包的源码目录，
而是 [docs/PLAN.md](../docs/PLAN.md) 中 C1 / C2 / C3（契约与 fixture）与 D2 / B5（CSV 交换与数据管道）的落点。

## 职责与边界

**收**：跨语言契约（语义名、最小参数集、各语言符号）、fixture 生成与读取工具、跨语言对拍入口。

**不收**：任何语言的算子实现（在各语言源码目录）、模块级语义文档（在代码同目录的 `.md`）、
一次性验证脚本（放 `temp/`）。

## 目录

| 路径 | 用途 |
| --- | --- |
| `contract/` | C1 符号与签名契约表 + 覆盖校验（`go test ./tests/contract/`） |
| `fixtures/` | C2 生成器、C3 三语言读取器、Go / C++ 对拍（`conform/`、`conform_cpp/`）与 `CPP-GAPS.md`，按 `<分组>/<用例>.npz` 存放样本 |
| `io/` | D2 跨语言 CSV 互读互写与 B5 按主键合并的驱动器与工具，见 `io/README.md` |

## 与上下层的关系

- 上游：语义来自 Python（np / pd）；Python 未落地期间以 NumPy 语义 + Go `Xxx_Go` + 模块 `.md` 作临时基准。
- 下游：`N1 ~ N4` 的实现与 C++ / Rust / Go 的发布物都按本目录的 fixture 判定一致性。
- 契约表是**唯一事实源**：新增或改名功能先改 `contract/NAME-MAP.tsv`，再改实现。

## 跨语言对应

| 资产 | python | cpp | rust | go |
| --- | --- | --- | --- | --- |
| 契约表 | `—`（读取方） | 读取方 | 读取方 | `contract/name_map_test.go` |
| fixture 生成 | `fixtures/gen/` | — | — | — |
| fixture 读取 | — | `fixtures/read_cpp/` | `fixtures/read_rs/` | `fixtures/read_go/` |
| CSV 读写 / 合并 | `io/roundtrip.py`、`io/pipeline.py` | `io/cpp_*` | `io/rs` | `io/` + `io/go_*` |
