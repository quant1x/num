# tests/install —— 安装包冒烟（D3）

验证 C++ 产物「装出来就能被别人用」这条链路：

- `smoke.py`：`cmake --install` 到 `temp/install`，再配置 `consumer/`（只用 `find_package(num)`），
  构建并运行，断言输出 `10 2.5`；
- `consumer/`：最小消费方工程，示范两个容易踩的前提（见 [RELEASE.md](../../docs/RELEASE.md) §3.1）：
  1. `CMAKE_PREFIX_PATH` 要同时包含 num 的安装前缀与 **xtensor** 的包目录；
  2. MSVC 要显式 `/utf-8` —— 头文件含中文注释，中文 Windows 默认按 GBK 读会报语法错误。

```bash
python tests/install/smoke.py
```

前置：先按 [BUILDING.md](../../BUILDING.md) §6.2 构建 `build/release-static`。
