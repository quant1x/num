#pragma once
#ifndef NUM_NEG_H
#define NUM_NEG_H 1

// spec: x32/vek.go Neg / x64/vek.go Neg

#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> neg(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = -x[i];
        }
        return d;
    }

}  // namespace num

#endif  // NUM_NEG_H
