#pragma once
#ifndef NUM_SCATTER_H
#define NUM_SCATTER_H 1

// spec: x32/vek.go Scatter / x64/vek.go Scatter
//
// 在长度为 size 的零序列上按下标写入；x 与 idx 长度必须一致，越界抛 RangeError。

#include "errors.h"
#include "repeat.h"
#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> scatter(const Vector<T> &x, const Vector<i64> &index, std::size_t size) {
        if (x.size() != index.size()) {
            throw LengthMismatch();
        }
        Vector<T> d = repeat(T{0}, size);
        for (std::size_t i = 0; i < index.size(); ++i) {
            const i64 position = index[i];
            if (position < 0 || static_cast<std::size_t>(position) >= size) {
                throw RangeError();
            }
            d[static_cast<std::size_t>(position)] = x[i];
        }
        return d;
    }

}  // namespace num

#endif  // NUM_SCATTER_H
