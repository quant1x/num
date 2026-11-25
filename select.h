#pragma once
#ifndef NUM_SELECT_H
#define NUM_SELECT_H 1

// spec: x32/vek.go Select / x64/vek.go Select
//
// 保留 mask 为真的元素，顺序不变；x 与 mask 长度必须一致。

#include "errors.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> select(const Vector<T> &x, const Vector<bool> &mask) {
        if (x.size() != mask.size()) {
            throw LengthMismatch();
        }
        std::size_t n = 0;
        for (std::size_t i = 0; i < mask.size(); ++i) {
            if (mask[i]) {
                ++n;
            }
        }
        Vector<T> d = Vector<T>::from_shape({n});
        std::size_t j = 0;
        for (std::size_t i = 0; i < x.size(); ++i) {
            if (mask[i]) {
                d[j++] = x[i];
            }
        }
        return d;
    }

}  // namespace num

#endif  // NUM_SELECT_H
