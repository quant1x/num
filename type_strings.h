#pragma once
#ifndef NUM_TYPE_STRINGS_H
#define NUM_TYPE_STRINGS_H 1

// spec: type_strings.go
//
// 字符串判据与数值转字符串：
//
//   - `is_empty`：去掉首尾空白后为空即算空（Go 的 strings.TrimSpace）；
//   - `string_is_nan`：去空白、转小写后等于 "nan"；
//   - `to_string`：整型走十进制，浮点走 Go 的 `FormatFloat(v, 'G', -1, bits)` 口径
//     （最短往返表示；十进制指数 < -4 或 >= 6 时用大写 E 的科学计数法，否则定点），
//     bool 走 "true" / "false"，字符串原样返回；
//   - `slice_to_string`：逐元素转字符串。
//
// 与 Go 的差异：
//
//   - Go 的 `AnyToString` / `PrintString` 先过 `ExtractValueFromPointer`（reflect 解指针），
//     C++ 侧没有这一层，`to_string` 直接吃值；
//   - `strings.TrimSpace` 认 Unicode 空白，这里只认 ASCII 空白（空格 / \t / \n / \r / \v / \f）；
//   - Go 的 `SliceToString` 靠 any 分派 14 种切片类型，C++ 侧是模板，类型在编译期确定；
//     字符串序列走 `std::vector<std::string>` 的重载。

#include <array>
#include <charconv>
#include <cmath>
#include <cstddef>
#include <cstdint>
#include <limits>
#include <string>
#include <string_view>
#include <type_traits>
#include <vector>

#include "type.h"

namespace num {

    // 字符串形式的 NaN
    inline constexpr std::string_view string_nan = "NaN";
    // 空指针转字符串的结果
    inline constexpr std::string_view nil2string = "NaN";
    // bool 转字符串
    inline constexpr std::string_view true2string = "true";
    inline constexpr std::string_view false2string = "false";

    // 有可能出现的 NaN 字符串的全部选项
    inline constexpr std::array<std::string_view, 4> possible_na_of_string = {"NA", "NaN", "nan", "<nil>"};

    namespace detail {

        inline bool is_space(char c) {
            return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f';
        }

        inline std::string_view trim_space(std::string_view s) {
            while (!s.empty() && is_space(s.front())) {
                s.remove_prefix(1);
            }
            while (!s.empty() && is_space(s.back())) {
                s.remove_suffix(1);
            }
            return s;
        }

        // Go 的 strconv.FormatFloat(v, 'G', -1, bits)：
        // 有效数字取最短往返表示，十进制指数（首位的指数）< -4 或 >= 6 时走大写 E 的科学计数法。
        // eprec 固定为 6 是 Go 在 prec = -1（最短表示）时的取值。
        template<typename Float>
        inline std::string format_general(Float v) {
            if (std::isnan(v)) {
                return std::string(string_nan);
            }
            if (std::isinf(v)) {
                return v < 0 ? "-Inf" : "+Inf";
            }
            // 零要先判：正负零的指数写法由 to_chars 决定，Go 给的是 "0" / "-0"
            if (v == 0) {
                return std::signbit(v) ? "-0" : "0";
            }

            char buf[64];
            const auto res = std::to_chars(buf, buf + sizeof(buf), v, std::chars_format::scientific);
            const std::string sci(buf, static_cast<std::size_t>(res.ptr - buf));

            bool negative = false;
            std::size_t pos = 0;
            if (!sci.empty() && sci[0] == '-') {
                negative = true;
                pos = 1;
            }
            const std::size_t e_pos = sci.find('e');
            std::string digits;
            for (std::size_t i = pos; i < e_pos; ++i) {
                if (sci[i] != '.') {
                    digits.push_back(sci[i]);
                }
            }
            const int exp = std::stoi(sci.substr(e_pos + 1));  // 首位数字的十进制指数
            const int nd = static_cast<int>(digits.size());

            std::string out;
            if (negative) {
                out.push_back('-');
            }
            if (exp < -4 || exp >= 6) {
                out.push_back(digits[0]);
                if (nd > 1) {
                    out.push_back('.');
                    out.append(digits, 1, std::string::npos);
                }
                out.push_back('E');
                out.push_back(exp < 0 ? '-' : '+');
                const int abs_exp = exp < 0 ? -exp : exp;
                if (abs_exp < 10) {
                    out.push_back('0');
                }
                out += std::to_string(abs_exp);
                return out;
            }

            const int dp = exp + 1;  // 小数点位置（相对 digits）
            if (dp <= 0) {
                out += "0.";
                out.append(static_cast<std::size_t>(-dp), '0');
                out += digits;
            } else if (dp >= nd) {
                out += digits;
                out.append(static_cast<std::size_t>(dp - nd), '0');
            } else {
                out.append(digits, 0, static_cast<std::size_t>(dp));
                out.push_back('.');
                out.append(digits, static_cast<std::size_t>(dp), std::string::npos);
            }
            return out;
        }

        // 字符串 → 整型的内核，口径对齐 Go 的 strconv.ParseInt(s, base, bits)：
        // 整串必须消费完、base 16 时不认 "0x" 前缀（Go 只在 base = 0 时认前缀）、认 +/- 符号、
        // 越界一律算失败（Go 返回 ErrRange，调用方按解析失败处理）。
        template<typename Int>
        inline bool parse_int_strict(std::string_view s, int base, Int &out) {
            bool negative = false;
            if (!s.empty() && (s.front() == '+' || s.front() == '-')) {
                negative = s.front() == '-';
                s.remove_prefix(1);
            }
            if (s.empty()) {
                return false;
            }
            if (base == 16 && s.size() >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')) {
                return false;
            }
            i64 v = 0;
            const char *first = s.data();
            const char *last = s.data() + s.size();
            const auto res = std::from_chars(first, last, v, base);
            if (res.ec != std::errc{} || res.ptr != last) {
                return false;
            }
            if (negative) {
                v = -v;
            }
            if (v < static_cast<i64>(std::numeric_limits<Int>::min())
                || v > static_cast<i64>(std::numeric_limits<Int>::max())) {
                return false;
            }
            out = static_cast<Int>(v);
            return true;
        }

        // 字符串 → 浮点的内核，口径对齐 Go 的 strconv.ParseFloat：
        //   - 先剥符号，再按词识别 "nan" / "inf" / "infinity"（大小写不敏感，from_chars 不认这些词）；
        //   - 十六进制浮点（"0x1p-2"）剥掉 0x 前缀后按 hex 格式解析，与 Go 一致；
        //   - 其余按十进制解析，整串必须消费完；越界（ErrRange）算失败，调用方按解析失败处理。
        template<typename Float>
        inline bool parse_float_strict(std::string_view s, Float &out) {
            bool negative = false;
            if (!s.empty() && (s.front() == '+' || s.front() == '-')) {
                negative = s.front() == '-';
                s.remove_prefix(1);
            }
            const auto lower = [](std::string_view t) {
                std::string r;
                r.reserve(t.size());
                for (char c : t) {
                    r.push_back((c >= 'A' && c <= 'Z') ? static_cast<char>(c - 'A' + 'a') : c);
                }
                return r;
            };
            const std::string word = lower(s);
            if (word == "nan") {
                out = std::numeric_limits<Float>::quiet_NaN();
                return true;
            }
            if (word == "inf" || word == "infinity") {
                out = negative ? -std::numeric_limits<Float>::infinity()
                               : std::numeric_limits<Float>::infinity();
                return true;
            }
            // 十六进制浮点：Go 的 ParseFloat 认 "0x1p-2"，from_chars 不认前缀，先剥掉再按 hex 解析
            std::chars_format fmt = std::chars_format::general;
            if (s.size() >= 3 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')) {
                s.remove_prefix(2);
                fmt = std::chars_format::hex;
            }
            Float v{};
            const char *first = s.data();
            const char *last = s.data() + s.size();
            const auto res = std::from_chars(first, last, v, fmt);
            if (res.ec != std::errc{} || res.ptr != last) {
                return false;
            }
            out = negative ? -v : v;
            return true;
        }

    }  // namespace detail

    // 去掉首尾空白后是否为空
    inline bool is_empty(std::string_view s) {
        return detail::trim_space(s).empty();
    }

    // 去空白、转小写后是否等于 "nan"
    inline bool string_is_nan(std::string_view s) {
        s = detail::trim_space(s);
        if (s.size() != 3) {
            return false;
        }
        const auto lower = [](char c) {
            return (c >= 'A' && c <= 'Z') ? static_cast<char>(c - 'A' + 'a') : c;
        };
        return lower(s[0]) == 'n' && lower(s[1]) == 'a' && lower(s[2]) == 'n';
    }

    inline std::string to_string(std::string_view v) {
        return std::string(v);
    }

    // 必须显式给 const char*：否则字符串字面量会走"指针 → bool"的标准转换，落进 bool 重载
    inline std::string to_string(const char *v) {
        return std::string(v);
    }

    inline std::string to_string(bool v) {
        return std::string(v ? true2string : false2string);
    }

    inline std::string to_string(f32 v) {
        return detail::format_general(v);
    }

    inline std::string to_string(f64 v) {
        return detail::format_general(v);
    }

    // 整型（含 i8 / u8 这类窄类型）走十进制
    template<typename T, std::enable_if_t<std::is_integral_v<T> && !std::is_same_v<T, bool>, int> = 0>
    inline std::string to_string(T v) {
        return std::to_string(static_cast<i64>(v));
    }

    template<typename T>
    inline std::vector<std::string> slice_to_string(const Vector<T> &v) {
        std::vector<std::string> d;
        d.reserve(v.size());
        for (std::size_t i = 0; i < v.size(); ++i) {
            d.push_back(to_string(v[i]));
        }
        return d;
    }

    // Go 的 SliceToString 对 []string 原样返回
    inline std::vector<std::string> slice_to_string(const std::vector<std::string> &v) {
        return v;
    }

}  // namespace num

#endif  // NUM_TYPE_STRINGS_H
