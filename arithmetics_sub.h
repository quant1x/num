#pragma once
#ifndef NUM_ARITHMETICS_SUB_H
#define NUM_ARITHMETICS_SUB_H 1

// spec: arithmetics_sub.go
//
// 逐元素减 z = x - y；sub_inplace 为原地形态（Go SubInplace）。

#include <num/arithmetics.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> sub(const Vector<T> &x, const Vector<T> &y) {
        Vector<T> d = x - binary_operand(y, x.size());
        return d;
    }

    template<typename T>
    inline Vector<T> sub(const Vector<T> &x, T y) {
        Vector<T> d = x - y;
        return d;
    }

    template<typename T>
    inline void sub_inplace(Vector<T> &x, const Vector<T> &y) {
        x -= binary_operand(y, x.size());
    }

    template<typename T>
    inline void sub_inplace(Vector<T> &x, T y) {
        x -= y;
    }

}  // namespace num

#endif  // NUM_ARITHMETICS_SUB_H
