#pragma once
#ifndef NUM_TRIG_H
#define NUM_TRIG_H 1

// spec: x32/vek.go Sin / Cos / SinCos_Into
//
// Go 侧目前只有 x32 入口；C++ 模板同时覆盖 f32 / f64，是同语义的两条宽度。

#include <cmath>

#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> sin(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = std::sin(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<T> cos(const Vector<T> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = std::cos(x[i]);
        }
        return d;
    }

    // 同时写正弦与余弦：两个输出都与输入等长，结果与分别调 sin / cos 一致。
    template<typename T>
    inline void sincos(const Vector<T> &x, Vector<T> &sin_out, Vector<T> &cos_out) {
        sin_out = Vector<T>::from_shape({x.size()});
        cos_out = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            sin_out[i] = std::sin(x[i]);
            cos_out[i] = std::cos(x[i]);
        }
    }

}  // namespace num

#endif  // NUM_TRIG_H
