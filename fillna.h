#pragma once
#ifndef NUM_FILLNA_H
#define NUM_FILLNA_H 1

// spec: fillna.go
//
// fill / fill_na 把无效值替换成给定值。无效值的口径是 is_nan 那一条：**NaN 与 ±Inf
// 都算**（Go 的 Float32IsNaN / Float64IsNaN 都含 Inf），不是只判 NaN。
//
// 整型与其它非浮点类型没有无效值（0 是有效值），故与 Go 的 default 分支相同 ——
// 原样返回，不做任何替换。
//
// 两个入口的差别只在 Go 的函数签名（Fill 取同类型值、FillNa 取 any 再转换），
// 语义相同；C++ 侧类型在编译期已知，两者都直接按目标类型的值传入。
//
// Go 的第三个形参 args ...any 里那个 inplace 开关在 C++ 拆成 `*_inplace` 入口
// （与 arithmetics.h 的 add/add_inplace 同一写法）：值语义下返回新容器，要改调用方
// 的容器就走 `*_inplace`。

#include <cmath>
#include <cstddef>
#include <type_traits>

#include "type.h"

namespace num {

    namespace detail {

        // 无效值口径：NaN 与 ±Inf 都算（Go Float32IsNaN / Float64IsNaN）
        template<typename T>
        inline bool is_invalid(T v) {
            if constexpr (std::is_floating_point_v<T>) {
                return std::isnan(v) || std::isinf(v);
            } else {
                return false;
            }
        }

        template<typename T>
        inline void fill_missing(Vector<T> &dest, T d) {
            if constexpr (std::is_floating_point_v<T>) {
                for (std::size_t i = 0; i < dest.size(); ++i) {
                    if (is_invalid(dest[i])) {
                        dest[i] = d;
                    }
                }
            } else {
                (void) dest;
                (void) d;
            }
        }

    }  // namespace detail

    // 返回填充后的新容器（入参不变）
    template<typename T>
    inline Vector<T> fill(const Vector<T> &v, T d) {
        Vector<T> dest = v;
        detail::fill_missing(dest, d);
        return dest;
    }

    // 原地填充
    template<typename T>
    inline void fill_inplace(Vector<T> &v, T d) {
        detail::fill_missing(v, d);
    }

    // 同 fill，对应 Go 的 FillNa（值为 any，C++ 侧按目标类型传入）
    template<typename T>
    inline Vector<T> fill_na(const Vector<T> &x, T v) {
        return fill(x, v);
    }

    // 原地填充，对应 Go 的 FillNa(..., inplace=true)
    template<typename T>
    inline void fill_na_inplace(Vector<T> &x, T v) {
        detail::fill_missing(x, v);
    }

}  // namespace num

#endif  // NUM_FILLNA_H
