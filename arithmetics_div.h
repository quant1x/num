#pragma once
#ifndef NUM_ARITHMETICS_DIV_H
#define NUM_ARITHMETICS_DIV_H 1

// spec: arithmetics_div.go
//
// 逐元素除。除零按 IEEE 走（±Inf / NaN），不抛异常 —— 与 Go 的浮点除零一致。

#include <num/arithmetics.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> div(const Vector<T> &x, const Vector<T> &y) {
        Vector<T> d = x / binary_operand(y, x.size());
        return d;
    }

    template<typename T>
    inline Vector<T> div(const Vector<T> &x, T y) {
        Vector<T> d = x / y;
        return d;
    }

    template<typename T>
    inline void div_inplace(Vector<T> &x, const Vector<T> &y) {
        x /= binary_operand(y, x.size());
    }

    template<typename T>
    inline void div_inplace(Vector<T> &x, T y) {
        x /= y;
    }

}  // namespace num

#endif  // NUM_ARITHMETICS_DIV_H
