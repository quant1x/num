#pragma once
#ifndef NUM_MATRIX_H
#define NUM_MATRIX_H 1

// spec: matrix.go
//
// 二维整形与求逆：
//
//   - `concat1d(a, b, c)` 把三条等长序列按行拼成 n×3 矩阵（Go 的 Concat1D，列数写死 3）；
//   - `transpose2d` 转置；
//   - `inverse` 用 Gauss-Jordan 消元求逆 —— **不做主元搜索**，奇异矩阵得到 Inf/NaN 而不报错
//     （与 Go 一致；要报错得两侧一起改，登记在 num.md §5）。
//
// 形状不合法（不是方阵、拼列时长度不等）抛 LengthMismatch —— Go 侧是下标越界 panic，
// C++ 的越界是 UB，必须先拦住。

#include <cstddef>

#include "errors.h"
#include "type.h"

namespace num {

    // 矩阵乘：顺序累加，与 Go 的 __matmul 同口径（归约顺序敏感，见 CONTRIBUTING §5.6.3）。
    // 维度不匹配（a 的列数 != b 的行数）抛 LengthMismatch，Go 侧是下标越界 panic。
    // Go 的 Dot / Dot_v1 就是这一份语义（两者逐字相同），故不再另设同义符号。
    template<typename T>
    inline Matrix<T> matmul(const Matrix<T> &a, const Matrix<T> &b) {
        if (a.shape(1) != b.shape(0)) {
            throw LengthMismatch();
        }
        Matrix<T> d = Matrix<T>::from_shape({a.shape(0), b.shape(1)});
        for (std::size_t i = 0; i < a.shape(0); ++i) {
            for (std::size_t j = 0; j < b.shape(1); ++j) {
                T acc{};
                for (std::size_t k = 0; k < b.shape(0); ++k) {
                    acc += a(i, k) * b(k, j);
                }
                d(i, j) = acc;
            }
        }
        return d;
    }

    // 三条等长序列按行拼成 n×3
    template<typename T>
    inline Matrix<T> concat1d(const Vector<T> &a, const Vector<T> &b, const Vector<T> &c) {
        if (b.size() < a.size() || c.size() < a.size()) {
            throw LengthMismatch();
        }
        const std::size_t n = a.size();
        Matrix<T> d = Matrix<T>::from_shape({n, 3});
        for (std::size_t i = 0; i < n; ++i) {
            d(i, 0) = a[i];
            d(i, 1) = b[i];
            d(i, 2) = c[i];
        }
        return d;
    }

    // 转置
    template<typename T>
    inline Matrix<T> transpose2d(const Matrix<T> &x) {
        Matrix<T> d = Matrix<T>::from_shape({x.shape(1), x.shape(0)});
        for (std::size_t i = 0; i < x.shape(0); ++i) {
            for (std::size_t j = 0; j < x.shape(1); ++j) {
                d(j, i) = x(i, j);
            }
        }
        return d;
    }

    // Gauss-Jordan 求逆（无主元搜索）
    inline Matrix<f64> inverse(const Matrix<f64> &a) {
        const std::size_t n = a.shape(0);
        if (a.shape(1) != n) {
            throw LengthMismatch();
        }
        // from_shape 只分配、不清零（见 type.h）：增广矩阵的右半（单位阵那块）必须逐格写满，
        // 否则非对角位置吃的是堆上的残留值，求逆结果会随内存内容漂移。
        Matrix<f64> augmented = Matrix<f64>::from_shape({n, 2 * n});
        for (std::size_t i = 0; i < n; ++i) {
            for (std::size_t j = 0; j < n; ++j) {
                augmented(i, j) = a(i, j);
                augmented(i, j + n) = (i == j) ? 1.0 : 0.0;
            }
        }

        // 前向消元
        for (std::size_t i = 0; i < n; ++i) {
            const f64 pivot = augmented(i, i);
            for (std::size_t j = i + 1; j < n; ++j) {
                const f64 factor = augmented(j, i) / pivot;
                for (std::size_t k = i; k < 2 * n; ++k) {
                    augmented(j, k) -= factor * augmented(i, k);
                }
            }
        }

        // 回代
        for (std::size_t i = n; i-- > 0;) {
            const f64 pivot = augmented(i, i);
            for (std::size_t j = i; j-- > 0;) {
                const f64 factor = augmented(j, i) / pivot;
                for (std::size_t k = i; k < 2 * n; ++k) {
                    augmented(j, k) -= factor * augmented(i, k);
                }
            }
        }

        // 行归一
        for (std::size_t i = 0; i < n; ++i) {
            const f64 pivot = augmented(i, i);
            for (std::size_t j = 0; j < 2 * n; ++j) {
                augmented(i, j) /= pivot;
            }
        }

        Matrix<f64> d = Matrix<f64>::from_shape({n, n});
        for (std::size_t i = 0; i < n; ++i) {
            for (std::size_t j = 0; j < n; ++j) {
                d(i, j) = augmented(i, j + n);
            }
        }
        return d;
    }

}  // namespace num

#endif  // NUM_MATRIX_H
