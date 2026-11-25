#pragma once
#ifndef NUM_NUMB_H
#define NUM_NUMB_H 1

// spec: docs/DATA-FORMAT.md §3.2（版本 2，提案 0001）
//
// 头部 9 字节 + 8×ndim：magic `NUMB` | version u16 LE | dtype 2 字节 ASCII | ndim u8 | dims u64 LE × ndim。
// 数据区小端、紧密排列；文件长度必须精确等于 `9 + 8×ndim + 元素个数×宽度`。
// 版本 1（dtype 占 1 字节）不兼容，遇到即报错。

#include <cstdint>
#include <cstring>
#include <fstream>
#include <stdexcept>
#include <string>
#include <vector>

namespace num {

    constexpr std::uint16_t kNumbVersion = 2;
    constexpr std::size_t kNumbHeaderSize = 9;

    // InvalidNumb 表示头部、版本、dtype 或长度不自洽。
    class InvalidNumb : public std::runtime_error {
    public:
        explicit InvalidNumb(const std::string &message) : std::runtime_error(message) {}
    };

    // Numb 是 `.numb` 的内存形态。
    struct Numb {
        std::string dtype;
        std::vector<std::uint64_t> shape;
        std::vector<std::uint8_t> payload;
    };

    // 各 dtype 码的元素宽度；未知码返回 0。
    inline std::size_t numb_width(const std::string &dtype) {
        if (dtype == "b1" || dtype == "i1" || dtype == "u1") return 1;
        if (dtype == "i2" || dtype == "u2") return 2;
        if (dtype == "i4" || dtype == "u4" || dtype == "f4") return 4;
        if (dtype == "i8" || dtype == "u8" || dtype == "f8") return 8;
        return 0;
    }

    // 形状描述的元素个数；维度为 0 时按标量计 1。
    inline std::uint64_t numb_count(const std::vector<std::uint64_t> &shape) {
        std::uint64_t count = 1;
        for (const auto length : shape) {
            if (length != 0 && count > UINT64_MAX / length) {
                throw InvalidNumb("numb: 形状过大，元素个数溢出");
            }
            count *= length;
        }
        return count;
    }

    inline void numb_put_u16(std::vector<std::uint8_t> &out, std::uint16_t value) {
        out.push_back(static_cast<std::uint8_t>(value & 0xFF));
        out.push_back(static_cast<std::uint8_t>((value >> 8) & 0xFF));
    }

    inline void numb_put_u64(std::vector<std::uint8_t> &out, std::uint64_t value) {
        for (int shift = 0; shift < 64; shift += 8) {
            out.push_back(static_cast<std::uint8_t>((value >> shift) & 0xFF));
        }
    }

    inline std::uint16_t numb_get_u16(const std::vector<std::uint8_t> &data, std::size_t offset) {
        return static_cast<std::uint16_t>(data[offset]) |
               static_cast<std::uint16_t>(data[offset + 1] << 8);
    }

    inline std::uint64_t numb_get_u64(const std::vector<std::uint8_t> &data, std::size_t offset) {
        std::uint64_t value = 0;
        for (int index = 0; index < 8; ++index) {
            value |= static_cast<std::uint64_t>(data[offset + index]) << (8 * index);
        }
        return value;
    }

    // 校验 dtype、形状与载荷长度。
    inline void numb_validate(const Numb &numb) {
        const std::size_t width = numb_width(numb.dtype);
        if (width == 0) {
            throw InvalidNumb("numb: 未知 dtype 码 \"" + numb.dtype + "\"");
        }
        const std::uint64_t count = numb_count(numb.shape);
        if (static_cast<std::uint64_t>(numb.payload.size()) != count * width) {
            throw InvalidNumb("numb: 载荷字节数与形状不符");
        }
    }

    // 编码成完整文件字节。
    inline std::vector<std::uint8_t> encode_numb(const Numb &numb) {
        numb_validate(numb);
        std::vector<std::uint8_t> out;
        out.reserve(kNumbHeaderSize + 8 * numb.shape.size() + numb.payload.size());
        out.insert(out.end(), {'N', 'U', 'M', 'B'});
        numb_put_u16(out, kNumbVersion);
        out.push_back(static_cast<std::uint8_t>(numb.dtype[0]));
        out.push_back(static_cast<std::uint8_t>(numb.dtype[1]));
        out.push_back(static_cast<std::uint8_t>(numb.shape.size()));
        for (const auto length : numb.shape) {
            numb_put_u64(out, length);
        }
        out.insert(out.end(), numb.payload.begin(), numb.payload.end());
        return out;
    }

    // 解析完整文件字节。
    inline Numb decode_numb(const std::vector<std::uint8_t> &data) {
        if (data.size() < kNumbHeaderSize) {
            throw InvalidNumb("numb: 文件过短");
        }
        if (data[0] != 'N' || data[1] != 'U' || data[2] != 'M' || data[3] != 'B') {
            throw InvalidNumb("numb: magic 不是 NUMB");
        }
        if (numb_get_u16(data, 4) != kNumbVersion) {
            throw InvalidNumb("numb: 版本不受支持");
        }
        const std::string dtype(reinterpret_cast<const char *>(data.data() + 6), 2);
        const std::size_t width = numb_width(dtype);
        if (width == 0) {
            throw InvalidNumb("numb: 未知 dtype 码 \"" + dtype + "\"");
        }
        const std::size_t ndim = data[8];
        const std::size_t header = kNumbHeaderSize + 8 * ndim;
        if (data.size() < header) {
            throw InvalidNumb("numb: ndim 超出文件长度");
        }
        Numb numb;
        numb.dtype = dtype;
        numb.shape.resize(ndim);
        for (std::size_t i = 0; i < ndim; ++i) {
            numb.shape[i] = numb_get_u64(data, kNumbHeaderSize + 8 * i);
        }
        const std::uint64_t count = numb_count(numb.shape);
        if (data.size() != header + static_cast<std::size_t>(count) * width) {
            throw InvalidNumb("numb: 文件长度与形状不符");
        }
        numb.payload.assign(data.begin() + static_cast<std::ptrdiff_t>(header), data.end());
        return numb;
    }

    inline Numb read_numb_file(const std::string &path) {
        std::ifstream input(path, std::ios::binary);
        if (!input) {
            throw InvalidNumb("numb: 无法打开 " + path);
        }
        std::vector<std::uint8_t> data((std::istreambuf_iterator<char>(input)),
                                       std::istreambuf_iterator<char>());
        return decode_numb(data);
    }

    inline void write_numb_file(const std::string &path, const Numb &numb) {
        const auto data = encode_numb(numb);
        std::ofstream output(path, std::ios::binary);
        if (!output) {
            throw InvalidNumb("numb: 无法写入 " + path);
        }
        output.write(reinterpret_cast<const char *>(data.data()), static_cast<std::streamsize>(data.size()));
    }

    // 用 float64 序列构造 Numb。
    inline Numb numb_f64(const std::vector<std::uint64_t> &shape, const std::vector<double> &values) {
        Numb numb;
        numb.dtype = "f8";
        numb.shape = shape;
        numb.payload.resize(8 * values.size());
        for (std::size_t i = 0; i < values.size(); ++i) {
            std::uint64_t bits = 0;
            std::memcpy(&bits, &values[i], 8);
            for (int shift = 0; shift < 64; shift += 8) {
                numb.payload[8 * i + static_cast<std::size_t>(shift / 8)] =
                    static_cast<std::uint8_t>((bits >> shift) & 0xFF);
            }
        }
        return numb;
    }

    // 按 f8 解码。
    inline std::vector<double> numb_floats(const Numb &numb) {
        if (numb.dtype != "f8") {
            throw InvalidNumb("numb: dtype 不是 f8");
        }
        std::vector<double> out(numb.payload.size() / 8);
        for (std::size_t i = 0; i < out.size(); ++i) {
            std::uint64_t bits = 0;
            for (int index = 0; index < 8; ++index) {
                bits |= static_cast<std::uint64_t>(numb.payload[8 * i + static_cast<std::size_t>(index)])
                        << (8 * index);
            }
            std::memcpy(&out[i], &bits, 8);
        }
        return out;
    }

    // 按 f4 解码。
    inline std::vector<float> numb_floats32(const Numb &numb) {
        if (numb.dtype != "f4") {
            throw InvalidNumb("numb: dtype 不是 f4");
        }
        std::vector<float> out(numb.payload.size() / 4);
        for (std::size_t i = 0; i < out.size(); ++i) {
            std::uint32_t bits = 0;
            for (int index = 0; index < 4; ++index) {
                bits |= static_cast<std::uint32_t>(numb.payload[4 * i + static_cast<std::size_t>(index)])
                        << (8 * index);
            }
            std::memcpy(&out[i], &bits, 4);
        }
        return out;
    }

}  // namespace num

#endif  // NUM_NUMB_H
