#pragma once
#ifndef NUM_ALIGN_H
#define NUM_ALIGN_H 1

// spec: align.go
//
// 把向量对齐到长度 n：比 n 长则截断，比 n 短则用 a 补齐。
// 缺位填什么由调用方决定 —— 二元算子传 type_default<T>()（浮点是 NaN），
// 这是 Go 侧 Align 的用法，不要在这里改成 0。

#include "type.h"

namespace num {

    // 对齐到长度 n（值语义：总是返回新向量，不改入参）
    template<typename T>
    inline Vector<T> align(const Vector<T> &x, T a, std::size_t n) {
        if (x.size() == n) {
            return x;
        }
        Vector<T> d = Vector<T>::from_shape({n});
        const std::size_t m = x.size() < n ? x.size() : n;
        for (std::size_t i = 0; i < m; ++i) {
            d[i] = x[i];
        }
        for (std::size_t i = m; i < n; ++i) {
            d[i] = a;
        }
        return d;
    }

}  // namespace num

#endif  // NUM_ALIGN_H
