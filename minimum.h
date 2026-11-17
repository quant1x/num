#pragma once
#ifndef NUM_MINIMUM_H
#define NUM_MINIMUM_H 1

// spec: minimum.go
//
// 两个序列横向取小：长度按最长算、短的一侧用 type_default<T>() 补齐，
// 遇 NaN 保持原值（只有 y[i] < x[i] 成立才覆盖）。

#include <num/align.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> minimum(const Vector<T> &x, const Vector<T> &y) {
        const std::size_t n = x.size() > y.size() ? x.size() : y.size();
        Vector<T> d = align(x, type_default<T>(), n);
        const Vector<T> b = align(y, type_default<T>(), n);
        for (std::size_t i = 0; i < n; ++i) {
            if (b[i] < d[i]) {
                d[i] = b[i];
            }
        }
        return d;
    }

    template<typename T>
    inline Vector<T> minimum(const Vector<T> &x, T y) {
        Vector<T> d = x;
        for (std::size_t i = 0; i < d.size(); ++i) {
            if (y < d[i]) {
                d[i] = y;
            }
        }
        return d;
    }

}  // namespace num

#endif  // NUM_MINIMUM_H
