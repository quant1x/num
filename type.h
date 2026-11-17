#pragma once
#ifndef NUM_TYPE_H
#define NUM_TYPE_H 1

// spec: type.go / type_dtype.go
//
// 标量别名与容器别名的跨语言口径见 CONTRIBUTING.md §4.3 / §4.6：
// i8..u64 / f32 / f64 与 Rust 原生标量名逐字一致，容器名用 Array / Vector / Matrix。

#include <cmath>
#include <cstddef>
#include <cstdint>
#include <limits>
#include <type_traits>

#include <xtensor/containers/xtensor.hpp>

namespace num {

    using i8 = std::int8_t;
    using u8 = std::uint8_t;
    using i16 = std::int16_t;
    using u16 = std::uint16_t;
    using i32 = std::int32_t;
    using u32 = std::uint32_t;
    using i64 = std::int64_t;
    using u64 = std::uint64_t;
    using f32 = float;
    using f64 = double;

    // 数值运算的标量口径（CONTRIBUTING.md §4.3）
    using DType = f64;

    // 一维容器：Go 侧 []T
    template<typename T>
    using Vector = xt::xtensor<T, 1>;

    // 二维及其以上：Go 侧 [][]T / [][][]T
    template<typename T>
    using Matrix = xt::xtensor<T, 2>;

    // 容器约定：`from_shape({...})` **只分配、不清零**（xtensor 的 shape 构造函数不初始化
    // 元素），与 Go 的 `make`（零值）不同 —— 创建之后必须逐格写完，或者先 `fill(T{})`。
    // 需要"零填充"语义时要显式写，不要指望 from_shape 自带零值。

    template<typename T, std::size_t N>
    using Array = xt::xtensor<T, N>;

    // 该类型的无效值（null / 缺失）：浮点是 NaN，整型是 0，布尔是 false。
    // 与 Go TypeDefault[T]() 同口径 —— 填充类算子（align / shift / where）用它的值
    // 决定"缺位长什么样"，不要在这里改成 0 或哨兵值。
    template<typename T>
    inline T type_default() {
        if constexpr (std::is_floating_point_v<T>) {
            return std::numeric_limits<T>::quiet_NaN();
        } else {
            return T{};
        }
    }

    // DType 的无效值
    inline DType nan() { return std::numeric_limits<DType>::quiet_NaN(); }

    // 判无效值：NaN 与 ±Inf 都算（与 Go Float64IsNaN 同口径）。
    // 判 NaN 一律用它，不要直接与 NaN 比较。
    inline bool is_nan(DType d) { return std::isnan(d) || std::isinf(d); }

    // 类型名：只用于异常文案（Go TypeError 的 "unsupported type: <name>"）。
    template<typename T>
    inline const char *type_name() {
        if constexpr (std::is_same_v<T, bool>) return "bool";
        else if constexpr (std::is_same_v<T, i8>) return "int8";
        else if constexpr (std::is_same_v<T, u8>) return "uint8";
        else if constexpr (std::is_same_v<T, i16>) return "int16";
        else if constexpr (std::is_same_v<T, u16>) return "uint16";
        else if constexpr (std::is_same_v<T, i32>) return "int32";
        else if constexpr (std::is_same_v<T, u32>) return "uint32";
        else if constexpr (std::is_same_v<T, i64>) return "int64";
        else if constexpr (std::is_same_v<T, u64>) return "uint64";
        else if constexpr (std::is_same_v<T, f32>) return "float32";
        else if constexpr (std::is_same_v<T, f64>) return "float64";
        else return "unknown";
    }

}  // namespace num

#endif  // NUM_TYPE_H
