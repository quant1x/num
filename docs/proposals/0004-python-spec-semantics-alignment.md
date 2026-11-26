# 提案 0004：Python spec 与 NumPy 的语义对齐（第一批）

- 状态：**部分定案**——§2.1 `median` 已采纳 NumPy 并落地（#66）；§2.2 ~ §2.5 待定
- 前置：[0003](0003-python-numpy-package-layout.md) 只搬家与改名、不改数值；本提案处理它 §2.5 留下列出的语义债
- 影响面：`num/core.py`、`num/linalg.py`、`tests/fixtures/gen/gen_fixtures.py` 与重生成的 npz、`NAME-MAP.tsv` 的 `note`；对齐实现侧 `median.go` / `x32/vek.go` / `x64/vek.go` / `median.h` / `rounding.h` / `median.rs` / `vek.rs` / `median.zig` / `vek.zig`

## 1. 审计口径

以 NumPy 2.x 的同名函数为基准，f64 输入，逐例比对返回值与类型。下表结果均为实测，不含推断。

## 2. 逐项结果

### 2.1 median：已定案，采纳 NumPy（#66）

| 输入 | 改动前 `num.median` | `np.median` |
| --- | --- | --- |
| `[3.0]` | `0.0` | `3.0` |
| `[3.0, 1.0, 2.0]` | `1.0` | `2.0` |
| `[4.0, 1.0, 3.0, 2.0]` | `2.0` | `2.5` |
| `[nan, 1.0, 2.0]` | `1.0` | `nan` |

`num` 原先按入参位置取中间元素、不排序，单元素固定返回 `0`。这不是中位数：未排序输入取到的是任意一个元素，单元素输入直接丢掉真实值。定案采纳 NumPy 口径——先排序、偶数长度取中间两值均值、单元素返回其本身、NaN 参与即传播；落地后 `num.median` 与 `np.median` 逐例一致。

影响：`median-unsorted-f64` 改名为 `median-basic-f64` 并按 NumPy 期望重算，`median-even-f64` 重算，另补 `median-single-f64` 与 `median-nan-f64`；四门实现改为排序后取值。本仓没有内部调用方依赖"不排序"，排序开销可接受。

Go 侧其实已有排序路径：私有 `__median_go` 走 `internal/partial.TopK` 快速选择、单元素返回 `x[0]`、偶数长度走 `__quantile_go(x, 0.5)`，只是公开的 `Median` 没接线；x32 / x64 的 `Median` 早已走排序内核。C++ / Rust / Zig 按公开入口逐行照抄，`median.h` 里还留了"不要在 C++ 侧顺手修正"的注释。落地时公开入口统一改为直接排序实现（`(a+b)/2` 与 spec 逐位一致），`__median_go` 保持未接线；`median.h` 的旧注释一并删除。

### 2.2 round / ceil / floor：返回类型建议立即改，半数口径待定

| 表达式 | `num` | NumPy |
| --- | --- | --- |
| `round([0.5, 1.5, 2.5, -0.5, -1.5, -2.5])` | `[1, 2, 3, -1, -2, -3]`（int） | `[0, 2, 2, -0, -2, -2]`（float64） |
| `round([1.4, -1.6])` 的返回类型 | `int` | `float64` |
| `ceil([1.4])` / `floor([1.4])` 的返回类型 | `int` | `float64` |

两件事要分开：

1. **返回类型**：`num` 返回 Python `int`，是 `math.floor` / `math.ceil` 的实现产物；Go / C++ / Rust / Zig 的同名实现都返回浮点，说明 int 是 Python 侧瑕疵而非契约。建议改为返回 `float`，现有 `round-f64` fixture 的数值不变。
2. **半数口径**：`num` 半程远离零，`np.round` 半程取偶。两者只在恰好 `.5` 时不同，现有 fixture 只含 `2.25`，该点在两种口径下都是 `2.0`，所以 fixture 覆盖不到。需要定案：
   - 方案 A：采纳 `np.round`（半程取偶），四门实现同步改；
   - 方案 B：保留四舍五入（域内公式常按此口径），文档与 `note` 明确写出与 `np.round` 的差异，并补 `0.5 / 1.5 / 2.5` 的反例用例把差异钉住。

### 2.3 diff：建议保留

`num.diff([1, 3, 6]) == [nan, -2, -3]`，长度与输入一致；`np.diff` 为 `[2, 3]`，长度 `n-1`。方向是前值减当前值、前 `n` 个补 NaN，属金融时序常用口径，且定长输出便于与输入逐元素对拍。本仓无内部调用方依赖 NumPy 口径。建议保留，`NAME-MAP` 的 `note` 补"与 `np.diff` 相反、不缩短长度"。

### 2.4 covariance：建议保留，只补文档

`num.covariance` 是总体协方差（`ddof=0`），与 `np.cov(..., ddof=0)` 逐值一致，与 `np.cov` 默认 `ddof=1` 不同；这与 `num.var` / `num.std` 的总体口径自洽。建议保留，`note` 补 `ddof=0`。

### 2.5 空输入：建议保留，只补文档

| 调用 | `num` | NumPy |
| --- | --- | --- |
| `max([])` / `min([])` | `0.0` | 抛 `ValueError` |
| `var([])` / `std([])` | `0.0` | `nan` |
| `sum([])` / `prod([])` | `0.0` / `1.0` | 一致 |

`0.0` 是流式统计的默认值，Go 侧归约也按此写；改动会波及全部语言的空序列契约，收益低。建议保留，在模块文档写明"空输入返回 0"，避免与 NumPy 混淆。

### 2.6 已核对无偏差

- `maximum` / `minimum` 的 NaN 传播与 NumPy 一致（NaN 比较为假，保留左值，因此 NaN 传播）；
- 长度不一致直接抛 `ValueError`，标量广播行为与 NumPy 一致；
- `sum` / `prod` / `count` / `argmax` / `quantile` 与 NumPy 逐值一致。

## 3. 建议的落地批次

| 批次 | 内容 | 是否需定案 | 验收 |
| --- | --- | --- | --- |
| P1 | `median` 采纳 NumPy；新增单元素与已排序用例 | 是（§2.1） | 新 fixture 全过；四门 `median` 排序实现 |
| P2 | `round` / `ceil` / `floor` 返回 `float` | 否（类型瑕疵） | 现有 fixture 数值不变，类型断言通过 |
| P3 | `round` 半数口径 | 是（§2.2 方案 A / B） | 选 A 则四门同步，选 B 则补反例用例 |
| P4 | `diff` / `covariance` / 空输入 只补文档与 `note` | 否 | 文档一致，fixture 无变化 |

## 4. 验收命令

```bash
python -m unittest discover -s tests/python -v
go test ./tests/contract/ ./internal/functions/ ./x32/ ./x64/
cargo test
zig test num_fixture_test.zig
ctest --preset release-static
```

## 5. 不做

- 不新增 NumPy 里没有的函数；`cross` / `fit_line` 的口径复核仍按 `num.md` §5 单列。
- 不引入 `nanmean` / `nanmax` 之类的 NaN 感知家族。
- 不改 `authority` 列，P1 / P3 落地时仍由 Python 作基准、其余语言跟随。
