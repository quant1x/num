# C++ 相对 Go 的缺口清单（M1）

本清单是 M1 的出口物：按 [tests/contract/NAME-MAP.tsv](contract/NAME-MAP.tsv) 逐行核对 C++ 侧是否有同语义实现，
并记录**已有实现但口径需要留意**的项。Go 是当前临时基准，不是 spec；两边都与 NumPy 不同的口径单列在 §3。

**核对范围**：N1（标量与类型边界）、N2（逐元素）、N3（归约）、N6（向量与容器）。N4 / N5 属 M2，不在本表。

## 0. 对拍状态

`tests/fixtures/conform_cpp` 已把 C++ 已有算子全部接入 fixture：**139 个样本逐值通过，0 跳过**。

最近补齐：`rolling`（含 `min_periods`）/ `rolling_apply` / `cumprod` / `bool_to_numeric`，以及上一批的
一元、舍入、三角与指数对数、逻辑异或、选择/索引、prod/norm/quantile/none、浮点幂与浮点区间。

## 1. 无 fixture 的接口（不在对拍范围）

| 语义 | 现状 | 说明 |
| --- | --- | --- |
| `from_bool` / `from_int32` / `from_int64` / `from_float` | 无（C++ 直接用模板构造） | 转换族没有 fixture |
| `to_bool` / `to_int32` / `to_int64` / `to_float32` / `to_float64` | 无（C++ 用 `xt::cast`） | 同上 |
| `dtype2int` | 无 | 同上 |
| `any_to_scalar` | 只有 `to_bool` / `to_string` 的值重载 | C++ 无运行期 `any`，可判为不适用 |

## 2. Go 侧不对称（先定 Go 落点，再谈 C++ 对位）

这些不是 C++ 的缺失，而是 Go 自身两条宽度或门面签名不一致；C++ 若按现状补，会固化一套只存在于一门的语义。

| 项 | 现状 | 影响 |
| --- | --- | --- |
| 三角 / 指数对数 | `Sin` / `Cos` / `Exp` / `Log` / `Log2` / `Log10` / `SinCos` 只有 `x32`，`x64` 无同名实现，门面也没有入口 | C++ 模板已覆盖 f32 / f64，但 Go 缺 f64 对照；待 Go 补齐后对齐 |
| 方差 | `Variance` 只接受 `[]float64`，没有 f32 路径；`Std` 有两条宽度 | f32 方差无法对拍 |
| `Where` 条件 | `Where[T StatType]` 要求条件与取值**同类型**，不能传 `[]bool` | 与 NumPy 的 `where(bool, x, y)` 签名不同，C++ 已按同类型实现 |
| `Ones` | 按入参切片长度生成，不接受形状整数 | 与 `Zeros(shape)` 不对称 |
| `Expm1` / `Log1p` | `Expm1` 返回单值，`Log1p` 返回 `(值, error)` | 同类算子失败约定不统一 |
| `find` 未命中 | 返回 `-1`（已按代码修正契约表） | 早前契约写的 `len(x)` 是错的 |

## 3. 已有实现但口径需注意

以下 C++ 已落地，且**与 Go 一致**，但两者都和 NumPy 不同；fixture 只用双方一致的取值，不把这些差异写成"标准"。

| 语义 | 口径 | 与 NumPy 的差异 |
| --- | --- | --- |
| `median` | 按入参顺序取中间元素，不排序 | NumPy 会排序 |
| `diff` | `s[i-n] - s[i]` | pandas / NumPy 方向相反 |
| `fill` / `fill_na` | `±Inf` 与 `NaN` 都视为缺失 | NumPy 只填 NaN |
| `maximum` / `minimum` | 只在 `y > x` 时覆盖，NaN 不传播 | NumPy 传播 NaN |
| `concat1d` | 逐元素拼成 n×3（列式） | NumPy `concatenate` 按行拼接 |
| `equal` / `not_equal` / `logic_xor` | C++ 统一按最长对齐、缺位补 0；Go 的 `Equal` / `NotEqual` / `Xor` 按第一个操作数长度 | 样本等长，未触发；长度不齐时结果不同 |
| `shift` 非整数常量 | Go 运行期 panic；C++ 用 `= delete` 编译期拒绝 | 机制不同，语义相同 |
| `decimal` | 半程远离零，NaN 按 0 参与 | 与 NumPy `round` 的银行家舍入不同 |

## 4. 下一步

1. C++ 补齐优先按 §1 的顺序走：先 N2 逐元素（`neg` / `inv` / 舍入 / 三角与指数对数），再 N3（`prod` / `norm` / `quantile` / `none`），最后 N6（`gather` / `scatter` / `select` / `mat4mul`）。
2. 补 C++ 之前先定 §2 的两条宽度问题（尤其 f64 三角 / 指数对数的落点），否则补出来的符号没有 Go 对照。
3. `num_fixture_conform` 已全量接入 129 个样本；后续新增 fixture 时同步扩展分派。
