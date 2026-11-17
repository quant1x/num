# num：门面数值 API —— 按泛型类型分发到 `x32` / `x64` 与加速内核

**边界**：本包（`github.com/quant1x/num`，`package num`）是仓库的**门面层**：把同一语义的强类型实现
（`x32/`、`x64/`）与加速内核（`internal/functions/`）收在一条泛型签名后面，负责**类型分发**与**回退选择**。
它**不**定义新的数值语义，也**不**承载 IO、序列化、时间序列数据容器；算法本体在 `x32/` / `x64/` /
`internal/functions/` / `internal/partial/` / `math32/` / `binary/` / `vectors/` 各自成包。

**文件组织**：逐语义一文件平铺在仓库根（`abs.go`、`sum.go`、`dot.go`、`rolling.go` …），本文件是整包的
功能级文档，§2 一行一个功能对应到具体符号；本包无 `num.go`，文档以包名命名。

## 1. 路线图

**当前位置与下一步**：Go 侧门面 API 已落地（64 个实现文件 + 49 个测试文件），覆盖算术、归约、逻辑比较、
数组构造、滑动窗口、线代与几何、正态分布表、类型转换九类。Python spec、C++、Rust 三侧均未落地，按
[CONTRIBUTING.md §6.4](CONTRIBUTING.md#64-落地状态与配套纪律) 现阶段整包判为
**部分落地（缺 python / c++ / rust）**。下一个里程碑：先补 Python spec 作为语义准绳，再按
[§2.3](CONTRIBUTING.md#23-语言角色与演进顺序) 顺序推广到 C++ 与 Rust。

**临时基准**：Python 未落地期间，本文件与 Go 用例共同充当临时基准
（[§9.3](CONTRIBUTING.md#93-spec--实现--验证)），Python 落地后以其为准。

| 阶段 | 语言（角色） | 载体 / 文件 | 状态 | 说明 |
| --- | --- | --- | --- | --- |
| spec 定语义 | Python（spec） | 仓库根 `*.py`（与 `.go` 同名） | 未落地 | 语义源头，见 §5 TODO |
| 生产实现 | C++（生产） | 仓库根 `*.h` / `*.cpp` | 未落地 | 见 §5 TODO |
| 对齐实现 | Rust（对齐） | 仓库根 `*.rs` | 未落地 | 见 §5 TODO |
| 对齐实现 | Go（对齐） | 仓库根 `*.go`（`package num`，64 个实现文件） | 已落地 | 用例见 §4.1 |

## 2. 功能清单

**统一口径**：

- **类型分发**：浮点（`[]float32`）走 `x32/`，浮点（`[]float64`）走 `x64/`，两者再下沉到
  `internal/functions/` 的 `Xxx_AVX2_F32/F64` 内核；整型、无符号与其它类型走包内 `__go_*` 回退路径。
  未命中分支的处理按函数分别是**原样返回**（`Abs` 对无符号 / 其它类型）、**panic `ErrUnsupportedType`**
  （`Sqrt` 的非浮点）或 **panic `TypeError(v)`**（`Mean2` 的非数值），逐行见下表备注。
- **无效值（NaN / nil）**：浮点用 `math.NaN()`（`NaN()` / `Float64NaN()` / `Float32NaN()`）；`int32` / `int64`
  的无效值是 `0`；`bool` 的无效值是 `false`；`string` 的无效值是 `"NaN"`（`StringNaN`）。判 NaN 一律用
  `DTypeIsNaN` / `Float32IsNaN` / `Float64IsNaN` / `StringIsNaN`，不直接与 `NaN` 比较。
- **缺失值填充**：`Fill` / `FillNa` 默认**复制**入参再填（`args[0] == true` 时原地改）；`Rolling` 窗口不足
  时**用空切片占位**，不补值；`Shift` 越界位填 `TypeDefault[T]()`。
- **长度不齐**：`Where`、`Align`、`BinaryOperations` 按最长输入对齐，缺位用 `TypeDefault[T]()` 填。
- **舍入**：`Decimal` 默认 2 位，负 `digits` 视为 0，入参为 `NaN` 时按 0 参与舍入。
- **精度与平台**：浮点结果以 `x32` / `x64` 的加速实现为准，`internal/functions` 在 amd64 走 AVX2、arm64 走
  NEON（内核同名 `Xxx_AVX2_F32/F64`），开关关闭或架构不支持时回退 Go 路径（[§5.6](CONTRIBUTING.md#56-向量化内核规范)）。
- **并发**：`SetAvx2Enabled` **非线程安全**，只允许在用例与基准里切换，生产路径不得调用。

| 功能 | python（spec） | c++ | rust | go | zig | java | 状态 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 逐元素绝对值 | — | — | — | `Abs` | — | — | 部分落地（缺 python / c++ / rust） |
| 标量整数绝对值 | — | — | — | `AbsInt` / `AbsInt8` / `AbsInt16` / `AbsInt32` / `AbsInt64` | — | — | 部分落地（缺 python / c++ / rust） |
| 标量浮点绝对值 | — | — | — | `AbsFloat32` / `AbsFloat64` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素平方根 | — | — | — | `Sqrt` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素整数次幂 | — | — | — | `Pow` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素 `log(1+x)` | — | — | — | `Log1p` / `Log1pUnrolled` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素 `e^x-1` | — | — | — | `Expm1` / `Expm1Unrolled` | — | — | 部分落地（缺 python / c++ / rust） |
| 二进制位舍入 | — | — | — | `Decimal` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素加 | — | — | — | `Add` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素减 | — | — | — | `Sub` / `SubInplace` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素乘 | — | — | — | `Mul` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素除 | — | — | — | `Div` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素取大 / 取小 | — | — | — | `Maximum` / `Minimum` | — | — | 部分落地（缺 python / c++ / rust） |
| 变化率 / 净变化率 | — | — | — | `ChangeRate` / `NetChangeRate` | — | — | 部分落地（缺 python / c++ / rust） |
| 一元运算分发内核 | — | — | — | `UnaryOperations` / `UnaryOperations1` / `UnaryOperations2` | — | — | 部分落地（缺 python / c++ / rust） |
| 二元运算分发内核 | — | — | — | `BinaryOperations` / `BinaryOperations2` | — | — | 部分落地（缺 python / c++ / rust） |
| 求和 | — | — | — | `Sum` | — | — | 部分落地（缺 python / c++ / rust） |
| 均值 | — | — | — | `Mean` / `Mean2` | — | — | 部分落地（缺 python / c++ / rust） |
| 中位数 | — | — | — | `Median` | — | — | 部分落地（缺 python / c++ / rust） |
| 最大值 / 最小值 | — | — | — | `Max` / `Max2` / `Min` / `Min2` | — | — | 部分落地（缺 python / c++ / rust） |
| 极值下标 | — | — | — | `ArgMax` / `ArgMax2` / `ArgMin` / `ArgMin2` | — | — | 部分落地（缺 python / c++ / rust） |
| 标准差 / 方差 / 协方差 | — | — | — | `Std` / `Variance` / `Covariance` | — | — | 部分落地（缺 python / c++ / rust） |
| 非零计数 | — | — | — | `Count` | — | — | 部分落地（缺 python / c++ / rust） |
| 全真 / 任一为真 | — | — | — | `All` / `Any` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素逻辑与 / 或 / 非 | — | — | — | `And` / `Or` / `Not` | — | — | 部分落地（缺 python / c++ / rust） |
| 逐元素比较 | — | — | — | `Gt` / `Gte` / `Lt` / `Lte` / `Equal` / `NotEqual` | — | — | 部分落地（缺 python / c++ / rust） |
| 条件选择 | — | — | — | `Where` | — | — | 部分落地（缺 python / c++ / rust） |
| 等差序列 | — | — | — | `Arange` / `Linspace` / `Range` | — | — | 部分落地（缺 python / c++ / rust） |
| 全零 / 全一 | — | — | — | `Zeros` / `Ones` | — | — | 部分落地（缺 python / c++ / rust） |
| 按最长对齐 | — | — | — | `Align` | — | — | 部分落地（缺 python / c++ / rust） |
| 形状 | — | — | — | `Shape` | — | — | 部分落地（缺 python / c++ / rust） |
| 重复 / 重复写入 | — | — | — | `Repeat` / `RepeatInto` | — | — | 部分落地（缺 python / c++ / rust） |
| 一维拼接 | — | — | — | `Concat1D` | — | — | 部分落地（缺 python / c++ / rust） |
| 二维转置 | — | — | — | `Transpose2D` | — | — | 部分落地（缺 python / c++ / rust） |
| 缺失值填充 | — | — | — | `Fill` / `FillNa` | — | — | 部分落地（缺 python / c++ / rust） |
| 累计和 | — | — | — | `CumSum` | — | — | 部分落地（缺 python / c++ / rust） |
| 一阶离散差分 | — | — | — | `Diff` / `V2Diff` | — | — | 部分落地（缺 python / c++ / rust） |
| 索引平移 | — | — | — | `Shift` | — | — | 部分落地（缺 python / c++ / rust） |
| 滑动窗口 | — | — | — | `Rolling` / `RollingApply` | — | — | 部分落地（缺 python / c++ / rust） |
| 滚动窗口对象 | — | — | — | `RollingWindow` / `Any2Window` / `Window.At` / `Window.IsConst` | — | — | 部分落地（缺 python / c++ / rust） |
| 周期序列 | — | — | — | `Periods` / `AnyToPeriod` | — | — | 部分落地（缺 python / c++ / rust） |
| 波峰波谷 | — | — | — | `PeeksAndValleys` / `ArgRelExtrema` | — | — | 部分落地（缺 python / c++ / rust） |
| 两序列交叉 | — | — | — | `Cross` | — | — | 部分落地（缺 python / c++ / rust） |
| 频率解析 / 日期序列 | — | — | — | `ParseFreq` / `DateRange` | — | — | 部分落地（缺 python / c++ / rust） |
| 点积 | — | — | — | `Dot` / `Dot1D` / `Dot2D` / `Dot2D1` / `Dot2D1_v2` / `Dot2D_V1` / `Dot_v1` | — | — | 部分落地（缺 python / c++ / rust） |
| 矩阵求逆 | — | — | — | `Inverse` | — | — | 部分落地（缺 python / c++ / rust） |
| 多项式拟合 / 求值 | — | — | — | `PolyFit` / `PolyVal` | — | — | 部分落地（缺 python / c++ / rust） |
| 直线方程与角度 | — | — | — | `Slope` / `SlopeToDegrees` / `DegreesToSlope` / `CalculateLineEquation` / `TriangleBevel` | — | — | 部分落地（缺 python / c++ / rust） |
| 趋势线拟合 | — | — | — | `FitLine` / `Line.*`（`Analyze` / `Extend` / `Angle` / `Incr` / `X` / `Y` …） | — | — | 部分落地（缺 python / c++ / rust） |
| 点 / 数据点 | — | — | — | `Point` / `NewPoint` / `DataPoint` / `Point.Add` | — | — | 部分落地（缺 python / c++ / rust） |
| 关联趋势 | — | — | — | `LinerTrend` | — | — | 部分落地（缺 python / c++ / rust） |
| 正态分布分位 | — | — | — | `FindPercent` / `FindZScore` | — | — | 部分落地（缺 python / c++ / rust） |
| 置信区间换算 | — | — | — | `ConfidenceIntervalToZscore` / `ZscoreToConfidenceInterval` | — | — | 部分落地（缺 python / c++ / rust） |
| 类型零值 / 无效值 | — | — | — | `TypeDefault` / `NaN` / `DTypeIsNaN` / `Float32NaN` / `Float64NaN` | — | — | 部分落地（缺 python / c++ / rust） |
| `any` 转数值 | — | — | — | `AnyToInt32` / `AnyToInt64` / `AnyToFloat32` / `AnyToFloat64` / `AnyToBool` / `AnyToString` | — | — | 部分落地（缺 python / c++ / rust） |
| 字符串解析 | — | — | — | `ParseBool` / `ParseInt32` / `ParseInt64` / `ParseFloat32` / `ParseFloat64` / `GenericParse` | — | — | 部分落地（缺 python / c++ / rust） |
| 切片类型转换 | — | — | — | `SliceToBool` / `SliceToInts` / `SliceToFloat32` / `SliceToFloat64` / `SliceToString` / `SliceFloatToInt32` / `SliceFloatToInt64` / `Slice2DType` / `DType2Int` | — | — | 部分落地（缺 python / c++ / rust） |
| 布尔与数值互转 | — | — | — | `BoolToInt` / `BoolToInt32` / `BoolToInt64` / `BoolToFloat32` / `BoolToFloat64` / `BoolToString` | — | — | 部分落地（缺 python / c++ / rust） |
| 字符串判空 / 判真 | — | — | — | `IsEmpty` / `StringIsTrue` / `StringIsFalse` / `StringIsNaN` | — | — | 部分落地（缺 python / c++ / rust） |
| 泛型与反射工具 | — | — | — | `AnyToGeneric` / `AnyToSlice` / `Any2DType` / `CheckoutRawType` / `ExtractValueFromPointer` | — | — | 部分落地（缺 python / c++ / rust） |
| 格式化与诊断 | — | — | — | `DefaultFormatter` / `PrintString` / `PanicTrace` / `TypeError` / `ErrUnsupportedType` / `ErrRange` / `ErrInvalidWindow` | — | — | 部分落地（缺 python / c++ / rust） |
| 加速开关 | — | — | — | `GetAvx2Enabled` / `SetAvx2Enabled` | — | — | 部分落地（缺 python / c++ / rust） |

**状态判据**：[CONTRIBUTING.md §6.4](CONTRIBUTING.md#64-落地状态与配套纪律)。四门必做语言
（python / c++ / rust / go）中目前只有 Go 落地，故全表统一为「部分落地（缺 python / c++ / rust）」；
`zig` / `java` 为扩展列，不参与判定。本表**未**逐函数审计 Go 用例覆盖度，包级覆盖率见 §4.2。

### 2.1 两层 API 对照

门面层（本包）与强类型层（`x32/` / `x64/`）是**同一语义的两套入口**：门面层按运行时类型分发，强类型层是
分发的落点，也是加速内核的直接调用者。二者用例各在自己目录，互不替代。

| 层 | 定位 | python（spec） | c++ | rust | go |
| --- | --- | --- | --- | --- | --- |
| 门面 API | 调用方默认入口，泛型 + `any`，跨类型统一签名 | — | — | — | 本包（`package num`） |
| 强类型 API | 需要固定 dtype 或直接控制加速路径时使用 | — | — | — | `x32/` / `x64/`（`package x32` / `package x64`） |

## 3. 依赖

**出向依赖（本模块依赖的东西）**：

| 依赖项 | 版本 / 来源 | 用途 | 适用语言 | 约束与备注 |
| --- | --- | --- | --- | --- |
| `github.com/quant1x/num/x32` | 本仓 `x32/` | `[]float32` 的强类型实现与加速内核调用 | go | 门面层的 f32 落点，`Abs` / `Sqrt` / 算术 / 归约等均经它下沉 |
| `github.com/quant1x/num/x64` | 本仓 `x64/` | `[]float64` 的强类型实现与加速内核调用 | go | 门面层的 f64 落点；`Count` / `Shape` 等也走它 |
| `github.com/quant1x/num/internal/functions` | 本仓 `internal/functions/` | `SetAcceleration` 与 `UseAVX2` 开关 | go | 只经 `builtin.go` 读写加速状态 |
| `github.com/quant1x/num/internal/partial` | 本仓 `internal/partial/` | 部分排序 / 选择（`Median` 使用） | go | 未排序路径的最坏复杂度控制 |
| `github.com/quant1x/num/math32` | 本仓 `math32/` | `float32` 的 `Sqrt` 回退实现 | go | 标准库无 `float32` 版 `Sqrt`，故必须自备 |
| `github.com/quant1x/num/binary` | 本仓 `binary/` | `Repeat` 的批量写入 | go | 逐元素写入过慢，批量复制才达性能预期 |
| `github.com/quant1x/num/vectors` | 本仓 `vectors/` | `Repeat` 的向量化路径 | go | 与 `binary/` 同为 `Repeat` 的性能路径 |
| `math` / `slices` / `reflect` / `errors` / `fmt` / `strconv` / `strings` / `time` | Go 标准库 | 数学函数、切片克隆、类型反射、错误与频率解析 | go | 标准库足够，不引入三方依赖 |
| `github.com/stretchr/testify` | v1.11.1 | 测试断言（仅 `*_test.go`） | go | 只用于用例，不进运行时路径 |

**入向依赖（谁依赖本模块）**：

| 调用方 | 用途 | 备注 |
| --- | --- | --- |
| 仓库外部使用者 | 直接以 `github.com/quant1x/num` 作为门面调用 | 仓库内**无**其它包 import 本包（`x32/` / `x64/` / `internal/` 均不反向依赖），无循环依赖风险 |

## 4. 测试

### 4.1 单元测试

| 语言 | 测试文件 | 用例数 | 运行命令 |
| --- | --- | --- | --- |
| Python（spec） | 未落地 | — | — |
| C++ | 未落地 | — | — |
| Rust | 未落地 | — | — |
| Go | 仓库根 `*_test.go`（49 个文件，与实现同名） | 82 个 `Test` + 46 个 `Benchmark` | `go test .`（实测 `ok github.com/quant1x/num 0.491s`） |

用例与实现同目录同包；`window_test.go` 只有基准、无用例，见 §5 TODO。本包用例**未**按
[§7](CONTRIBUTING.md#7-测试规范) 的 `accel` 两遍跑（`functions.Info().Acceleration` 在 82 个 `Test` 里无引用，
只有 `align` / `repeat` / `diff` 的基准用 `SetAvx2Enabled` 切换），见 §4.6 与 §5 TODO。

### 4.2 覆盖率

口径为 Go 语句覆盖，只在 Go 内比较，无跨语言可比性。

| 语言 | 工具 | 门槛 | 运行命令 |
| --- | --- | --- | --- |
| Go | `go test -cover` | 仓库未设门槛 | `go test -cover .` |

**实测**：`42.1% of statements`（`go test -count=1 -cover .`，HEAD `04e2d1f` + 工作树未提交改动）。
覆盖率缺口与未设门槛一并登记在 §5 TODO。

### 4.3 重复度

本模块不涉及：仓库尚未引入重复度检测工具（无 `.jscpd.json` 等配置），按 §6.3 保留本节并说明。引入后须
在此登记工具、阈值与扫描范围。

### 4.4 静态检查

| 语言 | 工具 | 运行命令 | 备注 |
| --- | --- | --- | --- |
| Go | `gofmt -l` ＋ `go vet` | `gofmt -l .`、`go vet .` | `gofmt -l` 实测无输出；`go vet .` **未通过**，报 5 处 `unreachable code`，见 §5 TODO |

`go vet .` 当前输出：`type.go:211`、`type.go:296`、`type_float64.go:134`、`type_integer.go:75`、
`type_strings.go:188`。

### 4.5 基准测试

| 语言 | 基准路径 | 指标 | 基线（含环境与 commit） | 运行命令 |
| --- | --- | --- | --- | --- |
| Go | 仓库根 `*_test.go` 内的 46 个 `Benchmark*`（`align` / `diff` / `repeat` / `shift` / 算术等） | `ns/op`、`B/op`、`allocs/op` | 未登记（见 §5 TODO） | `go test -run=NONE -bench=. -benchmem .` |

基准按 `avx2` / `noavx2` 分组命名（如 `BenchmarkRepeat_release_avx2_float64` /
`BenchmarkRepeat_release_noavx2_float64`），切换时用 `SetAvx2Enabled` 并在 `defer` 里复位。无基线值时
只看同机同提交内的相对差异。

### 4.6 契约 / 一致性测试

本模块不涉及跨语言互操作（只有 Go 落地）；跨**架构**一致性由 `x32/` / `x64/` 与 `internal/functions/`
的用例承担，不在本包：

| 场景 | 测试位置 | 断言内容 | 运行命令 |
| --- | --- | --- | --- |
| AVX2（amd64） / NEON（arm64） 与 Go 路径等价 | `x32/*_test.go`、`x64/*_test.go`、`internal/functions/*_test.go` | 逐元素相等；长度 1、块大小 ±1；`NaN` / `±Inf` / `±0` | `go test ./x32/ ./x64/ ./internal/functions/` |
| 门面层两条路径 | 本包用例（**缺**） | `UseAVX2 = false` 与 `true` 下结果一致 | 见 §5 TODO |

## 5. TODO

| 待办项 | 缺什么 | 卡在哪 | 关联章节 | 关联 issue | 优先级 |
| --- | --- | --- | --- | --- | --- |
| 补 Python spec | 整包无 `.py`，语义准绳缺失，§2 全部 62 行均缺 python 列 | 无技术阻塞，待排期 | §1 / §2 | 待建单 | 高 |
| 补 C++ 生产实现 | 整包无 `.h` / `.cpp` | 无技术阻塞，待排期 | §1 / §2 | 待建单 | 高 |
| 补 Rust 对齐实现 | 整包无 `.rs` | 无技术阻塞，待排期 | §1 / §2 | 待建单 | 高 |
| 清 `go vet` 报错 | 5 处 `unreachable code`：`type.go:211`、`type.go:296`、`type_float64.go:134`、`type_integer.go:75`、`type_strings.go:188` | 无技术阻塞；需逐个确认是死分支还是 `panic` 之后的正常不可达 | §4.4 | 待建单 | 中 |
| 根包用例按 `accel` 两遍跑 | 82 个 `Test` 均未引用 `functions.Info().Acceleration`，`UseAVX2 = false` 分支在根包层无用例（`x32` / `x64` 的用例已按 §7 跑两遍） | 无技术阻塞，待排期 | §4.1 / §4.6 | 待建单 | 中 |
| 提升覆盖率并设门槛 | 语句覆盖 42.1%，仓库未设门槛；`window_test.go` 只有基准、无 `Test` | 无技术阻塞，待排期 | §4.2 | 待建单 | 中 |
| 登记基准基线 | 46 个 `Benchmark*` 均无基线值与环境说明 | 需固定一台基准机，并写明 CPU / OS / commit | §4.5 | 待建单 | 中 |
| 引入重复度检测 | 仓库无重复度工具与阈值配置 | 无技术阻塞，待排期 | §4.3 | 待建单 | 低 |
| 逐功能 `.md` 是否拆分 | 当前整包共用 `num.md`；若后续按语义拆出 `abs.md` / `dot.md` 等，须同步删减 §2 对应行 | 待评审确定粒度 | §2 | 待建单 | 低 |
