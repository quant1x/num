# io —— D1 CSV 读写（D2）

`package io` 按 [docs/DATA-FORMAT.md](../docs/DATA-FORMAT.md) 的 **CSV 基线**读写表格数据：
UTF-8 无 BOM、纯列名表头、空字段表缺失、`inf` / `-inf`、布尔 `true` / `false`、
时间 ISO 8601（无时区按 UTC）、LF 写出、按主键覆盖。

包名与标准库 `io` 相同，包内不导入标准库 `io`；使用方需要时可起别名：

```go
import numio "github.com/quant1x/num/io"
```

## 职责与边界

- **收**：带 schema 的 CSV 读取与写出、表头与字段校验、缺失值与无穷的文本映射、时间归一化。
- **不收**：CSV 之外的交换格式（JSON Lines / 二进制见 D1 的可选扩展，尚未落地）、
  业务主键合并、按列统计。上层数据管道（B5）在这里之上做覆盖写。

## API

```go
schema := []numio.Column{
    {Name: "date", Kind: numio.String},
    {Name: "close", Kind: numio.Float64},
    {Name: "volume", Kind: numio.Int64},
    {Name: "flag", Kind: numio.Bool},
}
table, err := numio.ReadFile("in.csv", schema)
rows, err := table.Rows()          // 校验各列长度一致
err = numio.WriteFile("out.csv", table)
```

`Column` 同时承载 schema（`Name` / `Kind`）与数据（`Floats` / `Ints` / `Bools` / `Strings` / `Times`），
`Missing` 按行标记缺失；`Rows()` 以 `Missing` 的长度为准。

## 按主键覆盖（B5）

```go
merged, err := numio.MergeByKey(existing, incoming, "account")
err = numio.WriteByKey("data.csv", incoming, "account")
```

- 主键列只支持 `String` / `Int64` / `Time`；浮点相等受 NaN / 精度影响，不允许当键。
- 已有键更新整行、新键按增量顺序追加，已有行顺序不变。
- 主键缺失、增量内重复键、已有文件内重复键都返回错误（不做隐式去重）。
- `WriteByKey` 读取已有文件 → 合并 → 整文件覆盖；同一增量重复执行，产物逐字节一致。

## 口径

| 项 | 规则 |
| --- | --- |
| 表头 | `^[a-z][a-z0-9_]*$`，长度 ≤ 64，同文件唯一；顺序可与 schema 不同，按列名对齐 |
| 缺失 / NaN | 写空字段；读取空字段得到 `Missing = true`，零值入列 |
| 无穷 | `inf` / `-inf`；拒绝 `NaN` / `Infinity` 等字面量 |
| 时间 | `YYYY-MM-DD` 或 RFC3339；读取归一化到 UTC，写出为 RFC3339 |
| 行尾 | 写 LF；读接受 LF 与 CRLF |
| 幂等 | 同一张表重复写出逐字节一致 |

## 测试

```bash
go test ./io/
```

覆盖读写往返、幂等、BOM 拒绝、非法表头、字段数不符、NaN 字面量拒绝、CRLF、时间归一化，
以及按主键合并的更新 / 追加 / 重复键 / 缺键 / 浮点键拒绝 / 幂等。
跨语言互读互写见 `tests/io/roundtrip.py`，按主键覆盖见 `tests/io/pipeline.py`。
