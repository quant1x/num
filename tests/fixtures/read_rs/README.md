# tests/fixtures/read_rs —— Rust 读取器（C3）

独立 crate，只依赖 Rust 标准库，读取 C2 产出的 npz。**不是**生产 crate（生产 crate 属 M4），
因此不进仓库根的 `Cargo.toml`，按需用 `--manifest-path` 调用。

## 文件

| 文件 | 内容 |
| --- | --- |
| `Cargo.toml` | `num-read-fixtures` 二进制，无依赖 |
| `src/main.rs` | zip 中央目录、本地头、`.npy` 头解析、dtype 与长度校验、FNV-1a |

## 用法

```bash
cargo run --quiet --release --manifest-path tests/fixtures/read_rs/Cargo.toml -- tests/fixtures
```

## 边界

只读不写；不支持 deflate、`fortran_order`、大端 dtype。错误走 `Result`，不使用 `unwrap` / `expect`。
