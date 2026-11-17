#pragma once
#ifndef NUM_MEAN_H
#define NUM_MEAN_H 1

// spec: mean.go
//
// 均值：空输入返回 0（Go 门面层 Mean 经 UnaryOperations1，长度 0 时返回零值，
// 不是抛错 —— 不要照 x32/x64 的 panic 改）。

#include <num/sum.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline T mean(const Vector<T> &x) {
        if (x.size() == 0) {
            return T{};
        }
        return sum(x) / static_cast<T>(x.size());
    }

}  // namespace num

#endif  // NUM_MEAN_H
