#pragma once
#ifndef NUM_SHAPE_H
#define NUM_SHAPE_H 1

// spec: shape.go
//
// 形状查询：`(行数, 列数)`。一维容器的第二维按 Go 的习惯报 0（不是 1）。
//
// Go 的 Shape 靠 any 分派，另外两条分支在 C++ 没有对应：标量（any 为单值）报
// (0, 0)、非切片类型报 (-1, -1) —— C++ 的类型在编译期已知，用不上这两条。

#include <utility>

#include "type.h"

namespace num {

    // 一维：(长度, 0)
    template<typename T>
    inline std::pair<i64, i64> shape(const Vector<T> &x) {
        return {static_cast<i64>(x.size()), 0};
    }

    // 二维：(行, 列)
    template<typename T>
    inline std::pair<i64, i64> shape(const Matrix<T> &x) {
        return {static_cast<i64>(x.shape(0)), static_cast<i64>(x.shape(1))};
    }

}  // namespace num

#endif  // NUM_SHAPE_H
