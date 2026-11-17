#pragma once
#ifndef NUM_MEDIAN_H
#define NUM_MEDIAN_H 1

// spec: median.go
//
// 中位数：取入参序列的中点。口径与 Go 的 Median 逐行一致，不要在 C++ 侧"顺手修正"：
//
//   - 空序列返回 NaN()；
//   - 只有一个元素返回 0（Go 侧是 DType(0)，不是该元素本身）；
//   - 偶数长度取中间两个元素的算术平均（加法在 T 上做，窄整型可能溢出，与 Go 同），
//     奇数长度直接取中间那个；
//   - **不排序**：按入参给定的顺序取下标，"序列已排序"是调用方的责任。
//
// Go 侧的 __median_go / __quantile_go 是未接线的私有实现（依赖 internal/partial.TopK
// 的快速选择），没有调用方，C++ 侧不移植。

#include <cstddef>

#include "type.h"

namespace num {

    template<typename T>
    inline DType median(const Vector<T> &values) {
        const std::size_t n = values.size();
        if (n == 0) {
            return nan();
        }
        if (n == 1) {
            return 0.0;
        }
        const std::size_t i = n / 2;
        if (n % 2 == 0) {
            return static_cast<DType>(values[i - 1] + values[i]) / 2.0;
        }
        return static_cast<DType>(values[i]);
    }

}  // namespace num

#endif  // NUM_MEDIAN_H
