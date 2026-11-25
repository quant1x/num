# cmake —— 安装导出的包配置模板（D3）

`numConfig.cmake.in` 由 `CMakeLists.txt` 的 `configure_package_config_file` 展开，随 `cmake --install` 落到
`lib/cmake/num/`，供下游 `find_package(num CONFIG)` 使用。

- 导出目标：`num::num_core`（静态库）。
- 依赖：安装时按 `find_dependency(xtensor)` 解析，下游需要把 xtensor 放进 `CMAKE_PREFIX_PATH`
  （本仓的开发构建由 conan 生成器负责）。
- 版本：`numConfigVersion.cmake` 的版本从**精确匹配的 git tag** 派生（`vX.Y.Z` → `X.Y.Z`）；
  未打 tag 的开发树为 `0.0.0`。不要在源码里写版本常量。

用法见 [docs/RELEASE.md](../docs/RELEASE.md)。
