# C++ 相对 Go 的缺口清单（M1）

本清单是 M1 的出口物：按 [tests/contract/NAME-MAP.tsv](contract/NAME-MAP.tsv) 逐行核对 C++ 侧是否有同语义实现，
并记录**已有实现但口径需要留意**的项。Go 是当前临时基准，不是 spec；两边都与 NumPy 不同的口径单列在 §3。

**核对范围**：N1（标量与类型边界）、N2（逐元素）、N3（归约）、N6（向量与容器）。N4 / N5 属 M2，不在本表。

## 1. 未落地（C++ 无同语义实现）

| 语义 | Go 符号 | 缺口 | 影响 |
| --- | --- | --- | --- |
| `pow_float` | `x32:Pow` / `x64:Pow` | C++ `pow` 只接受整数指数；缺逐元素幂 | N2 对位 |
| `neg` | `x32:Neg` / `x64:Neg` | 无 | N2 对位 |
| `inv` | `x32:Inv` / `x64:Inv` | 无 | N2 对位 |
| `ceil` / `floor` / `round` | `x32` / `x64` 同名 | 无；舍入族只有 `decimal` | N2 对位 |
| `sin` / `cos` | `x32:Sin` / `x32:Cos` | 无 | N2 对位；Go 侧 x64 也缺，见 §2 |
| `exp` / `log` / `log2` / `log10` | `x32` 同名 | 无；指数对数族只有 `log1p` / `expm1` | N2 对位 |
| `logic_xor` | `x32:Xor` / `x64:Xor` | 无 | N2 对位 |
| `select` | `x32:Select` / `x64:Select` | 无 | N2 对位 |
| `find` | `x32:Find` / `x64:Find` | 无（`type_strings.h` 只有 `std::string::find` 调用，`ztable.h` 的 `find_percent` / `find_zscore` 是另一语义） | N2 对位 |
| `gather` / `scatter` | `x32` / `x64` 同名 | 无 | N6 对位 |
| `prod` | `x32:Prod` / `x64:Prod` | 无 | N3 对位 |
| `norm` / `manhattan_norm` | `x32` / `x64` 同名 | 无 | N3 对位 |
| `quantile` | `x32:Quantile` / `x64:Quantile` | 无 | N3 对位 |
| `none` | `x32:None` / `x64:None` | 无（只有 `all` / `any`） | N3 对位 |
| `range_float` | `x32:Range` / `x64:Range` | C++ `range` 是按整数个数的入口，缺浮点区间重载 | N6 对位 |
| `mat4mul` | `x32:Mat4Mul` / `x64:Mat4Mul` | 无 | N6 对位 |
| `from_bool` / `from_int32` / `from_int64` / `from_float` | `x32` / `x64` 同名 | 无（C++ 直接用模板构造） | N1 对位 |
| `to_bool` / `to_int32` / `to_int64` / `to_float32` / `to_float64` | `x32` / `x64` 同名 | 无（C++ 用 `xt::cast`） | N1 对位 |
| `dtype2int` | `num:DType2Int` | 无 | N1 对位 |
| `any_to_scalar` | `num:AnyToInt32` 等 | 只有 `to_bool` / `to_string` 的值重载，缺 `any` 到数值的入口 | N1 对位；C++ 无运行期 `any`，可判为不适用 |

## 2. Go 侧不对称（先定 Go 落点，再谈 C++ 对位）

这些不是 C++ 的缺失，而是 Go 自身两条宽度或门面签名不一致；C++ 若按现状补，会固化一套只存在于一门的语义。

| 项 | 现状 | 影响 |
| --- | --- | --- |
| 三角 / 指数对数 | `Sin` / `Cos` / `Exp` / `Log` / `Log2` / `Log10` / `SinCos` 只有 `x32`，`x64` 无同名实现，门面也没有入口 | f64 侧无对照，C++ 补齐前先定 Go 的 f64 落点 |
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
| `shift` 非整数常量 | Go 运行期 panic；C++ 用 `= delete` 编译期拒绝 | 机制不同，语义相同 |
| `decimal` | 半程远离零，NaN 按 0 参与 | 与 NumPy `round` 的银行家舍入不同 |

## 4. 下一步

1. C++ 补齐优先按 §1 的顺序走：先 N2 逐元素（`neg` / `inv` / 舍入 / 三角与指数对数），再 N3（`prod` / `norm` / `quantile` / `none`），最后 N6（`gather` / `scatter` / `select` / `mat4mul`）。
2. 补 C++ 之前先定 §2 的两条宽度问题（尤其 f64 三角 / 指数对数的落点），否则补出来的符号没有 Go 对照。
3. C++ 补一个 `num_fixture_conform`（与 Go 的 `tests/fixtures/conform` 同构）：读同一批 npz、调用 C++ 实现、逐值比对；完成后本清单逐行销账。
