#pragma once
#ifndef NUM_ARITHMETICS_MUL_H
#define NUM_ARITHMETICS_MUL_H 1

// spec: arithmetics_mul.go
//
// 逐元素乘。

#include "arithmetics.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> mul(const Vector<T> &x, const Vector<T> &y) {
        Vector<T> d = x * binary_operand(y, x.size());
        return d;
    }

    template<typename T>
    inline Vector<T> mul(const Vector<T> &x, T y) {
        Vector<T> d = x * y;
        return d;
    }

    template<typename T>
    inline void mul_inplace(Vector<T> &x, const Vector<T> &y) {
        x *= binary_operand(y, x.size());
    }

    template<typename T>
    inline void mul_inplace(Vector<T> &x, T y) {
        x *= y;
    }

}  // namespace num

#endif  // NUM_ARITHMETICS_MUL_H
