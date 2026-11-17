#pragma once
#ifndef NUM_TYPE_INT64_H
#define NUM_TYPE_INT64_H 1

// spec: type_int64.go
//
// int64 的取值口径与解析：
//
//   - 常量 `max_int64` / `min_int64` / `nil2int64` / `int64_nan` / `true2int64` / `false2int64` /
//     `string_bad2int64` / `string_true2int64` / `string_false2int64`；
//   - `parse_int64`：先判空 / true 词 / false 词 / NaN 词，再按 10 进制解析，失败接着按 16 进制解析；
//     两次都失败时按 builtin.h 的 `ignore_parse_exceptions` 返回 0 或抛 DomainError；
//   - `slice_float_to_int64`：浮点序列转 int64，截断取整（Go 的 x64.ToInt64）。
//
// 与 Go 的差异：
//
//   - Go 的签名是 `ParseInt64(s string, v any)`，第二个参数只用来打 panic 日志，C++ 侧去掉；
//   - Go 的 `int64ToString` 是私有辅助，没被任何导出符号使用，不移植；
//   - 浮点转整型的边界：NaN / ±Inf / 越界一律给 MinInt64（实测 Go 在 amd64 上就是这个结果），
//     C++ 侧显式判断后返回，不触发 float → int 的未定义行为。

#include <cstddef>
#include <limits>
#include <string>
#include <string_view>
#include <vector>

#include "builtin.h"
#include "errors.h"
#include "type.h"
#include "type_integer.h"
#include "type_bool.h"
#include "type_strings.h"

namespace num {

    inline constexpr i64 max_int64 = std::numeric_limits<i64>::max();
    inline constexpr i64 min_int64 = std::numeric_limits<i64>::min();
    inline constexpr i64 nil2int64 = 0;          // 空指针转 int64
    inline constexpr i64 int64_nan = 0;          // int64 无效值
    inline constexpr i64 true2int64 = 1;
    inline constexpr i64 false2int64 = 0;
    inline constexpr i64 string_bad2int64 = 0;   // 字符串解析 int64 异常
    inline constexpr i64 string_true2int64 = 1;
    inline constexpr i64 string_false2int64 = 0;

    inline i64 parse_int64(std::string_view s) {
        if (is_empty(s)) {
            return nil2int64;
        }
        if (string_is_true(s)) {
            return string_true2int64;
        }
        if (string_is_false(s)) {
            return string_false2int64;
        }
        if (string_is_nan(s)) {
            return string_false2int64;
        }
        i64 v = 0;
        if (detail::parse_int_strict(s, 10, v)) {
            return v;
        }
        // 10 进制失败后按 16 进制再试一次
        if (detail::parse_int_strict(s, 16, v)) {
            return v;
        }
        if (ignore_parse_exceptions) {
            return string_bad2int64;
        }
        throw DomainError(std::string("strconv.ParseInt: parsing \"") + std::string(s)
                          + "\": invalid syntax");
    }

    // 浮点转 int64：截断取零方向；NaN / ±Inf / 越界给 MinInt64（与 Go 在 amd64 上一致）
    inline i64 float_to_int64(f64 v) {
        constexpr f64 k_min = -9223372036854775808.0;   // -2^63
        constexpr f64 k_max = 9223372036854775808.0;    // 2^63（开区间）
        if (!(v >= k_min && v < k_max)) {
            return std::numeric_limits<i64>::min();
        }
        return static_cast<i64>(v);
    }

    template<typename T, std::enable_if_t<std::is_floating_point_v<T>, int> = 0>
    inline Vector<i64> slice_float_to_int64(const Vector<T> &x) {
        Vector<i64> d = Vector<i64>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = float_to_int64(static_cast<f64>(x[i]));
        }
        return d;
    }

    // Go 的 SliceToInts 对 []string 逐个走 AnyToGeneric[int]（即 ParseInt64）
    inline Vector<i64> slice_to_ints(const std::vector<std::string> &values) {
        Vector<i64> d = Vector<i64>::from_shape({values.size()});
        for (std::size_t i = 0; i < values.size(); ++i) {
            d[i] = parse_int64(values[i]);
        }
        return d;
    }

}  // namespace num

#endif  // NUM_TYPE_INT64_H
