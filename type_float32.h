#pragma once
#ifndef NUM_TYPE_FLOAT32_H
#define NUM_TYPE_FLOAT32_H 1

// spec: type_float32.go
//
// 与 type_float64.h 同一套口径，宽度换成 32 位：
//
//   - 常量 `max_float32` / `min_float32` / `string_true2float32` / `string_false2float32`；
//   - `float32_nan()` / `float32_is_nan(f)`（NaN 与 ±Inf 都算）；
//   - `parse_float32`：空串 → NaN；true 词 → 1、false 词 → 0；解析按 32 位精度，**越界算解析失败**
//     （Go 的 ParseFloat(s, 32) 在越界时返回 ErrRange，代码随即落到 NaN 分支）；
//   - `slice_to_float32`：数值 / bool / 字符串序列转 float32（Go 的 SliceToFloat32）。
//
// 与 Go 的差异：签名去掉 `v any`（只用于 panic 日志）；Go 的 []float32 分支原样返回同一底层数组，
// C++ 侧返回副本。

#include <cstddef>
#include <limits>
#include <string>
#include <string_view>
#include <type_traits>
#include <vector>

#include "builtin.h"
#include "errors.h"
#include "type.h"
#include "type_float64.h"
#include "type_bool.h"
#include "type_strings.h"

namespace num {

    inline constexpr f32 max_float32 = std::numeric_limits<f32>::max();
    inline constexpr f32 min_float32 = std::numeric_limits<f32>::denorm_min();
    inline constexpr f32 string_true2float32 = 1.0f;
    inline constexpr f32 string_false2float32 = 0.0f;

    // 空串 / 解析失败时的返回值（Go 的 __nilToFloat32）
    inline f32 nil_to_float32() {
        return std::numeric_limits<f32>::quiet_NaN();
    }

    inline f32 float32_nan() {
        return std::numeric_limits<f32>::quiet_NaN();
    }

    inline bool float32_is_nan(f32 f) {
        return std::isnan(f) || std::isinf(f);
    }

    inline f32 parse_float32(std::string_view s) {
        if (is_empty(s)) {
            return nil_to_float32();
        }
        if (string_is_true(s)) {
            return string_true2float32;
        }
        if (string_is_false(s)) {
            return string_false2float32;
        }
        f32 v = 0;
        if (detail::parse_float_strict(s, v)) {
            return v;
        }
        if (ignore_parse_exceptions) {
            return nil_to_float32();
        }
        throw DomainError(std::string("strconv.ParseFloat: parsing \"") + std::string(s)
                          + "\": invalid syntax");
    }

    template<typename T, std::enable_if_t<std::is_arithmetic_v<T> && !std::is_same_v<T, bool>, int> = 0>
    inline Vector<f32> slice_to_float32(const Vector<T> &v) {
        Vector<f32> d = Vector<f32>::from_shape({v.size()});
        for (std::size_t i = 0; i < v.size(); ++i) {
            d[i] = static_cast<f32>(v[i]);
        }
        return d;
    }

    // Go 的 SliceToFloat32 对 []bool 走 x32.FromBool：true → 1、false → 0
    inline Vector<f32> slice_to_float32(const Vector<bool> &v) {
        Vector<f32> d = Vector<f32>::from_shape({v.size()});
        for (std::size_t i = 0; i < v.size(); ++i) {
            d[i] = v[i] ? 1.0f : 0.0f;
        }
        return d;
    }

    // Go 的 SliceToFloat32 对 []string 逐个走 float32(AnyToFloat64(iv))
    inline Vector<f32> slice_to_float32(const std::vector<std::string> &values) {
        Vector<f32> d = Vector<f32>::from_shape({values.size()});
        for (std::size_t i = 0; i < values.size(); ++i) {
            d[i] = static_cast<f32>(parse_float64(values[i]));
        }
        return d;
    }

}  // namespace num

#endif  // NUM_TYPE_FLOAT32_H
