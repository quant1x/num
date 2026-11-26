# 提案 0003：Python spec 按 NumPy 模块归集

- 状态：**已定案（方案 C：独立 `num/` 包）**，2026-11-26 落地
- 取代：[0002](0002-python-package-layout.md) 的布局部分（「仓库根即包」）；0002 的 `requires-python` 与「不依赖 NumPy」结论继续有效
- 影响面：`num/` 包、`pyproject.toml`、`tests/python`、`NAME-MAP.tsv` 的 python 列、`CONTRIBUTING.md` §3.1 / §3.2、`BUILDING.md` §4.5、`docs/PLAN.md`

## 1. 背景

Python 已是 spec，但源码仍按 Go 文件平铺：50 个短 `.py` 与 `.go` 同名同目录，模块边界由 Go 的文件粒度决定。Go 侧因为语言习惯把功能拆得很细，这个粒度不适合作为 spec 的组织方式；NumPy 的公开命名空间才是 Python 使用者看到的边界。

## 2. 决定

1. Python 源码移入真正的 `num/` 包，仓库根直接 `import num`；不再使用 `package-dir = {num = "."}`，也不再要求把仓库父目录放进 `sys.path`。
2. 按 NumPy 公开命名空间组织：顶层 `num`、`num.linalg`、`num.polynomial`、`num.lib`；num 专有领域能力单独成 `num.ta`、`num.ztable`、`num.freq`、`num.lib.format`。
3. canonical 函数名取 NumPy 名（`var` / `greater` / `negative` / `reciprocal` / `subtract` / `multiply` / `divide` / `power` / `dot` / `transpose` / `linalg.inv` / `polyfit` / `polyval` / `take` / `compress` / `isnan`）；旧名只在 `num/__init__.py` 保留一版别名，不进 `__all__`。
4. 旧模块路径（`num.sum`、`num.matrix`、`num.x32.vek`）不保留。
5. 本次只搬家与改名，不改数值语义；与 NumPy 的已知差异（`median` 不排序、`diff` 方向、`round` 半数口径、`covariance` ddof、NaN 传播、长度对齐等）另立 issue 分批处理。

## 3. 目标布局

```text
num/
  __init__.py
  core.py
  linalg.py
  polynomial.py
  lib/{__init__,convert,parse,format,array}.py
  ta.py
  ztable.py
  _ztable_data.py
  freq.py
```

`num/__init__.py` 只做转出；`num/ztable.py` 是唯一导入 `_ztable_data.py` 的模块，保证 `import num` 不加载 428 KB 表。

## 4. 与「同目录同名」的关系

`CONTRIBUTING.md` §3.2 的同目录同名规则继续约束 C++ / Rust / Go / Zig；**Python 是唯一例外**：spec 按 NumPy 命名空间组织，跨语言对位由 `NAME-MAP.tsv` 的 python 列承担。其他语言不得援引此例外。

## 5. 落地与验证

- `NAME-MAP.tsv` 的 python 列改为相对 `num/` 的 `module.function`；`tests/contract/name_map_test.go` 的路径改为 `repoRoot/num/<module>.py`。
- 验收：122 个 python 符号可导入、163 个 npz 全过、wheel 安装冒烟通过；`authority=python` 的 122 行不翻转。
