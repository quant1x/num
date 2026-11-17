#pragma once
#ifndef NUM_TYPE_BOOL_H
#define NUM_TYPE_BOOL_H 1

// spec: type_bool.go
//
// bool 的判据与转换：
//
//   - `string_is_true` / `string_is_false`：字符串的 true / false 词表（大小写混写与中文词），
//     两边都不命中时**两个函数都返回 false**，不要把它当互斥判断；
//   - `bool_to_int` / `bool_to_int32` / `bool_to_int64` / `bool_to_float32` / `bool_to_float64`：
//     true 对应 1，false 对应 0；
//   - `bool_to_string`："true" / "false"（常量与 type_strings.h 同源）；
//   - `parse_bool`：Go 的 strconv.ParseBool（认 "1"/"t"/"T"/"true" 等与 "0"/"f"/"F"/"false" 等，
//     **不去空白**）；解析失败按 builtin.h 的 `ignore_parse_exceptions` 处理 —— true 时返回 false，
//     false 时抛 DomainError（Go 侧是故意 panic）；
//   - `to_bool`：Go 的 AnyToBool 的可移植部分（bool / 整型 / 浮点 / 字符串），字符串先看空、
//     再看 true / false 词表，最后落到 parse_bool；
//   - `slice_to_bool`：逐元素转 bool。
//
// 与 Go 的差异：
//
//   - Go 的 AnyToBool 还认指针（`*bool` / `*int64` / `*string` …）与 default 分支的
//     `fmt.Sprintf("%v")`，C++ 侧没有指针解引用与运行期格式化这一层，只保留值重载；
//   - Go 的 `SliceToBool` 对 []float64 走 x64.ToBool 加速、对 []bool 原样返回；C++ 侧统一走
//     模板循环（结果相同，加速内核属另一条交付线）；
//   - `bool_to_int32` / `bool_to_int64` 在 Go 里返回 type_int32.go / type_int64.go 的
//     True2Int32 / True2Int64 常量，那两个文件尚未移植，这里直接给同值 1 / 0。

#include <cstddef>
#include <string>
#include <string_view>
#include <type_traits>
#include <vector>

#include "builtin.h"
#include "errors.h"
#include "type.h"
#include "type_strings.h"

namespace num {

    inline constexpr bool nil2bool = false;    // 空指针转 bool
    inline constexpr bool bool_nan = false;    // bool 无效值
    inline constexpr bool true2bool = true;
    inline constexpr bool false2bool = false;
    inline constexpr f32 true2float32 = 1.0f;
    inline constexpr f32 false2float32 = 0.0f;

    inline constexpr bool string_bad2bool = false;  // 字符串解析 bool 异常
    inline constexpr bool string_true2bool = true;
    inline constexpr bool string_false2bool = false;

    inline bool string_is_true(std::string_view s) {
        if (s == "true" || s == "TRUE" || s == "True") {
            return true;
        }
        if (s == "y" || s == "Y" || s == "yes" || s == "Yes" || s == "YES"
            || s == "on" || s == "On" || s == "ON") {
            return true;
        }
        if (s == "1") {
            return true;
        }
        if (s == "真" || s == "对" || s == "好") {
            return true;
        }
        return false;
    }

    inline bool string_is_false(std::string_view s) {
        if (s == "false" || s == "FALSE" || s == "False") {
            return true;
        }
        if (s == "n" || s == "N" || s == "no" || s == "No" || s == "NO"
            || s == "off" || s == "Off" || s == "OFF") {
            return true;
        }
        if (s == "0") {
            return true;
        }
        if (s == "假" || s == "错" || s == "坏") {
            return true;
        }
        return false;
    }

    inline i8 bool_to_int(bool b) {
        return static_cast<i8>(b ? 1 : 0);
    }

    inline i32 bool_to_int32(bool b) {
        return b ? 1 : 0;
    }

    inline i64 bool_to_int64(bool b) {
        return b ? 1 : 0;
    }

    inline f32 bool_to_float32(bool b) {
        return b ? true2float32 : false2float32;
    }

    inline f64 bool_to_float64(bool b) {
        return b ? 1.0 : 0.0;
    }

    inline std::string bool_to_string(bool b) {
        return std::string(b ? true2string : false2string);
    }

    // Go 的 strconv.ParseBool：不去空白，只认固定词表
    inline bool parse_bool(std::string_view s) {
        if (s == "1" || s == "t" || s == "T" || s == "TRUE" || s == "true" || s == "True") {
            return true;
        }
        if (s == "0" || s == "f" || s == "F" || s == "FALSE" || s == "false" || s == "False") {
            return false;
        }
        if (ignore_parse_exceptions) {
            return nil2bool;
        }
        throw DomainError(std::string("strconv.ParseBool: parsing \"") + std::string(s)
                          + "\": invalid syntax");
    }

    inline bool to_bool(bool v) {
        return v ? true2bool : false2bool;
    }

    template<typename T, std::enable_if_t<std::is_integral_v<T> && !std::is_same_v<T, bool>, int> = 0>
    inline bool to_bool(T v) {
        return v != 0;
    }

    template<typename T, std::enable_if_t<std::is_floating_point_v<T>, int> = 0>
    inline bool to_bool(T v) {
        return v != 0;
    }

    inline bool to_bool(std::string_view v) {
        if (is_empty(v)) {
            return nil2bool;
        }
        if (string_is_true(v)) {
            return string_true2bool;
        }
        if (string_is_false(v)) {
            return string_false2bool;
        }
        return parse_bool(v);
    }

    // 必须显式给 const char*：否则字符串字面量会走"指针 → bool"的标准转换，落进 bool 重载
    inline bool to_bool(const char *v) {
        return to_bool(std::string_view(v));
    }

    template<typename T>
    inline Vector<bool> slice_to_bool(const Vector<T> &v) {
        Vector<bool> d = Vector<bool>::from_shape({v.size()});
        for (std::size_t i = 0; i < v.size(); ++i) {
            d[i] = to_bool(v[i]);
        }
        return d;
    }

    // Go 的 SliceToBool 对 []bool 原样返回
    inline Vector<bool> slice_to_bool(const Vector<bool> &v) {
        return v;
    }

    // Go 的 SliceToBool 对 []string 逐个走 AnyToBool
    inline Vector<bool> slice_to_bool(const std::vector<std::string> &values) {
        Vector<bool> d = Vector<bool>::from_shape({values.size()});
        for (std::size_t i = 0; i < values.size(); ++i) {
            d[i] = to_bool(std::string_view(values[i]));
        }
        return d;
    }

}  // namespace num

#endif  // NUM_TYPE_BOOL_H
