#pragma once
#ifndef NUM_TYPE_INTEGER_H
#define NUM_TYPE_INTEGER_H 1

// spec: type_integer.go
//
// 整型序列的转换：
//
//   - `slice_to_ints`：数值序列逐个转 int64（Go 的 `sliceNumberConvert[F, int]`，Go 的 `int` 在
//     amd64 上是 64 位，故 C++ 侧落 i64）；bool 序列转 1 / 0。
//
// 字符串序列的那条分支（Go 里逐个走 AnyToGeneric[int]，落到 ParseInt64）放在 type_int64.h ——
// 它需要 `parse_int64`，放这里会形成头文件循环。

#include <cstddef>
#include <type_traits>

#include "type.h"

namespace num {

    template<typename T, std::enable_if_t<std::is_arithmetic_v<T> && !std::is_same_v<T, bool>, int> = 0>
    inline Vector<i64> slice_to_ints(const Vector<T> &v) {
        Vector<i64> d = Vector<i64>::from_shape({v.size()});
        for (std::size_t i = 0; i < v.size(); ++i) {
            d[i] = static_cast<i64>(v[i]);
        }
        return d;
    }

    // Go 的 SliceToInts 对 []bool 逐个走 AnyToGeneric[int]：true → 1、false → 0
    inline Vector<i64> slice_to_ints(const Vector<bool> &v) {
        Vector<i64> d = Vector<i64>::from_shape({v.size()});
        for (std::size_t i = 0; i < v.size(); ++i) {
            d[i] = v[i] ? 1 : 0;
        }
        return d;
    }

}  // namespace num

#endif  // NUM_TYPE_INTEGER_H
