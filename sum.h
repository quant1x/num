#pragma once
#ifndef NUM_SUM_H
#define NUM_SUM_H 1

// spec: sum.go
//
// 求和：按顺序累加（与 Go Sum_Go 同序），空输入返回 0。
// 归约顺序是口径，不要改成 xt::sum 之类的重排实现 —— 严格容差的用例会漂移
// （CONTRIBUTING.md §5.6.3）。

#include <num/type.h>

namespace num {

    template<typename T>
    inline T sum(const Vector<T> &x) {
        T s{};
        for (std::size_t i = 0; i < x.size(); ++i) {
            s += x[i];
        }
        return s;
    }

}  // namespace num

#endif  // NUM_SUM_H
