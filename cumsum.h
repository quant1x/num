#pragma once
#ifndef NUM_CUMSUM_H
#define NUM_CUMSUM_H 1

// spec: cumsum.go
//
// 累计和：第 i 个元素等于 x[0..i] 之和（顺序累加，不重排）。

#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> cumsum(const Vector<T> &x) {
        Vector<T> d = x;
        T s{};
        for (std::size_t i = 0; i < d.size(); ++i) {
            s += d[i];
            d[i] = s;
        }
        return d;
    }

    template<typename T>
    inline void cumsum_inplace(Vector<T> &x) {
        T s{};
        for (std::size_t i = 0; i < x.size(); ++i) {
            s += x[i];
            x[i] = s;
        }
    }

}  // namespace num

#endif  // NUM_CUMSUM_H
