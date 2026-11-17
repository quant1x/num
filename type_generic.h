#pragma once
#ifndef NUM_TYPE_GENERIC_H
#define NUM_TYPE_GENERIC_H 1

// spec: type.go（`GenericParse` 的可移植部分）/ type_generic.go
//
// `generic_parse<T>(text)`：按目标类型 T 解析一个字符串，对应 Go 的 `GenericParse[T]`。
// 口径照抄 Go 的 `AnyToGeneric` → `__anyToNumber`：
//
//   - T = string：原样返回；
//   - T = bool：走 `to_bool`（空串 → false，true / false 词表，其余落到 parse_bool）；
//   - T 是浮点：**先按 float64 解析**，再收窄到 T（`"1e50"` 转 float32 得 +Inf，与 Go 一致）；
//   - T 是整型：同样先按 float64 解析（**不是** `parse_int64`，所以 `"1e5"` 得 100000、
//     `"1F"` / `"abc"` 得 0）；解析结果为 NaN / ±Inf 时给 0，否则先取整到 i64 再按目标宽度
//     按位截断（`int8("300")` 得 44、`uint8("-1")` 得 255、`int64("1e20")` 得 MinInt64，
//     都与 Go 在 amd64 上的实测一致）。
//
// 与 Go 的差异：
//
//   - Go 的 `AnyToGeneric[T]` 与 `GenericParse[T]` 共用 any 分派，C++ 侧只有字符串入口
//     `generic_parse<T>`；数值 → 目标类型请直接用 `slice_to_*` / `to_bool` / `to_string`；
//   - `type_generic.go` 本身只有 `CheckoutRawType`（reflect 取切片元素类型），C++ 侧类型在
//     编译期已知，不移植；
//   - `AnyToSlice`（vectors.go，按 any 分派 17 种切片）同样不移植，需要时直接调 `repeat` / `align`。

#include <limits>
#include <string>
#include <string_view>
#include <type_traits>

#include "type.h"
#include "type_bool.h"
#include "type_float64.h"

namespace num {

    template<typename T>
    inline T generic_parse(std::string_view text) {
        if constexpr (std::is_same_v<T, std::string>) {
            return std::string(text);
        } else if constexpr (std::is_same_v<T, bool>) {
            return to_bool(text);
        } else if constexpr (std::is_floating_point_v<T>) {
            return static_cast<T>(parse_float64(text));
        } else if constexpr (std::is_integral_v<T>) {
            const f64 vt = parse_float64(text);
            if (is_nan(vt)) {
                // Go 的 Float64IsNaN 把 ±Inf 也算无效值：目标不是浮点时一律给 0
                return T{};
            }
            // 先取整到 i64（越界时 Go 在 amd64 上得到 MinInt64），再按目标宽度按位截断
            i64 wide = std::numeric_limits<i64>::min();
            constexpr f64 k_min = -9223372036854775808.0;  // -2^63
            constexpr f64 k_max = 9223372036854775808.0;   // 2^63（开区间）
            if (vt >= k_min && vt < k_max) {
                wide = static_cast<i64>(vt);
            }
            return static_cast<T>(static_cast<std::make_unsigned_t<T>>(static_cast<u64>(wide)));
        } else {
            static_assert(std::is_same_v<T, void>,
                          "generic_parse 只支持 string / bool / 整型 / 浮点");
        }
    }

}  // namespace num

#endif  // NUM_TYPE_GENERIC_H
