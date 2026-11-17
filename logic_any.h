#pragma once
#ifndef NUM_LOGIC_ANY_H
#define NUM_LOGIC_ANY_H 1

// spec: logic_any.go
//
// 任一为真：数值按「不为 0」判（NaN 不为 0，故算真），布尔按值判。空输入为假。

#include <num/type.h>

namespace num {

    template<typename T>
    inline bool any(const Vector<T> &x) {
        for (std::size_t i = 0; i < x.size(); ++i) {
            if (x[i] != T{0}) {
                return true;
            }
        }
        return false;
    }

}  // namespace num

#endif  // NUM_LOGIC_ANY_H
