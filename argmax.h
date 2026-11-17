#pragma once
#ifndef NUM_ARGMAX_H
#define NUM_ARGMAX_H 1

// spec: argmax.go
//
// 最大值下标：返回**首次**出现的下标（严格大于才更新）；空输入返回 0。
// NaN 参与比较恒为 false：首元素是 NaN 时下标停在 0，中间的 NaN 不会被选中。

#include <cstddef>

#include <num/type.h>

namespace num {

    template<typename T>
    inline std::size_t argmax(const Vector<T> &x) {
        if (x.size() == 0) {
            return 0;
        }
        std::size_t idx = 0;
        T m = x[0];
        for (std::size_t i = 1; i < x.size(); ++i) {
            if (x[i] > m) {
                m = x[i];
                idx = i;
            }
        }
        return idx;
    }

}  // namespace num

#endif  // NUM_ARGMAX_H
