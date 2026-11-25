#pragma once
#ifndef NUM_ROUNDING_H
#define NUM_ROUNDING_H 1

// spec: x32/vek.go Ceil / Floor / Round
//
// round 为半数远离零（std::round），不是银行家舍入。

#include <cmath>

#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> ceil(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = std::ceil(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<T> floor(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = std::floor(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<T> round(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = std::round(x[i]);
        }
        return d;
    }

}  // namespace num

#endif  // NUM_ROUNDING_H
