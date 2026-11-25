# 提案 0002：Python spec 的包布局与模块可导入方式

- 状态：**已定案（方案 A）**，2026-11-25 落地；结论已回写 `CONTRIBUTING.md` §3.1 / §3.2 与 `BUILDING.md` §4.5
- 影响面：`pyproject.toml`、Python 模块文件位置、`tests/` 的导入方式、`CONTRIBUTING.md` §3.1 / §3.2、
  `docs/PLAN.md` 的 `authority` 切换（go → python）
- 归属：新增权威实现前的结构决定，须走 §9.2

## 1. 背景

`CONTRIBUTING.md` 同时要求两件互相拉扯的事：

- §3.2：同一功能的各语言实现**放在同一目录、同名不同后缀**（`sum.go` / `sum.py` / `sum.h` / `sum.rs`），
  **不得**按语言拆目录把实现物理隔离；
- §3.1：仓库根只放构建入口，将来补 `pyproject.toml`。

若 Python 文件按 §3.2 平铺在仓库根（`sum.py` 与 `sum.go` 并列），那么「包」的目录就是仓库根本身。
Python 的导入模型要求包名与目录名一致，于是包名只能是仓库目录名 `num`，且**从源码运行时必须把仓库的
父目录放进 `sys.path`**（`import num` 而不是 `import sum`）。这带来三点麻烦：

1. 在仓库根直接 `python -c "import num"` 会失败，测试脚本要自己改 `sys.path`；
2. 仓库根同时是 Go 包、C++ 头文件目录、Python 包目录，`setuptools` 打包时要用
   `package-dir = {num = "."}` 这种反直觉映射；
3. 若有人把仓库目录改名，包名就变了。

## 2. 候选方案

### 方案 A：仓库根即包（`package-dir = {num = "."}`，推荐）

```toml
[tool.setuptools]
packages = ["num"]
package-dir = {num = "."}
```

- 文件位置：`sum.py` 与 `sum.go` 同目录（严格符合 §3.2）
- 导入：`import num` → `num.sum`；`num/__init__.py` 即根目录的 `__init__.py`
- 测试：用例里显式把**父目录**加入 `sys.path`（一行），或统一用
  `python -m pytest tests/python` 配合 `conftest.py` 里的同一行
- 代价：打包配置不直观；仓库目录名成为包名

### 方案 B：根目录平铺 `py-modules`

```toml
[tool.setuptools]
py-modules = ["sum", "mean", "dot"]
```

- 导入：`import sum` —— **直接污染顶层命名空间**（`sum` 是内建名），不可接受
- 结论：排除

### 方案 C：独立 `python/`（或 `num/`）子目录

- 结构清爽、`import num` 天生可用、`pyproject.toml` 常规写法
- 但**违反 §3.2**（把 Python 实现从同名文件旁挪走），跨语言对照要跳目录
- 若选它，必须同时修改 §3.2 的措辞，并接受「四门语言不再同目录并列」的长期影响

### 方案 D：根目录平铺 + 安装期重映射（`package-dir` 指向别处）

- 例如源文件在 `py/`、打包成 `num`：本质与 C 相同，只是多一层映射，不解决 §3.2

## 3. 推荐与影响

推荐 **方案 A**：它是唯一同时满足「同目录并列」与「可安装为 `num` 包」的写法，代价只在打包配置与测试导入，
且集中在两处（`pyproject.toml`、测试的 `sys.path` 一行）。

定案后需要同步：

1. `CONTRIBUTING.md` §3.1 / §3.2 补一段「Python 包即仓库根」的说明与导入约束；
2. 新增 `pyproject.toml`（`requires-python`、`packages`、`package-dir`）与根目录 `__init__.py`；
3. `BUILDING.md` 增加 Python 的构建 / 测试命令；
4. 首批实现的函数（建议从 `sum` / `mean` / `abs` / `dot1d` 起）与 `tests/fixtures` 的 npz 对拍，
   通过后才把 `NAME-MAP.tsv` 的 `authority` 由 `go` 改为 `python`（分批改，不一次切）。

## 4. 定案

1. **采用方案 A**：根目录即包，`packages = ["num"]` + `package-dir = { num = "." }`；
2. 最低 Python 版本 **3.11**；
3. Python 实现**不依赖 NumPy**，NumPy 只用于生成器与对拍用例。

落地物：根目录 `pyproject.toml` / `__init__.py` 与 `sum.py` / `mean.py` / `abs.py` / `dot.py` /
`cumsum.py` / `stddev.py` / `covariance.py`；`tests/python/test_conform.py` 读 16 个 npz 样本对拍。
`NAME-MAP.tsv` 的 `authority` 逐批切换，未切的语义仍以 Go 为临时基准。

## 5. 附：首批建议实现顺序

`sum` → `mean` → `abs` → `dot1d` → `cumsum` → `std`（总体）→ `var`（总体）。
七个函数覆盖标量归约、逐元素、累计与二阶统计，且 npz fixture 已有对应样本（`reduction` 组），
可以先做「Python 读 npz 自证」，再切 `authority`。
