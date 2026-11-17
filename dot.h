#pragma once
#ifndef NUM_DOT_H
#define NUM_DOT_H 1

// spec: dot.go
//
// 一维点积：逐元素相乘后**顺序**累加。累加顺序与 Go 一致，不换成并行 / 分块归约 ——
// 归约顺序敏感，改了会让严格容差的用例漂移（CONTRIBUTING.md §5.6.3）。
// 累加器类型就是元素类型 T（f32 入参按 f32 累加），与 Go 的 T(0) 同口径。
//
// Go 的 Dot1D 以 len(a) 为上界、不校验 b 的长度（b 短则越界 panic、b 长则静默忽略
// 尾巴）；C++ 侧越界是 UB，故要求两序列等长，不等长抛 LengthMismatch。
//
// 二维族与 Go 的对应关系（口径照抄，**不要**按标准线代改写）：
//
//   - Go 的 `Dot` / `Dot_v1`：标准矩阵乘，即 matrix.h 的 `matmul`（两者逐字相同，不另设符号）；
//   - `dot2d`（Go 的 `Dot2D`）：行数取 **a 的列数**、列数取 b 的列数，结果
//     min(a.cols, b.cols) × b.cols，累加上界是 a.cols —— 有效行只有前 min(a.cols, b.cols) 行；
//   - `dot2d_v1`（Go 的 `Dot2D_V1`）：两矩阵各取**右下 3×3** 再相乘，语义等同 Go 的 __align_go；
//   - `dot2d1`（Go 的 `Dot2D1`）：矩阵乘向量，要求 a 的列数 == b 的长度；
//     Go 的 `Dot2D1_v2` 结果与之相同（多打了几行调试输出），不另设符号；
//   - Go 的 `__dot2d_go` / `__align_go` 是内部辅助，不单独导出。
//
// Go 侧 __align_go / Dot2D1_v2 里的 fmt.Println 是调试输出，C++ 侧不打印。
// 越界（如 dot2d_v1 的矩阵不足 3×3）在 Go 是 panic，C++ 抛 LengthMismatch。

#include <cstddef>

#include "errors.h"
#include "type.h"

namespace num {

    template<typename T>
    inline T dot1d(const Vector<T> &a, const Vector<T> &b) {
        if (a.size() != b.size()) {
            throw LengthMismatch();
        }
        T res{};
        for (std::size_t i = 0; i < a.size(); ++i) {
            res += a[i] * b[i];
        }
        return res;
    }

    // Go 的 Dot2D：行数取 a 的列数、列数取 b 的列数
    template<typename T>
    inline Matrix<T> dot2d(const Matrix<T> &a, const Matrix<T> &b) {
        const std::size_t r_len = a.shape(1);
        const std::size_t c_len = b.shape(1);
        const std::size_t x_len = r_len < c_len ? r_len : c_len;
        if (b.shape(0) < r_len || a.shape(0) < x_len) {
            throw LengthMismatch();
        }
        Matrix<T> d = Matrix<T>::from_shape({x_len, c_len});
        for (std::size_t i = 0; i < x_len; ++i) {
            for (std::size_t j = 0; j < c_len; ++j) {
                T acc{};
                for (std::size_t k = 0; k < r_len; ++k) {
                    acc += a(i, k) * b(k, j);
                }
                d(i, j) = acc;
            }
        }
        return d;
    }

    // Go 的 Dot2D_V1：先各取右下 3×3，再按 Dot2D 相乘
    template<typename T>
    inline Matrix<T> dot2d_v1(const Matrix<T> &a, const Matrix<T> &b) {
        constexpr std::size_t w = 3;
        if (a.shape(0) < w || a.shape(1) < w || b.shape(0) < w || b.shape(1) < w) {
            throw LengthMismatch();
        }
        Matrix<T> x = Matrix<T>::from_shape({w, w});
        Matrix<T> y = Matrix<T>::from_shape({w, w});
        for (std::size_t i = 0; i < w; ++i) {
            for (std::size_t j = 0; j < w; ++j) {
                x(i, j) = a(a.shape(0) - w + i, a.shape(1) - w + j);
                y(i, j) = b(b.shape(0) - w + i, b.shape(1) - w + j);
            }
        }
        return dot2d(x, y);
    }

    // Go 的 Dot2D1：矩阵乘向量
    template<typename T>
    inline Vector<T> dot2d1(const Matrix<T> &a, const Vector<T> &b) {
        if (a.shape(1) != b.size()) {
            throw LengthMismatch();
        }
        Vector<T> d = Vector<T>::from_shape({a.shape(0)});
        for (std::size_t i = 0; i < a.shape(0); ++i) {
            T acc{};
            for (std::size_t r = 0; r < b.size(); ++r) {
                acc += a(i, r) * b[r];
            }
            d[i] = acc;
        }
        return d;
    }

}  // namespace num

#endif  // NUM_DOT_H
