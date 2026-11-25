#pragma once
#ifndef NUM_NORM_H
#define NUM_NORM_H 1

// spec: x32/vek.go Norm / ManhattanNorm

#include <cmath>

#include "type.h"

namespace num {

    template<typename T>
    inline T norm(const Vector<T> &x) {
        T total = T{0};
        for (std::size_t i = 0; i < x.size(); ++i) {
            total += x[i] * x[i];
        }
        return std::sqrt(total);
    }

    template<typename T>
    inline T manhattan_norm(const Vector<T> &x) {
        T total = T{0};
        for (std::size_t i = 0; i < x.size(); ++i) {
            total += std::abs(x[i]);
        }
        return total;
    }

}  // namespace num

#endif  // NUM_NORM_H
