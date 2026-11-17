#pragma once
#ifndef NUM_TYPE_FLOAT64_H
#define NUM_TYPE_FLOAT64_H 1

// spec: type_float64.go
//
// float64 的无效值、解析与序列转换：
//
//   - `float64_nan()`：IEEE 754 的 NaN；`float64_is_nan(f)`：NaN 或 ±Inf 都算（与 type.h 的
//     `is_nan` 同一口径，这里保留 Go 的独立入口）；
//   - `parse_float64`：空串 → NaN；true 词 → 1、false 词 → 0；NaN / Inf 词按 strconv 的口径识别；
//     其余按十进制解析，失败按 `ignore_parse_exceptions` 返回 NaN 或抛 DomainError；
//   - `slice_to_float64`：数值 / bool / 字符串序列转 float64（Go 的 SliceToFloat64）。
//
// 与 Go 的差异：
//
//   - Go 的签名带 `v any`（只用于 panic 日志），C++ 侧去掉；
//   - 十六进制浮点（"0x1p-2"）Go 的 ParseFloat 认，这里不认（见 type_strings.h 的已知差异）；
//   - Go 的 `SliceToFloat64` 靠 any 分派，C++ 侧是重载；[]float64 在 Go 里原样返回同一底层数组，
//     C++ 侧返回副本。

#include <cstddef>
#include <limits>
#include <string>
#include <string_view>
#include <type_traits>
#include <vector>

#include "builtin.h"
#include "errors.h"
#include "type.h"
#include "type_bool.h"
#include "type_strings.h"

namespace num {

    inline constexpr f64 max_float64 = std::numeric_limits<f64>::max();
    inline constexpr f64 min_float64 = std::numeric_limits<f64>::denorm_min();
    inline constexpr f64 true2float64 = 1.0;
    inline constexpr f64 false2float64 = 0.0;
    inline constexpr f64 string_nil2float = 0.0;    // deprecated：字符串空指针转 float64
    inline constexpr f64 string_bad2float = 0.0;    // deprecated：字符串解析 float64 异常
    inline constexpr f64 string_true2float64 = 1.0;
    inline constexpr f64 string_false2float64 = 0.0;

    // 空串 / 解析失败时的返回值（Go 的 __nilToFloat64）
    inline f64 nil_to_float64() {
        return std::numeric_limits<f64>::quiet_NaN();
    }

    inline f64 float64_nan() {
        return std::numeric_limits<f64>::quiet_NaN();
    }

    // NaN 与 ±Inf 都算无效值
    inline bool float64_is_nan(f64 f) {
        return std::isnan(f) || std::isinf(f);
    }

    inline f64 parse_float64(std::string_view s) {
        if (is_empty(s)) {
            return nil_to_float64();
        }
        if (string_is_true(s)) {
            return string_true2float64;
        }
        if (string_is_false(s)) {
            return string_false2float64;
        }
        f64 v = 0;
        if (detail::parse_float_strict(s, v)) {
            return v;
        }
        if (ignore_parse_exceptions) {
            return nil_to_float64();
        }
        throw DomainError(std::string("strconv.ParseFloat: parsing \"") + std::string(s)
                          + "\": invalid syntax");
    }

    template<typename T, std::enable_if_t<std::is_arithmetic_v<T> && !std::is_same_v<T, bool>, int> = 0>
    inline Vector<f64> slice_to_float64(const Vector<T> &v) {
        Vector<f64> d = Vector<f64>::from_shape({v.size()});
        for (std::size_t i = 0; i < v.size(); ++i) {
            d[i] = static_cast<f64>(v[i]);
        }
        return d;
    }

    // Go 的 SliceToFloat64 对 []bool 走 x64.FromBool：true → 1、false → 0
    inline Vector<f64> slice_to_float64(const Vector<bool> &v) {
        Vector<f64> d = Vector<f64>::from_shape({v.size()});
        for (std::size_t i = 0; i < v.size(); ++i) {
            d[i] = v[i] ? 1.0 : 0.0;
        }
        return d;
    }

    // Go 的 SliceToFloat64 对 []string 逐个走 AnyToFloat64（即 ParseFloat64）
    inline Vector<f64> slice_to_float64(const std::vector<std::string> &values) {
        Vector<f64> d = Vector<f64>::from_shape({values.size()});
        for (std::size_t i = 0; i < values.size(); ++i) {
            d[i] = parse_float64(values[i]);
        }
        return d;
    }

}  // namespace num

#endif  // NUM_TYPE_FLOAT64_H
