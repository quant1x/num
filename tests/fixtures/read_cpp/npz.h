// npz 读取器（C3）：只依赖 C++ 标准库，读取 C2 产出的 fixture。
#pragma once

#include <cstddef>
#include <cstdint>
#include <map>
#include <string>
#include <vector>

namespace num::fixture {

// Array 是一个 .npy 条目：dtype、形状与原始数据。
struct Array {
    std::string name;
    std::string descr;
    std::vector<std::uint64_t> shape;
    std::vector<std::uint8_t> data;

    std::size_t itemsize() const;
    std::size_t element_count() const;

    // 数值解码：只支持小端 f8 / f4、i8..i1 / u8..u1 与 b1；其它 dtype 抛 std::runtime_error。
    std::vector<double> floats() const;
    std::vector<std::int64_t> ints() const;
    std::vector<bool> bools() const;
    double scalar_f64() const;
    std::int64_t scalar_i64() const;
};

// File 是一份 npz 的全部条目。
struct File {
    std::string path;
    std::map<std::string, Array> arrays;
};

// read_npz 解析 npz，失败抛 std::runtime_error。
File read_npz(const std::string& path);

// verify 校验 dtype / shape 与数据长度一致，并核对 meta_checksum。
void verify(const File& file);

// walk 读取目录下所有 .npz 并校验，返回文件数。
int walk(const std::string& root);

}  // namespace num::fixture
