#pragma once
#ifndef NUM_MAT4MUL_H
#define NUM_MAT4MUL_H 1

// spec: x32/vek.go Mat4Mul / x64/vek.go Mat4Mul
//
// 两个 4×4 行主序矩阵相乘；输入长度必须为 16。

#include "errors.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> mat4mul(const Vector<T> &x, const Vector<T> &y) {
        if (x.size() != 16 || y.size() != 16) {
            throw LengthMismatch();
        }
        Vector<T> d = Vector<T>::from_shape({16});
        for (std::size_t row = 0; row < 4; ++row) {
            for (std::size_t col = 0; col < 4; ++col) {
                T value = T{0};
                for (std::size_t k = 0; k < 4; ++k) {
                    value += x[row * 4 + k] * y[k * 4 + col];
                }
                d[row * 4 + col] = value;
            }
        }
        return d;
    }

}  // namespace num

#endif  // NUM_MAT4MUL_H
