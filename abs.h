#pragma once
#ifndef NUM_ABS_H
#define NUM_ABS_H 1

// spec: abs.go
//
// 逐元素绝对值 + 标量绝对值。浮点按位清符号位（与 Go 一致：abs(-0.0) = +0.0，
// abs(NaN) 仍是 NaN），整型用 Go 的 (n ^ shifted) - shifted 口径（最小值溢出仍返回
// 自身，与 Go 相同）。

#include <cmath>
#include <type_traits>

#include <xtensor/core/xmath.hpp>

#include "type.h"

namespace num {

    // 逐元素绝对值；无符号类型原样返回（Go Abs 同口径）
    template<typename T>
    inline Vector<T> abs(const Vector<T> &x) {
        if constexpr (std::is_unsigned_v<T>) {
            return x;
        } else {
            Vector<T> d = xt::abs(x);
            return d;
        }
    }

    template<typename T>
    inline void abs_inplace(Vector<T> &x) {
        if constexpr (!std::is_unsigned_v<T>) {
            x = xt::abs(x);
        }
    }

    inline i8 abs_int8(i8 n) {
        const i8 shifted = static_cast<i8>(n >> 7);
        return static_cast<i8>((n ^ shifted) - shifted);
    }

    inline i16 abs_int16(i16 n) {
        const i16 shifted = static_cast<i16>(n >> 15);
        return static_cast<i16>((n ^ shifted) - shifted);
    }

    inline i32 abs_int32(i32 n) {
        const i32 shifted = n >> 31;
        return (n ^ shifted) - shifted;
    }

    inline i64 abs_int64(i64 n) {
        const i64 shifted = n >> 63;
        return (n ^ shifted) - shifted;
    }

    inline i64 abs_int(i64 n) {
        const i64 shifted = n >> 63;
        return (n ^ shifted) - shifted;
    }

    inline f32 abs_float32(f32 n) { return std::fabs(n); }

    inline f64 abs_float64(f64 n) { return std::fabs(n); }

}  // namespace num

#endif  // NUM_ABS_H
