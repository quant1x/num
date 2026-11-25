# tests/io —— D2 跨语言 CSV 互读互写

Python 写出一份 D1 格式的 CSV，Go / C++ / Rust 各自读入并按同一规则重写，再由 Python 读回三份产物逐值比较。
这是 D2 的验收入口，与 `tests/fixtures` 的 npz 对拍相互独立。

## 布局

| 路径 | 用途 |
| --- | --- |
| `roundtrip.py` | 驱动器：写样本、调用三种语言、比较值、检查幂等 |
| `go_roundtrip/` | Go 命令行工具（`go run`，使用 `io/` 包） |
| `cpp_roundtrip/` | C++ 命令行工具（`num_csv_roundtrip`，链接 `num_core`） |
| `rs/` | Rust 独立工具 crate（无三方依赖；生产 crate 属 M4） |

样本与临时产物落在 `temp/tests/io/roundtrip/`，不入库。

## 运行

```bash
# 1) 先构建 C++ 工具
cmake --build --preset release-static --target num_csv_roundtrip

# 2) 跑互读互写
python tests/io/roundtrip.py
```

输出 `ALL PASS` 表示：

- Go / C++ / Rust 都能读 Python 产物，并按 D1 重写；
- Python 能读 Go 产物；
- 每门语言重复写出的文件逐字节一致。

## schema

命令行统一用 `name:kind` 逗号串：`date:s,close:f,volume:i,flag:b,note:s`
（`f` = f64，`i` = i64，`b` = bool，`s` = string，`t` = time）。三个工具共用这套 spec，
新增用例只改 `roundtrip.py` 的 `SCHEMA` 与 `ROWS`。

## 边界

样本覆盖缺失值、`±inf`、带逗号 / 引号的字符串与 UTF-8；不覆盖浮点文本精度以外的格式化差异。
值比较而非字节比较——三门的浮点文本写法允许不同，但必须逐值还原。

时间列：Go 读取时归一化到 UTC 并按 RFC3339 写出；C++ / Rust 目前只做 ISO 8601 格式校验、原样保留文本。
跨语言时间往返样本待时间类型在两侧真正落地后再加。
