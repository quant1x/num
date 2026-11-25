# tests/fixtures/read_cpp —— C++ 读取器（C3）

只依赖 C++20 标准库，读取 C2 产出的 npz。构建目标是根 `CMakeLists.txt` 里的 `num_fixture_check`，
由 ctest 以 `fixtures_npz_read_cpp` 注册。

## 文件

| 文件 | 内容 |
| --- | --- |
| `npz.h` / `npz.cpp` | `read_npz` / `verify` / `walk`：zip 中央目录、本地头、`.npy` 头解析、dtype 与长度校验、FNV-1a |
| `main.cpp` | 命令行入口，接受 fixture 根目录，退出码非 0 表示校验失败 |

## 用法

```bash
cmake --build --preset release-static --target num_fixture_check
./build/release-static/Release/num_fixture_check.exe tests/fixtures
ctest --preset release-static -R fixtures_npz_read_cpp
```

## 边界

只读不写；不支持 deflate、`fortran_order`、大端 dtype。解析失败抛 `std::runtime_error`，
由 `main` 转成非 0 退出码，调用方按异常类型处理，不解析文案。
