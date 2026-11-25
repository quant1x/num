#pragma once
#ifndef NUM_PROD_H
#define NUM_PROD_H 1

// spec: x32/vek.go Prod / x64/vek.go Prod

#include "type.h"

namespace num {

    template<typename T>
    inline T prod(const Vector<T> &x) {
        T value = T{1};
        for (std::size_t i = 0; i < x.size(); ++i) {
            value *= x[i];
        }
        return value;
    }

}  // namespace num

#endif  // NUM_PROD_H
