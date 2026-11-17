#pragma once
#ifndef NUM_COUNT_H
#define NUM_COUNT_H 1

// spec: count.go
//
// 非零计数：数值按「不为 0」计入（NaN 不为 0，故**计入**），布尔按「为真」计入。

#include <cstddef>
#include <type_traits>

#include <num/type.h>

namespace num {

    template<typename T>
    inline std::size_t count(const Vector<T> &x) {
        std::size_t n = 0;
        for (std::size_t i = 0; i < x.size(); ++i) {
            if constexpr (std::is_same_v<T, bool>) {
                if (x[i]) {
                    ++n;
                }
            } else {
                if (x[i] != T{0}) {
                    ++n;
                }
            }
        }
        return n;
    }

}  // namespace num

#endif  // NUM_COUNT_H
