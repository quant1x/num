#pragma once
#ifndef NUM_SHIFT_H
#define NUM_SHIFT_H 1

// spec: shift.go
//
// 按周期数平移索引：n > 0 向右（前部留空），n < 0 向左（尾部留空），
// 越界与缺失的位置填 type_default<T>()（浮点是 NaN）。|n| >= len 时整条都是默认值。
//
// n == 0 返回**原序列**（位移 0 位不改变序列，见 num.md §2 缺失值填充）。

#include <cstddef>
#include <type_traits>

#include "type.h"

namespace num {

    // 常量位移
    template<typename T>
    inline Vector<T> shift(const Vector<T> &s, int n) {
        const std::size_t len = s.size();
        Vector<T> d = Vector<T>::from_shape({len});
        const auto m = static_cast<std::size_t>(n < 0 ? -static_cast<i64>(n) : static_cast<i64>(n));
        if (m == 0) {
            d = s;
            return d;
        }
        if (m >= len) {
            d.fill(type_default<T>());
            return d;
        }
        if (n > 0) {
            for (std::size_t i = m; i < len; ++i) {
                d[i] = s[i - m];
            }
            for (std::size_t i = 0; i < m; ++i) {
                d[i] = type_default<T>();
            }
        } else {
            for (std::size_t i = 0; i + m < len; ++i) {
                d[i] = s[i + m];
            }
            for (std::size_t i = len - m; i < len; ++i) {
                d[i] = type_default<T>();
            }
        }
        return d;
    }

    // 位移周期必须是整数：非 int 的算术常量（0.5 / 1L / 1u …）一律拒绝，不静默截断
    template<typename T, typename N,
             typename = std::enable_if_t<std::is_arithmetic_v<N> && !std::is_same_v<N, int>>>
    Vector<T> shift(const Vector<T> &, N) = delete;

    // 逐元素位移：n[i] 为 NaN 或越界时该位置填默认值
    template<typename T>
    inline Vector<T> shift(const Vector<T> &s, const Vector<DType> &n) {
        const std::size_t len = s.size();
        Vector<T> d = Vector<T>::from_shape({len});
        for (std::size_t i = 0; i < len; ++i) {
            const DType period = i < n.size() ? n[i] : nan();
            if (is_nan(period)) {
                d[i] = type_default<T>();
                continue;
            }
            const auto pos = static_cast<i64>(period);
            const i64 src = static_cast<i64>(i) - pos;
            if (src < 0 || src >= static_cast<i64>(len)) {
                d[i] = type_default<T>();
                continue;
            }
            d[i] = s[static_cast<std::size_t>(src)];
        }
        return d;
    }

}  // namespace num

#endif  // NUM_SHIFT_H
