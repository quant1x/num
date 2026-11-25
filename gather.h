#pragma once
#ifndef NUM_GATHER_H
#define NUM_GATHER_H 1

// spec: x32/vek.go Gather / x64/vek.go Gather
//
// 按下标取值；越界抛 RangeError（Go 是下标越界 panic）。

#include "errors.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> gather(const Vector<T> &x, const Vector<i64> &index) {
        Vector<T> d = Vector<T>::from_shape({index.size()});
        for (std::size_t i = 0; i < index.size(); ++i) {
            const i64 position = index[i];
            if (position < 0 || static_cast<std::size_t>(position) >= x.size()) {
                throw RangeError();
            }
            d[i] = x[static_cast<std::size_t>(position)];
        }
        return d;
    }

}  // namespace num

#endif  // NUM_GATHER_H
