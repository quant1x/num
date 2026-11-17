#pragma once
#ifndef NUM_MIN_H
#define NUM_MIN_H 1

// spec: min.go
//
// 最小值：空输入返回 0；NaN 参与比较恒为 false，故首个元素是 NaN 时结果就是 NaN。

#include <num/type.h>

namespace num {

    template<typename T>
    inline T min(const Vector<T> &x) {
        if (x.size() == 0) {
            return T{};
        }
        T m = x[0];
        for (std::size_t i = 1; i < x.size(); ++i) {
            if (x[i] < m) {
                m = x[i];
            }
        }
        return m;
    }

}  // namespace num

#endif  // NUM_MIN_H
