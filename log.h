#pragma once
#ifndef NUM_LOG_H
#define NUM_LOG_H 1

// spec: log.go
//
// log1p / expm1 是为小参数保精度而设的对数 / 指数族：x ≈ 0 时直接算 log(1+x) 会把
// 1+x 舍成 1、直接算 e^x - 1 会把 e^x 舍成 1，有效数字全丢，故走 std::log1p /
// std::expm1。
//
// 输出固定为 Vector<f64>：与 Go 返回 []float64 同口径，不随入参宽度变化（f32 入参
// 也返回 f64），不要改成同宽度返回。
//
// log1p 在 1 + x <= 0 处无定义，抛 DomainError（Go 返回 error）；NaN 一路传出去 ——
// 1 + NaN <= 0 为假，不触发定义域检查。
//
// Go 的 Log1pUnrolled / Expm1Unrolled 是手写四路展开的性能路径，C++ 侧由 xtensor
// 表达式向量化承担，不另设同名入口。

#include <cmath>
#include <cstddef>
#include <string>
#include <type_traits>

#include "errors.h"
#include "type.h"

namespace num {

    // 逐元素 log(1 + x)
    template<typename T>
    inline Vector<f64> log1p(const Vector<T> &x) {
        static_assert(std::is_arithmetic_v<T> && !std::is_same_v<T, bool>, "log1p 只接受数值类型");
        const std::size_t n = x.size();
        Vector<f64> d = Vector<f64>::from_shape({n});
        for (std::size_t i = 0; i < n; ++i) {
            const f64 v = static_cast<f64>(x[i]);
            if (1.0 + v <= 0.0) {
                throw DomainError(std::string("Log1p undefined for x[") + std::to_string(i) + "]=" +
                                  std::to_string(v) + ": 1 + x <= 0");
            }
            d[i] = std::log1p(v);
        }
        return d;
    }

    // 逐元素 e^x - 1
    template<typename T>
    inline Vector<f64> expm1(const Vector<T> &x) {
        static_assert(std::is_arithmetic_v<T> && !std::is_same_v<T, bool>, "expm1 只接受数值类型");
        const std::size_t n = x.size();
        Vector<f64> d = Vector<f64>::from_shape({n});
        for (std::size_t i = 0; i < n; ++i) {
            d[i] = std::expm1(static_cast<f64>(x[i]));
        }
        return d;
    }

}  // namespace num

#endif  // NUM_LOG_H
