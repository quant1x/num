#pragma once
#ifndef NUM_TYPE_INT32_H
#define NUM_TYPE_INT32_H 1

// spec: type_int32.go
//
// 与 type_int64.h 同一套口径，宽度换成 32 位：
//
//   - 常量 `max_int32` / `min_int32` / `nil2int32` / `int32_nan` / `true2int32` / `false2int32` /
//     `string_bad2int32` / `string_true2int32` / `string_false2int32`；
//   - `parse_int32`：先判空 / 词表，再按 10 进制、失败转 16 进制；都失败按 `ignore_parse_exceptions`
//     返回 0 或抛 DomainError；
//   - `slice_float_to_int32`：浮点序列转 int32，截断取整（Go 的 x64.ToInt32 / x32.ToInt32）。
//
// 与 Go 的差异：
//
//   - Go 的签名带 `v any`（只用于 panic 日志），C++ 侧去掉；私有辅助 `int32ToString` 未移植；
//   - 浮点转整型的边界：NaN / ±Inf / 越界一律给 MinInt32（实测 Go 在 amd64 上就是这个结果）。

#include <cstddef>
#include <limits>
#include <string>
#include <string_view>

#include "builtin.h"
#include "errors.h"
#include "type.h"
#include "type_bool.h"
#include "type_strings.h"

namespace num {

    inline constexpr i32 max_int32 = std::numeric_limits<i32>::max();
    inline constexpr i32 min_int32 = std::numeric_limits<i32>::min();
    inline constexpr i32 nil2int32 = 0;
    inline constexpr i32 int32_nan = 0;
    inline constexpr i32 true2int32 = 1;
    inline constexpr i32 false2int32 = 0;
    inline constexpr i32 string_bad2int32 = 0;
    inline constexpr i32 string_true2int32 = 1;
    inline constexpr i32 string_false2int32 = 0;

    inline i32 parse_int32(std::string_view s) {
        if (is_empty(s)) {
            return nil2int32;
        }
        if (string_is_true(s)) {
            return string_true2int32;
        }
        if (string_is_false(s)) {
            return string_false2int32;
        }
        if (string_is_nan(s)) {
            return string_false2int32;
        }
        i32 v = 0;
        if (detail::parse_int_strict(s, 10, v)) {
            return v;
        }
        if (detail::parse_int_strict(s, 16, v)) {
            return v;
        }
        if (ignore_parse_exceptions) {
            return string_bad2int32;
        }
        throw DomainError(std::string("strconv.ParseInt: parsing \"") + std::string(s)
                          + "\": invalid syntax");
    }

    // 浮点转 int32：截断取零方向；NaN / ±Inf / 越界给 MinInt32（与 Go 在 amd64 上一致）
    inline i32 float_to_int32(f64 v) {
        constexpr f64 k_min = -2147483649.0;  // 截断后落在 -2^31 上的最小值
        constexpr f64 k_max = 2147483648.0;   // 2^31（开区间）
        if (!(v > k_min && v < k_max)) {
            return std::numeric_limits<i32>::min();
        }
        return static_cast<i32>(v);
    }

    template<typename T, std::enable_if_t<std::is_floating_point_v<T>, int> = 0>
    inline Vector<i32> slice_float_to_int32(const Vector<T> &x) {
        Vector<i32> d = Vector<i32>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = float_to_int32(static_cast<f64>(x[i]));
        }
        return d;
    }

}  // namespace num

#endif  // NUM_TYPE_INT32_H
