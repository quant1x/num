#pragma once
#ifndef NUM_ARITHMETICS_H
#define NUM_ARITHMETICS_H 1

// spec: arithmetics.go
//
// Go 的 arithmetics.go 是"按运行时类型分派到 x32 / x64 / 包内回退"的那一层，
// C++ 用模板在编译期完成同一件事，故本文件只保留二元算子共用的**操作数对齐**：
// 第二操作数按第一个操作数的长度对齐 —— 标量广播，向量比目标短则用默认值补齐、
// 长则截断（与 Go AnyToSlice(y, len(x)) 同口径）。

#include "align.h"
#include "repeat.h"
#include "type.h"

namespace num {

    // 向量操作数对齐到长度 n
    template<typename T>
    inline Vector<T> binary_operand(const Vector<T> &y, std::size_t n) {
        return align(y, type_default<T>(), n);
    }

    // 标量操作数广播成长度 n
    template<typename T>
    inline Vector<T> binary_operand(T y, std::size_t n) { return repeat(y, n); }

}  // namespace num

#endif  // NUM_ARITHMETICS_H
