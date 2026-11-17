#pragma once
#ifndef NUM_ZTABLE_DATA_H
#define NUM_ZTABLE_DATA_H 1

// spec: ztable_data.go
//
// 正态分布表的原始数据：两张各 10000 项的表，下标即万分之一刻度（percent × 10000）。
//
//   - `percent_to_zscore`：累计概率 → Z 分值，单调递增；
//   - `z_table`：置信区间 → Z 分值，首项是 -Inf（对应 0 置信区间）。
//
// 数据不在头文件里展开：定义落在 ztable_data.cpp 这一个编译单元，包含方只拿到声明，
// 免得每个 TU 都背一份 20000 个浮点字面量。表内容由 ztable_data.go 逐值转写而来，
// 两边改动必须同步。

#include <cstddef>

#include "type.h"

namespace num {

    namespace detail {

        // 表的长度（Go 的 ztableMaxScale）
        inline constexpr std::size_t ztable_size = 10000;

        extern const f64 percent_to_zscore[ztable_size];
        extern const f64 z_table[ztable_size];

    }  // namespace detail

}  // namespace num

#endif  // NUM_ZTABLE_DATA_H
