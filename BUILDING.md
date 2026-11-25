# num 构建指南

本文说明 num 仓库的依赖、构建、测试、跨架构验证与产物位置。

num 是**多语言平行实现**的数值计算库：Python / C++ / Rust / Go 各自独立构建、独立运行，**不存在
FFI / Cgo / SWIG 等跨语言调用**。因此每种语言有自己的一条构建链路，互不依赖；跨语言之间只通过同名接口、
统一的数据契约和同一组黄金测试向量保持一致。

**现阶段 Go 全落地、C++ 与 Rust 部分落地**：本文 §2 ~ §5 是可执行的 Go 部分，§6.2 是可执行的 C++ 部分，
§6.3 是 Rust 部分；其余语言章节保留骨架，命令随落地补齐。

| 构建单元 | 构建入口 | 产物 | 当前状态 |
| --- | --- | --- | --- |
| Go | `go.mod` | 库（无命令行入口） | **已落地** |
| 汇编内核（amd64 / arm64） | 随 Go 包编译 | `.s` 编进库 | **已落地** |
| Python | `pyproject.toml`（未建） | 包 | 未落地 |
| C/C++ | `conanfile.py` + `CMakeLists.txt` + `CMakePresets.json` | 静态库 `num_core` + ctest 用例 | 部分落地（59 个功能行） |
| Rust | `Cargo.toml` | 库 + 160 个 fixture 的对拍用例 | **已落地（fixture 覆盖子集）** |
| Zig / Java | 无清单文件 | — | 未落地（扩展语言） |

---

## 1. 通用前置条件

### 1.1 版本要求

| 语言 | 版本来源文件 | 当前要求 |
| --- | --- | --- |
| Go | `go.mod` | **1.27.3**（以 `go` 指令为准） |
| Python | `pyproject.toml`（未建） | 落地后以 `requires-python` 为准（与 quant1x 同口径：3.8+，建议 3.12+） |
| C/C++ | `CMakeLists.txt`（`cmake_minimum_required`） | **C++20**（`CMAKE_CXX_STANDARD 20`，`CXX_EXTENSIONS OFF`）；**CMake ≥ 3.30**；Windows 用 MSVC v143（19.30+，统一 `/MT` 静态 CRT），Linux 用 GCC 13+，macOS 用 AppleClang。基线不符**配置期**直接 `FATAL_ERROR`，不是告警 |
| Rust | `rust-toolchain.toml` | **1.100.0**，Edition 2024 |
| Zig | 无清单文件 | 落地后按 0.17+ |
| Java | `pom.xml`（未建） | 落地后按 Java 8 字节码目标 |

- `README.md` 的徽章可能滞后（当前写 1.21.8+），**一律以 `go.mod` 为准**；升级 Go 版本时同步
  `go.mod`、[CONTRIBUTING.md](CONTRIBUTING.md) §2.1 与 README。

### 1.2 目录与产物约定

| 路径 | 性质 | 说明 |
| --- | --- | --- |
| `temp/` | 临时产物 | 验证脚本、交叉编译二进制（`*.test`）、导出数据；层级复刻源码结构，不入库 |
| `coverage.txt` / `coverage.out` | 覆盖率产物 | `go.test.sh` 生成，已在 `.gitignore` |
| `*.test` | 测试二进制 | `go test -c` 的产物，已在 `.gitignore`；不要留在包目录 |
| `build/` | C++ 构建产物 | conan 的 `generators/` 与 CMake 的中间件、用例二进制；**整目录不入库**（已在 `.gitignore`） |
| `tests/fixtures/read_rs/target/` | Rust 工具 crate 产物 | C3 读取器的 cargo 构建目录，已在 `.gitignore` |
| `target/` | Rust 根 crate 产物 | `cargo build/test` 的输出，已在 `.gitignore` |
| `CMakeUserPresets.json` | C++ 预设台账 | conan 默认写在仓库根的台账文件，本仓不用（预设收进各自的 `build/<配置>/generators/`），已在 `.gitignore` |
| Go build cache / `GOMODCACHE` | 工具缓存 | 由工具管理，不手工清理挪动 |

临时验证脚本与一次性产物**放 `temp/`**，不要污染源码树；产物与产生它的脚本同目录。

---

## 2. Go 构建

### 2.1 模块与依赖

- 模块路径 `github.com/quant1x/num`（`go.mod`）。
- 运行时依赖只有两个：`github.com/stretchr/testify`（用例断言）与 `golang.org/x/sys/cpu`
  （运行时指令集探测，`cpu_info.go`）。其余为它们的间接依赖。
- **不提交 `vendor/`**；依赖升级用 `go get <module>@<version>` + `go mod tidy`，并在提交信息里写明原因。

### 2.2 构建

```bash
# 编译全部包（含 internal/*、x32、x64、math32、asm/* 样例）
go build ./...

# 只编译公开包
go build . ./x32/ ./x64/ ./math32/ ./math64/ ./binary/ ./kind/ ./vectors/
```

库代码不产出命令行程序；需要试算时写临时用例（`temp/` 下）或用 `go test -run`。唯一的命令行入口是
`cmd/numctl`（CSV 交换与统计），构建方式：

```bash
CGO_ENABLED=0 go build -o numctl ./cmd/numctl
```

平台与架构：

```bash
# 交叉编译（示例：linux/arm64）
GOOS=linux GOARCH=arm64 go build ./...
```

- 加速路径由**手写汇编**承担，不依赖编译器自动向量化，因此**不需要**设置 `GOAMD64=v3`；
  保持默认 `GOAMD64=v1` 即可，运行时按 CPU 特性选路（见 §5.1）。
- 汇编内核按架构分文件（`*_amd64.s` / `*_arm64.s`），交叉编译时只会编进目标架构的那一份。

### 2.3 依赖代理

网络受限时可通过环境变量配置：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=github.com/quant1x
```

---

## 3. 包清单与新增包

### 3.1 公开包

| 包 | 路径 | 职责 |
| --- | --- | --- |
| `num` | 仓库根 | NumPy 风格门面 API（`Sum` / `Dot` / `Matrix` / `Rolling` / `Repeat` / `Where` …），按 `DType` 分派 |
| `x32` / `x64` | `x32/` `x64/` | `float32` / `float64` 强类型 API：校验 + `Op` / `Op_Inplace` / `Op_Into` 三种形态，按 `UseAVX2` 选路 |
| `math32` | `math32/` | `float32` 标量数学（含 `exp_amd64.s` / `exp_arm64.s` / `sqrt_*.s` 等） |
| `math64` | `math64/` | `float64` 位级与指针工具 |
| `binary` | `binary/` | 切片头与内存视图（`SliceHeader` / endian） |
| `kind` | `kind/` | 类型判别与默认值 |
| `vectors` | `vectors/` | 小工具集 |
| `io` | `io/` | D1 CSV 读写：schema、缺失值、ISO 8601、幂等 |
| `labs` | `labs/` | 用例辅助（testify 适配） |

### 3.2 内部包

| 包 | 路径 | 职责 |
| --- | --- | --- |
| `functions` | `internal/functions/` | 内核层：`Xxx_Go` 参考实现、`Xxx_AVX2_*` 加速实现、`accel_*` 探测与声明、汇编内核 |
| `partial` | `internal/partial/` | TopK / 部分排序 |
| `constraints` | `internal/constraints/` | 泛型约束 |
| `rand` | `internal/rand/` | 可复现随机数（用例专用） |

### 3.3 实验目录

`asm/`（`asm2plan9` / `c2goasm` 的 fork、`demo`、`src` 与生成样例）**不属于支持面**，见
[CONTRIBUTING.md §3.7](CONTRIBUTING.md#37-asm实验目录)。

### 3.4 新增包的登记

新增一个包时，同一次提交里改三处：

1. 本节 §3.1 / §3.2 加一行（路径 + 职责 + 边界）；
2. 该目录补 `README.md`（[CONTRIBUTING.md §6.1](CONTRIBUTING.md#61-每个目录必须有-readme-md)）；
3. 若它属于新的一层（如将来的 `tests/`、`docs/`），同步 [CONTRIBUTING.md §3.1](CONTRIBUTING.md#31-仓库根与实现根)。

C++ 侧另有一条：新增算子必须在 `CMakeLists.txt` 里**逐文件显式登记**（`num_core` 的源清单与
`num_cpp_tests` 的用例清单各一处），不用 `file(GLOB ...)` —— glob 的结果依赖配置时刻的目录快照，
新增文件不重配就静默漏编。

### 3.5 C++ 目标

| 目标 | 类型 | 内容 | 说明 |
| --- | --- | --- | --- |
| `num_core`（别名 `num::num_core`） | `STATIC` | 仓库根的 `*.h` + `decimals.cpp` | 门面实现。以 header-only 模板为主，只有非模板实现（`decimal`）落编译单元 |
| `num_cpp_tests` | `EXECUTABLE` | 仓库根的 13 个 `*_test.cpp` | 链接 `num_core` 与 `GTest::gtest_main`，由 `gtest_discover_tests` 逐 `TEST` 注册到 ctest |

包含根是**仓库根的上一级**，故公共头的包含形态为 `<num/<模块>/<头文件>.h>`（根目录的 `sum.h` →
`#include <num/sum.h>`）；与 [CONTRIBUTING.md §3.1](CONTRIBUTING.md#31-仓库根与实现根) 登记的三段式对位，
将来补安装导出时布局照此形（`include/num/...`）。

`num_fixture_check` 是 C3 的 fixture 读取器（`tests/fixtures/read_cpp/`），只依赖标准库、不链接
`num_core`；`NUM_BUILD_TESTS=ON` 时由 ctest 以 `fixtures_npz_read_cpp` 注册。

### 3.6 验收目录（`tests/`）

`tests/` 不属于实现包，收跨语言验收资产；Go 侧只有 `tests/contract` 是包（校验用例），其余是
工具与数据。

| 路径 | 形式 | 职责 |
| --- | --- | --- |
| `tests/contract` | Go 包 `contract` | C1 契约表与覆盖校验（`go test ./tests/contract/`） |
| `tests/fixtures/gen` | Python 脚本 | C2 生成器，需要 NumPy |
| `tests/fixtures/conform` | Go 包 `conform` | M1 Go 对拍用例：读样本 → 调 Go 实现 → 逐值比较 |
| `tests/fixtures/read_go` | Go 包 `read_go` | C3 Go 读取器与用例 |
| `tests/fixtures/read_cpp` | C++ 源文件 | C3 C++ 读取器，目标是 `num_fixture_check` |
| `tests/fixtures/conform_cpp` | C++ 源文件 | C++ 对拍，目标是 `num_fixture_conform`（未实现语义跳过） |
| `tests/fixtures/read_rs` | 独立 crate | C3 Rust 读取器，不进生产 crate |
| `tests/io` | Python + 三个工具 | D2 互读互写与 B5 按主键合并，见 `tests/io/README.md` |

### 3.7 命令行入口

| 目标 | 路径 | 职责 | 构建 |
| --- | --- | --- | --- |
| `numctl` | `cmd/numctl/` | CSV 往返 / 按主键合并 / 数值统计，单二进制、无 Python 运行时 | `CGO_ENABLED=0 go build -o numctl ./cmd/numctl` |

新增验收目录时同批补 `README.md`，并在本节登记。

---

## 4. 测试与回归判定

### 4.1 跑法（回归判定的准绳）

```bash
# 核心包：日常回归用这一组，必须全绿
go test ./internal/functions/ ./internal/partial/ ./x32/ ./x64/ ./vectors/ ./kind/ ./binary/ . 

# 契约与 fixture（C1 / C3）
go test ./tests/contract/ ./tests/fixtures/conform/ ./tests/fixtures/read_go/

# CSV 读写（D2）与数据管道（B5）
go test ./io/
python tests/io/roundtrip.py   # 需先构建 num_csv_roundtrip，见 tests/io/README.md
python tests/io/pipeline.py    # 需先构建 num_csv_merge

# 单条用例
go test ./internal/functions/ -run TestSumProd -v

# 静态检查（支持面；asm/ 是实验目录，见 §4.2）
go vet . ./binary/ ./math32/ ./internal/functions/ ./x32/ ./x64/ ./vectors/ ./kind/ ./io/ ./cmd/numctl/
gofmt -l .
```

C++ 侧：

```bash
ctest --preset release-static      # 247 个用例，Release + /MT
ctest --preset debug-static        # 同一批用例，Debug + /MTd
```

**不要用 `go test ./...` 的失败项下结论**：`asm/` 是实验目录，它的用例本来就是红的（见 §4.2）。
判断是否回归，只能看上面这组明确包列表。

### 4.2 已知失败与存量告警

`asm/` 实验目录的既有问题（**不是回归**，也不许顺手修）：

| 位置 | 现象 |
| --- | --- |
| `asm/c2goasm/cgocmp` | `undefined: CgoMultiplyAndAdd`（用例引用了未导出的 C 侧符号），包级 `go test` 直接 build failed |
| `asm/c2goasm/test` | `*_test.go` 引用未定义符号；`clib_amd64.s` 报 `wrong argument size` |
| `asm/asm2plan9` | `TestLongInstruction` 等断言与当前反汇编输出不一致 |
| `asm/floats_avx_test.go` | `undefined: ___mm256_mul_to`；`go vet` 报 `possible misuse of unsafe.Pointer` |

`go vet` 在支持面（除 `asm/` 外的全部包）已无告警：根包 5 处 `unreachable code`、
`binary` 3 处 `SliceHeader` 误用、`math32/exp.go` 的 `unreachable code` 已清理。
仍存在的是跨架构汇编声明问题：**arm64 目标下** `math32/exp_arm64.s` / `Exp2` 报 `wrong argument size`，
改动 `math32` 时一并核对。

`gofmt -l .` 目前无输出 —— 出现输出就是新引入的问题，提交前修掉。

### 4.3 覆盖率

```bash
bash go.test.sh     # Linux / WSL / Git Bash；逐包 -race -coverprofile，汇总到 coverage.txt / coverage.out
```

Windows PowerShell 下若没有 bash，用等价命令：

```powershell
go test -race -coverprofile=coverage.out -covermode=atomic ./internal/functions/ ./x32/ ./x64/ .
go tool cover -func=coverage.out
```

汇编内核是覆盖率盲区（`.s` 不计入行覆盖），内核的正确性靠 [§4.1](#41-跑法回归判定的准绳) 的对拍用例
与 [§5.6](#56-跨架构验证) 的跨架构执行保证，不要用覆盖率当内核的验收指标。

### 4.4 基准

```bash
go test -bench . -benchmem ./internal/functions
go test -bench BenchmarkSumProd -benchmem ./internal/functions
```

`internal/functions` 的基准覆盖 `sizes = {10, 100, 1_000, 10_000, 100_000}`。约定：

- 基线对比只允许**同口径实现**之间进行（`go_*` vs `avx2_*`），不要把自己临时改的变体混进仓库基准；
- 汇编改动要给出改动前后的同一命令输出（含 `-benchmem`），否则无法判断是真实收益还是噪声；
- 基准在目标架构上跑，`arm64` 的结果必须注明运行方式（真机 / qemu ——qemu 的绝对耗时不可信，
  只看相对趋势）。

---

## 5. 汇编内核与架构

### 5.1 架构矩阵与运行时选路

| 架构 | 加速实现 | 判定条件 | 落点 |
| --- | --- | --- | --- |
| `amd64` | AVX2 + FMA | `cpu.X86.HasAVX2 && cpu.X86.HasFMA && runtime.GOOS != "darwin"` | `accel_amd64.go` |
| `arm64` | NEON（ASIMD） | `cpu.ARM64.HasASIMD` | `accel_arm64.go` |
| 其余 | 无 | 恒为 `false` | `accel_noasm.go`（`//go:build !amd64 && !arm64`） |

- `UseAVX2` 是包级布尔值，**非线程安全**，只在初始化与用例里改写；`SetAcceleration(enabled)` 只允许
  "关"或在支持的平台上"开"，反过来的请求会被拒绝：
  `if enabled && !simdSupported() { UseAVX2 = false; return }`。
- 门面层通过 `SetAvx2Enabled` / `GetAvx2Enabled` 暴露同一个开关（仅供调试与用例）。
- `UseAVX2 == false` 时，`x32` / `x64` 全部改调 `Xxx_Go`；其余平台不会触达 `accel_noasm.go` 的
  `panic("not implemented")` 桩。
- `darwin/amd64` 被显式排除（历史原因：macOS 上的 AVX2 路径未验证），改动前需先补验证。

### 5.2 amd64 内核的生成链路

`internal/functions/accel_avx2_amd64.s` 与 `math32/*_amd64.s` 属于**生成的汇编**，来源是 C 源码经
clang 编成 Intel 语法汇编，再由 `c2goasm` 转成 Plan 9 语法：

```bash
# 1) c2goasm（工具在 asm/ 下有 fork，也可 go install github.com/minio/c2goasm@latest）
# 2) C 源码转 cpp 后生成汇编（c2goasm 对纯 C 支持不佳）
cp src/floats_avx.c src/floats_avx.cpp
clang -O3 -mavx -mfma -masm=intel -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti \
      -S src/floats_avx.cpp -o src/floats_avx.s
c2goasm -a src/floats_avx.s floats_avx_amd64.s
```

- 完整参数（`-mno-red-zone -mstackrealign -ffast-math` 等）见 `asm/makefile`。
- **本仓没有 `internal/functions` 内核的 C 源与生成脚本**：`asm/src/floats_avx.c` 只是上述链路的样例
  （导出 `asm/floats_avx_amd64.s`），不含 `Accel` 这一批内核。因此 `accel_avx2_amd64.s` 目前按
  **手工维护**处理 —— 其中 `Find_AVX2_*` / `Maximum_AVX2_*` 等本就是手写内核（标签是 `loop` / `tail`
  而不是 clang 的 `LBB*`）。改动时按 §5.5 验证，并在提交信息里说明改的是内核本身。
- 生成的汇编自带标量尾部处理，任意长度都能直接调用（与 arm64 内核的分工不同，见 §5.3）；**标量尾部与
  向量段受同一套特殊值口径约束**（[CONTRIBUTING.md §5.6.3](CONTRIBUTING.md#563-浮点归约的顺序敏感)），
  即 `NaN` 参与比较一律 false、`±Inf` 与 `±0` 按 IEEE，`ucomis*` 之后的跳转必须把 unordered（PF=1）算进去。
- 符号命名由 c2goasm 决定（`·Xxx_AVX2_F64(SB)`），Go 侧声明在 `accel_avx2_amd64.go`，两侧必须同名。

### 5.3 arm64 NEON 内核（手写）

`internal/functions/accel_neon_arm64.s` 是**手写**的 Plan 9 汇编，约定与 amd64 不同：

- **只处理整块**：f64 每块 4 个、f32 每块 8 个；`n` 必须是块大小的整数倍，不满一块的尾部由 Go 包装层
  （`accel_neon_arm64.go`）补齐。
- 每个内核一组寄存器约定：累加器与广播值不得共用寄存器（矩阵乘法用 `V20` / `V21` 作累加器，输入广播用
  `V4` / `V12`）。
- 内核入口参数按 ABI0 布局读取（`MOVD a+0(FP), R0`），帧大小 `$0-<字节数>` 必须与实际参数 + 返回值
  一致；Go 侧声明加 `//go:noescape`。
- 常量表用 `DATA` + `GLOBL ... RODATA` 声明（`·neonOneF64` / `·neonOneF32`），不要在内核里现算。
- 验证必须跑 `GOARCH=arm64` 的用例（[§5.6](#56-跨架构验证)）；x86 机器上无法直跑。

### 5.4 `asm/` 实验链路的结论

| 结论 | 说明 |
| --- | --- |
| c2goasm 只适合**无异常 / 无 RTTI / 无内建函数**的纯计算代码 | 生成物里没有栈展开信息，混入运行时依赖会崩 |
| c2goasm 对纯 C 支持不佳 | 需先转 `.cpp` 再编（`asm/README.md`） |
| asm2plan9s 已停止维护 | 与当前反汇编输出存在断言级差异，不要拿它做正确性判据 |
| 有生成源时不要手改生成物 | 见 §5.2；要改就改源与脚本。`internal/functions/accel_avx2_amd64.s` 本仓无生成源，按手工维护 |

### 5.5 修改内核后的必查项

1. `go build ./...` 与 `go vet ./internal/functions/` 通过（汇编参数大小写错会在这里暴露）。
2. `go test ./internal/functions/` 全绿（对拍用例覆盖 `Xxx_Go` vs `Xxx_AVX2_*`）。
3. 长度边界：`1, 块大小-1, 块大小, 块大小+1, 2*块大小` 都要跑到（NEON 尾部处理的高发区）。
4. 两种架构都验：amd64 直跑 + arm64 交叉编译经 qemu 执行。
5. 改到归约（`Sum` / `Prod` / `Dot` / `Norm`）时，确认累加顺序与门槛未被顺手改动
   （[CONTRIBUTING.md §5.6.3](CONTRIBUTING.md#563-浮点归约的顺序敏感)）。

### 5.6 跨架构验证

x86 主机上验证 arm64 实现：**交叉编译 + qemu-user**。以 WSL 的 Debian 为例（需 `qemu-user-static`）：

```bash
# 1) 主机侧交叉编译测试二进制，产物放 temp/（不要落系统临时目录）
GOOS=linux GOARCH=arm64 go test -c -o temp/internal/functions/arm64/functions.test ./internal/functions

# 2) 在 WSL 里用 qemu 执行
wsl -d Debian bash -lc "cd /mnt/d/projects/quant1x/num && chmod +x temp/internal/functions/arm64/functions.test \
  && qemu-aarch64-static temp/internal/functions/arm64/functions.test"
```

PowerShell 写法（等价）：

```powershell
$env:GOOS = "linux"; $env:GOARCH = "arm64"
go test -c -o temp/internal/functions/arm64/functions.test ./internal/functions
wsl -d Debian bash -lc "cd /mnt/d/projects/quant1x/num && chmod +x temp/internal/functions/arm64/functions.test && qemu-aarch64-static temp/internal/functions/arm64/functions.test"
```

注意事项：

- `-test.v` / `-test.run <正则>` 照常可用（如 `... functions.test -test.run TestSumProd -test.v`）。
- **qemu 下的耗时不准**，不要用它出基准结论；只用于功能验证。
- 验证完删掉 `temp/` 下的二进制；不要提交。
- 需要真机结论（性能、指令退化的实际影响）时，在 arm64 机器上直接 `go test`，并在结论里注明环境。

---

## 6. 其它语言

各语言的构建章节随着落地补齐，骨架与口径先固定：

### 6.1 Python

落地后：`python -m venv .venv` → `pip install -e .` → `python -m pytest tests -q`。
Python 是 spec 锚点，**新语义先在这里定案**。

### 6.2 C/C++

依赖走 Conan 2（`conanfile.py`：`xtensor/0.27.1` 带 `xsimd`，`gtest/1.18.0` 走 `test_requires`），构建走
CMake 预设（`CMakePresets.json`）。Windows 用 MSVC + `/MT`，Linux 用 GCC 13+，C++20。

```bash
# 1) 取依赖并生成 toolchain：两套配置各一次
conan install . -s build_type=Release -s compiler.runtime=static
conan install . -s build_type=Debug   -s compiler.runtime=static

# 2) 配置 + 编译
cmake --preset release-static
cmake --build --preset release-static

# 3) 用例
ctest --preset release-static

# 4) 安装导出（D3）：静态库 + include/num + CMake 包
cmake --install build/release-static --prefix dist/cpp

# 5) 编辑器索引用的编译数据库（可选，见下条）
cmake --preset ninja-compile-db
```

- 安装包版本从**精确匹配的 git tag** 派生（`vX.Y.Z` → `X.Y.Z`），未打 tag 为 `0.0.0`；
  详见 [docs/RELEASE.md](docs/RELEASE.md)。

- 预设名与 `conanfile.py` 的 `layout()` 一一对应：`Release + static` → `build/release-static`，
  `Debug + static` → `build/debug-static`。`compiler.runtime` 必须与 `CMakeLists.txt` 固化的
  `CMAKE_MSVC_RUNTIME_LIBRARY` 同档，否则链接期报 `LNK2038`（gtest 与 `num_core` 的 CRT 不一致）。
- Linux / macOS 上去掉 `-s compiler.runtime=static`（无 CRT 之分），预设名与目录不变。
- 不加 `--build=missing` 时只用缓存里的二进制包；本机 `CONAN_HOME` 已有三者的静态 CRT 包，可离线完成。
- `conan install` **不要**把预设写回仓库根的 `CMakeUserPresets.json`：VS 是多配置生成器，两套构建目录的
  预设同名（`conan-default`），先后装两个 `build_type` 后 CMake 一读档就报同名冲突。本仓已让 conan 把台账
  收进各自的 `build/<配置>/generators/`，入口只有根目录的 `CMakePresets.json`（直接指向
  `conan_toolchain.cmake`）。改 `conanfile.py` 的 `generate()` 时不要动这一行。
- `ninja-compile-db` 是**只为编辑器索引**的预设：Ninja 生成器 + 复用 `build/release-static/generators`
  的 conan 配置文件，产出 `build/ninja-compile-db/compile_commands.json`，不出产物、不影响 VS 预设。
  VS 生成器不产出编译数据库，clangd 拿不到包含路径就会把 `xtensor/containers/xtensor.hpp` 报成
  "file not found"。仓库根 `.clangd` 已指向该目录，VSCode / Zed 用 clangd 时无需另配包含路径。
  MSVC 要在 VS 开发者环境里跑（`Enter-VsDevShell`，或 VSCode CMake Tools 的 kit），否则 `cl.exe`
  不在 PATH、配置阶段就报 `CMAKE_CXX_COMPILER not set`。依赖或头文件变动后重跑这条即可。
- 新增算子按 §3.4 的 C++ 条款在 `CMakeLists.txt` 显式登记。
- 打包（`conan create` / `cmake --install` 导出）与 C++ 基准（`benchmark`）未落地，见 `num.md` §5 TODO。

### 6.3 Rust

版本由 `rust-toolchain.toml` 锁定（1.100.0 / Edition 2024），根 crate 为 `quant1x-num`（`publish = false`）。

```bash
cargo build
cargo test                                   # 含 tests/conform.rs：160 个 fixture 对拍
cargo clippy --all-targets -- -D warnings
cargo fmt --check
```

- `tests/fixtures/read_rs` 与 `tests/io/rs` 是同仓的**工具 crate**，各自有 `Cargo.toml`，用
  `--manifest-path` 单独构建，不参与根 crate 的发布。
- 根 `Cargo.toml` 的 `version` 是 Cargo 必填的占位 `0.0.0`；正式版本由 tag / autochangelog 注入。
- 语义以 Go 临时基准为准；与 Go 的已知差异（`median` 不排序、`fill` 连 Inf 填、`concat1d` n×3）
  照搬，不按 NumPy 口径"修正"。

### 6.4 Zig / Java

扩展语言，只承载小功能与对照实现；Zig 无 `build.zig` 时直接 `zig test <文件>`，Java 用 Maven wrapper。

---

## 7. 跨语言一致性验证

同一输入在各语言实现上应给出 epsilon 内一致的结果：

1. 输入与黄金输出在 `tests/fixtures/`（期望值由 Python 生成）；
2. Go 用 `tests/fixtures/conform`、Rust 用 `tests/conform.rs` 读同一批 npz 并逐值对拍；
3. 浮点按 fixture 的 `meta_tolerance` 比较（f64 默认 1e-12、f32 1e-6），布尔 / 整型严格相等。

Go 侧另有**包内对拍**：`Xxx_Go` 作为期望值，`Xxx_AVX2_F32/F64` 作为被测对象
（`internal/functions/*_test.go`）。C++ 的计算侧对拍已接入 `tests/fixtures/conform_cpp`（160 个样本全量通过），
剩余差异见 `tests/fixtures/CPP-GAPS.md`。

---

## 8. 版本与发布约束

- 版本唯一由 git tag `vMAJOR.MINOR.PATCH` 决定，各语言共享同一版本号。
- `autochangelog` 是唯一发布入口：tag 自增、主次版本递增、需要手改版本号的技术栈配置、`CHANGELOG.md`
  生成，全部由它完成。**AI 与协作者一律不碰**（[CONTRIBUTING.md §10](CONTRIBUTING.md#10-版本与发布)）。
- `go.mod` 没有 `version` 字段，Go 侧不参与版本号维护 —— 不要为了"对齐版本"往 Go 代码里塞版本常量。

---

## 9. 常见问题

| 现象 | 原因与处置 |
| --- | --- |
| `go test ./...` 报 `asm/*` 失败 | 实验目录的既有问题，不是回归；用 §4.1 的包列表判定 |
| `go vet ./...` 报 `asm/*` 的 `unsafe.Pointer` | 实验目录的存量告警（§4.2）；支持面用 §4.1 的包列表 |
| `go vet` 在 arm64 目标上报 `math32/exp_arm64.s` 参数大小 | 存量问题，汇编参数声明与实际不符；改动 `math32` 时一并核对 |
| 交叉编译的 `.test` 在 Windows 上跑不了 | 需要 qemu（[§5.6](#56-跨架构验证)）；Windows 自身无法执行 ELF |
| `go test -c` 产物进了源码目录、被 `git status` 看到 | `*.test` 已在 `.gitignore`，但仍应输出到 `temp/` 并清理 |
| 中文注释显示成「姹傚拰」 | 历史遗留的编码损坏（UTF-8 → GBK 误解码）；新增注释写正常中文，只在正在改的函数内顺手修 |
| README 徽章说 1.21+，`go.mod` 是 1.27.3 | 徽章滞后，以 `go.mod` 为准 |
| Windows 下 `git commit -m "中文"` 变乱码 | 改走 `-F` 文件（[CONTRIBUTING.md §8.2](CONTRIBUTING.md#82-windows-下提交信息的传递必须走文件)） |
| C++ 编译出成片语法错误（`C4819` 后面跟着 `C4430` / `C2143`） | 源文件是无 BOM 的 UTF-8，而中文 Windows 上 MSVC 默认按 GBK 读，中文注释被切碎并吃掉换行。`CMakeLists.txt` 已固化 `/utf-8`，别删；手写编译命令时也要带 |
| `cmake --preset` 报 `Invalid "configurePreset": "conan-debug"` 或预设同名冲突 | 仓库根出现了 conan 写的 `CMakeUserPresets.json`（把两套构建目录的 `conan-default` 都 include 进来了）。删掉该文件，并按 §6.2 让 conan 把台账收进构建目录 |
| `ctest --preset <名字>` 报找不到 test preset | 先跑同名的 `cmake --preset` 配置；预设的 `binaryDir` 不存在时 ctest 没有可读的 CTest 清单 |
| conan 报 `unexpected keyword argument 'auto_refresh'` 之类 | 依赖版本与预期不符：本仓只钉 `xtensor/0.27.1` / `gtest/1.18.0`，先看 `conanfile.py` 而不是改调用方 |

---

## 10. 相关文档

- [CONTRIBUTING.md](CONTRIBUTING.md) — 规范源：分层、命名、内核契约、测试、提交
- [docs/RELEASE.md](docs/RELEASE.md) — 编译产物、安装导出与版本派生（D3）
- [AGENTS.md](AGENTS.md) — AI 入口：读什么、硬约束、交付前自检
- [README.md](README.md) — 项目介绍与 `Rolling` / `min_period` 用法
- `asm/README.md` — 实验链路的原始记录（clang + c2goasm 命令）
- `math32/README.md` — `math32` 的来源说明（`chewxy/math32` 的 fork）
