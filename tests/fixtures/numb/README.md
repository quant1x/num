# tests/fixtures/numb —— `.numb` 黄金样本（D1 §3.2，版本 2）

`gen_numb.py` 生成 6 个样本，覆盖三种宽度、二维形状、标量与空维度：

| 文件 | dtype | 形状 | 说明 |
| --- | --- | --- | --- |
| `f8-vector.numb` | `f8` | `[5]` | float64 一维 |
| `f4-matrix.numb` | `f4` | `[2,3]` | float32 二维，行主序 |
| `i8-vector.numb` | `i8` | `[4]` | int64，含负值 |
| `b1-vector.numb` | `b1` | `[3]` | bool，0/1 各一字节 |
| `f8-scalar.numb` | `f8` | `[]` | ndim=0 的标量 |
| `f8-empty.numb` | `f8` | `[0]` | 空维度 |

三端用例（Go `io/numb_test.go`、C++ `numb_test.cpp`、Rust `tests/numb.rs`）做三件事：

1. 读出 dtype / 形状 / 值并与预期比对；
2. 把读到的内容**重新编码**，要求与原文件逐字节一致；
3. 各自构造含 NaN / Inf / -0.0 的样本往返，并覆盖 magic、版本、dtype、截断、尾部多余的拒绝路径。

重新生成：

```bash
python tests/fixtures/numb/gen_numb.py
```

固定输入、无随机数，重复执行逐字节一致。
