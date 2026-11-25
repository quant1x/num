# cmd/numctl —— Go 单二进制入口（G1）

`numctl` 是 num 仓库唯一的命令行程序：读写 D1 基线 CSV，并按主键合并；纯 Go，`CGO_ENABLED=0` 可静态编译，
不依赖 Python 运行时。

## 构建

```bash
CGO_ENABLED=0 go build -o numctl ./cmd/numctl
```

## 子命令

```bash
numctl csv roundtrip <schema> <in.csv> <out.csv>
numctl csv merge     <schema> <path> <delta.csv> <key>
numctl csv stats     <schema> <file.csv> <column>
```

- `schema`：`name:f|i|b|s|t,...`（f=f64，i=i64，b=bool，s=string，t=time），见 `io/README.md`。
- `roundtrip`：读入再按 D1 写出，用于验证交换格式。
- `merge`：读已有文件（不存在则直接写），按 `key` 更新已有行、追加新键；重复执行幂等。
- `stats`：对数值列输出 `count / sum / mean / std / min / max`（总体标准差，跳过缺失行）。

## 退出码与边界

- `0` 成功；`1` 数据或参数错误（错误信息写 stderr）；`2` 留给将来的用法错误细分。
- 不提供 `--version`：版本由 tag / autochangelog 管理，代码里不出现版本字符串。
- 只处理 D1 的 CSV 基线；Arrow / Parquet 与远端服务不在本命令范围。

## 测试

```bash
go test ./cmd/numctl/
```

覆盖三个子命令、缺失值跳过、非数值列拒绝与用法错误。
