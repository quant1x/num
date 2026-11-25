# num：门面数值 API —— 按泛型类型分发到 `x32` / `x64` 与加速内核

**边界**：本包（Go 侧 `github.com/quant1x/num`，`package num`；C++ 侧 `namespace num`）是仓库的**门面层**：把同一语义的强类型实现（`x32/`、`x64/`）与加速内核（`internal/functions/`）收在一条泛型签名后面，负责**类型分发**与**回退选择**。它**不**定义新的数值语义，也**不**承载 IO、序列化、时间序列数据容器；算法本体在 `x32/` / `x64/` / `internal/functions/` / `internal/partial/` / `math32/` / `binary/` / `vectors/` 各自成包。

C++ 侧只有**一层**：`Vector<T>`（`xt::xtensor<T, 1>`）用模板覆盖 f32 / f64 两条宽度，向量化由 xtensor 的 xsimd 后端承担，参考实现与加速实现是同一份代码，故不设 `x32` / `x64` / `internal/functions` 的对应文件（见 §2.1）。

**文件组织**：逐语义一文件平铺在仓库根，Go 侧 `abs.go` / `sum.go` / `dot.go` / `rolling.go` …，C++ 侧同名不同后缀（`abs.h` / `sum.h` …，不用 `include/` 与 `src/` 两级目录）；本文件是整包的功能级文档，§2 一行一个功能对应到具体符号；本包无 `num.go` / `num.h`，文档以包名命名。

## 1. 路线图

**当前位置与下一步**：Go 侧门面 API 已落地（64 个实现文件 + 49 个测试文件），覆盖算术、归约、逻辑比较、数组构造、滑动窗口、线代与几何、正态分布表、类型转换九类。C++ 侧已立起构建链路并落地**第一批算子**（算术、归约、逻辑比较、构造与整形、舍入、指数对数、填充、滑动窗口、线代与几何，§2 中 59 个功能行有 C++ 实现），余下只有反射族（`AnyToGeneric` 的反射部分 / `AnyToSlice` / `CheckoutRawType` / `ExtractValueFromPointer`）与 Go 侧本身就是空实现的 `PeeksAndValleys`。Python spec 已落地 113 个语义（163 个 npz 样本全过，`NAME-MAP` 的 authority 逐行标 python）；Rust 已落地 fixture 覆盖的 N1 ~ N4、N6 语义（根 `Cargo.toml`，`*.rs` 与 `.go` 平铺同目录），按 [CONTRIBUTING.md §6.4](CONTRIBUTING.md#64-落地状态与配套纪律) 整包仍判为**部分落地（Python 覆盖 fixture 子集、Rust 仅覆盖子集）**。下一个里程碑：按 §5 TODO 的顺序补齐 C++ 剩余算子、扩 Rust 覆盖面；Python 侧先给剩余无语义的算子补 fixture 再继续切 `authority`。

**临时基准**：本仓对标 NumPy，算法与数值口径以它为准。`NAME-MAP` 标 `authority=python` 的语义以 Python 实现为准；未标 python 的仍以本文件 + NumPy 语义 + Go 用例作临时基准（[§9.3](CONTRIBUTING.md#93-spec--实现--验证)）。Go 是最早落地的一份实现，可作可执行对照，但不是 spec。

| 阶段 | 语言（角色） | 载体 / 文件 | 状态 | 说明 |
| --- | --- | --- | --- | --- |
| spec 定语义 | Python（spec） | 仓库根 `*.py`（与 `.go` 同名） | 未落地 | 语义源头，见 §5 TODO |
| 生产实现 | C++（生产） | 仓库根 `*.h` / `*.cpp`（与 `.go` 同名） | 部分落地 | 59 个功能行，见 §2 / §5 TODO |
| 对齐实现 | Rust（对齐） | 根 `Cargo.toml` + 根目录 `*.rs`（与 `.go` 同名平铺） | 部分落地 | N1 ~ N4、N6 的 fixture 语义，见 §5 TODO |
| 对齐实现 | Go（对齐） | 仓库根 `*.go`（`package num`，64 个实现文件） | 已落地 | 用例见 §4.1 |

## 2. 功能清单

**统一口径**：

- **类型分发（Go）**：浮点（`[]float32`）走 `x32/`，浮点（`[]float64`）走 `x64/`，两者再下沉到 `internal/functions/` 的 `Xxx_AVX2_F32/F64` 内核；整型、无符号与其它类型走包内 `__go_*` 回退路径。未命中分支的处理按函数分别是**原样返回**（`Abs` 对无符号 / 其它类型）、**panic `ErrUnsupportedType`**（`Sqrt` 的非浮点）或 **panic `TypeError(v)`**（`Mean2` 的非数值），逐行见下表备注。
- **类型分发（C++）**：模板在**编译期**完成同一件事，非浮点类型在编译期即被 `static_assert` 拒绝（`sqrt` / `stddev`），没有运行期分支，也没有加速开关。
- **无效值（NaN / nil）**：浮点用 `math.NaN()`（`NaN()` / `Float64NaN()` / `Float32NaN()`）；`int32` / `int64` 的无效值是 `0`；`bool` 的无效值是 `false`；`string` 的无效值是 `"NaN"`（`StringNaN`）。判 NaN 一律用 `DTypeIsNaN` / `Float32IsNaN` / `Float64IsNaN` / `StringIsNaN`，不直接与 `NaN` 比较。C++ 侧同口径收在 `type.h`：`type_default<T>()`（浮点 `NaN`、整型 `0`、`bool` `false`）与 `is_nan(DType)`。
- **缺失值填充**：`Fill` / `FillNa` 默认**复制**入参再填（`args[0] == true` 时原地改）；`Rolling` 窗口不足时**用空切片占位**，不补值；`Shift` 越界位填 `TypeDefault[T]()`（C++ 侧 `shift` 同此口径）；`n == 0` 返回**原序列**（[#14](https://codeberg.org/quant1x/num/issues/14) 定案），`Diff` 随之全 0；窗口参数只接受 `int`（常量位移）与 `[]DType`（逐位置位移），其余类型 panic `ErrInvalidWindow`（[#15](https://codeberg.org/quant1x/num/issues/15) 定案），C++ 侧的非 `int` 算术常量由 `= delete` 重载在编译期拒绝。
- **长度不齐**：`Align`、`Where`、`Maximum` / `Minimum`、`Gt` / `Gte` / `Lt` / `Lte` / `And` / `Or` 按**最长**输入对齐；`Add` / `Sub` / `Mul` / `Div` 与 `Equal` / `NotEqual` 按**第一个操作数**的长度对齐（第二个短则补默认值、长则截断）。两侧都用 `TypeDefault[T]()` 填充，C++ 照此实现，不要统一成一种。
- **舍入**：`Decimal` 默认 2 位，负 `digits` 视为 0，入参为 `NaN` 时按 0 参与舍入；进位是**半程远离零**（`decimal(0.5, 0) = 1`、`decimal(-0.5, 0) = -1`），不是银行家舍入。
- **精度与平台**：Go 侧浮点结果以 `x32` / `x64` 的加速实现为准，`internal/functions` 在 amd64 走 AVX2、arm64 走 NEON（内核同名 `Xxx_AVX2_F32/F64`），开关关闭或架构不支持时回退 Go 路径（[§5.6](CONTRIBUTING.md#56-向量化内核规范)）。C++ 侧由 xtensor 表达式经 xsimd 后端向量化，指令集按编译器基线下发，**不**加 `/arch:AVX2`（见 §5 TODO）。
- **并发**：`SetAvx2Enabled` **非线程安全**，只允许在用例与基准里切换，生产路径不得调用。C++ 侧无此开关。
- **失败约定（C++）**：Go 的 `panic` 与 `error` 在 C++ 侧换成异常，类型与文案一一对应（`num::Error` 为基类，另有 `UnsupportedType` / `TypeError` / `RangeError` / `InvalidWindow` / `LengthMismatch` / `EmptyInput` / `DomainError`），调用方按类型 `catch`，不解析文案。
- **C++ 侧的三个语言机制差异**（不是语义差异）：① 多返回值用 `std::tuple`（顺序与 Go 一致），需要保留"半成品 + 失败位"的用 `std::pair`（如 `fit_line`），只有"有 / 没有"的用 `std::optional`（如 `Line::intersect`）；② 结果元素是结构体时容器用 `std::vector`（`cross` 的 `LinerTrend`），xtensor 容器只装数值；③ 比较器按 `template` 形参收（`arg_rel_extrema`），Go 的 `func(a, b float64) bool` 直译成 lambda 即可。
- **命名偏差（C++ 语言约束，非语义差异）**：`Std` → `stddev`（`num::std` 会把同命名空间内所有 `std::` 限定名遮蔽）；`And` / `Or` / `Not` → `logic_and` / `logic_or` / `logic_not`（`and` / `or` / `not` 是 C++ 关键字）。

| 功能 | python（spec） | c++ | rust | go | zig | java | 状态 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 逐元素绝对值 | `abs.abs` | `abs` | `abs` | `Abs` | — | — | 已落地（zig / java 为扩展） |
| 标量整数绝对值 | — | `abs_int` / `abs_int8` / `abs_int16` / `abs_int32` / `abs_int64` | — | `AbsInt` / `AbsInt8` / `AbsInt16` / `AbsInt32` / `AbsInt64` | — | — | 部分落地（缺 python / rust） |
| 标量浮点绝对值 | — | `abs_float32` / `abs_float64` | — | `AbsFloat32` / `AbsFloat64` | — | — | 部分落地（缺 python / rust） |
| 逐元素平方根 | `sqrt.sqrt` | `sqrt` | `sqrt` | `Sqrt` | — | — | 已落地（zig / java 为扩展） |
| 逐元素整数次幂 | `x32.vek.pow_integer` | `pow` | `pow_integer` | `Pow` | — | — | 已落地（zig / java 为扩展） |
| 逐元素 `log(1+x)` | `log.log1p` | `log1p` | `log1p` | `Log1p` / `Log1pUnrolled` | — | — | 部分落地（子集；zig / java 为扩展） |
| 逐元素 `e^x-1` | `log.expm1` | `expm1` | `expm1` | `Expm1` / `Expm1Unrolled` | — | — | 部分落地（子集；zig / java 为扩展） |
| 二进制位舍入 | `decimals.decimal` | `decimal` | `decimal` | `Decimal` | — | — | 已落地（zig / java 为扩展） |
| 逐元素加 | `arithmetics_add.add` | `add` / `add_inplace` | `add` | `Add` | — | — | 已落地（zig / java 为扩展） |
| 逐元素减 | `arithmetics_sub.sub` | `sub` / `sub_inplace` | `sub` | `Sub` / `SubInplace` | — | — | 部分落地（子集；zig / java 为扩展） |
| 逐元素乘 | `arithmetics_mul.mul` | `mul` / `mul_inplace` | `mul` | `Mul` | — | — | 已落地（zig / java 为扩展） |
| 逐元素除 | `arithmetics_div.div` | `div` / `div_inplace` | `div` | `Div` | — | — | 已落地（zig / java 为扩展） |
| 逐元素取大 / 取小 | `maximum.maximum` / `minimum.minimum` | `maximum` / `minimum` | `maximum` / `minimum` | `Maximum` / `Minimum` | — | — | 已落地（zig / java 为扩展） |
| 变化率 / 净变化率 | `decimals.change_rate` / `decimals.net_change_rate` | `change_rate` / `net_change_rate` | `change_rate` / `net_change_rate` | `ChangeRate` / `NetChangeRate` | — | — | 已落地（zig / java 为扩展） |
| 一元运算分发内核 | — | —（模板在编译期分发，不设此层） | — | `UnaryOperations` / `UnaryOperations1` / `UnaryOperations2` | — | — | 部分落地（缺 python / rust） |
| 二元运算分发内核 | — | —（模板在编译期分发，不设此层） | — | `BinaryOperations` / `BinaryOperations2` | — | — | 部分落地（缺 python / rust） |
| 求和 | `sum.sum` | `sum` | `sum` | `Sum` | — | — | 已落地（zig / java 为扩展） |
| 均值 | `mean.mean` | `mean` | `mean` | `Mean` / `Mean2` | — | — | 部分落地（子集；zig / java 为扩展） |
| 中位数 | `median.median` | `median` | `median` | `Median` | — | — | 已落地（zig / java 为扩展） |
| 最大值 / 最小值 | `max.max` / `min.min` | `max` / `min` | `max` / `min` | `Max` / `Max2` / `Min` / `Min2` | — | — | 部分落地（子集；zig / java 为扩展） |
| 极值下标 | `argmax.argmax` / `argmin.argmin` | `argmax` / `argmin` | `argmax` / `argmin` | `ArgMax` / `ArgMax2` / `ArgMin` / `ArgMin2` | — | — | 部分落地（子集；zig / java 为扩展） |
| 标准差 / 方差 / 协方差 | `stddev.std` / `covariance.variance` / `covariance.covariance` | `stddev` / `variance` / `covariance` | `std` / `variance` / `covariance` | `Std` / `Variance` / `Covariance` | — | — | 已落地（zig / java 为扩展） |
| 非零计数 | `count.count` | `count` | `count, count_true` | `Count` | — | — | 已落地（zig / java 为扩展） |
| 全真 / 任一为真 | `logic_all.all` / `logic_any.any` | `all` / `any` | `all` / `any` | `All` / `Any` | — | — | 已落地（zig / java 为扩展） |
| 逐元素逻辑与 / 或 / 非 | `logic.logic_and` / `logic.logic_or` / `logic.logic_not` | `logic_and` / `logic_or` / `logic_not` | `logic_and` / `logic_or` / `logic_not` | `And` / `Or` / `Not` | — | — | 已落地（zig / java 为扩展） |
| 逐元素比较 | `logic.gt` / `logic.gte` / `logic.lt` / `logic.lte` / `logic_equal.equal` / `logic_equal.not_equal` | `gt` / `gte` / `lt` / `lte` / `equal` / `not_equal` | `greater` / `greater_equal` / `less` / `less_equal` / `equal` / `not_equal` | `Gt` / `Gte` / `Lt` / `Lte` / `Equal` / `NotEqual` | — | — | 已落地（zig / java 为扩展） |
| 条件选择 | `where.where` | `where` | `where_select` | `Where` | — | — | 已落地（zig / java 为扩展） |
| 等差序列 | `range.arange` / `range.linspace` / `range.range` | `arange` / `linspace` / `range` | `arange_f64, arange_i64` / `linspace` / `range_i64` | `Arange` / `Linspace` / `Range` | — | — | 已落地（zig / java 为扩展） |
| 全零 / 全一 | `zeros.zeros` / `zeros.ones` | `zeros` / `ones` | `zeros` / `ones` | `Zeros` / `Ones` | — | — | 已落地（zig / java 为扩展） |
| 按最长对齐 | `align.align` | `align` | `align` | `Align` | — | — | 已落地（zig / java 为扩展） |
| 形状 | `shape.shape` | `shape` | `shape, shape_2d` | `Shape` | — | — | 已落地（zig / java 为扩展） |
| 重复 / 重复写入 | `repeat.repeat` | `repeat` / `repeat_into` | `repeat` | `Repeat` / `RepeatInto` | — | — | 部分落地（子集；zig / java 为扩展） |
| 一维拼接 | `matrix.concat1d` | `concat1d` | `concat1d` | `Concat1D` | — | — | 已落地（zig / java 为扩展） |
| 二维转置 | `matrix.transpose2d` | `transpose2d` | `transpose2d` | `Transpose2D` | — | — | 已落地（zig / java 为扩展） |
| 缺失值填充 | `fillna.fill` / `fillna.fill_na` | `fill` / `fill_inplace` / `fill_na` / `fill_na_inplace` | `fill` / `fill_na` | `Fill` / `FillNa` | — | — | 已落地（zig / java 为扩展） |
| 累计和 | `cumsum.cumsum` | `cumsum` / `cumsum_inplace` | `cumsum` | `CumSum` | — | — | 已落地（zig / java 为扩展） |
| 一阶离散差分 | `diff.diff` | `diff` | `diff` | `Diff` / `V2Diff` | — | — | 部分落地（子集；zig / java 为扩展） |
| 索引平移 | `shift.shift` | `shift` | `shift` | `Shift` | — | — | 已落地（zig / java 为扩展） |
| 滑动窗口 | `rolling.rolling` / `rolling.rolling_apply` | `rolling` / `rolling_apply` | `rolling` / `rolling_apply` | `Rolling` / `RollingApply` | — | — | 已落地（zig / java 为扩展） |
| 滚动窗口对象 | — | `Window` / `rolling_window`（`at` / `is_const`） | — | `RollingWindow` / `Any2Window` / `Window.At` / `Window.IsConst` | — | — | 部分落地（缺 python / rust） |
| 周期序列 | — | `Periods`（`at` / `is_const`） | — | `Periods` / `AnyToPeriod` | — | — | 部分落地（缺 python / rust） |
| 波峰波谷 | — | `arg_rel_extrema`（`peeks_and_valleys` 在 Go 侧是空实现，未移植） | — | `PeeksAndValleys` / `ArgRelExtrema` | — | — | 部分落地（缺 python / rust） |
| 两序列交叉 | `cross.cross` | `cross` / `LinerTrend` / `TrendState` | `cross` | `Cross` | — | — | 已落地（zig / java 为扩展） |
| 频率解析 / 日期序列 | `freq.parse_freq` / `freq.date_range` | `parse_freq` / `date_range` | `parse_freq` / `date_range` | `ParseFreq` / `DateRange` | — | — | 已落地（zig / java 为扩展） |
| 点积 | `dot.dot1d` / `matrix.dot2d` / `matrix.dot2d1` | `dot1d` / `matmul` / `dot2d` / `dot2d_v1` / `dot2d1` | `dot1d` / `dot2d` / `dot2d1` | `Dot` / `Dot1D` / `Dot2D` / `Dot2D1` / `Dot2D1_v2` / `Dot2D_V1` / `Dot_v1` | — | — | 部分落地（子集；zig / java 为扩展） |
| 矩阵求逆 | `matrix.inverse` | `inverse` | `inverse` | `Inverse` | — | — | 已落地（zig / java 为扩展） |
| 多项式拟合 / 求值 | `polynomial.poly_fit` / `polynomial.poly_val` | `poly_fit` / `poly_val` | `poly_fit` / `poly_val` | `PolyFit` / `PolyVal` | — | — | 已落地（zig / java 为扩展） |
| 直线方程与角度 | `slope.slope` / `slope.slope_to_degrees` / `slope.degrees_to_slope` / `line.calculate_line_equation` / `slope.triangle_bevel` | `slope` / `triangle_bevel` / `degrees_to_slope` / `slope_to_degrees` / `calculate_line_equation` / `Line`（`radian` / `degrees` / `angle` / `equation` / `distance` / `horizontal_distance` / `vertical_distance` / `shortest_distance` / `parallel_line` / `symmetric_parallel_line` / `x` / `y` / `incr` / `intersect`） | `slope` / `slope_to_degrees` / `degrees_to_slope` / `calculate_line_equation` / `triangle_bevel` | `Slope` / `SlopeToDegrees` / `DegreesToSlope` / `CalculateLineEquation` / `TriangleBevel` | — | — | 已落地（zig / java 为扩展） |
| 趋势线拟合 | `line.fit_line` | `fit_line`（`Line::extend` / `Line::analyze`） | `fit_line` | `FitLine` / `Line.*`（`Analyze` / `Extend` / `Angle` / `Incr` / `X` / `Y` …） | — | — | 部分落地（子集；zig / java 为扩展） |
| 点 / 数据点 | — | `Point` / `DataPoint` / `new_point`（`add` / `to_data_point` / `to_point`） | — | `Point` / `NewPoint` / `DataPoint` / `Point.Add` | — | — | 部分落地（缺 python / rust） |
| 关联趋势 | — | `LinerTrend`（`x` / `state`） | — | `LinerTrend` | — | — | 部分落地（缺 python / rust） |
| 正态分布分位 | `ztable.find_percent` / `ztable.find_zscore` | `find_percent` / `find_zscore` | — | `FindPercent` / `FindZScore` | — | — | 部分落地（缺 rust） |
| 置信区间换算 | `ztable.confidence_interval_to_zscore` / `ztable.zscore_to_confidence_interval` | `confidence_interval_to_zscore` / `zscore_to_confidence_interval` | — | `ConfidenceIntervalToZscore` / `ZscoreToConfidenceInterval` | — | — | 部分落地（缺 rust） |
| 类型零值 / 无效值 | `type_dtype.type_default` / `decimals.is_nan` | `type_default` / `nan` / `is_nan` | `type_default` / `is_nan` | `TypeDefault` / `NaN` / `DTypeIsNaN` / `Float32NaN` / `Float64NaN` | — | — | 部分落地（子集；zig / java 为扩展） |
| `any` 转数值 | — | `to_bool` / `to_string`（值重载，无 `any` 分派） | — | `AnyToInt32` / `AnyToInt64` / `AnyToFloat32` / `AnyToFloat64` / `AnyToBool` / `AnyToString` | — | — | 部分落地（缺 python / rust） |
| 字符串解析 | `type_generic.parse_bool, type_generic.parse_int32, type_generic.parse_int64, type_generic.parse_float32, type_generic.parse_float64` | `parse_bool` / `parse_int32` / `parse_int64` / `parse_float32` / `parse_float64` / `generic_parse`（目标类型由模板形参给出，`GenericParse[T]` 的对位） | — | `ParseBool` / `ParseInt32` / `ParseInt64` / `ParseFloat32` / `ParseFloat64` / `GenericParse` | — | — | 部分落地（缺 rust） |
| 切片类型转换 | `type_bool.slice_to_bool` / `type_integer.slice_to_int` / `type_float32.slice_to_float32` / `x32.vek.slice_to_float64` / `x32.vek.slice_float_to_int` / `type_dtype.dtype2int` | `slice_to_bool` / `slice_to_string` / `slice_to_ints` / `slice_to_float32` / `slice_to_float64` / `slice_float_to_int32` / `slice_float_to_int64`（`Slice2DType` / `DType2Int` 分别是 `slice_to_float64` / `slice_float_to_int32` 的同义包装，按「一个功能一个导出符号」不另设） | `floats_to_int32, floats_to_int64` / `dtype2int` | `SliceToBool` / `SliceToInts` / `SliceToFloat32` / `SliceToFloat64` / `SliceToString` / `SliceFloatToInt32` / `SliceFloatToInt64` / `Slice2DType` / `DType2Int` | — | — | 部分落地（子集；zig / java 为扩展） |
| 布尔与数值互转 | `type_bool.bool_to_int64, type_bool.bool_to_float64` | `bool_to_int` / `bool_to_int32` / `bool_to_int64` / `bool_to_float32` / `bool_to_float64` / `bool_to_string` | `bool_to_i64, bool_to_f64, bools_to_f64` | `BoolToInt` / `BoolToInt32` / `BoolToInt64` / `BoolToFloat32` / `BoolToFloat64` / `BoolToString` | — | — | 部分落地（子集；zig / java 为扩展） |
| 字符串判空 / 判真 | `type_strings.is_empty, type_strings.string_is_true, type_strings.string_is_false, type_strings.string_is_nan` | `is_empty` / `string_is_true` / `string_is_false` / `string_is_nan` | — | `IsEmpty` / `StringIsTrue` / `StringIsFalse` / `StringIsNaN` | — | — | 部分落地（缺 rust） |
| 泛型与反射工具 | — | `generic_parse`（`GenericParse` 的字符串入口；`AnyToGeneric` 的数值入口由 `slice_to_*` / `to_bool` / `to_string` 覆盖，反射部分不设） | — | `AnyToGeneric` / `AnyToSlice` / `Any2DType` / `CheckoutRawType` / `ExtractValueFromPointer` | — | — | 部分落地（缺 python / rust） |
| 格式化与诊断 | — | 异常族 `Error` / `UnsupportedType` / `TypeError` / `RangeError` / `InvalidWindow` / `LengthMismatch` / `EmptyInput` / `DomainError` | — | `DefaultFormatter` / `PrintString` / `PanicTrace` / `TypeError` / `ErrUnsupportedType` / `ErrRange` / `ErrInvalidWindow` | — | — | 部分落地（缺 python / rust） |
| 加速开关 | — | —（无运行期开关） | — | `GetAvx2Enabled` / `SetAvx2Enabled` | — | — | 部分落地（缺 python / c++ / rust） |

**状态判据**：[CONTRIBUTING.md §6.4](CONTRIBUTING.md#64-落地状态与配套纪律)。四门必做语言（python / c++ / rust / go）中 Go 全落地、C++ 落地第一批 59 个功能行、Rust 落地 fixture 覆盖的 N1 ~ N4 / N6 子集；本表按**全量 API** 记状态，故带 c++ 符号的行判「部分落地（缺 python / rust）」，其余判「部分落地（缺 python / c++ / rust）」，不因 Rust 子集改判。`zig` / `java` 为扩展列，不参与判定。本表**未**逐函数审计用例覆盖度，包级覆盖率见 §4.2。

### 2.1 两层 API 对照

门面层（本包）与强类型层（`x32/` / `x64/`）是**同一语义的两套入口**：门面层按运行时类型分发，强类型层是分发的落点，也是加速内核的直接调用者。二者用例各在自己目录，互不替代。C++ 侧只有一层 —— 模板同时覆盖 f32 / f64 两条宽度，强类型层的存在理由（避开运行期分派成本）由编译期实例化天然满足。

| 层 | 定位 | python（spec） | c++ | rust | go |
| --- | --- | --- | --- | --- | --- |
| 门面 API | 调用方默认入口，泛型 + `any`，跨类型统一签名 | — | 本包（`namespace num`，仓库根 `*.h` / `*.cpp`） | 根 crate（snake_case，无运行期分派） | 本包（`package num`） |
| 强类型 API | 需要固定 dtype 或直接控制加速路径时使用 | — | —（模板已覆盖两条宽度，不另设层） | —（泛型 `Float` 覆盖两条宽度） | `x32/` / `x64/`（`package x32` / `package x64`） |

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
| `xtensor` | conan `xtensor/0.27.1` | `xtensor` / `xarray` 容器与逐元素表达式 | c++ | 仓库根 `conanfile.py` 固化版本；`Vector<T>` = `xt::xtensor<T, 1>` |
| `xsimd` | conan `xsimd/13.2.0`（随 `xtensor`） | 表达式层的 SIMD 后端 | c++ | 由 `xtensor/*:xsimd=True` 打开，编译期下发 `XTENSOR_USE_XSIMD` |
| `gtest` | conan `gtest/1.18.0`（`test_requires`） | C++ 用例断言 | c++ | 只用于 `*_test.cpp`，不进库产物 |
| 标准库 `<cmath>` / `<cstdint>` / `<stdexcept>` | C++20 | 数学函数、定宽标量、异常 | c++ | 除 xtensor 外不引入第三方运行时依赖 |
| `github.com/stretchr/testify` | v1.11.1 | 测试断言（仅 `*_test.go`） | go | 只用于用例，不进运行时路径 |

**入向依赖（谁依赖本模块）**：

| 调用方 | 用途 | 备注 |
| --- | --- | --- |
| 仓库外部使用者 | 直接以 `github.com/quant1x/num`（Go）或 `num_core`（C++）作为门面调用 | 仓库内**无**其它包 import 本包（`x32/` / `x64/` / `internal/` 均不反向依赖），无循环依赖风险 |

## 4. 测试

### 4.1 单元测试

| 语言 | 测试文件 | 用例数 | 运行命令 |
| --- | --- | --- | --- |
| Python（spec） | 未落地 | — | — |
| C++ | 仓库根 `*_test.cpp`（23 个文件，与语义族同名） | 221 个 `TEST` | `ctest --preset release-static`（实测 `100% tests passed out of 221`，Debug 预设同数） |
| Rust | `tests/conform.rs`（读 `tests/fixtures`） | 160 个 fixture | `cargo test` |
| Go | 仓库根 `*_test.go`（49 个文件，与实现同名） | 84 个 `Test` + 46 个 `Benchmark` | `go test .`（实测 `ok github.com/quant1x/num 0.491s`） |

用例与实现同目录；`window_test.go` 只有基准、无用例，见 §5 TODO。本包 Go 用例**未**按 [§7](CONTRIBUTING.md#7-测试规范) 的 `accel` 两遍跑（`functions.Info().Acceleration` 在 84 个 `Test` 里无引用，只有 `align` / `repeat` / `diff` 的基准用 `SetAvx2Enabled` 切换），见 §4.6 与 §5 TODO。

C++ 用例按**语义族**合并成 23 个文件（`arithmetics_test.cpp` / `aggregates_test.cpp` / `logic_test.cpp` / `construct_test.cpp` / `ndarray_test.cpp` / `decimals_test.cpp` / `log_test.cpp` / `fillna_test.cpp` / `dot_test.cpp` / `cross_test.cpp` / `peak_finding_test.cpp` / `slope_test.cpp` / `rolling_test.cpp` / `median_test.cpp` / `freq_test.cpp` / `shape_test.cpp` / `matrix_test.cpp` / `polynomial_test.cpp` / `ztable_test.cpp` / `type_strings_test.cpp` / `type_bool_test.cpp` / `type_integer_test.cpp` / `type_generic_test.cpp`），每个 `TEST` 由 `gtest_discover_tests` 注册成一条 ctest 用例；覆盖 f32 / f64 两条宽度、空输入、长度不齐、`NaN` / `±Inf` / `±0`。C++ 侧无 `accel` 两遍跑可言（没有运行期加速开关）。

### 4.2 覆盖率

口径为 Go 语句覆盖，只在 Go 内比较，无跨语言可比性；C++ 侧尚未接入覆盖率工具。

| 语言 | 工具 | 门槛 | 运行命令 |
| --- | --- | --- | --- |
| Go | `go test -cover` | 仓库未设门槛 | `go test -cover .` |
| C++ | 未接入（见 §5 TODO） | 仓库未设门槛 | — |

**实测**：`42.1% of statements`（`go test -count=1 -cover .`，HEAD `04e2d1f` + 工作树未提交改动）。覆盖率缺口与未设门槛一并登记在 §5 TODO。

### 4.3 重复度

本模块不涉及：仓库尚未引入重复度检测工具（无 `.jscpd.json` 等配置），按 §6.3 保留本节并说明。引入后须在此登记工具、阈值与扫描范围。

### 4.4 静态检查

| 语言 | 工具 | 运行命令 | 备注 |
| --- | --- | --- | --- |
| Go | `gofmt -l` ＋ `go vet` | `gofmt -l .`、`go vet .` | `gofmt -l` 实测无输出；`go vet` 支持面已无告警 |
| C++ | MSVC `/W4`（构建期即检查） | `cmake --build --preset release-static` | 实测 0 警告；`/utf-8` 必开（无 BOM 的 UTF-8 源文件在中文 Windows 上会被按 GBK 误读）；59 个根头文件另各自单独成 TU 编译（`cl /Zs /W4`），实测 0 error / 0 warning —— 没有"靠别的头间接带入 include"的隐患 |
| Rust | `cargo clippy` ＋ `cargo fmt --check` | `cargo clippy --all-targets -- -D warnings`、`cargo fmt --check` | 实测 0 告警 |


### 4.5 基准测试

| 语言 | 基准路径 | 指标 | 基线（含环境与 commit） | 运行命令 |
| --- | --- | --- | --- | --- |
| Go | 仓库根 `*_test.go` 内的 46 个 `Benchmark*`（`align` / `diff` / `repeat` / `shift` / 算术等） | `ns/op`、`B/op`、`allocs/op` | 未登记（见 §5 TODO） | `go test -run=NONE -bench=. -benchmem .` |
| C++ | 未落地（见 §5 TODO） | `ns/op`（Google Benchmark） | 未登记 | — |

基准按 `avx2` / `noavx2` 分组命名（如 `BenchmarkRepeat_release_avx2_float64` / `BenchmarkRepeat_release_noavx2_float64`），切换时用 `SetAvx2Enabled` 并在 `defer` 里复位。无基线值时只看同机同提交内的相对差异。

### 4.6 契约 / 一致性测试

跨**架构**一致性由 `x32/` / `x64/` 与 `internal/functions/` 的用例承担，不在本包；跨**语言**一致性由 `tests/fixtures/` 的 npz 样本承担 —— Go 有 `tests/fixtures/conform`、Rust 有根 crate 的 `tests/conform.rs`、C++ 有 `tests/fixtures/conform_cpp`（160 个样本全量通过）：

| 场景 | 测试位置 | 断言内容 | 运行命令 |
| --- | --- | --- | --- |
| AVX2（amd64） / NEON（arm64） 与 Go 路径等价 | `x32/*_test.go`、`x64/*_test.go`、`internal/functions/*_test.go` | 逐元素相等；长度 1、块大小 ±1；`NaN` / `±Inf` / `±0` | `go test ./x32/ ./x64/ ./internal/functions/` |
| 门面层两条路径 | 本包 Go 用例（**缺**） | `UseAVX2 = false` 与 `true` 下结果一致 | 见 §5 TODO |
| C++ 与 Go 逐值一致 | `tests/fixtures/conform_cpp`（160 个样本） | 已覆盖语义逐值相等；未覆盖语义跳过并列在 `CPP-GAPS.md` | `ctest --preset release-static -R fixtures_conform_cpp` |

## 5. TODO

| 待办项 | 缺什么 | 卡在哪 | 关联章节 | 关联 issue | 优先级 |
| --- | --- | --- | --- | --- | --- |
| 补齐 C++ 剩余算子 | 已落地 59 个功能行（§2 有 c++ 符号的行）；只剩 `PeeksAndValleys`（Go 侧就是空实现，返回 nil） | 该函数在 Go 里没有语义，等 Go 侧定案后再同步两侧；`AnyToGeneric` 的反射部分 / `AnyToSlice` / `CheckoutRawType` / `ExtractValueFromPointer` 属 `any` 反射，不移植 | §1 / §2 | #11 | 高 |
| C++ rolling 缺 min_periods | ~~Go 已接线，C++ 仍只有完整窗口入口~~ | 已随 #27 补齐并通过 fixture | §2 | 已解决 |
| `cross` 的相等口径 | 判据是**前一根与当前根都严格分居两侧**，任一格相等即不记；"先触碰再穿越"的常见写法会多记一根 | 需在补 Python spec 时确认这是意图还是历史写法 | §2 | 待建单 | 低 |
| `fit_line` 的失败返回 | 单点输入返回 `intercept = y[0]` 的水平线同时标记失败（`std::pair<Line,bool>`），C++ 照此保留 | 补 spec 时确认失败位与半成品是否都要保留 | §2 | 待建单 | 低 |
| 补 Python spec | 整包无 `.py`，语义准绳缺失 | 无技术阻塞，待排期 | §1 / §2 | 待建单 | 高 |
| 扩 Rust 覆盖面 | crate 已通过 160 个 fixture；线性代数、几何、时间与公式族未移植 | 按 `NAME-MAP.tsv` 逐行补，B3 在 pandas 仓库 | §1 / §2 | #21 | 中 |
| C++ 对拍覆盖 | `num_fixture_conform` 已全量接入 160 个样本；无 fixture 的转换族与 mat4mul 不在范围 | 新增 fixture 时同步扩展分派 | §4.6 | #27 | 低 |
| C++ 打包与导出 | `conan create` 与 `cmake --install` 未落地，`num_core` 只能源码内使用 | 安装布局要与扁平目录一起设计（头文件按 `num/` 前缀导出） | §3 | #11 | 中 |
| C++ 指令集分发 | 只按编译器基线下发，未做运行期多版本分发（Go 侧有 `UseAVX2` 开关与两套内核） | 需要 CMake 多目标编译 + 运行期 `__builtin_cpu_supports` / `IsProcessorFeaturePresent` 探测；会牺牲老 CPU 兼容，须先定策略 | §2 精度与平台 | #11 | 中 |
| `diff` 方向文档化核对 | C++ 已按 Go 口径实现为 `s[i-n] - s[i]`（与 pandas 相反），但 Go 侧无文档说明 | 需确认是设计意图还是历史包袱 | §2 | #11 | 低 |
| `Median` 是否内置排序 | Go 与 C++ 均按入参顺序取中间元素、**不排序**（空输入 `NaN`、单元素 `0`），排序责任在调用方 | 两侧已逐行一致（`median.h` 照抄 Go），待补 spec 时确认「不排序」是意图；定案后一起改 | §2 | 待建单 | 低 |
| C++ 覆盖率与基准 | 无覆盖率工具接入、无 Google Benchmark 用例与基线 | 需要在 CMake 侧引入 `llvm-cov` / `benchmark` | §4.2 / §4.5 | #11 | 低 |
| 根包用例按 `accel` 两遍跑 | 84 个 `Test` 均未引用 `functions.Info().Acceleration`，`UseAVX2 = false` 分支在根包层无用例（`x32` / `x64` 的用例已按 §7 跑两遍） | 无技术阻塞，待排期 | §4.1 / §4.6 | 待建单 | 中 |
| 提升覆盖率并设门槛 | 语句覆盖 42.1%，仓库未设门槛；`window_test.go` 只有基准、无 `Test` | 无技术阻塞，待排期 | §4.2 | 待建单 | 中 |
| 登记基准基线 | 46 个 `Benchmark*` 均无基线值与环境说明 | 需固定一台基准机，并写明 CPU / OS / commit | §4.5 | 待建单 | 中 |
| 引入重复度检测 | 仓库无重复度工具与阈值配置 | 无技术阻塞，待排期 | §4.3 | 待建单 | 低 |
| 逐功能 `.md` 是否拆分 | 当前整包共用 `num.md`；若后续按语义拆出 `abs.md` / `dot.md` 等，须同步删减 §2 对应行 | 待评审确定粒度 | §2 | 待建单 | 低 |
