#pragma once
#ifndef NUM_FIND_H
#define NUM_FIND_H 1

// spec: x32/vek.go Find / x64/vek.go Find
//
// 返回首个匹配下标，未命中返回 -1（与 Go 一致）。

#include "type.h"

namespace num {

    template<typename T>
    inline i64 find(const Vector<T> &x, T value) {
        for (std::size_t i = 0; i < x.size(); ++i) {
            if (x[i] == value) {
                return static_cast<i64>(i);
            }
        }
        return -1;
    }

}  // namespace num

#endif  // NUM_FIND_H
