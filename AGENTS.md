# AGENTS.md — AI 助手入口

> **本仓库的规范源是 [CONTRIBUTING.md](CONTRIBUTING.md)**：规则、约定与判据全部在那里。
> 本文件只做三件事 —— 告诉 AI 该读什么、哪些线不能踩、改完怎么验。**与 CONTRIBUTING.md 冲突时以
> 后者为准**；本文件不复制细则，只给指针。构建与依赖见 [BUILDING.md](BUILDING.md)。

## 1. 开工前必读

按顺序读，读完再动手：

| 顺序 | 读什么 | 为什么 |
| --- | --- | --- |
| 1 | [CONTRIBUTING.md](CONTRIBUTING.md) | 规范全集：语言角色、目录布局、命名、注释、文档、测试、提交、工作准则 |
| 2 | 目标目录的 `README.md` 与同名模块 `.md` | 该模块的口径、数据契约、跨语言对照、已知缺口（§5 TODO） |
| 3 | **NumPy 的对应语义 + Python spec 与用例** | 本仓对标 NumPy，算法与数值口径以它为准。**现阶段 Python 侧尚未落地**，此时以 NumPy 语义 + Go 参考实现 `Xxx_Go` 与同包用例作可执行对照（Go 是最早的一份实现，不是 spec） |
| 4 | 受影响的清单（`go.mod` / [BUILDING.md](BUILDING.md) 的包清单） | 新增包必须在本文件与 BUILDING.md 登记，见 BUILDING.md §3.4 |

## 2. 硬约束

| # | 约束 | 细则 |
| --- | --- | --- |
| 1 | **Python 是 spec**，语义、API、默认行为、失败约定一律以它为准；生产实现为 **C++**，Rust / Go 是对齐实现，Zig 为扩展，Java 为可选对照。**现阶段 Go 全落地、C++ 与 Rust 部分落地**，落地期间任何一侧都不得为自己加 Python 里没有的语义 | [§2.3](CONTRIBUTING.md#23-语言角色与演进顺序) |
| 2 | **不自己造东西**：不造 spec 里没有的结构 / 函数名 / 逻辑；需要新能力先补 spec，再推广 | [§4.1](CONTRIBUTING.md#41-总纲snake_case-优先) |
| 3 | **读 spec 代码，不读注释猜**；怀疑 spec 有错时先跑验证，再与用户商定是否改 spec | [§2.3](CONTRIBUTING.md#23-语言角色与演进顺序) |
| 4 | **类型名查词表**：数值与容器类型用 NumPy / BLAS / LAPACK 等国际通用术语，不自造别名、不用拼音 | [§4.3](CONTRIBUTING.md#43-领域类型词表type-vocabulary) |
| 5 | **注释只说"是什么 / 怎么做 / 为什么这样"**，不写"本轮改了啥、上一版是什么" | [§5.1](CONTRIBUTING.md#51-注释只说是什么怎么做) |
| 6 | **通用能力先查 `internal/`、`math32`、`math64`**，不写私有副本；缺就补到对应包 | [§5.2](CONTRIBUTING.md#52-通用能力优先落-internal) |
| 7 | **不擅自加开关**；调试开关一律 `debug_` 前缀；`UseAVX2` 只允许在用例里切换 | [§5.3](CONTRIBUTING.md#53-调试开关) |
| 8 | **同一批输入跑多遍产物一致**；导出 / 落盘按主键覆盖，不追加 | [§5.4](CONTRIBUTING.md#54-重复触发与幂等) |
| 9 | **文档与代码同批**：新目录带 `README.md`，跨语言模块带同名 `.md` | [§6](CONTRIBUTING.md#6-文档规范) |
| 10 | **编码与行尾**：UTF-8 无 BOM、LF、末尾一个换行 | [§5.8](CONTRIBUTING.md#58-文件编码与行尾) |
| 11 | **不擅自 `git commit` / `push` / `rebase`**；提交信息按 `type[language]：summary` 模板，Windows 下用 `-F` 传文件 | [§8](CONTRIBUTING.md#8-git-提交规范) |
| 12 | **对照表列序固定** `python → c++ → rust → go → zig → java`；未落地的列写 `—`，不删行删列 | [§2.3](CONTRIBUTING.md#23-语言角色与演进顺序) |
| 13 | **`temp/` 才是垃圾站**：一次性脚本、交叉编译产物、中间数据落 `temp/`，层级复刻源码结构（不许扁平堆在根），不入库、文档与提交信息也不引用其路径 | [§3.6](CONTRIBUTING.md#36-temp仓库的垃圾站) |
| 14 | **只在项目目录内写文件**：临时文件、导出物、调试脚本一律落 `temp/`；只有用户明确要求时才可写到项目外，并在回复里点名路径 | [§3.6](CONTRIBUTING.md#36-temp仓库的垃圾站) |
| 15 | **Go 包名铁律**：`package <目录名>`（`internal/functions/` → `package functions`） | [§3.3](CONTRIBUTING.md#33-go-的目录与包名) |
| 16 | **先建单**：每项功能或任务先用 `labs-git-issue` 建 Issue（不传 `--host` / `--repo`，按 `origin` 自动识别 Codeberg），完成后关单 | [§9.1](CONTRIBUTING.md#91-先建单再动手) |
| 17 | **版本号与 tag 一律不碰**：不得改版本字符串、不建 tag、不调 `autochangelog`、不手工改 `CHANGELOG.md` | [§10](CONTRIBUTING.md#10-版本与发布) |
| 18 | **加速实现的 API 名固定为 `Xxx_AVX2_F32/F64`**（跨架构同名：amd64 走 AVX2，arm64 走 NEON），不得改名为 `_NEON_*` / `_SIMD_*` | [§4.4](CONTRIBUTING.md#44-固定后缀与跨架构符号名) |
| 19 | **两套架构都要落地**：只落 amd64 或只落 arm64 视为未完成；`UseAVX2=false` 的 Go 回退是降级路径，不是交付 | [§5.6](CONTRIBUTING.md#56-向量化内核规范) |
| 20 | **归约顺序与向量化门槛不随手改**：乘积归约按 amd64 口径（f64 ≥16、f32 ≥32 才向量化），改了会让严格容差的用例漂移 | [§5.6.3](CONTRIBUTING.md#563-浮点归约的顺序敏感) |
| 21 | **`asm/` 是实验目录**（c2goasm / asm2plan9s 链路与样例）：不纳入支持面，其用例失败不算回归，也不许顺手"修好" | [§3.7](CONTRIBUTING.md#37-asm实验目录) |
| 22 | **C++ 侧单层扁平**：不设 `include/` / `src/`，文件与同名 `.go` 同目录；不建 `x32` / `x64` / `internal/functions` 的对应文件（模板 + xsimd 已覆盖两条宽度与向量化） | [§3.8](CONTRIBUTING.md#38-c单层扁平无-include--src) |
| 23 | **C++ 新增算子必须同时在 `CMakeLists.txt` 的 `num_core` 与 `num_cpp_tests` 里逐文件登记**，不用 `file(GLOB ...)`；`/utf-8` 与 `/MT` 两项编译设置不得删改 | [§3.8](CONTRIBUTING.md#38-c单层扁平无-include--src)、[BUILDING.md](BUILDING.md) §3.4 |

## 3. 主题索引

| 我要… | 去哪 |
| --- | --- |
| 查语言最低版本 / 编译器基线 | [§2.1](CONTRIBUTING.md#21-最低版本)、[§2.2](CONTRIBUTING.md#22-c-编译器基线) |
| 判断"同构 / 异构"这类术语怎么用（语义 vs 实现） | [§2.4](CONTRIBUTING.md#24-术语语义与实现) |
| 决定文件放哪、叫什么、包怎么登记 | [§3](CONTRIBUTING.md#3-目录与多语言布局) |
| 查跨语言导入 / 包含路径怎么对位 | [§3.1](CONTRIBUTING.md#31-仓库根与实现根) |
| 区分架构专属文件与构建标签 | [§3.4](CONTRIBUTING.md#34-架构专属文件) |
| 临时脚本 / 一次性复现 / 中间产物放哪 | [§3.6](CONTRIBUTING.md#36-temp仓库的垃圾站) |
| 定函数名 / 命名空间 / 平台宏 | [§4](CONTRIBUTING.md#4-命名与数据契约) |
| 查领域类型标准名（`Array` / `DType` / `Vector` …） | [§4.3](CONTRIBUTING.md#43-领域类型词表type-vocabulary) |
| 定整型 / 浮点的跨语言名与宽度 | [§4.6](CONTRIBUTING.md#46-标量类型) |
| 定落盘字段名 / Go struct tag | [§4.7](CONTRIBUTING.md#47-边界态io-与序列化铁律) |
| 写注释、判断该不该写某句 | [§5.1](CONTRIBUTING.md#51-注释只说是什么怎么做) |
| 找某个通用能力该不该落 `internal/` | [§5.2](CONTRIBUTING.md#52-通用能力优先落-internal) |
| 处理重复触发 / 幂等 | [§5.4](CONTRIBUTING.md#54-重复触发与幂等) |
| 跑 Go、写 import | [§5.5](CONTRIBUTING.md#55-go-包形态与导入) |
| 写向量化内核 / 加一段汇编 | [§5.6](CONTRIBUTING.md#56-向量化内核规范) |
| 遇到缺依赖 / spec 未定 | [§5.9](CONTRIBUTING.md#59-卡点先抛不支持) |
| 写模块文档 / 判"算不算完成" | [§6](CONTRIBUTING.md#6-文档规范) |
| 写测试 | [§7](CONTRIBUTING.md#7-测试规范) |
| 写提交信息 | [§8.1](CONTRIBUTING.md#81-提交信息) |
| 判断要不要立提案 | [§9.2](CONTRIBUTING.md#92-提案rfc变更受理与裁决) |
| 定工作方式（假设 / 简洁 / 精准 / 目标驱动） | [§11](CONTRIBUTING.md#11-工作准则) |
| 构建 / 依赖 / 发布 | [BUILDING.md](BUILDING.md) |

## 4. 常用命令

| 目的 | 命令 |
| --- | --- |
| Go 编译 / 用例 / 检查 | `go build ./...`、`go test ./...`、`go vet ./...`、`gofmt -l .` |
| Go 跑核心包用例 | `go test ./internal/functions/ ./x32/ ./x64/ ./vectors/ ./kind/ .` |
| Go 基准 | `go test -bench . -benchmem ./internal/functions` |
| Go 覆盖率 | `bash go.test.sh`（Linux / WSL / Git Bash） |
| arm64 交叉验证 | 交叉编译 `go test -c` 后用 `qemu-aarch64-static` 执行，见 [BUILDING.md](BUILDING.md) §5.6 |
| 编码自查（BOM / CRLF） | [§5.8](CONTRIBUTING.md#58-文件编码与行尾) 的 PowerShell 片段 |
| 提交（Windows） | 写入 `.git/commit-msg.tmp` 后 `git commit -F .git/commit-msg.tmp` |
| 建 / 查 / 关工单 | `labs-git-issue`（不传 `--host` / `--repo`，按 `origin` 自动识别），正文模板见 [§9.1](CONTRIBUTING.md#91-先建单再动手) |
| C++ 取依赖 / 配置 / 编译 | `conan install . -s build_type=Release -s compiler.runtime=static`、`cmake --preset release-static`、`cmake --build --preset release-static` |
| C++ 用例 | `ctest --preset release-static`（230 个用例，Debug 用 `debug-static`） |
| Rust 编译 / 用例 / 检查 | `cargo build`、`cargo test`、`cargo clippy --all-targets -- -D warnings`、`cargo fmt --check` |
| Python / Zig / Java | **尚未落地**：命令随各语言落地时补进 [BUILDING.md](BUILDING.md) 对应章节，此处不写未验证的命令 |

## 5. 交付前自检

1. **语义**：与基准实现（现阶段是 Go 的 `Xxx_Go`；C++ 侧按 `num.md` §2 的口径逐条核对）逐值对拍过（不是"看起来一样"）；f32 / f64 两条都覆盖。
2. **架构**：amd64 与 arm64 两侧都落地且都能编译；新增内核在 Go 声明与 `.s` 符号两处都登记。
3. **构建与用例**：`go build ./...`、`ctest --preset release-static`、`cargo test` 通过；跳过项（如需 arm64 硬件 / qemu）要能说明原因。
4. **一致性**：新增 / 修改的接口在已落地的语言里命名对位；`_AVX2_` 后缀名未被改动；C++ 的关键字避让（`stddev` / `logic_and` 等）登记在对照表。
5. **注释**：没有"本轮 / 上一版 / 原来 / 曾经"这类变更叙事（[§5.1](CONTRIBUTING.md#51-注释只说是什么怎么做)）。
6. **文档**：目录 `README.md`、模块 `.md` 的路线图 / 功能表 / 测试 / TODO 与实际一致，无陈旧条目。
7. **编码**：UTF-8 无 BOM、LF、末尾换行（[§5.8](CONTRIBUTING.md#58-文件编码与行尾)）。
8. **提交信息**：符合 `type[language]：summary`，正文分条含验证命令，无版本号 / 代理署名。
9. **工单**：这轮任务有对应 Issue，状态与正文和实际落地一致（[§9.1](CONTRIBUTING.md#91-先建单再动手)）。
