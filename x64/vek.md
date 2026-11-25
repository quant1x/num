# vek：`float64` 向量算子门面

**边界**：本包（`github.com/quant1x/num/x64`，`package x64`）是 `[]float64` 的**强类型门面**：把切片算子以三种
形式暴露给调用方，并按 `functions.UseAVX2` 把每个算子分派到 `internal/functions` 的加速内核或 `Xxx_Go_*`
参考路径。它**不**做类型分派（那是根包 `num` 的活）、**不**提供标量数学函数（f64 直接用标准库 `math`）、
**不**承载容器与 IO。

**三种 API 形式**（同一算子三个入口，一套实现）：

| 形式 | 签名形态 | 是否克隆入参 | 返回值 |
| --- | --- | --- | --- |
| `Xxx` | `Xxx(x, ...) []float64` | 是（`slices.Clone` 后原地算） | 新切片 |
| `Xxx_Inplace` | `Xxx_Inplace(x, ...)` | 否，改第一个切片 | 无 |
| `Xxx_Into` | `Xxx_Into(dst, x, ...) []float64` | 否，结果写 `dst` | `dst[:len(x)]` |

**与 `x32` 的不对称**：本包**没有**超越函数（`Sin` / `Cos` / `Exp` / `Log*` 都不在 x64，f64 场景直接用标准库
`math`），独有 `FromFloat32` / `ToFloat32` 两个与 f32 的互转；`x32` 侧反过来独有 6 个 f32 超越函数与
`FromFloat64` / `ToFloat64`。

## 1. 路线图

**当前位置与下一步**：Go 侧已落地 —— 单文件 `vek.go`（1527 行）导出 156 个函数，全部经 `functions.UseAVX2`
分支到 AVX2（amd64）/ NEON（arm64）内核或 Go 参考路径。Python / Rust 未落地；C++ 侧**不设本层**
（强类型层与内核层的分工由模板与 xtensor 表达式吸收，见
[§3.8](../../CONTRIBUTING.md#38-c单层扁平无-include--src)），故本表 c++ 列留 `—`。按
[§6.4](../../CONTRIBUTING.md#64-落地状态与配套纪律) 判 **部分落地（缺 python / c++ / rust）**。下一个里程碑：
按 [§2.3](../../CONTRIBUTING.md#23-语言角色与演进顺序) 补 Python spec，把三形式契约与失败约定写进 spec。

**基准**：`tests/contract/NAME-MAP.tsv` 的 `authority` 逐行标注现行基准 —— 标 `python` 的以根目录同名 `.py` 为准，
其余以本文件 + `vek_test.go` 为准（[§9.3](../../CONTRIBUTING.md#93-spec--实现--验证)）。

| 阶段 | 语言（角色） | 载体 / 文件 | 状态 | 说明 |
| --- | --- | --- | --- | --- |
| spec 定语义 | Python（spec） | 不适用 | 不适用 | Python 只有一种浮点类型，x32 / x64 两个宽度合并到 `x32/vek.py`，不另建 `x64/` |
| 生产实现 | C++（生产） | —（不设本层） | — | 对应语义在仓库根门面层，见 [§3.8](../../CONTRIBUTING.md#38-c单层扁平无-include--src) |
| 对齐实现 | Rust（对齐） | `x64/vek.rs` | 未落地 | 见 §5 TODO |
| 对齐实现 | Go（对齐） | `x64/vek.go` | 已落地 | 156 个导出函数，用例见 §4.1 |

## 2. 功能清单

**统一口径**：

- **克隆语义**：只有不带后缀的 `Xxx` 克隆入参；`_Inplace` 与 `_Into` 都不克隆，调用方自己保证入参独占。
- **容量**：`_Into` 要求 `cap(dst) >= len(x)`，返回值为 `dst[:len(x)]`（长度被截到 `len(x)`，不是 `dst` 原长）。
- **失败约定**：一律 `panic`，不返回 `error` —— 长度不等（`slices must be of equal length`）、`dst` 与入参有
  重叠（`destination slice must not overlap input slice`）、`dst` 容量不足（`destination slice not large enough to hold result`）、
  归约类输入为空（`slice must not be empty`，`Mean` / `Median` / `Dot` 等 13 处）。
- **两条路径**：`functions.UseAVX2` 为真走向量内核，否则走 `Xxx_Go_*`；两条路径语义相同，差异只在性能。
- **比较类**：返回 `[]bool`；`NaN` 与任何值比较一律 `false`（含 `Neq`），不按位判等。
- **未命中**：`Find` 对外统一返回 `-1`；加速内核返回的 `len(x)` 由本层归一（内核层口径见
  `internal/functions/functions.md` §2）。
- **归约精度**：`Sum` / `Prod` / `Dot` / `Norm` 走内核时按 amd64 口径重排，用例按量级设容差
  （[§5.6.3](../../CONTRIBUTING.md#563-浮点归约的顺序敏感)）。

| 功能 | python（spec） | c++ | rust | go | zig | java | 状态 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 逐元素加 / 减 / 乘 / 除 | — | — | — | `Add*` / `Sub*` / `Mul*` / `Div*` | — | — | 部分落地（缺 python / c++ / rust） |
| 与标量加 / 减 / 乘 / 除 | — | — | — | `AddNumber*` / `SubNumber*` / `MulNumber*` / `DivNumber*` | — | — | 部分落地（缺 python / c++ / rust） |
| 取负 / 取倒数 / 绝对值 | — | — | — | `Neg*` / `Inv*` / `Abs*` | — | — | 部分落地（缺 python / c++ / rust） |
| 开方 / 取整 | — | — | — | `Sqrt*` / `Round*` / `Ceil*` / `Floor*` | — | — | 部分落地（缺 python / c++ / rust） |
| 幂 | — | — | — | `Pow*` | — | — | 部分落地（缺 python / c++ / rust） |
| 求和 / 求积 / 点积 | — | — | — | `Sum` / `Prod` / `Dot` | — | — | 部分落地（缺 python / c++ / rust） |
| 范数与距离 | — | — | — | `Norm` / `Distance` / `ManhattanNorm` / `ManhattanDistance` / `CosineSimilarity` | — | — | 部分落地（缺 python / c++ / rust） |
| 累计和 / 累计积 | — | — | — | `CumSum*` / `CumProd*` | — | — | 部分落地（缺 python / c++ / rust） |
| 最值与极值下标 | — | — | — | `Min` / `Max` / `ArgMin` / `ArgMax` / `Minimum*` / `Maximum*` / `MinimumNumber*` / `MaximumNumber*` | — | — | 部分落地（缺 python / c++ / rust） |
| 均值 / 中位数 / 分位 | — | — | — | `Mean` / `Median` / `Quantile` | — | — | 部分落地（缺 python / c++ / rust） |
| 比较 | — | — | — | `Lt*` / `Lte*` / `Gt*` / `Gte*` / `Eq*` / `Neq*` 及 `*Number*` 变体 | — | — | 部分落地（缺 python / c++ / rust） |
| 布尔逐元素 | — | — | — | `And*` / `Or*` / `Not*` / `Xor*` | — | — | 部分落地（缺 python / c++ / rust） |
| 布尔归约 | — | — | — | `All` / `Any` / `None` / `Count` | — | — | 部分落地（缺 python / c++ / rust） |
| 查找 | — | — | — | `Find`（未命中 `-1`） | — | — | 部分落地（缺 python / c++ / rust） |
| 构造与填充 | — | — | — | `Zeros*` / `Ones*` / `Repeat*` / `Range*` | — | — | 部分落地（缺 python / c++ / rust） |
| 选择 / 聚集 / 散布 | — | — | — | `Select*` / `Gather*` / `Scatter*` | — | — | 部分落地（缺 python / c++ / rust） |
| 转换 | — | — | — | `FromBool*` / `FromInt32*` / `FromInt64*` / `FromFloat32*` / `ToBool*` / `ToInt32*` / `ToInt64*` / `ToFloat32*` | — | — | 部分落地（缺 python / c++ / rust） |
| 矩阵 | — | — | — | `MatMul*` / `Mat4Mul*` | — | — | 部分落地（缺 python / c++ / rust） |

**单元格记法**：`*` 表示该算子族的三种形式（`Xxx` / `Xxx_Inplace` / `Xxx_Into`），并非通配符语法。
**状态判据**：[§6.4](../../CONTRIBUTING.md#64-落地状态与配套纪律)；四门必做语言中只有 Go 落地，故全表统一
「部分落地（缺 python / c++ / rust）」。

### 2.1 两层 API 对照

本包是 §2.1 所说的**强类型层**，与根包 `num` 的门面层构成同一语义的两套入口：门面层按运行时类型分发到本包
与 `x32`，本包直接面对 `internal/functions` 内核。跨语言列在本表不另立（只有 Go 落地）。

## 3. 依赖

**出向依赖（本模块依赖的东西）**：

| 依赖项 | 版本 / 来源 | 用途 | 适用语言 | 约束与备注 |
| --- | --- | --- | --- | --- |
| `github.com/quant1x/num/internal/functions` | 本仓 | 全部算子的加速内核与 Go 参考路径，以及 `UseAVX2` 开关 | go | 本包是它的主要适配点之一（70 余处 `UseAVX2` 判定） |
| `github.com/quant1x/num/internal/constraints` | 本仓 | 转换函数的泛型约束 | go | 与 `functions` 共用同一套约束 |
| `math` | Go 标准库 | 标量数学（`Pow` 的标量路径、`Median` 的插值等） | go | f64 无需自备数学函数，这是与 `x32` 的关键差别 |
| `slices` / `unsafe` | Go 标准库 | 克隆、重叠检测（切片头比较） | go | `unsafe` 只用于判断两个切片是否共享底层数组 |

**入向依赖（谁依赖本模块）**：

| 调用方 | 用途 | 备注 |
| --- | --- | --- |
| 根包 `num`（29 个文件） | 门面层的 f64 落点：`Abs` / `Sqrt` / 算术 / 归约 / 比较 / 类型转换等 | 反向依赖禁止：本包不得 import 根包或 `x32` |

## 4. 测试

### 4.1 单元测试

| 语言 | 测试文件 | 用例数 | 运行命令 |
| --- | --- | --- | --- |
| Python（spec） | 未落地 | — | — |
| C++ | 未落地 | — | — |
| Rust | 未落地 | — | — |
| Go | `x64/vek_test.go`（671 行） | 8 个 `Test`、0 个 `Benchmark` | `go test ./x64/` |

### 4.2 覆盖率

| 语言 | 工具 | 门槛 | 运行命令 |
| --- | --- | --- | --- |
| Go | `go test -cover` | 仓库未设门槛 | `go test -cover ./x64/` |

**实测**：`93.8% of statements`（`go test -count=1 -cover ./x64/`，HEAD `b05114e`）。

### 4.3 重复度

本模块不涉及：仓库尚未引入重复度检测工具（无 `.jscpd.json` 等配置），按 §6.3 保留本节并说明。

### 4.4 静态检查

| 语言 | 工具 | 运行命令 | 备注 |
| --- | --- | --- | --- |
| Go | `gofmt -l` ＋ `go vet` | `gofmt -l x64`、`go vet ./x64/` | 实测均无输出 |

### 4.5 基准测试

本模块不涉及：`x64/` 无 `Benchmark*`，性能口径由下层 `internal/functions/` 的基准承担
（`go test -bench . ./internal/functions/`）。

### 4.6 契约 / 一致性测试

| 场景 | 测试位置 | 断言内容 | 运行命令 |
| --- | --- | --- | --- |
| 与 `x32` 同构算子的行为一致 | 两侧 `vek_test.go` | 同一算子族在 f64 / f32 上的返回形态与失败约定一致 | `go test ./x64/ ./x32/` |
| 加速 / 回退两条路径 | `vek_test.go` 顶部 `accel []bool{false, functions.Info().Acceleration}` | 两种开关下同一算子结果一致 | `go test ./x64/` |

## 5. TODO

| 待办项 | 缺什么 | 卡在哪 | 关联章节 | 关联 issue | 优先级 |
| --- | --- | --- | --- | --- | --- |
| 补 Python spec | 无 `x64/vek.py`，三形式契约与失败约定无正式 spec | 无技术阻塞，待排期 | §1 | 待建单 | 高 |
| 与 `x32` 的 API 不对称 | x64 无超越函数（`Sin` / `Cos` / `Exp` / `Log*`），独有 `FromFloat32` / `ToFloat32`；是否需要跨宽度对齐 | 取决于 `x32` / `x64` 的 spec 定案 | §2 | 待建单 | 低 |
| 补目录 `README.md` | 本目录缺 `README.md`（[§6.1](../../CONTRIBUTING.md#61-每个目录必须有-readmemd)） | 无技术阻塞，待排期 | §1 | 待建单 | 中 |
| 设覆盖率门槛 | 仓库未设统一门槛，当前 93.8% | 无技术阻塞，待排期 | §4.2 | 待建单 | 低 |
