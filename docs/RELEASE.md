# RELEASE —— 编译产物与版本（D3）

发布版本唯一由 git tag `vMAJOR.MINOR.PATCH` 决定，tag 与 `CHANGELOG.md` 由 autochangelog 管理
（[CONTRIBUTING.md](../CONTRIBUTING.md) §10）。**代码与构建脚本不维护版本常量**：构建时从 tag 派生。

## 1. 产物矩阵

| 语言 | 产物 | 入口 | 版本来源 |
| --- | --- | --- | --- |
| Go | `numctl` 单二进制 | `cmd/numctl` | 编译期由 Go 工具链写入 build info |
| C++ | `num_core` 静态库 + `include/num/*.h` + CMake 包 | `CMakeLists.txt` | 精确匹配的 git tag；未打 tag 为 `0.0.0` |
| Rust | `quant1x-num` 库 | `Cargo.toml` | `Cargo.toml` 的 `version`，发布时由 autochangelog 改写 |

## 2. Go

```bash
CGO_ENABLED=0 go build -trimpath -o numctl ./cmd/numctl
./numctl version
```

- `numctl version` 读 `runtime/debug.ReadBuildInfo()`，在打 tag 的提交上构建时得到 tag 版本；
  开发树是伪版本（`vX.Y.Z-0.<时间>-<commit>` 或 `(devel)`），属于预期。
- 不提供硬编码的 `--version`；源码里没有版本字符串。

## 3. C++

```bash
conan install . -s build_type=Release -s compiler.runtime=static
cmake --preset release-static
cmake --build --preset release-static
cmake --install build/release-static --prefix dist/cpp
```

安装结果：

```text
dist/cpp/
├── include/num/*.h
├── lib/num_core.lib            （Linux / macOS 为 libnum_core.a）
└── lib/cmake/num/
    ├── numConfig.cmake
    ├── numConfigVersion.cmake
    └── numTargets.cmake
```

- 版本派生：`git describe --tags --exact-match` 命中 `vX.Y.Z` 时写入 `X.Y.Z`，否则 `0.0.0`。
  未打 tag 的开发树装出 `0.0.0` 是预期行为。
- 下游用法：把 xtensor 放进 `CMAKE_PREFIX_PATH`，再 `find_package(num CONFIG REQUIRED)`，链接 `num::num_core`，
  包含形态为 `<num/abs.h>`。

## 4. Rust

```bash
cargo build --release
```

- 产物版本即 `Cargo.toml` 的 `version`；发布时由 autochangelog 改写，本仓不手工维护第二个版本来源。
- 工具 crate（`tests/fixtures/read_rs`、`tests/io/rs`）的 `0.0.0` 是内部占位，不随发布。

## 5. 发版校验

1. tag 已由 autochangelog 建好，工作树干净；
2. `go test ./internal/functions/ ./x32/ ./x64/ ./vectors/ ./kind/ ./binary/ .`、`ctest --preset release-static`、
   `cargo test` 全绿；
3. `CGO_ENABLED=0 go build` 后 `numctl version` 显示 tag 版本；
4. `cmake --install` 后 `numConfigVersion.cmake` 的 `PACKAGE_VERSION` 等于 tag 版本；
5. 归档 `dist/` 下的产物并记录校验和（上传、签名不在本仓脚本范围）。

## 6. 不做什么

- 不手工改 `CHANGELOG.md`、不建 tag、不调 autochangelog；
- 不在源码或构建脚本里写死版本号；
- 不做制品上传、签名与容器镜像。

相关：[BUILDING.md](../BUILDING.md) §8、[docs/PLAN.md](PLAN.md) 的 D3。
