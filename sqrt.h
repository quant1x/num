#pragma once
#ifndef NUM_SQRT_H
#define NUM_SQRT_H 1

// spec: sqrt.go
//
// 逐元素平方根。特殊值按 IEEE：sqrt(+Inf) = +Inf、sqrt(-Inf) = NaN、sqrt(±0) = ±0，
// 次正规输入不归零。

#include <cmath>
#include <type_traits>

#include <xtensor/core/xmath.hpp>

#include "type.h"

namespace num {

    // 非浮点类型在编译期即拒绝（Go 侧是同义的运行期 ErrUnsupportedType）
    template<typename T>
    inline Vector<T> sqrt(const Vector<T> &x) {
        static_assert(std::is_floating_point_v<T>, "sqrt 只支持 f32 / f64");
        Vector<T> d = xt::sqrt(x);
        return d;
    }

    template<typename T>
    inline void sqrt_inplace(Vector<T> &x) {
        static_assert(std::is_floating_point_v<T>, "sqrt 只支持 f32 / f64");
        x = xt::sqrt(x);
    }

}  // namespace num

#endif  // NUM_SQRT_H
