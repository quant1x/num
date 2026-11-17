#pragma once
#ifndef NUM_MAX_H
#define NUM_MAX_H 1

// spec: max.go
//
// 最大值：空输入返回 0；NaN 参与比较恒为 false，故首个元素是 NaN 时结果就是 NaN。

#include "type.h"

namespace num {

    template<typename T>
    inline T max(const Vector<T> &x) {
        if (x.size() == 0) {
            return T{};
        }
        T m = x[0];
        for (std::size_t i = 1; i < x.size(); ++i) {
            if (x[i] > m) {
                m = x[i];
            }
        }
        return m;
    }

}  // namespace num

#endif  // NUM_MAX_H
