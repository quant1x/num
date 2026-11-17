#pragma once
#ifndef NUM_ONES_H
#define NUM_ONES_H 1

// spec: ones.go
//
// 与 zeros.h 的不对称来自 Go：Zeros 取长度（shape），Ones 取一个样本向量，
// 返回**同样长度**、元素全为 1 的新向量。

#include <num/repeat.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline Vector<T> ones(const Vector<T> &v) { return repeat(T{1}, v.size()); }

}  // namespace num

#endif  // NUM_ONES_H
