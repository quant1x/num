#pragma once
#ifndef NUM_CONVERSIONS_H
#define NUM_CONVERSIONS_H 1

// spec: x32/vek.go / x64/vek.go 的 From* / To*，以及 DType2Int
//
// Vector 级的类型转换：from_* 从 bool / 整型 / 浮点转目标类型；to_* 从数值转 bool / 整型 / 浮点。
// 整型转换按向零截断（C++ 的 static_cast 语义），与 Go 一致。

#include <type_traits>

#include "type.h"

namespace num {

    template<typename T>
    inline Vector<T> from_bool(const Vector<bool> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = x[i] ? T{1} : T{0};
        }
        return d;
    }

    template<typename T>
    inline Vector<T> from_int32(const Vector<i32> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<T>(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<T> from_int64(const Vector<i64> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<T>(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<T> from_float32(const Vector<f32> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<T>(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<T> from_float64(const Vector<f64> &x) {
        Vector<T> d = Vector<T>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<T>(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<bool> to_bool(const Vector<T> &x) {
        Vector<bool> d = Vector<bool>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = x[i] != T{0};
        }
        return d;
    }

    template<typename T>
    inline Vector<i32> to_int32(const Vector<T> &x) {
        Vector<i32> d = Vector<i32>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<i32>(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<i64> to_int64(const Vector<T> &x) {
        Vector<i64> d = Vector<i64>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<i64>(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<f32> to_float32(const Vector<T> &x) {
        Vector<f32> d = Vector<f32>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<f32>(x[i]);
        }
        return d;
    }

    template<typename T>
    inline Vector<f64> to_float64(const Vector<T> &x) {
        Vector<f64> d = Vector<f64>::from_shape({x.size()});
        for (std::size_t i = 0; i < x.size(); ++i) {
            d[i] = static_cast<f64>(x[i]);
        }
        return d;
    }

    inline Vector<i32> dtype2int(const Vector<f64> &x) { return to_int32(x); }

}  // namespace num

#endif  // NUM_CONVERSIONS_H
