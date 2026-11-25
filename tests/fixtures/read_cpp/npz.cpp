#include "npz.h"

#include <algorithm>
#include <bit>
#include <cctype>
#include <cstdint>
#include <filesystem>
#include <fstream>
#include <sstream>
#include <stdexcept>
#include <string>
#include <utility>
#include <vector>

namespace num::fixture {
namespace {

constexpr std::uint32_t kLocalHeader = 0x04034b50U;
constexpr std::uint32_t kCentralHeader = 0x02014b50U;
constexpr std::uint32_t kEndOfCentral = 0x06054b50U;
constexpr std::uint64_t kFnvOffset = 14695981039346656037ULL;
constexpr std::uint64_t kFnvPrime = 1099511628211ULL;
constexpr char kMetaPrefix[] = "meta_";

std::vector<std::uint8_t> read_file(const std::string& path) {
    std::ifstream in(path, std::ios::binary);
    if (!in) {
        throw std::runtime_error(path + ": 打开失败");
    }
    in.seekg(0, std::ios::end);
    const std::streamoff size = in.tellg();
    if (size < 0) {
        throw std::runtime_error(path + ": 无法取文件长度");
    }
    in.seekg(0, std::ios::beg);
    std::vector<std::uint8_t> data(static_cast<std::size_t>(size));
    if (size > 0 && !in.read(reinterpret_cast<char*>(data.data()), size)) {
        throw std::runtime_error(path + ": 读取失败");
    }
    return data;
}

std::uint32_t read_u16(const std::vector<std::uint8_t>& data, std::size_t offset, const std::string& what) {
    if (offset + 2 > data.size()) {
        throw std::runtime_error(what + ": 越界");
    }
    return static_cast<std::uint32_t>(data[offset]) |
           (static_cast<std::uint32_t>(data[offset + 1]) << 8);
}

std::uint32_t read_u32(const std::vector<std::uint8_t>& data, std::size_t offset, const std::string& what) {
    if (offset + 4 > data.size()) {
        throw std::runtime_error(what + ": 越界");
    }
    return static_cast<std::uint32_t>(data[offset]) |
           (static_cast<std::uint32_t>(data[offset + 1]) << 8) |
           (static_cast<std::uint32_t>(data[offset + 2]) << 16) |
           (static_cast<std::uint32_t>(data[offset + 3]) << 24);
}

std::uint64_t read_u64(const std::vector<std::uint8_t>& data, std::size_t offset, const std::string& what) {
    if (offset + 8 > data.size()) {
        throw std::runtime_error(what + ": 越界");
    }
    std::uint64_t value = 0;
    for (int i = 7; i >= 0; --i) {
        value = (value << 8) | data[offset + static_cast<std::size_t>(i)];
    }
    return value;
}

std::string trim(const std::string& text) {
    std::size_t begin = 0;
    std::size_t end = text.size();
    while (begin < end && std::isspace(static_cast<unsigned char>(text[begin])) != 0) {
        ++begin;
    }
    while (end > begin && std::isspace(static_cast<unsigned char>(text[end - 1])) != 0) {
        --end;
    }
    return text.substr(begin, end - begin);
}

std::string quoted(const std::string& header, const std::string& key) {
    for (const char quote : {'\'', '"'}) {
        const std::string one(1, quote);
        const std::string marker = one + key + one + ": " + one;
        const std::size_t pos = header.find(marker);
        if (pos == std::string::npos) {
            continue;
        }
        const std::size_t begin = pos + marker.size();
        const std::size_t end = header.find(quote, begin);
        if (end != std::string::npos) {
            return header.substr(begin, end - begin);
        }
    }
    throw std::runtime_error(std::string("npy 头缺少 ") + key);
}

std::vector<std::uint64_t> parse_shape(const std::string& header) {
    std::size_t begin = std::string::npos;
    const std::size_t single = header.find("'shape': (");
    if (single != std::string::npos) {
        begin = single + 10;
    } else {
        const std::size_t dual = header.find("\"shape\": (");
        if (dual == std::string::npos) {
            throw std::runtime_error("npy 头缺少 shape");
        }
        begin = dual + 10;
    }
    const std::size_t end = header.find(')', begin);
    if (end == std::string::npos) {
        throw std::runtime_error("npy 头的 shape 未闭合");
    }
    std::vector<std::uint64_t> shape;
    std::stringstream stream(header.substr(begin, end - begin));
    std::string item;
    while (std::getline(stream, item, ',')) {
        item = trim(item);
        if (item.empty()) {
            continue;
        }
        shape.push_back(static_cast<std::uint64_t>(std::stoull(item)));
    }
    return shape;
}

std::size_t descriptor_size(const std::string& descr) {
    if (descr.empty()) {
        throw std::runtime_error("dtype 为空");
    }
    std::size_t begin = 0;
    if (descr[0] == '<' || descr[0] == '|' || descr[0] == '=') {
        begin = 1;
    } else if (descr[0] == '>') {
        throw std::runtime_error("不支持大端 dtype: " + descr);
    }
    const std::string body = descr.substr(begin);
    if (body.empty()) {
        throw std::runtime_error("dtype 非法: " + descr);
    }
    if (body[0] == 'S' || body[0] == 'U') {
        const std::size_t width = static_cast<std::size_t>(std::stoull(body.substr(1)));
        return body[0] == 'U' ? width * 4 : width;
    }
    if (body == "b1") {
        return 1;
    }
    if (body == "f2" || body == "i2" || body == "u2") {
        return 2;
    }
    if (body == "f4" || body == "i4" || body == "u4") {
        return 4;
    }
    if (body == "f8" || body == "i8" || body == "u8" || body == "c8") {
        return 8;
    }
    if (body == "c16") {
        return 16;
    }
    throw std::runtime_error("不支持的 dtype: " + descr);
}

Array parse_npy(const std::string& name, const std::vector<std::uint8_t>& raw) {
    if (raw.size() < 10) {
        throw std::runtime_error(name + ": .npy 头被截断");
    }
    const std::string magic(reinterpret_cast<const char*>(raw.data()), 6);
    if (magic != "\x93NUMPY") {
        throw std::runtime_error(name + ": 不是 .npy（缺少 magic）");
    }
    const std::uint8_t major = raw[6];
    std::size_t header_length = 0;
    std::size_t offset = 0;
    if (major == 1 || major == 2) {
        header_length = read_u16(raw, 8, name + ": npy 头长度");
        offset = 10;
    } else if (major == 3) {
        header_length = read_u32(raw, 8, name + ": npy 头长度");
        offset = 12;
    } else {
        throw std::runtime_error(name + ": 不支持的 .npy 版本 " + std::to_string(major));
    }
    if (offset + header_length > raw.size()) {
        throw std::runtime_error(name + ": npy 头长度越界");
    }
    const std::string header(reinterpret_cast<const char*>(raw.data() + offset), header_length);
    if (header.find("'fortran_order': True") != std::string::npos ||
        header.find("\"fortran_order\": True") != std::string::npos) {
        throw std::runtime_error(name + ": 不支持 fortran_order");
    }

    Array array;
    array.name = name;
    array.descr = quoted(header, "descr");
    array.shape = parse_shape(header);
    const std::size_t begin = offset + header_length;
    array.data.assign(raw.begin() + static_cast<std::ptrdiff_t>(begin), raw.end());
    return array;
}

std::uint64_t checksum(const File& file) {
    std::uint64_t hash = kFnvOffset;
    for (const auto& [name, array] : file.arrays) {
        if (name.rfind(kMetaPrefix, 0) == 0) {
            continue;
        }
        for (const std::uint8_t byte : array.data) {
            hash ^= static_cast<std::uint64_t>(byte);
            hash *= kFnvPrime;
        }
    }
    return hash;
}

}  // namespace

std::size_t Array::itemsize() const {
    try {
        return descriptor_size(descr);
    } catch (const std::exception& error) {
        throw std::runtime_error(name + ": " + error.what());
    }
}

std::size_t Array::element_count() const {
    std::size_t count = 1;
    for (const std::uint64_t dim : shape) {
        count *= static_cast<std::size_t>(dim);
    }
    return count;
}


std::vector<double> Array::floats() const {
    const std::size_t count = element_count();
    std::vector<double> out;
    out.reserve(count);
    if (descr == "<f8") {
        for (std::size_t i = 0; i < count; ++i) {
            out.push_back(std::bit_cast<double>(read_u64(data, i * 8, name)));
        }
        return out;
    }
    if (descr == "<f4") {
        for (std::size_t i = 0; i < count; ++i) {
            out.push_back(static_cast<double>(std::bit_cast<float>(read_u32(data, i * 4, name))));
        }
        return out;
    }
    throw std::runtime_error(name + ": 不是浮点数组（" + descr + "）");
}

std::vector<std::int64_t> Array::ints() const {
    const std::size_t count = element_count();
    std::vector<std::int64_t> out;
    out.reserve(count);
    for (std::size_t i = 0; i < count; ++i) {
        std::int64_t value = 0;
        if (descr == "<i8") {
            value = static_cast<std::int64_t>(read_u64(data, i * 8, name));
        } else if (descr == "<i4") {
            value = static_cast<std::int32_t>(read_u32(data, i * 4, name));
        } else if (descr == "<i2") {
            value = static_cast<std::int16_t>(read_u16(data, i * 2, name));
        } else if (descr == "<i1") {
            value = static_cast<std::int8_t>(data[i]);
        } else if (descr == "<u8") {
            value = static_cast<std::int64_t>(read_u64(data, i * 8, name));
        } else if (descr == "<u4") {
            value = static_cast<std::int64_t>(read_u32(data, i * 4, name));
        } else if (descr == "<u2") {
            value = static_cast<std::int64_t>(read_u16(data, i * 2, name));
        } else if (descr == "<u1") {
            value = static_cast<std::int64_t>(data[i]);
        } else {
            throw std::runtime_error(name + ": 不是整数数组（" + descr + "）");
        }
        out.push_back(value);
    }
    return out;
}

std::vector<bool> Array::bools() const {
    if (descr != "|b1") {
        throw std::runtime_error(name + ": 不是 bool 数组（" + descr + "）");
    }
    std::vector<bool> out;
    out.reserve(data.size());
    for (const std::uint8_t byte : data) {
        out.push_back(byte != 0);
    }
    return out;
}

double Array::scalar_f64() const {
    const auto values = floats();
    if (values.size() != 1) {
        throw std::runtime_error(name + ": 应为标量");
    }
    return values[0];
}

std::int64_t Array::scalar_i64() const {
    const auto values = ints();
    if (values.size() != 1) {
        throw std::runtime_error(name + ": 应为标量");
    }
    return values[0];
}

File read_npz(const std::string& path) {
    const std::vector<std::uint8_t> raw = read_file(path);

    std::int64_t last = static_cast<std::int64_t>(raw.size()) - 22;
    std::int64_t first = last - 65535;
    if (first < 0) {
        first = 0;
    }
    std::int64_t eocd = -1;
    for (std::int64_t pos = last; pos >= first; --pos) {
        const std::size_t offset = static_cast<std::size_t>(pos);
        if (read_u32(raw, offset, path) == kEndOfCentral) {
            eocd = pos;
            break;
        }
    }
    if (eocd < 0) {
        throw std::runtime_error(path + ": 找不到 zip 目录末尾记录");
    }
    const std::size_t eocd_offset = static_cast<std::size_t>(eocd);
    const std::uint32_t total = read_u16(raw, eocd_offset + 10, path);
    const std::uint32_t central_size = read_u32(raw, eocd_offset + 12, path);
    const std::uint32_t central_offset = read_u32(raw, eocd_offset + 16, path);
    if (static_cast<std::size_t>(central_offset) + central_size > raw.size()) {
        throw std::runtime_error(path + ": zip 中央目录越界");
    }

    File file;
    file.path = path;
    std::size_t pos = central_offset;
    for (std::uint32_t index = 0; index < total; ++index) {
        if (read_u32(raw, pos, path) != kCentralHeader) {
            throw std::runtime_error(path + ": 中央目录条目签名不符");
        }
        const std::uint32_t name_length = read_u16(raw, pos + 28, path);
        const std::uint32_t extra_length = read_u16(raw, pos + 30, path);
        const std::uint32_t comment_length = read_u16(raw, pos + 32, path);
        const std::uint32_t local_offset = read_u32(raw, pos + 42, path);
        if (pos + 46 + name_length > raw.size()) {
            throw std::runtime_error(path + ": 中央目录条目名越界");
        }
        const std::string entry_name(reinterpret_cast<const char*>(raw.data() + pos + 46), name_length);
        pos += 46 + name_length + extra_length + comment_length;

        if (entry_name.size() < 4 || entry_name.compare(entry_name.size() - 4, 4, ".npy") != 0) {
            continue;
        }
        if (read_u32(raw, local_offset, path) != kLocalHeader) {
            throw std::runtime_error(path + ": 本地文件头签名不符");
        }
        const std::uint32_t method = read_u16(raw, local_offset + 8, path);
        const std::uint32_t compressed = read_u32(raw, local_offset + 18, path);
        const std::uint32_t uncompressed = read_u32(raw, local_offset + 22, path);
        const std::uint32_t local_name = read_u16(raw, local_offset + 26, path);
        const std::uint32_t local_extra = read_u16(raw, local_offset + 28, path);
        if (method != 0 || compressed != uncompressed) {
            throw std::runtime_error(path + ": " + entry_name + " 不是 ZIP_STORED");
        }
        const std::size_t data_offset =
            static_cast<std::size_t>(local_offset) + 30 + local_name + local_extra;
        if (data_offset + compressed > raw.size()) {
            throw std::runtime_error(path + ": " + entry_name + " 数据越界");
        }
        std::vector<std::uint8_t> payload(raw.begin() + static_cast<std::ptrdiff_t>(data_offset),
                                          raw.begin() + static_cast<std::ptrdiff_t>(data_offset + compressed));
        Array array = parse_npy(entry_name.substr(0, entry_name.size() - 4), payload);
        std::string key = array.name;
        file.arrays.emplace(std::move(key), std::move(array));
    }
    if (file.arrays.empty()) {
        throw std::runtime_error(path + ": 没有 .npy 条目");
    }
    return file;
}

void verify(const File& file) {
    for (const auto& [name, array] : file.arrays) {
        const std::size_t expected = array.element_count() * array.itemsize();
        if (expected != array.data.size()) {
            throw std::runtime_error(file.path + ": " + name + " 数据长度 " +
                                     std::to_string(array.data.size()) + "，按 shape 应为 " +
                                     std::to_string(expected));
        }
    }

    const auto checksum_array = file.arrays.find("meta_checksum");
    if (checksum_array == file.arrays.end()) {
        throw std::runtime_error(file.path + ": 缺少 meta_checksum");
    }
    if (checksum_array->second.descr != "<u8" || !checksum_array->second.shape.empty()) {
        throw std::runtime_error(file.path + ": meta_checksum 必须是标量 u8");
    }
    const std::uint64_t want = read_u64(checksum_array->second.data, 0, file.path + ": meta_checksum");
    const std::uint64_t got = checksum(file);
    if (want != got) {
        throw std::runtime_error(file.path + ": 校验和不符");
    }
}

int walk(const std::string& root) {
    namespace fs = std::filesystem;
    int count = 0;
    for (const auto& group : fs::directory_iterator(root)) {
        if (!group.is_directory()) {
            continue;
        }
        const std::string name = group.path().filename().string();
        if (name == "gen" || name.rfind("read_", 0) == 0) {
            continue;
        }
        for (const auto& entry : fs::directory_iterator(group.path())) {
            if (!entry.is_regular_file() || entry.path().extension() != ".npz") {
                continue;
            }
            const File file = read_npz(entry.path().string());
            verify(file);
            ++count;
        }
    }
    return count;
}

}  // namespace num::fixture
