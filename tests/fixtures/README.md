# tests/fixtures —— 跨语言对拍样本（C2 / C3）

本目录存**语言无关的黄金样本**：由 Python（NumPy）生成，Go / C++ / Rust 读取后判定一致。
样本是跨语言一致性的唯一判据（[CONTRIBUTING.md](../../CONTRIBUTING.md) §7）。

## 布局

```text
tests/fixtures/
├── gen/gen_fixtures.py     # 生成器：NumPy 产出输入与期望输出
├── elementwise/            # N2 逐元素样本
├── reduction/              # N3 归约样本
├── read_go/                # Go 读取器（package read_go，读取器用例在 npz_test.go）
├── read_cpp/               # C++ 读取器（num_fixture_check 目标）
└── read_rs/                # Rust 读取器（独立 crate，无三方依赖）
```

样本路径为 `<分组>/<用例>.npz`，分组对应 N 层（`elementwise` = N2，`reduction` = N3）。

## npz 约定

| 条目 | 内容 |
| --- | --- |
| `x` / `y` / `condition` … | 输入数组，名字取自 `tests/contract/NAME-MAP.tsv` 的 `min_params` |
| `expected` | 期望输出；返回值多于一个时用 `expected_0`、`expected_1` |
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
go test ./tests/contract/ ./tests/fixtures/read_go/
cargo run --quiet --release --manifest-path tests/fixtures/read_rs/Cargo.toml -- tests/fixtures
ctest --preset release-static -R fixtures_npz_read_cpp
```

三套读取器对同一批样本都通过，才算 C3 达标。新增样本时三处都要重跑，不只看自己那门语言。

## 与上下层

- 期望值只由 **Python（NumPy）** 生成；其它语言只读不写。Python spec 落地前，
  `NAME-MAP.tsv` 的 `note` 列与样本注释共同记录与 NumPy 的已知差异（如 `median` 不排序）。
- 样本是 M1 内核对齐的输入；算子实现不得反向修改样本去迁就当前输出。
