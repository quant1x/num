#pragma once
#ifndef NUM_RANGE_H
#define NUM_RANGE_H 1

// spec: range.go
//
// 等差序列：arange（起止步长）/ linspace（起止个数）/ range（0..n-1）。
// 步长为 0 时 arange 无解，抛 RangeError（Go 侧同等输入是死循环，属缺陷，不复制）。

#include "errors.h"
#include "type.h"

namespace num {

    // [start, end) 上按 step 等距取值；step 全为正值时取 v < end，负值时取 v > end
    template<typename T>
    inline Vector<T> arange(T start, T end, T step = T{1}) {
        if (step == T{0}) {
            throw RangeError();
        }
        const bool forward = step > T{0};
        std::size_t n = 0;
        for (T v = start; forward ? v < end : v > end; v += step) {
            ++n;
        }
        Vector<T> d = Vector<T>::from_shape({n});
        std::size_t i = 0;
        for (T v = start; i < n; v += step, ++i) {
            d[i] = v;
        }
        return d;
    }

    // [start, end] 上 n 个等距点；endpoint 为 false 时不含 end
    template<typename T>
    inline Vector<T> linspace(T start, T end, std::size_t n, bool endpoint = true) {
        if (n == 0) {
            return Vector<T>::from_shape({0});
        }
        if (n == 1) {
            Vector<T> one = Vector<T>::from_shape({1});
            one[0] = start;
            return one;
        }
        const std::size_t num = endpoint ? n - 1 : n;
        const T step = (end - start) / static_cast<T>(num);
        Vector<T> d = Vector<T>::from_shape({n});
        for (std::size_t i = 0; i < n; ++i) {
            d[i] = start + static_cast<T>(i) * step;
        }
        return d;
    }

    // 0, 1, ..., n-1
    template<typename E>
    inline Vector<E> range(std::size_t n) {
        Vector<E> d = Vector<E>::from_shape({n});
        for (std::size_t i = 0; i < n; ++i) {
            d[i] = static_cast<E>(i);
        }
        return d;
    }

}  // namespace num

#endif  // NUM_RANGE_H
