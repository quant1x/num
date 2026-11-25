#pragma once
#ifndef NUM_POW_H
#define NUM_POW_H 1

// spec: pow.go
//
// 逐元素整数次幂。求值口径与 Go 一致：**先升到 f64 求幂再落回 T**，
// 因此 f32 输入走的是 f64 精度的 pow，不要改成 f32 精度以求快。

#include <cmath>

#include "errors.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> pow(const Vector<T> &x, int n) {
        const f64 e = static_cast<f64>(n);
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<T>(std::pow(static_cast<f64>(x[i]), e));
        }
        return d;
    }

    template<typename T>
    inline void pow_inplace(Vector<T> &x, int n) {
        x = pow(x, n);
    }


    // 逐元素幂：两个序列等长；浮点求值升到 f64 再落回 T（与整数幂同口径）
    template<typename T>
    inline Vector<T> pow(const Vector<T> &x, const Vector<T> &y) {
        if (x.size() != y.size()) {
            throw LengthMismatch();
        }
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<T>(std::pow(static_cast<f64>(x[i]), static_cast<f64>(y[i])));
        }
        return d;
    }

}  // namespace num

#endif  // NUM_POW_H
