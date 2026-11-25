#pragma once
#ifndef NUM_INV_H
#define NUM_INV_H 1

// spec: x32/vek.go Inv / x64/vek.go Inv

#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> inv(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = T{1} / x[i];
        }
        return d;
    }

}  // namespace num

#endif  // NUM_INV_H
