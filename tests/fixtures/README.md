# tests/fixtures —— 跨语言对拍样本（C2 / C3）

本目录存**语言无关的黄金样本**：由 Python（NumPy）生成，Go / C++ / Rust 读取后判定一致。
样本是跨语言一致性的唯一判据（[CONTRIBUTING.md](../../CONTRIBUTING.md) §7）。

## 布局

```text
tests/fixtures/
├── gen/gen_fixtures.py     # 生成器：NumPy 产出输入与期望输出
├── scalar/                 # N1 标量与类型边界样本
├── elementwise/            # N2 逐元素样本
├── reduction/              # N3 归约样本
├── window/                 # N4 滚动族样本（rolling / shift / diff）
├── container/              # N6 向量与容器样本
├── linalg/                 # N6 线代样本（dot2d / matmul / inverse）
├── geometry/               # 几何与多项式（slope / poly_val）
├── conform/                # Go 实现对拍用例（M1 准入）
├── conform_cpp/            # C++ 对拍（未实现语义跳过，缺口见 CPP-GAPS）
├── read_go/                # Go 读取器（package read_go，读取器用例在 npz_test.go）
├── read_cpp/               # C++ 读取器（num_fixture_check 目标）
├── read_rs/                # Rust 读取器（独立 crate，无三方依赖）
└── CPP-GAPS.md             # C++ 相对 Go 的缺口清单（M1 出口）
```

样本路径为 `<分组>/<用例>.npz`，分组对应 N 层（`scalar` = N1、`elementwise` = N2、`reduction` = N3、
`window` = N4、`container` = N6、`linalg` = N6 线代）。当前共 139 个样本，覆盖 N1 ~ N4、N6 的数值语义；字符串解析、`any` 分派与
反射类转换不在样本范围内（其口径见契约表，落地时另补）。

## npz 约定

| 条目 | 内容 |
| --- | --- |
| `x` / `y` / `condition` … | 输入数组，名字取自 `tests/contract/NAME-MAP.tsv` 的 `min_params` |
| `expected` | 期望输出；返回值多于一个时用 `expected_0`、`expected_1` |
| `expected_flat` / `expected_lengths` | 嵌套或不等长输出（`Rolling` 的窗口列表）用展平值 + 每段长度表示 |
| `meta_semantic` | `S32` 定宽字节串标量（numpy 的 `S32` dtype），语义名，对应契约表 `semantic` 列 |
| `meta_tolerance` | `f8` 标量；0 表示逐位相等，非 0 为该用例的浮点容差 |
| `meta_checksum` | `u8` 标量；所有非 `meta_` 数组按名字排序串联后的 FNV-1a 64 |

- 采用 `np.savez` 的等价形态（ZIP_STORED），读取器**不实现 deflate**，避免为解压引入三方依赖。
- zip 时间戳与条目顺序固定，重复生成逐字节一致；样本入库，不靠运行期生成。
- 数组 dtype 与 shape 由 `.npy` 头承载，读取器必须断言 `元素个数 × 元素宽度 == 数据长度`，并核对校验和。

## 生成

```bash
python tests/fixtures/gen/gen_fixtures.py
```

需要 NumPy；生成器不联网、不使用随机数，输入全部是字面量。

## 校验

```bash
# C1 契约 + Go 对拍
go test ./tests/contract/ ./tests/fixtures/conform/ ./tests/fixtures/read_go/
# 三套读取器
cargo run --quiet --release --manifest-path tests/fixtures/read_rs/Cargo.toml -- tests/fixtures
ctest --preset release-static -R fixtures_npz_read_cpp
ctest --preset release-static -R fixtures_conform_cpp
```

Go / C++ 对拍用例把样本交给各自实现计算并与 `expected` 比较（C++ 已全量 139 个通过）；三套读取器保证同一批样本能被三门语言读出。
新增样本时三处都要重跑，不只看自己那门语言。

C++ 相对 Go 的缺口（缺哪些算子、哪些口径待核）见 [CPP-GAPS.md](CPP-GAPS.md)。

## 与上下层

- 期望值只由 **Python（NumPy）** 生成；其它语言只读不写。Python spec 落地前，
  `NAME-MAP.tsv` 的 `note` 列与样本注释共同记录与 NumPy 的已知差异（如 `median` 不排序）。
- 样本是 M1 内核对齐的输入；Go 对拍在 `conform/`，Rust 对拍在根 crate 的 `tests/conform.rs`，
  C++ 对拍待补。算子实现不得反向修改样本去迁就当前输出。
