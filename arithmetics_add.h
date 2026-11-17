#pragma once
#ifndef NUM_ARITHMETICS_ADD_H
#define NUM_ARITHMETICS_ADD_H 1

// spec: arithmetics_add.go
//
// 逐元素加：结果长度与第一个操作数相同；第二个操作数可以是标量或向量。

#include <num/arithmetics.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> add(const Vector<T> &x, const Vector<T> &y) {
        Vector<T> d = x + binary_operand(y, x.size());
        return d;
    }

    template<typename T>
    inline Vector<T> add(const Vector<T> &x, T y) {
        Vector<T> d = x + y;
        return d;
    }

    template<typename T>
    inline void add_inplace(Vector<T> &x, const Vector<T> &y) {
        x += binary_operand(y, x.size());
    }

    template<typename T>
    inline void add_inplace(Vector<T> &x, T y) {
        x += y;
    }

}  // namespace num

#endif  // NUM_ARITHMETICS_ADD_H
