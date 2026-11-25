# tests/fixtures/gen —— fixture 生成器（C2）

`gen_fixtures.py` 用 NumPy 产出输入与期望输出，落 `tests/fixtures/<分组>/<用例>.npz`。

## 职责

- 固定字面量输入 + NumPy 计算，保证可复现；不使用随机数、不联网。
- 按 `tests/fixtures/README.md` 的 npz 约定写元数据（`meta_semantic` / `meta_tolerance` / `meta_checksum`）。
- 覆写同名文件，不追加；zip 时间戳固定，重复执行产物逐字节一致。

## 用法

```bash
python tests/fixtures/gen/gen_fixtures.py            # 输出到 tests/fixtures
python tests/fixtures/gen/gen_fixtures.py --out DIR  # 输出到别处（仅调试用，不入库）
```

## 新增样本

1. 在 `elementwise()` 或 `reduction()` 里加一次 `save(...)` 调用：语义名必须已在
   `tests/contract/NAME-MAP.tsv` 登记；
2. 期望值用 NumPy 计算；与 Go 口径有差异的（`median` 不排序、`diff` 方向相反）在函数里注明并按 Go 口径手算；
3. 重跑生成器与三套读取器，确认产物一致。

## 边界

不写算子实现、不做跨语言比较；比较由各语言读取器与 M1 的对拍用例承担。
