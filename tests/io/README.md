# tests/io —— D2 跨语言 CSV 互读互写

Python 写出一份 D1 格式的 CSV，Go / C++ / Rust 各自读入并按同一规则重写，再由 Python 读回三份产物逐值比较。
这是 D2 的验收入口，与 `tests/fixtures` 的 npz 对拍相互独立。

## 布局

| 路径 | 用途 |
| --- | --- |
| `roundtrip.py` | D2 驱动器：写样本、调用三种语言、比较值、检查幂等 |
| `pipeline.py` | M6 驱动器：基线 + 增量 → 三门按主键合并 → 比较值、检查幂等 |
| `go_roundtrip/` | Go 往返工具（`go run`，使用 `io/` 包） |
| `go_merge/` | Go 合并工具（`io.WriteByKey`） |
| `cpp_roundtrip/` | C++ 往返工具（`num_csv_roundtrip`，链接 `num_core`） |
| `cpp_merge/` | C++ 合并工具（`num_csv_merge`） |
| `rs/` | Rust 独立工具 crate：lib + `num-csv-roundtrip` / `num-csv-merge` 两个 bin |

样本与临时产物落在 `temp/tests/io/roundtrip/`，不入库。

## 运行

```bash
# 1) 先构建 C++ 工具
cmake --build --preset release-static --target num_csv_roundtrip num_csv_merge

# 2) 互读互写
python tests/io/roundtrip.py

# 3) 按主键覆盖 + 幂等
python tests/io/pipeline.py
```

`roundtrip.py` 的 `ALL PASS` 表示：

- Go / C++ / Rust 都能读 Python 产物，并按 D1 重写；
- Python 能读 Go 产物；
- 每门语言重复写出的文件逐字节一致。

`pipeline.py` 的 `ALL PASS` 表示三门都能把增量按主键合并进基线（已有键更新、新键追加），
且重复合并同一增量时产物逐字节一致。

## schema

命令行统一用 `name:kind` 逗号串：`date:s,close:f,volume:i,flag:b,note:s`
（`f` = f64，`i` = i64，`b` = bool，`s` = string，`t` = time）。三个工具共用这套 spec，
新增用例只改 `roundtrip.py` 的 `SCHEMA` 与 `ROWS`。

## 边界

样本覆盖缺失值、`±inf`、带逗号 / 引号的字符串与 UTF-8；不覆盖浮点文本精度以外的格式化差异。
值比较而非字节比较——三门的浮点文本写法允许不同，但必须逐值还原。

时间列：Go 读取时归一化到 UTC 并按 RFC3339 写出；C++ / Rust 目前只做 ISO 8601 格式校验、原样保留文本。
跨语言时间往返样本待时间类型在两侧真正落地后再加。
