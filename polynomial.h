#pragma once
#ifndef NUM_POLYNOMIAL_H
#define NUM_POLYNOMIAL_H 1

// spec: polynomial.go
//
// 旧版（numpy 1.4 之前）多项式 API 的两个算子：
//
//   - `poly_fit(x, y, deg, increasing)` 最小二乘拟合 `p[0]*x**deg + ... + p[deg]`，
//     走正规方程 XTX⁻¹XT y：先造范德蒙矩阵 X（x 的幂次），再转置、矩阵乘、求逆、矩阵乘向量；
//   - `poly_val(p, x)` 用 Horner 方案逐点求值（逐系数乘 x 再加常数项）。
//
// 系数默认按幂次降序（p[0] 是最高次），`increasing=true` 时改升序，与 Go 的 args[0] 布尔开关同义。
//
// 与 Go 的差异（两侧一起改时登记到 num.md §5）：
//
//   - y 长度必须与 x 相同：Go 的 __matvec 按 len(y) 少算也不报错，C++ 侧抛 LengthMismatch；
//   - deg < 0 抛 RangeError（Go 侧 make 负长度直接 panic）；
//   - x 为空抛 EmptyInput（Go 侧 __transpose 取 a[0] 越界 panic）。
//
// 求逆不做主元搜索（与 Go 一致），范德蒙矩阵病态时系数会带 Inf/NaN，拟合质量由调用方判断。
// `increasing` 只换范德蒙矩阵的列序，但正规方程的消元过程不是简单换列，病态数据上两种列序
// 的浮点结果不保证互为倒序 —— 高次拟合的稳定性由调用方负责（与 Go 同）。

#include <cmath>
#include <cstddef>

#include "errors.h"
#include "matrix.h"
#include "repeat.h"
#include "type.h"

namespace num {

    // 最小二乘多项式拟合，返回 deg+1 个系数
    inline Vector<DType> poly_fit(const Vector<DType> &x, const Vector<DType> &y, i64 deg, bool increasing = false) {
        if (deg < 0) {
            throw RangeError();
        }
        if (x.size() == 0) {
            throw EmptyInput();
        }
        if (y.size() != x.size()) {
            throw LengthMismatch();
        }

        const std::size_t m = x.size();
        const std::size_t n = static_cast<std::size_t>(deg) + 1;

        // 范德蒙矩阵：X[i][j] = x[i] ** k
        Matrix<DType> xv = Matrix<DType>::from_shape({m, n});
        for (std::size_t i = 0; i < m; ++i) {
            for (std::size_t j = 0; j < n; ++j) {
                const std::size_t k = increasing ? j : n - 1 - j;
                xv(i, j) = std::pow(x[i], static_cast<DType>(k));
            }
        }

        const Matrix<DType> xt = transpose2d(xv);            // (deg+1) × M
        const Matrix<DType> xtx = detail::matmul(xt, xv);    // (deg+1) × (deg+1)
        const Matrix<DType> xtxinv = inverse(xtx);
        const Matrix<DType> xtxinvxt = detail::matmul(xtxinv, xt);  // (deg+1) × M

        // 与 y 相乘：顺序累加，与 Go 的 __matvec 同口径
        Vector<DType> coef = Vector<DType>::from_shape({n});
        for (std::size_t j = 0; j < n; ++j) {
            DType acc = 0.0;
            for (std::size_t k = 0; k < m; ++k) {
                acc += xtxinvxt(j, k) * y[k];
            }
            coef[j] = acc;
        }
        return coef;
    }

    // Horner 求值：p 由高次到低次，逐点算 p(x)
    inline Vector<DType> poly_val(const Vector<DType> &p, const Vector<DType> &x) {
        Vector<DType> y = repeat<DType>(0.0, x.size());
        for (std::size_t i = 0; i < p.size(); ++i) {
            for (std::size_t j = 0; j < y.size(); ++j) {
                y[j] = y[j] * x[j] + p[i];
            }
        }
        return y;
    }

}  // namespace num

#endif  // NUM_POLYNOMIAL_H
