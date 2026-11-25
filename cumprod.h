#pragma once
#ifndef NUM_CUMPROD_H
#define NUM_CUMPROD_H 1

// spec: x32/vek.go CumProd / x64/vek.go CumProd

#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> cumprod(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        T total = T{1};
        for (std::size_t i = 0; i < x.size(); ++i) {
            total *= x[i];
            d[i] = total;
        }
        return d;
    }

}  // namespace num

#endif  // NUM_CUMPROD_H
