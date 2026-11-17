# math32：`float32` 版标准库数学函数

**边界**：本包（`github.com/quant1x/num/math32`，`package math32`）提供 Go 标准库 `math` 的 `float32` 对应物 ——
标量函数与常量。它**不**提供切片 / 矩阵批量运算（那属 `x32/` 与 `internal/functions/`），也**不**做类型转换与
容器操作（那属根包 `num`）。

**来源**：本目录是上游 `chewxy/math32` 的引入（BSD 2-Clause，`Copyright (c) 2021, Xuanyi Chew and the Go Authors`），
许可证随目录 `LICENSE`，上游说明保留在 `README.md`；代码源自 Go 标准库的 `float32` 复刻。语义基线与 Go
标准库 `math` 对齐，本仓不做语义改写的承诺在 §2「统一口径」逐条列出。

## 1. 路线图

**当前位置与下一步**：Go 侧已落地 —— 62 个导出函数、3 组常量，含 amd64 / arm64 汇编实现
（`exp_*.s` / `log_amd64.s` / `remainder_amd64.s` / `sqrt_*.s`）与 7 个其它平台的桩（`stubs_*.s`）。Python / C++ /
Rust 无对应实现，而且本目录属**外部实现引入**，不在 [§2.3](../CONTRIBUTING.md#23-语言角色与演进顺序) 的四门
推广序列内：其它语言若要 `float32` 数学函数，应由各自的 spec 定语义后自行落地。故 §2 的状态按
[§6.4](../CONTRIBUTING.md#64-落地状态与配套纪律) 记为「部分落地（缺 python / c++ / rust）」，缺口的性质是
**未移植**，不是**未开工**。下一个里程碑：先把本目录与上游的对应关系（版本 / commit）登记清楚，再评估是否
随仓移植。

**临时基准**：Python 未落地期间，本文件与 `math32` 用例共同充当临时基准
（[§9.3](../CONTRIBUTING.md#93-spec--实现--验证)），Python 落地后以其为准。

| 阶段 | 语言（角色） | 载体 / 文件 | 状态 | 说明 |
| --- | --- | --- | --- | --- |
| spec 定语义 | Python（spec） | `math32/*.py` | 未落地 | 语义沿用 Go 标准库 `math`，不另立口径；见 §5 TODO |
| 生产实现 | C++（生产） | `math32/*.h` / `*.cpp` | 未落地 | 见 §5 TODO |
| 对齐实现 | Rust（对齐） | `math32/*.rs` | 未落地 | 见 §5 TODO |
| 对齐实现 | Go（对齐） | `math32/*.go`（48 个实现文件 + 2 个测试文件）+ 13 个 `.s` | 已落地 | 上游实现引入，用例见 §4.1 |

## 2. 功能清单

**统一口径**：

- **返回值**：全部为 `float32`（`Frexp` / `Modf` / `Lgamma` / `Sincos` 另带 `int` 或第二个返回值）；
  位操作函数返回 `uint32` / `uint64`。
- **实现分两类**：`doc.go` 声明的 25 个函数是 `float32(math.XXX)` 薄包装（`Acos` / `Acosh` / `Asin` / `Asinh` /
  `Atan` / `Atan2` / `Atanh` / `Cbrt` / `Cos` / `Cosh` / `Erfc` / `Gamma` / `J0` / `J1` / `Jn` / `Log10` / `Log1p` /
  `Log2` / `Logb` / `Pow10` / `Sin` / `Sinh` / `Tan` / `Y0` / `Y1`），其余为 `float32` 原生实现。薄包装存在
  **二次舍入**（`float64` 计算后转 `float32`），精度判据以此为准。
- **失败约定**：不 panic、不返回 `error`。定义域外按 IEEE 754 语义给 `NaN` / `±Inf` / `±0`，与标准库 `math`
  一致；判定一律用 `IsNaN` / `IsInf` / `Signbit`，不用 `==` 与 `NaN` 比较。
- **平台差异**：amd64 与 arm64 有手写汇编（`exp_*.s`、`log_amd64.s`、`remainder_amd64.s`、`sqrt_*.s`）；
  386 / arm / ppc64le / riscv64 / s390x / wasm 由 `stubs_*.s` 转到 Go 路径。两条路径对同一组参考值断言，
  见 §4.6。
- **无批量语义**：本包只收标量；切片版本在 `x32/`，加速内核在 `internal/functions/`。

| 功能 | python（spec） | c++ | rust | go | zig | java | 状态 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 数学常量 | — | — | — | `E` / `Pi` / `Phi` / `Sqrt2` / `SqrtE` / `SqrtPi` / `SqrtPhi` | — | — | 部分落地（缺 python / c++ / rust） |
| 浮点与整数极值常量 | — | — | — | `MaxFloat32` / `SmallestNonzeroFloat32` / `MaxInt8` … `MinInt64` / `MaxUint8` … `MaxUint64` | — | — | 部分落地（缺 python / c++ / rust） |
| 绝对值 / 取大 / 取小 | — | — | — | `Abs` / `Max` / `Min` | — | — | 部分落地（缺 python / c++ / rust） |
| 取整 | — | — | — | `Ceil` / `Floor` / `Round` / `Trunc` | — | — | 部分落地（缺 python / c++ / rust） |
| 根式 | — | — | — | `Sqrt` / `Cbrt` / `Hypot` | — | — | 部分落地（缺 python / c++ / rust） |
| 幂与指数 | — | — | — | `Pow` / `Pow10` / `Exp` / `Exp2` / `Expm1` | — | — | 部分落地（缺 python / c++ / rust） |
| 对数 | — | — | — | `Log` / `Log10` / `Log2` / `Log1p` / `Logb` / `Ilogb` | — | — | 部分落地（缺 python / c++ / rust） |
| 三角函数 | — | — | — | `Sin` / `Cos` / `Sincos` / `Tan` / `Asin` / `Acos` / `Atan` / `Atan2` | — | — | 部分落地（缺 python / c++ / rust） |
| 双曲函数 | — | — | — | `Sinh` / `Cosh` / `Tanh` / `Asinh` / `Acosh` / `Atanh` | — | — | 部分落地（缺 python / c++ / rust） |
| 误差函数与伽马函数 | — | — | — | `Erf` / `Erfc` / `Gamma` / `Lgamma` | — | — | 部分落地（缺 python / c++ / rust） |
| 贝塞尔函数 | — | — | — | `J0` / `J1` / `Jn` / `Y0` / `Y1` / `Yn` | — | — | 部分落地（缺 python / c++ / rust） |
| 分解与合成 | — | — | — | `Frexp` / `Ldexp` / `Modf` / `Mod` / `Remainder` / `Dim` | — | — | 部分落地（缺 python / c++ / rust） |
| 相邻浮点与符号 | — | — | — | `Nextafter` / `Copysign` / `Signbit` | — | — | 部分落地（缺 python / c++ / rust） |
| 特殊值与判定 | — | — | — | `NaN` / `Inf` / `IsNaN` / `IsInf` | — | — | 部分落地（缺 python / c++ / rust） |
| 位表示互转 | — | — | — | `Float32bits` / `Float32frombits` / `Float64bits` / `Float64frombits` | — | — | 部分落地（缺 python / c++ / rust） |

**状态判据**：[§6.4](../CONTRIBUTING.md#64-落地状态与配套纪律)。四门必做语言中目前只有 Go 落地，且本目录属
外部实现引入（见 §1），故全表统一为「部分落地（缺 python / c++ / rust）」；`zig` / `java` 为扩展列，不参与判定。

### 2.1 两层 API 对照

本包 API 只有一层（标量函数 + 常量），与调用方 `x32/` 的「强类型切片 API」不是同一层语义，故本节不涉及。

## 3. 依赖

**出向依赖（本模块依赖的东西）**：

| 依赖项 | 版本 / 来源 | 用途 | 适用语言 | 约束与备注 |
| --- | --- | --- | --- | --- |
| `math` | Go 标准库 | 薄包装函数（25 个）的计算落点 | go | 标准库 `math` 只有 `float64`，`float32` 版本必须自备 |
| `math/bits` | Go 标准库 | 位分解 / `Ilogb` 等位级实现 | go | 与 `Float32bits` 配合避免额外分配 |
| `unsafe` | Go 标准库 | `Float32bits` / `Float32frombits` / `Float64bits` / `Float64frombits` 的零分配位互转 | go | 四个函数无法用纯算术表达，指针重解释是唯一无拷贝做法 |

无三方运行时依赖。

**入向依赖（谁依赖本模块）**：

| 调用方 | 用途 | 备注 |
| --- | --- | --- |
| `x32/vek.go` | `Max` / `Ceil` | 向量级取大与取整需要 `float32` 版本 |
| 根包 `sqrt.go` | `Sqrt` | 根包 `Sqrt` 的非加速回退路径 |
| `internal/functions/special*.go` / `distance.go` / `accel_neon_arm64.go` | `Sqrt` / `Abs` / `Round` / `Floor` / `Ceil` / `Pow` / `Sin` / `Cos` / `Exp` / `Log` / `Log2` / `Log10` | f32 加速内核的标量兜底与 tail 处理 |

## 4. 测试

### 4.1 单元测试

| 语言 | 测试文件 | 用例数 | 运行命令 |
| --- | --- | --- | --- |
| Python（spec） | 未落地 | — | — |
| C++ | 未落地 | — | — |
| Rust | 未落地 | — | — |
| Go | `math32/all_test.go`（外部测试包 `math32_test`）、`math32/export_test.go`（导出内部 `sqrt` / `exp` / `hypot` 供对拍） | 60 个 `Test` + 57 个 `Benchmark` | `go test ./math32/` |

`all_test.go` 沿用 Go 标准库的做法：输入 `vf` 与期望值取自 26 位高精度计算器（`keisan.casio.com`）的参考值，
期望值硬编码在用例里。

### 4.2 覆盖率

| 语言 | 工具 | 门槛 | 运行命令 |
| --- | --- | --- | --- |
| Go | `go test -cover` | 仓库未设门槛 | `go test -cover ./math32/` |

**实测**：`81.0% of statements`（`go test -count=1 -cover ./math32/`，HEAD `a2f3fc2` + 工作树未提交改动）。

### 4.3 重复度

本模块不涉及：仓库尚未引入重复度检测工具（无 `.jscpd.json` 等配置），按 §6.3 保留本节并说明。

### 4.4 静态检查

| 语言 | 工具 | 运行命令 | 备注 |
| --- | --- | --- | --- |
| Go | `gofmt -l` ＋ `go vet` | `gofmt -l math32`、`go vet ./math32` | `gofmt -l math32` 实测无输出；`go vet ./math32` **未通过**，报 `math32/exp.go:34: unreachable code`，见 §5 TODO |

### 4.5 基准测试

| 语言 | 基准路径 | 指标 | 基线（含环境与 commit） | 运行命令 |
| --- | --- | --- | --- | --- |
| Go | `math32/all_test.go` 内 57 个 `Benchmark*`（`BenchmarkAcos` / `BenchmarkExp` / `BenchmarkSqrt` …） | `ns/op`、`B/op`、`allocs/op` | 未登记（见 §5 TODO） | `go test -run=NONE -bench=. -benchmem ./math32/` |

### 4.6 契约 / 一致性测试

本包没有跨语言契约（只有 Go 落地），一致性体现在**汇编路径与 Go 路径对拍**：

| 场景 | 测试位置 | 断言内容 | 运行命令 |
| --- | --- | --- | --- |
| `Sqrt` 汇编路径 vs 纯 Go `sqrt` | `math32/all_test.go`（`SqrtGo` 来自 `export_test.go`） | 同一组高精度参考值下两路径都命中，且 `Sqrt` 对 `±Inf` / `NaN` 的返回与标准库一致 | `go test -run 'TestSqrt' ./math32/` |
| `Exp` / `Hypot` 的 Go 路径对拍 | 用例与基准当前被注释（`all_test.go` 中 `testExp(t, ExpGo, …)`、`BenchmarkExpGo`） | 未生效，见 §5 TODO | — |

## 5. TODO

| 待办项 | 缺什么 | 卡在哪 | 关联章节 | 关联 issue | 优先级 |
| --- | --- | --- | --- | --- | --- |
| 登记上游对应版本 | 目录内无上游 `chewxy/math32` 的 commit / tag 记录，无法判断本目录相对上游改了什么 | 无技术阻塞，待排期；需比对上游历史 | §1 / §3 | 待建单 | 中 |
| 清 `go vet` 报错 | `math32/exp.go:34: unreachable code` | 无技术阻塞；需确认是死分支还是 `panic` 之后的正常不可达 | §4.4 | 待建单 | 中 |
| 恢复 `Exp` / `Hypot` 的 Go 路径对拍 | `testExp(t, ExpGo, …)` 与 `BenchmarkExpGo`、`BenchmarkSqrtGo` 被注释，纯 Go 路径在 `Exp` / `Hypot` 上无断言 | 无技术阻塞，待排期 | §4.6 | 待建单 | 中 |
| 决定是否随仓移植 | 无 `math32/*.py` / `*.h` / `*.rs`，四门其余语言缺对应实现 | 属外部实现引入，需先定「是否纳入本仓推广序列」 | §1 / §2 | 待建单 | 中 |
| 设覆盖率门槛 | 81.0%，仓库未设统一门槛 | 无技术阻塞，待排期 | §4.2 | 待建单 | 低 |
| 登记基准基线 | 57 个 `Benchmark*` 无基线值与环境说明 | 需固定基准机并写明 CPU / OS / commit | §4.5 | 待建单 | 低 |
