#pragma once
#ifndef NUM_DOT_H
#define NUM_DOT_H 1

// spec: dot.go
//
// 一维点积：逐元素相乘后**顺序**累加。累加顺序与 Go 一致，不换成并行 / 分块归约 ——
// 归约顺序敏感，改了会让严格容差的用例漂移（CONTRIBUTING.md §5.6.3）。
// 累加器类型就是元素类型 T（f32 入参按 f32 累加），与 Go 的 T(0) 同口径。
//
// Go 的 Dot1D 以 len(a) 为上界、不校验 b 的长度（b 短则越界 panic、b 长则静默忽略
// 尾巴）；C++ 侧越界是 UB，故要求两序列等长，不等长抛 LengthMismatch。
//
// 二维族（Dot / Dot2D / Dot2D1 / Dot2D1_v2 / Dot2D_V1 / Dot_v1）依赖二维矩阵语义，
// 尚未移植，见 num.md §5 TODO。

#include <cstddef>

#include <num/errors.h>
#include <num/type.h>

namespace num {

    template<typename T>
    inline T dot1d(const Vector<T> &a, const Vector<T> &b) {
        if (a.size() != b.size()) {
            throw LengthMismatch();
        }
        T res{};
        for (std::size_t i = 0; i < a.size(); ++i) {
            res += a[i] * b[i];
        }
        return res;
    }

}  // namespace num

#endif  // NUM_DOT_H
