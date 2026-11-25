#pragma once
#ifndef NUM_QUANTILE_H
#define NUM_QUANTILE_H 1

// spec: x32/vek.go Quantile / x64/vek.go Quantile
//
// 线性插值：位置为 (n-1)*q，q 截断到 [0,1]；空输入抛 EmptyInput。

#include <algorithm>
#include <cmath>
#include <vector>

#include "errors.h"
#include "type.h"

namespace num {

    template<typename T>
    inline T quantile(const Vector<T> &x, f64 q) {
        if (x.size() == 0) {
            throw EmptyInput();
        }
        if (x.size() == 1) {
            return x[0];
        }
        if (q < 0.0) {
            q = 0.0;
        }
        if (q > 1.0) {
            q = 1.0;
        }
        std::vector<T> sorted(x.begin(), x.end());
        std::sort(sorted.begin(), sorted.end());
        const f64 position = static_cast<f64>(sorted.size() - 1) * q;
        const std::size_t lower = static_cast<std::size_t>(position);
        const std::size_t upper = static_cast<std::size_t>(std::ceil(position));
        if (lower == upper) {
            return sorted[lower];
        }
        const f64 weight = position - static_cast<f64>(lower);
        return sorted[lower] + static_cast<T>((sorted[upper] - sorted[lower]) * weight);
    }

}  // namespace num

#endif  // NUM_QUANTILE_H
