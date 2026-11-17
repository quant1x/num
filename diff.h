#pragma once
#ifndef NUM_DIFF_H
#define NUM_DIFF_H 1

// spec: diff.go
//
// 一阶离散差分：diff(s, n)[i] = shift(s, n)[i] - s[i] = s[i-n] - s[i]。
//
// 方向是「前一格减当前格」，递增序列得到负值 —— 这是 Go 的口径，与 pandas 的
// diff（当前减前一）符号相反，不要按 pandas 的直觉改。
// 平移引入的越界位是默认值（浮点为 NaN），故这些位置的结果也是 NaN。

#include "shift.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> diff(const Vector<T> &s, int n = 1) {
        Vector<T> d = shift(s, n);
        for (std::size_t i = 0; i < d.size(); ++i) {
            d[i] = static_cast<T>(d[i] - s[i]);
        }
        return d;
    }

}  // namespace num

#endif  // NUM_DIFF_H
