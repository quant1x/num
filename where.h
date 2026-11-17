#pragma once
#ifndef NUM_WHERE_H
#define NUM_WHERE_H 1

// spec: where.go
//
// 条件选择：condition 为真取 x，否则取 y。三个入参按**最长**对齐，缺位用
// type_default<T>() 补齐；条件按「不为 0」判真，故 NaN 条件算真。

#include <num/align.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> where(const Vector<T> &condition, const Vector<T> &x, const Vector<T> &y) {
        std::size_t n = condition.size();
        if (x.size() > n) {
            n = x.size();
        }
        if (y.size() > n) {
            n = y.size();
        }
        const Vector<T> c = align(condition, type_default<T>(), n);
        const Vector<T> a = align(x, type_default<T>(), n);
        const Vector<T> b = align(y, type_default<T>(), n);

        Vector<T> d = Vector<T>::from_shape({n});
        for (std::size_t i = 0; i < n; ++i) {
            d[i] = (c[i] != T{0}) ? a[i] : b[i];
        }
        return d;
    }

}  // namespace num

#endif  // NUM_WHERE_H
