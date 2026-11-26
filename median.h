#pragma once
#ifndef NUM_MEDIAN_H
#define NUM_MEDIAN_H 1

// spec: median.go
//
// 中位数：先排序再取中点。空序列 NaN；单元素返回其本身；偶长取中间两数的平均；
// 奇长取中间那个；任一 NaN 传播。

#include <algorithm>
#include <cmath>
#include <cstddef>
#include <type_traits>
#include <vector>

#include "type.h"

namespace num {

    template<typename T>
    inline DType median(const Vector<T> &values) {
        const std::size_t n = values.size();
        if (n == 0) {
            return nan();
        }
        std::vector<T> ordered(values.cbegin(), values.cend());
        if constexpr (std::is_floating_point_v<T>) {
            for (const T value : ordered) {
                if (std::isnan(value)) {
                    return nan();
                }
            }
        }
        std::sort(ordered.begin(), ordered.end());
        const std::size_t i = n / 2;
        if (n % 2 == 0) {
            return static_cast<DType>(ordered[i - 1] + ordered[i]) / 2.0;
        }
        return static_cast<DType>(ordered[i]);
    }

}  // namespace num

#endif  // NUM_MEDIAN_H
