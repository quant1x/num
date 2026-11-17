#pragma once
#ifndef NUM_LOGIC_EQUAL_H
#define NUM_LOGIC_EQUAL_H 1

// spec: logic_equal.go
//
// 逐元素相等 / 不相等。与 logic.h 的比较算子有一处口径差异（源自 Go）：
// 结果长度取**第一个操作数**，第二个按 type_default<T>() 补齐（浮点是 NaN），
// 不是取最长。故 float 场景下超出的那一侧会与 NaN 比较，结果恒为 false、
// 不相等恒为 true。

#include "align.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<bool> equal(const Vector<T> &x, const Vector<T> &y) {
        const Vector<T> b = align(y, type_default<T>(), x.size());
        Vector<bool> d = Vector<bool>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<f64>(x[i]) == static_cast<f64>(b[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<bool> not_equal(const Vector<T> &x, const Vector<T> &y) {
        const Vector<T> b = align(y, type_default<T>(), x.size());
        Vector<bool> d = Vector<bool>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<f64>(x[i]) != static_cast<f64>(b[i]);
        }
        return d;
    }

}  // namespace num

#endif  // NUM_LOGIC_EQUAL_H
