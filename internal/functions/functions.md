# functions：向量化内核与架构分发（`internal`）

**边界**：本包（`github.com/quant1x/num/internal/functions`，`package functions`）是**加速层**：每个算子都有
两套实现 —— Go 参考路径 `Xxx_Go_*` 与架构内核 `Xxx_AVX2_F32/F64`，后者在 amd64 走 AVX2、在 arm64 走 NEON，
**同一符号名跨架构不区分**；再外加能力判定与开关（`UseAVX2` / `SetAcceleration` / `Info`）。它**不**提供
对外门面（那属根包 `num` 与 `x32/` / `x64/`），也**不**做类型分派、容器与 IO。包在 `internal/` 下，仓外不可导入。

**三层分工**：

| 层 | 文件 | 职责 |
| --- | --- | --- |
| 参考路径 | `*.go` 里的 `Xxx_Go_*`（如 `sum.go` 的 `Sum_Go`） | 纯 Go、与架构无关，是内核的对拍基准 |
| 架构内核 | `accel_avx2_amd64.go` / `accel_avx2_amd64.s`、`accel_neon_arm64.go` / `accel_neon_arm64.s` | 同名 `Xxx_AVX2_F32/F64`；汇编只处理整块，尾块回落 Go |
| 兜底与开关 | `accel_noasm.go`、`accel_amd64.go`、`accel_arm64.go`、`cpu_info.go` | 其它平台一律 `panic("not implemented")` 桩；`simdSupported()` 判定与 `SetAcceleration` 开关 |

## 1. 路线图

**当前位置与下一步**：Go 侧两套架构都已落地 —— amd64 的 AVX2（`accel_avx2_amd64.s`，11579 行）与 arm64 的
NEON（`accel_neon_arm64.s`，1695 行 / 28 个 `TEXT` 符号），其余平台由 `accel_noasm.go` 提供桩并在
`UseAVX2 = false` 下不被调用。Python / C++ / Rust 均未落地，按
[§6.4](../../CONTRIBUTING.md#64-落地状态与配套纪律) 现判 **部分落地（缺 python / c++ / rust）**。下一个里程碑：
补 Python spec 作为内核语义准绳（归约顺序、未命中返回值、NaN 传播），再按
[§2.3](../../CONTRIBUTING.md#23-语言角色与演进顺序) 推广。

**临时基准**：Python 未落地期间，本文件与 `Xxx_Go_*` 参考路径共同充当临时基准
（[§9.3](../../CONTRIBUTING.md#93-spec--实现--验证)），Python 落地后以其为准。

| 阶段 | 语言（角色） | 载体 / 文件 | 状态 | 说明 |
| --- | --- | --- | --- | --- |
| spec 定语义 | Python（spec） | `internal/functions/*.py` | 未落地 | 见 §5 TODO |
| 生产实现 | C++（生产） | `internal/functions/*.h` / `*.cpp` | 未落地 | 见 §5 TODO |
| 对齐实现 | Rust（对齐） | `internal/functions/*.rs` | 未落地 | 见 §5 TODO |
| 对齐实现 | Go（对齐） | `internal/functions/*.go`（39 个）＋ 2 个 `.s` | 已落地 | amd64 AVX2 ＋ arm64 NEON ＋ 其余平台桩 |

## 2. 功能清单

**统一口径**：

- **命名**：参考实现 `Xxx_Go_*`，内核 `Xxx_AVX2_F32` / `Xxx_AVX2_F64`。`_F32` / `_F64` 指元素宽度，
  `_AVX2_` 是**跨架构固定后缀**（arm64 上由 NEON 承载同名符号），不得改名为 `_NEON_*` / `_SIMD_*`
  （[§4.4](../../CONTRIBUTING.md#44-固定后缀与跨架构符号名)）。
- **块大小与尾块**：内核按整块处理（f64 每块 4 个、f32 每块 8 个），不足一块或余下的尾元素由同名 Go 逻辑
  处理；`Exp` / `Log` / `Log2` / `Log10` 在 arm64 上也只对 `len &^ 7` 的部分走向量，尾部用 `math32`。
- **归约顺序**：乘积归约的向量化门槛对齐 amd64 —— f64 至少 16 个元素、f32 至少 32 个元素才走向量内核，
  不足门槛走顺序乘法；门槛是口径而非性能旋钮，改了严格容差用例立即漂移。求和 / 点积在用例里用
  `InDelta` / `InDeltaSlice` 比较，不得为过用例放宽全局容差
  （[§5.6.3](../../CONTRIBUTING.md#563-浮点归约的顺序敏感)）。
- **特殊值语义**：比较类算子对 `NaN` 一律返回 false，`±Inf` 按 IEEE 处理，不按位判等。
- **未命中语义**：口径是 `Find_AVX2_*` 未命中返回 `len(x)`（[§5.6.3](../../CONTRIBUTING.md#563-浮点归约的顺序敏感)），
  适配层据此判定；`ArgMin` / `ArgMax` 按 `len(x)` 判未命中并归一为 `-1` 对外返回。当前参考路径
  `Find_Go` 仍返回 `-1`，与口径不一致，见 §5 TODO。
- **布尔归约签名**：`All_AVX2` / `Any_AVX2` / `None_AVX2` 返回 `int`（1 / 0），对应的 `All_Go` / `Any_Go` /
  `None_Go` 返回 `bool`；调用方按整数判真值。
- **失败约定**：除 amd64 / arm64 之外的平台由 `accel_noasm.go` 提供 `panic("not implemented")` 桩，缺口显式；
  这些桩在 `UseAVX2 = false` 时不会被适配层调用（[§5.9](../../CONTRIBUTING.md#59-卡点先抛不支持)）。
- **开关**：`UseAVX2` 由 `simdSupported()` 决定 —— amd64 为 `AVX2 && FMA && GOOS != "darwin"`，arm64 为
  `cpu.ARM64.HasASIMD`，其余为 `false`。`SetAcceleration` **非线程安全**，只允许在用例与基准里切换
  （[§5.3](../../CONTRIBUTING.md#53-调试开关)）。

| 功能 | python（spec） | c++ | rust | go | zig | java | 状态 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 逐元素加 / 减 / 乘 / 除 | — | — | — | `Add_AVX2_F*` / `Sub_AVX2_F*` / `Mul_AVX2_F*` / `Div_AVX2_F*` ＋ `*_Go` | — | — | 部分落地（缺 python / c++ / rust） |
| 与标量运算 | — | — | — | `AddNumber_*` / `SubNumber_*` / `MulNumber_*` / `DivNumber_*` | — | — | 部分落地（缺 python / c++ / rust） |
| 取负 / 取倒数 / 绝对值 | — | — | — | `Neg_AVX2_F*` / `Inv_AVX2_F*` / `Abs_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 求和 / 求积 | — | — | — | `Sum_AVX2_F*` / `Prod_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 累计和 / 累计积 | — | — | — | `CumSum_AVX2_F*` / `CumProd_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 点积 | — | — | — | `Dot_AVX2_F*` / `Dot_Go` | — | — | 部分落地（缺 python / c++ / rust） |
| 范数与距离 | — | — | — | `Norm_AVX2_F*` / `Distance_AVX2_F*` / `ManhattanNorm_AVX2_F*` / `ManhattanDistance_AVX2_F*` / `CosineSimilarity_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 最值 | — | — | — | `Min_AVX2_F*` / `Max_AVX2_F*` / `Minimum_AVX2_F*` / `Maximum_AVX2_F*` / `MinNumber_*` / `MaxNumber_*` | — | — | 部分落地（缺 python / c++ / rust） |
| 极值下标 | — | — | — | `ArgMin_AVX2_F*` / `ArgMax_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 聚合统计 | — | — | — | `Mean_AVX2_F*` / `Median_AVX2_F*` / `Quantile_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 比较（逐元素） | — | — | — | `Lt_*` / `Lte_*` / `Gt_*` / `Gte_*` / `Eq_*` / `Neq_*` 及 `*Number_*` 变体 | — | — | 部分落地（缺 python / c++ / rust） |
| 布尔逐元素 | — | — | — | `And_AVX2` / `Or_AVX2` / `Not_AVX2` / `Xor_AVX2` | — | — | 部分落地（缺 python / c++ / rust） |
| 布尔归约 | — | — | — | `All_AVX2` / `Any_AVX2` / `None_AVX2` / `Count_AVX2`（返回 `int`） | — | — | 部分落地（缺 python / c++ / rust） |
| 查找 | — | — | — | `Find_AVX2_F*`（未命中返回 `len(x)`） | — | — | 部分落地（缺 python / c++ / rust） |
| 取整与开方 | — | — | — | `Round_AVX2_F*` / `Floor_AVX2_F*` / `Ceil_AVX2_F*` / `Sqrt_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 指数与对数 | — | — | — | `Exp_AVX2_F32` / `Log_AVX2_F32` / `Log2_AVX2_F32` / `Log10_AVX2_F32`（arm64 走 NEON，f64 走 Go） | — | — | 部分落地（缺 python / c++ / rust） |
| 幂与三角 | — | — | — | `Pow_AVX2_F*` / `Sin_AVX2_F32` / `Cos_AVX2_F32` / `SinCos_AVX2_F32`（arm64 暂走 Go） | — | — | 部分落地（缺 python / c++ / rust） |
| 矩阵乘 | — | — | — | `MatMul_AVX2_F*` / `MatMulVec_AVX2_F*` / `MatMulTiled_AVX2_F*` / `Mat4Mul_AVX2_F*` / `MatMul_Parallel_AVX2_F*` | — | — | 部分落地（缺 python / c++ / rust） |
| 构造与填充 | — | — | — | `Zeros_Go` / `Repeat_AVX2_F*` / `Range_AVX2_F*` / `Zeros_AVX2_F64` / `Ones_AVX2_F64` | — | — | 部分落地（缺 python / c++ / rust） |
| 类型转换 | — | — | — | `ToBool_AVX2_F*` / `ToInt32_AVX2_F*` / `ToInt64_AVX2_F*` / `FromBool_AVX2_F*` / `FromInt32_AVX2_F*` / `FromInt64_AVX2_F*` / `FromFloat32_AVX2_F64` / `FromFloat64_AVX2_F32` | — | — | 部分落地（缺 python / c++ / rust） |
| 选择 / 散布 / 聚集 | — | — | — | `Select_Go` / `Gather_Go` / `Scatter_Go` / `Clone_Go` / `Clone_Go_Loop` | — | — | 部分落地（缺 python / c++ / rust） |
| 随机数据 | — | — | — | `Random` / `RandomBool` / `RandomRange` | — | — | 部分落地（缺 python / c++ / rust） |
| 能力与开关 | — | — | — | `UseAVX2` / `SetAcceleration` / `Info` / `simdSupported` | — | — | 部分落地（缺 python / c++ / rust） |

**状态判据**：[§6.4](../../CONTRIBUTING.md#64-落地状态与配套纪律)。四门必做语言中只有 Go 落地，故全表统一为
「部分落地（缺 python / c++ / rust）」；`zig` / `java` 为扩展列，不参与判定。

### 2.1 两层 API 对照

本包本身就是「两层」的实现处，但它是**内部层**，不对外：上层的 `x32/` / `x64/` 才是公开门面，
本包与它们的关系见 §3 入向依赖。故本节的跨语言对照不另立。

## 3. 依赖

**出向依赖（本模块依赖的东西）**：

| 依赖项 | 版本 / 来源 | 用途 | 适用语言 | 约束与备注 |
| --- | --- | --- | --- | --- |
| `github.com/quant1x/num/internal/constraints` | 本仓 | `Float` / `Integer` 等类型集约束 | go | 泛型参考实现的类型参数来源 |
| `github.com/quant1x/num/internal/partial` | 本仓 | `TopK`（`Median` / `Quantile` 的选择） | go | 避免全排序，控制最坏复杂度 |
| `github.com/quant1x/num/internal/rand` | 本仓 | 用例与基准的随机数据（固定种子可复现） | go | [§7](../../CONTRIBUTING.md#7-测试规范) 要求失败可重跑 |
| `github.com/quant1x/num/math32` | 本仓 `math32/` | `Exp` / `Log` / `Sqrt` 等尾块与 Go 回退的标量实现 | go | 标准库无 `float32` 版数学函数 |
| `github.com/quant1x/num/vectors` | 本仓 | 比较结果打包等向量辅助 | go | 与内核共享同一批位运算语义 |
| `golang.org/x/sys/cpu` | v0.50.0 | `HasAVX2` / `HasFMA` / `HasASIMD` / `HasSVE` 能力探测 | go | 标准库无运行时 CPU 特性探测 |
| `math` / `runtime` / `slices` / `sync` | Go 标准库 | 标量数学、`GOOS` 判定、克隆与并行矩阵乘的同步 | go | 标准库足够 |

**入向依赖（谁依赖本模块）**：

| 调用方 | 用途 | 备注 |
| --- | --- | --- |
| `x32/vek.go` | 近 300 处 `functions.Xxx_AVX2_*` 调用与 `functions.UseAVX2` 判定 | f32 与 f64 两条路径都经它分发，是主要适配点 |
| `x64/*.go` | 70 余处 `functions.UseAVX2` 判定与各算子内核调用 | f64 强类型层 |
| 根包 `builtin.go` | `SetAcceleration` / `UseAVX2`（对外暴露为 `SetAvx2Enabled` / `GetAvx2Enabled`） | 门面层唯一的开关入口 |

## 4. 测试

### 4.1 单元测试

| 语言 | 测试文件 | 用例数 | 运行命令 |
| --- | --- | --- | --- |
| Python（spec） | 未落地 | — | — |
| C++ | 未落地 | — | — |
| Rust | 未落地 | — | — |
| Go | `internal/functions/` 下 11 个 `*_test.go`（`arithmetic` / `aggregates` / `boolean` / `comparison` / `construct` / `distance` / `find` / `matrix` / `max` / `min` / `special`） | 25 个 `Test` + 26 个 `Benchmark` | `go test ./internal/functions/` |

用例按「`Xxx_Go` vs `Xxx_AVX2_*` 逐元素对拍」组织，长度覆盖 `1`、块大小 `±1` 与多个块；随机数据取自
`internal/rand` 并固定 `rand.Seed(...)`。

### 4.2 覆盖率

| 语言 | 工具 | 门槛 | 运行命令 |
| --- | --- | --- | --- |
| Go | `go test -cover` | 仓库未设门槛 | `go test -cover ./internal/functions/` |

**实测**：`64.5% of statements`（`go test -count=1 -cover ./internal/functions/`，HEAD `e80a04e`）。

### 4.3 重复度

本模块不涉及：仓库尚未引入重复度检测工具（无 `.jscpd.json` 等配置），按 §6.3 保留本节并说明。

### 4.4 静态检查

| 语言 | 工具 | 运行命令 | 备注 |
| --- | --- | --- | --- |
| Go | `gofmt -l` ＋ `go vet` | `gofmt -l internal/functions`、`go vet ./internal/functions/` | 实测两者在 amd64 与 `GOARCH=arm64` 下均无输出 |

### 4.5 基准测试

| 语言 | 基准路径 | 指标 | 基线（含环境与 commit） | 运行命令 |
| --- | --- | --- | --- | --- |
| Go | `internal/functions/*_test.go` 内 26 个 `Benchmark*` | `ns/op`、`B/op`、`allocs/op` | 未登记（见 §5 TODO） | `go test -run=NONE -bench=. -benchmem ./internal/functions/` |

基准覆盖 `sizes = {10, 100, 1_000, 10_000, 100_000}`（[BUILDING.md §4.2](../../BUILDING.md)）。arm64 结果必须在
真机上取；qemu 下的绝对耗时不可信，只能看相对趋势。

### 4.6 契约 / 一致性测试

本包没有跨语言契约（只有 Go 落地），一致性体现在**架构之间**与**两条路径之间**：

| 场景 | 测试位置 | 断言内容 | 运行命令 |
| --- | --- | --- | --- |
| `Xxx_Go` vs `Xxx_AVX2_*` 对拍 | 各 `*_test.go` | 逐元素相等；长度 1、块大小 ±1、多块；`NaN` / `±Inf` / `±0` | `go test ./internal/functions/` |
| arm64 NEON vs Go 路径 | 交叉编译后经 qemu 执行 | 同一批用例在 arm64 上全绿 | 见下 |
| 两条路径开关 | 由适配层承担：`x32/vek_test.go` 与 `x64/vek_test.go` 用 `accel := []bool{false, functions.Info().Acceleration}` 跑两遍 | 本包用例直接对拍 `Xxx_Go` 与 `Xxx_AVX2_*`，不需要开关 | `go test ./x32/ ./x64/` |

arm64 验证流程（[BUILDING.md §5.6](../../BUILDING.md)）：

```powershell
$env:GOOS = "linux"; $env:GOARCH = "arm64"
go test -c -o temp/internal/functions/arm64/functions.test ./internal/functions
wsl -d Debian bash -lc "cd /mnt/d/projects/quant1x/num && chmod +x temp/internal/functions/arm64/functions.test && qemu-aarch64-static temp/internal/functions/arm64/functions.test"
```

## 5. TODO

| 待办项 | 缺什么 | 卡在哪 | 关联章节 | 关联 issue | 优先级 |
| --- | --- | --- | --- | --- | --- |
| 统一 `Find` 未命中返回值 | 参考路径 `Find_Go` 返回 `-1`，内核 `Find_AVX2_*` 返回 `len(x)`；`TestFind` 只挑数组内已存在的值，未命中分支无用例 | 属语义口径问题，需先定 spec 再改两侧与用例 | §2 / §4.1 | 待建单 | 高 |
| 统一布尔归约签名 | `All_AVX2` / `Any_AVX2` / `None_AVX2` 返回 `int`，`All_Go` / `Any_Go` / `None_Go` 返回 `bool` | 需确认适配层调用方是否需要布尔语义 | §2 | 待建单 | 中 |
| arm64 补 `Pow` / `Sin` / `Cos` / `SinCos` 内核 | `special_arm64.go` 里这四个直接调 Go 路径，无 NEON 实现 | 无技术阻塞，待排期 | §2 | 待建单 | 中 |
| 清理 `Zeros_AVX2_F64` / `Ones_AVX2_F64` | arm64 与 noasm 有定义、amd64 无定义，且全仓无调用方 | 需决定是补 amd64 实现还是删掉这三处 | §2 | 待建单 | 中 |
| 补 Python spec | 无 `*.py`，归约顺序、未命中返回值等口径无正式 spec | 无技术阻塞，待排期 | §1 | 待建单 | 高 |
| 补目录 `README.md` | 本目录缺 `README.md`（[§6.1](../../CONTRIBUTING.md#61-每个目录必须有-readmemd)），文件清单与层级关系只在本文档里 | 无技术阻塞，待排期 | §1 | 待建单 | 中 |
| 设覆盖率门槛 | 64.5%，仓库未设统一门槛 | 无技术阻塞，待排期 | §4.2 | 待建单 | 低 |
| 登记基准基线 | 26 个 `Benchmark*` 无基线值；arm64 侧需真机数据 | 需固定基准机并写明 CPU / OS / commit | §4.5 | 待建单 | 低 |
