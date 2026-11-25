# tests/contract —— 符号与签名契约表（C1）

`NAME-MAP.tsv` 是 C++ / Rust / Go / Python 之间**语义名与最小参数集**的唯一事实源
（[docs/PLAN.md](../../docs/PLAN.md) §2.1）。`NAME-MAP-EXCLUDED.tsv` 登记**有意不纳入契约**的公开符号及理由。

## 文件

| 文件 | 内容 |
| --- | --- |
| `NAME-MAP.tsv` | 每行一个功能：语义名、最小参数集、默认值、NaN 约定、返回形状、各语言符号、权威实现语言 |
| `NAME-MAP-EXCLUDED.tsv` | 未纳入契约的符号：包、符号（`*` 表示整包）、理由 |
| `name_map_test.go` | 校验列数、语义名唯一、取值合法、编码，以及 `go/ast` 枚举的公开符号是否全覆盖 |

## 列与取值

| 列 | 含义 |
| --- | --- |
| `semantic` | 语义名，snake_case，优先取 NumPy 惯用名 |
| `min_params` | 最小参数集：业务必需参数，增删改名即视为语义变更 |
| `defaults` | 可省略参数的默认值；无默认值写 `—` |
| `nan_convention` | `propagate` / `comparison_false` / `no_special` / `nan_to_zero` / `default` / `—` |
| `return_shape` | `scalar` / `tuple` / `index`（未命中返回 `-1`）/ `first_len` / `max_len` / `same_as(x)` / `bool[...]` / `[n]` / `[n][m]` 等 |
| `python` `cpp` `rust` `go` `zig` `java` | 各语言符号名；`go` 列带包限定符 `num:` / `x32:` / `x64:` |
| `authority` | 权威实现语言：已由 Python spec 实现的写 `python`，其余写 `go`（临时基准）；标 `python` 的行必须在 `python` 列有符号 |
| `note` | 边界、失败约定与已知缺口 |

## 覆盖规则

- `go` 列登记**基础符号**即可；`_Inplace` / `_Into` / `Inplace` / `Number` 等写入与标量变体由
  校验用例按后缀视为同一功能。
- 校验只枚举**导出顶层函数**；方法（`Line.Angle`、`Window.At` 等）不在函数粒度契约内。
- `rust` 列登记 `src/*.rs` 里的 `pub fn` 名称（逗号分隔），由 `TestNameMapRustSymbols` 校验存在性；
  Rust 未实现的语义写 `—`。
- `python` 列登记 `模块.函数`（模块是仓库根的 `.py` 文件名），由 `TestNameMapPythonSymbols`
  校验「文件存在且含 `def 函数名(`」；Python 未实现的语义写 `—`。
- python 列写 `—` 的两种情形都要在 `note` 里写明：**未实现**（待排期）或**不适用**
  （如 `avx2_switch` 的运行期开关、`abs_int` / `abs_float` 的定宽类型）。
- 新增公开函数时二选一：在 `NAME-MAP.tsv` 补一行，或在 `NAME-MAP-EXCLUDED.tsv` 写明排除理由。
- `tests/contract` 只读当前工作树，不联网、不读写 `temp/`。

## 校验

```bash
go test ./tests/contract/
```

失败信息会列出四类问题：`未登记契约表`（新增符号漏登记）、`登记的符号不存在`（Go 改名或删除后未同步）、
`rust 符号在 src/ 中不存在`、`python 符号在根目录 .py 里不存在`。
