# tests/fixtures/read_go —— Go 读取器（C3）

`package read_go`，只依赖 Go 标准库（`archive/zip`、`encoding/binary`），读取 C2 产出的 npz。

## 文件

| 文件 | 内容 |
| --- | --- |
| `npz.go` | `ReadFile` / `Verify` / `Walk`；zip 与 `.npy` 头解析、dtype 与长度校验、FNV-1a 校验和 |
| `npz_test.go` | 读取 `tests/fixtures` 下全部样本；单条 `add-f64` 逐值核对；缺校验和必须报错 |

## API

```go
fixture, err := read_go.ReadFile("tests/fixtures/elementwise/add-f64.npz")
err = fixture.Verify()
values, err := fixture.Arrays["x"].Float64()
```

- `Array.Descr` 保留 `.npy` 的 dtype 写法（`<f8` / `<f4` / `|b1` / `|S32` …）。
- `Float64()` / `Int64()` 是对拍用的便捷解包；不支持的类型返回 error，不做隐式转换。

## 边界

只读不写；不支持 deflate 压缩的 npz（生成器固定 ZIP_STORED），不支持 `fortran_order`、大端 dtype。
