#pragma once
#ifndef NUM_LOGIC_NONE_H
#define NUM_LOGIC_NONE_H 1

// spec: x32/vek.go None / x64/vek.go None
//
// 全部为 false 时返回 true；空输入为 true。

#include "logic_any.h"
#include "type.h"

namespace num {

    template<typename T>
    inline bool none(const Vector<T> &x) {
        return !any(x);
    }

}  // namespace num

#endif  // NUM_LOGIC_NONE_H
