#pragma once
#ifndef NUM_REPEAT_H
#define NUM_REPEAT_H 1

// spec: repeat.go
//
// 重复构造：把标量 a 铺成 n 个元素的向量。n < 1 返回空向量（Go 同口径）。

#include "type.h"

namespace num {

    // 构造长度 n、元素全为 a 的向量
    template<typename T>
    inline Vector<T> repeat(T a, std::size_t n) {
        if (n < 1) {
            return Vector<T>::from_shape({0});
        }
        Vector<T> d = Vector<T>::from_shape({n});
        d.fill(a);
        return d;
    }

    // 把前 n 个元素写成 a（其余不动），等价 Go RepeatInto
    template<typename T>
    inline void repeat_into(Vector<T> &s, T a, std::size_t n) {
        const std::size_t m = n < s.size() ? n : s.size();
        for (std::size_t i = 0; i < m; ++i) {
            s[i] = a;
        }
    }

}  // namespace num

#endif  // NUM_REPEAT_H
