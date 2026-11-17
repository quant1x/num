#pragma once
#ifndef NUM_LOGIC_ALL_H
#define NUM_LOGIC_ALL_H 1

// spec: logic_all.go
//
// 全部为真：数值按「不为 0」判（NaN 不为 0，故算真），布尔按值判。空输入为真。

#include "type.h"

namespace num {

    template<typename T>
    inline bool all(const Vector<T> &x) {
        for (std::size_t i = 0; i < x.size(); ++i) {
            if (!(x[i] != T{0})) {
                return false;
            }
        }
        return true;
    }

}  // namespace num

#endif  // NUM_LOGIC_ALL_H
