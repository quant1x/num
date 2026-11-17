#pragma once
#ifndef NUM_STDDEV_H
#define NUM_STDDEV_H 1

// spec: stddev.go
//
// 标准差（总体口径：除以 n，不是 n-1 的样本口径）。空输入返回 0。
//
// 命名：Go 的 Std 在 C++ 侧不能叫 std —— `num::std` 会把 `std::` 限定名遮蔽，
// 同命名空间内所有 `std::xxx` 立刻编不过；故统一叫 stddev，对照表已登记。

#include <cmath>
#include <type_traits>

#include "mean.h"
#include "type.h"

namespace num {

    template<typename T>
    inline T stddev(const Vector<T> &x) {
        static_assert(std::is_floating_point_v<T>, "stddev 只支持 f32 / f64");
        if (x.size() == 0) {
            return T{};
        }
        const T m = mean(x);
        T acc{};
        for (std::size_t i = 0; i < x.size(); ++i) {
            const T dv = x[i] - m;
            acc += dv * dv;
        }
        return static_cast<T>(std::sqrt(static_cast<f64>(acc / static_cast<T>(x.size()))));
    }

}  // namespace num

#endif  // NUM_STDDEV_H
