#pragma once
#ifndef NUM_ZEROS_H
#define NUM_ZEROS_H 1

// spec: zeros.go
//
// 给定长度、元素全为 0 的向量。

#include <num/repeat.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> zeros(std::size_t shape) { return repeat(T{0}, shape); }

}  // namespace num

#endif  // NUM_ZEROS_H
